package cfg

import (
	"log"
	"testing"
)

func TestRead(t *testing.T) {
	cfgManager := NewConfigManager()
	config := cfgManager.ReadConfig()
	log.Println("RCON Endpoint:", config.RCONSettings.Endpoint)
	log.Println("RCON Password:", config.RCONSettings.Password)
	log.Println("Archive time_seconds:", config.ArchiveTimeSeconds)
}
