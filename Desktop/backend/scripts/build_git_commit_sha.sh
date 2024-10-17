#!/bin/bash

# Check if working directory is clean
if [ ! -z "$(git status --porcelain | grep -v "BUILD.bazel" | grep -v "WORKSPACE" | grep -v "Taskfile.yml")" ]; then
  echo "Working directory is not clean! Commit your changes."
  exit 1
fi

# Set bazel image version variable based on git commit sha
git rev-parse --verify HEAD
