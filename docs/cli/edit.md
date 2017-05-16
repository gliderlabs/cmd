---
date: 2017-01-31T18:00:00-06:00
title: edit
menu: cli
type: cli
weight: 50
---
#### Edits a command

```sh
$ cat ./script | ssh alpha.cmd.io :edit <name> -
```

`:edit` will rebuild the command `<name>` with the script passed via STDIN.
The builtin requires the `-` second argument to inform it to read from
STDIN, as future versions may introduce an interactive mode.

See [:create](../create/) for the expected script format.
