package util

import (
	"net"
	"p2p/lib/errors"
	//	"os"
)

func GetHostInfo(listener net.Listener) (string, string) {
	
	/*
	name, err := os.Hostname()
	errors.PrintIfError(err, "Error while getting host name from os")

	addrs, _ := net.LookupHost(name)
	errors.PrintIfError(err, "Error while looking up address of host")

	*/
	_, port, err := net.SplitHostPort(listener.Addr().String())
	errors.PrintIfError(err, "Error while getting port of tcp listener")

	ip := "127.0.0.1"
	return ip, port //addrs[0]
}
