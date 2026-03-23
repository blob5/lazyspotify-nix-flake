package librespot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/dubeyKartikay/lazyspotify/librespot/models"
)

const (
  healthPath = "/"
	playPath = "/player/play"
	playpausePath = "/playpause"
)

type LibrespotApiServer struct {
	host string
	port string
}

type LibrespotApiClient struct {
	server *LibrespotApiServer
  client *http.Client
}

func NewLibrespotApiServer(host string, port string) *LibrespotApiServer {
  return &LibrespotApiServer{
    host: host,
    port: port,
  }
}

func (l *LibrespotApiServer) GetServerUrl() string {
	return fmt.Sprintf("http://%s:%s", l.host, l.port)
}

func NewLibrespotApiClient(server *LibrespotApiServer) *LibrespotApiClient {
	client := http.Client{
		Timeout: 30*time.Second,
	}
	return &LibrespotApiClient{
		client: &client,
    server: server,
	}
}

func (l *LibrespotApiClient) GetHealth() (*models.HealthResponse,error) {
	url := l.server.GetServerUrl() + healthPath;
	req, err := http.NewRequest("GET", url, nil)
	fmt.Println("Requesting", url)
	if err != nil {
		return nil,err
	}
	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resData, err := io.ReadAll(resp.Body)
	if err != nil {
    return nil,err
  }
	healthRes, err := models.DecodeHealthResponse(resData)
	if err != nil {
    return nil, err
  }
	return &healthRes, nil
}

func (l *LibrespotApiClient) Play(ctx context.Context ,uri string, skip_to_uri string, paused bool) int{
  url := l.server.GetServerUrl() + playPath;
	playRequestJson,err := models.NewPlayRequest(uri, skip_to_uri, paused)
	if (err != nil) {
    return 500
	}
  req, err := http.NewRequestWithContext(ctx,"POST", url, bytes.NewReader(playRequestJson))
	req.Header.Set("Content-Type", "application/json")
  fmt.Printf("Requesting %+v\n", req)
  if err != nil {
    return 500
  }
  resp, err := DoWithRetry(l.client, req, 3, 100*time.Millisecond)
  if err != nil {
    fmt.Println(err)
    return 500
  }
  defer resp.Body.Close()
  return resp.StatusCode
}

func DoWithRetry(client *http.Client, req *http.Request, maxRetries int, retryDelay time.Duration) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= maxRetries; i++ {

		if req.GetBody != nil {
			req.Body, _ = req.GetBody()
		}

		resp, err = client.Do(req)

		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		} else{
			fmt.Println(err)
		}

		if resp != nil {
			fmt.Printf("%+v\b",resp)
			resp.Body.Close()
		}

		if i >= maxRetries {
			break
		}

		backoffDuration := time.Duration(math.Pow(2, float64(i))) * retryDelay
		fmt.Printf("Request failed. Retrying in %v...\n", backoffDuration)
		time.Sleep(backoffDuration)
	}

	return resp, fmt.Errorf("request failed after %d retries. Last error: %v", maxRetries, err)
}
