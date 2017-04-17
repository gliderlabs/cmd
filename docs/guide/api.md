---
date: 2017-01-31T18:00:00-06:00
title: Run API
weight: 50
---

You can run commands without SSH as well. The Run API gives you an HTTP and
WebSocket API.

## Authentication

Since we can't use your SSH key easily for HTTP authentication, it's best to use
tokens. You can also use your GitHub username and password via Basic Auth in
development, as we ensure the connection uses TLS. However this should never be
used in production integrations.

## Tokens

The builtin command `cmd-tokens` has a number of commands for managing access
tokens that you can give access to like a regular user.

```
$ ssh alpha.cmd.io cmd-tokens
Usage:
  ssh <user>@cmd.io cmd-tokens [command]

Available Commands:
  ls          List tokens
  new         Create a token
  rm          Delete a token

Additional help topics:

Use "[command] --help" for help about a meta command.
```

Creating a token will generate a UUID based token:

```
$ ssh alpha.cmd.io cmd-tokens new
Token created: 84e67956-0ec3-4e49-a588-8788719dea95
```

Now we can add this to a command:

```
$ ssh alpha.cmd.io welcome:access add 84e67956-0ec3-4e49-a588-8788719dea95
Access granted.
```

Be careful with these tokens. This token now lets anybody that has it run your
command over SSH as well as the Run API. When using over SSH, the token is
used as the username and no key or password is required.

## HTTP API

Now that you have a token you can use the run API endpoint:

```
https://alpha.cmd.io/run/<user>/<command>
```

The `user` path part is the owner of the command, in this case, your username.
Performing a `GET` request with the token will capture all the output (combined
`STDOUT` and `STDERR`) and return it in the response:

```
$ curl https://84e67956-0ec3-4e49-a588-8788719dea95@alpha.cmd.io/run/progrium/welcome

Hello! This is a cmd.io command. All it does is display this message.
However, cmd.io commands can do lots more. They can pretty much do
anything you can do in a Docker container, except for long-running
processes like daemons.

You can install cmd.io commands from a number of sources, including
anything off Docker Hub. Once you have a command installed, you can
configure it and share access to it. Anybody that has access to your
command can run it from anywhere they have an SSH client.

```

With the query parameter `?stream=1`, the output will be streamed in chunked
transfer encoding. There is no support for `STDIN` or separate `STDOUT` and `STDERR`
output via this API.

## WebSocket API

The same endpoint can be upgraded to a WebSocket connection. Output will
be streamed back as text.
