## Deprecated use v3

v1 and v2 of `egoscale` are deprecated, please use [v3](https://pkg.go.dev/github.com/sauterp/egoscale/v3)

---
title: Egoscale
description: the Go library for Exoscale
---

<a href="https://gopherize.me/gopher/9c1bc7cfe1d84cf43e477dbfc4aa86332065f1fd"><img src="gopher.png" align="right" alt=""></a>

[![Actions Status](https://github.com/exoscale/egoscale/workflows/CI/badge.svg?branch=master)](https://github.com/exoscale/egoscale/actions?query=workflow%3ACI+branch%3Amaster)
[![GoDoc](https://godoc.org/github.com/exoscale/egoscale?status.svg)](https://godoc.org/github.com/exoscale/egoscale/v2) [![Go Report Card](https://goreportcard.com/badge/github.com/exoscale/egoscale)](https://goreportcard.com/report/github.com/exoscale/egoscale)

A wrapper for the [Exoscale public cloud](https://www.exoscale.com) API.

Actively maintained version is **v2**, it can be imported as:

```
	import "github.com/exoscale/egoscale/v2"
```

**Version v3 is in alpha state and breaking changes are expected.**

## Known users

- [Exoscale CLI](https://github.com/exoscale/cli)
- [Exoscale Terraform provider](https://github.com/exoscale/terraform-provider-exoscale)
- [ExoIP](https://github.com/exoscale/exoip): IP Watchdog
- [Lego](https://github.com/go-acme/lego): Let's Encrypt and ACME library
- Kubernetes Incubator: [External DNS](https://github.com/kubernetes-incubator/external-dns)
- [Docker machine](https://docs.docker.com/machine/drivers/exoscale/)
- [etc.](https://godoc.org/github.com/exoscale/egoscale?importers)

## License

Licensed under the Apache License, Version 2.0 (the "License"); you
may not use this file except in compliance with the License. You may
obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
