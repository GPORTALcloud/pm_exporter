package api

import (
	"encoding/base64"
	"fmt"
)

type HpeApi struct {
	auth string
	host string
}

func NewHpeApi(host string) *HpeApi {
	r := HpeApi{host: host}
	return &r
}

func (r *HpeApi) SetUser(user string, pass string) *HpeApi {
	r.auth = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pass)))
	return r
}

func (r *HpeApi) Request(path string) (*HpeSystemMetric, error) {
	// TODO implement
	return nil, nil
}
func (r *HpeApi) GetHost() string {
	return r.host
}

func (r *HpeApi) FetchInventoryMetrics() error {
	return nil
}

type HpeSystemMetric struct {
}
