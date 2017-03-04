package stream

import (
	"errors"

	getstream "github.com/GetStream/stream-go"
)

var Client *getstream.Client

func Connect(apiKey string, apiSecret string, apiRegion string) error {
	var err error
	if apiKey == "" || apiSecret == "" || apiRegion == "" {
		return errors.New("Config not complete")
	}

	Client, err = getstream.New(&getstream.Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		Location:  apiRegion,
	})

	return err
}
