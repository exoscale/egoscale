package v3

// DNSAPI provides access to [Exoscale DNS] API resources.
//
// [Exoscale DNS]: https://community.exoscale.com/documentation/dns/
type DNSAPI struct {
	client *Client
}

//func (a *DNS) Domains() *Domains {
//return NewDomains(a.oapiClient)
//}
