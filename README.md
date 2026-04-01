# Turms IM (Go Edition)

> A high-performance Golang rewrite/refactor of the Turms Instant Messaging engine.

## 📌 Repository Structure

This repository is the pure Go implementation and its architectural context. The original Java source code is mounted as a Git submodule (`turms-orig`) to provide a complete, side-by-side context for the refactoring process.

```text
.
├── turms-orig/              # Original Java Source Code (Git Submodule pointing to turms-im/turms)
├── docs/                    # Architecture analysis, BDD/TDD strategies, and Go mapping rules 
├── AGENTS.md                # 🧠 Main Context Entrypoint (Crucial for AI Agents)
├── cmd/                     # Go application entrypoints
├── internal/                # Private application and library code
└── pkg/                     # Public library code
```

---

## 🤖 For AI Agents & Developers: Initialization Guide

Because the original Java code (`turms-orig`) is mounted as a Git Submodule to prevent duplicate storage and ensure we can track upstream updates, **you must initialize the submodules correctly when cloning this repository**.

### 1. Initial Clone
When cloning this repository for the first time, use the `--recursive` flag to fetch the submodule data automatically:

```bash
git clone --recursive https://github.com/zhangxaochen/turms-im-go.git
```

**If you already cloned without the flag**, the `turms-orig` directory will be empty. Run the following command at the repository root to fetch the original code:

```bash
git submodule update --init --recursive
```

### 2. Following Upstream (Updating `turms-orig`)
To compare against the latest original Java source code, you can update the `turms-orig` submodule to track the newest upstream commit from the official repository:

```bash
# Pull the latest commit from the official turms repository
git submodule update --remote turms-orig

# Commit the updated submodule reference to this monorepo
git add turms-orig
git commit -m "chore: sync original turms to latest upstream"
```

---

## 🚀 Development Workflow

1. **Read `AGENTS.md` First**: Any AI Agent dropping into this repository must read `AGENTS.md` and the documents under `docs/` to understand the domain models, API protocol constraints, and translation rules.
2. **Generate Protobufs**:
   ```bash
   make generate
   ```

## 📜 Original Project
- Upstream Project: [turms-im/turms](https://github.com/turms-im/turms) 
- License: Apache 2.0
