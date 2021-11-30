# Rport CLI (v1)
Rport CLI is a tool to help you managing [rport API](https://github.com/cloudradar-monitoring/rport) directly from your terminal.

## Installation

### As a compiled binary

Jump to [our release page](https://github.com/cloudradar-monitoring/rportcli/releases/tag/v0.0.1pre1) and download a binary for your host OS. Don't forget to download a corresponding md5 file as well.


    # On MacOS
    wget https://github.com/cloudradar-monitoring/rportcli/releases/download/v0.0.1pre1/rportcli-v0.0.1pre1-darwin-amd64.tar.gz
    
    # On linux
    wget https://github.com/cloudradar-monitoring/rportcli/releases/download/v0.0.1pre1/rportcli-v0.0.1pre1-linux-386.tar.gz
    
    # On Windows
    Just download https://github.com/cloudradar-monitoring/rportcli/releases/download/v0.0.1pre1/rportcli-v0.0.1pre1-windows-amd64.zip
    Also download https://github.com/cloudradar-monitoring/rportcli/releases/download/v0.0.1pre1/rportcli-v0.0.1pre1-windows-amd64.zip.md5
     
     
Verify the checksum:

    
    #On MacOS
    curl -Ls https://github.com/cloudradar-monitoring/rportcli/releases/download/v0.0.1pre1/rportcli-v0.0.1pre1-darwin-amd64.tar.gz.md5 | sed 's:$: rportcli-v0.0.1pre1-darwin-amd64.tar.gz:' | md5sum -c
    
    #On linux
    curl -Ls https://github.com/cloudradar-monitoring/rportcli/releases/download/v0.0.1pre1/rportcli-v0.0.1pre1-linux-386.tar.gz.md5 | sed 's:$: rportcli-v0.0.1pre1-linux-386.tar.gz:' | md5sum -c
     
    #On Windows assuming you're in the directory with the donwloaded file
    CertUtil -hashfile rportcli-v0.0.1pre1-windows-amd64.zip MD5
    
    #The output will be :
    MD5 hash of tacoscript-0.0.4pre-windows-amd64.zip:
    7103fcda170a54fa39cf92fe816833d1
    CertUtil: -hashfile command completed successfully.
    
    #Compare the command output to the contents of file rportcli-v0.0.1pre1-windows-amd64.zip.md5 they should match
  

_Note: if the checksums didn't match please skip the installation!_

Unpack and install the rportcli binary on your host machine

    
    #On linux/MacOS
    tar -xzvf rportcli-v0.0.1pre1-darwin-amd64.tar.gz 
    mv rportcli /usr/local/bin/rportcli
    chmod +x /usr/local/bin/rportcli
    

For Windows

- Extract file contents

- Create a `RportCLI` folder in `C:\Program Files`

- Copy the rportcli.exe binary to `C:\Program Files\RportCLI`

- Double click on the rportcli.exe and allow it's execution

## Install as a go binary:

    go get github.com/cloudradar-monitoring/rportcli

## Config

Rportcli looks for a config file at $HOME/.config/rportcli/config.json (for Linux and MacOS) or C:\Users\<CurrentUserName>\.config\rportcli\config.json (for Windows).
If current user has no home folder, RportCli will look for a config file next to the current binary location.

You can override config path by providing env variable CONFIG_PATH, e.g.


    CONFIG_PATH=/tmp/config.json rportcli init
    
    
You can generate config by running:


    rportcli init


Rportcli will interactively ask for config options and validate the result:
You'll get request for following parameters:

**server address** 

address of rport server, e.g. `http://localhost:3000`

**login**

basic auth login to access rport server, e.g. `admin`

**password**

basic auth password to access rport server, e.g. `foobaz`

You can skip the interactive wizard by providing parameters as cli options , e.g.


     rportcli init -s http://localhost:3000 -l admin -p foobaz


If you prefer to use environment variables instead you can do the following: 


    export RPORT_USER=admin
    export RPORT_PASSWORD=foobaz
    export RPORT_SERVER_URL=http://localhost:3000
    #now you can run any rportcli command without config
    rportcli client list


You can also use a hybrid variant, where e.g. user and server url are provided as config options and password as an environment variable.


     rportcli init -s http://localhost:3000 -l admin
     export RPORT_SERVER_URL=http://localhost:3000
     rportcli client list
 

After the config initialisation, Rportcli will check the provided options by calling the rport [status API](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/cloudradar-monitoring/rport/master/api-doc.yml#/default/get_status).


## Cli

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
    <td>initial lifetime of an interactive command session in seconds. Max value is 90 days</td>
    <td>10(minutes) * 60</td>
    <td>SESSION_VALIDITY_SECONDS=1800 rportcli command -i</td>
    </tr>
    <tr>
    <td>RPORT_USER</td>
    <td>basic auth login to access rport server</td>
    <td></td>
    <td>RPORT_USER=admin rportcli client list</td>
    </tr>
    <tr>
    <td>RPORT_PASSWORD</td>
    <td>basic auth password to access rport server</td>
    <td></td>
    <td>RPORT_PASSWORD=foobaz rportcli client list</td>
    </tr>
    <tr>
    <td>RPORT_SERVER_URL</td>
    <td>address of rport server</td>
    <td>http://localhost:3000</td>
    <td>RPORT_SERVER_URL=http://localhost:3000 rportcli client list</td>
    </tr>
    <tr>
    <td>TOTP_SECRET_APP_NAME</td>
    <td>account name visible in Authenticator app</td>
    <td>rport</td>
    <td>TOTP_SECRET_APP_NAME=rport rportcli init</td>
    </tr>
</table>
