---
date: 2017-01-31T18:00:00-06:00
title: Shell commands in the cloud
weight: 10
---

Cmd is a service that runs commands and scripts in the cloud. The primary
interface to Cmd is SSH, meaning you can use cloud commands from any host
with an SSH client. A couple of use cases for Cmd include:

 * **CLI utilities as a service.** Example: bring `jq` wherever you go, without installing it.
 * **Avoid installing big CLI programs.** Example: `latex` is huge and difficult to build but used infrequently.
 * **Scripts you can use from anywhere.** Trigger and orchestrate other systems using any language.
 * **Share and control access to automation.** Build tools for your team or Rickroll your friends.

A very common use case is for sysadmins and operators to build cloud scripts to deploy an app or website. Whether you use Ansible or Terraform or any other command line tool, you can wrap it all up into a cloud command and share it with your team. They can run it from the command line without installing anything, and without needing to have security credentials the command is configured with.

The Run API increases the use cases by letting you easily write scripts in any language and use them as webhook handlers.

If you want to jump in, get started with the [Quickstart](/quickstart/) guide!

## Command Environment

Cloud commands can run any x86 64-bit program in a Linux Docker container.
However, Cmd is designed for commands and scripts rather than long-running
daemons. In fact, commands can only run for a limited time before they timeout,
so Cmd is also not suitable for replacing your shell. Kudos for such a clever
idea, though.

Commands are unable to listen on addressable ports. Another important constraint
of commands is that they are stateless. Besides configured environment
variables, nothing persists across command runs. If the command makes a file, it
will not be there the next time you run it. Utilities that expect files on the
filesystem will need to be wrapped so they can receive the file(s) via STDIN.
