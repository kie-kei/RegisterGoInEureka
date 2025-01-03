package registerGoInEureka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("ошибка при получении интерфейсов: %v", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}

	return "", fmt.Errorf("не удалось найти локальный IP-адрес")
}

type EurekaConfig struct {
	EurekaURL        string
	InstanceId       string
	AppName          string
	HostName         string
	IPAddr           string
	Port             int
	SecurePort       int
	StatusPageUrl    string
	HomePageUrl      string
	HealthCheckUrl   string
	VipAddress       string
	SecureVipAddress string
	DataCenterName   string
}

type EurekaInstance struct {
	Instance struct {
		InstanceId       string `json:"instanceId"`
		HostName         string `json:"hostName"`
		App              string `json:"app"`
		IpAddr           string `json:"ipAddr"`
		VipAddress       string `json:"vipAddress"`
		SecureVipAddress string `json:"secureVipAddress"`
		StatusPageUrl    string `json:"statusPageUrl"`
		HomePageUrl      string `json:"homePageUrl"`
		HealthCheckUrl   string `json:"healthCheckUrl"`
		Port             struct {
			Value int `json:"$"`
		} `json:"port"`
		SecurePort struct {
			Value int `json:"$"`
		} `json:"securePort,omitempty"`
		DataCenterInfo struct {
			Class string `json:"@class"`
			Name  string `json:"name"`
		} `json:"dataCenterInfo"`
	} `json:"instance"`
}

func RegisterInstance(config EurekaConfig) error {
	if config.IPAddr == "" {
		ip, err := getLocalIP()
		if err != nil {
			return fmt.Errorf("не удалось получить локальный IP-адрес: %v", err)
		}
		config.IPAddr = ip
	}
	baseUrl := "http://" + config.IPAddr + ":" + strconv.Itoa(config.Port)
	config.StatusPageUrl = baseUrl + config.StatusPageUrl
	if config.HomePageUrl == "" {
		config.HomePageUrl = config.StatusPageUrl
	} else {
		config.HomePageUrl = baseUrl + config.HomePageUrl
	}
	if config.HealthCheckUrl == "" {
		config.HealthCheckUrl = config.StatusPageUrl
	} else {
		config.HealthCheckUrl = baseUrl + config.HealthCheckUrl
	}
	if config.VipAddress == "" {
		config.VipAddress = config.AppName
	}
	if config.SecureVipAddress == "" {
		config.SecureVipAddress = config.AppName
	}
	if config.DataCenterName == "" {
		config.DataCenterName = "MyOwn"
	}

	instance := EurekaInstance{}
	instance.Instance.InstanceId = config.InstanceId
	instance.Instance.App = config.AppName
	instance.Instance.HostName = config.HostName
	instance.Instance.IpAddr = config.IPAddr
	instance.Instance.StatusPageUrl = config.StatusPageUrl
	instance.Instance.HomePageUrl = config.HomePageUrl
	instance.Instance.HealthCheckUrl = config.HealthCheckUrl
	instance.Instance.VipAddress = config.VipAddress
	instance.Instance.SecureVipAddress = config.SecureVipAddress
	instance.Instance.Port.Value = config.Port
	if config.SecurePort != 0 {
		instance.Instance.SecurePort.Value = config.SecurePort
	}
	instance.Instance.DataCenterInfo.Class = "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo"
	instance.Instance.DataCenterInfo.Name = config.DataCenterName

	data, err := json.Marshal(instance)
	if err != nil {
		return fmt.Errorf("Ошибка при маршалинге данных: %v", err)
	}

	resp, err := http.Post(config.EurekaURL+"/apps/"+config.AppName, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Ошибка при отправке запроса на регистрацию: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Ошибка при регистрации, статус: %d", resp.StatusCode)
	}

	return nil
}
