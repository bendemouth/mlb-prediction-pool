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

# Build-time arguments for React environment variables.
# These are baked into the JS bundle at build time by Create React App.
ARG REACT_APP_API_URL
ARG REACT_APP_COGNITO_USER_POOL_ID
ARG REACT_APP_COGNITO_APP_CLIENT_ID
ARG REACT_APP_COGNITO_USER_POOL_CLIENT_ID

ENV REACT_APP_API_URL=$REACT_APP_API_URL
ENV REACT_APP_COGNITO_USER_POOL_ID=$REACT_APP_COGNITO_USER_POOL_ID
ENV REACT_APP_COGNITO_APP_CLIENT_ID=$REACT_APP_COGNITO_APP_CLIENT_ID
ENV REACT_APP_COGNITO_USER_POOL_CLIENT_ID=$REACT_APP_COGNITO_USER_POOL_CLIENT_ID

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