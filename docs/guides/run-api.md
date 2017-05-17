---
date: 2017-01-31T18:00:00-06:00
title: Using the Run API
menu: guides
weight: 50
---

What if you could turn a shell script into a web API in seconds?

<img src="https://cdn-images-1.medium.com/max/800/1*ITLxMPLRoXrEBsd_VxHw8w.gif" />

Every Cmd.io command you make can be run via SSH *and* our Run API. All you have to do is write the script and create the command as usual, and you get this right out of the box.

The Run API exposes a secure HTTP endpoint to execute and return the output of your command, with an option to stream the response in real-time. The same endpoint can also be upgraded to a WebSocket connection in order to stream the output into a browser or to a WebSocket client.

Running a command this way is restricted to those with Access Tokens, which you create and assign to your commands. You can revoke the token for a command&mdash;or delete the token entirely to prevent access&mdash;at any time.

Combined with environment configuration, you can start exposing secure APIs for capabilities that would otherwise only be reasonable to run from a workstation shell.

Let’s run through an example. We’ll take `git log`, showing us changes to a repository since a recent point in time, and expose it as a command with an API endpoint. The period of time to review will be a human-friendly argument of the command. Here is our Cmd script consisting of two lines of Bash after two shebang lines:

```sh
#!cmd alpine bash git
#!/bin/bash
git clone --depth "${DEPTH:-10}" "${REPO?}" . &> /dev/null
git log --no-merges --raw --since="$*"
```

This otherwise normal shell script with a special Cmd shebang does a few things. First, it sets up an Alpine Linux environment with `bash` and `git` to run in, then using Bash it clones a repository 10 commits back but hiding the output, and then it runs `git log` using an argument specifying how far back in time to look.

You might notice the clone depth of 10 can optionally be configured with `$DEPTH`. We’ll leave that, but we *will* need to configure `$REPO` before we run the command, or it will run and immediately exit in error.

We’ll pipe this script into Cmd.io as a new command called `recent`:

```sh
$ cat ./script | ssh alpha.cmd.io :create recent
Creating command... done
```

Now we need to set the `$REPO` environment variable to a repository URL. Let’s use my [go-basher](https://github.com/progrium/go-basher) project. This is actually one command written on two lines using `\` for readability:

```sh
$ ssh alpha.cmd.io :env recent set \
    REPO=https://github.com/progrium/go-basher
Setting REPO on recent... done
```

That was it! We can now run `recent` via SSH, asking it to show changes from the last 2 weeks:

```sh
$ ssh alpha.cmd.io recent "4 weeks ago"
commit 1bb28e5ef958827fc687a61f2f5bba71dca62c79
Author: Leon van Kammen <leonvankammen@gmail.com>
Date:   Wed Feb 1 22:33:44 2017 +0100

    Update README.md
    * removed jsonPointer from Application()-example (broke the app)
    * replaced long sentences with short instructions
:100644 100644 347ef69… 6933945… M README.md
```

The Run API uses a different method of authentication than SSH does; we need to create an access token and add it to our `recent` command’s access control list:

```sh
$ ssh alpha.cmd.io :tokens new
Token created: a17671fb-b2f1–4286–89de
$ ssh alpha.cmd.io :access recent grant a17671fb-b2f1–4286–89de
Granting a17671fb-b2f1–4286–89de access to recent... done
```

Now we can access the command via HTTP using `curl`, again using a one line command shown on two lines:

```sh
$ curl -u "a17671fb-b2f1–4286–89de:" \
    https://alpha.cmd.io/run/progrium/recent/2+weeks+ago
commit 1bb28e5ef958827fc687a61f2f5bba71dca62c79
Author: Leon van Kammen <leonvankammen@gmail.com>
Date:   Wed Feb 1 22:33:44 2017 +0100

    Update README.md
    * removed jsonPointer from Application()-example (broke the app)
    * replaced long sentences with short instructions
:100644 100644 347ef69… 6933945… M README.md
```

Ta-da! Note that we pass the token as the username with no password. We also have the option to pass the token via the `access_token` query parameter. Possessing a token is all you need to execute the commands associated with it, so be careful with them!

You might be thinking, “Cool, but what good is an API that returns mostly unstructured text?” Well, that’s up to you as the command author! To see a version of this command that returns JSON, as well as a WebSocket demo, check out this quick video:

{{< youtube ndlxB1KmQgc >}}
