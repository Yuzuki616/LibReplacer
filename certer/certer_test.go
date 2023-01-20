package certer

import (
	"LibReplacer/conf"
	"log"
	"os"
	"testing"
)

var l *Certer

func init() {
	var err error
	l, err = New(&conf.CertConfig{
		Mode:     "dns",
		Email:    "test@test.com",
		Domain:   "test.test.com",
		Provider: "cloudflare",
		DnsEnv: map[string]string{
			"CF_DNS_API_TOKEN": "123",
		},
		CertFile: "./cert/1.pem",
		KeyFile:  "./cert/1.key",
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func TestLego_CreateCertByDns(t *testing.T) {
	err := l.CreateCert()
	if err != nil {
		t.Error(err)
	}
}

func TestLego_RenewCert(t *testing.T) {
	log.Println(l.RenewCert())
}
