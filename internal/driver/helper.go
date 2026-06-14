package driver

import "time"

func (driver *ModbusTCPDriver) ToDuration(d string) time.Duration {
	if dt, err := time.ParseDuration(d); err != nil {
		driver.logger.Warn("Timeout invalide dans la config, valeur par défaut 2s", "error", err)
		return 2 * time.Second
	} else {
		return dt
	}
}
