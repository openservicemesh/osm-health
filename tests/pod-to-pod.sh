#!/bin/bash


if [[ -f "./osm" ]]
then
    echo "Found OSM version $(./osm version)"
else
    OSM_RELEASE="v0.9.1"
    curl -L https://github.com/openservicemesh/osm/releases/download/${OSM_RELEASE}/osm-${OSM_RELEASE}-linux-amd64.tar.gz | tar -vxzf -

    mkdir -p ./bin
    mv ./linux-amd64/osm ./bin
    rm -rf ./linux-amd64
fi


kubectl create namespace bookstore
./bin/osm namespace add bookstore

for x in bookbuyer bookstore; do
    while [[ $(kubectl get pods -n "$x" -l "app=${x}" -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}') != "True" ]]; do
        echo "waiting for pod" && sleep 1
    done
done

POD1=$(kubectl get pod -n bookbuyer --selector app=bookbuyer --no-headers | awk '{print $1}')
POD2=$(kubectl get pod -n bookstore --selector app=bookstore --no-headers | awk '{print $1}')

./bin/osm-health connectivity pod-to-pod \
                 "bookbuyer/${POD1}" \
                 "bookstore/${POD2}"
