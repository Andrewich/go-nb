package main

import (
	"errors"
	"log"
	"net/url"
	"os"
	"strconv"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "go-nb",
		Usage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "nb_host",
				Usage:    "Address NetBox",
				Required: true,
				EnvVars:  []string{"NETBOX_HOST"},
			},
			&cli.StringFlag{
				Name:     "nb_token",
				Usage:    "Token for NetBox",
				Required: true,
				EnvVars:  []string{"NETBOX_TOKEN"},
			},
			&cli.BoolFlag{
				Name:    "debug",
				Usage:   "Debug output",
				Value:   false,
				Aliases: []string{"d"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "prefix",
				Usage: "Commands for prefixes (add, del, list ...)",
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "list prefixes",
						Action:  list_prefixes,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "vrfid", Value: "17", Aliases: []string{"i"}},
							&cli.StringFlag{Name: "mask", Value: "32", Aliases: []string{"l"}},
							&cli.BoolFlag{Name: "plain", Usage: "Output Plain text", Aliases: []string{"p"}},
						},
					},
					{
						Name:   "add",
						Usage:  "add prefixe",
						Action: add_prefix,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "vrfid", Value: "17", Aliases: []string{"i"}},
							&cli.BoolFlag{Name: "prefix", Usage: "Value Prefix", Aliases: []string{"p"}},
						},
					},
					{
						Name:   "del",
						Usage:  "Delete prefix",
						Action: del_prefix,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "vrfid", Value: "17", Aliases: []string{"i"}},
							&cli.StringFlag{Name: "prefix", Required: true, Usage: "Value Prefix", Aliases: []string{"p"}},
						},
					},
				},
			},
			{
				Name:  "ip",
				Usage: "Commands for IP addresses (add, del, list, search ...)",
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "list IP Addresses",
						Action:  list_ip,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "vrfid", Value: "17", Aliases: []string{"i"}},
							&cli.BoolFlag{Name: "plain", Usage: "Output Plain text", Aliases: []string{"p"}},
						},
					},
					{
						Name:    "search",
						Aliases: []string{"s"},
						Usage:   "search IP Addresses",
						Action:  search_ip,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "vrfid", Value: "17", Aliases: []string{"i"}},
							&cli.StringFlag{Name: "address", Aliases: []string{"a"}},
						},
					},
					{
						Name:    "add",
						Aliases: []string{"a"},
						Usage:   "add IP Addresses",
						Action:  add_ip,
						Flags: []cli.Flag{
							&cli.Int64Flag{Name: "vrfid", Value: 17, Aliases: []string{"i"}},
							&cli.StringFlag{Name: "address", Required: true, Usage: "address with mask: eg. /32", Aliases: []string{"a"}},
							&cli.StringFlag{Name: "dns", Value: "", Aliases: []string{"d"}},
							&cli.StringFlag{Name: "description", Value: "", Aliases: []string{"c"}},
						},
					},
					{
						Name:   "delete",
						Usage:  "Delete IP Address",
						Action: del_ip,
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "vrfid", Value: "17", Aliases: []string{"i"}},
							&cli.StringFlag{Name: "address", Aliases: []string{"a"}},
						},
					},
				},
			},
			{
				Name:   "vrf",
				Usage:  "Commands for prefixes (list ...)",
				Action: list_vrfs,
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "list VRFs",
						Action:  list_vrfs,
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "plain", Usage: "Output Plain text", Aliases: []string{"p"}},
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func list_prefixes(context *cli.Context) error {
	u, err := url.Parse(context.String("nb_host"))
	if err != nil {
		return err
	}
	transport := httptransport.New(u.Host, client.DefaultBasePath, []string{u.Scheme})
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+context.String("nb_token"))

	transport.SetDebug(context.Bool("debug"))

	c := client.New(transport, nil)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Change table lines
	if context.Bool("plain") {
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
	} else {
		table.SetHeader([]string{"Prefix"})
		table.SetCenterSeparator("*")
		table.SetColumnSeparator("|")
		//table.SetRowSeparator("-")
	}

	vrf_id := context.String("vrfid")
	mask_length := context.String("mask")
	var limit int64 = 1000

	req := ipam.NewIpamPrefixesListParams()
	if vrf_id != "1" {
		req.VrfID = &vrf_id
	}

	req.MaskLength = &mask_length
	req.Limit = &limit
	pl, err := c.Ipam.IpamPrefixesList(req, nil)
	if err != nil {
		return err
	}

	total := 0
	for _, v := range pl.Payload.Results {
		total += 1
		table.Append([]string{*v.Prefix})
	}

	if !context.Bool("plain") {
		table.SetFooter([]string{strconv.Itoa(total)})
	}

	table.Render()

	return nil
}

func add_prefix(context *cli.Context) error {
	return errors.New("Unimplemented")
}

func del_prefix(context *cli.Context) error {
	u, err := url.Parse(context.String("nb_host"))
	if err != nil {
		return err
	}
	transport := httptransport.New(u.Host, client.DefaultBasePath, []string{u.Scheme})
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+context.String("nb_token"))

	transport.SetDebug(context.Bool("debug"))

	c := client.New(transport, nil)

	vrf_id := context.String("vrfid")
	prefix := context.String("prefix")

	req := ipam.NewIpamPrefixesListParams()
	if vrf_id != "1" {
		req.VrfID = &vrf_id
	}
	req.Prefix = &prefix
	pl, err := c.Ipam.IpamPrefixesList(req, nil)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.SetCenterSeparator("")
	table.SetColumnSeparator("")

	table.SetHeader([]string{"VRF ID", "VRF", "Status"})

	for _, v := range pl.Payload.Results {
		req := ipam.NewIpamPrefixesDeleteParams()
		req.ID = v.ID
		_, err := c.Ipam.IpamPrefixesDelete(req, nil)
		if err != nil {
			return err
		}

		table.Append([]string{strconv.FormatInt(v.ID, 10), *v.Prefix, "Deleted"})
	}

	table.Render()

	return nil

}

func list_ip(context *cli.Context) error {
	return errors.New("Unimplemented")
}

func search_ip(context *cli.Context) error {
	u, err := url.Parse(context.String("nb_host"))
	if err != nil {
		return err
	}
	transport := httptransport.New(u.Host, client.DefaultBasePath, []string{u.Scheme})
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+context.String("nb_token"))

	c := client.New(transport, nil)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.SetCenterSeparator("")
	table.SetColumnSeparator("")

	// Change table lines

	vrf_id := context.String("vrfid")
	ip_address := context.String("address")
	var limit int64 = 1000

	req := ipam.NewIpamIPAddressesListParams()
	req.VrfID = &vrf_id
	req.Address = &ip_address
	//req.MaskLength = &mask_length
	req.Limit = &limit
	pl, err := c.Ipam.IpamIPAddressesList(req, nil)
	if err != nil {
		return err
	}

	//total := 0
	for _, v := range pl.Payload.Results {
		//total += 1
		table.Append([]string{*v.Address})
	}

	// if !context.Bool("plain") {
	// 	table.SetFooter([]string{strconv.Itoa(total)})
	// }

	table.Render()

	return nil
}

func add_ip(context *cli.Context) error {
	u, err := url.Parse(context.String("nb_host"))
	if err != nil {
		return err
	}
	transport := httptransport.New(u.Host, client.DefaultBasePath, []string{u.Scheme})
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+context.String("nb_token"))

	transport.SetDebug(context.Bool("debug"))

	c := client.New(transport, nil)

	var vrf_id *int64
	if context.Int64("vrfid") == 0 {
		vrf_id = nil
	} else {
		v := context.Int64("vrfid")
		vrf_id = &v
	}

	ip_address := context.String("address")
	dns_name := context.String("dns")
	description := context.String("description")

	data := &models.WritableIPAddress{
		Vrf:         vrf_id,
		Address:     &ip_address,
		Status:      "active",
		Description: description,
		DNSName:     dns_name,
		Tags:        []*models.NestedTag{},
	}

	req := ipam.NewIpamIPAddressesCreateParams().WithData(data)

	pl, err := c.Ipam.IpamIPAddressesCreate(req, nil)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.SetCenterSeparator("")
	table.SetColumnSeparator("")

	table.SetHeader([]string{"VRF", "Address", "DNS", "Status"})

	table.Append([]string{*pl.Payload.Vrf.Name, *pl.Payload.Address, pl.Payload.DNSName, "Created"})

	table.Render()

	return nil
}

func del_ip(context *cli.Context) error {
	return errors.New("Unimplemented")
}

func list_vrfs(context *cli.Context) error {
	u, err := url.Parse(context.String("nb_host"))
	if err != nil {
		return err
	}
	transport := httptransport.New(u.Host, client.DefaultBasePath, []string{u.Scheme})
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+context.String("nb_token"))

	transport.SetDebug(context.Bool("debug"))

	c := client.New(transport, nil)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "VRF"})

	vrfs, err := c.Ipam.IpamVrfsList(nil, nil)
	if err != nil {
		return err
	}

	for _, v := range vrfs.Payload.Results {
		table.Append([]string{strconv.FormatInt(v.ID, 10), *v.Name})
	}

	table.Render()

	return nil
}
