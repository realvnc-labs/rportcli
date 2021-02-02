# Rport CLI (v1)
Rport CLI is a tool to help you managing [rport API](https://github.com/cloudradar-monitoring/rport) directly from your terminal.

## Installation

### As a compiled binary

Jump to [our release page](https://github.com/cloudradar-monitoring/tacoscript/releases/tag/latest) and download a binary for your host OS. Don't forget to download a corresponding md5 file as well.


        # On MacOS
        wget https://github.com/cloudradar-monitoring/tacoscript/releases/download/latest/tacoscript-latest-darwin-amd64.tar.gz
        
        # On linux
        wget https://github.com/cloudradar-monitoring/tacoscript/releases/download/latest/tacoscript-latest-linux-amd64.tar.gz
        
        # On Windows
        Just download https://github.com/cloudradar-monitoring/tacoscript/releases/download/latest/tacoscript-latest-windows-amd64.zip
        Also download https://github.com/cloudradar-monitoring/tacoscript/releases/download/latest/tacoscript-latest-windows-amd64.zip.md5
     
     
Verify the checksum:

    
        #On MacOS
        curl -Ls https://github.com/cloudradar-monitoring/tacoscript/releases/download/latest/tacoscript-latest-darwin-amd64.tar.gz.md5 | sed 's:$: tacoscript-latest-darwin-amd64.tar.gz:' | md5sum -c
        
        #On linux
         curl -Ls https://github.com/cloudradar-monitoring/tacoscript/releases/download/latest/tacoscript-latest-linux-386.tar.gz.md5 | sed 's:$: tacoscript-latest-linux-amd64.tar.gz:' | md5sum -c
         
        #On Windows assuming you're in the directory with the donwloaded file
         CertUtil -hashfile tacoscript-latest-linux-amd64.tar.gz MD5
        
        #The output will be :
         MD5 hash of tacoscript-0.0.4pre-windows-amd64.zip:
         7103fcda170a54fa39cf92fe816833d1
         CertUtil: -hashfile command completed successfully.
        
        #Compare the command output to the contents of file tacoscript-latest-windows-amd64.zip.md5 they should match
  
  
    
_Note: if the checksums didn't match please don't continue the installation!_

Unpack and install the tacoscript binary on your host machine

    
        #On linux/MacOS
        tar -xzvOf tacoscript-0.0.4pre-darwin-amd64.tar.gz >> /usr/local/bin/tacoscript
        chmod +x /usr/local/bin/tacoscript
    

For Windows

Extract file contents:
![C:\Downloads](docs/Extract.png?raw=true "Extract")

Create a `Tacoscript` folder in `C:\Program Files`

![C:\Program Files](docs/ProgramFiles.png?raw=true "ProgramFiles")

Copy the tacoscript.exe binary to `C:\Program Files\Tacoscript`
![C:\Program Files\Tacoscript\tacoscript.exe](docs/ProgramFilesWithTacoscript.png?raw=true "ProgramFilesWithTacoscript")

Double click on the tacoscript.exe and allow it's execution:

![C:\Program Files\Tacoscript\tacoscript.exe](docs/AllowRun.png?raw=true "AllowRun")

## Install as a go binary:

    go get github.com/cloudradar-monitoring/tacoscript

## Program execution

