package service

import (
	"strings"

	"github.com/godeps/newapi/common"
	"github.com/godeps/newapi/setting/system_setting"
)

func PaymentReturnURL(suffix string) string {
	base := strings.TrimRight(system_setting.ServerAddress, "/")
	return base + common.ThemeAwarePath(suffix)
}
