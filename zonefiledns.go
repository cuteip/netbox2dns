package netbox2dns

import (
	"context"

	"github.com/scottlaird/netbox2dns/zonefile"
)

// ZoneFileDNS provides an implementation of DNS using traditional
// BIND-style zone files.
type ZoneFileDNS struct {
	zone *zonefile.Zone
}

// NewZoneFileDNS creates a new ZoneFileDNS object.
func NewZoneFileDNS(ctx context.Context, cz *ConfigZone) (*ZoneFileDNS, error) {
	zone, err := zonefile.New(cz.Filename)
	if err != nil {
		return nil, err
	}

	zfd := &ZoneFileDNS{
		zone: zone,
	}

	return zfd, nil
}

func (zfd *ZoneFileDNS) rrFromRecord(cz *ConfigZone, r *Record) zonefile.ResourceRecord {
	return zonefile.ResourceRecord{
		Name:  r.Name,
		Type:  r.Type,
		Class: "IN",
		TTL:   uint32(r.TTL),
		Rdata: r.Rrdatas,
	}
}

// WriteRecord writes a Record to the zonefile behind the ZoneFileDNS.
// Note that this won't actually be written until 'Save()' is called.
func (zfd *ZoneFileDNS) WriteRecord(cz *ConfigZone, r *Record) error {
	entry := zfd.rrFromRecord(cz, r)

	err := zfd.zone.Add(entry)
	if err != nil {
		return err
	}

	return nil
}

// Save flushes the current zonefile to disk.  Without this, no
// changes will be written out.
func (zfd *ZoneFileDNS) Save(cz *ConfigZone) error {
	return zfd.zone.Save()
}
