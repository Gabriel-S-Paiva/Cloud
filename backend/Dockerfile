FROM golang:1.26.6-alpine AS builder
WORKDIR /src
COPY backend/go.mod backend/go.sum ./backend/
ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=$GOPROXY
RUN cd backend && go mod download
COPY backend ./backend
COPY migration ./migration
WORKDIR /src/backend
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server .

FROM alpine:3.20
RUN adduser -D appuser
WORKDIR /app/backend
COPY --from=builder /out/server ./server
COPY --from=builder /src/migration /app/migration
RUN mkdir -p /app/backend/data && chown -R appuser /app
USER appuser
EXPOSE 8080
CMD ["./server"]