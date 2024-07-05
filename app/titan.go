package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
	"titan/app/config"
	"titan/internel/client"
)

type Titan struct {
	dcMap        map[string]string
	conf         config.Config
	lastCallTime map[string]time.Time
}

func NewTitan(conf config.Config) (*Titan, error) {
	baseDir := conf.BaseDir // Assuming conf has a BaseDir field
	dcTokenContent, err := os.ReadFile(baseDir + "/dc.txt")
	if err != nil {
		slog.Error("failed to read dc.txt", err)
		return nil, err
	}
	titanAddressContent, err := os.ReadFile(baseDir + "/address.txt")
	if err != nil {
		slog.Error("failed to read address.txt", err)
		return nil, err
	}
	dcToken := strings.Split(string(dcTokenContent), "\n")
	titanAddress := strings.Split(string(titanAddressContent), "\n")
	if len(dcToken) != len(titanAddress) {
		return nil, fmt.Errorf("dc.txt and address.txt do not have the same number of lines")

	}
	var dcMap = make(map[string]string)
	for i, dc := range dcToken {
		dcMap[dc] = titanAddress[i]
	}
	return &Titan{
		dcMap:        dcMap,
		conf:         conf,
		lastCallTime: make(map[string]time.Time),
	}, nil
}

func (t *Titan) Start() {
	interval := 2 * time.Hour
	ticker := time.NewTicker(interval + 1*time.Minute)
	for {
		for dc, address := range t.dcMap {
			slog.Info("Starting Titan", "dc", dc, "address", address)
			elapsed := time.Since(t.lastCallTime[dc])
			if elapsed < interval {
				time.Sleep(interval - elapsed)
			}
			t.callRequestDcMessage(dc, address)
		}
		<-ticker.C
	}
}

func (t *Titan) callRequestDcMessage(dc string, address string) {
	err := client.RequestDcMessage(dc, address)
	if err != nil {
		slog.Error("failed to send message dc:", "dc", dc, "address:", address)
	} else {
		t.lastCallTime[dc] = time.Now()
	}
	// Random delay between 4 to 7 seconds
	delay := time.Duration(4) * time.Second
	time.Sleep(delay)
}
