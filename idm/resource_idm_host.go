package idm

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	sdk "github.com/tehwalris/go-freeipa/freeipa"
)

func resourceIDMHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceIDMHostCreate,
		Read:   resourceIDMHostRead,
		Update: resourceIDMHostUpdate,
		Delete: resourceIDMHostDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIDMHostImport,
		},

		Schema: map[string]*schema.Schema{
			"fqdn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"random": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"userpassword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"randompassword": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceIDMHostCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO][IDM] Creating Host: %s", d.Id())
	client, err := meta.(*Config).NewClient()
	if err != nil {
		return err
	}

	fqdn := d.Get("fqdn").(string)
	description := d.Get("description").(string)
	random := d.Get("random").(bool)
	userpassword := d.Get("userpassword").(string)
	force := d.Get("force").(bool)

	optArgs := sdk.HostAddOptionalArgs{
		Description: &description,
		Random:      &random,
		Force:       &force,
	}

	if userpassword != "" {
		optArgs.Userpassword = &userpassword
	}

	res, err := client.HostAdd(
		&sdk.HostAddArgs{
			Fqdn: fqdn,
		},
		&optArgs,
	)
	if err != nil {
		return err
	}

	d.SetId(fqdn)

	// randompassword is not returned by HostShow
	if d.Get("random").(bool) {
		d.Set("randompassword", *res.Result.Randompassword)
	}

	sleepDelay := 1 * time.Second
	for {
		err := resourceIDMHostRead(d, meta)
		if err == nil {
			return nil
		}
		time.Sleep(sleepDelay)
		sleepDelay = sleepDelay * 2
	}
}

func resourceIDMHostUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Host: %s", d.Id())
	client, err := meta.(*Config).NewClient()
	if err != nil {
		return err
	}

	fqdn := d.Get("fqdn").(string)
	description := d.Get("description").(string)
	random := d.Get("random").(bool)
	userpassword := d.Get("userpassword").(string)

	optArgs := sdk.HostModOptionalArgs{
		Description: &description,
		Random:      &random,
	}

	if userpassword != "" {
		optArgs.Userpassword = &userpassword
	}

	res, err := client.HostMod(
		&sdk.HostModArgs{
			Fqdn: fqdn,
		},
		&optArgs,
	)
	if err != nil {
		return err
	}

	// randompassword is not returned by HostShow
	if d.Get("random").(bool) {
		d.Set("randompassword", *res.Result.Randompassword)
	}

	return resourceIDMHostRead(d, meta)
}

func resourceIDMHostRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Host: %s", d.Id())
	client, err := meta.(*Config).NewClient()
	if err != nil {
		return err
	}

	fqdn := d.Get("fqdn").(string)

	res, err := client.HostShow(
		&sdk.HostShowArgs{
			Fqdn: fqdn,
		},
		&sdk.HostShowOptionalArgs{},
	)
	if err != nil {
		return err
	}

	if res.Result.Description != nil {
		d.Set("description", *res.Result.Description)
	}
	if res.Result.Userpassword != nil {
		d.Set("userpassword", *res.Result.Userpassword)
	}

	return nil
}

func resourceIDMHostDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Host: %s", d.Id())
	client, err := meta.(*Config).NewClient()
	if err != nil {
		return err
	}

	fqdn := d.Get("fqdn").(string)

	_, err = client.HostDel(
		&sdk.HostDelArgs{
			Fqdn: []string{fqdn},
		},
		&sdk.HostDelOptionalArgs{},
	)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceIDMHostImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.SetId(d.Id())
	d.Set("fqdn", d.Id())

	err := resourceIDMHostRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
