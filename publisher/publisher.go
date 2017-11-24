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

package publisher

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"

	"github.com/intelsdi-x/snap-plugin-publisher-kairosdb/kairos"
	"github.com/intelsdi-x/snap/core/serror"
)

const (
	name        = "kairos"
	version     = 3
	pluginType  = plugin.PublisherPluginType
	publishPath = "/api/v1/datapoints"
)

// Meta returns a plugin meta data
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		name,
		version,
		pluginType,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
	)
}

// New returns an instance of the KairosDB publisher
func New() *publisher {
	return &publisher{}
}

func (pub *publisher) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()

	r1, err := cpolicy.NewStringRule("host", true)
	if err != nil {
		fields := map[string]interface{}{"StringRule": "host"}
		return nil, serror.New(err, fields)
	}
	r1.Description = "KairosDB host"
	config.Add(r1)

	r2, err := cpolicy.NewIntegerRule("port", true)
	if err != nil {
		fields := map[string]interface{}{"StringRule": "port"}
		return nil, serror.New(err, fields)
	}
	r2.Description = "KairosDB port"
	config.Add(r2)

	r3, err := cpolicy.NewBoolRule("useDynamic", true)
	if err != nil {
		fields := map[string]interface{}{"StringRule": "useDynamic"}
		return nil, serror.New(err, fields)
	}
	r3.Description = "Use dynamic namespaces for metric and tags"
	config.Add(r3)

	cp.Add([]string{""}, config)
	return cp, nil
}

// Publish publishes metric data to Kairosdb
func (pub *publisher) Publish(contentType string, content []byte, config map[string]ctypes.ConfigValue) error {
	logger := getLogger(config)
	var metrics []plugin.MetricType
	useDynamic := config["useDynamic"].(ctypes.ConfigValueBool).Value

	// decode content to metrics type
	switch contentType {
	case plugin.SnapGOBContentType:
		dec := gob.NewDecoder(bytes.NewBuffer(content))
		if err := dec.Decode(&metrics); err != nil {
			logger.WithFields(log.Fields{
				"err": err,
			}).Error("decoding error")
			fields := map[string]interface{}{"Decode": "metrics"}
			return serror.New(err, fields)
		}
	default:
		logger.Errorf("unknown content type '%v'", contentType)
		fields := map[string]interface{}{"contentType": "unknown"}
		return serror.New(fmt.Errorf("Unknown content type '%s'", contentType), fields)
	}

	// translate metrics to KairosDB publishing format
	points := []kairos.DataPoint{}
	for _, metric := range metrics {
		isDynamic, indexes := metric.Namespace().IsDynamic()
		ns := metric.Namespace().Strings()
		tags := metric.Tags()

		// create KairosDB data point
		if useDynamic {
			tempTags := make(map[string]string)
			if isDynamic {
				for i, j := range indexes {
					// The second return value from IsDynamic(), in this case `indexes`, is the index of
					// the dynamic element in the unmodified namespace. However, here we're deleting
					// elements, which is problematic when the number of dynamic elements in a namespace is
					// greater than 1. Therefore, we subtract i (the loop iteration) from j
					// (the original index) to compensate.
					//
					// Remove "data" from the namespace and create a tag for it
					ns = append(ns[:j-i], ns[j-i+1:]...)
					tempTags[metric.Namespace()[j].Name] = string(metric.Namespace()[j].Value)
				}
			}
			for k, v := range tags {
				tempTags[k] = string(v)
			}
			tempTags["host"] = string(tags[core.STD_TAG_PLUGIN_RUNNING_ON])
			point := kairos.DataPoint{
				Name:      "/" + strings.Join(ns, "/"),
				Value:     metric.Data(),
				TimeStamp: metric.Timestamp().UnixNano() / int64(time.Millisecond),
				Tags:      tempTags,
			}
			points = append(points, point)
		} else {
			point := kairos.DataPoint{
				Name:      metric.Namespace().String(),
				Value:     metric.Data(),
				TimeStamp: metric.Timestamp().UnixNano() / int64(time.Millisecond),
				Tags:      metric.Tags(),
			}
			points = append(points, point)
		}
	}

	// serialization
	rendered, err := json.Marshal(points)
	if err != nil {
		logger.WithFields(log.Fields{
			"err": err,
		}).Error("Serialization error")
		fields := map[string]interface{}{"Marshal": "points"}
		return serror.New(err, fields)

	}

	// prepare publishing request
	u, err := url.Parse(
		fmt.Sprintf("http://%s:%d%s",
			config["host"].(ctypes.ConfigValueStr).Value,
			config["port"].(ctypes.ConfigValueInt).Value,
			publishPath,
		),
	)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(rendered))
	req.Header.Set("Content-Type", "application/json")

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.WithFields(log.Fields{
			"err": err,
		}).Error("Request error")
		fields := map[string]interface{}{"Request": "send"}
		return serror.New(err, fields)
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusNoContent {
		logger.WithFields(log.Fields{
			"err": err,
		}).Error("Response error ", resp.StatusCode)
		fields := map[string]interface{}{"Response": resp.Status}
		return serror.New(err, fields)
	}

	return nil
}

type publisher struct{}

func getLogger(config map[string]ctypes.ConfigValue) *log.Entry {
	logger := log.WithFields(log.Fields{
		"plugin-name":    name,
		"plugin-version": version,
		"plugin-type":    pluginType.String(),
	})

	// default
	log.SetLevel(log.WarnLevel)

	if debug, ok := config["debug"]; ok {
		switch v := debug.(type) {
		case ctypes.ConfigValueBool:
			if v.Value {
				log.SetLevel(log.DebugLevel)
				return logger
			}
		default:
			logger.WithFields(log.Fields{
				"field":         "debug",
				"type":          v,
				"expected type": "ctypes.ConfigValueBool",
			}).Error("invalid config type")
		}
	}

	if loglevel, ok := config["log-level"]; ok {
		switch v := loglevel.(type) {
		case ctypes.ConfigValueStr:
			switch strings.ToLower(v.Value) {
			case "warn":
				log.SetLevel(log.WarnLevel)
			case "error":
				log.SetLevel(log.ErrorLevel)
			case "debug":
				log.SetLevel(log.DebugLevel)
			case "info":
				log.SetLevel(log.InfoLevel)
			default:
				log.WithFields(log.Fields{
					"value":             strings.ToLower(v.Value),
					"acceptable values": "warn, error, debug, info",
				}).Warn("invalid config value")
			}
		default:
			logger.WithFields(log.Fields{
				"field":         "log-level",
				"type":          v,
				"expected type": "ctypes.ConfigValueStr",
			}).Error("invalid config type")
		}
	}

	return logger
}
