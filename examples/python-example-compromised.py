"""Read the NetAPI compromised-domain / compromised-IP feed.

The feed is free and needs no token. The first line is a "#" comment with the
list name and date; every other line is one domain or IP address.

Requires: requests (pip install requests)
"""

import requests

API_URL = "https://netapi.com/api2/"

params = {
    "method": "compromised",
    "dataset_type": "url",  # "url", "ip", "url-all" or "ip-all"
}

resp = requests.get(API_URL, params=params, timeout=120)
if resp.status_code != 200:
    raise SystemExit(f"HTTP {resp.status_code}: {resp.text}")

entries = [line.strip() for line in resp.text.splitlines() if line and not line.startswith("#")]
print(f"{len(entries):,} entries in the current list")
print("First 10:", entries[:10])

blocked = set(entries)
for domain in ("example.com", entries[0]):
    print(domain, "-> LISTED" if domain in blocked else "-> not listed")
