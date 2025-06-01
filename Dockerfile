FROM golang:1.24.3-bookworm AS build

WORKDIR /usr/src/app
COPY go.mod go.sum ./

RUN go mod download && go mod verify
RUN go install github.com/a-h/templ/cmd/templ@v0.2.501

COPY . .
RUN go generate ./... && go build -v -o /usr/local/bin/app ./cmd


FROM alpine:3.22.0

ENV LEDGER_FILE=/etc/hledger.journal
EXPOSE 8080

RUN apk update && \
  apk add gcompat hledger

COPY --from=build /usr/local/bin/app /usr/local/bin/app
CMD ["app"]
