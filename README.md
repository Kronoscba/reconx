# 🛡️ ReconX

![Go Version](https://img.shields.io/badge/language-go-00ADDR?logo=go)
![Python Version](https://img.shields.io/badge/language-python-37695E?logo=python)
![CI](https://github.com/Kronoscba/reconx/actions/workflows/ci.yml/badge.svg)
![License](https://img.shields.io/badge/license-MIT-blue)
![Status](https://img.shields.io/badge/status-production-green)
![Security](https://img.shields.io/badge/type-Security%20Framework-red)

**ReconX** is a professional-grade reconnaissance framework designed for high-scale asset discovery, attack surface mapping, and continuous monitoring. Developed as a high-performance Go engine, it orchestrates a curated ecosystem of industry-standard security tools—primarily written in Go and supplemented by specialized Python modules—to provide a scalable and production-ready platform for Bug Bounty hunters, Red Teams, and Security Researchers.

---

## 🚀 Core Industrial Features

Unlike traditional recon scripts, ReconX is built as an orchestrated engine:

- **🧩 Task-Based Architecture**: Every module implements a strict `Task` interface, ensuring a decoupled and extensible pipeline.
- **⚡ Memory-First Session**: Modules communicate via a shared in-memory `Session` object, drastically reducing disk I/O and increasing speed.
- **🔄 Resumable Scans**: Use the `-resume` flag to pick up exactly where a previous scan left off.
- **👁️ Continuous Monitoring**: Scheduled execution via `-monitor -interval <time>` to detect new assets in real-time.
- **🔐 Secure API Management**: API keys are stored outside the project directory (`~/.config/reconx/`) and managed via a dedicated `--setup-api` CLI.
- **🛠️ Self-Healing Dependencies**: Automatically detects and attempts to install missing Go/Python tools from official sources.

---

## 🗺️ The Reconnaissance Pipeline

ReconX orchestrates 17 specialized phases to provide a complete view of the target's attack surface:

| Phase | Focus | Primary Tools | Output |
| :--- | :--- | :--- | :--- |
| **1-2** | **Asset Discovery** | `subfinder`, `chaos`, `github-subdomains` | `subdomains.json` |
| **3** | **DNS Resolution** | `dnsx` | `dns.json` |
| **4** | **Port Discovery** | `naabu`, `nmap` | `ports.json` |
| **5** | **HTTP Discovery** | `httpx` | `http.json`, `alive.txt` |
| **6** | **Content Discovery** | `katana` | `crawl.json`, `urls.txt` |
| **7-8** | **Intel & JS** | `gau`, `waymore`, `xnLinkFinder` | `historical.json`, `javascript.json` |
| **9-11** | **Secrets & Cloud** | `trufflehog`, `cloud-enum`, `uncover` | `secrets.json`, `cloud.json`, `osint.json` |
| **12** | **Visual Recon** | `gowitness` | `screenshots/`, `index.html` |
| **13** | **Vuln Discovery** | `nuclei` | `nuclei.json` |
| **14-17** | **Analysis & Ops** | Internal Engine | `report.md`, `diff.json`, `history.db` |

---

## 🛠️ Installation

### 📦 Quick Install (Recommended)
Download the latest pre-compiled binary for your OS from the [Releases](https://github.com/Kronoscba/reconx/releases) page.

### 💻 Build from Source
**Prerequisites**: Go (latest stable), Python 3.x (with `uv`).

```bash
# Clone the repository
git clone https://github.com/Kronoscba/reconx.git
cd reconx

# Build the binary
go build -o reconx cmd/reconx/main.go

# Configure API Keys (Shodan, Censys, etc.)
./reconx --setup-api
```

---

## 🤖 CI/CD & Automation

ReconX utilizes GitHub Actions to maintain industrial software standards:

- **Continuous Integration**: Every push and PR is automatically validated through linting, unit testing, and build checks to ensure stability.
- **Automated Multi-Platform Releases**: On every version tag (`v*`), the framework is cross-compiled for **Linux, macOS, and Windows (amd64 & arm64)** and uploaded automatically to GitHub Releases.

## 📖 Usage Examples

### Full Pipeline Execution
Run the entire 17-phase pipeline against a target:
```bash
./reconx -target example.com
```

### Targeted Task Execution
Execute only a specific module (e.g., just the DNS resolution phase):
```bash
./reconx -target example.com -task dnsres
```

### Resuming an Interrupted Scan
```bash
./reconx -resume
```

### Continuous Monitoring
Run the pipeline every 24 hours and alert on new findings:
```bash
./reconx -target example.com -monitor -interval 24h
```

---

## 🏗️ Architecture

ReconX follows a clean, modular Go structure:
- `cmd/`: Entry points and CLI logic.
- `internal/engine/`: The orchestrator managing task sequencing, state, and resumability.
- `internal/discovery`, `internal/dnsres`, etc.: Decoupled modules implementing the `Task` interface.
- `pkg/`: Shared utility libraries.

---

## 🤝 Contributing

ReconX is designed to be extended. To add a new module:
1. Implement the `Task` interface in a new package within `internal/`.
2. Register the task in the `engine` pipeline.
3. Add necessary tool dependencies to the `deps` manager.

---

## 📜 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
