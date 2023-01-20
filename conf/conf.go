package conf

import (
	"fmt"
	"github.com/goccy/go-json"
	"os"
)

type Conf struct {
	EmbyUrl     string      `json:"EmbyUrl"`
	LibraryPath string      `json:"LibraryPath"`
	BlackList   []string    `json:"BlackList"`
	Addr        string      `json:"Addr"`
	EnableSsl   bool        `json:"EnableSsl"`
	CertConfig  *CertConfig `json:"CertConfig"`
}

type CertConfig struct {
	Mode      string            `json:"Mode"` // file, http, dns
	Domain    string            `json:"Domain"`
	DataPath  string            `json:"DataPath"`
	CertFile  string            `json:"CertFile"`
	KeyFile   string            `json:"KeyFile"`
	Provider  string            `json:"Provider"` // alidns, cloudflare, gandi, godaddy....
	Email     string            `json:"Email"`
	DnsEnv    map[string]string `json:"DnsEnv"`
	CheckTime int               `json:"CheckTime"`
}

func New() *Conf {
	return &Conf{
		EmbyUrl:     "http://localhost:8096",
		LibraryPath: "./lib.json",
		Addr:        ":8097",
		CertConfig: &CertConfig{
			Mode:     "none",
			Domain:   "localhost",
			DataPath: "./cert",
			CertFile: "./cert/1.pem",
			KeyFile:  "./cert/1.key",
			Provider: "alidns",
			Email:    "",
		},
	}
}

func (p *Conf) LoadFromPath(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config file error: %s", err)
	}
	err = json.NewDecoder(f).Decode(p)
	if err != nil {
		return fmt.Errorf("decode config error: %s", err)
	}
	return nil
}
