.PHONY: build dynamo-local dynamo-local-rm network-local network-local-rm run-local-scrape docs-setup

venv = .venv
python = $(venv)/bin/python
pip = $(venv)/bin/pip

$(venv):
	python3 -m venv $(venv)

docs-setup: $(venv) docs/requirements.txt
	$(pip) install -r docs/requirements.txt

diagram: docs-setup
	$(python) docs/architecture.py
	mv kiwibuild_notifier.png docs/

build:
	sam build

run-local-scrape: dynamo-local build
	sam local invoke \
		--docker-network lambda-local \
		--env-vars local.json \
		"KiwiBuildScrapeFunction"

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
	-@docker stop dynamodb &> /dev/null
	-@docker rm dynamodb &> /dev/null

network-local: network-local-rm
	docker network create lambda-local

network-local-rm:
	-@docker network rm lambda-local &> /dev/null
