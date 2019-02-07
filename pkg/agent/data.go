package agent

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/nicktitle/rc/pkg/fileformat"
	"github.com/pkg/errors"
)

func sendData(connSettings fileformat.ConnectionSettings) (Event, error) {
	var e Event

	client, err := makeHTTPClient(connSettings.Source)
	if err != nil {
		return e, errors.Wrap(err, "building client for data sending")
	}

	postBytes := make([]byte, connSettings.PayloadSize)
	if _, err := rand.Read(postBytes); err != nil {
		return e, errors.Wrap(err, "failed to read random bytes for payload")
	}

	postURL := url.URL{
		Scheme: connSettings.Scheme,
		Host:   fmt.Sprintf("%s:%v", connSettings.Destination.Host, connSettings.Destination.Port),
	}

	resp, err := client.Post(postURL.String(), "", bytes.NewBuffer(postBytes))
	if err != nil {
		return e, errors.Wrapf(err, "failed to post to url: %s", postURL.String())
	}
	defer resp.Body.Close()

	return Event{
		"activity":            fileformat.SendData,
		"destination_address": connSettings.Destination.Host,
		"destination_port":    connSettings.Destination.Port,
		"source_address":      connSettings.Source.Host,
		"source_port":         connSettings.Source.Port,
		"bytes_sent":          connSettings.PayloadSize,
		"protocol":            connSettings.Scheme,
	}, nil
}

// this was a tricky one to figure out. in order to pick your TCP address,
// you need to replace the default dialer that you'd get with a standard http
// client with one that is specifically targeting the host/port combo you want
//
// my solution was accompished by pulling the strings i found at the following
// url, and then jumping into source to figure out my specifics
//
// blog post: https://joshrendek.com/2015/09/using-a-custom-http-dialer-in-go/
// relevant src: https://golang.org/src/net/tcpsock.go?s=1596:1658#L68
//
func makeHTTPClient(addr fileformat.ConnectionAddress) (http.Client, error) {
	var c http.Client

	// first ensure that you can resolve this address
	localAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%v", addr.Host, addr.Port))
	if err != nil {
		return c, errors.Wrap(err, "unable to resolve TCP address with runbook settings")
	}

	// then provide this address to your dialer. note that the other settings there
	// are not required, but will cause hangs if using initialized values (0's are ignored)
	// literature on custom dialers showed me that we need to handle our own timeouts at a minimum
	customDialer := &net.Dialer{
		LocalAddr: localAddr,
		KeepAlive: 1 * time.Second,
		Timeout:   1 * time.Second,
	}

	// finally, override the normal transport with a new one, which will use the addr you
	// specify, instead of pushing traffic on a random available port. note again that i needed
	// to supply a tls handshake value here to prevent the server hanging on an unresponsive tls request
	return http.Client{
		Transport: &http.Transport{
			Dial:                customDialer.Dial,
			TLSHandshakeTimeout: 1 * time.Second,
		}}, nil
}
