package router

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	dnsCacheTimeout = 5 * time.Minute
)

var (
	dnsCache  = map[string]lookup{}
	dnsCacheM sync.Mutex
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

	dnsCacheM.Lock()
	defer dnsCacheM.Unlock()

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

	log.Debugf("dialling %s:%s", ips[i], parts[1])

	ip := ips[i]
	if ip := ip.To4(); ip != nil {
		return net.Dial(network, fmt.Sprintf("%s:%s", ip, parts[1]))
	} else if ip := ip.To16(); ip != nil {
		return net.Dial(network, fmt.Sprintf("[%s]:%s", ip, parts[1]))
	} else {
		return nil, fmt.Errorf("invalid ip address '%s'", ip.String())
	}
}
