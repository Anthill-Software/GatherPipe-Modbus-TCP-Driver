# GatherPipe Modbus-TCP Driver Plugin

Ce plugin permet à GatherPipe de collecter des données à partir d'équipements utilisant le protocole Modbus-TCP (automates, stations météo, capteurs industriels) sur un bus partagé.

## Caractéristiques

- Support multi-esclaves (Slave ID / Unit ID configurable par métrique).

- Types de données supportés : `uint16`, `int16`, `int32`, `float32`.

- Persistance automatique de la configuration en YAML.

- Configuration via la console SSH de GatherPipe.

## Types de Registres

- `holding` : Holding Registers (Lecture)

- `input` : Input Registers (Lecture)

## Configuration Console (CLI)

Une fois le plugin chargé, les commandes suivantes sont disponibles dans la console GatherPipe :

| Commande | Description | Exemple / Usage |
| :--- | :--- | :--- |
| `status` | État de la connexion et fichier config | `config plugin modbus-driver status` |
| `set-endpoint` | Modifie l'adresse IP et le port | `config plugin modbus-driver set-endpoint 192.168.1.50:502` |
| `list-metrics` | Liste les métriques actives | `config plugin modbus-driver list-metrics` |
| `add-metric` | Ajoute un capteur au cycle de sonde | `config plugin modbus-driver add-metric outTemp 10 holding 1 int16 0.1 °C` |
| `del-metric` | Supprime une métrique par son nom | `config plugin modbus-driver del-metric outTemp` |
