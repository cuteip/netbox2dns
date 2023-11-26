config: {
	netbox: {
		host:  "netbox.example.com"
		token: "changeme"
	}

	defaults: {
		ttl:     300
	}

	zones: [
		{name: "internal.example.com"
			filename: "internal-example-com.zone"
 			zonetype: "zonefile"
		},
		{name: "example.com"
			filename: "example-com.zone"
 			zonetype: "zonefile"
		},
		{name: "10.in-addr.arpa"
			filename:       "reverse-v4-10.zone"
			delete_entries: true
 			zonetype: "zonefile"
		},
		{name: "0.0.0.0.ip6.arpa"
			filename:       "reverse-v6-0000.zone"
			delete_entries: true
 			zonetype: "zonefile"
		},
	]
}
