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

// common contains shared functions for general purposes, like Authentication, choosing version etc.

package openstack

import (
	"fmt"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/apiversions"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tenants"

	apiversionsintel "github.com/intelsdi-x/snap-plugin-collector-cinder/openstack/apiversions"
	openstackintel "github.com/intelsdi-x/snap-plugin-collector-cinder/openstack/v2"
	"github.com/intelsdi-x/snap-plugin-collector-cinder/types"
)

var apiPriority = map[string]int{
	"v1.0": 1,
	"v2.0": 2,
}

// Commoner provides abstraction for shared functions mainly for mocking
type Commoner interface {
	GetTenants(endpoint, user, password string) ([]types.Tenant, error)
	GetApiVersions(provider *gophercloud.ProviderClient) ([]string, error)
}

// Common is a receiver for Commoner interface
type Common struct{}

// GetTenants is used to retrieve list of available tenant for authenticated user
// List of tenants can then be used to authenticate user for each given tenant
func (c Common) GetTenants(endpoint, user, password string) ([]types.Tenant, error) {
	tnts := []types.Tenant{}

	provider, err := Authenticate(endpoint, user, password, "")
	if err != nil {
		return nil, err
	}

	client := openstack.NewIdentityV2(provider)

	opts := tenants.ListOpts{}
	pager := tenants.List(client, &opts)

	page, err := pager.AllPages()
	if err != nil {
		return tnts, err
	}

	tenantList, err := tenants.ExtractTenants(page)
	if err != nil {
		return tnts, err
	}

	for _, t := range tenantList {
		tnts = append(tnts, types.Tenant{Name: t.Name, ID: t.ID})
	}

	return tnts, nil
}

// GetApiVersions is used to retrieve list of available Cinder API versions
// List of api version is then used to dispatch calls to proper API version based on defined priority
func (c Common) GetApiVersions(provider *gophercloud.ProviderClient) ([]string, error) {
	apis := []string{}

	client, err := openstackintel.NewBlockStorageV2(provider, gophercloud.EndpointOpts{})

	if err != nil {
		return apis, err
	}

	pager := apiversionsintel.List(client)
	page, err := pager.AllPages()
	if err != nil {
		fmt.Println("err ", err)
		return apis, err
	}

	apiVersions, err := apiversions.ExtractAPIVersions(page)
	if err != nil {
		return apis, err
	}

	for _, apiVersion := range apiVersions {
		apis = append(apis, apiVersion.ID)
	}

	return apis, nil
}

// Authenticate is used to authenticate user for given tenant. Request is send to provided Keystone endpoint
// Returns authenticated provider client, which is used as a base for service clients.
func Authenticate(endpoint, user, password, tenant string) (*gophercloud.ProviderClient, error) {
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: endpoint,
		Username:         user,
		Password:         password,
		TenantName:       tenant,
		AllowReauth:      true,
	}

	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// ChooseVersion returns chosen Cinder API version based on defined priority
func ChooseVersion(recognized []string) (string, error) {
	if len(recognized) < 1 {
		return "", fmt.Errorf("No recognized API versions provided")
	}
	chosen := recognized[0]
	for _, ver := range recognized[1:] {
		chosenPriority, ok1 := apiPriority[chosen]
		verPriority, ok2 := apiPriority[ver]
		if ok1 && ok2 {
			if chosenPriority < verPriority {
				chosen = ver
			}
		}
	}
	return chosen, nil
}
