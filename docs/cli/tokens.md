---
date: 2017-01-31T18:00:00-06:00
title: tokens
menu: cli
type: cli
weight: 90
---
#### Manages access tokens

```sh
$ ssh alpha.cmd.io :tokens [<subcommand>]
```

`:tokens` provides management of access tokens that can be used to provide access to apps and users without a user account. The
builtin has subcommands for creating and deleting tokens.

By default, if no subcommand is provided, it will list access tokens.

## Subcommands

### ls

Lists access tokens

```sh
$ ssh alpha.cmd.io :tokens ls
```

The `ls` subcommand will display existing access tokens. This is the default subcommand to `:tokens`.

### new

Creates a new access token

```sh
$ ssh alpha.cmd.io :tokens new
```

The `new` subcommand will create and display a new access token that can be used
with [:access](../access/).

### rm

Deletes an access token

```sh
$ ssh alpha.cmd.io :tokens rm <token>
```

The `rm` subcommand will revoke all access for and delete the token `<token>`.
