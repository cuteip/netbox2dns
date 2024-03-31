package zonefile

// "github.com/shuLhan/share/lib/dns" をベースに netbox2dns に必要なもののみに絞る

import (
	"fmt"
	"os"
)

type Zone struct {
	File            *os.File
	ResourceRecords []ResourceRecord
}

type ResourceRecord struct {
	Name  string
	Type  string
	Class string
	TTL   uint32
	Rdata []string
}

func New(filename string) (*Zone, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return &Zone{File: f, ResourceRecords: []ResourceRecord{}}, nil
}

func (z *Zone) Add(r ResourceRecord) error {
	z.ResourceRecords = append(z.ResourceRecords, r)
	return nil
}

func (z *Zone) Save() error {
	str := ""
	for _, rr := range z.ResourceRecords {
		for _, rd := range rr.Rdata {
			str += fmt.Sprintf("%s %d %s %s %s\n", rr.Name, rr.TTL, rr.Class, rr.Type, rd)
		}
	}

	_, err := z.File.WriteString(str)
	if err != nil {
		return err
	}
	return z.File.Close()
}
