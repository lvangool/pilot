package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
)

type SysInfo struct {
	Hostname string
	IP       []string
	CallerIP string
	UpSince  time.Time
	Version  string
}

const VERSION = "0.2"

func main() {
	sysInfo := &SysInfo{
		UpSince: time.Now().UTC(),
		Version: VERSION,
	}
	sysInfo.Hostname, _ = os.Hostname()
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				sysInfo.IP = append(sysInfo.IP, v.IP.String())
			case *net.IPAddr:
				sysInfo.IP = append(sysInfo.IP, v.IP.String())
			}
		}
	}
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request) {
		sysInfo.CallerIP = r.RemoteAddr
		w.WriteJson(sysInfo)
	}))
	log.Fatal(http.ListenAndServe(":8050", api.MakeHandler()))
}
