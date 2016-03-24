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

package publisher

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
)

type KairosSuite struct {
	suite.Suite
	server *httptest.Server
}

func (s *KairosSuite) SetupSuite() {
	router := mux.NewRouter()
	s.server = httptest.NewServer(router)

	registerPublish(router)
}

func (s *KairosSuite) TearDownSuite() {
	s.server.Close()
}

func (s *KairosSuite) TestMeta() {
	Convey("Meta should return metadata for the plugin", s.T(), func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, name)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, plugin.PublisherPluginType)
	})
}

func (s *KairosSuite) TestGetConfigPolicy() {

	Convey("Create KairosPublisher", s.T(), func() {
		kp := New()

		Convey("So ip should not be nil", func() {
			So(kp, ShouldNotBeNil)
		})

		Convey("So ip should be of publisher type", func() {
			So(kp, ShouldHaveSameTypeAs, &publisher{})
		})

		configPolicy, err := kp.GetConfigPolicy()
		Convey("ip.GetConfigPolicy() should return a config policy", func() {

			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})

			Convey("So we should not get an err retreiving the config policy", func() {
				So(err, ShouldBeNil)
			})

			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})

			testConfig := make(map[string]ctypes.ConfigValue)
			testConfig["host"] = ctypes.ConfigValueStr{Value: "localhost"}
			testConfig["port"] = ctypes.ConfigValueInt{Value: 8083}

			cfg, errs := configPolicy.Get([]string{""}).Process(testConfig)

			Convey("So config policy should process testConfig and return a config", func() {
				So(cfg, ShouldNotBeNil)
			})

			Convey("So testConfig processing should return no errors", func() {
				So(errs.HasErrors(), ShouldBeFalse)
			})

			testConfig = make(map[string]ctypes.ConfigValue)
			testConfig["endpoint"] = ctypes.ConfigValueStr{Value: "http://localhost:8080"}
			cfg, errs = configPolicy.Get([]string{""}).Process(testConfig)

			Convey("So config policy should not return a config after processing invalid testConfig", func() {
				So(cfg, ShouldBeNil)
			})

			Convey("So testConfig processing should return errors", func() {
				So(errs.HasErrors(), ShouldBeTrue)
			})
		})
	})
}

func (s *KairosSuite) TestPublish() {
	Convey("Given snap KairosDB publisher testing", s.T(), func() {
		var buf bytes.Buffer
		config := make(map[string]ctypes.ConfigValue)
		u, _ := url.Parse(s.server.URL)
		hostPort := strings.Split(u.Host, ":")
		config["host"] = ctypes.ConfigValueStr{Value: hostPort[0]}
		port, _ := strconv.Atoi(hostPort[1])
		config["port"] = ctypes.ConfigValueInt{Value: port}

		kp := New()
		cp, _ := kp.GetConfigPolicy()
		cfg, _ := cp.Get([]string{""}).Process(config)

		Convey("Publish provided metrics", func() {
			metrics := []plugin.PluginMetricType{
				plugin.PluginMetricType{
					Namespace_: []string{"foo"},
					Timestamp_: time.Now(),
					Source_:    "127.0.0.1",
					Data_:      43,
				},
				plugin.PluginMetricType{
					Namespace_: []string{"bar"},
					Timestamp_: time.Now(),
					Source_:    "127.0.0.1",
					Data_:      44,
				},
			}
			buf.Reset()
			enc := gob.NewEncoder(&buf)
			enc.Encode(metrics)
			err := kp.Publish(plugin.SnapGOBContentType, buf.Bytes(), *cfg)
			So(err, ShouldBeNil)
		})

	})
}

func TestKairosSuite(t *testing.T) {
	kairos := new(KairosSuite)
	suite.Run(t, kairos)
}

func registerPublish(r *mux.Router) {
	r.HandleFunc("/api/v1/datapoints", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "")
	}).Methods("POST")
}
