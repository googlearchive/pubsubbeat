[![Build Status](https://travis-ci.org/GoogleCloudPlatform/pubsubbeat.svg?branch=master)](https://travis-ci.org/GoogleCloudPlatform/pubsubbeat) [![Go Report Card](https://goreportcard.com/badge/github.com/GoogleCloudPlatform/pubsubbeat)](https://goreportcard.com/report/github.com/GoogleCloudPlatform/pubsubbeat) [![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Pubsubbeat

Pubsubbeat is an elastic [Beat](https://www.elastic.co/products/beats) for [Google Cloud Pub/Sub](https://cloud.google.com/pubsub/).
This Beat subscribes to a topic and ingest messages.

The main motivation behind the development of this plugin is to ingest [Stackdriver Logs](https://cloud.google.com/stackdriver/)
via the [Exported Logs](https://cloud.google.com/logging/docs/export/using_exported_logs) feature and send them
directly to Elasticsearch ingest nodes.

This is not an officially supported Google product.

## Getting Started with Pubsubbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.12

### Build

To build the binary for Pubsubbeat run the command below. This will generate a binary
in the same directory with the name pubsubbeat.

```
make
```

### Run

To run Pubsubbeat with debugging output enabled, run:

```
./pubsubbeat -c pubsubbeat.yml -e -d "*"
```

### Test

To test Pubsubbeat, run the following command:

```
make test
```

### Cleanup

To clean  Pubsubbeat source code, run the following commands:

```
make pre-commit
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


## Packaging

To build releases for available platforms:

```
make release
```

This will fetch and create binaries for all Linux, Windows and OSX
