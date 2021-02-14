//+build mage

package main

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
)

func SetEnvFromKVConf(fs string) error {
	f, err := os.Open(fs)
	if err != nil {
		return err
	}
	kv := map[string]string{}
	err = toml.NewDecoder(f).Decode(&kv)
	if err != nil {
		return err
	}
	for k, v := range kv {
		fmt.Printf("%s=%s", k, v)
		os.Setenv(k, v)
	}
	return nil
}
