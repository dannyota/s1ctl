# Identity

Manage Identity AD Service configuration, connectors, and ISPM (Identity
Security Posture Management). The `identity` command group covers AD
configuration management, connector operations, onboarding status, and
ISPM exposure management. The existing `ranger-ad` group covers posture
reads (status, exposures, affected-objects) and assessment triggers.

## Prerequisites

- s1ctl [installed](install.md) and [configured](configure.md)
- `S1_CONSOLE_URL` and `S1_TOKEN` set
- Identity module enabled on your console

## AD configuration

Manage Active Directory configurations for identity posture scanning.

### List configurations

```bash
s1ctl identity config get
s1ctl identity config get --json
```

Credentials are redacted in the output.

### Add a configuration

```bash
s1ctl identity config add --from-file ad-config.json --yes
```

The add command is **dry-run by default**; pass `--yes` to apply.

### Delete a configuration

```bash
s1ctl identity config delete 000000 --yes
```

## Connectors

Manage AD connectors (Cloudlink agents) that bridge between the console and
your Active Directory environment.

```bash
s1ctl identity connector list
s1ctl identity connector get
s1ctl identity connector agents
```

Replace the connector with a different agent:

```bash
s1ctl identity connector replace --agent-id 000000 --yes
```

## Onboarding status

Check the AD service onboarding status:

```bash
s1ctl identity onboard
```

## Reference data

List available AD features and timezones for configuration:

```bash
s1ctl identity features
s1ctl identity timezones
```

List discovered AD domains:

```bash
s1ctl identity domains
```

## ISPM exposure management

Skip or acknowledge ISPM exposures to manage your identity posture backlog.
These commands work with exposure IDs from `ranger-ad exposures`.

### Skip exposures

Mark exposures as skipped (accepted risk):

```bash
s1ctl identity skip-exposures --ids 000000,000001 --yes
```

### Acknowledge exposures

Mark exposures as acknowledged:

```bash
s1ctl identity ack-exposures --ids 000000,000001 --yes
```

Both skip and acknowledge are **dry-run by default**; pass `--yes` to apply.

## Workflows

### Set up identity posture scanning

1. Check onboarding status: `s1ctl identity onboard`
2. Add AD configuration: `s1ctl identity config add --from-file ad.json --yes`
3. Verify connectors: `s1ctl identity connector list`
4. Run an assessment: `s1ctl ranger-ad assess --site-id 000000 --yes`
5. Review exposures: `s1ctl ranger-ad exposures --site-id 000000`

### Triage ISPM exposures

```bash
# List exposures
s1ctl ranger-ad exposures --site-id 000000

# Skip known-safe exposures
s1ctl identity skip-exposures --ids 000000 --yes

# Acknowledge reviewed exposures
s1ctl identity ack-exposures --ids 000000 --yes
```

## See also

- [Ranger AD](../commands/ranger-ad.md) — posture reads and assessment triggers
- [Settings](settings.md) — AD settings via the platform settings surface
