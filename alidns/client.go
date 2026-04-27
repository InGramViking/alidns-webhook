package alidns

import (
	"fmt"

	alidnsclient "github.com/alibabacloud-go/alidns-20150109/v5/client"
	openapiutil "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
)

type Client struct {
	dnsc *alidnsclient.Client
}

func newClient(region, accessKey, secretKey string) (*Client, error) {
	config := &openapiutil.Config{
		AccessKeyId:     &accessKey,
		AccessKeySecret: &secretKey,
		RegionId:        &region,
		Type:            strPtr("access_key"),
	}
	client, err := alidnsclient.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Client{dnsc: client}, nil
}

func strPtr(s string) *string {
	return &s
}

func (c *Client) getHostedZone(zone string) (string, error) {
	request := &alidnsclient.DescribeDomainsRequest{}
	request.SetKeyWord(util.UnFqdn(zone))
	request.SetSearchMode("EXACT")

	response, err := c.dnsc.DescribeDomains(request)
	if err != nil {
		return "", err
	}

	domains := response.Body.Domains.Domain
	if len(domains) == 0 {
		return "", fmt.Errorf("zone %s does not exist", zone)
	}

	return *domains[0].DomainName, nil
}

func (c *Client) addTxtRecord(zone, rr, value string) error {
	request := &alidnsclient.AddDomainRecordRequest{}
	request.SetDomainName(zone)
	request.SetRR(rr)
	request.SetType("TXT")
	request.SetValue(value)

	_, err := c.dnsc.AddDomainRecord(request)
	return err
}

func (c *Client) getTxtRecord(zone, rr string) (*alidnsclient.DescribeDomainRecordsResponseBodyDomainRecordsRecord, error) {
	request := &alidnsclient.DescribeDomainRecordsRequest{}
	request.SetDomainName(zone)
	request.SetType("TXT")
	request.SetRRKeyWord(rr)

	response, err := c.dnsc.DescribeDomainRecords(request)
	if err != nil {
		return nil, err
	}

	records := response.Body.DomainRecords.Record
	for _, r := range records {
		if *r.RR == rr {
			return r, nil
		}
	}

	return nil, fmt.Errorf("txt record does not exist: %v.%v", rr, zone)
}

func (c *Client) deleteDomainRecord(id string) error {
	request := &alidnsclient.DeleteDomainRecordRequest{}
	request.SetRecordId(id)

	_, err := c.dnsc.DeleteDomainRecord(request)
	return err
}
