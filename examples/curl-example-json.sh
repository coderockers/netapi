#!/usr/bin/env bash
# NetAPI JSON API (/api-json/): JSON answers for integrations and AI tools.
# Free methods need no token; pass your token as a Bearer header for the paid ones and higher limits.
set -euo pipefail

TOKEN="YOUR_API_TOKEN"
BASE="https://netapi.com/api-json/"

# 1. Is a domain in the compromised feed? (free)
curl -sS "${BASE}?method=compromised-check&value=example.com"
echo

# 2. Domains registered in the last 3 days that contain "crypto" (free, 100 results per call without a token)
curl -sS "${BASE}?method=search-new&q=crypto&days=3&limit=5"
echo

# 3. Domain lookup (active plan required; Bearer token)
curl -sS -H "Authorization: Bearer ${TOKEN}" "${BASE}?method=lookup-domain&domain=example.com"
echo

# Errors come as {"error": {"code": "...", "message": "..."}} with the matching HTTP status,
# e.g. 429 {"error": {"code": "rate_limited", ..., "retry_after": 60}} when a limit is reached.
