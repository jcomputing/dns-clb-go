package randomclb

import (
	"errors"
	"github.com/jcomputing/dns-clb-go/dns"
	"math/rand"
)

func NewRandomClb(lib dns.Lookup) *RandomClb {
	lb := new(RandomClb)
	lb.dnsLib = lib
	return lb
}

type RandomClb struct {
	dnsLib dns.Lookup
}

func (lb *RandomClb) GetAddress(name string) (dns.Address, error) {
	add := dns.Address{}

	srvs, err := lb.dnsLib.LookupSRV(name)
	if err != nil {
		return add, err
	}
	//	log.Printf("%+v", srvs)
	if len(srvs) == 0 {
		return add, errors.New("No SRV records found.")
	}

	srv := srvs[rand.Intn(len(srvs))]

	ip, err := lb.dnsLib.LookupA(srv.Target)
	if err != nil {
		return add, err
	}

	return dns.Address{Address: ip, Port: srv.Port}, nil
}
