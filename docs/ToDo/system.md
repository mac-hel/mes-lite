## update to latest versions:
- go, go deps, postgresql
- anything else

## local development
- Containerize local app (currently only DB is in container)
- corelate with production builds

## CI/CD pipeline
inspect and make it useful

## GitHub Actions, pre-commit hooks

## golangci-lint
- improve configuration of used linters
- improve usage of this tool (what other useful capabilities it has)
- what linters are worth to add?
    - their configuration

## air
- improve config

## <ADDED> Sanity tests
Should we introduce sanity tests? Either implement new ones or convert some existing tests.

## <DONE> cross-cutting/infra dir?
Currently `internal/` contains mix of business, infra (postgres, server, auth, version), config and helper packages.
I am considering to separate business slices by moving non-business packages into separate directory.
What is your take on this?
HOW: `internal/platform/` dir where non-buss packages were moved

## config from environment
Analyze what settings from code should be read from env on start

## Config/Feature Flag - by layer and scope

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

> opencode -s ses_fa7ce8e2dffex2B3VRTrLQXyg2

I would drop persisted runtime config.
The same effect should be achieveable by changing config file on the fly and triggering re-read at application runtime.

**SCRATCH**
```
// Business settings
Code: default value / required              - for all app instances (all Companies/Factories)
Company Settings file: Company overrides    - each app instance (Company) has separate file
    - re-read can be triggered run-time, without app restart
Factory Settings file: Factory overrides    - each app instance can have multiple files (for multiple Factories)
    - re-read can be triggered run-time, without app restart

// Environment (Technical?) settings
Code: default value / required              - for all app instances (all Companies/Factories)
Env variables: Company overrides            - each app instance (Company) has separate env
    - re-read only on restart
```
In the code, settings are behind facade/abstraction - this will allow to implement DB backed or UI-driven settings in the future.

## Unified observability
- info/error logs
- panics/startup failures
- metrics - technical and business
- tracing

* info/error types
* how each type should be handled and where logged?
* E.G.: worker's attempt to register production fails - should it go to metrics or logs? Should it raise immediate alarm (not only be reported to worker)????????????????

By default all goes to stdout, but where should be exported for each Company/Factory?
1 Collector per Factory?

## Multiple Factories
What if company runs multiple factories?
Should production be managed separately, but in single system?

## <SEE CHATGPT CONVERSATION> Multiple Clients
- Separate or single repository for each client?
- Should they have any common core?
- How organize development and deploys?
- For cloud-based application - single domain with subdomains (e.g. companyA.meslite.app) or separate domain (e.g. companyA.app)?

## Performance requirements - define for AI
- expected users/requests num
- expected load (csv imports etc.)
- what else?
