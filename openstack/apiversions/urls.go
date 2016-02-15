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

// urls generation for Cinder API versions

package apiversions

import (
	"net/url"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/apiversions"
)

// Exactly the same as original, but kept here to achieve proper call to modified listURL() until PR is merged
func List(c *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, listURL(c), func(r pagination.PageResult) pagination.Page {
		return apiversions.APIVersionPage{pagination.SinglePageBase(r)}
	})
}

func listURL(c *gophercloud.ServiceClient) string {
	u, _ := url.Parse(c.ServiceURL(""))
	u.Path = "/"
	return u.String()
}
