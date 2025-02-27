package utils

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

// GPSData 结构体保存解析后的 GPS 信息
type GPSData struct {
	LocalTime    time.Time // UTC 时间
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
func readGPSData(portName string) (GPSData, error) {
	port, err := serial.Open(portName, &serial.Mode{BaudRate: 9600})
	if err != nil {
		return GPSData{}, fmt.Errorf("无法打开串口 %s: %v", portName, err)
	}
	defer port.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dataCh := make(chan GPSData)
	errCh := make(chan error)

	go func() {
		reader := bufio.NewReader(port)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				errCh <- fmt.Errorf("读取串口数据时出错: %v", err)
				return
			}

			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "$GPRMC") {
				data, err := parseGPRMC(line)
				if err != nil {
					errCh <- fmt.Errorf("解析 GPRMC 消息失败: %v", err)
					return
				}
				data.LocalTimeStr = data.LocalTime.Format("2006-01-02 15:04:05")
				dataCh <- data
				return
			}
		}
	}()

	select {
	case data := <-dataCh:
		logrus.Debugf("UTC时间: %s", data.LocalTimeStr)
		logrus.Debugf("纬度: %.6f°", data.Latitude)
		logrus.Debugf("经度: %.6f°", data.Longitude)
		return data, nil
	case err := <-errCh:
		return GPSData{}, err
	case <-ctx.Done():
		return GPSData{}, fmt.Errorf("读取 GPS 数据超时")
	}
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
		return GPSData{}, fmt.Errorf("GPRMC 消息字段不足: %s", sentence)
	}

	// 提取 UTC 时间（hhmmss.sss）
	rawTime := fields[1] // 格式: HHMMSS.SS
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
	// 获取本地时区（可以替换为其他时区，如："Asia/Shanghai"）
	localLocation, err := time.LoadLocation("Local")
	if err != nil {
		log.Fatalf("无法加载本地时区: %v", err)
	}

	// 将 UTC 时间转换为本地时间
	localTime := utcTime.In(localLocation)

	// 提取纬度和经度
	latitude, err := convertToDecimalDegrees(fields[3], fields[4])
	if err != nil {
		return GPSData{}, fmt.Errorf("纬度解析失败: %v", err)
	}

	longitude, err := convertToDecimalDegrees(fields[5], fields[6])
	if err != nil {
		return GPSData{}, fmt.Errorf("经度解析失败: %v", err)
	}

	// 返回封装好的 GPS 数据
	return GPSData{
		LocalTime: localTime,
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
	return nil
}
func GetNtpInfo() (GPSData, error) {
	// 配置串口名称（注意替换为实际串口）
	gpsPort := "/dev/ttyUSB1" // 用于接收 GPS 数据的串口
	// 开始读取 GPS 数据
	data, err := readGPSData(gpsPort)
	if err != nil {
		logrus.Errorf("读取 GPS 数据失败: %v", err)
		return GPSData{}, nil
	}
	return data, nil
}

// 将 NMEA 格式的经纬度转换为十进制格式
func convertToDecimalDegrees(value, direction string) (float64, error) {
	if value == "" || direction == "" {
		return 0, fmt.Errorf("无效的值或方向: value=%s, direction=%s", value, direction)
	}

	// 将纬度或经度的 NMEA 格式 (ddmm.mmmm) 转换为十进制格式
	degrees, err := strconv.ParseFloat(value[:2], 64) // 取前两位作为度数
	if err != nil {
		return 0, fmt.Errorf("无法解析度: %v", err)
	}

	minutes, err := strconv.ParseFloat(value[2:], 64) // 后面的部分是分
	if err != nil {
		return 0, fmt.Errorf("无法解析分: %v", err)
	}

	// 转换为十进制
	decimal := degrees + (minutes / 60.0)

	// 根据方向调整符号
	if direction == "S" || direction == "W" {
		decimal = -decimal
	}

	return decimal, nil
}
