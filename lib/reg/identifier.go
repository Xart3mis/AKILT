package reg

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

func GetUniqueSystemId() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	s, _, err := k.GetStringValue("ProductId")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Windows SystemId = %q\n", s)
	return s
}
