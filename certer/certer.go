package certer

import (
	"LibReplacer/conf"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"os"
	"path"
	"time"
)

type Certer struct {
	client *lego.Client
	config *conf.CertConfig
}

func New(config *conf.CertConfig) (*Certer, error) {
	user, err := NewUser(path.Join(config.DataPath,
		fmt.Sprintf("user-%s.json", config.Email)),
		config.Email)
	if err != nil {
		return nil, fmt.Errorf("create user error: %s", err)
	}
	c := lego.NewConfig(user)
	//c.CADirURL = "http://192.168.99.100:4000/directory"
	c.Certificate.KeyType = certcrypto.RSA2048
	client, err := lego.NewClient(c)
	if err != nil {
		return nil, err
	}
	l := Certer{
		client: client,
		config: config,
	}
	err = l.SetProvider()
	if err != nil {
		return nil, fmt.Errorf("set provider error: %s", err)
	}
	return &l, nil
}

func checkPath(p string) error {
	if !isExist(path.Dir(p)) {
		err := os.MkdirAll(path.Dir(p), 0755)
		if err != nil {
			return fmt.Errorf("create dir error: %s", err)
		}
	}
	return nil
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func (l *Certer) Start() error {
	err := l.SetProvider()
	if err != nil {
		return fmt.Errorf("set provider error: %s", err)
	}
	if !isExist(l.config.CertFile) ||
		!isExist(l.config.KeyFile) {
		err := l.CreateCert()
		if err != nil {
			fmt.Printf("get cert error: %s", err)
		}
	}
	go func() {
		for range time.Tick(time.Second * time.Duration(l.config.CheckTime)) {
			err := l.RenewCert()
			if err != nil {
				fmt.Printf("renew cert error: %s", err)
			}
		}
	}()
	return nil
}
