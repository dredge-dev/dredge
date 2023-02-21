package providers

import "fmt"

func checkConfig(config map[string]string, configKeys []string) error {
	for _, key := range configKeys {
		if _, ok := config[key]; !ok {
			return fmt.Errorf("could not find field %s in config", key)
		}
	}
	return nil
}
