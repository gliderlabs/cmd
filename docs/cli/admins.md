---
date: 2017-01-31T18:00:00-06:00
title: admins
menu: cli
type: cli
weight: 80
---
##### Manages command admins

```sh
$ ssh alpha.cmd.io :admins <name> [<subcommand>]
```

`:admins` allows you to control who can administer your command `<name>`. The
builtin has subcommands for granting and revoking command admin to users.

Admins have the ability to view and change environment and access, and edit a
command script.

If no subcommand is provided, it will list command admins by default.

## Subcommands

### ls

##### Lists command admins

```sh
$ ssh alpha.cmd.io :admins <name> ls
```

The `ls` subcommand will display what users have admin access to
the command `<name>`. This is the default subcommand to `:admins`.

### grant

##### Grants command admin to a user

```sh
$ ssh alpha.cmd.io :admins <name> grant <user>...
```

The `grant` subcommand will allow a user to manage the command `<name>` in your namespace.

More than one user can be provided as extra arguments.

### revoke

##### Revokes command admin from a user

```sh
$ ssh alpha.cmd.io :admins <name> revoke <user>...
```

The `revoke` subcommand will revoke admin privilege from a user for the command `<name>` in your namespace.

More than one user can be provided as extra arguments.
