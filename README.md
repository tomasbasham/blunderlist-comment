# blunderlist-comment

This README outlines the details of collaborating on this Go application. A
short introduction of this app could easily go here.

## Prerequisites

You will need the following things properly installed on your computer.

* [Git](https://git-scm.com/)
* [Go](https://golang.org/)
* [Docker](https://www.docker.com/)

## Installation

* `git clone <repository-url>` this repository
* `cd blunderlist-comment`
* `docker build -t comment .`

## Running / Development

* `docker run --rm -it -p 8080:8080 --env-file .env comment`
* Visit your app at [http://localhost:8080](http://localhost:8080).

## Further Reading / Useful Links

* [Go](https://golang.org/)
* [gRPC](https://grpc.io/)
* [protobuf](https://developers.google.com/protocol-buffers/)
