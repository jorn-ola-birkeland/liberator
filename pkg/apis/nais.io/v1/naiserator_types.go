package nais_io_v1

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	DeploymentCorrelationIDAnnotation = "nais.io/deploymentCorrelationID"
	SkipDeploymentMessageAnnotation   = "nais.io/skipDeploymentMessage"
	DefaultVaultMountPath             = "/var/run/secrets/nais.io/vault"
)

type Azure struct {
	// Configures an Azure AD client for this application.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/"
	Application *AzureApplication `json:"application"`
	// Sidecar configures a sidecar that intercepts every HTTP request, and performs the OIDC flow if necessary.
	// All requests to ingress + `/oauth2` will be processed only by the sidecar, whereas all other requests
	// will be proxied to the application.
	//
	// If the client is authenticated with Azure AD, the `Authorization` header will be set to `Bearer <JWT>`.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/sidecar/"
	Sidecar *AzureSidecar `json:"sidecar,omitempty"`
}

type AzureNaisJob struct {
	// Configures an Azure AD client for this application.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/"
	Application *AzureApplication `json:"application"`
}

type AzureApplication struct {
	// Whether to enable provisioning of an Azure AD application.
	// If enabled, an Azure AD application will be provisioned.
	Enabled bool `json:"enabled"`
	// ReplyURLs is a list of allowed redirect URLs used when performing OpenID Connect flows for authenticating end-users.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/configuration#reply-urls"
	ReplyURLs []AzureAdReplyUrlString `json:"replyURLs,omitempty"`
	// A Tenant represents an organization in Azure AD.
	//
	// If unspecified, will default to `trygdeetaten.no` for development clusters and `nav.no` for production clusters.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/concepts#tenants"
	// +kubebuilder:validation:Enum=nav.no;trygdeetaten.no
	Tenant string         `json:"tenant,omitempty"`
	Claims *AzureAdClaims `json:"claims,omitempty"`
	// SinglePageApplication denotes whether this Azure AD application should be registered as a single-page-application for usage in client-side applications without access to secrets.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/configuration#single-page-application"
	// +nais:doc:Default="false"
	SinglePageApplication *bool `json:"singlePageApplication,omitempty"`
	// AllowAllUsers denotes whether all users within the tenant should be allowed to access this AzureAdApplication.
	// +nais:doc:Default="false"
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/access-policy#users"
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/access-policy#groups"
	AllowAllUsers *bool `json:"allowAllUsers,omitempty"`
}

// +kubebuilder:validation:Pattern=`^https?:\/\/.+$`
type AzureAdReplyUrlString string

type AzureSidecar struct {
	Wonderwall `json:",inline"`
}

type OpenSearch struct {
	// Configure your application to access your OpenSearch instance.
	// Use the `instance_name` that you specified in the [navikt/aiven-iac](https://github.com/navikt/aiven-iac) repository.
	Instance string `json:"instance"`
	// Access level for OpenSearch user
	// +kubebuilder:validation:Enum=read;write;readwrite;admin
	Access string `json:"access,omitempty"`
}

type Redis struct {
	// The last part of the name used when creating the instance (ie. redis-<team>-<instance>)
	Instance string `json:"instance,omitempty"`
	// Access level for redis user
	// +kubebuilder:validation:Enum=read;write;readwrite;admin
	Access string `json:"access,omitempty"`
}

type Influx struct {
	// Provisions an InfluxDB instance and configures your application to access it.
	// Use the prefix: `influx-` + `team` that you specified in the [navikt/aiven-iac](https://github.com/navikt/aiven-iac) repository.
	Instance string `json:"instance"`
}

// +kubebuilder:validation:Pattern=`^https:\/\/.+$`
type Ingress string

type IDPorten struct {
	// Whether to enable provisioning of an ID-porten client.
	// If enabled, an ID-porten client be provisioned.
	// +nais:doc:Availability="team namespaces"
	Enabled bool `json:"enabled"`
	// AccessTokenLifetime is the lifetime in seconds for any issued access token from ID-porten.
	//
	// If unspecified, defaults to `3600` seconds (1 hour).
	// +nais:doc:Default="3600"
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3600
	AccessTokenLifetime *int `json:"accessTokenLifetime,omitempty"`
	// ClientURI is the URL shown to the user at ID-porten when displaying a 'back' button or on errors.
	ClientURI IDPortenURI `json:"clientURI,omitempty"`
	// FrontchannelLogoutPath is a valid path for your application where ID-porten sends a request to whenever the user has
	// initiated a logout elsewhere as part of a single logout (front channel logout) process.
	//
	// +nais:doc:Link="https://doc.nais.io/security/auth/idporten/#front-channel-logout";"https://docs.digdir.no/oidc_func_sso.html#2-h%C3%A5ndtere-utlogging-fra-id-porten"
	// +nais:doc:Default="/oauth2/logout"
	// +kubebuilder:validation:Pattern=`^\/.*$`
	FrontchannelLogoutPath string `json:"frontchannelLogoutPath,omitempty"`
	// IntegrationType is used to make sensible choices for your client.
	// Which type of integration you choose will provide guidance on which scopes you can use with the client.
	// A client can only have one integration type.
	//
	// NB! It is not possible to change the integration type after creation.
	//
	// +nais:doc:Immutable=true
	// +nais:doc:Default=idporten
	// +nais:doc:Link="https://docs.digdir.no/oidc_protocol_scope.html#scope-limitations"
	// +nais:doc:Link="https://docs.digdir.no/oidc_func_clientreg.html"
	// +kubebuilder:validation:Enum=krr;idporten;api_klient
	IntegrationType string `json:"integrationType,omitempty" nais:"immutable"`
	// PostLogoutRedirectURIs are valid URIs that ID-porten will allow redirecting the end-user to after a single logout
	// has been initiated and performed by the application.
	//
	// +nais:doc:Default="https://www.nav.no"
	// +nais:doc:Link="https://doc.nais.io/security/auth/idporten/#self-initiated-logout";"https://docs.digdir.no/oidc_func_sso.html#1-utlogging-fra-egen-tjeneste"
	PostLogoutRedirectURIs []IDPortenURI `json:"postLogoutRedirectURIs,omitempty"`
	// RedirectPath is a valid path that ID-porten redirects back to after a successful authorization request.
	//
	// +nais:doc:Default="/oauth2/callback"
	// +kubebuilder:validation:Pattern=`^\/.*$`
	RedirectPath string `json:"redirectPath,omitempty"`
	// Register different oauth2 Scopes on your client.
	// You will not be able to add a scope to your client that conflicts with the client's IntegrationType.
	// For example, you can not add a scope that is limited to the IntegrationType `krr` of IntegrationType `idporten`, and vice versa.
	//
	// Default for IntegrationType `krr` = ("krr:global/kontaktinformasjon.read", "krr:global/digitalpost.read")
	// Default for IntegrationType `idporten` = ("openid", "profile")
	// IntegrationType `api_klient` have no Default, checkout Digdir documentation.
	//
	// +nais:doc:Link="https://docs.digdir.no/oidc_func_clientreg.html?h=api_klient#scopes"
	Scopes []string `json:"scopes,omitempty"`
	// SessionLifetime is the maximum lifetime in seconds for any given user's session in your application.
	// The timeout starts whenever the user is redirected from the `authorization_endpoint` at ID-porten.
	//
	// If unspecified, defaults to `7200` seconds (2 hours).
	// Note: Attempting to refresh the user's `access_token` beyond this timeout will yield an error.
	//
	// +nais:doc:Default="7200"
	// +kubebuilder:validation:Minimum=3600
	// +kubebuilder:validation:Maximum=7200
	SessionLifetime *int `json:"sessionLifetime,omitempty"`
	// Sidecar configures a sidecar that intercepts every HTTP request, and performs the OIDC flow if necessary.
	// All requests to ingress + `/oauth2` will be processed only by the sidecar, whereas all other requests
	// will be proxied to the application.
	//
	// If the client is authenticated with IDPorten, the `Authorization` header will be set to `Bearer <JWT>`.
	// +nais:doc:Link="https://doc.nais.io/security/auth/idporten/"
	Sidecar *IDPortenSidecar `json:"sidecar,omitempty"`
}

type IDPortenSidecar struct {
	Wonderwall `json:",inline"`
	// Default security level for all authentication requests.
	// +nais:doc:Default="Level4"
	// +nais:doc:Link="https://doc.nais.io/security/auth/idporten#security-levels"
	// +kubebuilder:validation:Enum=Level3;Level4;idporten-loa-substantial;idporten-loa-high
	Level string `json:"level,omitempty"`
	// Default user interface locale for all authentication requests.
	// +nais:doc:Default="nb"
	// +nais:doc:Link="https://doc.nais.io/security/auth/idporten#locales"
	// +kubebuilder:validation:Enum=nb;nn;en;se
	Locale string `json:"locale,omitempty"`
}

type GCP struct {
	// Provision BigQuery datasets and give your application's pod mountable secrets for connecting to each dataset.
	// Datasets are immutable and cannot be changed.
	// +nais:doc:Link="https://cloud.google.com/bigquery/docs"
	// +nais:doc:Availability=GCP
	BigQueryDatasets []CloudBigQueryDataset `json:"bigQueryDatasets,omitempty"`
	// Provision cloud storage buckets and connect them to your application.
	// +nais:doc:Link="https://doc.nais.io/persistence/buckets/"
	// +nais:doc:Availability=GCP
	Buckets []CloudStorageBucket `json:"buckets,omitempty"`
	// Provision database instances and connect them to your application.
	// +nais:doc:Link="https://doc.nais.io/persistence/postgres/";"https://cloud.google.com/sql/docs/postgres/instance-settings#impact"
	// +nais:doc:Availability=GCP
	SqlInstances []CloudSqlInstance `json:"sqlInstances,omitempty"`
	// List of _additional_ permissions that should be granted to your application for accessing external GCP resources that have not been provisioned through NAIS.
	// +nais:doc:Link="https://doc.nais.io/nais-application/permissions-in-gcp/"
	// +nais:doc:Availability=GCP
	Permissions []CloudIAMPermission `json:"permissions,omitempty"`
}

type EnvVars []EnvVar

type EnvVar struct {
	// Environment variable name. May only contain letters, digits, and the underscore `_` character.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Environment variable value. Numbers and boolean values must be quoted.
	// Required unless `valueFrom` is specified.
	Value string `json:"value,omitempty"`
	// Dynamically set environment variables based on fields found in the Pod spec.
	// +nais:doc:Link="https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/"
	ValueFrom *EnvVarSource `json:"valueFrom,omitempty"`
}

type EnvVarSource struct {
	FieldRef ObjectFieldSelector `json:"fieldRef"`
}

type EnvFrom struct {
	// Name of the `ConfigMap` where environment variables are specified.
	// Required unless `secret` is set.
	ConfigMap string `json:"configmap,omitempty"`
	// Name of the `Secret` where environment variables are specified.
	// Required unless `configMap` is set.
	Secret string `json:"secret,omitempty"`
}

type ObjectFieldSelector struct {
	// Field value from the `Pod` spec that should be copied into the environment variable.
	// +kubebuilder:validation:Enum="";metadata.name;metadata.namespace;metadata.labels;metadata.annotations;spec.nodeName;spec.serviceAccountName;status.hostIP;status.podIP
	FieldPath string `json:"fieldPath"`
}

type FilesFrom struct {
	// Name of the `ConfigMap` that contains files that should be mounted into the container.
	// Required unless `secret` or `persistentVolumeClaim` is set.
	ConfigMap string `json:"configmap,omitempty"`
	// Name of the `Secret` that contains files that should be mounted into the container.
	// Required unless `configMap` or `persistentVolumeClaim` is set.
	// If mounting multiple secrets, `mountPath` *MUST* be set to avoid collisions.
	Secret string `json:"secret,omitempty"`
	// Specification of an empty directory
	EmptyDir *EmptyDir `json:"emptyDir,omitempty"`
	// Name of the `PersistentVolumeClaim` that should be mounted into the container.
	// Required unless `configMap` or `secret` is set.
	// This feature requires coordination with the NAIS team.
	PersistentVolumeClaim string `json:"persistentVolumeClaim,omitempty"`
	// Filesystem path inside the pod where files are mounted.
	// The directory will be created if it does not exist. If the directory exists,
	// any files in the directory will be made unaccessible.
	//
	// Defaults to `/var/run/configmaps/<NAME>`, `/var/run/secrets`, or `/var/run/pvc/<NAME>`, depending on which of them is specified.
	// For EmptyDir, MountPath must be set.
	MountPath string `json:"mountPath,omitempty"`
}

type ExecAction struct {
	// Command is the command line to execute inside the container before the pod is shut down.
	// The command is not run inside a shell, so traditional shell instructions (pipes, redirects, etc.) won't work.
	// To use a shell, you need to explicitly call out to that shell.
	//
	// If the exit status is non-zero, the pod will still be shut down, and marked as `Failed`.
	Command []string `json:"command,omitempty"`
}

type HttpGetAction struct {
	// Path to access on the HTTP server.
	Path string `json:"path"`
	// Port to access on the container.
	// Defaults to application port, as defined in `.spec.port`.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port *int `json:"port,omitempty"`
}

type PreStopHook struct {
	// Command that should be run inside the main container just before the pod is shut down by Kubernetes.
	Exec *ExecAction `json:"exec,omitempty"`
	// HTTP GET request that is called just before the pod is shut down by Kubernetes.
	Http *HttpGetAction `json:"http,omitempty"`
}

// Liveness probe and readiness probe definitions.
type Probe struct {
	// HTTP endpoint path that signals 200 OK if the application has started successfully.
	Path string `json:"path"`
	// Port for the startup probe.
	Port int `json:"port,omitempty"`
	// Number of seconds after the container has started before startup probes are initiated.
	InitialDelay int `json:"initialDelay,omitempty"`
	// How often (in seconds) to perform the probe.
	PeriodSeconds int `json:"periodSeconds,omitempty"`
	// When a Pod starts, and the probe fails, Kubernetes will try _failureThreshold_ times before giving up.
	// Giving up in case of a startup probe means restarting the Pod.
	FailureThreshold int `json:"failureThreshold,omitempty"`
	// Number of seconds after which the probe times out.
	Timeout int `json:"timeout,omitempty"`
}

type SecretPath struct {
	// File system path that the secret will be mounted into.
	MountPath string `json:"mountPath"`
	// Path to Vault key/value store that should be mounted into the file system.
	KvPath string `json:"kvPath"`
	// Format of the secret that should be processed.
	// +kubebuilder:validation:Enum=flatten;json;yaml;env;properties;""
	Format string `json:"format,omitempty"`
}

type Vault struct {
	// If set to true, fetch secrets from Vault and inject into the pods.
	Enabled bool `json:"enabled,omitempty"`
	// If enabled, the sidecar will automatically refresh the token's Time-To-Live before it expires.
	Sidecar bool `json:"sidecar,omitempty"`
	// List of secret paths to be read from Vault and injected into the pod's filesystem.
	// Overriding the `paths` array is optional, and will give you fine-grained control over which Vault paths that will be mounted on the file system.
	//
	// By default, the list will contain an entry with
	//
	// `kvPath: /kv/<environment>/<zone>/<application>/<namespace>`
	// `mountPath: /var/run/secrets/nais.io/vault`
	//
	// that will always be attempted to be mounted.
	Paths []SecretPath `json:"paths,omitempty"`
}

type Strategy struct {
	// Specifies the strategy used to replace old Pods by new ones.
	// `RollingUpdate` is the default value.
	// +kubebuilder:validation:Enum=Recreate;RollingUpdate
	Type          string                      `json:"type,omitempty"`
	RollingUpdate *v1.RollingUpdateDeployment `json:"rollingUpdate,omitempty"`
}

type Service struct {
	// +kubebuilder:validation:Enum=http;redis;tcp;grpc
	// Which protocol the backend service runs on. Default is `http`.
	Protocol string `json:"protocol,omitempty"`
	// Port for the default service. Default port is 80.
	Port int32 `json:"port"`
}

type TokenX struct {
	// If enabled, will provision and configure a TokenX client and inject an accompanying secret.
	Enabled bool `json:"enabled"`
	// If enabled, secrets for TokenX will be mounted as files only, i.e. not as environment variables.
	MountSecretsAsFilesOnly bool `json:"mountSecretsAsFilesOnly,omitempty"`
}

type Kafka struct {
	// Configures your application to access an Aiven Kafka cluster.
	Pool string `json:"pool"`

	// Allow this app to use kafka streams
	// +nais:doc:Link="https://doc.nais.io/persistence/kafka/application/#using-kafka-streams-with-internal-topics"
	// +nais:doc:Availability=GCP
	// +nais:doc:Default="false"
	Streams bool `json:"streams,omitempty"`
}

type CloudIAMResource struct {
	// Kubernetes _APIVersion_.
	APIVersion string `json:"apiVersion"`
	// Kubernetes _Kind_.
	Kind string `json:"kind"`
	// Kubernetes _Name_.
	Name string `json:"name,omitempty"`
}

type CloudIAMPermission struct {
	// Name of the GCP role to bind the resource to.
	Role string `json:"role"`
	// IAM resource to bind the role to.
	Resource CloudIAMResource `json:"resource"`
}

type Maskinporten struct {
	// If enabled, provisions and configures a Maskinporten client with consumed scopes and/or Exposed scopes with DigDir.
	// +nais:doc:Availability="team namespaces"
	// +nais:doc:Default="false"
	Enabled bool `json:"enabled"`
	// Schema to configure Maskinporten clients with consumed scopes and/or exposed scopes.
	Scopes MaskinportenScope `json:"scopes,omitempty"`
}

type SecureLogs struct {
	// Whether to enable a sidecar container for secure logging.
	// If enabled, a volume is mounted in the pods where secure logs can be saved.
	Enabled bool `json:"enabled"`
}

type PrometheusConfig struct {
	Enabled bool   `json:"enabled,omitempty"`
	Port    string `json:"port,omitempty"`
	Path    string `json:"path,omitempty"`
}

type Replicas struct {
	// The minimum amount of running replicas for a deployment.
	Min *int `json:"min,omitempty"`
	// The pod autoscaler will increase replicas when required up to the maximum.
	Max *int `json:"max,omitempty"`
	// Amount of CPU usage before the autoscaler kicks in.
	CpuThresholdPercentage int `json:"cpuThresholdPercentage,omitempty"`
	// Disable autoscaling
	// +nais:doc:Default="false"
	DisableAutoScaling bool `json:"disableAutoScaling,omitempty"`
}

type ResourceSpec struct {
	// +kubebuilder:validation:Pattern=^\d+m?$
	Cpu string `json:"cpu,omitempty"`
	// +kubebuilder:validation:Pattern=^\d+[KMG]i$
	Memory string `json:"memory,omitempty"`
}

type ResourceRequirements struct {
	// Limit defines the maximum amount of resources a container can use before getting evicted.
	Limits *ResourceSpec `json:"limits,omitempty"`
	// Request defines the amount of resources a container is allocated on startup.
	Requests *ResourceSpec `json:"requests,omitempty"`
}

// BigQueryPermission defines access level
type BigQueryPermission string

const (
	BigQueryPermissionRead      BigQueryPermission = "READ"
	BigQueryPermissionReadWrite BigQueryPermission = "READWRITE"
)

func (b BigQueryPermission) String() string {
	return string(b)
}

func (b BigQueryPermission) GoogleType() string {
	switch b {
	case BigQueryPermissionRead:
		return "READER"
	case BigQueryPermissionReadWrite:
		return "WRITER"
	}
	return ""
}

type CloudBigQueryDataset struct {
	// Name of the BigQuery Dataset.
	// The canonical name of the dataset will be `<TEAM_PROJECT_ID>:<NAME>`.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[a-z0-9][a-z0-9_]+$`
	// +nais:doc:Immutable=true
	Name string `json:"name" nais:"immutable,key"`
	// Permission level given to application.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=READ;READWRITE
	// +nais:doc:Immutable=true
	Permission BigQueryPermission `json:"permission" nais:"immutable"`
	// When set to true will delete the dataset, when the application resource is deleted.
	// NB: If no tables exist in the bigquery dataset, it _will_ delete the dataset even if this value is set/defaulted to `false`.
	// Default value is `false`.
	// +nais:doc:Immutable=true
	CascadingDelete bool `json:"cascadingDelete,omitempty" nais:"immutable"`
	// Human-readable description of what this BigQuery dataset contains, or is used for.
	// Will be visible in the GCP Console.
	// +nais:doc:Immutable=true
	Description string `json:"description,omitempty" nais:"immutable"`
}

type CloudStorageBucket struct {
	// The name of the bucket
	Name string `json:"name" nais:"immutable,key"`
	// Allows deletion of bucket. Set to true if you want to delete the bucket.
	CascadingDelete bool `json:"cascadingDelete,omitempty"`
	// The number of days to hold objects in the bucket before it is allowed to delete them.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=36500
	RetentionPeriodDays *int `json:"retentionPeriodDays,omitempty"`
	// Conditions for the bucket to use when selecting objects to delete in cleanup.
	// +nais:doc:Link="https://cloud.google.com/storage/docs/lifecycle"
	LifecycleCondition *LifecycleCondition `json:"lifecycleCondition,omitempty"`
	// Allows you to uniformly control access to your Cloud Storage resources.
	// When you enable uniform bucket-level access on a bucket, Access Control Lists (ACLs) are disabled, and only bucket-level Identity
	// and Access Management (IAM) permissions grant access to that bucket and the objects it contains.
	//
	// Uniform access control can not be reversed after 90 days! This is controlled by Google.
	// +nais:doc:Link="https://cloud.google.com/storage/docs/uniform-bucket-level-access"
	// +nais:doc:Default="false"
	UniformBucketLevelAccess bool `json:"uniformBucketLevelAccess,omitempty"`
	// Public access prevention allows you to prevent public access to your bucket.
	// +nais:doc:Link="https://cloud.google.com/storage/docs/public-access-prevention"
	// +nais:doc:Default="false"
	PublicAccessPrevention bool `json:"publicAccessPrevention,omitempty"`
}

type LifecycleCondition struct {
	// Condition is satisfied when the object reaches the specified age in days. These will be deleted.
	Age int `json:"age,omitempty"`
	// Condition is satisfied when the object is created before midnight on the specified date. These will be deleted.
	CreatedBefore string `json:"createdBefore,omitempty"`
	// Condition is satisfied when the object has the specified number of newer versions.
	// The older versions will be deleted.
	NumNewerVersions int `json:"numNewerVersions,omitempty"`
	// Condition is satisfied when the object has the specified state.
	// +kubebuilder:validation:Enum="";LIVE;ARCHIVED;ANY
	WithState string `json:"withState,omitempty"`
}

type CloudSqlInstanceType string

const (
	CloudSqlInstanceTypePostgres11 CloudSqlInstanceType = "POSTGRES_11"
	CloudSqlInstanceTypePostgres12 CloudSqlInstanceType = "POSTGRES_12"
	CloudSqlInstanceTypePostgres13 CloudSqlInstanceType = "POSTGRES_13"
	CloudSqlInstanceTypePostgres14 CloudSqlInstanceType = "POSTGRES_14"
	CloudSqlInstanceTypePostgres15 CloudSqlInstanceType = "POSTGRES_15"
)

type CloudSqlInstanceDiskType string

func (c CloudSqlInstanceDiskType) GoogleType() string {
	return fmt.Sprintf("PD_%s", c)
}

const (
	CloudSqlInstanceDiskTypeSSD CloudSqlInstanceDiskType = "SSD"
	CloudSqlInstanceDiskTypeHDD CloudSqlInstanceDiskType = "HDD"
)

type CloudSqlDatabase struct {
	// Database name.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Prefix to add to environment variables made available for database connection.
	// If switching to `EnvVarPrefix` you need to [reset database credentials](https://docs.nais.io/persistence/postgres/#reset-database-credentials).
	EnvVarPrefix string `json:"envVarPrefix,omitempty"`
	// Add extra users for database access. These users need to be manually given access to database tables.
	Users []CloudSqlDatabaseUser `json:"users,omitempty"`
}

type CloudSqlDatabaseUser struct {
	// User name.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^[_a-zA-Z][-_a-zA-Z0-9]+$"
	Name string `json:"name"`
}

type CloudSqlFlag struct {
	// Name of the flag.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Value of the flag.
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

type CloudSqlInstance struct {
	// PostgreSQL version.
	// +kubebuilder:validation:Enum=POSTGRES_11;POSTGRES_12;POSTGRES_13;POSTGRES_14;POSTGRES_15
	// +kubebuilder:validation:Required
	// +nais:doc:Link="https://cloud.google.com/sql/docs/postgres/instance-settings"
	Type CloudSqlInstanceType `json:"type"`
	// The name of the instance, if omitted the application name will be used.
	Name string `json:"name,omitempty"`
	// Server tier, i.e. how much CPU and memory allocated.
	// Available tiers are `db-f1-micro`, `db-g1-small` and custom `db-custom-CPU-RAM`.
	// Custom memory must be mulitple of 256 MB and at least 3.75 GB (e.g. `db-custom-1-3840` for 1 cpu, 3840 MB ram)
	// Also check out [sizing your database](../../persistence/postgres/#sizing-your-database).
	// +kubebuilder:validation:Pattern="db-.+"
	// +nais:doc:Default="db-f1-micro"
	Tier string `json:"tier,omitempty"`
	// Disk type to use for storage in the database.
	// +kubebuilder:validation:Enum=SSD;HDD
	DiskType CloudSqlInstanceDiskType `json:"diskType,omitempty"`
	// When set to true this will set up standby database for failover.
	HighAvailability bool `json:"highAvailability,omitempty"`
	// How much hard drive space to allocate for the SQL server, in gigabytes.
	// This parameter is used when first provisioning a server.
	// Disk size can be changed using this field _only when diskAutoresize is set to false_.
	// +kubebuilder:validation:Minimum=10
	DiskSize int `json:"diskSize,omitempty"`
	// When set to true, GCP will automatically increase storage by XXX for the database when
	// disk usage is above the high water mark. Setting this field to true also disables
	// manual control over disk size, i.e. the `diskSize` parameter will be ignored.
	// +nais:doc:Link="https://cloud.google.com/sql/docs/postgres/instance-settings#threshold"
	DiskAutoresize bool `json:"diskAutoresize,omitempty"`
	// If specified, run automatic backups of the SQL database at the given hour.
	// Note that this will backup the whole SQL instance, and not separate databases.
	// Restores are done using the Google Cloud Console.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	AutoBackupHour *int `json:"autoBackupHour,omitempty"`
	// Number of daily backups to retain. Defaults to 7 backups.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=365
	// +nais:doc:Default="7"
	// +nais:doc:Link="https://cloud.google.com/sql/docs/postgres/backup-recovery/backups"
	RetainedBackups *int `json:"retainedBackups,omitempty"`
	// Desired maintenance window for database updates.
	Maintenance *Maintenance `json:"maintenance,omitempty"`
	// List of databases that should be created on this Postgres server.
	// +kubebuilder:validation:Required
	Databases []CloudSqlDatabase `json:"databases,omitempty"`
	// Remove the entire Postgres server including all data when the Kubernetes resource is deleted.
	// *THIS IS A DESTRUCTIVE OPERATION*! Set cascading delete only when you want to remove data forever.
	CascadingDelete bool `json:"cascadingDelete,omitempty"`
	// Sort order for `ORDER BY ...` clauses.
	Collation string `json:"collation,omitempty"`
	// Enables point-in-time recovery for sql instances using write-ahead logs.
	PointInTimeRecovery bool `json:"pointInTimeRecovery,omitempty"`
	// Configures query insights which are now default for new sql instances.
	Insights *InsightsConfiguration `json:"insights,omitempty"`
	// Set flags to control the behavior of the instance.
	// Be aware that NAIS _does not validate_ these flags, so take extra care
	// to make sure the values match against the specification, otherwise your deployment
	// will seemingly work OK, but the database flags will not function as expected.
	// +nais:doc:Link="https://cloud.google.com/sql/docs/postgres/flags#list-flags-postgres"
	// +nais:doc:Experimental=true
	Flags []CloudSqlFlag `json:"flags,omitempty"`
}

type InsightsConfiguration struct {
	// True if Query Insights feature is enabled.
	// +nais:doc:Default="true"
	Enabled *bool `json:"enabled,omitempty"`
	// Maximum query length stored in bytes. Between 256 and 4500. Default to 1024.
	// +kubebuilder:validation:Minimum=256
	// +kubebuilder:validation:Maximum=4500
	QueryStringLength int `json:"queryStringLength,omitempty"`
	// True if Query Insights will record application tags from query when enabled.
	RecordApplicationTags bool `json:"recordApplicationTags,omitempty"`
	// True if Query Insights will record client address when enabled.
	RecordClientAddress bool `json:"recordClientAddress,omitempty"`
}

// IsEnabled returns true if Enabled is true, nil or if InsightsConfiguration is nil.
func (i *InsightsConfiguration) IsEnabled() bool {
	if i == nil || i.Enabled == nil {
		return true
	}

	return *i.Enabled
}

type MediumType string

const (
	MediumTypeMemory MediumType = "Memory"
	MediumTypeDisk   MediumType = "Disk"
)

type EmptyDir struct {
	// +kubebuilder:validation:Enum=Memory;Disk
	Medium MediumType `json:"medium,omitempty"`
}

type Maintenance struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=7
	Day int `json:"day,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	Hour *int `json:"hour,omitempty"` // must use pointer here to be able to distinguish between no value and value 0 from user.
}

func (envVar EnvVar) ToKubernetes() corev1.EnvVar {
	if envVar.ValueFrom != nil && envVar.ValueFrom.FieldRef.FieldPath != "" {
		return corev1.EnvVar{
			Name: envVar.Name,
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{FieldPath: envVar.ValueFrom.FieldRef.FieldPath},
			},
		}
	} else {
		return corev1.EnvVar{Name: envVar.Name, Value: envVar.Value}
	}
}

// Maps environment variables from ApplicationSpec to the ones we use in CreateSpec
func (envVars EnvVars) ToKubernetes() []corev1.EnvVar {
	var newEnvVars []corev1.EnvVar
	for _, envVar := range envVars {
		newEnvVars = append(newEnvVars, envVar.ToKubernetes())
	}

	return newEnvVars
}

type Wonderwall struct {
	// Automatically redirect the user to login for all proxied GET requests.
	// +nais:doc:Default="false"
	// +nais:doc:Link="https://doc.nais.io/appendix/wonderwall/#12-autologin"
	AutoLogin bool `json:"autoLogin,omitempty"`
	// Comma separated list of absolute paths to ignore when auto-login is enabled.
	// +nais:doc:Link="https://doc.nais.io/appendix/wonderwall/#12-autologin"
	AutoLoginIgnorePaths []WonderwallIgnorePaths `json:"autoLoginIgnorePaths,omitempty"`
	// Enable the sidecar.
	Enabled bool `json:"enabled"`
	// Resource requirements for the sidecar container.
	// +nais:doc:Link="https://doc.nais.io/appendix/wonderwall/#4-resource-requirements"
	Resources *ResourceRequirements `json:"resources,omitempty"`
}

// +kubebuilder:validation:Pattern=`^\/.*$`
type WonderwallIgnorePaths string

type FrontendGeneratedConfig struct {
	// If specified, a Javascript file with application specific frontend configuration variables
	// will be generated and mounted into the pod file system at the specified path.
	// You can import this file directly from your Javascript application.
	// +nais:doc:Link="https://doc.nais.io/observability/frontend/#auto-configuration"
	MountPath string `json:"mountPath"`
}

type Frontend struct {
	GeneratedConfig *FrontendGeneratedConfig `json:"generatedConfig,omitempty"`
}

type Tracing struct {
	Enabled bool `json:"enabled,omitempty"`
}

type Observability struct {
	// Enable application performance monitoring with traces collected using OpenTelemetry and the OTLP exporter.
	// +nais:doc:Availability="GCP"
	// +nais:doc:Experimental=true
	// +nais:doc:Link="https://doc.nais.io/observability/tracing/"
	Tracing *Tracing `json:"tracing,omitempty"`
}
