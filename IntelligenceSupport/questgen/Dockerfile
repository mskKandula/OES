FROM python:3.10.5

# Create app directory
WORKDIR /app/OES/IntelligenceSupport

# Install app dependencies
COPY ./IntelligenceSupport/questgen/requirements.txt .

RUN pip install -r requirements.txt

# Copy code files
COPY ./IntelligenceSupport/questgen/ .

EXPOSE 50051

CMD [ "python", "server.py" ]