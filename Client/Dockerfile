# Stage 1 as builder
FROM node:16 as builder

WORKDIR /app/OES/client

# Copy rest of the files
COPY ./Client/oes/ .

RUN npm install

RUN npm run build

########## Stage 2 ##########
FROM nginx:stable-alpine as production-build

WORKDIR /usr/share/nginx/html

# Delete default files
RUN rm -rf ./*

# Copy from the stage 1
COPY --from=builder /app/OES/client/dist .

# Copy configuration
COPY ./Client/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]