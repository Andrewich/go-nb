package main

import (
	"fmt"
	"os"
	"strconv"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/olekukonko/tablewriter"
)

func main() {
	token := os.Getenv("NETBOX_TOKEN")
	if token == "" {
		log.Fatalf("Please provide netbox API token via env var NETBOX_TOKEN")
	}

	netboxHost := os.Getenv("NETBOX_HOST")
	if netboxHost == "" {
		log.Fatalf("Please provide netbox host via env var NETBOX_HOST")
	}

	transport := httptransport.New(netboxHost, client.DefaultBasePath, []string{"https"})
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+token)

	c := client.New(transport, nil)

	vrf, err := list_vrf(c)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "IP addresses", "/32", "/31", "/30", "/29", "/28", "/27", "/26", "/25", "/24", "/23"})

	for _, v := range vrf.Payload.Results {

		lp32, _ := list_prefixes(c, v.ID, 32)
		lp31, _ := list_prefixes(c, v.ID, 31)
		lp30, _ := list_prefixes(c, v.ID, 30)
		lp29, _ := list_prefixes(c, v.ID, 29)
		lp28, _ := list_prefixes(c, v.ID, 28)
		lp27, _ := list_prefixes(c, v.ID, 27)
		lp26, _ := list_prefixes(c, v.ID, 26)
		lp25, _ := list_prefixes(c, v.ID, 25)
		lp24, _ := list_prefixes(c, v.ID, 24)
		lp23, _ := list_prefixes(c, v.ID, 23)

		table.Append([]string{strconv.FormatInt(v.ID, 10), *v.Name, strconv.FormatInt(v.IpaddressCount, 10), lp32, lp31, lp30, lp29, lp28, lp27, lp26, lp25, lp24, lp23})
		//fmt.Printf("%d) ID: %d\tName: %v\t\tIP addresses: %d\tPrefixes /32: %d\n", i, v.ID, *v.Name, v.IpaddressCount, *(lp.Payload.Count))
	}

	//table.SetFooter([]string{"", "", "Total /32 Prefixes", strconv.FormatInt(count_prefixes, 10)})
	//table.SetFooter([]string{"", "", "Total Prefixes", strconv.FormatInt(*(pl.Payload.Count), 10)})
	table.Render()
}

func list_vrf(c *client.NetBoxAPI) (*ipam.IpamVrfsListOK, error) {
	req := ipam.NewIpamVrfsListParams()
	pl, err := c.Ipam.IpamVrfsList(req, nil)
	if err != nil {
		return nil, err
	}
	return pl, nil
}

func list_prefixes(c *client.NetBoxAPI, vrf_id int64, mask int64) (string, error) {
	vrf_id_s := strconv.FormatInt(vrf_id, 10)
	mask_length := strconv.FormatInt(mask, 10)
	req := ipam.NewIpamPrefixesListParams()
	req.VrfID = &(vrf_id_s)
	req.MaskLength = &(mask_length)
	pl, err := c.Ipam.IpamPrefixesList(req, nil)
	if err != nil {
		return "", err
	}

	// n, err := strconv.ParseInt(*(pl.Payload.Count), 10, 64)
	// if err != nil {
	// 	return "", err
	// }

	return strconv.FormatInt(*(pl.Payload.Count), 10), nil

	//return pl, nil
}
