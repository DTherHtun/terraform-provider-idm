package idm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	sdk "github.com/tehwalris/go-freeipa/freeipa"
)

func resourceIDMDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceIDMDNSRecordCreate,
		Read:   resourceIDMDNSRecordRead,
		Update: resourceIDMDNSRecordUpdate,
		Delete: resourceIDMDNSRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIDMDNSRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"idm_dns_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"idm_dns_zone_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"records": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				// Set:      schema.HashString,
			},
			"dnsttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dnsclass": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIDMDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO][Redhat IDM] Creating DNS Record: %s", d.Id())

	client, err := meta.(*Config).NewClient()
	if err != nil {
		return err
	}

	idmdnsname := d.Get("idm_dns_name").(string)
	idmdnszonename := d.Get("idm_dns_zone_name").(string)

	args := sdk.DnsrecordAddArgs{
		Idnsname: idmdnsname,
	}

	optArgs := sdk.DnsrecordAddOptionalArgs{
		Dnszoneidnsname: &idmdnszonename,
	}

	_type := d.Get("type")
	_records := d.Get("records").(*schema.Set).List()
	records := make([]string, len(_records))
	for i, d := range _records {
		records[i] = d.(string)
	}
	switch _type {
	case "A":
		optArgs.Arecord = &records
	case "SRV":
		optArgs.Srvrecord = &records
	}

	if _dnsttl, ok := d.GetOkExists("dnsttl"); ok {
		dnsttl := _dnsttl.(int)
		optArgs.Dnsttl = &dnsttl
	}

	if _dnsclass, ok := d.GetOkExists("dnsclass"); ok {
		dnsclass := _dnsclass.(string)
		optArgs.Dnsclass = &dnsclass
	}

	_, err = client.DnsrecordAdd(&args, &optArgs)
	if err != nil {
		return err
	}

	// TODO: use aws_route53_records' way to generate ID
	d.SetId(fmt.Sprintf("%s.%s", idmdnsname, idmdnszonename))

	return resourceIDMDNSRecordRead(d, meta)
}

func resourceIDMDNSRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating DNS Record: %s", d.Id())

	client, err := meta.(*Config).NewClient()
	if err != nil {
		return err
	}

	args := sdk.DnsrecordModArgs{
		Idnsname: d.Get("idm_dns_name").(string),
	}

	idmdnszonename := d.Get("idm_dns_zone_name").(string)
	optArgs := sdk.DnsrecordModOptionalArgs{
		Dnszoneidnsname: &idmdnszonename,
	}

	_type := d.Get("type")
	_records := d.Get("records").(*schema.Set).List()
	records := make([]string, len(_records))
	for i, d := range _records {
		records[i] = d.(string)
	}
	switch _type {
	case "A":
		optArgs.Arecord = &records
	case "SRV":
		optArgs.Srvrecord = &records
	}

	if _dnsttl, ok := d.GetOkExists("dnsttl"); ok {
		dnsttl := _dnsttl.(int)
		optArgs.Dnsttl = &dnsttl
	}

	if _dnsclass, ok := d.GetOkExists("dnsclass"); ok {
		dnsclass := _dnsclass.(string)
		optArgs.Dnsclass = &dnsclass
	}

	_, err = client.DnsrecordMod(&args, &optArgs)
	if err != nil {
		return err
	}

	return resourceIDMDNSRecordRead(d, meta)
}

func resourceIDMDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing DNS Record: %s", d.Id())

	client, err := meta.(*Config).NewClient()
	if err != nil {
		return err
	}

	args := sdk.DnsrecordShowArgs{
		Idnsname: d.Get("idm_dns_name").(string),
	}

	idmdnszonename := d.Get("idm_dns_name").(string)
	all := true
	optArgs := sdk.DnsrecordShowOptionalArgs{
		Dnszoneidnsname: &idmdnszonename,
		All:             &all,
	}

	res, err := client.DnsrecordShow(&args, &optArgs)
	if err != nil {
		return err
	}

	if res.Result.Arecord != nil {
		d.Set("records", *res.Result.Arecord)
	}

	if res.Result.Srvrecord != nil {
		d.Set("records", *res.Result.Srvrecord)
	}

	if res.Result.Dnsttl != nil {
		d.Set("dnsttl", *res.Result.Dnsttl)
	}

	if res.Result.Dnsclass != nil {
		d.Set("dnsclass", *res.Result.Dnsclass)
	}

	return nil
}

func resourceIDMDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting DNS Record: %s", d.Id())

	client, err := meta.(*Config).NewClient()
	if err != nil {
		return err
	}

	args := sdk.DnsrecordDelArgs{
		Idnsname: d.Get("idm_dns_name").(string),
	}

	idmdnszonename := d.Get("idm_dns_zone_name").(string)
	optArgs := sdk.DnsrecordDelOptionalArgs{
		Dnszoneidnsname: &idmdnszonename,
	}

	_type := d.Get("type")
	_records := d.Get("records").(*schema.Set).List()
	records := make([]string, len(_records))
	for i, d := range _records {
		records[i] = d.(string)
	}
	switch _type {
	case "A":
		optArgs.Arecord = &records
	case "SRV":
		optArgs.Srvrecord = &records
	}

	_, err = client.DnsrecordDel(&args, &optArgs)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceIDMDNSRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[INFO] Importing DNS Record: %s", d.Id())

	d.SetId(d.Id())
	d.Set("idm_dns_name", d.Id())

	err := resourceIDMDNSRecordRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
