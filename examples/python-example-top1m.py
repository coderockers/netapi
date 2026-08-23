"""Load the NetAPI Top 1 Million ranking and look up a few domains.

The list is free (CC BY 4.0), needs no token and is rebuilt daily:
https://netapi.com/top-1-mln-websites/

Requires: requests (pip install requests)
"""

import csv
import io

import requests

TOP1M_URL = "https://netapi.com/netapi_top1mln.csv"

resp = requests.get(TOP1M_URL, timeout=120)
resp.raise_for_status()

rank_by_domain = {
    row["Domain"]: int(row["NetAPI_Rank"])
    for row in csv.DictReader(io.StringIO(resp.text))
}
print(f"{len(rank_by_domain):,} domains loaded")

for domain in ("google.com", "wikipedia.org", "github.com", "example.com"):
    rank = rank_by_domain.get(domain)
    print(f"{domain:20} {rank if rank else 'not in Top 1M'}")
