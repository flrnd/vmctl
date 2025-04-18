package vm

type VMConfig struct {
	Name     string `mapstructure:"name"`
	Memory   int    `mapstructure:"memory"`
	VCPUs    int    `mapstructure:"vcpus"`
	DiskSize int    `mapstructure:"disk_size_gb"`
	Network  string `mapstructure:"network"`
	SSHKey   string `mapstructure:"ssh_key"`
	ImageURL string `mapstructure:"image_url"`
	OSType   string `mapstructure:"os_type"`
}
