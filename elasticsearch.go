package elasticsearch

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	plugin "github.com/mackerelio/go-mackerel-plugin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ElasticsearchPlugin struct {
	Client Client
	prefix string
}

func (p ElasticsearchPlugin) FetchMetrics() (map[string]float64, error) {
	clusterHealth := &ClusterHealth{}
	client := p.Client
	body, err := client.GetClusterHealth()
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	json.Unmarshal(body, &clusterHealth)
	return map[string]float64{
		"active_primary_shards":            clusterHealth.ActivePrimaryShards,
		"active_shards_percent_as_number":  clusterHealth.ActiveShardsPercentAsNumber,
		"active_shards":                    clusterHealth.ActiveShards,
		"delayed_unassigned_shards":        clusterHealth.DelayedUnassignedShards,
		"initializing_shards":              clusterHealth.InitializingShards,
		"number_of_data_nodes":             clusterHealth.NumberOfDataNodes,
		"number_of_in_flight_fetch":        clusterHealth.NumberOfInFlightFetch,
		"number_of_nodes":                  clusterHealth.NumberOfNodes,
		"number_of_pending_tasks":          clusterHealth.NumberOfPendingTasks,
		"relocating_shards":                clusterHealth.RelocatingShards,
		"task_max_waiting_in_queue_millis": clusterHealth.TaskMaxWaitingInQueueMillis,
		"unassigned_shards":                clusterHealth.UnassignedShards,
	}, nil
}

func (p ElasticsearchPlugin) GraphDefinition() map[string]plugin.Graphs {
	var (
		labelPrefix  = cases.Title(language.Und).String(p.MetricKeyPrefix())
		shardsPrefix = strings.Join([]string{labelPrefix, " Shards"}, "")
	)
	graphdef := map[string]plugin.Graphs{
		"shards": {
			Label: shardsPrefix,
			Unit:  "integer",
			Metrics: []plugin.Metrics{
				{
					Name:    "active_primary_shards",
					Label:   "active primary shards",
					Diff:    false,
					Stacked: false,
					Scale:   0,
				},
				{
					Name:    "active_shards",
					Label:   "active shards",
					Diff:    false,
					Stacked: true,
					Scale:   0,
				},
				{
					Name:    "delayed_unassigned_shards",
					Label:   "delayed unassigned shards",
					Diff:    false,
					Stacked: true,
					Scale:   0,
				},
				{
					Name:    "initializing_shards",
					Label:   "initializing shards",
					Diff:    false,
					Stacked: true,
					Scale:   0,
				},
				{
					Name:    "relocating_shards",
					Label:   "relocating shards",
					Diff:    false,
					Stacked: true,
					Scale:   0,
				},
				{
					Name:    "unassigned_shards",
					Label:   "unassigned shards",
					Diff:    false,
					Stacked: true,
					Scale:   0,
				},
			},
		},
	}
	return graphdef
}

func (p ElasticsearchPlugin) MetricKeyPrefix() string {
	if p.prefix == "" {
		return "elasticsearch"
	}
	return p.prefix
}

func Do() {
	var (
		prefix   = flag.String("metric-key-prefix", "", "Metric key prefix")
		endpoint = flag.String("endpoint", "", "endpoint of Elasticsearch")
		username = flag.String("username", "", "username of Elasticsearch")
		password = flag.String("password", "", "password of Elasticsearch")
	)
	r := regexp.MustCompile("https?://(.*)/?")
	flag.Parse()
	if !r.MatchString(*endpoint) {
		flag.Usage()
		fmt.Println("endpoint must be set")
		os.Exit(1)
	} else {
		parsedUrl, err := url.Parse(*endpoint)
		if err != nil {
			fmt.Println(err)
			return
		}
		*endpoint = strings.Join([]string{parsedUrl.Scheme, "://", parsedUrl.Host}, "")
	}
	plugin.NewMackerelPlugin(&ElasticsearchPlugin{
		Client: &ElasticsearchClient{
			endpoint: *endpoint,
			username: *username,
			password: *password,
		},
		prefix: *prefix,
	}).Run()
}
