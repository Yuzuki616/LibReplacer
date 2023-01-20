package main

import (
	"LibReplacer/certer"
	"LibReplacer/conf"
	"LibReplacer/library"
	"LibReplacer/router"
	"flag"
	"log"
)

var config = flag.String("config", "config.json", "config file path")

func main() {
	log.Println("LibReplacer Start")
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
