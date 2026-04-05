package address

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// IpDetector maps to IpDetector in Java.
type IpDetector struct {
	mu                         sync.RWMutex
	cachedPrivateIp            string
	privateIpLastUpdated       time.Time
	cachedPublicIp             string
	publicIpLastUpdated        time.Time
	
	client                     *http.Client
}

func NewIpDetector() *IpDetector {
	return &IpDetector{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// QueryPrivateIp maps to queryPrivateIp() in Java.
func (d *IpDetector) QueryPrivateIp(expireAfterMillis int) (string, error) {
	d.mu.RLock()
	if expireAfterMillis > 0 && d.cachedPrivateIp != "" && time.Since(d.privateIpLastUpdated) < time.Duration(expireAfterMillis)*time.Millisecond {
		ip := d.cachedPrivateIp
		d.mu.RUnlock()
		return ip, nil
	}
	d.mu.RUnlock()

	conn, err := net.Dial("udp", "8.8.8.8:10002")
	if err != nil {
		return "", errors.New("failed to detect the local IP: " + err.Error())
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP
	if !ip.IsPrivate() {
		return "", errors.New("the IP address (" + ip.String() + ") is not a site local IP address")
	}
	
	ipStr := ip.String()

	d.mu.Lock()
	d.cachedPrivateIp = ipStr
	d.privateIpLastUpdated = time.Now()
	d.mu.Unlock()

	return ipStr, nil
}

// QueryPublicIp maps to queryPublicIp() in Java.
func (d *IpDetector) QueryPublicIp(ctx context.Context, detectorAddresses []string, expireAfterMillis int) (string, error) {
	d.mu.RLock()
	if expireAfterMillis > 0 && d.cachedPublicIp != "" && time.Since(d.publicIpLastUpdated) < time.Duration(expireAfterMillis)*time.Millisecond {
		ip := d.cachedPublicIp
		d.mu.RUnlock()
		return ip, nil
	}
	d.mu.RUnlock()

	if len(detectorAddresses) == 0 {
		return "", errors.New("failed to detect the public IP of the local node because no IP detector address is specified")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultCh := make(chan string, len(detectorAddresses))
	
	var wg sync.WaitGroup
	for _, addr := range detectorAddresses {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return
			}
			resp, err := d.client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				body, err := io.ReadAll(resp.Body)
				if err == nil {
					ip := strings.TrimSpace(string(body))
					if net.ParseIP(ip) != nil {
						select {
						case resultCh <- ip:
						case <-ctx.Done():
						}
					}
				}
			}
		}(addr)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	select {
	case ip, ok := <-resultCh:
		if !ok {
			return "", errors.New("failed to detect the public IP of the local node because there is no available IP")
		}
		// First successful IP
		d.mu.Lock()
		d.cachedPublicIp = ip
		d.publicIpLastUpdated = time.Now()
		d.mu.Unlock()
		return ip, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
