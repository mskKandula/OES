# Stage 1 — builder
FROM python:3.11-slim AS builder

# Install build dependencies needed by chromadb / sentence-transformers
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    gcc \
    g++ \
    cmake \
    git \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app/OES/IntelligenceSupport

# Install Python dependencies into a target directory for clean copy
COPY ./IntelligenceSupport/questgen/requirements.txt .
RUN pip install --no-cache-dir --target=/app/OES/IntelligenceSupport/dependencies -r requirements.txt

# Copy application source
COPY ./IntelligenceSupport/questgen/ .

# Stage 2 — runtime
FROM python:3.11-slim

# Install only runtime shared libraries
RUN apt-get update && apt-get install -y --no-install-recommends \
    libstdc++6 \
    libgomp1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app/OES/IntelligenceSupport

COPY --from=builder /app/OES/IntelligenceSupport/ .

ENV PYTHONPATH="${PYTHONPATH}:/app/OES/IntelligenceSupport/dependencies"

EXPOSE 50051

CMD ["python", "server.py"]
