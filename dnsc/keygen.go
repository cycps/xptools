package dnsc

import (
	"os/exec"
)

func Keygen(dnsname string) error {
	cmd := exec.Command(
		"dnssec-keygen", "-r", "/dev/urandom",
		"-a", "HMAC-MD5", "-b", "512",
		"-n", "USER",
		dnsname+".")
	err := cmd.Run()
	if err != nil {
		return err
	}

	exec.Command("bash", "-c", "rm K"+dnsname+"*").Run()
	return nil
}
