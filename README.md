## iso-downloader

CLI/TUI tool to download Linux ISOs by distro and version.

### Install (one-liner)

```bash
curl -fsSL https://raw.githubusercontent.com/dhruvmistry2000/iso-downloader/refs/heads/main/install.sh | bash
```

Set `REPO=<owner>/<repo>` to install from your fork:

```bash
curl -fsSL https://raw.githubusercontent.com/dhruvmistry2000/iso-downloader/refs/heads/main/install.sh | REPO=<owner>/<repo> bash
```

### Usage

Interactive TUI:

```bash
iso-downloader
```

Non-interactive:

```bash
iso-downloader --non-interactive --distro ubuntu --versions 24.04,22.04 --output ./downloads
iso-downloader --non-interactive --distro ubuntu --versions all
```

### Config

TOML config is bundled at `iso_downloader/data/distros.toml`. You can override with `--config` or `ISO_DOWNLOADER_CONFIG`.

Schema:

```toml
[distros.<key>]
base_url = "..."                  # optional
url_template = "...{version}..."  # required
filename_template = "...{version}.iso"  # optional
versions = ["<version>", ...]
```

### Develop

```bash
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
pip install -e .
iso-downloader --help
```

### Build single binary

```bash
bash scripts/build.pyinstaller.sh
ls -l dist/iso-downloader
```

### Release via GitHub Actions

Tag a release:

```bash
git tag v0.1.0 && git push origin v0.1.0
```


