#!/bin/bash

ADDR=$(kubectl get svc istio-ingressgateway -n istio-system --output jsonpath="{.status.loadBalancer.ingress[0].ip}")
#ADDR="localhost:8080"

counterGO=0
counterJAVA=0
i=0
j=0


while true; do \
    output=$(curl  http://$ADDR/sentiment  -H "Content-type: application/json" \
    -d '{"sentence": "I love microservices"}' \
    -s -w "\t Time: %{time_total}s \t Status: %{http_code} \n" -o -); \
    if echo "$output" | grep -q "go"; then \
    version=${output:50:14} ; response=${output:68} ; time=${response:10:4} ;\
    i=$((i+1)) ; \
    time=$((10#`echo $time`))
    counterGO=$((counterGO + time)) ;\
    duree=$(calc $counterGO/1000); \
    a=$(jq -n $duree/$i) ;\
    result="\e[31m${response:0:17} \e[0m${response:18} \e[41m$version\e[0m " ; \
    echo -e $result "\e[31m Time_Avg: " $a "ms\\n"; sleep 0.8; \
    else \
    version=${output:50:16} ; response=${output:69} ; time=${response:9:5} \
    j=$((j+1)) ;\
    time=$((10#`echo $time`))
    counterJAVA=$((counterJAVA + time)) ;\
    duree=$(calc $counterJAVA/1000); \
    a=$(jq -n $duree/$j) ;\
    result="\e[34m${response:0:17} \e[0m${response:18} \e[44m$version\e[0m " ; \
    echo -e $result "\e[34m Time_Avg: " $a "ms\\n"; sleep 0.8; \
    
    fi \
    done

