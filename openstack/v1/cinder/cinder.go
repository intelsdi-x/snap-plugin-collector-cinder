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

// Cinder package contains wrapper functions designed to collect required metrics
// All functions are dependant on OpenStack BlockStorage API Version 1
package cinder

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"

	limitsintel "github.com/intelsdi-x/snap-plugin-collector-cinder/openstack/limits"
	"github.com/intelsdi-x/snap-plugin-collector-cinder/types"
)

// ServiceV1 serves as dispatcher for Cinder API version 1.0
type ServiceV1 struct{}

// GetLimits collects tenant limits by sending REST call to cinderhost:8776/v1/tenant_id/limits
func (s ServiceV1) GetLimits(provider *gophercloud.ProviderClient) (types.Limits, error) {
	limits := types.Limits{}

	client, err := openstack.NewBlockStorageV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return limits, err
	}

	tenantLimits, err := limitsintel.Get(client, "limits").Extract()
	if err != nil {
		return limits, err
	}

	limits.MaxTotalVolumes = tenantLimits.MaxTotalVolumes
	limits.MaxTotalVolumeGigabytes = tenantLimits.MaxTotalVolumeGigabytes

	return limits, nil
}

// GetVolumes collects volumes data by sending REST call to cinderhost:8776/v1/tenant_id/volumes
func (s ServiceV1) GetVolumes(provider *gophercloud.ProviderClient) (types.Volumes, error) {
	vols := types.Volumes{}

	client, err := openstack.NewBlockStorageV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return vols, err
	}

	//opts := volumes.ListOpts{AllTenants: true}
	opts := volumes.ListOpts{}

	pager := volumes.List(client, opts)
	page, err := pager.AllPages()
	if err != nil {
		return vols, err
	}

	volumeList, err := volumes.ExtractVolumes(page)
	if err != nil {
		return vols, err
	}

	for _, volume := range volumeList {
		vols.Count += 1
		vols.Bytes += volume.Size
	}

	return vols, nil
}

// GetSnapshots collects snapshot data by sending REST call to cinderhost:8776/v1/tenant_id/snapshots
func (s ServiceV1) GetSnapshots(provider *gophercloud.ProviderClient) (types.Snapshots, error) {
	snaps := types.Snapshots{}

	client, err := openstack.NewBlockStorageV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return snaps, err
	}

	opts := snapshots.ListOpts{}

	pager := snapshots.List(client, opts)
	page, err := pager.AllPages()
	if err != nil {
		return snaps, err
	}

	snapshotList, err := snapshots.ExtractSnapshots(page)
	if err != nil {
		return snaps, err
	}

	for _, snapshot := range snapshotList {
		snaps.Count += 1
		snaps.Bytes += snapshot.Size
	}

	return snaps, nil
}
