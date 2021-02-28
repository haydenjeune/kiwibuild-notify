# kiwibuild-notify

This project exists because new properties come up for sale quite often on the [KiwiBuild website](https://www.kiwibuild.govt.nz/), however there is no way to configure notifications! This project fixes that by checking for new properties once a day, and sending a notification via email.

## Requirements

* [Docker](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
* [AWS CLI](https://aws.amazon.com/cli/)

## Setup process

### Installing dependencies & building the target 

To automatically download all the dependencies and package build targets, use 

```bash
make build
```

### Local development

**Invoking locally**

The function that scrapes the KiwiBuild website and stores the result in DynamoDB can be run locally with:

```bash
make run-local-scrape
```

Note that this will run a local version of DynamoDB at `http://localhost:8000`

## Packaging and deployment

To deploy the last built version of the application, run

```bash
sam deploy
```
