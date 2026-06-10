# jsondiff

[简体中文](README.zh-CN.md)

`jsondiff` is a Go CLI for comparing two JSON values. It supports recursive
field-level diffs, array move detection, object matching, JSON path filtering,
and colored terminal output.

## Features

- Recursive comparison of objects, arrays, primitives, and `null`
- Colored output for added, removed, changed, and moved values
- LCS-based array move detection
- Optional recursive array comparison without considering element order
- Automatic array object matching using Jaccard key-set similarity
- Exact array object matching with `--key`
- JSON path filtering with `--path`
- File paths and inline JSON inputs
- Number preservation with `json.Number`

## Requirements

- Go 1.26.4 or later

## Build

Windows:

```powershell
.\build.bat
```

Linux and macOS:

```bash
chmod +x build.sh
./build.sh
```

The scripts can be called from any working directory and always build in the
project root. `build.bat` updates `jsondiff.exe`; `build.sh` produces
`jsondiff`, or `jsondiff.exe` when Go targets Windows.

The equivalent direct Go command is:

```bash
go build -o jsondiff .
```

On Windows, use `go build -o jsondiff.exe .`.

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

# Treat arrays as unordered multisets
jsondiff old.json new.json --ignore-array-order

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
| `--ignore-array-order` | Compare arrays without considering element order | `false` |
| `--color` | Color mode: `auto`, `always`, or `never` | `auto` |

When `--key` is set, unmatched objects are not paired by similarity. Key values
also preserve their JSON types, so string `"1"` and number `1` are distinct.

With `--ignore-array-order`, arrays use recursive multiset semantics:

- `[1, 2]` and `[2, 1]` are equal.
- Duplicate counts remain significant.
- Nested arrays also ignore order.
- Move differences are suppressed; only changes, additions, and removals remain.

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

When `--ignore-array-order` is enabled, the LCS and move-detection stages are
replaced by recursive unordered matching. The default mode remains
order-sensitive.

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
