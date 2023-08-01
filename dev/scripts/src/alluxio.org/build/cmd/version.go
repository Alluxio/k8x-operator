/*
 * The Alluxio Open Foundation licenses this work under the Apache License, version 2.0
 * (the "License"). You may not use this work except in compliance with the License, which is
 * available at www.apache.org/licenses/LICENSE-2.0
 *
 * This software is distributed on an "AS IS" basis, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied, as more fully set forth in the License.
 *
 * See the NOTICE file distributed with this work for information regarding copyright ownership.
 */

package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/palantir/stacktrace"
	"gopkg.in/yaml.v3"

	"alluxio.org/common/repo"
)

const (
	// using a specific Chart.yaml since all charts and docker images in this repo will be same version
	alluxioChartYaml = "deploy/charts/alluxio/Chart.yaml"
)

func VersionF() error {
	ver, err := versionFromChartYaml()
	if err != nil {
		return stacktrace.Propagate(err, "error finding version from %v", alluxioChartYaml)
	}
	fmt.Println()
	fmt.Printf("Alluxio version from %v: %v\n", alluxioChartYaml, ver)
	return nil
}

func versionFromChartYaml() (string, error) {
	chartYaml := filepath.Join(repo.FindRepoRoot(), alluxioChartYaml)
	contents, err := ioutil.ReadFile(chartYaml)
	if err != nil {
		return "", stacktrace.Propagate(err, "error reading %v", chartYaml)
	}
	type yamlObj struct {
		Version string `yaml:"version"`
	}
	var chart yamlObj
	if err := yaml.Unmarshal(contents, &chart); err != nil {
		return "", stacktrace.Propagate(err, "error unmarshalling %v", alluxioChartYaml)
	}
	return chart.Version, nil
}
