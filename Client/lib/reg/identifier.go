// package reg deals with registry manipulation (r/w)
package reg

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

// gets unique systemid by reading the MachineGuid string from registry
func GetUniqueSystemId() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	s, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Windows MachineGuid = %q\n", s)
	return s
}
