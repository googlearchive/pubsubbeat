# CI Readme

This project performs integration testing using [Concourse](https://concourse-ci.org/).

## Integration Test Setup

If you want to be able to perform integration tests on your own concourse pipeline you will need
to create three things:

* A GCP project for the integration tests.
* A service account in that project that has permissions to read/write/view Pub/Sub topics and
  subscriptions. **NOTE** this service account should **only** have permissions to do this.
* A Pub/Sub topic and subscription for testing.
  These are set to `pubsubbeat-integration-topic` and `pubsubbeat-integration-subscription` in
  the pipeline configuration by default.

Here is a script that will set up your project:

```bash
PROJECT=my-project
TOPIC=pubsubbeat-integration-topic
SUBSCRIPTION=pubsubbeat-integration-subscription
SERVICE_ACCT=pubsub-integration-test-sa

# Create the topic and subscription
gcloud pubsub topics create $TOPIC
gcloud pubsub subscriptions create $SUBSCRIPTION --topic=$TOPIC --topic-project=$PROJECT

# Create the service account and key
gcloud iam service-accounts create $SERVICE_ACCT
gcloud iam service-accounts keys create key.json --iam-account $SERVICE_ACCT@$PROJECT.iam.gserviceaccount.com

# Grant the service account the proper permissions
gcloud projects add-iam-policy-binding $PROJECT --member serviceAccount:$SERVICE_ACCT@$PROJECT.iam.gserviceaccount.com \
  --role roles/pubsub.viewer --role roles/pubsub.publisher --role roles/pubsub.subscriber
```

## Integration Test Teardown

If you want to stop using concourse and want to tear down your platform, you can do that using the
following commands:

```bash
PROJECT=my-project
TOPIC=pubsubbeat-integration-topic
SUBSCRIPTION=pubsubbeat-integration-subscription
SERVICE_ACCT=pubsub-integration-test-sa

# Delete the topic and subscription
gcloud pubsub topics delete $TOPIC

# Delete the service account and key
gcloud iam service-accounts delete $SERVICE_ACCT@$PROJECT.iam.gserviceaccount.com
```

