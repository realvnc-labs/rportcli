---
title: "Authorization"
draft: false
weight: 1
slug: authorization
---
{{< toc >}}

You have two options to authorize the rportcli on the API of your rport server.

1. Using a username, password and second factor (if applicable).
   These are the same credentials you use for the webinterface.
   A token (JWT) will be created and stored in your home directory.
   The stored token has a limited lifetime and after the expiry you will be asked for your credentials again.  
   This authentication method is not recommended, if you want to integrate rportcli into unattended script like
   cronjobs.
2. Using an API token. This requires you have created API token for your user account first.
   The token is stored as an environment variable. No JWT is generated and no local configuration file is created.
   API tokens do not have an expiry date.
   This authentication method is suitable for unattended scripting. Store the API token securely.

## by username, password and 2FA

Just execute `rportcli init` and enter the URL of your RPort server, username, password and, if applicable,
the second factor of the two-factor authentication. Just enter the base URL of rport e.g.
`https://rport.example.com` without the api path.

If you do that for the first time, you will be informed about the new configuration file created.

Test the connection by executing `rportcli me`.

Rportcli looks for the config file at `\$HOME/.config/rportcli/config.json` (for Linux and macOS) or
`C:\Users\<CurrentUserName>\.config\rportcli\config.json` (for Windows).
If the current user has no home folder, RportCli will look for a config file next to the current binary location.

You can override config path by providing an environment variable CONFIG_PATH, e.g.
`CONFIG_PATH=/tmp/config.json rportcli init`.

With the environment variable `SESSION_VALIDITY_SECONDS` you can set initial lifetime of an interactive command session 
in seconds. Max value is 90 days.

## by API token

The easiest and most flexible way to authenticate with the rportd server is to use an API
token stored as environment variable. Using `RPORT_API_TOKEN` will bypass two-factor authentication,
allowing the rport cli to be used in automated scenarios.

Also, if using RPORT_API_TOKEN then the config file will be ignored completed, so the
`RPORT_API_URL` must be used.

If you have created an API token for your user account, put the following variables to your environment:
`RPORT_API_TOKEN`, `RPORT_API_URL`, `RPORT_API_USER`.

{{< tabs "env" >}}
{{< tab "macOS/Linux" >}}
    export RPORT_API_URL=<https://rport.example.com>
    export RPORT_API_TOKEN=1234abc
    export RPORT_API_USER=john
{{< /tab >}}
{{< tab "Windows" >}}
    $env:RPORT_API_URL="https://rport.example.com"
    $env:RPORT_API_TOKEN="1234abc"
    $env:RPORT_API_USER="john"
{{< /tab >}}
{{< /tabs >}}

## deprecation notes

`RPORT_API_USER`, `RPORT_API_PASSWORD` and `RPORT_API_URL` replace the previous `RPORT_USER`,
`RPORT_PASSWORD` and `RPORT_SERVER_URL` environment variables. Please update any scripts
accordingly. Support for `RPORT_USER`, `RPORT_PASSWORD` and `RPORT_SERVER_URL` will be removed
in a future release.
