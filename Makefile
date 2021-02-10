.PHONY: build dynamo-local dynamo-local-rm network-local network-local-rm

build:
	sam build

run-local: dynamo-local build
	sam local invoke \
		--docker-network lambda-local

dynamo-local: dynamo-local-rm network-local
	docker run -d -p 8000:8000 --network lambda-local --name dynamodb amazon/dynamodb-local
	@sleep 1
	aws dynamodb create-table \
		--table-name Property \
		--attribute-definitions \
			AttributeName=Title,AttributeType=S \
			AttributeName=Type,AttributeType=S \
		--key-schema \
			AttributeName=Title,KeyType=HASH \
			AttributeName=Type,KeyType=SORT \
		--billing-mode PAY_PER_REQUEST \
		--endpoint-url http://localhost:8000 \
		> /dev/null

dynamo-local-rm:
	-@docker stop dynamodb > /dev/null
	-@docker rm dynamodb > /dev/null

network-local: network-local-rm
	docker network create lambda-local

network-local-rm:
	-@docker network rm lambda-local > /dev/null
