# dank-data Catalog Specification

This document defines the machine-readable contract between the `dank-data` repository and its consumers.

## Purpose

`snapshots/catalog.json` is the single source of truth for discovering and verifying dataset snapshots published in this repository. It is generated automatically by the repository's CI pipeline and must never be hand-edited. Consumers resolve a dataset identifier to a download URL and SHA-256 checksum via this file.

## Schema (Version 1)

```json
{
  "version": 1,
  "datasets": {
    "us/ct": {
      "description": "Cannabis brand registry, applications, credentials, lottery, retail locations, tax, weekly sales, and zoning.",
      "duckdb_url": "https://raw.githubusercontent.com/AgentDank/dank-data/main/snapshots/us/ct/dank-data.duckdb.zst",
      "sha256": "aabbccdd00112233...",
      "title": "United States — Connecticut",
      "updated_at": "2026-04-19T00:00:00Z"
    }
  }
}
```

### Top-level fields

| Field    | Required | Meaning                                                                   |
|----------|----------|---------------------------------------------------------------------------|
| version  | yes      | Catalog format version. Bump only on breaking changes. Consumers must reject unknown versions. |
| datasets | yes      | Map keyed by dataset id.                                                  |

### Per-dataset fields

| Field      | Required | Meaning                                                                   |
|------------|----------|---------------------------------------------------------------------------|
| description | yes     | One-sentence human description of what's in the dataset.                  |
| duckdb_url  | yes     | Absolute URL to the zstd-compressed DuckDB file (`dank-data.duckdb.zst`). Absolute so the file is self-describing when read independently. |
| sha256      | yes     | Hex-encoded SHA-256 of the compressed `.duckdb.zst` bytes. Verified by consumers after download. |
| title       | yes     | Short human-readable name.                                                |
| updated_at  | no      | ISO-8601 timestamp of the last snapshot. Informational only.              |

Unknown fields must be permitted at any level for forward compatibility. Consumers must ignore fields they do not recognize.

## Dataset Identifiers

Dataset ids are the keys of the `datasets` object. They must match the regular expression:

```
^[a-z]{2}/[a-z0-9_-]+$
```

- Two lowercase letters representing the region.
- A forward slash (`/`).
- One or more lowercase alphanumeric characters, hyphens, or underscores representing the dataset name.

This convention mirrors the directory layout under `snapshots/<region>/<dataset>/`.

## Per-Dataset Metadata

Human-readable fields (`title` and `description`) cannot be derived from the data files. They live in a per-dataset `metadata.json` co-located with the snapshot:

```
snapshots/<region>/<dataset>/metadata.json
```

### Schema

```json
{
  "title": "United States — Connecticut",
  "description": "Cannabis brand registry, applications, credentials, lottery, retail locations, tax, weekly sales, and zoning."
}
```

| Field       | Required | Meaning                                      |
|-------------|----------|----------------------------------------------|
| title       | yes      | Short human-readable name.                   |
| description | yes      | One-sentence description of the dataset.     |

Unknown fields are permitted for forward compatibility but are ignored by the current generator.

The generator reads `metadata.json` alongside `dank-data.duckdb.zst` when building `catalog.json`. If either file is missing or a required metadata field is empty, the generator exits with a non-zero status.

## Serialization Rules

The JSON output must be deterministic so that regenerating the catalog produces identical bytes when inputs have not changed. This keeps git diffs minimal and reviewable.

- Object keys are sorted alphabetically at every level (both the top-level object and each dataset entry).
- Dataset ids under `datasets` are sorted alphabetically.
- Indentation is exactly two spaces.
- Lines end with `\n` (Unix line endings).
- The file ends with a single trailing newline.
- No trailing whitespace on any line.

## Version Bump Policy

The `version` field is incremented only when a breaking change is made to the schema or serialization rules. Examples of breaking changes:

- Removing a required field.
- Changing the meaning of an existing field.
- Changing the dataset-id validation rules.

Adding new optional fields or permitting unknown fields is **not** a breaking change and does not require a version bump.

## Consumer Expectations

A compliant consumer should:

1. Fetch `https://raw.githubusercontent.com/AgentDank/dank-data/main/snapshots/catalog.json`.
2. Reject the catalog if `version` is not a recognized value.
3. Validate that the requested dataset id exists in `datasets`.
4. Validate the dataset id against the regex `^[a-z]{2}/[a-z0-9_-]+$` for path safety.
5. Download the file at `duckdb_url`.
6. Verify the downloaded bytes against `sha256` using SHA-256.
7. Decompress the verified bytes with ZStandard to obtain the DuckDB database file.

## Example

A minimal consumer in Python:

```python
import hashlib, json, urllib.request, zstandard

CATALOG_URL = "https://raw.githubusercontent.com/AgentDank/dank-data/main/snapshots/catalog.json"

def fetch_dataset(dataset_id: str) -> bytes:
    assert re.match(r"^[a-z]{2}/[a-z0-9_-]+$", dataset_id)

    with urllib.request.urlopen(CATALOG_URL) as r:
        catalog = json.load(r)

    if catalog.get("version") != 1:
        raise ValueError(f"Unsupported catalog version: {catalog.get('version')}")

    ds = catalog["datasets"].get(dataset_id)
    if not ds:
        raise ValueError(f"Dataset not found: {dataset_id}")

    with urllib.request.urlopen(ds["duckdb_url"]) as r:
        data = r.read()

    if hashlib.sha256(data).hexdigest() != ds["sha256"]:
        raise ValueError("SHA-256 mismatch")

    return zstandard.decompress(data)
```
