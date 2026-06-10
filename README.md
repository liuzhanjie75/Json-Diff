# jsondiff

[简体中文](README.zh-CN.md)

`jsondiff` is a Go CLI for comparing two JSON values. It supports recursive
field-level diffs, array move detection, object matching, JSON path filtering,
and colored terminal output.

## Features

- Recursive comparison of objects, arrays, primitives, and `null`
- Colored output for added, removed, changed, and moved values
- LCS-based array move detection
- Automatic array object matching using Jaccard key-set similarity
- Exact array object matching with `--key`
- JSON path filtering with `--path`
- File paths and inline JSON inputs
- Number preservation with `json.Number`

## Requirements

- Go 1.26.4 or later

## Build

```bash
go build -o jsondiff .
```

On Windows, the output is `jsondiff.exe`.

## Usage

```bash
# Compare two files
jsondiff old.json new.json

# Compare inline JSON values
jsondiff '{"a":1}' '{"a":2}'

# Compare a sub-path
jsondiff old.json new.json --path "database.connection"

# Match array objects by an exact identity field
jsondiff old.json new.json --key "id"

# Control terminal colors
jsondiff old.json new.json --color always
jsondiff old.json new.json --color never
```

When an argument names an existing file, `jsondiff` reads that file before
considering inline JSON syntax. This allows filenames such as `2024.json`.
Inputs must contain exactly one valid JSON value; trailing text and additional
JSON values are rejected.

## Options

| Option | Description | Default |
|---|---|---|
| `--path` | Compare only the selected GJSON path | Entire document |
| `--key` | Match array objects by an exact field value | Similarity matching |
| `--color` | Color mode: `auto`, `always`, or `never` | `auto` |

When `--key` is set, unmatched objects are not paired by similarity. Key values
also preserve their JSON types, so string `"1"` and number `1` are distinct.

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | The JSON values are equal |
| `1` | Differences were found |
| `2` | An input, parsing, path, or argument error occurred |

## Example Output

```text
[CHANGED]  $.version  : "1.0.0"  →  "2.0.0"
[ADDED]    $.config.verbose  : true
[REMOVED]  $.config.retries  : 3
[MOVED]    $.features[2]  [0] → [2]  : "cache"
```

## Matching Behavior

Array comparison runs in three stages:

1. Match unchanged elements with an LCS over canonical JSON hashes.
2. Match remaining objects either by the configured `--key` field or, when no
   key is configured, by Jaccard key-set similarity.
3. Detect remaining identical values at different indices as moves, then report
   all unmatched values as additions or removals.

The automatic similarity threshold is `0.5`.

## Development

```bash
go test ./...
go vet ./...
```

Tests live beside the packages they cover under `internal/`. The `diff` tests
use the same package so private algorithms can be tested without production
test wrappers.

## Dependencies

- [Cobra](https://github.com/spf13/cobra) for the CLI
- [fatih/color](https://github.com/fatih/color) for terminal colors
- [GJSON](https://github.com/tidwall/gjson) for path extraction

## License

MIT
