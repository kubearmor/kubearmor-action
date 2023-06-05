#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright 2023 Authors of KubeArmor


# Define the platforms to build for
PLATFORMS="linux/amd64 linux/arm64 windows/amd64 windows/arm64"

# Define the output directory
OUTPUT_DIR="../_output/bin"

# Define the version to use
VERSION=$(git rev-parse HEAD)

# Update the output directory
rm -rf $OUTPUT_DIR
mkdir -p $OUTPUT_DIR

# Loop through each platform and build the executable file
for PLATFORM in $PLATFORMS; do
  # Split the platform string into OS and ARCH
  IFS="/" read -r -a PARTS <<< "$PLATFORM"
  OS="${PARTS[0]}"
  ARCH="${PARTS[1]}"

  # Set the output file name based on the platform
  # Add .exe extension for Windows
  if [ "$OS" = "windows" ]; then
    FILENAME="$OS-$ARCH-$VERSION.exe"
  else
    FILENAME="$OS-$ARCH-$VERSION"
  fi

  # Set the environment variables for cross-compiling
  export GOOS="$OS"
  export GOARCH="$ARCH"

  # Build the executable file and move it to the output directory
  go build -o "$OUTPUT_DIR/$FILENAME" ../main.go
done