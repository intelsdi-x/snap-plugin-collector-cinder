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

type CinderV1Suite struct {
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

func (s *CinderV1Suite) SetupSuite() {
	th.SetupHTTP()
	registerRoot()
	registerAuthentication(s)
	registerLimits(s)
	registerVolumes(s)
	s.SnapShotSize = 3
	registerSnapshots(s, s.SnapShotSize)
}

func (suite *CinderV1Suite) TearDownSuite() {
	defer th.TeardownHTTP()
}

func TestRunSuite(t *testing.T) {
	cinderTestSuite := new(CinderV1Suite)
	suite.Run(t, cinderTestSuite)
}

func (s *CinderV1Suite) TestGetLimits() {
	Convey("Given Cinder absolute limits are requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, err := openstackintel.Authenticate(th.Endpoint(), "me", "secret", "tenant")
			th.AssertNoErr(s.T(), err)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)

			Convey("and GetLimits called", func() {
				dispatch := ServiceV1{}
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

func (s *CinderV1Suite) TestGetVolumes() {
	Convey("Given Cinder volumes are requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, err := openstackintel.Authenticate(th.Endpoint(), "me", "secret", "tenant")
			th.AssertNoErr(s.T(), err)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)

			Convey("and GetVolumes called", func() {
				dispatch := ServiceV1{}
				volumes, err := dispatch.GetVolumes(provider)

				Convey("Then proper limits values are returned", func() {
					So(volumes.Count, ShouldEqual, 2)
					So(volumes.Bytes, ShouldEqual, s.Vol1Size+s.Vol2Size)
				})

				Convey("and no error reported", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})
}

func (s *CinderV1Suite) TestGetSnapshots() {
	Convey("Given Cinder snapshots are requested", s.T(), func() {

		Convey("When authentication is required", func() {
			provider, err := openstackintel.Authenticate(th.Endpoint(), "me", "secret", "tenant")
			th.AssertNoErr(s.T(), err)
			th.CheckEquals(s.T(), s.Token, provider.TokenID)

			Convey("and GetSnapshots called", func() {
				dispatch := ServiceV1{}
				snapshots, err := dispatch.GetSnapshots(provider)

				Convey("Then proper limits values are returned", func() {
					So(snapshots.Count, ShouldEqual, 1)
					So(snapshots.Bytes, ShouldEqual, s.SnapShotSize)
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

func registerAuthentication(s *CinderV1Suite) {
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

func registerLimits(s *CinderV1Suite) {
	s.LimitsV1 = "/" + s.V1 + "/limits"
	s.MaxTotalVolumeGigabytes = 800
	s.MaxTotalVolumes = 7
	th.Mux.HandleFunc(s.LimitsV1, func(w http.ResponseWriter, r *http.Request) {
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

func registerVolumes(s *CinderV1Suite) {
	s.VolMeta = "/" + s.V1 + "/volumes"
	s.Vol1 = s.V1 + "/volumes/vol1cccc"
	s.Vol2 = s.V1 + "/volumes/vol2cccc"
	s.Vol1Size = 12
	s.Vol2Size = 23
	th.Mux.HandleFunc(s.VolMeta, func(w http.ResponseWriter, r *http.Request) {
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
							"created_at": "2016-02-21T17:02:21.000000",
							"display_description": null,
							"display_name": "vol2_d",
							"encrypted": false,
							"id": "%s",
							"metadata": {},
							"multiattach": "false",
							"size": %d,
							"snapshot_id": null,
							"source_volid": null,
							"status": "available",
							"volume_type": "lvmdriver-1"
						},
						{
							"attachments": [],
							"availability_zone": "nova",
							"bootable": "true",
							"created_at": "2016-02-21T17:00:12.000000",
							"display_description": null,
							"display_name": "vol1_d",
							"encrypted": false,
							"id": "%s",
							"metadata": {},
							"multiattach": "false",
							"size": %d,
							"snapshot_id": null,
							"source_volid": null,
							"status": "available",
							"volume_type": "lvmdriver-1"
						}
					]
				}
		`, th.Endpoint()+s.Vol2, s.Vol2Size, th.Endpoint()+s.Vol1, s.Vol1Size)
	})

}

func registerSnapshots(s *CinderV1Suite, snapSize int) {
	snapshots := "/" + s.V1 + "/snapshots"
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
