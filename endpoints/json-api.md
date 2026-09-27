# JSON API — `/api-json/`

The same data as the CSV API, as JSON documents. Built for integrations and AI tools; the [NetAPI MCP server](https://netapi.com/help/mcp/) (`https://mcp.netapi.com/mcp`, source: [coderockers/netapi-mcp](https://github.com/coderockers/netapi-mcp)) is a thin layer over these methods.

```
GET https://netapi.com/api-json/?method=<method>&<parameters>
Authorization: Bearer YOUR_API_TOKEN        (optional for the free methods; or token=... in the query)
```

The method can also be the path: `GET https://netapi.com/api-json/tld-stats/?tld=de` is the same call (one path per method; [`openapi-gpt.yaml`](../openapi-gpt.yaml) describes the JSON API this way for GPT Actions and other OpenAPI tools).

Answers are JSON (`Content-Type: application/json`). Errors are JSON as well, with the matching HTTP status:

```json
{"error": {"code": "rate_limited", "message": "Too many requests: limit of 6 per minute for this method group.", "retry_after": 60}}
```

| HTTP | `code` | Meaning |
|---|---|---|
| 400 | `missing_parameter`, `invalid_value`, `invalid_domain`, `invalid_ip`, `invalid_tld`, `invalid_query`, `query_too_short`, `invalid_parameter` | a parameter is missing or wrong |
| 401 | `missing_token` | the method needs a token; `invalid_token` for an expired OAuth token |
| 403 | `invalid_token`, `no_plan`, `plan_too_low`, `payment_under_review` | the token is wrong or the plan does not cover the method |
| 404 | `unknown_method`, `invalid_tld`, `not_found`, `file_not_found`, `not_available` | nothing to return |
| 429 | `rate_limited` | a request limit is reached; `retry_after` seconds, `Retry-After` header |
| 503 | `unavailable` | the new-domains index is being rebuilt; retry later |

Public statistics answer with `Cache-Control: public, max-age=3600`; everything else is `no-store`. CORS is open for GET.

## Free methods (no token; a token raises the limits)

| method | parameters | answer |
|---|---|---|
| `compromised-check` | `value` = domain, hostname, URL or IP | `listed_now` (seen in the last 24 h), `listed_ever`, `first_seen`, `last_seen`, `block_reason`, `registrar`; parent domains of a hostname are tried too (`matched_value`) |
| `search-new` | `q` (4+ characters), `days` 1-7 (default 1), `tld`, `match` = `contains` \| `starts`, `limit` | `results[]` (`domain`, `tld`, `added_on`), `truncated` |
| `tld-stats` | `tld` | registry, `domains`, `new_24h`, `deleted_24h`, `change` over 7/30/90/365 days, `rank_by_size`, `compromised`, `top_1m`, `policies`, `description` |
| `domain-rank` | `domain` | `rank` in the NetAPI Top 1M, `rank_in_tld` |
| `top-websites` | `tld`, `limit`, `offset` (offset + limit ≤ 1,000) | most popular domains of the TLD with global rank |
| `top-1m` | `limit`, `offset` | slice of the Top 1M |
| `registrar-info` | `q` = name, brand or id | `domains`, `rank_by_size`, `abuse` (compromised domains, rate, overall rate), `page`, `other_matches` |
| `dns-provider-info` | `q` = alias, name or nameserver root | `domains`, `market_share` (share, rank, 30-day and 1-year change), `page`, `other_matches` |
| `zones` | `min_domains`, `cctld` = `1` \| `0` | every zone: `domains`, `new_24h`, `deleted_24h`, `updated_daily`, `page` |
| `me` | token required | `plan`, `tier`, `limits`, `auth` (`api_token` or `oauth`) |

Examples:

```
https://netapi.com/api-json/?method=compromised-check&value=example.com
https://netapi.com/api-json/?method=search-new&q=crypto&days=3&limit=20
https://netapi.com/api-json/?method=tld-stats&tld=de
```

## Methods that need an active plan

| method | parameters | answer |
|---|---|---|
| `lookup-domain` | `domain` | `dns[]`, `hostname`, `ip`, `ip_country`, `registered_on`, `expires_on`, `registrar` (`found: false` when unknown) |
| `lookup-ip` | `ip` | `domains[]` (up to 3) with `hostname` and `dns[]` |
| `download-url` | `zone_tld` (or `all-zones`), `dataset_type` = `list` \| `dataset`, `filter_type` = `active` \| `new` \| `deleted`, `format` = `gz` \| `plain` | `download_url` valid 24 hours, `file_name`, `size_bytes`, `file_date` — the same rules as [download](download.md) (datasets need Plus or Pro; deleted lists only for zones whose registry publishes the complete zone daily) |

## Request limits

Counted per account (per IP address without a token), reset at 00:00 UTC, shared with the CSV API and the MCP server.

| Access | Free methods | `top-websites`, `top-1m` | `search-new` | Lookups | Downloads |
|---|---|---|---|---|---|
| No token | 30 / min, 1,000 / day | 10 / min, 200 / day | 6 / min, 100 / day, 100 rows | — | — |
| Token, no plan | 60 / min, 5,000 / day | 30 / min, 1,000 / day | 30 / min, 500 / day, 500 rows | — | — |
| Basic | 60 / min | 60 / min | 2,000 / day, 1,000 rows | 2,000 / day | 300 / day |
| Plus | 60 / min | 60 / min | 5,000 / day, 5,000 rows | 5,000 / day | 1,000 / day |
| Pro | 120 / min | 120 / min | no daily limit | no daily limit | no daily limit |

Examples in [examples/](../examples/): `curl-example-json.sh`, `python-example-json.py`, `node-example-json.js`, `go-example-json.go`, `php-example-json.php`.
