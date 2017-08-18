/*
Package com provides a registry for component objects, which are used with Go
packages to create component packages. A component package is a Go package that
registers a struct with com, typically named Component, that can implement hook
interfaces of other components. This allows component packages to extend and
hook into each other in a loosely-coupled way.

For example, imagine an issue tracking system like GitHub issues built as a
package. It's completely separate from the source control system, but needs to
hook into certain events to implement features like automatically closing issues
mentioned in commit messages.

  package issues

  import "github.com/gliderlabs/comlab/pkg/com"

  func init() {
    com.Register("issues", &Component{})
  }

  type Component struct {}

  func (c *Component) GitPostCommit(commit git.Commit) {
    if issue, ok := MentionsIssue(commit.Message, []string{"closes", "fixes"}); ok {
      issue.Close()
    }
  }

Registering a component with com also lets you define component-specific options
or configuration that can easily be accessed from within that package. These
options are populated by a configuration provider such as Viper that can read
the environment or a single configuration file to configure all components. In
other words, component packages can define and access their own configuration
in a unified way.

  package issues

  import "github.com/gliderlabs/comlab/pkg/com"

  func init() {
    com.Register("issues", &Component{},
      com.Option("keywords", []string{"closes", "fixes"}, "Keywords to identify issue mentions"}))
  }

  type Component struct {}

  func (c *Component) GitPostCommit(commit gitscm.Commit) {
    if issue, ok := MentionsIssue(commit.Message, com.GetStrings("keywords")); ok {
      issue.Close()
    }
  }

Components can be enabled or disabled, either through configuration or any
system that implements a component context that says whether a component is
enabled. This means when you implement a feature as a component, you have an
easy way to turn it on and off at boot-time or at runtime, potentially per user.
Components give you the foundation for feature flags.

The way you query the component registry is by interface. If your component
package declared an interface, and other registered components implemented that
interface, then you can get back all the enabled component objects implementing
that interface. This is intended to be used as a way to expose hooks or
extension points between components that can interact however you define in your
interfaces.

  package gitscm

  import "github.com/gliderlabs/comlab/pkg/com"

  func init() {
    com.Register("gitscm", &Component{})
  }

  type Component struct {}

  type CommitObserver interface {
    GitPostCommit(commit Commit)
  }

  // ... rest of the gitscm package implemention here, then ...

  func HandleCommit(commit Commit) {

    // ... code to handle the commit, then ...

    for _, observer := range com.Enabled(new(CommitObserver), nil) {
      observer.(CommitObserver).GitPostCommit(commit)
    }
  }

You can also use the registry for pluggable modules like drivers or backends. A
host component package can provide an interface that multiple registered driver
components can implement, then you can select which component to use by name.
Perhaps specified by the user via configuration.

  package gitscm

  import "github.com/gliderlabs/comlab/pkg/com"

  func init() {
    com.Register("gitscm", &Component{},
      com.Option("storage_backend", "gitscm.filesystem", "Git storage backend"))
  }

  type Component struct {
    store StorageBackend
  }

  type StorageBackend interface {
    // ...
  }

  // ...

  func (c *Component) AppInitialize() {
    backend := com.Select(com.GetString("storage_backend"), new(StorageBackend))
    if backend == nil {
      log.Fatal("Storage backend not registered:", com.GetString("storage_backend"))
    } else {
      c.store = backend.(StorageBackend)
    }
  }


These minimal features work together to provide a sort of lightweight component
framework, allowing you to build applications that are component-oriented,
highly modular, and highly extensible.

Although com is a standalone package, it's maintained as part of Comlab, which
expands on this package with tooling and patterns for building
component-oriented applications with it.

*/
package com
