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
        flag := false
	switch mtype := inData.(type) {
	   case map[string]interface{}:
	        hostName = mtype["server"].(string)
		//check for IPV4
		if strings.Count(hostName, ".") == 3 {
			subStr := strings.Split(hostName, ":")
			hName, err := net.LookupAddr(subStr[0])
			if err == nil {
				hostName = strings.Join(hName, ".")
                                flag = true
			}
		} else {
			if strings.Contains(hostName, "::") == true {
				subStr := strings.Split(hostName, "]")
				tStr := strings.TrimLeft(subStr[0], "[")
				hName, err := net.LookupAddr(tStr)
				if err == nil {
					hostName = strings.Join(hName, ".")
                                        flag = true
				}
			}
		}
	}
        if flag == false {
	  //Default hostname with port will be encoded to md5
	  hostName = fmt.Sprintf("host_id_%s", getMD5Hash(hostName))
        } else {
	  hostName = fmt.Sprintf("host_id_%s", hostName)
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
	ns = core.NewNamespace(ss...)
	return ns
}

func switchType(outMetric *[]plugin.MetricType, mval interface{}, ak string) {
	switch mtype := mval.(type) {
	case bool:
		if checkIgnoreMetric(ak) == true {
			return
		}
                ns := getNamespace(ak)
		tmp := plugin.MetricType{}
		tmp.Namespace_= ns
                if mval.(bool) == false {
		   tmp.Data_= 0 
                } else {
		   tmp.Data_= 1
                }
		tmp.Timestamp_= time.Now()
		*outMetric = append(*outMetric, tmp)
	case int, int64, float64, string:
		if checkIgnoreMetric(ak) == true {
			return
		}
                ns := getNamespace(ak)
		tmp := plugin.MetricType{}
		tmp.Namespace_= ns
		tmp.Data_=      mval
		tmp.Timestamp_= time.Now()
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

func (n *Nginx) GetMetricTypes(cfg plugin.ConfigType) (mts []plugin.MetricType, err error) {
	webservercfg := cfg.Table()["nginx_status_url"]
	if webservercfg == nil {
		return nil, errConfigReadError
	}

	webserver, ok := webservercfg.(ctypes.ConfigValueStr)
	if !ok {
		return nil, errBadWebserver
	}

	mts, err = getMetrics(webserver.Value, []string{})

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
