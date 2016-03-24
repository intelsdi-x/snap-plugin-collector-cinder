/*
http://www.apache.org/licenses/LICENSE-2.0.txt
Copyright 2016 Intel Corporation
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

// service contains interface and dispatcher methods for Cinder API versions

package services

import (
	"github.com/rackspace/gophercloud"

	openstackintel "github.com/intelsdi-x/snap-plugin-collector-cinder/openstack"
	cinderv1 "github.com/intelsdi-x/snap-plugin-collector-cinder/openstack/v1/cinder"
	cinderv2 "github.com/intelsdi-x/snap-plugin-collector-cinder/openstack/v2/cinder"
	"github.com/intelsdi-x/snap-plugin-collector-cinder/types"
)

// Cinderer allows usage of different Cinder API versions for metric collection
type Cinderer interface {
	GetLimits(provider *gophercloud.ProviderClient) (types.Limits, error)
	GetVolumes(provider *gophercloud.ProviderClient) (map[string]types.Volumes, error)
	GetSnapshots(provider *gophercloud.ProviderClient) (map[string]types.Snapshots, error)
}

// Services serves as a API calls dispatcher
type Service struct {
	cinder Cinderer
}

// Set allows to set proper API version implementation
func (c *Service) Set(new Cinderer) {
	c.cinder = new
}

// GetLimits dispatches call to proper API version calls to collect limits metrics
func (s Service) GetLimits(provider *gophercloud.ProviderClient) (types.Limits, error) {
	return s.cinder.GetLimits(provider)
}

// GetVolumes dispatches call to proper API version calls to collect volumes metrics
func (s Service) GetVolumes(provider *gophercloud.ProviderClient) (map[string]types.Volumes, error) {
	return s.cinder.GetVolumes(provider)
}

// GetSnapshots dispatches call to proper API version calls to collect snapshot metrics
func (s Service) GetSnapshots(provider *gophercloud.ProviderClient) (map[string]types.Snapshots, error) {
	return s.cinder.GetSnapshots(provider)
}

// Dispatch redirects to selected Cinder API version based on priority
func Dispatch(provider *gophercloud.ProviderClient) Service {
	cmn := openstackintel.Common{}
	versions, err := cmn.GetApiVersions(provider)
	if err != nil {
		panic(err)
	}

	chosen, err := openstackintel.ChooseVersion(versions)
	if err != nil {
		panic(err)
	}

	service := Service{}
	switch chosen {
	case "v1.0":
		service.Set(cinderv1.ServiceV1{})
	case "v2.0":
		service.Set(cinderv2.ServiceV2{})
	default:
		panic("Could not select dispatcher")
	}

	return service
}
