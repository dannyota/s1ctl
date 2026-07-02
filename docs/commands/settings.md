# settings

Manage platform settings

## settings cancel-pending-emails

Cancel queued pending email notifications

```text
s1ctl settings cancel-pending-emails [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope to account IDs |
| `--site-id` | stringSlice | - | scope to site IDs |
| `--yes` | bool | false | apply the action (default: dry-run) |

## settings delete-recipient

Delete a notification recipient

```text
s1ctl settings delete-recipient <id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## settings get

Get settings configuration

```text
s1ctl settings get <type> [flags]
```

Get configuration for a specific settings type.

Types: notifications, sso, smtp, syslog, sms, recipients, ad, ad-scope-mapping

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--site-id` | stringSlice | - | filter by site ID |

## settings list

List settings categories

```text
s1ctl settings list
```

## settings sso-cert

Show or download the SSO service-provider signing certificate

```text
s1ctl settings sso-cert [flags]
```

Show the SAML service-provider signing certificate. The certificate is
public key material, not a secret. With --out, download the raw certificate
file to disk; otherwise print its metadata and PEM.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--out` | string | - | write the downloaded certificate to this file |
| `--site-id` | stringSlice | - | filter by site ID |

## settings test

Test settings connectivity

```text
s1ctl settings test <type> [flags]
```

Test connectivity for SMTP, syslog, or Active Directory settings.

Types: smtp, syslog, ad

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--site-id` | stringSlice | - | filter by site ID |
| `--yes` | bool | false | apply the action (default: dry-run) |

## settings update

Update settings from a JSON file (pull with 'settings get', edit, push back)

```text
s1ctl settings update
```
