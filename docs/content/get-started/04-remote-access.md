---
title: "Remote Access"
weight: 4
slug: remote-access
---
{{< toc >}}

## At a glance

With rportcli you can create tunnels to your remote machines. Optionally, you can directly launch an application
such as ssh or remote desktop to instantly log in. Examples:

* Create an SSH tunnel to a client called "Juan-Ford" and execute the openSSH client on it

    ```shell
    rportcli tunnel create -n Juan-Ford -s ssh -b "-l root"
    ```

* Create an RDP tunnel to a client called ABRAHAM and start the remote desktop using the Administrator user.

    ```shell
    rportcli tunnel create -n ABRAHAM -s rdp -d -u Administrator
    ```

## Create tunnels

With `rportcli tunnel create` you can establish a tunnel to a client.
You have three mutual exclusive options to specify the client the tunnel is created for. Also called **targeting**.

By name `-n, --name string`
: The client is identified by its name. Example:
: `rportcli tunnel create -n My-Remote-Machine`
: 🧙‍♂️ *Wildcards are supported.* If you don't want to type in a long name, use `-n "Alvin*"` for example.
: If more than one client matches the wildcard search, you will get an error.
: Don't omit the quotation marks. Otherwise, the wildcard sign `*` is resolved by your shell.  

By ID `-c, --client string`
: The client is identified by its ID. Example:
: `rportcli tunnel create -c 0658aeabf8e04f759dd2b9dbbec53068`

You can either copy the client's name or ID from the web user interface or use `rportcli client list`.

Once you have a client identified, specify for which port you want to create the tunnel.

The remote port, `-r, --remote string`
: This is the port on the remote machine you wish to get access to.

The other end of the tunnel, `-l, --local string`
: This is the port on your rport server you will use to access the tunnel.
: If none is given, a random free port is selected.

## Examples

Create a **tunnel to the port 22 (SSH)** of the client called "Holly-Harris".
A random free port on the rport server is used.

```shell
$ rportcli tunnel create -n Holly-Harris -r 22
Tunnel
KEY                VALUE                                  
ID:                1                                      
CLIENT_ID:         Holly-Harris                           
LOCAL_HOST:        0.0.0.0                                
LOCAL_PORT:        29853                                  
REMOTE_HOST:       127.0.0.1                              
REMOTE_PORT:       22                                     
LOCAL_PORT RANDOM: true                                   
SCHEME:            ssh                                    
IDLE TIMEOUT MINS: 5                                      
ACL:               89.0.79.56                             
USAGE:             ssh -p 29853 rport.example.com -l ${USER}
```

Create a so-called **service forwarding**. The host "ANTMAN" is used as a bridge host to create a tunnel for RDP to a
remote machine where rport is not installed.

```shell
rportcli tunnel create --name ANTMAN -r 192.168.219.46:3389 -s rdp
Tunnel
KEY                VALUE                                  
ID:                2                                      
CLIENT_ID:         ANTMAN                                 
LOCAL_HOST:        0.0.0.0                                
LOCAL_PORT:        27139                                  
REMOTE_HOST:       192.168.219.46                         
REMOTE_PORT:       3389                                   
LOCAL_PORT RANDOM: true                                   
SCHEME:            rdp                                    
IDLE TIMEOUT MINS: 5                                      
ACL:               89.0.79.56                             
USAGE:             rdp://rport.example.com:27139 
```

See `rportcli tunnel create -h` for all options.

{{< hint type=important icon=gdoc_shield title="Secured by default" >}}
All tunnels are protected with a **tight access control list** (ACL) by default. Only the current public IP Address of
the rportcli host will be allowed to access the port of the tunnel. Use `-a, --acl` to create custom ACLs.

Also, all tunnels are **closed automatically** after an inactivity of 5 minutes. Use `-m, --idle-timeout-minutes int`
to change this behaviour.
{{< /hint >}}

## Close tunnels

Use `rportcli tunnel list` to display the list of active tunnels.
Then use `rportcli tunnel delete -c <CLIENT-ID> -u <TUNNEL-ID>` to delete a tunnel for a client.

See `rportcli tunnel delete -h` for all options.

## Time-saving shortcuts: Create and launch tunnels 🏎

For the two most widely used remote access protocols, SSH and RDP, rportcli has built-in shortcuts.
After a tunnel is created, openSSH or Microsoft Remote Desktop will automatically start a session.

### Remote Desktop

Create a tunnel for RDP to the host identified by its name. The Remote Desktop Client will automatically start
using the username "Administrator" and a geometry of 1024*768.

```shell
rportcli tunnel create -n ABRAHAM -s rdp -d -u Administrator -w 1024 -i 768
```

Behind the scenes rportcli creates a temporary `.rdp` file and then the default app for this file type is launched.

### SSH

Create a tunnel for SSH to the host identified by its name. The openSSH client is started with the ssh
options `-l root -A`, meaning that the SSH user is `root` and the ssh-agent is passed into the session.

```shell
rportcli tunnel create -n Juan-Ford -b "-l root -A"
```

Rportcli will directly launch `ssh` - the open-ssh client - with the port of the tunnel appended.
Other SSH client such as Putty for example are not supported.

The tunnel is closed automatically after the ssh client is closed. You don't need to do this manually.

*This applies only to ssh. For other apps a tunnel close can't be triggered on app close.
Tunnels will close after 5 minutes without network activity.*  

### URI Open

To access `http`, `https`, `vnc` or `realvnc` services via a tunnel, you can use the generic URI launcher.
For example:

```shell
rportcli tunnel create -n MyMachine -s http --launch-uri
```

After the tunnel is created, your default app for the specified URI is launched.
