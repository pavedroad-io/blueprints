package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type Customers struct {
	Customers []Customer `yaml:"customers"`
}

type Customer struct {
	// ID unique ID for this customer
	ID uuid.UUID `yaml:"id"`

	// Customer Name
	Name string `yaml:"name"`

	// Logs to monitor
	Logs []LogBuckets `yaml:"logs"`

	// Providers associated with logs
	Providers Providers `yaml:"providers"`
}

func (c *Customer) LoadFromDisk(file string) ([]Customer, error) {
	cl := []Customer{}

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("failed to open:", file, ", error:", err)
	}
	defer f.Close()

	byteValue, e := ioutil.ReadAll(f)
	if e != nil {
		fmt.Println("read failed for ", file)
		return nil, err
	}

	err = yaml.Unmarshal([]byte(byteValue), &cl)
	if err != nil {
		fmt.Println("Unmarshal faild", err)
		return nil, err
	}

	return cl, nil
}
