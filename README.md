# netbox2dns

original: https://github.com/scottlaird/netbox2dns

netbox2dns is a tool for publishing DNS records from [NetBox](http://netbox.dev) data.

NetBox provides a reasonable interface for managing and documenting IP
addresses and network devices, but out of the box there's no good way
to publish NetBox's data into DNS.  This tool is designed to publish
A, AAAA, and PTR records from NetBox into zonefile.  It should
be possible to add other DNS providers without too much work, as long
as they're able to handle incremental record additions and removals.

## fork 元との差分

- netbox2dns にて forward zone が管理されていない場合でもエラーにしない
- ゾーンファイルが存在しない場合にエラーにせず空ファイルを作成する
- Tag `netbox2dns_exclude` が付与された IP アドレスを除外する（ゾーンファイルに反映しない）
- 求められる機能を減らし、NetBox の情報をゾーンファイルに書き込むことだけに
  - `push` オプションのみにして、設定ファイルにて指定されるゾーンファイルは netbox2dns によって管理されている前提の構成に

## Compiling

Check out a copy of the `netbox2dns` code from GitHub using `git clone
https://github.com/cuteip/netbox2dns.git`.  Then, run `go build
cmd/netbox2dns/netbox2dns.go`, and it should generate a `netbox2dns`
binary.  This can be copied to other directories or other systems as
needed.

## Configuration

Edit `netbox2dns.yaml`.  Here is an example config:

```yaml
config:
  netbox:
    host:  "netbox.example.com"
    token: "01234567890abcdef"

  defaults:
    ttl: 300

  zones:
    - name: "internal.example.com"
      zonetype: "zonefile"
      filename: "/etc/dns/internal.example.com.zone"
    - name: "example.com"
      zonetype: "zonefile"
      filename: "/etc/dns/example.com.zone"
    - name: "10.in-addr.arpa"
      zonetype: "zonefile"
      filename: "/etc/dns/10.in-addr.arpa.zone"
    - name: "0.0.0.0.ip6.arpa"
      zonetype: "zonefile"
      filename: "/etc/dns/0.0.0.0.ip6.arpa.zone"
```

Each zone needs to specify a name and a zonetype.  Currently supported
zonetype is `zonefile` for text
zone files.  See `config.cue` for an authoratative list of parameters
per zone.

To talk to NetBox, you'll need to provide your NetBox host, a NetBox
API token with (at a minimum) read access to NetBox's IP Address data.

Finally, list your zones. When adding new records, netbox2dns will add
records to the *longest* matching zone name.  For the example above,
with `internal.example.com` and `example.com`, if NetBox has a record
for `router1.internal.example.com`, then it will be added to
`internal.example.com`.  Any records that don't fix into a listed zone
will be ignored.

By default, netbox2dns will search in `/etc/netbox2dns/`,
`/usr/local/etc/netbox2dns/`, and the correct directory for its config
file.  Config files can be in YAML (shown above), JSON, or CUE format.
Examples in [all 3
formats](https://github.com/scottlaird/netbox2dns/tree/main/testdata/config4)
are available.

## Use

Short version: create a configuration file (see previous section),
then run `netbox2dns push`.

Upon startup, netbox2dns will fetch all IP Address records from NetBox
*and* all A/AAAA/PTR records from the listed zones.  netbox2dns
ignores other record types, including SOA, NS, and CNAME.

For each active IP address in NetBox that has a DNS name, netbox2dns
will try to add both forward and reverse DNS records.  Both IPv4
and IPv6 should be handled automatically.

When you run `netbox2dns push`, netbox2dns will generate zone files.
At that time, the contents already written in the zone file will be deleted.
Generated zone files are expected to be included by `$INCLUDE` directive.
