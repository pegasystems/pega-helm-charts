#!/bin/bash

exitCode=0
# Nested for loops to iterate over the chart names, values files, and providers
for chart_name in "addons" "backingservices" "pega"
do
    for values_file in "test_tls.yml" "values.yml"
    do
        for provider in "k8s" "openshift" "eks" "gke" "pks" "aks"
        do
            echo "Running helm lint --with-subcharts --values lint/$values_file --set-string global.provider=$provider --strict charts/$chart_name"
            helm lint --with-subcharts --values "lint/$values_file" --set-string "global.provider=$provider" --strict "charts/$chart_name" || exitCode=1
        done
    done
done
exit $exitCode
