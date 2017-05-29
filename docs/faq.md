---
date: 2017-01-31T18:00:00-06:00
title: FAQ
weight: 10
---

#### Where can I go for help using Cmd?

Besides this documentation, you can [join our Slack](/community/) (the `#cmd`
channel) and the community will do their best to help. This is preferred before
submitting a Github issue for general help, however if you believe you found a
bug feel to submit an issue.

#### Can config be used for secrets?

Config is intended to be used for secrets, however it's important to understand
the risks. Config values are encrypted at rest, but are available unencrypted
in-memory when commands are run. This is protected by our container isolation
configuration around Docker which itself is open source and continuously
improved.

When sharing access to commands, the ability for those users to read config
exposed as environment variables depends on your command and the programs your
command might use. If a program allows for arbitrary shell execution or displaying
the entire environment, the user of the command may be able to read your config.
