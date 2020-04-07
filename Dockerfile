# specify the base image we need for our go application
FROM golang:1.12.0-alpine3.9

RUN apk update && apk upgrade && apk add --no-cache git

# create an /app directory within our image that will hold our application source files
RUN mkdir /app

# specify that we now wish to execute any further commands inside our /app directory
WORKDIR /app

# get dependency libraries
RUN go get -v -u github.com/gorilla/mux github.com/gorilla/handlers github.com/nullseed/logruseq github.com/sirupsen/logrus

# copy source into /app
ADD . /app

# build
RUN go build -o restsvc .

# our newly created binary executable
CMD ["/app/restsvc"]
