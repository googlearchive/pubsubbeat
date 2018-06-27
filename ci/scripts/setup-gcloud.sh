#!/bin/bash

bininfo() {
	echo "Binary information for $1"
	echo "Location:"
	which $1
	echo -e "\nVersion:"
	$1 version
	echo -e "\n\n"
}

SDK_TARGZ=google-cloud-sdk-204.0.0-linux-x86_64.tar.gz


echo "Pulling down gcloud"
curl https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/$SDK_TARGZ > gcloud.tar.gz

echo "Extracting gcloud"
tar xfz gcloud.tar.gz

echo "Configuring path"
export PATH=$PWD/google-cloud-sdk/bin:$PATH

echo "Updating to the most recent components"
gcloud components update --quiet

bininfo gcloud
bininfo gsutil
bininfo bq

echo "Setting up authentication"
echo $SERVICE_ACCOUNT > $PWD/google-cloud-sdk/key.json
gcloud auth activate-service-account --key-file $PWD/google-cloud-sdk/key.json
