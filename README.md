# ssh-aliases
[![Build Status](https://travis-ci.org/dankraw/ssh-aliases.svg?branch=master)](https://travis-ci.org/dankraw/ssh-aliases) 
[![Go Report Card](https://goreportcard.com/badge/github.com/dankraw/ssh-aliases)](https://goreportcard.com/report/github.com/dankraw/ssh-aliases)

`ssh-aliases` is a command line tool that brings ease to living with `~/.ssh/config`.

In short, `ssh-aliases`:
* combines multiple [human friendly config files](#configuration-files) into a single `ssh` config file
* is able to generate a list of hosts out of a single entry by using [expanding expressions](#expanding-expressions), 
like `instance[1..3].example.com` or `[master|slave].example.com`
* creates aliases for hosts by compiling [templates](#alias-templates) (no need to write regexps)
* allows multiple hosts reuse the same `ssh` configuration
* is a single binary file

## Table of contents

* [Installation](#installation)
* [Configuration files](#configuration-files)
    * [Scanned directories](#scanned-directories)
    * [Components](#components)
        * [Host definitions](#host-definitions)
        * [Config properties](#config-properties)
            * [Extending configurations](#extending-configurations)
        * [Variables](#variables)
    * [Expanding hosts](#expanding-hosts)
        * [Expanding expressions](#expanding-expressions)
        * [Alias templates](#alias-templates)
    * [Tips and tricks](#tips-and-tricks)
* [Usage (CLI)](#usage-cli)
    * [`compile`](#compile---generating-configuration-for-ssh) - generating configuration for `ssh`
    * [`list`](#list---listing-aliases-definitions) - listing aliases definitions
* [License](#license)


## Installation

### Binary distribution

Binary releases for Linux and MacOS can be found on [GitHub releases page](https://github.com/dankraw/ssh-aliases/releases).

### Homebrew tap

MacOS users can install `ssh-aliases` easily using [Homebrew](https://brew.sh):

``` console
brew tap dankraw/ssh-aliases
brew install ssh-aliases
```

### Source code

If you are familiar with Go: 

``` console
go get github.com/dankraw/ssh-aliases
```

There is a `Makefile`, so you can use it as well:

``` console
make test # run tests
make fmt  # format code
make lint # run linters
make      # build binary to ./target/ssh-aliases
```

## Configuration files

`ssh-aliases` is a tool for `~/.ssh/config` file generation. 
The input for `ssh-aliases` are [HCL](https://github.com/hashicorp/hcl) config files. 
HCL was designed to be written and modified by humans. 
In some way it is similar to JSON, but is more expressive and concise at the same time, 
allows using comments, etc. 

Looking at examples below will be enough to become familiar with HCL format.

### Scanned directories

`ssh-aliases` allows you to divide your `ssh` configuration into multiple files depending on your needs.
When running `ssh-aliases` you point it to a directory (by default it's `~/.ssh_aliases`) 
containing any number of HCL config files. The directory will be scanned for files with `.hcl` extension.
Keep in mind it does not scan recursively - child directories won't be considered.

### Components

A single config file may contain any number of components defined in it. 
Currently there are two types of components:
* [Host definitions](#host-definitions)
* [Config properties](#config-properties)

#### Host definitions

A host definition consists of a `host` keyword and it's globally unique (among all scanned files) name.
Each `host` should contain following attributes:
* `hostname` - is a target hostname, possibly containing [expanding expressions](#expanding-hosts)
* `alias` - is an alias template for the destination hostname
* `config` - an embedded [config properties](#config-properties) definition, or a name (a `string`) that points 
to existing properties definition in the same or any other configuration file

An example host definition looks like:

``` hcl
host "my-service" {
  hostname = "instance[1..2].my-service.example.com",
  alias = "myservice{#1}"
  config = {
    user = "ubuntu"
    identity_file = "~/.ssh/my_service.pem"
    port = 22
    // etc.
  }
}
```

or (when pointing an external config named `my-service-config`)

``` hcl
host "my-service" {
  hostname = "instance[1..2].my-service.example.com",
  alias = "myservice{#1}"
  config = "my-service-config"
}
```

#### Config properties

A config properties definition consists of a `config` keyword and it's globally unique (among all scanned files) name.
It's body is a list of properties that map to `ssh_config` keywords and their values. 
A complete list of ssh config keywords can be seen [here](https://linux.die.net/man/5/ssh_config) 
or listed via `man ssh_config` in your terminal.

Each property may contain an underscore (`_`) in its keyword for clarity, 
all underscores are removed during config compilation, first character and all letters that follow underscores are capitalized - 
this makes generated file easier to read. For example, `identity_file` will become `IdentityFile` in the destination config file.
By design `ssh_config` keywords are case insensitive, and their values are case sensitive.

Provided properties are not validated by `ssh-aliases`, so it should work even if you have a custom built `ssh` command.

An example config properties definition may look like:

``` hcl
config "some-config" {
  user = "ubuntu"
  identity_file = "id_rsa.pem"
  port = 22
  // etc.
}
```

##### Extending configurations

A special property `_extend` can be used in order to include properties from other configurations.
When importing top level properties override lower level properties. 

```hcl
config "top-level" {
  user = "eden"
  _extend = "lower-level"
  # port = 2222 (will be included from below configuration)
  # ...
}

config "lower-level" {
  user = "helix"
  port = 2222
  # etc.
}
```

A single configuration may extend multiple configurations, an array of configuration names should be provided, for example:

```hcl
config "top-level" {
  user = "eden"
  _extend = ["lower-level1", "lower-level2"] 
  # configurations are included from left to right (and overridden in this order)
  # port = 4444
}

config "lower-level1" {
  user = "helix"
  port = 2222
  # etc.
}

config "lower-level2" {
  user = "torus"
  port = 4444
  # etc.
}
```

#### Variables

Variables are declared in object blocks marked with `var` keyword. There may be many `var` blocks distributed along multiple files, but variable names have global scope, so each one can be declared only once.

Example variables block may look like:

```hcl
var {
    dc1 = "my.domain1.example.com"
    dc2 = "some.other.domain2.net"
    keys {
        service_a = "/path/to/a_key.pem"
        service_b = "/path/to/b_key.pem"
    }
    nodes {
      service_a = 5
    }
    users {
      a = "eden"
      b = "helix"
    }
}
```

Variables can be nested, their lookup names are flattened during processing with `.` separator.
In this example we have defined following variables: `dc1`, `dc2`, `keys.service_a`, `keys.service_b`, `nodes.service_a`, `users.a`, `users.b`.

Variables can be used in:
* Aliases 
* Hostnames
* Config property values

String interpolation with variables is done by using a `${name}` placeholder, for example:

```hcl
host "service-a" {
  hostname = "instance[1..${nodes.service_a}].${dc1}",
  alias = "myservice{#1}"
  config = {
    user = "${users.a}"
    identity_file = "${keys.service_a}"
  }
}
```

### Expanding hosts

The most important feature of `ssh-aliases` is *hosts expansion*.
It's a mechanism of generating multiple `Host ...` entries in the destination `ssh_config` out of a single [host definition](#host-definitions).
It is done by using *range expressions* in hostnames and compiling host aliases from templates.

#### Expanding expressions

There are two types of *expanding expressions* available:
* ranges
* sets

A **range** is represented as `[m..n]`, where `m` and `n` are positive integers that `m < n`. 
For example, a hostname `instance[1..3].example.com` will be expanded to:

``` console
instance1.example.com
instance2.example.com
instance3.example.com
```

A **set** is represented as `[a]`, `[a|b]`, `[a|b|c]` and so on, 
where `a`, `b`, `c`... are some arbitrary strings of characters allowed in hostnames.
For example, a hostname `server.[dev|test|prod].example.com` will be expanded to:

``` console
server.dev.example.com
server.test.example.com
server.prod.example.com
```

Of course ranges and sets can be used together multiple times each. 
A final result will be a [cartesian product](https://en.wikipedia.org/wiki/Cartesian_product) of all expanding expressions provided.
For example, a hostname `server[1..2].[dev|test].example.com` would be expanded to:

``` console
server1.dev.example.com
server1.test.example.com
server2.dev.example.com
server2.test.example.com
```

#### Alias templates

Each generated `Host ...` entry needs to have an alias that is provided in [host definition](#host-definitions). 
If hostnames are [expanded](#expanding-expressions) it is required to provide **placeholders** 
for all expanding expressions used. Otherwise `ssh-aliases` would generate 
the same alias for multiple hostnames, and that simply makes no sense.

An expanding expression placeholder is represented as `{#n}`, where `n=1,2...k`, 
`n` points the `n-th` expression used in hostname (sequence from left to right) 
and so `k` is the number of expanding expressions used in total. 

For example, `{#1}` points the first expression used in hostname, `{#2}` points the second, and so on.

If we look at the hostname example from above section `server[1..2].[dev|test].example.com`, we have two expressions used:
1. `[1..2]`
2. `[dev|test]`

We can declare an alias template like `{#2}.server{#1}`, 
which would compile following aliases for the generated hostnames:

``` console
dev.server1
test.server1
dev.server2
test.server2
```

### Tips and tricks

* Generated `ssh_config` configuration can be used not only with `ssh` command, but with other OpenSSH client commands, like `scp` and `sftp`
* Multiple alias templates may be provided for the same host definition, for example:

```hcl
host "my-service" {
  hostname = "instance[1..2].myservice.example.com",
  alias = "myservice{#1} ms{#1}" # separated with space
  config = "some-config"
}
```

* `ssh_config` (v7.2+) ships with `Include` directive ([see docs](https://man.openbsd.org/ssh_config.5#Include)) that can be used to include other files. This can be useful for mixing `ssh-aliases` generated configs with pure `ssh_config` files:

```ssh_config
Include path/to/ssh-aliases/generated/ssh_config

# below some legacy ssh config that one day may be migrated to ssh-aliases
Host myservice
    HostName myservice.example.com
    User myself
# ...
```

* `config` properties are optional when `alias` is provided:

```hcl
host "example" {
    hostname = "my.service[1..2].example.com"
    alias = "myservice{#1}"
}
```

* `hostname` is optional when `config` properties are provided. This can be useful for creating wildcard (`*`) configurations that match any host:

```hcl
host "all-hosts" {
    alias = "*"
    config {
        # ...
    }
}
```

## Usage (CLI)

Run `ssh-aliases --help` to see available options of the `ssh-aliases` command line interface (CLI).

In general, there are only two commands available:
* `compile` - prints (or saves to a file) compiled `ssh` config
* `list` - prints preview of generates aliases and hostnames

Both commands share the same global option: `--scan` or `-s` which should point to the directory 
containing [input config files](#configuration-files).
If omitted, `ssh-aliases` will look for `~/.ssh_aliases` directory.
This option should be passed *before* the selected command name.

### `compile` - generating configuration for `ssh`

`compile` is the primary command in `ssh-aliases` - it combines together all input config files 
and compiles configuration for `ssh`.  

Options for `compile`

* `--save` - adding this option makes `ssh-aliases` save the output to the file instead of printing to `stdout`, 
asks for confirmation if the file exists (unless `--force` is used) and overwrites its contents if accepted
* `--file <PATH>` - when using `--save` it tells where should the file be saved, defaults to `~/.ssh/config`
* `--force` - when using `--save` it will overwrite possibly existing file without confirmation
* `--help` - shows command usage

Example command run with all options provided:

``` console
$ ssh-aliases --scan ~/my_custom_dir compile --save --file ~/.ssh/ssh_aliases_config --force 
```

Now, let's suppose we have `./examples/readme` directory that contains two files:

``` hcl
# ./example/readme/example_service_1.hcl
host "abc" {
  hostname = "node[1..2].abc.[dev|test].example.com"
  alias = "{#2}.abc{#1}"
  config = "abc-config"
}

config "abc-config" {
  user = "ubuntu"
  identity_file = "~/.ssh/abc.pem"
  port = 22
}
```

``` hcl
# ./example/readme/example_service_2.hcl
host "other" {
  hostname = "other[1..2].example.com"
  alias = "other{#1}"
  config = {
    user = "lurker"
    identity_file = "~/.ssh/other.pem"
    port = 22
  }
}
```

Let's run `compile` command for that directory:

``` console
$ ssh-aliases --scan ./examples/readme compile
```

`ssh-aliases` will print:

``` console
Host dev.abc1
     HostName node1.abc.dev.example.com
     IdentityFile ~/.ssh/abc.pem
     Port 22
     User ubuntu

Host dev.abc2
     HostName node2.abc.dev.example.com
     IdentityFile ~/.ssh/abc.pem
     Port 22
     User ubuntu

Host test.abc1
     HostName node1.abc.test.example.com
     IdentityFile ~/.ssh/abc.pem
     Port 22
     User ubuntu

Host test.abc2
     HostName node2.abc.test.example.com
     IdentityFile ~/.ssh/abc.pem
     Port 22
     User ubuntu

Host other1
     HostName other1.example.com
     IdentityFile ~/.ssh/other.pem
     Port 22
     User lurker

Host other2
     HostName other2.example.com
     IdentityFile ~/.ssh/other.pem
     Port 22
     User lurker
```

### `list` - listing aliases definitions

`list` command should be used to check correctness of declared [hostname patterns](#host-definitions) 
and [alias templates](#alias-templates). 
It will print a concise list of compiled results, yet omitting linked [config properties](#config-properties). 

For example, let's run `list` for `./examples/readme` directory from previous paragraph:
 
``` console 
$ ssh-aliases --scan ./examples/readme list
```
The printed result will be:

``` console
readme/example_service_1.hcl (1):

 abc (4):
  dev.abc1: node1.abc.dev.example.com
  dev.abc2: node2.abc.dev.example.com
  test.abc1: node1.abc.test.example.com
  test.abc2: node2.abc.test.example.com

readme/example_service_2.hcl (1):

 other (2):
  other1: other1.example.com
  other2: other2.example.com
```

## License

`ssh-aliases` is published under [MIT License](LICENSE).
