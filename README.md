
## Development

Development is optimized for OSX but you should be able to find equivalents
of everything for [Linux](#linux-dependencies) and probably Windows.

### Go Dependencies

The cmd project is written in Golang, so 
[install the latest version of go](https://golang.org/doc/install#install) for
the platform. From here, make sure your `GOPATH` is set and the cmd project is
cloned to `$GOPATH/src/gliderlabs/cmd` (this location is important). In some
cases you need to add $GOPATH/bin to your $PATH variable so the built binaries
are executable by following commands.

### OSX Dependencies

 * [Docker for Mac](https://docs.docker.com/docker-for-mac/install/) (we develop against the Edge channel)
 * [Homebrew](https://brew.sh/) so you can `brew bundle` the rest

### Linux Dependencies

Since each distribution of Linux has different packaging tools will need to
translate these dependencies for your own distribution.

 * `make` is required
 * [Docker Community Edition](https://www.docker.com/community-edition)
  * `docker-compose` is also required and if not available as a package for
    your environment, can be installed via `pip`.
 * `hugo` is required to make the 'www-dev' target

## Steps

1. make setup  - (Installs minor utilities)
2. make build

If these steps fail or you need to run them again `make clobber` between builds.
