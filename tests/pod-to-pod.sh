#!/bin/bash



kubectl create namespace bookstore
./osm namespace add bookstore

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
