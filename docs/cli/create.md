---
date: 2017-01-31T18:00:00-06:00
title: create
menu: cli
type: cli
weight: 20
---
##### Creates a new command

```sh
$ cat ./script | ssh alpha.cmd.io :create <name> -
```

`:create` will build and add a new command to your namespace with the given
`<name>` using the script passed via STDIN. The `-` argument is required to
inform the command to read from STDIN, as future versions may introduce an
interactive mode.

### Script format

Cmd expects a shell script with at least two shebang lines. The first line
must be a special `#!cmd` line in this format:

```text
#!cmd <base> [<package>...]
```

The first argument `<base>` is either `alpine` or `ubuntu`. The following optional
arguments are packages available in their respective package repositories. Examples:

```text
#!cmd alpine
```

```text
#!cmd alpine bash curl
```

```text
#!cmd ubuntu git build-essential
```

The next line must be a regular shebang line, either defining an interpreter
for the rest of the script or the binary to run for the command. Examples:

```text
#!cmd alpine bash
#!/bin/bash
echo "Hello world"
```

```text
#!cmd alpine curl
#!/usr/bin/curl
```

### Example
This simple example will install a package and use the binary as the interpreter. Our script will be the following:

```text
#!cmd alpine speedtest-cli
#!/usr/bin/speedtest-cli
```

Let's create our `speedtest` command by redirecting a heredoc (an alternative to piping the text file) to the `:create` built-in command:

```sh
$ ssh alpha.cmd.io :create speedtest - <<EOF
#!cmd alpine speedtest-cli
#!/usr/bin/speedtest-cli
EOF
Creating command... done
```

We can then run and pass arguments to the command:

```sh
$ ssh alpha.cmd.io speedtest --simple
Ping: 13.78 ms
Download: 552.79 Mbit/s
Upload: 126.91 Mbit/s
```

The output is a test of the Cmd infrastructure (not your local Internet connection). While not very practical, it demonstrates how we can run simple utilities entirely remotely.
