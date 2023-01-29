package cmd

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sooprox/types"
	"sooprox/utils"
	"strconv"
	"strings"
)

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello there! You seem to be lost!"))
	})
}

func Listen(defaultConfig types.Config, config <-chan types.Config) {

	host := defaultConfig.Host + ":" + strconv.Itoa(defaultConfig.Port)
	log.Printf("Listening on %s", host)

	// Set the custom 404 handler
	server := &http.Server{
		Addr: host,
	}

	if defaultConfig.Secure {
		log.Printf("Using TLS hence checking for CA configurations")
		if _, err := os.Stat(utils.GetCAPath("ca.pem")); os.IsNotExist(err) {
			if err := utils.GenerateCA(); err != nil {
				panic(err)
			}
		}

		tls, err := utils.MakeTlsConfig(utils.GetCAPath("ca.pem"), utils.GetCAPath("ca.key"), defaultConfig.Host)
		if err != nil {
			log.Fatalf("Either ca.pem or ca.key is missing or the files are corrupted: %s", err)
			log.Printf("You can generate a new CA by deleting the existing ca.pem and ca.key files. The server will automatically generate a new CA")
		}
		server.TLSConfig = tls
	}
	makeRoutes(server, defaultConfig, http.NewServeMux())

	go func(config <-chan types.Config) {
		for {
			c := <-config
			// Reload config
			log.Println("Reloading config")
			makeRoutes(server, c, http.NewServeMux())
		}
	}(config)

	if defaultConfig.Secure {
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		log.Fatal(server.ListenAndServe())
	}

}

// func xListen(host string, port int, config string) {
// 	c := types.Config{}

// 	if config != "" {
// 		c = ReadConfig(config)
// 	}

// 	configChanged := make(chan interface{})
// 	defer close(configChanged)

// 	go utils.WatchConfig(configChanged, config)

// 	if _, err := os.Stat("ca.pem"); os.IsNotExist(err) {
// 		if err := utils.GenerateCA(); err != nil {
// 			panic(err)
// 		}
// 	}
// 	tls, err := utils.MakeTlsConfig("ca.pem", "ca.key", c.Host)
// 	if err != nil {
// 		panic(err)
// 	}

// 	host = host + ":" + strconv.Itoa(port)

// 	server := &http.Server{
// 		Addr:      host,
// 		Handler:   NotFoundHandler(),
// 		TLSConfig: tls,
// 	}

// 	makeRoutes(server, c, http.NewServeMux())

// 	go func(config <-chan interface{}) {
// 		for {

// 			<-config
// 			// Reload config
// 			log.Println("Reloading config")
// 			c = ReadConfig()
// 			makeRoutes(server, c, http.NewServeMux())
// 		}
// 	}(configChanged)

// 	log.Fatal(server.ListenAndServeTLS("", ""))
// }

func makeRoutes(server *http.Server, c types.Config, mux *http.ServeMux) {
	log.Println("Bootstrapping routes")
	server.Handler = mux

	prefixes := []string{}

g:
	for _, p := range c.Proxies {
		if !strings.HasPrefix(p.Prefix, "/") {
			p.Prefix = "/" + p.Prefix
		}
		if !strings.HasSuffix(p.Prefix, "/") {
			p.Prefix = p.Prefix + "/"
		}
		// check if p.Prefix in prefixes
		for _, v := range prefixes {
			if v == p.Prefix {
				log.Printf("Proxy prefix %s pointed to %s is already in use", p.Prefix, p.Host)
				continue g
			}
		}

		prefixes = append(prefixes, p.Prefix)
		log.Printf("Registering %s to %s", p.Prefix, p.Host)
		register(p.Prefix, p.Host, mux)
	}

	// check if "/" is in prefixes
	for _, v := range prefixes {
		if v == "/" {
			log.Println("Found / in prefixes")
			return
		}
	}

	mux.HandleFunc("/", NotFoundHandler().ServeHTTP)
}

func register(prefix, target string, mux *http.ServeMux) {
	url, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Director = func(req *http.Request) {
		path := strings.TrimPrefix(req.URL.Path, prefix)
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		req.URL.Path = path
		req.URL.Host = url.Host
		req.Host = url.Host
		req.URL.Scheme = url.Scheme

	}
	mux.HandleFunc(prefix, proxy.ServeHTTP)
}
