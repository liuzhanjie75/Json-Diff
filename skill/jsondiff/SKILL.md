---
name: "jsondiff"
description: "Use when comparing JSON files or inline JSON values, explaining structural JSON differences, matching array objects by identity fields, filtering a JSON path, or deciding whether array order should matter."
---

# JSON Diff

Use the bundled `jsondiff` binary for deterministic JSON comparison. Do not
reimplement the diff in ad hoc scripts.

## Run

Resolve the executable relative to this `SKILL.md`:

- Windows: `bin/jsondiff.exe`
- Linux/macOS: `bin/jsondiff`

After resolving the directory containing this `SKILL.md` to an absolute path,
invoke the bundled binary explicitly.

Windows PowerShell:

```powershell
& "$skillDir\bin\jsondiff.exe" <source> <target> [options]
```

Linux/macOS:

```bash
"$skill_dir/bin/jsondiff" <source> <target> [options]
```

Each input may be an existing file path or an inline JSON value.

## Choose Options

| Need | Option |
|---|---|
| Compare only a sub-value | `--path "users.0.profile"` |
| Match array objects by identity | `--key "id"` |
| Treat `[1,2]` and `[2,1]` as equal | `--ignore-array-order` |
| Produce stable non-colored output | `--color never` |

Use `--key` when array objects have a reliable identity field. Do not invent a
key when none exists. Use `--ignore-array-order` only when order is not part of
the data's meaning; duplicate counts still matter.

## Interpret Results

- Exit `0`: equal under the selected options.
- Exit `1`: differences found; summarize the rendered paths and operations.
- Exit `2`: input, parsing, path, or argument error; report the error instead
  of describing the JSON as different.

`ADDED`, `REMOVED`, `CHANGED`, and `MOVED` describe the target relative to the
source. Object field order is always ignored. Array order is significant unless
`--ignore-array-order` is set.

## Common Mistakes

- Quote paths that contain spaces.
- Pass `--color never` when capturing output.
- Do not treat exit `1` as command failure; it is the normal "different" result.
- Do not use `--ignore-array-order` when array position carries meaning.
