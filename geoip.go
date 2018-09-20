// GeoIP generator
//
// Before running this file, the GeoIP database must be downloaded and present.
// To download GeoIP database: https://dev.maxmind.com/geoip/geoip2/geolite2/
// Inside you will find block files for IPv4 and IPv6 and country code mapping.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gogo/protobuf/proto"
	"v2ray.com/core/app/router"
	"v2ray.com/ext/tools/conf"
	"v2ray.com/core/common"
)

func getIPsList(fileName string) []*router.CIDR {
	

	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	ips := strings.Split(string(d), "\n")

	cidr := make([]*router.CIDR, 0, len(ips))

	for _, ip := range ips {
		c, err := conf.ParseIP(ip)
		common.Must(err)
		cidr = append(cidr, c)
	}
	return cidr
}

func main() {
	datfile := flag.String("dat", "geoip.dat", "datfile's name")
	ipfile := flag.String("dir", "ips", "ips file folder")

	flag.CommandLine.Usage = func() {
		fmt.Println(`./v2ipdat -dat geoip.dat -dir ips`)
	}
	flag.Parse()

	geoIPList := new(router.GeoIPList)

	rulefiles, err := ioutil.ReadDir(*ipfile)
	if err != nil {
		panic(err)
	}

	for _, rf := range rulefiles {
		filename := rf.Name()
		fmt.Println(filename)
		geoIPList.Entry = append(geoIPList.Entry, &router.GeoIP{
			CountryCode: strings.ToUpper(filename),
			Cidr:      getIPsList(*ipfile + "/" + filename),
		})
	}

	geoIPBytes, err := proto.Marshal(geoIPList)
	if err != nil {
		fmt.Println("Error marshalling geoip list:", err)
	}

	if err := ioutil.WriteFile("geoip.dat", geoIPBytes, 0777); err != nil {
		fmt.Println("Error writing geoip to file:", err)
	}

	fmt.Println(*ipfile, "-->", *datfile)
}

