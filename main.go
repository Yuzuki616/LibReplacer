package main

import (
	"flag"
	"github.com/Yuzuki616/LibReplacer/certer"
	"github.com/Yuzuki616/LibReplacer/conf"
	"github.com/Yuzuki616/LibReplacer/library"
	"github.com/Yuzuki616/LibReplacer/router"
	"log"
)

var config = flag.String("config", "config.json", "config file path")
var version = "temp"

func main() {
	log.Println("LibReplacer Start")
	log.Println("Version: ", version)
	flag.Parse()
	c := conf.New()
	err := c.LoadFromPath(*config)
	if err != nil {
		log.Fatalln("load config error: ", err)
	}
	lib := library.New(c.LibraryPath)
	err = lib.Load()
	if err != nil {
		log.Fatalln("Load library file error: ", err)
	}
	err = lib.StartWatcher()
	if err != nil {
		log.Fatalln("Start library watcher error: ", err)
	}
	if c.EnableSsl {
		var cert *certer.Certer
		cert, err = certer.New(c.CertConfig)
		if err != nil {
			log.Fatalln("Create certer obj error", err)
		}
		err = cert.Start()
		if err != nil {
			log.Fatalln("Start certer error", err)
		}
	}
	r := router.New(c, lib)
	err = r.Start()
	if err != nil {
		log.Fatalln("Start router error: ", err)
	}
}
