# SQLcl Docker Extension

sqlcl extension for Docker Desktop

## Manual Installation

Until this extension is ready at Docker Extension Hub you can install just by executing:

```bash
$ docker extension install mochoa/sqlcl-docker-extension:25.1.0
Extensions can install binaries, invoke commands and access files on your machine. 
Are you sure you want to continue? [y/N] y
Installing new extension "mochoa/sqlcl-docker-extension:25.1.0"
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

Note that an Oracle server running at Docker Desktop or externals is not localhost, the IP for OracleXE running at Docker Desktop is available at the menu, Settings -> Resources -> Network -> Docker subnet, in my case is 192.168.65.0/24 so an internal IP for reaching OracleXE container running at Docker Desktop will be 192.168.65.2, also there is an internal DNS name that resolve above IP named **host.docker.internal**.

![Docker Desktop Subnet](https://miro.medium.com/max/700/0*m4e0OEQprx_GgUA7)

Let see an example of OracleXE started using OracleXE Docker Desktop Extension:

SQL> connect scott/tiger@host.docker.internal:1521/xepdb1

which means for a sqlcl Add New Server:

- Hostname/address: host.docker.internal
- Port: 1521
- PDB: xepdb1
- Username: scott
- Password: tiger

the extension have a persistent history of command, if you choose exit and start again SQLcl Docker Extension just use cursor up and down to navigate on history entries.

## Upload/Download files

If you want to import/export large files from/to Oracle RDBMs using some of the functionality of SQLcl such as dp (Data Pump) or load command is necessary to provide that files to the extension persistent volume, you can check it by using Volumes menu:

![Docker Desktop Volumes](https://miro.medium.com/max/1400/0*DZdlhN995x75t5a9)

Persistent volume is **mochoa_sqlcl-docker-extension-desktop-extension_sqlcl_home**.

A sample usage taken from Jeff Smith blog but using SQLcl Docker Extension is:

```sql
SQL> set feedback off
SQL> spool objects.csv
SQL> SELECT * FROM all_objects fetch FIRST 100 ROWS ONLY;
SQL> spool off
```

copy exported data to you local file system and upload again in another directory, run this in a command line shell:

```bash
docker cp mochoa_sqlcl-docker-extension-desktop-extension-service:/home/sqlcl/objects.csv .
docker cp objects.csv mochoa_sqlcl-docker-extension-desktop-extension-service:/tmp
```

finally import again using SQLcl LOAD command:

```sql
SQL> CREATE TABLE demo_load AS SELECT * FROM all_objects WHERE 1=2;
Table DEMO_LOAD created.
SQL> load demo_load /tmp/objects.csv
Load data into table SCOTT.DEMO_LOAD
csv
column_names on
delimiter ,
enclosures ""
encoding UTF8
row_limit off
row_terminator default
skip_rows 0
skip_after_names
#INFO Number of rows processed: 100
#INFO Number of rows in error: 0
#INFO No rows committed
SUCCESS: Processed without errors
SQL> set sqlformat default
SQL> select count(*) from demo_load;
  COUNT(*)
----------
       100
```

## Connect to Oracle Cloud Autonomous DB

SQLcl Desktop Extension is able to connect to your autonomous DB, you have to download first your Wallet file from OCI cloud console as is described into [Download Database Connection Information](https://docs.oracle.com/en-us/iaas/autonomous-database-shared/doc/connect-download-wallet.html#GUID-B06202D2-0597-41AA-9481-3B174F75D4B1), you will see a tab like:

![ATP DB Connection Wallet Download Tab](https://miro.medium.com/max/1352/1*wendOYTkXqtViKsxHy70BQ.png)

Once you have a Zip file upload it to SQLcl persistent volume using docker cp command:

```bash
docker cp Wallet_DBparquet.zip mochoa_sqlcl-docker-extension-desktop-extension-service:/home/sqlcl
```

finally use you Wallet inside SQLcl Docker Extension using:

```sql
SQL> set cloudconfig /home/sqlcl/Wallet_DBparquet.zip
SQL> connect admin/MyStrongSecretPwd@dbparquet_high
SQL> desc demo_tab
SQL> select count(*) from demo_tab;
```

## Using Liquibase with Scott user

By default scott user created during OracleXE installation only have connect,resource roles; to use latest Liquibase distributon included on SQLcl you have to grant another extra role:

```sql
SQL> grant CREATE VIEW to scott;
```

## Uninstall

To uninstall the extension just execute:

```bash
$ docker extension uninstall mochoa/sqlcl-docker-extension:25.1.0
Extension "Oracle SQLcl client tool" uninstalled successfully
```

## Source Code

As usual the code of this extension is at [GitHub](https://github.com/marcelo-ochoa/sqlcl-docker-extension), feel free to suggest changes and make contributions, note that I am a beginner developer of React and TypeScript so contributions to make this UI better are welcome.
