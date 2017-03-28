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
	"net/http"
	"net/url"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/apiversions"
	"github.com/rackspace/gophercloud/pagination"
)

// According to official documentation it's not paged
func Get(c *gophercloud.ServiceClient) apiversions.APIVersionPage {
	u, _ := url.Parse(c.ServiceURL(""))
	u.Path = "/"

	res := gophercloud.Result{Body: map[string]interface{}{}}
	var resp *http.Response
	resp, res.Err = c.Get(u.String(), &res.Body, &gophercloud.RequestOpts{OkCodes: []int{200, 300}})
	if res.Err == nil {
		res.Header = resp.Header
	}
	return apiversions.APIVersionPage{pagination.SinglePageBase{Result: res, URL: *u}}
}
