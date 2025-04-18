package config

import (
	"fmt"
	"vmctl/pkg/vm"

	"github.com/spf13/viper"
)

func Load() (*vm.VMConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/homelab/")

	fmt.Println("🔍 Loading configuration from config.yaml...")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("❌ error reading config file: %w", err)
	}

	var config vm.VMConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("❌ unable to decode into struct: %w", err)
	}

	fmt.Println("✅ Config loaded successfully.")
	return &config, nil
}
