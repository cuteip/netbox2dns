package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	log "github.com/golang/glog"
	nb "github.com/scottlaird/netbox2dns"
	"github.com/scottlaird/netbox2dns/netboxlib"
)

var (
	config = flag.String("config", "", "Path of a config file, with a .yaml, .json, or .cue extension")
)

func usage() {
	fmt.Printf("Usage: netbox2dns [--config=FILE] push\n")
	os.Exit(1)
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		usage()
	}

	switch args[0] {
	case "push":
	default:
		usage()
	}

	var err error

	// Load config file
	file := *config
	if file == "" {
		file, err = nb.FindConfig("netbox2dns")
		if err != nil {
			log.Fatal(err)
		}
	}
	cfg, err := nb.ParseConfig(file)
	if err != nil {
		log.Fatalf("Failed to parse config: %v")
	}
	log.Infof("Config read: %+v", cfg)

	ctx := context.Background()

	// Create new zones using data from Netbox
	newZones := nb.NewZones()
	for _, cz := range cfg.ZoneMap {
		newZones.NewZone(cz)
	}

	netboxClient := netboxlib.NewClient(cfg.Netbox.Host, cfg.Netbox.Token)
	addrs, err := netboxClient.GetNetboxIPAddresses(nil)
	if err != nil {
		log.Fatalf("Unable to fetch IP Addresses from Netbox: %v", err)
	}

	fmt.Printf("Found %d IP Addresses in %d zones\n", len(addrs), len(newZones.Zones))

	// Add Netbox IPs to our new zones
	err = newZones.AddAddrs(addrs)
	if err != nil {
		log.Fatalf("Unable to add IP addresses: %v", err)
	}

	log.Infof("Created %d zones", len(newZones.Zones))

	for _, zone := range newZones.Zones {
		provider, err := nb.NewDNSProvider(ctx, cfg.ZoneMap[zone.Name])
		if err != nil {
			log.Fatalf("Failed to create DNS provider for %q: %v", zone.Name, err)
		}

		for _, rec := range zone.Records {
			err = provider.WriteRecord(cfg.ZoneMap[zone.Name], rec)
			if err != nil {
				log.Errorf("Failed to update record: %v", err)
			}
		}
		err = provider.Save(cfg.ZoneMap[zone.Name])
		if err != nil {
			log.Fatalf("Failed to save: %v", err)
		}
	}
}
