test:
	go test -cover -tags=unit ./... 
int_test:
	#run docker compose
	docker-compose -f resources/mongo.yml up
	go test -cover -tags=integration ./...
	docker-composer -f resources/mongo.yml down
build:
	#build contianer
run:
	docker-composer -f resources/mongo.yml up
	go run 
	