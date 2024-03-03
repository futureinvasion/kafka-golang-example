# kafka-golang-example
Built this repository to demonstrate how to write producer/consumer for Kafka with Golang.

## dependencies
- Go v1.15 or higher
- Kafka

## usage
```
go mod tidy
go run ./kafka.go
```

Post message to this API to produce/consume the message on Kafka.

```
curl -X "POST" "http://127.0.0.1:1000/registration" \
     -H 'Content-Type: application/json' \
     -d $'{
  "userId": 9,
  "username": "pete"
}'
```