# OpenCode Agent FlowKit

A CLI tool to install a pre-configured OpenCode agent workflow in any repository.

## What It Installs

This tool sets up a structured agent workflow for OpenCode with the following agents:

| Agent | Mode | Purpose |
|-------|------|---------|
| **plan** | primary | Architecture, constraints, sequencing. Writes to `.opencode/plan.md` |
| **execution** | primary | Orchestrates subagents to implement plans |
| **build** | primary | Default development agent with full access |
| **coder** | subagent | Writes code following instructions from execution |
| **review** | subagent | Reviews implementation against plan |
| **test-runner** | subagent | Executes tests and reports results |
| **testing-strategy** | subagent | Defines what to test and what not to test |
| **ask** | primary | Read-only agent for answering questions |
| **merge** | primary | Detect and resolve git merge conflicts interactively |

## Installation

### From Source

```bash
go install github.com/GH-Jaider/opencode-agent-flowkit@latest
```

### Build Locally

```bash
git clone https://github.com/GH-Jaider/opencode-agent-flowkit.git
cd opencode-agent-flowkit
make install
```

## Usage

```bash
# Install in current directory
opencode-agent-flowkit

# Install in a specific directory
opencode-agent-flowkit ./my-project

# Force overwrite existing files
opencode-agent-flowkit --force

# Show version
opencode-agent-flowkit --version

# Show help
opencode-agent-flowkit --help
```

## Files Created

```
your-repo/
├── opencode.json              # Main OpenCode configuration
└── .opencode/
    ├── config.json            # Model configuration (editable)
    ├── plan.md                # Plan template
    └── agents/
        ├── execution.md       # Orchestrator agent
        ├── coder.md           # Code writing agent
        ├── review.md          # Review/critic agent
        ├── test-runner.md     # Test execution agent
        ├── testing-strategy.md # Testing philosophy agent
        ├── ask.md             # Read-only Q&A agent
        └── merge.md           # Merge conflict resolution agent
```

## Configuration

### Model Configuration

Edit `.opencode/config.json` to configure which LLM model each agent uses:

```json
{
  "models": {
    "default": null,
    "execution": null,
    "coder": null,
    "review": "anthropic/claude-sonnet-4-20250514",
    "test-runner": null,
    "testing-strategy": null,
    "ask": null
  }
}
```

- `null` = uses OpenCode's default model
- Specify a model string to override for that agent

### Code Style Guidelines

This tool does NOT install an `AGENTS.md` file. Code style guidelines are language-specific and should be:

1. Created manually per project
2. Installed via language-specific skills/superpowers
3. Added to OpenCode's global configuration

## Workflow

The intended workflow is:

1. **Plan** (`plan` mode): Define what you're building in `.opencode/plan.md`
2. **Execute** (`execution` mode): Orchestrator reads plan and delegates to subagents
3. **Review** (`@review`): Compare implementation against plan
4. **Iterate**: Fix issues identified in review

```
Plan → Execute → Review → Fix → Done
         ↑                   ↓
         └───────────────────┘
```

### Merge Conflict Resolution

When you encounter merge conflicts, switch to the **merge** agent (Tab):

1. **Detect** - Agent runs `git status` and identifies all conflicted files
2. **Present** - For each conflict, shows both sides with context
3. **Choose** - You pick a resolution: keep ours, keep theirs, keep both, or a smart merge
4. **Apply** - Agent edits the file (with your approval) and offers to stage it
5. **Commit** - You review and commit when ready

```
git merge <branch> → conflicts? → Switch to merge agent (Tab)
                                        │
                              For each conflict:
                              ├── View ours vs theirs
                              ├── Choose resolution
                              └── Approve edit
                                        │
                              git add → git commit
```

## Updating

To update to the latest version, run the installer again with `--force`:

```bash
opencode-agent-flowkit --force
```

This will overwrite all configuration files with the latest templates.

## License

MIT
