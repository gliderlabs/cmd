---
date: 2017-01-31T18:00:00-06:00
title: Overview
weight: 1
---

Cmd.io is a service that remotely runs commands and scripts. The primary
interface to Cmd.io is SSH, meaning you can use Cmd.io commands from any host
with an SSH client. A couple of use cases for Cmd.io include:

 * **CLI utilities as a service.** Example, bring `jq` wherever you go, without installing it.
 * **Avoid installing big CLI programs.** Example, `latex` is huge and difficult to build but used infrequently.
 * **Scripts you can use from anywhere.** Trigger and orchestrate other systems using any language.
 * **Share and control access to automation.** Build tools for your team or Rickroll your friends.

Cmd.io commands can run any x86 64-bit program in an Linux Docker container.
However, it's designed for commands and scripts as opposed to long-running
daemons. In fact, commands can only run for a limited time before they timeout.
So Cmd.io is also not suitable for replacing your shell. Kudos for such a clever
idea, though.

Commands are unable to listen on routable ports, which is not common with most
commands anyway. In the rare cases you need to, such as debugging with netcat,
consider using a tool like [ngrok](https://ngrok.com/).

Another important constraint of commands is that they are stateless. Besides
configured environment variables, nothing persists across command runs. If the
command makes a file, it will not be there the next time you run it. Utilities
that expect files on the filesystem will need to be wrapped so they can receive
the file(s) via STDIN.
