# go-audit

[![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](http://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/travis/slackhq/go-audit.svg?style=flat-square)](https://travis-ci.org/slackhq/go-audit)
[![codecov](https://codecov.io/gh/slackhq/go-audit/branch/master/graph/badge.svg)](https://codecov.io/gh/slackhq/go-audit)

## About

go-audit is an alternative to the auditd daemon that ships with many distros.
After having created an [auditd audisp](https://people.redhat.com/sgrubb/audit/) plugin to convert audit logs to json,
I became interested in creating a replacement for the existing daemon.

##### Goals

* Safe : Written in a modern language that is type safe and performant
* Fast : Never ever ever ever block if we can avoid it
* Outputs json : Yay
* Pluggable pipelines : Can write to syslog, local file, or stdout. Additional outputs are easily written.
* Connects to the linux kernel via netlink (info [here](https://git.kernel.org/cgit/linux/kernel/git/stable/linux-stable.git/tree/kernel/audit.c?id=refs/tags/v3.14.56) and [here](https://git.kernel.org/cgit/linux/kernel/git/stable/linux-stable.git/tree/include/uapi/linux/audit.h?h=linux-3.14.y))

## Usage

##### Installation

1. Install [golang](https://golang.org/doc/install), version 1.7 or greater is required
2. Install [`govendor`](https://github.com/kardianos/govendor) if you haven't already

    ```go get -u github.com/kardianos/govendor```

2. Clone the repo

    ```
    git clone (this repo)
    cd go-audit
    ```

2. Build the binary

    ```
    make
    ```

3. Copy the binary `go-audit` to wherever you'd like

##### Testing

- `make test` - run the unit test suite
- `make test-cov-html` - run the unit tests and open up the code coverage results
- `make bench` - run the benchmark test suite
- `make bench-cpu` - run the benchmark test suite with cpu profiling
- `make bench-cpulong` - run the benchmark test suite with cpu profiling and try to get some gc collection

##### Running as a service

Check the [contrib](contrib) folder, it contains examples for how to run `go-audit` as a proper service on your machine.

##### Example Config

See [go-audit.yaml.example](go-audit.yaml.example)

## FAQ

#### I am seeing `Error during message receive: no buffer space available` in the logs

This is because `go-audit` is not receiving data as quickly as your system is generating it. You can increase
the receive buffer system wide and maybe it will help. Best to try and reduce the amount of data `go-audit` has
to handle.

If reducing audit velocity is not an option you can try increasing `socket_buffer.receive` in your config.
See [Example Config](#example-config) for more information

```
socket_buffer:
    receive: <some number bigger than (the current value * 2)>
```

#### Sometime files don't have a `name`, only `inode`, what gives?

The kernel doesn't always know the filename for file access. Figuring out the filename from an inode is expensive and
error prone.

You can map back to a filename, possibly not *the* filename, that triggured the audit line though.

```
sudo debugfs -R "ncheck <inode to map>" /dev/<your block device here>
```

#### I don't like math and want you to tell me the syslog priority to use

Use the default, or consult this handy table.

Wikipedia has a pretty good [page](https://en.wikipedia.org/wiki/Syslog) on this
