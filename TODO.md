## <IGNORE FOR NOW> Employee state/status
Employee.isActive -> multiple statuses (sick, vacation, out-of-work etc.)

## <POSTPONED> User-configurable product categories
Distinct entity stored in DB

## <PARTIALY> Encapsulation - hide private properties of entities (structs)
Currently all domain Entities have public fields. Would it be feasible to make fields private at current phase of project?
- <DONE> Order

## <PARTIALY> Fields as Collections
What do you think about using collections as types of fields that contain multiple values, e.g. fields like:
- <DONE> `Order.lines` (currently `[]OrderLine`, could be collection)
- <REJECTED FOR NOW> `Order.assignedEmployees` (currently `[]string`, could be collection)

## <ADDED> (as separate reports) Daily Production Report lacks num of products made by employee
Currently Daily Production Report provides:
- total number of manufactured units for product
- number of production entries for product
Shouldn't it also provide number of units manufactured by employee?
**Reports**
*Daily Production*      - daily, How much of each product was made?
*Employee Productivity* - during period, How much products did each employee made overall?
*Product Statistics*    - during period, How much of each product was made?
**Added:**
*Daily Employee Production* - daily, How much of each product did each employee made?
*Employee Productivity for Products* - during period, How much of each product did each employee made?

## <REJECTED> Each package builds all sqlc queries
Currently each package (slice) builds all sqlc queries across all packages.
Isn't this approach:
- confusing for developers
- crossing package boundaries
Shouldn't each package build only own `sqlc` queries and create only own query models?
see `docs/sqlc-boundaries.md`

## <REJECTED> Double id - UUID for user-facing, int for internal reference
What do you think about using UUID and int ID simultaneously? UUID for user-facing things, int for internal references (for speed).
ANSWER: Do not introduce int IDs just for speed right now. Keep one canonical identifier per entity. If performance becomes real later, benchmark and revisit with data.

## <ADDED> Sanity tests
Should we introduce sanity tests? Either implement new ones or convert some existing tests.

## cross-cutting/infra dir?
Is this good idea to move packages not related to business domain into one or two common directories?
Cross-cutting or infrastructure related packages that could be moved: auth, config, postgres, server, version.

## validation library
Is there a place for validation library in the project, in any package, in any layer?

## pipeline
inspect and make it useful

## golangci-lint
- improve configuration of used linters
- improve usage of this tool (what other useful capabilities it has)
- what linters are worth to add?
    - their configuration

## air
- improve config

## GitHub Actions, pre-commit hooks

## Multiple Factories
What if company runs multiple factories?
Should production be managed separately, but in single system?

## Multiple Clients
- Separate or single repository for each client?
- Should they have any common core?
- How organize development and deploys?
- For cloud-based application - single domain with subdomains (e.g. companyA.meslite.app) or separate domain (e.g. companyA.app)?

## remove protections for existing data (in reality no data exist yet)
e.g. remove `request_id <> ''` from `production_entries.request_id` index
remove `NOT VALID` from migrations/0004_add_production_reference_foreign_keys.sql
check other migrations




## Conversion Learn->Real Project
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
2. finish TODOs
3. review whole project for problems, bugs, inconsistencies, gaps, improvements, refactors.
    - docs
    - code
    - devops and tooling
4. developer's guide: whole path/route from README.md -> docs (about GO, about Project) -> proj arch/struct + packages relations -> code



I have few questions regarding this lesson's implementation:
1. Should `func (c Correction) Validate() error` be not a method but stand-alone function, that can be used without `Correction` instance?
2. Shouldn't `SaveCorrection/ListCorrections` methods be placed in own interface, not in Store?
3. Similar question related to `Registrar interface` - shouldn't `CorrectEntry/ListCorrections` methods be placed in own interface, not in Registrar?
4. `ListProductionEntriesResponse` contains `[]Entry` field which is domain entity. Is this acceptable that response is coupled to domain this way?
5. Shouldn't `POST /production-entries/{id}/corrections` be renamed? It is bit misleading now. But - is there other REST-like name for this route? E.g. `POST /production-entries/{id}/correct`.
