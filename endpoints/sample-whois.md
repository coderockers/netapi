# `sample-whois` — free sample of a registrar list or dataset

Returns the first rows of a registrar's domain list or detailed dataset so you can check the format before subscribing. The sample has the same columns as the file served by [`download-whois`](download-whois.md).

- **Authentication:** none
- **Response:** plain CSV (not compressed) with a header row and 10 data rows, served as a file download

## Request

```
GET https://netapi.com/api2/?method=sample-whois&registrar_id={id}&dataset_type={list|dataset}
```

| Parameter | Required | Description |
|---|---|---|
| `registrar_id` | yes | Numeric registrar ID from [`registrars`](registrars.md), e.g. `146` for GoDaddy.com, LLC. |
| `dataset_type` | yes | `list` — domains only · `dataset` — domains with registration data and metadata |

## Response

`dataset_type=list` (filename `{registrar_id}_list.csv`):

```
url
themis.abogado
abogada.abogado
```

`dataset_type=dataset` (filename `{registrar_id}_dataset.csv`):

```
registrar,url,registered_at,expiring_at,majestic_rank,emails,phones,ip,ip_country
"GoDaddy.com, LLC",themis.abogado,2025-03-12,2027-03-12,,info@themis.abogado,,34.174.62.87,US
"GoDaddy.com, LLC",abogada.abogado,2024-12-21,2026-12-21,,law@law.law,,64.23.224.186,US
```

Column descriptions are in [download-whois.md](download-whois.md#dataset-columns).

## Errors

| HTTP | Message |
|---|---|
| 405 | `405 Method Not Allowed: Missing Registrar ID.` |
| 405 | `405 Method Not Allowed: Invalid Registrar ID.` |
| 405 | `405 Method Not Allowed: Missing dataset type.` |
| 405 | `405 Method Not Allowed: Invalid dataset type.` |
| 404 | `404 Not Found: File not found. …` — no sample has been generated for this registrar yet |

## Examples

```bash
curl "https://netapi.com/api2/?method=sample-whois&registrar_id=146&dataset_type=dataset"
```

```python
import requests

resp = requests.get(
    "https://netapi.com/api2/",
    params={"method": "sample-whois", "registrar_id": 146, "dataset_type": "dataset"},
    timeout=30,
)
resp.raise_for_status()
print(resp.text)
```

## Notes

- Sample rows are taken from the live dataset and change when the registrar lists are rebuilt.
- There is no sample endpoint for DNS-provider lists; their columns are identical to the zone datasets, so use [`sample`](sample.md) to see the format.
