package gonepm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Registry holds the registry information.
type Registry struct {
	BaseURL string  `json:"baseurl"`
	Info    RegInfo `json:"reginfo"`
}

// RegInfo holds the registry information from GET /.
type RegInfo struct {
	DbName             string `json:"db_name"`
	DocCount           int    `json:"doc_count"`
	DocDelCount        int    `json:"doc_del_count"`
	UpdateSeq          int    `json:"update_seq"`
	PurgeSeq           int    `json:"purge_seq"`
	CompactRunning     bool   `json:"compact_running"`
	DiskSize           int    `json:"disk_size"`
	DataSize           int    `json:"data_size"`
	InstanceStartTime  string `json:"instance_start_time"`
	DiskFormatVersion  int    `json:"disk_format_version"`
	CommittedUpdateSeq int    `json:"committed_update_seq"`
}

// NewRegistry creates a new registry pointing to the URL.
func NewRegistry(baseURL string) (*Registry, error) {

	// Remove trailing slash.
	baseURL = strings.TrimRight(baseURL, "/")

	reg := &Registry{BaseURL: baseURL}
	if baseURL == "" {
		return reg, fmt.Errorf("gonepm.NewRegistry: empty baseURL")
	}

	log.Printf("Connecting to %s\n", baseURL)

	if err := reg.RegistryInfo(); err != nil {
		return reg, fmt.Errorf("gonpm:NewRegistry: %v", err.Error())
	}
	log.Printf("Successfully connected to %s\n", baseURL)
	log.Printf("Registry info: %s\n", reg.String())
	return reg, nil
}

// MakeRegistries creates a slice of registries.
func MakeRegistries(registryAddresses []string) ([]*Registry, error) {
	var regs []*Registry
	for _, addr := range registryAddresses {
		tempReg, err := NewRegistry(addr)
		if err != nil {
			return regs, fmt.Errorf("gonepm.MakeRegistries: %s", err.Error())
		}
		regs = append(regs, tempReg)
	}
	return regs, nil
}

// RegistryInfo performs a GET / and populates the registry object.
func (r *Registry) RegistryInfo() error {
	resp, err := http.Get(r.BaseURL)
	if err != nil {
		return fmt.Errorf("gonepm.RegistryInfo: %s", err.Error())
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&r.Info)
}

// String is the stringer for a registry.
func (r *Registry) String() string {

	return fmt.Sprintf("BaseURL: %s - db_name: %s - doc_count: %d",
		r.BaseURL, r.Info.DbName, r.Info.DocCount)
}

// RegInfo marshals registry.Info and returns the JSON text.
func (r *Registry) RegInfo() string {
	js, err := json.MarshalIndent(r.Info, "", "  ")
	if err != nil {
		return ""
	}
	return string(js)
}
