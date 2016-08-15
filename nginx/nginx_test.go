package nginx

import (
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNginxPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, pluginName)
		So(meta.Version, ShouldResemble, pluginVersion)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})

	Convey("Create nginx Collector", t, func() {
		nginxCol := NewNginxCollector()
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
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
		})
	})
}
