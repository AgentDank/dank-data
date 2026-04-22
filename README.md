# dank-data

`dank-data` is AgentDank's public dataset repository.  We take public cannabis datasets, openly clean them, and publish them for public benefit.

This is one of the data sources for [AgentDank](https://agentdank.com).  You can connect to it yourself using the [AgentDank MCP server](https://github.com/AgentDank/dank-mcp) named `dank-mcp`.

The data is sourced and cleaned using [`dank-extract`](https://github.com/AgentDank/dank-extract), our CLI tool for fetching, cleaning, and exporting cannabis data.

 * [Snapshots](#snapshots)
 * [Data Cleaning](#data-cleaning)
 * [Contribution and Conduct](#contribution-and-conduct)
 * [Credits and License](#credits-and-license)

## Datasets

| Dataset | Location | README |
|---------|----------|--------|
| United States — Connecticut | [`snapshots/us/ct/`](./snapshots/us/ct/) | [README](./snapshots/us/ct/README.md) |

Each dataset directory has its own `README.md` describing source tables, files, and download instructions.

## Snapshots

Data snapshots are stored in [`snapshots/`](./snapshots/), organized as `snapshots/<region>/<dataset>/`. CSV and JSON files are stored uncompressed for efficient Git delta compression. Only [DuckDB](https://duckdb.org) files are compressed with [ZStandard](https://en.wikipedia.org/wiki/Zstd) (`.zst`).

Snapshots are updated **weekly** via [GitHub Actions](https://github.com/AgentDank/dank-data/actions). Git history provides access to previous snapshots.

### Catalog

[`snapshots/catalog.json`](./snapshots/catalog.json) is the machine-readable index of all published datasets. It maps dataset identifiers to download URLs and SHA-256 checksums. Consumers such as the [`dank-mcp`](https://github.com/AgentDank/dank-mcp) server use it to discover and verify snapshots.

The catalog is **generated** — do not edit it by hand. Run `go run scripts/generate-catalog.go` locally to refresh it, or let CI regenerate it automatically with each snapshot.

See [`docs/catalog-spec.md`](./docs/catalog-spec.md) for the full specification (schema, version policy, and consumer expectations).

### Download

See the per-dataset README for download examples. In general, a snapshot directory contains:

 * `dank-data.duckdb.zst` — zstd-compressed DuckDB database (decompress with `zstd -d`)
 * `*.csv` / `*.json` — per-table exports, stored uncompressed
 * `metadata.json` — human-readable title and description used by the catalog generator

### Historical Data

Git history is your time machine:

```bash
git clone https://github.com/AgentDank/dank-data.git
cd dank-data

# View snapshot history
git log --oneline -- snapshots/

# Get data from a specific date
git checkout $(git rev-list -n1 --before="2026-01-01" main) -- snapshots/
```

## Data Cleaning

Because the upstream datasets are not perfect, we have to "clean" them. Such practices are opinionated, but since this is Open Source, you can inspect how we clean the data by examining [`dank-extract`](https://github.com/AgentDank/dank-extract).

Generally, we simply remove weird characters. We treat detected "trace" amounts as 0 -- and even with that, there's a decision about what field entry actually means "trace".

We also remove rows with ridiculous data -- as much as I'd love a `90,385% THC product`, I don't think it really exists. In that case, it was an incorrect decimal point (by looking at the label picture), but not every picture validates every error. It's also a pain -- but maybe a multi-modal vision LLM can help with that?

---

## Contribution and Conduct

As [AgentDank](https://github.com/AgentDank) curates these datasets, pull requests are generally not welcome here.  If you want to contribute, please create an issue.

Either way, obey our [Code of Conduct](./CODE_OF_CONDUCT.md).

## Credits and License

Copyright (c) 2026 Neomantra Corp. Authored by Evan Wies for [AgentDank](https://github.com/agentdank).

Released under the [MIT License](https://en.wikipedia.org/wiki/MIT_License), see [LICENSE.txt](./LICENSE.txt).

----
Made with :herb: and :fire: by the team behind [AgentDank](https://github.com/AgentDank).
