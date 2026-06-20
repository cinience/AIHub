package operation_setting

import (
	"os"
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/config"
)

// 额度展示类型
const (
	QuotaDisplayTypeUSD    = "USD"
	QuotaDisplayTypeCNY    = "CNY"
	QuotaDisplayTypeTokens = "TOKENS"
	QuotaDisplayTypeCustom = "CUSTOM"
)

type GeneralSetting struct {
	DocsLink            string `json:"docs_link"`
	PingIntervalEnabled bool   `json:"ping_interval_enabled"`
	PingIntervalSeconds int    `json:"ping_interval_seconds"`
	// 当前站点额度展示类型：USD / CNY / TOKENS
	QuotaDisplayType string `json:"quota_display_type"`
	// 自定义货币符号，用于 CUSTOM 展示类型
	CustomCurrencySymbol string `json:"custom_currency_symbol"`
	// 自定义货币与美元汇率（1 USD = X Custom）
	CustomCurrencyExchangeRate       float64           `json:"custom_currency_exchange_rate"`
	RequestUpstreamOverrideEnabled   bool              `json:"request_upstream_override_enabled"`
	RequestUpstreamOverrideAllowlist []string          `json:"request_upstream_override_allowlist"`
	RequestUpstreamProxyMap          map[string]string `json:"request_upstream_proxy_map"`
}

// 默认配置
var generalSetting = GeneralSetting{
	DocsLink:                   "https://docs.newapi.pro",
	PingIntervalEnabled:        false,
	PingIntervalSeconds:        60,
	QuotaDisplayType:           QuotaDisplayTypeUSD,
	CustomCurrencySymbol:       "¤",
	CustomCurrencyExchangeRate: 1.0,
	RequestUpstreamProxyMap:    map[string]string{},
}

func init() {
	// 注册到全局配置管理器
	config.GlobalConfig.Register("general_setting", &generalSetting)
}

func GetGeneralSetting() *GeneralSetting {
	applyGeneralSettingEnvOverrides()
	return &generalSetting
}

func IsRequestUpstreamOverrideEnabled() bool {
	return GetGeneralSetting().RequestUpstreamOverrideEnabled
}

func GetRequestUpstreamOverrideAllowlist() []string {
	return GetGeneralSetting().RequestUpstreamOverrideAllowlist
}

func GetRequestUpstreamProxyMap() map[string]string {
	return GetGeneralSetting().RequestUpstreamProxyMap
}

func applyGeneralSettingEnvOverrides() {
	if value, ok := os.LookupEnv("REQUEST_UPSTREAM_OVERRIDE_ENABLED"); ok {
		enabled, err := strconv.ParseBool(value)
		if err == nil {
			generalSetting.RequestUpstreamOverrideEnabled = enabled
		}
	}
	if value, ok := os.LookupEnv("REQUEST_UPSTREAM_OVERRIDE_ALLOWLIST"); ok {
		var allowlist []string
		if err := common.Unmarshal([]byte(value), &allowlist); err == nil {
			generalSetting.RequestUpstreamOverrideAllowlist = allowlist
		}
	}
	if value, ok := os.LookupEnv("REQUEST_UPSTREAM_PROXY_MAP"); ok {
		var proxyMap map[string]string
		if err := common.Unmarshal([]byte(value), &proxyMap); err == nil {
			generalSetting.RequestUpstreamProxyMap = proxyMap
		}
	}
}

// IsCurrencyDisplay 是否以货币形式展示（美元或人民币）
func IsCurrencyDisplay() bool {
	return generalSetting.QuotaDisplayType != QuotaDisplayTypeTokens
}

// IsCNYDisplay 是否以人民币展示
func IsCNYDisplay() bool {
	return generalSetting.QuotaDisplayType == QuotaDisplayTypeCNY
}

// GetQuotaDisplayType 返回额度展示类型
func GetQuotaDisplayType() string {
	return generalSetting.QuotaDisplayType
}

// GetCurrencySymbol 返回当前展示类型对应符号
func GetCurrencySymbol() string {
	switch generalSetting.QuotaDisplayType {
	case QuotaDisplayTypeUSD:
		return "$"
	case QuotaDisplayTypeCNY:
		return "¥"
	case QuotaDisplayTypeCustom:
		if generalSetting.CustomCurrencySymbol != "" {
			return generalSetting.CustomCurrencySymbol
		}
		return "¤"
	default:
		return ""
	}
}

// GetUsdToCurrencyRate 返回 1 USD = X <currency> 的 X（TOKENS 不适用）
func GetUsdToCurrencyRate(usdToCny float64) float64 {
	switch generalSetting.QuotaDisplayType {
	case QuotaDisplayTypeUSD:
		return 1
	case QuotaDisplayTypeCNY:
		return usdToCny
	case QuotaDisplayTypeCustom:
		if generalSetting.CustomCurrencyExchangeRate > 0 {
			return generalSetting.CustomCurrencyExchangeRate
		}
		return 1
	default:
		return 1
	}
}
