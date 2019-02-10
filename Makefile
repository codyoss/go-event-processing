.PHONY: clean

build-ectl: 
	cd cmd/ectl && go build && mv ectl ../..
build-kafka-consumer: 
	cd cmd/kafka-consumer-2 && go build && mv kafka-consumer ../..
build-kafka-app: 
	cd cmd/kafka-app && go build && mv kafka-app ../..
build-goka: 
	cd cmd/goka && go build && mv goka ../..
build-nats: 
	cd cmd/nats && go build && mv nats ../..
test:
	go test -count=1 ./...

up:
	docker-compose up -d
down:
	docker-compose down

clean:
	rm -rf /tmp/goka