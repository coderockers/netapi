# `zones` — list of supported domain zones

Returns every TLD available in NetAPI together with two flags: whether the zone is refreshed daily and whether it is a country-code TLD.

Use it to discover valid `zone_tld` values for [`download`](download.md) and [`compromised-zone`](compromised-zone.md), and to find out which zones support the `new` and `deleted` filters.

- **Authentication:** none
- **Parameters:** none
- **Response:** plain CSV with a header row, about 1,500 lines

## Request

```
GET https://netapi.com/api2/?method=zones
```

## Response

```
zoneTLD,isUpdatedDaily,isCountryCode
```

| Column | Type | Description |
|---|---|---|
| `zoneTLD` | string | The TLD without the leading dot: `com`, `de`, `co.uk`. Internationalised TLDs are given in punycode, e.g. `xn--p1ai` for `.рф`. |
| `isUpdatedDaily` | `0` / `1` | `1` — the zone is refreshed every day and supports the `new` and `deleted` filters in `download`; `0` — the zone is refreshed monthly and supports `active` only. |
| `isCountryCode` | `0` / `1` | `1` — country-code TLD (ccTLD); `0` — generic TLD (gTLD). |

Example (the order of the lines is not significant):

```
zoneTLD,isUpdatedDaily,isCountryCode
com,1,0
de,1,1
xn--p1ai,1,1
io,0,1
```

## Examples

cURL:

```bash
curl "https://netapi.com/api2/?method=zones"
```

Python — collect the zones that support daily `new` lists:

```python
import csv
import io
import requests

resp = requests.get("https://netapi.com/api2/", params={"method": "zones"}, timeout=30)
resp.raise_for_status()

daily_zones = [
    row["zoneTLD"]
    for row in csv.DictReader(io.StringIO(resp.text))
    if row["isUpdatedDaily"] == "1"
]
print(len(daily_zones), "zones are updated daily")
```

## Notes

- Counts of registered domains per zone are not part of this response; they are shown on the zone pages, e.g. [netapi.com/com-domains/](https://netapi.com/com-domains/).
- The same list can be downloaded as a file from the [API reference](https://netapi.com/help/api/#zone-api).
