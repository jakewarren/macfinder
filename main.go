package main

//Download latest OUI db from Wireshark
//go:generate wget -O oui.txt "https://code.wireshark.org/review/gitweb?p=wireshark.git;a=blob_plain;f=manuf"

import (
	"fmt"
	"os"
	"path"
	"runtime"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/jakewarren/go-ouitools"
)

var db *ouidb.OuiDb
var (
	app = kingpin.New("macfinder", "Look-up the manufacturer for a MAC address")
	mac = app.Arg("mac", "the MAC address to lookup.").Required().String()
)

func main() {
	app.Version("0.1").VersionFlag.Short('V')
	app.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.SeparateOptionalFlagsUsageTemplate)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	log.SetHandler(cli.New(os.Stderr))

	//load the oui.txt file from the library location
	_, filename, _, _ := runtime.Caller(0)
	db := ouidb.New(path.Join(path.Dir(filename), "oui.txt"))

	if db == nil {
		log.Fatal("database not initialized")
	}

	v, err := db.VendorLookup(*mac)
	if err != nil {
	    log.WithError(err).Fatalf("error parsing: %s", *mac)
	}

	fmt.Printf("%s => %s\n", *mac, v)
}
