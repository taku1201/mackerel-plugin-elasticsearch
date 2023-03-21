package elasticsearch

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client interface {
	GetClusterHealth() ([]byte, error)
}

type ElasticsearchClient struct {
	endpoint string
	username string
	password string
}

type ElasticsearchClientMock struct {
	endpoint string
	username string
	password string
}

func (c *ElasticsearchClient) GetClusterHealth() ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, _ := http.NewRequest("GET", strings.Join([]string{c.endpoint, "_cluster", "health"}, "/"), nil)
	req.SetBasicAuth(c.username, c.password)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return body, nil
}

func (c *ElasticsearchClientMock) GetClusterHealth() ([]byte, error) {
	if strings.Contains(c.endpoint, "fail") {
		err := errors.New("error")
		return nil, err
	}
	res := []byte(`{"active_primary_shards":1,"active_shards_percent_as_number":100,"active_shards":1,"delayed_unassigned_shards":0,"initializing_shards":0,"number_of_data_nodes":1,"number_of_in_flight_fetch":0,"number_of_nodes":1,"number_of_pending_tasks":0,"relocating_shards":0,"task_max_waiting_in_queue_millis":0,"unassigned_shards":0}`)
	return res, nil
}
