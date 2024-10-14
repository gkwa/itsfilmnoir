#!/bin/bash

# Function to list and sort tags
list_tags() {
    local image="$1"
    local filter="$2"
    echo "Listing tags for $image..."
    tags=$(skopeo list-tags "docker://docker.io/$image" | jq --raw-output '.Tags[]' | sort --version-sort)
    
    if [ -n "$filter" ]; then
        echo "$tags" | grep -i "$filter"
    else
        echo "$tags"
    fi
    echo "------------------------"
}

# Check if skopeo and jq are installed
if ! command -v skopeo &> /dev/null || ! command -v jq &> /dev/null; then
    echo "Error: This script requires skopeo and jq to be installed."
    exit 1
fi

# List tags for different images
list_tags "amazon/aws-cli"
list_tags "node" "alpine"
list_tags "homebrew/brew"

echo "Done!"