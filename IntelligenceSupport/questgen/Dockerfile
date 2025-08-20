# Stage 1 as builder
FROM python:3.10.5-alpine as builder

# Install necessary build dependencies and libraries
RUN apk add --no-cache \
    build-base \
    gcc \
    g++ \
    libstdc++ \
    libatomic \
    linux-headers

# Create app directory
WORKDIR /app/OES/IntelligenceSupport

# Install app dependencies
COPY ./IntelligenceSupport/questgen/requirements.txt .

RUN pip install --target=/app/OES/IntelligenceSupport/dependencies -r requirements.txt

# Copy code files
COPY ./IntelligenceSupport/questgen/ .

# Stage 2
FROM python:3.10.5-alpine

# Install runtime dependencies
RUN apk add --no-cache libstdc++ libatomic

# Create app directory
WORKDIR /app/OES/IntelligenceSupport

COPY --from=builder	/app/OES/IntelligenceSupport/ .

ENV PYTHONPATH="${PYTHONPATH}:/app/OES/IntelligenceSupport/dependencies"

EXPOSE 50051

CMD ["python", "server.py"]
