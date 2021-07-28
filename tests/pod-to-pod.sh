#!/bin/bash



kubectl create namespace bookstore
./osm namespace add bookstore


POD1=$(kubectl get pod -n bookbuyer --selector app=bookbuyer --no-headers | awk '{print $1}')
POD2=$(kubectl get pod -n bookstore --selector app=bookstore --no-headers | awk '{print $1}')

./bin/osm-health connectivity pod-to-pod \
                 "bookbuyer/${POD1}" \
                 "bookstore/${POD2}"
