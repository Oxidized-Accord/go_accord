FROM golang:1.24-alpine

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy project
COPY . .

WORKDIR /app/src/services/messaging_service
RUN go mod download

RUN go build -o messaging_service .

EXPOSE 50051

CMD [ "./messaging_service" ]
