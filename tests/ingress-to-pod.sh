#!/bin/bash



kubectl create namespace bookstore
./osm namespace add bookstore


POD=$(kubectl get pod -n bookbuyer --selector app=bookbuyer --no-headers | awk '{print $1}')

./bin/osm-health ingress to-pod "bookbuyer/${POD}"
