package vm

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-libvirt/sdk/go/libvirt"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreatePulumiVM(config *VMConfig) {
	pulumi.Run(func(ctx *pulumi.Context) error {
		pool, err := libvirt.NewPool(ctx, "default", &libvirt.PoolArgs{
			Name: pulumi.String("default"),
			Type: pulumi.String("dir"),
			Path: pulumi.String("/var/lib/libvirt/images"),
		})
		if err != nil {
			return fmt.Errorf("failed to get/create pool: %w", err)
		}

		volume, err := libvirt.NewVolume(ctx, config.Name+"-base", &libvirt.VolumeArgs{
			Name:   pulumi.String(config.Name + ".qcow2"),
			Pool:   pool.Name,
			Source: pulumi.String(config.ImageURL),
			Format: pulumi.String("qcow2"),
		})
		if err != nil {
			return fmt.Errorf("failed to create base volume: %w", err)
		}

		ignData, err := generateIgnitionJSON(config)
		if err != nil {
			return fmt.Errorf("failed to generate ignition: %w", err)
		}

		ignition, err := libvirt.NewIgnition(ctx, config.Name+"-ign", &libvirt.IgnitionArgs{
			Name:    pulumi.String(config.Name + "-ignition"),
			Pool:    pool.Name,
			Content: pulumi.String(string(ignData)),
		})
		if err != nil {
			return fmt.Errorf("failed to upload ignition: %w", err)
		}

		_, err = libvirt.NewDomain(ctx, config.Name, &libvirt.DomainArgs{
			Name:      pulumi.String(config.Name),
			Memory:    pulumi.IntPtr(int(config.Memory)),
			Vcpu:      pulumi.IntPtr(int(config.VCPUs)),
			QemuAgent: pulumi.Bool(true),
			Disks: libvirt.DomainDiskArray{
				libvirt.DomainDiskArgs{
					VolumeId: volume.ID(),
					Scsi:     pulumi.Bool(true),
				},
			},
			NetworkInterfaces: libvirt.DomainNetworkInterfaceArray{
				libvirt.DomainNetworkInterfaceArgs{
					NetworkName:  pulumi.String(config.Network),
					WaitForLease: pulumi.Bool(true),
				},
			},
			CoreosIgnition: ignition.ID().ToStringOutput().ToStringPtrOutput(),
		})
		if err != nil {
			return fmt.Errorf("failed to create domain: %w", err)
		}

		ctx.Export("vmName", pulumi.String(config.Name))
		return nil
	})
}

func generateIgnitionJSON(config *VMConfig) ([]byte, error) {
	ign := map[string]interface{}{
		"ignition": map[string]string{"version": "3.4.0"},
		"passwd": map[string]interface{}{
			"users": []map[string]interface{}{
				{
					"name":              "core",
					"sshAuthorizedKeys": []string{config.SSHKey},
					"groups":            []string{"sudo", "docker"},
				},
			},
		},
	}
	return json.MarshalIndent(ign, "", "  ")
}
