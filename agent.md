# Project: reconx | Type: Security
## Role: Expert Security Developer
## SDD Protocol: Spec-First
## Skills: ./skills/new_spec.sh "Feature Name"
# ReconX Framework Development Agent

## Mission

You are the lead engineer of ReconX.

ReconX is a professional reconnaissance framework written in Go.

The goal is to create a modern, modular, scalable and production-ready reconnaissance platform for:

* Bug Bounty
* Red Team Operations
* External Pentesting
* Attack Surface Management
* Continuous Asset Discovery

The framework is intended for public GitHub release and community adoption.

The final product must resemble professional security software rather than a collection of scripts.

---

# Core Principles

Prioritize:

* Maintainability
* Extensibility
* Reproducibility
* Evidence collection
* Automation
* Performance

Avoid:

* Giant Bash scripts
* Script-kiddie workflows
* Hardcoded values
* Tool-specific assumptions
* Monolithic architecture

---

# Technology Stack

Primary language:

Go

Allowed secondary languages:

* Bash (installation only)
* Python (optional analysis plugins)

Core framework must remain in Go.

---

# Architecture

Repository structure:

reconx/
  cmd/
  internal/
  pkg/
  configs/
  tests/
  docs/
  examples/
  .github/

## Industrialization Standard

The framework must move from a sequential script to an orchestrated engine:

1. **Task Interface**: Every module must implement a `Task` interface (Name, Execute).
2. **Session Model**: Modules communicate via a shared `Session` object in memory, reducing disk I/O. Disk persistence occurs only at the end of each task.
3. **Pipeline Engine**: A central orchestrator in `internal/engine` that manages the execution sequence, state, and resumability.
4. **Worker Pool**: Controlled concurrency for tool execution to prevent system exhaustion.

### Execution Modes
- **Full Pipeline**: Executes all registered tasks in sequence.
- **Targeted Execution**: Ability to run a specific task via CLI (`-task <task_name>`) for development and debugging.
- **Resume Mode**: Resume an interrupted scan from the last completed task (`-resume`).
- **Monitoring Mode**: Run the full pipeline continuously on a schedule (`-monitor`, `-interval`).

### API Key Management
- **Secure Storage**: OSINT API keys are stored in `~/.config/uncover/provider-config.yaml`, outside the project directory.
- **Setup CLI**: Keys are configured interactively via `./reconx --setup-api`.

### Dependency Management
- **Automated Verification**: The framework must verify the presence of all required tools in the system PATH before execution.
- **Self-Healing/Auto-Install**: If a required tool is missing and a `go install` or `uv sync` command is known, the framework should attempt to install it automatically.

Each module must be independent and decoupled from the orchestration logic.

---

# Development Roadmap

## Phase 1

Project Bootstrap

Create:

* CLI
* Logger
* Configuration Loader
* Tool Detection Engine
* Output Engine

Deliverables:

* main.go
* config package
* logging package
* version package

---

## Phase 2

Asset Discovery

Implement:

* Subfinder
* Chaos
* Github-Subdomains
* Amass (optional)

Output:

subdomains.json

subdomains.txt

Goals:

* deduplication
* validation
* concurrency

---

## Phase 3

DNS Resolution

Implement:

* dnsx

Collect:

* A
* AAAA
* CNAME
* MX
* TXT

Generate:

dns.json

---

## Phase 4

Port Discovery

Implement:

* naabu

Optional:

* nmap integration

Generate:

ports.json

---

## Phase 5

HTTP Discovery

Implement:

* httpx

Collect:

* title
* technologies
* redirects
* status
* CDN
* ASN
* IP
* favicon hashes

Generate:

http.json

alive.txt

---

## Phase 6

Content Discovery

Implement:

* katana

Support:

* headless crawling
* JS parsing
* forms
* API discovery

Generate:

crawl.json

urls.txt

---

## Phase 7

Historical Collection

Implement:

* gau
* waymore

Avoid:

* waybackurls

Generate:

historical.json

---

## Phase 8

JavaScript Intelligence

Implement:

* xnLinkFinder

Extract:

* endpoints
* secrets
* tokens
* API references
* cloud references

Generate:

javascript.json

---

## Phase 9

Secret Discovery

Implement:

* trufflehog

Detect:

* AWS
* Azure
* GCP
* GitHub
* GitLab
* JWT
* API Keys

Generate:

secrets.json

---

## Phase 10

Cloud Enumeration

Implement:

* cloudlist
* s3scanner
* bucket discovery modules

Discover:

* S3
* Azure Storage
* GCP Storage

Generate:

cloud.json

---

## Phase 11

OSINT Expansion

Implement:

* uncover

Support:

* Shodan
* FOFA
* Censys
* Hunter

Generate:

osint.json

---

## Phase 12

Visual Reconnaissance

Implement:

* gowitness

Generate:

screenshots/

index.html

---

## Phase 13

Vulnerability Discovery

Implement:

* nuclei

Requirements:

* template version tracking
* severity filtering
* tagging

Generate:

nuclei.json

---

## Phase 14

Reporting Engine

Generate:

* report.md
* report.html
* report.json

Include:

* statistics
* findings
* screenshots
* evidence

---

## Phase 15

Workspace Management

Support:

* resumable scans
* caching
* incremental scans
* historical comparisons

Generate:

scan history

---

## Phase 16

Attack Surface Diffing

Compare:

Current Scan

vs

Previous Scan

Identify:

* new hosts
* removed hosts
* new ports
* new technologies

Generate:

diff reports

---

## Phase 17

Continuous Monitoring

Implement:

scheduled scans

alerting

new asset detection

---

# Tool Priority List (Current Implementation)

**Go binaries (auto-install via `go install`):**

- subfinder
- dnsx
- naabu
- httpx
- katana
- gau
- uncover
- **gowitness** (github.com/sensepost/gowitness)
- **nuclei** (github.com/projectdiscovery/nuclei/v3)

**Python tools (managed via `uv sync`):**

- **cloud-enum** (github.com/initstring/cloud_enum)

**Manual installation required:**

- trufflehog (binary)
- waymore (go install)
- assetfinder (go install)
- xnLinkFinder (python)

**Optional:**

- Amass
- Nmap

**Avoid:**

- Sublist3r
- Hakrawler
- Waybackurls
- cloudlist (deprecated, replaced by cloud-enum)

---

# Performance Requirements

Must support:

* 10 domains
* 1,000 domains
* enterprise scopes

Implement:

* worker pools
* context cancellation
* rate limiting
* retries
* timeouts

---

# Testing Requirements

Every package requires:

* unit tests
* integration tests
* error handling tests

Minimum:

80% coverage

---

# Logging Requirements

Every module must log:

* start time
* end time
* duration
* errors
* warnings

Support:

JSON logging

Structured logging

---

# Security Requirements

Never store:

* credentials
* secrets
* API keys

Use:

environment variables

Validate all user input.

Prevent command injection.

Sanitize filenames.

---

# GitHub Standards

Repository must include:

README.md

LICENSE

CHANGELOG.md

CONTRIBUTING.md

SECURITY.md

GitHub Actions

Release Workflow

Issue Templates

Pull Request Templates

---

# Documentation Requirements

Generate:

Architecture Documentation

Installation Guide

Developer Guide

Module Documentation

Examples

Migration Guides

---

# Long-Term Vision

ReconX should evolve into a complete reconnaissance platform capable of:

* Asset Discovery
* Attack Surface Mapping
* Cloud Enumeration
* Continuous Monitoring
* Vulnerability Discovery
* Reporting

while remaining modular and easy to maintain.

Every implementation decision should move the project toward that vision.

---

# Project Status (June 2026)

All 17 phases are implemented and verified.

| Phase | Module | Status |
|-------|--------|--------|
| 1. Bootstrap | `cmd/reconx`, `internal/config`, `internal/logging`, `pkg/version` | ✅ |
| 2. Asset Discovery | `internal/discovery` | ✅ |
| 3. DNS Resolution | `internal/dnsres` | ✅ |
| 4. Port Discovery | `internal/ports` | ✅ |
| 5. HTTP Discovery | `internal/httpdisc` | ✅ |
| 6. Content Discovery | `internal/content` | ✅ |
| 7. Historical Collection | `internal/historical` | ✅ |
| 8. JS Intelligence | `internal/jsintel` | ✅ |
| 9. Secret Discovery | `internal/secrets` | ✅ |
| 10. Cloud Enumeration | `internal/cloud` (wraps cloud-enum) | ✅ |
| 11. OSINT Expansion | `internal/osint` (wraps uncover) | ✅ |
| 12. Visual Reconnaissance | `internal/visual` (wraps gowitness) | ✅ |
| 13. Vulnerability Discovery | `internal/vuln` (wraps nuclei) | ✅ |
| 14. Reporting Engine | `internal/reporting` | ✅ |
| 15. Workspace Management | `internal/engine` (resume, history) | ✅ |
| 16. Attack Surface Diffing | `internal/diffing` | ✅ |
| 17. Continuous Monitoring | `cmd/reconx` (-monitor flag) | ✅ |

**Industrial Features:**
- Task-based pipeline engine (`internal/engine`)
- Modular Go packages with `Task` interface
- Auto-dependency management (`internal/deps`)
- Secure API key storage (`--setup-api`)
- Resumable scans (`-resume`)
- Targeted task execution (`-task`)
- Continuous monitoring (`-monitor -interval N`)
- Reporting (Markdown, HTML, JSON)
- Attack surface diffing
