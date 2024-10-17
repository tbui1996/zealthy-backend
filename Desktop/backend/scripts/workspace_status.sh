#!/bin/bash

# Import .env file
export $(egrep -v '^#' .env | xargs)

# Set bazel status variables
echo "BUILD_AWS_ACCOUNT_ID ${AWS_ACCOUNT_ID}"

# Check if working directory is clean
bash scripts/build_git_commit_sha.sh

echo "BUILD_GIT_COMMIT_SHA $(git rev-parse --verify HEAD)"
