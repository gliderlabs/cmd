---
date: 2017-01-31T18:00:00-06:00
title: API Reference
menu: api
type: api
weight: 10
---

Commands can be run over secure HTTP and WebSocket using [access tokens](/cli/tokens/). This can be used to run commands when
you don't have SSH or need something more programmatic.

When commands are run via this API, they receive a representation of the HTTP request using [CGI](https://en.wikipedia.org/wiki/Common_Gateway_Interface). This means request metadata will be added to the environment and the body of the request will be passed as STDIN.

### Authentication

The Run API requires the use of [access tokens](/cli/tokens/), which can be created and given access to one or more commands. The token can then be used as the user in Basic Auth or as the query param `access_token`.

### Endpoint

```
[https|wss]://alpha.cmd.io/run/<username>/<command>[query]
```
Query parameters:

* `access_token` - Token with access to this command. Required if not provided by Basic Auth.
* `args` - Optional string of arguments to send to command, should be `+` separated and/or URL encoded (ex. https://alpha.cmd.io/run/hansgruber/shoot?args=the+glass).
* Additional query parameters - You may include any other query parameters as you wish, they will be injected as environment variables into your command via [CGI](https://en.wikipedia.org/wiki/Common_Gateway_Interface).

### `printenv` command example

This is a simple example command that will just echo `env` to stdout and exit. It's useful for debugging and will show exactly what environment variables are set for you to use in your scripts.

```
$ cat <<'EOD'
#!cmd.io alpine bash
#!/bin/bash
env && echo "Args: $@"
EOD
|ssh <username>@alpha.cmd.io :create printenv && unset test_cmd

Creating command... done
```

Next, grant access to a token in order to use the API endpoint. Make a note of the token output.
```
$ ssh <username>@alpha.cmd.io :access printenv grant $(ssh <username>@alpha.cmd.io :tokens new printenv)

Granting 18e370e1-6eae-4e2d-9b60-7c74d08c4333 access to printenv... done
```

Now issue curl commands to execute the `printenv` command:
```
$ curl -u "${TOKEN}:" https://alpha.cmd.io/run/<username>/printenv

HOSTNAME=alpha.cmd.io
SERVER_PORT=443
HTTP_HOST=alpha.cmd.io
USER=18e370e1-6eae-4e2d-9b60-7c74d08c433
CMD_NAME=printenv
REQUEST_URI=/run/chiefy/printenv
SCRIPT_NAME=printenv
PATH_INFO=/run/chiefy/printenv
CMD_CHANNEL=alpha
CMD_VERSION=alpha
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
PWD=/cmd
REMOTE_ADDR=
SHLVL=1
HOME=/root
SERVER_NAME=alpha.cmd.io
CONTENT_LENGTH=0
SERVER_SOFTWARE=cmd.io
QUERY_STRING=
GATEWAY_INTERFACE=CGI/1.1
SERVER_PROTOCOL=HTTP/1.1
CONTENT_TYPE=
REQUEST_METHOD=GET
_=/usr/bin/env
Args:
```
There were no query parameters sent, so these variables are mostly empty. Note the args array is also empty since there is no `args` query parameter.

Next, use the `args` query parameter to send a command's arguments. The command can access these arguments the same as a normal shell script (${1}, ${2}, etc.).
```
$ curl -u "${TOKEN}:" https://alpha.cmd.io/run/<username>/printenv?args=hi+two+three+four

HOSTNAME=alpha.cmd.io
SERVER_PORT=443
HTTP_HOST=alpha.cmd.io
USER=18e370e1-6eae-4e2d-9b60-7c74d08c433
CMD_NAME=printenv
REQUEST_URI=/run/chiefy/printenv?args=hi+two+three+four
SCRIPT_NAME=printenv
PATH_INFO=/run/chiefy/printenv
CMD_CHANNEL=alpha
CMD_VERSION=alpha
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
PWD=/cmd
REMOTE_ADDR=
SHLVL=1
HOME=/root
SERVER_NAME=alpha.cmd.io
CONTENT_LENGTH=0
SERVER_SOFTWARE=cmd.io
QUERY_STRING=args=hi+two+three+four
GATEWAY_INTERFACE=CGI/1.1
SERVER_PROTOCOL=HTTP/1.1
CONTENT_TYPE=
REQUEST_METHOD=GET
_=/usr/bin/env
Args: hi two three four
```

Finally, adding arbitrary query string parameters is allowed and are exposed as environment variables per CGI.

```
$ curl -u "${TOKEN}:" https://alpha.cmd.io/run/<username>/printenv?val=true&something=else&whatever=false+maybe"

HOSTNAME=alpha.cmd.io
SERVER_PORT=443
HTTP_HOST=alpha.cmd.io
USER=18e370e1-6eae-4e2d-9b60-7c74d08c433
CMD_NAME=printenv
REQUEST_URI=/run/chiefy/printenv?val=true&something=else&whatever=false+maybe
SCRIPT_NAME=printenv
PATH_INFO=/run/chiefy/printenv
CMD_CHANNEL=alpha
CMD_VERSION=alpha
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
PWD=/cmd
REMOTE_ADDR=
SHLVL=1
HOME=/root
SERVER_NAME=alpha.cmd.io
CONTENT_LENGTH=0
SERVER_SOFTWARE=cmd.io
QUERY_STRING=val=true&something=else&whatever=false+maybe
GATEWAY_INTERFACE=CGI/1.1
SERVER_PROTOCOL=HTTP/1.1
CONTENT_TYPE=
REQUEST_METHOD=GET
_=/usr/bin/env
Args:
```

In this example, your command would need to parse the `QUERY_STRING` variable programatically in order to access those variables.
