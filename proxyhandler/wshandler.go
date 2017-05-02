package proxyhandler

import (
	"crypto/x509"
	"net/url"
	"log"
	"net/http/httputil"
	"html/template"
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"bufio"
	xj "github.com/basgys/goxml2json"
)

// Transform config.
type Transform struct {
	RequestTemplate string `json:"request-template"`
	ResponseSection string `json:"response-section"`
	Template *template.Template
}

// Reverse Proxy Handler (route controller)
type ProxyHandler struct {
	Proxy *httputil.ReverseProxy
	Route Route
}

// Creates a new Proxy Handler
func NewWSHandler(route Route, certPool *x509.CertPool) Handler {
	url, err := url.Parse(route.Service)

	if err != nil {
		log.Fatal("%v\n", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	if route.CertFile != "" {
		proxy.Transport = transport(route, certPool)
	}

	proxy.Director = director(route.Service)
	proxyHandler := &ProxyHandler{Route: route, Proxy: proxy}

	if route.Transform != nil {
		log.Printf("Transformer: %v -> %v\n", route.Name, route.Transform.RequestTemplate)
		route.Transform.Template = template.Must(template.New(route.Name).ParseFiles(route.Transform.RequestTemplate))
		route.Transform.Template.Option("missingkey=zero")
		proxy.ModifyResponse = proxyHandler.transformResponse
	}

	return proxyHandler
}

//
func (p *ProxyHandler) Path() (string) {
	return p.Route.Path
}

// Main HTTP Request Handler
func (p *ProxyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Route %v -- %v\n", p.Route.Path, r.Header)
	w.Header().Set("X-GoProxy", "GoProxy")
	r.URL.Path = ""
	if p.Route.Transform != nil {
		log.Printf("Transform %v\n", p.Route.Transform.Template)
		err := p.transformRequest(r)
		if err != nil {
			log.Println(err)
		}
	}
	p.Proxy.ServeHTTP(w, r)
}

// Transforms requests from JSON to SOAP
func (p *ProxyHandler) transformRequest(r *http.Request) error {
	sourceBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	log.Printf("Before Transform: %v\n", string(sourceBytes))


	var f interface{}
	err = json.Unmarshal(sourceBytes, &f)
	if err != nil {
		return err
	}

	var destBuffer bytes.Buffer
	wr := bufio.NewWriter(&destBuffer)

	err = p.Route.Transform.Template.ExecuteTemplate(wr, p.Route.Transform.RequestTemplate, f)
	if err != nil {
		return err
	}

	err = wr.Flush()
	if err != nil {
		return err
	}
	destBytes := destBuffer.Bytes()

	log.Printf("After Transform: %v\n", string(destBytes))

	r.ContentLength = int64(len(destBytes))
	r.Body = ioutil.NopCloser(bytes.NewReader(destBytes))
	r.Header.Set("Content-Type", "application/xml")
	r.Header.Set("SOAPAction", "")

	return nil
}

// Transforms responses from SOAP to JSON
func (p *ProxyHandler) transformResponse(response *http.Response) error {
	jsonBody, err := xj.Convert(response.Body)
	if err != nil {
		return err
	}

	var f interface{}
	err = json.Unmarshal(jsonBody.Bytes(), &f)
	if err != nil {
		return err
	}

	section := findSection(p.Route.Transform.ResponseSection, f)

	var body []byte
	body, err = json.Marshal(section)
	if err != nil {
		return err
	}

	response.ContentLength = int64(len(body))
	response.Body = ioutil.NopCloser(bytes.NewReader(body))
	response.Header.Set("Content-Type", "application/json; charset=utf-8")

	return nil
}

//
func findSection(name string, f interface{}) interface{} {
	m := f.(map[string]interface{})

	for k, v := range m {
		if k == name {
			return m[name]
		}
		if _, ok := m[k]; ok {
			switch v.(type) {
			case string, int, int64, float64, float32, bool, []interface{}:
			default:
				value := findSection(name, m[k])
				if value != nil {
					return value
				}
			}
		}
	}
	return nil
}

// Transport security to backend services
func transport(route Route, pool *x509.CertPool) *http.Transport {
	cert := route.loadClientCert()

	c := tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: pool}
	c.BuildNameToCertificate()

	return &http.Transport{TLSClientConfig: &c}
}

// Proxy Request Director
func director(serviceUrl string) func(req *http.Request) {
	target, err := url.Parse(serviceUrl)
	if err != nil {
		log.Fatal(err)
	}

	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		targetQuery := target.RawQuery
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		req.Header.Set("User-Agent", "RP/1.0")
	}

	return director
}


