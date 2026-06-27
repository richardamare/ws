# Azure CLI / SP usage

Read `../security/README.md` first — this file is the mechanics; that file is the rules.

## Scoped login

Everyday sessions log in as the project's **Reader SP** into an **isolated config dir**:

```bash
export AZURE_CONFIG_DIR=~/.azure-proj1
az login --service-principal -u <appId> -p ~/.config/ws/certs/proj1.pem --tenant <tenantId>
```

`ws` sets `AZURE_CONFIG_DIR` for every tab it opens, so the session can never see the personal admin
token in `~/.azure`.

## Idempotent auth

`ws auth` / the auth step of `ws up`:

1. `AZURE_CONFIG_DIR=<dir> az account show` → if it returns the expected SP, **skip** (already logged in).
2. Else `az login --service-principal …`.
3. If the cert is expired/near-expiry, surface it and suggest `ws rotate`.

## Creating the SP (`ws new`)

```bash
az ad sp create-for-rbac \
  --name sp-ws-<proj>-reader \
  --role Reader \
  --scopes /subscriptions/<sub>/resourceGroups/<rg> \
  --create-cert
az ad sp update --id <appId> --set notes="ws reader, <proj>"   # for easy cleanup
```

## Cert handling

- Store at `~/.config/ws/certs/<proj>.pem`, **`chmod 600`**.
- SP certs expire (~1 yr). `ws auth` detects; `ws rotate` issues a new cert and updates the SP.

## Verifying confinement

```bash
AZURE_CONFIG_DIR=~/.azure-proj1 az group show  -n <rg>           # ✅ works (Reader)
AZURE_CONFIG_DIR=~/.azure-proj1 az group delete -n other-rg --yes # ❌ AuthorizationFailed
```

## Never

- Never use the SP for write. Never grant it more than Reader on the one RG.
- Never run `az group delete`, `az role assignment create/delete`, `az ad sp delete` from automation —
  they're in the `.claude/settings.json` deny list.
- Write/Terraform → `ws elevate` (personal `az login`), human-approved.
