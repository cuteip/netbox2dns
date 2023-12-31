package netbox2dns

import (
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/netbox-community/go-netbox/v3/netbox/client"
	"github.com/scottlaird/netboxlib/netbox"
)

// GetNetboxIPAddresses fetches a list of IP Addresses from a Netbox server.
func GetNetboxIPAddresses(host, token string) (netbox.IPAddrs, error) {
	transport := httptransport.New(host, client.DefaultBasePath, []string{"https"})
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+token)
	c := client.New(transport, nil)

	return netbox.ListIPAddrs(c)
}
