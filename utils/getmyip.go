package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
)

type Service struct {
	Name     string
	URL      string
	IsJSON   bool
	ParseKey string
	IPv6Only bool
}

var defaultServices = []Service{
	{"AWS CheckIP", "https://checkip.amazonaws.com", false, "", false},
	{"ifconfig.me", "https://ifconfig.me", false, "", false},
	{"ifconfig.co", "https://ifconfig.co", false, "", false},
	{"ipify", "https://api.ipify.org", false, "", false},
	{"icanhazip", "https://icanhazip.com", false, "", false},
	{"ident.me", "https://ident.me", false, "", false},
	{"ipinfo.io", "https://ipinfo.io/json", true, "ip", false},
	{"ip-api", "http://ip-api.com/json", true, "query", false},
	{"ipwho.is", "https://ipwho.is", true, "ip", false},
	{"wtfismyip", "https://wtfismyip.com/text", false, "", false},
}

func GetMyIP(ctx context.Context, ipv6 bool) string {
	log.Debug("GetMyIP")
	var (
		ip   string
		name string
		err  error
	)
	for _, svc := range defaultServices {
		ip, name, err = fetchIP(ctx, svc, ipv6, 3)
		if err != nil {
			log.Error("fetchIP", err)
			continue
		}
		log.Debugf("  From '%s' got '%s'", name, ip)
		break
	}

	return ip
}

func fetchIP(ctx context.Context, svc Service, preferV6 bool, retries int) (string, string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	transport := &http.Transport{}

	if preferV6 || svc.IPv6Only {
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return (&net.Dialer{Timeout: 5 * time.Second, DualStack: true}).DialContext(ctx, "tcp6", addr)
		}
	}
	client.Transport = transport

	var lastErr error
	for i := 0; i <= retries; i++ {
		req, err := http.NewRequestWithContext(ctx, "GET", svc.URL, nil)
		if err != nil {
			log.Error("http.NewRequestWithContext()", err)
			return "", svc.Name, err
		}

		req.Header.Add("User-Agent", "curl/8.12")

		resp, err := client.Do(req)
		if err != nil {
			log.Error("client.Do()", err)
			lastErr = err
			time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
			continue
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		if svc.IsJSON {
			var result map[string]interface{}
			if err := json.Unmarshal(body, &result); err != nil {
				log.Error("json.Unmarshal()", err)
				return "", svc.Name, fmt.Errorf("JSON parse error: %w", err)
			}
			if val, ok := result[svc.ParseKey]; ok {
				if ipStr, ok := val.(string); ok {
					return ipStr, svc.Name, nil
				}
				err := errors.New("invalid type for IP")
				log.Error("IP type", err)
				return "", svc.Name, err
			}
			err := fmt.Errorf("missing '%s' key in JSON", svc.ParseKey)
			log.Error("JSON Key", err)
			return "", svc.Name, err
		}

		return strings.TrimSpace(string(body)), svc.Name, nil
	}

	return "", svc.Name, lastErr
}
