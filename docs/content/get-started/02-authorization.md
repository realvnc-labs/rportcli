---
title: "Authorization"
draft: false
weight: 1
slug: authorization
---

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

## by API token

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
