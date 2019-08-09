* Implement MustGet* methods (e.g. MustGetZone(<GetZone() args>)) that don't return errors but panic() instead
* Disable internal/egoscale environment variable support
* Check custom HTTP client User-Agent
* Use https://github.com/hashicorp/go-cleanhttp (in Pooled mode)
