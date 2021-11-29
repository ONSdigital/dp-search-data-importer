FROM golang:1.17

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/ONSDigital/dp-search-data-importer/

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Build the Go app
RUN go build -o ./docker-dp-search-data-importer .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 25800 to the outside world
EXPOSE 25800

# Run the executable
CMD ["docker-dp-search-data-importer"]

# Create Kafka Topics ##############
ADD topics.sh ./topics.sh

RUN chmod +x ./topics.sh

CMD ./topics.sh