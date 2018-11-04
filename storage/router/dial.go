package router

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	dnsCacheTimeout = 5 * time.Minute
)

var (
	dnsCache = map[string]lookup{}
)

type lookup struct {
	ips []net.IP
	on  time.Time
}

//dial wrapper around net.Dial that provide dns lookup caching
func dial(network, address string) (net.Conn, error) {
	parts := strings.SplitN(address, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("mallformed address expected format <host>:<port>")
	}

	var ips []net.IP
	if lookup, ok := dnsCache[parts[0]]; ok {
		if time.Since(lookup.on) < dnsCacheTimeout {
			ips = lookup.ips
		}
	}

	if len(ips) == 0 {
		var err error
		ips, err = net.LookupIP(parts[0])
		if err != nil {
			return nil, err
		}

		dnsCache[parts[0]] = lookup{ips, time.Now()}
	}

	i := rand.Intn(len(ips))

	return net.Dial(network, fmt.Sprintf("%s:%s", ips[i], parts[1]))
}
