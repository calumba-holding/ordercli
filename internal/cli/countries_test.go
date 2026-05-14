package cli

import (
	"strings"
	"testing"
)

func TestFindPreset(t *testing.T) {
	for _, tc := range []countryPreset{
		{Code: "CZ", BaseURL: "https://cz.fd-api.com/api/v5/", GlobalEntityID: "DJ_CZ", TargetISO: "CZ"},
		{Code: "SE", BaseURL: "https://se.fd-api.com/api/v5/", GlobalEntityID: "OP_SE", TargetISO: "SE"},
	} {
		t.Run(tc.Code, func(t *testing.T) {
			p, ok := findPreset(tc.Code)
			if !ok {
				t.Fatalf("expected %s preset", tc.Code)
			}
			if p.BaseURL != tc.BaseURL {
				t.Fatalf("BaseURL=%q", p.BaseURL)
			}
			if p.GlobalEntityID != tc.GlobalEntityID {
				t.Fatalf("GlobalEntityID=%q", p.GlobalEntityID)
			}
			if p.TargetISO != tc.TargetISO {
				t.Fatalf("TargetISO=%q", p.TargetISO)
			}
		})
	}
}

func TestCountriesCmd_IncludesBundledPresets(t *testing.T) {
	out, _, err := runCLI(t.TempDir()+"/config.json", []string{"foodora", "countries"}, "")
	if err != nil {
		t.Fatalf("countries: %v", err)
	}
	for _, want := range []string{
		"CZ\tDJ_CZ\thttps://cz.fd-api.com/api/v5/",
		"SE\tOP_SE\thttps://se.fd-api.com/api/v5/",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %s preset in output:\n%s", want[:2], out)
		}
	}
}

func TestConfigSetCountry(t *testing.T) {
	for _, tc := range []struct {
		code string
		want []string
	}{
		{
			code: "CZ",
			want: []string{
				"base_url=https://cz.fd-api.com/api/v5/",
				"global_entity_id=DJ_CZ",
				"target_country_iso=CZ",
			},
		},
		{
			code: "SE",
			want: []string{
				"base_url=https://se.fd-api.com/api/v5/",
				"global_entity_id=OP_SE",
				"target_country_iso=SE",
			},
		},
	} {
		t.Run(tc.code, func(t *testing.T) {
			cfgPath := t.TempDir() + "/config.json"
			if _, _, err := runCLI(cfgPath, []string{"foodora", "config", "set", "--country", tc.code}, ""); err != nil {
				t.Fatalf("config set --country %s: %v", tc.code, err)
			}
			out, _, err := runCLI(cfgPath, []string{"foodora", "config", "show"}, "")
			if err != nil {
				t.Fatalf("config show: %v", err)
			}
			for _, want := range tc.want {
				if !strings.Contains(out, want) {
					t.Fatalf("missing %q in config output:\n%s", want, out)
				}
			}
		})
	}
}
