// generate-catalog.go generates snapshots/catalog.json from local snapshot data.
//
// Usage: go run scripts/generate-catalog.go
//
// Walks snapshots/<region>/<dataset>/, reads per-dataset metadata.json,
// computes SHA-256 of dank-data.duckdb.zst, and writes a deterministic
// snapshots/catalog.json with alphabetically sorted keys.
//
// Exits non-zero if any snapshot directory is missing required files or metadata.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	baseURL      = "https://raw.githubusercontent.com/AgentDank/dank-data/main/snapshots"
	snapshotsDir = "snapshots"
)

var datasetIDRe = regexp.MustCompile(`^[a-z]{2}/[a-z0-9_-]+$`)

// DatasetEntry is a single dataset in the catalog.
// Unknown fields are permitted for forward compatibility.
type DatasetEntry struct {
	Description string `json:"description"`
	DuckdbURL   string `json:"duckdb_url"`
	Sha256      string `json:"sha256"`
	Title       string `json:"title"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// Catalog is the top-level catalog structure.
type Catalog struct {
	Version  int                      `json:"version"`
	Datasets map[string]*DatasetEntry `json:"datasets"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}

	snapshotsPath := filepath.Join(repoRoot, snapshotsDir)
	entries, err := collectDatasets(snapshotsPath)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("no datasets found in %s", snapshotsPath)
	}

	// Sort dataset IDs alphabetically for deterministic output.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].id < entries[j].id
	})

	datasets := make(map[string]*DatasetEntry, len(entries))
	for _, e := range entries {
		datasets[e.id] = e.entry
	}

	catalog := &Catalog{
		Version:  1,
		Datasets: datasets,
	}

	// json.MarshalIndent sorts object keys alphabetically by default.
	b, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal catalog: %w", err)
	}

	// Ensure single trailing newline.
	os.Stdout.Write(b)
	os.Stdout.WriteString("\n")
	return nil
}

type collectedEntry struct {
	id    string
	entry *DatasetEntry
}

func collectDatasets(snapshotsPath string) ([]collectedEntry, error) {
	regions, err := os.ReadDir(snapshotsPath)
	if err != nil {
		return nil, fmt.Errorf("read snapshots dir: %w", err)
	}

	var entries []collectedEntry
	for _, region := range regions {
		if !region.IsDir() {
			continue
		}
		regionPath := filepath.Join(snapshotsPath, region.Name())
		datasets, err := os.ReadDir(regionPath)
		if err != nil {
			return nil, fmt.Errorf("read region dir %s: %w", region.Name(), err)
		}
		for _, dataset := range datasets {
			if !dataset.IsDir() {
				continue
			}
			id := region.Name() + "/" + dataset.Name()
			if !datasetIDRe.MatchString(id) {
				return nil, fmt.Errorf("invalid dataset id %q — must match %s", id, datasetIDRe.String())
			}

			datasetPath := filepath.Join(regionPath, dataset.Name())
			entry, err := buildEntry(id, datasetPath)
			if err != nil {
				return nil, fmt.Errorf("dataset %s: %w", id, err)
			}
			entries = append(entries, collectedEntry{id: id, entry: entry})
		}
	}
	return entries, nil
}

func buildEntry(id, datasetPath string) (*DatasetEntry, error) {
	duckdbFile := filepath.Join(datasetPath, "dank-data.duckdb.zst")
	if _, err := os.Stat(duckdbFile); err != nil {
		return nil, fmt.Errorf("missing required file %s", duckdbFile)
	}

	metadataFile := filepath.Join(datasetPath, "metadata.json")
	metadataRaw, err := os.ReadFile(metadataFile)
	if err != nil {
		return nil, fmt.Errorf("missing required metadata %s", metadataFile)
	}

	var metadata struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(metadataRaw, &metadata); err != nil {
		return nil, fmt.Errorf("parse metadata: %w", err)
	}
	if strings.TrimSpace(metadata.Title) == "" {
		return nil, fmt.Errorf("missing required field 'title' in %s", metadataFile)
	}
	if strings.TrimSpace(metadata.Description) == "" {
		return nil, fmt.Errorf("missing required field 'description' in %s", metadataFile)
	}

	hash, err := fileSHA256(duckdbFile)
	if err != nil {
		return nil, fmt.Errorf("compute sha256: %w", err)
	}

	info, err := os.Stat(duckdbFile)
	if err != nil {
		return nil, fmt.Errorf("stat duckdb file: %w", err)
	}

	return &DatasetEntry{
		Description: metadata.Description,
		DuckdbURL:   fmt.Sprintf("%s/%s/dank-data.duckdb.zst", baseURL, id),
		Sha256:      hash,
		Title:       metadata.Title,
		UpdatedAt:   info.ModTime().UTC().Format(time.RFC3339),
	}, nil
}

func fileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func findRepoRoot() (string, error) {
	// If invoked from the repo root, cwd is correct.
	// Also support being run from scripts/ by walking up.
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, snapshotsDir)); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("could not find repo root (no %s directory)", snapshotsDir)
}
