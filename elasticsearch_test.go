package elasticsearch

import (
	"reflect"
	"testing"
)

func TestMetricKeyPrefix_Default(t *testing.T) {
	plugin := &ElasticsearchPlugin{
		Client: &ElasticsearchClientMock{
			endpoint: "https://localhost:9200",
			username: "elastic",
			password: "password",
		},
		prefix: "",
	}
	actual := plugin.MetricKeyPrefix()
	if actual != "elasticsearch" {
		t.Errorf("expected: elasticsearch, actual: %v", actual)
	}
}

func TestMetricKeyPrefix_Custom(t *testing.T) {
	plugin := &ElasticsearchPlugin{
		Client: &ElasticsearchClientMock{
			endpoint: "https://localhost:9200",
			username: "elastic",
			password: "password",
		},
		prefix: "test",
	}
	expected := plugin.prefix
	actual := plugin.MetricKeyPrefix()
	if actual != expected {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestFetchMetrics(t *testing.T) {
	expected := map[string]float64{
		"active_primary_shards":            1,
		"active_shards_percent_as_number":  100,
		"active_shards":                    1,
		"delayed_unassigned_shards":        0,
		"initializing_shards":              0,
		"number_of_data_nodes":             1,
		"number_of_in_flight_fetch":        0,
		"number_of_nodes":                  1,
		"number_of_pending_tasks":          0,
		"relocating_shards":                0,
		"task_max_waiting_in_queue_millis": 0,
		"unassigned_shards":                0,
	}
	plugin := &ElasticsearchPlugin{
		Client: &ElasticsearchClientMock{
			endpoint: "https://localhost:9200",
			username: "elastic",
			password: "password",
		},
		prefix: "",
	}
	actual, err := plugin.FetchMetrics()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v actual: %v\n", expected, actual)
	}
}
