/*
Copyright 2016 Staples, Inc.

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

package nginx

import (
	"fmt"
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNginxPlugin(t *testing.T) {
	nginxCol := NewNginxCollector()
	Convey("Create nginx Collector", t, func() {
		Convey("So nginxCol should not be nil", func() {
			So(nginxCol, ShouldNotBeNil)
		})
		Convey("So nginxCol should be of nginx type", func() {
			So(nginxCol, ShouldHaveSameTypeAs, &NginxCollector{})
		})
		Convey("nginxCol.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := nginxCol.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, plugin.ConfigPolicy{})
			})
		})
	})
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost/status",
		httpmock.NewStringResponder(200, exampleStatus))

	cfg := plugin.Config{
		"nginx_server_url": "http://localhost/status",
	}

	Convey("Get Nginx Metric Catalog", t, func() {
		metrics, err := nginxCol.GetMetricTypes(cfg)
		So(err, ShouldBeNil)
		for _, mt := range metrics {
			fmt.Println(mt.Namespace.Strings())
		}
	})
}
