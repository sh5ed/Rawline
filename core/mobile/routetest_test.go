package mobile

import (
	"encoding/json"
	"testing"

	"RawlineCore/internal/routetest"
)

func TestRouteSelectsConfiguredInbound(t *testing.T) {
	const config = `{
		"inbounds": [
			{"type":"mixed","tag":"mixed-in","listen":"127.0.0.1","listen_port":2080},
			{"type":"http","tag":"http-in","listen":"127.0.0.1","listen_port":2081},
			{"type":"tun","tag":"tun-in","address":["172.19.0.1/30"]}
		],
		"outbounds": [
			{"type":"direct","tag":"proxy-egress"},
			{"type":"direct","tag":"tun-egress"},
			{"type":"direct","tag":"fallback"}
		],
		"route": {
			"rules": [
				{"inbound":["mixed-in","http-in"],"domain":["example.test"],"outbound":"proxy-egress"},
				{"inbound":["tun-in"],"domain":["example.test"],"outbound":"tun-egress"}
			],
			"final":"fallback"
		}
	}`
	for _, tc := range []struct {
		name, inbound, wantInbound, wantOutbound, wantError string
	}{
		{name: "automatic", wantInbound: "tun-in", wantOutbound: "tun-egress"},
		{name: "socks", inbound: "mixed-in", wantInbound: "mixed-in", wantOutbound: "proxy-egress"},
		{name: "http", inbound: "http-in", wantInbound: "http-in", wantOutbound: "proxy-egress"},
		{name: "tun", inbound: "tun-in", wantInbound: "tun-in", wantOutbound: "tun-egress"},
		{name: "unknown", inbound: "missing", wantError: "unknown inbound: missing"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, err := TestRoute(nil, &RouteTestRequest{Config: config, Domain: "example.test", Inbound: tc.inbound})
			if err != nil {
				t.Fatal(err)
			}
			var result routetest.Result
			if err := json.Unmarshal([]byte(data), &result); err != nil {
				t.Fatal(err)
			}
			if result.Inbound != tc.wantInbound || result.Outbound != tc.wantOutbound || result.Error != tc.wantError {
				t.Fatalf("unexpected route test response: %s", data)
			}
		})
	}
}
