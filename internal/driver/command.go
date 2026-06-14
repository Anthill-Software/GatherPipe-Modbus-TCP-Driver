package driver

import (
	"fmt"
	"strconv"
	"strings"

	GatherPipe "github.com/Anthill-Software/GatherPipe/core"
)

func (driver *ModbusTCPDriver) SupportedCommands() ([]GatherPipe.CommandArg, error) {
	return []GatherPipe.CommandArg{
		{
			Name:        "status",
			Usage:       "status",
			Description: "Affiche l'état du driver et le résumé de la configuration",
		},
		{
			Name:        "set-endpoint",
			Usage:       "set-endpoint [host:port]",
			Description: "Met à jour l'adresse IP et le port du serveur Modbus (ex: 192.168.1.10:502)",
		},
		{
			Name:        "list-metrics",
			Usage:       "list-metrics",
			Description: "Liste toutes les métriques Modbus actuellement configurées",
		},
		{
			Name:        "add-metric",
			Usage:       "add-metric [name] [register_type] [address] [data_type] [scale] [unit]",
			Description: "Ajoute une nouvelle métrique (ex: add-metric temp holding 0 int16 0.1 °C)",
		},
		{
			Name:        "del-metric",
			Usage:       "del-metric [name]",
			Description: "Supprime une métrique de la configuration par son nom",
		},
	}, nil
}

func (driver *ModbusTCPDriver) ExecuteCommand(cmd string, args []string) (string, error) {
	switch cmd {

	case "status":
		metricsCount := len(driver.ListMetrics())
		return fmt.Sprintf(
			"=== Statut Driver Modbus TCP ===\n"+
				"Configuration : %s\n"+
				"Endpoint      : %s\n"+
				"Timeout       : %s\n"+
				"Métriques     : %d",
			driver.configPath, driver.config.Endpoint, driver.config.Timeout, metricsCount,
		), nil

	case "set-endpoint":
		if len(args) < 1 {
			return "", fmt.Errorf("Argument manquant: l'endpoint au format host:port est requis")
		}
		newEndpoint := args[0]
		if err := driver.SetEndPoint(newEndpoint); err != nil {
			return "", fmt.Errorf("Echec de la mise à jour de l'endpoint: %w", err)
		}
		return fmt.Sprintf("Succès : l'endpoint est maintenant %s", newEndpoint), nil

	case "list-metrics":
		metrics := driver.ListMetrics()
		if len(metrics) == 0 {
			return "Aucune métrique configurée.", nil
		}
		var sb strings.Builder
		sb.WriteString("=== Métriques configurées ===\n")
		for _, m := range metrics {
			sb.WriteString(fmt.Sprintf("- %s | %s register | Adresse: %d | Type: %s | Scale: %g | Unité: %s\n",
				m.Name, m.RegisterType, m.Address, m.DataType, m.Scale, m.Unit))
		}
		return strings.TrimSuffix(sb.String(), "\n"), nil

	case "add-metric":
		// Syntaxe attendue : Name, SlaveID, RegisterType, Address, DataType, Scale, Unit
		if len(args) < 6 {
			return "", fmt.Errorf("Arguments manquants. Usage: AddMetric [slave_id] [name] [register_type] [address] [data_type] [scale] [unit]")
		}

		// Parsing du slave id
		var slaveID byte
		if _, err := fmt.Sscanf(args[1], "%d", &slaveID); err != nil {
			return "", fmt.Errorf("slave_id invalide (0-255): %w", err)
		}

		// Parsing du numéro d'adresse
		addr, err := strconv.ParseUint(args[3], 10, 16)
		if err != nil {
			return "", fmt.Errorf("Adresse invalide (doit être un entier 16-bits): %w", err)
		}

		// Parsing du scale (float64)
		scale, err := strconv.ParseFloat(args[5], 64)
		if err != nil {
			return "", fmt.Errorf("Scale invalide (doit être un nombre décimal): %w", err)
		}

		newMetric := MetricConfig{
			Name:         args[0],
			SlaveID:      slaveID,
			RegisterType: args[2], // holding, input, etc.
			Address:      uint16(addr),
			DataType:     args[4], // int16, float32, etc.
			Scale:        scale,
			Unit:         args[6],
		}

		if err := driver.AddMetric(newMetric); err != nil {
			return "", fmt.Errorf("Impossible d'ajouter la métrique: %w", err)
		}
		return fmt.Sprintf("Succès : métrique '%s' ajoutée et sauvegardée", newMetric.Name), nil

	case "del-metric":
		if len(args) < 1 {
			return "", fmt.Errorf("Argument manquant: le nom de la métrique à supprimer est requis")
		}
		targetName := args[0]
		if err := driver.DelMetric(targetName); err != nil {
			return "", fmt.Errorf("échec de la suppression: %w", err)
		}
		return fmt.Sprintf("Succès : métrique '%s' supprimée de la configuration", targetName), nil

	default:
		return "", fmt.Errorf("Commande '%s' inconnue pour le plugin Modbus", cmd)
	}
}
