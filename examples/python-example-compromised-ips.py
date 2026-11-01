"""Read the NetAPI compromised IP feed with threat data (JSON Lines).

The addresses were caught attacking NetAPI Nightwatch, NetAPI's own honeypot network.
The feed is free, needs no token and is rebuilt every hour. Each line is one JSON
object: ip, first_seen, last_seen, category (A infection attempt, B password guessing,
C exploit), top_threat, threats (type, name, family, cve, mitre, first, last, days),
country, asn, as_name, sensors.

The response is streamed line by line, nothing is buffered in memory.

Requires: requests (pip install requests)
"""

import json
from collections import Counter

import requests

API_URL = "https://netapi.com/api2/"

params = {
    "method": "compromised",
    "dataset_type": "ip",  # "ip" = last 24 hours, "ip-all" = every address ever listed
    "format": "jsonl",
    # Optional filters, the same as the list builder on https://netapi.com/compromised-ips/
    # "category": "B,C",
    # "country_not": "us",
    # "min_sensors": 2,
}

categories = Counter()
threats = Counter()
countries = Counter()
listed = set()

with requests.get(API_URL, params=params, stream=True, timeout=120) as resp:
    if resp.status_code != 200:
        # API errors are plain text, e.g. "405 Method Not Allowed: Invalid category: ..."
        raise SystemExit(f"HTTP {resp.status_code}: {resp.text}")
    for line in resp.iter_lines(decode_unicode=True):
        if not line:
            continue
        record = json.loads(line)
        listed.add(record["ip"])
        categories[record["category"]] += 1
        countries[record["country"] or "??"] += 1
        for threat in record["threats"]:
            threats[threat["type"]] += 1

print(f"{len(listed):,} addresses in the list")
print("By category:", dict(sorted(categories.items())))
print("Top attack types:", threats.most_common(5))
print("Top countries:", countries.most_common(5))

for ip in ("8.8.8.8", next(iter(listed), "")):
    if ip:
        print(ip, "-> LISTED" if ip in listed else "-> not listed")
