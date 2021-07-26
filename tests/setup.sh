#!/bin/bash

set -auexo pipefail

RELEASE="v0.9.1"
curl -L https://github.com/openservicemesh/osm/releases/download/${RELEASE}/osm-${RELEASE}-linux-amd64.tar.gz | tar -vxzf -
mv ./linux-amd64/osm ./
rm -rf ./linux-amd64/osm

./osm install || true

for ns in bookbuyer bookstore; do
    kubectl create namespace "${ns}"
    ./osm namespace add "${ns}"
done

kubectl apply -f https://raw.githubusercontent.com/openservicemesh/osm/release-v0.9/docs/example/manifests/apps/bookbuyer.yaml
kubectl apply -f https://raw.githubusercontent.com/openservicemesh/osm/release-v0.9/docs/example/manifests/apps/bookstore.yaml
./tests/apply-policy.sh
