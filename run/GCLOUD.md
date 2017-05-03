# GCE Notes

## Requirements

* Stripe account and API keys (you can use the test API keys for testing).
* Auth0 account (see [Auth0 Configuration](#auth0-configuration) section for configuration).
* Google Cloud SDK and `gcloud` utility (`brew cask install google-cloud-sdk` or see https://cloud.google.com/sdk/downloads).

## Auth0 Configuration

Sign up for a new Auth0 account if you do not already have one (the free account will work). Then we need to create the client, create an OAuth key, and enable the management API.

### Social Connection
1. Navigate to Connections and then to Social.
1. Click the enable GitHub.
1. Click the _How to obtain a ClientID?_ link and follow the instructions to create a GitHub OAuth application for Cmd.
1. Under _Attributes_ be sure to select _Email address_ and also check _read:org_ under _Permissions_.

## Client
1. Navigate to the Clients section and Create Client.
1. The client name can be whatever you like (we will use `Cmd`) and it will be a Regular Web Application.
1. Click the Settings tab for this new client and save the Client ID and Client Secret in a safe place (we will need to populate these in our Kubernetes secret later on).


## Google Cloud Platform

###
1. `g config set project cmdio-166321`
1. `g config set compute/zone us-central1-b`
1. `g config set container/cluster cmdio`
1. `g projects create --set-as-default --name cmdio`
1. `g container clusters create cmdio --disk-size 64 --image-type COS --num-nodes 2 --scopes default --enable-cloud-logging --enable-cloud-monitoring`
1. `g container node-pools create sandbox --disk-size 64 --image-type COS --num-nodes 2 --scopes default`
1. `g container clusters get-credentials cmdio`
1. `g compute disks create --size 64GB --zone us-central1-b dynamodb`
1. `k create -f run/dynamodb/dynamodb.yaml`
1. `k create -f run/channels/gcloud.yaml`
1. `k create -f run/dynamodb/dynamodb.yaml`
1. `g compute addresses create sandbox1 sandbox2 --addresses $(gcloud compute instances list --filter='name:cmdio-sandbox' --format='value[terminator=","](networkInterfaces[0].accessConfigs[0].natIP)')`
