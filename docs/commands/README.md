# Command reference

Auto-generated from the command tree. Run `s1ctl docs generate` to update.

| Group | Commands | Description |
|-------|----------|-------------|
| [accounts](accounts.md) | count, expire, get, list, reactivate, uninstall-password | Manage accounts |
| [activities](activities.md) | count, export, list, types | View activity log |
| [agents](agents.md) | abort-scan, approve-uninstall, broadcast, count, decommission, disable, enable, fetch-files, fetch-firewall-rules, fetch-installed-apps, fetch-logs, firewall-logging, get, health, isolate, list, local-upgrade, local-upgrade-status, mark-up-to-date, move, move-to-site, outdated, passphrases, randomize-uuid, ranger, reconnect, reject-uninstall, reset-config, reset-passphrase, restart, scan, set-external-id, shutdown, uninstall, upgrade, versions | Manage endpoint agents |
| [alerts](alerts.md) | add-note, count, counts, delete-note, export, get, history, list, notes, resolve, stats, status, timeline, update-note, verdict | Manage unified alerts (GraphQL UAM) |
| [applications](applications.md) | cves, list, risks | Manage application inventory and risk |
| [assets](assets.md) | categories, overview | Manage XDR asset inventory |
| [audit](audit.md) | list | View mutation audit log |
| [blocklist](blocklist.md) | create, delete, export, list, pull, push, update, validate | Manage the blocklist (blocked file hashes) |
| [cloud-policies](cloud-policies.md) | delete, disable, enable, get, list, pull, push | Manage cloud security policies (CNS rules) |
| [cloud-rules](cloud-rules.md) | create, delete, disable, enable, evaluate, get, list, types, update | Manage CNS custom cloud rules (Cloud Native Security) |
| [config](config.md) | init, show | Manage s1ctl configuration |
| [datalake](datalake.md) | dashboards, facet, files, ingest, numeric, powerquery, query, saved-queries, timeseries | Query Singularity Data Lake (SDL) |
| [detection-library](detection-library.md) | data-sources, disable, enable, list, surfaces | Manage platform detection rules (detection library) |
| [devicecontrol](devicecontrol.md) | copy, delete, disable, enable, events, get, list, pull, push, reorder | Manage device control rules |
| [dlp](dlp.md) | classifications, rules, settings | Manage Data Loss Prevention (DLP) rules and classifications |
| [exclusions](exclusions.md) | create, delete, get, list, pull, push, update | Manage exclusions and blocklist |
| [filters](filters.md) | create, delete, list, update | Manage saved endpoint filters |
| [firewall](firewall.md) | copy, delete, disable, enable, export, get, import, list, protocols, pull, push, reorder | Manage firewall control rules |
| [global](global.md) | commands, completion, doctor, drift, help, version | Top-level commands |
| [groups](groups.md) | count, create, delete, get, list, pull, push, update | Manage groups |
| [iocs](iocs.md) | config, create, delete, list | Manage threat intelligence IOCs |
| [locations](locations.md) | create, delete, list, pull, push, update | Manage firewall locations |
| [maintenance](maintenance.md) | export, get, get-flexible, set, set-flexible | Manage task maintenance-window configuration |
| [mcp](mcp.md) | install, serve | Run Model Context Protocol server |
| [misconfigurations](misconfigurations.md) | add-note, assign, delete-note, export, get, history, list, notes, related-assets, status, update-note, verdict | Manage xSPM misconfigurations |
| [network](network.md) | configuration, copy, delete, disable, enable, export, get, import, list, move, protocols, pull, push, reorder, set-location, tags | Manage network quarantine rules |
| [policies](policies.md) | diff, get, list, pull, push, revert | Manage endpoint policies |
| [ranger-ad](ranger-ad.md) | affected-objects, assess, exposures, status | Manage Ranger AD exposure assessments (ISPM) |
| [remoteops](remoteops.md) | content, get, guardrails, list, pending, results, run, update, upload-limits | Manage remote operations and scripts |
| [reports](reports.md) | create, download, list, tasks, types | Manage reports and report tasks |
| [roles](roles.md) | create, delete, get, list, template, update | Manage RBAC roles |
| [rules](rules.md) | detections, diff, disable, enable, get, health, list, pull, push, trends, validate | Manage custom detection rules (STAR) |
| [service-users](service-users.md) | bulk-delete, create, delete, export, generate-token, get, list, update | Manage service users (API-token identities) |
| [settings](settings.md) | cancel-pending-emails, delete-recipient, get, list, sso-cert, test, update | Manage platform settings |
| [sites](sites.md) | count, create, delete, duplicate, expire, get, licenses, list, pull, push, reactivate, regenerate-key, token, update | Manage sites |
| [status](status.md) | capabilities, enums, surfaces | Show environment health summary |
| [system](system.md) | info | Show console system information |
| [tag-rules](tag-rules.md) | create, delete, list, test, update | Manage dynamic asset tag rules |
| [tags](tags.md) | create, delete, get, list, pull, push, update | Manage tags |
| [threats](threats.md) | add-note, add-to-exclusions, blacklist, count, exclusion-options, export, fetch-file, get, list, mitigate, mitigate-alerts, notes, quarantined-files, resolve, set-ticket, status, timeline, verdict | Manage threats |
| [unified-exclusions](unified-exclusions.md) | create, export, list | Manage unified exclusions |
| [updates](updates.md) | get, list | Manage agent update packages |
| [upgrade-policies](upgrade-policies.md) | activate, create, deactivate, delete, get, list, packages, update | Manage agent auto-upgrade policies |
| [users](users.md) | 2fa, delete, generate-token, get, list, revoke-token, token-details, update | Manage users |
| [visibility](visibility.md) | query | Run Deep Visibility queries |
| [vulnerabilities](vulnerabilities.md) | add-note, assign, cve, cves, delete-note, export, get, health, history, list, notes, related-assets, stats, status, update-note, verdict | Manage xSPM vulnerabilities |
