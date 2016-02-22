// +build unit

/*
http://www.apache.org/licenses/LICENSE-2.0.txt
Copyright 2015 Intel Corporation
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

package collector

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"

	"github.com/intelsdi-x/snap-plugin-collector-cinder/types"
	str "github.com/intelsdi-x/snap-plugin-utilities/strings"
	"github.com/rackspace/gophercloud"
)

type CollectorSuite struct {
	suite.Suite
	Token                                    string
	V1, V2                                   string
	LimitsV2                                 string
	Tenant1, Tenant2                         string
	MaxTotalVolumeGigabytes, MaxTotalVolumes int
	Vol1, Vol2                               string
	Vol1Size, Vol2Size                       int
	VolMeta                                  string
	SnapShotSize                             int
}

func (s *CollectorSuite) SetupSuite() {
	th.SetupHTTP()
	registerRoot()
	registerAuthentication(s)
	registerTenants(s, "1fffff", "2eeeee")
	registerLimits(s)
	registerVolMeta(s)
	s.Vol1Size = 11
	s.Vol2Size = 22
	registerVolume(s, s.Vol1, 11)
	registerVolume(s, s.Vol2, 22)
	s.SnapShotSize = 5
	registerSnapshots(s, s.SnapShotSize)
}

func (s *CollectorSuite) TearDownSuite() {
	th.TeardownHTTP()
}

func (s *CollectorSuite) TestGetMetricTypes() {
	Convey("Given config with enpoint, user and password defined", s.T(), func() {
		cfg := setupCfg(th.Endpoint(), "me", "secret")

		Convey("When GetMetricTypes() is called", func() {
			collector := New()
			mts, err := collector.GetMetricTypes(cfg)

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("and proper metric types are returned", func() {
				metricNames := []string{}
				for _, m := range mts {
					metricNames = append(metricNames, strings.Join(m.Namespace(), "/"))
				}

				So(len(mts), ShouldEqual, 12)
				So(str.Contains(metricNames, "intel/openstack/cinder/test_tenant/snapshots/count"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/test_tenant/snapshots/bytes"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/test_tenant/volumes/count"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/test_tenant/volumes/bytes"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/test_tenant/limits/MaxTotalVolumeGigabytes"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/test_tenant/limits/MaxTotalVolumes"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/admin/volumes/count"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/admin/volumes/bytes"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/admin/limits/MaxTotalVolumeGigabytes"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/admin/limits/MaxTotalVolumes"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/admin/snapshots/count"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/admin/snapshots/bytes"), ShouldBeTrue)
			})
		})
	})
}

//func (s *CollectorSuite) TestCollectMetrics() {
//	Convey("Given set of metric types", s.T(), func() {
//		cfg := setupCfg(th.Endpoint(), "me", "secret")
//		m1 := plugin.PluginMetricType{
//			Namespace_: []string{"intel", "openstack", "cinder", "demo", "limits", "MaxTotalVolumeGigabytes"},
//			Config_: cfg.ConfigDataNode}
//		//m2 := plugin.PluginMetricType{
//		//	Namespace_: []string{"intel", "openstack", "cinder", "demo", "volumes", "count"},
//		//	Config_: &cfg.ConfigDataNode}
//		//m3 := plugin.PluginMetricType{
//		//	Namespace_: []string{"intel", "openstack", "cinder", "demo", "snapshots", "bytes"},
//		//	Config_: &cfg.ConfigDataNode}
//		//
//
//
//		servMock := ServicesMock{}
//		limits := types.Limits{
//			MaxTotalVolumeGigabytes: 333,
//			MaxTotalVolumes: 111,
//		}
//		servMock.On("GetLimits", mock.AnythingOfType("*gophercloud.ProviderClient")).Return(limits, nil)
//
//
//		cmnMock := CommonMock{}
//		tenants := []types.Tenant{types.Tenant{ID: "1fffff", Name: "demo"}}
//		cmnMock.On("GetTenants", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(tenants, nil)
//		cmnMock.On("GetApiVersions", mock.AnythingOfType("*gophercloud.ProviderClient")).Return([]string{"v1.0", "v2.0"}, nil)
//
//		Convey("When ColelctMetrics() is called", func() {
//			collector := New()
//
//			collector.common = cmnMock
//			collector.service = servMock
//
//			mts, err := collector.CollectMetrics([]plugin.PluginMetricType{m1})
//
//			Convey("Then no error should be reported", func() {
//				So(err, ShouldBeNil)
//			})
//
//			Convey("and proper metric types are returned", func() {
//				metricNames := map[string]interface{}{}
//				for _, m := range mts {
//					ns := strings.Join(m.Namespace(), "/")
//					metricNames[ns] = m.Data()
//				}
//
//				So(len(mts), ShouldEqual, 1)
//
//			})
//		})
//	})
//}

func TestCollectorSuite(t *testing.T) {
	collectorTestSuite := new(CollectorSuite)
	suite.Run(t, collectorTestSuite)
}

type ServicesMock struct {
	mock.Mock
}

func (servMock ServicesMock) GetLimits(provider *gophercloud.ProviderClient) (types.Limits, error) {
	ret := servMock.Mock.Called(provider)
	return ret.Get(0).(types.Limits), ret.Error(1)
}

func (servMock ServicesMock) GetVolumes(provider *gophercloud.ProviderClient) (types.Volumes, error) {
	ret := servMock.Mock.Called(provider)
	return ret.Get(0).(types.Volumes), ret.Error(1)
}

func (servMock ServicesMock) GetSnapshots(provider *gophercloud.ProviderClient) (types.Snapshots, error) {
	ret := servMock.Mock.Called(provider)
	return ret.Get(0).(types.Snapshots), ret.Error(1)
}

type CommonMock struct {
	mock.Mock
}

func (cmnMock CommonMock) GetTenants(endpoint, user, password string) ([]types.Tenant, error) {
	ret := cmnMock.Mock.Called(endpoint, user, password)
	return ret.Get(0).([]types.Tenant), ret.Error(1)
}

func (cmnMock CommonMock) GetApiVersions(provider *gophercloud.ProviderClient) ([]string, error) {
	ret := cmnMock.Mock.Called(provider)
	return ret.Get(0).([]string), ret.Error(1)
}

func setupCfg(endpoint, user, password string) plugin.PluginConfigType {
	node := cdata.NewNode()
	node.AddItem("endpoint", ctypes.ConfigValueStr{Value: endpoint})
	node.AddItem("user", ctypes.ConfigValueStr{Value: user})
	node.AddItem("password", ctypes.ConfigValueStr{Value: password})
	return plugin.PluginConfigType{ConfigDataNode: node}
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

func registerAuthentication(s *CollectorSuite) {
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

func registerTenants(s *CollectorSuite, tenant1 string, tenant2 string) {
	s.Tenant1 = tenant1
	s.Tenant2 = tenant2
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
						"name": "test_tenant"
					},
					{
						"description": "admin tenant",
						"enabled": true,
						"id": "%s",
						"name": "admin"
					}
				],
				"tenants_links": []
			}
		`, s.Tenant1, s.Tenant2)
	})
}

func registerLimits(s *CollectorSuite) {
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

func registerVolMeta(s *CollectorSuite) {
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

func registerVolume(s *CollectorSuite, volID string, volSize int) {
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

func registerSnapshots(s *CollectorSuite, snapSize int) {
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
