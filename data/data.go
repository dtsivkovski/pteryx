package data

import _ "embed"

// FileSigsJSON is Gary C. Kessler's file signature data.
//
// Source: https://www.garykessler.net/library/file_sigs_GCK_latest.html
// Copyright (c) 2002-2026 Gary C. Kessler. Used with attribution.
//
//go:embed file_sigs_normalized.json
var FileSigsJSON []byte
