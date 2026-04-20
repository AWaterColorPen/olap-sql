# Architecture

> **Work in progress.** This file is a placeholder — a full architecture deep-dive will be added in a follow-up PR.

## High-level overview

olap-sql is organized around four core abstractions:

| Component | Responsibility |
|-----------|----------------|
| **Schema / Configuration** | Declares the virtual schema: sets, sources, dimensions, metrics, and joins |
| **Manager** | Loads configuration and owns the compiled schema at runtime |
| **Query** | Describes what data the caller wants (filters, grouping, ordering, pagination) |
| **Result** | Delivers the query output as structured rows |

See [Getting Started](./getting-started.md) for a walkthrough of how these fit together.

## Coming soon

- Component interaction diagrams
- Dependency graph design
- Extension points and custom adapters
