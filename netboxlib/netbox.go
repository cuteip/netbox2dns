package netboxlib

import (
	"net/netip"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/netbox-community/go-netbox/v3/netbox/client"
	"github.com/netbox-community/go-netbox/v3/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/v3/netbox/models"
)

type IpamIPAddress struct {
	Address netip.Addr
	DNSName string
	Status  string
}

type Client struct {
	api *client.NetBoxAPI
}

func NewClient(host, token string) *Client {
	transport := httptransport.New(host, client.DefaultBasePath, []string{"https"})
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+token)
	c := client.New(transport, nil)
	return &Client{
		api: c,
	}
}

func (c *Client) GetNetboxIPAddresses(queryParameters []string) ([]IpamIPAddress, error) {
	param := ipam.NewIpamIPAddressesListParams()
	var limit int64
	limit = 0
	param.SetLimit(&limit)

	falseStrPtr := "false"
	param.SetDNSNameEmpty(&falseStrPtr)

	tagn := "netbox2dns_exclude"
	param.SetTagn(&tagn)

	res, err := c.api.Ipam.IpamIPAddressesList(param, nil)
	if err != nil {
		return nil, err
	}

	var iipAddresses []IpamIPAddress
	for _, result := range res.Payload.Results {
		iipAddress, err := covertModelsIPAddressToIpamIPAddress(*result)
		if err != nil {
			return nil, err
		}
		iipAddresses = append(iipAddresses, iipAddress)
	}
	return iipAddresses, nil
}

func covertModelsIPAddressToIpamIPAddress(m models.IPAddress) (IpamIPAddress, error) {
	// m.Address example: 192.0.2.1/32, 192.0.2.2/24, ...
	prefix, err := netip.ParsePrefix(*m.Address)
	if err != nil {
		return IpamIPAddress{}, err
	}
	return IpamIPAddress{
		Address: prefix.Addr(),
		DNSName: m.DNSName,
		Status:  *m.Status.Value,
	}, nil
}
