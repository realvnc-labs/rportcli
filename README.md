
[![License](https://img.shields.io/github/license/cloudradar-monitoring/rportcli?style=for-the-badge)](https://github.com/cloudradar-monitoring/rportcli/blob/main/LICENSE)
# Rport CLI (v1)

Rport CLI is a tool to help you managing your remote machines connected 
to the [rport API](https://github.com/cloudradar-monitoring/rport) directly from your terminal.

The rport command line interface `rportcli` is a great addition to the rport server. Executing some tasks on the
command line can be much faster and more efficient than doing it on the web user interface.

Rportcli does not cover all functions of the user interface. On the other hand, you can do things on the command line,
you cannot do on the user interface.

Integrates well with:

![](https://img.shields.io/badge/Powershell-2CA5E0?style=for-the-badge&logo=powershell&logoColor=white) 
![](https://img.shields.io/badge/GNU%20Bash-4EAA25?style=for-the-badge&logo=GNU%20Bash&logoColor=white)
![](https://img.shields.io/badge/GIT-E44C30?style=for-the-badge&logo=git&logoColor=white)

## Documentation

[![documentation](https://img.shields.io/badge/Documentation-Read_Now-green?style=for-the-badge&logo=Gitbook)](https://cli.rport.io)

Read the documentation on [https://cli.rport.io/](https://cli.rport.io/)

## Using the Cli

Trigger this command to see all available commands and their options:

    rportcli --help

You can also display help for a certain command:

    rportcli init --help

## Additional configuration with environment variables

<table>
    <tr>
    <th>Variable</th>
    <th>Description</th>
    <th>Default Value</th>
    <th>Example</th>
    </tr>
    <tr>
    <td>CONFIG_PATH</td>
    <td>Changes default config location to the provided value</td>
    <td>${HOME}/.config/rportcli/config.json</td>
    <td>CONFIG_PATH=/tmp/config.json rportcli init</td>
    </tr>
    <tr>
    <td>CONN_TIMEOUT_SEC</td>
    <td>Connection timeout to call rport server (seconds)</td>
    <td>10 seconds</td>
    <td>CONN_TIMEOUT_SEC=20 rportcli client list</td>
    </tr>
    <tr>
    <td>SESSION_VALIDITY_SECONDS</td>
    <td>Initial lifetime of an interactive command session in seconds. Max value is 90 days</td>
    <td>10(minutes) * 60</td>
    <td>SESSION_VALIDITY_SECONDS=1800 rportcli command -i</td>
    </tr>
    <tr>
    <td>RPORT_API_USER</td>
    <td>Basic auth login to access rport server</td>
    <td></td>
    <td>RPORT_API_USER=admin rportcli client list</td>
    </tr>
    <tr>
    <td>RPORT_API_PASSWORD</td>
    <td>Basic auth password to access rport server</td>
    <td></td>
    <td>RPORT_API_PASSWORD=foobaz rportcli client list</td>
    </tr>
    <tr>
    <td>RPORT_API_URL</td>
    <td>Address of rport server</td>
    <td></td>
    <td>RPORT_API_URL=http://localhost:3000 rportcli client list</td>
    </tr>
    <tr>
    <td>RPORT_API_TOKEN</td>
    <td>Api token for accessing the rport server. Must be specified with RPORT_API_USER and RPORT_API_URL.</td>
    <td></td>
    <td>RPORT_API_TOKEN=xxxxxxxx rportcli client list</td>
    </tr>
</table>

## Using YAML input for Options (command and script only)

For `command execute` and `script execute`, the CLI supports reading of the required options / parameters
from yaml files via the `--read-yaml` (short form `-y`) command line option. This will allow users to set the CLI options
for execution in files, rather than always being required on the command line itself.

Any related options specified on the command line will have precedence over options set in YAML files.

Multiple yaml files are supported but only to set the options for a single execution. The last yaml file
on the command line with have precedence and any duplication options will overwrite options set in previously included
yaml files.

For example

```yaml
# check-clients.yaml
cids:
  - cdeb33642b4b43caa13b73ce0045d388
  - 7ca5718bd76f1bca7a5ee72660d3120c
  - 42560923b8414a519c7a42047f251fb3
conc: true
full-command-response: true
command: |
  ls
```

Can be executed using

```bash
$ rport_cli -y "check-clients.yml
```

And the targeted clients can be overridden using:

```bash
$ `rport_cli -y "check-clients.yml --name "test-server"`
```
