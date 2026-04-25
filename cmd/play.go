package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/rad/radio"
)

func play(input string, player string) error {
	url := input

	stations := loadStations()
	if strings.HasPrefix(input, "@") {
		if stations == nil {
			return fmt.Errorf("no stations found in config")
		}

		alias := strings.TrimPrefix(input, "@")
		station, ok := stations[alias]
		if !ok {
			return fmt.Errorf("station alias not found: %s", alias)
		}
		url = station.URL
	} else if stations != nil {
		if station, ok := stations[input]; ok {
			url = station.URL
		}
	}

	r := radio.NewRadio(url)
	return r.Play(player)
}
