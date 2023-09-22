FROM golang:1.21.1

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest


COPY . /app/

RUN go mod tidy

CMD ["tail", "-f", "/dev/null"]