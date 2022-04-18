# def-generator

Create Go struct for vela X-definition 

for example:

1. annotation:
```cue
patch: {
		metadata: {
			annotations: {
				for k, v in parameter {
					"\(k)": v
				}
			}
		}
		if context.output.spec != _|_ && context.output.spec.template != _|_ {
			spec: template: metadata: annotations: {
				for k, v in parameter {
					"\(k)": v
				}
			}
		}
	}
parameter: [string]: string | null
```

generate result:

```go
package main
// Item -
type Item string

// Parameter -
type Parameter map[string]Item
```

2. websevice:


```cue
// webservice
import (
	"strconv"
)

// output...
// emmit for briefty

parameter: {
	// +usage=Specify the labels in the workload
	labels?: [string]: string

	// +usage=Specify the annotations in the workload
	annotations?: [string]: string

	// +usage=Which image would you like to use for your service
	// +short=i
	image: string

	// +usage=Specify image pull policy for your service
	imagePullPolicy?: "Always" | "Never" | "IfNotPresent"

	// +usage=Specify image pull secrets for your service
	imagePullSecrets?: [...string]

	// +ignore
	// +usage=Deprecated field, please use ports instead
	// +short=p
	port?: int

	// +usage=Which ports do you want customer traffic sent to, defaults to 80
	ports?: [...{
		// +usage=Number of port to expose on the pod's IP address
		port: int
		// +usage=Name of the port
		name?: string
		// +usage=Protocol for port. Must be UDP, TCP, or SCTP
		protocol: *"TCP" | "UDP" | "SCTP"
		// +usage=Specify if the port should be exposed
		expose: *false | bool
	}]

	// +ignore
	// +usage=Specify what kind of Service you want. options: "ClusterIP", "NodePort", "LoadBalancer", "ExternalName"
	exposeType: *"ClusterIP" | "NodePort" | "LoadBalancer" | "ExternalName"

	// +ignore
	// +usage=If addRevisionLabel is true, the appRevision label will be added to the underlying pods
	addRevisionLabel: *false | bool

	// +usage=Commands to run in the container
	cmd?: [...string]

	// +usage=Define arguments by using environment variables
	env?: [...{
		// +usage=Environment variable name
		name: string
		// +usage=The value of the environment variable
		value?: string
		// +usage=Specifies a source the value of this var should come from
		valueFrom?: {
			// +usage=Selects a key of a secret in the pod's namespace
			secretKeyRef?: {
				// +usage=The name of the secret in the pod's namespace to select from
				name: string
				// +usage=The key of the secret to select from. Must be a valid secret key
				key: string
			}
			// +usage=Selects a key of a config map in the pod's namespace
			configMapKeyRef?: {
				// +usage=The name of the config map in the pod's namespace to select from
				name: string
				// +usage=The key of the config map to select from. Must be a valid secret key
				key: string
			}
		}
	}]

	// +usage=Number of CPU units for the service
	cpu?: string

	// +usage=Specifies the attributes of the memory resource required for the container.
	memory?: string

	volumeMounts?: {
		// +usage=Mount PVC type volume
		pvc?: [...{
			name:      string
			mountPath: string
			// +usage=The name of the PVC
			claimName: string
		}]
		// +usage=Mount ConfigMap type volume
		configMap?: [...{
			name:        string
			mountPath:   string
			defaultMode: *420 | int
			cmName:      string
			items?: [...{
				key:  string
				path: string
				mode: *511 | int
			}]
		}]
		// +usage=Mount Secret type volume
		secret?: [...{
			name:        string
			mountPath:   string
			defaultMode: *420 | int
			secretName:  string
			items?: [...{
				key:  string
				path: string
				mode: *511 | int
			}]
		}]
		// +usage=Mount EmptyDir type volume
		emptyDir?: [...{
			name:      string
			mountPath: string
			medium:    *"" | "Memory"
		}]
		// +usage=Mount HostPath type volume
		hostPath?: [...{
			name:      string
			mountPath: string
			path:      string
		}]
	}

	// +usage=Deprecated field, use volumeMounts instead.
	volumes?: [...{
		name:      string
		mountPath: string
		// +usage=Specify volume type, options: "pvc","configMap","secret","emptyDir"
		type: "pvc" | "configMap" | "secret" | "emptyDir"
		if type == "pvc" {
			claimName: string
		}
		if type == "configMap" {
			defaultMode: *420 | int
			cmName:      string
			items?: [...{
				key:  string
				path: string
				mode: *511 | int
			}]
		}
		if type == "secret" {
			defaultMode: *420 | int
			secretName:  string
			items?: [...{
				key:  string
				path: string
				mode: *511 | int
			}]
		}
		if type == "emptyDir" {
			medium: *"" | "Memory"
		}
	}]

	// +usage=Instructions for assessing whether the container is alive.
	livenessProbe?: #HealthProbe

	// +usage=Instructions for assessing whether the container is in a suitable state to serve traffic.
	readinessProbe?: #HealthProbe

	// +usage=Specify the hostAliases to add
	hostAliases: [...{
		ip: string
		hostnames: [...string]
	}]
}

#HealthProbe: {

	// +usage=Instructions for assessing container health by executing a command. Either this attribute or the httpGet attribute or the tcpSocket attribute MUST be specified. This attribute is mutually exclusive with both the httpGet attribute and the tcpSocket attribute.
	exec?: {
		// +usage=A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.
		command: [...string]
	}

	// +usage=Instructions for assessing container health by executing an HTTP GET request. Either this attribute or the exec attribute or the tcpSocket attribute MUST be specified. This attribute is mutually exclusive with both the exec attribute and the tcpSocket attribute.
	httpGet?: {
		// +usage=The endpoint, relative to the port, to which the HTTP GET request should be directed.
		path: string
		// +usage=The TCP socket within the container to which the HTTP GET request should be directed.
		port: int
		httpHeaders?: [...{
			name:  string
			value: string
		}]
	}

	// +usage=Instructions for assessing container health by probing a TCP socket. Either this attribute or the exec attribute or the httpGet attribute MUST be specified. This attribute is mutually exclusive with both the exec attribute and the httpGet attribute.
	tcpSocket?: {
		// +usage=The TCP socket within the container that should be probed to assess container health.
		port: int
	}

	// +usage=Number of seconds after the container is started before the first probe is initiated.
	initialDelaySeconds: *0 | int

	// +usage=How often, in seconds, to execute the probe.
	periodSeconds: *10 | int

	// +usage=Number of seconds after which the probe times out.
	timeoutSeconds: *1 | int

	// +usage=Minimum consecutive successes for the probe to be considered successful after having failed.
	successThreshold: *1 | int

	// +usage=Number of consecutive failures required to determine the container is not alive (liveness probe) or not ready (readiness probe).
	failureThreshold: *3 | int
}
```

```go
package main
// SecretKeyRef Selects a key of a secret in the pod's namespace
type SecretKeyRef struct {
    Name string `json:"name"`
    Key string `json:"key"`
}

// ConfigMapKeyRef Selects a key of a config map in the pod's namespace
type ConfigMapKeyRef struct {
    Name string `json:"name"`
    Key string `json:"key"`
}

// ValueFrom Specifies a source the value of this var should come from
type ValueFrom struct {
    SecretKeyRef SecretKeyRef `json:"secretKeyRef"`
    ConfigMapKeyRef ConfigMapKeyRef `json:"configMapKeyRef"`
}

// Env -
type Env struct {
    Name string `json:"name"`
    Value string `json:"value"`
    ValueFrom ValueFrom `json:"valueFrom"`
}

// PVC -
type PVC struct {
    Name string `json:"name"`
    MountPath string `json:"mountPath"`
    ClaimName string `json:"claimName"`
}

// Items -
type Items struct {
    Path string `json:"path"`
    Key string `json:"key"`
    Mode int `json:"mode"`
}

// ConfigMap -
type ConfigMap struct {
    Name string `json:"name"`
    MountPath string `json:"mountPath"`
    DefaultMode int `json:"defaultMode"`
    CmName string `json:"cmName"`
    Items []Items `json:"items"`
}

// Items -
type Items struct {
    Path string `json:"path"`
    Key string `json:"key"`
    Mode int `json:"mode"`
}

// Secret -
type Secret struct {
    Name string `json:"name"`
    MountPath string `json:"mountPath"`
    DefaultMode int `json:"defaultMode"`
    Items []Items `json:"items"`
    SecretName string `json:"secretName"`
}

// EmptyDir -
type EmptyDir struct {
    Name string `json:"name"`
    MountPath string `json:"mountPath"`
    Medium string `json:"medium"`
}

// HostPath -
type HostPath struct {
    Path string `json:"path"`
    Name string `json:"name"`
    MountPath string `json:"mountPath"`
}

// VolumeMounts -
type VolumeMounts struct {
    PVC []PVC `json:"pvc"`
    ConfigMap []ConfigMap `json:"configMap"`
    Secret []Secret `json:"secret"`
    EmptyDir []EmptyDir `json:"emptyDir"`
    HostPath []HostPath `json:"hostPath"`
}

// Ports -
type Ports struct {
    Name string `json:"name"`
    Port int `json:"port"`
    Protocol string `json:"protocol"`
    Expose bool `json:"expose"`
}

// Volumes -
type Volumes struct {
    Name string `json:"name"`
    MountPath string `json:"mountPath"`
    Type string `json:"type"`
}

// Exec Instructions for assessing container health by executing a command. Either this attribute or the httpGet attribute or the tcpSocket attribute MUST be specified. This attribute is mutually exclusive with both the httpGet attribute and the tcpSocket attribute.
type Exec struct {
    Command []string `json:"command"`
}

// HTTPHeaders -
type HTTPHeaders struct {
    Name string `json:"name"`
    Value string `json:"value"`
}

// HTTPGet Instructions for assessing container health by executing an HTTP GET request. Either this attribute or the exec attribute or the tcpSocket attribute MUST be specified. This attribute is mutually exclusive with both the exec attribute and the tcpSocket attribute.
type HTTPGet struct {
    Path string `json:"path"`
    Port int `json:"port"`
    HTTPHeaders []HTTPHeaders `json:"httpHeaders"`
}

// TcpSocket Instructions for assessing container health by probing a TCP socket. Either this attribute or the exec attribute or the httpGet attribute MUST be specified. This attribute is mutually exclusive with both the exec attribute and the httpGet attribute.
type TcpSocket struct {
    Port int `json:"port"`
}

// LivenessProbe Instructions for assessing whether the container is alive.
type LivenessProbe struct {
    Exec Exec `json:"exec"`
    HTTPGet HTTPGet `json:"httpGet"`
    TcpSocket TcpSocket `json:"tcpSocket"`
    InitialDelaySeconds int `json:"initialDelaySeconds"`
    PeriodSeconds int `json:"periodSeconds"`
    TimeoutSeconds int `json:"timeoutSeconds"`
    SuccessThreshold int `json:"successThreshold"`
    FailureThreshold int `json:"failureThreshold"`
}

// Exec Instructions for assessing container health by executing a command. Either this attribute or the httpGet attribute or the tcpSocket attribute MUST be specified. This attribute is mutually exclusive with both the httpGet attribute and the tcpSocket attribute.
type Exec struct {
    Command []string `json:"command"`
}

// HTTPHeaders -
type HTTPHeaders struct {
    Name string `json:"name"`
    Value string `json:"value"`
}

// HTTPGet Instructions for assessing container health by executing an HTTP GET request. Either this attribute or the exec attribute or the tcpSocket attribute MUST be specified. This attribute is mutually exclusive with both the exec attribute and the tcpSocket attribute.
type HTTPGet struct {
    Path string `json:"path"`
    Port int `json:"port"`
    HTTPHeaders []HTTPHeaders `json:"httpHeaders"`
}

// TcpSocket Instructions for assessing container health by probing a TCP socket. Either this attribute or the exec attribute or the httpGet attribute MUST be specified. This attribute is mutually exclusive with both the exec attribute and the httpGet attribute.
type TcpSocket struct {
    Port int `json:"port"`
}

// ReadinessProbe Instructions for assessing whether the container is in a suitable state to serve traffic.
type ReadinessProbe struct {
    Exec Exec `json:"exec"`
    HTTPGet HTTPGet `json:"httpGet"`
    TcpSocket TcpSocket `json:"tcpSocket"`
    InitialDelaySeconds int `json:"initialDelaySeconds"`
    PeriodSeconds int `json:"periodSeconds"`
    TimeoutSeconds int `json:"timeoutSeconds"`
    SuccessThreshold int `json:"successThreshold"`
    FailureThreshold int `json:"failureThreshold"`
}

// HostAliases -
type HostAliases struct {
    Ip string `json:"ip"`
    Hostnames []string `json:"hostnames"`
}

// Parameter -
type Parameter struct {
    Cmd []string `json:"cmd"`
    Env []Env `json:"env"`
    VolumeMounts VolumeMounts `json:"volumeMounts"`
    Labels map[string]string `json:"labels"`
    AddRevisionLabel bool `json:"addRevisionLabel"`
    Annotations map[string]string `json:"annotations"`
    Image string `json:"image"`
    Ports []Ports `json:"ports"`
    Port int `json:"port"`
    ImagePullPolicy string `json:"imagePullPolicy"`
    CPU string `json:"cpu"`
    Memory string `json:"memory"`
    Volumes []Volumes `json:"volumes"`
    LivenessProbe LivenessProbe `json:"livenessProbe"`
    ReadinessProbe ReadinessProbe `json:"readinessProbe"`
    HostAliases []HostAliases `json:"hostAliases"`
    ImagePullSecrets []string `json:"imagePullSecrets"`
    ExposeType string `json:"exposeType"`
}


```


