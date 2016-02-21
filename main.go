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

package main

import (
	//"os"

	"github.com/intelsdi-x/snap/control/plugin"

	"github.com/intelsdi-x/snap-plugin-collector-cinder/collector"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	"fmt"
)

const (
	endpoint = "http://localhost:5000/v2.0"
	user = "admin"
	password = "openstack"
)

func main() {
	plg := collector.New()

//	plugin.Start(
//		collector.Meta(),
//		plg,
//		os.Args[1],
//	)

	node := cdata.NewNode()
	node.AddItem("endpoint", ctypes.ConfigValueStr{Value: endpoint})
	node.AddItem("user", ctypes.ConfigValueStr{Value: user})
	node.AddItem("password", ctypes.ConfigValueStr{Value: password})
	node.AddItem("tenant", ctypes.ConfigValueStr{Value: "demo"})
	cfg := plugin.PluginConfigType{ConfigDataNode: node}
	//m1 := plugin.PluginMetricType{Namespace_: []string{"intel", "openstack", "cinder", "demo", "limits", "MaxTotalVolumeGigabytes"}, Config_: cfg}
	//m2 := plugin.PluginMetricType{Namespace_: []string{"intel", "openstack", "cinder", "demo", "volumes", "count"}, Config_: cfg}
	//m3 := plugin.PluginMetricType{Namespace_: []string{"intel", "openstack", "cinder", "demo", "snapshots", "bytes"}, Config_: cfg}
	mts , err := plg.GetMetricTypes(cfg)
	if err != nil {panic(err)}
	for _, m := range mts {
		fmt.Println(m.Namespace())
	}
	//panic("getmet")
	mmm, _ := plg.CollectMetrics(mts)
	if err != nil {panic(err)}
	for _, m := range mmm {
		fmt.Println(m.Namespace()[3:], "=", m.Data())
	}

}
