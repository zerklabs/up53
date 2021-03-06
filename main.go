// Updates an Amazon Route53 resource record with the public IP address of the network that this
// program runs on
package main

import (
	"bytes"
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/route53"

	"github.com/zerklabs/auburn/log"
)

var (
	zoneId     string
	recordName string
	recordType string
	ttl        int
	interval   time.Duration
	client     *route53.Route53
)

func init() {
	flag.StringVar(&zoneId, "zoneid", "", "Route53 Hosted Zone ID")
	flag.StringVar(&recordName, "name", "", "Record name to update")
	flag.StringVar(&recordType, "type", "A", "Record type to update")
	flag.IntVar(&ttl, "ttl", 3600, "TTL of record to update")
	flag.DurationVar(&interval, "interval", time.Hour*1, "How often to check and update the record")
}

func loop(i time.Duration) {
	c := time.Tick(i)

	for _ = range c {
		ip, err := getPubIp()
		if err != nil {
			log.Errorf("Error fetching public IP: %v", err)
			return
		}

		updateRecord(ip)
	}

}

func main() {
	flag.Parse()
	var err error

	if zoneId == "" {
		if os.Getenv("ZONE_ID") == "" {
			log.Error("Zone ID required")
			return
		} else {
			zoneId = os.Getenv("ZONE_ID")
		}
	}

	if recordName == "" {
		if os.Getenv("RECORD_NAME") == "" {
			log.Error("Record name required")
			return
		} else {
			recordName = os.Getenv("RECORD_NAME")
		}
	}

	if os.Getenv("RECORD_TYPE") != "" {
		recordType = os.Getenv("RECORD_TYPE")
	}

	if os.Getenv("RECORD_TTL") != "" {
		ttl, _ = strconv.Atoi(os.Getenv("RECORD_TTL"))
	}

	if os.Getenv("INTERVAL") != "" {
		interval, err = time.ParseDuration(os.Getenv("INTERVAL"))
		if err != nil {
			log.Error(err)
			return
		}
	}

	log.Debugf("Checking IP address every: %s and updating RR: %s", interval, recordName)

	auth, err := aws.EnvAuth()
	if err != nil {
		log.Error(err)
		return
	}
	client = route53.New(auth, aws.USEast)

	loop(interval)
}

func fetchRecords() {
	/**
	resp, err := client.ListResourceRecordSets(zoneId, nil)
	if err != nil {
		log.Error(err)
		return
	}
	for _, v := range resp.Records {
		for _, vv := range v.Records {
			log.Infof("%s: %s", v.Name, vv)
		}
	}
	*/

}

func getPubIp() (string, error) {
	client := &http.Client{}
	resp, err := client.Get("http://icanhazip.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	log.Debugf("Fetched public IP: %s", buf.String())

	return buf.String(), nil
}

func updateRecord(recordValue string) error {
	log.Debugf("Updating resource record: %s %s %s %s %d", zoneId, recordName, recordType, recordValue, ttl)

	resp, err := client.ChangeResourceRecordSets(zoneId, &route53.ChangeResourceRecordSetsRequest{
		Changes: []route53.Change{
			route53.Change{
				Action: "UPSERT",
				Record: route53.ResourceRecordSet{
					Name:    recordName,
					TTL:     ttl,
					Type:    recordType,
					Records: []string{recordValue},
				},
			},
		},
	})

	if err != nil {
		return err
	}

	log.Debugf("ID: %s, Status: %s, Submitted At: %s", resp.ChangeInfo.ID, resp.ChangeInfo.Status, resp.ChangeInfo.SubmittedAt)

	return nil
}
