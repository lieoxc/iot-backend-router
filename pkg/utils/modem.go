package utils

import (
	"context"
	"errors"
	"fmt"
	"project/pkg/global"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

// GPSData 结构体保存解析后的 GPS 信息
type modemData struct {
	signal       string
	modelVersion string
	network      string
}

var globalModemData modemData

// 打开 GPS 功能的函数
func enableGPS(portName string) error {
	port, err := serial.Open(portName, &serial.Mode{
		BaudRate: 9600,
	})
	if err != nil {
		return fmt.Errorf("无法打开串口 %s: %v", portName, err)
	}
	defer port.Close()

	// 发送开启 GPS 的指令
	SendATCommand(portName, "AT+QGPS=1\r\n")
	return nil
}
func UpdateModemInfoToRedis() {
	global.REDIS.HSet(context.Background(), "modem_data", "signal", globalModemData.signal,
		"modelVersion", globalModemData.modelVersion, "network", globalModemData.network)
}
func modemInfoLoop(portName string, ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			logrus.Debugln("modemInfoLoop Done.")
		case <-ticker.C:
			query4GModemInfo(portName)
		}
	}
}

type modemCommanSt struct {
	command string                       // 寄存器名称
	Handler func(string) (string, error) // 处理读取数据的函数
	Update  func(string)                 // 处理读取数据的函数
}

func query4GModemInfo(portName string) {
	commands := map[string]modemCommanSt{
		"信号强度": {"AT+CSQ\r\n", ParseCSQ, UpdateCsq},
		"模块信息": {"ATI\r\n", ParseRevision, UpdateRevision},
		"网络制式": {"AT+COPS?\r\n", ParseOperator, UpdateNetwork},
	}
	for name, cmd := range commands {
		//logrus.Debugf("query: %s...", name)
		resp, err := SendATCommand(portName, cmd.command)
		if err != nil {
			logrus.Debugf("%s 查询失败: %v", name, err)
			continue
		}
		//logrus.Debugf("resp: %s", resp)
		result, err := cmd.Handler(resp)
		if err != nil {
			//logrus.Debugf("%s 解析失败: %v", name, err)
			continue
		}
		//logrus.Debugf("update: %s", result)
		cmd.Update(result)
		//logrus.Debugf("end --: %s", result)

	}
	UpdateModemInfoToRedis()
}

// SendATCommand 发送AT指令到指定串口并获取响应
// 参数：
//   portName: 串口设备路径 (如 "/dev/ttyUSB2")
//   baudRate: 波特率 (如 9600)
//   command:  要发送的AT指令 (需要包含回车换行符，如 "AT+CSQ\r\n")
//   timeout:  读取超时时间
// 返回值：
//   string: 完整响应内容
//   error:  错误信息（成功时为nil）
func SendATCommand(portName string, command string) (string, error) {
	// 配置串口参数
	port, err := serial.Open(portName, &serial.Mode{
		BaudRate: 9600,
	})
	defer port.Close() // 确保函数返回前关闭串口

	// 发送AT指令
	_, err = port.Write([]byte(command))
	if err != nil {
		return "", fmt.Errorf("发送指令失败: %v", err)
	}
	// 设置串口读写超时（500ms间隔检查）
	port.SetReadTimeout(1000 * time.Millisecond)
	// 读取响应
	buf := make([]byte, 128)
	response := ""
	for {
		n, err := port.Read(buf)
		if err != nil {
			return response, fmt.Errorf("读取数据失败: %v", err)
		}
		if n == 0 {
			break // 无数据可读时退出
		}
		response += string(buf[:n])

		// 检查终止条件：收到OK或ERROR（根据模块实际响应调整）
		if len(response) >= 4 {
			last4 := response[len(response)-4:]
			if last4 == "OK\r\n" || last4 == "ERROR\r\n" {
				break
			}
		}
	}

	return response, nil
}
func UpdateCsq(csq string) {
	globalModemData.signal = csq
}
func UpdateRevision(revision string) {
	globalModemData.modelVersion = revision
}
func UpdateNetwork(network string) {
	globalModemData.network = network
}

func ParseCSQ(response string) (string, error) {
	// 使用正则表达式匹配模式
	re := regexp.MustCompile(`\+CSQ:\s*(\d+),`)
	matches := re.FindStringSubmatch(response)
	if len(matches) < 2 {
		return "", errors.New("未找到有效的信号强度数据")
	}
	return matches[1], nil
}

// ParseRevision 解析模块版本信息
// 示例输入：
// Quectel
// EC20F
// Revision: EC20CEFILGR06A03M1G
// OK
func ParseRevision(response string) (string, error) {
	// 按行分割响应内容
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		// 查找包含Revision的行
		if strings.Contains(line, "Revision:") {
			// 分割键值对
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}
			// 去除前后空白字符
			version := strings.TrimSpace(parts[1])
			if version == "" {
				return "", errors.New("版本号为空")
			}
			return version, nil
		}
	}
	return "", errors.New("未找到版本信息")
}

// ParseOperator 解析运营商信息
// 示例输入：
// +COPS: 0,0,"CHINA MOBILE",7
// OK
func ParseOperator(response string) (string, error) {
	// 使用正则表达式匹配带引号的运营商名称
	re := regexp.MustCompile(`\+COPS:\s*\d+,\d+,"([^"]+)"`)
	matches := re.FindStringSubmatch(response)
	if len(matches) < 2 {
		// 尝试匹配不带引号的格式（某些模块可能不同）
		reAlt := regexp.MustCompile(`\+COPS:\s*\d+,\d+,([^,]+)`)
		matchesAlt := reAlt.FindStringSubmatch(response)
		if len(matchesAlt) < 2 {
			return "", errors.New("未找到运营商信息")
		}
		return strings.TrimSpace(matchesAlt[1]), nil
	}

	// 去除前后空白字符
	operator := strings.TrimSpace(matches[1])
	if operator == "" {
		return "", errors.New("运营商名称为空")
	}

	return operator, nil
}
