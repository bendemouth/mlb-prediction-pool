# ─── Stage 1: Build Go backend ────────────────────────────────────────────────
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main ./cmd/api

# ─── Stage 2: Build React frontend ────────────────────────────────────────────
FROM node:18-alpine AS frontend-builder

WORKDIR /app

# Build-time arguments — baked into the JS bundle by Create React App
ARG REACT_APP_API_URL=/api
ARG REACT_APP_COGNITO_USER_POOL_ID
ARG REACT_APP_COGNITO_APP_CLIENT_ID
ARG REACT_APP_COGNITO_USER_POOL_CLIENT_ID

ENV REACT_APP_API_URL=$REACT_APP_API_URL
ENV REACT_APP_COGNITO_USER_POOL_ID=$REACT_APP_COGNITO_USER_POOL_ID
ENV REACT_APP_COGNITO_APP_CLIENT_ID=$REACT_APP_COGNITO_APP_CLIENT_ID
ENV REACT_APP_COGNITO_USER_POOL_CLIENT_ID=$REACT_APP_COGNITO_USER_POOL_CLIENT_ID

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ .
RUN npm run build

# ─── Stage 3: Production image ─────────────────────────────────────────────────
FROM nginx:alpine AS production

# ca-certificates needed for outbound HTTPS calls from the Go binary
RUN apk add --no-cache ca-certificates wget

# Copy Go binary
COPY --from=backend-builder /app/main /app/main

# Copy React build output
COPY --from=frontend-builder /app/build /usr/share/nginx/html

# Copy nginx config and startup script
COPY nginx.combined.conf /etc/nginx/conf.d/default.conf
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Create non-root user for the Go process
RUN addgroup -g 1001 appuser && adduser -D -u 1001 -G appuser appuser
RUN chown appuser:appuser /app/main

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost/health || exit 1

CMD ["/start.sh"]
