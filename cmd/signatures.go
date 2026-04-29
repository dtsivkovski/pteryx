package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"pteryx/data"
)

// signature and gck struct definitions
type Signature struct {
	Description string
	Header      []byte
	Offset      int
}

type SignatureMap map[string][]Signature

type gckFile struct {
	FileSigs []gckSignature `json:"filesigs"`
}

type gckSignature struct {
	Description  string `json:"File description"`
	HeaderHex    string `json:"Header (hex)"`
	Extensions   string `json:"File extension"`
	HeaderOffset string `json:"Header offset"`
	TrailerHex   string `json:"Trailer (hex)"`
}

var ( // load signatures and max read length on init
	fileSignatures = mustLoadSignatures()
	maxFileSignatureReadBytes = maxKnownSignatureReadLength(fileSignatures)
)

func mustLoadSignatures() SignatureMap {
	signatures, err := loadSignatures(data.FileSigsJSON)
	if err != nil {
		panic(err) // if data is invalid
	}

	return signatures
}

func loadSignatures(rawJSON []byte) (SignatureMap, error) {
	// get file and parse JSON
	var gck gckFile
	if err := json.Unmarshal(rawJSON, &gck); err != nil {
		return nil, fmt.Errorf("parse GCK file signatures: %w", err)
	}

	// convert to program signature format
	signatures := make(SignatureMap)
	for _, entry := range gck.FileSigs {
		// header hex
		header, err := parseHexBytes(entry.HeaderHex)
		if err != nil {
			return nil, fmt.Errorf("parse header for %q: %w", entry.Description, err)
		}
		if len(header) == 0 {
			continue
		}

		// get offset for headers that don't start at 0
		offset, err := parseHeaderOffset(entry.HeaderOffset)
		if err != nil {
			return nil, fmt.Errorf("parse header offset for %q: %w", entry.Description, err)
		}

		// create signature and append to map
		signature := Signature{
			Description: entry.Description,
			Header:      header,
			Offset:      offset,
		}

		for _, ext := range parseExtensions(entry.Extensions) {
			signatures[ext] = append(signatures[ext], signature)
		}
	}

	return signatures, nil
}

// parses hex string into byte array accounting for spaces
func parseHexBytes(hexString string) ([]byte, error) {
	var bytes []byte
	for _, field := range strings.Fields(hexString) {
		value, err := strconv.ParseUint(field, 16, 8) // parse string into hex byte
		if err != nil {
			return nil, err
		}

		bytes = append(bytes, byte(value))
	}

	return bytes, nil
}

// parses the offset string into an integer
func parseHeaderOffset(offsetString string) (int, error) {
	offsetString = strings.TrimSpace(offsetString)
	var digits strings.Builder

	// stop at non-digit characters
	for _, char := range offsetString {
		if char < '0' || char > '9' {
			break
		}

		digits.WriteRune(char)
	}

	if digits.Len() == 0 {
		return 0, fmt.Errorf("missing numeric offset in %q", offsetString)
	}

	return strconv.Atoi(digits.String())
}

// parses the extensions string into a into a slice of extensions
func parseExtensions(extensions string) []string {
	parts := strings.Split(extensions, "|")
	seen := make(map[string]bool)
	var parsed []string

	for _, part := range parts {
		ext := strings.TrimSpace(part)
		if ext == "" || strings.EqualFold(ext, "(none)") {
			continue
		}

		ext = "." + strings.TrimPrefix(strings.ToLower(ext), ".")
		if seen[ext] {
			continue
		}

		seen[ext] = true
		parsed = append(parsed, ext)
	}

	return parsed
}

// checks signature lengths to determine max bytes needed to read sigs
func maxSignatureReadLength(signatures []Signature) int {
	maxLength := 0
	for _, signature := range signatures {
		end := signature.Offset + len(signature.Header)
		if end > maxLength {
			maxLength = end
		}
	}

	return maxLength
}

// checks all signatures to determine max bytes needed to read sigs by ext
func maxKnownSignatureReadLength(signaturesByExt SignatureMap) int {
	maxLength := 0
	for _, signatures := range signaturesByExt {
		if length := maxSignatureReadLength(signatures); length > maxLength {
			maxLength = length
		}
	}

	return maxLength
}

// checks if fyle buffer matches the signature header
func signatureMatches(buffer []byte, signature Signature) bool {
	end := signature.Offset + len(signature.Header)
	if len(buffer) < end {
		return false
	}

	return bytes.Equal(buffer[signature.Offset:end], signature.Header)
}

func signatureSpecificityScore(signature Signature) int {
	score := len(signature.Header)

	// boost score for signatures that start with "ftyp" at offset 4
	if signature.Offset == 4 && len(signature.Header) >= 8 && bytes.Equal(signature.Header[:4], []byte("ftyp")) {
		score += 1000
	}

	// penalize score for signatures that start with "ftyp" at offset 0
	if signature.Offset == 0 && len(signature.Header) >= 8 && bytes.Equal(signature.Header[4:8], []byte("ftyp")) {
		score -= 1000
	}

	return score
}

func sortExtensions(extensions []string) {
	sort.Strings(extensions)
}
