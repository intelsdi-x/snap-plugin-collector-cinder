// +build small

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

package openstack

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"

	th "github.com/rackspace/gophercloud/testhelper"
)

type CommonSuite struct {
	suite.Suite
	Token                    string
	BlockStorageEndpoint     string
	V1, V2                   string
	Tenant1ID, Tenant2ID     string
	Tenant1Name, Tenant2Name string
}

func (s *CommonSuite) SetupSuite() {
	th.SetupHTTP()
	registerRoot()
	registerAuthentication(s)
	s.Tenant1Name = "admin"
	s.Tenant2Name = "demo"
	s.Tenant1ID = "3e3e3e"
	s.Tenant2ID = "4f4f4f"
	registerTenants(s)
}

func (s *CommonSuite) TearDownSuite() {
	th.TeardownHTTP()
}

func (s *CommonSuite) TestGetTenants() {
	Convey("Given tenants are requested", s.T(), func() {
		c := Common{}
		Convey("When Gettenants is called", func() {
			tenants, err := c.GetTenants(th.Endpoint(), "me", "secret", "", "")

			Convey("Then list of available tenats is returned", func() {
				So(len(tenants), ShouldEqual, 2)
				So(tenants[s.Tenant1ID], ShouldEqual, s.Tenant1Name)
				So(tenants[s.Tenant2ID], ShouldEqual, s.Tenant2Name)
				So(err, ShouldBeNil)
			})
		})
	})
}

func (s *CommonSuite) TestGetAPI() {
	Convey("Given api versions are requested", s.T(), func() {
		c := Common{}
		Convey("When GetAPIVersions is called", func() {
			provider, err := Authenticate(th.Endpoint(), "me", "secret", "tenant", "", "")
			th.AssertNoErr(s.T(), err)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				fmt.Fprintf(w, `
					{
						"versions": [
							{
								"id": "v1.0",
								"links": [
									{
										"href": "http://192.168.20.2:8776/v1/",
										"rel": "self"
									}
								],
								"status": "SUPPORTED",
								"updated": "2014-06-28T12:20:21Z"
							},
							{
								"id": "v2.0",
								"links": [
									{
										"href": "http://192.168.20.2:8776/v2/",
										"rel": "self"
									}
								],
								"status": "CURRENT",
								"updated": "2012-11-21T11:33:21Z"
							}
						]
					}
				`)
			}))
			defer server.Close()
			transport := &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return url.Parse(server.URL)
				},
			}

			httpClient := http.Client{Transport: transport}
			provider.HTTPClient = httpClient
			apis, err := c.GetApiVersions(provider)

			Convey("Then list of available versions is returned", func() {
				So(len(apis), ShouldEqual, 2)
				So(apis[0], ShouldEqual, "v1.0")
				So(apis[1], ShouldEqual, "v2.0")
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestCommonSuite(t *testing.T) {
	commonTestSuite := new(CommonSuite)
	suite.Run(t, commonTestSuite)
}

func registerRoot() {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
				{
					"versions": {
						"values": [
							{
								"status": "experimental",
								"id": "v3.0",
								"links": [
									{ "href": "%s", "rel": "self" }
								]
							},
							{
								"status": "stable",
								"id": "v2.0",
								"links": [
									{ "href": "%s", "rel": "self" }
								]
							}
						]
					}
				}
				`, th.Endpoint()+"v3/", th.Endpoint()+"v2.0/")
	})
}

func registerAuthentication(s *CommonSuite) {
	s.Token = "2ed210f132564f21b178afb197ee99e3"
	s.BlockStorageEndpoint = "http://127.0.0.1:8080"
	th.Mux.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
				{
					"access": {
						"metadata": {
							"is_admin": 0,
							"roles": [
								"3083d61996d648ca88d6ff420542f324"
							]
						},
						"serviceCatalog": [
							{
								"endpoints": [
									{
										"adminURL": "%s",
										"id": "3ffe125aa59547029ed774c10b932349",
										"internalURL": "%s",
										"publicURL": "%s",
										"region": "RegionOne"
									}
								],
								"endpoints_links": [],
								"name": "cinderv2",
								"type": "volumev2"
							},
							{
								"endpoints": [
									{
										"adminURL": "%s",
										"id": "a056ce874d414393a946e42e920ce157",
										"internalURL": "%s",
										"publicURL": "%s",
										"region": "RegionOne"
									}
								],
								"endpoints_links": [],
								"name": "cinder",
								"type": "volume"
							}

						],
						"token": {
							"expires": "2016-02-21T14:28:30Z",
							"id": "%s",
							"issued_at": "2016-02-21T13:28:30.656527",
							"tenant": {
								"description": null,
								"enabled": true,
								"id": "97ea299c37bb4e04b3779039ea8aba44",
								"name": "tenant"
							}
						}
					}
				}
			`,
			s.BlockStorageEndpoint,
			s.BlockStorageEndpoint,
			s.BlockStorageEndpoint,
			s.BlockStorageEndpoint,
			s.BlockStorageEndpoint,
			s.BlockStorageEndpoint,
			s.Token)
	})
}

func registerTenants(s *CommonSuite) {
	th.Mux.HandleFunc("/v2.0/tenants", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"tenants": [
					{
						"description": "Test tenat",
						"enabled": true,
						"id": "%s",
						"name": "%s"
					},
					{
						"description": "admin tenant",
						"enabled": true,
						"id": "%s",
						"name": "%s"
					}
				],
				"tenants_links": []
			}
		`, s.Tenant1ID, s.Tenant1Name, s.Tenant2ID, s.Tenant2Name)
	})
}

func registerAPI(s *CommonSuite) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"versions": [
					{
						"id": "v1.0",
						"links": [
							{
								"href": "http://192.168.20.2:8776/v1/",
								"rel": "self"
							}
						],
						"status": "SUPPORTED",
						"updated": "2014-06-28T12:20:21Z"
					},
					{
						"id": "v2.0",
						"links": [
							{
								"href": "http://192.168.20.2:8776/v2/",
								"rel": "self"
							}
						],
						"status": "CURRENT",
						"updated": "2012-11-21T11:33:21Z"
					}
				]
			}
		`)
	})
}
