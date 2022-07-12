# Docker Extension

sqlcl extension for Docker Desktop

## Manual Installation

Until this extension is ready at Docker Extension Hub you can install just by executing:

```bash
$ docker extension install mochoa/sqlcl-docker-extension:22.2.0
Extensions can install binaries, invoke commands and access files on your machine. 
Are you sure you want to continue? [y/N] y
Installing new extension "mochoa/sqlcl-docker-extension:22.2.0"
Installing service in Desktop VM...
Setting additional compose attributes
VM service started
Installing Desktop extension UI for tab "sqlcl"...
Extension UI tab "sqlcl" added.
Extension "Oracle SQLcl client tool" installed successfully
```

**Note**: Docker Extension CLI is required to execute above command, follow the instructions at [Extension SDK (Beta) -> Prerequisites](https://docs.docker.com/desktop/extensions-sdk/#prerequisites) page for instructions on how to add it.

## Using SQLcl Docker Extension

Once the extension is installed a new extension is listed at the pane Extension (Beta) of Docker Desktop.

By clicking at SQLcl icon the extension main window will display a progress bar for a few second and finally SQLcl is launched.

![Progress bar indicator](docs/images/screenshot1.png?raw=true)

SQLcl is not logged into the Oracle RDBMS you should log using connect command, for example for scott user on OracleXE is.

![Connect sample](docs/images/screenshot2.png?raw=true)

Note that an Oracle server running at Docker Desktop or externals is not localhost, the IP for OracleXE running at Docker Desktop is available at the menu, Settings -> Resources -> Network -> Docker subnet, in my case is 192.168.65.0/24 so an internal IP for reaching OracleXE container running at Docker Desktop will be 192.168.65.2.

![Docker Desktop Subnet](https://miro.medium.com/max/700/0*m4e0OEQprx_GgUA7)

Let see an example of OracleXE started using OracleXE Docker Desktop Extension:

SQL> connect scott/tiger@192.168.65.2:1521/xepdb1

which means for a sqlcl Add New Server:

- Hostname/address: 192.168.65.2
- Port: 1521
- PDB: xepdb1
- Username: scott
- Password: tiger

the extension have a persistent history of command, if you choose exit and start again SQLcl Docker Extension just use cursor up and down to navigate on history entries.

### Knowns caveats

For some reasons that I don't know if you are using SQLcl Docker Desktop Extension and want to switch to another extension such as Disk Usage by clicking at left side pane Extensions (Beta) -> Disk Usage main windows focus still at SQLcl Extension, as a workaround just click on some left main pane options such as Home or Containers and go back to the desire extension. This problem is visible also when using Logs Explorer and the extension is showing a lot of logs.
