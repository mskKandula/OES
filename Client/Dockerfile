# Stage 1 as builder
FROM node:16 as builder

WORKDIR /app/OES/client

# Copy rest of the files
COPY ./Client/oes/ .

RUN npm install

RUN npm run build

########## Stage 2 ##########
FROM nginx:stable-alpine as production-build

RUN rm /etc/nginx/nginx.conf /etc/nginx/conf.d/default.conf
# Delete default files
# RUN rm -rf ./*

WORKDIR /usr/share/nginx/html

RUN apk update && \
    apk add --no-cache openssl && \
    openssl req -x509 -nodes -days 365 \
    -subj  "/C=CA/ST=QC/O=Company Inc/CN=example.com" \
    -newkey rsa:2048 -keyout /etc/ssl/private/nginx-selfsigned.key \
    -out /etc/ssl/certs/nginx-selfsigned.crt;

# Copy from the stage 1
COPY --from=builder /app/OES/client/dist .

# Copy configuration
COPY ./Client/nginx.conf /etc/nginx/

EXPOSE 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]