package bundles

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed "assets/server-cert.pem"
var server_cert []byte

//go:embed "assets/server-key.pem"
var server_key []byte

// Write Server Certificate file
func WriteCert() {
	f, err := os.Create("server-cert.pem")
	if err != nil {
		fmt.Println(err)
		return
	}

	n2, err := f.Write(server_cert)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(n2, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Write Server Key file
func WriteCertKey() {
	f, err := os.Create("server-key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}

	n2, err := f.Write(server_key)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	fmt.Println(n2, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
