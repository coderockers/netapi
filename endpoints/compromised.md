# `compromised` — compromised IP addresses and domains (free)

Two free threat-intelligence feeds behind one method:

- **Compromised IPs** (`dataset_type=ip`, `ip-all`): IP addresses caught attacking [NetAPI Nightwatch](https://netapi.com/compromised-ips/), NetAPI's own honeypot network (password guessing, exploits, malware downloads). Every address on the list was seen attacking. Rebuilt every hour. Available as a plain list, as CSV or JSON Lines with threat data, in firewall formats, and as filtered custom lists.
- **Compromised domains** (`dataset_type=url`, `url-all`): domain names reported as compromised or malicious by public feeds and notifications. Rebuilt daily. Available as a plain list or in hosts-file format.

- **Authentication:** none
- **Response:** plain text (or CSV / JSON Lines), served as a file download
- **License:** [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/), free for any use with attribution to NetAPI

## Request

```
GET https://netapi.com/api2/?method=compromised&dataset_type={ip|ip-all|url|url-all}
```

| `dataset_type` | Contents | Size (Nov 2026) |
|---|---|---|
| `ip` | IP addresses that attacked NetAPI Nightwatch in the last 24 hours, the **live** blocklist | ~7,500 |
| `ip-all` | Every IP address ever listed in the feed | ~12,000 |
| `url` | Domains seen on threat feeds in the last 24 hours, the **current** blocklist | ~143,000 |
| `url-all` | Every domain ever listed since the feed started | ~1.2 million |

Statistics, a map of attack origins, a list builder and the attack types of the IP feed are on [netapi.com/compromised-ips/](https://netapi.com/compromised-ips/). The domain feed is on [netapi.com/compromised-urls/](https://netapi.com/compromised-urls/).

## Compromised IPs

### Plain list

Without `format` (or with `format=plain`) the response is one IPv4 address per line after a `#` comment line. There is no CSV header.

`dataset_type=ip` (filename `compromised_ip.csv`):

```
# NetAPI.com - currently compromised IPs (2026-11-01).
1.10.132.54
1.10.214.100
```

`dataset_type=ip-all` (filename `compromised_ip_history.csv`) starts with `# NetAPI.com - compromised IP addresses (historical).`

### Output formats

| `format` | Contents | Filename |
|---|---|---|
| `plain` (default) | One IP address per line, as above | `compromised_ip.csv`, `compromised_ip_history.csv` |
| `csv` | One row per address with threat data, header row first (fields below) | `compromised_ips_live.csv`, `compromised_ips_history.csv` |
| `jsonl` | One JSON object per address (JSON Lines, fields below) | `compromised_ips_live.jsonl`, `compromised_ips_history.jsonl` |
| `ipset` | `ipset save` layout for the set `netapi-compromised`. Load with `ipset restore -exist < file` | `compromised_ips_custom.ipset` |
| `nftables` | The set `compromised_ips` in the table `inet netapi`. Load with `nft -f file`, drop with a rule such as `ip saddr @compromised_ips drop` | `compromised_ips_custom.nft` |
| `iptables` | The chain `NETAPI-COMPROMISED` with one DROP rule per address. Load with `iptables-restore --noflush < file` and jump to the chain from INPUT | `compromised_ips_custom.rules` |
| `mikrotik` | A RouterOS script that replaces the address list `netapi-compromised`. Run it with `/import` | `compromised_ips_custom.rsc` |
| `cidr` | The addresses collapsed into the fewest CIDR blocks, one per line | `compromised_ips_custom.txt` |

Any other `format` value is ignored and the plain list is returned, as before. `format=hosts` is for the domain lists only (405 on the IP lists).

### Filters

The parameters of the [list builder](https://netapi.com/compromised-ips/#build). Combine them freely with any format.

| Parameter | Values |
|---|---|
| `window` | With `dataset_type=ip-all`: `7d` or `30d`, the addresses that attacked within the last 7 or 30 days. `dataset_type=ip` is always the 24-hour list |
| `category` | `A` infection attempts, `B` password guessing, `C` exploits, `D` form spam (messages posted through the contact or comment forms of sites that have no real visitors). Several separated by commas: `category=B,C` |
| `threat` | Attack types, for example `ssh-bruteforce,rdp-bruteforce`. An address matches when it did one of them within the window |
| `country` | Two-letter country codes of the addresses to keep, for example `country=cn,ru` |
| `country_not` | Two-letter country codes of the addresses to leave out. Use `country` or `country_not`, not both |
| `min_sensors` | Only addresses seen by at least this many sensors, `1` to `100` |

Common attack types: `ms17-010-exploit`, `ssh-bruteforce`, `telnet-bruteforce`, `mssql-bruteforce`, `rdp-bruteforce`, `wp-login-bruteforce`, `mirai-loader`, `malware-download`, `bot-dropper`, `shell-recon`, `ssh-key-implant-mdrfckr`, `env-file-hunting`, `git-exposure`. The full set (about 80 types) appears in the `threats` field of the CSV / JSON Lines lists.

Example, addresses that guessed passwords in the last 7 days and were seen by at least two sensors, as an nftables set:

```
https://netapi.com/api2/?method=compromised&dataset_type=ip-all&window=7d&category=B&min_sensors=2&format=nftables
```

A filtered list and the `nftables`, `iptables`, `mikrotik` and `cidr` formats start with a comment line such as:

```
# NetAPI.com - compromised IPs, custom list (2026-11-01 16:20 UTC), 1549 addresses. Source: https://netapi.com/compromised-ips/ (CC BY 4.0).
```

The `ipset` format has no comment line (`ipset restore` does not accept one).

Every link returns the list as of the last hourly rebuild, so a cron job or a firewall can poll it, for example once an hour.

### CSV and JSON Lines fields

| Field | Meaning |
|---|---|
| `ip` | IPv4 address |
| `first_seen`, `last_seen` | First and latest attack on our sensors, UTC, ISO 8601 |
| `category` | `A` infection attempt, `B` password guessing, `C` exploit, `D` form spam (the most severe one seen) |
| `top_threat` | The most severe attack type of the address |
| `threats` | CSV: attack types joined with `\|`. JSON: objects with `type`, `name`, `family`, `cve`, `mitre` (MITRE ATT&CK technique IDs), `first`, `last`, `days` |
| `country`, `asn`, `as_name` | Country code and network of the address |
| `sensors` | Number of NetAPI Nightwatch sensors that saw it |

CSV:

```
ip,first_seen,last_seen,category,top_threat,threats,country,asn,as_name,sensors
1.10.132.54,2026-11-01T05:11:30Z,2026-11-01T05:14:18Z,C,ms17-010-exploit,ms17-010-exploit,TH,23969,"TOT-NET - TOT Public Company Limited",1
```

JSON Lines (one object per line, shown wrapped here):

```json
{"ip":"1.10.132.54","category":"C","top_threat":"ms17-010-exploit",
 "first_seen":"2026-11-01T05:11:30Z","last_seen":"2026-11-01T05:14:18Z",
 "country":"TH","asn":23969,"as_name":"TOT-NET - TOT Public Company Limited","sensors":1,
 "threats":[{"type":"ms17-010-exploit","name":"EternalBlue (MS17-010) exploit","family":"service-exploit",
   "cve":"CVE-2017-0144","mitre":["T1210","T1190"],"first":"2026-11-01T05:11:30Z","last":"2026-11-01T05:14:18Z","days":1}]}
```

## Compromised domains

`dataset_type=url` (filename `compromised_url.csv`):

```
# NetAPI.com - currently compromised URLs (2026-11-01).
000359.xyz
00zyku.com
```

`dataset_type=url-all` (filename `compromised_url_history.csv`) starts with `# NetAPI.com - compromised URLs (historical).`

`format=hosts` returns the domain lists in hosts-file format, one `0.0.0.0 domain` line per domain, for hosts-based blockers (dnsmasq `addn-hosts`, AdAway, Pi-hole). Filenames `compromised_domains_hosts.txt` and `compromised_domains_history_hosts.txt`:

```
https://netapi.com/api2/?method=compromised&dataset_type=url&format=hosts
```

A domain listed by mistake can be reported through the [removal form](https://netapi.com/compromised-urls/removal/).

## Errors

All errors are plain text.

| HTTP | Message |
|---|---|
| 405 | `405 Method Not Allowed: Missing list type.` (`dataset_type` is missing) |
| 405 | `405 Method Not Allowed: Invalid list type.` |
| 405 | `405 Method Not Allowed: format=hosts is available for the domain lists only (dataset_type=url or url-all).` |
| 405 | `405 Method Not Allowed: window=7d and window=30d work with dataset_type=ip-all (dataset_type=ip is the 24-hour list).` |
| 405 | `405 Method Not Allowed: Invalid window: use 24h, 7d or 30d.` |
| 405 | `405 Method Not Allowed: Invalid category: use A, B, C or D, separated by commas.` (likewise for `threat`, `country`, `country_not`) |
| 405 | `405 Method Not Allowed: Use country or country_not, not both.` |
| 405 | `405 Method Not Allowed: Invalid min_sensors: a number from 1 to 100.` |
| 503 | `503 Service Unavailable: The list is being rebuilt, try again in a few minutes.` |

## Examples

cURL, the live IP list as an nftables set:

```bash
curl --fail -s "https://netapi.com/api2/?method=compromised&dataset_type=ip&format=nftables" -o netapi.nft
sudo nft -f netapi.nft
```

Python, today's domain blocklist into a set:

```python
import requests

resp = requests.get(
    "https://netapi.com/api2/", params={"method": "compromised", "dataset_type": "url"}, timeout=120
)
resp.raise_for_status()

blocked = {line.strip() for line in resp.text.splitlines() if line and not line.startswith("#")}
print(len(blocked), "domains currently listed")
```

Complete examples in five languages: [examples/](../examples/) (`*-example-compromised.*` for the domain feed, `*-example-compromised-ips.*` for the IP feed with threat data).

## Notes

- The IP feed lists single IPv4 addresses only, no network ranges. `format=cidr` merges neighbouring addresses into blocks.
- An address stays on the live list (`ip`) for 24 hours after its last attack. Scanners of known research and security companies are not listed.
- The `*-all` lists are historical. An address or domain there was listed at some point and is not necessarily malicious now.
- The `url` lists contain host names only (registered domains and, where reported, subdomains), no schemes or paths.
- To check a single domain or IP address, use [`compromised-check`](json-api.md) of the JSON API. Per-zone domain lists: [`compromised-zone`](compromised-zone.md). Abuse research by TLD and registrar: [netapi.com/research/](https://netapi.com/research/).
- Web reference: [netapi.com/help/api/#compromised-ips-api](https://netapi.com/help/api/#compromised-ips-api) and [#compromised-domains-api](https://netapi.com/help/api/#compromised-domains-api).
