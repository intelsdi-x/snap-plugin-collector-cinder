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
// All functions are dependant on OpenStack BlockStorage API Version 2

package cinder

import (
	"github.com/rackspace/gophercloud"

	limitsintel "github.com/intelsdi-x/snap-plugin-collector-cinder/openstack/limits"
	openstackintel "github.com/intelsdi-x/snap-plugin-collector-cinder/openstack/v2"
	snapshotsintel "github.com/intelsdi-x/snap-plugin-collector-cinder/openstack/v2/snapshots"
	volumesintel "github.com/intelsdi-x/snap-plugin-collector-cinder/openstack/v2/volumes"
	"github.com/intelsdi-x/snap-plugin-collector-cinder/types"
)

// ServiceV2 serves as dispatcher for Cinder API version 2.0
type ServiceV2 struct{}

// GetLimits collects tenant limits by sending REST call to cinderhost:8776/v2/tenant_id/limits
func (s ServiceV2) GetLimits(provider *gophercloud.ProviderClient) (types.Limits, error) {
	limits := types.Limits{}

	client, err := openstackintel.NewBlockStorageV2(provider, gophercloud.EndpointOpts{})
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

// GetVolumes collects volumes data by sending REST call to cinderhost:8776/v2/tenant_id/volumes/detail?all_tenants=true
func (s ServiceV2) GetVolumes(provider *gophercloud.ProviderClient) (map[string]types.Volumes, error) {
	vols := map[string]types.Volumes{}

	client, err := openstackintel.NewBlockStorageV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, err
	}

	opts := volumesintel.ListOpts{AllTenants: true}

	pager := volumesintel.List(client, opts)
	page, err := pager.AllPages()
	if err != nil {
		return nil, err
	}

	volumes, err := volumesintel.ExtractVolumes(page)
	if err != nil {
		return nil, err
	}

	for _, volume := range volumes {
		volCounts := vols[volume.OsVolTenantAttrTenantID]
		volCounts.Count += 1
		volCounts.Bytes += volume.Size * 1024 * 1024 * 1024
		vols[volume.OsVolTenantAttrTenantID] = volCounts
	}

	return vols, nil
}

// GetSnapshots collects snapshot data by sending REST call to cinderhost:8776/v2/tenant_id/snapshots/detail?all_tenants=true
func (s ServiceV2) GetSnapshots(provider *gophercloud.ProviderClient) (map[string]types.Snapshots, error) {
	snaps := map[string]types.Snapshots{}

	client, err := openstackintel.NewBlockStorageV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return snaps, err
	}

	opts := snapshotsintel.ListOpts{AllTenants: true}
	pager := snapshotsintel.List(client, opts)
	page, err := pager.AllPages()
	if err != nil {
		return snaps, err
	}

	snapshotList, err := snapshotsintel.ExtractSnapshots(page)
	if err != nil {
		return snaps, err
	}

	for _, snapshot := range snapshotList {
		snapCounts := snaps[snapshot.OsExtendedSnapshotAttributesProjectID]
		snapCounts.Count += 1
		snapCounts.Bytes += snapshot.Size * 1024 * 1024 * 1024
		snaps[snapshot.OsExtendedSnapshotAttributesProjectID] = snapCounts
	}

	return snaps, nil
}
