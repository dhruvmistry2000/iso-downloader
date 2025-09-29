## iso-downloader (Go)

CLI/TUI to quickly download Linux ISOs, organized by distro family → distro → versions.

### Features

- Select family → distro → versions via interactive prompts
- Directory picker: uses `yazi --chooser-dir` if available, else simple prompt
- Organized downloads: `downloads/<distro>/<filename>.iso`
- Config is fetched from GitHub JSON on every run (no local config required)
- Supports Debian family (Debian, Ubuntu), Fedora, Arch (extensible)

### Run (no install)

```bash
curl -fsSL https://raw.githubusercontent.com/dhruvmistry2000/iso-downloader/refs/heads/main/start.sh | bash
```

Use a different repo fork:

```bash
curl -fsSL https://raw.githubusercontent.com/dhruvmistry2000/iso-downloader/refs/heads/main/start.sh | REPO=<owner>/<repo> bash
```

### How it works

- The app fetches `data/distros.json` from the repo at runtime by default.
- Override with `ISO_DOWNLOADER_CONFIG` to point to a different JSON URL.
- For distros with dynamic filenames, it lists directory contents and matches a glob.

### Build locally

```bash
bash scripts/build_run.sh
```

This builds `dist/iso-downloader` and runs it.

### JSON schema (simplified)

```json
{
  "families": {
    "debian": {
      "distros": {
        "debian": { "base_url": "...", "list_url_template": "...", "filename_glob": "...", "versions": ["..."] },
        "ubuntu": { "base_url": "...", "url_template": "...", "filename_template": "...", "versions": ["..."] }
      }
    },
    "fedora": { "distros": { "fedora": { "base_url": "...", "list_url_template": "...", "filename_glob": "...", "versions": ["..."] } } },
    "arch":   { "distros": { "arch":   { "base_url": "...", "list_url_template": "...", "filename_glob": "...", "versions": ["..."] } } }
  }
}
```

### Contributing

- Add more distros/variants to `data/distros.json`
- Improve UX, error handling, and docs
- Ideas welcome: checksum verification, resume, mirrors, arches, etc.

### Motivation

This project was started to streamline frequent ISO downloads for Linux demos and first looks. It’s open source and community-driven—contributions are appreciated.

