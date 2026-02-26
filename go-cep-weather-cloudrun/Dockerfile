FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/server ./cmd/server

FROM alpine:3.20
RUN adduser -D -H appuser
USER appuser
COPY --from=build /bin/server /app/server
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/app/server"]
