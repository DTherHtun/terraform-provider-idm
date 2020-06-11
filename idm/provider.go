package idm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

//Provider for main func
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"idm_server": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IDM_SERVER", nil),
				Description: descriptions["idm_server"],
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IDM_USER", nil),
				Description: descriptions["user"],
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IDM_PASSWORD", nil),
				Description: descriptions["password"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				DefaultFunc: schema.EnvDefaultFunc("IDM_INSECURE", false),
				Description: descriptions["insecure"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"idm_host":       resourceIDMHost(),
			"idm_dns_record": resourceIDMDNSRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc:  providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"idm_server": "The FreeIPA host",

		"user": "Username to use for connection",

		"password": "Password to use for connection",

		"insecure": "Whether to verify the server's SSL certificate",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return &Config{
		Host:     d.Get("idm_server").(string),
		Username: d.Get("user").(string),
		Password: d.Get("password").(string),
		Insecure: d.Get("insecure").(bool),
	}, nil
}
