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

## <DONE> cross-cutting/infra dir?
Currently `internal/` contains mix of business, infra (postgres, server, auth, version), config and helper packages.
I am considering to separate business slices by moving non-business packages into separate directory.
What is your take on this?
HOW: `internal/platform/` dir where non-buss packages were moved

## validation library
Is there a place for validation library in the project, in any package, in any layer?

## CI/CD pipeline
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

## <SEE CHATGPT CONVERSATION> Multiple Clients
- Separate or single repository for each client?
- Should they have any common core?
- How organize development and deploys?
- For cloud-based application - single domain with subdomains (e.g. companyA.meslite.app) or separate domain (e.g. companyA.app)?

## remove protections for existing data (in reality no data exist yet)
e.g. remove `request_id <> ''` from `production_entries.request_id` index
remove `NOT VALID` from migrations/0004_add_production_reference_foreign_keys.sql
check other migrations

## config from environment
Analyze what settings from code should be read from env on start

## move all "Lesson Scope" and "Lesson Completion Notes" out of ROADMAP
to lessons dir

## Feature Flag

I prepared proposal of application configuration in file `docs/config-synthesis.md`.
First analyze and asses it.
Provide **concise summary of your assesment**.

Agree to your recommendations. Now review proposal again in context of your assesment.
Provide a **short summary of your findings** and identify **all questions I need to answer** to close the remaining:
* information gaps,
* ambiguities,
* inconsistencies,
* contradictions,
* unclear decisions,
* and assumptions that would otherwise need to be made.

opencode -s ses_fa7ce8e2dffex2B3VRTrLQXyg2




Questions:
1. What `cancel()` does in following context (WorkerPool.execute method):
```
	jobCtx, cancel := context.WithCancel(ctx)
	p.mu.Lock()
	p.running[job.ID] = cancel
	p.mu.Unlock()
	defer func() {
		cancel()
		p.mu.Lock()
		delete(p.running, job.ID)
		p.mu.Unlock()
	}()
```
2. How job handler knows that it is cancelled and how it handles cleanup? Give example.

Please add comments explaining above to code where appropriate - I want code reader to understand cancellation flow.
Also add markdown file with examples how to use worker pool and queue. Place this file in `docs/04-development/` directory.

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




## System design documents
I added two new files to this ChatGPT Project's sources:
- `MESLite-marketing-concept.md`
- `MESLite-product-description.md`

Read and analyze these files together with the sources and conversations you considered previously.
Then **re-evaluate your previous assessment** of whether there is enough information to create `product-context.md`.
Provide a **short summary of your updated findings**. Do not create or modify `product-context.md` yet.

Most importantly, identify **all questions I need to answer** to close the remaining:
* information gaps,
* ambiguities,
* inconsistencies,
* contradictions,
* unclear product decisions,
* and assumptions that would otherwise need to be made.

Prioritize the questions that are most important for creating an accurate `product-context.md`.
For each question, briefly explain **why the answer matters** or what part of the product context it affects.

Do not ask questions that can already be answered confidently from the available sources. Where the sources provide conflicting information, explicitly point out the conflict and ask me to resolve it.




I will use an AI Agent to develop the **MESLite** project.

In this ChatGPT Project's sources, there is a file named `MESLite-documentation-guide.md`. It defines the set of documents intended to support AI-driven development of MESLite.

I want to eventually create a `product-context.md` document.

## Your task

Before creating anything, analyze whether there is **enough reliable and sufficiently detailed information** available to create a high-quality product-context.md.

Use and cross-reference the following sources:

1. **This ChatGPT Project's conversation history**, where relevant.
2. **Primarily this specific chat**, which contains the main discussion about the MESLite product:
`https://chatgpt.com/g/g-p-6a8336554548819186fb62fa5638ce06-mes-lite/c/6a833b59-d6f8-83eb-b4a0-19fd9e47e219`
3. **All relevant Markdown files uploaded to this ChatGPT Project's sources**.

## What I want you to assess

Determine:
* What information about MESLite is already available.
* Whether that information is sufficient to produce `product-context.md` according to the documentation guidance in `MESLite-documentation-guide.md`.
* Which important pieces of product context are well established and which are ambiguous, incomplete, or missing.
* Whether there are contradictions between the available sources that should be resolved before creating the document.
* What additional information, if any, I should provide before `product-context.md` can be created confidently.
* Whether some information should be inferred from the existing material, and clearly distinguish such inferences from explicitly stated facts.

## Important constraints

* **Do NOT create or write** `product-context.md` yet.
* Do **not** fill gaps by inventing assumptions.
* Do **not** treat uncertain or inferred information as established facts.
* Focus on assessing the available information and identifying gaps.
* If the available information is sufficient, explain **why** it is sufficient and what the resulting document would be able to cover.
* If it is insufficient, provide a **concrete list of missing or unclear information**, preferably organized by priority.

## Expected output

Give me an assessment structured roughly as:
1. **Overall assessment** — whether we have enough information to create `product-context.md`.
2. **What is already known** — key product context supported by the sources.
3. **Confidence / evidence** — distinguish explicit facts, strong inferences, and uncertainties.
4. **Gaps and unresolved questions** — information that is missing or needs clarification.
5. **Contradictions or inconsistencies** — if any exist between sources.
6. **Recommendation** — whether to proceed with creating `product-context.md` now or gather more information first.

Again: **do not create `product-context.md` in this step**.


I will use AI Agent to develop MESLite project.
Create markdown document that will contain all above arrangements.
Keep documents precise and as concise and condensed as possible.
Documents must provide sufficient information about this project for AI to develop and maintain this project.
Based on these documents, the AI should be able to:
- understand general project idea
- understand business context
- understand business requirements
- understand system architecture
- understand code architecture and how to stick to it
- anything else to develop and maintain project successfuly and according to best practices
