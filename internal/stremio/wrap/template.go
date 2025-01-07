package stremio_wrap

import (
	"github.com/MunifTanjim/stremthru/internal/stremio/configure"
)

func getTemplateData() *configure.TemplateData {
	return &configure.TemplateData{
		Base: configure.Base{
			Title:       "StremThru Wrap",
			Description: "Stremio Addon to Wrap another Addon with StremThru",
			NavTitle:    "Wrap",
		},
		Configs: []configure.Config{
			configure.Config{
				Key:         "manifest_url",
				Type:        "url",
				Default:     "",
				Title:       "Upstream Manifest URL",
				Description: "Manifest URL for the Upstream Addon",
				Required:    true,
				Action: configure.ConfigAction{
					Label:   "Configure",
					OnClick: "onUpstreamManifestConfigure()",
				},
			},
			getStoreNameConfig(),
			configure.Config{
				Key:         "token",
				Type:        "password",
				Default:     "",
				Title:       "Store Token",
				Description: "",
				Required:    true,
			},
			configure.Config{
				Key:     "cached",
				Type:    configure.ConfigTypeCheckbox,
				Title:   "Only Show Cached Content",
				Options: []configure.ConfigOption{},
			},
		},
		Script: configure.GetScriptStoreTokenDescription("store", "token") + `
function onUpstreamManifestConfigure() {
  const url = document.querySelector("input[name='manifest_url']").value.replace(/\/manifest.json$/,'') + "/configure";
  window.open(url, "_blank");
}
`,
	}
}
