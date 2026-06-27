# Session bookmarks

## The problem

Claude Code's `--resume` picker becomes unusable with many sessions — hard to tell which has the good
context. The user currently keeps ~6 cmux workspaces open **just so he doesn't lose** specific sessions.

## The fix

A short, **curated, named** list per project — stored in the project YAML — instead of scrolling
Claude's full session history. Closing a workspace no longer means losing a session.

```yaml
sessions:
  - { label: auth-refactor,       id: 3ee3..., note: "good RBAC context" }
  - { label: terraform-bootstrap, id: 9ab1..., note: "infra setup" }
```

## Flow

```bash
ws save proj1 auth-refactor    # bookmark the current session
ws sessions proj1              # see your short labeled list
ws resume proj1 auth-refactor  # → claude --resume <id>
ws resume proj1                # no label → interactive picker
```

## Getting the session id to save

cmux already tracks the focused agent's resume id:

```bash
cmux surface resume show
```

`ws save` reads that and writes the bookmark. (Fallback: Claude Code writes session ids to its own
project history; `ws` can read the latest if cmux can't supply it.)

## Payoff

The 6 always-open workspaces collapse to a few lines in a file. Reopen any good-context session on
demand with `ws resume`. The user only ever sees **his own labels**, never Claude's raw resume list.
