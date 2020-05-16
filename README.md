# Juniper Archive Agent

The project based <https://www.juniper.net/documentation/en_US/junos/topics/task/configuration/junos-software-system-management-router-configuration-archiving.html>

## JunOS Configuration

```text
set system archival configuration transfer-on-commit
set system archival configuration archive-sites ftp://<username>@<host>:<port>/<url-path> password <password>
```

e.q:

```text
set system archival configuration transfer-on-commit
set system archival configuration archive-sites ftp://juniper@example.com password "your password"
```
