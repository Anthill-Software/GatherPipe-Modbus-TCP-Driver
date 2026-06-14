package driver

import (
	"path/filepath"
	"time"

	GatherPipe "github.com/Anthill-Software/GatherPipe/core"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/simonvetter/modbus"
)

type ModbusTCPDriver struct {
	logger       hclog.Logger
	configPath   string
	config       ModbusConfig
	clientModbus *modbus.ModbusClient
	filters      []GatherPipe.MetricFilter
}

func (driver *ModbusTCPDriver) Init(config GatherPipe.Config) error {
	driver.configPath = filepath.Join(config.Server.Plugin.Dir, "modbus", "config.yaml")
	driver.GetConfiguration()
	driver.filters = append(driver.filters, GatherPipe.MetricFilter{
		Key:   "application",
		Value: "weather",
	})
	driver.logger.Info("Modbus TCP initialisé avec succès")
	return nil
}

func (p *ModbusTCPDriver) Name() (string, error) {
	return "ModbusTCPDriver", nil
}

func (driver *ModbusTCPDriver) Fetch() ([]GatherPipe.Metric, error) {
	var metrics []GatherPipe.Metric
	now := time.Now()

	if client, err := driver.NewModbusClient(); err == nil {
		for _, mCfg := range driver.config.Metrics {
			if value, err := driver.ReadValue(client, mCfg); err == nil {
				metrics = append(metrics, GatherPipe.Metric{
					Timestamp: now,
					ID:        mCfg.Name,
					Value:     value,
					Format:    mCfg.DataType,
					Unit:      mCfg.Unit,
					Filters:   driver.filters,
				})
			} else {
				continue
			}
		}
		driver.CloseModbusClient(client)
	} else {
		return metrics, err
	}

	driver.logger.Debug("Collecte réussie", "count", len(metrics))
	return metrics, nil
}

func Start() {
	pluginLogger := hclog.New(&hclog.LoggerOptions{
		Level:       hclog.Debug,
		DisableTime: true,
	})

	driver := &ModbusTCPDriver{
		logger: pluginLogger,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: GatherPipe.Handshake,
		Plugins: map[string]plugin.Plugin{
			"driver":    &GatherPipe.DriverPlugin{Impl: driver},
			"commander": &GatherPipe.CommanderPlugin{Impl: driver},
		},
		Logger: pluginLogger,
	})
}
