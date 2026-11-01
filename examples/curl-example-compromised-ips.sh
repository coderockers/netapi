#!/usr/bin/env sh
# Compromised IP addresses caught by NetAPI Nightwatch, NetAPI's own honeypot network
# (free, no token required, rebuilt every hour).
#
# 1. The live list (last 24 hours) as an nftables set, ready for `nft -f`.
# 2. A filtered list with threat data: addresses that guessed passwords (category B)
#    in the last 7 days and were seen by at least two sensors, as CSV.
#
# Formats: plain, csv, jsonl, ipset, nftables, iptables, mikrotik, cidr.
# Usage: sh curl-example-compromised-ips.sh

API="https://netapi.com/api2/?method=compromised"

# --fail-with-body keeps the plain-text error message of the API on HTTP errors
curl --fail-with-body --silent --show-error \
  "$API&dataset_type=ip&format=nftables" \
  --output netapi-compromised.nft || exit 1
head -n 1 netapi-compromised.nft
# Load it (needs root):  nft -f netapi-compromised.nft
# Then drop the traffic: nft add rule inet netapi input ip saddr @compromised_ips drop
#   (after creating a chain "input" with: nft add chain inet netapi input '{ type filter hook input priority 0; }')

curl --fail-with-body --silent --show-error \
  "$API&dataset_type=ip-all&window=7d&category=B&min_sensors=2&format=csv" \
  --output netapi-compromised-b.csv || exit 1

# Number of addresses (the first line is the CSV header) and the first rows:
echo "$(($(wc -l < netapi-compromised-b.csv) - 1)) addresses"
head -n 4 netapi-compromised-b.csv
