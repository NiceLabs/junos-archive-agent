# JunOS Archive Agent

The project based <https://www.juniper.net/documentation/en_US/junos/topics/task/configuration/junos-software-system-management-router-configuration-archiving.html>

based `transfer-on-commit` command automate backup junos configuration to github repository

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

## Configure file

```json
{
  "ftp": {
    "port": 21,
    "host": "::",
    "username": "juniper",
    "password": "[your password]"
  },
  "github": {
    "owner": "[your owner]",
    "repo": "[your repo]",
    "token": "[your personal token]"
  }
}
```

## Logging

```text
# docker-compose logs
Attaching to junos-archive-agent_agent_1
agent_1  | 2020/05/16 06:01:31   JunOS Archive listening on 21
agent_1  | 2020/05/16 06:06:32 2cd57795b3724c991be7  Connection Established
agent_1  | 2020/05/16 06:06:32 2cd57795b3724c991be7 < 220 JunOS FTP Archive Server
agent_1  | 2020/05/16 06:06:32 2cd57795b3724c991be7 > USER juniper
agent_1  | 2020/05/16 06:06:32 2cd57795b3724c991be7 < 331 User name ok, password required
agent_1  | 2020/05/16 06:06:33 2cd57795b3724c991be7 > PASS ****
agent_1  | 2020/05/16 06:06:33 2cd57795b3724c991be7 < 230 Password ok, continue
agent_1  | 2020/05/16 06:06:33 2cd57795b3724c991be7 > TYPE I
agent_1  | 2020/05/16 06:06:33 2cd57795b3724c991be7 < 200 Type set to binary
agent_1  | 2020/05/16 06:06:33 2cd57795b3724c991be7 > CWD .
agent_1  | 2020/05/16 06:06:33 2cd57795b3724c991be7 < 250 Directory changed to /
agent_1  | 2020/05/16 06:06:33 2cd57795b3724c991be7 > PORT 218,81,2,86,201,52
agent_1  | 2020/05/16 06:06:33 2cd57795b3724c991be7  Opening active data connection to *.*.*.*:51508
agent_1  | 2020/05/16 06:06:34 2cd57795b3724c991be7 < 200 Connection established (51508)
agent_1  | 2020/05/16 06:06:34 2cd57795b3724c991be7 > STOR gateway_20200516_060553_juniper.conf.gz
agent_1  | 2020/05/16 06:06:34 2cd57795b3724c991be7 < 150 Data transfer starting
agent_1  | 2020/05/16 06:06:38 2cd57795b3724c991be7 < 226 OK, received 8651 bytes
agent_1  | 2020/05/16 06:06:39 2cd57795b3724c991be7  Connection Terminated
```
