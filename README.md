<div align="center">
  <h1>pteryx</h1>
  <img width="655" height="268" alt="pteryx-banner" src="https://github.com/user-attachments/assets/36e03011-8e61-41a0-9848-bcd37c1c8b19" />
</div>

Pteryx is a tool for checking file signatures (magic numbers) to verify file types. It can be used to identify files with incorrect extensions or ones that are intentionally obscuring themselves.

## Installation

You have two options:

1. Install using brew

```bash
brew install dtsivkovski/tap/pteryx
```

2. Build from source 

```bash
git clone https://github.com/dtsivkovski/pteryx.git
cd pteryx
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

By default, `sig` lists failures and the summary. Add `-V` to include files that pass signature checks.

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

Any function can accept `-l` to enable logging. The log output file is currently labeled `pteryx.log` in the directory in which the command was run.

## Attribution

This project is developed by Daniel Tsivkovski and licensed under the MIT License. See [LICENSE](LICENSE) for more details.

### Outside Data Used

File signature data derived from [Gary C. Kessler's File Signature Table](https://www.garykessler.net/library/file_sigs_GCK_latest.html).

The original GCK file signature JSON is preserved in [data/file_sigs.json](https://github.com/dtsivkovski/pteryx/blob/main/data/file_sigs.json).
[data/file_sigs.normalized.json](https://github.com/dtsivkovski/pteryx/blob/main/data/file_sigs_normalized.json) contains parser-oriented cleanup of apparent
field formatting issues, such as Header offset values like "0(null)". Copyright © 2002-2026 Gary C. Kessler. Used with attribution.

## Why Pteryx?

I chose pteryx to name it after flying dinosaurs, but specifically after this cool one called the [Hatzegopteryx](https://en.wikipedia.org/wiki/Hatzegopteryx). I love how it's a flying dinosaur but comparable to the size of a giraffe.
