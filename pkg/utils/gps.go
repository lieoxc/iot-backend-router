package utils

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

// 全局变量定义
var (
	currentGPSData GPSData
	gpsDataMutex   sync.RWMutex
)

// 安全更新全局GPS数据
func updateGlobalData(data GPSData) {
	gpsDataMutex.Lock()
	defer gpsDataMutex.Unlock()

	currentGPSData = data
}

// GetCurrentGPSData 安全获取当前GPS数据
func GetCurrentGPSData() GPSData {
	gpsDataMutex.RLock()
	defer gpsDataMutex.RUnlock()
	return currentGPSData
}

// GPSData 结构体保存解析后的 GPS 信息
type GPSData struct {
	UtcTime      time.Time // UTC 时间
	LocalTimeStr string
	Latitude     float64 // 纬度
	Longitude    float64 // 经度
}

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
	command := "AT+QGPS=1\r\n" // 注意指令中包含换行符
	_, err = port.Write([]byte(command))
	if err != nil {
		return fmt.Errorf("发送开启 GPS 指令失败: %v", err)
	}
	return nil
}

// 读取 GPS 数据的函数
func readGPSData() (GPSData, error) {
	return GetCurrentGPSData(), nil
}

// 解析 GPRMC 消息以获取 UTC 时间和经纬度
func parseGPRMC(sentence string) (GPSData, error) {
	// 检查是否为 GPRMC 消息
	if !strings.HasPrefix(sentence, "$GPRMC") {
		return GPSData{}, fmt.Errorf("不是有效的 GPRMC 消息: %s", sentence)
	}

	// 拆分字段
	fields := strings.Split(sentence, ",")
	if len(fields) < 10 {
		return GPSData{}, fmt.Errorf("GPRMC fields no finsh")
	}

	// 提取 UTC 时间（hhmmss.sss）
	rawTime := fields[1] // 格式: HHMMSS.SS
	if len(rawTime) < 7 {
		return GPSData{}, fmt.Errorf("GPS  rawTime no finsh")
	}
	hh, mm, ss := rawTime[:2], rawTime[2:4], rawTime[4:6]

	// 提取日期（DDMMYY）
	rawDate := fields[9] // 格式: DDMMYY
	dd, month, year := rawDate[:2], rawDate[2:4], rawDate[4:6]

	// 将年份补全到 4 位
	fullYear := "20" + year

	// 构建完整的 UTC 时间字符串
	utcTimeStr := fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", fullYear, month, dd, hh, mm, ss)

	// 解析为 UTC 时间
	utcTime, err := time.Parse(time.RFC3339, utcTimeStr)
	if err != nil {
		return GPSData{}, fmt.Errorf("无法解析 UTC 时间: %v", err)
	}

	// 提取纬度和经度
	latitude, err := parseCoordinate(fields[3], fields[4], 2)
	if err != nil {
		return GPSData{}, fmt.Errorf("纬度解析失败: %v", err)
	}

	longitude, err := parseCoordinate(fields[5], fields[6], 3)
	if err != nil {
		return GPSData{}, fmt.Errorf("经度解析失败: %v", err)
	}

	// 返回封装好的 GPS 数据
	return GPSData{
		UtcTime:   utcTime,
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}
func GPSInit() error {
	controlPort := "/dev/ttyUSB2" // 用于发送 AT 指令的串口
	logrus.Debugln("Start GPS Init.")
	if err := enableGPS(controlPort); err != nil {
		logrus.Errorln("GPS Init Failed.", err)
		return err
	}
	logrus.Debugln("Finsh GPS Init.")
	// 配置串口名称（注意替换为实际串口）
	gpsPort := "/dev/ttyUSB1" // 用于接收 GPS 数据的串口
	go gpsReadloop(gpsPort, context.Background())
	return nil
}
func gpsReadloop(portName string, ctx context.Context) {
	port, err := serial.Open(portName, &serial.Mode{BaudRate: 9600})
	if err != nil {
		logrus.Errorf("无法打开串口 %s: %v", portName, err)
		return
	}
	defer port.Close()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	reader := bufio.NewReader(port)
	for {
		select {
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				logrus.Errorf("读取串口数据时出错: %v", err)
				break
			}

			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "$GPRMC") {
				data, err := parseGPRMC(line)
				if err != nil {
					logrus.Errorf("解析 GPRMC 消息失败: %v", err)
					break
				}
				localTime := data.UtcTime.Add(8 * time.Hour)
				data.LocalTimeStr = localTime.Format("2006-01-02 15:04:05")
				updateGlobalData(data)
				break
			}
		}
	}
}
func GetNtpInfo() (GPSData, error) {
	// 开始读取 GPS 数据
	data, err := readGPSData()
	if err != nil {
		logrus.Errorf("读取 GPS 数据失败: %v", err)
		return GPSData{}, err
	}
	return data, nil
}

// 将 NMEA 格式的经纬度转换为十进制格式
func parseCoordinate(value string, hemisphere string, degLen int) (float64, error) {
	// 基本格式校验
	if len(value) < degLen+1 {
		return 0, fmt.Errorf("invalid coordinate format: %s", value)
	}

	// 分离度数和分钟
	degStr := value[:degLen]
	minStr := value[degLen:]

	// 转换为浮点数
	degrees, err := strconv.ParseFloat(degStr, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid degrees: %s", degStr)
	}

	minutes, err := strconv.ParseFloat(minStr, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %s", minStr)
	}

	// 计算最终坐标
	coord := degrees + minutes/60

	// 处理半球方向
	switch hemisphere {
	case "S", "W":
		coord = -coord
	case "N", "E":
	default:
		return 0, fmt.Errorf("invalid hemisphere: %s", hemisphere)
	}

	return float64(coord), nil
}
