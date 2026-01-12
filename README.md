# dank-data

`dank-data` is AgentDank's public dataset repository.  We take public cannabis datasets, openly clean them, and publish them for public benefit.

This is the one of the data sources for [AgentDank](https://agentdank.com).  You can connect with it yourself using the [AgentDank MCP server](https://github.com/AgentDank/dank-mcp) named `dank-mcp`].

The data is sourced and cleaned using [`dank-extract`](https://github.com/AgentDank/dank-extract), our CLI tool for fetching, cleaning, and exporting cannabis data.

 * [Snapshots](#snapshots)
 * [Latest Snapshot](#latest-snapshot)
 * [Data Cleaning](#data-cleaning)
 * [Contribution and Conduct](#contribution-and-conduct)
 * [Credits and License](#credits-and-license)


Currently the following datasets are snapshotted:

 * [US CT Medical Marijuana and Adult Use Cannabis Brand Registry](https://data.ct.gov/Health-and-Human-Services/Medical-Marijuana-and-Adult-Use-Cannabis-Brand-Reg/egd5-wb6r/about_data)
 * [US CT Cannabis Credentials](https://data.ct.gov/Government/Cannabis-Credential-Counts-and-Type/tjfe-s2x9/about_data)
 * [US CT Cannabis Applications](https://data.ct.gov/Government/Cannabis-Application-Status/bqby-dyzr/about_data)
 * [US CT Cannabis Weekly Sales](https://data.ct.gov/Business/Cannabis-Retail-Sales/ucaf-96h6/about_data)
 * [US CT Cannabis Tax Revenue](https://data.ct.gov/Government/Cannabis-Tax-Revenue/jey2-vq68/about_data)

## Snapshots

We have included our ["cleaned"](#data-cleaning) data snapshots of public datasets in the repository for the researcher's convenience. These may be found in the [`snapshots`](./snapshots/) directory. Due to their size, the CSV and JSON and [DuckDB](https://duckdb.org) files are compressed with [ZStandard, `.zst`](https://en.wikipedia.org/wiki/Zstd).

Snapshots are taken **weekly** at 4:20PM Pacific via [GitHub Actions](https://github.com/AgentDank/dank-data/actions).

### Latest Snapshot

The [`snapshots/us/ct/latest`](./snapshots/us/ct/latest) symlink always points to the most recent snapshot.

```bash
# Download latest DuckDB directly via symlink
curl -LO "https://github.com/AgentDank/dank-data/raw/main/snapshots/us/ct/latest/dank-data.duckdb.zst"

# Decompress
zstd -d dank-data.duckdb.zst
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
