package registerGoInEureka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type EurekaConfig struct {
	EurekaURL        string
	InstanceId       string
	AppName          string
	HostName         string
	IPAddr           string
	Port             int
	SecurePort       int // Теперь это необязательный параметр
	StatusPageUrl    string
	HomePageUrl      string
	HealthCheckUrl   string
	VipAddress       string
	SecureVipAddress string
	DataCenterName   string // Новый параметр
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
		} `json:"securePort,omitempty"` // Убираем обязательность для SecurePort
		DataCenterInfo struct {
			Class string `json:"@class"`
			Name  string `json:"name"`
		} `json:"dataCenterInfo"`
	} `json:"instance"`
}

func RegisterInstance(config EurekaConfig) error {
	if config.HomePageUrl == "" {
		config.HomePageUrl = config.StatusPageUrl
	}
	if config.HealthCheckUrl == "" {
		config.HealthCheckUrl = config.StatusPageUrl
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

	// Формируем структуру запроса
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

	// Отправляем запрос на регистрацию в Eureka
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
