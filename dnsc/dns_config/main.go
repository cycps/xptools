package main

import (
	"encoding/json"
	"fmt"
	"github.com/cycps/xptools/dnsc"
	"io/ioutil"
	"os"
	"text/template"
)

var templateDir string
var rspec dnsc.RouterSpec

func copyFile(lfn string) {
	fn := templateDir + "/" + lfn
	f, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"failed to read %s template (from GOPATH)\n", fn)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	ioutil.WriteFile(lfn, f, 0644)
}

func namedConf() {
	copyFile("named.conf")
}

func applyTemplate(lfn, ofn string) {

	f, _ := ioutil.ReadFile(templateDir + "/" + lfn)
	tmpl, _ := template.New(lfn).Parse(string(f))
	of, _ := os.Create(ofn)
	tmpl.Execute(of, rspec)
	of.Close()

}

func namedConfLocal() {
	applyTemplate("named.conf.local", "named.conf.local")
}

func namedConfOptions() {
	copyFile("named.conf.options")
}

func keys() {
	applyTemplate("keys.conf", "keys.conf")
}

func apparmor() {
	copyFile("usr.sbin.named")
}

func zoneConf() {
	applyTemplate(
		"db.__XPNAME__.cypress.net",
		fmt.Sprintf("db.%s.cypress.net", rspec.Xpname))
}

func setupScript() {
	fn := "setup_dns.sh"
	applyTemplate(fn, fn)
	os.Chmod(fn, 0755)
}

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: dns_config <router spec>\n")
		os.Exit(1)
	}

	fmt.Printf("dns_config v0.1\n\n")

	templateDir = os.Getenv("GOPATH") +
		"/src/github.com/cycps/xptools/dnsc/server_templates"

	fmt.Printf("template dir: %s\n", templateDir)

	rspec_src, _ := ioutil.ReadFile(os.Args[1])
	err := json.Unmarshal(rspec_src, &rspec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", "bad rspec file")
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	namedConf()
	namedConfLocal()
	namedConfOptions()
	keys()
	zoneConf()
	apparmor()
	setupScript()

}
