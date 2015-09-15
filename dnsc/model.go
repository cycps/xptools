package dnsc

type RouterSpec struct {
	Xpname  string       `json:"xpname"`
	FQDN    string       `json:"fqdn"`
	Addr    string       `json:"addr"`
	Clients []ClientSpec `json:"clients"`
}

type ClientSpec struct {
	FQDN string `json:"fqdn"`
	Key  string `json:"key"`
}

type XPClientSpec struct {
	Xpname string `json:"xpname"`
	NSaddr string `nsaddr:"nsaddr"`
	ClientSpec
}
