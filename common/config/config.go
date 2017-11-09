package config

import (
	"fmt"
	"github.com/guilhebl/go-props"
	"github.com/guilhebl/xcrypto"
	"log"
	"strconv"
	"strings"
)

var properties props.Properties

func init() {
	log.Printf("%s", "Init Config")

	propsFile, err := props.ReadPropertiesFile("common/config/app-config.properties")
	if err != nil {
		log.Println("Error while reading config properties file")
	}

	properties = propsFile
}

func GetProperty(p string) string {
	return properties[p]
}

func GetIntProperty(p string) int64 {
	prop, _ := strconv.ParseInt(GetProperty(p), 10, 0)
	return prop
}

func getHost() string {
	return properties["protocol"] + properties["host"] + ":" + properties["port"] + "/"
}

func getImageFolderUrl() string {
	return getHost() + "assets/images/"
}

// builds an img source from an external server should proxy if http is used to avoid security warns (HTTP/S MIXED MODE)
// if empty string return img placeholder default
func BuildImgUrlExternal(s string, proxyRequired bool) string {
	if s == "" {
		return fmt.Sprintf(getImageFolderUrl() + "image-placeholder.png")
	}

	if proxyRequired {
		key := GetProperty("privateKeyAES")
		hash, err := xcrypto.Encrypt([]byte(key), []byte(s))
		if err != nil {
			return ""
		}

		url := fmt.Sprintf(GetProperty("proxyHost") + "?hash=" + string(hash))
		return url
	} else {
		// case when provider has an https available service just switch to that if not already https
		url := fmt.Sprintf(s)
		if strings.Index(s, "http://") != -1 {
			url = strings.Replace(s, "http://", "https://", 1)
		}
		return url
	}

	return ""
}

// builds img from local assets folder
func BuildImgUrl(s string) string {
	img := s
	if img == "" {
		img = "image-placeholder.png"
	}

	return fmt.Sprintf(getImageFolderUrl() + img)
}

// returns max number of providers
func CountMarketplaceProviders(country string) int {
	var size int

	switch country {

	//Canada
	case "can":
		{
			arr := strings.Split(GetProperty("marketplaceProvidersCanada"), ",")
			size = len(arr)
		}
	default:
		{
			arr := strings.Split(GetProperty("marketplaceProviders"), ",")
			size = len(arr)
		}
	}

	return size
}
