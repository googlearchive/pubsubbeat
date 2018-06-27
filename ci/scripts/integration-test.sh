#!/bin/bash

require() {
	if [ "$1" == "" ]
	then
		echo "The environment variable $2 must be present and contain $3"
		exit 4
	fi
}

section() {
	echo -e "\n============================\n= $1\n============================\n"
}

SCRIPTDIR=$(dirname $0)
SCRIPTDIR=$(realpath $SCRIPTDIR)

section "Checking config"
require "$SERVICE_ACCOUNT" "SERVICE_ACCOUNT" "the service account JSON pubsubbeat will use"
require "$PROJECT_ID" "PROJECT_ID" "the GCP project ID"
require "$TOPIC_NAME" "TOPIC_NAME" "the name of the integration test pubsub topic"
require "$SUBSCRIPTION_NAME" "SUBSCRIPTION_NAME" "the name of the integration test subscription"

section "Building"
source $SCRIPTDIR/setup-go.sh
make

section "Pulling down gcloud"
source $SCRIPTDIR/setup-gcloud.sh

section "Setting up environment"
export TESTID=`date +%s`
export TESTDIR="/tmp/$TESTID"
export KEYPATH="$TESTDIR/key.json"

mkdir -p $TESTDIR

echo "Test ID: $TESTID"
echo $SERVICE_ACCOUNT > $KEYPATH

section "Clearing any existing messages from $SUBSCRIPTION_NAME"
gcloud pubsub subscriptions pull --limit 100 --auto-ack "projects/$PROJECT_ID/subscriptions/$SUBSCRIPTION_NAME"

section "Publishing test messages"

publish_test_message() {
	echo "$1" >> "$TESTDIR/expected.txt"
	gcloud --project=$PROJECT_ID pubsub topics publish $TOPIC_NAME --message "$1"
}

publish_test_message "000) First Message"
publish_test_message "001) Second Message"

section "Running pubsubbeat"
./pubsubbeat -e -v -c ci/fixtures/integration-test-config.yml &
sleep 30

section "Checking Results"
cat "$TESTDIR/actual.txt" | sort > "$TESTDIR/actual-sorted.txt"
echo "Expected"
cat -n "$TESTDIR/expected.txt"
echo "Actual"
cat -n "$TESTDIR/actual-sorted.txt"
ls -lah "$TESTDIR"

cmp "$TESTDIR/actual-sorted.txt" "$TESTDIR/expected.txt"
TEST_RESULT=$?
echo "Result: $TEST_RESULT"


section "Tearing down environment"
gcloud pubsub subscriptions pull --limit 100 --auto-ack "projects/$PROJECT_ID/subscriptions/$SUBSCRIPTION_NAME"
rm -rf $TESTDIR

exit $TEST_RESULT
