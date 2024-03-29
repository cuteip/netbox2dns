package netbox2dns

import (
	"context"
	"os"
	"strings"

	"github.com/shuLhan/share/lib/dns"
)

// ZoneFileDNS provides an implementation of DNS using traditional
// BIND-style zone files.
type ZoneFileDNS struct {
	zone *dns.Zone
}

// NewZoneFileDNS creates a new ZoneFileDNS object.
func NewZoneFileDNS(ctx context.Context, cz *ConfigZone) (*ZoneFileDNS, error) {
	// dns.ParseZoneFile() は既にファイルが存在していることを期待している
	_, err := os.Stat(cz.Filename)
	if os.IsNotExist(err) {
		_, err = os.Create(cz.Filename)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	zone, err := dns.ParseZoneFile(cz.Filename, cz.Name, uint32(cz.TTL))
	if err != nil {
		return nil, err
	}

	zfd := &ZoneFileDNS{
		zone: zone,
	}

	return zfd, nil
}

// ImportZone reads DNS entries from a zone file on disk (as specified
// as part of the zone config in the netbox2dns config file) and
// populates the ZoneFileDNS with them.
func (zfd *ZoneFileDNS) ImportZone(cz *ConfigZone) (*Zone, error) {
	zone := &Zone{
		Name:     cz.Name,
		Filename: cz.Filename,
		TTL:      cz.TTL,
		Records:  make(map[string][]*Record),
	}

	for _, i := range zfd.zone.Records {
		for _, entry := range i {
			r := Record{
				Name: strings.TrimRight(entry.Name, ".") + ".",
				Type: dns.RecordTypeNames[entry.Type],
				TTL:  int64(entry.TTL),
			}

			s := entry.Value.(string)
			if r.Type == "A" || r.Type == "AAAA" {
				r.Rrdatas = []string{s}
			} else {
				// PTR, etc ... ?
				r.Rrdatas = []string{strings.TrimRight(s, ".") + "."}
			}

			zone.AddRecord(&r)
		}
	}

	return zone, nil
}

// rrFromRecord creates a DNS ResourceRecord from a netbox2dns Record.
func (zfd *ZoneFileDNS) rrFromRecord(cz *ConfigZone, r *Record) *dns.ResourceRecord {
	return &dns.ResourceRecord{
		Name:  r.NameNoDot(),
		Type:  dns.RecordTypes[r.Type],
		Class: dns.RecordClassIN,
		TTL:   uint32(r.TTL),
		Value: r.RrdataNoDot(),
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

// RemoveRecord removes a Record from the zonefile behind the ZoneFileDNS.
// Note that this won't actually be written until 'Save()' is called.
func (zfd *ZoneFileDNS) RemoveRecord(cz *ConfigZone, r *Record) error {
	entry := zfd.rrFromRecord(cz, r)

	err := zfd.zone.Remove(entry)
	if err != nil {
		return err
	}
	return nil
}

// Save flushes the current zonefile to disk.  Without this, no
// changes will be written out.
func (zfd *ZoneFileDNS) Save(cz *ConfigZone) error {
	newserial, err := IncrementSerial(cz, zfd.zone.SOA.Serial)
	if err != nil {
		return err
	}
	zfd.zone.SOA.Serial = newserial

	return zfd.zone.Save()
}
