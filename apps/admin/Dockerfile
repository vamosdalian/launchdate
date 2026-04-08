# Multi-stage build for optimized production image

# Stage 1: Build the application
FROM node:20 AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Configure npm to handle SSL issues in build environments
# Note: This is required for environments with corporate proxies/SSL inspection
# In production GitHub Actions, consider using a custom CA bundle instead
RUN npm config set strict-ssl false

# Install dependencies
RUN npm install

# Copy source code  
COPY . .

# Build the application
RUN npm run build

# Stage 2: Production image with nginx
FROM nginx:alpine

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built application from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
