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
// - Snapshot structure:
//   - renamed Metadata field to Meta
//   - renamed CreatedAt field to Created
//   - removed SnaphotID field
//   - removed VolumeType field
//   - removed Bootable field
//   - removed AvailabilityZone field
//   - removed Attachments field
//   - removed original field comments
//   - added OsExtendedSnapshotAttributesProgress field
//   - added OsExtendedSnapshotAttributesProjectID field
package snapshots

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// Snapshot contains information associated with an OpenStack Snapshot.
type Snapshot struct {
	Created                               string                 `mapstructure:"created_at"`
	Description                           string                 `mapstructure:"description"`
	ID                                    string                 `mapstructure:"id"`
	Meta                                  map[string]interface{} `mapstructure:"metadata"`
	Name                                  string                 `mapstructure:"name"`
	OsExtendedSnapshotAttributesProgress  string                 `mapstructure:"os-extended-snapshot-attributes:progress"`
	OsExtendedSnapshotAttributesProjectID string                 `mapstructure:"os-extended-snapshot-attributes:project_id"`
	Status                                string                 `mapstructure:"status"`
	Size                                  int                    `mapstructure:"size"`
	VolumeID                              string                 `mapstructure:"volume_id"`
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// ListResult is a pagination.Pager that is returned from a call to the List function.
type ListResult struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no Snapshots.
func (r ListResult) IsEmpty() (bool, error) {
	volumes, err := ExtractSnapshots(r)
	if err != nil {
		return true, err
	}
	return len(volumes) == 0, nil
}

// ExtractSnapshots extracts and returns Snapshots. It is used while iterating over a snapshots.List call.
func ExtractSnapshots(page pagination.Page) ([]Snapshot, error) {
	var response struct {
		Snapshots []Snapshot `json:"snapshots"`
	}

	err := mapstructure.Decode(page.(ListResult).Body, &response)
	return response.Snapshots, err
}

type commonResult struct {
	gophercloud.Result
}
