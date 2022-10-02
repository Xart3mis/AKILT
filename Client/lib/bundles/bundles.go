package bundles

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed assets/MedusaGothic.ttf
var MedusaGothic []byte

//go:embed "assets/Fira Code Regular Nerd Font.ttf"
var FiraCodeNerd []byte

func WriteMedusaGothic() {
	f, err := os.Create("MedusaGothic.ttf")
	if err != nil {
		fmt.Println(err)
		return
	}

	n2, err := f.Write(MedusaGothic)
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
