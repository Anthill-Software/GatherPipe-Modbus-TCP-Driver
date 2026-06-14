package driver

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func (d *ModbusTCPDriver) GetConfiguration() {
	if _, err := os.Stat(d.configPath); os.IsNotExist(err) {
		d.logger.Debug("Création du fichier de configuration")
		if err := d.createConfiguration(); err != nil {
			d.logger.Error(err.Error())
		}
	}

	if err := d.readConfiguration(); err != nil {
		d.logger.Error("Échec du parsing de la configuration Modbus", "error", err)
	}
}

func (d *ModbusTCPDriver) persistConfiguration() error {
	// Encoder et écrire le fichier
	data, err := yaml.Marshal(&d.config)
	if err != nil {
		return fmt.Errorf("Erreur lors de la sérialisation de la config par défaut: %w", err)
	}

	if err := os.WriteFile(d.configPath, data, 0644); err != nil {
		return fmt.Errorf("Impossible d'écrire la config par défaut: %w", err)
	}

	return nil
}

func (d *ModbusTCPDriver) createConfiguration() error {
	// S'assurer que les répertoires parents existent
	if err := os.MkdirAll(filepath.Dir(d.configPath), 0755); err != nil {
		return fmt.Errorf("Impossible de créer le dossier de configuration: %w", err)
	}

	// Initialiser une configuration par défaut
	d.config = ModbusConfig{
		Endpoint: "127.0.0.1:502",
		Timeout:  "2s",
		Metrics:  []MetricConfig{},
	}

	if err := d.persistConfiguration(); err != nil {
		return err
	}

	return nil
}

func (d *ModbusTCPDriver) readConfiguration() error {
	data, err := os.ReadFile(d.configPath)
	if err != nil {
		return fmt.Errorf("Impossible de lire le fichier de configuration: %w", err)
	}

	if err := yaml.Unmarshal(data, &d.config); err != nil {
		return fmt.Errorf("Erreur de parsing du fichier YAML: %w", err)
	}

	return nil
}

func (d *ModbusTCPDriver) GetEndPoint() string {
	return d.config.Endpoint
}

func (d *ModbusTCPDriver) SetEndPoint(endpoint string) error {
	d.config.Endpoint = endpoint
	if err := d.persistConfiguration(); err != nil {
		return err
	}
	return nil
}

func (d *ModbusTCPDriver) ListMetrics() []MetricConfig {
	return d.config.Metrics
}

func (d *ModbusTCPDriver) AddMetric(metric MetricConfig) error {
	for _, m := range d.config.Metrics {
		if m.Name == metric.Name {
			return fmt.Errorf("La métrique %s existe déjà", metric.Name)
		}
	}

	d.config.Metrics = append(d.config.Metrics, metric)

	if err := d.persistConfiguration(); err != nil {
		return fmt.Errorf("Impossible de sauvegarder après l'ajout de la métrique: %w", err)
	}

	d.logger.Debug("Métrique ajoutée avec succès", "name", metric.Name)
	return nil
}

func (d *ModbusTCPDriver) DelMetric(name string) error {
	found := false

	for i, m := range d.config.Metrics {
		if m.Name == name {
			d.config.Metrics = append(d.config.Metrics[:i], d.config.Metrics[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Métrique %s introuvable", name)
	}

	if err := d.persistConfiguration(); err != nil {
		return fmt.Errorf("Impossible de sauvegarder après la suppression de la métrique: %w", err)
	}

	d.logger.Debug("Métrique supprimée avec succès", "name", name)
	return nil
}
