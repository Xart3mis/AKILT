package bundles

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed "assets/Fira Code Regular Nerd Font.ttf"
var FiraCodeNerd []byte

//go:embed "assets/ca-cert.pem"
var ca_cert []byte

func WriteCaCertPem() {
	f, err := os.Create("ca-cert.pem")
	if err != nil {
		fmt.Println(err)
		return
	}

	n2, err := f.Write(ca_cert)
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

func WriteFiraCodeNerd() {
	f, err := os.Create("FiraCodeNerd.ttf")
	if err != nil {
		fmt.Println(err)
		return
	}

	n2, err := f.Write(FiraCodeNerd)
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
