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
  echo -e "Validation failed for the following files:"
  echo -e "$failed_files"
  echo -e "This indicates that the following files are out of sync with the definitions in the main charts/pega supplemental file."
  echo -e "Since all subcharts should mirror the definitions of the superchart, consider running sync_supplementals.sh to fix this issue."
  exit 1
else
  echo "Validation passed: All files match the master file."
  exit 0
fi

