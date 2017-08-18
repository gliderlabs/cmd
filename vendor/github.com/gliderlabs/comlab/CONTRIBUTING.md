# Please contribute!

Welcome to Comlab and the start of your journey as a contributor. We're glad
you're excited and ready to help out. Here are some ways you can contribute:

 * Talk about the project to others!
 * Submit issues
 * Garden issues
 * Write or fix docs
 * Fix issues
 * Contribute minor improvements
 * Contribute to roadmap items
 * Become a maintainer

## Submitting Issues

No issue or bug is too little to go unrecorded. Here are some tips for good issues:

  * Make reasonable effort to be sure it is not a duplicate of another issue.
  * If your bug has a stack trace, please include it.
  * If there are appropriate labels to apply, please do.
  * If a UI bug, try to include a screenshot or short video reproducing it.
  * If you know the relevant component, prefix the issue title `<component>: `.

## Making Changes

The rest of this document explains how to contribute changes to Comlab.
It assumes you have followed the installation instructions and have written and
tested your code.

It's mostly adapted from the [contribution
guidelines](https://golang.org/doc/contribute.html) for the Go programming
language.

### Discuss your design

The project welcomes submissions but please let [everyone know](
https://slack.gliderlabs.com/) what you're working on to prevent
collisions.

Before writing something new for the project, please file an issue (or claim an
existing issue). Significant changes must go through the change proposal process
before they can be accepted.

This process gives everyone a chance to validate the design, helps prevent
duplication of effort, and ensures that the idea fits inside the goals for the
service. It also checks that the design is sound before code is written; code
review is not the place for high-level discussions.

### Make the change

Once you have edited files, you must tell Git that they have been modified. You
must also tell Git about any files that are added, removed, or renamed files.
These operations are done with the usual Git commands, `git add`, `git rm`, and
`git mv`.

For any work, create a feature branch with

```
$ git checkout -b <branch-name>
```

The name `<branch-name>` is an arbitrary one you choose to identify the local
branch containing your changes. Make sure it clearly describes the work you are
doing. Avoid long lived branches because these become difficult to integrate
back into master.

When you are ready to submit a Pull Request on GitHub, rebase and squash commits
into meaningful chunks. Below is a template for a good Pull Request message:

```
com/events: Prevent users from submitting bad AWS credentials.
ui: Add form validations around AWS credentials.

The existing implementation allowed users to change their AWS credentials to
be invalid when they had already set up valid ones.

Fixes #159
```

Ideally there is one summary as your commit shouldn't cross too many top level
directories. Feel free to use bullet points for your detailed description of
changes.

The special notation "Fixes #159" or "Closes #159" associates the change with
issue 159 in the GitHub issue tracker. When this change is eventually submitted,
the issue tracker will automatically mark the issue as fixed.

Once you have finished writing the commit message, save the file and exit the
editor. And be sure to sign-off your commits.

### Sign your work

The sign-off is a simple line at the end of the explanation for the patch. Your
signature certifies that you wrote the patch or otherwise have the right to pass
it on as an open-source patch. The rules are pretty simple: if you can certify
the below (from [developercertificate.org](http://developercertificate.org/)):

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
660 York Street, Suite 102,
San Francisco, CA 94110 USA

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.

Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

Then you just add a line to every git commit message:

    Signed-off-by: Joe Smith <joe.smith@email.com>

Use your real name (sorry, no pseudonyms or anonymous contributions.)

If you set your `user.name` and `user.email` git configs, you can sign your
commit automatically with `git commit -s`, which you can add to any aliases you
might have.

### Code review and merging into master

No human merges code into master manually. Instead, work through Pull Requests.
The change is eligible to be merged into master with two contributors giving a
comment of "LGTM" (Looks Good to Me).

You can either submit a PR from a fork, or if you have commit access you can
push to a branch name then submit a PR from it.
