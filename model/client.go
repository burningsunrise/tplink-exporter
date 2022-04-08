package model

import (
	"crypto/tls"
	"net/http"
	"time"
)

func HttpClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	return client
}
