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

// requests contains Cinder API requests for limits/quotas

package limits

import (
	"github.com/rackspace/gophercloud"
)

// Get prepares http GET call on Cinder endpoint
func Get(client *gophercloud.ServiceClient, limits string) GetResult {
	var res GetResult
	_, err := client.Get(client.ResourceBaseURL()+limits, &res.Body, nil)
	res.Err = err
	return res
}
