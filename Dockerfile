FROM golang:1.20.3-alpine3.17 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./

RUN go mod download
RUN go build -o /app/floppa-bot

FROM alpine:3.17.3 as final

WORKDIR /app

COPY --from=builder /app/floppa-bot /app/floppa-bot

COPY ./.env /app/
ADD floppa /app/floppa/
ADD video  /app/video/

CMD [ "/app/floppa-bot" ]