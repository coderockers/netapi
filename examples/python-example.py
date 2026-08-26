"""Download a zone dataset from NetAPI and stream it row by row.

The response is a gzip-compressed CSV. This example decompresses it on the
fly, so even the .com dataset can be processed without buffering it in memory
or unpacking it to disk.

Requires: requests (pip install requests)
"""

import csv
import gzip

import requests

API_URL = "https://netapi.com/api2/"
API_TOKEN = "YOUR_API_TOKEN"  # https://netapi.com/dashboard/

params = {
    "method": "download",
    "zone_tld": "net",          # any TLD from ?method=zones, or "all-zones"
    "dataset_type": "dataset",  # "list" (domains only) or "dataset" (with metadata)
    "filter_type": "active",    # "active" or "new"
    "token": API_TOKEN,
}

with requests.get(API_URL, params=params, stream=True, timeout=600) as resp:
    if resp.status_code != 200:
        # Errors come back as a short plain-text message, e.g.
        # "403 Forbidden: No active/paid plan."
        raise SystemExit(f"HTTP {resp.status_code}: {resp.text}")

    with gzip.open(resp.raw, mode="rt", encoding="utf-8", newline="") as f:
        reader = csv.DictReader(f)
        print("Columns:", reader.fieldnames)

        for i, row in enumerate(reader, start=1):
            print(row["url"], row["dns1"], row["ip"], row["ip_country"])
            if i == 10:
                break  # remove to process the whole file
