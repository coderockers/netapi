#!/usr/bin/env sh
# Download the current compromised-domain feed (free, no token required).
# dataset_type: url (current) | url-all (historical)
# Add &format=hosts for "0.0.0.0 domain" lines (hosts-based blockers).
#
# Usage: sh curl-example-compromised.sh

curl --fail --silent \
  "https://netapi.com/api2/?method=compromised&dataset_type=url" \
  --output compromised_url.csv

# Number of listed domains (the first line is a "#" comment):
grep -vc '^#' compromised_url.csv
