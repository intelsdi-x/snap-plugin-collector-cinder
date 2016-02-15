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

package volumes

import (
    "github.com/rackspace/gophercloud"
    "github.com/rackspace/gophercloud/pagination"
    "github.com/mitchellh/mapstructure"
)

// Volume contains information associated with an OpenStack Volume
type Volume struct{
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
}

// Volume contains information associated with an OpenStack Volume metadata
type VolumeMeta struct {
	ID string `mapstructure:"id"`
	Links []map[string]interface{}  `mapstructure:"links"`
	Name string `mapstructure:"name"`

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
	volumes, err := ExtractVolumesMeta(r)
	if err != nil {
		return true, err
	}
	return len(volumes) == 0, nil
}

// ExtractVolumes extracts and returns Volumes. It is used while iterating over a volumes.List call.
func ExtractVolumesMeta(page pagination.Page) ([]VolumeMeta, error) {
	var response struct {
		VolumesMeta []VolumeMeta `mapstructure:"volumes"`
	}

	err := mapstructure.Decode(page.(ListResult).Body, &response)
	return response.VolumesMeta, err
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