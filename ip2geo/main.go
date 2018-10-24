//go:generate go-bindata-assetfs -debug -o assets.go assets/...
// Adding -debug before -o above and run `go generate`, go-bindata-assetfs doesn't embed the asset files inside the executable but access
// them from filessytem at runtime. This allows for faster builds during development.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/ww9/misc/ip2geo/ip2location"
)

func main() {
	var err error
	c := &Controller{}
	c.DbIP, err = ip2location.OpenFromBytes(MustAsset("assets/IP2LOCATION-LITE-DB11.BIN"))
	if err != nil {
		log.Fatal(err)
	}
	defer c.DbIP.Close()

	http.HandleFunc("/ip2geo", func(w http.ResponseWriter, req *http.Request) {
		ip := req.URL.Query().Get("i")
		location := c.SearchIP(ip)
		locationJSON, err := json.Marshal(location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(locationJSON)
	})
	http.Handle("/", http.FileServer(assetFS()))
	fmt.Println("Open http://127.0.0.1:8081 in browser")
	log.Fatalln(http.ListenAndServe("127.0.0.1:8081", nil))
}

func searchIP(dbIP *ip2location.Database) {
	record := dbIP.Get_all("8.8.8.8")
	fmt.Printf("country_short: %s\n", record.CountryShort)
	fmt.Printf("country_long: %s\n", record.CountryLong)
	fmt.Printf("region: %s\n", record.Region)
	fmt.Printf("city: %s\n", record.City)
	fmt.Printf("latitude: %f\n", record.Latitude)
	fmt.Printf("longitude: %f\n", record.Longitude)
	fmt.Printf("zipcode: %s\n", record.Zipcode)
	fmt.Printf("timezone: %s\n", record.Timezone)
	fmt.Printf("api version: %s\n", ip2location.APIVersion)
}

// Controller is a simple example of automatic Go-to-JS data binding
type Controller struct {
	DbIP *ip2location.Database
	lock sync.Mutex
}

// SearchIP queries database for an IP and returns its geolocation informaton
func (c *Controller) SearchIP(ip string) ip2location.Record {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.DbIP.Get_all(ip)
}
