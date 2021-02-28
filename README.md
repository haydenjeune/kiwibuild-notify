# kiwibuild-notify

This project exists because new properties come up for sale quite often on the [KiwiBuild website](https://www.kiwibuild.govt.nz/), however there is no way to configure notifications! This project fixes that by checking for new properties once a day, and sending a notification via email.

## Design

The design is described by the diagram below. The application has been split into two Lambda functions. One reads and stores the properties from the website, and is triggered once a day. The other is responsible for sending the notifications, and is triggered on changes to the properties in the database via a DynamoDB table stream. A email is then sent to me via a SNS topic.

![Design chart](docs/kiwibuild_notifier.png?raw=true)

## Requirements

- [Docker](https://www.docker.com/community-edition)
- [Golang](https://golang.org)
- [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
- [AWS CLI](https://aws.amazon.com/cli/)
- [Graphviz](https://graphviz.gitlab.io/download/) (only for building docs)

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

**Building diagram**

The diagram above is created using the [`diagrams`](https://diagrams.mingrammer.com/) Python package in `docs/architecture.py`. To build a new diagram after changes to this file, run

```bash
make diagram
```

_Note: The Graphviz library must be installed for this to work._

## Packaging and deployment

To deploy the last built version of the application, run

```bash
sam deploy
```
