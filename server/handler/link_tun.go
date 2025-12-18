package handler

import (
	"fmt"

	"github.com/bjdgyc/anylink/base"
	"github.com/bjdgyc/anylink/pkg/utils"
	"github.com/bjdgyc/anylink/sessdata"
	"github.com/coreos/go-iptables/iptables"
	gosysctl "github.com/lorenzosaino/go-sysctl"
	"github.com/songgao/water"
)

const (
	iptablesComment = "SSLVPN"
)

func cleanupIptablesRules() error {
	ipt, err := iptables.New()
	if err != nil {
		return err
	}

	base.Info("清理旧的iptables规则...")

	natRules, err := ipt.List("nat", "POSTROUTING")
	if err != nil {
		base.Warn("获取NAT规则失败:", err)
	} else {
		for _, rule := range natRules {
			if len(rule) > 0 && !isChainRule(rule) && containsComment(rule, iptablesComment) {
				ruleSpec := parseRuleSpec(rule)
				if len(ruleSpec) > 0 {
					err = ipt.Delete("nat", "POSTROUTING", ruleSpec...)
					if err != nil {
						base.Warn("删除NAT规则失败:", err)
					} else {
						base.Info("已删除NAT规则:", rule)
					}
				}
			}
		}
	}

	forwardRules, err := ipt.List("filter", "FORWARD")
	if err != nil {
		base.Warn("获取FORWARD规则失败:", err)
	} else {
		for _, rule := range forwardRules {
			if len(rule) > 0 && !isChainRule(rule) && containsComment(rule, iptablesComment) {
				ruleSpec := parseRuleSpec(rule)
				if len(ruleSpec) > 0 {
					err = ipt.Delete("filter", "FORWARD", ruleSpec...)
					if err != nil {
						base.Warn("删除FORWARD规则失败:", err)
					} else {
						base.Info("已删除FORWARD规则:", rule)
					}
				}
			}
		}
	}

	return nil
}

func isChainRule(rule string) bool {
	return len(rule) > 0 && (rule[0] == '-' && rule[1] == 'P' || rule[0] == '-' && rule[1] == 'N')
}

func containsComment(rule, comment string) bool {
	commentFlag := "/* " + comment + " */"
	for i := 0; i < len(rule)-len(commentFlag)+1; i++ {
		if rule[i:i+len(commentFlag)] == commentFlag {
			return true
		}
	}
	return false
}

func parseRuleSpec(rule string) []string {
	if len(rule) == 0 || rule[0] != '-' {
		return nil
	}
	
	var specs []string
	var current string
	inQuote := false
	
	for i := 0; i < len(rule); i++ {
		c := rule[i]
		if c == '"' {
			inQuote = !inQuote
		} else if c == ' ' && !inQuote {
			if len(current) > 0 {
				specs = append(specs, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if len(current) > 0 {
		specs = append(specs, current)
	}
	
	return specs
}

func checkTun() {
	base.CheckModOrLoad("tun")

	cfg := water.Config{
		DeviceType: water.TUN,
	}

	ifce, err := water.New(cfg)
	if err != nil {
		base.Fatal("open tun err: ", err)
	}
	defer ifce.Close()

	cmdstr1 := fmt.Sprintf("ip link set dev %s up mtu %s multicast off", ifce.Name(), "1399")
	err = execCmd([]string{cmdstr1})
	if err != nil {
		base.Fatal("testTun err: ", err)
	}

	if base.Cfg.IptablesNat {
		err = cleanupIptablesRules()
		if err != nil {
			base.Error("清理iptables规则失败:", err)
		}

		ipt, err := iptables.New()
		if err != nil {
			base.Fatal(err)
			return
		}

		base.CheckModOrLoad("iptable_filter")
		base.CheckModOrLoad("iptable_nat")

		natRule := []string{"-s", base.Cfg.Ipv4CIDR, "-o", base.Cfg.Ipv4Master, "-m", "comment",
			"--comment", iptablesComment, "-j", "MASQUERADE"}
		if base.InContainer {
			natRule = []string{"-s", base.Cfg.Ipv4CIDR, "-o", base.Cfg.Ipv4Master, "-j", "MASQUERADE"}
		}
		err = ipt.AppendUnique("nat", "POSTROUTING", natRule...)
		if err != nil {
			base.Error(err)
		}

		forwardRule := []string{"-m", "comment", "--comment", iptablesComment, "-j", "ACCEPT"}
		if base.InContainer {
			forwardRule = []string{"-j", "ACCEPT"}
		}
		err = ipt.AppendUnique("filter", "FORWARD", forwardRule...)
		if err != nil {
			base.Error(err)
		}

		base.Info(ipt.List("nat", "POSTROUTING"))
		base.Info(ipt.List("filter", "FORWARD"))
	}
}

// 创建tun网卡
func LinkTun(cSess *sessdata.ConnSession) error {
	cfg := water.Config{
		DeviceType: water.TUN,
	}

	ifce, err := water.New(cfg)
	if err != nil {
		base.Error(err)
		return err
	}
	// log.Printf("Interface Name: %s\n", ifce.Name())
	cSess.SetIfName(ifce.Name())

	// 通过 ip link show  查看 alias 信息
	alias := utils.ParseName(cSess.Group.Name + "." + cSess.Username)
	cmdstr1 := fmt.Sprintf("ip link set dev %s up mtu %d multicast off alias %s", ifce.Name(), cSess.Mtu, alias)
	cmdstr2 := fmt.Sprintf("ip addr add dev %s local %s peer %s/32",
		ifce.Name(), base.Cfg.Ipv4Gateway, cSess.IpAddr)
	err = execCmd([]string{cmdstr1, cmdstr2})
	if err != nil {
		base.Error(err)
		_ = ifce.Close()
		return err
	}

	// cmdstr3 := fmt.Sprintf("sysctl -w net.ipv6.conf.%s.disable_ipv6=1", ifce.Name())
	// execCmd([]string{cmdstr3})
	err = gosysctl.Set(fmt.Sprintf("net.ipv6.conf.%s.disable_ipv6", ifce.Name()), "1")
	if err != nil {
		base.Warn(err)
	}

	go tunRead(ifce, cSess)
	go tunWrite(ifce, cSess)
	return nil
}

func tunWrite(ifce *water.Interface, cSess *sessdata.ConnSession) {
	defer func() {
		base.Debug("LinkTun return", cSess.IpAddr)
		cSess.Close()
		_ = ifce.Close()
	}()

	var (
		err error
		pl  *sessdata.Payload
	)

	for {
		select {
		case pl = <-cSess.PayloadIn:
		case <-cSess.CloseChan:
			return
		}

		_, err = ifce.Write(pl.Data)
		if err != nil {
			base.Error("tun Write err", err)
			return
		}

		putPayloadInBefore(cSess, pl)
	}
}

func tunRead(ifce *water.Interface, cSess *sessdata.ConnSession) {
	defer func() {
		base.Debug("tunRead return", cSess.IpAddr)
		_ = ifce.Close()
	}()
	var (
		err error
		n   int
	)

	for {
		// data := make([]byte, BufferSize)
		pl := getPayload()
		n, err = ifce.Read(pl.Data)
		if err != nil {
			base.Error("tun Read err", n, err)
			return
		}

		// 更新数据长度
		pl.Data = (pl.Data)[:n]

		// data = data[:n]
		// ip_src := waterutil.IPv4Source(data)
		// ip_dst := waterutil.IPv4Destination(data)
		// ip_port := waterutil.IPv4DestinationPort(data)
		// fmt.Println("sent:", ip_src, ip_dst, ip_port)
		// packet := gopacket.NewPacket(data, layers.LayerTypeIPv4, gopacket.Default)
		// fmt.Println("read:", packet)

		if payloadOut(cSess, pl) {
			return
		}
	}
}