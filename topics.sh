#!/bin/sh
echo "Start: Sleep 15 seconds"
sleep 15;
wait;
echo "Begin creating topic"
docker exec dp-search-data-importer-stream_kafka-1 kafka-topics --create --if-not-exists --zookeeper zookeeper:2181 --partitions 1 --replication-factor 1 --topic search-data-import
echo "Done creating topics"