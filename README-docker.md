Building a docker Image
    docker build -t "dp-search-data-importer:latest" .

    docker images

Run the docker container
    docker-compose up
    docker-compose down

useful commands
    docker ps
    docker ps -a | grep <container-id>

create a topic
docker exec dp-search-data-importer-stream_kafka-1 kafka-topics --create --if-not-exists --zookeeper zookeeper:2181 --partitions 1 --replication-factor 1 --topic search-data-import


Run this app
    Either use make debug
    or run
    %GOPATH/bin/dp-search-data-importer

+++++++++++++++++++++++

If testing this locally then set up dependencies (see following steps).

In dp-compose run MongoDB on port 27017 as follows:

docker-compose up -d

In any directory run Vault as follows:

vault server -dev

In the zebedee directory run Zebedee as follows:

./run.sh

Then in the dp-search-data-importer run:

make debug


The unit tests can be run using this command: 
    make test

The component tests can be run using this 
    command: go test -component