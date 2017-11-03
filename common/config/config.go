package config

import (
	"log"
	"github.com/guilhebl/fileutil"
)

var props fileutil.AppConfigProperties

func init() {
	log.Printf("%s","Init Config")

	propsFile, err := fileutil.ReadPropertiesFile("app-config.properties")
	if err != nil {
		log.Println("Error while reading properties file")
	}

	props = propsFile
}

func GetProperty(p string) string {
	return props[p]
}

// gets backend endpoint
func GetEndpoint() string {
	return props["protocol"] + props["host"] + ":" + props["port"] + "/" + props["productEndpointPath"]
}
