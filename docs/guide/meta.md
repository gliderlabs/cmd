---
date: 2017-01-31T18:00:00-06:00
title: Meta Commands
weight: 40
---

Some of the real power in Cmd.io comes from what you can do with meta commands.
Meta command are commands that operate on a command. They look like
`<command>:<metacommand>`. Let's use the `:help` meta command on our `netpoll`
command from before to see what we can do:

```
$ ssh alpha.cmd.io netpoll:help
Usage:
  ssh <user>@cmd.io netpoll:[command]

Available Commands:
  :access     Manage command access
  :admins     Manage command admins
  :env        Manage command environment

Use "[command] --help" for help about a meta command.
```

You can explore `:access` and `:admins` on your own. In short, access lets you
share a command with others by adding and removing GitHub usernames. They can
run it prefixed with your username and a slash. For example, if I shared
`netpoll` with you, you could run it with `ssh cmd.io progrium/netpoll`. Access
also lets you make a command public, letting any user run it, or make it private
again.

The difference between access and admins is currently that admins have access to
these meta commands just like you. Users with access don't. A useful dynamic
this allows is storing credentials in configuration, and users with just access
can use your commands that use those credentials, but aren't able to see the
credential themselves.

> **However, if there is a way to display the environment in your command, those
credentials will be plainly visible. Use at your own risk.**

### Setting environment variables

Let's look at the `:env` meta command:

```
$ ssh alpha.cmd.io netpoll:env --help
Manage command environment

Usage:
  demo :config [command]

Available Commands:
  set         Manage command environment
  unset       Manage command environment

Use "[command] --help" for help about a meta command.

```

What this help is not telling you is that running `:env` without `--help` will
list current configuration values. However if you do that now, there will be no
output since there is no configuration. Let's set configuration on `netpoll`.

Configuration is exposed to commands as environment variables, so you can
usually just think of `:config` as managing environment variables. If you look
back, our `netpoll` script actually uses a variable if set called `TIMEOUT`,
which defaults to `10`. We can change that value by setting `TIMEOUT`:

```
$ ssh alpha.cmd.io netpoll:env set TIMEOUT=30
Config updated.
$ ssh alpha.cmd.io netpoll:env
TIMEOUT=30
```

Now if you use `netpoll` against a port that's not accepting connections, it's
going to loop 30 times instead of 10. This will apply to any user that has
access to this command.
