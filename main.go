package main

import (	
	"os"
	"strconv"
	"log"
	
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{		
		Name: "go-nb",
		Usage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
			  Name:    "nb_host",
			  //Aliases: []string{"h"},			  
			  Usage:   "Address NetBox",
			  EnvVars: []string{"NETBOX_HOST"},
			},
			&cli.StringFlag{
				Name:    "nb_token",
				//Aliases: []string{"t"},			  
				Usage:   "Token for NetBox",
				EnvVars: []string{"NETBOX_TOKEN"},
			  },
		  },
		Commands: []*cli.Command{
			{
			  Name:    "prefixes",
			  Aliases: []string{"p"},
			  Usage:   "list prefixes /32",
			  Action:  list_prefixes,
			  Flags: []cli.Flag{				
				&cli.StringFlag{Name: "vrfid", Value: "17", Aliases: []string{"i"}},
				&cli.StringFlag{Name: "mask_length", Value: "32", Aliases: []string{"l"}},
			  },
			},
			{
				Name:    "vrfs",
				Aliases: []string{"v"},
				Usage:   "list VRFs",
				Action:  list_vrfs,				
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func list_prefixes(context *cli.Context) error {
	transport := httptransport.New(context.String("nb_host"), client.DefaultBasePath, []string{"https"})
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+context.String("nb_token"))

	c := client.New(transport, nil)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Prefix"})

	vrf_id := context.String("vrfid")
	mask_length := context.String("mask_length")
	var limit int64 = 1000

	req := ipam.NewIpamPrefixesListParams()
	req.VrfID = &vrf_id
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

	table.SetFooter([]string{strconv.Itoa(total)})
	table.Render()

	return nil
}

func list_vrfs(context *cli.Context) error {
	transport := httptransport.New(context.String("nb_host"), client.DefaultBasePath, []string{"https"})
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+context.String("nb_token"))

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
