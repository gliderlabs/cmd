---
date: 2017-01-31T18:00:00-06:00
title: env
menu: cli
type: cli
weight: 60
---
##### Manages command environment

```sh
$ ssh alpha.cmd.io :env <name> [<subcommand>]
```

`:env` provides management of environment variables for your command `<name>`. The
builtin has subcommands for setting and unsetting variables.

By default, if no subcommand is provided, it will list environment variables.

## Subcommands

### ls

##### Lists command environment variables

```sh
$ ssh alpha.cmd.io :env <name> ls
```

The `ls` subcommand will display environment variables configured for
the command `<name>`. This is the default subcommand to `:env`.

### set

##### Sets an environment variable

```sh
$ ssh alpha.cmd.io :env <name> set <key=value>...
```

The `set` subcommand will set one or more variables for the command `<name>` in your namespace. Variables are set with key value pairs like `FOO=bar`.

More than one variable can be set with extra arguments.

### unset

##### Unsets an environment variable

```sh
$ ssh alpha.cmd.io :env <name> unset <key>...
```

The `unset` subcommand will delete one or more environment variables for the command `<name>` in your namespace.

More than one variable key can be provided as extra arguments.
