package boxmain

import (
	"testing"

	"RawlineCore/internal/torproc"
	"github.com/sagernet/sing-box/option"
)

func TestTorConfigUsesManagedPaths(t *testing.T) {
	raw := []byte(`{"outbounds":[{"type":"tor","tag":"exit","executable_path":"custom/tor","data_directory":"custom/data","extra_args":["--Quiet"]}]}`)
	options, err := parseConfig(newBoxContext(), raw)
	if err != nil {
		t.Fatal(err)
	}
	settings := options.Outbounds[0].Options.(*option.TorOutboundOptions)
	if settings.ExecutablePath != torproc.ExecutablePath() || settings.DataDirectory != torproc.DataDirectoryPath() || len(settings.ExtraArgs) != 0 {
		t.Fatalf("unexpected Tor settings: %+v", settings)
	}
}
