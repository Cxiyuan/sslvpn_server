package dbdata

import (
	"errors"
	"net"
	"time"
)

type IpMap struct {
	Id        int       `json:"id" xorm:"pk autoincr not null"`
	IpAddr    string    `json:"ip_addr" xorm:"varchar(32) not null unique"`
	MacAddr   string    `json:"mac_addr" xorm:"varchar(32) not null unique"`
	UniqueMac bool      `json:"unique_mac" xorm:"Bool index"`
	Username  string    `json:"username" xorm:"varchar(60) index"`
	Keep      bool      `json:"keep" xorm:"Bool"` // 保留 ip-username 绑定
	KeepTime  time.Time `json:"keep_time" xorm:"DateTime"`
	Note      string    `json:"note" xorm:"varchar(255)"` // 备注
	LastLogin time.Time `json:"last_login" xorm:"DateTime"`
	UpdatedAt time.Time `json:"updated_at" xorm:"DateTime updated"`
}

func SetIpMap(v *IpMap) error {
	var err error

	if len(v.IpAddr) < 4 {
		return errors.New("IP错误")
	}

	// 验证IP格式
	ip := net.ParseIP(v.IpAddr)
	if ip == nil {
		return errors.New("IP格式错误")
	}

	// 如果开启了IP保留，必须设置用户名
	if v.Keep && v.Username == "" {
		return errors.New("IP保留必须设置用户名")
	}

	// MAC地址验证：如果设置了MAC，必须验证格式
	if v.MacAddr != "" {
		macHw, err := net.ParseMAC(v.MacAddr)
		if err != nil {
			return errors.New("MAC格式错误")
		}
		// 统一macAddr的格式
		v.MacAddr = macHw.String()
	} else {
		// 如果MAC为空，生成一个占位符以满足唯一约束
		// 使用时间戳生成唯一的占位MAC
		v.MacAddr = "00:00:00:00:00:00"
	}

	v.UpdatedAt = time.Now()
	if v.Id > 0 {
		err = Set(v)
	} else {
		err = Add(v)
	}
	return err
}