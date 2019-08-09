# Egoscale: official Go library for the Exoscale API

[![Build Status](https://travis-ci.com/exoscale/egoscale.svg?branch=master)](https://travis-ci.org/exoscale/egoscale) [![GoDoc](https://godoc.org/github.com/exoscale/egoscale?status.svg)](https://godoc.org/github.com/exoscale/egoscale)

This library allows developpers to use the [Exoscale] cloud platform API with high-level Go bindings.

## Usage

```go
package main

import (
    "fmt"

    "github.com/exoscale/egoscale"
)

function main() {
	exo, err := egoscale.NewClient(egoscale.ConfigFromEnv())
	if err != nil {
		// ...
	}

	zones, err := exo.Compute.ListZones()
	if err != nil {
		// ...
    }

	for _, zone := range zones {
		fmt.Println(zone.ID, zone.Name)
	}
```

## Client configuration

The library `Client` object can be instantiated by passing the `NewClient()` function a [`ConfigFunc`][configfunc]
closure returning a configuration profile from a literal configuration profile, a configuration file or environment
variables. If no `ConfigFunc` argument is provided, the client initialization looks up the default configuration file
location `$HOME/.exoscale/config.toml` then environment variables (see below for the environment variables list).
In case multiple ConfigFunc are provided, they are processed in the provided order and the first successful execution
halts the configuration lookup and the returned configuration profile is used to initialize the client.

### From a literal configuration profile

To instantiate an `egoscale.Client` object from a literal configuration profile:

```go
exo, err := egoscale.NewClient(egoscale.ConfigFromProfile(egoscale.ConfigProfile{
    APIKey: "EXOa1b2c3...",
    APISecret: "...",
}))
```

### From a configuration file

To instantiate an `egoscale.Client` object from a configuration file:

```go
exo, err := egoscale.NewClient(egoscale.ConfigFromFile("/path/to/exoscale/config.toml"))
```

The configuration file format is [TOML][toml], the expected structure is:

```toml
default_profile = "alice"

[[profiles]]
name = "alice"
api_key = "EXOa1b2c3..."
api_secret = "..."

[[profiles]]
name = "bob"
api_key = "EXOd5e6f7..."
api_secret = "..."
```

A `[[profiles]]` entry is a dictionary supporting the following key/values (keys marked with a "*" are required):

* `name`*: the name of the profile
* `api_key`*: the profile Exoscale client API key
* `api_secret`*: the profile Exoscale client API secret
* `storage_zone`: an Exoscale Object Storage zone (required for using the Storage API)
* `compute_api_endpoint`: an alternative Exoscale Compute API endpoint
* `dns_api_endpoint`: an alternative Exoscale DNS API endpoint
* `storage_api_endpoint`: an alternative Exoscale Storage API endpoint
* `runstatus_api_endpoint`: an alternative Exoscale Runstatus API endpoint

### From environment variables

To instantiate an `egoscale.Client` object from environment variables:

```go
exo, err := egoscale.NewClient(egoscale.ConfigFromEnv())
```

The following environment variables can be used in place of a configuration file:

* `EXOSCALE_API_KEY`: a Exoscale client API key
* `EXOSCALE_API_SECRET`: a Exoscale client API secret
* `EXOSCALE_COMPUTE_API_ENDPOINT`: an alternative Exoscale Compute API endpoint
* `EXOSCALE_DNS_API_ENDPOINT`: an alternative Exoscale DNS API endpoint
* `EXOSCALE_RUNSTATUS_API_ENDPOINT`: an alternative Exoscale Runstatus API
  endpoint
* `EXOSCALE_STORAGE_API_ENDPOINT`: an alternative Exoscale Storage API endpoint
* `EXOSCALE_STORAGE_ZONE`: an Exoscale Storage zone
* `EXOSCALE_CONFIG_FILE`: an alternative configuration file location (default:
  `$HOME/.exoscale/config.toml`)

## License

Licensed under the Apache License, Version 2.0 (the "License"); you
may not use this file except in compliance with the License. You may
obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

[toml]: https://github.com/toml-lang/toml
[exoscale]: https://www.exoscale.com/
[configfunc]: https://godoc.org/github.com/exoscale/egoscale#ConfigFunc