package nginx

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	// Name of plugin
	Name = "nginx"
	// Version of plugin
	Version = 1
	// Type of plugin
	Type = plugin.CollectorPluginType
)

var (
	errNoWebserver     = errors.New("nginx_status_url config required. Check your config JSON file")
	errBadWebserver    = errors.New("Failed to parse given nginx_status_url")
	errReqFailed       = errors.New("Request to nginx webserver failed")
	errConfigReadError = errors.New("Config read Error")
)

// make sure that we actually satisify required interface
var _ plugin.CollectorPlugin = (*Nginx)(nil)

type Nginx struct{}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getHostName(inData interface{}, hostName string) string {
	hostName = fmt.Sprintf("host_id_%s", hostName)
	switch mtype := inData.(type) {
	case map[string]interface{}:
		val := mtype["server"].(string)
		//Default hostname with port will be encoded to md5
		//If LookUpAddr get valid hostname it will be replace
		hostName = fmt.Sprintf("host_id_%s", getMD5Hash(val))

		//check for IPV4
		if strings.Count(val, ".") == 3 {
			subStr := strings.Split(val, ":")
			hName, err := net.LookupAddr(subStr[0])
			if err == nil {
				hostName = strings.Join(hName, ".")
			}
		} else {
			if strings.Contains(val, "::") == true {
				subStr := strings.Split(val, "]")
				tStr := strings.TrimLeft(subStr[0], "[")
				hName, err := net.LookupAddr(tStr)
				if err == nil {
					hostName = strings.Join(hName, ".")
				}
			}
		}
	}
	hostName = strings.TrimRight(hostName, ".")
	replacer := strings.NewReplacer(".", "_", "/", "_", "\\", "_", ":", "_", "%", "_")
	hostName = replacer.Replace(hostName)
	return hostName
}

func checkIgnoreMetric(mkey string) bool {
	IgnoreChildMetric := "server"
	IgnoreMetric := ""
	ret := false
	if strings.EqualFold(IgnoreChildMetric, "nil") == false {
		subMetric := strings.Split(mkey, "/")
		if strings.Contains(IgnoreChildMetric, subMetric[len(subMetric)-1]) == true {
			ret = true
		}
	}
	if strings.EqualFold(IgnoreMetric, "nil") == false {
		if strings.Contains(IgnoreMetric, mkey) == true {
			ret = true
		}
	}

	return ret
}

func getNamespace(mkey string) (ns core.Namespace) {

	rc := strings.Replace(mkey, ".", "-", -1)
	ss := strings.Split(rc, "/")
	ns = core.NewNamespace(ss[0])
	for i := 1; i < len(ss); i++ {
		ns = ns.AddStaticElement(ss[i])
	}
	return ns
}

func switchType(outMetric *[]plugin.MetricType, mval interface{}, ak string) {

	switch mtype := mval.(type) {
	case bool:
		if checkIgnoreMetric(ak) == true {
			return
		}
		if mval == false {
			ns := getNamespace(ak)
			tmp := plugin.MetricType{
				Namespace_: ns,
				Data_:      0,
				Timestamp_: time.Now(),
			}
			*outMetric = append(*outMetric, tmp)
		} else {
			ns := getNamespace(ak)
			tmp := plugin.MetricType{
				Namespace_: ns,
				Data_:      1,
				Timestamp_: time.Now(),
			}
			*outMetric = append(*outMetric, tmp)
		}
	case int, int64, float64, string:
		if checkIgnoreMetric(ak) == true {
			return
		}
		ns := getNamespace(ak)
		tmp := plugin.MetricType{
			Namespace_: ns,
			Data_:      mval,
			Timestamp_: time.Now(),
		}
		*outMetric = append(*outMetric, tmp)
	case map[string]interface{}:
		parseMetrics(outMetric, mtype, ak)
	case []interface{}:
		parseArrMetrics(outMetric, mtype, ak)
	default:
		log.Println("In default missing type =", reflect.TypeOf(mval))
	}
	return
}

func parseArrMetrics(outMetric *[]plugin.MetricType, inData []interface{}, parentKey string) {
	for mkey, mval := range inData {
		subMetric := strings.Split(parentKey, "/")
		if subMetric[len(subMetric)-1] == "peers" {
			hostName := getHostName(mval, strconv.Itoa(mkey))
			switchType(outMetric, mval, parentKey+"/"+hostName)
		} else {
			switchType(outMetric, mval, parentKey+"/"+strconv.Itoa(mkey))
		}
	}
	return
}

func parseMetrics(outMetric *[]plugin.MetricType, inData map[string]interface{}, parentKey string) {

	for mkey, mval := range inData {
		switchType(outMetric, mval, parentKey+"/"+mkey)
	}
	return
}

func getMetrics(webserver string, metrics []string) (mts []plugin.MetricType, err error) {
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	resp, err1 := client.Get(webserver)
	if err1 != nil {
		return nil, err1
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		return nil, errReqFailed
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return nil, err2
	}

	jFmt := make(map[string]interface{})

	err = json.Unmarshal(body, &jFmt)
	if err != nil {
		return nil, err
	}

	pk := "staples" + "/" + "nginx"
	parseMetrics(&mts, jFmt, pk)

	return mts, nil
}

func (n *Nginx) CollectMetrics(inmts []plugin.MetricType) ([]plugin.MetricType, error) {
	webservercfg := inmts[0].Config().Table()["nginx_status_url"]

	webserver, ok := webservercfg.(ctypes.ConfigValueStr)
	if !ok {
		return nil, errBadWebserver
	}

	metrics := make([]string, len(inmts))

	for i, m := range inmts {
		metrics[i] = m.Namespace().String()
	}

	mts, err := getMetrics(webserver.Value, metrics)

	if err != nil {
		log.Println("Error in Get Metric =", err)
	}

        return mts, nil   
}

func (n *Nginx) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	webservercfg := cfg.Table()["nginx_status_url"]
	if webservercfg == nil {
		return nil, errConfigReadError
	}

	webserver, ok := webservercfg.(ctypes.ConfigValueStr)
	if !ok {
		return nil, errBadWebserver
	}

	mts, err := getMetrics(webserver.Value, []string{})

	if err != nil {
		log.Println("Error in Get Metric =", err)
	}

        return mts, nil 
}

func (n *Nginx) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cfg := cpolicy.New()
	nginxrule, _ := cpolicy.NewStringRule("nginx_status_url", false, "http://demo.nginx.com/status")
	policy := cpolicy.NewPolicyNode()
	policy.Add(nginxrule)
	cfg.Add([]string{"staples", "nginx"}, policy)
	return cfg, nil
}

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		Name,
		Version,
		Type,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.Unsecure(true),
		plugin.RoutingStrategy(plugin.DefaultRouting),
		plugin.CacheTTL(1100*time.Millisecond),
	)
}
