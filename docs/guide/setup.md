---
date: 2017-01-31T18:00:00-06:00
title: Setup
weight: 10
---

## User Notes

We'd love it if you took notes about your experience getting started and then
using Cmd.io. You can then share on
[UserNotes](https://github.com/gliderlabs/cmd/wiki/UserNotes) and we can use
them as living documents to improve the experience. Also, if something is wrong
on this page or any other, feel free to edit this wiki!

## Setup

> This document is currently for the alpha release channel (alpha.cmd.io) ...
any accidental references to the cmd.io domain should be read as alpha.cmd.io

Cmd.io uses your GitHub user for authentication. It also relies on the public
keys stored with your GitHub account. If you haven't uploaded a public key to
GitHub, you can easily add one in Settings under [SSH and GPG
keys](https://github.com/settings/keys).

Now you can connect to Cmd.io over SSH:

```
$ ssh progrium@alpha.cmd.io
Usage:
  ssh <user>@cmd.io [command]

Available Commands:
  cmd-add             Install a command
  cmd-ls              List installed commands
  cmd-rm              Uninstall a command

Use "[command] --help" for help about a meta command.
```

If you're using a local user with the same username as your GitHub username, SSH
will use it by default. If not, you can also avoid specifying a username with
some configuration added to your `~/.ssh/config` file:

```
Host alpha.cmd.io
    User progrium
```
