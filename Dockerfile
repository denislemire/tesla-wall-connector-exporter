# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS build
RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/tesla-wall-connector-exporter ./cmd/tesla-wall-connector-exporter

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/tesla-wall-connector-exporter /tesla-wall-connector-exporter
EXPOSE 9859
USER nonroot:nonroot
ENTRYPOINT ["/tesla-wall-connector-exporter"]
