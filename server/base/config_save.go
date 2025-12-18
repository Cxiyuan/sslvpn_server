package base

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// SaveAdminPass 保存管理员密码到配置文件
func SaveAdminPass(newPassHash string) error {
	configFile := Cfg.Conf
	if configFile == "" {
		return fmt.Errorf("配置文件路径为空")
	}

	// 读取配置文件
	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("打开配置文件失败: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	updated := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// 查找 admin_pass 配置行
		if strings.HasPrefix(trimmed, "admin_pass") && !strings.HasPrefix(trimmed, "#") {
			// 替换为新密码
			lines = append(lines, fmt.Sprintf("admin_pass = \"%s\"", newPassHash))
			updated = true
		} else {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	if !updated {
		return fmt.Errorf("配置文件中未找到 admin_pass 配置项")
	}

	// 写回配置文件
	output := strings.Join(lines, "\n")
	err = os.WriteFile(configFile, []byte(output), 0600)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	Info("管理员密码已保存到配置文件:", configFile)
	return nil
}
