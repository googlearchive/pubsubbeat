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

* [Golang](https://golang.org/dl/) 1.9

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
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean  Pubsubbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Pubsubbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/GoogleCloudPlatform/pubsubbeat
git clone https://github.com/GoogleCloudPlatform/pubsubbeat ${GOPATH}/src/github.com/GoogleCloudPlatform/pubsubbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
