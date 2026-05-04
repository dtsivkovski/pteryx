# pteryx

Pteryx is a tool for checking file signatures (magic numbers) to verify file types. It can be used to identify files with incorrect extensions or ones that are intentionally obscuring themselves.

## Installation

```go
go install
go build -o pteryx
```

## Usage

Check file extensions against magic number signatures
```bash
pteryx sig <file>
pteryx sig <directory> -d
pteryx sig <directory> -d -r
```

Create hash baseline
```bash
pteryx hash create <file>
pteryx hash create <directory> -d -r -o pteryx.hash
```

Compare files against saved hash baseline
```bash
pteryx hash compare <file> -i pteryx.hash
pteryx hash compare <directory> -d -r -i pteryx.hash
```

## Attribution

This project is developed by Daniel Tsivkovski and licensed under the MIT License. See [LICENSE](LICENSE) for more details.

### Outside Data Used

File signature data derived from [Gary C. Kessler's File Signature Table](https://www.garykessler.net/library/file_sigs_GCK_latest.html).

The original GCK file signature JSON is preserved in [data/file_sigs.json](https://github.com/dtsivkovski/pteryx/blob/main/data/file_sigs.json).
[data/file_sigs.normalized.json](https://github.com/dtsivkovski/pteryx/blob/main/data/file_sigs_normalized.json) contains parser-oriented cleanup of apparent
field formatting issues, such as Header offset values like "0(null)".

Copyright © 2002-2026 Gary C. Kessler. Used with attribution.

