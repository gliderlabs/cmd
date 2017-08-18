## com

Directory of component packages that make up Comlab CLI. Standard reusable Go
packages should live under `pkg` if they are not Comlab specific, or under `lib`
if they are Comlab specific.

### How to organize component packages

When a component package has a component object, its type must be `Component`
and should be defined and registered in `com.go`. Methods on the component that
are not part of an extension interface can also live in `com.go`. Extension
interface types provided by this component should also be defined in this file.

Extension interfaces the component *implements* should be defined in a
file by the name of the package that defines the extension interface. For
example, if a component implements the `web.Handler` interface, this would live
in `web.go`.

Exported package-level functions, variables, and types should live in
`<package>.go`. However, types that have a lot of code can live in their own
file.

Non-exported functions, variables, and types should live in the file closest to
where they are used. If they're used by several files, they should go in
`<package>.go`. If there are a lot of them, they can live in their own file.

## When to use sub-packages

Generally, component packages should live in the same namespace to promote a
more clear wide-over-deep component structure. However, there are cases when
sub-packages of a component make sense:

 * When the component has built-in backends or drivers (`users/auth0`)
 * When there a number of components of the same "class" (`app/...`)
 * When the components represent sub-functionality (`automation/...`)

## Other directories inside packages

There are two other directories you may find under a package directory:

 * `ui` - Web UI / frontend source files
 * `data` - Data files relevant to package
