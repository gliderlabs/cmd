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
[https|wss]://alpha.cmd.io/run/<username>/<command>
```
Query parameters:

* `access_token` - Token with access to this command. Required if not provided by Basic Auth.
* `args` - Optional string of arguments. They can also be provided as another path part.

### Example using curl

```
$ curl -u "$TOKEN:" https://alpha.cmd.io/run/progrium/helloworld
```
