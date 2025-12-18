package base

import (
	"errors"
	"fmt"
	"net"
)

// ValidateConfig 验证配置的合法性
func ValidateConfig(cfg *ServerConfig) error {
	// 验证网络配置
	if err := validateNetworkConfig(cfg); err != nil {
		return fmt.Errorf("网络配置错误: %w", err)
	}

	// 验证连接数限制
	if err := validateConnectionLimits(cfg); err != nil {
		return fmt.Errorf("连接限制配置错误: %w", err)
	}

	// 验证时间配置
	if err := validateTimeouts(cfg); err != nil {
		return fmt.Errorf("超时配置错误: %w", err)
	}

	// 验证数据库配置
	if err := validateDatabaseConfig(cfg); err != nil {
		return fmt.Errorf("数据库配置错误: %w", err)
	}

	return nil
}

func validateNetworkConfig(cfg *ServerConfig) error {
	// 验证 CIDR
	_, ipnet, err := net.ParseCIDR(cfg.Ipv4CIDR)
	if err != nil {
		return fmt.Errorf("ipv4_cidr 格式错误: %s", cfg.Ipv4CIDR)
	}

	// 验证 Gateway
	gateway := net.ParseIP(cfg.Ipv4Gateway)
	if gateway == nil {
		return fmt.Errorf("ipv4_gateway 格式错误: %s", cfg.Ipv4Gateway)
	}
	if !ipnet.Contains(gateway) {
		return fmt.Errorf("ipv4_gateway 不在 ipv4_cidr 范围内")
	}

	// 验证 Start IP
	startIP := net.ParseIP(cfg.Ipv4Start)
	if startIP == nil {
		return fmt.Errorf("ipv4_start 格式错误: %s", cfg.Ipv4Start)
	}
	if !ipnet.Contains(startIP) {
		return fmt.Errorf("ipv4_start 不在 ipv4_cidr 范围内")
	}

	// 验证 End IP
	endIP := net.ParseIP(cfg.Ipv4End)
	if endIP == nil {
		return fmt.Errorf("ipv4_end 格式错误: %s", cfg.Ipv4End)
	}
	if !ipnet.Contains(endIP) {
		return fmt.Errorf("ipv4_end 不在 ipv4_cidr 范围内")
	}

	// 验证 Start <= End
	if compareIP(startIP, endIP) > 0 {
		return fmt.Errorf("ipv4_start 必须小于等于 ipv4_end")
	}

	// 验证 MTU
	if cfg.Mtu < 576 || cfg.Mtu > 9000 {
		return fmt.Errorf("mtu 必须在 576-9000 之间，当前值: %d", cfg.Mtu)
	}

	// 验证 link_mode
	validModes := []string{LinkModeTUN, LinkModeTAP, LinkModeMacvtap, LinkModeIpvtap}
	if !contains(validModes, cfg.LinkMode) {
		return fmt.Errorf("link_mode 必须是 %v 之一，当前值: %s", validModes, cfg.LinkMode)
	}

	return nil
}

func validateConnectionLimits(cfg *ServerConfig) error {
	if cfg.MaxClient <= 0 {
		return fmt.Errorf("max_client 必须大于 0，当前值: %d", cfg.MaxClient)
	}

	if cfg.MaxUserClient <= 0 {
		return fmt.Errorf("max_user_client 必须大于 0，当前值: %d", cfg.MaxUserClient)
	}

	if cfg.MaxUserClient > cfg.MaxClient {
		return fmt.Errorf("max_user_client(%d) 不能大于 max_client(%d)", cfg.MaxUserClient, cfg.MaxClient)
	}

	return nil
}

func validateTimeouts(cfg *ServerConfig) error {
	// IP租期至少60秒
	if cfg.IpLease < 60 {
		return fmt.Errorf("ip_lease 必须至少为 60 秒，当前值: %d", cfg.IpLease)
	}

	// keepalive必须大于0
	if cfg.CstpKeepalive <= 0 {
		return fmt.Errorf("cstp_keepalive 必须大于 0，当前值: %d", cfg.CstpKeepalive)
	}

	if cfg.MobileKeepalive <= 0 {
		return fmt.Errorf("mobile_keepalive 必须大于 0，当前值: %d", cfg.MobileKeepalive)
	}

	// DPD必须大于keepalive
	if cfg.CstpDpd <= cfg.CstpKeepalive {
		return fmt.Errorf("cstp_dpd(%d) 必须大于 cstp_keepalive(%d)", cfg.CstpDpd, cfg.CstpKeepalive)
	}

	if cfg.MobileDpd <= cfg.MobileKeepalive {
		return fmt.Errorf("mobile_dpd(%d) 必须大于 mobile_keepalive(%d)", cfg.MobileDpd, cfg.MobileKeepalive)
	}

	// 空闲超时和会话超时可以为0（表示禁用）
	if cfg.IdleTimeout < 0 {
		return fmt.Errorf("idle_timeout 不能为负数，当前值: %d", cfg.IdleTimeout)
	}

	if cfg.SessionTimeout < 0 {
		return fmt.Errorf("session_timeout 不能为负数，当前值: %d", cfg.SessionTimeout)
	}

	return nil
}

func validateDatabaseConfig(cfg *ServerConfig) error {
	validDBTypes := []string{"sqlite3", "mysql", "postgres"}
	if !contains(validDBTypes, cfg.DbType) {
		return fmt.Errorf("db_type 必须是 %v 之一，当前值: %s", validDBTypes, cfg.DbType)
	}

	if cfg.DbSource == "" {
		return fmt.Errorf("db_source 不能为空")
	}

	return nil
}

// compareIP 比较两个IP地址的大小
func compareIP(ip1, ip2 net.IP) int {
	ip1 = ip1.To4()
	ip2 = ip2.To4()
	
	for i := 0; i < len(ip1); i++ {
		if ip1[i] < ip2[i] {
			return -1
		}
		if ip1[i] > ip2[i] {
			return 1
		}
	}
	return 0
}

// contains 检查字符串切片是否包含指定字符串
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// ValidateConfigUpdate 验证可热更新的配置项
func ValidateConfigUpdate(key string, value interface{}) error {
	// 检查是否为敏感配置
	for _, cfg := range configs {
		if cfg.Name == key && cfg.Sensitive {
			return errors.New("此配置项为敏感配置，不允许热更新")
		}
	}

	// 根据配置项进行特定验证
	switch key {
	case "ipv4_cidr":
		if strVal, ok := value.(string); ok {
			if _, _, err := net.ParseCIDR(strVal); err != nil {
				return fmt.Errorf("ipv4_cidr 格式错误: %s", strVal)
			}
		}
	case "ipv4_gateway", "ipv4_start", "ipv4_end":
		if strVal, ok := value.(string); ok {
			if net.ParseIP(strVal) == nil {
				return fmt.Errorf("%s 格式错误: %s", key, strVal)
			}
		}
	case "mtu":
		if intVal, ok := value.(int); ok {
			if intVal < 576 || intVal > 9000 {
				return fmt.Errorf("mtu 必须在 576-9000 之间")
			}
		}
	case "max_client", "max_user_client":
		if intVal, ok := value.(int); ok {
			if intVal <= 0 {
				return fmt.Errorf("%s 必须大于 0", key)
			}
		}
	case "ip_lease":
		if intVal, ok := value.(int); ok {
			if intVal < 60 {
				return fmt.Errorf("ip_lease 必须至少为 60 秒")
			}
		}
	case "link_mode":
		if strVal, ok := value.(string); ok {
			validModes := []string{LinkModeTUN, LinkModeTAP, LinkModeMacvtap, LinkModeIpvtap}
			if !contains(validModes, strVal) {
				return fmt.Errorf("link_mode 必须是 %v 之一", validModes)
			}
		}
	}

	return nil
}