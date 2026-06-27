# Release engineering

How `ws` is versioned and shipped.

## Versioning

[Semantic Versioning](https://semver.org). Pre-1.0, breaking changes bump the
minor. The version is injected at build time into `main.version` via
`-ldflags "-X main.version=<v>"` (GoReleaser does this from the git tag);
unversioned builds report `dev`.

## Cutting a release

1. Make sure `master` is green and everything is merged.
2. Move the `## [Unreleased]` items in [`../../CHANGELOG.md`](../../CHANGELOG.md)
   under a new `## [X.Y.Z]` heading with today's date; update the compare links.
3. Commit (`chore: release vX.Y.Z`) via PR and merge.
4. Tag the merge commit and push the tag:

   ```bash
   git tag vX.Y.Z
   git push origin vX.Y.Z
   ```

5. The **Release** workflow (`.github/workflows/release.yml`) runs on the tag:
   it tests, then GoReleaser builds darwin/linux × amd64/arm64 archives +
   checksums and publishes a GitHub Release with auto-generated notes.

## Configuration

- **`.goreleaser.yaml`** — build matrix, archive naming, changelog grouping
  (Features / Fixes), `main.version` injection.
- **Tag protection** — the `release tags v*` ruleset blocks deleting or
  force-updating release tags (see repo rulesets).

## Local dry-run

```bash
make snapshot   # goreleaser release --snapshot --clean (no publish)
```

## Conventions

Commit messages follow Conventional Commits (`docs/patterns/git.md`); GoReleaser
groups `feat`/`fix` in the release notes and drops `docs`/`test`/`chore`/`ci`.
