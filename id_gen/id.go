package id_gen

import (
	"fmt"

	"github.com/sony/sonyflake"
)

var (
	sonyFlake     *sonyflake.Sonyflake
	sonyMachineID uint16
)

func getMachineID() (uint16, error) {
	return uint16(sonyMachineID), nil
}

func Init(machineId uint16) (err error) {
	sonyMachineID = machineId
	settings := sonyflake.Settings{}
	settings.MachineID = getMachineID
	sonyFlake = sonyflake.NewSonyflake(settings)
	return
}

func GetId() (id uint64, err error) {
	if sonyFlake == nil {
		err = fmt.Errorf("sony flake not inted")
		return
	}
	id, err = sonyFlake.NextID()
	return
}
