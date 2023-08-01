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
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/palantir/stacktrace"
	"gopkg.in/yaml.v3"

	"alluxio.org/build/artifact"
)

const (
	defaultPresetsYmlFilePath = "src/alluxio.org/build/presets.yml"
)

type preset struct {
	Docker []string `yaml:"docker,omitempty"`
	Helm   []string `yaml:"helm,omitempty"`
}

func PresetsF(args []string) error {
	cmd := flag.NewFlagSet(Presets, flag.ExitOnError)
	// preset flags
	var flagArtifactOutput, flagPresetName, flagPresetsYmlFile string
	cmd.StringVar(&flagArtifactOutput, "artifact", "", "If set, writes object representing the tarball to YAML output file")
	cmd.StringVar(&flagPresetsYmlFile, "presetsYmlFile", defaultPresetsYmlFilePath, "Path to presets.yml file")
	cmd.StringVar(&flagPresetName, "name", "", "Choose the preset to build. See available presets in presets.yml")

	//parse flags
	if err := cmd.Parse(args); err != nil {
		return stacktrace.Propagate(err, "error parsing flags")
	}
	// parse presets.yml
	var presets map[string]preset
	{
		wd, err := os.Getwd()
		if err != nil {
			return stacktrace.Propagate(err, "error getting current working directory")
		}
		dockerYmlPath := filepath.Join(wd, flagPresetsYmlFile)
		content, err := ioutil.ReadFile(dockerYmlPath)
		if err != nil {
			return stacktrace.Propagate(err, "error reading file at %v", dockerYmlPath)
		}
		if err := yaml.Unmarshal(content, &presets); err != nil {
			return stacktrace.Propagate(err, "error unmarshalling presets from:\n%v", string(content))
		}
	}

	p, ok := presets[flagPresetName]
	if !ok {
		return stacktrace.NewError("error finding preset named %v", flagPresetName)
	}

	alluxioVersion, err := versionFromChartYaml()
	if err != nil {
		return stacktrace.Propagate(err, "error parsing version string")
	}

	a, err := artifact.NewArtifactGroup(alluxioVersion)
	if err != nil {
		return stacktrace.Propagate(err, "error creating artifact group")
	}
	for _, dArgs := range p.Docker {
		dOpts, err := newDockerBuildOpts(strings.Split(dArgs, " "))
		if err != nil {
			return stacktrace.Propagate(err, "error creating docker build opts")
		}
		image, ok := dOpts.dockerImages[dOpts.image]
		if !ok {
			return stacktrace.NewError("must provide valid 'image' arg")
		}
		a.Add(artifact.DockerArtifact, dOpts.outputDir, image.TargetName, dOpts.metadata)
	}
	for _, hArgs := range p.Helm {
		hOpts, err := newHelmBuildOpts(strings.Split(hArgs, " "))
		if err != nil {
			return stacktrace.Propagate(err, "error creating helm build opts")
		}
		image, ok := hOpts.helmCharts[hOpts.chartName]
		if !ok {
			return stacktrace.NewError("must provide valid 'chartName' arg")
		}
		a.Add(artifact.HelmArtifact, hOpts.outputDir, image.TargetName, nil)
	}

	if flagArtifactOutput != "" {
		return a.WriteToFile(flagArtifactOutput)
	}

	for _, dArgs := range p.Docker {
		if err := DockerF(strings.Split(dArgs, " ")); err != nil {
			return stacktrace.Propagate(err, "error building docker image")
		}
	}
	for _, hArgs := range p.Helm {
		if err := HelmF(strings.Split(hArgs, " ")); err != nil {
			return stacktrace.Propagate(err, "error building helm chart")
		}
	}

	return nil
}
