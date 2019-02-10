# go-event-processing

An example of how to communicate with a couple of different event messaging systems in go

## Useful Kafka Commands

- `./kafka-console-producer --broker-list localhost:9092 --topic input-topic`
- `./kafka-topics --zookeeper localhost:2181 --list`
- `./kafka-consumer-groups --bootstrap-server localhost:9092 --group local-group-name --describe`

## ectl Commands

- `./ectl nats producer --cluster=test-cluster --channel=input-channel --msg=hello`
- `./ectl nats consumer --cluster=test-cluster --channel=input-channel --group=local-group`
- `./ectl kafka producer --broker=localhost:9092 --topic=input-topic --msg=hello`
- `./ectl kafka consumer --broker=localhost:9092 --topic=input-topic --group=local-group`

### Commands to inject data

- `./ectl nats producer --cluster=test-cluster --channel=input-channel --msg="{\"name\":\"Sven\",\"age\":1,\"type\":\"Dog\"}"`

- `./ectl kafka producer --broker=localhost:9092 --topic=input-topic --msg="{\"name\":\"Treestump\",\"age\":16,\"type\":\"Cat\"}"`
- `./ectl kafka producer --broker=localhost:9092 --topic=input-topic --msg="{\"name\":\"Sven\",\"age\":1,\"type\":\"Dog\"}"`
- `./ectl kafka producer --broker=localhost:9092 --topic=input-topic --msg="{\"name\":\"Meeko\",\"age\":5,\"type\":\"Dog\"}"`
