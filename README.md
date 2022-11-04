# AKILT

[![Join the chat at https://gitter.im/_AKILT/community](https://badges.gitter.im/_AKILT/community.svg)](https://gitter.im/_AKILT/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
![Code Grade](https://api.codiga.io/project/34798/status/svg)
![top language](https://img.shields.io/github/languages/top/Xart3mis/AKILT)
![Lines of code](https://img.shields.io/tokei/lines/github/Xart3mis/AKILT)

#### AKILT _(pronounced ay kilt)_ is an [undetectable](https://www.virustotal.com/gui/file/42673f19cf40d15b1f38235b5bb952c36647c42c64b209c353cd7978a1ddb555/detection) windows \*botnet _??_\* written in golang with a cross-platform C&C Server

**AKILT** aims to help security enthusiasts and malware analysts better understand how botnets work by providing an open source example of an advanced botnet.

### Setup
#### You can download one of the prebuilt binaries to test it.
#### Or build from source by following these steps:
###### Install the dependencies:
  - a gcc compiler [(download TDM-GCC from here)](https://jmeubank.github.io/tdm-gcc/download/)  
  - make (download [make](https://community.chocolatey.org/packages/make) using [chocolatey](https://chocolatey.org/install))
  - Go [(download the golang installer from here)](https://go.dev/dl/)
  
  
###### After installing the dependencies above. Run these commands.
```bash
git clone github.com/Xart3mis/AKILT
```  

```bash
cd AKILT
```  

```bash
make
```  

**NOTE: Compiling the client only works on windows 64-bit**  
###### If you'd like to compile the client or server seperately: 
instead of running `make` directly run `make client` or `make server`  


##### the compiled binaries will appear in `Client/bin/` and `Server/bin/`
### Server

![help](https://github.com/Xart3mis/AKILT/blob/master/help.gif)

## Features

**[X]** Capture Client Screen  
**[X]** Display Dialog on Client PC  
**[X]** Take Picture from Client webcam  
**[X]** DDOS (Slowloris, httpflood, Udpflood)  
**[X]** Remote Command Execution  
**[X]** Display text on client screen  
**[WIP]** Download files to client  
**[WIP]** Upload files from client  
**[WIP]** Play Audio on client pc  
**[X]** Keylogger  
**[WIP]** Spreading mechanism _i am yet to figure it out_  
**[WIP]** Persistence  
**[WIP]** [Bot and C2 builder](https://github.com/Xart3mis/akiltbuilder/)  
**[WIP]** Reverse Shell  
**[WIP]** Hardcoded login system  
**[WIP]** Get System Information  
**[WIP]** Retrieve user password hashes  

###### feel free to contribute. _please do_
