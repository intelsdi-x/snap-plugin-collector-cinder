// +build unit

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
}

func (s *CinderV2Suite) SetupSuite() {
	th.SetupHTTP()
	registerRoot()
	registerAuthentication(s)
	registerLimits(s)
	registerVolMeta(s)
	s.Vol1Size = 11
	s.Vol2Size = 22
	registerVolume(s, s.Vol1, 11)
	registerVolume(s, s.Vol2, 22)
	s.SnapShotSize = 5
	registerSnapshots(s, s.SnapShotSize)
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
					So(volumes.Count, ShouldEqual, 2)
					So(volumes.Bytes, ShouldEqual, (s.Vol1Size+s.Vol2Size)*1024*1024*1024)
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
					So(snapshots.Count, ShouldEqual, 1)
					So(snapshots.Bytes, ShouldEqual, s.SnapShotSize*1024*1024*1024)
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
	s.V1 = "v1/v1ffff"
	s.V2 = "v2/v2ffff"
	s.Token = "2ed210f132564f21b178afb197ee99e3"
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

func registerVolMeta(s *CinderV2Suite) {
	s.VolMeta = "/" + s.V2 + "/volumes"
	s.Vol1 = s.V2 + "/volumes/vol1cccc"
	s.Vol2 = s.V2 + "/volumes/vol2cccc"
	th.Mux.HandleFunc(s.VolMeta, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
				{
					"volumes": [
						{
							"id": "vol2cccc",
							"links": [
								{
									"href": "%s",
									"rel": "self"
								},
								{
									"href": "%s",
									"rel": "bookmark"
								}
							],
							"name": "vol2"
						},
						{
							"id": "vol1cccc",
							"links": [
								{
									"href": "%s",
									"rel": "self"
								},
								{
									"href": "%s",
									"rel": "bookmark"
								}
							],
							"name": "vol1"
						}
					]
				}
			`, th.Endpoint()+s.Vol2, th.Endpoint()+s.Vol2, th.Endpoint()+s.Vol1, th.Endpoint()+s.Vol1)
	})

}

func registerVolume(s *CinderV2Suite, volID string, volSize int) {
	th.Mux.HandleFunc("/"+volID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
			  "volume": {
				"attachments": [],
				"availability_zone": "nova",
				"bootable": "true",
				"consistencygroup_id": null,
				"created_at": "2016-02-21T17:00:12.000000",
				"description": null,
				"encrypted": false,
				"id": "vol1cccc",
				"links": [
				  {
					"href": "%s",
					"rel": "self"
				  },
				  {
					"href": "%s",
					"rel": "bookmark"
				  }
				],
				"metadata": {},
				"multiattach": false,
				"name": "vol1",
				"os-vol-host-attr:host": "devstack@lvmdriver-1#lvmdriver-1",
				"os-vol-mig-status-attr:migstat": null,
				"os-vol-mig-status-attr:name_id": null,
				"os-vol-tenant-attr:tenant_id": "97ea299c37bb4e04b3779039ea8aba44",
				"os-volume-replication:driver_data": null,
				"os-volume-replication:extended_status": null,
				"replication_status": "disabled",
				"size": %d,
				"snapshot_id": null,
				"source_volid": null,
				"status": "available",
				"user_id": "7379713c4da04e88af09fe8c7f2077dc",
				"volume_image_metadata": {
				  "checksum": "eb9139e4942121f22bbc2afc0400b2a4",
				  "container_format": "ami",
				  "disk_format": "ami",
				  "image_id": "550f1f21-4b58-4fc3-9158-105487b2d5e8",
				  "image_name": "cirros-0.3.4-x86_64-uec",
				  "kernel_id": "30e61305-a024-4182-9a8c-697c08b3d73d",
				  "min_disk": "0",
				  "min_ram": "0",
				  "ramdisk_id": "f78c9363-0fde-4a48-9899-2071505da7d5",
				  "size": "25165824"
				},
				"volume_type": "lvmdriver-1"
			  }
			}
		`, th.Endpoint()+volID, th.Endpoint()+volID, volSize)
	})

}

func registerSnapshots(s *CinderV2Suite, snapSize int) {
	snapshots := "/" + s.V2 + "/snapshots"
	th.Mux.HandleFunc(snapshots, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

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
						"size": %d,
						"status": "available",
						"volume_id": "495a1698-ca2f-4e84-8d34-fa544c65ae3d"
					}
				]
			}
		`, snapSize)
	})
}
