---
title: "RPort CLI"
draft: false
---

## At a glance

The rport command line interface `rportcli` is a great addition to the rport server. Executing some tasks on the
command line can be much faster and more efficient than doing it on the web user interface.

Rportcli does not cover all functions of the user interface. On the other hand, you can do things on the command line,
you cannot do on the user interface.

Rportcli can be integrated into scripts (bash, zsh, PowerShell, etc.) giving you endless options for automation.

The rport server has a built-in library for storing scripts. But with rportcli you can also store, share and execute
scripts using a version control system like Git.

## Built-in help

Rportcli comes with comprehensive help built-in. Type in `rportcli help`. Each sub-command has its own help.
For example, type in `rportcli tunnel -h` and drill down with `rportcli tunnel create -h`.

```text
$ rportcli help
rportcli

Usage:
  rportcli [command]

Available Commands:
  client      manage rport clients
  command     command management
  help        Help about any command
  init        initialize your connection to the rportd API
  me          show current user info
  script      scripts management
  tunnel      manage tunnels of connected clients
  version     print the version number of rportcli

Flags:
  -h, --help             help for rportcli
  -j, --json-pretty      in combination with json format this flag will pretty print the json data
  -o, --output string    Output format: json, yaml or human (default "human")
  -t, --timeout string   Timeout value as seconds, e.g. 10s, minutes e.g. 1m or hours e.g. 2h, 
                         if not provided no timeout will be set
  -v, --verbose          verbose output

Use "rportcli [command] --help" for more information about a command.
```
