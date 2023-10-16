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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/palantir/stacktrace"
	"gopkg.in/yaml.v3"

	"alluxio.org/build/artifact"
	"alluxio.org/common/command"
	"alluxio.org/common/repo"
)

const (
	defaultHelmYmlFilePath = "src/alluxio.org/build/helm.yml"
)

type HelmChart struct {
	ChartPath  string `yaml:"chartPath"`
	Flags      string `yaml:"flags"`
	TargetName string `yaml:"targetName,omitempty"`

	outputTarball string `yaml:"-"`
}

type helmBuildOpts struct {
	artifactOutput string
	chartName      string
	helmCharts     map[string]*HelmChart
	helmYmlFile    string
	outputDir      string
}

func newHelmBuildOpts(args []string) (*helmBuildOpts, error) {
	cmd := flag.NewFlagSet(Helm, flag.ExitOnError)
	opts := &helmBuildOpts{}
	// helm flags
	cmd.StringVar(&opts.artifactOutput, "artifact", "", "If set, writes object representing the tarball to YAML output file")
	cmd.StringVar(&opts.outputDir, "outputDir", repo.FindRepoRoot(), "Set output dir for generated tarball")
	cmd.StringVar(&opts.chartName, "chartName", "", "Choose the helm chart to build. See available charts in helm.yml")
	cmd.StringVar(&opts.helmYmlFile, "helmYmlFile", defaultHelmYmlFilePath, "Path to helm.yml file")

	// parse flags
	if err := cmd.Parse(args); err != nil {
		return nil, stacktrace.Propagate(err, "error parsing flags")
	}

	alluxioVersion, err := versionFromChartYaml()
	if err != nil {
		return nil, stacktrace.Propagate(err, "error parsing version string")
	}

	// parse available helm charts in helm.yml
	{
		wd, err := os.Getwd()
		if err != nil {
			return nil, stacktrace.Propagate(err, "error getting current working directory")
		}
		helmYmlPath := filepath.Join(wd, opts.helmYmlFile)
		content, err := ioutil.ReadFile(helmYmlPath)
		if err != nil {
			return nil, stacktrace.Propagate(err, "error reading file at %v", helmYmlPath)
		}
		if err := yaml.Unmarshal(content, &opts.helmCharts); err != nil {
			return nil, stacktrace.Propagate(err, "error unmarshalling helm charts from:\n%v", string(content))
		}
		for _, chart := range opts.helmCharts {
			chart.init(alluxioVersion)
		}
	}

	return opts, nil
}

func HelmF(args []string) error {
	opts, err := newHelmBuildOpts(args)
	if err != nil {
		return stacktrace.Propagate(err, "error creating helm build opts")
	}

	alluxioVersion, err := versionFromChartYaml()
	if err != nil {
		return stacktrace.Propagate(err, "error parsing version string")
	}

	chart, ok := opts.helmCharts[opts.chartName]
	if !ok {
		return stacktrace.NewError("must provide valid 'chartName' arg")
	}
	if opts.artifactOutput != "" {
		a, err := artifact.NewArtifactGroup(alluxioVersion)
		if err != nil {
			return stacktrace.Propagate(err, "error creating artifact group")
		}
		a.Add(artifact.HelmArtifact,
			opts.outputDir,
			chart.TargetName,
			map[string]string{
				artifact.HelmChartName: opts.chartName,
			},
		)
		return a.WriteToFile(opts.artifactOutput)
	}

	if err := chart.build(opts); err != nil {
		return stacktrace.Propagate(err, "error building helm chart %v", opts.chartName)
	}

	return nil
}

func (h *HelmChart) init(alluxioVersion string) {
	h.TargetName = strings.ReplaceAll(h.TargetName, versionPlaceholder, alluxioVersion)
}

func (h *HelmChart) build(opts *helmBuildOpts) error {
	h.outputTarball = fmt.Sprintf("%v/%v", opts.outputDir, h.TargetName)
	cmds := []string{
		fmt.Sprintf("helm package %v --destination %v %v", h.ChartPath, opts.outputDir, h.Flags),
	}
	for _, c := range cmds {
		log.Printf("Running: %v", c)
		if out, err := command.New(c).WithDir(repo.FindRepoRoot()).CombinedOutput(); err != nil {
			return stacktrace.Propagate(err, "error from running cmd: %v", string(out))
		}
	}
	if _, err := os.Stat(h.outputTarball); os.IsNotExist(err) {
		return stacktrace.Propagate(err, "expected tarball to be created at %v", h.outputTarball)
	}
	return nil
}
