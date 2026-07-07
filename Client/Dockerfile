# Stage 1 as builder
FROM node:24-alpine AS builder

WORKDIR /build

# Copy manifests first to leverage layer cache
COPY ./oes/package*.json .
RUN npm ci

# Copy rest of the files
COPY ./oes .
RUN npm run build

########## Stage 2 ##########
FROM nginx:stable-alpine

RUN apk add --no-cache openssl && addgroup -S oesgroup && adduser -S oesuser -G oesgroup && \
    mkdir -p /etc/ssl/private && \
    openssl req -x509 -nodes -days 365 \
    -subj  "/C=CA/ST=QC/O=Company Inc/CN=example.com" \
    -newkey rsa:2048 \
    -keyout /etc/ssl/private/nginx-selfsigned.key \
    -out /etc/ssl/certs/nginx-selfsigned.crt && \
    rm -f /etc/nginx/conf.d/default.conf

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Copy from the stage 1
COPY --from=builder /build/dist /usr/share/nginx/html

RUN chown -R oesuser:oesgroup \
    /usr/share/nginx/html \
    /var/cache/nginx \
    /var/log/nginx \
    /etc/nginx

EXPOSE 80

# nginx master process requires root to bind port 80; worker processes drop to nginx user via config
ENTRYPOINT ["nginx", "-g", "daemon off;"]
