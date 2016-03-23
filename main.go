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
	"os"

	"github.com/intelsdi-x/snap/control/plugin"

	"bytes"
	"encoding/gob"
	"github.com/intelsdi-x/snap-plugin-publisher-kairosdb/publisher"
	"github.com/intelsdi-x/snap/core/ctypes"
	"time"
	"fmt"
)

func main() {

	var buf bytes.Buffer
	config := make(map[string]ctypes.ConfigValue)
	config["host"] = ctypes.ConfigValueStr{Value: "10.91.97.189"}
	config["port"] = ctypes.ConfigValueInt{Value: 8083}

	kp := publisher.New()
	cp, _ := kp.GetConfigPolicy()
	cfg, _ := cp.Get([]string{""}).Process(config)
	metrics := []plugin.PluginMetricType{
		plugin.PluginMetricType{
			Namespace_: []string{"kromar"},
			Timestamp_: time.Now(),
			Source_:    "127.0.0.1",
			Data_:      47,
		},
		plugin.PluginMetricType{
			Namespace_: []string{"intel"},
			Timestamp_: time.Now(),
			Source_:    "127.0.0.1",
			Data_:      47,
		},
	}
	buf.Reset()
	enc := gob.NewEncoder(&buf)
	enc.Encode(metrics)
	err := kp.Publish(plugin.SnapGOBContentType, buf.Bytes(), *cfg)
	if err != nil {
		fmt.Println("Mam err ", err)
	}
	panic("dupa")
	meta := publisher.Meta()
	plugin.Start(meta, publisher.New(), os.Args[1])
}
