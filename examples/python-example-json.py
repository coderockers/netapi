#!/usr/bin/env python3
"""NetAPI JSON API (/api-json/): JSON answers for integrations and AI tools.

Free methods need no token; the paid ones (and higher limits) take the token as a Bearer header.
Standard library only.
"""

import json
import sys
import urllib.error
import urllib.parse
import urllib.request

TOKEN = "YOUR_API_TOKEN"
BASE = "https://netapi.com/api-json/"


def call(method: str, token: str | None = None, **params) -> dict:
    query = urllib.parse.urlencode({"method": method, **{k: v for k, v in params.items() if v is not None}})
    request = urllib.request.Request(BASE + "?" + query, headers={"Accept": "application/json"})
    if token:
        request.add_header("Authorization", f"Bearer {token}")
    try:
        with urllib.request.urlopen(request, timeout=60) as response:
            return json.load(response)
    except urllib.error.HTTPError as error:
        # errors are JSON too: {"error": {"code": "...", "message": "..."}}
        body = json.load(error)
        code = body.get("error", {}).get("code")
        message = body.get("error", {}).get("message")
        if error.code == 429:
            print(f"rate limited, retry after {body['error'].get('retry_after')} s", file=sys.stderr)
        raise SystemExit(f"{error.code} {code}: {message}")


# 1. compromised check (free)
check = call("compromised-check", value="example.com")
print("example.com listed now:", check["listed_now"], "| ever:", check["listed_ever"])

# 2. newly registered domains containing "crypto" (free)
new = call("search-new", q="crypto", days=3, limit=5)
for row in new["results"]:
    print(row["added_on"], row["domain"])
if new["truncated"]:
    print("more results available:", new["note"])

# 3. TLD statistics (free)
stats = call("tld-stats", tld="de")
print(".de domains:", stats["domains"], "| registry:", stats["registry"])

# 4. domain lookup (active plan required)
lookup = call("lookup-domain", token=TOKEN, domain="example.com")
print("nameservers:", lookup.get("dns"), "| expires:", lookup.get("expires_on"))
