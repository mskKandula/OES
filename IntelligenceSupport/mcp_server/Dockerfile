# =============================================================================
# OES Intelligence MCP Server — Multi-stage Docker build
#
# Exposes 3 generic tools over Streamable HTTP (FastMCP) on port 8080:
#   - query_database        (read-only MySQL)
#   - execute_kubectl       (read-only k8s via in-cluster RBAC)
#   - search_knowledge_base (ChromaDB similarity search)
#
# Build context must be the repo root.
# =============================================================================

# Stage 1 — builder
FROM python:3.12-slim AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    git \
    curl \
    libffi-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY ./IntelligenceSupport/mcp_server/requirements.txt .
RUN pip install --no-cache-dir --target=/build/dependencies -r requirements.txt

COPY ./IntelligenceSupport/mcp_server/ .

# Stage 2 — runtime
FROM python:3.12-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libstdc++6 \
    libgomp1 \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -r oesgroup && useradd -r -g oesgroup oesuser

WORKDIR /app

COPY --from=builder --chown=oesuser:oesgroup /build .

USER oesuser

EXPOSE 8080

CMD ["python", "mcp_server.py"]