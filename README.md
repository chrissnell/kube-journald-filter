# kube-journald-filter
A simple journalctl-like reader for journald that adds Kubernetes metadata to logged messages

## About
```kube-journald-filter``` tails messages from [journald](http://www.freedesktop.org/software/systemd/man/systemd-journald.service.html) and extracts Kubernetes metadata from the entries, which are then printed to stdout.

## Example
Here is an example log message printed by `kube-journald-filter`:

`Dec 18 19:32:34 coreos-005.srv.example.com docker[2121] NS=webby POD=webby-ez9uk 2015/12/18 19:32:33 / 10.85.93.0:34332 curl/7.43.0`

`NS=<...>`  - The namespace that's logging

`POD=<...>` - The pod ID that's logging

Non-Kubernetes logged messages will be printed as usual, without the `NS=` and `POD=` tags.

## Installation
Download a binary from the releases or build one yourself.

### Building your own binary.

> NOTE: v1.1 depends on some new functionality in [github.com/coreos/go-systemd](https://github.com/coreos/go-systemd) that I have sent upstream as PRs.  Since I'm not vendoring this library and these PRs haven't been accepted yet, you will not be able to build this latest version from source unless you manually merge the PRs ([1](https://github.com/coreos/go-systemd/pull/128), [2](https://github.com/coreos/go-systemd/pull/129)) into go-systemd.  I recommend you download the binary release instead, which should work so long as you have libsystemd installed (```apt-get -y install libsystemd-dev``` on Ubuntu).  Alternatively, you can check out the older v1.0 tag and build without the go-systemd changes.

You will need the systemd shared libraries.  For Ubuntu, you can install them like so:
```
% sudo apt-get install libsystemd-dev
```

You will also need a working [Go](http://golang.org) development environment with GOROOT and GOPATH set.

Install the kubernetes libraries:
```
go get k8s.io/kubernetes/...
```

If go balks about a docker "units" library, go into your ```$GOPATH/src/github.com/docker/docker``` directory and run this:
```
% git checkout v1.9.0
```
and then re-run ```go get k8s.io/kubernetes/...```

Finally, fetch and build this package:
```
go get github.com/chrissnell/kube-journald-filter
```

# Running
Simple.  Just run the ```kube-journald-filter``` executable.  No flags necessary.  Pipe to your favorite log collection utility.
