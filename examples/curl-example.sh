#!/usr/bin/env sh
# Download the full .net domain list (gzip-compressed CSV) from NetAPI.
# Replace YOUR_API_TOKEN with the token from https://netapi.com/dashboard/
#
# Usage: sh curl-example.sh

curl --fail --location \
  "https://netapi.com/api2/?method=download&zone_tld=net&dataset_type=list&filter_type=active&token=YOUR_API_TOKEN" \
  --output net_active_list.csv.gz

# Peek at the first lines without unpacking the file to disk:
gunzip -c net_active_list.csv.gz | head -n 5
