# Repositories

Shared by every daml-escrow repo; the source is `daml-escrow`. Edit it there, then copy it over unchanged.

| Repository | Visibility | Owns |
|---|---|---|
| `daml-escrow` | private | The escrow platform: DAML contracts, the Go escrow-api, the Astro portal, the Helm charts. The source of the shared files below |
| `daml-escrow-commons` | public | Shared Go libraries (schema validation, HMAC signing, DTO validation, logging, Canton ids, API errors), imported by tagged version |
| `daml-escrow-identity` | private | The T1 identity directory shared across products: OIDC subject, email and a stable cross-product identity hash |
| `daml-escrow-cms` | private | Contract lifecycle management: contract authoring, templates and negotiation before an escrow exists |
| `daml-escrow-platform` | private | Tenancy, subscriptions and usage metering for the platform's own operating costs, and the institution directory |

Each repo is a sibling checkout (`../<name>`), merges only through PRs to `main`, and releases by `RELEASING.md`.

## Shared files

These are identical in every repo. The `Shared files` job in `.github/workflows/repo-policy.yml` fails a PR (and a
weekly run) in a sibling whose copy differs from `daml-escrow`'s `main`:

- `.claude/skills/github-ci/SKILL.md`: the Git and CI rules
- `.github/pull_request_template.md`
- `.github/workflows/repo-policy.yml`: signed commits and this check
- `docs/REPOSITORIES.md`: this file

To change one, change it in `daml-escrow` first, merge, then copy it unchanged to each sibling. The siblings read
`daml-escrow` with the `ESCROW_READ_TOKEN` org secret (a fine-grained token, Contents read-only on `daml-escrow`
only); without it the job warns and passes.
