# Pathfinder: Data Platform Sandbox

**Project Description**  
Pathfinder is a sandbox project for experimenting with diverse data technologies, including brokers, databases, data models, and architectural patterns. The goal is to create a flexible data platform that allows for rapid prototyping and evaluation of new ideas in the data space.

## Purpose

This repository serves as a playground for:
- Testing integration patterns between data brokers and databases
- Modeling data using various paradigms (relational, NoSQL, event-driven, etc.)
- Evaluating scalability, reliability, and performance of different data technologies
- Learning and documenting best practices and lessons learned

## Technologies Used

The codebase is composed of several languages and frameworks:

- **Smarty (58.2%)** — Used primarily for templating and dynamic content generation
- **Mustache (36.8%)** — Provides logic-less templates for clean data rendering
- **Shell (3.2%)** — Scripts for automation, setup, and orchestration
- **Go (1.7%)** — Backend services, data connectors, or tooling
- **Makefile (0.1%)** — Build automation and workflow management

## Getting Started

> This repo is experimental and subject to frequent changes.  
> To get started, clone the repository and explore the `/docs` or `/examples` directories for sample setups.

```sh
git clone https://github.com/MGYOSBEL/pathfinder.git
cd pathfinder
```

### Setup

Depending on the experiment, you may need:
- Docker and Docker Compose
- Go (>=1.18)
- Bash

Refer to specific directories or documentation for instructions on running individual components.

## Directory Structure

```
.
├── templates/      # Smarty/Mustache templates
├── scripts/        # Shell scripts for automation
├── go/             # Go services and tools
├── Makefile        # Automation rules
├── README.md
└── docs/           # Documentation and notes
```

## Contributing

This is a personal sandbox, but feedback and suggestions are welcome.  
If you have ideas for data technologies to explore, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.

---

*Data Platform project sandbox by [MGYOSBEL](https://github.com/MGYOSBEL)*
