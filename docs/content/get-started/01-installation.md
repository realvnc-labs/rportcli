---
title: "Installation"
weight: 1
slug: installation
---
{{< toc >}}

RPort cli is a single static binary available for almost any operating system.
Basically, the installation is just downloading a file.

## On macOS

```shell
wget https://downloads.rport.io/rportcli/stable/?arch=Darwin_x86_64 -O rportcli.tar.gz
tar xzf rportcli.tar.gz rportcli
sudo mv rportcli /usr/local/bin
rm rportcli.tar.gz
```

{{< hint type=note title="Apple M1" >}}
Currently, no binaries are available for the Darwin arm64 architecture aka the M1 CPU.
We are working on the set-up an apple-compliant build process and the required code signing.
Meanwhile, you must [build the binary from the sources](https://github.com/cloudradar-monitoring/rportcli#install-as-a-go-binary).
{{< /hint >}}

## On Linux

```shell
wget https://downloads.rport.io/rportcli/stable/?arch=Linux_$(uname -m) -O rportcli.tar.gz
tar xzf rportcli.tar.gz rportcli
sudo mv rportcli /usr/local/bin
rm rportcli.tar.gz
```

## On Windows

Because on Windows a default folder for command line utilities does not exist,
we will create a new folder in `C:\Program Files`.

```powershell
New-Item "C:\Program Files\rportcli" -itemType Directory
cd "C:\Program Files\rportcli"
iwr https://downloads.rport.io/rportcli/stable/?arch=Windows_x86_64.zip -OutFile rportcli.zip
Expand-Archive -Path .\rportcli.zip -DestinationPath .
New-Item bin -itemType Directory
Move-Item .\rportcli.exe .\bin\
Remove-Item .\rportcli.zip
# Add RPort CLI to your path
$Env:PATH="$Env:PATH;C:\Program Files\rportcli\bin"
[Environment]::SetEnvironmentVariable(
    "PATH", $Env:PATH + ";C:\Program Files\rportcli\bin", [EnvironmentVariableTarget]::Machine
)
& rportcli --version
```

## Install as a go binary

If you have go installed, try the following

```shell
go get github.com/cloudradar-monitoring/rportcli
```

## Update

To update rportcli to the latest version, perform the above steps to overwrite your existing version
with the latest one.
