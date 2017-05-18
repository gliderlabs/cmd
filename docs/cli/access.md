---
date: 2017-01-31T18:00:00-06:00
title: access
menu: cli
type: cli
weight: 70
---
##### Manages command access

```sh
$ ssh alpha.cmd.io :access <name> [<subcommand>]
```

`:access` allows you to control who can run your command `<name>`. The
builtin has subcommands for granting and revoking access to users or access tokens.

If no subcommand is provided, it will list command access by default.

## Subcommands

### ls

##### Lists command access

```sh
$ ssh alpha.cmd.io :access <name> ls
```

The `ls` subcommand will display what users and access tokens have access to
run the command `<name>`. This is the default subcommand to `:access`.

### grant

##### Grants command access to a subject

```sh
$ ssh alpha.cmd.io :access <name> grant <subject>...
```

The `grant` subcommand will allow a subject, either a username or access token,
to run the command `<name>` in your namespace. The command will be configured with the same environment, but they will not have the ability to view or change environment
or access to the command. For this, see [:admins](../admins/).

More than one subject can be provided as extra arguments.

### revoke

##### Revokes command access from a subject

```sh
$ ssh alpha.cmd.io :access <name> revoke <subject>...
```

The `revoke` subcommand will revoke access from subject, either a username or access token, to run the command `<name>` in your namespace.

More than one subject can be provided as extra arguments.
