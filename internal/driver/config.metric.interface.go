package driver

type MetricConfig struct {
	Name         string  `yaml:"name" json:"name"`
	SlaveID      byte    `yaml:"slave_id" json:"slave_id"`
	RegisterType string  `yaml:"register_type" json:"register_type"` // holding, input, coil, discrete
	Address      uint16  `yaml:"address" json:"address"`
	DataType     string  `yaml:"data_type" json:"data_type"` // int16, uint32, float32...
	Scale        float64 `yaml:"scale" json:"scale"`
	Unit         string  `yaml:"unit" json:"unit"`
}
