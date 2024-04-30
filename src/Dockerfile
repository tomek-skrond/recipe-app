FROM golang:1.22.2-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN mkdir -p ./template

COPY *.go ./
COPY template ./template

RUN CGO_ENABLED=0 GOOS=linux go build -o /recipes 

ENTRYPOINT ["/recipes"]