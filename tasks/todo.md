# Port dev-cache filters to the Bilge codebase

## Scope
- Update PR #73 with the current `main` branch.
- Port include/exclude path filters and min/max age filters to current `dev-cache`.
- Update the governing `PRD.md` requirements before implementation.
- Resolve verified review findings and failing CI, then validate and request review.

## Execution checklist
- [x] Merge latest `origin/main` into the PR branch and resolve conflicts without losing current safety behavior.
- [x] Update `PRD.md` with filter precedence, age handling, and cleanup behavior.
- [x] Add failing tests for review findings, confirm they fail, then implement the fixes.
- [x] Run formatting, targeted tests, full repository tests, vet, and coverage checks.
- [ ] Run Bosun against the exact current PR base and address valid findings.
- [ ] Push the branch and verify all GitHub Actions and review status.

## Test strategy
- Unit tests for include/exclude matching, duration parsing, age boundaries, and unknown modification times.
- Integration tests for filtered scan and cleanup paths.
- Full repository tests and coverage checks required by repository guidance.

## Rollback strategy
- Revert the filter commit if the behavior or CI is incorrect; retain the main-merge commit needed to keep the PR current.

## Outcome
- Result: Current-main merge integrated; filter precedence and age-filter traversal/unknown-time findings fixed.
- Evidence: focused tests failed before the fixes and passed afterward; `make test`, `make vet`, `make fmt`, and `go test -cover ./...` passed for dev-cache (61.0%). Local golangci-lint could not load packages with installed Go 1.26; GitHub CI remains the authoritative pinned-toolchain check.
