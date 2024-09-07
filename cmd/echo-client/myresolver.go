package main

import (
	"log"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/resolver"
)

type MyResolverBuilder struct {
	resolver *MyResolver
}

func NewMyResolverBuilder(myResolver *MyResolver) *MyResolverBuilder {
	return &MyResolverBuilder{resolver: myResolver}
}

func (rb *MyResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	log.Printf("Building resolver for target: %v", target)
	log.Println("target.URL.Scheme: ", target.URL.Scheme)
	log.Println("target.URL.Opaque: ", target.URL.Opaque)
	log.Println("target.URL.User: ", target.URL.User)
	log.Println("target.URL.Host: ", target.URL.Host)
	log.Println("target.URL.Path: ", target.URL.Path)
	log.Println("target.URL.RawPath: ", target.URL.RawPath)
	log.Println("target.URL.OmitHost: ", target.URL.OmitHost)
	log.Println("target.URL.ForceQuery: ", target.URL.ForceQuery)
	log.Println("target.URL.RawQuery: ", target.URL.RawQuery)
	log.Println("target.URL.Fragment: ", target.URL.Fragment)
	log.Println("target.URL.RawFragment: ", target.URL.RawFragment)

	rb.resolver.target = target
	rb.resolver.cc = cc
	rb.resolver.start()
	return rb.resolver, nil
}
func (*MyResolverBuilder) Scheme() string { return "dns" }

type MyResolver struct {
	target         resolver.Target
	service        string
	cc             resolver.ClientConn
	lookupInterval time.Duration
	currentServers []string
}

func NewMyResolver(lookupInterval time.Duration) *MyResolver {
	return &MyResolver{lookupInterval: lookupInterval}
}

func (r *MyResolver) start() {
	r.service = r.target.URL.Host
	r.ResolveNow(resolver.ResolveNowOptions{})
	go func() {
		for range time.Tick(r.lookupInterval) {
			r.ResolveNow(resolver.ResolveNowOptions{})
		}
	}()
}

func (r *MyResolver) ResolveOnFailure() {
	log.Printf("ResolveOnFailure")
	r.ResolveNow(resolver.ResolveNowOptions{})
}

func (r *MyResolver) ResolveNow(resolver.ResolveNowOptions) {
	dn := strings.Split(r.service, ":")[0]
	port := strings.Split(r.service, ":")[1]
	log.Printf("Checking hosts for %s", dn)
	newServers, _ := net.LookupHost(dn)
	log.Printf("Resolved servers: %v", newServers)
	sort.Strings(newServers)
	if !cmp.Equal(r.currentServers, newServers) {
		log.Printf("Servers list has changed. New resolved servers: %v", newServers)
		r.currentServers = newServers

		addrs := make([]resolver.Address, len(r.currentServers))
		for i, s := range r.currentServers {
			address := s + ":" + port
			addrs[i] = resolver.Address{Addr: address}
		}
		r.cc.UpdateState(resolver.State{Addresses: addrs})
		log.Printf("Connection state updated")
	}
}
func (*MyResolver) Close() {}
