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

package v1alpha1

import (
	"go.uber.org/config"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type AlluxioClusterHelmChartSpec struct {
	NameOverride       *string             `json:"nameOverride,omitempty" yaml:"nameOverride,omitempty"`
	Dataset            *DatasetConf        `json:"dataset" yaml:"dataset"`
	DatasetName        *string             `json:"datasetName" yaml:"datasetName"`
	Image              *string             `json:"image,omitempty" yaml:"image,omitempty"`
	ImageTag           *string             `json:"imageTag,omitempty" yaml:"imageTag,omitempty"`
	ImagePullPolicy    *string             `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	ImagePullSecrets   []string            `json:"imagePullSecrets,omitempty" yaml:"imagePullSecrets,omitempty"`
	User               *int                `json:"user,omitempty" yaml:"user,omitempty"`
	Group              *int                `json:"group,omitempty" yaml:"group,omitempty"`
	FsGroup            *string             `json:"fsGroup,omitempty" yaml:"fsGroup,omitempty"`
	HostNetwork        *bool               `json:"hostNetwork,omitempty" yaml:"hostNetwork,omitempty"`
	DnsPolicy          *string             `json:"dnsPolicy,omitempty" yaml:"dnsPolicy,omitempty"`
	ServiceAccountName *string             `json:"serviceAccountName,omitempty" yaml:"serviceAccountName,omitempty"`
	HostAliases        []*HostAlias        `json:"hostAliases,omitempty" yaml:"hostAliases,omitempty"`
	HostPathForLogging *bool               `json:"hostPathForLogging,omitempty" yaml:"hostPathForLogging,omitempty"`
	NodeSelector       map[string]string   `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	Tolerations        []*Toleration       `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	Properties         map[string]string   `json:"properties,omitempty" yaml:"properties,omitempty"`
	JvmOptions         []string            `json:"jvmOptions,omitempty" yaml:"jvmOptions,omitempty"`
	PvcMounts          *MountSpec          `json:"pvcMounts,omitempty" yaml:"pvcMounts,omitempty"`
	ConfigMaps         *MountSpec          `json:"configMaps,omitempty" yaml:"configMaps,omitempty"`
	Secrets            *MountSpec          `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	Master             *MasterSpec         `json:"master,omitempty" yaml:"master,omitempty"`
	Journal            *JournalSpec        `json:"journal,omitempty" yaml:"journal,omitempty"`
	Worker             *WorkerSpec         `json:"worker,omitempty" yaml:"worker,omitempty"`
	Pagestore          *PagestoreSpec      `json:"pagestore,omitempty" yaml:"pagestore,omitempty"`
	Metastore          *MetastoreSpec      `json:"metastore,omitempty" yaml:"metastore,omitempty"`
	Proxy              *ProxySpec          `json:"proxy,omitempty" yaml:"proxy,omitempty"`
	Fuse               *FuseSpec           `json:"fuse,omitempty" yaml:"fuse,omitempty"`
	Csi                *CsiSpec            `json:"csi,omitempty" yaml:"csi,omitempty"`
	Metrics            *MetricsSpec        `json:"metrics,omitempty" yaml:"metrics,omitempty"`
	AlluxioMonitor     *AlluxioMonitorSpec `json:"alluxio-monitor,omitempty" yaml:"alluxio-monitor,omitempty"`
	Etcd               *EtcdSpec           `json:"etcd,omitempty" yaml:"etcd,omitempty"`
}

type HostAlias struct {
	Ip        *string  `json:"ip" yaml:"ip"`
	Hostnames []string `json:"hostnames" yaml:"hostnames"`
}

type Toleration struct {
	Key      *string `json:"key" yaml:"key"`
	Operator *string `json:"operator" yaml:"operator"`
	Value    *string `json:"value" yaml:"value"`
	Effect   *string `json:"effect" yaml:"effect"`
}

type MountSpec struct {
	Master map[string]string `json:"master,omitempty" yaml:"master,omitempty"`
	Worker map[string]string `json:"worker,omitempty" yaml:"worker,omitempty"`
	Fuse   map[string]string `json:"fuse,omitempty" yaml:"fuse,omitempty"`
	Proxy  map[string]string `json:"proxy,omitempty" yaml:"proxy,omitempty"`
}

type MasterSpec struct {
	Affinity        *corev1.Affinity  `json:"affinity,omitempty" yaml:"affinity,omitempty"`
	Count           *int              `json:"count,omitempty" yaml:"count,omitempty"`
	Enabled         *bool             `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Env             map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	HostPathForLogs *string           `json:"hostPathForLogs,omitempty" yaml:"hostPathForLogs,omitempty"`
	JvmOptions      []string          `json:"jvmOptions,omitempty" yaml:"jvmOptions,omitempty"`
	LivenessProbe   *ProbeSpec        `json:"livenessProbe,omitempty" yaml:"livenessProbe,omitempty"`
	NodeSelector    map[string]string `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	PodAnnotations  map[string]string `json:"podAnnotations,omitempty" yaml:"podAnnotations,omitempty"`
	Ports           map[string]int    `json:"ports,omitempty" yaml:"ports,omitempty"`
	ReadinessProbe  *ProbeSpec        `json:"readinessProbe,omitempty" yaml:"readinessProbe,omitempty"`
	Resources       *ResourcesSpec    `json:"resources,omitempty" yaml:"resources,omitempty"`
	StartupProbe    *ProbeSpec        `json:"startupProbe,omitempty" yaml:"startupProbe,omitempty"`
	Tolerations     []*Toleration     `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
}

type JournalSpec struct {
	HostPath     *string `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`
	RunFormat    *bool   `json:"runFormat,omitempty" yaml:"runFormat,omitempty"`
	Size         *string `json:"size,omitempty" yaml:"size,omitempty"`
	StorageClass *string `json:"storageClass,omitempty" yaml:"storageClass,omitempty"`
	Type         *string `json:"type,omitempty" yaml:"type,omitempty"`
}

type WorkerSpec struct {
	Affinity              *corev1.Affinity  `json:"affinity,omitempty" yaml:"affinity,omitempty"`
	Count                 int               `json:"count,omitempty" yaml:"count,omitempty"`
	Env                   map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	HostPathForLogs       *string           `json:"hostPathForLogs,omitempty" yaml:"hostPathForLogs,omitempty"`
	JvmOptions            []string          `json:"jvmOptions,omitempty" yaml:"jvmOptions,omitempty"`
	LimitOneWorkerPerNode *bool             `json:"limitOneWorkerPerNode,omitempty" yaml:"limitOneWorkerPerNode,omitempty"`
	LivenessProbe         *ProbeSpec        `json:"livenessProbe,omitempty" yaml:"livenessProbe,omitempty"`
	NodeSelector          map[string]string `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	PodAnnotations        map[string]string `json:"podAnnotations,omitempty" yaml:"podAnnotations,omitempty"`
	Ports                 map[string]int    `json:"ports,omitempty" yaml:"ports,omitempty"`
	ReadinessProbe        *ProbeSpec        `json:"readinessProbe,omitempty" yaml:"readinessProbe,omitempty"`
	Resources             *ResourcesSpec    `json:"resources,omitempty" yaml:"resources,omitempty"`
	StartupProbe          *ProbeSpec        `json:"startupProbe,omitempty" yaml:"startupProbe,omitempty"`
	Tolerations           []*Toleration     `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
}

type PagestoreSpec struct {
	HostPath     *string `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`
	MemoryBacked *bool   `json:"memoryBacked,omitempty" yaml:"memoryBacked,omitempty"`
	Quota        *string `json:"quota,omitempty" yaml:"quota,omitempty"`
	StorageClass *string `json:"storageClass,omitempty" yaml:"storageClass,omitempty"`
	Type         *string `json:"type,omitempty" yaml:"type,omitempty"`
}

type MetastoreSpec struct {
	HostPath     *string `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`
	Size         *string `json:"size,omitempty" yaml:"size,omitempty"`
	StorageClass *string `json:"storageClass,omitempty" yaml:"storageClass,omitempty"`
	Type         *string `json:"type,omitempty" yaml:"type,omitempty"`
}

type ProxySpec struct {
	Affinity        *corev1.Affinity  `json:"affinity,omitempty" yaml:"affinity,omitempty"`
	Enabled         *bool             `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Env             map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	HostPathForLogs *string           `json:"hostPathForLogs,omitempty" yaml:"hostPathForLogs,omitempty"`
	JvmOptions      []string          `json:"jvmOptions,omitempty" yaml:"jvmOptions,omitempty"`
	NodeSelector    map[string]string `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	PodAnnotations  map[string]string `json:"podAnnotations,omitempty" yaml:"podAnnotations,omitempty"`
	Ports           map[string]int    `json:"ports,omitempty" yaml:"ports,omitempty"`
	Resources       *ResourcesSpec    `json:"resources,omitempty" yaml:"resources,omitempty"`
	Tolerations     []*Toleration     `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
}

type FuseSpec struct {
	Affinity        *corev1.Affinity  `json:"affinity,omitempty" yaml:"affinity,omitempty"`
	Enabled         *bool             `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Env             map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	Group           *int              `json:"group,omitempty" yaml:"group,omitempty"`
	HostPathForLogs *string           `json:"hostPathForLogs,omitempty" yaml:"hostPathForLogs,omitempty"`
	JvmOptions      []string          `json:"jvmOptions,omitempty" yaml:"jvmOptions,omitempty"`
	MountOptions    []string          `json:"mountOptions,omitempty" yaml:"mountOptions,omitempty"`
	NodeSelector    map[string]string `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	PodAnnotations  map[string]string `json:"podAnnotations,omitempty" yaml:"podAnnotations,omitempty"`
	Resources       *ResourcesSpec    `json:"resources,omitempty" yaml:"resources,omitempty"`
	Tolerations     []*Toleration     `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	User            *int              `json:"user,omitempty" yaml:"user,omitempty"`
}

type CsiSpec struct {
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

type ResourcesSpec struct {
	Limits   *CpuMemSpec `json:"limits,omitempty" yaml:"limits,omitempty"`
	Requests *CpuMemSpec `json:"requests,omitempty" yaml:"requests,omitempty"`
}

type CpuMemSpec struct {
	Cpu    *string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory *string `json:"memory,omitempty" yaml:"memory,omitempty"`
}

type ProbeSpec struct {
	FailureThreshold    *int `json:"failureThreshold,omitempty" yaml:"failureThreshold,omitempty"`
	InitialDelaySeconds *int `json:"initialDelaySeconds,omitempty" yaml:"initialDelaySeconds,omitempty"`
	PeriodSeconds       *int `json:"periodSeconds,omitempty" yaml:"periodSeconds,omitempty"`
	SuccessThreshold    *int `json:"successThreshold,omitempty" yaml:"successThreshold,omitempty"`
	TimeoutSeconds      *int `json:"timeoutSeconds,omitempty" yaml:"timeoutSeconds,omitempty"`
}

type MetricsSpec struct {
	ConsoleSink              *ConsoleSinkSpec              `json:"consoleSink,omitempty" yaml:"consoleSink,omitempty"`
	CsvSink                  *CsvSinkSpec                  `json:"csvSink,omitempty" yaml:"csvSink,omitempty"`
	GraphiteSink             *GraphiteSinkSpec             `json:"graphiteSink,omitempty" yaml:"graphiteSink,omitempty"`
	JmxSink                  *JmxSinkSpec                  `json:"jmxSink,omitempty" yaml:"jmxSink,omitempty"`
	PrometheusMetricsServlet *PrometheusMetricsServletSpec `json:"prometheusMetricsServlet,omitempty" yaml:"prometheusMetricsServlet,omitempty"`
	Slf4jSink                *Slf4jSinkSpec                `json:"slf4jSink,omitempty" yaml:"slf4jSink,omitempty"`
}

type ConsoleSinkSpec struct {
	Enabled *bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Period  *int    `json:"period,omitempty" yaml:"period,omitempty"`
	Unit    *string `json:"unit,omitempty" yaml:"unit,omitempty"`
}

type CsvSinkSpec struct {
	Directory *string `json:"directory,omitempty" yaml:"directory,omitempty"`
	Enabled   *bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Period    *int    `json:"period,omitempty" yaml:"period,omitempty"`
	Unit      *string `json:"unit,omitempty" yaml:"unit,omitempty"`
}

type GraphiteSinkSpec struct {
	Enabled  *bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Hostname *string `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Period   *int    `json:"period,omitempty" yaml:"period,omitempty"`
	Port     *int    `json:"port,omitempty" yaml:"port,omitempty"`
	Prefix   *string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Unit     *string `json:"unit,omitempty" yaml:"unit,omitempty"`
}

type JmxSinkSpec struct {
	Enabled *bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Domain  *string `json:"domain,omitempty" yaml:"domain,omitempty"`
}

type PrometheusMetricsServletSpec struct {
	Enabled        *bool             `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	PodAnnotations map[string]string `json:"podAnnotations,omitempty" yaml:"podAnnotations,omitempty"`
}

type Slf4jSinkSpec struct {
	Enabled     *bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	FilterClass *string `json:"filterClass,omitempty" yaml:"filterClass,omitempty"`
	FilterRegex *string `json:"filterRegex,omitempty" yaml:"filterRegex,omitempty"`
	Period      *int    `json:"period,omitempty" yaml:"period,omitempty"`
	Unit        *string `json:"unit,omitempty" yaml:"unit,omitempty"`
}

type AlluxioMonitorSpec struct {
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

type EtcdSpec struct {
	Enabled *bool `json:"enabled" yaml:"enabled"`

	Auth         *EtcdAuthSpec     `json:"auth,omitempty" yaml:"auth,omitempty"`
	Image        *EtcdImageSpec    `json:"image,omitempty" yaml:"image,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	ReplicaCount *int              `json:"replicaCount,omitempty" yaml:"replicaCount,omitempty"`
	Resources    ResourcesSpec     `json:"resources,omitempty" yaml:"resources,omitempty"`
}

type EtcdAuthSpec struct {
	Client *EtcdAuthClientSpec `json:"client,omitempty" yaml:"client,omitempty"`
}

type EtcdAuthClientSpec struct {
	EnableAuthentication *bool `json:"enableAuthentication,omitempty" yaml:"enableAuthentication,omitempty"`
}

type EtcdImageSpec struct {
	Registry   *string `json:"registry,omitempty" yaml:"registry,omitempty"`
	Repository *string `json:"repository,omitempty" yaml:"repository,omitempty"`
	Tag        *string `json:"tag,omitempty" yaml:"tag,omitempty"`
}

func (a *AlluxioClusterHelmChartSpec) YAMLToBytes(configYaml *config.YAML) ([]byte, error) {
	if err := configYaml.Get(config.Root).Populate(a); err != nil {
		return nil, err
	}
	return yaml.Marshal(a)
}
