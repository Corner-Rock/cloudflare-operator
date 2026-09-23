package cf

import "testing"

func TestDnsUpToDate(t *testing.T) {
	api := &API{ValidTunnelId: "tunnel-1"}
	target := api.TunnelDomain()
	if target != "tunnel-1.cfargotunnel.com" {
		t.Fatalf("unexpected tunnel domain %q", target)
	}

	managed := DnsManagedRecordTxt{DnsId: "cname-1", TunnelId: "tunnel-1", TunnelName: "k3s"}
	current := DnsCNameRecord{Id: "cname-1", Content: target, Proxied: true}

	cases := []struct {
		name     string
		existing DnsCNameRecord
		txtId    string
		txt      DnsManagedRecordTxt
		want     bool
	}{
		{"records match", current, "txt-1", managed, true},
		{"no cname", DnsCNameRecord{}, "txt-1", managed, false},
		{"no txt marker", current, "", DnsManagedRecordTxt{}, false},
		{"txt points at another cname", current, "txt-1", DnsManagedRecordTxt{DnsId: "cname-old", TunnelId: "tunnel-1"}, false},
		{"txt owned by another tunnel", current, "txt-1", DnsManagedRecordTxt{DnsId: "cname-1", TunnelId: "tunnel-2"}, false},
		{"cname targets another tunnel", DnsCNameRecord{Id: "cname-1", Content: "tunnel-2.cfargotunnel.com", Proxied: true}, "txt-1", managed, false},
		{"cname not proxied", DnsCNameRecord{Id: "cname-1", Content: target, Proxied: false}, "txt-1", managed, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := api.DnsUpToDate(tc.existing, tc.txtId, tc.txt); got != tc.want {
				t.Errorf("DnsUpToDate() = %v, want %v", got, tc.want)
			}
		})
	}
}
