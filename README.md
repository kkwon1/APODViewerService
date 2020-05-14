# APOD Viewer Service

A back-end service for APOD Viewer written in Go. Communicates with NASA's API to serve image and allow various actions by users.

## Description

APOD (Astronomy Picture of the Day) is published by NASA daily, and can be found [here](https://apod.nasa.gov/apod/astropix.html). APOD Viewer uses NASA's API to service a web-app that allows users to:

- Browse
- Search
- Like
- Save
- Comment

The front-end React app repo can be found [here](https://github.com/kkwon1/APODViewer).

## Prerequisites

This project can be run using docker or running the service locally on your machine. Make sure you've installed all the correct software before running the service. Please create a `.env` file in your root directory with the correct environment variables.

### Docker

You will need to install docker and docker-compose

- You can find the guide to install docker [here](https://docs.docker.com/get-docker/)
- You can read more about docker on [their website](https://www.docker.com/)

### Local

The service is written in Go and uses MongoDB for data persistence. You will need to install both onto your machine.

- You can find the official binary releases [here](https://golang.org/dl/)
- Find the correct installation steps for your OS [here](https://golang.org/doc/install)

This project is using go modules. After pulling the repo, simply run

```bash
go get -v -d ./...
```

### NASA API

NASA let's you sign up for a [free API key](https://api.nasa.gov/). Once you have the key, add that as an environment variable

```
NASA_API_KEY=<YOUR_API_KEY>
```

### Firebase

User authentication is done via Google's Firebase. You will need to set up your own [firebase project](https://firebase.google.com/) and generate a private Admin SDK key (used to verify user token from front end) stored as a `.json` file, which you can set the path as an environment variable.

```
GOOGLE_APPLICATION_CREDENTIALS=<PATH_TO_SECRETS>
```

## Usage

### Docker

To start the mongodb and go service containers.

```bash
docker-compose up
```

Once the container is running, you can ssh into the container by using the command

```bash
docker exec -it <CONTAINER_NAME> sh
```

### Local

From the root directory you can build and run the binary file

```bash
go build main.go
main.exe
```

Or if you'd like to directly run the service without building

```bash
go run cmd/apodviewer/main.go
```

## Running Tests

To run all tests

```bash
go test -v ./...
```

To run a single test file, specify the test file path

```bash
go test -v ./test/dir/<FILE_NAME>
```

## Linting

Make sure you have golangci-lint installed on your computer. You can check out their [github page](https://github.com/golangci/golangci-lint)

```bash
golangci-lint run
```

## Roadmap

- Once all the core features are ready, I'm hoping to implement a simple ML algorithm to recommend images to users depending on what they have searched, liked and saved.

- Maybe I can implement a section where users can upload and share their own astronomy images instead of only relying on NASA as the source

- Mobile applications on iOS and android to reach a larger audience

## Authors

- **Kevin Kwon** - [portfolio](https://kkwon1.github.io/portfolio/)
