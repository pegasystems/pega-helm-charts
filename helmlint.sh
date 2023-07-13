#!/bin/bash

chart_names="addons" "backingservices" "pega"
values_files="test_tls.yml" "values.yml"
providers="k8s" "openshift" "eks" "gke" "pks" "aks"

# Nested for loops to iterate over the chart names, values files, and providers
for chart_name in "${chart_names[@]}"
do
    for values_file in "${values_files[@]}"
    do
        for provider in "${providers[@]}"
        do
            helm lint --with-subcharts --values "lint/$values_file" --set-string "global.provider=$provider" --strict "charts/$chart_name"
        done
    done
done
