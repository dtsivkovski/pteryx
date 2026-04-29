# pteryx

Pteryx is a tool for checking file signatures (magic numbers) to verify file types. It can be used to identify files with incorrect extensions or ones that are intentionally obscuring themselves.

## Usage

```bash
pteryx [flags] <file or directory path>
```

## Attribution

This project is developed by Daniel Tsivkovski and licensed under the MIT License. See [LICENSE](LICENSE) for more details.

### Outside Data Used

File signature data derived from Gary C. Kessler's File Signature Table:
https://www.garykessler.net/library/file_sigs_GCK_latest.html

The original GCK file signature JSON is preserved in data/file_sigs.json.
data/file_sigs.normalized.json contains parser-oriented cleanup of apparent
field formatting issues, such as Header offset values like "0(null)".

Copyright © 2002-2026 Gary C. Kessler. Used with attribution.


