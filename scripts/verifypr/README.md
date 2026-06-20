# verifypr

Simple script that is called by the [verifypr](../../.github/workflows/ci-verifypr.yml) github action
that verifies story PRs.

It checks two things:

1. The PR **title** follows the [conventional commit](https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-conventional#type-enum)
   style, e.g. `feat: ...` or `fix(scope): ...`. Casing, punctuation and scope are not restricted beyond
   what the conventional commit spec requires; the title is capped at 100 characters.
2. The PR **body footer** contains a single `issue:` line referencing a github issue, e.g. `issue: #1234`,
   `issue: piplabs/story-geth#1559` or `issue: none`. The footer is the last paragraph of the body (the
   block after the final blank line), so anything above it (e.g. bullet lists) is free-form.
