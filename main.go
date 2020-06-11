package main

import (
	idm "github.com/DTherHtun/terraform-provider-idm/idm"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: idm.Provider,
	})
}
