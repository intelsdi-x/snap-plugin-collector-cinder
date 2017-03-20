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

This file incorporates work covered by the following copyright and permission notice:

Copyright 2012-2013 Rackspace, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

// Package contains code from Rackspace Gophercloud (https://github.com/rackspace/gophercloud) with following changes:
// - Volume structure:
//   - changed field order
//   - added VolImageMeta field
//   - added Links field
//   - added OsVolHostAttrHost field
//   - added OsVolMigStatusAttrMigstat field
//   - added OsVolMigStatusAttrNameID field
//   - added OsVolTenantAttrTenantID field
//   - added OsVolumeReplicationDriverData field
//   - added OsVolumeReplicationExtendedStatus field
package volumes

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// Volume contains information associated with an OpenStack Volume
type Volume struct {
	// Current status of the volume.
	Status string `mapstructure:"status"`

	// Human-readable display name for the volume.
	Name string `mapstructure:"name"`

	// Instances onto which the volume is attached.
	Attachments []map[string]interface{} `mapstructure:"attachments"`

	// This parameter is no longer used.
	AvailabilityZone string `mapstructure:"availability_zone"`

	// Indicates whether this is a bootable volume.
	Bootable string `mapstructure:"bootable"`

	// The date when this volume was created.
	CreatedAt string `mapstructure:"created_at"`

	// Human-readable description for the volume.
	Description string `mapstructure:"description"`

	// The type of volume to create, either SATA or SSD.
	VolumeType string `mapstructure:"volume_type"`

	// The ID of the snapshot from which the volume was created
	SnapshotID string `mapstructure:"snapshot_id"`

	// The ID of another block storage volume from which the current volume was created
	SourceVolID string `mapstructure:"source_volid"`

	// Arbitrary key-value pairs defined by the user.
	Metadata map[string]string `mapstructure:"metadata"`

	// Unique identifier for the volume.
	ID string `mapstructure:"id"`

	// Size of the volume in GB.
	Size int `mapstructure:"size"`

	// Volume Image key-value pairs meta data
	VolImageMeta map[string]string `mapstructure:"volume_image_metadata"`

	// User ID
	UserID string `mapstructure:"user_id"`

	// If true volume is encrypted
	Encrypted bool `json:"encrypted" mapstructure:"encrypted"`

	// Volume links
	Links []map[string]interface{} `json:"links" mapstructure:"links"`

	// If true, this volume can attach to more than one instance.
	MultiAttach bool `json:"multiattach" mapstructure:"multiattach"`

	// The UUID of the consistency group
	ConsistencyGroupId string `json:"consistencygroup_id" mapstructure:"consistencygroup_id"`

	// Current back-end of the volume
	OsVolHostAttrHost string `json:"os-vol-host-attr:host" mapstructure:"os-vol-host-attr:host"`

	// The status of this volume migratio
	OsVolMigStatusAttrMigstat string `json:"os-vol-mig-status-attr:migstat" mapstructure:"os-vol-mig-status-attr:migstat"`

	// The volume ID that this volume name on the back-end is based on
	OsVolMigStatusAttrNameID string `json:"os-vol-mig-status-attr:name_id" mapstructure:"os-vol-mig-status-attr:name_id"`

	// The tenant ID which the volume belongs to
	OsVolTenantAttrTenantID string `json:"os-vol-tenant-attr:tenant_id" mapstructure:"os-vol-tenant-attr:tenant_id"`

	// The name of the volume replication driver
	OsVolumeReplicationDriverData string `json:"os-volume-replication:driver_data" mapstructure:"os-volume-replication:driver_data"`

	// The volume replication status managed by the driver of backend storage
	OsVolumeReplicationExtendedStatus string `json:"os-volume-replication:extended_status" mapstructure:"os-volume-replication:extended_status"`

	// The volume replication status
	ReplicationStatus string `json:"replication_status" mapstructure:"replication_status"`
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// ListMetaResult is a pagination.pager that is returned from a call to the ListMeta function.
type ListResult struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no Volumes.
func (r ListResult) IsEmpty() (bool, error) {
	volumes, err := ExtractVolumes(r)
	if err != nil {
		return true, err
	}
	return len(volumes) == 0, nil
}

// ExtractVolumes extracts and returns Volumes. It is used while iterating over a volumes.List call.
func ExtractVolumes(page pagination.Page) ([]Volume, error) {
	var response struct {
		Volumes []Volume `mapstructure:"volumes"`
	}

	err := mapstructure.Decode(page.(ListResult).Body, &response)

	return response.Volumes, err
}

// Extract will get the Volume object out of the commonResult object.
func (r commonResult) Extract() (*Volume, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Volume *Volume `json:"volume"`
	}

	err := mapstructure.Decode(r.Body, &res)

	return res.Volume, err
}

type commonResult struct {
	gophercloud.Result
}
