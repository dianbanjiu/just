package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type JustMySocks struct {
	conf *JustConf
}
type JustConf struct {
	UrlWidthCounter string `json:"url_width_counter"`
}

func NewJust(conf JustConf) JustMySocks {
	return JustMySocks{conf: &conf}
}

func (s *JustMySocks) client() *http.Client {
	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   3 * time.Second,
	}
	return client
}

type BandwidthCounter struct {
	MonthlyBWLimitB   int64 `json:"monthly_bw_limit_b"`    // 每月总流量，单位B
	BWCounterB        int64 `json:"bw_counter_b"`          // 已经使用的流量，单位B
	BWResetDayOfMonth int   `json:"bw_reset_day_of_month"` // 每月流量重置日，洛杉矶时区
}

func (s *JustMySocks) GetWidthCounter(ctx context.Context) (*BandwidthCounter, error) {
	request, err := http.NewRequest(http.MethodGet, s.conf.UrlWidthCounter, nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	response, err := s.client().Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response code: %d, body: %s", response.StatusCode, string(body))
	}

	var data = new(BandwidthCounter)
	err = json.Unmarshal(body, &data)
	return data, err
}
