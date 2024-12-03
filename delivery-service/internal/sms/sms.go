package sms

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// Клиент для отправки SMS
type SMSClient struct {
	Login    string
	Password string
}

func NewSMSClient() (*SMSClient, error) {
	login := os.Getenv("SMSC_LOGIN")
	password := os.Getenv("SMSC_PASSWORD")

	if login == "" || password == "" {
		return nil, errors.New("SMSC credentials are not set")
	}

	return &SMSClient{
		Login:    login,
		Password: password,
	}, nil
}

// Отправляет SMS-сообщение (метод класса)
func (c *SMSClient) SendSMS(phone, message string) error {
	if c.Login == "" || c.Password == "" {
		return errors.New("SMSC credentials are not set")
	}

	// Данные для отправки
	params := url.Values{}
	params.Set("login", c.Login)
	params.Set("psw", c.Password)
	params.Set("phones", phone)
	params.Set("mes", message)
	params.Set("fmt", "3")
	params.Set("sender", "SMSC")

	apiURL := "https://smsc.ru/sys/send.php?" + params.Encode()

	// СМС не отправляются, тк это платно, но в лк smsc можно увидеть что попытка отправки былв
	// HTTP-клиент с отключенной проверкой сертификатов
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Запрос
	resp, err := client.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMS API returned status: %d", resp.StatusCode)
	}

	return nil
}

// Отправляет SMS-сообщение (функция вне класса)
func SendSMS(phone, message string) error {
	login := os.Getenv("SMSC_LOGIN")
	password := os.Getenv("SMSC_PASSWORD")

	if login == "" || password == "" {
		return errors.New("SMSC credentials are not set")
	}

	// Данные для отправки
	params := url.Values{}
	params.Set("login", login)
	params.Set("psw", password)
	params.Set("phones", phone)
	params.Set("mes", message)
	params.Set("fmt", "3")
	params.Set("sender", "SMSC")

	apiURL := "https://smsc.ru/sys/send.php?" + params.Encode()

	// HTTP-клиент с отключенной проверкой сертификатов
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMS API returned status: %d", resp.StatusCode)
	}

	return nil
}
