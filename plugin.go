package traefik_plugin_rename_header

import (
	"context"
	"errors"
	"net/http"
)

type Config struct {
	OldHeader string `json:"oldHeader,omitempty"`
	NewHeader string `json:"newHeader,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type RenameHeader struct {
	next      http.Handler
	oldHeader string
	newHeader string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.OldHeader == "" || config.NewHeader == "" {
		return nil, errors.New("both oldHeader and newHeader must be specified")
	}

	return &RenameHeader{
		next:      next,
		oldHeader: config.OldHeader,
		newHeader: config.NewHeader,
	}, nil
}

func (r *RenameHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if val := req.Header.Get(r.oldHeader); val != "" {
		req.Header.Set(r.newHeader, val)
		req.Header.Del(r.oldHeader)
	}

	r.next.ServeHTTP(rw, req)
}
