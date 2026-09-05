## Conversion to next phase
> current: Learning project -> Proof of Concept / Demo (production ready)
0. Learning project
1. Proof of Concept / Demo (PoC)
2. Prototype (Proto)
3. Most Valuable Product (MVP)

At every phase (point) whole Process should be made from begining to end.

### ONE-TIME Process prerequisites
- refine `docs/*`
    - <current> docs/01-product/product-context.md
- move all "Lesson Scope" and "Lesson Completion Notes" out of ROADMAP to lessons dir
    - replace with references to lesson file

### Process prerequisites
- define what is project's scope for current phase
- create/update `requirements.md` and `domain.md` docs.
- check todos `docs/ToDo/*`
- Review whole project, steps:
    - code review?
    - architectural review?
    - documentation review?

### Process
1. documentation/knowledge for AI (AI-driven development)
What is needed (documentation, knowledge, skills, commands, prompts, AGENTS.md)
project business requirements? (+summary/compaction?)
project technical description? (+summary/compaction?)
project architecture? (+summary/compaction?)
project design? (+summary/compaction?)
should all those documents be live and updated periodically?
Is this knowledge-base?
input:
    - ROADMAPs
    - code in repo
    - summary from GPT (for me an Piter - marketing, concept)
    - docs/system-design/architecture-and-mvp.md or/and docs/system-design-ai-conversation.md
    - docs/adr/*
    - docs/sqlc-boundaries.md,validation.md
2. review whole project for problems, bugs, inconsistencies, gaps, improvements, refactors.
    - docs
    - code
    - devops and tooling
3. developer's guide: whole path/route from README.md -> docs (about GO, about Project) -> proj arch/struct + packages relations -> code
