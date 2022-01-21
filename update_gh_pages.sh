#!/usr/bin/env bash

set -eu

repo_uri="https://x-access-token:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"
remote_name="origin"
main_branch="master"
target_branch="gh-pages"
tmp_build_dir="/tmp/build_workspace"

cd "$GITHUB_WORKSPACE"

git config user.name "$GITHUB_ACTOR"
git config user.email "${GITHUB_ACTOR}@bots.github.com"

echo "Creating a temporary directory to build"
mkdir -p "$tmp_build_dir"

echo "clone a single branch gh-pages"
git clone --quiet --branch="$target_branch" --depth=1 '$repo_uri' '$tmp_build_dir' > /dev/null 2>&1

cd "$tmp_build_dir"

echo "Copying the files to the temporary build directory"
rsync -rl --exclude .git --delete "$GITHUB_WORKSPACE/" .

echo "preparing to commit to gh-pages"
git add -A
git commit  -qm 'Deploy ${GITHUB_REPOSITORY} to ${GITHUB_REPOSITORY}:${target_branch}'
git show --stat-count=10 HEAD

echo "Pushing to gh-pages"
git push -q "$remote_name" "$target_branch" > /dev/null 2>&1

git status

echo "Pushing to gh-pages complete"