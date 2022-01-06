package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

type Client struct {
	client  influxdb2.Client
	orgID   string
	orgName string
}

func New() *Client {
	opts := influxdb2.DefaultOptions()
	opts.SetHTTPRequestTimeout(600)
	client := influxdb2.NewClientWithOptions(
		"http://hackhost:8086",
		"Qpw669YsF-LYXz2UHa8PIqxVRasfLJEBSV684ZSBLdytjnG9w0HtGBO66yAudXwhGlZK0PjKyqtgEsCgvJrcnQ==",
		opts,
	)
	orgID := "3744ad65272d71de"
	orgName := "Manual Pilot"

	return &Client{client: client, orgID: orgID, orgName: orgName}
}

func (c *Client) Close() {
	c.client.Close()
}

func (c *Client) ListFragments(ctx context.Context) ([]string, error) {
	bucketAPI := c.client.BucketsAPI()
	buckets, err := bucketAPI.GetBuckets(ctx)
	if err != nil {
		return nil, err
	}

	frags := []string{}
	for _, bucket := range *buckets {
		if *bucket.Type != domain.BucketTypeSystem {
			frags = append(frags, bucket.Name)
		}
	}
	return frags, nil
}

func (c *Client) ImportFragment(ctx context.Context, frag string) (api.WriteAPI, error) {
	frags, err := c.ListFragments(ctx)
	if err != nil {
		return nil, err
	}
	for _, f := range frags {
		if f == frag {
			return c.client.WriteAPI(c.orgName, frag), nil
		}
	}

	bucketAPI := c.client.BucketsAPI()
	if _, err := bucketAPI.CreateBucketWithNameWithID(ctx, c.orgID, frag); err != nil {
		return nil, err
	}
	return c.client.WriteAPI(c.orgName, frag), nil
}

func (c *Client) DeleteFragment(ctx context.Context, frag string) error {
	bucketAPI := c.client.BucketsAPI()

	buckets, err := bucketAPI.GetBuckets(ctx)
	if err != nil {
		return err
	}

	for _, bucket := range *buckets {
		if frag == bucket.Name {
			return bucketAPI.DeleteBucket(ctx, &bucket)
		}
	}
	return nil
}

func (c *Client) GetFieldStats(ctx context.Context, frag string, start, stop time.Time, filters []string, event, field string) (map[string]int64, error) {
	queryAPI := c.client.QueryAPI(c.orgID)

	if !strings.HasPrefix(field, "f_") {
		field = "f_" + field
	}

	tr, err := queryAPI.Query(ctx, fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s) %s %s
			|> filter(fn: (r) => r["_measurement"] =~ /[0-9]{5}/) 
			|> group(columns: ["%s"])
			|> sum()
	`, frag, start.Format(time.RFC3339), stop.Format(time.RFC3339), buildEventFilter([]string{event}), buildFieldFilter(filters), field))
	if err != nil {
		return nil, err
	}
	defer tr.Close()

	fmap := map[string]int64{}
	for tr.Next() {
		rec := tr.Record()
		name := ""
		if n, ok := rec.Values()[field].(string); ok {
			name = n
		}
		fmap[name] = fmap[name] + rec.Value().(int64)
	}
	return fmap, tr.Err()
}

func (c *Client) GetStats(ctx context.Context, frag string, start, stop time.Time) (map[string]int64, map[string]string, error) {
	queryAPI := c.client.QueryAPI(c.orgID)

	tr, err := queryAPI.Query(ctx, fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] =~ /[0-9]{5}/) 
			|> group(columns: ["_measurement", "name"])
			|> sum()
	`, frag, start.Format(time.RFC3339), stop.Format(time.RFC3339)))
	if err != nil {
		return nil, nil, err
	}
	defer tr.Close()

	countMap := map[string]int64{}
	messageMap := map[string]string{}
	for tr.Next() {
		rec := tr.Record()
		if name, ok := rec.Values()["name"].(string); ok {
			messageMap[rec.Measurement()] = name
		} else {
			continue
		}
		if cnt, ok := rec.Value().(int64); ok {
			countMap[rec.Measurement()] = countMap[rec.Measurement()] + cnt
		} else if cnt, ok := rec.Value().(float64); ok {
			countMap[rec.Measurement()] = int64(float64(countMap[rec.Measurement()]) + cnt)
		}
	}
	return countMap, messageMap, tr.Err()
}

func (c *Client) Search(ctx context.Context, event string, start, stop time.Time) ([]string, error) {
	queryAPI := c.client.QueryAPI(c.orgID)

	frags, err := c.ListFragments(ctx)
	if err != nil {
		return nil, err
	}

	xs := []string{}
	for _, frag := range frags {
		tr, err := queryAPI.Query(ctx, fmt.Sprintf(`
			from(bucket: "%s")
				|> range(start: %s, stop: %s)
				|> filter(fn: (r) => r["_measurement"] =~ /[0-9]{5}/) 
				|> filter(fn: (r) => r._measurement == "%s")
				|> group(columns: ["_measurement"])
				|> limit(n: 1)
		`, frag, start.Format(time.RFC3339), stop.Format(time.RFC3339), event))
		if err != nil {
			return nil, err
		}
		defer tr.Close()

		if tr.Next() {
			xs = append(xs, frag)
		}
	}

	return xs, nil
}
