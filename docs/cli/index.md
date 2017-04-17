---
date: 2017-01-31T18:00:00-06:00
title: Overview
menu: cli
type: cli
weight: 10
---

## Using the CLI via SSH

There is no Cmd-specific client. Instead, you use an SSH client to interact with Cmd. Most of our documentation assumes OpenSSH, which is available on all Linux and MacOS systems.

Typically using Cmd via SSH looks like this in the terminal:

```sh
$ ssh <username>@alpha.cmd.io <command>
```

This can be shortened a number of ways. The most common is eliminating the need
for specifying the username by setting it in your SSH configuration. You can also go so far as to create an alias:

```sh
$ alias cmd="ssh <username>@alpha.cmd.io"
$ cmd <command>
```

## Builtin Commands

The point of Cmd is to run commands you create, but there are builtin
commands to manage and configure your commands. These commands are prefixed with
a colon.

 |
--- | ---
[:access](./access) &nbsp;|&nbsp; Manage command access
[:admins](./admins) &nbsp;|&nbsp; Manage command admins
[:create](./create) &nbsp;|&nbsp; Create a command
[:delete](./delete) &nbsp;|&nbsp; Delete a command
[:edit](./edit)     &nbsp;|&nbsp; Edit a command
[:env](./env)       &nbsp;|&nbsp; Manage command environment
[:ls](./ls)         &nbsp;|&nbsp; List available commands
[:tokens](./tokens) &nbsp;|&nbsp; Manage access tokens
:help               &nbsp;|&nbsp; Help about any command
