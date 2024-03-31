package netbox2dns

import (
	"context"
	"fmt"
)

// DNSProvider is an interface to a DNS provider backend, such a ZoneFile.
type DNSProvider interface {
	WriteRecord(cz *ConfigZone, r *Record) error
	Save(cz *ConfigZone) error
}

// NewDNSProvider creates a provider of the correct type for the described zone.
func NewDNSProvider(ctx context.Context, cz *ConfigZone) (DNSProvider, error) {
	switch cz.ZoneType {
	case "zonefile":
		return NewZoneFileDNS(ctx, cz)
	default:
		return nil, fmt.Errorf("Unknown DNS provider type %q", cz.ZoneType)
	}
}
