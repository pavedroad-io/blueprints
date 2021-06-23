package main

// Provider for creating new client
type Provider struct {
	// Name AWS, GCP, etc
	Name string `yaml:"name"`

	// Credentials key name
	Credentials string `yaml:"credentials"`
	// Key value
	Key string `yaml:"key"`

	// Region us-west-1, etc
	Region string `yaml:"region"`

	// Endpoint s3.aws.com
	Endpoint string `yaml:"endpoint"`
}

type Providers []Provider

func (ps *Providers) Lookup(pName string) (Provider, error) {
	rp := Provider{}
	for _, p := range *ps {
		if p.Name == pName {
			return p, nil
		}
	}
	return rp, nil
}
