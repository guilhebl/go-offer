package config

import (
	"fmt"
	"github.com/guilhebl/go-props"
	"github.com/guilhebl/xcrypto"
	"log"
	"strings"
	"sync"
)

var instance *props.Properties
var once sync.Once

func GetInstance() *props.Properties {
	once.Do(func() {
		instance = newProperties()
	})
	return instance
}

func newProperties() *props.Properties {
	log.Printf("%s", "Init Config")

	propsFile, err := props.ReadPropertiesFile("common/config/app-config.properties")
	if err != nil {
		log.Fatal(err)
	}

	return &propsFile
}

func GetProperty(p string) string {
	return GetInstance().GetProperty(p)
}

func GetIntProperty(p string) int64 {
	return GetInstance().GetIntProperty(p)
}

func getHost() string {
	return GetProperty("protocol") + GetProperty("host") + ":" + GetProperty("port") + "/"
}

func getImageFolderUrl() string {
	return getHost() + "assets/images/"
}

// returns if this provider requires image proxy (HTTP/HTTPS)
func IsProxyRequired(provider string) bool {
	return strings.Index(GetProperty("marketplaceProvidersImageProxyRequired"), provider) != -1
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

func CountMarketplaceProviderListSize() int {
	arr := strings.Split(GetProperty("marketplaceProviders"), ",")
	return len(arr)
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
