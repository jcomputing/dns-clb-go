package roundrobinclb

import (
	"fmt"
	"github.com/jcomputing/dns-clb-go/dns"
	"net"
	"sort"
)

type ByTarget []net.SRV

func (a ByTarget) Len() int           { return len(a) }
func (a ByTarget) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTarget) Less(i, j int) bool { return a[i].Target < a[j].Target }

func NewRoundRobinClb(lib dns.Lookup) *RoundRobinClb {
	lb := new(RoundRobinClb)
	lb.dnsLib = lib
	lb.i = make(map[string]int)

	return lb
}

type RoundRobinClb struct {
	dnsLib dns.Lookup
	i      map[string]int
}

func (lb *RoundRobinClb) GetAddress(name string) (dns.Address, error) {
	add := dns.Address{}

	srvs, err := lb.dnsLib.LookupSRV(name)
	if err != nil {
		return add, err
	}

	if len(srvs) == 0 {
		return add, fmt.Errorf("no SRV records found")
	}

	sort.Sort(ByTarget(srvs))

	//	log.Printf("%+v", srvs)
	i := lb.i[name]
	i = i % len(srvs)
	srv := srvs[i]
	lb.i[name] = i + 1

	ip, err := lb.dnsLib.LookupA(srv.Target)
	if err != nil {
		return add, err
	}

	return dns.Address{Address: ip, Port: srv.Port}, nil
}
