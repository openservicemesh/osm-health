#!/bin/bash



kubectl create namespace bookstore
./osm namespace add bookstore


while [[ $(kubectl get pods -n "bookstore" -l "app=bookstore" -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}') != "True" ]]; do
    echo "waiting for pod" && sleep 1
done

POD=$(kubectl get pod -n bookbuyer --selector app=bookbuyer --no-headers | awk '{print $1}')

./bin/osm-health ingress to-pod "bookbuyer/${POD}"
