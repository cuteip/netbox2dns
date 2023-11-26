package netbox2dns

import (
	"testing"
)

func TestFindConfig(t *testing.T) {
	dirs := []string{"./testdata/config1", "./testdata/config2", "./testdata/config3"}

	want1 := "./testdata/config1/conf.yaml"
	got, err := findConfig("conf", dirs, configExtensions)
	if err != nil {
		t.Fatalf("findConfig() returned an error: %v", err)
	}
	if got != want1 {
		t.Errorf("findConfig(): got %q, want %q", got, want1)

	}

	want2 := "./testdata/config2/conf.json"
	got, err = findConfig("conf", dirs, []string{"json"})
	if err != nil {
		t.Fatalf("findConfig() returned an error: %v", err)
	}
	if got != want2 {
		t.Errorf("findConfig(): got %q, want %q", got, want2)

	}

	want3 := "./testdata/config3/conf.cue"
	got, err = findConfig("conf", dirs, []string{"cue"})
	if err != nil {
		t.Fatalf("findConfig() returned an error: %v", err)
	}
	if got != want3 {
		t.Errorf("findConfig(): got %q, want %q", got, want3)

	}

	want4 := "./testdata/config2/conf.json"
	got, err = findConfig("conf", dirs, []string{"json", "cue"})
	if err != nil {
		t.Fatalf("findConfig() returned an error: %v", err)
	}
	if got != want4 {
		t.Errorf("findConfig(): got %q, want %q", got, want4)

	}
}

func TestParseYaml(t *testing.T) {
	cfg, err := ParseConfig("testdata/config4/conf.yaml")
	valueTest(t, cfg, err)
}

func TestParseJSON(t *testing.T) {
	cfg, err := ParseConfig("testdata/config4/conf.json")
	valueTest(t, cfg, err)
}

func TestParseCUE(t *testing.T) {
	cfg, err := ParseConfig("testdata/config4/conf.json")
	valueTest(t, cfg, err)
}

func valueTest(t *testing.T, cfg *Config, err error) {
	if err != nil {
		t.Fatalf("Unable to parse config: %v", err)
	}

	if cfg.Netbox.Host != "netbox.example.com" {
		t.Errorf("cfg.Netbox.Host wrong; got %q want %q", cfg.Netbox.Host, "netbox.example.com")
	}
	if cfg.Netbox.Token != "changeme" {
		t.Errorf("cfg.Netbox.Host wrong; got %q want %q", cfg.Netbox.Token, "changeme")
	}
	if len(cfg.Zones) != 4 {
		t.Errorf("len(cfg.Zones) wrong; got %d want 4", len(cfg.Zones))
	}
	if len(cfg.ZoneMap) != 4 {
		t.Errorf("len(cfg.ZoneMap) wrong; got %d want 4", len(cfg.Zones))
	}
	z := cfg.ZoneMap["0.0.0.0.ip6.arpa"]
	if z == nil {
		t.Fatalf("Failed to find zone for 0.0.0.0.ip6.arpa")
	}
	if z.Name != "0.0.0.0.ip6.arpa" {
		t.Errorf("z.Name wrong; got %q want %q", z.Name, "0.0.0.0.ip6.arpa")
	}
	if z.ZoneType != "zonefile" {
		t.Errorf("z.ZoneType wrong; got %q want %q", z.ZoneType, "zonefile")
	}
	if z.Filename != "reverse-v6-0000.zone" {
		t.Errorf("z.Filename wrong; got %q want %q", z.Filename, "reverse-6-0000.zone")
	}
	if z.TTL != 300 {
		t.Errorf("z.TTL wrong; got %d want 300", z.TTL)
	}
	if z.DeleteEntries != true {
		t.Errorf("z.DeleteEntries wrong; want true")
	}
}

func TestValidateYaml(t *testing.T) {
	_, err := ParseConfig("testdata/config5/conf1.yaml")
	if err == nil {
		t.Errorf("Should have failed validation, but succeeded.")
	}
}
