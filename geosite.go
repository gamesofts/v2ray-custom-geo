package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gogo/protobuf/proto"
	"v2ray.com/core/app/router"
)

var (
	ruleSitesDomains []*router.Domain
)

func GetSitesList(fileName string) []*router.Domain {

	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	domains := strings.Split(string(d), "\n")

	ruleSitesDomains = make([]*router.Domain, len(domains))
	for idx, pattern := range domains {
		ruleSitesDomains[idx] = &router.Domain{
			Type:  router.Domain_Domain,
			Value: pattern,
		}
	}
	return ruleSitesDomains
}

func main() {
	datfile := flag.String("dat", "geosite.dat", "datfile's name")
	sitefile := flag.String("dir", "sites", "sites file folder")

	flag.CommandLine.Usage = func() {
		fmt.Println(`./v2sitedat -dat geosite.dat -dir sites`)
	}
	flag.Parse()

	v2SiteList := new(router.GeoSiteList)

	rulefiles, err := ioutil.ReadDir(*sitefile)
	if err != nil {
		panic(err)
	}

	for _, rf := range rulefiles {
		filename := rf.Name()
		fmt.Println(filename)
		v2SiteList.Entry = append(v2SiteList.Entry, &router.GeoSite{
			CountryCode: strings.ToUpper(filename),
			Domain:      GetSitesList(*sitefile + "/" + filename),
		})
	}

	v2SiteListBytes, err := proto.Marshal(v2SiteList)
	if err != nil {
		fmt.Println("failed to marshal v2sites:", err)
		return
	}
	if err := ioutil.WriteFile(*datfile, v2SiteListBytes, 0777); err != nil {
		fmt.Println("failed to write %s.", *datfile, err)
	}

	fmt.Println(*sitefile, "-->", *datfile)
}
