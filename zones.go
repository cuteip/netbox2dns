package netbox2dns

import (
	"fmt"
	"net/netip"
	"sort"
	"strings"

	log "github.com/golang/glog"
	"github.com/scottlaird/netbox2dns/netboxlib"
)

// ByLength is a wrapper for []string for sorting the string
// slice by length, from longest to shortest.
type ByLength []*Zone

func (a ByLength) Len() int {
	return len(a)
}
func (a ByLength) Less(i, j int) bool {
	return len(a[i].Name) > len(a[j].Name)
}
func (a ByLength) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Zones represents the set of all DNS zones known to netbox2dns.
type Zones struct {
	Zones       map[string]*Zone
	sortedZones []*Zone
}

// NewZones creates a new Zones structure and initializes it.
func NewZones() *Zones {
	return &Zones{
		Zones: make(map[string]*Zone),
	}
}

// AddRecord adds a record to the appropriate zone.  It finds the
// longest suffix match among all known zones and adds the new record
// there.  If no zones match, then an error is returned.
func (z *Zones) AddRecord(r *Record) error {
	for _, zone := range z.sortedZones {
		if strings.HasSuffix(r.Name, zone.Name+".") {
			zone.AddRecord(r)
			return nil
		}
	}
	return fmt.Errorf("Can't find zone matching record %q in %v", r.Name, z.sortedZones)
}

// AddZone adds a new Zone to Zones.
func (z *Zones) AddZone(zone *Zone) {
	z.Zones[zone.Name] = zone
	z.sortZones()
}

// NewZone creates a new Zone in Zones using the settings in the
// provided ConfigZone.  The resulting Zone is added to Zones
// automatically.
func (z *Zones) NewZone(cz *ConfigZone) {
	zone := Zone{
		Name:     cz.Name,
		Filename: cz.Filename,
		TTL:      cz.TTL,
		Records:  make(map[string][]*Record),
	}
	z.AddZone(&zone)
}

// sortZones sorts zones from longest to shortest and populates `sortedZones`.
func (z *Zones) sortZones() {
	zones := make([]*Zone, len(z.Zones))
	i := 0
	for _, zone := range z.Zones {
		zones[i] = zone
		i++
	}
	sort.Sort(ByLength(zones))

	z.sortedZones = zones
}

// Zone represents a single DNS zone on a single provider (fixed zone files, etc).
type Zone struct {
	Name     string
	Filename string
	TTL      int64
	Records  map[string][]*Record
}

// AddRecord adds a single record to this zone.  It does not check
// that this is the correct zone for the record.
func (z *Zone) AddRecord(r *Record) {
	if r.TTL == 0 {
		r.TTL = z.TTL
	}
	z.Records[r.Name] = append(z.Records[r.Name], r)
}

// ReverseName takes an IP address and returns the correct reverse DNS
// name for that IP.  It maps IPv4 addresses into `in-addr.arpa` and
// IPv6 addresses into `ip6.arpa`.
func ReverseName(addr netip.Addr) string {
	if addr.Is4() {
		return reverseName4(addr)
	}
	return reverseName6(addr)
}

func reverseName4(addr netip.Addr) string {
	b := addr.As4()
	return fmt.Sprintf("%d.%d.%d.%d.in-addr.arpa.", b[3], b[2], b[1], b[0])
}

func reverseName6(addr netip.Addr) string {
	ret := ""
	b := addr.As16()
	for i := 15; i >= 0; i-- {
		ret += fmt.Sprintf("%x.%x.", b[i]&0xf, (b[i]&0xf0)>>4)
	}
	return ret + "ip6.arpa."
}

// AddAddrs adds multiple addresses to a set of Zones.  This creates
// both forward and reverse DNS entries.
func (z *Zones) AddAddrs(addrs []netboxlib.IpamIPAddress) error {
	for _, addr := range addrs {
		if addr.DNSName != "" && addr.Status == "active" {
			forward := Record{
				Name:    addr.DNSName + ".",
				Rrdatas: []string{addr.Address.String()},
			}
			reverse := Record{
				Name:    ReverseName(addr.Address),
				Type:    "PTR",
				Rrdatas: []string{addr.DNSName + "."},
			}
			if addr.Address.Is4() {
				forward.Type = "A"
			} else {
				forward.Type = "AAAA"
			}

			err := z.AddRecord(&forward)
			if err != nil {
				log.Warningf("Unable to add forward record: %v", err)
			}
			err = z.AddRecord(&reverse)
			if err != nil {
				log.Warningf("Unable to add reverse record: %v", err)
			}
		}
	}
	return nil
}
