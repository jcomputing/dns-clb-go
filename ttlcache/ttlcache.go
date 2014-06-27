package ttlcache

import (
	"github.com/jcomputing/dns-clb-go/dns"
	"net"
	"time"
)

func NewTtlCache(lib dns.Lookup, ttl int32) *TtlCache {
	c := new(TtlCache)
	c.lib = lib
	c.ttl = ttl
	c.lastUpdate = 0

	return c
}

type TtlCache struct {
	lib        dns.Lookup
	ttl        int32
	lastUpdate int32
	srvs       map[string][]net.SRV
	as         map[string]string
}

func (l *TtlCache) LookupSRV(name string) ([]net.SRV, error) {
	err := l.checkCache()
	if err != nil {
		return nil, err
	}
	_, ok := l.srvs[name]
	if !ok {
		l.srvs[name], err = l.lib.LookupSRV(name)
		if err != nil {
			return nil, err
		}
	}
	return l.srvs[name], nil
}

func (l *TtlCache) LookupA(name string) (string, error) {
	err := l.checkCache()
	if err != nil {
		return "", err
	}

	_, ok := l.as[name]
	if !ok {
		l.as[name], err = l.lib.LookupA(name)
		if err != nil {
			return "", err
		}
	}

	return l.as[name], nil
}

func (l *TtlCache) checkCache() error {
	now := int32(time.Now().Unix())
	if l.lastUpdate+l.ttl < now {
		l.lastUpdate = now
		l.srvs = map[string][]net.SRV{}
		l.as = map[string]string{}
	}
	return nil
}
