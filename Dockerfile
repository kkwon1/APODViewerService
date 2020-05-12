FROM golang:1.14.2-alpine3.11

WORKDIR /DevWork/APODViewerService

COPY . .

ENV GOOGLE_APPLICATION_CREDENTIALS /DevWork/APODViewerService/secrets/firebase_secrets.json
ENV MONGODB_URI mongodb://mongo:27017

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build main.go

# This container exposes port 8081 to the outside world
EXPOSE 8081

CMD ["./main"]
