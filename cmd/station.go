package cmd

import "github.com/spf13/viper"

type Station struct {
	Alias string
	URL   string
}

type Stations map[string]Station

func loadStations() Stations {
	rawStations := viper.GetStringMapString("stations")
	if rawStations == nil {
		return nil
	}

	stations := make(Stations, len(rawStations))
	for alias, url := range rawStations {
		stations[alias] = Station{
			Alias: alias,
			URL:   url,
		}
	}

	return stations
}

func saveStation(station Station) error {
	rawStations := viper.GetStringMapString("stations")
	if rawStations == nil {
		rawStations = make(map[string]string)
	}

	rawStations[station.Alias] = station.URL
	viper.Set("stations", rawStations)

	return viper.WriteConfig()
}
