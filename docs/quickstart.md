---
date: 2017-01-31T18:00:00-06:00
title: Quickstart
weight: 1
---

{{< admonition title="Alpha users" type="note" >}}
You also have access to the [source code](https://github.com/gliderlabs/cmd) and our [development wiki](https://github.com/gliderlabs/cmd/wiki)! A great way
to help out is to take notes while playing with Cmd.io: just make a page for
yourself on our [UserNotes](https://github.com/gliderlabs/cmd/wiki/UserNotes) wiki page.
{{</ admonition >}}

With Cmd we can make shell scripts that live and run in the cloud. In this quick guide, we'll make a "Hello, world" script in Bash as a cloud command.

First, we'll need to make sure we've set up our GitHub user and SSH key to access Cmd over SSH. Take a quick look at [Using the CLI via SSH](/cli/) if you need to add your key to GitHub. Cmd uses your GitHub and SSH key to authenticate.

If it's your first time using Cmd and you haven't logged into the web-based [Console](https://alpha.cmd.io/console/), it may prompt you to do that first.

After that you should be able to just SSH to alpha.cmd.io to get a basic usage response, something like this:

```text
$ ssh alpha.cmd.io
⌘ Cmd by Glider Labs
  version: 473573c2

► Usage
  ssh alpha.cmd.io [ command | builtin ]

☰ Commands

☰ Builtins
  :access      Manage command access
  :admins      Manage command admins
  :create      Create a command
  :delete      Delete a command
  :edit        Edit a command
  :env         Manage command environment
  :help        Help about any command
  :ls          List available commands
  :tokens      Manage access tokens

⚑ Flags
  -h, --help   help for cmd

Use "ssh alpha.cmd.io [builtin] --help" to learn more about a builtin.
Connection to alpha.cmd.io closed.
```

No commands are listed, but we can make one using the builtin commands.

### Making a command

We can create a command with `:create`. First, we have to create a Cmd script locally to feed into it. Let's make a `hello.cmd` file:

```text
#!cmd alpine bash
#!/bin/bash
echo "Hello, ${1:-world}!"
```

You'll notice this looks like a standard shell script with the addition of an extra shebang line. This tells Cmd.io how to build the environment for the command. The first argument `alpine` represents Alpine Linux, currently the only supported Linux distro. Any following arguments are packages to install. You can search for packages [based on name](http://pkgs.alpinelinux.org/packages) or [based on contents](http://pkgs.alpinelinux.org/contents).

Let's create the command from the script:

```text
$ cat hello.cmd | ssh alpha.cmd.io :create hello
Creating command... done
```

Now we can run it:

```text
$ ssh alpha.cmd.io hello
Hello, world!
```

We made it use `world` if there is no argument, but we can also provide one when we run it:

```text
$ ssh alpha.cmd.io hello everybody
Hello, everybody!
```

That's it!

### Making a command for an existing utility

If you want to build a command based on an existing Alpine package, you can just install it and then specify it as the "interpreter" (in this case it's acting as an "entrypoint"). No need to specify a full path, but you can if you want. Here's a script to make a [jq](https://stedolan.github.io/jq/) command:

```text
#!cmd alpine jq
#!jq
```

Create the command as usual. Let's say that the above was called `jq.cmd`:

```text
$ cat jq.cmd | ssh alpha.cmd.io :create jq
```

Now you have a cloud `jq` you can use from anywhere. Here we'll curl the GitHub API for my user and get my name from the JSON by piping into our cloud `jq`:

```text
$ curl -s https://api.github.com/users/progrium \
  | ssh alpha.cmd.io jq .name
"Jeff Lindsay"
```

{{< admonition title="Exploring Alpine" type="note" >}}
If you want to play with Alpine Linux, the easiest way to interactively experiment with it is using [Docker](https://www.docker.com/). The official `alpine` image is the same image used by Cmd. Here's how you'd jump into an Alpine shell in Docker:

```
$ docker run -it alpine sh
```
{{</ admonition >}}

### What next?

See what else you can do by exploring the rest of the docs! If you need help, [join our Slack](http://slack.gliderlabs.com/) and say hi in the `#cmd` channel.
