package driver

type ModbusConfig struct {
	Endpoint string         `yaml:"endpoint" json:"endpoint"`
	Timeout  string         `yaml:"timeout" json:"timeout"`
	Metrics  []MetricConfig `yaml:"metrics" json:"metrics"`
}
