# VCD Claude Speckit

This repository is vibe coding experimental with Github Speckit, Serena MCP and Claude AI (With GLM 4.5 model)

## Overview

- Use Speckit to follow Spec Driven Development (SDD) principles.
- Use Serena MCP for semantic retrieval and editing capabilities (MCP server & other integrations)
- Use Claude AI to generate code based on requirements and specifications.

## Setup

### 1. Initialize a new Speckit project

-- Follow the instructions in the [Speckit documentation](https://github.com/github/spec-kit) to set up your project.

   ```bash
   uvx --from git+https://github.com/github/spec-kit.git specify init <PROJECT_NAME>
   ```

### 2. Setup Claude Code with GLM 4.5

- Register and get API access from [Anthropic](https://www.anthropic.com/claude-code).
- Configure your environment to use the GLM 4.5 model follow the instructions in the [Claude Code documentation](https://docs.z.ai/devpack/tool/claude) for setup.

### 3. Install Serena MCP

- Follow the instructions in the [Serena MCP repository](https://github.com/oraios/serena) for installation and setup.

## References

- [Github Speckit](https://github.com/github/spec-kit)
- [Serena MCP](https://github.com/oraios/serena)
- [Claude Code](https://www.anthropic.com/claude-code)
- [GLM 4.5](https://bigmodel.cn/)
