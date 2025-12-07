FROM golang:1.24-bookworm AS builder

WORKDIR /src

COPY ./ ./

WORKDIR /src/internal/examples

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o server/bin/opamp-server server/main.go

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /src/internal/examples/server/bin/opamp-server ./

EXPOSE 4320 4321

ENTRYPOINT ["./opamp-server"]

# docker build -t opamp-server:dev .

# docker run --rm -p 4320:4320 -p 4321:4321 opamp-server:dev
