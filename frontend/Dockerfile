# Development stage
FROM node:18-alpine AS development

WORKDIR /app

# Copy package files
COPY package*.json ./

RUN npm install 

COPY . .

EXPOSE 3000

CMD ["npm", "start"]

FROM node:18-alpine AS seeder

WORKDIR /app

COPY package*.json ./

RUN npm ci || npm install

COPY tsconfig.json ./
COPY scripts/ ./scripts/

RUN npm install --no-save @aws-sdk/client-dynamodb @aws-sdk/lib-dynamodb ts-node typescript @types/node

CMD ["npx", "ts-node", "scripts/seedDatabase.ts"]

# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

RUN npm ci 

COPY . .

RUN npm run build

# Production stage with nginx
FROM nginx:alpine AS production

# Install wget for health checks
RUN apk add --no-cache wget

# Copy build files
COPY --from=builder /app/build /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:80 || exit 1

CMD ["nginx", "-g", "daemon off;"]