package main

import (
	"net/http"
	"flag"
	"encoding/json"
	"io/ioutil"
	"log"
	ph "github.com/advptr/reverse-proxy/proxyhandler"
	"crypto/x509"
)

// Main config
type Config struct {
	TrustFiles []string `json:"trust-files"`
	Routes []ph.Route `json:"routes""`
}


// Program  Options
type Options struct {
	ServerAddress string
	ConfigFile    string
}


// Main Bootstrap
func main() {
	args := args()
	config := config(args.ConfigFile)

	for _, r := range config.Routes {
		handler := ph.NewWSHandler(r, config.trustedCertPool())
		http.HandleFunc(handler.Path(), handler.Handle)
	}
	log.Printf("reverse-proxy on %v\n", args.ServerAddress)

	err := http.ListenAndServe(args.ServerAddress, nil)
	if err != nil {
		panic(err)
	}
}


// Parse command args
func args() Options {
	const (
		defaultServerAddress = ":80"
		serverAddressUsage   = "server address: ':80', '0.0.0.0:8080'..."
		defaultRouteConfig   = "config.json"
		routeConfigUsage     = "configuration file: 'config.json'"
	)
	address := flag.String("address", defaultServerAddress, serverAddressUsage)
	config := flag.String("config", defaultRouteConfig, routeConfigUsage)
	flag.Parse()

	return Options{*address, *config}
}

//
func (c *Config) trustedCertPool() (*x509.CertPool) {
	trustedCertPool := x509.NewCertPool()
	for _, file := range c.TrustFiles {
		trustedCert, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		trustedCertPool.AppendCertsFromPEM(trustedCert);
	}

	return trustedCertPool
}


// Unmarshal routes from JSON configuration file
func config(configFile string) Config {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

