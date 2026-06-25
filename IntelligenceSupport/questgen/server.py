import sys
import os
import logging

# setting path
sys.path.append('./pb')

import grpc
from concurrent import futures

from pb.questgen_pb2 import QuestGenResponse, AskQuestionResponse
from pb.questgen_pb2_grpc import QuestGenServiceServicer, add_QuestGenServiceServicer_to_server

import chromadb
from langchain_ollama import OllamaLLM, OllamaEmbeddings
from langchain_chroma import Chroma
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.schema import Document

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# ── Configuration from environment (injected by docker-compose) ──────────────
OLLAMA_BASE_URL = os.environ.get("OLLAMA_BASE_URL", "http://localhost:11434")
CHROMA_HOST     = os.environ.get("CHROMA_HOST", "localhost")
CHROMA_PORT     = int(os.environ.get("CHROMA_PORT", "8000"))
LLM_MODEL       = os.environ.get("LLM_MODEL", "llama3.2:3b")
EMBED_MODEL     = os.environ.get("EMBED_MODEL", "nomic-embed-text")

# Collection name prefix — one collection per client: oes_knowledge_{client_id}
COLLECTION_PREFIX = "oes_knowledge"

# ── Prompt templates ─────────────────────────────────────────────────────────
QUESTION_GEN_PROMPT = """You are an expert exam question generator for an Online Examination System.

Based on the following context extracted from a document, generate exactly 5 high-quality exam questions.
The questions should test understanding, not just recall.
Format each question on a new line, numbered 1 to 5.

Context:
{context}

Generate 5 exam questions:"""

STUDENT_QA_PROMPT = """You are a helpful study assistant for an Online Examination System.
Answer the student's question based strictly on the provided course material context.
If the answer is not in the context, say "I don't have enough information on this topic in the course material."
Keep your answer clear, concise and educational.

Course Material Context:
{context}

Student Question: {question}

Answer:"""


def build_rag_components():
    """Initialise LLM, embeddings and ChromaDB client."""
    logger.info("Connecting to Ollama at %s | LLM: %s | Embed: %s", OLLAMA_BASE_URL, LLM_MODEL, EMBED_MODEL)

    llm = OllamaLLM(
        base_url=OLLAMA_BASE_URL,
        model=LLM_MODEL,
        temperature=0.7,
    )

    embeddings = OllamaEmbeddings(
        base_url=OLLAMA_BASE_URL,
        model=EMBED_MODEL,
    )

    logger.info("Connecting to ChromaDB at %s:%s", CHROMA_HOST, CHROMA_PORT)
    chroma_client = chromadb.HttpClient(host=CHROMA_HOST, port=CHROMA_PORT)

    return llm, embeddings, chroma_client


class QuestGenerator(QuestGenServiceServicer):

    def __init__(self, *args, **kwargs):
        logger.info("Initialising QuestGenerator RAG service...")
        self.llm, self.embeddings, self.chroma_client = build_rag_components()
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=500,
            chunk_overlap=50,
            separators=["\n\n", "\n", ".", " "],
        )
        logger.info("QuestGenerator RAG service ready.")

    def _get_client_store(self, client_id: str) -> Chroma:
        """Return the ChromaDB vector store for the given client.
        Collection name: oes_knowledge_{client_id}
        Both examiner uploads and student queries use this same collection,
        isolated per tenant by collection name.
        """
        collection_name = f"{COLLECTION_PREFIX}_{client_id}" if client_id else f"{COLLECTION_PREFIX}_default"
        logger.info("Using collection: %s", collection_name)
        return Chroma(
            client=self.chroma_client,
            collection_name=collection_name,
            embedding_function=self.embeddings,
            collection_metadata={"hnsw:space": "cosine"},
        )

    # ── RPC 1: Examiner — Generate exam questions from a paragraph ────────────
    def QuestGen(self, request, context):
        paragraph  = request.request
        client_id  = request.client_id  if request.client_id  else "default"
        context_id = request.context_id if request.context_id else ""

        if not paragraph or not paragraph.strip():
            return QuestGenResponse(response="Error: empty input received.")

        try:
            # 1. Split the incoming text into chunks
            raw_docs = self.text_splitter.create_documents([paragraph])
            logger.info("QuestGen [%s]: split input into %d chunk(s).", client_id, len(raw_docs))

            # 2. Attach metadata (context_id for optional student filtering later)
            docs = [
                Document(
                    page_content=d.page_content,
                    metadata={"context_id": context_id} if context_id else {},
                )
                for d in raw_docs
            ]

            # 3. Upsert chunks into the client's ChromaDB collection
            store = self._get_client_store(client_id)
            ids = [f"qg_{client_id}_{abs(hash(d.page_content))}" for d in docs]
            store.add_documents(documents=docs, ids=ids)
            logger.info("QuestGen [%s]: upserted %d chunk(s) to ChromaDB.", client_id, len(docs))

            # 4. Retrieve the top-k most relevant chunks as context
            k = min(3, len(docs))
            retriever = store.as_retriever(
                search_type="similarity",
                search_kwargs={"k": k},
            )
            relevant_docs = retriever.invoke(paragraph[:200])
            context_text = "\n\n".join(d.page_content for d in relevant_docs)
            logger.info("QuestGen [%s]: retrieved %d relevant chunk(s).", client_id, len(relevant_docs))

            # 5. Build prompt and call Ollama LLM
            prompt = QUESTION_GEN_PROMPT.format(context=context_text)
            logger.info("QuestGen [%s]: calling LLM (%s)...", client_id, LLM_MODEL)
            response_text = self.llm.invoke(prompt)
            logger.info("QuestGen [%s]: LLM response received.", client_id)

            return QuestGenResponse(response=response_text)

        except Exception as exc:
            logger.exception("QuestGen [%s] error: %s", client_id, exc)
            return QuestGenResponse(response=f"Error: {str(exc)}")

    # ── RPC 2: Student — Ask a question against indexed course material ────────
    def AskQuestion(self, request, context):
        question   = request.question
        context_id = request.context_id if request.context_id else ""
        client_id  = request.client_id  if request.client_id  else "default"

        if not question or not question.strip():
            return AskQuestionResponse(answer="Error: empty question received.")

        try:
            # Query the client's collection (same one examiner wrote to)
            store = self._get_client_store(client_id)

            # Apply context_id metadata filter if provided
            search_kwargs = {"k": 4}
            if context_id:
                search_kwargs["filter"] = {"context_id": context_id}

            retriever = store.as_retriever(
                search_type="similarity",
                search_kwargs=search_kwargs,
            )
            relevant_docs = retriever.invoke(question)

            if not relevant_docs:
                # No indexed material found — answer from model's general knowledge
                logger.info("AskQuestion [%s]: no course material found, using general LLM knowledge.", client_id)
                prompt = f"You are a helpful study assistant.\n\nQuestion: {question}\n\nAnswer:"
            else:
                context_text = "\n\n".join(d.page_content for d in relevant_docs)
                logger.info("AskQuestion [%s]: retrieved %d relevant chunk(s).", client_id, len(relevant_docs))
                prompt = STUDENT_QA_PROMPT.format(context=context_text, question=question)

            logger.info("AskQuestion [%s]: calling LLM (%s)...", client_id, LLM_MODEL)
            answer = self.llm.invoke(prompt)
            logger.info("AskQuestion [%s]: LLM response received.", client_id)

            return AskQuestionResponse(answer=answer)

        except Exception as exc:
            logger.exception("AskQuestion [%s] error: %s", client_id, exc)
            return AskQuestionResponse(answer=f"Error: {str(exc)}")


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_QuestGenServiceServicer_to_server(QuestGenerator(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    logger.info("gRPC server started on port 50051.")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
