# GCE Notes

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
