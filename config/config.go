// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Repositories []Repository `config:"repositories"`
	Token        string `config:"token"`
}

type Repository struct {
	Owner        string `config:"owner"`
	Name         string `config:"name"`
	DocumentType string `config:"document_type"`
	TimeInterval    time.Duration `config:"time_interval"`
}

var DefaultConfig = Config{
	Token:        "",
	Repositories: []Repository{},
}
