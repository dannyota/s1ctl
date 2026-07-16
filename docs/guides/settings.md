# Settings

Read and update platform settings: notifications, SSO, SMTP, syslog, SMS,
recipients, Active Directory, and AD scope mapping.

## Categories

List available settings categories:

```bash
s1ctl settings list
```

| Category | CLI token | Description |
|----------|-----------|-------------|
| notifications | `notifications` | Notification preferences and alert routing |
| sso | `sso` | SSO/SAML authentication configuration |
| smtp | `smtp` | SMTP email server configuration |
| syslog | `syslog` | Syslog forwarding configuration |
| sms | `sms` | SMS notification service configuration |
| recipients | `recipients` | Notification recipient list |
| active-directory | `ad` | Active Directory integration |
| (sub-category) | `ad-scope-mapping` | Active Directory scope mapping |

> **Note:** `settings list` shows `active-directory` as the category name, but
> the CLI token for `settings get` and `settings update` is `ad`. Similarly,
> AD scope mapping is accessed via the token `ad-scope-mapping`.

## Read settings

```bash
s1ctl settings get notifications
s1ctl settings get sso --json
s1ctl settings get smtp --site-id 000000
s1ctl settings get syslog --account-id 000000
s1ctl settings get sms
s1ctl settings get recipients
s1ctl settings get ad
s1ctl settings get ad-scope-mapping
```

| Flag | Description |
|------|-------------|
| `--site-id` | Scope to site ID (repeatable) |
| `--account-id` | Scope to account ID (repeatable) |
| `--json` | Machine-readable output |

Sensitive fields (passwords, tokens, certificate contents) are redacted in
the output.

## Update settings

The update workflow is pull-edit-push: read the current settings to a JSON
file, edit the file, then push it back.

```bash
# 1. Pull current settings
s1ctl settings get smtp --json > smtp.json

# 2. Edit smtp.json (set host, port, credentials, etc.)

# 3. Push (dry-run first)
s1ctl settings update smtp --from-file smtp.json
s1ctl settings update smtp --from-file smtp.json --yes
```

All update commands are **dry-run by default**; pass `--yes` to apply.

Available update subcommands:

```bash
s1ctl settings update notifications --from-file settings.json --yes
s1ctl settings update sso --from-file sso.json --yes
s1ctl settings update smtp --from-file smtp.json --yes
s1ctl settings update syslog --from-file syslog.json --yes
s1ctl settings update sms --from-file sms.json --yes
s1ctl settings update recipients --from-file recipients.json --yes
s1ctl settings update ad --from-file ad.json --yes
s1ctl settings update ad-scope-mapping --from-file mapping.json --yes
```

| Flag | Description |
|------|-------------|
| `--from-file` | JSON file with the settings payload (required) |
| `--site-id` | Scope to site ID (repeatable) |
| `--account-id` | Scope to account ID (repeatable) |
| `--yes` | Apply the update (default: dry-run) |

> **Warning:** The pulled JSON has sensitive fields redacted. Before pushing,
> re-enter any secret values (passwords, tokens, certificate contents) --
> otherwise the update writes them back empty.

## Test connectivity

Test SMTP, syslog, or Active Directory connectivity using the currently
configured settings:

```bash
s1ctl settings test smtp --yes
s1ctl settings test syslog --site-id 000000 --yes
s1ctl settings test ad --yes
```

The test action is **dry-run by default**; pass `--yes` to run the test.

## SSO certificate

Show or download the SAML service-provider signing certificate:

```bash
# Show certificate metadata and PEM
s1ctl settings sso-cert

# Download certificate to a file
s1ctl settings sso-cert --out sp-cert.pem
```

The certificate is public key material, not a secret.

## Recipients

List recipients:

```bash
s1ctl settings get recipients
```

Delete a recipient:

```bash
s1ctl settings delete-recipient 000000 --yes
```

## Cancel pending emails

Cancel queued email notifications that have not yet been sent:

```bash
s1ctl settings cancel-pending-emails --yes
s1ctl settings cancel-pending-emails --site-id 000000 --yes
```

## Config overrides

Config overrides change agent behavior at a selected scope (tenant, account,
site, or group). They are powerful — they override the agent's configuration
at that scope.

### List and inspect

```bash
s1ctl settings overrides list
s1ctl settings overrides list --site-id 000000
s1ctl settings overrides get 000000
```

### Create

```bash
s1ctl settings overrides create --name "Custom override" --scope site --scope-id 000000 \
  --os-type linux --config '{"agent.maxLogSize": 100}' --yes
```

### Update

```bash
s1ctl settings overrides update 000000 --config '{"agent.maxLogSize": 200}' --yes
```

### Delete

```bash
s1ctl settings overrides delete 000000 --yes
```

All mutations are **dry-run by default**; pass `--yes` to apply.

## Sentinel Deploy

Manage credential groups used by Sentinel Deploy (Ranger auto-deploy) to
install agents on unprotected endpoints. Accessed via `updates deploy`.

### List credential groups and details

```bash
s1ctl updates deploy list-groups
s1ctl updates deploy list-details --cred-group-id 000000
```

### Create and delete groups

```bash
s1ctl updates deploy create-group --group-name "Deploy creds" --scope-id 000000 \
  --target-os windows --yes
s1ctl updates deploy delete-group 000000 --yes
```

### Manage credential details

```bash
s1ctl updates deploy add-detail --cred-group-id 000000 --title "Admin" \
  --cred-type "User/Password" --yes
s1ctl updates deploy update-detail 000000 --title "Updated Admin" --yes
s1ctl updates deploy delete-detail 000000 --yes
```

All mutations are **dry-run by default**; pass `--yes` to apply.

## Workflows

### Audit SMTP configuration across sites

```bash
s1ctl settings get smtp --site-id 111111 --json > smtp-a.json
s1ctl settings get smtp --site-id 222222 --json > smtp-b.json
diff smtp-a.json smtp-b.json
```

### Update syslog forwarding

```bash
s1ctl settings get syslog --json > syslog.json
# edit syslog.json: change host, port, enable SSL
s1ctl settings update syslog --from-file syslog.json        # dry-run
s1ctl settings update syslog --from-file syslog.json --yes  # apply
s1ctl settings test syslog --yes                             # verify
```
