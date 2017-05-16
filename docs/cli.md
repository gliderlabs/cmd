---
date: 2017-01-31T18:00:00-06:00
title: CLI Reference
type: cli
weight: 10
---

## Using the CLI via SSH

There is no Cmd-specific client. Instead, you use an SSH client to interact with Cmd. Most of our documentation assumes OpenSSH, which is available on all Linux and MacOS systems.

Typically using Cmd via SSH looks like this in the terminal:

```sh
$ ssh <username>@alpha.cmd.io <command>
```

{{< admonition title="Quick Tip" type="note" >}}
The above can be shortened a number of ways. The most common is editing your `~/.ssh/config` file to add:

```
Host cmd
  HostName alpha.cmd.io
  User <username>
```
Then you can run commands with:
```sh
$ ssh cmd <command>
```
{{< /admonition >}}

### Authentication

Cmd uses your GitHub user for authentication. It also uses the SSH public keys stored with your GitHub account. If you haven't uploaded a public key to GitHub, you can easily add one in Settings under [SSH and GPG keys](https://github.com/settings/keys).

For more information you can read [Connecting to GitHub with SSH](https://help.github.com/articles/connecting-to-github-with-ssh/). If you can connect to GitHub via SSH, you can connect to Cmd.

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
