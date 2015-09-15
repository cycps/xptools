package main

import (
	"encoding/json"
	"fmt"
	"github.com/cycps/xptools/dnsc"
	"io/ioutil"
	"os"
)

var cspec dnsc.XPClientSpec
var templateDir string

func resolveConfHead() {
	dnsc.ApplyTemplate("resolve_conf_d_head", "head", templateDir, cspec)
}

func upfile() {
	dnsc.ApplyTemplate("upfile", "upfile", templateDir, cspec)
}

func dnskey() {
	dnsc.ApplyTemplate("dns.key", "dns.key", templateDir, cspec)
}

func setupScript() {
	dnsc.CopyFile("setup_client.sh", templateDir)
	os.Chmod("setup_client.sh", 0755)
}

func doUpdateScript() {
	dnsc.CopyFile("do-nsupdate.sh", templateDir)
	os.Chmod("do-nsupdate.sh", 0755)
}

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: dns_client <client spec>\n")
		os.Exit(1)
	}

	fmt.Printf("dns_client v0.1\n\n")

	templateDir = os.Getenv("GOPATH") +
		"/src/github.com/cycps/xptools/dnsc/client_templates"

	fmt.Printf("template dir: %s\n", templateDir)

	cspec_src, _ := ioutil.ReadFile(os.Args[1])
	err := json.Unmarshal(cspec_src, &cspec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", "bad cspec file")
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	resolveConfHead()
	upfile()
	dnskey()
	setupScript()
	doUpdateScript()

}
