//go:build with_tailscale

package boxmain

import (
	"testing"

	"github.com/sagernet/sing-box/option"
)

func TestTailscaleConfigUsesManagedDirectories(t *testing.T) {
	raw := []byte(`{"endpoints":[{"type":"tailscale","tag":"exit","state_directory":"/custom/state","taildrop_directory":"/custom/downloads"}]}`)
	options, err := parseConfig(newBoxContext(), raw)
	if err != nil {
		t.Fatal(err)
	}
	settings := options.Endpoints[0].Options.(*option.TailscaleEndpointOptions)
	if settings.StateDirectory != "tailscale" || settings.TaildropDirectory != "Taildrop" {
		t.Fatalf("unexpected Tailscale paths: state %q, Taildrop %q", settings.StateDirectory, settings.TaildropDirectory)
	}
}
