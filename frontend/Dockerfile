FROM oven/bun:1 AS builder 
WORKDIR /app

# Copy package files and install dependencies
COPY package.json bun.lockb ./
RUN bun install

# Copy source code
COPY . .

RUN bun run build

# Serve the built files with Nginx
FROM nginx:alpine

# Remove the default Nginx site contents
RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /app/dist /usr/share/nginx/html

# Nginx config
COPY default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
