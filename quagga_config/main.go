package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
)

func main() {

	fmt.Printf("quagga_config v0.1\n\n")
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: quagga_config <neighbors json>\n")
		os.Exit(1)
	}

	chart := readRouterChart(os.Args[1])
	cfg := initRouterConfig(&chart)

	resolveInterfaceInfo(&cfg, &chart)

	cfg.GenOspf6Conf()
	ioutil.WriteFile("ospf6d.conf", []byte(cfg.Ospf6Conf), 0644)

	cfg.GenZebraConf()
	ioutil.WriteFile("zebra.conf", []byte(cfg.ZebraConf), 0644)

}

//In toplogy (mathematical) a local topology is called a chart
type RouterChart struct {
	Id                         uint32
	DownstreamHosts, PeerHosts []string
}

type RouterConfig struct {
	BasePrefix, BaseAddr, CoreSubnet, DownstreamSubnet string
	DownstreamInterfaces, PeerInterfaces               []string
	Ospf6Conf, ZebraConf                               string
}

func ifxOspf(area, ifx, subnet string) string {
	src := ""
	src += "interface " + ifx + "\n"
	src += "area " + area + " range " + subnet + "\n"
	src += "interface " + ifx + " area " + area + "\n"
	return src
}

func ifxZebra(ifx, addr, prefix string) string {
	src := ""
	src += "interface " + ifx + "\n"
	src += " link-detect\n"
	src += " no ipv6 nd suppress-ra\n"
	src += " ipv6 nd ra-interval 10\n"
	src += " ipv6 address " + addr + "\n"
	src += " ipv6 nd prefix " + prefix + "\n"
	return src
}

func (cfg *RouterConfig) GenOspf6Conf() {

	src := "password muffins\n"
	src += "router ospf6\n"
	src += "redistribute static\n"
	src += "redistribute connected\n"

	src += "!\n!peer\n!\n"
	for i, ifx := range cfg.PeerInterfaces {
		area := fmt.Sprintf("0.0.0.%d", i)
		src += ifxOspf(area, ifx, cfg.CoreSubnet)
		src += "!\n"
	}

	src += "!\n!downstream\n!\n"
	for i, ifx := range cfg.DownstreamInterfaces {
		area := fmt.Sprintf("0.0.0.%d", i+len(cfg.PeerInterfaces)+1)
		src += ifxOspf(area, ifx, cfg.DownstreamSubnet)
		src += "!\n"
	}

	cfg.Ospf6Conf = src

}

func (cfg *RouterConfig) GenZebraConf() {

	src := "hostname r\n"
	src += "password muffins\n"

	src += "!\n!downstream\n!\n"
	for i, ifx := range cfg.DownstreamInterfaces {
		addr := fmt.Sprintf("%s::%d/64", cfg.BasePrefix, i+1)
		prefix := fmt.Sprintf("%s::/64", cfg.BasePrefix, i+1)
		src += ifxZebra(ifx, addr, prefix)
		src += "!\n"
	}

	src += "interface lo\n"
	src += " link-detect\n"

	src += "ipv6 forwarding\n"
	//peer interface routing is dynamic is determined by ospf
	/*
		for i, ifx := range cfg.PeerInterfaces {
			src += " ipv6 route " + fmt.Sprintf("%s::%d/64 %s\n", cfg.BasePrefix, i+1, ifx)
		}
	*/
	for i, ifx := range cfg.DownstreamInterfaces {
		src += " ipv6 route " + fmt.Sprintf("%s::%d/64 %s\n", cfg.BasePrefix, i+1, ifx)
	}

	cfg.ZebraConf = src

}

func readRouterChart(filename string) RouterChart {

	src, _ := ioutil.ReadFile(filename)

	info := RouterChart{}
	json.Unmarshal(src, &info)

	return info

}

func resolveHostIfx(host string) string {

	ip, _ := net.LookupIP(host)
	theIP := ip[0].String()

	out, _ := exec.Command("ip", "route", "get", theIP).Output()

	rx, _ := regexp.Compile("dev\\s+(\\S+)")
	m := rx.FindStringSubmatch(string(out))

	return m[len(m)-1]

}

func uint32ToHexPrefix(x uint32) string {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, x)
	s := hex.EncodeToString(buf)
	_s := s[:4] + ":" + s[4:]
	return _s
}

func initRouterConfig(ec *RouterChart) RouterConfig {

	rc := RouterConfig{}
	rc.CoreSubnet = "2001:cc::/32"
	rc.BasePrefix = "2001:cc:" + uint32ToHexPrefix(ec.Id)
	rc.BaseAddr = rc.BasePrefix + "::1"
	rc.DownstreamSubnet = rc.BasePrefix + "::/64"

	fmt.Printf("base-address: %s\n", rc.BaseAddr)
	fmt.Printf("core-subnet: %s\n", rc.CoreSubnet)
	fmt.Printf("downstream-subnet: %s\n\n", rc.DownstreamSubnet)

	return rc

}

func resolveInterfaceInfo(cfg *RouterConfig, chart *RouterChart) {

	fmt.Println("downstream")
	for _, x := range chart.DownstreamHosts {

		ifx := resolveHostIfx(x)
		fmt.Printf("%s -> %s\n", x, ifx)
		cfg.DownstreamInterfaces = append(cfg.DownstreamInterfaces, ifx)

	}

	fmt.Println("peer")
	for _, x := range chart.PeerHosts {

		ifx := resolveHostIfx(x)
		fmt.Printf("%s -> %s\n", x, ifx)
		cfg.PeerInterfaces = append(cfg.PeerInterfaces, ifx)

	}

}
