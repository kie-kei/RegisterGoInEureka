package registerGoInEureka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Конфигурация для регистрации в Eureka
type EurekaConfig struct {
	EurekaURL     string
	InstanceId    string
	AppName       string
	HostName      string
	IPAddr        string
	Port          int
	SecurePort    int
	StatusPageUrl string
}

// Структура для регистрации экземпляра в Eureka
type EurekaInstance struct {
	Instance struct {
		InstanceId string `json:"instanceId"`
		HostName   string `json:"hostName"`
		App        string `json:"app"`
		IpAddr     string `json:"ipAddr"`
		Port       struct {
			Value int `json:"$"`
		} `json:"port"`
		SecurePort struct {
			Value int `json:"$"`
		} `json:"securePort"`
		StatusPageUrl string `json:"statusPageUrl"`
	} `json:"instance"`
}

// Функция для регистрации экземпляра в Eureka
func RegisterInstance(config EurekaConfig) error {
	// Формируем структуру запроса
	instance := EurekaInstance{}
	instance.Instance.InstanceId = config.InstanceId
	instance.Instance.App = config.AppName
	instance.Instance.HostName = config.HostName
	instance.Instance.IpAddr = config.IPAddr
	instance.Instance.StatusPageUrl = config.StatusPageUrl
	instance.Instance.Port.Value = config.Port
	instance.Instance.SecurePort.Value = config.SecurePort

	// Преобразуем структуру в JSON
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
