/*
Copyright 2014 Rohith All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package marathon

// IpAddressInfo defines the IP address per container info
type IPAddressInfo struct {
	Groups        []string          `json:"groups,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	DiscoveryInfo *DiscoveryInfo    `json:"discovery,omitempty"`
	NetworkName   string            `json:"networkName,omitempty"`
}

type DiscoveryInfo struct {
	Ports *[]Port
}

type Port struct {
	Number   int               `json:"number,omitempty"`
	Name     string            `json:"name,omitempty"`
	Protocol string            `json:"protocol,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
}
