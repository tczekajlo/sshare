package driver

import (
	"sshare/driver/k8s"
	"sshare/types"

	"github.com/spf13/viper"
)

type Driver struct{}

func (d *Driver) New() types.DriverAdapter {

	switch source := viper.GetString("driver"); source {
	case "kubernetes":
		return k8s.New()
	default:
		return nil
	}
}
