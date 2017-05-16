---
date: 2017-01-31T18:00:00-06:00
title: Handling git push
menu: guides
weight: 60
draft: true
---


You can use Cmd.io commands over SSH as Git
remote endpoints to implement and react to `git push` from your repositories.
Specifically, you can use one command for Git remotes: `git-receive-pack`.

You see, all that happens when you `git push` to an SSH remote like
`git@github.com:gliderlabs/cmd.git` is it translates to the SSH command `ssh
git@github.com git-receive-pack 'gliderlabs/cmd.git'`. This opens a session to
`github.com` and runs the command `git-receive-pack` with the repo path
argument. The actual `git-receive-pack` command then reads packed Git data over
STDIN and applies it to the bare repository on the filesystem that was provided
as an argument.

What Heroku and later [gitreceive](https://github.com/progrium/gitreceive),
Dokku, Flynn, etc all do in some form effectively is wrap `git-receive-pack` and
install a pre-receive hook into the repository (perhaps created on the fly) that
does some task. Git is designed to display the output of that hook back to the
user during the push, so from that hook script you can `git archive` to get a
tar of what was pushed and deploy or do whatever you want with it.

Currently, you can do this with Cmd.io by creating a command named
`git-receive-pack`, which will handle all Git pushes to Cmd.io that authenticate
with your username. Cmd.io is doing nothing special to make this work, but now
you need a command that will properly handle the push.

Here's an example to get you started. It involves three files: `Dockerfile`,
`git-receive`, and `pre-receive`. Remember the last two need to have `chmod +x`
run on them.

#### Dockerfile
```
FROM alpine:3.4
RUN apk add --update --no-cache git sed bash
COPY ./git-receive /bin/git-receive
COPY ./pre-receive /hooks/
ENTRYPOINT ["/bin/git-receive"]
```
Note that anything your `pre-receive` script is going to use also needs to be
installed in this Dockerfile.

#### git-receive
```
#!/bin/bash
repo="$1"
if [[ "$repo" != /* ]]; then
  repo="/$repo"
fi
git init --quiet --bare "$repo"
cp /hooks/pre-receive $repo/hooks
git-shell -c "git-receive-pack '$repo'"
```
This normalizes the repository argument, creates a bare repo in that location,
installs `pre-receive` as a hook for that repository, and performs the actual
`git-receive-pack`. This will then trigger `pre-receive`.

#### pre-receive
```
#!/bin/bash

main() {
  # reads git push header data into variables
  read old new ref

  # use archive to tarpipe pushed branch files to a working directory
  git archive "$new" | (cd /tmp && tar -xpf -)

  # go to that directory
  cd /tmp

  # do something with the files!
  # exit non-zero and the push will fail.
}

delete-remote-prefix() {
  # this removes "remote: " that git prefixes hook output with client side
  sed -u "s/^/"$'\e[1G'"/"
}

main | delete-remote-prefix
```
This is just a template but from here you could deploy, do builds, run checks,
or something more creative.

#### git remote

Once you've made a command from the above and installed it as
`git-receive-pack`, you can add a remote to the repositories you want to push
from like this:

```
$ git remote add cmd ssh://progrium@cmd.io/repo/path
```

The name of the remote can be anything you like, here it's `cmd`. You can also
drop the username if you were previously. The repo path can be anything as well,
perhaps use it to determine what to do in your pre-receive script.

Unfortunately, `git-receive-pack` is not sharable. Perhaps in the near future
there will be more integrated support in Cmd.io for this pattern.
