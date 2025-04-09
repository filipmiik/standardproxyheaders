package standardproxyheaders

import (
	"context"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	ForwardedByHostname bool   `yaml:"forwardedByHostname"`
	ForwardedByHeader   string `yaml:"forwardedByHeader"`
	ForwardedByStatic   string `yaml:"forwardedByStatic"`
	ForwardedForRemote  bool   `yaml:"forwardedForRemote"`
	ForwardedForHeader  string `yaml:"forwardedForHeader"`
	ForwardedForStatic  string `yaml:"forwardedForStatic"`
}

func CreateConfig() *Config {
	return &Config{}
}

type Plugin struct {
	name                string
	next                http.Handler
	forwardedByHostname bool
	forwardedByHeader   string
	forwardedByStatic   string
	forwardedForRemote  bool
	forwardedForHeader  string
	forwardedForStatic  string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Plugin{
		name:                name,
		next:                next,
		forwardedByHostname: config.ForwardedByHostname,
		forwardedByHeader:   strings.TrimSpace(config.ForwardedByHeader),
		forwardedByStatic:   strings.TrimSpace(config.ForwardedByStatic),
		forwardedForRemote:  config.ForwardedForRemote,
		forwardedForHeader:  strings.TrimSpace(config.ForwardedForHeader),
		forwardedForStatic:  strings.TrimSpace(config.ForwardedForStatic),
	}, nil
}

func (plugin *Plugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var ForwardedFields []string

	// Get OS hostname
	Hostname, _ := os.Hostname()
	if len(Hostname) == 0 {
		Hostname = "traefik"
	}

	// Set Forwarded by field
	var ForwardedBy string
	if plugin.forwardedByHostname {
		ForwardedBy = Hostname
	} else if len(plugin.forwardedByHeader) > 0 {
		ForwardedBy = strings.TrimSpace(req.Header.Get(plugin.forwardedByHeader))
	} else if len(plugin.forwardedByStatic) > 0 {
		ForwardedBy = plugin.forwardedByStatic
	}

	// Set Forwarded for field
	var ForwardedFor string
	if plugin.forwardedForRemote {
		ForwardedFor = req.RemoteAddr
	} else if len(plugin.forwardedForHeader) > 0 {
		ForwardedFor = strings.TrimSpace(req.Header.Get(plugin.forwardedForHeader))
	} else if len(plugin.forwardedForStatic) > 0 {
		ForwardedFor = plugin.forwardedForStatic
	}

	// Set other Forwarded fields
	ForwardedHost := req.Host
	ForwardedProto := req.URL.Scheme

	// Construct Forwarded header
	if len(ForwardedBy) > 0 {
		ForwardedFields = append(ForwardedFields, "by="+ForwardedBy)
	}
	if len(ForwardedFor) > 0 {
		ForwardedFields = append(ForwardedFields, "for="+ForwardedFor)
	}
	if len(ForwardedHost) > 0 {
		ForwardedFields = append(ForwardedFields, "host="+ForwardedHost)
	}
	if len(ForwardedProto) > 0 {
		ForwardedFields = append(ForwardedFields, "proto="+ForwardedProto)
	}

	// Append Forwarded header
	Forwarded := req.Header.Get("Forwarded")
	if len(Forwarded) > 0 {
		Forwarded += ", "
	}
	req.Header.Set("Forwarded", Forwarded+strings.Join(ForwardedFields, ";"))

	// Append Via header
	Via := req.Header.Get("Via")
	if len(Via) > 0 {
		Via += ", "
	}
	req.Header.Set("Via", Via+req.Proto+" "+Hostname)

	// Continue processing request
	plugin.next.ServeHTTP(rw, req)
}
