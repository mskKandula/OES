import sys
import os
import logging
import asyncio

# setting path for pb imports
sys.path.append('./pb')

import grpc
from concurrent import futures

from pb.intelligence_pb2 import IntelligenceResponse
from pb.intelligence_pb2_grpc import IntelligenceServiceServicer, add_IntelligenceServiceServicer_to_server

from langchain_ollama import OllamaLLM
from langchain_core.messages import HumanMessage, SystemMessage
from langchain_mcp_adapters.client import MultiServerMCPClient
from langgraph.prebuilt import create_react_agent

from agent_context import build_system_prompt

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# ── Configuration from environment ───────────────────────────────────────────
OLLAMA_BASE_URL = os.environ.get("OLLAMA_BASE_URL", "http://localhost:11434")
LLM_MODEL       = os.environ.get("LLM_MODEL", "llama3.2:1b")
MCP_SERVER_URL  = os.environ.get("MCP_SERVER_URL", "http://localhost:8080/mcp")

logger.info("Agent config — Ollama: %s | Model: %s | MCP: %s", OLLAMA_BASE_URL, LLM_MODEL, MCP_SERVER_URL)


def build_llm() -> OllamaLLM:
    """Initialise the Ollama LLM for tool-calling."""
    return OllamaLLM(
        base_url=OLLAMA_BASE_URL,
        model=LLM_MODEL,
        temperature=0.2,   # low temp for deterministic tool selection
    )


async def run_agent(query: str, role: str, client_id: str, user_id: str, context_id: str) -> tuple[str, str]:
    """
    Run the ReAct agent for a single request.

    Connects to the MCP server over Streamable HTTP, fetches all tool schemas,
    builds the agent, and invokes it with the user's query + system prompt.

    Returns:
        (answer, tool_used) — final answer string and comma-separated tool names used
    """
    llm = build_llm()
    system_prompt = build_system_prompt(role, client_id, user_id, context_id)

    async with MultiServerMCPClient(
        {
            "oes-tools": {
                "url": MCP_SERVER_URL,
                "transport": "streamable_http",
            }
        }
    ) as mcp_client:
        tools = mcp_client.get_tools()
        logger.info("Agent [%s/%s]: loaded %d tool(s) from MCP server.", role, client_id, len(tools))

        agent = create_react_agent(llm, tools)

        messages = [
            SystemMessage(content=system_prompt),
            HumanMessage(content=query),
        ]

        logger.info("Agent [%s/%s]: invoking with query: %.100s...", role, client_id, query)

        result = await agent.ainvoke({"messages": messages})

        # Extract final answer from last AI message
        final_message = result["messages"][-1]
        answer = final_message.content if hasattr(final_message, "content") else str(final_message)

        # Collect tool names used during the run
        tool_calls = [
            msg.name
            for msg in result["messages"]
            if hasattr(msg, "name") and msg.name
        ]
        tool_used = ", ".join(tool_calls) if tool_calls else "none"

        logger.info("Agent [%s/%s]: tool(s) used: %s", role, client_id, tool_used)

        return answer, tool_used


class IntelligenceAgent(IntelligenceServiceServicer):
    """gRPC servicer — wraps the async ReAct agent in a synchronous gRPC handler."""

    def Ask(self, request, context):
        query      = request.query
        role       = request.role       if request.role       else "Student"
        client_id  = request.client_id  if request.client_id  else "default"
        user_id    = request.user_id    if request.user_id    else ""
        context_id = request.context_id if request.context_id else ""

        if not query or not query.strip():
            return IntelligenceResponse(
                answer="Error: empty query received.",
                tool_used="none",
                success=False,
            )

        try:
            # Run the async agent in a new event loop (gRPC uses a thread pool)
            answer, tool_used = asyncio.run(
                run_agent(query, role, client_id, user_id, context_id)
            )
            return IntelligenceResponse(
                answer=answer,
                tool_used=tool_used,
                success=True,
            )

        except Exception as exc:
            logger.exception("Agent [%s/%s] error: %s", role, client_id, exc)
            return IntelligenceResponse(
                answer=f"Error: {str(exc)}",
                tool_used="none",
                success=False,
            )


def serve():
    # max_workers=4: Ollama processes one inference at a time on CPU.
    # 4 threads match the practical concurrency ceiling — extra threads
    # would just queue behind Ollama anyway, wasting memory.
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    add_IntelligenceServiceServicer_to_server(IntelligenceAgent(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    logger.info("Intelligence Agent gRPC server started on port 50051.")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()