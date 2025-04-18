package main

import (
	"fmt"
	"vmctl/pkg/config"
	"vmctl/pkg/vm"
)

func main() {
	config, err := config.Load()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Provisioning VM: %s\n", config.Name)

	vm.CreatePulumiVM(config)
}
