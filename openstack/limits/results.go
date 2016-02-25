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

// results contains Cinder API responses and their processing for limits/quotas

package limits

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

// GetResult contains the response body and error from a Get request
type GetResult struct {
	commonResult
}

// Extract will get the limit object out of the commonResult object
func (r commonResult) Extract() (limits, error) {
	tenantLimits := limits{}
	if r.Err != nil {
		return tenantLimits, r.Err
	}

	var res struct {
		Absolute absolute `mapstructure:"limits"`
	}

	err := mapstructure.Decode(r.Body, &res)
	if err != nil {
		return tenantLimits, r.Err
	}

	return res.Absolute.Limits, err
}

type limits struct {
	TotalSnapshotsUsed       int `mapstructure:"totalSnapshotsUsed"`
	MaxTotalBackups          int `mapstructure:"maxTotalBackups"`
	MaxTotalVolumeGigabytes  int `mapstructure:"maxTotalVolumeGigabytes"`
	MaxTotalSnapshots        int `mapstructure:"maxTotalSnapshots"`
	MaxTotalBackupGigabytes  int `mapstructure:"maxTotalBackupGigabytes"`
	TotalBackupGigabytesUsed int `mapstructure:"totalBackupGigabytesUsed"`
	MaxTotalVolumes          int `mapstructure:"maxTotalVolumes"`
	TotalVolumesUsed         int `mapstructure:"totalVolumesUsed"`
	TotalBackupsUsed         int `mapstructure:"totalBackupsUsed"`
	TotalGigabytesUsed       int `mapstructure:"totalGigabytesUsed"`
}

type absolute struct {
	Limits limits `mapstructure:"absolute"`
	Rate   []rate `mapstructure:"rate"`
}

type rate struct{}

type commonResult struct {
	gophercloud.Result
}
