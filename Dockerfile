# We specify the base image we need for our
# go application
FROM golang:1.12.0-alpine3.9

RUN apk update && apk upgrade && apk add --no-cache git

# We create an /app directory within our
# image that will hold our application source
# files
RUN mkdir /app
# We specify that we now wish to execute 
# any further commands inside our /app
# directory
WORKDIR /app

# we run go build to compile the binary
# executable of our Go program
RUN go get -v -u github.com/gorilla/mux
RUN go get -v -u github.com/gorilla/handlers
RUN go get -v -u github.com/nullseed/logruseq
RUN go get -v -u github.com/sirupsen/logrus

# We copy everything in the root directory
# into our /app directory
ADD . /app

RUN go build -o restsvc .

# Our start command which kicks off
# our newly created binary executable
CMD ["/app/restsvc"]
