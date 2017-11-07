package config

import (
	"log"
	"github.com/guilhebl/fileutil"
	"strconv"
)

var props fileutil.AppConfigProperties

func init() {
	log.Printf("%s","Init Config")

	propsFile, err := fileutil.ReadPropertiesFile("app-config.properties")
	if err != nil {
		log.Println("Error while reading config properties file")
	}

	props = propsFile
}

func GetProperty(p string) string {
	return props[p]
}

func GetIntProperty(p string) int64 {
	prop, _ := strconv.ParseInt(GetProperty(p), 10, 0)
	return prop
}

// gets backend endpoint
func GetEndpoint() string {
	return props["protocol"] + props["host"] + ":" + props["port"] + "/" + props["productEndpointPath"]
}
