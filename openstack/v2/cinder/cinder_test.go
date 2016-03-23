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

package cinder

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"

	openstackintel "github.com/intelsdi-x/snap-plugin-collector-cinder/openstack"
)

type CinderV2Suite struct {
	suite.Suite
	MaxTotalVolumeGigabytes, MaxTotalVolumes int
	V1, V2                                   string
	LimitsV1, LimitsV2                       string
	VolMeta                                  string
	Vol1, Vol2                               string
	Vol1Size, Vol2Size                       int
	Token                                    string
	SnapShotSize                             int
	Tenant1ID, Tenant2ID                     string
}

func (s *CinderV2Suite) SetupSuite() {
	th.SetupHTTP()
	registerRoot()
	s.V1 = "v1/v1ffff"
	s.V2 = "v2/v2ffff"
	s.Token = "2ed210f132564f21b178afb197ee99e3"
	registerAuthentication(s)
	registerLimits(s)

	s.Tenant1ID = "admin_id123"
	s.Tenant2ID = "demo_id123"
	s.Vol1 = "vol1id_123"
	s.Vol2 = "vol2id_321"
	s.Vol1Size = 11
	s.Vol2Size = 22
	registerVolumes(s)
	s.SnapShotSize = 5
	registerSnapshots(s)
}

func (suite *CinderV2Suite) TearDownSuite() {
	th.TeardownHTTP()
}

func TestRunSuite(t *testing.T) {
	cinderTestSuite := new(CinderV2Suite)
	suite.Run(t, cinderTestSuite)
}

func (s *CinderV2Suite) TestGetLimits() {
	Convey("Given Cinder absolute limits are requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, err := openstackintel.Authenticate(th.Endpoint(), "me", "secret", "tenant")
			th.AssertNoErr(s.T(), err)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)

			Convey("and GetLimits called", func() {
				dispatch := ServiceV2{}
				limits, err := dispatch.GetLimits(provider)

				Convey("Then proper limits values are returned", func() {
					So(limits.MaxTotalVolumes, ShouldEqual, s.MaxTotalVolumes)
					So(limits.MaxTotalVolumeGigabytes, ShouldEqual, s.MaxTotalVolumeGigabytes)
				})

				Convey("and no error reported", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})
}

func (s *CinderV2Suite) TestGetVolumes() {
	Convey("Given Cinder volumes are requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, err := openstackintel.Authenticate(th.Endpoint(), "me", "secret", "tenant")
			th.AssertNoErr(s.T(), err)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)

			Convey("and GetVolumes called", func() {
				dispatch := ServiceV2{}
				volumes, err := dispatch.GetVolumes(provider)

				Convey("Then proper limits values are returned", func() {
					So(len(volumes), ShouldEqual, 2)
					So(volumes[s.Tenant1ID].Bytes, ShouldEqual, s.Vol1Size*1024*1024*1024)
					So(volumes[s.Tenant2ID].Bytes, ShouldEqual, s.Vol2Size*1024*1024*1024)
					So(volumes[s.Tenant1ID].Count, ShouldEqual, 1)
					So(volumes[s.Tenant2ID].Count, ShouldEqual, 1)
				})

				Convey("and no error reported", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})
}

func (s *CinderV2Suite) TestGetSnapshots() {
	Convey("Given Cinder snapshots are requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, err := openstackintel.Authenticate(th.Endpoint(), "me", "secret", "tenant")
			th.AssertNoErr(s.T(), err)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)

			Convey("and GetSnapshots called", func() {
				dispatch := ServiceV2{}
				snapshots, err := dispatch.GetSnapshots(provider)

				Convey("Then proper limits values are returned", func() {
					So(len(snapshots), ShouldEqual, 1)
					So(snapshots[s.Tenant1ID].Count, ShouldEqual, 1)
					So(snapshots[s.Tenant1ID].Bytes, ShouldEqual, s.SnapShotSize*1024*1024*1024)
				})

				Convey("and no error reported", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})
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

func registerAuthentication(s *CinderV2Suite) {
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
			th.Endpoint()+s.V2,
			th.Endpoint()+s.V2,
			th.Endpoint()+s.V2,
			th.Endpoint()+s.V1,
			th.Endpoint()+s.V1,
			th.Endpoint()+s.V1,
			s.Token)
	})
}

func registerLimits(s *CinderV2Suite) {
	s.LimitsV2 = "/" + s.V2 + "/limits"
	s.MaxTotalVolumeGigabytes = 1000
	s.MaxTotalVolumes = 10
	th.Mux.HandleFunc(s.LimitsV2, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
				{
					"limits": {
						"absolute": {
							"maxTotalBackupGigabytes": 1000,
							"maxTotalBackups": 10,
							"maxTotalSnapshots": 10,
							"maxTotalVolumeGigabytes": %d,
							"maxTotalVolumes": %d,
							"totalBackupGigabytesUsed": 3,
							"totalBackupsUsed": 1,
							"totalGigabytesUsed": 4,
							"totalSnapshotsUsed": 5,
							"totalVolumesUsed": 2
						},
						"rate": []
					}
				}
			`, s.MaxTotalVolumeGigabytes, s.MaxTotalVolumes)
	})
}

func registerVolumes(s *CinderV2Suite) {
	url := "/v2/v2ffff/volumes/detail" //?all_tenants=true
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestFormValues(s.T(), r, map[string]string{"all_tenants": "true"})
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"volumes": [
					{
						"attachments": [],
						"availability_zone": "nova",
						"bootable": "true",
						"consistencygroup_id": null,
						"created_at": "2016-02-12T10:04:27.000000",
						"description": "Volume for test tenant",
						"encrypted": false,
						"id": "%s",
						"links": [
							{
								"href": "http://192.168.20.2:8776/v2/d98e06adf5db49ad9f372625cad7840b/volumes/1877e478-56bd-4993-80f0-8da9a7e06290",
								"rel": "self"
							},
							{
								"href": "http://192.168.20.2:8776/d98e06adf5db49ad9f372625cad7840b/volumes/1877e478-56bd-4993-80f0-8da9a7e06290",
								"rel": "bookmark"
							}
						],
						"metadata": {},
						"multiattach": false,
						"name": "test_tenant_volume",
						"os-vol-host-attr:host": "rbd:volumes#DEFAULT",
						"os-vol-mig-status-attr:migstat": null,
						"os-vol-mig-status-attr:name_id": null,
						"os-vol-tenant-attr:tenant_id": "%s",
						"os-volume-replication:driver_data": null,
						"os-volume-replication:extended_status": null,
						"replication_status": "disabled",
						"size": %d,
						"snapshot_id": null,
						"source_volid": null,
						"status": "available",
						"user_id": "a3edd7a918fc4373981051c975295dc8",
						"volume_image_metadata": {
							"checksum": "ee1eca47dc88f4879d8a229cc70a07c6",
							"container_format": "bare",
							"disk_format": "qcow2",
							"image_id": "e256d524-bbd7-40af-9bfa-463d86917459",
							"image_name": "TestVM",
							"min_disk": "0",
							"min_ram": "64",
							"size": "13287936"
						},
						"volume_type": null
					},
					{
						"attachments": [],
						"availability_zone": "nova",
						"bootable": "true",
						"consistencygroup_id": null,
						"created_at": "2016-02-09T15:24:27.000000",
						"description": null,
						"encrypted": false,
						"id": "%s",
						"links": [
							{
								"href": "http://192.168.20.2:8776/v2/d98e06adf5db49ad9f372625cad7840b/volumes/ff3e438c-250d-4b03-82ce-3bec50a6c858",
								"rel": "self"
							},
							{
								"href": "http://192.168.20.2:8776/d98e06adf5db49ad9f372625cad7840b/volumes/ff3e438c-250d-4b03-82ce-3bec50a6c858",
								"rel": "bookmark"
							}
						],
						"metadata": {},
						"multiattach": false,
						"name": "test-volume",
						"os-vol-host-attr:host": "rbd:volumes#DEFAULT",
						"os-vol-mig-status-attr:migstat": null,
						"os-vol-mig-status-attr:name_id": null,
						"os-vol-tenant-attr:tenant_id": "%s",
						"os-volume-replication:driver_data": null,
						"os-volume-replication:extended_status": null,
						"replication_status": "disabled",
						"size": %d,
						"snapshot_id": null,
						"source_volid": null,
						"status": "available",
						"user_id": "a3edd7a918fc4373981051c975295dc8",
						"volume_image_metadata": {
							"checksum": "ee1eca47dc88f4879d8a229cc70a07c6",
							"container_format": "bare",
							"disk_format": "qcow2",
							"image_id": "e256d524-bbd7-40af-9bfa-463d86917459",
							"image_name": "TestVM",
							"min_disk": "0",
							"min_ram": "64",
							"size": "13287936"
						},
						"volume_type": null
					}
    			]
       		 }
		`, s.Vol1, s.Tenant1ID, s.Vol1Size, s.Vol2, s.Tenant2ID, s.Vol2Size)
	})
}

func registerSnapshots(s *CinderV2Suite) {
	snapshots := "/v2/v2ffff/snapshots/detail"
	th.Mux.HandleFunc(snapshots, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)
		th.TestFormValues(s.T(), r, map[string]string{"all_tenants": "true"})
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"snapshots": [
					{
						"created_at": "2016-02-21T19:59:15.000000",
						"description": "description",
						"id": "snap1cccc",
						"metadata": {},
						"name": "snapshot_1",
						"os-extended-snapshot-attributes:progress": "100",
            			"os-extended-snapshot-attributes:project_id": "%s",
						"size": %d,
						"status": "available",
						"volume_id": "495a1698-ca2f-4e84-8d34-fa544c65ae3d"
					}
				]
			}
		`, s.Tenant1ID, s.SnapShotSize)
	})
}
