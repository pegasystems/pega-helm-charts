#!/bin/bash

master_file="charts/pega/templates/_supplemental.tpl"
validation_failed=0
failed_files=""

# Loop through each TPL file
for tpl_file in "charts/pega/charts/installer/templates/_supplemental.tpl" "charts/pega/charts/hazelcast/templates/_supplemental.tpl" "charts/pega/charts/pegasearch/templates/_supplemental.tpl"; do
  # Compare the TPL file with the master file
  if ! cmp -s "$master_file" "$tpl_file"; then
    validation_failed=1
    failed_files="$failed_files\n$tpl_file"
  fi
done

#Check validation status
if [ $validation_failed -eq 1 ]; then
  echo -e "Validatoin failed for the following files:"
  echo -e "$failed_files"
  exit 1
else
  echo "Validation passed: All files match the master file."
  exit 0
fi

