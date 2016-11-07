// +build linux

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

	"github.com/gorilla/mux"
	th "github.com/rackspace/gophercloud/testhelper"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"

	"github.com/intelsdi-x/snap-plugin-utilities/str"

	"net/http/httptest"
)

type CollectorSuite struct {
	suite.Suite
	Token                                    string
	V1, V2                                   string
	LimitsV2                                 string
	Tenant1Name, Tenant2Name                 string
	Tenant1ID, Tenant2ID                     string
	MaxTotalVolumeGigabytes, MaxTotalVolumes int
	Vol1, Vol2                               string
	Vol1Size, Vol2Size                       int
	VolMeta                                  string
	SnapShotSize                             int
	server                                   *httptest.Server
}

func (s *CollectorSuite) SetupSuite() {
	// for cinder calls
	th.SetupHTTP()
	// for identity calls
	router := mux.NewRouter()
	s.server = httptest.NewServer(router)

	registerIdentityRoot(s, router)
	s.V1 = "v1/v1ffff"
	s.V2 = "v2/v2ffff"
	s.Token = "2ed210f132564f21b178afb197ee99e3"
	registerIdentityToken(s, router)
	s.Tenant1Name = "admin"
	s.Tenant2Name = "demo"
	s.Tenant1ID = "admin_id123"
	s.Tenant2ID = "demo_id123"
	registerIdentityTenants(s, router)

	registerCinderApi(s)
	registerCinderLimits(s)
	s.Vol1 = "vol1id_123"
	s.Vol2 = "vol2id_321"
	s.Vol1Size = 11
	s.Vol2Size = 22
	registerCinderVolumes(s)
	s.SnapShotSize = 5
	registerCinderSnapshots(s)
}

func (s *CollectorSuite) TearDownSuite() {
	th.TeardownHTTP()
	s.server.Close()
}

func (s *CollectorSuite) TestGetMetricTypes() {
	Convey("Given config with enpoint, user and password defined", s.T(), func() {
		cfg := setupCfg(s.server.URL, "me", "secret", "admin")

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
				So(str.Contains(metricNames, "intel/openstack/cinder/demo/snapshots/count"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/demo/snapshots/bytes"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/demo/volumes/count"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/demo/volumes/bytes"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/demo/limits/MaxTotalVolumeGigabytes"), ShouldBeTrue)
				So(str.Contains(metricNames, "intel/openstack/cinder/demo/limits/MaxTotalVolumes"), ShouldBeTrue)
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

func (s *CollectorSuite) TestCollectMetrics() {

	Convey("Given set of metric types", s.T(), func() {
		cfg := setupCfg(s.server.URL, "me", "secret", "admin")
		m1 := plugin.PluginMetricType{
			Namespace_: []string{"intel", "openstack", "cinder", "demo", "limits", "MaxTotalVolumeGigabytes"},
			Config_:    cfg.ConfigDataNode}
		m2 := plugin.PluginMetricType{
			Namespace_: []string{"intel", "openstack", "cinder", "demo", "volumes", "count"},
			Config_:    cfg.ConfigDataNode}
		m3 := plugin.PluginMetricType{
			Namespace_: []string{"intel", "openstack", "cinder", "demo", "snapshots", "bytes"},
			Config_:    cfg.ConfigDataNode}

		Convey("When ColelctMetrics() is called", func() {
			collector := New()

			mts, err := collector.CollectMetrics([]plugin.PluginMetricType{m1, m2, m3})

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("and proper metric types are returned", func() {
				metricNames := map[string]interface{}{}
				for _, m := range mts {
					ns := strings.Join(m.Namespace(), "/")
					metricNames[ns] = m.Data()
					fmt.Println(ns, "=", m.Data())
				}

				So(len(mts), ShouldEqual, 3)

				val, ok := metricNames["intel/openstack/cinder/demo/limits/MaxTotalVolumeGigabytes"]
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, s.MaxTotalVolumeGigabytes)

				val, ok = metricNames["intel/openstack/cinder/demo/volumes/count"]
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, 1)

				val, ok = metricNames["intel/openstack/cinder/demo/snapshots/bytes"]
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, s.SnapShotSize*1024*1024*1024)

			})
		})
	})
}

func TestCollectorSuite(t *testing.T) {
	collectorTestSuite := new(CollectorSuite)
	suite.Run(t, collectorTestSuite)
}

func setupCfg(endpoint, user, password, tenant string) plugin.PluginConfigType {
	node := cdata.NewNode()
	node.AddItem("endpoint", ctypes.ConfigValueStr{Value: endpoint})
	node.AddItem("user", ctypes.ConfigValueStr{Value: user})
	node.AddItem("password", ctypes.ConfigValueStr{Value: password})
	node.AddItem("tenant", ctypes.ConfigValueStr{Value: tenant})
	return plugin.PluginConfigType{ConfigDataNode: node}
}

func registerIdentityRoot(s *CollectorSuite, r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
				`, s.server.URL+"/v3/", s.server.URL+"/v2.0/")
	})
}

func registerIdentityToken(s *CollectorSuite, r *mux.Router) {
	r.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
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

func registerIdentityTenants(s *CollectorSuite, r *mux.Router) {
	r.HandleFunc("/v2.0/tenants", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Auth-Token", s.Token)

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
	}).Methods("GET")
}

func registerCinderApi(s *CollectorSuite) {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"versions": [
					{
						"id": "v1.0",
						"links": [
							{
								"href": "%s",
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
								"href": "%s",
								"rel": "self"
							}
						],
						"status": "CURRENT",
						"updated": "2012-11-21T11:33:21Z"
					}
				]
			}
			`, th.Endpoint()+"v1", th.Endpoint()+"v2")
	})
}

func registerCinderLimits(s *CollectorSuite) {
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

func registerCinderVolumes(s *CollectorSuite) {
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

func registerCinderSnapshots(s *CollectorSuite) {
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
		`, s.Tenant2ID, s.SnapShotSize)
	})
}
