#!/bin/bash

operator-sdk new pega-operator --type=helm \
    --api-version=platform.pega.com/v1alpha1 \
    --kind=Pega \
    --helm-chart ./charts/pega