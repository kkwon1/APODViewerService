FROM golang:1.14.2-alpine3.11

WORKDIR /DevWork/APODViewerService
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build cmd/apodviewer/main.go

# This container exposes port 8080 to the outside world
EXPOSE 8081

CMD ["./main"]
