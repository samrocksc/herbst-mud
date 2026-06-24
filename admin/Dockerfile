# HerbSt Admin Panel Dockerfile
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . ./
RUN npx vite build

# Production nginx stage
FROM nginx:1.31.2-alpine3.23

# Copy built files to nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Custom nginx config: SPA + API reverse proxy (external file)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]