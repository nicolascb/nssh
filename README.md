# nssh - An easy way to manage SSH connections in Linux

Command line tool for manage ~/.ssh/config written in Go.

Inspired by [storm](https://github.com/emre/storm) project.

![](_images/nssh.gif)

# Index

- [Install](#install)
- [Commands](#commands)
- [Features](#features)
- [Donate](#donate)

## Install

### With go get

```
go get -u github.com/nicolascb/nssh
```

### Manual

[Download binary from releases](https://github.com/nicolascb/nssh/releases)

### Compilation instructions

```
$ git clone git@github.com:nicolascb/nssh.git
$ cd nssh/
$ make dep
$ make build
```

Binary location ./bin

## Commands

```
USAGE:
   nssh [global options] command [command options] [arguments...]

VERSION:
   1.0.2

DESCRIPTION:
   An easy way to manage SSH config, see more: github.com/nicolascb/nssh

AUTHOR:
   Nicolas Barbosa <ndevbarbosa@gmail.com>

COMMANDS:
     add         Add a new SSH alias to ~/.ssh/config
     backup      Backup SSH alias
     del         Delete SSH by alias name
     edit        Edit SSH alias by name
     export-csv  Export SSH alias to csv file (Only name, user, hostname and port)
     list        List SSH alias
     search      Search SSH alias by given search text
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Features

- [Add](#add)
- [Edit](#edit)
- [List](#list)
- [Delete](#delete)
- [Search](#search)
- [Backup](#backup)
- [Export to CSV](#export-to-csv)

## Add

```
NAME:
   add - Add a new SSH alias to ~/.ssh/config

USAGE:
   nssh add [-h] [--key SSH_KEY_FILE] name user@host:port -o option=value -o option2=value2

OPTIONS:
   --key value  SSH Key
   -o value     Option config
```

User and port is optional.

Use `-o` to set custom options.

Examples:

```
nssh add prod 192.168.0.99
nssh add dev dev@10.16.1.100 -o connecttimeout=60
nssh add nicolascb root@nicolascb.com.br:5122 --key ~/mykey -o connecttimeout=60 -o loglevel=info
```

## Edit

```
NAME:
   edit - Edit SSH alias by name

USAGE:
   nssh edit [-h] [--key SSH_KEY_FILE] name user@host:port -o option=value -o option2=value2 -p

OPTIONS:
   --key value  SSH Key
   -r value     Rename host
   -o value     Option config
   -p           Preserve options
   -f           Force edit, don't ask to preserve another options
```

Use `-o` to set custom options.

Use `-r` to rename host.

**IMPORTANT**

Use `-p` to preserve another options.

Use `-f` to force edit and don't ask to preserve another options.

Examples:

```
# Only rename host and preserve another options
nssh edit prod -r prod_cloud -p

## Set connectimeout and a new hostname and not preserve options
nssh edit dev -o connectimeout=60 -o hostname=10.10.2.2 -f

## Update
nssh edit nicolascb nicolas@nicolascb.com.br:22 --key ~/mykey2 -p
```

## List

```
NAME:
   list - List SSH alias

USAGE:
   nssh list [-h]
```

Examples:

```
nssh list
```

## Delete

```
NAME:
   del - Delete SSH by alias name

USAGE:
   nssh del [-h] name
```

Examples:

```
nssh del dev_teste
```

## Search

```
NAME:
   search - Search SSH alias by given search text

USAGE:
   nssh search [-h] text
```

Examples:

```
nssh search nicolas
```

## Backup

```
NAME:
   backup - Backup SSH alias

USAGE:
   nssh backup [-h] -file=/tmp/backup.sshconfig

OPTIONS:
   --file value  Output backup file
```

Examples:

```
nssh backup --file ~/nicolas.backup
```

## Export to CSV

```
NAME:
   export-csv - Export SSH alias to csv file (Only name, user, hostname and port)

USAGE:
   nssh export [-h] -file /tmp/sshconfig.csv

OPTIONS:
   --file value  Output CSV file
```

Examples:

```
nssh export-csv --file ~/nicolas.backup
```

## Donate

```
XMR: 48yKKJWPmJq1MW3BFAcyor6vD8RHT85EkXV6D1D7xGPjJksLmfoF7AcNFGPFYbAiYk999Pga6NCNQMZT6uaYqvPvNeSCF8t
```
