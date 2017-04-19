---
date: 2017-01-31T18:00:00-06:00
title: Builtin Commands
weight: 20
---

Cmd.io is about running and managing user commands. Builtin commands let us add and
remove user commands. They're prefixed with `cmd-`:

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

## Installing a command

We can install a command with `cmd-add`. We can see how with `--help`:

```
$ ssh alpha.cmd.io cmd-add --help
Install a command

Usage:
  ssh <user>@cmd.io cmd-add <name> <source>

```

Although this help could be more helpful, it does at least tell us the arguments
to `cmd-add`. The first, `name`, is a name you choose for the command. The
second, `source`, is where to get the command. Right now, the only supported
sources are public Docker registries like Docker Hub. Let's use this mysterious
Docker image, [progrium/welcome](https://hub.docker.com/r/progrium/welcome/):

```
$ ssh alpha.cmd.io cmd-add demo progrium/welcome
Command installed
$ ssh alpha.cmd.io
Usage:
  ssh <user>@cmd.io [command]

Available Commands:
  cmd-add             Install a command
  cmd-ls              List installed commands
  cmd-rm              Uninstall a command
  demo

Use "[command] --help" for help about a meta command.

```

Although it lacks a description, the command `demo` is now available. Run it:

```
$ ssh alpha.cmd.io demo

Hello! This is a cmd.io command. All it does is display this message.
However, cmd.io commands can do lots more. They can pretty much do
anything you can do in a Docker container, except for long-running
processes like daemons.

You can install cmd.io commands from a number of sources, including
anything off Docker Hub. Once you have a command installed, you can
configure it and share access to it. Anybody that has access to your
command can run it from anywhere they have an SSH client.

```

Ignoring the fact this message is lying to you (Docker Hub is currently the
*only* source for commands), you've run your first Cmd.io command!

The container images you install don't need to be specific to Cmd.io. Pretty
much any CLI tool in a container can be installed. Here's a
[netcat](https://hub.docker.com/r/gophernet/netcat/) container we can use to
show that `cmd.io` port `22` is open:

```
$ ssh alpha.cmd.io cmd-add nc gophernet/netcat
Command installed
$ ssh alpha.cmd.io nc -z -v cmd.io 22
cmd.io (159.203.159.60:22) open
```

## Removing commands

This is pretty straightforward:

```
$ ssh alpha.cmd.io cmd-rm nc
Command uninstalled
```
