# AGENTS.md

This file defines non-negotiable rules for all Coding Agents.
Violations are considered incorrect behavior

---

## 1. Ownership Rules

- One task = one owning Agent.
- Either the Agent completes the task fully, or the human does.
- Do NOT mix partial human edits with Agent-written code.

---

## 2. Planning Is Mandatory

Before writing any production code, you MUST:

- Produce a concrete plan
- Decompose tasks into TODOs

Persist them in `todo.md`.

No plan → no code.

---

## 3. Result-First Constraint

- Humans specify desired outcomes, not implementations.
- You MUST propose solutions and designs.
- Do NOT wait for step-by-step instructions unless blocked.

---

## 4. Research Phase Rules

When investigating:

- Assume unlimited time and budget.
- Maximize reasoning depth.
- Over-research is preferred to under-research.

Human input is optional unless assumptions break.

---

## 5. Implementation Rules

- Implement end-to-end according to the plan.
- Avoid premature refactoring.
- Unit tests are required but never sufficient alone.

---

## 6. Testing Defines Completion

- “There’s a test, there’s a feature.”
- A feature is incomplete without validation.

Priority order:

1. Integration tests
2. End-to-end tests
3. Unit tests

Build a one-command test runner if missing.

---

## 7. Test Context Isolation

- Test generation must NOT share context with implementation.
- Treat tests as an external adversary.

Document test usage in `agents.md`.

---

## 8. Review Is Non-Optional

- Code MUST be reviewed by a different Agent.
- Reviewer receives:
  - Code
  - Design docs
- Reviewer must NOT receive:
  - Prior conversation context
  - Author rationale

Iterate until reviewer and author converge.

---

## 9. Knowledge Persistence

After merge, you MUST:

- Add design docs to `.modix/knowledge/`.
- Update `agents.md` with lessons learned.

Unwritten knowledge is lost knowledge.

---

## 10. Refactoring Threshold

- ~50k LOC per module is a hard warning sign.
- Agents do NOT autonomously redefine architecture.

Humans define:

- Module boundaries
- New structure

Agents execute refactors.

---

## 11. Parallelism Rules

- Parallel Agents must work on:
  - Different modules, or
  - Different branches, or
  - Different git worktrees

Avoid shared mutable context.

---

## 12. Core Principle

Agents optimize execution.
Humans optimize direction and boundaries.

Follow these rules strictly.

---
