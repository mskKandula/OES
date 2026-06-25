# Stage 1 as builder
FROM node:20-alpine AS builder

WORKDIR /app/OES/client

# Copy manifests first to leverage layer cache
COPY ./Client/oes/package.json ./Client/oes/package-lock.json ./
RUN npm ci

# Copy rest of the files
COPY ./Client/oes/ .
RUN npm run build

########## Stage 2 ##########
FROM nginx:stable-alpine AS production-build

RUN apk add --no-cache openssl && \
    openssl req -x509 -nodes -days 365 \
    -subj  "/C=CA/ST=QC/O=Company Inc/CN=example.com" \
    -newkey rsa:2048 \
    -keyout /etc/ssl/private/nginx-selfsigned.key \
    -out /etc/ssl/certs/nginx-selfsigned.crt && \
    rm /etc/nginx/nginx.conf /etc/nginx/conf.d/default.conf

WORKDIR /usr/share/nginx/html

# Copy from the stage 1
COPY --from=builder /app/OES/client/dist .

# Copy configuration
COPY ./Client/nginx.conf /etc/nginx/

EXPOSE 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]
