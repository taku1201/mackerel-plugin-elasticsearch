package elasticsearch

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	plugin "github.com/mackerelio/go-mackerel-plugin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ElasticsearchPlugin struct {
	prefix   string
	endpoint string
	username string
	password string
}

func (p ElasticsearchPlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, _ := http.NewRequest("GET", strings.Join([]string{p.endpoint, "_cluster", "health"}, "/"), nil)
	req.SetBasicAuth(p.username, p.password)
	// send request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	// read response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// unmarshall json
	json.Unmarshal(body, &stat)
	return stat, nil
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
					Name:    "initializing_shards",
					Label:   "initializing shards",
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
				{
					Name:    "delayed_unassigned_shards",
					Label:   "delayed unassigned shards",
					Diff:    false,
					Stacked: true,
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
					Name:    "active_primary_shards",
					Label:   "active primary shards",
					Diff:    false,
					Stacked: false,
					Scale:   0,
				},
				{
					Name:    "relocating_shards",
					Label:   "relocating shards",
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
		prefix:   *prefix,
		endpoint: *endpoint,
		username: *username,
		password: *password,
	}).Run()
}
