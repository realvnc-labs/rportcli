---
title: Command and Script execution
slug: command-and-script-execution
weight: 3
---

## At a glance

Rportcli is a fast and efficient way to execute commands and scripts on one or many remote machines. Scripts can be
loaded from your local filesystem. If you store all your scripts for remote execution on a Git repository, you create
a highly professional team collaboration. Examples:

* Execute the command `hostname` and a client called "ANTMAN":

    ```shell
    rportcli command execute -n ANTMAN -c hostname
    ```

* Execute a PowerShell script loaded from the homedir on the clients "ABRAHAM" and "HOMER":

    ```shell
    rportcli script execute -n ABRAHAM,HOMER -s ~/date.ps1
    ```  

{{< hint type=note title="Commands vs. Scripts">}}
RPort can execute commands and scripts. [📖 Learn more about the difference](https://kb.rport.io/digging-deeper/commands-and-scripts#the-difference-between-commands-and-scripts).

On the following guide and examples only scripts are mentioned. Commands are working almost identically.
Just use `rportcli command ...` instead of `rportcli script ...`.
{{< /hint >}}

## Targeting

Targeting is the selection of clients on which the command or script is executed.

{{< hint type=note title="Verify targeting first">}}
If you are uncertain, how many clients will be effected by your targeting,
execute a **harmless command** like `hostname` **first**.
{{< /hint >}}

By name, `-n, --names string`
: Specify one or many names of clients separated by a comma. Wildcards `*` are supported.

```shell
$ rportcli command execute -n "Ben*,Cecil*" -c 'echo my name is $(hostname)'
Cecil-Rodriguez
    my name is Cecil-Rodriguez
Benjamin-Rogers
    my name is Benjamin-Rogers
Ben-Ray
    my name is Ben-Ray
Ben-Little
    my name is Ben-Little
Cecil-Fox
    my name is Cecil-Fox
Cecil-Turner
    my name is Cecil-Turner
Ben-Harvey
    my name is Ben-Harvey
Cecil-Wright
    my name is Cecil-Wright
```

**Caution:** Always surround wildcards with quotation-marks.

By IDs, `-d, --cids string`
: Specify one or many IDs of clients separated by a comma.

```shell
$ rportcli command execute -d \
05dd97dc361e44688fc3dad6deaad657,\
07706ecb447843e6ad131ddd5996c695,\
062180aa3b9348b7953befa2310fd940 -c date
Danielle-Simpson
    Wed Jul 27 12:39:17 UTC 2022
Eddie-White
    Wed Jul 27 12:39:17 UTC 2022
Norman-Jordan
    Wed Jul 27 12:39:17 UTC 2022
```

By client group IDs, `-g, --gids string`
: Specify one or many group IDs separated by a comma.

## Read from Yaml

Instead of specifying all options for the command or script execution on the command line,
you can read from a yaml file.

Example:

```yaml
# List of client ids
cids:
  - 07b7795935b7446bbfdbfd961eccb86d
  - 07fb556b17a14c43ae0f3b7fa6e7c1ff
# Concurrency, false by default, optional
conc: true
# Working directory, optional
cwd: /tmp
# output detailed information of a script execution, optional
full-command-response: false
# Interpeter, optional
interpreter: /bin/sh
# Use sudo, default false, optional
is_sudo: false
# Script embedded
exec: |
    date
    ls -la
    pwd
    whoami
```

With the above saved to `my-job.yaml`, execute

```shell
rportcli script execute -y my-job.yaml 
```

### Supported fields

In the yaml file you can use the following fields.

`cids`
: type=list, List of client ids the command or script will be executed on
: mutual exclusive with `gids` and `names`

`gids`
: type=list, List of group ids the command or script will be executed on
: mutual exclusive with `cids` and `names`

`names`
: type=list, List of client names the command or script will be executed on
: mutual exclusive with `cids` and `gids`

`conc`
: type=boolean, default=false, Run concurrent an all clients

`cwd`
: type=string, default=system temp folder, set the working directory

`full-command-respone`
: type=boolean, default=false, output detailed information of a script execution

`interpreter`
: type=string,default=/bin/sh (macOS,Linux) cmd.exe (Windows), set the script interpreter

`is_sudo`
: type=boolean, default=false, use sudo to run with root rights, MacOS/Linux only

`script`
: type=string, path to a script to be executed,
: required if script is not embedded, mutual exclusive with `exec`

`exec`
: type=string, script embedded instead of loading from a file.
: required if script is not file-based, mutual exclusive with `script`
