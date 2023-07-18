#!/bin/bash

master_file="charts/pega/templates/_supplemental.tpl"
tpl_files="charts/pega/charts/installer/templates/_supplemental.tpl"
"charts/pega/charts/hazelcast/templates/_supplemental.tpl"
"charts/pega/charts/pegasearch/templates/_supplemental.tpl"

IFS='
'


# Loop through each TPL file
for tpl_file in "${tpl_files[@]}"; do
  # Compare the TPL file with the master file
  if ! cmp -s "$master_file" "$tpl_file"; then
    # Sync the master file to the TPL copy
    cp "$master_file" "$tpl_file"
    echo "Synced $tpl_file with master file."
  fi
done
