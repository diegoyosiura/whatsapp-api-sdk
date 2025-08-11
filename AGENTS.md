# AGENTS.md

## 1. Context and Purpose

* **Current state**: No intelligent agents implemented yet. The project aims to build a fully functional SDK for WhatsApp Business.
* **Agent objectives**: Support development by ensuring updated documentation, identifying bugs, optimizing code, creating automated tests, and verifying compliance with **SOLID** principles and **Hexagonal Architecture**.
* **Supervision**: All agents operate under human supervision. No changes are applied without explicit approval from the project owner.

## 2. Architecture and Technology Stack

* **Architecture**: Hexagonal.
* **Primary language**: Go (pure, no additional frameworks).
* **LLM**: OpenAI (Codex).

## 3. Types of Agents

> Reference: Common types of agents in software projects:
>
> * **Router**: Directs requests to specialized agents.
> * **Planner**: Defines an execution plan to achieve a goal.
> * **Executor**: Performs specific tasks such as API calls, code analysis, and documentation generation.
> * **Critic**: Reviews and validates generated results.
> * **Evaluator**: Measures quality, efficiency, and compliance.

For this project, the main roles will be **Critic** and **Evaluator**, with occasional **Executor** tasks for suggesting fixes or generating documentation.

* **Access**: GitHub repository of the project.
* **Sensitive data**: No handling of PII or PCI.

## 4. Observability and Metrics

* **Monitored metrics**: latency, cost, accuracy, and error rate.
* **Observability stack**: Not defined yet.

## 5. Development Process

* **Versioned prompts**: Yes, stored in version control.
* **AGENT.yaml per agent**: No. All interactions will be through ChatGPT.
* **Automated tests**: Yes, focused on unit tests.

## 6. Deployment and Governance

* **Environments**: DEV only for now.
* **Canary release and feature flags**: Not applicable.
* **Responsible for review/approval**: Project owner.

---

## Agent Rules and Guidelines

1. **Read-only code analysis**: Agents must not apply direct changes to the repository.
2. **Assisted operation**: All suggestions must be manually approved before integration.
3. **Mandatory documentation**: Any output generated (fix, bug report, optimization) must include a detailed explanation.
4. **Architectural compliance**: Continuous verification of SOLID principles and Hexagonal Architecture patterns.
5. **Activity logging**: Logs must capture the defined metrics and analysis context.

---

## Future Implementations

* Adopt an observability stack (e.g., Prometheus + Grafana).
* Expand scope to include automated tests generated and executed by agents.
* Implement agent audit policies to track suggestion and approval history.
