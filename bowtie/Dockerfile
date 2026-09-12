FROM golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS build

WORKDIR /source
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' \
    -o /out/bowtie-json-schema ./cmd/bowtie-json-schema

FROM scratch
COPY --from=build /out/bowtie-json-schema /bowtie-json-schema
USER 65532:65532
ENTRYPOINT ["/bowtie-json-schema"]
