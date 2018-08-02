package router

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

//PoolConfig is a map from hash-range to destination
type PoolConfig map[string]string

//Config defines config file format
type Config struct {
	Pools map[string]PoolConfig `yaml:"pools"`

	Lookup []string `yaml:"lookup"`
	Cache  []string `yaml:"cache"`
}

//Valid validate config structure
func (c *Config) Valid() error {
	var err Errors

	if len(c.Lookup) == 0 {
		err = err.Add(fmt.Errorf("no lookup table defined"))
	}

	for _, lookup := range c.Lookup {
		if _, ok := c.Pools[lookup]; !ok {
			err = err.Add(fmt.Errorf("no pool with name '%s'", lookup))
		}
	}

	for _, lookup := range c.Cache {
		if _, ok := c.Pools[lookup]; !ok {
			err = err.Add(fmt.Errorf("no pool with name '%s'", lookup))
		}
	}

	for _, pool := range c.Pools {
		for r, d := range pool {
			//validate range
			if _, rangeErr := NewRange(r); rangeErr != nil {
				err = err.Add(errors.Wrap(rangeErr, r))
			}

			if _, destErr := NewDestination(d); destErr != nil {
				err = err.Add(errors.Wrap(destErr, d))
			}
		}
	}

	if err.HasErrors() {
		return err
	}

	return nil
}

//Router returns a router that corresponds to configuration object
func (c *Config) Router() (*Router, error) {
	router := Router{
		Pools:  make(map[string]Pool),
		Lookup: c.Lookup,
		Cache:  c.Cache,
	}

	for name, cfg := range c.Pools {
		pool := new(ScanPool)

		for rangeStr, destStr := range cfg {
			hashRange, err := NewRange(rangeStr)
			if err != nil {
				return nil, errors.Wrap(err, rangeStr)
			}

			dest, err := NewDestination(destStr)
			if err != nil {
				return nil, errors.Wrap(err, destStr)
			}

			pool.Rules = append(
				pool.Rules,
				Rule{hashRange, dest},
			)
		}

		router.Pools[name] = pool
	}

	return &router, nil
}

//NewConfig loads config from reader, expecting yaml formatted config
func NewConfig(in io.Reader) (*Config, error) {
	buf, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(buf, &config); err != nil {
		return nil, err
	}

	return &config, config.Valid()
}

//NewConfigFromFile loads config from yaml file
func NewConfigFromFile(name string) (*Config, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, errors.Wrap(err, name)
	}

	defer file.Close()

	return NewConfig(file)
}
