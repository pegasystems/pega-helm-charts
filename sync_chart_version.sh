#!/bin/bash
set -e
tagVersion=""
if [ ${GITHUB_REF_TYPE} == "tag" ]
then
    tagVersion=${GITHUB_REF_NAME}
fi
export CHART_VERSION=$(expr ${tagVersion:1})

echo "${GITHUB_REF}"
echo "${GITHUB_REPOSITORY}"
echo "${GITHUB_ACTOR}"

repo_uri="https://x-access-token:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"
remote_name="origin"
target_branch="master"
tmp_build_dir="/tmp/build_dir"

cd "$GITHUB_WORKSPACE"

git config --global user.name "$GITHUB_ACTOR"
git config --global user.email "${GITHUB_ACTOR}@bots.github.com"

echo "Creating a temporary directory to build"
rm -rf "$tmp_build_dir"
mkdir -p "$tmp_build_dir"

echo "clone a single branch master"
git clone --quiet --branch="$target_branch" --depth=1 "$repo_uri" "$tmp_build_dir" > /dev/null

cd "$tmp_build_dir"

# Update version in charts/pega/Chart.yaml
awk -v new_version="${CHART_VERSION}" '/^version:/ {$2="\"" new_version "\""}1' charts/pega/Chart.yaml > temp && mv temp charts/pega/Chart.yaml
# Update version in charts/addons/Chart.yaml
awk -v new_version="${CHART_VERSION}" '/^version:/ {$2="\"" new_version "\""}1' charts/addons/Chart.yaml > temp && mv temp charts/addons/Chart.yaml
# Update version in charts/backingservices/Chart.yaml
awk -v new_version="${CHART_VERSION}" '/^version:/ {$2="\"" new_version "\""}1' charts/backingservices/Chart.yaml > temp && mv temp charts/backingservices/Chart.yaml

# Commit changes
git add charts/pega/Chart.yaml charts/addons/Chart.yaml charts/backingservices/Chart.yaml

echo "Updating chart versions to ${CHART_VERSION}"
git commit -m "Update chart versions to ${CHART_VERSION}"

echo "Pushing to master"
git push -q "$remote_name" "$target_branch" > /dev/null