package dnsc

import (
	"os/exec"
)

func Keygen(dnsname string) error {
	cmd := exec.Command("dnssec-keygen", "-a", "HMAC-MD5", "-n", "USER", dnsname+".")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
