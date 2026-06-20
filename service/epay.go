package service

import (
	"github.com/godeps/newapi/setting/operation_setting"
	"github.com/godeps/newapi/setting/system_setting"
)

func GetCallbackAddress() string {
	if operation_setting.CustomCallbackAddress == "" {
		return system_setting.ServerAddress
	}
	return operation_setting.CustomCallbackAddress
}
