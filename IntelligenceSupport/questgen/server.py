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

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# ── Configuration from environment (injected by docker-compose) ──────────────
OLLAMA_BASE_URL = os.environ.get("OLLAMA_BASE_URL", "http://localhost:11434")
CHROMA_HOST     = os.environ.get("CHROMA_HOST", "localhost")
CHROMA_PORT     = int(os.environ.get("CHROMA_PORT", "8000"))
LLM_MODEL       = os.environ.get("LLM_MODEL", "llama3.2:3b")
EMBED_MODEL     = os.environ.get("EMBED_MODEL", "llama3.2:3b")

# Collection names
QUESTGEN_COLLECTION = "oes_questgen"       # stores examiner-uploaded paragraphs
STUDENT_COLLECTION  = "oes_course_material" # stores indexed course material for student Q&A

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
    """Initialise LLM, embeddings and vector store connections."""
    logger.info("Connecting to Ollama at %s with model %s", OLLAMA_BASE_URL, LLM_MODEL)

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

    # Shared ChromaDB HTTP client
    chroma_client = chromadb.HttpClient(host=CHROMA_HOST, port=CHROMA_PORT)

    # Vector store for examiner question generation
    questgen_store = Chroma(
        client=chroma_client,
        collection_name=QUESTGEN_COLLECTION,
        embedding_function=embeddings,
        collection_metadata={"hnsw:space": "cosine"},
    )

    # Vector store for student Q&A against indexed course material
    student_store = Chroma(
        client=chroma_client,
        collection_name=STUDENT_COLLECTION,
        embedding_function=embeddings,
        collection_metadata={"hnsw:space": "cosine"},
    )

    return llm, questgen_store, student_store


class QuestGenerator(QuestGenServiceServicer):

    def __init__(self, *args, **kwargs):
        logger.info("Initialising QuestGenerator RAG service...")
        self.llm, self.questgen_store, self.student_store = build_rag_components()
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=500,
            chunk_overlap=50,
            separators=["\n\n", "\n", ".", " "],
        )
        logger.info("QuestGenerator RAG service ready.")

    # ── RPC 1: Examiner — Generate exam questions from a paragraph ────────────
    def QuestGen(self, request, context):
        paragraph = request.request

        if not paragraph or not paragraph.strip():
            return QuestGenResponse(response="Error: empty input received.")

        try:
            # 1. Split the incoming text into chunks
            docs = self.text_splitter.create_documents([paragraph])
            logger.info("QuestGen: split input into %d chunk(s).", len(docs))

            # 2. Upsert chunks into ChromaDB (idempotent via content hash)
            ids = [f"qg_chunk_{abs(hash(d.page_content))}" for d in docs]
            self.questgen_store.add_documents(documents=docs, ids=ids)
            logger.info("QuestGen: upserted %d chunk(s) to ChromaDB.", len(docs))

            # 3. Retrieve the top-3 most relevant chunks as context
            k = min(3, len(docs))
            retriever = self.questgen_store.as_retriever(
                search_type="similarity",
                search_kwargs={"k": k},
            )
            relevant_docs = retriever.invoke(paragraph[:200])
            context_text = "\n\n".join(d.page_content for d in relevant_docs)
            logger.info("QuestGen: retrieved %d relevant chunk(s).", len(relevant_docs))

            # 4. Build prompt and call Ollama LLM
            prompt = QUESTION_GEN_PROMPT.format(context=context_text)
            logger.info("QuestGen: calling LLM (%s)...", LLM_MODEL)
            response_text = self.llm.invoke(prompt)
            logger.info("QuestGen: LLM response received.")

            return QuestGenResponse(response=response_text)

        except Exception as exc:
            logger.exception("QuestGen error: %s", exc)
            return QuestGenResponse(response=f"Error: {str(exc)}")

    # ── RPC 2: Student — Ask a question against indexed course material ────────
    def AskQuestion(self, request, context):
        question   = request.question
        context_id = request.context_id if request.context_id else ""

        if not question or not question.strip():
            return AskQuestionResponse(answer="Error: empty question received.")

        try:
            # Retrieve relevant context from the student course material collection.
            # If context_id is provided, filter by that topic (metadata filter).
            search_kwargs = {"k": 4}
            if context_id:
                search_kwargs["filter"] = {"context_id": context_id}

            retriever = self.student_store.as_retriever(
                search_type="similarity",
                search_kwargs=search_kwargs,
            )
            relevant_docs = retriever.invoke(question)

            if not relevant_docs:
                # No indexed material found — answer from model's general knowledge
                logger.info("AskQuestion: no course material found, using general LLM knowledge.")
                prompt = f"You are a helpful study assistant.\n\nQuestion: {question}\n\nAnswer:"
            else:
                context_text = "\n\n".join(d.page_content for d in relevant_docs)
                logger.info("AskQuestion: retrieved %d relevant chunk(s).", len(relevant_docs))
                prompt = STUDENT_QA_PROMPT.format(context=context_text, question=question)

            logger.info("AskQuestion: calling LLM (%s)...", LLM_MODEL)
            answer = self.llm.invoke(prompt)
            logger.info("AskQuestion: LLM response received.")

            return AskQuestionResponse(answer=answer)

        except Exception as exc:
            logger.exception("AskQuestion error: %s", exc)
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
