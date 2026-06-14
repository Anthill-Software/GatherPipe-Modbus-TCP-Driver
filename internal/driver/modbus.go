package driver

import (
	"fmt"

	"github.com/simonvetter/modbus"
)

func (driver *ModbusTCPDriver) NewModbusClient() (*modbus.ModbusClient, error) {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://" + driver.config.Endpoint,
		Timeout: driver.ToDuration(driver.config.Timeout),
	})
	if err != nil {
		return nil, fmt.Errorf("Impossible d'initialiser un client Modbus: %w", err)
	}

	if err := client.Open(); err != nil {
		return nil, fmt.Errorf("Impossible d'ouvrir la connexion Modbus sur %s: %w", driver.config.Endpoint, err)
	}

	if err := client.SetEncoding(modbus.BIG_ENDIAN, modbus.HIGH_WORD_FIRST); err != nil {
		driver.logger.Warn("Impossible de configurer l'encodage Modbus", "error", err)
	}

	return client, nil
}

func (driver *ModbusTCPDriver) CloseModbusClient(client *modbus.ModbusClient) error {
	if err := client.Close(); err != nil {
		return err
	}
	return nil
}

func (driver *ModbusTCPDriver) ReadValue(client *modbus.ModbusClient, mCfg MetricConfig) (string, error) {
	var rawValue float64
	var errRead error

	// Détermination du type de registre (Holding ou Input)
	var regType modbus.RegType
	if mCfg.RegisterType == "input" {
		regType = modbus.INPUT_REGISTER
	} else {
		regType = modbus.HOLDING_REGISTER
	}

	client.SetUnitId(mCfg.SlaveID)

	switch mCfg.DataType {
	case "float32":
		val, err := client.ReadFloat32(mCfg.Address, regType)
		rawValue = float64(val)
		errRead = err

	case "uint16":
		val, err := client.ReadRegister(mCfg.Address, regType)
		rawValue = float64(val)
		errRead = err

	case "int16":
		val, err := client.ReadRegister(mCfg.Address, regType)
		rawValue = float64(int16(val))
		errRead = err

	case "int32":
		val, err := client.ReadRegisters(mCfg.Address, 2, regType)
		combined := uint32(val[0])<<16 | uint32(val[1])
		rawValue = float64(int32(combined))
		errRead = err

	default:
		driver.logger.Warn("Format de données non supporté", "format", mCfg.DataType, "metric", mCfg.Name)
		return "", fmt.Errorf("")
	}

	if errRead != nil {
		driver.logger.Error("Erreur de lecture du registre Modbus", "metric", mCfg.Name, "address", mCfg.Address, "error", errRead)
		return "", fmt.Errorf("")
	}

	if mCfg.Scale != 0 {
		rawValue = rawValue * mCfg.Scale
	}

	return fmt.Sprintf("%f", rawValue), nil
}
