#!/usr/bin/env sh
# Download the current compromised-domain feed (free, no token required).
# dataset_type: url | ip | url-all | ip-all
#
# Usage: sh curl-example-compromised.sh

curl --fail --silent \
  "https://netapi.com/api2/?method=compromised&dataset_type=url" \
  --output compromised_url.csv

# Number of listed domains (the first line is a "#" comment):
grep -vc '^#' compromised_url.csv
