FROM --platform=$BUILDPLATFORM docker.io/golang:1.22 AS server-builder
ARG TARGETPLATFORM
WORKDIR /usr/src/app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=$(echo $TARGETPLATFORM | sed 's/linux\///') \
    go build -o dist/main cmd/website/main.go

FROM docker.io/debian:stable-slim AS runner
RUN apt-get update -y && apt-get install -y ca-certificates
WORKDIR /app
COPY assets assets
COPY --from=server-builder /usr/src/app/dist /app
EXPOSE 8080
CMD ["/app/main"]
