all:
	protoc -I.\
		-I./third_party\
		--go_out=plugins=grpc,paths=source_relative:.\
		./v1/exercise_service.proto
	protoc -I.\
		-I./third_party\
		--grpc-gateway_out=logtostderr=true,paths=source_relative:.\
		./v1/exercise_service.proto
	protoc -I.\
		-I./third_party\
		--swagger_out=logtostderr=true:. \
		./v1/exercise_service.proto
	go generate ./v1/