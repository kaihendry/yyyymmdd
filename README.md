# yyyymmdd

Keep Downloads calm: loose files are placed into folders named after their modification date.

```text
Downloads/
├── 2026-08-25/
│   └── receipt.pdf
└── 2026-08-26/
    └── photo.jpg
```

## Install

```bash
go install github.com/kaihendry/yyyymmdd@latest
```

## Use

Preview the changes first (the default directory is `~/Downloads`):

```bash
yyyymmdd --dry-run
```

Then organise, with confirmation for each file:

```bash
yyyymmdd
```

Or organise everything without prompts, while leaving very recent downloads alone:

```bash
yyyymmdd --yes --older-than 10m
```

You can pass another directory as the final argument. Existing destination files are never overwritten, hidden files and subdirectories are ignored, and a dry run never changes the filesystem.

Run `yyyymmdd --help` for all options.

[Watch the original demo](http://www.youtube.com/watch?v=CYgu-N2xkwI).
