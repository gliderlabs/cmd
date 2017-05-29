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

#### Why are there links to Github that return 404 Not Found?

The Cmd source and Github project are currently not public, but anybody can request
access after signing up for Cmd by asking for it in the `#cmd` Slack channel.

#### When will the Github project and source code be public?

Either after our stable release or during our public beta. We want to grow that
community slowly, and also separate source availability from service launch
milestones. This is to avoid confusion from the prevailing assumption that 
open source means "run it yourself." While that is possible, our focus is on the
community-run service that happens to be open source by necessity. 
