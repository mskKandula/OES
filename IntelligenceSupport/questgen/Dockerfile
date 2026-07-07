# Stage 1 — builder
FROM python:3.14-slim AS builder

# Install build dependencies needed by chromadb / sentence-transformers
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    git \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# Install Python dependencies into a target directory for clean copy
COPY ./IntelligenceSupport/questgen/requirements.txt .
RUN pip install --no-cache-dir --target=/build/dependencies -r requirements.txt

# Copy application source
COPY ./IntelligenceSupport/questgen/ .

# Stage 2 — runtime
FROM python:3.14-slim

# Install only runtime shared libraries
RUN apt-get update && apt-get install -y --no-install-recommends \
    libstdc++6 \
    libgomp1 \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user and group
RUN groupadd -r oesgroup && useradd -r -g oesgroup oesuser

WORKDIR /app

COPY --from=builder --chown=oesuser:oesgroup /build .

# Switch to non-root user
USER oesuser

EXPOSE 50051

CMD ["python", "server.py"]
