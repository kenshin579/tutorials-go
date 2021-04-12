KAFKA_SERVER_0 := go-kafka-0.go-kafka-headless.default.svc.cluster.local:9092
TOPIC=go-kafka-topic

build-producer:
	go build -o bin/producer.go producer/main.go

build-consumer:
	go build -o bin/consumer.go consumer/main.go


run-producer: build-producer
	@bin/producer -brokers=$(KAFKA_SERVER_0) -topic=$(TOPIC)

run-consumer: build-consumer
	@bin/producer -brokers=$(KAFKA_SERVER_0) -topic=$(TOPIC) -message=$(message)

