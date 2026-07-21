# =============================================================================
# OES Intelligence Agent — Multi-stage Docker build
#
# Stage 1: Install Python dependencies into an isolated target directory.
# Stage 2: Minimal runtime image with non-root user.
#
# Build context must be the repo root so both stages can reference
# IntelligenceSupport/agent/ paths correctly.
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

COPY ./IntelligenceSupport/agent/requirements.txt .
RUN pip install --no-cache-dir --target=/build/dependencies -r requirements.txt

COPY ./IntelligenceSupport/agent/ .

# Stage 2 — runtime
FROM python:3.12-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libstdc++6 \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -r oesgroup && useradd -r -g oesgroup oesuser

WORKDIR /app

COPY --from=builder --chown=oesuser:oesgroup /build .

USER oesuser

EXPOSE 50051

CMD ["python", "server.py"]