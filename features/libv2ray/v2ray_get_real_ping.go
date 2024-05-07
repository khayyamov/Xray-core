package libv2ray

import (
	"context"
	"errors"
	"fmt"
	v2net "github.com/xtls/xray-core/common/net"
	v2core "github.com/xtls/xray-core/core"
	v2serial "github.com/xtls/xray-core/infra/conf/serial"
	_ "github.com/xtls/xray-core/main/distro/all"
	"net"
	"net/http"
	"strings"
	"time"
)

func MeasureOutboundDelay(ConfigureFileContent string) (int64, error) {
	config, err := v2serial.LoadJSONConfig(strings.NewReader(ConfigureFileContent))
	if err != nil {
		return -1, err
	}

	// dont listen to anything for test purpose
	config.Inbound = nil
	// config.App: (fakedns), log, dispatcher, InboundConfig, OutboundConfig, (stats), router, dns, (policy)
	// keep only basic features
	config.App = config.App[:5]

	inst, err := v2core.New(config)
	if err != nil {
		return -1, err
	}

	inst.Start()
	delay, err := measureInstDelay(context.Background(), inst)
	inst.Close()
	return delay, err
}

func measureInstDelay(ctx context.Context, inst *v2core.Instance) (int64, error) {
	if inst == nil {
		return -1, errors.New("core instance nil")
	}

	tr := &http.Transport{
		TLSHandshakeTimeout: 5 * time.Second,
		DisableKeepAlives:   true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dest, err := v2net.ParseDestination(fmt.Sprintf("%s:%s", network, addr))
			if err != nil {
				return nil, err
			}
			return v2core.Dial(ctx, inst, dest)
		},
	}

	c := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", "https://www.google.com/generate_204", nil)
	start := time.Now()
	resp, err := c.Do(req)
	if err != nil {
		return -1, err
	}
	if resp.StatusCode != http.StatusNoContent {
		return -1, fmt.Errorf("status != 204: %s", resp.Status)
	}
	resp.Body.Close()
	return time.Since(start).Milliseconds(), nil
}
