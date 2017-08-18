/*
Package log implements an opinionated, structured logging API and data model.
It's not actually a logger as it requires you to provide an Observer to actually
handle log events, which defines what it means to "log" (print to stdout, write
to disk, send to aggregator, etc).

Instead of providing nestable logger instances, you use the package-level API
and it will determine the calling package for context. The calling package is
added as a field to the Event produced by the call.

There are 4 types of log events you can create: Info, Debug, Local, Fatal. They
are more semantic tags than log levels. Fatal is a special log event type that
will cause the program to terminate non-zero after logging. Local and Debug are
tied to modes.

Local is a mode that implies local development. Use Local for caveman debugging
or temporary log events. Observers can choose to highlight these events so they
stand out to make development and local debugging easier. In many cases, Local
is unnecessary and you can just use Debug.

Debug is a mode that implies verbose logging for better debugging. This is often
used with local development, but can also be enabled in certain deployments such
as dev or staging. In some cases it might be turned on temporarily in
production.

When you log with Local or Debug, whether their mode is enabled will determine
if those log events will be processed. Otherwise, you log with Info, which is
intended for useful log events and messages.

Info, Debug, Local, and Fatal all take an arbitrary number of arguments of any
type. It's the job of the FieldProcessor to convert a value into appropriate
key-value fields. These fields are added to an Event object with the type and
timestamp.

This Event object is then passed to registered Observers. Observers can then
print, or colorize and print, or send to a log aggregator, or detect certain
errors and send to exception trackers, etc.

The package also comes with a wrapper for http.ResponseWriter to use this
package for HTTP logging.

*/
package log
