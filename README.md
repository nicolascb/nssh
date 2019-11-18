[![Go Report Card](https://goreportcard.com/badge/github.com/nicolascb/nssh)](https://goreportcard.com/report/github.com/nicolascb/nssh)
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

## Commands

```
USAGE:
   nssh [global options] command [command options] [arguments...]

VERSION:
   2.0

DESCRIPTION:
   An easy way to manage SSH config, see more: github.com/nicolascb/nssh

AUTHOR:
   Nicolas Barbosa <ndevbarbosa@gmail.com>

COMMANDS:
     add      Add a new SSH alias to ~/.ssh/config
     del      Delete SSH by alias name
     edit     Edit SSH alias by name
     list     List SSH alias
     search   Search SSH alias by given search text
     help, h  Shows a list of commands or help for one command

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

