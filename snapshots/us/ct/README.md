# United States — Connecticut

Cannabis brand registry, applications, credentials, lottery, retail locations, tax, weekly sales, and zoning.

Source: [data.ct.gov](https://data.ct.gov/). Snapshots are produced by [`dank-extract`](https://github.com/AgentDank/dank-extract) and published weekly.

## Tables

Each table is exported as both CSV and JSON, and is also a table in `dank-data.duckdb.zst`.

| Table | Files | Source |
|-------|-------|--------|
| `us_ct_applications` | `us_ct_applications.csv`, `us_ct_applications.json` | [Cannabis Application Status](https://data.ct.gov/Government/Cannabis-Application-Status/bqby-dyzr/about_data) |
| `us_ct_brands` | `us_ct_brands.csv`, `us_ct_brands.json` | [Medical Marijuana and Adult Use Cannabis Brand Registry](https://data.ct.gov/Health-and-Human-Services/Medical-Marijuana-and-Adult-Use-Cannabis-Brand-Reg/egd5-wb6r/about_data) |
| `us_ct_credentials` | `us_ct_credentials.csv`, `us_ct_credentials.json` | [Cannabis Credential Counts and Type](https://data.ct.gov/Government/Cannabis-Credential-Counts-and-Type/tjfe-s2x9/about_data) |
| `us_ct_lottery` | `us_ct_lottery.csv`, `us_ct_lottery.json` | data.ct.gov |
| `us_ct_retail_locations` | `us_ct_retail_locations.csv`, `us_ct_retail_locations.json` | data.ct.gov |
| `us_ct_tax` | `us_ct_tax.csv`, `us_ct_tax.json` | [Cannabis Tax Revenue](https://data.ct.gov/Government/Cannabis-Tax-Revenue/jey2-vq68/about_data) |
| `us_ct_weekly_sales` | `us_ct_weekly_sales.csv`, `us_ct_weekly_sales.json` | [Cannabis Retail Sales](https://data.ct.gov/Business/Cannabis-Retail-Sales/ucaf-96h6/about_data) |
| `us_ct_zoning` | `us_ct_zoning.csv`, `us_ct_zoning.json` | data.ct.gov |

## Files in this directory

 * `dank-data.duckdb.zst` — zstd-compressed DuckDB database containing every table above.
 * `*.csv` / `*.json` — per-table exports, stored uncompressed so Git delta-compresses them efficiently across snapshots.
 * `metadata.json` — title and description consumed by the catalog generator (see [`docs/catalog-spec.md`](../../../docs/catalog-spec.md)).

## Download

```bash
# DuckDB (all tables in one file)
curl -LO "https://github.com/AgentDank/dank-data/raw/main/snapshots/us/ct/dank-data.duckdb.zst"
zstd -d dank-data.duckdb.zst

# Or a single table as CSV / JSON (no decompression needed)
curl -LO "https://github.com/AgentDank/dank-data/raw/main/snapshots/us/ct/us_ct_brands.csv"
curl -LO "https://github.com/AgentDank/dank-data/raw/main/snapshots/us/ct/us_ct_brands.json"
```

The DuckDB file's SHA-256 is published in [`snapshots/catalog.json`](../../catalog.json) under dataset id `us/ct`. Verify downloads against it before use.

## Querying

```bash
duckdb dank-data.duckdb "SELECT COUNT(*) FROM us_ct_brands;"
```
