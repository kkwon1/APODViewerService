# APOD Viewer Service

The back-end service for my APOD Viewer project. It is written in Go with MongoDB for data persistence. NASA's APOD website can be found [here](https://apod.nasa.gov/apod/astropix.html).

The front-end React app repo can be found [here](https://github.com/kkwon1/APODViewer).

The back-end service is responsible for serving the front-end app with APOD images, and eventually will allow users to search, like and save for images. It will also be communicating with the NASA API.

## Dependencies
This project is using go modules. After pulling the repo, simply run `go get -v -d ./...`

## Usage
From the root directory, type `go build cmd/apodviewer/main.go` to build a binary file. Or if you'd like to directly run
the service, type `go run cmd/apodviewer/main.go`. Ensure you have a `.env` file with the appropriate environment variables
before you try to run the service. Make sure MongoDB is installed locally.

## Tests
Type `go test -v ./...` to run all tests. To run a single test file, specify the test file path `go test -v ./test/dir/filename`

## Linting
Type `golangci-lint run` to run a linter. Check out the github page [here](https://github.com/golangci/golangci-lint)