# Grift

Grift is a very simple library that allows you to write simple "task" scripts in Go and run them by name without having to write big `main` type of wrappers. Grift is similar to, and inspired by, [Rake](http://rake.rubyforge.org).

## Why?

Excellent question! When building applications there comes a point where you need different scripts to do different things. For example, you might want a script to seed your database, or perhaps a script to parse some logs, etc...

Grift lets you write these scripts using Go in a really simple and extensible way.

## Installation

Installation is really easy using `go get`.

```text
$ go get github.com/markbates/grift
```

You can confirm the installation by running:

```text
$ grift jim
```

## Usage/Getting Started

Apart from having the binary installed, the only other requirement is that the package you place your grifts in is called `grift`. That's it.

By running the following command:

```text
$ grift init
```

When you run the `init` sub-command Grift will generate a new `grifts` package and create a couple of simple grifts for you.

#### List available grifts

```text
$ grift list
```

#### Say Hello!

```text
$ grift hello
```

##

That's really it! Grift is meant to be simple. Write your grifts, use the full power of Go to do it.

For more information I would highly recommend checking out the [docs](https://godoc.org/github.com/markbates/grift/grift).
