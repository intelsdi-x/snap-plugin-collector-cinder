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
	"strings"

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

// GetVolumes collects volumes data by sending REST call to cinderhost:8776/v2/tenant_id/volumes
func (s ServiceV2) GetVolumes(provider *gophercloud.ProviderClient) (types.Volumes, error) {
	vols := types.Volumes{}

	client, err := openstackintel.NewBlockStorageV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return vols, err
	}

	opts := volumesintel.ListOpts{}

	pager := volumesintel.List(client, opts)
	page, err := pager.AllPages()
	if err != nil {
		return vols, err
	}

	volumesMeta, err := volumesintel.ExtractVolumesMeta(page)
	if err != nil {
		return vols, err
	}

	for _, volumeMeta := range volumesMeta {
		volURL := client.ResourceBaseURL() + strings.Join([]string{"volumes", volumeMeta.ID}, "/")
		volume, err := volumesintel.Get(client, volURL).Extract()
		if err != nil {
			panic(err)
		}
		vols.Count += 1
		vols.Bytes += volume.Size * 1024 * 1024 * 1024
	}

	return vols, nil
}

// GetSnapshots collects snapshot data by sending REST call to cinderhost:8776/v2/tenant_id/snapshots
func (s ServiceV2) GetSnapshots(provider *gophercloud.ProviderClient) (types.Snapshots, error) {
	snaps := types.Snapshots{}

	client, err := openstackintel.NewBlockStorageV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return snaps, err
	}

	opts := snapshotsintel.ListOpts{}
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
		snaps.Count += 1
		snaps.Bytes += snapshot.Size * 1024 * 1024 * 1024
	}

	return snaps, nil
}
