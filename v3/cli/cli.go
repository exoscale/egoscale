package main

import (
	"context"
	"flag"
	"fmt"
	v3 "github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/credentials"
	"os"
	"strings"
)

func printUsage() {
	commands := []struct {
		Name, Doc string
	}{
		{"list-anti-affinity-groups", "List Anti-affinity Groups"},
		{"create-anti-affinity-group", "Create an Anti-affinity Group"},
		{"delete-anti-affinity-group", "Delete an Anti-affinity Group"},
		{"get-anti-affinity-group", "Retrieve Anti-affinity Group details"},
		{"list-api-keys", "List API keys"},
		{"create-api-key", "Create a new API key"},
		{"delete-api-key", "Delete an API key"},
		{"get-api-key", "Get API key"},
		{"list-block-storage-volumes", "List block storage volumes"},
		{"create-block-storage-volume", "Create a block storage volume"},
		{"list-block-storage-snapshots", "List block storage snapshots"},
		{"delete-block-storage-snapshot", "Delete a block storage snapshot, data will be unrecoverable"},
		{"get-block-storage-snapshot", "Retrieve block storage snapshot details"},
		{"update-block-storage-snapshot", "Update block storage volume snapshot"},
		{"delete-block-storage-volume", "Delete a block storage volume, data will be unrecoverable"},
		{"get-block-storage-volume", "Retrieve block storage volume details"},
		{"update-block-storage-volume", "Update block storage volume"},
		{"attach-block-storage-volume-to-instance", "Attach block storage volume to an instance"},
		{"create-block-storage-snapshot", "Create a block storage snapshot"},
		{"detach-block-storage-volume", "Detach block storage volume"},
		{"resize-block-storage-volume", "Resize a block storage volume"},
		{"get-console-proxy-url", "Retrieve signed url valid for 60 seconds to connect via console-proxy websocket to VM VNC console."},
		{"get-dbaas-ca-certificate", "Get DBaaS CA Certificate"},
		{"delete-dbaas-external-endpoint-datadog", ""},
		{"get-dbaas-external-endpoint-datadog", ""},
		{"update-dbaas-external-endpoint-datadog", ""},
		{"create-dbaas-external-endpoint-datadog", ""},
		{"delete-dbaas-external-endpoint-elasticsearch", ""},
		{"get-dbaas-external-endpoint-elasticsearch", ""},
		{"update-dbaas-external-endpoint-elasticsearch", ""},
		{"create-dbaas-external-endpoint-elasticsearch", ""},
		{"delete-dbaas-external-endpoint-opensearch", ""},
		{"get-dbaas-external-endpoint-opensearch", ""},
		{"update-dbaas-external-endpoint-opensearch", ""},
		{"create-dbaas-external-endpoint-opensearch", ""},
		{"delete-dbaas-external-endpoint-prometheus", ""},
		{"get-dbaas-external-endpoint-prometheus", ""},
		{"update-dbaas-external-endpoint-prometheus", ""},
		{"create-dbaas-external-endpoint-prometheus", ""},
		{"delete-dbaas-external-endpoint-rsyslog", ""},
		{"get-dbaas-external-endpoint-rsyslog", ""},
		{"update-dbaas-external-endpoint-rsyslog", ""},
		{"create-dbaas-external-endpoint-rsyslog", ""},
		{"list-dbaas-external-endpoint-types", ""},
		{"attach-dbaas-service-to-endpoint", ""},
		{"detach-dbaas-service-from-endpoint", ""},
		{"list-dbaas-external-endpoints", ""},
		{"get-dbaas-external-integration-settings-datadog", ""},
		{"update-dbaas-external-integration-settings-datadog", ""},
		{"get-dbaas-external-integration", ""},
		{"list-dbaas-external-integrations", ""},
		{"delete-dbaas-service-grafana", "Delete a Grafana service"},
		{"get-dbaas-service-grafana", "Get a DBaaS Grafana service"},
		{"create-dbaas-service-grafana", ""},
		{"update-dbaas-service-grafana", "Update a DBaaS Grafana service"},
		{"start-dbaas-grafana-maintenance", "Initiate Grafana maintenance update"},
		{"reset-dbaas-grafana-user-password", "Reset the credentials of a DBaaS Grafana user"},
		{"reveal-dbaas-grafana-user-password", "Reveal the secrets of a DBaaS Grafana user"},
		{"create-dbaas-integration", ""},
		{"list-dbaas-integration-settings", ""},
		{"list-dbaas-integration-types", ""},
		{"delete-dbaas-integration", ""},
		{"get-dbaas-integration", ""},
		{"update-dbaas-integration", ""},
		{"delete-dbaas-service-kafka", "Delete a Kafka service"},
		{"get-dbaas-service-kafka", "Get a DBaaS Kafka service"},
		{"create-dbaas-service-kafka", "Create a DBaaS Kafka service"},
		{"update-dbaas-service-kafka", "Update a DBaaS Kafka service"},
		{"get-dbaas-kafka-acl-config", "Get DBaaS kafka ACL configuration"},
		{"start-dbaas-kafka-maintenance", "Initiate Kafka maintenance update"},
		{"create-dbaas-kafka-schema-registry-acl-config", "Add a Kafka Schema Registry ACL entry"},
		{"delete-dbaas-kafka-schema-registry-acl-config", "Delete a Kafka ACL entry"},
		{"create-dbaas-kafka-topic-acl-config", "Add a Kafka topic ACL entry"},
		{"delete-dbaas-kafka-topic-acl-config", "Delete a Kafka ACL entry"},
		{"reveal-dbaas-kafka-connect-password", "Reveal the secrets for DBaaS Kafka Connect"},
		{"create-dbaas-kafka-user", "Create a DBaaS Kafka user"},
		{"delete-dbaas-kafka-user", "Delete a DBaaS kafka user"},
		{"reset-dbaas-kafka-user-password", "Reset the credentials of a DBaaS Kafka user"},
		{"reveal-dbaas-kafka-user-password", "Reveal the secrets of a DBaaS Kafka user"},
		{"get-dbaas-migration-status", "Get a DBaaS migration status"},
		{"delete-dbaas-service-mysql", "Delete a MySQL service"},
		{"get-dbaas-service-mysql", "Get a DBaaS MySQL service"},
		{"create-dbaas-service-mysql", "Create a DBaaS MySQL service"},
		{"update-dbaas-service-mysql", "Update a DBaaS MySQL service"},
		{"enable-dbaas-mysql-writes", "Temporarily enable writes for MySQL services in read-only mode due to filled up storage"},
		{"start-dbaas-mysql-maintenance", "Initiate MySQL maintenance update"},
		{"stop-dbaas-mysql-migration", "Stop a DBaaS MySQL migration"},
		{"create-dbaas-mysql-database", "Create a DBaaS MySQL database"},
		{"delete-dbaas-mysql-database", "Delete a DBaaS MySQL database"},
		{"create-dbaas-mysql-user", "Create a DBaaS MySQL user"},
		{"delete-dbaas-mysql-user", "Delete a DBaaS MySQL user"},
		{"reset-dbaas-mysql-user-password", "Reset the credentials of a DBaaS mysql user"},
		{"reveal-dbaas-mysql-user-password", "Reveal the secrets of a DBaaS MySQL user"},
		{"delete-dbaas-service-opensearch", "Delete a OpenSearch service"},
		{"get-dbaas-service-opensearch", "Get a DBaaS OpenSearch service"},
		{"create-dbaas-service-opensearch", "Create a DBaaS OpenSearch service"},
		{"update-dbaas-service-opensearch", "Update a DBaaS OpenSearch service"},
		{"get-dbaas-opensearch-acl-config", "Get DBaaS OpenSearch ACL configuration"},
		{"update-dbaas-opensearch-acl-config", "Create a DBaaS OpenSearch ACL configuration"},
		{"start-dbaas-opensearch-maintenance", "Initiate OpenSearch maintenance update"},
		{"create-dbaas-opensearch-user", "Create a DBaaS OpenSearch user"},
		{"delete-dbaas-opensearch-user", "Delete a DBaaS OpenSearch user"},
		{"reset-dbaas-opensearch-user-password", "Reset the credentials of a DBaaS OpenSearch user"},
		{"reveal-dbaas-opensearch-user-password", "Reveal the secrets of a DBaaS OpenSearch user"},
		{"delete-dbaas-service-pg", "Delete a Postgres service"},
		{"get-dbaas-service-pg", "Get a DBaaS PostgreSQL service"},
		{"create-dbaas-service-pg", "Create a DBaaS PostgreSQL service"},
		{"update-dbaas-service-pg", "Update a DBaaS PostgreSQL service"},
		{"start-dbaas-pg-maintenance", "Initiate PostgreSQL maintenance update"},
		{"stop-dbaas-pg-migration", "Stop a DBaaS PostgreSQL migration"},
		{"create-dbaas-pg-connection-pool", "Create a DBaaS PostgreSQL connection pool"},
		{"delete-dbaas-pg-connection-pool", "Delete a DBaaS PostgreSQL connection pool"},
		{"update-dbaas-pg-connection-pool", "Update a DBaaS PostgreSQL connection pool"},
		{"create-dbaas-pg-database", "Create a DBaaS Postgres database"},
		{"delete-dbaas-pg-database", "Delete a DBaaS Postgres database"},
		{"create-dbaas-postgres-user", "Create a DBaaS Postgres user"},
		{"delete-dbaas-postgres-user", "Delete a DBaaS Postgres user"},
		{"update-dbaas-postgres-allow-replication", "Update access control for one service user"},
		{"reset-dbaas-postgres-user-password", "Reset the credentials of a DBaaS Postgres user"},
		{"reveal-dbaas-postgres-user-password", "Reveal the secrets of a DBaaS Postgres user"},
		{"create-dbaas-pg-upgrade-check", ""},
		{"delete-dbaas-service-redis", "Delete a Redis service"},
		{"get-dbaas-service-redis", "Get a DBaaS Redis service"},
		{"create-dbaas-service-redis", "Create a DBaaS Redis service"},
		{"update-dbaas-service-redis", "Update a DBaaS Redis service"},
		{"start-dbaas-redis-maintenance", "Initiate Redis maintenance update"},
		{"stop-dbaas-redis-migration", "Stop a DBaaS Redis migration"},
		{"start-dbaas-redis-to-valkey-upgrade", "Initiate Redis upgrade to Valkey"},
		{"create-dbaas-redis-user", "Create a DBaaS Redis user"},
		{"delete-dbaas-redis-user", "Delete a DBaaS Redis user"},
		{"reset-dbaas-redis-user-password", "Reset the credentials of a DBaaS Redis user"},
		{"reveal-dbaas-redis-user-password", "Reveal the secrets of a DBaaS Redis user"},
		{"list-dbaas-services", "List DBaaS services"},
		{"get-dbaas-service-logs", "Get logs of DBaaS service"},
		{"get-dbaas-service-metrics", "Get metrics of DBaaS service"},
		{"list-dbaas-service-types", "DBaaS Service Types"},
		{"get-dbaas-service-type", "Get a DBaaS service type"},
		{"delete-dbaas-service", "Delete a DBaaS service"},
		{"get-dbaas-settings-grafana", "Get DBaaS Grafana settings"},
		{"get-dbaas-settings-kafka", "Get DBaaS Kafka settings"},
		{"get-dbaas-settings-mysql", "Get DBaaS MySQL settings"},
		{"get-dbaas-settings-opensearch", "Get DBaaS OpenSearch settings"},
		{"get-dbaas-settings-pg", "Get DBaaS PostgreSQL settings"},
		{"get-dbaas-settings-redis", "Get DBaaS Redis settings"},
		{"get-dbaas-settings-valkey", "Get DBaaS Valkey settings"},
		{"create-dbaas-task-migration-check", ""},
		{"get-dbaas-task", "Get a DBaaS task"},
		{"delete-dbaas-service-valkey", "Delete a Valkey service"},
		{"get-dbaas-service-valkey", ""},
		{"create-dbaas-service-valkey", "Create a DBaaS Valkey service"},
		{"update-dbaas-service-valkey", ""},
		{"start-dbaas-valkey-maintenance", "Initiate Valkey maintenance update"},
		{"stop-dbaas-valkey-migration", "Stop a DBaaS Valkey migration"},
		{"create-dbaas-valkey-user", "Create a DBaaS Valkey user"},
		{"delete-dbaas-valkey-user", "Delete a DBaaS Valkey user"},
		{"reset-dbaas-valkey-user-password", "Reset the credentials of a DBaaS Valkey user"},
		{"reveal-dbaas-valkey-user-password", "Reveal the secrets of a DBaaS Valkey user"},
		{"list-deploy-targets", "List Deploy Targets"},
		{"get-deploy-target", "Retrieve Deploy Target details"},
		{"list-dns-domains", "List DNS domains"},
		{"create-dns-domain", "Create DNS domain"},
		{"list-dns-domain-records", "List DNS domain records"},
		{"create-dns-domain-record", "Create DNS domain record"},
		{"delete-dns-domain-record", "Delete DNS domain record"},
		{"get-dns-domain-record", "Retrieve DNS domain record details"},
		{"update-dns-domain-record", "Update DNS domain record"},
		{"delete-dns-domain", "Delete DNS Domain"},
		{"get-dns-domain", "Retrieve DNS domain details"},
		{"get-dns-domain-zone-file", "Retrieve DNS domain zone file"},
		{"list-elastic-ips", "List Elastic IPs"},
		{"create-elastic-ip", "Create an Elastic IP"},
		{"delete-elastic-ip", "Delete an Elastic IP"},
		{"get-elastic-ip", "Retrieve Elastic IP details"},
		{"update-elastic-ip", "Update an Elastic IP"},
		{"reset-elastic-ip-field", "Reset an Elastic IP field to its default value"},
		{"attach-instance-to-elastic-ip", "Attach a Compute instance to an Elastic IP"},
		{"detach-instance-from-elastic-ip", "Detach a Compute instance from an Elastic IP"},
		{"get-env-impact", "[BETA] Retrieve organization environmental impact reports"},
		{"list-events", "List Events"},
		{"get-iam-organization-policy", "Retrieve IAM Organization Policy"},
		{"update-iam-organization-policy", "Update IAM Organization Policy"},
		{"reset-iam-organization-policy", "Reset IAM Organization Policy"},
		{"list-iam-roles", "List IAM Roles"},
		{"create-iam-role", "Create IAM Role"},
		{"delete-iam-role", "Delete IAM Role"},
		{"get-iam-role", "Retrieve IAM Role"},
		{"update-iam-role", "Update IAM Role"},
		{"update-iam-role-policy", "Update IAM Role Policy"},
		{"list-instances", "List Compute instances"},
		{"create-instance", "Create a Compute instance"},
		{"list-instance-pools", "List Instance Pools"},
		{"create-instance-pool", "Create an Instance Pool"},
		{"delete-instance-pool", "Delete an Instance Pool"},
		{"get-instance-pool", "Retrieve Instance Pool details"},
		{"update-instance-pool", "Update an Instance Pool"},
		{"reset-instance-pool-field", "Reset an Instance Pool field to its default value"},
		{"evict-instance-pool-members", "Evict Instance Pool members"},
		{"scale-instance-pool", "Scale an Instance Pool"},
		{"list-instance-types", "List Compute instance Types"},
		{"get-instance-type", "Retrieve Instance Type details"},
		{"delete-instance", "Delete a Compute instance"},
		{"get-instance", "Retrieve Compute instance details"},
		{"update-instance", "Update a Compute instance"},
		{"reset-instance-field", "Reset Instance field"},
		{"add-instance-protection", "Set instance destruction protection"},
		{"create-snapshot", "Create a Snapshot of a Compute instance"},
		{"enable-tpm", "[Beta] Enable tpm for the instance."},
		{"reveal-instance-password", "Reveal the password used during instance creation or the latest password reset."},
		{"reboot-instance", "Reboot a Compute instance"},
		{"remove-instance-protection", "Remove instance destruction protection"},
		{"reset-instance", "Reset a Compute instance to a base/target template"},
		{"reset-instance-password", "Reset a compute instance password"},
		{"resize-instance-disk", "Resize a Compute instance disk"},
		{"scale-instance", "Scale a Compute instance to a new Instance Type"},
		{"start-instance", "Start a Compute instance"},
		{"stop-instance", "Stop a Compute instance"},
		{"revert-instance-to-snapshot", "Revert a snapshot for an instance"},
		{"list-load-balancers", "List Load Balancers"},
		{"create-load-balancer", "Create a Load Balancer"},
		{"delete-load-balancer", "Delete a Load Balancer"},
		{"get-load-balancer", "Retrieve Load Balancer details"},
		{"update-load-balancer", "Update a Load Balancer"},
		{"add-service-to-load-balancer", "Add a Load Balancer Service"},
		{"delete-load-balancer-service", "Delete a Load Balancer Service"},
		{"get-load-balancer-service", "Retrieve Load Balancer Service details"},
		{"update-load-balancer-service", "Update a Load Balancer Service"},
		{"reset-load-balancer-service-field", "Reset a Load Balancer Service field to its default value"},
		{"reset-load-balancer-field", "Reset a Load Balancer field to its default value"},
		{"get-operation", "Retrieve Operation details"},
		{"get-organization", "Retrieve an organization"},
		{"list-private-networks", "List Private Networks"},
		{"create-private-network", "Create a Private Network"},
		{"delete-private-network", "Delete a Private Network"},
		{"get-private-network", "Retrieve Private Network details"},
		{"update-private-network", "Update a Private Network"},
		{"reset-private-network-field", "Reset Private Network field"},
		{"attach-instance-to-private-network", "Attach a Compute instance to a Private Network"},
		{"detach-instance-from-private-network", "Detach a Compute instance from a Private Network"},
		{"update-private-network-instance-ip", "Update the IP address of an instance attached to a managed private network"},
		{"list-quotas", "List Organization Quotas"},
		{"get-quota", "Retrieve Resource Quota"},
		{"delete-reverse-dns-elastic-ip", "Delete the PTR DNS record for an elastic IP"},
		{"get-reverse-dns-elastic-ip", "Query the PTR DNS records for an elastic IP"},
		{"update-reverse-dns-elastic-ip", "Update/Create the PTR DNS record for an elastic IP"},
		{"delete-reverse-dns-instance", "Delete the PTR DNS record for an instance"},
		{"get-reverse-dns-instance", "Query the PTR DNS records for an instance"},
		{"update-reverse-dns-instance", "Update/Create the PTR DNS record for an instance"},
		{"list-security-groups", "List Security Groups."},
		{"create-security-group", "Create a Security Group"},
		{"delete-security-group", "Delete a Security Group"},
		{"get-security-group", "Retrieve Security Group details"},
		{"add-rule-to-security-group", "Create a Security Group rule"},
		{"delete-rule-from-security-group", "Delete a Security Group rule"},
		{"add-external-source-to-security-group", "Add an external source as a member of a Security Group"},
		{"attach-instance-to-security-group", "Attach a Compute instance to a Security Group"},
		{"detach-instance-from-security-group", "Detach a Compute instance from a Security Group"},
		{"remove-external-source-from-security-group", "Remove an external source from a Security Group"},
		{"list-sks-clusters", "List SKS clusters"},
		{"create-sks-cluster", "Create an SKS cluster"},
		{"list-sks-cluster-deprecated-resources", "Resources that are scheduled to be removed in future kubernetes releases"},
		{"generate-sks-cluster-kubeconfig", "Generate a new Kubeconfig file for a SKS cluster"},
		{"list-sks-cluster-versions", "List available versions for SKS clusters"},
		{"delete-sks-cluster", "Delete an SKS cluster"},
		{"get-sks-cluster", "Retrieve SKS cluster details"},
		{"update-sks-cluster", "Update an SKS cluster"},
		{"get-sks-cluster-authority-cert", "Get the certificate for a SKS cluster authority"},
		{"get-sks-cluster-inspection", "Get the latest inspection result"},
		{"create-sks-nodepool", "Create a new SKS Nodepool"},
		{"delete-sks-nodepool", "Delete an SKS Nodepool"},
		{"get-sks-nodepool", "Retrieve SKS Nodepool details"},
		{"update-sks-nodepool", "Update an SKS Nodepool"},
		{"reset-sks-nodepool-field", "Reset an SKS Nodepool field to its default value"},
		{"evict-sks-nodepool-members", "Evict Nodepool members"},
		{"scale-sks-nodepool", "Scale a SKS Nodepool"},
		{"rotate-sks-ccm-credentials", "Rotate Exoscale CCM credentials"},
		{"rotate-sks-csi-credentials", "Rotate Exoscale CSI credentials"},
		{"rotate-sks-operators-ca", "Rotate operators certificate authority"},
		{"upgrade-sks-cluster", "Upgrade an SKS cluster"},
		{"upgrade-sks-cluster-service-level", "Upgrade a SKS cluster to pro"},
		{"reset-sks-cluster-field", "Reset an SKS cluster field to its default value"},
		{"list-snapshots", "List Snapshots"},
		{"delete-snapshot", "Delete a Snapshot"},
		{"get-snapshot", "Retrieve Snapshot details"},
		{"export-snapshot", "Export a Snapshot"},
		{"promote-snapshot-to-template", "Promote a Snapshot to a Template"},
		{"list-sos-buckets-usage", "List SOS Buckets Usage"},
		{"get-sos-presigned-url", "Retrieve Presigned Download URL for SOS object"},
		{"list-ssh-keys", "List SSH keys"},
		{"register-ssh-key", "Import SSH key"},
		{"delete-ssh-key", "Delete a SSH key"},
		{"get-ssh-key", "Retrieve SSH key details"},
		{"list-templates", "List Templates"},
		{"register-template", "Register a Template"},
		{"delete-template", "Delete a Template"},
		{"get-template", "Retrieve Template details"},
		{"copy-template", "Copy a Template from a zone to another"},
		{"update-template", "Update template attributes"},
		{"get-usage-report", "Retrieve organization usage reports"},
		{"list-users", "List Users"},
		{"create-user", "Create a User"},
		{"delete-user", "Delete User"},
		{"update-user-role", "Update a User's IAM role"},
		{"list-zones", "List Zones"},
	}
	maxLen := 0
	for _, c := range commands {
		if l := len(c.Name); l > maxLen {
			maxLen = l
		}
	}
	fmt.Println("Usage: " + os.Args[0] + " <command>")
	fmt.Println("Available commands:")
	for _, c := range commands {
		fmt.Printf("  %-*s %s\n", maxLen, c.Name, c.Doc)
	}
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// TODO: Make credentials configurable via flags/env
	client, err := v3.NewClient(credentials.NewEnvCredentials())
	if err != nil {
		fmt.Println("failed to create client:", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list-anti-affinity-groups":
		ListAntiAffinityGroupsCmd(client)
	case "create-anti-affinity-group":
		CreateAntiAffinityGroupCmd(client)
	case "delete-anti-affinity-group":
		DeleteAntiAffinityGroupCmd(client)
	case "get-anti-affinity-group":
		GetAntiAffinityGroupCmd(client)
	case "list-api-keys":
		ListAPIKeysCmd(client)
	case "create-api-key":
		CreateAPIKeyCmd(client)
	case "delete-api-key":
		DeleteAPIKeyCmd(client)
	case "get-api-key":
		GetAPIKeyCmd(client)
	case "list-block-storage-volumes":
		ListBlockStorageVolumesCmd(client)
	case "create-block-storage-volume":
		CreateBlockStorageVolumeCmd(client)
	case "list-block-storage-snapshots":
		ListBlockStorageSnapshotsCmd(client)
	case "delete-block-storage-snapshot":
		DeleteBlockStorageSnapshotCmd(client)
	case "get-block-storage-snapshot":
		GetBlockStorageSnapshotCmd(client)
	case "update-block-storage-snapshot":
		UpdateBlockStorageSnapshotCmd(client)
	case "delete-block-storage-volume":
		DeleteBlockStorageVolumeCmd(client)
	case "get-block-storage-volume":
		GetBlockStorageVolumeCmd(client)
	case "update-block-storage-volume":
		UpdateBlockStorageVolumeCmd(client)
	case "attach-block-storage-volume-to-instance":
		AttachBlockStorageVolumeToInstanceCmd(client)
	case "create-block-storage-snapshot":
		CreateBlockStorageSnapshotCmd(client)
	case "detach-block-storage-volume":
		DetachBlockStorageVolumeCmd(client)
	case "resize-block-storage-volume":
		ResizeBlockStorageVolumeCmd(client)
	case "get-console-proxy-url":
		GetConsoleProxyURLCmd(client)
	case "get-dbaas-ca-certificate":
		GetDBAASCACertificateCmd(client)
	case "delete-dbaas-external-endpoint-datadog":
		DeleteDBAASExternalEndpointDatadogCmd(client)
	case "get-dbaas-external-endpoint-datadog":
		GetDBAASExternalEndpointDatadogCmd(client)
	case "update-dbaas-external-endpoint-datadog":
		UpdateDBAASExternalEndpointDatadogCmd(client)
	case "create-dbaas-external-endpoint-datadog":
		CreateDBAASExternalEndpointDatadogCmd(client)
	case "delete-dbaas-external-endpoint-elasticsearch":
		DeleteDBAASExternalEndpointElasticsearchCmd(client)
	case "get-dbaas-external-endpoint-elasticsearch":
		GetDBAASExternalEndpointElasticsearchCmd(client)
	case "update-dbaas-external-endpoint-elasticsearch":
		UpdateDBAASExternalEndpointElasticsearchCmd(client)
	case "create-dbaas-external-endpoint-elasticsearch":
		CreateDBAASExternalEndpointElasticsearchCmd(client)
	case "delete-dbaas-external-endpoint-opensearch":
		DeleteDBAASExternalEndpointOpensearchCmd(client)
	case "get-dbaas-external-endpoint-opensearch":
		GetDBAASExternalEndpointOpensearchCmd(client)
	case "update-dbaas-external-endpoint-opensearch":
		UpdateDBAASExternalEndpointOpensearchCmd(client)
	case "create-dbaas-external-endpoint-opensearch":
		CreateDBAASExternalEndpointOpensearchCmd(client)
	case "delete-dbaas-external-endpoint-prometheus":
		DeleteDBAASExternalEndpointPrometheusCmd(client)
	case "get-dbaas-external-endpoint-prometheus":
		GetDBAASExternalEndpointPrometheusCmd(client)
	case "update-dbaas-external-endpoint-prometheus":
		UpdateDBAASExternalEndpointPrometheusCmd(client)
	case "create-dbaas-external-endpoint-prometheus":
		CreateDBAASExternalEndpointPrometheusCmd(client)
	case "delete-dbaas-external-endpoint-rsyslog":
		DeleteDBAASExternalEndpointRsyslogCmd(client)
	case "get-dbaas-external-endpoint-rsyslog":
		GetDBAASExternalEndpointRsyslogCmd(client)
	case "update-dbaas-external-endpoint-rsyslog":
		UpdateDBAASExternalEndpointRsyslogCmd(client)
	case "create-dbaas-external-endpoint-rsyslog":
		CreateDBAASExternalEndpointRsyslogCmd(client)
	case "list-dbaas-external-endpoint-types":
		ListDBAASExternalEndpointTypesCmd(client)
	case "attach-dbaas-service-to-endpoint":
		AttachDBAASServiceToEndpointCmd(client)
	case "detach-dbaas-service-from-endpoint":
		DetachDBAASServiceFromEndpointCmd(client)
	case "list-dbaas-external-endpoints":
		ListDBAASExternalEndpointsCmd(client)
	case "get-dbaas-external-integration-settings-datadog":
		GetDBAASExternalIntegrationSettingsDatadogCmd(client)
	case "update-dbaas-external-integration-settings-datadog":
		UpdateDBAASExternalIntegrationSettingsDatadogCmd(client)
	case "get-dbaas-external-integration":
		GetDBAASExternalIntegrationCmd(client)
	case "list-dbaas-external-integrations":
		ListDBAASExternalIntegrationsCmd(client)
	case "delete-dbaas-service-grafana":
		DeleteDBAASServiceGrafanaCmd(client)
	case "get-dbaas-service-grafana":
		GetDBAASServiceGrafanaCmd(client)
	case "create-dbaas-service-grafana":
		CreateDBAASServiceGrafanaCmd(client)
	case "update-dbaas-service-grafana":
		UpdateDBAASServiceGrafanaCmd(client)
	case "start-dbaas-grafana-maintenance":
		StartDBAASGrafanaMaintenanceCmd(client)
	case "reset-dbaas-grafana-user-password":
		ResetDBAASGrafanaUserPasswordCmd(client)
	case "reveal-dbaas-grafana-user-password":
		RevealDBAASGrafanaUserPasswordCmd(client)
	case "create-dbaas-integration":
		CreateDBAASIntegrationCmd(client)
	case "list-dbaas-integration-settings":
		ListDBAASIntegrationSettingsCmd(client)
	case "list-dbaas-integration-types":
		ListDBAASIntegrationTypesCmd(client)
	case "delete-dbaas-integration":
		DeleteDBAASIntegrationCmd(client)
	case "get-dbaas-integration":
		GetDBAASIntegrationCmd(client)
	case "update-dbaas-integration":
		UpdateDBAASIntegrationCmd(client)
	case "delete-dbaas-service-kafka":
		DeleteDBAASServiceKafkaCmd(client)
	case "get-dbaas-service-kafka":
		GetDBAASServiceKafkaCmd(client)
	case "create-dbaas-service-kafka":
		CreateDBAASServiceKafkaCmd(client)
	case "update-dbaas-service-kafka":
		UpdateDBAASServiceKafkaCmd(client)
	case "get-dbaas-kafka-acl-config":
		GetDBAASKafkaAclConfigCmd(client)
	case "start-dbaas-kafka-maintenance":
		StartDBAASKafkaMaintenanceCmd(client)
	case "create-dbaas-kafka-schema-registry-acl-config":
		CreateDBAASKafkaSchemaRegistryAclConfigCmd(client)
	case "delete-dbaas-kafka-schema-registry-acl-config":
		DeleteDBAASKafkaSchemaRegistryAclConfigCmd(client)
	case "create-dbaas-kafka-topic-acl-config":
		CreateDBAASKafkaTopicAclConfigCmd(client)
	case "delete-dbaas-kafka-topic-acl-config":
		DeleteDBAASKafkaTopicAclConfigCmd(client)
	case "reveal-dbaas-kafka-connect-password":
		RevealDBAASKafkaConnectPasswordCmd(client)
	case "create-dbaas-kafka-user":
		CreateDBAASKafkaUserCmd(client)
	case "delete-dbaas-kafka-user":
		DeleteDBAASKafkaUserCmd(client)
	case "reset-dbaas-kafka-user-password":
		ResetDBAASKafkaUserPasswordCmd(client)
	case "reveal-dbaas-kafka-user-password":
		RevealDBAASKafkaUserPasswordCmd(client)
	case "get-dbaas-migration-status":
		GetDBAASMigrationStatusCmd(client)
	case "delete-dbaas-service-mysql":
		DeleteDBAASServiceMysqlCmd(client)
	case "get-dbaas-service-mysql":
		GetDBAASServiceMysqlCmd(client)
	case "create-dbaas-service-mysql":
		CreateDBAASServiceMysqlCmd(client)
	case "update-dbaas-service-mysql":
		UpdateDBAASServiceMysqlCmd(client)
	case "enable-dbaas-mysql-writes":
		EnableDBAASMysqlWritesCmd(client)
	case "start-dbaas-mysql-maintenance":
		StartDBAASMysqlMaintenanceCmd(client)
	case "stop-dbaas-mysql-migration":
		StopDBAASMysqlMigrationCmd(client)
	case "create-dbaas-mysql-database":
		CreateDBAASMysqlDatabaseCmd(client)
	case "delete-dbaas-mysql-database":
		DeleteDBAASMysqlDatabaseCmd(client)
	case "create-dbaas-mysql-user":
		CreateDBAASMysqlUserCmd(client)
	case "delete-dbaas-mysql-user":
		DeleteDBAASMysqlUserCmd(client)
	case "reset-dbaas-mysql-user-password":
		ResetDBAASMysqlUserPasswordCmd(client)
	case "reveal-dbaas-mysql-user-password":
		RevealDBAASMysqlUserPasswordCmd(client)
	case "delete-dbaas-service-opensearch":
		DeleteDBAASServiceOpensearchCmd(client)
	case "get-dbaas-service-opensearch":
		GetDBAASServiceOpensearchCmd(client)
	case "create-dbaas-service-opensearch":
		CreateDBAASServiceOpensearchCmd(client)
	case "update-dbaas-service-opensearch":
		UpdateDBAASServiceOpensearchCmd(client)
	case "get-dbaas-opensearch-acl-config":
		GetDBAASOpensearchAclConfigCmd(client)
	case "update-dbaas-opensearch-acl-config":
		UpdateDBAASOpensearchAclConfigCmd(client)
	case "start-dbaas-opensearch-maintenance":
		StartDBAASOpensearchMaintenanceCmd(client)
	case "create-dbaas-opensearch-user":
		CreateDBAASOpensearchUserCmd(client)
	case "delete-dbaas-opensearch-user":
		DeleteDBAASOpensearchUserCmd(client)
	case "reset-dbaas-opensearch-user-password":
		ResetDBAASOpensearchUserPasswordCmd(client)
	case "reveal-dbaas-opensearch-user-password":
		RevealDBAASOpensearchUserPasswordCmd(client)
	case "delete-dbaas-service-pg":
		DeleteDBAASServicePGCmd(client)
	case "get-dbaas-service-pg":
		GetDBAASServicePGCmd(client)
	case "create-dbaas-service-pg":
		CreateDBAASServicePGCmd(client)
	case "update-dbaas-service-pg":
		UpdateDBAASServicePGCmd(client)
	case "start-dbaas-pg-maintenance":
		StartDBAASPGMaintenanceCmd(client)
	case "stop-dbaas-pg-migration":
		StopDBAASPGMigrationCmd(client)
	case "create-dbaas-pg-connection-pool":
		CreateDBAASPGConnectionPoolCmd(client)
	case "delete-dbaas-pg-connection-pool":
		DeleteDBAASPGConnectionPoolCmd(client)
	case "update-dbaas-pg-connection-pool":
		UpdateDBAASPGConnectionPoolCmd(client)
	case "create-dbaas-pg-database":
		CreateDBAASPGDatabaseCmd(client)
	case "delete-dbaas-pg-database":
		DeleteDBAASPGDatabaseCmd(client)
	case "create-dbaas-postgres-user":
		CreateDBAASPostgresUserCmd(client)
	case "delete-dbaas-postgres-user":
		DeleteDBAASPostgresUserCmd(client)
	case "update-dbaas-postgres-allow-replication":
		UpdateDBAASPostgresAllowReplicationCmd(client)
	case "reset-dbaas-postgres-user-password":
		ResetDBAASPostgresUserPasswordCmd(client)
	case "reveal-dbaas-postgres-user-password":
		RevealDBAASPostgresUserPasswordCmd(client)
	case "create-dbaas-pg-upgrade-check":
		CreateDBAASPGUpgradeCheckCmd(client)
	case "delete-dbaas-service-redis":
		DeleteDBAASServiceRedisCmd(client)
	case "get-dbaas-service-redis":
		GetDBAASServiceRedisCmd(client)
	case "create-dbaas-service-redis":
		CreateDBAASServiceRedisCmd(client)
	case "update-dbaas-service-redis":
		UpdateDBAASServiceRedisCmd(client)
	case "start-dbaas-redis-maintenance":
		StartDBAASRedisMaintenanceCmd(client)
	case "stop-dbaas-redis-migration":
		StopDBAASRedisMigrationCmd(client)
	case "start-dbaas-redis-to-valkey-upgrade":
		StartDBAASRedisToValkeyUpgradeCmd(client)
	case "create-dbaas-redis-user":
		CreateDBAASRedisUserCmd(client)
	case "delete-dbaas-redis-user":
		DeleteDBAASRedisUserCmd(client)
	case "reset-dbaas-redis-user-password":
		ResetDBAASRedisUserPasswordCmd(client)
	case "reveal-dbaas-redis-user-password":
		RevealDBAASRedisUserPasswordCmd(client)
	case "list-dbaas-services":
		ListDBAASServicesCmd(client)
	case "get-dbaas-service-logs":
		GetDBAASServiceLogsCmd(client)
	case "get-dbaas-service-metrics":
		GetDBAASServiceMetricsCmd(client)
	case "list-dbaas-service-types":
		ListDBAASServiceTypesCmd(client)
	case "get-dbaas-service-type":
		GetDBAASServiceTypeCmd(client)
	case "delete-dbaas-service":
		DeleteDBAASServiceCmd(client)
	case "get-dbaas-settings-grafana":
		GetDBAASSettingsGrafanaCmd(client)
	case "get-dbaas-settings-kafka":
		GetDBAASSettingsKafkaCmd(client)
	case "get-dbaas-settings-mysql":
		GetDBAASSettingsMysqlCmd(client)
	case "get-dbaas-settings-opensearch":
		GetDBAASSettingsOpensearchCmd(client)
	case "get-dbaas-settings-pg":
		GetDBAASSettingsPGCmd(client)
	case "get-dbaas-settings-redis":
		GetDBAASSettingsRedisCmd(client)
	case "get-dbaas-settings-valkey":
		GetDBAASSettingsValkeyCmd(client)
	case "create-dbaas-task-migration-check":
		CreateDBAASTaskMigrationCheckCmd(client)
	case "get-dbaas-task":
		GetDBAASTaskCmd(client)
	case "delete-dbaas-service-valkey":
		DeleteDBAASServiceValkeyCmd(client)
	case "get-dbaas-service-valkey":
		GetDBAASServiceValkeyCmd(client)
	case "create-dbaas-service-valkey":
		CreateDBAASServiceValkeyCmd(client)
	case "update-dbaas-service-valkey":
		UpdateDBAASServiceValkeyCmd(client)
	case "start-dbaas-valkey-maintenance":
		StartDBAASValkeyMaintenanceCmd(client)
	case "stop-dbaas-valkey-migration":
		StopDBAASValkeyMigrationCmd(client)
	case "create-dbaas-valkey-user":
		CreateDBAASValkeyUserCmd(client)
	case "delete-dbaas-valkey-user":
		DeleteDBAASValkeyUserCmd(client)
	case "reset-dbaas-valkey-user-password":
		ResetDBAASValkeyUserPasswordCmd(client)
	case "reveal-dbaas-valkey-user-password":
		RevealDBAASValkeyUserPasswordCmd(client)
	case "list-deploy-targets":
		ListDeployTargetsCmd(client)
	case "get-deploy-target":
		GetDeployTargetCmd(client)
	case "list-dns-domains":
		ListDNSDomainsCmd(client)
	case "create-dns-domain":
		CreateDNSDomainCmd(client)
	case "list-dns-domain-records":
		ListDNSDomainRecordsCmd(client)
	case "create-dns-domain-record":
		CreateDNSDomainRecordCmd(client)
	case "delete-dns-domain-record":
		DeleteDNSDomainRecordCmd(client)
	case "get-dns-domain-record":
		GetDNSDomainRecordCmd(client)
	case "update-dns-domain-record":
		UpdateDNSDomainRecordCmd(client)
	case "delete-dns-domain":
		DeleteDNSDomainCmd(client)
	case "get-dns-domain":
		GetDNSDomainCmd(client)
	case "get-dns-domain-zone-file":
		GetDNSDomainZoneFileCmd(client)
	case "list-elastic-ips":
		ListElasticIPSCmd(client)
	case "create-elastic-ip":
		CreateElasticIPCmd(client)
	case "delete-elastic-ip":
		DeleteElasticIPCmd(client)
	case "get-elastic-ip":
		GetElasticIPCmd(client)
	case "update-elastic-ip":
		UpdateElasticIPCmd(client)
	case "reset-elastic-ip-field":
		ResetElasticIPFieldCmd(client)
	case "attach-instance-to-elastic-ip":
		AttachInstanceToElasticIPCmd(client)
	case "detach-instance-from-elastic-ip":
		DetachInstanceFromElasticIPCmd(client)
	case "get-env-impact":
		GetEnvImpactCmd(client)
	case "list-events":
		ListEventsCmd(client)
	case "get-iam-organization-policy":
		GetIAMOrganizationPolicyCmd(client)
	case "update-iam-organization-policy":
		UpdateIAMOrganizationPolicyCmd(client)
	case "reset-iam-organization-policy":
		ResetIAMOrganizationPolicyCmd(client)
	case "list-iam-roles":
		ListIAMRolesCmd(client)
	case "create-iam-role":
		CreateIAMRoleCmd(client)
	case "delete-iam-role":
		DeleteIAMRoleCmd(client)
	case "get-iam-role":
		GetIAMRoleCmd(client)
	case "update-iam-role":
		UpdateIAMRoleCmd(client)
	case "update-iam-role-policy":
		UpdateIAMRolePolicyCmd(client)
	case "list-instances":
		ListInstancesCmd(client)
	case "create-instance":
		CreateInstanceCmd(client)
	case "list-instance-pools":
		ListInstancePoolsCmd(client)
	case "create-instance-pool":
		CreateInstancePoolCmd(client)
	case "delete-instance-pool":
		DeleteInstancePoolCmd(client)
	case "get-instance-pool":
		GetInstancePoolCmd(client)
	case "update-instance-pool":
		UpdateInstancePoolCmd(client)
	case "reset-instance-pool-field":
		ResetInstancePoolFieldCmd(client)
	case "evict-instance-pool-members":
		EvictInstancePoolMembersCmd(client)
	case "scale-instance-pool":
		ScaleInstancePoolCmd(client)
	case "list-instance-types":
		ListInstanceTypesCmd(client)
	case "get-instance-type":
		GetInstanceTypeCmd(client)
	case "delete-instance":
		DeleteInstanceCmd(client)
	case "get-instance":
		GetInstanceCmd(client)
	case "update-instance":
		UpdateInstanceCmd(client)
	case "reset-instance-field":
		ResetInstanceFieldCmd(client)
	case "add-instance-protection":
		AddInstanceProtectionCmd(client)
	case "create-snapshot":
		CreateSnapshotCmd(client)
	case "enable-tpm":
		EnableTpmCmd(client)
	case "reveal-instance-password":
		RevealInstancePasswordCmd(client)
	case "reboot-instance":
		RebootInstanceCmd(client)
	case "remove-instance-protection":
		RemoveInstanceProtectionCmd(client)
	case "reset-instance":
		ResetInstanceCmd(client)
	case "reset-instance-password":
		ResetInstancePasswordCmd(client)
	case "resize-instance-disk":
		ResizeInstanceDiskCmd(client)
	case "scale-instance":
		ScaleInstanceCmd(client)
	case "start-instance":
		StartInstanceCmd(client)
	case "stop-instance":
		StopInstanceCmd(client)
	case "revert-instance-to-snapshot":
		RevertInstanceToSnapshotCmd(client)
	case "list-load-balancers":
		ListLoadBalancersCmd(client)
	case "create-load-balancer":
		CreateLoadBalancerCmd(client)
	case "delete-load-balancer":
		DeleteLoadBalancerCmd(client)
	case "get-load-balancer":
		GetLoadBalancerCmd(client)
	case "update-load-balancer":
		UpdateLoadBalancerCmd(client)
	case "add-service-to-load-balancer":
		AddServiceToLoadBalancerCmd(client)
	case "delete-load-balancer-service":
		DeleteLoadBalancerServiceCmd(client)
	case "get-load-balancer-service":
		GetLoadBalancerServiceCmd(client)
	case "update-load-balancer-service":
		UpdateLoadBalancerServiceCmd(client)
	case "reset-load-balancer-service-field":
		ResetLoadBalancerServiceFieldCmd(client)
	case "reset-load-balancer-field":
		ResetLoadBalancerFieldCmd(client)
	case "get-operation":
		GetOperationCmd(client)
	case "get-organization":
		GetOrganizationCmd(client)
	case "list-private-networks":
		ListPrivateNetworksCmd(client)
	case "create-private-network":
		CreatePrivateNetworkCmd(client)
	case "delete-private-network":
		DeletePrivateNetworkCmd(client)
	case "get-private-network":
		GetPrivateNetworkCmd(client)
	case "update-private-network":
		UpdatePrivateNetworkCmd(client)
	case "reset-private-network-field":
		ResetPrivateNetworkFieldCmd(client)
	case "attach-instance-to-private-network":
		AttachInstanceToPrivateNetworkCmd(client)
	case "detach-instance-from-private-network":
		DetachInstanceFromPrivateNetworkCmd(client)
	case "update-private-network-instance-ip":
		UpdatePrivateNetworkInstanceIPCmd(client)
	case "list-quotas":
		ListQuotasCmd(client)
	case "get-quota":
		GetQuotaCmd(client)
	case "delete-reverse-dns-elastic-ip":
		DeleteReverseDNSElasticIPCmd(client)
	case "get-reverse-dns-elastic-ip":
		GetReverseDNSElasticIPCmd(client)
	case "update-reverse-dns-elastic-ip":
		UpdateReverseDNSElasticIPCmd(client)
	case "delete-reverse-dns-instance":
		DeleteReverseDNSInstanceCmd(client)
	case "get-reverse-dns-instance":
		GetReverseDNSInstanceCmd(client)
	case "update-reverse-dns-instance":
		UpdateReverseDNSInstanceCmd(client)
	case "list-security-groups":
		ListSecurityGroupsCmd(client)
	case "create-security-group":
		CreateSecurityGroupCmd(client)
	case "delete-security-group":
		DeleteSecurityGroupCmd(client)
	case "get-security-group":
		GetSecurityGroupCmd(client)
	case "add-rule-to-security-group":
		AddRuleToSecurityGroupCmd(client)
	case "delete-rule-from-security-group":
		DeleteRuleFromSecurityGroupCmd(client)
	case "add-external-source-to-security-group":
		AddExternalSourceToSecurityGroupCmd(client)
	case "attach-instance-to-security-group":
		AttachInstanceToSecurityGroupCmd(client)
	case "detach-instance-from-security-group":
		DetachInstanceFromSecurityGroupCmd(client)
	case "remove-external-source-from-security-group":
		RemoveExternalSourceFromSecurityGroupCmd(client)
	case "list-sks-clusters":
		ListSKSClustersCmd(client)
	case "create-sks-cluster":
		CreateSKSClusterCmd(client)
	case "list-sks-cluster-deprecated-resources":
		ListSKSClusterDeprecatedResourcesCmd(client)
	case "generate-sks-cluster-kubeconfig":
		GenerateSKSClusterKubeconfigCmd(client)
	case "list-sks-cluster-versions":
		ListSKSClusterVersionsCmd(client)
	case "delete-sks-cluster":
		DeleteSKSClusterCmd(client)
	case "get-sks-cluster":
		GetSKSClusterCmd(client)
	case "update-sks-cluster":
		UpdateSKSClusterCmd(client)
	case "get-sks-cluster-authority-cert":
		GetSKSClusterAuthorityCertCmd(client)
	case "get-sks-cluster-inspection":
		GetSKSClusterInspectionCmd(client)
	case "create-sks-nodepool":
		CreateSKSNodepoolCmd(client)
	case "delete-sks-nodepool":
		DeleteSKSNodepoolCmd(client)
	case "get-sks-nodepool":
		GetSKSNodepoolCmd(client)
	case "update-sks-nodepool":
		UpdateSKSNodepoolCmd(client)
	case "reset-sks-nodepool-field":
		ResetSKSNodepoolFieldCmd(client)
	case "evict-sks-nodepool-members":
		EvictSKSNodepoolMembersCmd(client)
	case "scale-sks-nodepool":
		ScaleSKSNodepoolCmd(client)
	case "rotate-sks-ccm-credentials":
		RotateSKSCcmCredentialsCmd(client)
	case "rotate-sks-csi-credentials":
		RotateSKSCsiCredentialsCmd(client)
	case "rotate-sks-operators-ca":
		RotateSKSOperatorsCACmd(client)
	case "upgrade-sks-cluster":
		UpgradeSKSClusterCmd(client)
	case "upgrade-sks-cluster-service-level":
		UpgradeSKSClusterServiceLevelCmd(client)
	case "reset-sks-cluster-field":
		ResetSKSClusterFieldCmd(client)
	case "list-snapshots":
		ListSnapshotsCmd(client)
	case "delete-snapshot":
		DeleteSnapshotCmd(client)
	case "get-snapshot":
		GetSnapshotCmd(client)
	case "export-snapshot":
		ExportSnapshotCmd(client)
	case "promote-snapshot-to-template":
		PromoteSnapshotToTemplateCmd(client)
	case "list-sos-buckets-usage":
		ListSOSBucketsUsageCmd(client)
	case "get-sos-presigned-url":
		GetSOSPresignedURLCmd(client)
	case "list-ssh-keys":
		ListSSHKeysCmd(client)
	case "register-ssh-key":
		RegisterSSHKeyCmd(client)
	case "delete-ssh-key":
		DeleteSSHKeyCmd(client)
	case "get-ssh-key":
		GetSSHKeyCmd(client)
	case "list-templates":
		ListTemplatesCmd(client)
	case "register-template":
		RegisterTemplateCmd(client)
	case "delete-template":
		DeleteTemplateCmd(client)
	case "get-template":
		GetTemplateCmd(client)
	case "copy-template":
		CopyTemplateCmd(client)
	case "update-template":
		UpdateTemplateCmd(client)
	case "get-usage-report":
		GetUsageReportCmd(client)
	case "list-users":
		ListUsersCmd(client)
	case "create-user":
		CreateUserCmd(client)
	case "delete-user":
		DeleteUserCmd(client)
	case "update-user-role":
		UpdateUserRoleCmd(client)
	case "list-zones":
		ListZonesCmd(client)
	default:
		fmt.Println("unknown command:", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func ListAntiAffinityGroupsCmd(client *v3.Client) {
	resp, err := client.ListAntiAffinityGroups(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateAntiAffinityGroupCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-anti-affinity-group", flag.ExitOnError)
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Anti-affinity Group description")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Anti-affinity Group name")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateAntiAffinityGroupRequest
	req.Name = reqNameFlag
	req.Description = reqDescriptionFlag

	resp, err := client.CreateAntiAffinityGroup(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteAntiAffinityGroupCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-anti-affinity-group", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteAntiAffinityGroup(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetAntiAffinityGroupCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-anti-affinity-group", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetAntiAffinityGroup(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListAPIKeysCmd(client *v3.Client) {
	resp, err := client.ListAPIKeys(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateAPIKeyCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-api-key", flag.ExitOnError)
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "IAM API Key Name")
	var reqRoleIDFlag string
	flagset.StringVar(&reqRoleIDFlag, "role-id", "", "IAM API Key Role ID")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateAPIKeyRequest
	req.RoleID = reqRoleIDFlag
	req.Name = reqNameFlag

	resp, err := client.CreateAPIKey(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteAPIKeyCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-api-key", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteAPIKey(context.Background(), idFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetAPIKeyCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-api-key", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetAPIKey(context.Background(), idFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListBlockStorageVolumesCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("list-block-storage-volumes", flag.ExitOnError)
	var instanceIDFlag string
	flagset.StringVar(&instanceIDFlag, "InstanceID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ListBlockStorageVolumes(context.Background(), v3.UUID(instanceIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateBlockStorageVolumeCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-block-storage-volume", flag.ExitOnError)
	var reqBlockStorageSnapshotIDFlag string
	flagset.StringVar(&reqBlockStorageSnapshotIDFlag, "block-storage-snapshot.id", "", "Block storage snapshot ID")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Volume name")
	var reqSizeFlag int64
	flagset.Int64Var(&reqSizeFlag, "size", 0, "Volume size in GiB.                             When a snapshot ID is supplied, this defaults to the size of the source volume, but can be set to a larger value.")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateBlockStorageVolumeRequest
	req.Size = reqSizeFlag
	req.Name = reqNameFlag
	if v3.UUID(reqBlockStorageSnapshotIDFlag) != "" {
		req.BlockStorageSnapshot = &v3.BlockStorageSnapshotTarget{}
		req.BlockStorageSnapshot.ID = v3.UUID(reqBlockStorageSnapshotIDFlag)
	}

	resp, err := client.CreateBlockStorageVolume(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListBlockStorageSnapshotsCmd(client *v3.Client) {
	resp, err := client.ListBlockStorageSnapshots(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteBlockStorageSnapshotCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-block-storage-snapshot", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteBlockStorageSnapshot(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetBlockStorageSnapshotCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-block-storage-snapshot", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetBlockStorageSnapshot(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateBlockStorageSnapshotCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-block-storage-snapshot", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Snapshot name")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateBlockStorageSnapshotRequest
	req.Name = reqNameFlag

	resp, err := client.UpdateBlockStorageSnapshot(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteBlockStorageVolumeCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-block-storage-volume", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteBlockStorageVolume(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetBlockStorageVolumeCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-block-storage-volume", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetBlockStorageVolume(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateBlockStorageVolumeCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-block-storage-volume", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Volume name")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateBlockStorageVolumeRequest
	req.Name = reqNameFlag

	resp, err := client.UpdateBlockStorageVolume(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func AttachBlockStorageVolumeToInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("attach-block-storage-volume-to-instance", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqInstanceIDFlag string
	flagset.StringVar(&reqInstanceIDFlag, "instance.id", "", "Instance ID")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.AttachBlockStorageVolumeToInstanceRequest
	if v3.UUID(reqInstanceIDFlag) != "" {
		req.Instance = &v3.InstanceTarget{}
		req.Instance.ID = v3.UUID(reqInstanceIDFlag)
	}

	resp, err := client.AttachBlockStorageVolumeToInstance(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateBlockStorageSnapshotCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-block-storage-snapshot", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Snapshot name")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateBlockStorageSnapshotRequest
	req.Name = reqNameFlag

	resp, err := client.CreateBlockStorageSnapshot(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DetachBlockStorageVolumeCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("detach-block-storage-volume", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DetachBlockStorageVolume(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResizeBlockStorageVolumeCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("resize-block-storage-volume", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqSizeFlag int64
	flagset.Int64Var(&reqSizeFlag, "size", 0, "Volume size in GiB")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ResizeBlockStorageVolumeRequest
	req.Size = reqSizeFlag

	resp, err := client.ResizeBlockStorageVolume(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetConsoleProxyURLCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-console-proxy-url", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetConsoleProxyURL(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASCACertificateCmd(client *v3.Client) {
	resp, err := client.GetDBAASCACertificate(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASExternalEndpointDatadogCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-external-endpoint-datadog", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASExternalEndpointDatadog(context.Background(), v3.UUID(endpointIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASExternalEndpointDatadogCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-external-endpoint-datadog", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASExternalEndpointDatadog(context.Background(), v3.UUID(endpointIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASExternalEndpointDatadogCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-external-endpoint-datadog", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")
	var reqSettingsDatadogAPIKeyFlag string
	flagset.StringVar(&reqSettingsDatadogAPIKeyFlag, "settings.datadog-api-key", "", "Datadog API key")
	var reqSettingsDisableConsumerStatsFlag bool
	flagset.BoolVar(&reqSettingsDisableConsumerStatsFlag, "settings.disable-consumer-stats", false, "Disable kafka consumer group metrics. Applies only when attached to kafka services.")
	var reqSettingsKafkaConsumerCheckInstancesFlag int64
	flagset.Int64Var(&reqSettingsKafkaConsumerCheckInstancesFlag, "settings.kafka-consumer-check-instances", 0, "Number of separate instances to fetch kafka consumer statistics with. Applies only when attached to kafka services.")
	var reqSettingsKafkaConsumerStatsTimeoutFlag int64
	flagset.Int64Var(&reqSettingsKafkaConsumerStatsTimeoutFlag, "settings.kafka-consumer-stats-timeout", 0, "Number of seconds that datadog will wait to get consumer statistics from brokers. Applies only when attached to kafka services.")
	var reqSettingsMaxPartitionContextsFlag int64
	flagset.Int64Var(&reqSettingsMaxPartitionContextsFlag, "settings.max-partition-contexts", 0, "Maximum number of partition contexts to send. Applies only when attached to kafka services.")
	var reqSettingsSiteFlag string
	flagset.StringVar(&reqSettingsSiteFlag, "settings.site", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASExternalEndpointDatadogRequest
	if reqSettingsSiteFlag != "" {
		req.Settings = &v3.UpdateDBAASExternalEndpointDatadogRequestSettings{}
		req.Settings.Site = reqSettingsSiteFlag
	}
	if reqSettingsMaxPartitionContextsFlag != 0 {

		req.Settings.MaxPartitionContexts = reqSettingsMaxPartitionContextsFlag
	}
	if reqSettingsKafkaConsumerStatsTimeoutFlag != 0 {

		req.Settings.KafkaConsumerStatsTimeout = reqSettingsKafkaConsumerStatsTimeoutFlag
	}
	if reqSettingsKafkaConsumerCheckInstancesFlag != 0 {

		req.Settings.KafkaConsumerCheckInstances = reqSettingsKafkaConsumerCheckInstancesFlag
	}
	if reqSettingsDisableConsumerStatsFlag != false {

		req.Settings.DisableConsumerStats = reqSettingsDisableConsumerStatsFlag
	}
	if reqSettingsDatadogAPIKeyFlag != "" {

		req.Settings.DatadogAPIKey = reqSettingsDatadogAPIKeyFlag
	}

	resp, err := client.UpdateDBAASExternalEndpointDatadog(context.Background(), v3.UUID(endpointIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASExternalEndpointDatadogCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-external-endpoint-datadog", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqSettingsDatadogAPIKeyFlag string
	flagset.StringVar(&reqSettingsDatadogAPIKeyFlag, "settings.datadog-api-key", "", "Datadog API key")
	var reqSettingsDisableConsumerStatsFlag bool
	flagset.BoolVar(&reqSettingsDisableConsumerStatsFlag, "settings.disable-consumer-stats", false, "Disable kafka consumer group metrics. Applies only when attached to kafka services.")
	var reqSettingsKafkaConsumerCheckInstancesFlag int64
	flagset.Int64Var(&reqSettingsKafkaConsumerCheckInstancesFlag, "settings.kafka-consumer-check-instances", 0, "Number of separate instances to fetch kafka consumer statistics with. Applies only when attached to kafka services.")
	var reqSettingsKafkaConsumerStatsTimeoutFlag int64
	flagset.Int64Var(&reqSettingsKafkaConsumerStatsTimeoutFlag, "settings.kafka-consumer-stats-timeout", 0, "Number of seconds that datadog will wait to get consumer statistics from brokers. Applies only when attached to kafka services.")
	var reqSettingsMaxPartitionContextsFlag int64
	flagset.Int64Var(&reqSettingsMaxPartitionContextsFlag, "settings.max-partition-contexts", 0, "Maximum number of partition contexts to send. Applies only when attached to kafka services.")
	var reqSettingsSiteFlag string
	flagset.StringVar(&reqSettingsSiteFlag, "settings.site", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASExternalEndpointDatadogRequest
	if reqSettingsSiteFlag != "" {
		req.Settings = &v3.CreateDBAASExternalEndpointDatadogRequestSettings{}
		req.Settings.Site = reqSettingsSiteFlag
	}
	if reqSettingsMaxPartitionContextsFlag != 0 {

		req.Settings.MaxPartitionContexts = reqSettingsMaxPartitionContextsFlag
	}
	if reqSettingsKafkaConsumerStatsTimeoutFlag != 0 {

		req.Settings.KafkaConsumerStatsTimeout = reqSettingsKafkaConsumerStatsTimeoutFlag
	}
	if reqSettingsKafkaConsumerCheckInstancesFlag != 0 {

		req.Settings.KafkaConsumerCheckInstances = reqSettingsKafkaConsumerCheckInstancesFlag
	}
	if reqSettingsDisableConsumerStatsFlag != false {

		req.Settings.DisableConsumerStats = reqSettingsDisableConsumerStatsFlag
	}
	if reqSettingsDatadogAPIKeyFlag != "" {

		req.Settings.DatadogAPIKey = reqSettingsDatadogAPIKeyFlag
	}

	resp, err := client.CreateDBAASExternalEndpointDatadog(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASExternalEndpointElasticsearchCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-external-endpoint-elasticsearch", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASExternalEndpointElasticsearch(context.Background(), v3.UUID(endpointIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASExternalEndpointElasticsearchCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-external-endpoint-elasticsearch", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASExternalEndpointElasticsearch(context.Background(), v3.UUID(endpointIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASExternalEndpointElasticsearchCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-external-endpoint-elasticsearch", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")
	var reqSettingsCAFlag string
	flagset.StringVar(&reqSettingsCAFlag, "settings.ca", "", "PEM encoded CA certificate")
	var reqSettingsIndexDaysMaxFlag int64
	flagset.Int64Var(&reqSettingsIndexDaysMaxFlag, "settings.index-days-max", 0, "Maximum number of days of logs to keep")
	var reqSettingsIndexPrefixFlag string
	flagset.StringVar(&reqSettingsIndexPrefixFlag, "settings.index-prefix", "", "Elasticsearch index prefix")
	var reqSettingsTimeoutFlag int64
	flagset.Int64Var(&reqSettingsTimeoutFlag, "settings.timeout", 0, "Elasticsearch request timeout limit")
	var reqSettingsURLFlag string
	flagset.StringVar(&reqSettingsURLFlag, "settings.url", "", "Elasticsearch connection URL")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASExternalEndpointElasticsearchRequest
	if reqSettingsURLFlag != "" {
		req.Settings = &v3.UpdateDBAASExternalEndpointElasticsearchRequestSettings{}
		req.Settings.URL = reqSettingsURLFlag
	}
	if reqSettingsTimeoutFlag != 0 {

		req.Settings.Timeout = reqSettingsTimeoutFlag
	}
	if reqSettingsIndexPrefixFlag != "" {

		req.Settings.IndexPrefix = reqSettingsIndexPrefixFlag
	}
	if reqSettingsIndexDaysMaxFlag != 0 {

		req.Settings.IndexDaysMax = reqSettingsIndexDaysMaxFlag
	}
	if reqSettingsCAFlag != "" {

		req.Settings.CA = reqSettingsCAFlag
	}

	resp, err := client.UpdateDBAASExternalEndpointElasticsearch(context.Background(), v3.UUID(endpointIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASExternalEndpointElasticsearchCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-external-endpoint-elasticsearch", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqSettingsCAFlag string
	flagset.StringVar(&reqSettingsCAFlag, "settings.ca", "", "PEM encoded CA certificate")
	var reqSettingsIndexDaysMaxFlag int64
	flagset.Int64Var(&reqSettingsIndexDaysMaxFlag, "settings.index-days-max", 0, "Maximum number of days of logs to keep")
	var reqSettingsIndexPrefixFlag string
	flagset.StringVar(&reqSettingsIndexPrefixFlag, "settings.index-prefix", "", "Elasticsearch index prefix")
	var reqSettingsTimeoutFlag int64
	flagset.Int64Var(&reqSettingsTimeoutFlag, "settings.timeout", 0, "Elasticsearch request timeout limit")
	var reqSettingsURLFlag string
	flagset.StringVar(&reqSettingsURLFlag, "settings.url", "", "Elasticsearch connection URL")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASExternalEndpointElasticsearchRequest
	if reqSettingsURLFlag != "" {
		req.Settings = &v3.CreateDBAASExternalEndpointElasticsearchRequestSettings{}
		req.Settings.URL = reqSettingsURLFlag
	}
	if reqSettingsTimeoutFlag != 0 {

		req.Settings.Timeout = reqSettingsTimeoutFlag
	}
	if reqSettingsIndexPrefixFlag != "" {

		req.Settings.IndexPrefix = reqSettingsIndexPrefixFlag
	}
	if reqSettingsIndexDaysMaxFlag != 0 {

		req.Settings.IndexDaysMax = reqSettingsIndexDaysMaxFlag
	}
	if reqSettingsCAFlag != "" {

		req.Settings.CA = reqSettingsCAFlag
	}

	resp, err := client.CreateDBAASExternalEndpointElasticsearch(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASExternalEndpointOpensearchCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-external-endpoint-opensearch", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASExternalEndpointOpensearch(context.Background(), v3.UUID(endpointIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASExternalEndpointOpensearchCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-external-endpoint-opensearch", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASExternalEndpointOpensearch(context.Background(), v3.UUID(endpointIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASExternalEndpointOpensearchCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-external-endpoint-opensearch", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")
	var reqSettingsCAFlag string
	flagset.StringVar(&reqSettingsCAFlag, "settings.ca", "", "PEM encoded CA certificate")
	var reqSettingsIndexDaysMaxFlag int64
	flagset.Int64Var(&reqSettingsIndexDaysMaxFlag, "settings.index-days-max", 0, "Maximum number of days of logs to keep")
	var reqSettingsIndexPrefixFlag string
	flagset.StringVar(&reqSettingsIndexPrefixFlag, "settings.index-prefix", "", "OpenSearch index prefix")
	var reqSettingsTimeoutFlag int64
	flagset.Int64Var(&reqSettingsTimeoutFlag, "settings.timeout", 0, "OpenSearch request timeout limit")
	var reqSettingsURLFlag string
	flagset.StringVar(&reqSettingsURLFlag, "settings.url", "", "OpenSearch connection URL")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASExternalEndpointOpensearchRequest
	if reqSettingsURLFlag != "" {
		req.Settings = &v3.UpdateDBAASExternalEndpointOpensearchRequestSettings{}
		req.Settings.URL = reqSettingsURLFlag
	}
	if reqSettingsTimeoutFlag != 0 {

		req.Settings.Timeout = reqSettingsTimeoutFlag
	}
	if reqSettingsIndexPrefixFlag != "" {

		req.Settings.IndexPrefix = reqSettingsIndexPrefixFlag
	}
	if reqSettingsIndexDaysMaxFlag != 0 {

		req.Settings.IndexDaysMax = reqSettingsIndexDaysMaxFlag
	}
	if reqSettingsCAFlag != "" {

		req.Settings.CA = reqSettingsCAFlag
	}

	resp, err := client.UpdateDBAASExternalEndpointOpensearch(context.Background(), v3.UUID(endpointIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASExternalEndpointOpensearchCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-external-endpoint-opensearch", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqSettingsCAFlag string
	flagset.StringVar(&reqSettingsCAFlag, "settings.ca", "", "PEM encoded CA certificate")
	var reqSettingsIndexDaysMaxFlag int64
	flagset.Int64Var(&reqSettingsIndexDaysMaxFlag, "settings.index-days-max", 0, "Maximum number of days of logs to keep")
	var reqSettingsIndexPrefixFlag string
	flagset.StringVar(&reqSettingsIndexPrefixFlag, "settings.index-prefix", "", "OpenSearch index prefix")
	var reqSettingsTimeoutFlag int64
	flagset.Int64Var(&reqSettingsTimeoutFlag, "settings.timeout", 0, "OpenSearch request timeout limit")
	var reqSettingsURLFlag string
	flagset.StringVar(&reqSettingsURLFlag, "settings.url", "", "OpenSearch connection URL")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASExternalEndpointOpensearchRequest
	if reqSettingsURLFlag != "" {
		req.Settings = &v3.CreateDBAASExternalEndpointOpensearchRequestSettings{}
		req.Settings.URL = reqSettingsURLFlag
	}
	if reqSettingsTimeoutFlag != 0 {

		req.Settings.Timeout = reqSettingsTimeoutFlag
	}
	if reqSettingsIndexPrefixFlag != "" {

		req.Settings.IndexPrefix = reqSettingsIndexPrefixFlag
	}
	if reqSettingsIndexDaysMaxFlag != 0 {

		req.Settings.IndexDaysMax = reqSettingsIndexDaysMaxFlag
	}
	if reqSettingsCAFlag != "" {

		req.Settings.CA = reqSettingsCAFlag
	}

	resp, err := client.CreateDBAASExternalEndpointOpensearch(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASExternalEndpointPrometheusCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-external-endpoint-prometheus", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASExternalEndpointPrometheus(context.Background(), v3.UUID(endpointIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASExternalEndpointPrometheusCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-external-endpoint-prometheus", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASExternalEndpointPrometheus(context.Background(), v3.UUID(endpointIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASExternalEndpointPrometheusCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-external-endpoint-prometheus", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")
	var reqSettingsBasicAuthPasswordFlag string
	flagset.StringVar(&reqSettingsBasicAuthPasswordFlag, "settings.basic-auth-password", "", "Prometheus basic authentication password")
	var reqSettingsBasicAuthUsernameFlag string
	flagset.StringVar(&reqSettingsBasicAuthUsernameFlag, "settings.basic-auth-username", "", "Prometheus basic authentication username")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASExternalEndpointPrometheusRequest
	if reqSettingsBasicAuthUsernameFlag != "" {
		req.Setting = &v3.UpdateDBAASExternalEndpointPrometheusRequestSettings{}
		req.Settings.BasicAuthUsername = reqSettingsBasicAuthUsernameFlag
	}
	if reqSettingsBasicAuthPasswordFlag != "" {

		req.Settings.BasicAuthPassword = reqSettingsBasicAuthPasswordFlag
	}

	resp, err := client.UpdateDBAASExternalEndpointPrometheus(context.Background(), v3.UUID(endpointIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASExternalEndpointPrometheusCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-external-endpoint-prometheus", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqSettingsBasicAuthPasswordFlag string
	flagset.StringVar(&reqSettingsBasicAuthPasswordFlag, "settings.basic-auth-password", "", "Prometheus basic authentication password")
	var reqSettingsBasicAuthUsernameFlag string
	flagset.StringVar(&reqSettingsBasicAuthUsernameFlag, "settings.basic-auth-username", "", "Prometheus basic authentication username")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASExternalEndpointPrometheusRequest
	if reqSettingsBasicAuthUsernameFlag != "" {
		req.Setting = &v3.CreateDBAASExternalEndpointPrometheusRequestSettings{}
		req.Settings.BasicAuthUsername = reqSettingsBasicAuthUsernameFlag
	}
	if reqSettingsBasicAuthPasswordFlag != "" {

		req.Settings.BasicAuthPassword = reqSettingsBasicAuthPasswordFlag
	}

	resp, err := client.CreateDBAASExternalEndpointPrometheus(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASExternalEndpointRsyslogCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-external-endpoint-rsyslog", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASExternalEndpointRsyslog(context.Background(), v3.UUID(endpointIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASExternalEndpointRsyslogCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-external-endpoint-rsyslog", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASExternalEndpointRsyslog(context.Background(), v3.UUID(endpointIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASExternalEndpointRsyslogCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-external-endpoint-rsyslog", flag.ExitOnError)
	var endpointIDFlag string
	flagset.StringVar(&endpointIDFlag, "EndpointID", "", "")
	var reqSettingsCAFlag string
	flagset.StringVar(&reqSettingsCAFlag, "settings.ca", "", "PEM encoded CA certificate")
	var reqSettingsCertFlag string
	flagset.StringVar(&reqSettingsCertFlag, "settings.cert", "", "PEM encoded client certificate")
	var reqSettingsFormatFlag string
	flagset.StringVar(&reqSettingsFormatFlag, "settings.format", "", "")
	var reqSettingsKeyFlag string
	flagset.StringVar(&reqSettingsKeyFlag, "settings.key", "", "PEM encoded client key")
	var reqSettingsLoglineFlag string
	flagset.StringVar(&reqSettingsLoglineFlag, "settings.logline", "", "Custom syslog message format")
	var reqSettingsMaxMessageSizeFlag int64
	flagset.Int64Var(&reqSettingsMaxMessageSizeFlag, "settings.max-message-size", 0, "Rsyslog max message size")
	var reqSettingsPortFlag int64
	flagset.Int64Var(&reqSettingsPortFlag, "settings.port", 0, "Rsyslog server port")
	var reqSettingsSDFlag string
	flagset.StringVar(&reqSettingsSDFlag, "settings.sd", "", "Structured data block for log message")
	var reqSettingsServerFlag string
	flagset.StringVar(&reqSettingsServerFlag, "settings.server", "", "Rsyslog server IP address or hostname")
	var reqSettingsTlsFlag bool
	flagset.BoolVar(&reqSettingsTlsFlag, "settings.tls", false, "Require TLS")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASExternalEndpointRsyslogRequest
	if reqSettingsTlsFlag != false {
		req.Setting = &v3.UpdateDBAASExternalEndpointRsyslogRequestSettings{}
		req.Settings.Tls = reqSettingsTlsFlag
	}
	if reqSettingsServerFlag != "" {

		req.Settings.Server = reqSettingsServerFlag
	}
	if reqSettingsSDFlag != "" {

		req.Settings.SD = reqSettingsSDFlag
	}
	if reqSettingsPortFlag != 0 {

		req.Settings.Port = reqSettingsPortFlag
	}
	if reqSettingsMaxMessageSizeFlag != 0 {

		req.Settings.MaxMessageSize = reqSettingsMaxMessageSizeFlag
	}
	if reqSettingsLoglineFlag != "" {

		req.Settings.Logline = reqSettingsLoglineFlag
	}
	if reqSettingsKeyFlag != "" {

		req.Settings.Key = reqSettingsKeyFlag
	}
	if reqSettingsFormatFlag != "" {

		req.Settings.Format = reqSettingsFormatFlag
	}
	if reqSettingsCertFlag != "" {

		req.Settings.Cert = reqSettingsCertFlag
	}
	if reqSettingsCAFlag != "" {

		req.Settings.CA = reqSettingsCAFlag
	}

	resp, err := client.UpdateDBAASExternalEndpointRsyslog(context.Background(), v3.UUID(endpointIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASExternalEndpointRsyslogCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-external-endpoint-rsyslog", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqSettingsCAFlag string
	flagset.StringVar(&reqSettingsCAFlag, "settings.ca", "", "PEM encoded CA certificate")
	var reqSettingsCertFlag string
	flagset.StringVar(&reqSettingsCertFlag, "settings.cert", "", "PEM encoded client certificate")
	var reqSettingsFormatFlag string
	flagset.StringVar(&reqSettingsFormatFlag, "settings.format", "", "")
	var reqSettingsKeyFlag string
	flagset.StringVar(&reqSettingsKeyFlag, "settings.key", "", "PEM encoded client key")
	var reqSettingsLoglineFlag string
	flagset.StringVar(&reqSettingsLoglineFlag, "settings.logline", "", "Custom syslog message format")
	var reqSettingsMaxMessageSizeFlag int64
	flagset.Int64Var(&reqSettingsMaxMessageSizeFlag, "settings.max-message-size", 0, "Rsyslog max message size")
	var reqSettingsPortFlag int64
	flagset.Int64Var(&reqSettingsPortFlag, "settings.port", 0, "Rsyslog server port")
	var reqSettingsSDFlag string
	flagset.StringVar(&reqSettingsSDFlag, "settings.sd", "", "Structured data block for log message")
	var reqSettingsServerFlag string
	flagset.StringVar(&reqSettingsServerFlag, "settings.server", "", "Rsyslog server IP address or hostname")
	var reqSettingsTlsFlag bool
	flagset.BoolVar(&reqSettingsTlsFlag, "settings.tls", false, "Require TLS")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASExternalEndpointRsyslogRequest
	if reqSettingsTlsFlag != false {
		req.Setting = &v3.CreateDBAASExternalEndpointRsyslogRequestSettings{}
		req.Settings.Tls = reqSettingsTlsFlag
	}
	if reqSettingsServerFlag != "" {

		req.Settings.Server = reqSettingsServerFlag
	}
	if reqSettingsSDFlag != "" {

		req.Settings.SD = reqSettingsSDFlag
	}
	if reqSettingsPortFlag != 0 {

		req.Settings.Port = reqSettingsPortFlag
	}
	if reqSettingsMaxMessageSizeFlag != 0 {

		req.Settings.MaxMessageSize = reqSettingsMaxMessageSizeFlag
	}
	if reqSettingsLoglineFlag != "" {

		req.Settings.Logline = reqSettingsLoglineFlag
	}
	if reqSettingsKeyFlag != "" {

		req.Settings.Key = reqSettingsKeyFlag
	}
	if reqSettingsFormatFlag != "" {

		req.Settings.Format = reqSettingsFormatFlag
	}
	if reqSettingsCertFlag != "" {

		req.Settings.Cert = reqSettingsCertFlag
	}
	if reqSettingsCAFlag != "" {

		req.Settings.CA = reqSettingsCAFlag
	}

	resp, err := client.CreateDBAASExternalEndpointRsyslog(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListDBAASExternalEndpointTypesCmd(client *v3.Client) {
	resp, err := client.ListDBAASExternalEndpointTypes(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func AttachDBAASServiceToEndpointCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("attach-dbaas-service-to-endpoint", flag.ExitOnError)
	var sourceServiceNameFlag string
	flagset.StringVar(&sourceServiceNameFlag, "SourceServiceName", "", "")
	var reqDestEndpointIDFlag string
	flagset.StringVar(&reqDestEndpointIDFlag, "dest-endpoint-id", "", "External endpoint id")
	var reqTypeFlag string
	flagset.StringVar(&reqTypeFlag, "type", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.AttachDBAASServiceToEndpointRequest
	req.Type = reqTypeFlag
	req.DestEndpointID = reqDestEndpointIDFlag

	resp, err := client.AttachDBAASServiceToEndpoint(context.Background(), sourceServiceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DetachDBAASServiceFromEndpointCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("detach-dbaas-service-from-endpoint", flag.ExitOnError)
	var sourceServiceNameFlag string
	flagset.StringVar(&sourceServiceNameFlag, "SourceServiceName", "", "")
	var reqIntegrationIDFlag string
	flagset.StringVar(&reqIntegrationIDFlag, "integration-id", "", "External Integration ID")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.DetachDBAASServiceFromEndpointRequest
	req.IntegrationID = reqIntegrationIDFlag

	resp, err := client.DetachDBAASServiceFromEndpoint(context.Background(), sourceServiceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListDBAASExternalEndpointsCmd(client *v3.Client) {
	resp, err := client.ListDBAASExternalEndpoints(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASExternalIntegrationSettingsDatadogCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-external-integration-settings-datadog", flag.ExitOnError)
	var integrationIDFlag string
	flagset.StringVar(&integrationIDFlag, "IntegrationID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASExternalIntegrationSettingsDatadog(context.Background(), v3.UUID(integrationIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASExternalIntegrationSettingsDatadogCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-external-integration-settings-datadog", flag.ExitOnError)
	var integrationIDFlag string
	flagset.StringVar(&integrationIDFlag, "IntegrationID", "", "")
	var reqSettingsDatadogDbmEnabledFlag bool
	flagset.BoolVar(&reqSettingsDatadogDbmEnabledFlag, "settings.datadog-dbm-enabled", false, "Database monitoring: view query metrics, explain plans, and execution details. Correlate queries with host metrics.")
	var reqSettingsDatadogPgbouncerEnabledFlag bool
	flagset.BoolVar(&reqSettingsDatadogPgbouncerEnabledFlag, "settings.datadog-pgbouncer-enabled", false, "Integrate PgBouncer with Datadog to track connection pool metrics and monitor application traffic.")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASExternalIntegrationSettingsDatadogRequest
	if reqSettingsDatadogPgbouncerEnabledFlag != false {
		req.Settings = &v3.DBAASIntegrationSettingsDatadog{}
		req.Settings.DatadogPgbouncerEnabled = reqSettingsDatadogPgbouncerEnabledFlag
	}
	if reqSettingsDatadogDbmEnabledFlag != false {

		req.Settings.DatadogDbmEnabled = reqSettingsDatadogDbmEnabledFlag
	}

	resp, err := client.UpdateDBAASExternalIntegrationSettingsDatadog(context.Background(), v3.UUID(integrationIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASExternalIntegrationCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-external-integration", flag.ExitOnError)
	var integrationIDFlag string
	flagset.StringVar(&integrationIDFlag, "IntegrationID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASExternalIntegration(context.Background(), v3.UUID(integrationIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListDBAASExternalIntegrationsCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("list-dbaas-external-integrations", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ListDBAASExternalIntegrations(context.Background(), serviceNameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASServiceGrafanaCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-service-grafana", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASServiceGrafana(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASServiceGrafanaCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-service-grafana", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASServiceGrafana(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASServiceGrafanaCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-service-grafana", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqForkFromServiceFlag string
	flagset.StringVar(&reqForkFromServiceFlag, "fork-from-service", "", "")
	var reqGrafanaSettingsAlertingEnabledFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAlertingEnabledFlag, "grafana-settings.alerting_enabled", false, "Enable or disable Grafana legacy alerting functionality. This should not be enabled with unified_alerting_enabled.")
	var reqGrafanaSettingsAlertingErrorORTimeoutFlag string
	flagset.StringVar(&reqGrafanaSettingsAlertingErrorORTimeoutFlag, "grafana-settings.alerting_error_or_timeout", "", "Default error or timeout setting for new alerting rules")
	var reqGrafanaSettingsAlertingMaxAnnotationsToKeepFlag int
	flagset.IntVar(&reqGrafanaSettingsAlertingMaxAnnotationsToKeepFlag, "grafana-settings.alerting_max_annotations_to_keep", 0, "Max number of alert annotations that Grafana stores. 0 (default) keeps all alert annotations.")
	var reqGrafanaSettingsAlertingNodataORNullvaluesFlag string
	flagset.StringVar(&reqGrafanaSettingsAlertingNodataORNullvaluesFlag, "grafana-settings.alerting_nodata_or_nullvalues", "", "Default value for 'no data or null values' for new alerting rules")
	var reqGrafanaSettingsAllowEmbeddingFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAllowEmbeddingFlag, "grafana-settings.allow_embedding", false, "Allow embedding Grafana dashboards with iframe/frame/object/embed tags. Disabled by default to limit impact of clickjacking")
	var reqGrafanaSettingsAuthAzureadAllowSignUPFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthAzureadAllowSignUPFlag, "grafana-settings.auth_azuread.allow_sign_up", false, "Automatically sign-up users on successful sign-in")
	var reqGrafanaSettingsAuthAzureadAuthURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthAzureadAuthURLFlag, "grafana-settings.auth_azuread.auth_url", "", "Authorization URL")
	var reqGrafanaSettingsAuthAzureadClientIDFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthAzureadClientIDFlag, "grafana-settings.auth_azuread.client_id", "", "Client ID from provider")
	var reqGrafanaSettingsAuthAzureadClientSecretFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthAzureadClientSecretFlag, "grafana-settings.auth_azuread.client_secret", "", "Client secret from provider")
	var reqGrafanaSettingsAuthAzureadTokenURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthAzureadTokenURLFlag, "grafana-settings.auth_azuread.token_url", "", "Token URL")
	var reqGrafanaSettingsAuthBasicEnabledFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthBasicEnabledFlag, "grafana-settings.auth_basic_enabled", false, "Enable or disable basic authentication form, used by Grafana built-in login")
	var reqGrafanaSettingsAuthGenericOauthAllowSignUPFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGenericOauthAllowSignUPFlag, "grafana-settings.auth_generic_oauth.allow_sign_up", false, "Automatically sign-up users on successful sign-in")
	var reqGrafanaSettingsAuthGenericOauthAPIURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGenericOauthAPIURLFlag, "grafana-settings.auth_generic_oauth.api_url", "", "API URL")
	var reqGrafanaSettingsAuthGenericOauthAuthURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGenericOauthAuthURLFlag, "grafana-settings.auth_generic_oauth.auth_url", "", "Authorization URL")
	var reqGrafanaSettingsAuthGenericOauthAutoLoginFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGenericOauthAutoLoginFlag, "grafana-settings.auth_generic_oauth.auto_login", false, "Allow users to bypass the login screen and automatically log in")
	var reqGrafanaSettingsAuthGenericOauthClientIDFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGenericOauthClientIDFlag, "grafana-settings.auth_generic_oauth.client_id", "", "Client ID from provider")
	var reqGrafanaSettingsAuthGenericOauthClientSecretFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGenericOauthClientSecretFlag, "grafana-settings.auth_generic_oauth.client_secret", "", "Client secret from provider")
	var reqGrafanaSettingsAuthGenericOauthNameFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGenericOauthNameFlag, "grafana-settings.auth_generic_oauth.name", "", "Name of the OAuth integration")
	var reqGrafanaSettingsAuthGenericOauthTokenURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGenericOauthTokenURLFlag, "grafana-settings.auth_generic_oauth.token_url", "", "Token URL")
	var reqGrafanaSettingsAuthGithubAllowSignUPFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGithubAllowSignUPFlag, "grafana-settings.auth_github.allow_sign_up", false, "Automatically sign-up users on successful sign-in")
	var reqGrafanaSettingsAuthGithubAutoLoginFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGithubAutoLoginFlag, "grafana-settings.auth_github.auto_login", false, "Allow users to bypass the login screen and automatically log in")
	var reqGrafanaSettingsAuthGithubClientIDFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGithubClientIDFlag, "grafana-settings.auth_github.client_id", "", "Client ID from provider")
	var reqGrafanaSettingsAuthGithubClientSecretFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGithubClientSecretFlag, "grafana-settings.auth_github.client_secret", "", "Client secret from provider")
	var reqGrafanaSettingsAuthGithubSkipOrgRoleSyncFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGithubSkipOrgRoleSyncFlag, "grafana-settings.auth_github.skip_org_role_sync", false, "Stop automatically syncing user roles")
	var reqGrafanaSettingsAuthGitlabAllowSignUPFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGitlabAllowSignUPFlag, "grafana-settings.auth_gitlab.allow_sign_up", false, "Automatically sign-up users on successful sign-in")
	var reqGrafanaSettingsAuthGitlabAPIURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGitlabAPIURLFlag, "grafana-settings.auth_gitlab.api_url", "", "API URL. This only needs to be set when using self hosted GitLab")
	var reqGrafanaSettingsAuthGitlabAuthURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGitlabAuthURLFlag, "grafana-settings.auth_gitlab.auth_url", "", "Authorization URL. This only needs to be set when using self hosted GitLab")
	var reqGrafanaSettingsAuthGitlabClientIDFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGitlabClientIDFlag, "grafana-settings.auth_gitlab.client_id", "", "Client ID from provider")
	var reqGrafanaSettingsAuthGitlabClientSecretFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGitlabClientSecretFlag, "grafana-settings.auth_gitlab.client_secret", "", "Client secret from provider")
	var reqGrafanaSettingsAuthGitlabTokenURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGitlabTokenURLFlag, "grafana-settings.auth_gitlab.token_url", "", "Token URL. This only needs to be set when using self hosted GitLab")
	var reqGrafanaSettingsAuthGoogleAllowSignUPFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGoogleAllowSignUPFlag, "grafana-settings.auth_google.allow_sign_up", false, "Automatically sign-up users on successful sign-in")
	var reqGrafanaSettingsAuthGoogleClientIDFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGoogleClientIDFlag, "grafana-settings.auth_google.client_id", "", "Client ID from provider")
	var reqGrafanaSettingsAuthGoogleClientSecretFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGoogleClientSecretFlag, "grafana-settings.auth_google.client_secret", "", "Client secret from provider")
	var reqGrafanaSettingsCookieSamesiteFlag string
	flagset.StringVar(&reqGrafanaSettingsCookieSamesiteFlag, "grafana-settings.cookie_samesite", "", "Cookie SameSite attribute: 'strict' prevents sending cookie for cross-site requests, effectively disabling direct linking from other sites to Grafana. 'lax' is the default value.")
	var reqGrafanaSettingsCustomDomainFlag string
	flagset.StringVar(&reqGrafanaSettingsCustomDomainFlag, "grafana-settings.custom_domain", "", "Serve the web frontend using a custom CNAME pointing to the Aiven DNS name")
	var reqGrafanaSettingsDashboardPreviewsEnabledFlag bool
	flagset.BoolVar(&reqGrafanaSettingsDashboardPreviewsEnabledFlag, "grafana-settings.dashboard_previews_enabled", false, "This feature is new in Grafana 9 and is quite resource intensive. It may cause low-end plans to work more slowly while the dashboard previews are rendering.")
	var reqGrafanaSettingsDashboardsMinRefreshIntervalFlag string
	flagset.StringVar(&reqGrafanaSettingsDashboardsMinRefreshIntervalFlag, "grafana-settings.dashboards_min_refresh_interval", "", "Signed sequence of decimal numbers, followed by a unit suffix (ms, s, m, h, d), e.g. 30s, 1h")
	var reqGrafanaSettingsDashboardsVersionsToKeepFlag int
	flagset.IntVar(&reqGrafanaSettingsDashboardsVersionsToKeepFlag, "grafana-settings.dashboards_versions_to_keep", 0, "Dashboard versions to keep per dashboard")
	var reqGrafanaSettingsDataproxySendUserHeaderFlag bool
	flagset.BoolVar(&reqGrafanaSettingsDataproxySendUserHeaderFlag, "grafana-settings.dataproxy_send_user_header", false, "Send 'X-Grafana-User' header to data source")
	var reqGrafanaSettingsDataproxyTimeoutFlag int
	flagset.IntVar(&reqGrafanaSettingsDataproxyTimeoutFlag, "grafana-settings.dataproxy_timeout", 0, "Timeout for data proxy requests in seconds")
	var reqGrafanaSettingsDateFormatsDefaultTimezoneFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsDefaultTimezoneFlag, "grafana-settings.date_formats.default_timezone", "", "Default time zone for user preferences. Value 'browser' uses browser local time zone.")
	var reqGrafanaSettingsDateFormatsFullDateFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsFullDateFlag, "grafana-settings.date_formats.full_date", "", "Moment.js style format string for cases where full date is shown")
	var reqGrafanaSettingsDateFormatsIntervalDayFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsIntervalDayFlag, "grafana-settings.date_formats.interval_day", "", "Moment.js style format string used when a time requiring day accuracy is shown")
	var reqGrafanaSettingsDateFormatsIntervalHourFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsIntervalHourFlag, "grafana-settings.date_formats.interval_hour", "", "Moment.js style format string used when a time requiring hour accuracy is shown")
	var reqGrafanaSettingsDateFormatsIntervalMinuteFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsIntervalMinuteFlag, "grafana-settings.date_formats.interval_minute", "", "Moment.js style format string used when a time requiring minute accuracy is shown")
	var reqGrafanaSettingsDateFormatsIntervalMonthFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsIntervalMonthFlag, "grafana-settings.date_formats.interval_month", "", "Moment.js style format string used when a time requiring month accuracy is shown")
	var reqGrafanaSettingsDateFormatsIntervalSecondFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsIntervalSecondFlag, "grafana-settings.date_formats.interval_second", "", "Moment.js style format string used when a time requiring second accuracy is shown")
	var reqGrafanaSettingsDateFormatsIntervalYearFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsIntervalYearFlag, "grafana-settings.date_formats.interval_year", "", "Moment.js style format string used when a time requiring year accuracy is shown")
	var reqGrafanaSettingsDisableGravatarFlag bool
	flagset.BoolVar(&reqGrafanaSettingsDisableGravatarFlag, "grafana-settings.disable_gravatar", false, "Set to true to disable gravatar. Defaults to false (gravatar is enabled)")
	var reqGrafanaSettingsEditorsCanAdminFlag bool
	flagset.BoolVar(&reqGrafanaSettingsEditorsCanAdminFlag, "grafana-settings.editors_can_admin", false, "Editors can manage folders, teams and dashboards created by them")
	var reqGrafanaSettingsGoogleAnalyticsUAIDFlag string
	flagset.StringVar(&reqGrafanaSettingsGoogleAnalyticsUAIDFlag, "grafana-settings.google_analytics_ua_id", "", "Google Analytics ID")
	var reqGrafanaSettingsMetricsEnabledFlag bool
	flagset.BoolVar(&reqGrafanaSettingsMetricsEnabledFlag, "grafana-settings.metrics_enabled", false, "Enable Grafana /metrics endpoint")
	var reqGrafanaSettingsOauthAllowInsecureEmailLookupFlag bool
	flagset.BoolVar(&reqGrafanaSettingsOauthAllowInsecureEmailLookupFlag, "grafana-settings.oauth_allow_insecure_email_lookup", false, "Enforce user lookup based on email instead of the unique ID provided by the IdP")
	var reqGrafanaSettingsServiceLogFlag bool
	flagset.BoolVar(&reqGrafanaSettingsServiceLogFlag, "grafana-settings.service_log", false, "Store logs for the service so that they are available in the HTTP API and console.")
	var reqGrafanaSettingsSMTPServerFromAddressFlag string
	flagset.StringVar(&reqGrafanaSettingsSMTPServerFromAddressFlag, "grafana-settings.smtp_server.from_address", "", "Address used for sending emails")
	var reqGrafanaSettingsSMTPServerFromNameFlag string
	flagset.StringVar(&reqGrafanaSettingsSMTPServerFromNameFlag, "grafana-settings.smtp_server.from_name", "", "Name used in outgoing emails, defaults to Grafana")
	var reqGrafanaSettingsSMTPServerHostFlag string
	flagset.StringVar(&reqGrafanaSettingsSMTPServerHostFlag, "grafana-settings.smtp_server.host", "", "Server hostname or IP")
	var reqGrafanaSettingsSMTPServerPasswordFlag string
	flagset.StringVar(&reqGrafanaSettingsSMTPServerPasswordFlag, "grafana-settings.smtp_server.password", "", "Password for SMTP authentication")
	var reqGrafanaSettingsSMTPServerPortFlag int
	flagset.IntVar(&reqGrafanaSettingsSMTPServerPortFlag, "grafana-settings.smtp_server.port", 0, "SMTP server port")
	var reqGrafanaSettingsSMTPServerSkipVerifyFlag bool
	flagset.BoolVar(&reqGrafanaSettingsSMTPServerSkipVerifyFlag, "grafana-settings.smtp_server.skip_verify", false, "Skip verifying server certificate. Defaults to false")
	var reqGrafanaSettingsSMTPServerStarttlsPolicyFlag string
	flagset.StringVar(&reqGrafanaSettingsSMTPServerStarttlsPolicyFlag, "grafana-settings.smtp_server.starttls_policy", "", "Either OpportunisticStartTLS, MandatoryStartTLS or NoStartTLS. Default is OpportunisticStartTLS.")
	var reqGrafanaSettingsSMTPServerUsernameFlag string
	flagset.StringVar(&reqGrafanaSettingsSMTPServerUsernameFlag, "grafana-settings.smtp_server.username", "", "Username for SMTP authentication")
	var reqGrafanaSettingsUnifiedAlertingEnabledFlag bool
	flagset.BoolVar(&reqGrafanaSettingsUnifiedAlertingEnabledFlag, "grafana-settings.unified_alerting_enabled", false, "Enable or disable Grafana unified alerting functionality. By default this is enabled and any legacy alerts will be migrated on upgrade to Grafana 9+. To stay on legacy alerting, set unified_alerting_enabled to false and alerting_enabled to true. See https://grafana.com/docs/grafana/latest/alerting/set-up/migrating-alerts/ for more details.")
	var reqGrafanaSettingsUserAutoAssignOrgFlag bool
	flagset.BoolVar(&reqGrafanaSettingsUserAutoAssignOrgFlag, "grafana-settings.user_auto_assign_org", false, "Auto-assign new users on signup to main organization. Defaults to false")
	var reqGrafanaSettingsUserAutoAssignOrgRoleFlag string
	flagset.StringVar(&reqGrafanaSettingsUserAutoAssignOrgRoleFlag, "grafana-settings.user_auto_assign_org_role", "", "Set role for new signups. Defaults to Viewer")
	var reqGrafanaSettingsViewersCanEditFlag bool
	flagset.BoolVar(&reqGrafanaSettingsViewersCanEditFlag, "grafana-settings.viewers_can_edit", false, "Users with view-only permission can edit but not save dashboards")
	var reqGrafanaSettingsWalFlag bool
	flagset.BoolVar(&reqGrafanaSettingsWalFlag, "grafana-settings.wal", false, "Setting to enable/disable Write-Ahead Logging. The default value is false (disabled).")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASServiceGrafanaRequest
	req.TerminationProtection = reqTerminationProtectionFlag
	req.Plan = reqPlanFlag
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.CreateDBAASServiceGrafanaRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}
	if reqGrafanaSettingsWalFlag != false {
		req.GrafanaSettings = &v3.JSONSchemaGrafana{}
		req.GrafanaSettings.Wal = reqGrafanaSettingsWalFlag
	}
	if reqGrafanaSettingsViewersCanEditFlag != false {

		req.GrafanaSettings.ViewersCanEdit = reqGrafanaSettingsViewersCanEditFlag
	}
	if reqGrafanaSettingsUserAutoAssignOrgRoleFlag != "" {

		req.GrafanaSettings.UserAutoAssignOrgRole = reqGrafanaSettingsUserAutoAssignOrgRoleFlag
	}
	if reqGrafanaSettingsUserAutoAssignOrgFlag != false {

		req.GrafanaSettings.UserAutoAssignOrg = reqGrafanaSettingsUserAutoAssignOrgFlag
	}
	if reqGrafanaSettingsUnifiedAlertingEnabledFlag != false {

		req.GrafanaSettings.UnifiedAlertingEnabled = reqGrafanaSettingsUnifiedAlertingEnabledFlag
	}
	if reqGrafanaSettingsSMTPServerUsernameFlag != "" {
		req.GrafanaSettingsSMTPServ = &v3.GrafanaSettingsSMTPServer{}
		req.GrafanaSettingsSMTPServer.Username = reqGrafanaSettingsSMTPServerUsernameFlag
	}
	if reqGrafanaSettingsSMTPServerStarttlsPolicyFlag != "" {

		req.GrafanaSettingsSMTPServer.StarttlsPolicy = reqGrafanaSettingsSMTPServerStarttlsPolicyFlag
	}
	if reqGrafanaSettingsSMTPServerSkipVerifyFlag != false {

		req.GrafanaSettingsSMTPServer.SkipVerify = reqGrafanaSettingsSMTPServerSkipVerifyFlag
	}
	if reqGrafanaSettingsSMTPServerPortFlag != 0 {

		req.GrafanaSettingsSMTPServer.Port = reqGrafanaSettingsSMTPServerPortFlag
	}
	if reqGrafanaSettingsSMTPServerPasswordFlag != "" {

		req.GrafanaSettingsSMTPServer.Password = reqGrafanaSettingsSMTPServerPasswordFlag
	}
	if reqGrafanaSettingsSMTPServerHostFlag != "" {

		req.GrafanaSettingsSMTPServer.Host = reqGrafanaSettingsSMTPServerHostFlag
	}
	if reqGrafanaSettingsSMTPServerFromNameFlag != "" {

		req.GrafanaSettingsSMTPServer.FromName = reqGrafanaSettingsSMTPServerFromNameFlag
	}
	if reqGrafanaSettingsSMTPServerFromAddressFlag != "" {

		req.GrafanaSettingsSMTPServer.FromAddress = reqGrafanaSettingsSMTPServerFromAddressFlag
	}
	if reqGrafanaSettingsServiceLogFlag != false {
		req.GrafanaSettings = &v3.JSONSchemaGrafana{}
		req.GrafanaSettings.ServiceLog = reqGrafanaSettingsServiceLogFlag
	}
	if reqGrafanaSettingsOauthAllowInsecureEmailLookupFlag != false {

		req.GrafanaSettings.OauthAllowInsecureEmailLookup = reqGrafanaSettingsOauthAllowInsecureEmailLookupFlag
	}
	if reqGrafanaSettingsMetricsEnabledFlag != false {

		req.GrafanaSettings.MetricsEnabled = reqGrafanaSettingsMetricsEnabledFlag
	}
	if reqGrafanaSettingsGoogleAnalyticsUAIDFlag != "" {

		req.GrafanaSettings.GoogleAnalyticsUAID = reqGrafanaSettingsGoogleAnalyticsUAIDFlag
	}
	if reqGrafanaSettingsEditorsCanAdminFlag != false {

		req.GrafanaSettings.EditorsCanAdmin = reqGrafanaSettingsEditorsCanAdminFlag
	}
	if reqGrafanaSettingsDisableGravatarFlag != false {

		req.GrafanaSettings.DisableGravatar = reqGrafanaSettingsDisableGravatarFlag
	}
	if reqGrafanaSettingsDateFormatsIntervalYearFlag != "" {
		req.GrafanaSettingsDateFormats = &v3.GrafanaSettingsDateFormats{}
		req.GrafanaSettingsDateFormats.IntervalYear = reqGrafanaSettingsDateFormatsIntervalYearFlag
	}
	if reqGrafanaSettingsDateFormatsIntervalSecondFlag != "" {

		req.GrafanaSettingsDateFormats.IntervalSecond = reqGrafanaSettingsDateFormatsIntervalSecondFlag
	}
	if reqGrafanaSettingsDateFormatsIntervalMonthFlag != "" {

		req.GrafanaSettingsDateFormats.IntervalMonth = reqGrafanaSettingsDateFormatsIntervalMonthFlag
	}
	if reqGrafanaSettingsDateFormatsIntervalMinuteFlag != "" {

		req.GrafanaSettingsDateFormats.IntervalMinute = reqGrafanaSettingsDateFormatsIntervalMinuteFlag
	}
	if reqGrafanaSettingsDateFormatsIntervalHourFlag != "" {

		req.GrafanaSettingsDateFormats.IntervalHour = reqGrafanaSettingsDateFormatsIntervalHourFlag
	}
	if reqGrafanaSettingsDateFormatsIntervalDayFlag != "" {

		req.GrafanaSettingsDateFormats.IntervalDay = reqGrafanaSettingsDateFormatsIntervalDayFlag
	}
	if reqGrafanaSettingsDateFormatsFullDateFlag != "" {

		req.GrafanaSettingsDateFormats.FullDate = reqGrafanaSettingsDateFormatsFullDateFlag
	}
	if reqGrafanaSettingsDateFormatsDefaultTimezoneFlag != "" {

		req.GrafanaSettingsDateFormats.DefaultTimezone = reqGrafanaSettingsDateFormatsDefaultTimezoneFlag
	}
	if reqGrafanaSettingsDataproxyTimeoutFlag != 0 {
		req.GrafanaSettings = &v3.JSONSchemaGrafana{}
		req.GrafanaSettings.DataproxyTimeout = reqGrafanaSettingsDataproxyTimeoutFlag
	}
	if reqGrafanaSettingsDataproxySendUserHeaderFlag != false {

		req.GrafanaSettings.DataproxySendUserHeader = reqGrafanaSettingsDataproxySendUserHeaderFlag
	}
	if reqGrafanaSettingsDashboardsVersionsToKeepFlag != 0 {

		req.GrafanaSettings.DashboardsVersionsToKeep = reqGrafanaSettingsDashboardsVersionsToKeepFlag
	}
	if reqGrafanaSettingsDashboardsMinRefreshIntervalFlag != "" {

		req.GrafanaSettings.DashboardsMinRefreshInterval = reqGrafanaSettingsDashboardsMinRefreshIntervalFlag
	}
	if reqGrafanaSettingsDashboardPreviewsEnabledFlag != false {

		req.GrafanaSettings.DashboardPreviewsEnabled = reqGrafanaSettingsDashboardPreviewsEnabledFlag
	}
	if reqGrafanaSettingsCustomDomainFlag != "" {

		req.GrafanaSettings.CustomDomain = reqGrafanaSettingsCustomDomainFlag
	}
	if reqGrafanaSettingsCookieSamesiteFlag != "" {

		req.GrafanaSettings.CookieSamesite = reqGrafanaSettingsCookieSamesiteFlag
	}
	if reqGrafanaSettingsAuthGoogleClientSecretFlag != "" {
		req.GrafanaSettingsAuthGoog = &v3.GrafanaSettingsAuthGoogle{}
		req.GrafanaSettingsAuthGoogle.ClientSecret = reqGrafanaSettingsAuthGoogleClientSecretFlag
	}
	if reqGrafanaSettingsAuthGoogleClientIDFlag != "" {

		req.GrafanaSettingsAuthGoogle.ClientID = reqGrafanaSettingsAuthGoogleClientIDFlag
	}
	if reqGrafanaSettingsAuthGoogleAllowSignUPFlag != false {

		req.GrafanaSettingsAuthGoogle.AllowSignUP = reqGrafanaSettingsAuthGoogleAllowSignUPFlag
	}
	if reqGrafanaSettingsAuthGitlabTokenURLFlag != "" {
		req.GrafanaSettingsAuthGitlab = &v3.GrafanaSettingsAuthGitlab{}
		req.GrafanaSettingsAuthGitlab.TokenURL = reqGrafanaSettingsAuthGitlabTokenURLFlag
	}
	if reqGrafanaSettingsAuthGitlabClientSecretFlag != "" {

		req.GrafanaSettingsAuthGitlab.ClientSecret = reqGrafanaSettingsAuthGitlabClientSecretFlag
	}
	if reqGrafanaSettingsAuthGitlabClientIDFlag != "" {

		req.GrafanaSettingsAuthGitlab.ClientID = reqGrafanaSettingsAuthGitlabClientIDFlag
	}
	if reqGrafanaSettingsAuthGitlabAuthURLFlag != "" {

		req.GrafanaSettingsAuthGitlab.AuthURL = reqGrafanaSettingsAuthGitlabAuthURLFlag
	}
	if reqGrafanaSettingsAuthGitlabAPIURLFlag != "" {

		req.GrafanaSettingsAuthGitlab.APIURL = reqGrafanaSettingsAuthGitlabAPIURLFlag
	}
	if reqGrafanaSettingsAuthGitlabAllowSignUPFlag != false {

		req.GrafanaSettingsAuthGitlab.AllowSignUP = reqGrafanaSettingsAuthGitlabAllowSignUPFlag
	}
	if reqGrafanaSettingsAuthGithubSkipOrgRoleSyncFlag != false {
		req.GrafanaSettingsAuthGithub = &v3.GrafanaSettingsAuthGithub{}
		req.GrafanaSettingsAuthGithub.SkipOrgRoleSync = reqGrafanaSettingsAuthGithubSkipOrgRoleSyncFlag
	}
	if reqGrafanaSettingsAuthGithubClientSecretFlag != "" {

		req.GrafanaSettingsAuthGithub.ClientSecret = reqGrafanaSettingsAuthGithubClientSecretFlag
	}
	if reqGrafanaSettingsAuthGithubClientIDFlag != "" {

		req.GrafanaSettingsAuthGithub.ClientID = reqGrafanaSettingsAuthGithubClientIDFlag
	}
	if reqGrafanaSettingsAuthGithubAutoLoginFlag != false {

		req.GrafanaSettingsAuthGithub.AutoLogin = reqGrafanaSettingsAuthGithubAutoLoginFlag
	}
	if reqGrafanaSettingsAuthGithubAllowSignUPFlag != false {

		req.GrafanaSettingsAuthGithub.AllowSignUP = reqGrafanaSettingsAuthGithubAllowSignUPFlag
	}
	if reqGrafanaSettingsAuthGenericOauthTokenURLFlag != "" {
		req.GrafanaSettingsAuthGenericOauth = &v3.GrafanaSettingsAuthGenericOauth{}
		req.GrafanaSettingsAuthGenericOauth.TokenURL = reqGrafanaSettingsAuthGenericOauthTokenURLFlag
	}
	if reqGrafanaSettingsAuthGenericOauthNameFlag != "" {

		req.GrafanaSettingsAuthGenericOauth.Name = reqGrafanaSettingsAuthGenericOauthNameFlag
	}
	if reqGrafanaSettingsAuthGenericOauthClientSecretFlag != "" {

		req.GrafanaSettingsAuthGenericOauth.ClientSecret = reqGrafanaSettingsAuthGenericOauthClientSecretFlag
	}
	if reqGrafanaSettingsAuthGenericOauthClientIDFlag != "" {

		req.GrafanaSettingsAuthGenericOauth.ClientID = reqGrafanaSettingsAuthGenericOauthClientIDFlag
	}
	if reqGrafanaSettingsAuthGenericOauthAutoLoginFlag != false {

		req.GrafanaSettingsAuthGenericOauth.AutoLogin = reqGrafanaSettingsAuthGenericOauthAutoLoginFlag
	}
	if reqGrafanaSettingsAuthGenericOauthAuthURLFlag != "" {

		req.GrafanaSettingsAuthGenericOauth.AuthURL = reqGrafanaSettingsAuthGenericOauthAuthURLFlag
	}
	if reqGrafanaSettingsAuthGenericOauthAPIURLFlag != "" {

		req.GrafanaSettingsAuthGenericOauth.APIURL = reqGrafanaSettingsAuthGenericOauthAPIURLFlag
	}
	if reqGrafanaSettingsAuthGenericOauthAllowSignUPFlag != false {

		req.GrafanaSettingsAuthGenericOauth.AllowSignUP = reqGrafanaSettingsAuthGenericOauthAllowSignUPFlag
	}
	if reqGrafanaSettingsAuthBasicEnabledFlag != false {
		req.GrafanaSetting = &v3.JSONSchemaGrafana{}
		req.GrafanaSettings.AuthBasicEnabled = reqGrafanaSettingsAuthBasicEnabledFlag
	}
	if reqGrafanaSettingsAuthAzureadTokenURLFlag != "" {
		req.GrafanaSettingsAuthAzuread = &v3.GrafanaSettingsAuthAzuread{}
		req.GrafanaSettingsAuthAzuread.TokenURL = reqGrafanaSettingsAuthAzureadTokenURLFlag
	}
	if reqGrafanaSettingsAuthAzureadClientSecretFlag != "" {

		req.GrafanaSettingsAuthAzuread.ClientSecret = reqGrafanaSettingsAuthAzureadClientSecretFlag
	}
	if reqGrafanaSettingsAuthAzureadClientIDFlag != "" {

		req.GrafanaSettingsAuthAzuread.ClientID = reqGrafanaSettingsAuthAzureadClientIDFlag
	}
	if reqGrafanaSettingsAuthAzureadAuthURLFlag != "" {

		req.GrafanaSettingsAuthAzuread.AuthURL = reqGrafanaSettingsAuthAzureadAuthURLFlag
	}
	if reqGrafanaSettingsAuthAzureadAllowSignUPFlag != false {

		req.GrafanaSettingsAuthAzuread.AllowSignUP = reqGrafanaSettingsAuthAzureadAllowSignUPFlag
	}
	if reqGrafanaSettingsAllowEmbeddingFlag != false {
		req.GrafanaSettings = &v3.JSONSchemaGrafana{}
		req.GrafanaSettings.AllowEmbedding = reqGrafanaSettingsAllowEmbeddingFlag
	}
	if reqGrafanaSettingsAlertingNodataORNullvaluesFlag != "" {

		req.GrafanaSettings.AlertingNodataORNullvalues = reqGrafanaSettingsAlertingNodataORNullvaluesFlag
	}
	if reqGrafanaSettingsAlertingMaxAnnotationsToKeepFlag != 0 {

		req.GrafanaSettings.AlertingMaxAnnotationsToKeep = reqGrafanaSettingsAlertingMaxAnnotationsToKeepFlag
	}
	if reqGrafanaSettingsAlertingErrorORTimeoutFlag != "" {

		req.GrafanaSettings.AlertingErrorORTimeout = reqGrafanaSettingsAlertingErrorORTimeoutFlag
	}
	if reqGrafanaSettingsAlertingEnabledFlag != false {

		req.GrafanaSettings.AlertingEnabled = reqGrafanaSettingsAlertingEnabledFlag
	}
	req.ForkFromService = reqForkFromServiceFlag

	resp, err := client.CreateDBAASServiceGrafana(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASServiceGrafanaCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-service-grafana", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqGrafanaSettingsAlertingEnabledFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAlertingEnabledFlag, "grafana-settings.alerting_enabled", false, "Enable or disable Grafana legacy alerting functionality. This should not be enabled with unified_alerting_enabled.")
	var reqGrafanaSettingsAlertingErrorORTimeoutFlag string
	flagset.StringVar(&reqGrafanaSettingsAlertingErrorORTimeoutFlag, "grafana-settings.alerting_error_or_timeout", "", "Default error or timeout setting for new alerting rules")
	var reqGrafanaSettingsAlertingMaxAnnotationsToKeepFlag int
	flagset.IntVar(&reqGrafanaSettingsAlertingMaxAnnotationsToKeepFlag, "grafana-settings.alerting_max_annotations_to_keep", 0, "Max number of alert annotations that Grafana stores. 0 (default) keeps all alert annotations.")
	var reqGrafanaSettingsAlertingNodataORNullvaluesFlag string
	flagset.StringVar(&reqGrafanaSettingsAlertingNodataORNullvaluesFlag, "grafana-settings.alerting_nodata_or_nullvalues", "", "Default value for 'no data or null values' for new alerting rules")
	var reqGrafanaSettingsAllowEmbeddingFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAllowEmbeddingFlag, "grafana-settings.allow_embedding", false, "Allow embedding Grafana dashboards with iframe/frame/object/embed tags. Disabled by default to limit impact of clickjacking")
	var reqGrafanaSettingsAuthAzureadAllowSignUPFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthAzureadAllowSignUPFlag, "grafana-settings.auth_azuread.allow_sign_up", false, "Automatically sign-up users on successful sign-in")
	var reqGrafanaSettingsAuthAzureadAuthURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthAzureadAuthURLFlag, "grafana-settings.auth_azuread.auth_url", "", "Authorization URL")
	var reqGrafanaSettingsAuthAzureadClientIDFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthAzureadClientIDFlag, "grafana-settings.auth_azuread.client_id", "", "Client ID from provider")
	var reqGrafanaSettingsAuthAzureadClientSecretFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthAzureadClientSecretFlag, "grafana-settings.auth_azuread.client_secret", "", "Client secret from provider")
	var reqGrafanaSettingsAuthAzureadTokenURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthAzureadTokenURLFlag, "grafana-settings.auth_azuread.token_url", "", "Token URL")
	var reqGrafanaSettingsAuthBasicEnabledFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthBasicEnabledFlag, "grafana-settings.auth_basic_enabled", false, "Enable or disable basic authentication form, used by Grafana built-in login")
	var reqGrafanaSettingsAuthGenericOauthAllowSignUPFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGenericOauthAllowSignUPFlag, "grafana-settings.auth_generic_oauth.allow_sign_up", false, "Automatically sign-up users on successful sign-in")
	var reqGrafanaSettingsAuthGenericOauthAPIURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGenericOauthAPIURLFlag, "grafana-settings.auth_generic_oauth.api_url", "", "API URL")
	var reqGrafanaSettingsAuthGenericOauthAuthURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGenericOauthAuthURLFlag, "grafana-settings.auth_generic_oauth.auth_url", "", "Authorization URL")
	var reqGrafanaSettingsAuthGenericOauthAutoLoginFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGenericOauthAutoLoginFlag, "grafana-settings.auth_generic_oauth.auto_login", false, "Allow users to bypass the login screen and automatically log in")
	var reqGrafanaSettingsAuthGenericOauthClientIDFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGenericOauthClientIDFlag, "grafana-settings.auth_generic_oauth.client_id", "", "Client ID from provider")
	var reqGrafanaSettingsAuthGenericOauthClientSecretFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGenericOauthClientSecretFlag, "grafana-settings.auth_generic_oauth.client_secret", "", "Client secret from provider")
	var reqGrafanaSettingsAuthGenericOauthNameFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGenericOauthNameFlag, "grafana-settings.auth_generic_oauth.name", "", "Name of the OAuth integration")
	var reqGrafanaSettingsAuthGenericOauthTokenURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGenericOauthTokenURLFlag, "grafana-settings.auth_generic_oauth.token_url", "", "Token URL")
	var reqGrafanaSettingsAuthGithubAllowSignUPFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGithubAllowSignUPFlag, "grafana-settings.auth_github.allow_sign_up", false, "Automatically sign-up users on successful sign-in")
	var reqGrafanaSettingsAuthGithubAutoLoginFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGithubAutoLoginFlag, "grafana-settings.auth_github.auto_login", false, "Allow users to bypass the login screen and automatically log in")
	var reqGrafanaSettingsAuthGithubClientIDFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGithubClientIDFlag, "grafana-settings.auth_github.client_id", "", "Client ID from provider")
	var reqGrafanaSettingsAuthGithubClientSecretFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGithubClientSecretFlag, "grafana-settings.auth_github.client_secret", "", "Client secret from provider")
	var reqGrafanaSettingsAuthGithubSkipOrgRoleSyncFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGithubSkipOrgRoleSyncFlag, "grafana-settings.auth_github.skip_org_role_sync", false, "Stop automatically syncing user roles")
	var reqGrafanaSettingsAuthGitlabAllowSignUPFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGitlabAllowSignUPFlag, "grafana-settings.auth_gitlab.allow_sign_up", false, "Automatically sign-up users on successful sign-in")
	var reqGrafanaSettingsAuthGitlabAPIURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGitlabAPIURLFlag, "grafana-settings.auth_gitlab.api_url", "", "API URL. This only needs to be set when using self hosted GitLab")
	var reqGrafanaSettingsAuthGitlabAuthURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGitlabAuthURLFlag, "grafana-settings.auth_gitlab.auth_url", "", "Authorization URL. This only needs to be set when using self hosted GitLab")
	var reqGrafanaSettingsAuthGitlabClientIDFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGitlabClientIDFlag, "grafana-settings.auth_gitlab.client_id", "", "Client ID from provider")
	var reqGrafanaSettingsAuthGitlabClientSecretFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGitlabClientSecretFlag, "grafana-settings.auth_gitlab.client_secret", "", "Client secret from provider")
	var reqGrafanaSettingsAuthGitlabTokenURLFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGitlabTokenURLFlag, "grafana-settings.auth_gitlab.token_url", "", "Token URL. This only needs to be set when using self hosted GitLab")
	var reqGrafanaSettingsAuthGoogleAllowSignUPFlag bool
	flagset.BoolVar(&reqGrafanaSettingsAuthGoogleAllowSignUPFlag, "grafana-settings.auth_google.allow_sign_up", false, "Automatically sign-up users on successful sign-in")
	var reqGrafanaSettingsAuthGoogleClientIDFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGoogleClientIDFlag, "grafana-settings.auth_google.client_id", "", "Client ID from provider")
	var reqGrafanaSettingsAuthGoogleClientSecretFlag string
	flagset.StringVar(&reqGrafanaSettingsAuthGoogleClientSecretFlag, "grafana-settings.auth_google.client_secret", "", "Client secret from provider")
	var reqGrafanaSettingsCookieSamesiteFlag string
	flagset.StringVar(&reqGrafanaSettingsCookieSamesiteFlag, "grafana-settings.cookie_samesite", "", "Cookie SameSite attribute: 'strict' prevents sending cookie for cross-site requests, effectively disabling direct linking from other sites to Grafana. 'lax' is the default value.")
	var reqGrafanaSettingsCustomDomainFlag string
	flagset.StringVar(&reqGrafanaSettingsCustomDomainFlag, "grafana-settings.custom_domain", "", "Serve the web frontend using a custom CNAME pointing to the Aiven DNS name")
	var reqGrafanaSettingsDashboardPreviewsEnabledFlag bool
	flagset.BoolVar(&reqGrafanaSettingsDashboardPreviewsEnabledFlag, "grafana-settings.dashboard_previews_enabled", false, "This feature is new in Grafana 9 and is quite resource intensive. It may cause low-end plans to work more slowly while the dashboard previews are rendering.")
	var reqGrafanaSettingsDashboardsMinRefreshIntervalFlag string
	flagset.StringVar(&reqGrafanaSettingsDashboardsMinRefreshIntervalFlag, "grafana-settings.dashboards_min_refresh_interval", "", "Signed sequence of decimal numbers, followed by a unit suffix (ms, s, m, h, d), e.g. 30s, 1h")
	var reqGrafanaSettingsDashboardsVersionsToKeepFlag int
	flagset.IntVar(&reqGrafanaSettingsDashboardsVersionsToKeepFlag, "grafana-settings.dashboards_versions_to_keep", 0, "Dashboard versions to keep per dashboard")
	var reqGrafanaSettingsDataproxySendUserHeaderFlag bool
	flagset.BoolVar(&reqGrafanaSettingsDataproxySendUserHeaderFlag, "grafana-settings.dataproxy_send_user_header", false, "Send 'X-Grafana-User' header to data source")
	var reqGrafanaSettingsDataproxyTimeoutFlag int
	flagset.IntVar(&reqGrafanaSettingsDataproxyTimeoutFlag, "grafana-settings.dataproxy_timeout", 0, "Timeout for data proxy requests in seconds")
	var reqGrafanaSettingsDateFormatsDefaultTimezoneFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsDefaultTimezoneFlag, "grafana-settings.date_formats.default_timezone", "", "Default time zone for user preferences. Value 'browser' uses browser local time zone.")
	var reqGrafanaSettingsDateFormatsFullDateFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsFullDateFlag, "grafana-settings.date_formats.full_date", "", "Moment.js style format string for cases where full date is shown")
	var reqGrafanaSettingsDateFormatsIntervalDayFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsIntervalDayFlag, "grafana-settings.date_formats.interval_day", "", "Moment.js style format string used when a time requiring day accuracy is shown")
	var reqGrafanaSettingsDateFormatsIntervalHourFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsIntervalHourFlag, "grafana-settings.date_formats.interval_hour", "", "Moment.js style format string used when a time requiring hour accuracy is shown")
	var reqGrafanaSettingsDateFormatsIntervalMinuteFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsIntervalMinuteFlag, "grafana-settings.date_formats.interval_minute", "", "Moment.js style format string used when a time requiring minute accuracy is shown")
	var reqGrafanaSettingsDateFormatsIntervalMonthFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsIntervalMonthFlag, "grafana-settings.date_formats.interval_month", "", "Moment.js style format string used when a time requiring month accuracy is shown")
	var reqGrafanaSettingsDateFormatsIntervalSecondFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsIntervalSecondFlag, "grafana-settings.date_formats.interval_second", "", "Moment.js style format string used when a time requiring second accuracy is shown")
	var reqGrafanaSettingsDateFormatsIntervalYearFlag string
	flagset.StringVar(&reqGrafanaSettingsDateFormatsIntervalYearFlag, "grafana-settings.date_formats.interval_year", "", "Moment.js style format string used when a time requiring year accuracy is shown")
	var reqGrafanaSettingsDisableGravatarFlag bool
	flagset.BoolVar(&reqGrafanaSettingsDisableGravatarFlag, "grafana-settings.disable_gravatar", false, "Set to true to disable gravatar. Defaults to false (gravatar is enabled)")
	var reqGrafanaSettingsEditorsCanAdminFlag bool
	flagset.BoolVar(&reqGrafanaSettingsEditorsCanAdminFlag, "grafana-settings.editors_can_admin", false, "Editors can manage folders, teams and dashboards created by them")
	var reqGrafanaSettingsGoogleAnalyticsUAIDFlag string
	flagset.StringVar(&reqGrafanaSettingsGoogleAnalyticsUAIDFlag, "grafana-settings.google_analytics_ua_id", "", "Google Analytics ID")
	var reqGrafanaSettingsMetricsEnabledFlag bool
	flagset.BoolVar(&reqGrafanaSettingsMetricsEnabledFlag, "grafana-settings.metrics_enabled", false, "Enable Grafana /metrics endpoint")
	var reqGrafanaSettingsOauthAllowInsecureEmailLookupFlag bool
	flagset.BoolVar(&reqGrafanaSettingsOauthAllowInsecureEmailLookupFlag, "grafana-settings.oauth_allow_insecure_email_lookup", false, "Enforce user lookup based on email instead of the unique ID provided by the IdP")
	var reqGrafanaSettingsServiceLogFlag bool
	flagset.BoolVar(&reqGrafanaSettingsServiceLogFlag, "grafana-settings.service_log", false, "Store logs for the service so that they are available in the HTTP API and console.")
	var reqGrafanaSettingsSMTPServerFromAddressFlag string
	flagset.StringVar(&reqGrafanaSettingsSMTPServerFromAddressFlag, "grafana-settings.smtp_server.from_address", "", "Address used for sending emails")
	var reqGrafanaSettingsSMTPServerFromNameFlag string
	flagset.StringVar(&reqGrafanaSettingsSMTPServerFromNameFlag, "grafana-settings.smtp_server.from_name", "", "Name used in outgoing emails, defaults to Grafana")
	var reqGrafanaSettingsSMTPServerHostFlag string
	flagset.StringVar(&reqGrafanaSettingsSMTPServerHostFlag, "grafana-settings.smtp_server.host", "", "Server hostname or IP")
	var reqGrafanaSettingsSMTPServerPasswordFlag string
	flagset.StringVar(&reqGrafanaSettingsSMTPServerPasswordFlag, "grafana-settings.smtp_server.password", "", "Password for SMTP authentication")
	var reqGrafanaSettingsSMTPServerPortFlag int
	flagset.IntVar(&reqGrafanaSettingsSMTPServerPortFlag, "grafana-settings.smtp_server.port", 0, "SMTP server port")
	var reqGrafanaSettingsSMTPServerSkipVerifyFlag bool
	flagset.BoolVar(&reqGrafanaSettingsSMTPServerSkipVerifyFlag, "grafana-settings.smtp_server.skip_verify", false, "Skip verifying server certificate. Defaults to false")
	var reqGrafanaSettingsSMTPServerStarttlsPolicyFlag string
	flagset.StringVar(&reqGrafanaSettingsSMTPServerStarttlsPolicyFlag, "grafana-settings.smtp_server.starttls_policy", "", "Either OpportunisticStartTLS, MandatoryStartTLS or NoStartTLS. Default is OpportunisticStartTLS.")
	var reqGrafanaSettingsSMTPServerUsernameFlag string
	flagset.StringVar(&reqGrafanaSettingsSMTPServerUsernameFlag, "grafana-settings.smtp_server.username", "", "Username for SMTP authentication")
	var reqGrafanaSettingsUnifiedAlertingEnabledFlag bool
	flagset.BoolVar(&reqGrafanaSettingsUnifiedAlertingEnabledFlag, "grafana-settings.unified_alerting_enabled", false, "Enable or disable Grafana unified alerting functionality. By default this is enabled and any legacy alerts will be migrated on upgrade to Grafana 9+. To stay on legacy alerting, set unified_alerting_enabled to false and alerting_enabled to true. See https://grafana.com/docs/grafana/latest/alerting/set-up/migrating-alerts/ for more details.")
	var reqGrafanaSettingsUserAutoAssignOrgFlag bool
	flagset.BoolVar(&reqGrafanaSettingsUserAutoAssignOrgFlag, "grafana-settings.user_auto_assign_org", false, "Auto-assign new users on signup to main organization. Defaults to false")
	var reqGrafanaSettingsUserAutoAssignOrgRoleFlag string
	flagset.StringVar(&reqGrafanaSettingsUserAutoAssignOrgRoleFlag, "grafana-settings.user_auto_assign_org_role", "", "Set role for new signups. Defaults to Viewer")
	var reqGrafanaSettingsViewersCanEditFlag bool
	flagset.BoolVar(&reqGrafanaSettingsViewersCanEditFlag, "grafana-settings.viewers_can_edit", false, "Users with view-only permission can edit but not save dashboards")
	var reqGrafanaSettingsWalFlag bool
	flagset.BoolVar(&reqGrafanaSettingsWalFlag, "grafana-settings.wal", false, "Setting to enable/disable Write-Ahead Logging. The default value is false (disabled).")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASServiceGrafanaRequest
	req.TerminationProtection = reqTerminationProtectionFlag
	req.Plan = reqPlanFlag
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.UpdateDBAASServiceGrafanaRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}
	if reqGrafanaSettingsWalFlag != false {
		req.GrafanaSettings = &v3.JSONSchemaGrafana{}
		req.GrafanaSettings.Wal = reqGrafanaSettingsWalFlag
	}
	if reqGrafanaSettingsViewersCanEditFlag != false {

		req.GrafanaSettings.ViewersCanEdit = reqGrafanaSettingsViewersCanEditFlag
	}
	if reqGrafanaSettingsUserAutoAssignOrgRoleFlag != "" {

		req.GrafanaSettings.UserAutoAssignOrgRole = reqGrafanaSettingsUserAutoAssignOrgRoleFlag
	}
	if reqGrafanaSettingsUserAutoAssignOrgFlag != false {

		req.GrafanaSettings.UserAutoAssignOrg = reqGrafanaSettingsUserAutoAssignOrgFlag
	}
	if reqGrafanaSettingsUnifiedAlertingEnabledFlag != false {

		req.GrafanaSettings.UnifiedAlertingEnabled = reqGrafanaSettingsUnifiedAlertingEnabledFlag
	}
	if reqGrafanaSettingsSMTPServerUsernameFlag != "" {
		req.GrafanaSettingsSMTPServ = &v3.GrafanaSettingsSMTPServer{}
		req.GrafanaSettingsSMTPServer.Username = reqGrafanaSettingsSMTPServerUsernameFlag
	}
	if reqGrafanaSettingsSMTPServerStarttlsPolicyFlag != "" {

		req.GrafanaSettingsSMTPServer.StarttlsPolicy = reqGrafanaSettingsSMTPServerStarttlsPolicyFlag
	}
	if reqGrafanaSettingsSMTPServerSkipVerifyFlag != false {

		req.GrafanaSettingsSMTPServer.SkipVerify = reqGrafanaSettingsSMTPServerSkipVerifyFlag
	}
	if reqGrafanaSettingsSMTPServerPortFlag != 0 {

		req.GrafanaSettingsSMTPServer.Port = reqGrafanaSettingsSMTPServerPortFlag
	}
	if reqGrafanaSettingsSMTPServerPasswordFlag != "" {

		req.GrafanaSettingsSMTPServer.Password = reqGrafanaSettingsSMTPServerPasswordFlag
	}
	if reqGrafanaSettingsSMTPServerHostFlag != "" {

		req.GrafanaSettingsSMTPServer.Host = reqGrafanaSettingsSMTPServerHostFlag
	}
	if reqGrafanaSettingsSMTPServerFromNameFlag != "" {

		req.GrafanaSettingsSMTPServer.FromName = reqGrafanaSettingsSMTPServerFromNameFlag
	}
	if reqGrafanaSettingsSMTPServerFromAddressFlag != "" {

		req.GrafanaSettingsSMTPServer.FromAddress = reqGrafanaSettingsSMTPServerFromAddressFlag
	}
	if reqGrafanaSettingsServiceLogFlag != false {
		req.GrafanaSettings = &v3.JSONSchemaGrafana{}
		req.GrafanaSettings.ServiceLog = reqGrafanaSettingsServiceLogFlag
	}
	if reqGrafanaSettingsOauthAllowInsecureEmailLookupFlag != false {

		req.GrafanaSettings.OauthAllowInsecureEmailLookup = reqGrafanaSettingsOauthAllowInsecureEmailLookupFlag
	}
	if reqGrafanaSettingsMetricsEnabledFlag != false {

		req.GrafanaSettings.MetricsEnabled = reqGrafanaSettingsMetricsEnabledFlag
	}
	if reqGrafanaSettingsGoogleAnalyticsUAIDFlag != "" {

		req.GrafanaSettings.GoogleAnalyticsUAID = reqGrafanaSettingsGoogleAnalyticsUAIDFlag
	}
	if reqGrafanaSettingsEditorsCanAdminFlag != false {

		req.GrafanaSettings.EditorsCanAdmin = reqGrafanaSettingsEditorsCanAdminFlag
	}
	if reqGrafanaSettingsDisableGravatarFlag != false {

		req.GrafanaSettings.DisableGravatar = reqGrafanaSettingsDisableGravatarFlag
	}
	if reqGrafanaSettingsDateFormatsIntervalYearFlag != "" {
		req.GrafanaSettingsDateFormats = &v3.GrafanaSettingsDateFormats{}
		req.GrafanaSettingsDateFormats.IntervalYear = reqGrafanaSettingsDateFormatsIntervalYearFlag
	}
	if reqGrafanaSettingsDateFormatsIntervalSecondFlag != "" {

		req.GrafanaSettingsDateFormats.IntervalSecond = reqGrafanaSettingsDateFormatsIntervalSecondFlag
	}
	if reqGrafanaSettingsDateFormatsIntervalMonthFlag != "" {

		req.GrafanaSettingsDateFormats.IntervalMonth = reqGrafanaSettingsDateFormatsIntervalMonthFlag
	}
	if reqGrafanaSettingsDateFormatsIntervalMinuteFlag != "" {

		req.GrafanaSettingsDateFormats.IntervalMinute = reqGrafanaSettingsDateFormatsIntervalMinuteFlag
	}
	if reqGrafanaSettingsDateFormatsIntervalHourFlag != "" {

		req.GrafanaSettingsDateFormats.IntervalHour = reqGrafanaSettingsDateFormatsIntervalHourFlag
	}
	if reqGrafanaSettingsDateFormatsIntervalDayFlag != "" {

		req.GrafanaSettingsDateFormats.IntervalDay = reqGrafanaSettingsDateFormatsIntervalDayFlag
	}
	if reqGrafanaSettingsDateFormatsFullDateFlag != "" {

		req.GrafanaSettingsDateFormats.FullDate = reqGrafanaSettingsDateFormatsFullDateFlag
	}
	if reqGrafanaSettingsDateFormatsDefaultTimezoneFlag != "" {

		req.GrafanaSettingsDateFormats.DefaultTimezone = reqGrafanaSettingsDateFormatsDefaultTimezoneFlag
	}
	if reqGrafanaSettingsDataproxyTimeoutFlag != 0 {
		req.GrafanaSettings = &v3.JSONSchemaGrafana{}
		req.GrafanaSettings.DataproxyTimeout = reqGrafanaSettingsDataproxyTimeoutFlag
	}
	if reqGrafanaSettingsDataproxySendUserHeaderFlag != false {

		req.GrafanaSettings.DataproxySendUserHeader = reqGrafanaSettingsDataproxySendUserHeaderFlag
	}
	if reqGrafanaSettingsDashboardsVersionsToKeepFlag != 0 {

		req.GrafanaSettings.DashboardsVersionsToKeep = reqGrafanaSettingsDashboardsVersionsToKeepFlag
	}
	if reqGrafanaSettingsDashboardsMinRefreshIntervalFlag != "" {

		req.GrafanaSettings.DashboardsMinRefreshInterval = reqGrafanaSettingsDashboardsMinRefreshIntervalFlag
	}
	if reqGrafanaSettingsDashboardPreviewsEnabledFlag != false {

		req.GrafanaSettings.DashboardPreviewsEnabled = reqGrafanaSettingsDashboardPreviewsEnabledFlag
	}
	if reqGrafanaSettingsCustomDomainFlag != "" {

		req.GrafanaSettings.CustomDomain = reqGrafanaSettingsCustomDomainFlag
	}
	if reqGrafanaSettingsCookieSamesiteFlag != "" {

		req.GrafanaSettings.CookieSamesite = reqGrafanaSettingsCookieSamesiteFlag
	}
	if reqGrafanaSettingsAuthGoogleClientSecretFlag != "" {
		req.GrafanaSettingsAuthGoog = &v3.GrafanaSettingsAuthGoogle{}
		req.GrafanaSettingsAuthGoogle.ClientSecret = reqGrafanaSettingsAuthGoogleClientSecretFlag
	}
	if reqGrafanaSettingsAuthGoogleClientIDFlag != "" {

		req.GrafanaSettingsAuthGoogle.ClientID = reqGrafanaSettingsAuthGoogleClientIDFlag
	}
	if reqGrafanaSettingsAuthGoogleAllowSignUPFlag != false {

		req.GrafanaSettingsAuthGoogle.AllowSignUP = reqGrafanaSettingsAuthGoogleAllowSignUPFlag
	}
	if reqGrafanaSettingsAuthGitlabTokenURLFlag != "" {
		req.GrafanaSettingsAuthGitlab = &v3.GrafanaSettingsAuthGitlab{}
		req.GrafanaSettingsAuthGitlab.TokenURL = reqGrafanaSettingsAuthGitlabTokenURLFlag
	}
	if reqGrafanaSettingsAuthGitlabClientSecretFlag != "" {

		req.GrafanaSettingsAuthGitlab.ClientSecret = reqGrafanaSettingsAuthGitlabClientSecretFlag
	}
	if reqGrafanaSettingsAuthGitlabClientIDFlag != "" {

		req.GrafanaSettingsAuthGitlab.ClientID = reqGrafanaSettingsAuthGitlabClientIDFlag
	}
	if reqGrafanaSettingsAuthGitlabAuthURLFlag != "" {

		req.GrafanaSettingsAuthGitlab.AuthURL = reqGrafanaSettingsAuthGitlabAuthURLFlag
	}
	if reqGrafanaSettingsAuthGitlabAPIURLFlag != "" {

		req.GrafanaSettingsAuthGitlab.APIURL = reqGrafanaSettingsAuthGitlabAPIURLFlag
	}
	if reqGrafanaSettingsAuthGitlabAllowSignUPFlag != false {

		req.GrafanaSettingsAuthGitlab.AllowSignUP = reqGrafanaSettingsAuthGitlabAllowSignUPFlag
	}
	if reqGrafanaSettingsAuthGithubSkipOrgRoleSyncFlag != false {
		req.GrafanaSettingsAuthGithub = &v3.GrafanaSettingsAuthGithub{}
		req.GrafanaSettingsAuthGithub.SkipOrgRoleSync = reqGrafanaSettingsAuthGithubSkipOrgRoleSyncFlag
	}
	if reqGrafanaSettingsAuthGithubClientSecretFlag != "" {

		req.GrafanaSettingsAuthGithub.ClientSecret = reqGrafanaSettingsAuthGithubClientSecretFlag
	}
	if reqGrafanaSettingsAuthGithubClientIDFlag != "" {

		req.GrafanaSettingsAuthGithub.ClientID = reqGrafanaSettingsAuthGithubClientIDFlag
	}
	if reqGrafanaSettingsAuthGithubAutoLoginFlag != false {

		req.GrafanaSettingsAuthGithub.AutoLogin = reqGrafanaSettingsAuthGithubAutoLoginFlag
	}
	if reqGrafanaSettingsAuthGithubAllowSignUPFlag != false {

		req.GrafanaSettingsAuthGithub.AllowSignUP = reqGrafanaSettingsAuthGithubAllowSignUPFlag
	}
	if reqGrafanaSettingsAuthGenericOauthTokenURLFlag != "" {
		req.GrafanaSettingsAuthGenericOauth = &v3.GrafanaSettingsAuthGenericOauth{}
		req.GrafanaSettingsAuthGenericOauth.TokenURL = reqGrafanaSettingsAuthGenericOauthTokenURLFlag
	}
	if reqGrafanaSettingsAuthGenericOauthNameFlag != "" {

		req.GrafanaSettingsAuthGenericOauth.Name = reqGrafanaSettingsAuthGenericOauthNameFlag
	}
	if reqGrafanaSettingsAuthGenericOauthClientSecretFlag != "" {

		req.GrafanaSettingsAuthGenericOauth.ClientSecret = reqGrafanaSettingsAuthGenericOauthClientSecretFlag
	}
	if reqGrafanaSettingsAuthGenericOauthClientIDFlag != "" {

		req.GrafanaSettingsAuthGenericOauth.ClientID = reqGrafanaSettingsAuthGenericOauthClientIDFlag
	}
	if reqGrafanaSettingsAuthGenericOauthAutoLoginFlag != false {

		req.GrafanaSettingsAuthGenericOauth.AutoLogin = reqGrafanaSettingsAuthGenericOauthAutoLoginFlag
	}
	if reqGrafanaSettingsAuthGenericOauthAuthURLFlag != "" {

		req.GrafanaSettingsAuthGenericOauth.AuthURL = reqGrafanaSettingsAuthGenericOauthAuthURLFlag
	}
	if reqGrafanaSettingsAuthGenericOauthAPIURLFlag != "" {

		req.GrafanaSettingsAuthGenericOauth.APIURL = reqGrafanaSettingsAuthGenericOauthAPIURLFlag
	}
	if reqGrafanaSettingsAuthGenericOauthAllowSignUPFlag != false {

		req.GrafanaSettingsAuthGenericOauth.AllowSignUP = reqGrafanaSettingsAuthGenericOauthAllowSignUPFlag
	}
	if reqGrafanaSettingsAuthBasicEnabledFlag != false {
		req.GrafanaSetting = &v3.JSONSchemaGrafana{}
		req.GrafanaSettings.AuthBasicEnabled = reqGrafanaSettingsAuthBasicEnabledFlag
	}
	if reqGrafanaSettingsAuthAzureadTokenURLFlag != "" {
		req.GrafanaSettingsAuthAzuread = &v3.GrafanaSettingsAuthAzuread{}
		req.GrafanaSettingsAuthAzuread.TokenURL = reqGrafanaSettingsAuthAzureadTokenURLFlag
	}
	if reqGrafanaSettingsAuthAzureadClientSecretFlag != "" {

		req.GrafanaSettingsAuthAzuread.ClientSecret = reqGrafanaSettingsAuthAzureadClientSecretFlag
	}
	if reqGrafanaSettingsAuthAzureadClientIDFlag != "" {

		req.GrafanaSettingsAuthAzuread.ClientID = reqGrafanaSettingsAuthAzureadClientIDFlag
	}
	if reqGrafanaSettingsAuthAzureadAuthURLFlag != "" {

		req.GrafanaSettingsAuthAzuread.AuthURL = reqGrafanaSettingsAuthAzureadAuthURLFlag
	}
	if reqGrafanaSettingsAuthAzureadAllowSignUPFlag != false {

		req.GrafanaSettingsAuthAzuread.AllowSignUP = reqGrafanaSettingsAuthAzureadAllowSignUPFlag
	}
	if reqGrafanaSettingsAllowEmbeddingFlag != false {
		req.GrafanaSettings = &v3.JSONSchemaGrafana{}
		req.GrafanaSettings.AllowEmbedding = reqGrafanaSettingsAllowEmbeddingFlag
	}
	if reqGrafanaSettingsAlertingNodataORNullvaluesFlag != "" {

		req.GrafanaSettings.AlertingNodataORNullvalues = reqGrafanaSettingsAlertingNodataORNullvaluesFlag
	}
	if reqGrafanaSettingsAlertingMaxAnnotationsToKeepFlag != 0 {

		req.GrafanaSettings.AlertingMaxAnnotationsToKeep = reqGrafanaSettingsAlertingMaxAnnotationsToKeepFlag
	}
	if reqGrafanaSettingsAlertingErrorORTimeoutFlag != "" {

		req.GrafanaSettings.AlertingErrorORTimeout = reqGrafanaSettingsAlertingErrorORTimeoutFlag
	}
	if reqGrafanaSettingsAlertingEnabledFlag != false {

		req.GrafanaSettings.AlertingEnabled = reqGrafanaSettingsAlertingEnabledFlag
	}

	resp, err := client.UpdateDBAASServiceGrafana(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StartDBAASGrafanaMaintenanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("start-dbaas-grafana-maintenance", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StartDBAASGrafanaMaintenance(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetDBAASGrafanaUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-dbaas-grafana-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")
	var reqPasswordFlag string
	flagset.StringVar(&reqPasswordFlag, "password", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ResetDBAASGrafanaUserPasswordRequest
	req.Password = reqPasswordFlag

	resp, err := client.ResetDBAASGrafanaUserPassword(context.Background(), serviceNameFlag, usernameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RevealDBAASGrafanaUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reveal-dbaas-grafana-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RevealDBAASGrafanaUserPassword(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASIntegrationCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-integration", flag.ExitOnError)
	var reqDestServiceFlag string
	flagset.StringVar(&reqDestServiceFlag, "dest-service", "", "")
	var reqIntegrationTypeFlag string
	flagset.StringVar(&reqIntegrationTypeFlag, "integration-type", "", "")
	var reqSourceServiceFlag string
	flagset.StringVar(&reqSourceServiceFlag, "source-service", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASIntegrationRequest
	req.SourceService = reqSourceServiceFlag
	req.IntegrationType = reqIntegrationTypeFlag
	req.DestService = reqDestServiceFlag

	resp, err := client.CreateDBAASIntegration(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListDBAASIntegrationSettingsCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("list-dbaas-integration-settings", flag.ExitOnError)
	var integrationTypeFlag string
	flagset.StringVar(&integrationTypeFlag, "IntegrationType", "", "")
	var sourceTypeFlag string
	flagset.StringVar(&sourceTypeFlag, "SourceType", "", "")
	var destTypeFlag string
	flagset.StringVar(&destTypeFlag, "DestType", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ListDBAASIntegrationSettings(context.Background(), integrationTypeFlag, sourceTypeFlag, destTypeFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListDBAASIntegrationTypesCmd(client *v3.Client) {
	resp, err := client.ListDBAASIntegrationTypes(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASIntegrationCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-integration", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASIntegration(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASIntegrationCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-integration", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASIntegration(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASIntegrationCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-integration", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASIntegrationRequest

	resp, err := client.UpdateDBAASIntegration(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASServiceKafkaCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-service-kafka", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASServiceKafka(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASServiceKafkaCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-service-kafka", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASServiceKafka(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASServiceKafkaCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-service-kafka", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqAuthenticationMethodsCertificateFlag bool
	flagset.BoolVar(&reqAuthenticationMethodsCertificateFlag, "authentication-methods.certificate", false, "Enable certificate/SSL authentication")
	var reqAuthenticationMethodsSaslFlag bool
	flagset.BoolVar(&reqAuthenticationMethodsSaslFlag, "authentication-methods.sasl", false, "Enable SASL authentication")
	var reqKafkaConnectEnabledFlag bool
	flagset.BoolVar(&reqKafkaConnectEnabledFlag, "kafka-connect-enabled", false, "Allow clients to connect to kafka_connect from the public internet for service nodes that are in a project VPC or another type of private network")
	var reqKafkaConnectSettingsConnectorClientConfigOverridePolicyFlag string
	flagset.StringVar(&reqKafkaConnectSettingsConnectorClientConfigOverridePolicyFlag, "kafka-connect-settings.connector_client_config_override_policy", "", "Defines what client configurations can be overridden by the connector. Default is None")
	var reqKafkaConnectSettingsConsumerAutoOffsetResetFlag string
	flagset.StringVar(&reqKafkaConnectSettingsConsumerAutoOffsetResetFlag, "kafka-connect-settings.consumer_auto_offset_reset", "", "What to do when there is no initial offset in Kafka or if the current offset does not exist any more on the server. Default is earliest")
	var reqKafkaConnectSettingsConsumerFetchMaxBytesFlag int
	flagset.IntVar(&reqKafkaConnectSettingsConsumerFetchMaxBytesFlag, "kafka-connect-settings.consumer_fetch_max_bytes", 0, "Records are fetched in batches by the consumer, and if the first record batch in the first non-empty partition of the fetch is larger than this value, the record batch will still be returned to ensure that the consumer can make progress. As such, this is not a absolute maximum.")
	var reqKafkaConnectSettingsConsumerIsolationLevelFlag string
	flagset.StringVar(&reqKafkaConnectSettingsConsumerIsolationLevelFlag, "kafka-connect-settings.consumer_isolation_level", "", "Transaction read isolation level. read_uncommitted is the default, but read_committed can be used if consume-exactly-once behavior is desired.")
	var reqKafkaConnectSettingsConsumerMaxPartitionFetchBytesFlag int
	flagset.IntVar(&reqKafkaConnectSettingsConsumerMaxPartitionFetchBytesFlag, "kafka-connect-settings.consumer_max_partition_fetch_bytes", 0, "Records are fetched in batches by the consumer.If the first record batch in the first non-empty partition of the fetch is larger than this limit, the batch will still be returned to ensure that the consumer can make progress. ")
	var reqKafkaConnectSettingsConsumerMaxPollIntervalMSFlag int
	flagset.IntVar(&reqKafkaConnectSettingsConsumerMaxPollIntervalMSFlag, "kafka-connect-settings.consumer_max_poll_interval_ms", 0, "The maximum delay in milliseconds between invocations of poll() when using consumer group management (defaults to 300000).")
	var reqKafkaConnectSettingsConsumerMaxPollRecordsFlag int
	flagset.IntVar(&reqKafkaConnectSettingsConsumerMaxPollRecordsFlag, "kafka-connect-settings.consumer_max_poll_records", 0, "The maximum number of records returned in a single call to poll() (defaults to 500).")
	var reqKafkaConnectSettingsOffsetFlushIntervalMSFlag int
	flagset.IntVar(&reqKafkaConnectSettingsOffsetFlushIntervalMSFlag, "kafka-connect-settings.offset_flush_interval_ms", 0, "The interval at which to try committing offsets for tasks (defaults to 60000).")
	var reqKafkaConnectSettingsOffsetFlushTimeoutMSFlag int
	flagset.IntVar(&reqKafkaConnectSettingsOffsetFlushTimeoutMSFlag, "kafka-connect-settings.offset_flush_timeout_ms", 0, "Maximum number of milliseconds to wait for records to flush and partition offset data to be committed to offset storage before cancelling the process and restoring the offset data to be committed in a future attempt (defaults to 5000).")
	var reqKafkaConnectSettingsProducerBatchSizeFlag int
	flagset.IntVar(&reqKafkaConnectSettingsProducerBatchSizeFlag, "kafka-connect-settings.producer_batch_size", 0, "This setting gives the upper bound of the batch size to be sent. If there are fewer than this many bytes accumulated for this partition, the producer will 'linger' for the linger.ms time waiting for more records to show up. A batch size of zero will disable batching entirely (defaults to 16384).")
	var reqKafkaConnectSettingsProducerBufferMemoryFlag int
	flagset.IntVar(&reqKafkaConnectSettingsProducerBufferMemoryFlag, "kafka-connect-settings.producer_buffer_memory", 0, "The total bytes of memory the producer can use to buffer records waiting to be sent to the broker (defaults to 33554432).")
	var reqKafkaConnectSettingsProducerCompressionTypeFlag string
	flagset.StringVar(&reqKafkaConnectSettingsProducerCompressionTypeFlag, "kafka-connect-settings.producer_compression_type", "", "Specify the default compression type for producers. This configuration accepts the standard compression codecs ('gzip', 'snappy', 'lz4', 'zstd'). It additionally accepts 'none' which is the default and equivalent to no compression.")
	var reqKafkaConnectSettingsProducerLingerMSFlag int
	flagset.IntVar(&reqKafkaConnectSettingsProducerLingerMSFlag, "kafka-connect-settings.producer_linger_ms", 0, "This setting gives the upper bound on the delay for batching: once there is batch.size worth of records for a partition it will be sent immediately regardless of this setting, however if there are fewer than this many bytes accumulated for this partition the producer will 'linger' for the specified time waiting for more records to show up. Defaults to 0.")
	var reqKafkaConnectSettingsProducerMaxRequestSizeFlag int
	flagset.IntVar(&reqKafkaConnectSettingsProducerMaxRequestSizeFlag, "kafka-connect-settings.producer_max_request_size", 0, "This setting will limit the number of record batches the producer will send in a single request to avoid sending huge requests.")
	var reqKafkaConnectSettingsScheduledRebalanceMaxDelayMSFlag int
	flagset.IntVar(&reqKafkaConnectSettingsScheduledRebalanceMaxDelayMSFlag, "kafka-connect-settings.scheduled_rebalance_max_delay_ms", 0, "The maximum delay that is scheduled in order to wait for the return of one or more departed workers before rebalancing and reassigning their connectors and tasks to the group. During this period the connectors and tasks of the departed workers remain unassigned. Defaults to 5 minutes.")
	var reqKafkaConnectSettingsSessionTimeoutMSFlag int
	flagset.IntVar(&reqKafkaConnectSettingsSessionTimeoutMSFlag, "kafka-connect-settings.session_timeout_ms", 0, "The timeout in milliseconds used to detect failures when using Kafka’s group management facilities (defaults to 10000).")
	var reqKafkaRestEnabledFlag bool
	flagset.BoolVar(&reqKafkaRestEnabledFlag, "kafka-rest-enabled", false, "Enable Kafka-REST service")
	var reqKafkaRestSettingsConsumerEnableAutoCommitFlag bool
	flagset.BoolVar(&reqKafkaRestSettingsConsumerEnableAutoCommitFlag, "kafka-rest-settings.consumer_enable_auto_commit", false, "If true the consumer's offset will be periodically committed to Kafka in the background")
	var reqKafkaRestSettingsConsumerRequestMaxBytesFlag int
	flagset.IntVar(&reqKafkaRestSettingsConsumerRequestMaxBytesFlag, "kafka-rest-settings.consumer_request_max_bytes", 0, "Maximum number of bytes in unencoded message keys and values by a single request")
	var reqKafkaRestSettingsConsumerRequestTimeoutMSFlag int
	flagset.IntVar(&reqKafkaRestSettingsConsumerRequestTimeoutMSFlag, "kafka-rest-settings.consumer_request_timeout_ms", 0, "The maximum total time to wait for messages for a request if the maximum number of messages has not yet been reached")
	var reqKafkaRestSettingsNameStrategyFlag string
	flagset.StringVar(&reqKafkaRestSettingsNameStrategyFlag, "kafka-rest-settings.name_strategy", "", "Name strategy to use when selecting subject for storing schemas")
	var reqKafkaRestSettingsNameStrategyValidationFlag bool
	flagset.BoolVar(&reqKafkaRestSettingsNameStrategyValidationFlag, "kafka-rest-settings.name_strategy_validation", false, "If true, validate that given schema is registered under expected subject name by the used name strategy when producing messages.")
	var reqKafkaRestSettingsProducerAcksFlag string
	flagset.StringVar(&reqKafkaRestSettingsProducerAcksFlag, "kafka-rest-settings.producer_acks", "", "The number of acknowledgments the producer requires the leader to have received before considering a request complete. If set to 'all' or '-1', the leader will wait for the full set of in-sync replicas to acknowledge the record.")
	var reqKafkaRestSettingsProducerCompressionTypeFlag string
	flagset.StringVar(&reqKafkaRestSettingsProducerCompressionTypeFlag, "kafka-rest-settings.producer_compression_type", "", "Specify the default compression type for producers. This configuration accepts the standard compression codecs ('gzip', 'snappy', 'lz4', 'zstd'). It additionally accepts 'none' which is the default and equivalent to no compression.")
	var reqKafkaRestSettingsProducerLingerMSFlag int
	flagset.IntVar(&reqKafkaRestSettingsProducerLingerMSFlag, "kafka-rest-settings.producer_linger_ms", 0, "Wait for up to the given delay to allow batching records together")
	var reqKafkaRestSettingsProducerMaxRequestSizeFlag int
	flagset.IntVar(&reqKafkaRestSettingsProducerMaxRequestSizeFlag, "kafka-rest-settings.producer_max_request_size", 0, "The maximum size of a request in bytes. Note that Kafka broker can also cap the record batch size.")
	var reqKafkaRestSettingsSimpleconsumerPoolSizeMaxFlag int
	flagset.IntVar(&reqKafkaRestSettingsSimpleconsumerPoolSizeMaxFlag, "kafka-rest-settings.simpleconsumer_pool_size_max", 0, "Maximum number of SimpleConsumers that can be instantiated per broker")
	var reqKafkaSettingsAutoCreateTopicsEnableFlag bool
	flagset.BoolVar(&reqKafkaSettingsAutoCreateTopicsEnableFlag, "kafka-settings.auto_create_topics_enable", false, "Enable auto creation of topics")
	var reqKafkaSettingsCompressionTypeFlag string
	flagset.StringVar(&reqKafkaSettingsCompressionTypeFlag, "kafka-settings.compression_type", "", "Specify the final compression type for a given topic. This configuration accepts the standard compression codecs ('gzip', 'snappy', 'lz4', 'zstd'). It additionally accepts 'uncompressed' which is equivalent to no compression; and 'producer' which means retain the original compression codec set by the producer.")
	var reqKafkaSettingsConnectionsMaxIdleMSFlag int
	flagset.IntVar(&reqKafkaSettingsConnectionsMaxIdleMSFlag, "kafka-settings.connections_max_idle_ms", 0, "Idle connections timeout: the server socket processor threads close the connections that idle for longer than this.")
	var reqKafkaSettingsDefaultReplicationFactorFlag int
	flagset.IntVar(&reqKafkaSettingsDefaultReplicationFactorFlag, "kafka-settings.default_replication_factor", 0, "Replication factor for autocreated topics")
	var reqKafkaSettingsGroupInitialRebalanceDelayMSFlag int
	flagset.IntVar(&reqKafkaSettingsGroupInitialRebalanceDelayMSFlag, "kafka-settings.group_initial_rebalance_delay_ms", 0, "The amount of time, in milliseconds, the group coordinator will wait for more consumers to join a new group before performing the first rebalance. A longer delay means potentially fewer rebalances, but increases the time until processing begins. The default value for this is 3 seconds. During development and testing it might be desirable to set this to 0 in order to not delay test execution time.")
	var reqKafkaSettingsGroupMaxSessionTimeoutMSFlag int
	flagset.IntVar(&reqKafkaSettingsGroupMaxSessionTimeoutMSFlag, "kafka-settings.group_max_session_timeout_ms", 0, "The maximum allowed session timeout for registered consumers. Longer timeouts give consumers more time to process messages in between heartbeats at the cost of a longer time to detect failures.")
	var reqKafkaSettingsGroupMinSessionTimeoutMSFlag int
	flagset.IntVar(&reqKafkaSettingsGroupMinSessionTimeoutMSFlag, "kafka-settings.group_min_session_timeout_ms", 0, "The minimum allowed session timeout for registered consumers. Longer timeouts give consumers more time to process messages in between heartbeats at the cost of a longer time to detect failures.")
	var reqKafkaSettingsLogCleanupAndCompactionLogCleanerDeleteRetentionMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogCleanupAndCompactionLogCleanerDeleteRetentionMSFlag, "kafka-settings.log-cleanup-and-compaction.log_cleaner_delete_retention_ms", 0, "How long are delete records retained?")
	var reqKafkaSettingsLogCleanupAndCompactionLogCleanerMaxCompactionLagMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogCleanupAndCompactionLogCleanerMaxCompactionLagMSFlag, "kafka-settings.log-cleanup-and-compaction.log_cleaner_max_compaction_lag_ms", 0, "The maximum amount of time message will remain uncompacted. Only applicable for logs that are being compacted")
	var reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCleanableRatioFlag float64
	flagset.Float64Var(&reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCleanableRatioFlag, "kafka-settings.log-cleanup-and-compaction.log_cleaner_min_cleanable_ratio", 0, "Controls log compactor frequency. Larger value means more frequent compactions but also more space wasted for logs. Consider setting log.cleaner.max.compaction.lag.ms to enforce compactions sooner, instead of setting a very high value for this option.")
	var reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCompactionLagMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCompactionLagMSFlag, "kafka-settings.log-cleanup-and-compaction.log_cleaner_min_compaction_lag_ms", 0, "The minimum time a message will remain uncompacted in the log. Only applicable for logs that are being compacted.")
	var reqKafkaSettingsLogCleanupAndCompactionLogCleanupPolicyFlag string
	flagset.StringVar(&reqKafkaSettingsLogCleanupAndCompactionLogCleanupPolicyFlag, "kafka-settings.log-cleanup-and-compaction.log_cleanup_policy", "", "The default cleanup policy for segments beyond the retention window")
	var reqKafkaSettingsLogFlushIntervalMessagesFlag int
	flagset.IntVar(&reqKafkaSettingsLogFlushIntervalMessagesFlag, "kafka-settings.log_flush_interval_messages", 0, "The number of messages accumulated on a log partition before messages are flushed to disk")
	var reqKafkaSettingsLogFlushIntervalMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogFlushIntervalMSFlag, "kafka-settings.log_flush_interval_ms", 0, "The maximum time in ms that a message in any topic is kept in memory before flushed to disk. If not set, the value in log.flush.scheduler.interval.ms is used")
	var reqKafkaSettingsLogIndexIntervalBytesFlag int
	flagset.IntVar(&reqKafkaSettingsLogIndexIntervalBytesFlag, "kafka-settings.log_index_interval_bytes", 0, "The interval with which Kafka adds an entry to the offset index")
	var reqKafkaSettingsLogIndexSizeMaxBytesFlag int
	flagset.IntVar(&reqKafkaSettingsLogIndexSizeMaxBytesFlag, "kafka-settings.log_index_size_max_bytes", 0, "The maximum size in bytes of the offset index")
	var reqKafkaSettingsLogLocalRetentionBytesFlag int
	flagset.IntVar(&reqKafkaSettingsLogLocalRetentionBytesFlag, "kafka-settings.log_local_retention_bytes", 0, "The maximum size of local log segments that can grow for a partition before it gets eligible for deletion. If set to -2, the value of log.retention.bytes is used. The effective value should always be less than or equal to log.retention.bytes value.")
	var reqKafkaSettingsLogLocalRetentionMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogLocalRetentionMSFlag, "kafka-settings.log_local_retention_ms", 0, "The number of milliseconds to keep the local log segments before it gets eligible for deletion. If set to -2, the value of log.retention.ms is used. The effective value should always be less than or equal to log.retention.ms value.")
	var reqKafkaSettingsLogMessageDownconversionEnableFlag bool
	flagset.BoolVar(&reqKafkaSettingsLogMessageDownconversionEnableFlag, "kafka-settings.log_message_downconversion_enable", false, "This configuration controls whether down-conversion of message formats is enabled to satisfy consume requests. ")
	var reqKafkaSettingsLogMessageTimestampDifferenceMaxMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogMessageTimestampDifferenceMaxMSFlag, "kafka-settings.log_message_timestamp_difference_max_ms", 0, "The maximum difference allowed between the timestamp when a broker receives a message and the timestamp specified in the message")
	var reqKafkaSettingsLogMessageTimestampTypeFlag string
	flagset.StringVar(&reqKafkaSettingsLogMessageTimestampTypeFlag, "kafka-settings.log_message_timestamp_type", "", "Define whether the timestamp in the message is message create time or log append time.")
	var reqKafkaSettingsLogPreallocateFlag bool
	flagset.BoolVar(&reqKafkaSettingsLogPreallocateFlag, "kafka-settings.log_preallocate", false, "Should pre allocate file when create new segment?")
	var reqKafkaSettingsLogRetentionBytesFlag int
	flagset.IntVar(&reqKafkaSettingsLogRetentionBytesFlag, "kafka-settings.log_retention_bytes", 0, "The maximum size of the log before deleting messages")
	var reqKafkaSettingsLogRetentionHoursFlag int
	flagset.IntVar(&reqKafkaSettingsLogRetentionHoursFlag, "kafka-settings.log_retention_hours", 0, "The number of hours to keep a log file before deleting it")
	var reqKafkaSettingsLogRetentionMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogRetentionMSFlag, "kafka-settings.log_retention_ms", 0, "The number of milliseconds to keep a log file before deleting it (in milliseconds), If not set, the value in log.retention.minutes is used. If set to -1, no time limit is applied.")
	var reqKafkaSettingsLogRollJitterMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogRollJitterMSFlag, "kafka-settings.log_roll_jitter_ms", 0, "The maximum jitter to subtract from logRollTimeMillis (in milliseconds). If not set, the value in log.roll.jitter.hours is used")
	var reqKafkaSettingsLogRollMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogRollMSFlag, "kafka-settings.log_roll_ms", 0, "The maximum time before a new log segment is rolled out (in milliseconds).")
	var reqKafkaSettingsLogSegmentBytesFlag int
	flagset.IntVar(&reqKafkaSettingsLogSegmentBytesFlag, "kafka-settings.log_segment_bytes", 0, "The maximum size of a single log file")
	var reqKafkaSettingsLogSegmentDeleteDelayMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogSegmentDeleteDelayMSFlag, "kafka-settings.log_segment_delete_delay_ms", 0, "The amount of time to wait before deleting a file from the filesystem")
	var reqKafkaSettingsMaxConnectionsPerIPFlag int
	flagset.IntVar(&reqKafkaSettingsMaxConnectionsPerIPFlag, "kafka-settings.max_connections_per_ip", 0, "The maximum number of connections allowed from each ip address (defaults to 2147483647).")
	var reqKafkaSettingsMaxIncrementalFetchSessionCacheSlotsFlag int
	flagset.IntVar(&reqKafkaSettingsMaxIncrementalFetchSessionCacheSlotsFlag, "kafka-settings.max_incremental_fetch_session_cache_slots", 0, "The maximum number of incremental fetch sessions that the broker will maintain.")
	var reqKafkaSettingsMessageMaxBytesFlag int
	flagset.IntVar(&reqKafkaSettingsMessageMaxBytesFlag, "kafka-settings.message_max_bytes", 0, "The maximum size of message that the server can receive.")
	var reqKafkaSettingsMinInsyncReplicasFlag int
	flagset.IntVar(&reqKafkaSettingsMinInsyncReplicasFlag, "kafka-settings.min_insync_replicas", 0, "When a producer sets acks to 'all' (or '-1'), min.insync.replicas specifies the minimum number of replicas that must acknowledge a write for the write to be considered successful.")
	var reqKafkaSettingsNumPartitionsFlag int
	flagset.IntVar(&reqKafkaSettingsNumPartitionsFlag, "kafka-settings.num_partitions", 0, "Number of partitions for autocreated topics")
	var reqKafkaSettingsOffsetsRetentionMinutesFlag int
	flagset.IntVar(&reqKafkaSettingsOffsetsRetentionMinutesFlag, "kafka-settings.offsets_retention_minutes", 0, "Log retention window in minutes for offsets topic")
	var reqKafkaSettingsProducerPurgatoryPurgeIntervalRequestsFlag int
	flagset.IntVar(&reqKafkaSettingsProducerPurgatoryPurgeIntervalRequestsFlag, "kafka-settings.producer_purgatory_purge_interval_requests", 0, "The purge interval (in number of requests) of the producer request purgatory(defaults to 1000).")
	var reqKafkaSettingsReplicaFetchMaxBytesFlag int
	flagset.IntVar(&reqKafkaSettingsReplicaFetchMaxBytesFlag, "kafka-settings.replica_fetch_max_bytes", 0, "The number of bytes of messages to attempt to fetch for each partition (defaults to 1048576). This is not an absolute maximum, if the first record batch in the first non-empty partition of the fetch is larger than this value, the record batch will still be returned to ensure that progress can be made.")
	var reqKafkaSettingsReplicaFetchResponseMaxBytesFlag int
	flagset.IntVar(&reqKafkaSettingsReplicaFetchResponseMaxBytesFlag, "kafka-settings.replica_fetch_response_max_bytes", 0, "Maximum bytes expected for the entire fetch response (defaults to 10485760). Records are fetched in batches, and if the first record batch in the first non-empty partition of the fetch is larger than this value, the record batch will still be returned to ensure that progress can be made. As such, this is not an absolute maximum.")
	var reqKafkaSettingsSaslOauthbearerExpectedAudienceFlag string
	flagset.StringVar(&reqKafkaSettingsSaslOauthbearerExpectedAudienceFlag, "kafka-settings.sasl_oauthbearer_expected_audience", "", "The (optional) comma-delimited setting for the broker to use to verify that the JWT was issued for one of the expected audiences.")
	var reqKafkaSettingsSaslOauthbearerExpectedIssuerFlag string
	flagset.StringVar(&reqKafkaSettingsSaslOauthbearerExpectedIssuerFlag, "kafka-settings.sasl_oauthbearer_expected_issuer", "", "Optional setting for the broker to use to verify that the JWT was created by the expected issuer.")
	var reqKafkaSettingsSaslOauthbearerJwksEndpointURLFlag string
	flagset.StringVar(&reqKafkaSettingsSaslOauthbearerJwksEndpointURLFlag, "kafka-settings.sasl_oauthbearer_jwks_endpoint_url", "", "OIDC JWKS endpoint URL. By setting this the SASL SSL OAuth2/OIDC authentication is enabled. See also other options for SASL OAuth2/OIDC. ")
	var reqKafkaSettingsSaslOauthbearerSubClaimNameFlag string
	flagset.StringVar(&reqKafkaSettingsSaslOauthbearerSubClaimNameFlag, "kafka-settings.sasl_oauthbearer_sub_claim_name", "", "Name of the scope from which to extract the subject claim from the JWT. Defaults to sub.")
	var reqKafkaSettingsSocketRequestMaxBytesFlag int
	flagset.IntVar(&reqKafkaSettingsSocketRequestMaxBytesFlag, "kafka-settings.socket_request_max_bytes", 0, "The maximum number of bytes in a socket request (defaults to 104857600).")
	var reqKafkaSettingsTransactionPartitionVerificationEnableFlag bool
	flagset.BoolVar(&reqKafkaSettingsTransactionPartitionVerificationEnableFlag, "kafka-settings.transaction_partition_verification_enable", false, "Enable verification that checks that the partition has been added to the transaction before writing transactional records to the partition")
	var reqKafkaSettingsTransactionRemoveExpiredTransactionCleanupIntervalMSFlag int
	flagset.IntVar(&reqKafkaSettingsTransactionRemoveExpiredTransactionCleanupIntervalMSFlag, "kafka-settings.transaction_remove_expired_transaction_cleanup_interval_ms", 0, "The interval at which to remove transactions that have expired due to transactional.id.expiration.ms passing (defaults to 3600000 (1 hour)).")
	var reqKafkaSettingsTransactionStateLogSegmentBytesFlag int
	flagset.IntVar(&reqKafkaSettingsTransactionStateLogSegmentBytesFlag, "kafka-settings.transaction_state_log_segment_bytes", 0, "The transaction topic segment bytes should be kept relatively small in order to facilitate faster log compaction and cache loads (defaults to 104857600 (100 mebibytes)).")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqSchemaRegistryEnabledFlag bool
	flagset.BoolVar(&reqSchemaRegistryEnabledFlag, "schema-registry-enabled", false, "Enable Schema-Registry service")
	var reqSchemaRegistrySettingsLeaderEligibilityFlag bool
	flagset.BoolVar(&reqSchemaRegistrySettingsLeaderEligibilityFlag, "schema-registry-settings.leader_eligibility", false, "If true, Karapace / Schema Registry on the service nodes can participate in leader election. It might be needed to disable this when the schemas topic is replicated to a secondary cluster and Karapace / Schema Registry there must not participate in leader election. Defaults to `true`.")
	var reqSchemaRegistrySettingsTopicNameFlag string
	flagset.StringVar(&reqSchemaRegistrySettingsTopicNameFlag, "schema-registry-settings.topic_name", "", "The durable single partition topic that acts as the durable log for the data. This topic must be compacted to avoid losing data due to retention policy. Please note that changing this configuration in an existing Schema Registry / Karapace setup leads to previous schemas being inaccessible, data encoded with them potentially unreadable and schema ID sequence put out of order. It's only possible to do the switch while Schema Registry / Karapace is disabled. Defaults to `_schemas`.")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")
	var reqVersionFlag string
	flagset.StringVar(&reqVersionFlag, "version", "", "Kafka major version")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASServiceKafkaRequest
	req.Version = reqVersionFlag
	req.TerminationProtection = reqTerminationProtectionFlag
	if reqSchemaRegistrySettingsTopicNameFlag != "" {
		req.SchemaRegistrySettings = &v3.JSONSchemaSchemaRegistry{}
		req.SchemaRegistrySettings.TopicName = reqSchemaRegistrySettingsTopicNameFlag
	}
	if reqSchemaRegistrySettingsLeaderEligibilityFlag != false {

		req.SchemaRegistrySettings.LeaderEligibility = reqSchemaRegistrySettingsLeaderEligibilityFlag
	}
	req.SchemaRegistryEnabled = reqSchemaRegistryEnabledFlag
	req.Plan = reqPlanFlag
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.CreateDBAASServiceKafkaRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}
	if reqKafkaSettingsTransactionStateLogSegmentBytesFlag != 0 {
		req.Kafk = &v3.JSONSchemaKafka{}
		req.KafkaSettings.TransactionStateLogSegmentBytes = reqKafkaSettingsTransactionStateLogSegmentBytesFlag
	}
	if reqKafkaSettingsTransactionRemoveExpiredTransactionCleanupIntervalMSFlag != 0 {

		req.KafkaSettings.TransactionRemoveExpiredTransactionCleanupIntervalMS = reqKafkaSettingsTransactionRemoveExpiredTransactionCleanupIntervalMSFlag
	}
	if reqKafkaSettingsTransactionPartitionVerificationEnableFlag != false {

		req.KafkaSettings.TransactionPartitionVerificationEnable = reqKafkaSettingsTransactionPartitionVerificationEnableFlag
	}
	if reqKafkaSettingsSocketRequestMaxBytesFlag != 0 {

		req.KafkaSettings.SocketRequestMaxBytes = reqKafkaSettingsSocketRequestMaxBytesFlag
	}
	if reqKafkaSettingsSaslOauthbearerSubClaimNameFlag != "" {

		req.KafkaSettings.SaslOauthbearerSubClaimName = reqKafkaSettingsSaslOauthbearerSubClaimNameFlag
	}
	if reqKafkaSettingsSaslOauthbearerJwksEndpointURLFlag != "" {

		req.KafkaSettings.SaslOauthbearerJwksEndpointURL = reqKafkaSettingsSaslOauthbearerJwksEndpointURLFlag
	}
	if reqKafkaSettingsSaslOauthbearerExpectedIssuerFlag != "" {

		req.KafkaSettings.SaslOauthbearerExpectedIssuer = reqKafkaSettingsSaslOauthbearerExpectedIssuerFlag
	}
	if reqKafkaSettingsSaslOauthbearerExpectedAudienceFlag != "" {

		req.KafkaSettings.SaslOauthbearerExpectedAudience = reqKafkaSettingsSaslOauthbearerExpectedAudienceFlag
	}
	if reqKafkaSettingsReplicaFetchResponseMaxBytesFlag != 0 {

		req.KafkaSettings.ReplicaFetchResponseMaxBytes = reqKafkaSettingsReplicaFetchResponseMaxBytesFlag
	}
	if reqKafkaSettingsReplicaFetchMaxBytesFlag != 0 {

		req.KafkaSettings.ReplicaFetchMaxBytes = reqKafkaSettingsReplicaFetchMaxBytesFlag
	}
	if reqKafkaSettingsProducerPurgatoryPurgeIntervalRequestsFlag != 0 {

		req.KafkaSettings.ProducerPurgatoryPurgeIntervalRequests = reqKafkaSettingsProducerPurgatoryPurgeIntervalRequestsFlag
	}
	if reqKafkaSettingsOffsetsRetentionMinutesFlag != 0 {

		req.KafkaSettings.OffsetsRetentionMinutes = reqKafkaSettingsOffsetsRetentionMinutesFlag
	}
	if reqKafkaSettingsNumPartitionsFlag != 0 {

		req.KafkaSettings.NumPartitions = reqKafkaSettingsNumPartitionsFlag
	}
	if reqKafkaSettingsMinInsyncReplicasFlag != 0 {

		req.KafkaSettings.MinInsyncReplicas = reqKafkaSettingsMinInsyncReplicasFlag
	}
	if reqKafkaSettingsMessageMaxBytesFlag != 0 {

		req.KafkaSettings.MessageMaxBytes = reqKafkaSettingsMessageMaxBytesFlag
	}
	if reqKafkaSettingsMaxIncrementalFetchSessionCacheSlotsFlag != 0 {

		req.KafkaSettings.MaxIncrementalFetchSessionCacheSlots = reqKafkaSettingsMaxIncrementalFetchSessionCacheSlotsFlag
	}
	if reqKafkaSettingsMaxConnectionsPerIPFlag != 0 {

		req.KafkaSettings.MaxConnectionsPerIP = reqKafkaSettingsMaxConnectionsPerIPFlag
	}
	if reqKafkaSettingsLogSegmentDeleteDelayMSFlag != 0 {

		req.KafkaSettings.LogSegmentDeleteDelayMS = reqKafkaSettingsLogSegmentDeleteDelayMSFlag
	}
	if reqKafkaSettingsLogSegmentBytesFlag != 0 {

		req.KafkaSettings.LogSegmentBytes = reqKafkaSettingsLogSegmentBytesFlag
	}
	if reqKafkaSettingsLogRollMSFlag != 0 {

		req.KafkaSettings.LogRollMS = reqKafkaSettingsLogRollMSFlag
	}
	if reqKafkaSettingsLogRollJitterMSFlag != 0 {

		req.KafkaSettings.LogRollJitterMS = reqKafkaSettingsLogRollJitterMSFlag
	}
	if reqKafkaSettingsLogRetentionMSFlag != 0 {

		req.KafkaSettings.LogRetentionMS = reqKafkaSettingsLogRetentionMSFlag
	}
	if reqKafkaSettingsLogRetentionHoursFlag != 0 {

		req.KafkaSettings.LogRetentionHours = reqKafkaSettingsLogRetentionHoursFlag
	}
	if reqKafkaSettingsLogRetentionBytesFlag != 0 {

		req.KafkaSettings.LogRetentionBytes = reqKafkaSettingsLogRetentionBytesFlag
	}
	if reqKafkaSettingsLogPreallocateFlag != false {

		req.KafkaSettings.LogPreallocate = reqKafkaSettingsLogPreallocateFlag
	}
	if reqKafkaSettingsLogMessageTimestampTypeFlag != "" {

		req.KafkaSettings.LogMessageTimestampType = reqKafkaSettingsLogMessageTimestampTypeFlag
	}
	if reqKafkaSettingsLogMessageTimestampDifferenceMaxMSFlag != 0 {

		req.KafkaSettings.LogMessageTimestampDifferenceMaxMS = reqKafkaSettingsLogMessageTimestampDifferenceMaxMSFlag
	}
	if reqKafkaSettingsLogMessageDownconversionEnableFlag != false {

		req.KafkaSettings.LogMessageDownconversionEnable = reqKafkaSettingsLogMessageDownconversionEnableFlag
	}
	if reqKafkaSettingsLogLocalRetentionMSFlag != 0 {

		req.KafkaSettings.LogLocalRetentionMS = reqKafkaSettingsLogLocalRetentionMSFlag
	}
	if reqKafkaSettingsLogLocalRetentionBytesFlag != 0 {

		req.KafkaSettings.LogLocalRetentionBytes = reqKafkaSettingsLogLocalRetentionBytesFlag
	}
	if reqKafkaSettingsLogIndexSizeMaxBytesFlag != 0 {

		req.KafkaSettings.LogIndexSizeMaxBytes = reqKafkaSettingsLogIndexSizeMaxBytesFlag
	}
	if reqKafkaSettingsLogIndexIntervalBytesFlag != 0 {

		req.KafkaSettings.LogIndexIntervalBytes = reqKafkaSettingsLogIndexIntervalBytesFlag
	}
	if reqKafkaSettingsLogFlushIntervalMSFlag != 0 {

		req.KafkaSettings.LogFlushIntervalMS = reqKafkaSettingsLogFlushIntervalMSFlag
	}
	if reqKafkaSettingsLogFlushIntervalMessagesFlag != 0 {

		req.KafkaSettings.LogFlushIntervalMessages = reqKafkaSettingsLogFlushIntervalMessagesFlag
	}
	if reqKafkaSettingsLogCleanupAndCompactionLogCleanupPolicyFlag != "" {
		req.KafkaSettingsLogCleanupAndCompact = &v3.KafkaSettingsLogCleanupAndCompaction{}
		req.KafkaSettingsLogCleanupAndCompaction.LogCleanupPolicy = reqKafkaSettingsLogCleanupAndCompactionLogCleanupPolicyFlag
	}
	if reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCompactionLagMSFlag != 0 {

		req.KafkaSettingsLogCleanupAndCompaction.LogCleanerMinCompactionLagMS = reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCompactionLagMSFlag
	}
	if reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCleanableRatioFlag != 0 {

		req.KafkaSettingsLogCleanupAndCompaction.LogCleanerMinCleanableRatio = reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCleanableRatioFlag
	}
	if reqKafkaSettingsLogCleanupAndCompactionLogCleanerMaxCompactionLagMSFlag != 0 {

		req.KafkaSettingsLogCleanupAndCompaction.LogCleanerMaxCompactionLagMS = reqKafkaSettingsLogCleanupAndCompactionLogCleanerMaxCompactionLagMSFlag
	}
	if reqKafkaSettingsLogCleanupAndCompactionLogCleanerDeleteRetentionMSFlag != 0 {

		req.KafkaSettingsLogCleanupAndCompaction.LogCleanerDeleteRetentionMS = reqKafkaSettingsLogCleanupAndCompactionLogCleanerDeleteRetentionMSFlag
	}
	if reqKafkaSettingsGroupMinSessionTimeoutMSFlag != 0 {
		req.KafkaSetting = &v3.JSONSchemaKafka{}
		req.KafkaSettings.GroupMinSessionTimeoutMS = reqKafkaSettingsGroupMinSessionTimeoutMSFlag
	}
	if reqKafkaSettingsGroupMaxSessionTimeoutMSFlag != 0 {

		req.KafkaSettings.GroupMaxSessionTimeoutMS = reqKafkaSettingsGroupMaxSessionTimeoutMSFlag
	}
	if reqKafkaSettingsGroupInitialRebalanceDelayMSFlag != 0 {

		req.KafkaSettings.GroupInitialRebalanceDelayMS = reqKafkaSettingsGroupInitialRebalanceDelayMSFlag
	}
	if reqKafkaSettingsDefaultReplicationFactorFlag != 0 {

		req.KafkaSettings.DefaultReplicationFactor = reqKafkaSettingsDefaultReplicationFactorFlag
	}
	if reqKafkaSettingsConnectionsMaxIdleMSFlag != 0 {

		req.KafkaSettings.ConnectionsMaxIdleMS = reqKafkaSettingsConnectionsMaxIdleMSFlag
	}
	if reqKafkaSettingsCompressionTypeFlag != "" {

		req.KafkaSettings.CompressionType = reqKafkaSettingsCompressionTypeFlag
	}
	if reqKafkaSettingsAutoCreateTopicsEnableFlag != false {

		req.KafkaSettings.AutoCreateTopicsEnable = reqKafkaSettingsAutoCreateTopicsEnableFlag
	}
	if reqKafkaRestSettingsSimpleconsumerPoolSizeMaxFlag != 0 {
		req.KafkaRestSetting = &v3.JSONSchemaKafkaRest{}
		req.KafkaRestSettings.SimpleconsumerPoolSizeMax = reqKafkaRestSettingsSimpleconsumerPoolSizeMaxFlag
	}
	if reqKafkaRestSettingsProducerMaxRequestSizeFlag != 0 {

		req.KafkaRestSettings.ProducerMaxRequestSize = reqKafkaRestSettingsProducerMaxRequestSizeFlag
	}
	if reqKafkaRestSettingsProducerLingerMSFlag != 0 {

		req.KafkaRestSettings.ProducerLingerMS = reqKafkaRestSettingsProducerLingerMSFlag
	}
	if reqKafkaRestSettingsProducerCompressionTypeFlag != "" {

		req.KafkaRestSettings.ProducerCompressionType = reqKafkaRestSettingsProducerCompressionTypeFlag
	}
	if reqKafkaRestSettingsProducerAcksFlag != "" {

		req.KafkaRestSettings.ProducerAcks = reqKafkaRestSettingsProducerAcksFlag
	}
	if reqKafkaRestSettingsNameStrategyValidationFlag != false {

		req.KafkaRestSettings.NameStrategyValidation = reqKafkaRestSettingsNameStrategyValidationFlag
	}
	if reqKafkaRestSettingsNameStrategyFlag != "" {

		req.KafkaRestSettings.NameStrategy = reqKafkaRestSettingsNameStrategyFlag
	}
	if reqKafkaRestSettingsConsumerRequestTimeoutMSFlag != 0 {

		req.KafkaRestSettings.ConsumerRequestTimeoutMS = reqKafkaRestSettingsConsumerRequestTimeoutMSFlag
	}
	if reqKafkaRestSettingsConsumerRequestMaxBytesFlag != 0 {

		req.KafkaRestSettings.ConsumerRequestMaxBytes = reqKafkaRestSettingsConsumerRequestMaxBytesFlag
	}
	if reqKafkaRestSettingsConsumerEnableAutoCommitFlag != false {

		req.KafkaRestSettings.ConsumerEnableAutoCommit = reqKafkaRestSettingsConsumerEnableAutoCommitFlag
	}
	req.KafkaRestEnabled = reqKafkaRestEnabledFlag
	if reqKafkaConnectSettingsSessionTimeoutMSFlag != 0 {
		req.KafkaConnectSetting = &v3.JSONSchemaKafkaConnect{}
		req.KafkaConnectSettings.SessionTimeoutMS = reqKafkaConnectSettingsSessionTimeoutMSFlag
	}
	if reqKafkaConnectSettingsScheduledRebalanceMaxDelayMSFlag != 0 {

		req.KafkaConnectSettings.ScheduledRebalanceMaxDelayMS = reqKafkaConnectSettingsScheduledRebalanceMaxDelayMSFlag
	}
	if reqKafkaConnectSettingsProducerMaxRequestSizeFlag != 0 {

		req.KafkaConnectSettings.ProducerMaxRequestSize = reqKafkaConnectSettingsProducerMaxRequestSizeFlag
	}
	if reqKafkaConnectSettingsProducerLingerMSFlag != 0 {

		req.KafkaConnectSettings.ProducerLingerMS = reqKafkaConnectSettingsProducerLingerMSFlag
	}
	if reqKafkaConnectSettingsProducerCompressionTypeFlag != "" {

		req.KafkaConnectSettings.ProducerCompressionType = reqKafkaConnectSettingsProducerCompressionTypeFlag
	}
	if reqKafkaConnectSettingsProducerBufferMemoryFlag != 0 {

		req.KafkaConnectSettings.ProducerBufferMemory = reqKafkaConnectSettingsProducerBufferMemoryFlag
	}
	if reqKafkaConnectSettingsProducerBatchSizeFlag != 0 {

		req.KafkaConnectSettings.ProducerBatchSize = reqKafkaConnectSettingsProducerBatchSizeFlag
	}
	if reqKafkaConnectSettingsOffsetFlushTimeoutMSFlag != 0 {

		req.KafkaConnectSettings.OffsetFlushTimeoutMS = reqKafkaConnectSettingsOffsetFlushTimeoutMSFlag
	}
	if reqKafkaConnectSettingsOffsetFlushIntervalMSFlag != 0 {

		req.KafkaConnectSettings.OffsetFlushIntervalMS = reqKafkaConnectSettingsOffsetFlushIntervalMSFlag
	}
	if reqKafkaConnectSettingsConsumerMaxPollRecordsFlag != 0 {

		req.KafkaConnectSettings.ConsumerMaxPollRecords = reqKafkaConnectSettingsConsumerMaxPollRecordsFlag
	}
	if reqKafkaConnectSettingsConsumerMaxPollIntervalMSFlag != 0 {

		req.KafkaConnectSettings.ConsumerMaxPollIntervalMS = reqKafkaConnectSettingsConsumerMaxPollIntervalMSFlag
	}
	if reqKafkaConnectSettingsConsumerMaxPartitionFetchBytesFlag != 0 {

		req.KafkaConnectSettings.ConsumerMaxPartitionFetchBytes = reqKafkaConnectSettingsConsumerMaxPartitionFetchBytesFlag
	}
	if reqKafkaConnectSettingsConsumerIsolationLevelFlag != "" {

		req.KafkaConnectSettings.ConsumerIsolationLevel = reqKafkaConnectSettingsConsumerIsolationLevelFlag
	}
	if reqKafkaConnectSettingsConsumerFetchMaxBytesFlag != 0 {

		req.KafkaConnectSettings.ConsumerFetchMaxBytes = reqKafkaConnectSettingsConsumerFetchMaxBytesFlag
	}
	if reqKafkaConnectSettingsConsumerAutoOffsetResetFlag != "" {

		req.KafkaConnectSettings.ConsumerAutoOffsetReset = reqKafkaConnectSettingsConsumerAutoOffsetResetFlag
	}
	if reqKafkaConnectSettingsConnectorClientConfigOverridePolicyFlag != "" {

		req.KafkaConnectSettings.ConnectorClientConfigOverridePolicy = reqKafkaConnectSettingsConnectorClientConfigOverridePolicyFlag
	}
	req.KafkaConnectEnabled = reqKafkaConnectEnabledFlag
	if reqAuthenticationMethodsSaslFlag != false {
		req.AuthenticationMethod = &v3.CreateDBAASServiceKafkaRequestAuthenticationMethods{}
		req.AuthenticationMethods.Sasl = reqAuthenticationMethodsSaslFlag
	}
	if reqAuthenticationMethodsCertificateFlag != false {

		req.AuthenticationMethods.Certificate = reqAuthenticationMethodsCertificateFlag
	}

	resp, err := client.CreateDBAASServiceKafka(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASServiceKafkaCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-service-kafka", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqAuthenticationMethodsCertificateFlag bool
	flagset.BoolVar(&reqAuthenticationMethodsCertificateFlag, "authentication-methods.certificate", false, "Enable certificate/SSL authentication")
	var reqAuthenticationMethodsSaslFlag bool
	flagset.BoolVar(&reqAuthenticationMethodsSaslFlag, "authentication-methods.sasl", false, "Enable SASL authentication")
	var reqKafkaConnectEnabledFlag bool
	flagset.BoolVar(&reqKafkaConnectEnabledFlag, "kafka-connect-enabled", false, "Allow clients to connect to kafka_connect from the public internet for service nodes that are in a project VPC or another type of private network")
	var reqKafkaConnectSettingsConnectorClientConfigOverridePolicyFlag string
	flagset.StringVar(&reqKafkaConnectSettingsConnectorClientConfigOverridePolicyFlag, "kafka-connect-settings.connector_client_config_override_policy", "", "Defines what client configurations can be overridden by the connector. Default is None")
	var reqKafkaConnectSettingsConsumerAutoOffsetResetFlag string
	flagset.StringVar(&reqKafkaConnectSettingsConsumerAutoOffsetResetFlag, "kafka-connect-settings.consumer_auto_offset_reset", "", "What to do when there is no initial offset in Kafka or if the current offset does not exist any more on the server. Default is earliest")
	var reqKafkaConnectSettingsConsumerFetchMaxBytesFlag int
	flagset.IntVar(&reqKafkaConnectSettingsConsumerFetchMaxBytesFlag, "kafka-connect-settings.consumer_fetch_max_bytes", 0, "Records are fetched in batches by the consumer, and if the first record batch in the first non-empty partition of the fetch is larger than this value, the record batch will still be returned to ensure that the consumer can make progress. As such, this is not a absolute maximum.")
	var reqKafkaConnectSettingsConsumerIsolationLevelFlag string
	flagset.StringVar(&reqKafkaConnectSettingsConsumerIsolationLevelFlag, "kafka-connect-settings.consumer_isolation_level", "", "Transaction read isolation level. read_uncommitted is the default, but read_committed can be used if consume-exactly-once behavior is desired.")
	var reqKafkaConnectSettingsConsumerMaxPartitionFetchBytesFlag int
	flagset.IntVar(&reqKafkaConnectSettingsConsumerMaxPartitionFetchBytesFlag, "kafka-connect-settings.consumer_max_partition_fetch_bytes", 0, "Records are fetched in batches by the consumer.If the first record batch in the first non-empty partition of the fetch is larger than this limit, the batch will still be returned to ensure that the consumer can make progress. ")
	var reqKafkaConnectSettingsConsumerMaxPollIntervalMSFlag int
	flagset.IntVar(&reqKafkaConnectSettingsConsumerMaxPollIntervalMSFlag, "kafka-connect-settings.consumer_max_poll_interval_ms", 0, "The maximum delay in milliseconds between invocations of poll() when using consumer group management (defaults to 300000).")
	var reqKafkaConnectSettingsConsumerMaxPollRecordsFlag int
	flagset.IntVar(&reqKafkaConnectSettingsConsumerMaxPollRecordsFlag, "kafka-connect-settings.consumer_max_poll_records", 0, "The maximum number of records returned in a single call to poll() (defaults to 500).")
	var reqKafkaConnectSettingsOffsetFlushIntervalMSFlag int
	flagset.IntVar(&reqKafkaConnectSettingsOffsetFlushIntervalMSFlag, "kafka-connect-settings.offset_flush_interval_ms", 0, "The interval at which to try committing offsets for tasks (defaults to 60000).")
	var reqKafkaConnectSettingsOffsetFlushTimeoutMSFlag int
	flagset.IntVar(&reqKafkaConnectSettingsOffsetFlushTimeoutMSFlag, "kafka-connect-settings.offset_flush_timeout_ms", 0, "Maximum number of milliseconds to wait for records to flush and partition offset data to be committed to offset storage before cancelling the process and restoring the offset data to be committed in a future attempt (defaults to 5000).")
	var reqKafkaConnectSettingsProducerBatchSizeFlag int
	flagset.IntVar(&reqKafkaConnectSettingsProducerBatchSizeFlag, "kafka-connect-settings.producer_batch_size", 0, "This setting gives the upper bound of the batch size to be sent. If there are fewer than this many bytes accumulated for this partition, the producer will 'linger' for the linger.ms time waiting for more records to show up. A batch size of zero will disable batching entirely (defaults to 16384).")
	var reqKafkaConnectSettingsProducerBufferMemoryFlag int
	flagset.IntVar(&reqKafkaConnectSettingsProducerBufferMemoryFlag, "kafka-connect-settings.producer_buffer_memory", 0, "The total bytes of memory the producer can use to buffer records waiting to be sent to the broker (defaults to 33554432).")
	var reqKafkaConnectSettingsProducerCompressionTypeFlag string
	flagset.StringVar(&reqKafkaConnectSettingsProducerCompressionTypeFlag, "kafka-connect-settings.producer_compression_type", "", "Specify the default compression type for producers. This configuration accepts the standard compression codecs ('gzip', 'snappy', 'lz4', 'zstd'). It additionally accepts 'none' which is the default and equivalent to no compression.")
	var reqKafkaConnectSettingsProducerLingerMSFlag int
	flagset.IntVar(&reqKafkaConnectSettingsProducerLingerMSFlag, "kafka-connect-settings.producer_linger_ms", 0, "This setting gives the upper bound on the delay for batching: once there is batch.size worth of records for a partition it will be sent immediately regardless of this setting, however if there are fewer than this many bytes accumulated for this partition the producer will 'linger' for the specified time waiting for more records to show up. Defaults to 0.")
	var reqKafkaConnectSettingsProducerMaxRequestSizeFlag int
	flagset.IntVar(&reqKafkaConnectSettingsProducerMaxRequestSizeFlag, "kafka-connect-settings.producer_max_request_size", 0, "This setting will limit the number of record batches the producer will send in a single request to avoid sending huge requests.")
	var reqKafkaConnectSettingsScheduledRebalanceMaxDelayMSFlag int
	flagset.IntVar(&reqKafkaConnectSettingsScheduledRebalanceMaxDelayMSFlag, "kafka-connect-settings.scheduled_rebalance_max_delay_ms", 0, "The maximum delay that is scheduled in order to wait for the return of one or more departed workers before rebalancing and reassigning their connectors and tasks to the group. During this period the connectors and tasks of the departed workers remain unassigned. Defaults to 5 minutes.")
	var reqKafkaConnectSettingsSessionTimeoutMSFlag int
	flagset.IntVar(&reqKafkaConnectSettingsSessionTimeoutMSFlag, "kafka-connect-settings.session_timeout_ms", 0, "The timeout in milliseconds used to detect failures when using Kafka’s group management facilities (defaults to 10000).")
	var reqKafkaRestEnabledFlag bool
	flagset.BoolVar(&reqKafkaRestEnabledFlag, "kafka-rest-enabled", false, "Enable Kafka-REST service")
	var reqKafkaRestSettingsConsumerEnableAutoCommitFlag bool
	flagset.BoolVar(&reqKafkaRestSettingsConsumerEnableAutoCommitFlag, "kafka-rest-settings.consumer_enable_auto_commit", false, "If true the consumer's offset will be periodically committed to Kafka in the background")
	var reqKafkaRestSettingsConsumerRequestMaxBytesFlag int
	flagset.IntVar(&reqKafkaRestSettingsConsumerRequestMaxBytesFlag, "kafka-rest-settings.consumer_request_max_bytes", 0, "Maximum number of bytes in unencoded message keys and values by a single request")
	var reqKafkaRestSettingsConsumerRequestTimeoutMSFlag int
	flagset.IntVar(&reqKafkaRestSettingsConsumerRequestTimeoutMSFlag, "kafka-rest-settings.consumer_request_timeout_ms", 0, "The maximum total time to wait for messages for a request if the maximum number of messages has not yet been reached")
	var reqKafkaRestSettingsNameStrategyFlag string
	flagset.StringVar(&reqKafkaRestSettingsNameStrategyFlag, "kafka-rest-settings.name_strategy", "", "Name strategy to use when selecting subject for storing schemas")
	var reqKafkaRestSettingsNameStrategyValidationFlag bool
	flagset.BoolVar(&reqKafkaRestSettingsNameStrategyValidationFlag, "kafka-rest-settings.name_strategy_validation", false, "If true, validate that given schema is registered under expected subject name by the used name strategy when producing messages.")
	var reqKafkaRestSettingsProducerAcksFlag string
	flagset.StringVar(&reqKafkaRestSettingsProducerAcksFlag, "kafka-rest-settings.producer_acks", "", "The number of acknowledgments the producer requires the leader to have received before considering a request complete. If set to 'all' or '-1', the leader will wait for the full set of in-sync replicas to acknowledge the record.")
	var reqKafkaRestSettingsProducerCompressionTypeFlag string
	flagset.StringVar(&reqKafkaRestSettingsProducerCompressionTypeFlag, "kafka-rest-settings.producer_compression_type", "", "Specify the default compression type for producers. This configuration accepts the standard compression codecs ('gzip', 'snappy', 'lz4', 'zstd'). It additionally accepts 'none' which is the default and equivalent to no compression.")
	var reqKafkaRestSettingsProducerLingerMSFlag int
	flagset.IntVar(&reqKafkaRestSettingsProducerLingerMSFlag, "kafka-rest-settings.producer_linger_ms", 0, "Wait for up to the given delay to allow batching records together")
	var reqKafkaRestSettingsProducerMaxRequestSizeFlag int
	flagset.IntVar(&reqKafkaRestSettingsProducerMaxRequestSizeFlag, "kafka-rest-settings.producer_max_request_size", 0, "The maximum size of a request in bytes. Note that Kafka broker can also cap the record batch size.")
	var reqKafkaRestSettingsSimpleconsumerPoolSizeMaxFlag int
	flagset.IntVar(&reqKafkaRestSettingsSimpleconsumerPoolSizeMaxFlag, "kafka-rest-settings.simpleconsumer_pool_size_max", 0, "Maximum number of SimpleConsumers that can be instantiated per broker")
	var reqKafkaSettingsAutoCreateTopicsEnableFlag bool
	flagset.BoolVar(&reqKafkaSettingsAutoCreateTopicsEnableFlag, "kafka-settings.auto_create_topics_enable", false, "Enable auto creation of topics")
	var reqKafkaSettingsCompressionTypeFlag string
	flagset.StringVar(&reqKafkaSettingsCompressionTypeFlag, "kafka-settings.compression_type", "", "Specify the final compression type for a given topic. This configuration accepts the standard compression codecs ('gzip', 'snappy', 'lz4', 'zstd'). It additionally accepts 'uncompressed' which is equivalent to no compression; and 'producer' which means retain the original compression codec set by the producer.")
	var reqKafkaSettingsConnectionsMaxIdleMSFlag int
	flagset.IntVar(&reqKafkaSettingsConnectionsMaxIdleMSFlag, "kafka-settings.connections_max_idle_ms", 0, "Idle connections timeout: the server socket processor threads close the connections that idle for longer than this.")
	var reqKafkaSettingsDefaultReplicationFactorFlag int
	flagset.IntVar(&reqKafkaSettingsDefaultReplicationFactorFlag, "kafka-settings.default_replication_factor", 0, "Replication factor for autocreated topics")
	var reqKafkaSettingsGroupInitialRebalanceDelayMSFlag int
	flagset.IntVar(&reqKafkaSettingsGroupInitialRebalanceDelayMSFlag, "kafka-settings.group_initial_rebalance_delay_ms", 0, "The amount of time, in milliseconds, the group coordinator will wait for more consumers to join a new group before performing the first rebalance. A longer delay means potentially fewer rebalances, but increases the time until processing begins. The default value for this is 3 seconds. During development and testing it might be desirable to set this to 0 in order to not delay test execution time.")
	var reqKafkaSettingsGroupMaxSessionTimeoutMSFlag int
	flagset.IntVar(&reqKafkaSettingsGroupMaxSessionTimeoutMSFlag, "kafka-settings.group_max_session_timeout_ms", 0, "The maximum allowed session timeout for registered consumers. Longer timeouts give consumers more time to process messages in between heartbeats at the cost of a longer time to detect failures.")
	var reqKafkaSettingsGroupMinSessionTimeoutMSFlag int
	flagset.IntVar(&reqKafkaSettingsGroupMinSessionTimeoutMSFlag, "kafka-settings.group_min_session_timeout_ms", 0, "The minimum allowed session timeout for registered consumers. Longer timeouts give consumers more time to process messages in between heartbeats at the cost of a longer time to detect failures.")
	var reqKafkaSettingsLogCleanupAndCompactionLogCleanerDeleteRetentionMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogCleanupAndCompactionLogCleanerDeleteRetentionMSFlag, "kafka-settings.log-cleanup-and-compaction.log_cleaner_delete_retention_ms", 0, "How long are delete records retained?")
	var reqKafkaSettingsLogCleanupAndCompactionLogCleanerMaxCompactionLagMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogCleanupAndCompactionLogCleanerMaxCompactionLagMSFlag, "kafka-settings.log-cleanup-and-compaction.log_cleaner_max_compaction_lag_ms", 0, "The maximum amount of time message will remain uncompacted. Only applicable for logs that are being compacted")
	var reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCleanableRatioFlag float64
	flagset.Float64Var(&reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCleanableRatioFlag, "kafka-settings.log-cleanup-and-compaction.log_cleaner_min_cleanable_ratio", 0, "Controls log compactor frequency. Larger value means more frequent compactions but also more space wasted for logs. Consider setting log.cleaner.max.compaction.lag.ms to enforce compactions sooner, instead of setting a very high value for this option.")
	var reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCompactionLagMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCompactionLagMSFlag, "kafka-settings.log-cleanup-and-compaction.log_cleaner_min_compaction_lag_ms", 0, "The minimum time a message will remain uncompacted in the log. Only applicable for logs that are being compacted.")
	var reqKafkaSettingsLogCleanupAndCompactionLogCleanupPolicyFlag string
	flagset.StringVar(&reqKafkaSettingsLogCleanupAndCompactionLogCleanupPolicyFlag, "kafka-settings.log-cleanup-and-compaction.log_cleanup_policy", "", "The default cleanup policy for segments beyond the retention window")
	var reqKafkaSettingsLogFlushIntervalMessagesFlag int
	flagset.IntVar(&reqKafkaSettingsLogFlushIntervalMessagesFlag, "kafka-settings.log_flush_interval_messages", 0, "The number of messages accumulated on a log partition before messages are flushed to disk")
	var reqKafkaSettingsLogFlushIntervalMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogFlushIntervalMSFlag, "kafka-settings.log_flush_interval_ms", 0, "The maximum time in ms that a message in any topic is kept in memory before flushed to disk. If not set, the value in log.flush.scheduler.interval.ms is used")
	var reqKafkaSettingsLogIndexIntervalBytesFlag int
	flagset.IntVar(&reqKafkaSettingsLogIndexIntervalBytesFlag, "kafka-settings.log_index_interval_bytes", 0, "The interval with which Kafka adds an entry to the offset index")
	var reqKafkaSettingsLogIndexSizeMaxBytesFlag int
	flagset.IntVar(&reqKafkaSettingsLogIndexSizeMaxBytesFlag, "kafka-settings.log_index_size_max_bytes", 0, "The maximum size in bytes of the offset index")
	var reqKafkaSettingsLogLocalRetentionBytesFlag int
	flagset.IntVar(&reqKafkaSettingsLogLocalRetentionBytesFlag, "kafka-settings.log_local_retention_bytes", 0, "The maximum size of local log segments that can grow for a partition before it gets eligible for deletion. If set to -2, the value of log.retention.bytes is used. The effective value should always be less than or equal to log.retention.bytes value.")
	var reqKafkaSettingsLogLocalRetentionMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogLocalRetentionMSFlag, "kafka-settings.log_local_retention_ms", 0, "The number of milliseconds to keep the local log segments before it gets eligible for deletion. If set to -2, the value of log.retention.ms is used. The effective value should always be less than or equal to log.retention.ms value.")
	var reqKafkaSettingsLogMessageDownconversionEnableFlag bool
	flagset.BoolVar(&reqKafkaSettingsLogMessageDownconversionEnableFlag, "kafka-settings.log_message_downconversion_enable", false, "This configuration controls whether down-conversion of message formats is enabled to satisfy consume requests. ")
	var reqKafkaSettingsLogMessageTimestampDifferenceMaxMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogMessageTimestampDifferenceMaxMSFlag, "kafka-settings.log_message_timestamp_difference_max_ms", 0, "The maximum difference allowed between the timestamp when a broker receives a message and the timestamp specified in the message")
	var reqKafkaSettingsLogMessageTimestampTypeFlag string
	flagset.StringVar(&reqKafkaSettingsLogMessageTimestampTypeFlag, "kafka-settings.log_message_timestamp_type", "", "Define whether the timestamp in the message is message create time or log append time.")
	var reqKafkaSettingsLogPreallocateFlag bool
	flagset.BoolVar(&reqKafkaSettingsLogPreallocateFlag, "kafka-settings.log_preallocate", false, "Should pre allocate file when create new segment?")
	var reqKafkaSettingsLogRetentionBytesFlag int
	flagset.IntVar(&reqKafkaSettingsLogRetentionBytesFlag, "kafka-settings.log_retention_bytes", 0, "The maximum size of the log before deleting messages")
	var reqKafkaSettingsLogRetentionHoursFlag int
	flagset.IntVar(&reqKafkaSettingsLogRetentionHoursFlag, "kafka-settings.log_retention_hours", 0, "The number of hours to keep a log file before deleting it")
	var reqKafkaSettingsLogRetentionMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogRetentionMSFlag, "kafka-settings.log_retention_ms", 0, "The number of milliseconds to keep a log file before deleting it (in milliseconds), If not set, the value in log.retention.minutes is used. If set to -1, no time limit is applied.")
	var reqKafkaSettingsLogRollJitterMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogRollJitterMSFlag, "kafka-settings.log_roll_jitter_ms", 0, "The maximum jitter to subtract from logRollTimeMillis (in milliseconds). If not set, the value in log.roll.jitter.hours is used")
	var reqKafkaSettingsLogRollMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogRollMSFlag, "kafka-settings.log_roll_ms", 0, "The maximum time before a new log segment is rolled out (in milliseconds).")
	var reqKafkaSettingsLogSegmentBytesFlag int
	flagset.IntVar(&reqKafkaSettingsLogSegmentBytesFlag, "kafka-settings.log_segment_bytes", 0, "The maximum size of a single log file")
	var reqKafkaSettingsLogSegmentDeleteDelayMSFlag int
	flagset.IntVar(&reqKafkaSettingsLogSegmentDeleteDelayMSFlag, "kafka-settings.log_segment_delete_delay_ms", 0, "The amount of time to wait before deleting a file from the filesystem")
	var reqKafkaSettingsMaxConnectionsPerIPFlag int
	flagset.IntVar(&reqKafkaSettingsMaxConnectionsPerIPFlag, "kafka-settings.max_connections_per_ip", 0, "The maximum number of connections allowed from each ip address (defaults to 2147483647).")
	var reqKafkaSettingsMaxIncrementalFetchSessionCacheSlotsFlag int
	flagset.IntVar(&reqKafkaSettingsMaxIncrementalFetchSessionCacheSlotsFlag, "kafka-settings.max_incremental_fetch_session_cache_slots", 0, "The maximum number of incremental fetch sessions that the broker will maintain.")
	var reqKafkaSettingsMessageMaxBytesFlag int
	flagset.IntVar(&reqKafkaSettingsMessageMaxBytesFlag, "kafka-settings.message_max_bytes", 0, "The maximum size of message that the server can receive.")
	var reqKafkaSettingsMinInsyncReplicasFlag int
	flagset.IntVar(&reqKafkaSettingsMinInsyncReplicasFlag, "kafka-settings.min_insync_replicas", 0, "When a producer sets acks to 'all' (or '-1'), min.insync.replicas specifies the minimum number of replicas that must acknowledge a write for the write to be considered successful.")
	var reqKafkaSettingsNumPartitionsFlag int
	flagset.IntVar(&reqKafkaSettingsNumPartitionsFlag, "kafka-settings.num_partitions", 0, "Number of partitions for autocreated topics")
	var reqKafkaSettingsOffsetsRetentionMinutesFlag int
	flagset.IntVar(&reqKafkaSettingsOffsetsRetentionMinutesFlag, "kafka-settings.offsets_retention_minutes", 0, "Log retention window in minutes for offsets topic")
	var reqKafkaSettingsProducerPurgatoryPurgeIntervalRequestsFlag int
	flagset.IntVar(&reqKafkaSettingsProducerPurgatoryPurgeIntervalRequestsFlag, "kafka-settings.producer_purgatory_purge_interval_requests", 0, "The purge interval (in number of requests) of the producer request purgatory(defaults to 1000).")
	var reqKafkaSettingsReplicaFetchMaxBytesFlag int
	flagset.IntVar(&reqKafkaSettingsReplicaFetchMaxBytesFlag, "kafka-settings.replica_fetch_max_bytes", 0, "The number of bytes of messages to attempt to fetch for each partition (defaults to 1048576). This is not an absolute maximum, if the first record batch in the first non-empty partition of the fetch is larger than this value, the record batch will still be returned to ensure that progress can be made.")
	var reqKafkaSettingsReplicaFetchResponseMaxBytesFlag int
	flagset.IntVar(&reqKafkaSettingsReplicaFetchResponseMaxBytesFlag, "kafka-settings.replica_fetch_response_max_bytes", 0, "Maximum bytes expected for the entire fetch response (defaults to 10485760). Records are fetched in batches, and if the first record batch in the first non-empty partition of the fetch is larger than this value, the record batch will still be returned to ensure that progress can be made. As such, this is not an absolute maximum.")
	var reqKafkaSettingsSaslOauthbearerExpectedAudienceFlag string
	flagset.StringVar(&reqKafkaSettingsSaslOauthbearerExpectedAudienceFlag, "kafka-settings.sasl_oauthbearer_expected_audience", "", "The (optional) comma-delimited setting for the broker to use to verify that the JWT was issued for one of the expected audiences.")
	var reqKafkaSettingsSaslOauthbearerExpectedIssuerFlag string
	flagset.StringVar(&reqKafkaSettingsSaslOauthbearerExpectedIssuerFlag, "kafka-settings.sasl_oauthbearer_expected_issuer", "", "Optional setting for the broker to use to verify that the JWT was created by the expected issuer.")
	var reqKafkaSettingsSaslOauthbearerJwksEndpointURLFlag string
	flagset.StringVar(&reqKafkaSettingsSaslOauthbearerJwksEndpointURLFlag, "kafka-settings.sasl_oauthbearer_jwks_endpoint_url", "", "OIDC JWKS endpoint URL. By setting this the SASL SSL OAuth2/OIDC authentication is enabled. See also other options for SASL OAuth2/OIDC. ")
	var reqKafkaSettingsSaslOauthbearerSubClaimNameFlag string
	flagset.StringVar(&reqKafkaSettingsSaslOauthbearerSubClaimNameFlag, "kafka-settings.sasl_oauthbearer_sub_claim_name", "", "Name of the scope from which to extract the subject claim from the JWT. Defaults to sub.")
	var reqKafkaSettingsSocketRequestMaxBytesFlag int
	flagset.IntVar(&reqKafkaSettingsSocketRequestMaxBytesFlag, "kafka-settings.socket_request_max_bytes", 0, "The maximum number of bytes in a socket request (defaults to 104857600).")
	var reqKafkaSettingsTransactionPartitionVerificationEnableFlag bool
	flagset.BoolVar(&reqKafkaSettingsTransactionPartitionVerificationEnableFlag, "kafka-settings.transaction_partition_verification_enable", false, "Enable verification that checks that the partition has been added to the transaction before writing transactional records to the partition")
	var reqKafkaSettingsTransactionRemoveExpiredTransactionCleanupIntervalMSFlag int
	flagset.IntVar(&reqKafkaSettingsTransactionRemoveExpiredTransactionCleanupIntervalMSFlag, "kafka-settings.transaction_remove_expired_transaction_cleanup_interval_ms", 0, "The interval at which to remove transactions that have expired due to transactional.id.expiration.ms passing (defaults to 3600000 (1 hour)).")
	var reqKafkaSettingsTransactionStateLogSegmentBytesFlag int
	flagset.IntVar(&reqKafkaSettingsTransactionStateLogSegmentBytesFlag, "kafka-settings.transaction_state_log_segment_bytes", 0, "The transaction topic segment bytes should be kept relatively small in order to facilitate faster log compaction and cache loads (defaults to 104857600 (100 mebibytes)).")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqSchemaRegistryEnabledFlag bool
	flagset.BoolVar(&reqSchemaRegistryEnabledFlag, "schema-registry-enabled", false, "Enable Schema-Registry service")
	var reqSchemaRegistrySettingsLeaderEligibilityFlag bool
	flagset.BoolVar(&reqSchemaRegistrySettingsLeaderEligibilityFlag, "schema-registry-settings.leader_eligibility", false, "If true, Karapace / Schema Registry on the service nodes can participate in leader election. It might be needed to disable this when the schemas topic is replicated to a secondary cluster and Karapace / Schema Registry there must not participate in leader election. Defaults to `true`.")
	var reqSchemaRegistrySettingsTopicNameFlag string
	flagset.StringVar(&reqSchemaRegistrySettingsTopicNameFlag, "schema-registry-settings.topic_name", "", "The durable single partition topic that acts as the durable log for the data. This topic must be compacted to avoid losing data due to retention policy. Please note that changing this configuration in an existing Schema Registry / Karapace setup leads to previous schemas being inaccessible, data encoded with them potentially unreadable and schema ID sequence put out of order. It's only possible to do the switch while Schema Registry / Karapace is disabled. Defaults to `_schemas`.")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")
	var reqVersionFlag string
	flagset.StringVar(&reqVersionFlag, "version", "", "Kafka major version")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASServiceKafkaRequest
	req.Version = reqVersionFlag
	req.TerminationProtection = reqTerminationProtectionFlag
	if reqSchemaRegistrySettingsTopicNameFlag != "" {
		req.SchemaRegistrySettings = &v3.JSONSchemaSchemaRegistry{}
		req.SchemaRegistrySettings.TopicName = reqSchemaRegistrySettingsTopicNameFlag
	}
	if reqSchemaRegistrySettingsLeaderEligibilityFlag != false {

		req.SchemaRegistrySettings.LeaderEligibility = reqSchemaRegistrySettingsLeaderEligibilityFlag
	}
	req.SchemaRegistryEnabled = reqSchemaRegistryEnabledFlag
	req.Plan = reqPlanFlag
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.UpdateDBAASServiceKafkaRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}
	if reqKafkaSettingsTransactionStateLogSegmentBytesFlag != 0 {
		req.Kafk = &v3.JSONSchemaKafka{}
		req.KafkaSettings.TransactionStateLogSegmentBytes = reqKafkaSettingsTransactionStateLogSegmentBytesFlag
	}
	if reqKafkaSettingsTransactionRemoveExpiredTransactionCleanupIntervalMSFlag != 0 {

		req.KafkaSettings.TransactionRemoveExpiredTransactionCleanupIntervalMS = reqKafkaSettingsTransactionRemoveExpiredTransactionCleanupIntervalMSFlag
	}
	if reqKafkaSettingsTransactionPartitionVerificationEnableFlag != false {

		req.KafkaSettings.TransactionPartitionVerificationEnable = reqKafkaSettingsTransactionPartitionVerificationEnableFlag
	}
	if reqKafkaSettingsSocketRequestMaxBytesFlag != 0 {

		req.KafkaSettings.SocketRequestMaxBytes = reqKafkaSettingsSocketRequestMaxBytesFlag
	}
	if reqKafkaSettingsSaslOauthbearerSubClaimNameFlag != "" {

		req.KafkaSettings.SaslOauthbearerSubClaimName = reqKafkaSettingsSaslOauthbearerSubClaimNameFlag
	}
	if reqKafkaSettingsSaslOauthbearerJwksEndpointURLFlag != "" {

		req.KafkaSettings.SaslOauthbearerJwksEndpointURL = reqKafkaSettingsSaslOauthbearerJwksEndpointURLFlag
	}
	if reqKafkaSettingsSaslOauthbearerExpectedIssuerFlag != "" {

		req.KafkaSettings.SaslOauthbearerExpectedIssuer = reqKafkaSettingsSaslOauthbearerExpectedIssuerFlag
	}
	if reqKafkaSettingsSaslOauthbearerExpectedAudienceFlag != "" {

		req.KafkaSettings.SaslOauthbearerExpectedAudience = reqKafkaSettingsSaslOauthbearerExpectedAudienceFlag
	}
	if reqKafkaSettingsReplicaFetchResponseMaxBytesFlag != 0 {

		req.KafkaSettings.ReplicaFetchResponseMaxBytes = reqKafkaSettingsReplicaFetchResponseMaxBytesFlag
	}
	if reqKafkaSettingsReplicaFetchMaxBytesFlag != 0 {

		req.KafkaSettings.ReplicaFetchMaxBytes = reqKafkaSettingsReplicaFetchMaxBytesFlag
	}
	if reqKafkaSettingsProducerPurgatoryPurgeIntervalRequestsFlag != 0 {

		req.KafkaSettings.ProducerPurgatoryPurgeIntervalRequests = reqKafkaSettingsProducerPurgatoryPurgeIntervalRequestsFlag
	}
	if reqKafkaSettingsOffsetsRetentionMinutesFlag != 0 {

		req.KafkaSettings.OffsetsRetentionMinutes = reqKafkaSettingsOffsetsRetentionMinutesFlag
	}
	if reqKafkaSettingsNumPartitionsFlag != 0 {

		req.KafkaSettings.NumPartitions = reqKafkaSettingsNumPartitionsFlag
	}
	if reqKafkaSettingsMinInsyncReplicasFlag != 0 {

		req.KafkaSettings.MinInsyncReplicas = reqKafkaSettingsMinInsyncReplicasFlag
	}
	if reqKafkaSettingsMessageMaxBytesFlag != 0 {

		req.KafkaSettings.MessageMaxBytes = reqKafkaSettingsMessageMaxBytesFlag
	}
	if reqKafkaSettingsMaxIncrementalFetchSessionCacheSlotsFlag != 0 {

		req.KafkaSettings.MaxIncrementalFetchSessionCacheSlots = reqKafkaSettingsMaxIncrementalFetchSessionCacheSlotsFlag
	}
	if reqKafkaSettingsMaxConnectionsPerIPFlag != 0 {

		req.KafkaSettings.MaxConnectionsPerIP = reqKafkaSettingsMaxConnectionsPerIPFlag
	}
	if reqKafkaSettingsLogSegmentDeleteDelayMSFlag != 0 {

		req.KafkaSettings.LogSegmentDeleteDelayMS = reqKafkaSettingsLogSegmentDeleteDelayMSFlag
	}
	if reqKafkaSettingsLogSegmentBytesFlag != 0 {

		req.KafkaSettings.LogSegmentBytes = reqKafkaSettingsLogSegmentBytesFlag
	}
	if reqKafkaSettingsLogRollMSFlag != 0 {

		req.KafkaSettings.LogRollMS = reqKafkaSettingsLogRollMSFlag
	}
	if reqKafkaSettingsLogRollJitterMSFlag != 0 {

		req.KafkaSettings.LogRollJitterMS = reqKafkaSettingsLogRollJitterMSFlag
	}
	if reqKafkaSettingsLogRetentionMSFlag != 0 {

		req.KafkaSettings.LogRetentionMS = reqKafkaSettingsLogRetentionMSFlag
	}
	if reqKafkaSettingsLogRetentionHoursFlag != 0 {

		req.KafkaSettings.LogRetentionHours = reqKafkaSettingsLogRetentionHoursFlag
	}
	if reqKafkaSettingsLogRetentionBytesFlag != 0 {

		req.KafkaSettings.LogRetentionBytes = reqKafkaSettingsLogRetentionBytesFlag
	}
	if reqKafkaSettingsLogPreallocateFlag != false {

		req.KafkaSettings.LogPreallocate = reqKafkaSettingsLogPreallocateFlag
	}
	if reqKafkaSettingsLogMessageTimestampTypeFlag != "" {

		req.KafkaSettings.LogMessageTimestampType = reqKafkaSettingsLogMessageTimestampTypeFlag
	}
	if reqKafkaSettingsLogMessageTimestampDifferenceMaxMSFlag != 0 {

		req.KafkaSettings.LogMessageTimestampDifferenceMaxMS = reqKafkaSettingsLogMessageTimestampDifferenceMaxMSFlag
	}
	if reqKafkaSettingsLogMessageDownconversionEnableFlag != false {

		req.KafkaSettings.LogMessageDownconversionEnable = reqKafkaSettingsLogMessageDownconversionEnableFlag
	}
	if reqKafkaSettingsLogLocalRetentionMSFlag != 0 {

		req.KafkaSettings.LogLocalRetentionMS = reqKafkaSettingsLogLocalRetentionMSFlag
	}
	if reqKafkaSettingsLogLocalRetentionBytesFlag != 0 {

		req.KafkaSettings.LogLocalRetentionBytes = reqKafkaSettingsLogLocalRetentionBytesFlag
	}
	if reqKafkaSettingsLogIndexSizeMaxBytesFlag != 0 {

		req.KafkaSettings.LogIndexSizeMaxBytes = reqKafkaSettingsLogIndexSizeMaxBytesFlag
	}
	if reqKafkaSettingsLogIndexIntervalBytesFlag != 0 {

		req.KafkaSettings.LogIndexIntervalBytes = reqKafkaSettingsLogIndexIntervalBytesFlag
	}
	if reqKafkaSettingsLogFlushIntervalMSFlag != 0 {

		req.KafkaSettings.LogFlushIntervalMS = reqKafkaSettingsLogFlushIntervalMSFlag
	}
	if reqKafkaSettingsLogFlushIntervalMessagesFlag != 0 {

		req.KafkaSettings.LogFlushIntervalMessages = reqKafkaSettingsLogFlushIntervalMessagesFlag
	}
	if reqKafkaSettingsLogCleanupAndCompactionLogCleanupPolicyFlag != "" {
		req.KafkaSettingsLogCleanupAndCompact = &v3.KafkaSettingsLogCleanupAndCompaction{}
		req.KafkaSettingsLogCleanupAndCompaction.LogCleanupPolicy = reqKafkaSettingsLogCleanupAndCompactionLogCleanupPolicyFlag
	}
	if reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCompactionLagMSFlag != 0 {

		req.KafkaSettingsLogCleanupAndCompaction.LogCleanerMinCompactionLagMS = reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCompactionLagMSFlag
	}
	if reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCleanableRatioFlag != 0 {

		req.KafkaSettingsLogCleanupAndCompaction.LogCleanerMinCleanableRatio = reqKafkaSettingsLogCleanupAndCompactionLogCleanerMinCleanableRatioFlag
	}
	if reqKafkaSettingsLogCleanupAndCompactionLogCleanerMaxCompactionLagMSFlag != 0 {

		req.KafkaSettingsLogCleanupAndCompaction.LogCleanerMaxCompactionLagMS = reqKafkaSettingsLogCleanupAndCompactionLogCleanerMaxCompactionLagMSFlag
	}
	if reqKafkaSettingsLogCleanupAndCompactionLogCleanerDeleteRetentionMSFlag != 0 {

		req.KafkaSettingsLogCleanupAndCompaction.LogCleanerDeleteRetentionMS = reqKafkaSettingsLogCleanupAndCompactionLogCleanerDeleteRetentionMSFlag
	}
	if reqKafkaSettingsGroupMinSessionTimeoutMSFlag != 0 {
		req.KafkaSetting = &v3.JSONSchemaKafka{}
		req.KafkaSettings.GroupMinSessionTimeoutMS = reqKafkaSettingsGroupMinSessionTimeoutMSFlag
	}
	if reqKafkaSettingsGroupMaxSessionTimeoutMSFlag != 0 {

		req.KafkaSettings.GroupMaxSessionTimeoutMS = reqKafkaSettingsGroupMaxSessionTimeoutMSFlag
	}
	if reqKafkaSettingsGroupInitialRebalanceDelayMSFlag != 0 {

		req.KafkaSettings.GroupInitialRebalanceDelayMS = reqKafkaSettingsGroupInitialRebalanceDelayMSFlag
	}
	if reqKafkaSettingsDefaultReplicationFactorFlag != 0 {

		req.KafkaSettings.DefaultReplicationFactor = reqKafkaSettingsDefaultReplicationFactorFlag
	}
	if reqKafkaSettingsConnectionsMaxIdleMSFlag != 0 {

		req.KafkaSettings.ConnectionsMaxIdleMS = reqKafkaSettingsConnectionsMaxIdleMSFlag
	}
	if reqKafkaSettingsCompressionTypeFlag != "" {

		req.KafkaSettings.CompressionType = reqKafkaSettingsCompressionTypeFlag
	}
	if reqKafkaSettingsAutoCreateTopicsEnableFlag != false {

		req.KafkaSettings.AutoCreateTopicsEnable = reqKafkaSettingsAutoCreateTopicsEnableFlag
	}
	if reqKafkaRestSettingsSimpleconsumerPoolSizeMaxFlag != 0 {
		req.KafkaRestSetting = &v3.JSONSchemaKafkaRest{}
		req.KafkaRestSettings.SimpleconsumerPoolSizeMax = reqKafkaRestSettingsSimpleconsumerPoolSizeMaxFlag
	}
	if reqKafkaRestSettingsProducerMaxRequestSizeFlag != 0 {

		req.KafkaRestSettings.ProducerMaxRequestSize = reqKafkaRestSettingsProducerMaxRequestSizeFlag
	}
	if reqKafkaRestSettingsProducerLingerMSFlag != 0 {

		req.KafkaRestSettings.ProducerLingerMS = reqKafkaRestSettingsProducerLingerMSFlag
	}
	if reqKafkaRestSettingsProducerCompressionTypeFlag != "" {

		req.KafkaRestSettings.ProducerCompressionType = reqKafkaRestSettingsProducerCompressionTypeFlag
	}
	if reqKafkaRestSettingsProducerAcksFlag != "" {

		req.KafkaRestSettings.ProducerAcks = reqKafkaRestSettingsProducerAcksFlag
	}
	if reqKafkaRestSettingsNameStrategyValidationFlag != false {

		req.KafkaRestSettings.NameStrategyValidation = reqKafkaRestSettingsNameStrategyValidationFlag
	}
	if reqKafkaRestSettingsNameStrategyFlag != "" {

		req.KafkaRestSettings.NameStrategy = reqKafkaRestSettingsNameStrategyFlag
	}
	if reqKafkaRestSettingsConsumerRequestTimeoutMSFlag != 0 {

		req.KafkaRestSettings.ConsumerRequestTimeoutMS = reqKafkaRestSettingsConsumerRequestTimeoutMSFlag
	}
	if reqKafkaRestSettingsConsumerRequestMaxBytesFlag != 0 {

		req.KafkaRestSettings.ConsumerRequestMaxBytes = reqKafkaRestSettingsConsumerRequestMaxBytesFlag
	}
	if reqKafkaRestSettingsConsumerEnableAutoCommitFlag != false {

		req.KafkaRestSettings.ConsumerEnableAutoCommit = reqKafkaRestSettingsConsumerEnableAutoCommitFlag
	}
	req.KafkaRestEnabled = reqKafkaRestEnabledFlag
	if reqKafkaConnectSettingsSessionTimeoutMSFlag != 0 {
		req.KafkaConnectSetting = &v3.JSONSchemaKafkaConnect{}
		req.KafkaConnectSettings.SessionTimeoutMS = reqKafkaConnectSettingsSessionTimeoutMSFlag
	}
	if reqKafkaConnectSettingsScheduledRebalanceMaxDelayMSFlag != 0 {

		req.KafkaConnectSettings.ScheduledRebalanceMaxDelayMS = reqKafkaConnectSettingsScheduledRebalanceMaxDelayMSFlag
	}
	if reqKafkaConnectSettingsProducerMaxRequestSizeFlag != 0 {

		req.KafkaConnectSettings.ProducerMaxRequestSize = reqKafkaConnectSettingsProducerMaxRequestSizeFlag
	}
	if reqKafkaConnectSettingsProducerLingerMSFlag != 0 {

		req.KafkaConnectSettings.ProducerLingerMS = reqKafkaConnectSettingsProducerLingerMSFlag
	}
	if reqKafkaConnectSettingsProducerCompressionTypeFlag != "" {

		req.KafkaConnectSettings.ProducerCompressionType = reqKafkaConnectSettingsProducerCompressionTypeFlag
	}
	if reqKafkaConnectSettingsProducerBufferMemoryFlag != 0 {

		req.KafkaConnectSettings.ProducerBufferMemory = reqKafkaConnectSettingsProducerBufferMemoryFlag
	}
	if reqKafkaConnectSettingsProducerBatchSizeFlag != 0 {

		req.KafkaConnectSettings.ProducerBatchSize = reqKafkaConnectSettingsProducerBatchSizeFlag
	}
	if reqKafkaConnectSettingsOffsetFlushTimeoutMSFlag != 0 {

		req.KafkaConnectSettings.OffsetFlushTimeoutMS = reqKafkaConnectSettingsOffsetFlushTimeoutMSFlag
	}
	if reqKafkaConnectSettingsOffsetFlushIntervalMSFlag != 0 {

		req.KafkaConnectSettings.OffsetFlushIntervalMS = reqKafkaConnectSettingsOffsetFlushIntervalMSFlag
	}
	if reqKafkaConnectSettingsConsumerMaxPollRecordsFlag != 0 {

		req.KafkaConnectSettings.ConsumerMaxPollRecords = reqKafkaConnectSettingsConsumerMaxPollRecordsFlag
	}
	if reqKafkaConnectSettingsConsumerMaxPollIntervalMSFlag != 0 {

		req.KafkaConnectSettings.ConsumerMaxPollIntervalMS = reqKafkaConnectSettingsConsumerMaxPollIntervalMSFlag
	}
	if reqKafkaConnectSettingsConsumerMaxPartitionFetchBytesFlag != 0 {

		req.KafkaConnectSettings.ConsumerMaxPartitionFetchBytes = reqKafkaConnectSettingsConsumerMaxPartitionFetchBytesFlag
	}
	if reqKafkaConnectSettingsConsumerIsolationLevelFlag != "" {

		req.KafkaConnectSettings.ConsumerIsolationLevel = reqKafkaConnectSettingsConsumerIsolationLevelFlag
	}
	if reqKafkaConnectSettingsConsumerFetchMaxBytesFlag != 0 {

		req.KafkaConnectSettings.ConsumerFetchMaxBytes = reqKafkaConnectSettingsConsumerFetchMaxBytesFlag
	}
	if reqKafkaConnectSettingsConsumerAutoOffsetResetFlag != "" {

		req.KafkaConnectSettings.ConsumerAutoOffsetReset = reqKafkaConnectSettingsConsumerAutoOffsetResetFlag
	}
	if reqKafkaConnectSettingsConnectorClientConfigOverridePolicyFlag != "" {

		req.KafkaConnectSettings.ConnectorClientConfigOverridePolicy = reqKafkaConnectSettingsConnectorClientConfigOverridePolicyFlag
	}
	req.KafkaConnectEnabled = reqKafkaConnectEnabledFlag
	if reqAuthenticationMethodsSaslFlag != false {
		req.AuthenticationMethod = &v3.UpdateDBAASServiceKafkaRequestAuthenticationMethods{}
		req.AuthenticationMethods.Sasl = reqAuthenticationMethodsSaslFlag
	}
	if reqAuthenticationMethodsCertificateFlag != false {

		req.AuthenticationMethods.Certificate = reqAuthenticationMethodsCertificateFlag
	}

	resp, err := client.UpdateDBAASServiceKafka(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASKafkaAclConfigCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-kafka-acl-config", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASKafkaAclConfig(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StartDBAASKafkaMaintenanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("start-dbaas-kafka-maintenance", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StartDBAASKafkaMaintenance(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASKafkaSchemaRegistryAclConfigCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-kafka-schema-registry-acl-config", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqIDFlag string
	flagset.StringVar(&reqIDFlag, "id", "", "")
	var reqPermissionFlag string
	flagset.StringVar(&reqPermissionFlag, "permission", "", "Kafka Schema Registry permission")
	var reqResourceFlag string
	flagset.StringVar(&reqResourceFlag, "resource", "", "Kafka Schema Registry name or pattern")
	var reqUsernameFlag string
	flagset.StringVar(&reqUsernameFlag, "username", "", "Kafka username or username pattern")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASKafkaSchemaRegistryAclConfigRequest
	req.Username = reqUsernameFlag
	req.Resource = reqResourceFlag
	req.Permission = reqPermissionFlag
	req.ID = reqIDFlag

	resp, err := client.CreateDBAASKafkaSchemaRegistryAclConfig(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASKafkaSchemaRegistryAclConfigCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-kafka-schema-registry-acl-config", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var aclIDFlag string
	flagset.StringVar(&aclIDFlag, "AclID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASKafkaSchemaRegistryAclConfig(context.Background(), nameFlag, aclIDFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASKafkaTopicAclConfigCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-kafka-topic-acl-config", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqIDFlag string
	flagset.StringVar(&reqIDFlag, "id", "", "")
	var reqPermissionFlag string
	flagset.StringVar(&reqPermissionFlag, "permission", "", "Kafka permission")
	var reqTopicFlag string
	flagset.StringVar(&reqTopicFlag, "topic", "", "Kafka topic name or pattern")
	var reqUsernameFlag string
	flagset.StringVar(&reqUsernameFlag, "username", "", "Kafka username or username pattern")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASKafkaTopicAclConfigRequest
	req.Username = reqUsernameFlag
	req.Topic = reqTopicFlag
	req.Permission = reqPermissionFlag
	req.ID = reqIDFlag

	resp, err := client.CreateDBAASKafkaTopicAclConfig(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASKafkaTopicAclConfigCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-kafka-topic-acl-config", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var aclIDFlag string
	flagset.StringVar(&aclIDFlag, "AclID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASKafkaTopicAclConfig(context.Background(), nameFlag, aclIDFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RevealDBAASKafkaConnectPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reveal-dbaas-kafka-connect-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RevealDBAASKafkaConnectPassword(context.Background(), serviceNameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASKafkaUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-kafka-user", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var reqUsernameFlag string
	flagset.StringVar(&reqUsernameFlag, "username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASKafkaUserRequest
	req.Username = reqUsernameFlag

	resp, err := client.CreateDBAASKafkaUser(context.Background(), serviceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASKafkaUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-kafka-user", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASKafkaUser(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetDBAASKafkaUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-dbaas-kafka-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")
	var reqPasswordFlag string
	flagset.StringVar(&reqPasswordFlag, "password", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ResetDBAASKafkaUserPasswordRequest
	req.Password = reqPasswordFlag

	resp, err := client.ResetDBAASKafkaUserPassword(context.Background(), serviceNameFlag, usernameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RevealDBAASKafkaUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reveal-dbaas-kafka-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RevealDBAASKafkaUserPassword(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASMigrationStatusCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-migration-status", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASMigrationStatus(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASServiceMysqlCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-service-mysql", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASServiceMysql(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASServiceMysqlCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-service-mysql", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASServiceMysql(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASServiceMysqlCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-service-mysql", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqAdminPasswordFlag string
	flagset.StringVar(&reqAdminPasswordFlag, "admin-password", "", "Custom password for admin user. Defaults to random string. This must be set only when a new service is being created.")
	var reqAdminUsernameFlag string
	flagset.StringVar(&reqAdminUsernameFlag, "admin-username", "", "Custom username for admin user. This must be set only when a new service is being created.")
	var reqBackupScheduleBackupHourFlag int64
	flagset.Int64Var(&reqBackupScheduleBackupHourFlag, "backup-schedule.backup-hour", 0, "The hour of day (in UTC) when backup for the service is started. New backup is only started if previous backup has already completed.")
	var reqBackupScheduleBackupMinuteFlag int64
	flagset.Int64Var(&reqBackupScheduleBackupMinuteFlag, "backup-schedule.backup-minute", 0, "The minute of an hour when backup for the service is started. New backup is only started if previous backup has already completed.")
	var reqBinlogRetentionPeriodFlag int64
	flagset.Int64Var(&reqBinlogRetentionPeriodFlag, "binlog-retention-period", 0, "The minimum amount of time in seconds to keep binlog entries before deletion. This may be extended for services that require binlog entries for longer than the default for example if using the MySQL Debezium Kafka connector.")
	var reqForkFromServiceFlag string
	flagset.StringVar(&reqForkFromServiceFlag, "fork-from-service", "", "")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqMigrationDbnameFlag string
	flagset.StringVar(&reqMigrationDbnameFlag, "migration.dbname", "", "Database name for bootstrapping the initial connection")
	var reqMigrationHostFlag string
	flagset.StringVar(&reqMigrationHostFlag, "migration.host", "", "Hostname or IP address of the server where to migrate data from")
	var reqMigrationIgnoreDbsFlag string
	flagset.StringVar(&reqMigrationIgnoreDbsFlag, "migration.ignore-dbs", "", "Comma-separated list of databases, which should be ignored during migration (supported by MySQL only at the moment)")
	var reqMigrationMethodFlag string
	flagset.StringVar(&reqMigrationMethodFlag, "migration.method", "", "")
	var reqMigrationPasswordFlag string
	flagset.StringVar(&reqMigrationPasswordFlag, "migration.password", "", "Password for authentication with the server where to migrate data from")
	var reqMigrationPortFlag int64
	flagset.Int64Var(&reqMigrationPortFlag, "migration.port", 0, "Port number of the server where to migrate data from")
	var reqMigrationSSLFlag bool
	flagset.BoolVar(&reqMigrationSSLFlag, "migration.ssl", false, "The server where to migrate data from is secured with SSL")
	var reqMigrationUsernameFlag string
	flagset.StringVar(&reqMigrationUsernameFlag, "migration.username", "", "User name for authentication with the server where to migrate data from")
	var reqMysqlSettingsConnectTimeoutFlag int
	flagset.IntVar(&reqMysqlSettingsConnectTimeoutFlag, "mysql-settings.connect_timeout", 0, "The number of seconds that the mysqld server waits for a connect packet before responding with Bad handshake")
	var reqMysqlSettingsDefaultTimeZoneFlag string
	flagset.StringVar(&reqMysqlSettingsDefaultTimeZoneFlag, "mysql-settings.default_time_zone", "", "Default server time zone as an offset from UTC (from -12:00 to +12:00), a time zone name, or 'SYSTEM' to use the MySQL server default.")
	var reqMysqlSettingsGroupConcatMaxLenFlag int
	flagset.IntVar(&reqMysqlSettingsGroupConcatMaxLenFlag, "mysql-settings.group_concat_max_len", 0, "The maximum permitted result length in bytes for the GROUP_CONCAT() function.")
	var reqMysqlSettingsInformationSchemaStatsExpiryFlag int
	flagset.IntVar(&reqMysqlSettingsInformationSchemaStatsExpiryFlag, "mysql-settings.information_schema_stats_expiry", 0, "The time, in seconds, before cached statistics expire")
	var reqMysqlSettingsInnodbChangeBufferMaxSizeFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbChangeBufferMaxSizeFlag, "mysql-settings.innodb_change_buffer_max_size", 0, "Maximum size for the InnoDB change buffer, as a percentage of the total size of the buffer pool. Default is 25")
	var reqMysqlSettingsInnodbFlushNeighborsFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbFlushNeighborsFlag, "mysql-settings.innodb_flush_neighbors", 0, "Specifies whether flushing a page from the InnoDB buffer pool also flushes other dirty pages in the same extent (default is 1): 0 - dirty pages in the same extent are not flushed, 1 - flush contiguous dirty pages in the same extent, 2 - flush dirty pages in the same extent")
	var reqMysqlSettingsInnodbFTMinTokenSizeFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbFTMinTokenSizeFlag, "mysql-settings.innodb_ft_min_token_size", 0, "Minimum length of words that are stored in an InnoDB FULLTEXT index. Changing this parameter will lead to a restart of the MySQL service.")
	var reqMysqlSettingsInnodbFTServerStopwordTableFlag string
	flagset.StringVar(&reqMysqlSettingsInnodbFTServerStopwordTableFlag, "mysql-settings.innodb_ft_server_stopword_table", "", "This option is used to specify your own InnoDB FULLTEXT index stopword list for all InnoDB tables.")
	var reqMysqlSettingsInnodbLockWaitTimeoutFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbLockWaitTimeoutFlag, "mysql-settings.innodb_lock_wait_timeout", 0, "The length of time in seconds an InnoDB transaction waits for a row lock before giving up. Default is 120.")
	var reqMysqlSettingsInnodbLogBufferSizeFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbLogBufferSizeFlag, "mysql-settings.innodb_log_buffer_size", 0, "The size in bytes of the buffer that InnoDB uses to write to the log files on disk.")
	var reqMysqlSettingsInnodbOnlineAlterLogMaxSizeFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbOnlineAlterLogMaxSizeFlag, "mysql-settings.innodb_online_alter_log_max_size", 0, "The upper limit in bytes on the size of the temporary log files used during online DDL operations for InnoDB tables.")
	var reqMysqlSettingsInnodbPrintAllDeadlocksFlag bool
	flagset.BoolVar(&reqMysqlSettingsInnodbPrintAllDeadlocksFlag, "mysql-settings.innodb_print_all_deadlocks", false, "When enabled, information about all deadlocks in InnoDB user transactions is recorded in the error log. Disabled by default.")
	var reqMysqlSettingsInnodbReadIoThreadsFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbReadIoThreadsFlag, "mysql-settings.innodb_read_io_threads", 0, "The number of I/O threads for read operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.")
	var reqMysqlSettingsInnodbRollbackOnTimeoutFlag bool
	flagset.BoolVar(&reqMysqlSettingsInnodbRollbackOnTimeoutFlag, "mysql-settings.innodb_rollback_on_timeout", false, "When enabled a transaction timeout causes InnoDB to abort and roll back the entire transaction. Changing this parameter will lead to a restart of the MySQL service.")
	var reqMysqlSettingsInnodbThreadConcurrencyFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbThreadConcurrencyFlag, "mysql-settings.innodb_thread_concurrency", 0, "Defines the maximum number of threads permitted inside of InnoDB. Default is 0 (infinite concurrency - no limit)")
	var reqMysqlSettingsInnodbWriteIoThreadsFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbWriteIoThreadsFlag, "mysql-settings.innodb_write_io_threads", 0, "The number of I/O threads for write operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.")
	var reqMysqlSettingsInteractiveTimeoutFlag int
	flagset.IntVar(&reqMysqlSettingsInteractiveTimeoutFlag, "mysql-settings.interactive_timeout", 0, "The number of seconds the server waits for activity on an interactive connection before closing it.")
	var reqMysqlSettingsInternalTmpMemStorageEngineFlag string
	flagset.StringVar(&reqMysqlSettingsInternalTmpMemStorageEngineFlag, "mysql-settings.internal_tmp_mem_storage_engine", "", "The storage engine for in-memory internal temporary tables.")
	var reqMysqlSettingsLogOutputFlag string
	flagset.StringVar(&reqMysqlSettingsLogOutputFlag, "mysql-settings.log_output", "", "The slow log output destination when slow_query_log is ON. To enable MySQL AI Insights, choose INSIGHTS. To use MySQL AI Insights and the mysql.slow_log table at the same time, choose INSIGHTS,TABLE. To only use the mysql.slow_log table, choose TABLE. To silence slow logs, choose NONE.")
	var reqMysqlSettingsLongQueryTimeFlag float64
	flagset.Float64Var(&reqMysqlSettingsLongQueryTimeFlag, "mysql-settings.long_query_time", 0, "The slow_query_logs work as SQL statements that take more than long_query_time seconds to execute. Default is 10s")
	var reqMysqlSettingsMaxAllowedPacketFlag int
	flagset.IntVar(&reqMysqlSettingsMaxAllowedPacketFlag, "mysql-settings.max_allowed_packet", 0, "Size of the largest message in bytes that can be received by the server. Default is 67108864 (64M)")
	var reqMysqlSettingsMaxHeapTableSizeFlag int
	flagset.IntVar(&reqMysqlSettingsMaxHeapTableSizeFlag, "mysql-settings.max_heap_table_size", 0, "Limits the size of internal in-memory tables. Also set tmp_table_size. Default is 16777216 (16M)")
	var reqMysqlSettingsNetBufferLengthFlag int
	flagset.IntVar(&reqMysqlSettingsNetBufferLengthFlag, "mysql-settings.net_buffer_length", 0, "Start sizes of connection buffer and result buffer. Default is 16384 (16K). Changing this parameter will lead to a restart of the MySQL service.")
	var reqMysqlSettingsNetReadTimeoutFlag int
	flagset.IntVar(&reqMysqlSettingsNetReadTimeoutFlag, "mysql-settings.net_read_timeout", 0, "The number of seconds to wait for more data from a connection before aborting the read.")
	var reqMysqlSettingsNetWriteTimeoutFlag int
	flagset.IntVar(&reqMysqlSettingsNetWriteTimeoutFlag, "mysql-settings.net_write_timeout", 0, "The number of seconds to wait for a block to be written to a connection before aborting the write.")
	var reqMysqlSettingsSlowQueryLogFlag bool
	flagset.BoolVar(&reqMysqlSettingsSlowQueryLogFlag, "mysql-settings.slow_query_log", false, "Slow query log enables capturing of slow queries. Setting slow_query_log to false also truncates the mysql.slow_log table. Default is off")
	var reqMysqlSettingsSortBufferSizeFlag int
	flagset.IntVar(&reqMysqlSettingsSortBufferSizeFlag, "mysql-settings.sort_buffer_size", 0, "Sort buffer size in bytes for ORDER BY optimization. Default is 262144 (256K)")
	var reqMysqlSettingsSQLModeFlag string
	flagset.StringVar(&reqMysqlSettingsSQLModeFlag, "mysql-settings.sql_mode", "", "Global SQL mode. Set to empty to use MySQL server defaults. When creating a new service and not setting this field Aiven default SQL mode (strict, SQL standard compliant) will be assigned.")
	var reqMysqlSettingsSQLRequirePrimaryKeyFlag bool
	flagset.BoolVar(&reqMysqlSettingsSQLRequirePrimaryKeyFlag, "mysql-settings.sql_require_primary_key", false, "Require primary key to be defined for new tables or old tables modified with ALTER TABLE and fail if missing. It is recommended to always have primary keys because various functionality may break if any large table is missing them.")
	var reqMysqlSettingsTmpTableSizeFlag int
	flagset.IntVar(&reqMysqlSettingsTmpTableSizeFlag, "mysql-settings.tmp_table_size", 0, "Limits the size of internal in-memory tables. Also set max_heap_table_size. Default is 16777216 (16M)")
	var reqMysqlSettingsWaitTimeoutFlag int
	flagset.IntVar(&reqMysqlSettingsWaitTimeoutFlag, "mysql-settings.wait_timeout", 0, "The number of seconds the server waits for activity on a noninteractive connection before closing it.")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqRecoveryBackupTimeFlag string
	flagset.StringVar(&reqRecoveryBackupTimeFlag, "recovery-backup-time", "", "ISO time of a backup to recover from for services that support arbitrary times")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")
	var reqVersionFlag string
	flagset.StringVar(&reqVersionFlag, "version", "", "MySQL major version")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASServiceMysqlRequest
	req.Version = reqVersionFlag
	req.TerminationProtection = reqTerminationProtectionFlag
	req.RecoveryBackupTime = reqRecoveryBackupTimeFlag
	req.Plan = reqPlanFlag
	if reqMysqlSettingsWaitTimeoutFlag != 0 {
		req.MysqlSettings = &v3.JSONSchemaMysql{}
		req.MysqlSettings.WaitTimeout = reqMysqlSettingsWaitTimeoutFlag
	}
	if reqMysqlSettingsTmpTableSizeFlag != 0 {

		req.MysqlSettings.TmpTableSize = reqMysqlSettingsTmpTableSizeFlag
	}
	if reqMysqlSettingsSQLRequirePrimaryKeyFlag != false {

		req.MysqlSettings.SQLRequirePrimaryKey = reqMysqlSettingsSQLRequirePrimaryKeyFlag
	}
	if reqMysqlSettingsSQLModeFlag != "" {

		req.MysqlSettings.SQLMode = reqMysqlSettingsSQLModeFlag
	}
	if reqMysqlSettingsSortBufferSizeFlag != 0 {

		req.MysqlSettings.SortBufferSize = reqMysqlSettingsSortBufferSizeFlag
	}
	if reqMysqlSettingsSlowQueryLogFlag != false {

		req.MysqlSettings.SlowQueryLog = reqMysqlSettingsSlowQueryLogFlag
	}
	if reqMysqlSettingsNetWriteTimeoutFlag != 0 {

		req.MysqlSettings.NetWriteTimeout = reqMysqlSettingsNetWriteTimeoutFlag
	}
	if reqMysqlSettingsNetReadTimeoutFlag != 0 {

		req.MysqlSettings.NetReadTimeout = reqMysqlSettingsNetReadTimeoutFlag
	}
	if reqMysqlSettingsNetBufferLengthFlag != 0 {

		req.MysqlSettings.NetBufferLength = reqMysqlSettingsNetBufferLengthFlag
	}
	if reqMysqlSettingsMaxHeapTableSizeFlag != 0 {

		req.MysqlSettings.MaxHeapTableSize = reqMysqlSettingsMaxHeapTableSizeFlag
	}
	if reqMysqlSettingsMaxAllowedPacketFlag != 0 {

		req.MysqlSettings.MaxAllowedPacket = reqMysqlSettingsMaxAllowedPacketFlag
	}
	if reqMysqlSettingsLongQueryTimeFlag != 0 {

		req.MysqlSettings.LongQueryTime = reqMysqlSettingsLongQueryTimeFlag
	}
	if reqMysqlSettingsLogOutputFlag != "" {

		req.MysqlSettings.LogOutput = reqMysqlSettingsLogOutputFlag
	}
	if reqMysqlSettingsInternalTmpMemStorageEngineFlag != "" {

		req.MysqlSettings.InternalTmpMemStorageEngine = reqMysqlSettingsInternalTmpMemStorageEngineFlag
	}
	if reqMysqlSettingsInteractiveTimeoutFlag != 0 {

		req.MysqlSettings.InteractiveTimeout = reqMysqlSettingsInteractiveTimeoutFlag
	}
	if reqMysqlSettingsInnodbWriteIoThreadsFlag != 0 {

		req.MysqlSettings.InnodbWriteIoThreads = reqMysqlSettingsInnodbWriteIoThreadsFlag
	}
	if reqMysqlSettingsInnodbThreadConcurrencyFlag != 0 {

		req.MysqlSettings.InnodbThreadConcurrency = reqMysqlSettingsInnodbThreadConcurrencyFlag
	}
	if reqMysqlSettingsInnodbRollbackOnTimeoutFlag != false {

		req.MysqlSettings.InnodbRollbackOnTimeout = reqMysqlSettingsInnodbRollbackOnTimeoutFlag
	}
	if reqMysqlSettingsInnodbReadIoThreadsFlag != 0 {

		req.MysqlSettings.InnodbReadIoThreads = reqMysqlSettingsInnodbReadIoThreadsFlag
	}
	if reqMysqlSettingsInnodbPrintAllDeadlocksFlag != false {

		req.MysqlSettings.InnodbPrintAllDeadlocks = reqMysqlSettingsInnodbPrintAllDeadlocksFlag
	}
	if reqMysqlSettingsInnodbOnlineAlterLogMaxSizeFlag != 0 {

		req.MysqlSettings.InnodbOnlineAlterLogMaxSize = reqMysqlSettingsInnodbOnlineAlterLogMaxSizeFlag
	}
	if reqMysqlSettingsInnodbLogBufferSizeFlag != 0 {

		req.MysqlSettings.InnodbLogBufferSize = reqMysqlSettingsInnodbLogBufferSizeFlag
	}
	if reqMysqlSettingsInnodbLockWaitTimeoutFlag != 0 {

		req.MysqlSettings.InnodbLockWaitTimeout = reqMysqlSettingsInnodbLockWaitTimeoutFlag
	}
	if reqMysqlSettingsInnodbFTServerStopwordTableFlag != "" {

		req.MysqlSettings.InnodbFTServerStopwordTable = reqMysqlSettingsInnodbFTServerStopwordTableFlag
	}
	if reqMysqlSettingsInnodbFTMinTokenSizeFlag != 0 {

		req.MysqlSettings.InnodbFTMinTokenSize = reqMysqlSettingsInnodbFTMinTokenSizeFlag
	}
	if reqMysqlSettingsInnodbFlushNeighborsFlag != 0 {

		req.MysqlSettings.InnodbFlushNeighbors = reqMysqlSettingsInnodbFlushNeighborsFlag
	}
	if reqMysqlSettingsInnodbChangeBufferMaxSizeFlag != 0 {

		req.MysqlSettings.InnodbChangeBufferMaxSize = reqMysqlSettingsInnodbChangeBufferMaxSizeFlag
	}
	if reqMysqlSettingsInformationSchemaStatsExpiryFlag != 0 {

		req.MysqlSettings.InformationSchemaStatsExpiry = reqMysqlSettingsInformationSchemaStatsExpiryFlag
	}
	if reqMysqlSettingsGroupConcatMaxLenFlag != 0 {

		req.MysqlSettings.GroupConcatMaxLen = reqMysqlSettingsGroupConcatMaxLenFlag
	}
	if reqMysqlSettingsDefaultTimeZoneFlag != "" {

		req.MysqlSettings.DefaultTimeZone = reqMysqlSettingsDefaultTimeZoneFlag
	}
	if reqMysqlSettingsConnectTimeoutFlag != 0 {

		req.MysqlSettings.ConnectTimeout = reqMysqlSettingsConnectTimeoutFlag
	}
	if reqMigrationUsernameFlag != "" {
		req.Migratio = &v3.CreateDBAASServiceMysqlRequestMigration{}
		req.Migration.Username = reqMigrationUsernameFlag
	}
	if reqMigrationSSLFlag != false {

		req.Migration.SSL = reqMigrationSSLFlag
	}
	if reqMigrationPortFlag != 0 {

		req.Migration.Port = reqMigrationPortFlag
	}
	if reqMigrationPasswordFlag != "" {

		req.Migration.Password = reqMigrationPasswordFlag
	}
	if reqMigrationMethodFlag != "" {

		req.Migration.Method = reqMigrationMethodFlag
	}
	if reqMigrationIgnoreDbsFlag != "" {

		req.Migration.IgnoreDbs = reqMigrationIgnoreDbsFlag
	}
	if reqMigrationHostFlag != "" {

		req.Migration.Host = reqMigrationHostFlag
	}
	if reqMigrationDbnameFlag != "" {

		req.Migration.Dbname = reqMigrationDbnameFlag
	}
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.CreateDBAASServiceMysqlRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}
	req.ForkFromService = reqForkFromServiceFlag
	req.BinlogRetentionPeriod = reqBinlogRetentionPeriodFlag
	if reqBackupScheduleBackupMinuteFlag != 0 {
		req.BackupSchedul = &v3.CreateDBAASServiceMysqlRequestBackupSchedule{}
		req.BackupSchedule.BackupMinute = reqBackupScheduleBackupMinuteFlag
	}
	if reqBackupScheduleBackupHourFlag != 0 {

		req.BackupSchedule.BackupHour = reqBackupScheduleBackupHourFlag
	}
	req.AdminUsername = reqAdminUsernameFlag
	req.AdminPassword = reqAdminPasswordFlag

	resp, err := client.CreateDBAASServiceMysql(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASServiceMysqlCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-service-mysql", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqBackupScheduleBackupHourFlag int64
	flagset.Int64Var(&reqBackupScheduleBackupHourFlag, "backup-schedule.backup-hour", 0, "The hour of day (in UTC) when backup for the service is started. New backup is only started if previous backup has already completed.")
	var reqBackupScheduleBackupMinuteFlag int64
	flagset.Int64Var(&reqBackupScheduleBackupMinuteFlag, "backup-schedule.backup-minute", 0, "The minute of an hour when backup for the service is started. New backup is only started if previous backup has already completed.")
	var reqBinlogRetentionPeriodFlag int64
	flagset.Int64Var(&reqBinlogRetentionPeriodFlag, "binlog-retention-period", 0, "The minimum amount of time in seconds to keep binlog entries before deletion. This may be extended for services that require binlog entries for longer than the default for example if using the MySQL Debezium Kafka connector.")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqMigrationDbnameFlag string
	flagset.StringVar(&reqMigrationDbnameFlag, "migration.dbname", "", "Database name for bootstrapping the initial connection")
	var reqMigrationHostFlag string
	flagset.StringVar(&reqMigrationHostFlag, "migration.host", "", "Hostname or IP address of the server where to migrate data from")
	var reqMigrationIgnoreDbsFlag string
	flagset.StringVar(&reqMigrationIgnoreDbsFlag, "migration.ignore-dbs", "", "Comma-separated list of databases, which should be ignored during migration (supported by MySQL only at the moment)")
	var reqMigrationMethodFlag string
	flagset.StringVar(&reqMigrationMethodFlag, "migration.method", "", "")
	var reqMigrationPasswordFlag string
	flagset.StringVar(&reqMigrationPasswordFlag, "migration.password", "", "Password for authentication with the server where to migrate data from")
	var reqMigrationPortFlag int64
	flagset.Int64Var(&reqMigrationPortFlag, "migration.port", 0, "Port number of the server where to migrate data from")
	var reqMigrationSSLFlag bool
	flagset.BoolVar(&reqMigrationSSLFlag, "migration.ssl", false, "The server where to migrate data from is secured with SSL")
	var reqMigrationUsernameFlag string
	flagset.StringVar(&reqMigrationUsernameFlag, "migration.username", "", "User name for authentication with the server where to migrate data from")
	var reqMysqlSettingsConnectTimeoutFlag int
	flagset.IntVar(&reqMysqlSettingsConnectTimeoutFlag, "mysql-settings.connect_timeout", 0, "The number of seconds that the mysqld server waits for a connect packet before responding with Bad handshake")
	var reqMysqlSettingsDefaultTimeZoneFlag string
	flagset.StringVar(&reqMysqlSettingsDefaultTimeZoneFlag, "mysql-settings.default_time_zone", "", "Default server time zone as an offset from UTC (from -12:00 to +12:00), a time zone name, or 'SYSTEM' to use the MySQL server default.")
	var reqMysqlSettingsGroupConcatMaxLenFlag int
	flagset.IntVar(&reqMysqlSettingsGroupConcatMaxLenFlag, "mysql-settings.group_concat_max_len", 0, "The maximum permitted result length in bytes for the GROUP_CONCAT() function.")
	var reqMysqlSettingsInformationSchemaStatsExpiryFlag int
	flagset.IntVar(&reqMysqlSettingsInformationSchemaStatsExpiryFlag, "mysql-settings.information_schema_stats_expiry", 0, "The time, in seconds, before cached statistics expire")
	var reqMysqlSettingsInnodbChangeBufferMaxSizeFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbChangeBufferMaxSizeFlag, "mysql-settings.innodb_change_buffer_max_size", 0, "Maximum size for the InnoDB change buffer, as a percentage of the total size of the buffer pool. Default is 25")
	var reqMysqlSettingsInnodbFlushNeighborsFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbFlushNeighborsFlag, "mysql-settings.innodb_flush_neighbors", 0, "Specifies whether flushing a page from the InnoDB buffer pool also flushes other dirty pages in the same extent (default is 1): 0 - dirty pages in the same extent are not flushed, 1 - flush contiguous dirty pages in the same extent, 2 - flush dirty pages in the same extent")
	var reqMysqlSettingsInnodbFTMinTokenSizeFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbFTMinTokenSizeFlag, "mysql-settings.innodb_ft_min_token_size", 0, "Minimum length of words that are stored in an InnoDB FULLTEXT index. Changing this parameter will lead to a restart of the MySQL service.")
	var reqMysqlSettingsInnodbFTServerStopwordTableFlag string
	flagset.StringVar(&reqMysqlSettingsInnodbFTServerStopwordTableFlag, "mysql-settings.innodb_ft_server_stopword_table", "", "This option is used to specify your own InnoDB FULLTEXT index stopword list for all InnoDB tables.")
	var reqMysqlSettingsInnodbLockWaitTimeoutFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbLockWaitTimeoutFlag, "mysql-settings.innodb_lock_wait_timeout", 0, "The length of time in seconds an InnoDB transaction waits for a row lock before giving up. Default is 120.")
	var reqMysqlSettingsInnodbLogBufferSizeFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbLogBufferSizeFlag, "mysql-settings.innodb_log_buffer_size", 0, "The size in bytes of the buffer that InnoDB uses to write to the log files on disk.")
	var reqMysqlSettingsInnodbOnlineAlterLogMaxSizeFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbOnlineAlterLogMaxSizeFlag, "mysql-settings.innodb_online_alter_log_max_size", 0, "The upper limit in bytes on the size of the temporary log files used during online DDL operations for InnoDB tables.")
	var reqMysqlSettingsInnodbPrintAllDeadlocksFlag bool
	flagset.BoolVar(&reqMysqlSettingsInnodbPrintAllDeadlocksFlag, "mysql-settings.innodb_print_all_deadlocks", false, "When enabled, information about all deadlocks in InnoDB user transactions is recorded in the error log. Disabled by default.")
	var reqMysqlSettingsInnodbReadIoThreadsFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbReadIoThreadsFlag, "mysql-settings.innodb_read_io_threads", 0, "The number of I/O threads for read operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.")
	var reqMysqlSettingsInnodbRollbackOnTimeoutFlag bool
	flagset.BoolVar(&reqMysqlSettingsInnodbRollbackOnTimeoutFlag, "mysql-settings.innodb_rollback_on_timeout", false, "When enabled a transaction timeout causes InnoDB to abort and roll back the entire transaction. Changing this parameter will lead to a restart of the MySQL service.")
	var reqMysqlSettingsInnodbThreadConcurrencyFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbThreadConcurrencyFlag, "mysql-settings.innodb_thread_concurrency", 0, "Defines the maximum number of threads permitted inside of InnoDB. Default is 0 (infinite concurrency - no limit)")
	var reqMysqlSettingsInnodbWriteIoThreadsFlag int
	flagset.IntVar(&reqMysqlSettingsInnodbWriteIoThreadsFlag, "mysql-settings.innodb_write_io_threads", 0, "The number of I/O threads for write operations in InnoDB. Default is 4. Changing this parameter will lead to a restart of the MySQL service.")
	var reqMysqlSettingsInteractiveTimeoutFlag int
	flagset.IntVar(&reqMysqlSettingsInteractiveTimeoutFlag, "mysql-settings.interactive_timeout", 0, "The number of seconds the server waits for activity on an interactive connection before closing it.")
	var reqMysqlSettingsInternalTmpMemStorageEngineFlag string
	flagset.StringVar(&reqMysqlSettingsInternalTmpMemStorageEngineFlag, "mysql-settings.internal_tmp_mem_storage_engine", "", "The storage engine for in-memory internal temporary tables.")
	var reqMysqlSettingsLogOutputFlag string
	flagset.StringVar(&reqMysqlSettingsLogOutputFlag, "mysql-settings.log_output", "", "The slow log output destination when slow_query_log is ON. To enable MySQL AI Insights, choose INSIGHTS. To use MySQL AI Insights and the mysql.slow_log table at the same time, choose INSIGHTS,TABLE. To only use the mysql.slow_log table, choose TABLE. To silence slow logs, choose NONE.")
	var reqMysqlSettingsLongQueryTimeFlag float64
	flagset.Float64Var(&reqMysqlSettingsLongQueryTimeFlag, "mysql-settings.long_query_time", 0, "The slow_query_logs work as SQL statements that take more than long_query_time seconds to execute. Default is 10s")
	var reqMysqlSettingsMaxAllowedPacketFlag int
	flagset.IntVar(&reqMysqlSettingsMaxAllowedPacketFlag, "mysql-settings.max_allowed_packet", 0, "Size of the largest message in bytes that can be received by the server. Default is 67108864 (64M)")
	var reqMysqlSettingsMaxHeapTableSizeFlag int
	flagset.IntVar(&reqMysqlSettingsMaxHeapTableSizeFlag, "mysql-settings.max_heap_table_size", 0, "Limits the size of internal in-memory tables. Also set tmp_table_size. Default is 16777216 (16M)")
	var reqMysqlSettingsNetBufferLengthFlag int
	flagset.IntVar(&reqMysqlSettingsNetBufferLengthFlag, "mysql-settings.net_buffer_length", 0, "Start sizes of connection buffer and result buffer. Default is 16384 (16K). Changing this parameter will lead to a restart of the MySQL service.")
	var reqMysqlSettingsNetReadTimeoutFlag int
	flagset.IntVar(&reqMysqlSettingsNetReadTimeoutFlag, "mysql-settings.net_read_timeout", 0, "The number of seconds to wait for more data from a connection before aborting the read.")
	var reqMysqlSettingsNetWriteTimeoutFlag int
	flagset.IntVar(&reqMysqlSettingsNetWriteTimeoutFlag, "mysql-settings.net_write_timeout", 0, "The number of seconds to wait for a block to be written to a connection before aborting the write.")
	var reqMysqlSettingsSlowQueryLogFlag bool
	flagset.BoolVar(&reqMysqlSettingsSlowQueryLogFlag, "mysql-settings.slow_query_log", false, "Slow query log enables capturing of slow queries. Setting slow_query_log to false also truncates the mysql.slow_log table. Default is off")
	var reqMysqlSettingsSortBufferSizeFlag int
	flagset.IntVar(&reqMysqlSettingsSortBufferSizeFlag, "mysql-settings.sort_buffer_size", 0, "Sort buffer size in bytes for ORDER BY optimization. Default is 262144 (256K)")
	var reqMysqlSettingsSQLModeFlag string
	flagset.StringVar(&reqMysqlSettingsSQLModeFlag, "mysql-settings.sql_mode", "", "Global SQL mode. Set to empty to use MySQL server defaults. When creating a new service and not setting this field Aiven default SQL mode (strict, SQL standard compliant) will be assigned.")
	var reqMysqlSettingsSQLRequirePrimaryKeyFlag bool
	flagset.BoolVar(&reqMysqlSettingsSQLRequirePrimaryKeyFlag, "mysql-settings.sql_require_primary_key", false, "Require primary key to be defined for new tables or old tables modified with ALTER TABLE and fail if missing. It is recommended to always have primary keys because various functionality may break if any large table is missing them.")
	var reqMysqlSettingsTmpTableSizeFlag int
	flagset.IntVar(&reqMysqlSettingsTmpTableSizeFlag, "mysql-settings.tmp_table_size", 0, "Limits the size of internal in-memory tables. Also set max_heap_table_size. Default is 16777216 (16M)")
	var reqMysqlSettingsWaitTimeoutFlag int
	flagset.IntVar(&reqMysqlSettingsWaitTimeoutFlag, "mysql-settings.wait_timeout", 0, "The number of seconds the server waits for activity on a noninteractive connection before closing it.")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASServiceMysqlRequest
	req.TerminationProtection = reqTerminationProtectionFlag
	req.Plan = reqPlanFlag
	if reqMysqlSettingsWaitTimeoutFlag != 0 {
		req.MysqlSettings = &v3.JSONSchemaMysql{}
		req.MysqlSettings.WaitTimeout = reqMysqlSettingsWaitTimeoutFlag
	}
	if reqMysqlSettingsTmpTableSizeFlag != 0 {

		req.MysqlSettings.TmpTableSize = reqMysqlSettingsTmpTableSizeFlag
	}
	if reqMysqlSettingsSQLRequirePrimaryKeyFlag != false {

		req.MysqlSettings.SQLRequirePrimaryKey = reqMysqlSettingsSQLRequirePrimaryKeyFlag
	}
	if reqMysqlSettingsSQLModeFlag != "" {

		req.MysqlSettings.SQLMode = reqMysqlSettingsSQLModeFlag
	}
	if reqMysqlSettingsSortBufferSizeFlag != 0 {

		req.MysqlSettings.SortBufferSize = reqMysqlSettingsSortBufferSizeFlag
	}
	if reqMysqlSettingsSlowQueryLogFlag != false {

		req.MysqlSettings.SlowQueryLog = reqMysqlSettingsSlowQueryLogFlag
	}
	if reqMysqlSettingsNetWriteTimeoutFlag != 0 {

		req.MysqlSettings.NetWriteTimeout = reqMysqlSettingsNetWriteTimeoutFlag
	}
	if reqMysqlSettingsNetReadTimeoutFlag != 0 {

		req.MysqlSettings.NetReadTimeout = reqMysqlSettingsNetReadTimeoutFlag
	}
	if reqMysqlSettingsNetBufferLengthFlag != 0 {

		req.MysqlSettings.NetBufferLength = reqMysqlSettingsNetBufferLengthFlag
	}
	if reqMysqlSettingsMaxHeapTableSizeFlag != 0 {

		req.MysqlSettings.MaxHeapTableSize = reqMysqlSettingsMaxHeapTableSizeFlag
	}
	if reqMysqlSettingsMaxAllowedPacketFlag != 0 {

		req.MysqlSettings.MaxAllowedPacket = reqMysqlSettingsMaxAllowedPacketFlag
	}
	if reqMysqlSettingsLongQueryTimeFlag != 0 {

		req.MysqlSettings.LongQueryTime = reqMysqlSettingsLongQueryTimeFlag
	}
	if reqMysqlSettingsLogOutputFlag != "" {

		req.MysqlSettings.LogOutput = reqMysqlSettingsLogOutputFlag
	}
	if reqMysqlSettingsInternalTmpMemStorageEngineFlag != "" {

		req.MysqlSettings.InternalTmpMemStorageEngine = reqMysqlSettingsInternalTmpMemStorageEngineFlag
	}
	if reqMysqlSettingsInteractiveTimeoutFlag != 0 {

		req.MysqlSettings.InteractiveTimeout = reqMysqlSettingsInteractiveTimeoutFlag
	}
	if reqMysqlSettingsInnodbWriteIoThreadsFlag != 0 {

		req.MysqlSettings.InnodbWriteIoThreads = reqMysqlSettingsInnodbWriteIoThreadsFlag
	}
	if reqMysqlSettingsInnodbThreadConcurrencyFlag != 0 {

		req.MysqlSettings.InnodbThreadConcurrency = reqMysqlSettingsInnodbThreadConcurrencyFlag
	}
	if reqMysqlSettingsInnodbRollbackOnTimeoutFlag != false {

		req.MysqlSettings.InnodbRollbackOnTimeout = reqMysqlSettingsInnodbRollbackOnTimeoutFlag
	}
	if reqMysqlSettingsInnodbReadIoThreadsFlag != 0 {

		req.MysqlSettings.InnodbReadIoThreads = reqMysqlSettingsInnodbReadIoThreadsFlag
	}
	if reqMysqlSettingsInnodbPrintAllDeadlocksFlag != false {

		req.MysqlSettings.InnodbPrintAllDeadlocks = reqMysqlSettingsInnodbPrintAllDeadlocksFlag
	}
	if reqMysqlSettingsInnodbOnlineAlterLogMaxSizeFlag != 0 {

		req.MysqlSettings.InnodbOnlineAlterLogMaxSize = reqMysqlSettingsInnodbOnlineAlterLogMaxSizeFlag
	}
	if reqMysqlSettingsInnodbLogBufferSizeFlag != 0 {

		req.MysqlSettings.InnodbLogBufferSize = reqMysqlSettingsInnodbLogBufferSizeFlag
	}
	if reqMysqlSettingsInnodbLockWaitTimeoutFlag != 0 {

		req.MysqlSettings.InnodbLockWaitTimeout = reqMysqlSettingsInnodbLockWaitTimeoutFlag
	}
	if reqMysqlSettingsInnodbFTServerStopwordTableFlag != "" {

		req.MysqlSettings.InnodbFTServerStopwordTable = reqMysqlSettingsInnodbFTServerStopwordTableFlag
	}
	if reqMysqlSettingsInnodbFTMinTokenSizeFlag != 0 {

		req.MysqlSettings.InnodbFTMinTokenSize = reqMysqlSettingsInnodbFTMinTokenSizeFlag
	}
	if reqMysqlSettingsInnodbFlushNeighborsFlag != 0 {

		req.MysqlSettings.InnodbFlushNeighbors = reqMysqlSettingsInnodbFlushNeighborsFlag
	}
	if reqMysqlSettingsInnodbChangeBufferMaxSizeFlag != 0 {

		req.MysqlSettings.InnodbChangeBufferMaxSize = reqMysqlSettingsInnodbChangeBufferMaxSizeFlag
	}
	if reqMysqlSettingsInformationSchemaStatsExpiryFlag != 0 {

		req.MysqlSettings.InformationSchemaStatsExpiry = reqMysqlSettingsInformationSchemaStatsExpiryFlag
	}
	if reqMysqlSettingsGroupConcatMaxLenFlag != 0 {

		req.MysqlSettings.GroupConcatMaxLen = reqMysqlSettingsGroupConcatMaxLenFlag
	}
	if reqMysqlSettingsDefaultTimeZoneFlag != "" {

		req.MysqlSettings.DefaultTimeZone = reqMysqlSettingsDefaultTimeZoneFlag
	}
	if reqMysqlSettingsConnectTimeoutFlag != 0 {

		req.MysqlSettings.ConnectTimeout = reqMysqlSettingsConnectTimeoutFlag
	}
	if reqMigrationUsernameFlag != "" {
		req.Migratio = &v3.UpdateDBAASServiceMysqlRequestMigration{}
		req.Migration.Username = reqMigrationUsernameFlag
	}
	if reqMigrationSSLFlag != false {

		req.Migration.SSL = reqMigrationSSLFlag
	}
	if reqMigrationPortFlag != 0 {

		req.Migration.Port = reqMigrationPortFlag
	}
	if reqMigrationPasswordFlag != "" {

		req.Migration.Password = reqMigrationPasswordFlag
	}
	if reqMigrationMethodFlag != "" {

		req.Migration.Method = reqMigrationMethodFlag
	}
	if reqMigrationIgnoreDbsFlag != "" {

		req.Migration.IgnoreDbs = reqMigrationIgnoreDbsFlag
	}
	if reqMigrationHostFlag != "" {

		req.Migration.Host = reqMigrationHostFlag
	}
	if reqMigrationDbnameFlag != "" {

		req.Migration.Dbname = reqMigrationDbnameFlag
	}
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.UpdateDBAASServiceMysqlRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}
	req.BinlogRetentionPeriod = reqBinlogRetentionPeriodFlag
	if reqBackupScheduleBackupMinuteFlag != 0 {
		req.BackupSchedul = &v3.UpdateDBAASServiceMysqlRequestBackupSchedule{}
		req.BackupSchedule.BackupMinute = reqBackupScheduleBackupMinuteFlag
	}
	if reqBackupScheduleBackupHourFlag != 0 {

		req.BackupSchedule.BackupHour = reqBackupScheduleBackupHourFlag
	}

	resp, err := client.UpdateDBAASServiceMysql(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func EnableDBAASMysqlWritesCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("enable-dbaas-mysql-writes", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.EnableDBAASMysqlWrites(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StartDBAASMysqlMaintenanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("start-dbaas-mysql-maintenance", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StartDBAASMysqlMaintenance(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StopDBAASMysqlMigrationCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("stop-dbaas-mysql-migration", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StopDBAASMysqlMigration(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASMysqlDatabaseCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-mysql-database", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var reqDatabaseNameFlag string
	flagset.StringVar(&reqDatabaseNameFlag, "database-name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASMysqlDatabaseRequest
	req.DatabaseName = reqDatabaseNameFlag

	resp, err := client.CreateDBAASMysqlDatabase(context.Background(), serviceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASMysqlDatabaseCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-mysql-database", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var databaseNameFlag string
	flagset.StringVar(&databaseNameFlag, "DatabaseName", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASMysqlDatabase(context.Background(), serviceNameFlag, databaseNameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASMysqlUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-mysql-user", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var reqAuthenticationFlag string
	flagset.StringVar(&reqAuthenticationFlag, "authentication", "", "")
	var reqUsernameFlag string
	flagset.StringVar(&reqUsernameFlag, "username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASMysqlUserRequest
	req.Username = reqUsernameFlag
	req.Authentication = reqAuthenticationFlag

	resp, err := client.CreateDBAASMysqlUser(context.Background(), serviceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASMysqlUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-mysql-user", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASMysqlUser(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetDBAASMysqlUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-dbaas-mysql-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")
	var reqAuthenticationFlag string
	flagset.StringVar(&reqAuthenticationFlag, "authentication", "", "")
	var reqPasswordFlag string
	flagset.StringVar(&reqPasswordFlag, "password", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ResetDBAASMysqlUserPasswordRequest
	req.Password = reqPasswordFlag
	req.Authentication = reqAuthenticationFlag

	resp, err := client.ResetDBAASMysqlUserPassword(context.Background(), serviceNameFlag, usernameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RevealDBAASMysqlUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reveal-dbaas-mysql-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RevealDBAASMysqlUserPassword(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASServiceOpensearchCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-service-opensearch", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASServiceOpensearch(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASServiceOpensearchCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-service-opensearch", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASServiceOpensearch(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASServiceOpensearchCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-service-opensearch", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqForkFromServiceFlag string
	flagset.StringVar(&reqForkFromServiceFlag, "fork-from-service", "", "")
	var reqIndexTemplateMappingNestedObjectsLimitFlag int64
	flagset.Int64Var(&reqIndexTemplateMappingNestedObjectsLimitFlag, "index-template.mapping-nested-objects-limit", 0, "The maximum number of nested JSON objects that a single document can contain across all nested types. This limit helps to prevent out of memory errors when a document contains too many nested objects. Default is 10000.")
	var reqIndexTemplateNumberOfReplicasFlag int64
	flagset.Int64Var(&reqIndexTemplateNumberOfReplicasFlag, "index-template.number-of-replicas", 0, "The number of replicas each primary shard has.")
	var reqIndexTemplateNumberOfShardsFlag int64
	flagset.Int64Var(&reqIndexTemplateNumberOfShardsFlag, "index-template.number-of-shards", 0, "The number of primary shards that an index should have.")
	var reqKeepIndexRefreshIntervalFlag bool
	flagset.BoolVar(&reqKeepIndexRefreshIntervalFlag, "keep-index-refresh-interval", false, "Aiven automation resets index.refresh_interval to default value for every index to be sure that indices are always visible to search. If it doesn't fit your case, you can disable this by setting up this flag to true.")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqMaxIndexCountFlag int64
	flagset.Int64Var(&reqMaxIndexCountFlag, "max-index-count", 0, "Maximum number of indexes to keep before deleting the oldest one")
	var reqOpensearchDashboardsEnabledFlag bool
	flagset.BoolVar(&reqOpensearchDashboardsEnabledFlag, "opensearch-dashboards.enabled", false, "Enable or disable OpenSearch Dashboards (default: true)")
	var reqOpensearchDashboardsMaxOldSpaceSizeFlag int64
	flagset.Int64Var(&reqOpensearchDashboardsMaxOldSpaceSizeFlag, "opensearch-dashboards.max-old-space-size", 0, "Limits the maximum amount of memory (in MiB) the OpenSearch Dashboards process can use. This sets the max_old_space_size option of the nodejs running the OpenSearch Dashboards. Note: the memory reserved by OpenSearch Dashboards is not available for OpenSearch. (default: 128)")
	var reqOpensearchDashboardsOpensearchRequestTimeoutFlag int64
	flagset.Int64Var(&reqOpensearchDashboardsOpensearchRequestTimeoutFlag, "opensearch-dashboards.opensearch-request-timeout", 0, "Timeout in milliseconds for requests made by OpenSearch Dashboards towards OpenSearch (default: 30000)")
	var reqOpensearchSettingsActionAutoCreateIndexEnabledFlag bool
	flagset.BoolVar(&reqOpensearchSettingsActionAutoCreateIndexEnabledFlag, "opensearch-settings.action_auto_create_index_enabled", false, "Explicitly allow or block automatic creation of indices. Defaults to true")
	var reqOpensearchSettingsActionDestructiveRequiresNameFlag bool
	flagset.BoolVar(&reqOpensearchSettingsActionDestructiveRequiresNameFlag, "opensearch-settings.action_destructive_requires_name", false, "Require explicit index names when deleting")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAllowedTriesFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAllowedTriesFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.allowed_tries", 0, "The number of login attempts allowed before login is blocked")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAuthenticationBackendFlag string
	flagset.StringVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAuthenticationBackendFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.authentication_backend", "", "The internal backend. Enter `internal`")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingBlockExpirySecondsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingBlockExpirySecondsFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.block_expiry_seconds", 0, "The duration of time that login remains blocked after a failed login")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxBlockedClientsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxBlockedClientsFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.max_blocked_clients", 0, "The maximum number of blocked IP addresses")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxTrackedClientsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxTrackedClientsFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.max_tracked_clients", 0, "The maximum number of tracked IP addresses that have failed login")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTimeWindowSecondsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTimeWindowSecondsFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.time_window_seconds", 0, "The window of time in which the value for `allowed_tries` is enforced")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTypeFlag string
	flagset.StringVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTypeFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.type", "", "The type of rate limiting")
	var reqOpensearchSettingsAuthFailureListenersIPRateLimitingAllowedTriesFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersIPRateLimitingAllowedTriesFlag, "opensearch-settings.auth_failure_listeners.ip_rate_limiting.allowed_tries", 0, "The number of login attempts allowed before login is blocked")
	var reqOpensearchSettingsAuthFailureListenersIPRateLimitingBlockExpirySecondsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersIPRateLimitingBlockExpirySecondsFlag, "opensearch-settings.auth_failure_listeners.ip_rate_limiting.block_expiry_seconds", 0, "The duration of time that login remains blocked after a failed login")
	var reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxBlockedClientsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxBlockedClientsFlag, "opensearch-settings.auth_failure_listeners.ip_rate_limiting.max_blocked_clients", 0, "The maximum number of blocked IP addresses")
	var reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxTrackedClientsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxTrackedClientsFlag, "opensearch-settings.auth_failure_listeners.ip_rate_limiting.max_tracked_clients", 0, "The maximum number of tracked IP addresses that have failed login")
	var reqOpensearchSettingsAuthFailureListenersIPRateLimitingTimeWindowSecondsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersIPRateLimitingTimeWindowSecondsFlag, "opensearch-settings.auth_failure_listeners.ip_rate_limiting.time_window_seconds", 0, "The window of time in which the value for `allowed_tries` is enforced")
	var reqOpensearchSettingsAuthFailureListenersIPRateLimitingTypeFlag string
	flagset.StringVar(&reqOpensearchSettingsAuthFailureListenersIPRateLimitingTypeFlag, "opensearch-settings.auth_failure_listeners.ip_rate_limiting.type", "", "The type of rate limiting")
	var reqOpensearchSettingsClusterMaxShardsPerNodeFlag int
	flagset.IntVar(&reqOpensearchSettingsClusterMaxShardsPerNodeFlag, "opensearch-settings.cluster_max_shards_per_node", 0, "Controls the number of shards allowed in the cluster per data node")
	var reqOpensearchSettingsClusterRoutingAllocationNodeConcurrentRecoveriesFlag int
	flagset.IntVar(&reqOpensearchSettingsClusterRoutingAllocationNodeConcurrentRecoveriesFlag, "opensearch-settings.cluster_routing_allocation_node_concurrent_recoveries", 0, "How many concurrent incoming/outgoing shard recoveries (normally replicas) are allowed to happen on a node. Defaults to 2.")
	var reqOpensearchSettingsEmailSenderEmailSenderNameFlag string
	flagset.StringVar(&reqOpensearchSettingsEmailSenderEmailSenderNameFlag, "opensearch-settings.email-sender.email_sender_name", "", "This should be identical to the Sender name defined in Opensearch dashboards")
	var reqOpensearchSettingsEmailSenderEmailSenderPasswordFlag string
	flagset.StringVar(&reqOpensearchSettingsEmailSenderEmailSenderPasswordFlag, "opensearch-settings.email-sender.email_sender_password", "", "Sender password for Opensearch alerts to authenticate with SMTP server")
	var reqOpensearchSettingsEmailSenderEmailSenderUsernameFlag string
	flagset.StringVar(&reqOpensearchSettingsEmailSenderEmailSenderUsernameFlag, "opensearch-settings.email-sender.email_sender_username", "", "Sender username for Opensearch alerts")
	var reqOpensearchSettingsEnableSecurityAuditFlag bool
	flagset.BoolVar(&reqOpensearchSettingsEnableSecurityAuditFlag, "opensearch-settings.enable_security_audit", false, "Enable/Disable security audit")
	var reqOpensearchSettingsHTTPMaxContentLengthFlag int
	flagset.IntVar(&reqOpensearchSettingsHTTPMaxContentLengthFlag, "opensearch-settings.http_max_content_length", 0, "Maximum content length for HTTP requests to the OpenSearch HTTP API, in bytes.")
	var reqOpensearchSettingsHTTPMaxHeaderSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsHTTPMaxHeaderSizeFlag, "opensearch-settings.http_max_header_size", 0, "The max size of allowed headers, in bytes")
	var reqOpensearchSettingsHTTPMaxInitialLineLengthFlag int
	flagset.IntVar(&reqOpensearchSettingsHTTPMaxInitialLineLengthFlag, "opensearch-settings.http_max_initial_line_length", 0, "The max length of an HTTP URL, in bytes")
	var reqOpensearchSettingsIndicesFielddataCacheSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesFielddataCacheSizeFlag, "opensearch-settings.indices_fielddata_cache_size", 0, "Relative amount. Maximum amount of heap memory used for field data cache. This is an expert setting; decreasing the value too much will increase overhead of loading field data; too much memory used for field data cache will decrease amount of heap available for other operations.")
	var reqOpensearchSettingsIndicesMemoryIndexBufferSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesMemoryIndexBufferSizeFlag, "opensearch-settings.indices_memory_index_buffer_size", 0, "Percentage value. Default is 10%. Total amount of heap used for indexing buffer, before writing segments to disk. This is an expert setting. Too low value will slow down indexing; too high value will increase indexing performance but causes performance issues for query performance.")
	var reqOpensearchSettingsIndicesMemoryMaxIndexBufferSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesMemoryMaxIndexBufferSizeFlag, "opensearch-settings.indices_memory_max_index_buffer_size", 0, "Absolute value. Default is unbound. Doesn't work without indices.memory.index_buffer_size. Maximum amount of heap used for query cache, an absolute indices.memory.index_buffer_size maximum hard limit.")
	var reqOpensearchSettingsIndicesMemoryMinIndexBufferSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesMemoryMinIndexBufferSizeFlag, "opensearch-settings.indices_memory_min_index_buffer_size", 0, "Absolute value. Default is 48mb. Doesn't work without indices.memory.index_buffer_size. Minimum amount of heap used for query cache, an absolute indices.memory.index_buffer_size minimal hard limit.")
	var reqOpensearchSettingsIndicesQueriesCacheSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesQueriesCacheSizeFlag, "opensearch-settings.indices_queries_cache_size", 0, "Percentage value. Default is 10%. Maximum amount of heap used for query cache. This is an expert setting. Too low value will decrease query performance and increase performance for other operations; too high value will cause issues with other OpenSearch functionality.")
	var reqOpensearchSettingsIndicesQueryBoolMaxClauseCountFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesQueryBoolMaxClauseCountFlag, "opensearch-settings.indices_query_bool_max_clause_count", 0, "Maximum number of clauses Lucene BooleanQuery can have. The default value (1024) is relatively high, and increasing it may cause performance issues. Investigate other approaches first before increasing this value.")
	var reqOpensearchSettingsIndicesRecoveryMaxBytesPerSecFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesRecoveryMaxBytesPerSecFlag, "opensearch-settings.indices_recovery_max_bytes_per_sec", 0, "Limits total inbound and outbound recovery traffic for each node. Applies to both peer recoveries as well as snapshot recoveries (i.e., restores from a snapshot). Defaults to 40mb")
	var reqOpensearchSettingsIndicesRecoveryMaxConcurrentFileChunksFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesRecoveryMaxConcurrentFileChunksFlag, "opensearch-settings.indices_recovery_max_concurrent_file_chunks", 0, "Number of file chunks sent in parallel for each recovery. Defaults to 2.")
	var reqOpensearchSettingsIsmHistoryIsmEnabledFlag bool
	flagset.BoolVar(&reqOpensearchSettingsIsmHistoryIsmEnabledFlag, "opensearch-settings.ism-history.ism_enabled", false, "Specifies whether ISM is enabled or not")
	var reqOpensearchSettingsIsmHistoryIsmHistoryEnabledFlag bool
	flagset.BoolVar(&reqOpensearchSettingsIsmHistoryIsmHistoryEnabledFlag, "opensearch-settings.ism-history.ism_history_enabled", false, "Specifies whether audit history is enabled or not. The logs from ISM are automatically indexed to a logs document.")
	var reqOpensearchSettingsIsmHistoryIsmHistoryMaxAgeFlag int
	flagset.IntVar(&reqOpensearchSettingsIsmHistoryIsmHistoryMaxAgeFlag, "opensearch-settings.ism-history.ism_history_max_age", 0, "The maximum age before rolling over the audit history index in hours")
	var reqOpensearchSettingsIsmHistoryIsmHistoryMaxDocsFlag int
	flagset.IntVar(&reqOpensearchSettingsIsmHistoryIsmHistoryMaxDocsFlag, "opensearch-settings.ism-history.ism_history_max_docs", 0, "The maximum number of documents before rolling over the audit history index.")
	var reqOpensearchSettingsIsmHistoryIsmHistoryRolloverCheckPeriodFlag int
	flagset.IntVar(&reqOpensearchSettingsIsmHistoryIsmHistoryRolloverCheckPeriodFlag, "opensearch-settings.ism-history.ism_history_rollover_check_period", 0, "The time between rollover checks for the audit history index in hours.")
	var reqOpensearchSettingsIsmHistoryIsmHistoryRolloverRetentionPeriodFlag int
	flagset.IntVar(&reqOpensearchSettingsIsmHistoryIsmHistoryRolloverRetentionPeriodFlag, "opensearch-settings.ism-history.ism_history_rollover_retention_period", 0, "How long audit history indices are kept in days.")
	var reqOpensearchSettingsKnnMemoryCircuitBreakerEnabledFlag bool
	flagset.BoolVar(&reqOpensearchSettingsKnnMemoryCircuitBreakerEnabledFlag, "opensearch-settings.knn_memory_circuit_breaker_enabled", false, "Enable or disable KNN memory circuit breaker. Defaults to true.")
	var reqOpensearchSettingsKnnMemoryCircuitBreakerLimitFlag int
	flagset.IntVar(&reqOpensearchSettingsKnnMemoryCircuitBreakerLimitFlag, "opensearch-settings.knn_memory_circuit_breaker_limit", 0, "Maximum amount of memory that can be used for KNN index. Defaults to 50% of the JVM heap size.")
	var reqOpensearchSettingsOverrideMainResponseVersionFlag bool
	flagset.BoolVar(&reqOpensearchSettingsOverrideMainResponseVersionFlag, "opensearch-settings.override_main_response_version", false, "Compatibility mode sets OpenSearch to report its version as 7.10 so clients continue to work. Default is false")
	var reqOpensearchSettingsPluginsAlertingFilterByBackendRolesFlag bool
	flagset.BoolVar(&reqOpensearchSettingsPluginsAlertingFilterByBackendRolesFlag, "opensearch-settings.plugins_alerting_filter_by_backend_roles", false, "Enable or disable filtering of alerting by backend roles. Requires Security plugin. Defaults to false")
	var reqOpensearchSettingsScriptMaxCompilationsRateFlag string
	flagset.StringVar(&reqOpensearchSettingsScriptMaxCompilationsRateFlag, "opensearch-settings.script_max_compilations_rate", "", "Script compilation circuit breaker limits the number of inline script compilations within a period of time. Default is use-context")
	var reqOpensearchSettingsSearchBackpressureModeFlag string
	flagset.StringVar(&reqOpensearchSettingsSearchBackpressureModeFlag, "opensearch-settings.search_backpressure.mode", "", "The search backpressure mode. Valid values are monitor_only, enforced, or disabled. Default is monitor_only")
	var reqOpensearchSettingsSearchBackpressureNodeDuressCPUThresholdFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureNodeDuressCPUThresholdFlag, "opensearch-settings.search_backpressure.node_duress.cpu_threshold", 0, "The CPU usage threshold (as a percentage) required for a node to be considered to be under duress. Default is 0.9")
	var reqOpensearchSettingsSearchBackpressureNodeDuressHeapThresholdFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureNodeDuressHeapThresholdFlag, "opensearch-settings.search_backpressure.node_duress.heap_threshold", 0, "The heap usage threshold (as a percentage) required for a node to be considered to be under duress. Default is 0.7")
	var reqOpensearchSettingsSearchBackpressureNodeDuressNumSuccessiveBreachesFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureNodeDuressNumSuccessiveBreachesFlag, "opensearch-settings.search_backpressure.node_duress.num_successive_breaches", 0, "The number of successive limit breaches after which the node is considered to be under duress. Default is 3")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationBurstFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationBurstFlag, "opensearch-settings.search_backpressure.search_shard_task.cancellation_burst", 0, "The maximum number of search tasks to cancel in a single iteration of the observer thread. Default is 10.0")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRateFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRateFlag, "opensearch-settings.search_backpressure.search_shard_task.cancellation_rate", 0, "The maximum number of tasks to cancel per millisecond of elapsed time. Default is 0.003")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRatioFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRatioFlag, "opensearch-settings.search_backpressure.search_shard_task.cancellation_ratio", 0, "The maximum number of tasks to cancel, as a percentage of successful task completions. Default is 0.1")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskCPUTimeMillisThresholdFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureSearchShardTaskCPUTimeMillisThresholdFlag, "opensearch-settings.search_backpressure.search_shard_task.cpu_time_millis_threshold", 0, "The CPU usage threshold (in milliseconds) required for a single search shard task before it is considered for cancellation. Default is 15000")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskElapsedTimeMillisThresholdFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureSearchShardTaskElapsedTimeMillisThresholdFlag, "opensearch-settings.search_backpressure.search_shard_task.elapsed_time_millis_threshold", 0, "The elapsed time threshold (in milliseconds) required for a single search shard task before it is considered for cancellation. Default is 30000")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapMovingAverageWindowSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapMovingAverageWindowSizeFlag, "opensearch-settings.search_backpressure.search_shard_task.heap_moving_average_window_size", 0, "The number of previously completed search shard tasks to consider when calculating the rolling average of heap usage. Default is 100")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapPercentThresholdFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapPercentThresholdFlag, "opensearch-settings.search_backpressure.search_shard_task.heap_percent_threshold", 0, "The heap usage threshold (as a percentage) required for a single search shard task before it is considered for cancellation. Default is 0.5")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapVarianceFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapVarianceFlag, "opensearch-settings.search_backpressure.search_shard_task.heap_variance", 0, "The minimum variance required for a single search shard task’s heap usage compared to the rolling average of previously completed tasks before it is considered for cancellation. Default is 2.0")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskTotalHeapPercentThresholdFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchShardTaskTotalHeapPercentThresholdFlag, "opensearch-settings.search_backpressure.search_shard_task.total_heap_percent_threshold", 0, "The heap usage threshold (as a percentage) required for the sum of heap usages of all search shard tasks before cancellation is applied. Default is 0.5")
	var reqOpensearchSettingsSearchBackpressureSearchTaskCancellationBurstFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchTaskCancellationBurstFlag, "opensearch-settings.search_backpressure.search_task.cancellation_burst", 0, "The maximum number of search tasks to cancel in a single iteration of the observer thread. Default is 5.0")
	var reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRateFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRateFlag, "opensearch-settings.search_backpressure.search_task.cancellation_rate", 0, "The maximum number of search tasks to cancel per millisecond of elapsed time. Default is 0.003")
	var reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRatioFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRatioFlag, "opensearch-settings.search_backpressure.search_task.cancellation_ratio", 0, "The maximum number of search tasks to cancel, as a percentage of successful search task completions. Default is 0.1")
	var reqOpensearchSettingsSearchBackpressureSearchTaskCPUTimeMillisThresholdFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureSearchTaskCPUTimeMillisThresholdFlag, "opensearch-settings.search_backpressure.search_task.cpu_time_millis_threshold", 0, "The CPU usage threshold (in milliseconds) required for an individual parent task before it is considered for cancellation. Default is 30000")
	var reqOpensearchSettingsSearchBackpressureSearchTaskElapsedTimeMillisThresholdFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureSearchTaskElapsedTimeMillisThresholdFlag, "opensearch-settings.search_backpressure.search_task.elapsed_time_millis_threshold", 0, "The elapsed time threshold (in milliseconds) required for an individual parent task before it is considered for cancellation. Default is 45000")
	var reqOpensearchSettingsSearchBackpressureSearchTaskHeapMovingAverageWindowSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureSearchTaskHeapMovingAverageWindowSizeFlag, "opensearch-settings.search_backpressure.search_task.heap_moving_average_window_size", 0, "The window size used to calculate the rolling average of the heap usage for the completed parent tasks. Default is 10")
	var reqOpensearchSettingsSearchBackpressureSearchTaskHeapPercentThresholdFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchTaskHeapPercentThresholdFlag, "opensearch-settings.search_backpressure.search_task.heap_percent_threshold", 0, "The heap usage threshold (as a percentage) required for an individual parent task before it is considered for cancellation. Default is 0.2")
	var reqOpensearchSettingsSearchBackpressureSearchTaskHeapVarianceFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchTaskHeapVarianceFlag, "opensearch-settings.search_backpressure.search_task.heap_variance", 0, "The heap usage variance required for an individual parent task before it is considered for cancellation. A task is considered for cancellation when taskHeapUsage is greater than or equal to heapUsageMovingAverage * variance. Default is 2.0")
	var reqOpensearchSettingsSearchBackpressureSearchTaskTotalHeapPercentThresholdFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchTaskTotalHeapPercentThresholdFlag, "opensearch-settings.search_backpressure.search_task.total_heap_percent_threshold", 0, "The heap usage threshold (as a percentage) required for the sum of heap usages of all search tasks before cancellation is applied. Default is 0.5")
	var reqOpensearchSettingsSearchMaxBucketsFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchMaxBucketsFlag, "opensearch-settings.search_max_buckets", 0, "Maximum number of aggregation buckets allowed in a single response. OpenSearch default value is used when this is not defined.")
	var reqOpensearchSettingsShardIndexingPressureEnabledFlag bool
	flagset.BoolVar(&reqOpensearchSettingsShardIndexingPressureEnabledFlag, "opensearch-settings.shard_indexing_pressure.enabled", false, "Enable or disable shard indexing backpressure. Default is false")
	var reqOpensearchSettingsShardIndexingPressureEnforcedFlag bool
	flagset.BoolVar(&reqOpensearchSettingsShardIndexingPressureEnforcedFlag, "opensearch-settings.shard_indexing_pressure.enforced", false, "Run shard indexing backpressure in shadow mode or enforced mode. In shadow mode (value set as false), shard indexing backpressure tracks all granular-level metrics, but it doesn’t actually reject any indexing requests. In enforced mode (value set as true), shard indexing backpressure rejects any requests to the cluster that might cause a dip in its performance. Default is false")
	var reqOpensearchSettingsShardIndexingPressureOperatingFactorLowerFlag float64
	flagset.Float64Var(&reqOpensearchSettingsShardIndexingPressureOperatingFactorLowerFlag, "opensearch-settings.shard_indexing_pressure.operating_factor.lower", 0, "Specify the lower occupancy limit of the allocated quota of memory for the shard. If the total memory usage of a shard is below this limit, shard indexing backpressure decreases the current allocated memory for that shard. Default is 0.75")
	var reqOpensearchSettingsShardIndexingPressureOperatingFactorOptimalFlag float64
	flagset.Float64Var(&reqOpensearchSettingsShardIndexingPressureOperatingFactorOptimalFlag, "opensearch-settings.shard_indexing_pressure.operating_factor.optimal", 0, "Specify the optimal occupancy of the allocated quota of memory for the shard. If the total memory usage of a shard is at this level, shard indexing backpressure doesn’t change the current allocated memory for that shard. Default is 0.85")
	var reqOpensearchSettingsShardIndexingPressureOperatingFactorUpperFlag float64
	flagset.Float64Var(&reqOpensearchSettingsShardIndexingPressureOperatingFactorUpperFlag, "opensearch-settings.shard_indexing_pressure.operating_factor.upper", 0, "Specify the upper occupancy limit of the allocated quota of memory for the shard. If the total memory usage of a shard is above this limit, shard indexing backpressure increases the current allocated memory for that shard. Default is 0.95")
	var reqOpensearchSettingsShardIndexingPressurePrimaryParameterNodeSoftLimitFlag float64
	flagset.Float64Var(&reqOpensearchSettingsShardIndexingPressurePrimaryParameterNodeSoftLimitFlag, "opensearch-settings.shard_indexing_pressure.primary_parameter.node.soft_limit", 0, "Define the percentage of the node-level memory threshold that acts as a soft indicator for strain on a node. Default is 0.7")
	var reqOpensearchSettingsShardIndexingPressurePrimaryParameterShardMinLimitFlag float64
	flagset.Float64Var(&reqOpensearchSettingsShardIndexingPressurePrimaryParameterShardMinLimitFlag, "opensearch-settings.shard_indexing_pressure.primary_parameter.shard.min_limit", 0, "Specify the minimum assigned quota for a new shard in any role (coordinator, primary, or replica). Shard indexing backpressure increases or decreases this allocated quota based on the inflow of traffic for the shard. Default is 0.001")
	var reqOpensearchSettingsThreadPoolAnalyzeQueueSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolAnalyzeQueueSizeFlag, "opensearch-settings.thread_pool_analyze_queue_size", 0, "Size for the thread pool queue. See documentation for exact details.")
	var reqOpensearchSettingsThreadPoolAnalyzeSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolAnalyzeSizeFlag, "opensearch-settings.thread_pool_analyze_size", 0, "Size for the thread pool. See documentation for exact details. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.")
	var reqOpensearchSettingsThreadPoolForceMergeSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolForceMergeSizeFlag, "opensearch-settings.thread_pool_force_merge_size", 0, "Size for the thread pool. See documentation for exact details. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.")
	var reqOpensearchSettingsThreadPoolGetQueueSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolGetQueueSizeFlag, "opensearch-settings.thread_pool_get_queue_size", 0, "Size for the thread pool queue. See documentation for exact details.")
	var reqOpensearchSettingsThreadPoolGetSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolGetSizeFlag, "opensearch-settings.thread_pool_get_size", 0, "Size for the thread pool. See documentation for exact details. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.")
	var reqOpensearchSettingsThreadPoolSearchQueueSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolSearchQueueSizeFlag, "opensearch-settings.thread_pool_search_queue_size", 0, "Size for the thread pool queue. See documentation for exact details.")
	var reqOpensearchSettingsThreadPoolSearchSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolSearchSizeFlag, "opensearch-settings.thread_pool_search_size", 0, "Size for the thread pool. See documentation for exact details. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.")
	var reqOpensearchSettingsThreadPoolSearchThrottledQueueSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolSearchThrottledQueueSizeFlag, "opensearch-settings.thread_pool_search_throttled_queue_size", 0, "Size for the thread pool queue. See documentation for exact details.")
	var reqOpensearchSettingsThreadPoolSearchThrottledSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolSearchThrottledSizeFlag, "opensearch-settings.thread_pool_search_throttled_size", 0, "Size for the thread pool. See documentation for exact details. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.")
	var reqOpensearchSettingsThreadPoolWriteQueueSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolWriteQueueSizeFlag, "opensearch-settings.thread_pool_write_queue_size", 0, "Size for the thread pool queue. See documentation for exact details.")
	var reqOpensearchSettingsThreadPoolWriteSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolWriteSizeFlag, "opensearch-settings.thread_pool_write_size", 0, "Size for the thread pool. See documentation for exact details. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqRecoveryBackupNameFlag string
	flagset.StringVar(&reqRecoveryBackupNameFlag, "recovery-backup-name", "", "Name of a backup to recover from for services that support backup names")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")
	var reqVersionFlag string
	flagset.StringVar(&reqVersionFlag, "version", "", "OpenSearch major version")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASServiceOpensearchRequest
	req.Version = reqVersionFlag
	req.TerminationProtection = reqTerminationProtectionFlag
	req.RecoveryBackupName = reqRecoveryBackupNameFlag
	req.Plan = reqPlanFlag
	if reqOpensearchSettingsThreadPoolWriteSizeFlag != 0 {
		req.OpensearchSettings = &v3.JSONSchemaOpensearch{}
		req.OpensearchSettings.ThreadPoolWriteSize = reqOpensearchSettingsThreadPoolWriteSizeFlag
	}
	if reqOpensearchSettingsThreadPoolWriteQueueSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolWriteQueueSize = reqOpensearchSettingsThreadPoolWriteQueueSizeFlag
	}
	if reqOpensearchSettingsThreadPoolSearchThrottledSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolSearchThrottledSize = reqOpensearchSettingsThreadPoolSearchThrottledSizeFlag
	}
	if reqOpensearchSettingsThreadPoolSearchThrottledQueueSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolSearchThrottledQueueSize = reqOpensearchSettingsThreadPoolSearchThrottledQueueSizeFlag
	}
	if reqOpensearchSettingsThreadPoolSearchSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolSearchSize = reqOpensearchSettingsThreadPoolSearchSizeFlag
	}
	if reqOpensearchSettingsThreadPoolSearchQueueSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolSearchQueueSize = reqOpensearchSettingsThreadPoolSearchQueueSizeFlag
	}
	if reqOpensearchSettingsThreadPoolGetSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolGetSize = reqOpensearchSettingsThreadPoolGetSizeFlag
	}
	if reqOpensearchSettingsThreadPoolGetQueueSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolGetQueueSize = reqOpensearchSettingsThreadPoolGetQueueSizeFlag
	}
	if reqOpensearchSettingsThreadPoolForceMergeSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolForceMergeSize = reqOpensearchSettingsThreadPoolForceMergeSizeFlag
	}
	if reqOpensearchSettingsThreadPoolAnalyzeSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolAnalyzeSize = reqOpensearchSettingsThreadPoolAnalyzeSizeFlag
	}
	if reqOpensearchSettingsThreadPoolAnalyzeQueueSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolAnalyzeQueueSize = reqOpensearchSettingsThreadPoolAnalyzeQueueSizeFlag
	}
	if reqOpensearchSettingsShardIndexingPressurePrimaryParameterShardMinLimitFlag != 0 {
		req.OpensearchSettingsShardIndexingPressurePrimaryParameterShard = &v3.OpensearchSettingsShardIndexingPressurePrimaryParameterShard{}
		req.OpensearchSettingsShardIndexingPressurePrimaryParameterShard.MinLimit = reqOpensearchSettingsShardIndexingPressurePrimaryParameterShardMinLimitFlag
	}
	if reqOpensearchSettingsShardIndexingPressurePrimaryParameterNodeSoftLimitFlag != 0 {
		req.OpensearchSettingsShardIndexingPressurePrimaryParameterNode = &v3.OpensearchSettingsShardIndexingPressurePrimaryParameterNode{}
		req.OpensearchSettingsShardIndexingPressurePrimaryParameterNode.SoftLimit = reqOpensearchSettingsShardIndexingPressurePrimaryParameterNodeSoftLimitFlag
	}
	if reqOpensearchSettingsShardIndexingPressureOperatingFactorUpperFlag != 0 {
		req.OpensearchSettingsShardIndexingPressureOperatingFacto = &v3.OpensearchSettingsShardIndexingPressureOperatingFactor{}
		req.OpensearchSettingsShardIndexingPressureOperatingFactor.Upper = reqOpensearchSettingsShardIndexingPressureOperatingFactorUpperFlag
	}
	if reqOpensearchSettingsShardIndexingPressureOperatingFactorOptimalFlag != 0 {

		req.OpensearchSettingsShardIndexingPressureOperatingFactor.Optimal = reqOpensearchSettingsShardIndexingPressureOperatingFactorOptimalFlag
	}
	if reqOpensearchSettingsShardIndexingPressureOperatingFactorLowerFlag != 0 {

		req.OpensearchSettingsShardIndexingPressureOperatingFactor.Lower = reqOpensearchSettingsShardIndexingPressureOperatingFactorLowerFlag
	}
	if reqOpensearchSettingsShardIndexingPressureEnforcedFlag != false {
		req.OpensearchSettingsShardIndexingPressu = &v3.OpensearchSettingsShardIndexingPressure{}
		req.OpensearchSettingsShardIndexingPressure.Enforced = reqOpensearchSettingsShardIndexingPressureEnforcedFlag
	}
	if reqOpensearchSettingsShardIndexingPressureEnabledFlag != false {

		req.OpensearchSettingsShardIndexingPressure.Enabled = reqOpensearchSettingsShardIndexingPressureEnabledFlag
	}
	if reqOpensearchSettingsSearchMaxBucketsFlag != 0 {
		req.OpensearchSetting = &v3.JSONSchemaOpensearch{}
		req.OpensearchSettings.SearchMaxBuckets = reqOpensearchSettingsSearchMaxBucketsFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskTotalHeapPercentThresholdFlag != 0 {
		req.OpensearchSettingsSearchBackpressureSearchTask = &v3.OpensearchSettingsSearchBackpressureSearchTask{}
		req.OpensearchSettingsSearchBackpressureSearchTask.TotalHeapPercentThreshold = reqOpensearchSettingsSearchBackpressureSearchTaskTotalHeapPercentThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskHeapVarianceFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.HeapVariance = reqOpensearchSettingsSearchBackpressureSearchTaskHeapVarianceFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskHeapPercentThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.HeapPercentThreshold = reqOpensearchSettingsSearchBackpressureSearchTaskHeapPercentThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskHeapMovingAverageWindowSizeFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.HeapMovingAverageWindowSize = reqOpensearchSettingsSearchBackpressureSearchTaskHeapMovingAverageWindowSizeFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskElapsedTimeMillisThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.ElapsedTimeMillisThreshold = reqOpensearchSettingsSearchBackpressureSearchTaskElapsedTimeMillisThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskCPUTimeMillisThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.CPUTimeMillisThreshold = reqOpensearchSettingsSearchBackpressureSearchTaskCPUTimeMillisThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRatioFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.CancellationRatio = reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRatioFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRateFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.CancellationRate = reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRateFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskCancellationBurstFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.CancellationBurst = reqOpensearchSettingsSearchBackpressureSearchTaskCancellationBurstFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskTotalHeapPercentThresholdFlag != 0 {
		req.OpensearchSettingsSearchBackpressureSearchShardTask = &v3.OpensearchSettingsSearchBackpressureSearchShardTask{}
		req.OpensearchSettingsSearchBackpressureSearchShardTask.TotalHeapPercentThreshold = reqOpensearchSettingsSearchBackpressureSearchShardTaskTotalHeapPercentThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapVarianceFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.HeapVariance = reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapVarianceFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapPercentThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.HeapPercentThreshold = reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapPercentThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapMovingAverageWindowSizeFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.HeapMovingAverageWindowSize = reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapMovingAverageWindowSizeFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskElapsedTimeMillisThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.ElapsedTimeMillisThreshold = reqOpensearchSettingsSearchBackpressureSearchShardTaskElapsedTimeMillisThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskCPUTimeMillisThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.CPUTimeMillisThreshold = reqOpensearchSettingsSearchBackpressureSearchShardTaskCPUTimeMillisThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRatioFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.CancellationRatio = reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRatioFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRateFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.CancellationRate = reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRateFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationBurstFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.CancellationBurst = reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationBurstFlag
	}
	if reqOpensearchSettingsSearchBackpressureNodeDuressNumSuccessiveBreachesFlag != 0 {
		req.OpensearchSettingsSearchBackpressureNodeD = &v3.OpensearchSettingsSearchBackpressureNodeDuress{}
		req.OpensearchSettingsSearchBackpressureNodeDuress.NumSuccessiveBreaches = reqOpensearchSettingsSearchBackpressureNodeDuressNumSuccessiveBreachesFlag
	}
	if reqOpensearchSettingsSearchBackpressureNodeDuressHeapThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureNodeDuress.HeapThreshold = reqOpensearchSettingsSearchBackpressureNodeDuressHeapThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureNodeDuressCPUThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureNodeDuress.CPUThreshold = reqOpensearchSettingsSearchBackpressureNodeDuressCPUThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureModeFlag != "" {
		req.OpensearchSettingsSearchBackpressur = &v3.OpensearchSettingsSearchBackpressure{}
		req.OpensearchSettingsSearchBackpressure.Mode = reqOpensearchSettingsSearchBackpressureModeFlag
	}
	if reqOpensearchSettingsScriptMaxCompilationsRateFlag != "" {
		req.OpensearchSetting = &v3.JSONSchemaOpensearch{}
		req.OpensearchSettings.ScriptMaxCompilationsRate = reqOpensearchSettingsScriptMaxCompilationsRateFlag
	}
	if reqOpensearchSettingsPluginsAlertingFilterByBackendRolesFlag != false {

		req.OpensearchSettings.PluginsAlertingFilterByBackendRoles = reqOpensearchSettingsPluginsAlertingFilterByBackendRolesFlag
	}
	if reqOpensearchSettingsOverrideMainResponseVersionFlag != false {

		req.OpensearchSettings.OverrideMainResponseVersion = reqOpensearchSettingsOverrideMainResponseVersionFlag
	}
	if reqOpensearchSettingsKnnMemoryCircuitBreakerLimitFlag != 0 {

		req.OpensearchSettings.KnnMemoryCircuitBreakerLimit = reqOpensearchSettingsKnnMemoryCircuitBreakerLimitFlag
	}
	if reqOpensearchSettingsKnnMemoryCircuitBreakerEnabledFlag != false {

		req.OpensearchSettings.KnnMemoryCircuitBreakerEnabled = reqOpensearchSettingsKnnMemoryCircuitBreakerEnabledFlag
	}
	if reqOpensearchSettingsIsmHistoryIsmHistoryRolloverRetentionPeriodFlag != 0 {
		req.OpensearchSetting = &v3.OpensearchSettingsIsmHistory{}
		req.OpensearchSettingsIsmHistory.IsmHistoryRolloverRetentionPeriod = reqOpensearchSettingsIsmHistoryIsmHistoryRolloverRetentionPeriodFlag
	}
	if reqOpensearchSettingsIsmHistoryIsmHistoryRolloverCheckPeriodFlag != 0 {

		req.OpensearchSettingsIsmHistory.IsmHistoryRolloverCheckPeriod = reqOpensearchSettingsIsmHistoryIsmHistoryRolloverCheckPeriodFlag
	}
	if reqOpensearchSettingsIsmHistoryIsmHistoryMaxDocsFlag != 0 {

		req.OpensearchSettingsIsmHistory.IsmHistoryMaxDocs = reqOpensearchSettingsIsmHistoryIsmHistoryMaxDocsFlag
	}
	if reqOpensearchSettingsIsmHistoryIsmHistoryMaxAgeFlag != 0 {

		req.OpensearchSettingsIsmHistory.IsmHistoryMaxAge = reqOpensearchSettingsIsmHistoryIsmHistoryMaxAgeFlag
	}
	if reqOpensearchSettingsIsmHistoryIsmHistoryEnabledFlag != false {

		req.OpensearchSettingsIsmHistory.IsmHistoryEnabled = reqOpensearchSettingsIsmHistoryIsmHistoryEnabledFlag
	}
	if reqOpensearchSettingsIsmHistoryIsmEnabledFlag != false {

		req.OpensearchSettingsIsmHistory.IsmEnabled = reqOpensearchSettingsIsmHistoryIsmEnabledFlag
	}
	if reqOpensearchSettingsIndicesRecoveryMaxConcurrentFileChunksFlag != 0 {
		req.OpensearchSetting = &v3.JSONSchemaOpensearch{}
		req.OpensearchSettings.IndicesRecoveryMaxConcurrentFileChunks = reqOpensearchSettingsIndicesRecoveryMaxConcurrentFileChunksFlag
	}
	if reqOpensearchSettingsIndicesRecoveryMaxBytesPerSecFlag != 0 {

		req.OpensearchSettings.IndicesRecoveryMaxBytesPerSec = reqOpensearchSettingsIndicesRecoveryMaxBytesPerSecFlag
	}
	if reqOpensearchSettingsIndicesQueryBoolMaxClauseCountFlag != 0 {

		req.OpensearchSettings.IndicesQueryBoolMaxClauseCount = reqOpensearchSettingsIndicesQueryBoolMaxClauseCountFlag
	}
	if reqOpensearchSettingsIndicesQueriesCacheSizeFlag != 0 {

		req.OpensearchSettings.IndicesQueriesCacheSize = reqOpensearchSettingsIndicesQueriesCacheSizeFlag
	}
	if reqOpensearchSettingsIndicesMemoryMinIndexBufferSizeFlag != 0 {

		req.OpensearchSettings.IndicesMemoryMinIndexBufferSize = reqOpensearchSettingsIndicesMemoryMinIndexBufferSizeFlag
	}
	if reqOpensearchSettingsIndicesMemoryMaxIndexBufferSizeFlag != 0 {

		req.OpensearchSettings.IndicesMemoryMaxIndexBufferSize = reqOpensearchSettingsIndicesMemoryMaxIndexBufferSizeFlag
	}
	if reqOpensearchSettingsIndicesMemoryIndexBufferSizeFlag != 0 {

		req.OpensearchSettings.IndicesMemoryIndexBufferSize = reqOpensearchSettingsIndicesMemoryIndexBufferSizeFlag
	}
	if reqOpensearchSettingsIndicesFielddataCacheSizeFlag != 0 {

		req.OpensearchSettings.IndicesFielddataCacheSize = reqOpensearchSettingsIndicesFielddataCacheSizeFlag
	}
	if reqOpensearchSettingsHTTPMaxInitialLineLengthFlag != 0 {

		req.OpensearchSettings.HTTPMaxInitialLineLength = reqOpensearchSettingsHTTPMaxInitialLineLengthFlag
	}
	if reqOpensearchSettingsHTTPMaxHeaderSizeFlag != 0 {

		req.OpensearchSettings.HTTPMaxHeaderSize = reqOpensearchSettingsHTTPMaxHeaderSizeFlag
	}
	if reqOpensearchSettingsHTTPMaxContentLengthFlag != 0 {

		req.OpensearchSettings.HTTPMaxContentLength = reqOpensearchSettingsHTTPMaxContentLengthFlag
	}
	if reqOpensearchSettingsEnableSecurityAuditFlag != false {

		req.OpensearchSettings.EnableSecurityAudit = reqOpensearchSettingsEnableSecurityAuditFlag
	}
	if reqOpensearchSettingsEmailSenderEmailSenderUsernameFlag != "" {
		req.OpensearchSetting = &v3.OpensearchSettingsEmailSender{}
		req.OpensearchSettingsEmailSender.EmailSenderUsername = reqOpensearchSettingsEmailSenderEmailSenderUsernameFlag
	}
	if reqOpensearchSettingsEmailSenderEmailSenderPasswordFlag != "" {

		req.OpensearchSettingsEmailSender.EmailSenderPassword = reqOpensearchSettingsEmailSenderEmailSenderPasswordFlag
	}
	if reqOpensearchSettingsEmailSenderEmailSenderNameFlag != "" {

		req.OpensearchSettingsEmailSender.EmailSenderName = reqOpensearchSettingsEmailSenderEmailSenderNameFlag
	}
	if reqOpensearchSettingsClusterRoutingAllocationNodeConcurrentRecoveriesFlag != 0 {
		req.OpensearchS = &v3.JSONSchemaOpensearch{}
		req.OpensearchSettings.ClusterRoutingAllocationNodeConcurrentRecoveries = reqOpensearchSettingsClusterRoutingAllocationNodeConcurrentRecoveriesFlag
	}
	if reqOpensearchSettingsClusterMaxShardsPerNodeFlag != 0 {

		req.OpensearchSettings.ClusterMaxShardsPerNode = reqOpensearchSettingsClusterMaxShardsPerNodeFlag
	}
	if reqOpensearchSettingsAuthFailureListenersIPRateLimitingTypeFlag != "" {
		req.OpensearchSettingsAuthFailureListenersIPRateLimiting = &v3.OpensearchSettingsAuthFailureListenersIPRateLimiting{}
		req.OpensearchSettingsAuthFailureListenersIPRateLimiting.Type = reqOpensearchSettingsAuthFailureListenersIPRateLimitingTypeFlag
	}
	if reqOpensearchSettingsAuthFailureListenersIPRateLimitingTimeWindowSecondsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersIPRateLimiting.TimeWindowSeconds = reqOpensearchSettingsAuthFailureListenersIPRateLimitingTimeWindowSecondsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxTrackedClientsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersIPRateLimiting.MaxTrackedClients = reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxTrackedClientsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxBlockedClientsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersIPRateLimiting.MaxBlockedClients = reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxBlockedClientsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersIPRateLimitingBlockExpirySecondsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersIPRateLimiting.BlockExpirySeconds = reqOpensearchSettingsAuthFailureListenersIPRateLimitingBlockExpirySecondsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersIPRateLimitingAllowedTriesFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersIPRateLimiting.AllowedTries = reqOpensearchSettingsAuthFailureListenersIPRateLimitingAllowedTriesFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTypeFlag != "" {
		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting = &v3.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting{}
		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.Type = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTypeFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTimeWindowSecondsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.TimeWindowSeconds = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTimeWindowSecondsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxTrackedClientsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.MaxTrackedClients = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxTrackedClientsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxBlockedClientsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.MaxBlockedClients = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxBlockedClientsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingBlockExpirySecondsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.BlockExpirySeconds = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingBlockExpirySecondsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAuthenticationBackendFlag != "" {

		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.AuthenticationBackend = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAuthenticationBackendFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAllowedTriesFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.AllowedTries = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAllowedTriesFlag
	}
	if reqOpensearchSettingsActionDestructiveRequiresNameFlag != false {
		req.OpensearchSetting = &v3.JSONSchemaOpensearch{}
		req.OpensearchSettings.ActionDestructiveRequiresName = reqOpensearchSettingsActionDestructiveRequiresNameFlag
	}
	if reqOpensearchSettingsActionAutoCreateIndexEnabledFlag != false {

		req.OpensearchSettings.ActionAutoCreateIndexEnabled = reqOpensearchSettingsActionAutoCreateIndexEnabledFlag
	}
	if reqOpensearchDashboardsOpensearchRequestTimeoutFlag != 0 {
		req.OpensearchDashboard = &v3.CreateDBAASServiceOpensearchRequestOpensearchDashboards{}
		req.OpensearchDashboards.OpensearchRequestTimeout = reqOpensearchDashboardsOpensearchRequestTimeoutFlag
	}
	if reqOpensearchDashboardsMaxOldSpaceSizeFlag != 0 {

		req.OpensearchDashboards.MaxOldSpaceSize = reqOpensearchDashboardsMaxOldSpaceSizeFlag
	}
	if reqOpensearchDashboardsEnabledFlag != false {

		req.OpensearchDashboards.Enabled = reqOpensearchDashboardsEnabledFlag
	}
	req.MaxIndexCount = reqMaxIndexCountFlag
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.CreateDBAASServiceOpensearchRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}
	req.KeepIndexRefreshInterval = reqKeepIndexRefreshIntervalFlag
	if reqIndexTemplateNumberOfShardsFlag != 0 {
		req.IndexTemplat = &v3.CreateDBAASServiceOpensearchRequestIndexTemplate{}
		req.IndexTemplate.NumberOfShards = reqIndexTemplateNumberOfShardsFlag
	}
	if reqIndexTemplateNumberOfReplicasFlag != 0 {

		req.IndexTemplate.NumberOfReplicas = reqIndexTemplateNumberOfReplicasFlag
	}
	if reqIndexTemplateMappingNestedObjectsLimitFlag != 0 {

		req.IndexTemplate.MappingNestedObjectsLimit = reqIndexTemplateMappingNestedObjectsLimitFlag
	}
	req.ForkFromService = reqForkFromServiceFlag

	resp, err := client.CreateDBAASServiceOpensearch(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASServiceOpensearchCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-service-opensearch", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqIndexTemplateMappingNestedObjectsLimitFlag int64
	flagset.Int64Var(&reqIndexTemplateMappingNestedObjectsLimitFlag, "index-template.mapping-nested-objects-limit", 0, "The maximum number of nested JSON objects that a single document can contain across all nested types. This limit helps to prevent out of memory errors when a document contains too many nested objects. Default is 10000.")
	var reqIndexTemplateNumberOfReplicasFlag int64
	flagset.Int64Var(&reqIndexTemplateNumberOfReplicasFlag, "index-template.number-of-replicas", 0, "The number of replicas each primary shard has.")
	var reqIndexTemplateNumberOfShardsFlag int64
	flagset.Int64Var(&reqIndexTemplateNumberOfShardsFlag, "index-template.number-of-shards", 0, "The number of primary shards that an index should have.")
	var reqKeepIndexRefreshIntervalFlag bool
	flagset.BoolVar(&reqKeepIndexRefreshIntervalFlag, "keep-index-refresh-interval", false, "Aiven automation resets index.refresh_interval to default value for every index to be sure that indices are always visible to search. If it doesn't fit your case, you can disable this by setting up this flag to true.")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqMaxIndexCountFlag int64
	flagset.Int64Var(&reqMaxIndexCountFlag, "max-index-count", 0, "Maximum number of indexes to keep before deleting the oldest one")
	var reqOpensearchDashboardsEnabledFlag bool
	flagset.BoolVar(&reqOpensearchDashboardsEnabledFlag, "opensearch-dashboards.enabled", false, "Enable or disable OpenSearch Dashboards (default: true)")
	var reqOpensearchDashboardsMaxOldSpaceSizeFlag int64
	flagset.Int64Var(&reqOpensearchDashboardsMaxOldSpaceSizeFlag, "opensearch-dashboards.max-old-space-size", 0, "Limits the maximum amount of memory (in MiB) the OpenSearch Dashboards process can use. This sets the max_old_space_size option of the nodejs running the OpenSearch Dashboards. Note: the memory reserved by OpenSearch Dashboards is not available for OpenSearch. (default: 128)")
	var reqOpensearchDashboardsOpensearchRequestTimeoutFlag int64
	flagset.Int64Var(&reqOpensearchDashboardsOpensearchRequestTimeoutFlag, "opensearch-dashboards.opensearch-request-timeout", 0, "Timeout in milliseconds for requests made by OpenSearch Dashboards towards OpenSearch (default: 30000)")
	var reqOpensearchSettingsActionAutoCreateIndexEnabledFlag bool
	flagset.BoolVar(&reqOpensearchSettingsActionAutoCreateIndexEnabledFlag, "opensearch-settings.action_auto_create_index_enabled", false, "Explicitly allow or block automatic creation of indices. Defaults to true")
	var reqOpensearchSettingsActionDestructiveRequiresNameFlag bool
	flagset.BoolVar(&reqOpensearchSettingsActionDestructiveRequiresNameFlag, "opensearch-settings.action_destructive_requires_name", false, "Require explicit index names when deleting")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAllowedTriesFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAllowedTriesFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.allowed_tries", 0, "The number of login attempts allowed before login is blocked")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAuthenticationBackendFlag string
	flagset.StringVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAuthenticationBackendFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.authentication_backend", "", "The internal backend. Enter `internal`")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingBlockExpirySecondsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingBlockExpirySecondsFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.block_expiry_seconds", 0, "The duration of time that login remains blocked after a failed login")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxBlockedClientsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxBlockedClientsFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.max_blocked_clients", 0, "The maximum number of blocked IP addresses")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxTrackedClientsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxTrackedClientsFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.max_tracked_clients", 0, "The maximum number of tracked IP addresses that have failed login")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTimeWindowSecondsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTimeWindowSecondsFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.time_window_seconds", 0, "The window of time in which the value for `allowed_tries` is enforced")
	var reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTypeFlag string
	flagset.StringVar(&reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTypeFlag, "opensearch-settings.auth_failure_listeners.internal_authentication_backend_limiting.type", "", "The type of rate limiting")
	var reqOpensearchSettingsAuthFailureListenersIPRateLimitingAllowedTriesFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersIPRateLimitingAllowedTriesFlag, "opensearch-settings.auth_failure_listeners.ip_rate_limiting.allowed_tries", 0, "The number of login attempts allowed before login is blocked")
	var reqOpensearchSettingsAuthFailureListenersIPRateLimitingBlockExpirySecondsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersIPRateLimitingBlockExpirySecondsFlag, "opensearch-settings.auth_failure_listeners.ip_rate_limiting.block_expiry_seconds", 0, "The duration of time that login remains blocked after a failed login")
	var reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxBlockedClientsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxBlockedClientsFlag, "opensearch-settings.auth_failure_listeners.ip_rate_limiting.max_blocked_clients", 0, "The maximum number of blocked IP addresses")
	var reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxTrackedClientsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxTrackedClientsFlag, "opensearch-settings.auth_failure_listeners.ip_rate_limiting.max_tracked_clients", 0, "The maximum number of tracked IP addresses that have failed login")
	var reqOpensearchSettingsAuthFailureListenersIPRateLimitingTimeWindowSecondsFlag int
	flagset.IntVar(&reqOpensearchSettingsAuthFailureListenersIPRateLimitingTimeWindowSecondsFlag, "opensearch-settings.auth_failure_listeners.ip_rate_limiting.time_window_seconds", 0, "The window of time in which the value for `allowed_tries` is enforced")
	var reqOpensearchSettingsAuthFailureListenersIPRateLimitingTypeFlag string
	flagset.StringVar(&reqOpensearchSettingsAuthFailureListenersIPRateLimitingTypeFlag, "opensearch-settings.auth_failure_listeners.ip_rate_limiting.type", "", "The type of rate limiting")
	var reqOpensearchSettingsClusterMaxShardsPerNodeFlag int
	flagset.IntVar(&reqOpensearchSettingsClusterMaxShardsPerNodeFlag, "opensearch-settings.cluster_max_shards_per_node", 0, "Controls the number of shards allowed in the cluster per data node")
	var reqOpensearchSettingsClusterRoutingAllocationNodeConcurrentRecoveriesFlag int
	flagset.IntVar(&reqOpensearchSettingsClusterRoutingAllocationNodeConcurrentRecoveriesFlag, "opensearch-settings.cluster_routing_allocation_node_concurrent_recoveries", 0, "How many concurrent incoming/outgoing shard recoveries (normally replicas) are allowed to happen on a node. Defaults to 2.")
	var reqOpensearchSettingsEmailSenderEmailSenderNameFlag string
	flagset.StringVar(&reqOpensearchSettingsEmailSenderEmailSenderNameFlag, "opensearch-settings.email-sender.email_sender_name", "", "This should be identical to the Sender name defined in Opensearch dashboards")
	var reqOpensearchSettingsEmailSenderEmailSenderPasswordFlag string
	flagset.StringVar(&reqOpensearchSettingsEmailSenderEmailSenderPasswordFlag, "opensearch-settings.email-sender.email_sender_password", "", "Sender password for Opensearch alerts to authenticate with SMTP server")
	var reqOpensearchSettingsEmailSenderEmailSenderUsernameFlag string
	flagset.StringVar(&reqOpensearchSettingsEmailSenderEmailSenderUsernameFlag, "opensearch-settings.email-sender.email_sender_username", "", "Sender username for Opensearch alerts")
	var reqOpensearchSettingsEnableSecurityAuditFlag bool
	flagset.BoolVar(&reqOpensearchSettingsEnableSecurityAuditFlag, "opensearch-settings.enable_security_audit", false, "Enable/Disable security audit")
	var reqOpensearchSettingsHTTPMaxContentLengthFlag int
	flagset.IntVar(&reqOpensearchSettingsHTTPMaxContentLengthFlag, "opensearch-settings.http_max_content_length", 0, "Maximum content length for HTTP requests to the OpenSearch HTTP API, in bytes.")
	var reqOpensearchSettingsHTTPMaxHeaderSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsHTTPMaxHeaderSizeFlag, "opensearch-settings.http_max_header_size", 0, "The max size of allowed headers, in bytes")
	var reqOpensearchSettingsHTTPMaxInitialLineLengthFlag int
	flagset.IntVar(&reqOpensearchSettingsHTTPMaxInitialLineLengthFlag, "opensearch-settings.http_max_initial_line_length", 0, "The max length of an HTTP URL, in bytes")
	var reqOpensearchSettingsIndicesFielddataCacheSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesFielddataCacheSizeFlag, "opensearch-settings.indices_fielddata_cache_size", 0, "Relative amount. Maximum amount of heap memory used for field data cache. This is an expert setting; decreasing the value too much will increase overhead of loading field data; too much memory used for field data cache will decrease amount of heap available for other operations.")
	var reqOpensearchSettingsIndicesMemoryIndexBufferSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesMemoryIndexBufferSizeFlag, "opensearch-settings.indices_memory_index_buffer_size", 0, "Percentage value. Default is 10%. Total amount of heap used for indexing buffer, before writing segments to disk. This is an expert setting. Too low value will slow down indexing; too high value will increase indexing performance but causes performance issues for query performance.")
	var reqOpensearchSettingsIndicesMemoryMaxIndexBufferSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesMemoryMaxIndexBufferSizeFlag, "opensearch-settings.indices_memory_max_index_buffer_size", 0, "Absolute value. Default is unbound. Doesn't work without indices.memory.index_buffer_size. Maximum amount of heap used for query cache, an absolute indices.memory.index_buffer_size maximum hard limit.")
	var reqOpensearchSettingsIndicesMemoryMinIndexBufferSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesMemoryMinIndexBufferSizeFlag, "opensearch-settings.indices_memory_min_index_buffer_size", 0, "Absolute value. Default is 48mb. Doesn't work without indices.memory.index_buffer_size. Minimum amount of heap used for query cache, an absolute indices.memory.index_buffer_size minimal hard limit.")
	var reqOpensearchSettingsIndicesQueriesCacheSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesQueriesCacheSizeFlag, "opensearch-settings.indices_queries_cache_size", 0, "Percentage value. Default is 10%. Maximum amount of heap used for query cache. This is an expert setting. Too low value will decrease query performance and increase performance for other operations; too high value will cause issues with other OpenSearch functionality.")
	var reqOpensearchSettingsIndicesQueryBoolMaxClauseCountFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesQueryBoolMaxClauseCountFlag, "opensearch-settings.indices_query_bool_max_clause_count", 0, "Maximum number of clauses Lucene BooleanQuery can have. The default value (1024) is relatively high, and increasing it may cause performance issues. Investigate other approaches first before increasing this value.")
	var reqOpensearchSettingsIndicesRecoveryMaxBytesPerSecFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesRecoveryMaxBytesPerSecFlag, "opensearch-settings.indices_recovery_max_bytes_per_sec", 0, "Limits total inbound and outbound recovery traffic for each node. Applies to both peer recoveries as well as snapshot recoveries (i.e., restores from a snapshot). Defaults to 40mb")
	var reqOpensearchSettingsIndicesRecoveryMaxConcurrentFileChunksFlag int
	flagset.IntVar(&reqOpensearchSettingsIndicesRecoveryMaxConcurrentFileChunksFlag, "opensearch-settings.indices_recovery_max_concurrent_file_chunks", 0, "Number of file chunks sent in parallel for each recovery. Defaults to 2.")
	var reqOpensearchSettingsIsmHistoryIsmEnabledFlag bool
	flagset.BoolVar(&reqOpensearchSettingsIsmHistoryIsmEnabledFlag, "opensearch-settings.ism-history.ism_enabled", false, "Specifies whether ISM is enabled or not")
	var reqOpensearchSettingsIsmHistoryIsmHistoryEnabledFlag bool
	flagset.BoolVar(&reqOpensearchSettingsIsmHistoryIsmHistoryEnabledFlag, "opensearch-settings.ism-history.ism_history_enabled", false, "Specifies whether audit history is enabled or not. The logs from ISM are automatically indexed to a logs document.")
	var reqOpensearchSettingsIsmHistoryIsmHistoryMaxAgeFlag int
	flagset.IntVar(&reqOpensearchSettingsIsmHistoryIsmHistoryMaxAgeFlag, "opensearch-settings.ism-history.ism_history_max_age", 0, "The maximum age before rolling over the audit history index in hours")
	var reqOpensearchSettingsIsmHistoryIsmHistoryMaxDocsFlag int
	flagset.IntVar(&reqOpensearchSettingsIsmHistoryIsmHistoryMaxDocsFlag, "opensearch-settings.ism-history.ism_history_max_docs", 0, "The maximum number of documents before rolling over the audit history index.")
	var reqOpensearchSettingsIsmHistoryIsmHistoryRolloverCheckPeriodFlag int
	flagset.IntVar(&reqOpensearchSettingsIsmHistoryIsmHistoryRolloverCheckPeriodFlag, "opensearch-settings.ism-history.ism_history_rollover_check_period", 0, "The time between rollover checks for the audit history index in hours.")
	var reqOpensearchSettingsIsmHistoryIsmHistoryRolloverRetentionPeriodFlag int
	flagset.IntVar(&reqOpensearchSettingsIsmHistoryIsmHistoryRolloverRetentionPeriodFlag, "opensearch-settings.ism-history.ism_history_rollover_retention_period", 0, "How long audit history indices are kept in days.")
	var reqOpensearchSettingsKnnMemoryCircuitBreakerEnabledFlag bool
	flagset.BoolVar(&reqOpensearchSettingsKnnMemoryCircuitBreakerEnabledFlag, "opensearch-settings.knn_memory_circuit_breaker_enabled", false, "Enable or disable KNN memory circuit breaker. Defaults to true.")
	var reqOpensearchSettingsKnnMemoryCircuitBreakerLimitFlag int
	flagset.IntVar(&reqOpensearchSettingsKnnMemoryCircuitBreakerLimitFlag, "opensearch-settings.knn_memory_circuit_breaker_limit", 0, "Maximum amount of memory that can be used for KNN index. Defaults to 50% of the JVM heap size.")
	var reqOpensearchSettingsOverrideMainResponseVersionFlag bool
	flagset.BoolVar(&reqOpensearchSettingsOverrideMainResponseVersionFlag, "opensearch-settings.override_main_response_version", false, "Compatibility mode sets OpenSearch to report its version as 7.10 so clients continue to work. Default is false")
	var reqOpensearchSettingsPluginsAlertingFilterByBackendRolesFlag bool
	flagset.BoolVar(&reqOpensearchSettingsPluginsAlertingFilterByBackendRolesFlag, "opensearch-settings.plugins_alerting_filter_by_backend_roles", false, "Enable or disable filtering of alerting by backend roles. Requires Security plugin. Defaults to false")
	var reqOpensearchSettingsScriptMaxCompilationsRateFlag string
	flagset.StringVar(&reqOpensearchSettingsScriptMaxCompilationsRateFlag, "opensearch-settings.script_max_compilations_rate", "", "Script compilation circuit breaker limits the number of inline script compilations within a period of time. Default is use-context")
	var reqOpensearchSettingsSearchBackpressureModeFlag string
	flagset.StringVar(&reqOpensearchSettingsSearchBackpressureModeFlag, "opensearch-settings.search_backpressure.mode", "", "The search backpressure mode. Valid values are monitor_only, enforced, or disabled. Default is monitor_only")
	var reqOpensearchSettingsSearchBackpressureNodeDuressCPUThresholdFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureNodeDuressCPUThresholdFlag, "opensearch-settings.search_backpressure.node_duress.cpu_threshold", 0, "The CPU usage threshold (as a percentage) required for a node to be considered to be under duress. Default is 0.9")
	var reqOpensearchSettingsSearchBackpressureNodeDuressHeapThresholdFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureNodeDuressHeapThresholdFlag, "opensearch-settings.search_backpressure.node_duress.heap_threshold", 0, "The heap usage threshold (as a percentage) required for a node to be considered to be under duress. Default is 0.7")
	var reqOpensearchSettingsSearchBackpressureNodeDuressNumSuccessiveBreachesFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureNodeDuressNumSuccessiveBreachesFlag, "opensearch-settings.search_backpressure.node_duress.num_successive_breaches", 0, "The number of successive limit breaches after which the node is considered to be under duress. Default is 3")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationBurstFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationBurstFlag, "opensearch-settings.search_backpressure.search_shard_task.cancellation_burst", 0, "The maximum number of search tasks to cancel in a single iteration of the observer thread. Default is 10.0")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRateFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRateFlag, "opensearch-settings.search_backpressure.search_shard_task.cancellation_rate", 0, "The maximum number of tasks to cancel per millisecond of elapsed time. Default is 0.003")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRatioFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRatioFlag, "opensearch-settings.search_backpressure.search_shard_task.cancellation_ratio", 0, "The maximum number of tasks to cancel, as a percentage of successful task completions. Default is 0.1")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskCPUTimeMillisThresholdFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureSearchShardTaskCPUTimeMillisThresholdFlag, "opensearch-settings.search_backpressure.search_shard_task.cpu_time_millis_threshold", 0, "The CPU usage threshold (in milliseconds) required for a single search shard task before it is considered for cancellation. Default is 15000")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskElapsedTimeMillisThresholdFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureSearchShardTaskElapsedTimeMillisThresholdFlag, "opensearch-settings.search_backpressure.search_shard_task.elapsed_time_millis_threshold", 0, "The elapsed time threshold (in milliseconds) required for a single search shard task before it is considered for cancellation. Default is 30000")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapMovingAverageWindowSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapMovingAverageWindowSizeFlag, "opensearch-settings.search_backpressure.search_shard_task.heap_moving_average_window_size", 0, "The number of previously completed search shard tasks to consider when calculating the rolling average of heap usage. Default is 100")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapPercentThresholdFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapPercentThresholdFlag, "opensearch-settings.search_backpressure.search_shard_task.heap_percent_threshold", 0, "The heap usage threshold (as a percentage) required for a single search shard task before it is considered for cancellation. Default is 0.5")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapVarianceFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapVarianceFlag, "opensearch-settings.search_backpressure.search_shard_task.heap_variance", 0, "The minimum variance required for a single search shard task’s heap usage compared to the rolling average of previously completed tasks before it is considered for cancellation. Default is 2.0")
	var reqOpensearchSettingsSearchBackpressureSearchShardTaskTotalHeapPercentThresholdFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchShardTaskTotalHeapPercentThresholdFlag, "opensearch-settings.search_backpressure.search_shard_task.total_heap_percent_threshold", 0, "The heap usage threshold (as a percentage) required for the sum of heap usages of all search shard tasks before cancellation is applied. Default is 0.5")
	var reqOpensearchSettingsSearchBackpressureSearchTaskCancellationBurstFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchTaskCancellationBurstFlag, "opensearch-settings.search_backpressure.search_task.cancellation_burst", 0, "The maximum number of search tasks to cancel in a single iteration of the observer thread. Default is 5.0")
	var reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRateFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRateFlag, "opensearch-settings.search_backpressure.search_task.cancellation_rate", 0, "The maximum number of search tasks to cancel per millisecond of elapsed time. Default is 0.003")
	var reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRatioFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRatioFlag, "opensearch-settings.search_backpressure.search_task.cancellation_ratio", 0, "The maximum number of search tasks to cancel, as a percentage of successful search task completions. Default is 0.1")
	var reqOpensearchSettingsSearchBackpressureSearchTaskCPUTimeMillisThresholdFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureSearchTaskCPUTimeMillisThresholdFlag, "opensearch-settings.search_backpressure.search_task.cpu_time_millis_threshold", 0, "The CPU usage threshold (in milliseconds) required for an individual parent task before it is considered for cancellation. Default is 30000")
	var reqOpensearchSettingsSearchBackpressureSearchTaskElapsedTimeMillisThresholdFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureSearchTaskElapsedTimeMillisThresholdFlag, "opensearch-settings.search_backpressure.search_task.elapsed_time_millis_threshold", 0, "The elapsed time threshold (in milliseconds) required for an individual parent task before it is considered for cancellation. Default is 45000")
	var reqOpensearchSettingsSearchBackpressureSearchTaskHeapMovingAverageWindowSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchBackpressureSearchTaskHeapMovingAverageWindowSizeFlag, "opensearch-settings.search_backpressure.search_task.heap_moving_average_window_size", 0, "The window size used to calculate the rolling average of the heap usage for the completed parent tasks. Default is 10")
	var reqOpensearchSettingsSearchBackpressureSearchTaskHeapPercentThresholdFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchTaskHeapPercentThresholdFlag, "opensearch-settings.search_backpressure.search_task.heap_percent_threshold", 0, "The heap usage threshold (as a percentage) required for an individual parent task before it is considered for cancellation. Default is 0.2")
	var reqOpensearchSettingsSearchBackpressureSearchTaskHeapVarianceFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchTaskHeapVarianceFlag, "opensearch-settings.search_backpressure.search_task.heap_variance", 0, "The heap usage variance required for an individual parent task before it is considered for cancellation. A task is considered for cancellation when taskHeapUsage is greater than or equal to heapUsageMovingAverage * variance. Default is 2.0")
	var reqOpensearchSettingsSearchBackpressureSearchTaskTotalHeapPercentThresholdFlag float64
	flagset.Float64Var(&reqOpensearchSettingsSearchBackpressureSearchTaskTotalHeapPercentThresholdFlag, "opensearch-settings.search_backpressure.search_task.total_heap_percent_threshold", 0, "The heap usage threshold (as a percentage) required for the sum of heap usages of all search tasks before cancellation is applied. Default is 0.5")
	var reqOpensearchSettingsSearchMaxBucketsFlag int
	flagset.IntVar(&reqOpensearchSettingsSearchMaxBucketsFlag, "opensearch-settings.search_max_buckets", 0, "Maximum number of aggregation buckets allowed in a single response. OpenSearch default value is used when this is not defined.")
	var reqOpensearchSettingsShardIndexingPressureEnabledFlag bool
	flagset.BoolVar(&reqOpensearchSettingsShardIndexingPressureEnabledFlag, "opensearch-settings.shard_indexing_pressure.enabled", false, "Enable or disable shard indexing backpressure. Default is false")
	var reqOpensearchSettingsShardIndexingPressureEnforcedFlag bool
	flagset.BoolVar(&reqOpensearchSettingsShardIndexingPressureEnforcedFlag, "opensearch-settings.shard_indexing_pressure.enforced", false, "Run shard indexing backpressure in shadow mode or enforced mode. In shadow mode (value set as false), shard indexing backpressure tracks all granular-level metrics, but it doesn’t actually reject any indexing requests. In enforced mode (value set as true), shard indexing backpressure rejects any requests to the cluster that might cause a dip in its performance. Default is false")
	var reqOpensearchSettingsShardIndexingPressureOperatingFactorLowerFlag float64
	flagset.Float64Var(&reqOpensearchSettingsShardIndexingPressureOperatingFactorLowerFlag, "opensearch-settings.shard_indexing_pressure.operating_factor.lower", 0, "Specify the lower occupancy limit of the allocated quota of memory for the shard. If the total memory usage of a shard is below this limit, shard indexing backpressure decreases the current allocated memory for that shard. Default is 0.75")
	var reqOpensearchSettingsShardIndexingPressureOperatingFactorOptimalFlag float64
	flagset.Float64Var(&reqOpensearchSettingsShardIndexingPressureOperatingFactorOptimalFlag, "opensearch-settings.shard_indexing_pressure.operating_factor.optimal", 0, "Specify the optimal occupancy of the allocated quota of memory for the shard. If the total memory usage of a shard is at this level, shard indexing backpressure doesn’t change the current allocated memory for that shard. Default is 0.85")
	var reqOpensearchSettingsShardIndexingPressureOperatingFactorUpperFlag float64
	flagset.Float64Var(&reqOpensearchSettingsShardIndexingPressureOperatingFactorUpperFlag, "opensearch-settings.shard_indexing_pressure.operating_factor.upper", 0, "Specify the upper occupancy limit of the allocated quota of memory for the shard. If the total memory usage of a shard is above this limit, shard indexing backpressure increases the current allocated memory for that shard. Default is 0.95")
	var reqOpensearchSettingsShardIndexingPressurePrimaryParameterNodeSoftLimitFlag float64
	flagset.Float64Var(&reqOpensearchSettingsShardIndexingPressurePrimaryParameterNodeSoftLimitFlag, "opensearch-settings.shard_indexing_pressure.primary_parameter.node.soft_limit", 0, "Define the percentage of the node-level memory threshold that acts as a soft indicator for strain on a node. Default is 0.7")
	var reqOpensearchSettingsShardIndexingPressurePrimaryParameterShardMinLimitFlag float64
	flagset.Float64Var(&reqOpensearchSettingsShardIndexingPressurePrimaryParameterShardMinLimitFlag, "opensearch-settings.shard_indexing_pressure.primary_parameter.shard.min_limit", 0, "Specify the minimum assigned quota for a new shard in any role (coordinator, primary, or replica). Shard indexing backpressure increases or decreases this allocated quota based on the inflow of traffic for the shard. Default is 0.001")
	var reqOpensearchSettingsThreadPoolAnalyzeQueueSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolAnalyzeQueueSizeFlag, "opensearch-settings.thread_pool_analyze_queue_size", 0, "Size for the thread pool queue. See documentation for exact details.")
	var reqOpensearchSettingsThreadPoolAnalyzeSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolAnalyzeSizeFlag, "opensearch-settings.thread_pool_analyze_size", 0, "Size for the thread pool. See documentation for exact details. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.")
	var reqOpensearchSettingsThreadPoolForceMergeSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolForceMergeSizeFlag, "opensearch-settings.thread_pool_force_merge_size", 0, "Size for the thread pool. See documentation for exact details. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.")
	var reqOpensearchSettingsThreadPoolGetQueueSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolGetQueueSizeFlag, "opensearch-settings.thread_pool_get_queue_size", 0, "Size for the thread pool queue. See documentation for exact details.")
	var reqOpensearchSettingsThreadPoolGetSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolGetSizeFlag, "opensearch-settings.thread_pool_get_size", 0, "Size for the thread pool. See documentation for exact details. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.")
	var reqOpensearchSettingsThreadPoolSearchQueueSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolSearchQueueSizeFlag, "opensearch-settings.thread_pool_search_queue_size", 0, "Size for the thread pool queue. See documentation for exact details.")
	var reqOpensearchSettingsThreadPoolSearchSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolSearchSizeFlag, "opensearch-settings.thread_pool_search_size", 0, "Size for the thread pool. See documentation for exact details. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.")
	var reqOpensearchSettingsThreadPoolSearchThrottledQueueSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolSearchThrottledQueueSizeFlag, "opensearch-settings.thread_pool_search_throttled_queue_size", 0, "Size for the thread pool queue. See documentation for exact details.")
	var reqOpensearchSettingsThreadPoolSearchThrottledSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolSearchThrottledSizeFlag, "opensearch-settings.thread_pool_search_throttled_size", 0, "Size for the thread pool. See documentation for exact details. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.")
	var reqOpensearchSettingsThreadPoolWriteQueueSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolWriteQueueSizeFlag, "opensearch-settings.thread_pool_write_queue_size", 0, "Size for the thread pool queue. See documentation for exact details.")
	var reqOpensearchSettingsThreadPoolWriteSizeFlag int
	flagset.IntVar(&reqOpensearchSettingsThreadPoolWriteSizeFlag, "opensearch-settings.thread_pool_write_size", 0, "Size for the thread pool. See documentation for exact details. Do note this may have maximum value depending on CPU count - value is automatically lowered if set to higher than maximum value.")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")
	var reqVersionFlag string
	flagset.StringVar(&reqVersionFlag, "version", "", "Version")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASServiceOpensearchRequest
	req.Version = reqVersionFlag
	req.TerminationProtection = reqTerminationProtectionFlag
	req.Plan = reqPlanFlag
	if reqOpensearchSettingsThreadPoolWriteSizeFlag != 0 {
		req.OpensearchSettings = &v3.JSONSchemaOpensearch{}
		req.OpensearchSettings.ThreadPoolWriteSize = reqOpensearchSettingsThreadPoolWriteSizeFlag
	}
	if reqOpensearchSettingsThreadPoolWriteQueueSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolWriteQueueSize = reqOpensearchSettingsThreadPoolWriteQueueSizeFlag
	}
	if reqOpensearchSettingsThreadPoolSearchThrottledSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolSearchThrottledSize = reqOpensearchSettingsThreadPoolSearchThrottledSizeFlag
	}
	if reqOpensearchSettingsThreadPoolSearchThrottledQueueSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolSearchThrottledQueueSize = reqOpensearchSettingsThreadPoolSearchThrottledQueueSizeFlag
	}
	if reqOpensearchSettingsThreadPoolSearchSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolSearchSize = reqOpensearchSettingsThreadPoolSearchSizeFlag
	}
	if reqOpensearchSettingsThreadPoolSearchQueueSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolSearchQueueSize = reqOpensearchSettingsThreadPoolSearchQueueSizeFlag
	}
	if reqOpensearchSettingsThreadPoolGetSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolGetSize = reqOpensearchSettingsThreadPoolGetSizeFlag
	}
	if reqOpensearchSettingsThreadPoolGetQueueSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolGetQueueSize = reqOpensearchSettingsThreadPoolGetQueueSizeFlag
	}
	if reqOpensearchSettingsThreadPoolForceMergeSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolForceMergeSize = reqOpensearchSettingsThreadPoolForceMergeSizeFlag
	}
	if reqOpensearchSettingsThreadPoolAnalyzeSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolAnalyzeSize = reqOpensearchSettingsThreadPoolAnalyzeSizeFlag
	}
	if reqOpensearchSettingsThreadPoolAnalyzeQueueSizeFlag != 0 {

		req.OpensearchSettings.ThreadPoolAnalyzeQueueSize = reqOpensearchSettingsThreadPoolAnalyzeQueueSizeFlag
	}
	if reqOpensearchSettingsShardIndexingPressurePrimaryParameterShardMinLimitFlag != 0 {
		req.OpensearchSettingsShardIndexingPressurePrimaryParameterShard = &v3.OpensearchSettingsShardIndexingPressurePrimaryParameterShard{}
		req.OpensearchSettingsShardIndexingPressurePrimaryParameterShard.MinLimit = reqOpensearchSettingsShardIndexingPressurePrimaryParameterShardMinLimitFlag
	}
	if reqOpensearchSettingsShardIndexingPressurePrimaryParameterNodeSoftLimitFlag != 0 {
		req.OpensearchSettingsShardIndexingPressurePrimaryParameterNode = &v3.OpensearchSettingsShardIndexingPressurePrimaryParameterNode{}
		req.OpensearchSettingsShardIndexingPressurePrimaryParameterNode.SoftLimit = reqOpensearchSettingsShardIndexingPressurePrimaryParameterNodeSoftLimitFlag
	}
	if reqOpensearchSettingsShardIndexingPressureOperatingFactorUpperFlag != 0 {
		req.OpensearchSettingsShardIndexingPressureOperatingFacto = &v3.OpensearchSettingsShardIndexingPressureOperatingFactor{}
		req.OpensearchSettingsShardIndexingPressureOperatingFactor.Upper = reqOpensearchSettingsShardIndexingPressureOperatingFactorUpperFlag
	}
	if reqOpensearchSettingsShardIndexingPressureOperatingFactorOptimalFlag != 0 {

		req.OpensearchSettingsShardIndexingPressureOperatingFactor.Optimal = reqOpensearchSettingsShardIndexingPressureOperatingFactorOptimalFlag
	}
	if reqOpensearchSettingsShardIndexingPressureOperatingFactorLowerFlag != 0 {

		req.OpensearchSettingsShardIndexingPressureOperatingFactor.Lower = reqOpensearchSettingsShardIndexingPressureOperatingFactorLowerFlag
	}
	if reqOpensearchSettingsShardIndexingPressureEnforcedFlag != false {
		req.OpensearchSettingsShardIndexingPressu = &v3.OpensearchSettingsShardIndexingPressure{}
		req.OpensearchSettingsShardIndexingPressure.Enforced = reqOpensearchSettingsShardIndexingPressureEnforcedFlag
	}
	if reqOpensearchSettingsShardIndexingPressureEnabledFlag != false {

		req.OpensearchSettingsShardIndexingPressure.Enabled = reqOpensearchSettingsShardIndexingPressureEnabledFlag
	}
	if reqOpensearchSettingsSearchMaxBucketsFlag != 0 {
		req.OpensearchSetting = &v3.JSONSchemaOpensearch{}
		req.OpensearchSettings.SearchMaxBuckets = reqOpensearchSettingsSearchMaxBucketsFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskTotalHeapPercentThresholdFlag != 0 {
		req.OpensearchSettingsSearchBackpressureSearchTask = &v3.OpensearchSettingsSearchBackpressureSearchTask{}
		req.OpensearchSettingsSearchBackpressureSearchTask.TotalHeapPercentThreshold = reqOpensearchSettingsSearchBackpressureSearchTaskTotalHeapPercentThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskHeapVarianceFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.HeapVariance = reqOpensearchSettingsSearchBackpressureSearchTaskHeapVarianceFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskHeapPercentThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.HeapPercentThreshold = reqOpensearchSettingsSearchBackpressureSearchTaskHeapPercentThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskHeapMovingAverageWindowSizeFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.HeapMovingAverageWindowSize = reqOpensearchSettingsSearchBackpressureSearchTaskHeapMovingAverageWindowSizeFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskElapsedTimeMillisThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.ElapsedTimeMillisThreshold = reqOpensearchSettingsSearchBackpressureSearchTaskElapsedTimeMillisThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskCPUTimeMillisThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.CPUTimeMillisThreshold = reqOpensearchSettingsSearchBackpressureSearchTaskCPUTimeMillisThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRatioFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.CancellationRatio = reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRatioFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRateFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.CancellationRate = reqOpensearchSettingsSearchBackpressureSearchTaskCancellationRateFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchTaskCancellationBurstFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchTask.CancellationBurst = reqOpensearchSettingsSearchBackpressureSearchTaskCancellationBurstFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskTotalHeapPercentThresholdFlag != 0 {
		req.OpensearchSettingsSearchBackpressureSearchShardTask = &v3.OpensearchSettingsSearchBackpressureSearchShardTask{}
		req.OpensearchSettingsSearchBackpressureSearchShardTask.TotalHeapPercentThreshold = reqOpensearchSettingsSearchBackpressureSearchShardTaskTotalHeapPercentThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapVarianceFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.HeapVariance = reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapVarianceFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapPercentThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.HeapPercentThreshold = reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapPercentThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapMovingAverageWindowSizeFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.HeapMovingAverageWindowSize = reqOpensearchSettingsSearchBackpressureSearchShardTaskHeapMovingAverageWindowSizeFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskElapsedTimeMillisThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.ElapsedTimeMillisThreshold = reqOpensearchSettingsSearchBackpressureSearchShardTaskElapsedTimeMillisThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskCPUTimeMillisThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.CPUTimeMillisThreshold = reqOpensearchSettingsSearchBackpressureSearchShardTaskCPUTimeMillisThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRatioFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.CancellationRatio = reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRatioFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRateFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.CancellationRate = reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationRateFlag
	}
	if reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationBurstFlag != 0 {

		req.OpensearchSettingsSearchBackpressureSearchShardTask.CancellationBurst = reqOpensearchSettingsSearchBackpressureSearchShardTaskCancellationBurstFlag
	}
	if reqOpensearchSettingsSearchBackpressureNodeDuressNumSuccessiveBreachesFlag != 0 {
		req.OpensearchSettingsSearchBackpressureNodeD = &v3.OpensearchSettingsSearchBackpressureNodeDuress{}
		req.OpensearchSettingsSearchBackpressureNodeDuress.NumSuccessiveBreaches = reqOpensearchSettingsSearchBackpressureNodeDuressNumSuccessiveBreachesFlag
	}
	if reqOpensearchSettingsSearchBackpressureNodeDuressHeapThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureNodeDuress.HeapThreshold = reqOpensearchSettingsSearchBackpressureNodeDuressHeapThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureNodeDuressCPUThresholdFlag != 0 {

		req.OpensearchSettingsSearchBackpressureNodeDuress.CPUThreshold = reqOpensearchSettingsSearchBackpressureNodeDuressCPUThresholdFlag
	}
	if reqOpensearchSettingsSearchBackpressureModeFlag != "" {
		req.OpensearchSettingsSearchBackpressur = &v3.OpensearchSettingsSearchBackpressure{}
		req.OpensearchSettingsSearchBackpressure.Mode = reqOpensearchSettingsSearchBackpressureModeFlag
	}
	if reqOpensearchSettingsScriptMaxCompilationsRateFlag != "" {
		req.OpensearchSetting = &v3.JSONSchemaOpensearch{}
		req.OpensearchSettings.ScriptMaxCompilationsRate = reqOpensearchSettingsScriptMaxCompilationsRateFlag
	}
	if reqOpensearchSettingsPluginsAlertingFilterByBackendRolesFlag != false {

		req.OpensearchSettings.PluginsAlertingFilterByBackendRoles = reqOpensearchSettingsPluginsAlertingFilterByBackendRolesFlag
	}
	if reqOpensearchSettingsOverrideMainResponseVersionFlag != false {

		req.OpensearchSettings.OverrideMainResponseVersion = reqOpensearchSettingsOverrideMainResponseVersionFlag
	}
	if reqOpensearchSettingsKnnMemoryCircuitBreakerLimitFlag != 0 {

		req.OpensearchSettings.KnnMemoryCircuitBreakerLimit = reqOpensearchSettingsKnnMemoryCircuitBreakerLimitFlag
	}
	if reqOpensearchSettingsKnnMemoryCircuitBreakerEnabledFlag != false {

		req.OpensearchSettings.KnnMemoryCircuitBreakerEnabled = reqOpensearchSettingsKnnMemoryCircuitBreakerEnabledFlag
	}
	if reqOpensearchSettingsIsmHistoryIsmHistoryRolloverRetentionPeriodFlag != 0 {
		req.OpensearchSetting = &v3.OpensearchSettingsIsmHistory{}
		req.OpensearchSettingsIsmHistory.IsmHistoryRolloverRetentionPeriod = reqOpensearchSettingsIsmHistoryIsmHistoryRolloverRetentionPeriodFlag
	}
	if reqOpensearchSettingsIsmHistoryIsmHistoryRolloverCheckPeriodFlag != 0 {

		req.OpensearchSettingsIsmHistory.IsmHistoryRolloverCheckPeriod = reqOpensearchSettingsIsmHistoryIsmHistoryRolloverCheckPeriodFlag
	}
	if reqOpensearchSettingsIsmHistoryIsmHistoryMaxDocsFlag != 0 {

		req.OpensearchSettingsIsmHistory.IsmHistoryMaxDocs = reqOpensearchSettingsIsmHistoryIsmHistoryMaxDocsFlag
	}
	if reqOpensearchSettingsIsmHistoryIsmHistoryMaxAgeFlag != 0 {

		req.OpensearchSettingsIsmHistory.IsmHistoryMaxAge = reqOpensearchSettingsIsmHistoryIsmHistoryMaxAgeFlag
	}
	if reqOpensearchSettingsIsmHistoryIsmHistoryEnabledFlag != false {

		req.OpensearchSettingsIsmHistory.IsmHistoryEnabled = reqOpensearchSettingsIsmHistoryIsmHistoryEnabledFlag
	}
	if reqOpensearchSettingsIsmHistoryIsmEnabledFlag != false {

		req.OpensearchSettingsIsmHistory.IsmEnabled = reqOpensearchSettingsIsmHistoryIsmEnabledFlag
	}
	if reqOpensearchSettingsIndicesRecoveryMaxConcurrentFileChunksFlag != 0 {
		req.OpensearchSetting = &v3.JSONSchemaOpensearch{}
		req.OpensearchSettings.IndicesRecoveryMaxConcurrentFileChunks = reqOpensearchSettingsIndicesRecoveryMaxConcurrentFileChunksFlag
	}
	if reqOpensearchSettingsIndicesRecoveryMaxBytesPerSecFlag != 0 {

		req.OpensearchSettings.IndicesRecoveryMaxBytesPerSec = reqOpensearchSettingsIndicesRecoveryMaxBytesPerSecFlag
	}
	if reqOpensearchSettingsIndicesQueryBoolMaxClauseCountFlag != 0 {

		req.OpensearchSettings.IndicesQueryBoolMaxClauseCount = reqOpensearchSettingsIndicesQueryBoolMaxClauseCountFlag
	}
	if reqOpensearchSettingsIndicesQueriesCacheSizeFlag != 0 {

		req.OpensearchSettings.IndicesQueriesCacheSize = reqOpensearchSettingsIndicesQueriesCacheSizeFlag
	}
	if reqOpensearchSettingsIndicesMemoryMinIndexBufferSizeFlag != 0 {

		req.OpensearchSettings.IndicesMemoryMinIndexBufferSize = reqOpensearchSettingsIndicesMemoryMinIndexBufferSizeFlag
	}
	if reqOpensearchSettingsIndicesMemoryMaxIndexBufferSizeFlag != 0 {

		req.OpensearchSettings.IndicesMemoryMaxIndexBufferSize = reqOpensearchSettingsIndicesMemoryMaxIndexBufferSizeFlag
	}
	if reqOpensearchSettingsIndicesMemoryIndexBufferSizeFlag != 0 {

		req.OpensearchSettings.IndicesMemoryIndexBufferSize = reqOpensearchSettingsIndicesMemoryIndexBufferSizeFlag
	}
	if reqOpensearchSettingsIndicesFielddataCacheSizeFlag != 0 {

		req.OpensearchSettings.IndicesFielddataCacheSize = reqOpensearchSettingsIndicesFielddataCacheSizeFlag
	}
	if reqOpensearchSettingsHTTPMaxInitialLineLengthFlag != 0 {

		req.OpensearchSettings.HTTPMaxInitialLineLength = reqOpensearchSettingsHTTPMaxInitialLineLengthFlag
	}
	if reqOpensearchSettingsHTTPMaxHeaderSizeFlag != 0 {

		req.OpensearchSettings.HTTPMaxHeaderSize = reqOpensearchSettingsHTTPMaxHeaderSizeFlag
	}
	if reqOpensearchSettingsHTTPMaxContentLengthFlag != 0 {

		req.OpensearchSettings.HTTPMaxContentLength = reqOpensearchSettingsHTTPMaxContentLengthFlag
	}
	if reqOpensearchSettingsEnableSecurityAuditFlag != false {

		req.OpensearchSettings.EnableSecurityAudit = reqOpensearchSettingsEnableSecurityAuditFlag
	}
	if reqOpensearchSettingsEmailSenderEmailSenderUsernameFlag != "" {
		req.OpensearchSetting = &v3.OpensearchSettingsEmailSender{}
		req.OpensearchSettingsEmailSender.EmailSenderUsername = reqOpensearchSettingsEmailSenderEmailSenderUsernameFlag
	}
	if reqOpensearchSettingsEmailSenderEmailSenderPasswordFlag != "" {

		req.OpensearchSettingsEmailSender.EmailSenderPassword = reqOpensearchSettingsEmailSenderEmailSenderPasswordFlag
	}
	if reqOpensearchSettingsEmailSenderEmailSenderNameFlag != "" {

		req.OpensearchSettingsEmailSender.EmailSenderName = reqOpensearchSettingsEmailSenderEmailSenderNameFlag
	}
	if reqOpensearchSettingsClusterRoutingAllocationNodeConcurrentRecoveriesFlag != 0 {
		req.OpensearchS = &v3.JSONSchemaOpensearch{}
		req.OpensearchSettings.ClusterRoutingAllocationNodeConcurrentRecoveries = reqOpensearchSettingsClusterRoutingAllocationNodeConcurrentRecoveriesFlag
	}
	if reqOpensearchSettingsClusterMaxShardsPerNodeFlag != 0 {

		req.OpensearchSettings.ClusterMaxShardsPerNode = reqOpensearchSettingsClusterMaxShardsPerNodeFlag
	}
	if reqOpensearchSettingsAuthFailureListenersIPRateLimitingTypeFlag != "" {
		req.OpensearchSettingsAuthFailureListenersIPRateLimiting = &v3.OpensearchSettingsAuthFailureListenersIPRateLimiting{}
		req.OpensearchSettingsAuthFailureListenersIPRateLimiting.Type = reqOpensearchSettingsAuthFailureListenersIPRateLimitingTypeFlag
	}
	if reqOpensearchSettingsAuthFailureListenersIPRateLimitingTimeWindowSecondsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersIPRateLimiting.TimeWindowSeconds = reqOpensearchSettingsAuthFailureListenersIPRateLimitingTimeWindowSecondsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxTrackedClientsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersIPRateLimiting.MaxTrackedClients = reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxTrackedClientsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxBlockedClientsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersIPRateLimiting.MaxBlockedClients = reqOpensearchSettingsAuthFailureListenersIPRateLimitingMaxBlockedClientsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersIPRateLimitingBlockExpirySecondsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersIPRateLimiting.BlockExpirySeconds = reqOpensearchSettingsAuthFailureListenersIPRateLimitingBlockExpirySecondsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersIPRateLimitingAllowedTriesFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersIPRateLimiting.AllowedTries = reqOpensearchSettingsAuthFailureListenersIPRateLimitingAllowedTriesFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTypeFlag != "" {
		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting = &v3.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting{}
		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.Type = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTypeFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTimeWindowSecondsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.TimeWindowSeconds = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingTimeWindowSecondsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxTrackedClientsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.MaxTrackedClients = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxTrackedClientsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxBlockedClientsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.MaxBlockedClients = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingMaxBlockedClientsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingBlockExpirySecondsFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.BlockExpirySeconds = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingBlockExpirySecondsFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAuthenticationBackendFlag != "" {

		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.AuthenticationBackend = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAuthenticationBackendFlag
	}
	if reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAllowedTriesFlag != 0 {

		req.OpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimiting.AllowedTries = reqOpensearchSettingsAuthFailureListenersInternalAuthenticationBackendLimitingAllowedTriesFlag
	}
	if reqOpensearchSettingsActionDestructiveRequiresNameFlag != false {
		req.OpensearchSetting = &v3.JSONSchemaOpensearch{}
		req.OpensearchSettings.ActionDestructiveRequiresName = reqOpensearchSettingsActionDestructiveRequiresNameFlag
	}
	if reqOpensearchSettingsActionAutoCreateIndexEnabledFlag != false {

		req.OpensearchSettings.ActionAutoCreateIndexEnabled = reqOpensearchSettingsActionAutoCreateIndexEnabledFlag
	}
	if reqOpensearchDashboardsOpensearchRequestTimeoutFlag != 0 {
		req.OpensearchDashboard = &v3.UpdateDBAASServiceOpensearchRequestOpensearchDashboards{}
		req.OpensearchDashboards.OpensearchRequestTimeout = reqOpensearchDashboardsOpensearchRequestTimeoutFlag
	}
	if reqOpensearchDashboardsMaxOldSpaceSizeFlag != 0 {

		req.OpensearchDashboards.MaxOldSpaceSize = reqOpensearchDashboardsMaxOldSpaceSizeFlag
	}
	if reqOpensearchDashboardsEnabledFlag != false {

		req.OpensearchDashboards.Enabled = reqOpensearchDashboardsEnabledFlag
	}
	req.MaxIndexCount = reqMaxIndexCountFlag
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.UpdateDBAASServiceOpensearchRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}
	req.KeepIndexRefreshInterval = reqKeepIndexRefreshIntervalFlag
	if reqIndexTemplateNumberOfShardsFlag != 0 {
		req.IndexTemplat = &v3.UpdateDBAASServiceOpensearchRequestIndexTemplate{}
		req.IndexTemplate.NumberOfShards = reqIndexTemplateNumberOfShardsFlag
	}
	if reqIndexTemplateNumberOfReplicasFlag != 0 {

		req.IndexTemplate.NumberOfReplicas = reqIndexTemplateNumberOfReplicasFlag
	}
	if reqIndexTemplateMappingNestedObjectsLimitFlag != 0 {

		req.IndexTemplate.MappingNestedObjectsLimit = reqIndexTemplateMappingNestedObjectsLimitFlag
	}

	resp, err := client.UpdateDBAASServiceOpensearch(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASOpensearchAclConfigCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-opensearch-acl-config", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASOpensearchAclConfig(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASOpensearchAclConfigCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-opensearch-acl-config", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqAclEnabledFlag bool
	flagset.BoolVar(&reqAclEnabledFlag, "acl-enabled", false, "Enable OpenSearch ACLs. When disabled authenticated service users have unrestricted access.")
	var reqExtendedAclEnabledFlag bool
	flagset.BoolVar(&reqExtendedAclEnabledFlag, "extended-acl-enabled", false, "Enable to enforce index rules in a limited fashion for requests that use the _mget, _msearch, and _bulk APIs")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASOpensearchAclConfigRequest
	req.ExtendedAclEnabled = reqExtendedAclEnabledFlag
	req.AclEnabled = reqAclEnabledFlag

	resp, err := client.UpdateDBAASOpensearchAclConfig(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StartDBAASOpensearchMaintenanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("start-dbaas-opensearch-maintenance", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StartDBAASOpensearchMaintenance(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASOpensearchUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-opensearch-user", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var reqUsernameFlag string
	flagset.StringVar(&reqUsernameFlag, "username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASOpensearchUserRequest
	req.Username = reqUsernameFlag

	resp, err := client.CreateDBAASOpensearchUser(context.Background(), serviceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASOpensearchUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-opensearch-user", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASOpensearchUser(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetDBAASOpensearchUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-dbaas-opensearch-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")
	var reqPasswordFlag string
	flagset.StringVar(&reqPasswordFlag, "password", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ResetDBAASOpensearchUserPasswordRequest
	req.Password = reqPasswordFlag

	resp, err := client.ResetDBAASOpensearchUserPassword(context.Background(), serviceNameFlag, usernameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RevealDBAASOpensearchUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reveal-dbaas-opensearch-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RevealDBAASOpensearchUserPassword(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASServicePGCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-service-pg", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASServicePG(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASServicePGCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-service-pg", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASServicePG(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASServicePGCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-service-pg", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqAdminPasswordFlag string
	flagset.StringVar(&reqAdminPasswordFlag, "admin-password", "", "Custom password for admin user. Defaults to random string. This must be set only when a new service is being created.")
	var reqAdminUsernameFlag string
	flagset.StringVar(&reqAdminUsernameFlag, "admin-username", "", "Custom username for admin user. This must be set only when a new service is being created.")
	var reqBackupScheduleBackupHourFlag int64
	flagset.Int64Var(&reqBackupScheduleBackupHourFlag, "backup-schedule.backup-hour", 0, "The hour of day (in UTC) when backup for the service is started. New backup is only started if previous backup has already completed.")
	var reqBackupScheduleBackupMinuteFlag int64
	flagset.Int64Var(&reqBackupScheduleBackupMinuteFlag, "backup-schedule.backup-minute", 0, "The minute of an hour when backup for the service is started. New backup is only started if previous backup has already completed.")
	var reqForkFromServiceFlag string
	flagset.StringVar(&reqForkFromServiceFlag, "fork-from-service", "", "")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqMigrationDbnameFlag string
	flagset.StringVar(&reqMigrationDbnameFlag, "migration.dbname", "", "Database name for bootstrapping the initial connection")
	var reqMigrationHostFlag string
	flagset.StringVar(&reqMigrationHostFlag, "migration.host", "", "Hostname or IP address of the server where to migrate data from")
	var reqMigrationIgnoreDbsFlag string
	flagset.StringVar(&reqMigrationIgnoreDbsFlag, "migration.ignore-dbs", "", "Comma-separated list of databases, which should be ignored during migration (supported by MySQL only at the moment)")
	var reqMigrationMethodFlag string
	flagset.StringVar(&reqMigrationMethodFlag, "migration.method", "", "")
	var reqMigrationPasswordFlag string
	flagset.StringVar(&reqMigrationPasswordFlag, "migration.password", "", "Password for authentication with the server where to migrate data from")
	var reqMigrationPortFlag int64
	flagset.Int64Var(&reqMigrationPortFlag, "migration.port", 0, "Port number of the server where to migrate data from")
	var reqMigrationSSLFlag bool
	flagset.BoolVar(&reqMigrationSSLFlag, "migration.ssl", false, "The server where to migrate data from is secured with SSL")
	var reqMigrationUsernameFlag string
	flagset.StringVar(&reqMigrationUsernameFlag, "migration.username", "", "User name for authentication with the server where to migrate data from")
	var reqPGSettingsAutovacuumAutovacuumAnalyzeScaleFactorFlag float64
	flagset.Float64Var(&reqPGSettingsAutovacuumAutovacuumAnalyzeScaleFactorFlag, "pg-settings.autovacuum.autovacuum_analyze_scale_factor", 0, "Specifies a fraction of the table size to add to autovacuum_analyze_threshold when deciding whether to trigger an ANALYZE. The default is 0.2 (20% of table size)")
	var reqPGSettingsAutovacuumAutovacuumAnalyzeThresholdFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumAnalyzeThresholdFlag, "pg-settings.autovacuum.autovacuum_analyze_threshold", 0, "Specifies the minimum number of inserted, updated or deleted tuples needed to trigger an ANALYZE in any one table. The default is 50 tuples.")
	var reqPGSettingsAutovacuumAutovacuumFreezeMaxAgeFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumFreezeMaxAgeFlag, "pg-settings.autovacuum.autovacuum_freeze_max_age", 0, "Specifies the maximum age (in transactions) that a table's pg_class.relfrozenxid field can attain before a VACUUM operation is forced to prevent transaction ID wraparound within the table. Note that the system will launch autovacuum processes to prevent wraparound even when autovacuum is otherwise disabled. This parameter will cause the server to be restarted.")
	var reqPGSettingsAutovacuumAutovacuumMaxWorkersFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumMaxWorkersFlag, "pg-settings.autovacuum.autovacuum_max_workers", 0, "Specifies the maximum number of autovacuum processes (other than the autovacuum launcher) that may be running at any one time. The default is three. This parameter can only be set at server start.")
	var reqPGSettingsAutovacuumAutovacuumNaptimeFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumNaptimeFlag, "pg-settings.autovacuum.autovacuum_naptime", 0, "Specifies the minimum delay between autovacuum runs on any given database. The delay is measured in seconds, and the default is one minute")
	var reqPGSettingsAutovacuumAutovacuumVacuumCostDelayFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumVacuumCostDelayFlag, "pg-settings.autovacuum.autovacuum_vacuum_cost_delay", 0, "Specifies the cost delay value that will be used in automatic VACUUM operations. If -1 is specified, the regular vacuum_cost_delay value will be used. The default value is 20 milliseconds")
	var reqPGSettingsAutovacuumAutovacuumVacuumCostLimitFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumVacuumCostLimitFlag, "pg-settings.autovacuum.autovacuum_vacuum_cost_limit", 0, "Specifies the cost limit value that will be used in automatic VACUUM operations. If -1 is specified (which is the default), the regular vacuum_cost_limit value will be used.")
	var reqPGSettingsAutovacuumAutovacuumVacuumScaleFactorFlag float64
	flagset.Float64Var(&reqPGSettingsAutovacuumAutovacuumVacuumScaleFactorFlag, "pg-settings.autovacuum.autovacuum_vacuum_scale_factor", 0, "Specifies a fraction of the table size to add to autovacuum_vacuum_threshold when deciding whether to trigger a VACUUM. The default is 0.2 (20% of table size)")
	var reqPGSettingsAutovacuumAutovacuumVacuumThresholdFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumVacuumThresholdFlag, "pg-settings.autovacuum.autovacuum_vacuum_threshold", 0, "Specifies the minimum number of updated or deleted tuples needed to trigger a VACUUM in any one table. The default is 50 tuples")
	var reqPGSettingsAutovacuumLogAutovacuumMinDurationFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumLogAutovacuumMinDurationFlag, "pg-settings.autovacuum.log_autovacuum_min_duration", 0, "Causes each action executed by autovacuum to be logged if it ran for at least the specified number of milliseconds. Setting this to zero logs all autovacuum actions. Minus-one (the default) disables logging autovacuum actions.")
	var reqPGSettingsBGWriterBgwriterDelayFlag int
	flagset.IntVar(&reqPGSettingsBGWriterBgwriterDelayFlag, "pg-settings.bg-writer.bgwriter_delay", 0, "Specifies the delay between activity rounds for the background writer in milliseconds. Default is 200.")
	var reqPGSettingsBGWriterBgwriterFlushAfterFlag int
	flagset.IntVar(&reqPGSettingsBGWriterBgwriterFlushAfterFlag, "pg-settings.bg-writer.bgwriter_flush_after", 0, "Whenever more than bgwriter_flush_after bytes have been written by the background writer, attempt to force the OS to issue these writes to the underlying storage. Specified in kilobytes, default is 512. Setting of 0 disables forced writeback.")
	var reqPGSettingsBGWriterBgwriterLruMaxpagesFlag int
	flagset.IntVar(&reqPGSettingsBGWriterBgwriterLruMaxpagesFlag, "pg-settings.bg-writer.bgwriter_lru_maxpages", 0, "In each round, no more than this many buffers will be written by the background writer. Setting this to zero disables background writing. Default is 100.")
	var reqPGSettingsBGWriterBgwriterLruMultiplierFlag float64
	flagset.Float64Var(&reqPGSettingsBGWriterBgwriterLruMultiplierFlag, "pg-settings.bg-writer.bgwriter_lru_multiplier", 0, "The average recent need for new buffers is multiplied by bgwriter_lru_multiplier to arrive at an estimate of the number that will be needed during the next round, (up to bgwriter_lru_maxpages). 1.0 represents a “just in time” policy of writing exactly the number of buffers predicted to be needed. Larger values provide some cushion against spikes in demand, while smaller values intentionally leave writes to be done by server processes. The default is 2.0.")
	var reqPGSettingsDeadlockTimeoutFlag int
	flagset.IntVar(&reqPGSettingsDeadlockTimeoutFlag, "pg-settings.deadlock_timeout", 0, "This is the amount of time, in milliseconds, to wait on a lock before checking to see if there is a deadlock condition.")
	var reqPGSettingsDefaultToastCompressionFlag string
	flagset.StringVar(&reqPGSettingsDefaultToastCompressionFlag, "pg-settings.default_toast_compression", "", "Specifies the default TOAST compression method for values of compressible columns (the default is lz4).")
	var reqPGSettingsIdleInTransactionSessionTimeoutFlag int
	flagset.IntVar(&reqPGSettingsIdleInTransactionSessionTimeoutFlag, "pg-settings.idle_in_transaction_session_timeout", 0, "Time out sessions with open transactions after this number of milliseconds")
	var reqPGSettingsJitFlag bool
	flagset.BoolVar(&reqPGSettingsJitFlag, "pg-settings.jit", false, "Controls system-wide use of Just-in-Time Compilation (JIT).")
	var reqPGSettingsLogErrorVerbosityFlag string
	flagset.StringVar(&reqPGSettingsLogErrorVerbosityFlag, "pg-settings.log_error_verbosity", "", "Controls the amount of detail written in the server log for each message that is logged.")
	var reqPGSettingsLogLinePrefixFlag string
	flagset.StringVar(&reqPGSettingsLogLinePrefixFlag, "pg-settings.log_line_prefix", "", "Choose from one of the available log-formats. These can support popular log analyzers like pgbadger, pganalyze etc.")
	var reqPGSettingsLogMinDurationStatementFlag int
	flagset.IntVar(&reqPGSettingsLogMinDurationStatementFlag, "pg-settings.log_min_duration_statement", 0, "Log statements that take more than this number of milliseconds to run, -1 disables")
	var reqPGSettingsLogTempFilesFlag int
	flagset.IntVar(&reqPGSettingsLogTempFilesFlag, "pg-settings.log_temp_files", 0, "Log statements for each temporary file created larger than this number of kilobytes, -1 disables")
	var reqPGSettingsMaxFilesPerProcessFlag int
	flagset.IntVar(&reqPGSettingsMaxFilesPerProcessFlag, "pg-settings.max_files_per_process", 0, "PostgreSQL maximum number of files that can be open per process")
	var reqPGSettingsMaxLocksPerTransactionFlag int
	flagset.IntVar(&reqPGSettingsMaxLocksPerTransactionFlag, "pg-settings.max_locks_per_transaction", 0, "PostgreSQL maximum locks per transaction")
	var reqPGSettingsMaxLogicalReplicationWorkersFlag int
	flagset.IntVar(&reqPGSettingsMaxLogicalReplicationWorkersFlag, "pg-settings.max_logical_replication_workers", 0, "PostgreSQL maximum logical replication workers (taken from the pool of max_parallel_workers)")
	var reqPGSettingsMaxParallelWorkersFlag int
	flagset.IntVar(&reqPGSettingsMaxParallelWorkersFlag, "pg-settings.max_parallel_workers", 0, "Sets the maximum number of workers that the system can support for parallel queries")
	var reqPGSettingsMaxParallelWorkersPerGatherFlag int
	flagset.IntVar(&reqPGSettingsMaxParallelWorkersPerGatherFlag, "pg-settings.max_parallel_workers_per_gather", 0, "Sets the maximum number of workers that can be started by a single Gather or Gather Merge node")
	var reqPGSettingsMaxPredLocksPerTransactionFlag int
	flagset.IntVar(&reqPGSettingsMaxPredLocksPerTransactionFlag, "pg-settings.max_pred_locks_per_transaction", 0, "PostgreSQL maximum predicate locks per transaction")
	var reqPGSettingsMaxPreparedTransactionsFlag int
	flagset.IntVar(&reqPGSettingsMaxPreparedTransactionsFlag, "pg-settings.max_prepared_transactions", 0, "PostgreSQL maximum prepared transactions")
	var reqPGSettingsMaxReplicationSlotsFlag int
	flagset.IntVar(&reqPGSettingsMaxReplicationSlotsFlag, "pg-settings.max_replication_slots", 0, "PostgreSQL maximum replication slots")
	var reqPGSettingsMaxStackDepthFlag int
	flagset.IntVar(&reqPGSettingsMaxStackDepthFlag, "pg-settings.max_stack_depth", 0, "Maximum depth of the stack in bytes")
	var reqPGSettingsMaxStandbyArchiveDelayFlag int
	flagset.IntVar(&reqPGSettingsMaxStandbyArchiveDelayFlag, "pg-settings.max_standby_archive_delay", 0, "Max standby archive delay in milliseconds")
	var reqPGSettingsMaxStandbyStreamingDelayFlag int
	flagset.IntVar(&reqPGSettingsMaxStandbyStreamingDelayFlag, "pg-settings.max_standby_streaming_delay", 0, "Max standby streaming delay in milliseconds")
	var reqPGSettingsMaxWorkerProcessesFlag int
	flagset.IntVar(&reqPGSettingsMaxWorkerProcessesFlag, "pg-settings.max_worker_processes", 0, "Sets the maximum number of background processes that the system can support")
	var reqPGSettingsPGPartmanBgwIntervalFlag int
	flagset.IntVar(&reqPGSettingsPGPartmanBgwIntervalFlag, "pg-settings.pg_partman_bgw.interval", 0, "Sets the time interval to run pg_partman's scheduled tasks")
	var reqPGSettingsPGPartmanBgwRoleFlag string
	flagset.StringVar(&reqPGSettingsPGPartmanBgwRoleFlag, "pg-settings.pg_partman_bgw.role", "", "Controls which role to use for pg_partman's scheduled background tasks.")
	var reqPGSettingsPGStatMonitorPgsmEnableQueryPlanFlag bool
	flagset.BoolVar(&reqPGSettingsPGStatMonitorPgsmEnableQueryPlanFlag, "pg-settings.pg_stat_monitor.pgsm_enable_query_plan", false, "Enables or disables query plan monitoring")
	var reqPGSettingsPGStatMonitorPgsmMaxBucketsFlag int
	flagset.IntVar(&reqPGSettingsPGStatMonitorPgsmMaxBucketsFlag, "pg-settings.pg_stat_monitor.pgsm_max_buckets", 0, "Sets the maximum number of buckets ")
	var reqPGSettingsPGStatStatementsTrackFlag string
	flagset.StringVar(&reqPGSettingsPGStatStatementsTrackFlag, "pg-settings.pg_stat_statements.track", "", "Controls which statements are counted. Specify top to track top-level statements (those issued directly by clients), all to also track nested statements (such as statements invoked within functions), or none to disable statement statistics collection. The default value is top.")
	var reqPGSettingsTempFileLimitFlag int
	flagset.IntVar(&reqPGSettingsTempFileLimitFlag, "pg-settings.temp_file_limit", 0, "PostgreSQL temporary file limit in KiB, -1 for unlimited")
	var reqPGSettingsTimezoneFlag string
	flagset.StringVar(&reqPGSettingsTimezoneFlag, "pg-settings.timezone", "", "PostgreSQL service timezone")
	var reqPGSettingsTrackActivityQuerySizeFlag int
	flagset.IntVar(&reqPGSettingsTrackActivityQuerySizeFlag, "pg-settings.track_activity_query_size", 0, "Specifies the number of bytes reserved to track the currently executing command for each active session.")
	var reqPGSettingsTrackCommitTimestampFlag string
	flagset.StringVar(&reqPGSettingsTrackCommitTimestampFlag, "pg-settings.track_commit_timestamp", "", "Record commit time of transactions.")
	var reqPGSettingsTrackFunctionsFlag string
	flagset.StringVar(&reqPGSettingsTrackFunctionsFlag, "pg-settings.track_functions", "", "Enables tracking of function call counts and time used.")
	var reqPGSettingsTrackIoTimingFlag string
	flagset.StringVar(&reqPGSettingsTrackIoTimingFlag, "pg-settings.track_io_timing", "", "Enables timing of database I/O calls. This parameter is off by default, because it will repeatedly query the operating system for the current time, which may cause significant overhead on some platforms.")
	var reqPGSettingsWalMaxSlotWalKeepSizeFlag int
	flagset.IntVar(&reqPGSettingsWalMaxSlotWalKeepSizeFlag, "pg-settings.wal.max_slot_wal_keep_size", 0, "PostgreSQL maximum WAL size (MB) reserved for replication slots. Default is -1 (unlimited). wal_keep_size minimum WAL size setting takes precedence over this.")
	var reqPGSettingsWalMaxWalSendersFlag int
	flagset.IntVar(&reqPGSettingsWalMaxWalSendersFlag, "pg-settings.wal.max_wal_senders", 0, "PostgreSQL maximum WAL senders")
	var reqPGSettingsWalWalSenderTimeoutFlag int
	flagset.IntVar(&reqPGSettingsWalWalSenderTimeoutFlag, "pg-settings.wal.wal_sender_timeout", 0, "Terminate replication connections that are inactive for longer than this amount of time, in milliseconds.")
	var reqPGSettingsWalWalWriterDelayFlag int
	flagset.IntVar(&reqPGSettingsWalWalWriterDelayFlag, "pg-settings.wal.wal_writer_delay", 0, "WAL flush interval in milliseconds. Note that setting this value to lower than the default 200ms may negatively impact performance")
	var reqPgbouncerSettingsAutodbIdleTimeoutFlag int
	flagset.IntVar(&reqPgbouncerSettingsAutodbIdleTimeoutFlag, "pgbouncer-settings.autodb_idle_timeout", 0, "If the automatically created database pools have been unused this many seconds, they are freed. If 0 then timeout is disabled. [seconds]")
	var reqPgbouncerSettingsAutodbMaxDBConnectionsFlag int
	flagset.IntVar(&reqPgbouncerSettingsAutodbMaxDBConnectionsFlag, "pgbouncer-settings.autodb_max_db_connections", 0, "Do not allow more than this many server connections per database (regardless of user). Setting it to 0 means unlimited.")
	var reqPgbouncerSettingsAutodbPoolModeFlag string
	flagset.StringVar(&reqPgbouncerSettingsAutodbPoolModeFlag, "pgbouncer-settings.autodb_pool_mode", "", "PGBouncer pool mode")
	var reqPgbouncerSettingsAutodbPoolSizeFlag int
	flagset.IntVar(&reqPgbouncerSettingsAutodbPoolSizeFlag, "pgbouncer-settings.autodb_pool_size", 0, "If non-zero then create automatically a pool of that size per user when a pool doesn't exist.")
	var reqPgbouncerSettingsMaxPreparedStatementsFlag int
	flagset.IntVar(&reqPgbouncerSettingsMaxPreparedStatementsFlag, "pgbouncer-settings.max_prepared_statements", 0, "PgBouncer tracks protocol-level named prepared statements related commands sent by the client in transaction and statement pooling modes when max_prepared_statements is set to a non-zero value. Setting it to 0 disables prepared statements. max_prepared_statements defaults to 100, and its maximum is 3000.")
	var reqPgbouncerSettingsMinPoolSizeFlag int
	flagset.IntVar(&reqPgbouncerSettingsMinPoolSizeFlag, "pgbouncer-settings.min_pool_size", 0, "Add more server connections to pool if below this number. Improves behavior when usual load comes suddenly back after period of total inactivity. The value is effectively capped at the pool size.")
	var reqPgbouncerSettingsServerIdleTimeoutFlag int
	flagset.IntVar(&reqPgbouncerSettingsServerIdleTimeoutFlag, "pgbouncer-settings.server_idle_timeout", 0, "If a server connection has been idle more than this many seconds it will be dropped. If 0 then timeout is disabled. [seconds]")
	var reqPgbouncerSettingsServerLifetimeFlag int
	flagset.IntVar(&reqPgbouncerSettingsServerLifetimeFlag, "pgbouncer-settings.server_lifetime", 0, "The pooler will close an unused server connection that has been connected longer than this. [seconds]")
	var reqPgbouncerSettingsServerResetQueryAlwaysFlag bool
	flagset.BoolVar(&reqPgbouncerSettingsServerResetQueryAlwaysFlag, "pgbouncer-settings.server_reset_query_always", false, "Run server_reset_query (DISCARD ALL) in all pooling modes")
	var reqPglookoutSettingsMaxFailoverReplicationTimeLagFlag int
	flagset.IntVar(&reqPglookoutSettingsMaxFailoverReplicationTimeLagFlag, "pglookout-settings.max_failover_replication_time_lag", 0, "Number of seconds of master unavailability before triggering database failover to standby")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqRecoveryBackupTimeFlag string
	flagset.StringVar(&reqRecoveryBackupTimeFlag, "recovery-backup-time", "", "ISO time of a backup to recover from for services that support arbitrary times")
	var reqSharedBuffersPercentageFlag int64
	flagset.Int64Var(&reqSharedBuffersPercentageFlag, "shared-buffers-percentage", 0, "Percentage of total RAM that the database server uses for shared memory buffers. Valid range is 20-60 (float), which corresponds to 20% - 60%. This setting adjusts the shared_buffers configuration value.")
	var reqSynchronousReplicationFlag string
	flagset.StringVar(&reqSynchronousReplicationFlag, "synchronous-replication", "", "")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")
	var reqTimescaledbSettingsMaxBackgroundWorkersFlag int
	flagset.IntVar(&reqTimescaledbSettingsMaxBackgroundWorkersFlag, "timescaledb-settings.max_background_workers", 0, "The number of background workers for timescaledb operations. You should configure this setting to the sum of your number of databases and the total number of concurrent background workers you want running at any given point in time.")
	var reqVariantFlag string
	flagset.StringVar(&reqVariantFlag, "variant", "", "")
	var reqVersionFlag string
	flagset.StringVar(&reqVersionFlag, "version", "", "")
	var reqWorkMemFlag int64
	flagset.Int64Var(&reqWorkMemFlag, "work-mem", 0, "Sets the maximum amount of memory to be used by a query operation (such as a sort or hash table) before writing to temporary disk files, in MB. Default is 1MB + 0.075% of total RAM (up to 32MB).")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASServicePGRequest
	req.WorkMem = reqWorkMemFlag
	req.Version = reqVersionFlag
	req.Variant = reqVariantFlag
	if reqTimescaledbSettingsMaxBackgroundWorkersFlag != 0 {
		req.TimescaledbSetti = &v3.JSONSchemaTimescaledb{}
		req.TimescaledbSettings.MaxBackgroundWorkers = reqTimescaledbSettingsMaxBackgroundWorkersFlag
	}
	req.TerminationProtection = reqTerminationProtectionFlag
	req.SynchronousReplication = reqSynchronousReplicationFlag
	req.SharedBuffersPercentage = reqSharedBuffersPercentageFlag
	req.RecoveryBackupTime = reqRecoveryBackupTimeFlag
	req.Plan = reqPlanFlag
	if reqPglookoutSettingsMaxFailoverReplicationTimeLagFlag != 0 {
		req.PglookoutSettings = &v3.JSONSchemaPglookout{}
		req.PglookoutSettings.MaxFailoverReplicationTimeLag = reqPglookoutSettingsMaxFailoverReplicationTimeLagFlag
	}
	if reqPgbouncerSettingsServerResetQueryAlwaysFlag != false {
		req.PgbouncerSetting = &v3.JSONSchemaPgbouncer{}
		req.PgbouncerSettings.ServerResetQueryAlways = reqPgbouncerSettingsServerResetQueryAlwaysFlag
	}
	if reqPgbouncerSettingsServerLifetimeFlag != 0 {

		req.PgbouncerSettings.ServerLifetime = reqPgbouncerSettingsServerLifetimeFlag
	}
	if reqPgbouncerSettingsServerIdleTimeoutFlag != 0 {

		req.PgbouncerSettings.ServerIdleTimeout = reqPgbouncerSettingsServerIdleTimeoutFlag
	}
	if reqPgbouncerSettingsMinPoolSizeFlag != 0 {

		req.PgbouncerSettings.MinPoolSize = reqPgbouncerSettingsMinPoolSizeFlag
	}
	if reqPgbouncerSettingsMaxPreparedStatementsFlag != 0 {

		req.PgbouncerSettings.MaxPreparedStatements = reqPgbouncerSettingsMaxPreparedStatementsFlag
	}
	if reqPgbouncerSettingsAutodbPoolSizeFlag != 0 {

		req.PgbouncerSettings.AutodbPoolSize = reqPgbouncerSettingsAutodbPoolSizeFlag
	}
	if reqPgbouncerSettingsAutodbPoolModeFlag != "" {

		req.PgbouncerSettings.AutodbPoolMode = reqPgbouncerSettingsAutodbPoolModeFlag
	}
	if reqPgbouncerSettingsAutodbMaxDBConnectionsFlag != 0 {

		req.PgbouncerSettings.AutodbMaxDBConnections = reqPgbouncerSettingsAutodbMaxDBConnectionsFlag
	}
	if reqPgbouncerSettingsAutodbIdleTimeoutFlag != 0 {

		req.PgbouncerSettings.AutodbIdleTimeout = reqPgbouncerSettingsAutodbIdleTimeoutFlag
	}
	if reqPGSettingsWalWalWriterDelayFlag != 0 {
		req.PGSettings = &v3.PGSettingsWal{}
		req.PGSettingsWal.WalWriterDelay = reqPGSettingsWalWalWriterDelayFlag
	}
	if reqPGSettingsWalWalSenderTimeoutFlag != 0 {

		req.PGSettingsWal.WalSenderTimeout = reqPGSettingsWalWalSenderTimeoutFlag
	}
	if reqPGSettingsWalMaxWalSendersFlag != 0 {

		req.PGSettingsWal.MaxWalSenders = reqPGSettingsWalMaxWalSendersFlag
	}
	if reqPGSettingsWalMaxSlotWalKeepSizeFlag != 0 {

		req.PGSettingsWal.MaxSlotWalKeepSize = reqPGSettingsWalMaxSlotWalKeepSizeFlag
	}
	if reqPGSettingsTrackIoTimingFlag != "" {
		req.PGSettings = &v3.JSONSchemaPG{}
		req.PGSettings.TrackIoTiming = reqPGSettingsTrackIoTimingFlag
	}
	if reqPGSettingsTrackFunctionsFlag != "" {

		req.PGSettings.TrackFunctions = reqPGSettingsTrackFunctionsFlag
	}
	if reqPGSettingsTrackCommitTimestampFlag != "" {

		req.PGSettings.TrackCommitTimestamp = reqPGSettingsTrackCommitTimestampFlag
	}
	if reqPGSettingsTrackActivityQuerySizeFlag != 0 {

		req.PGSettings.TrackActivityQuerySize = reqPGSettingsTrackActivityQuerySizeFlag
	}
	if reqPGSettingsTimezoneFlag != "" {

		req.PGSettings.Timezone = reqPGSettingsTimezoneFlag
	}
	if reqPGSettingsTempFileLimitFlag != 0 {

		req.PGSettings.TempFileLimit = reqPGSettingsTempFileLimitFlag
	}
	if reqPGSettingsPGStatStatementsTrackFlag != "" {

		req.PGSettings.PGStatStatementsTrack = reqPGSettingsPGStatStatementsTrackFlag
	}
	if reqPGSettingsPGStatMonitorPgsmMaxBucketsFlag != 0 {

		req.PGSettings.PGStatMonitorPgsmMaxBuckets = reqPGSettingsPGStatMonitorPgsmMaxBucketsFlag
	}
	if reqPGSettingsPGStatMonitorPgsmEnableQueryPlanFlag != false {

		req.PGSettings.PGStatMonitorPgsmEnableQueryPlan = reqPGSettingsPGStatMonitorPgsmEnableQueryPlanFlag
	}
	if reqPGSettingsPGPartmanBgwRoleFlag != "" {

		req.PGSettings.PGPartmanBgwRole = reqPGSettingsPGPartmanBgwRoleFlag
	}
	if reqPGSettingsPGPartmanBgwIntervalFlag != 0 {

		req.PGSettings.PGPartmanBgwInterval = reqPGSettingsPGPartmanBgwIntervalFlag
	}
	if reqPGSettingsMaxWorkerProcessesFlag != 0 {

		req.PGSettings.MaxWorkerProcesses = reqPGSettingsMaxWorkerProcessesFlag
	}
	if reqPGSettingsMaxStandbyStreamingDelayFlag != 0 {

		req.PGSettings.MaxStandbyStreamingDelay = reqPGSettingsMaxStandbyStreamingDelayFlag
	}
	if reqPGSettingsMaxStandbyArchiveDelayFlag != 0 {

		req.PGSettings.MaxStandbyArchiveDelay = reqPGSettingsMaxStandbyArchiveDelayFlag
	}
	if reqPGSettingsMaxStackDepthFlag != 0 {

		req.PGSettings.MaxStackDepth = reqPGSettingsMaxStackDepthFlag
	}
	if reqPGSettingsMaxReplicationSlotsFlag != 0 {

		req.PGSettings.MaxReplicationSlots = reqPGSettingsMaxReplicationSlotsFlag
	}
	if reqPGSettingsMaxPreparedTransactionsFlag != 0 {

		req.PGSettings.MaxPreparedTransactions = reqPGSettingsMaxPreparedTransactionsFlag
	}
	if reqPGSettingsMaxPredLocksPerTransactionFlag != 0 {

		req.PGSettings.MaxPredLocksPerTransaction = reqPGSettingsMaxPredLocksPerTransactionFlag
	}
	if reqPGSettingsMaxParallelWorkersPerGatherFlag != 0 {

		req.PGSettings.MaxParallelWorkersPerGather = reqPGSettingsMaxParallelWorkersPerGatherFlag
	}
	if reqPGSettingsMaxParallelWorkersFlag != 0 {

		req.PGSettings.MaxParallelWorkers = reqPGSettingsMaxParallelWorkersFlag
	}
	if reqPGSettingsMaxLogicalReplicationWorkersFlag != 0 {

		req.PGSettings.MaxLogicalReplicationWorkers = reqPGSettingsMaxLogicalReplicationWorkersFlag
	}
	if reqPGSettingsMaxLocksPerTransactionFlag != 0 {

		req.PGSettings.MaxLocksPerTransaction = reqPGSettingsMaxLocksPerTransactionFlag
	}
	if reqPGSettingsMaxFilesPerProcessFlag != 0 {

		req.PGSettings.MaxFilesPerProcess = reqPGSettingsMaxFilesPerProcessFlag
	}
	if reqPGSettingsLogTempFilesFlag != 0 {

		req.PGSettings.LogTempFiles = reqPGSettingsLogTempFilesFlag
	}
	if reqPGSettingsLogMinDurationStatementFlag != 0 {

		req.PGSettings.LogMinDurationStatement = reqPGSettingsLogMinDurationStatementFlag
	}
	if reqPGSettingsLogLinePrefixFlag != "" {

		req.PGSettings.LogLinePrefix = reqPGSettingsLogLinePrefixFlag
	}
	if reqPGSettingsLogErrorVerbosityFlag != "" {

		req.PGSettings.LogErrorVerbosity = reqPGSettingsLogErrorVerbosityFlag
	}
	if reqPGSettingsJitFlag != false {

		req.PGSettings.Jit = reqPGSettingsJitFlag
	}
	if reqPGSettingsIdleInTransactionSessionTimeoutFlag != 0 {

		req.PGSettings.IdleInTransactionSessionTimeout = reqPGSettingsIdleInTransactionSessionTimeoutFlag
	}
	if reqPGSettingsDefaultToastCompressionFlag != "" {

		req.PGSettings.DefaultToastCompression = reqPGSettingsDefaultToastCompressionFlag
	}
	if reqPGSettingsDeadlockTimeoutFlag != 0 {

		req.PGSettings.DeadlockTimeout = reqPGSettingsDeadlockTimeoutFlag
	}
	if reqPGSettingsBGWriterBgwriterLruMultiplierFlag != 0 {
		req.PGSettingsBGW = &v3.PGSettingsBGWriter{}
		req.PGSettingsBGWriter.BgwriterLruMultiplier = reqPGSettingsBGWriterBgwriterLruMultiplierFlag
	}
	if reqPGSettingsBGWriterBgwriterLruMaxpagesFlag != 0 {

		req.PGSettingsBGWriter.BgwriterLruMaxpages = reqPGSettingsBGWriterBgwriterLruMaxpagesFlag
	}
	if reqPGSettingsBGWriterBgwriterFlushAfterFlag != 0 {

		req.PGSettingsBGWriter.BgwriterFlushAfter = reqPGSettingsBGWriterBgwriterFlushAfterFlag
	}
	if reqPGSettingsBGWriterBgwriterDelayFlag != 0 {

		req.PGSettingsBGWriter.BgwriterDelay = reqPGSettingsBGWriterBgwriterDelayFlag
	}
	if reqPGSettingsAutovacuumLogAutovacuumMinDurationFlag != 0 {
		req.PGSettings = &v3.PGSettingsAutovacuum{}
		req.PGSettingsAutovacuum.LogAutovacuumMinDuration = reqPGSettingsAutovacuumLogAutovacuumMinDurationFlag
	}
	if reqPGSettingsAutovacuumAutovacuumVacuumThresholdFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumVacuumThreshold = reqPGSettingsAutovacuumAutovacuumVacuumThresholdFlag
	}
	if reqPGSettingsAutovacuumAutovacuumVacuumScaleFactorFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumVacuumScaleFactor = reqPGSettingsAutovacuumAutovacuumVacuumScaleFactorFlag
	}
	if reqPGSettingsAutovacuumAutovacuumVacuumCostLimitFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumVacuumCostLimit = reqPGSettingsAutovacuumAutovacuumVacuumCostLimitFlag
	}
	if reqPGSettingsAutovacuumAutovacuumVacuumCostDelayFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumVacuumCostDelay = reqPGSettingsAutovacuumAutovacuumVacuumCostDelayFlag
	}
	if reqPGSettingsAutovacuumAutovacuumNaptimeFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumNaptime = reqPGSettingsAutovacuumAutovacuumNaptimeFlag
	}
	if reqPGSettingsAutovacuumAutovacuumMaxWorkersFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumMaxWorkers = reqPGSettingsAutovacuumAutovacuumMaxWorkersFlag
	}
	if reqPGSettingsAutovacuumAutovacuumFreezeMaxAgeFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumFreezeMaxAge = reqPGSettingsAutovacuumAutovacuumFreezeMaxAgeFlag
	}
	if reqPGSettingsAutovacuumAutovacuumAnalyzeThresholdFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumAnalyzeThreshold = reqPGSettingsAutovacuumAutovacuumAnalyzeThresholdFlag
	}
	if reqPGSettingsAutovacuumAutovacuumAnalyzeScaleFactorFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumAnalyzeScaleFactor = reqPGSettingsAutovacuumAutovacuumAnalyzeScaleFactorFlag
	}
	if reqMigrationUsernameFlag != "" {
		req.Migratio = &v3.CreateDBAASServicePGRequestMigration{}
		req.Migration.Username = reqMigrationUsernameFlag
	}
	if reqMigrationSSLFlag != false {

		req.Migration.SSL = reqMigrationSSLFlag
	}
	if reqMigrationPortFlag != 0 {

		req.Migration.Port = reqMigrationPortFlag
	}
	if reqMigrationPasswordFlag != "" {

		req.Migration.Password = reqMigrationPasswordFlag
	}
	if reqMigrationMethodFlag != "" {

		req.Migration.Method = reqMigrationMethodFlag
	}
	if reqMigrationIgnoreDbsFlag != "" {

		req.Migration.IgnoreDbs = reqMigrationIgnoreDbsFlag
	}
	if reqMigrationHostFlag != "" {

		req.Migration.Host = reqMigrationHostFlag
	}
	if reqMigrationDbnameFlag != "" {

		req.Migration.Dbname = reqMigrationDbnameFlag
	}
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.CreateDBAASServicePGRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}
	req.ForkFromService = reqForkFromServiceFlag
	if reqBackupScheduleBackupMinuteFlag != 0 {
		req.BackupSchedul = &v3.CreateDBAASServicePGRequestBackupSchedule{}
		req.BackupSchedule.BackupMinute = reqBackupScheduleBackupMinuteFlag
	}
	if reqBackupScheduleBackupHourFlag != 0 {

		req.BackupSchedule.BackupHour = reqBackupScheduleBackupHourFlag
	}
	req.AdminUsername = reqAdminUsernameFlag
	req.AdminPassword = reqAdminPasswordFlag

	resp, err := client.CreateDBAASServicePG(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASServicePGCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-service-pg", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqBackupScheduleBackupHourFlag int64
	flagset.Int64Var(&reqBackupScheduleBackupHourFlag, "backup-schedule.backup-hour", 0, "The hour of day (in UTC) when backup for the service is started. New backup is only started if previous backup has already completed.")
	var reqBackupScheduleBackupMinuteFlag int64
	flagset.Int64Var(&reqBackupScheduleBackupMinuteFlag, "backup-schedule.backup-minute", 0, "The minute of an hour when backup for the service is started. New backup is only started if previous backup has already completed.")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqMigrationDbnameFlag string
	flagset.StringVar(&reqMigrationDbnameFlag, "migration.dbname", "", "Database name for bootstrapping the initial connection")
	var reqMigrationHostFlag string
	flagset.StringVar(&reqMigrationHostFlag, "migration.host", "", "Hostname or IP address of the server where to migrate data from")
	var reqMigrationIgnoreDbsFlag string
	flagset.StringVar(&reqMigrationIgnoreDbsFlag, "migration.ignore-dbs", "", "Comma-separated list of databases, which should be ignored during migration (supported by MySQL only at the moment)")
	var reqMigrationMethodFlag string
	flagset.StringVar(&reqMigrationMethodFlag, "migration.method", "", "")
	var reqMigrationPasswordFlag string
	flagset.StringVar(&reqMigrationPasswordFlag, "migration.password", "", "Password for authentication with the server where to migrate data from")
	var reqMigrationPortFlag int64
	flagset.Int64Var(&reqMigrationPortFlag, "migration.port", 0, "Port number of the server where to migrate data from")
	var reqMigrationSSLFlag bool
	flagset.BoolVar(&reqMigrationSSLFlag, "migration.ssl", false, "The server where to migrate data from is secured with SSL")
	var reqMigrationUsernameFlag string
	flagset.StringVar(&reqMigrationUsernameFlag, "migration.username", "", "User name for authentication with the server where to migrate data from")
	var reqPGSettingsAutovacuumAutovacuumAnalyzeScaleFactorFlag float64
	flagset.Float64Var(&reqPGSettingsAutovacuumAutovacuumAnalyzeScaleFactorFlag, "pg-settings.autovacuum.autovacuum_analyze_scale_factor", 0, "Specifies a fraction of the table size to add to autovacuum_analyze_threshold when deciding whether to trigger an ANALYZE. The default is 0.2 (20% of table size)")
	var reqPGSettingsAutovacuumAutovacuumAnalyzeThresholdFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumAnalyzeThresholdFlag, "pg-settings.autovacuum.autovacuum_analyze_threshold", 0, "Specifies the minimum number of inserted, updated or deleted tuples needed to trigger an ANALYZE in any one table. The default is 50 tuples.")
	var reqPGSettingsAutovacuumAutovacuumFreezeMaxAgeFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumFreezeMaxAgeFlag, "pg-settings.autovacuum.autovacuum_freeze_max_age", 0, "Specifies the maximum age (in transactions) that a table's pg_class.relfrozenxid field can attain before a VACUUM operation is forced to prevent transaction ID wraparound within the table. Note that the system will launch autovacuum processes to prevent wraparound even when autovacuum is otherwise disabled. This parameter will cause the server to be restarted.")
	var reqPGSettingsAutovacuumAutovacuumMaxWorkersFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumMaxWorkersFlag, "pg-settings.autovacuum.autovacuum_max_workers", 0, "Specifies the maximum number of autovacuum processes (other than the autovacuum launcher) that may be running at any one time. The default is three. This parameter can only be set at server start.")
	var reqPGSettingsAutovacuumAutovacuumNaptimeFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumNaptimeFlag, "pg-settings.autovacuum.autovacuum_naptime", 0, "Specifies the minimum delay between autovacuum runs on any given database. The delay is measured in seconds, and the default is one minute")
	var reqPGSettingsAutovacuumAutovacuumVacuumCostDelayFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumVacuumCostDelayFlag, "pg-settings.autovacuum.autovacuum_vacuum_cost_delay", 0, "Specifies the cost delay value that will be used in automatic VACUUM operations. If -1 is specified, the regular vacuum_cost_delay value will be used. The default value is 20 milliseconds")
	var reqPGSettingsAutovacuumAutovacuumVacuumCostLimitFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumVacuumCostLimitFlag, "pg-settings.autovacuum.autovacuum_vacuum_cost_limit", 0, "Specifies the cost limit value that will be used in automatic VACUUM operations. If -1 is specified (which is the default), the regular vacuum_cost_limit value will be used.")
	var reqPGSettingsAutovacuumAutovacuumVacuumScaleFactorFlag float64
	flagset.Float64Var(&reqPGSettingsAutovacuumAutovacuumVacuumScaleFactorFlag, "pg-settings.autovacuum.autovacuum_vacuum_scale_factor", 0, "Specifies a fraction of the table size to add to autovacuum_vacuum_threshold when deciding whether to trigger a VACUUM. The default is 0.2 (20% of table size)")
	var reqPGSettingsAutovacuumAutovacuumVacuumThresholdFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumAutovacuumVacuumThresholdFlag, "pg-settings.autovacuum.autovacuum_vacuum_threshold", 0, "Specifies the minimum number of updated or deleted tuples needed to trigger a VACUUM in any one table. The default is 50 tuples")
	var reqPGSettingsAutovacuumLogAutovacuumMinDurationFlag int
	flagset.IntVar(&reqPGSettingsAutovacuumLogAutovacuumMinDurationFlag, "pg-settings.autovacuum.log_autovacuum_min_duration", 0, "Causes each action executed by autovacuum to be logged if it ran for at least the specified number of milliseconds. Setting this to zero logs all autovacuum actions. Minus-one (the default) disables logging autovacuum actions.")
	var reqPGSettingsBGWriterBgwriterDelayFlag int
	flagset.IntVar(&reqPGSettingsBGWriterBgwriterDelayFlag, "pg-settings.bg-writer.bgwriter_delay", 0, "Specifies the delay between activity rounds for the background writer in milliseconds. Default is 200.")
	var reqPGSettingsBGWriterBgwriterFlushAfterFlag int
	flagset.IntVar(&reqPGSettingsBGWriterBgwriterFlushAfterFlag, "pg-settings.bg-writer.bgwriter_flush_after", 0, "Whenever more than bgwriter_flush_after bytes have been written by the background writer, attempt to force the OS to issue these writes to the underlying storage. Specified in kilobytes, default is 512. Setting of 0 disables forced writeback.")
	var reqPGSettingsBGWriterBgwriterLruMaxpagesFlag int
	flagset.IntVar(&reqPGSettingsBGWriterBgwriterLruMaxpagesFlag, "pg-settings.bg-writer.bgwriter_lru_maxpages", 0, "In each round, no more than this many buffers will be written by the background writer. Setting this to zero disables background writing. Default is 100.")
	var reqPGSettingsBGWriterBgwriterLruMultiplierFlag float64
	flagset.Float64Var(&reqPGSettingsBGWriterBgwriterLruMultiplierFlag, "pg-settings.bg-writer.bgwriter_lru_multiplier", 0, "The average recent need for new buffers is multiplied by bgwriter_lru_multiplier to arrive at an estimate of the number that will be needed during the next round, (up to bgwriter_lru_maxpages). 1.0 represents a “just in time” policy of writing exactly the number of buffers predicted to be needed. Larger values provide some cushion against spikes in demand, while smaller values intentionally leave writes to be done by server processes. The default is 2.0.")
	var reqPGSettingsDeadlockTimeoutFlag int
	flagset.IntVar(&reqPGSettingsDeadlockTimeoutFlag, "pg-settings.deadlock_timeout", 0, "This is the amount of time, in milliseconds, to wait on a lock before checking to see if there is a deadlock condition.")
	var reqPGSettingsDefaultToastCompressionFlag string
	flagset.StringVar(&reqPGSettingsDefaultToastCompressionFlag, "pg-settings.default_toast_compression", "", "Specifies the default TOAST compression method for values of compressible columns (the default is lz4).")
	var reqPGSettingsIdleInTransactionSessionTimeoutFlag int
	flagset.IntVar(&reqPGSettingsIdleInTransactionSessionTimeoutFlag, "pg-settings.idle_in_transaction_session_timeout", 0, "Time out sessions with open transactions after this number of milliseconds")
	var reqPGSettingsJitFlag bool
	flagset.BoolVar(&reqPGSettingsJitFlag, "pg-settings.jit", false, "Controls system-wide use of Just-in-Time Compilation (JIT).")
	var reqPGSettingsLogErrorVerbosityFlag string
	flagset.StringVar(&reqPGSettingsLogErrorVerbosityFlag, "pg-settings.log_error_verbosity", "", "Controls the amount of detail written in the server log for each message that is logged.")
	var reqPGSettingsLogLinePrefixFlag string
	flagset.StringVar(&reqPGSettingsLogLinePrefixFlag, "pg-settings.log_line_prefix", "", "Choose from one of the available log-formats. These can support popular log analyzers like pgbadger, pganalyze etc.")
	var reqPGSettingsLogMinDurationStatementFlag int
	flagset.IntVar(&reqPGSettingsLogMinDurationStatementFlag, "pg-settings.log_min_duration_statement", 0, "Log statements that take more than this number of milliseconds to run, -1 disables")
	var reqPGSettingsLogTempFilesFlag int
	flagset.IntVar(&reqPGSettingsLogTempFilesFlag, "pg-settings.log_temp_files", 0, "Log statements for each temporary file created larger than this number of kilobytes, -1 disables")
	var reqPGSettingsMaxFilesPerProcessFlag int
	flagset.IntVar(&reqPGSettingsMaxFilesPerProcessFlag, "pg-settings.max_files_per_process", 0, "PostgreSQL maximum number of files that can be open per process")
	var reqPGSettingsMaxLocksPerTransactionFlag int
	flagset.IntVar(&reqPGSettingsMaxLocksPerTransactionFlag, "pg-settings.max_locks_per_transaction", 0, "PostgreSQL maximum locks per transaction")
	var reqPGSettingsMaxLogicalReplicationWorkersFlag int
	flagset.IntVar(&reqPGSettingsMaxLogicalReplicationWorkersFlag, "pg-settings.max_logical_replication_workers", 0, "PostgreSQL maximum logical replication workers (taken from the pool of max_parallel_workers)")
	var reqPGSettingsMaxParallelWorkersFlag int
	flagset.IntVar(&reqPGSettingsMaxParallelWorkersFlag, "pg-settings.max_parallel_workers", 0, "Sets the maximum number of workers that the system can support for parallel queries")
	var reqPGSettingsMaxParallelWorkersPerGatherFlag int
	flagset.IntVar(&reqPGSettingsMaxParallelWorkersPerGatherFlag, "pg-settings.max_parallel_workers_per_gather", 0, "Sets the maximum number of workers that can be started by a single Gather or Gather Merge node")
	var reqPGSettingsMaxPredLocksPerTransactionFlag int
	flagset.IntVar(&reqPGSettingsMaxPredLocksPerTransactionFlag, "pg-settings.max_pred_locks_per_transaction", 0, "PostgreSQL maximum predicate locks per transaction")
	var reqPGSettingsMaxPreparedTransactionsFlag int
	flagset.IntVar(&reqPGSettingsMaxPreparedTransactionsFlag, "pg-settings.max_prepared_transactions", 0, "PostgreSQL maximum prepared transactions")
	var reqPGSettingsMaxReplicationSlotsFlag int
	flagset.IntVar(&reqPGSettingsMaxReplicationSlotsFlag, "pg-settings.max_replication_slots", 0, "PostgreSQL maximum replication slots")
	var reqPGSettingsMaxStackDepthFlag int
	flagset.IntVar(&reqPGSettingsMaxStackDepthFlag, "pg-settings.max_stack_depth", 0, "Maximum depth of the stack in bytes")
	var reqPGSettingsMaxStandbyArchiveDelayFlag int
	flagset.IntVar(&reqPGSettingsMaxStandbyArchiveDelayFlag, "pg-settings.max_standby_archive_delay", 0, "Max standby archive delay in milliseconds")
	var reqPGSettingsMaxStandbyStreamingDelayFlag int
	flagset.IntVar(&reqPGSettingsMaxStandbyStreamingDelayFlag, "pg-settings.max_standby_streaming_delay", 0, "Max standby streaming delay in milliseconds")
	var reqPGSettingsMaxWorkerProcessesFlag int
	flagset.IntVar(&reqPGSettingsMaxWorkerProcessesFlag, "pg-settings.max_worker_processes", 0, "Sets the maximum number of background processes that the system can support")
	var reqPGSettingsPGPartmanBgwIntervalFlag int
	flagset.IntVar(&reqPGSettingsPGPartmanBgwIntervalFlag, "pg-settings.pg_partman_bgw.interval", 0, "Sets the time interval to run pg_partman's scheduled tasks")
	var reqPGSettingsPGPartmanBgwRoleFlag string
	flagset.StringVar(&reqPGSettingsPGPartmanBgwRoleFlag, "pg-settings.pg_partman_bgw.role", "", "Controls which role to use for pg_partman's scheduled background tasks.")
	var reqPGSettingsPGStatMonitorPgsmEnableQueryPlanFlag bool
	flagset.BoolVar(&reqPGSettingsPGStatMonitorPgsmEnableQueryPlanFlag, "pg-settings.pg_stat_monitor.pgsm_enable_query_plan", false, "Enables or disables query plan monitoring")
	var reqPGSettingsPGStatMonitorPgsmMaxBucketsFlag int
	flagset.IntVar(&reqPGSettingsPGStatMonitorPgsmMaxBucketsFlag, "pg-settings.pg_stat_monitor.pgsm_max_buckets", 0, "Sets the maximum number of buckets ")
	var reqPGSettingsPGStatStatementsTrackFlag string
	flagset.StringVar(&reqPGSettingsPGStatStatementsTrackFlag, "pg-settings.pg_stat_statements.track", "", "Controls which statements are counted. Specify top to track top-level statements (those issued directly by clients), all to also track nested statements (such as statements invoked within functions), or none to disable statement statistics collection. The default value is top.")
	var reqPGSettingsTempFileLimitFlag int
	flagset.IntVar(&reqPGSettingsTempFileLimitFlag, "pg-settings.temp_file_limit", 0, "PostgreSQL temporary file limit in KiB, -1 for unlimited")
	var reqPGSettingsTimezoneFlag string
	flagset.StringVar(&reqPGSettingsTimezoneFlag, "pg-settings.timezone", "", "PostgreSQL service timezone")
	var reqPGSettingsTrackActivityQuerySizeFlag int
	flagset.IntVar(&reqPGSettingsTrackActivityQuerySizeFlag, "pg-settings.track_activity_query_size", 0, "Specifies the number of bytes reserved to track the currently executing command for each active session.")
	var reqPGSettingsTrackCommitTimestampFlag string
	flagset.StringVar(&reqPGSettingsTrackCommitTimestampFlag, "pg-settings.track_commit_timestamp", "", "Record commit time of transactions.")
	var reqPGSettingsTrackFunctionsFlag string
	flagset.StringVar(&reqPGSettingsTrackFunctionsFlag, "pg-settings.track_functions", "", "Enables tracking of function call counts and time used.")
	var reqPGSettingsTrackIoTimingFlag string
	flagset.StringVar(&reqPGSettingsTrackIoTimingFlag, "pg-settings.track_io_timing", "", "Enables timing of database I/O calls. This parameter is off by default, because it will repeatedly query the operating system for the current time, which may cause significant overhead on some platforms.")
	var reqPGSettingsWalMaxSlotWalKeepSizeFlag int
	flagset.IntVar(&reqPGSettingsWalMaxSlotWalKeepSizeFlag, "pg-settings.wal.max_slot_wal_keep_size", 0, "PostgreSQL maximum WAL size (MB) reserved for replication slots. Default is -1 (unlimited). wal_keep_size minimum WAL size setting takes precedence over this.")
	var reqPGSettingsWalMaxWalSendersFlag int
	flagset.IntVar(&reqPGSettingsWalMaxWalSendersFlag, "pg-settings.wal.max_wal_senders", 0, "PostgreSQL maximum WAL senders")
	var reqPGSettingsWalWalSenderTimeoutFlag int
	flagset.IntVar(&reqPGSettingsWalWalSenderTimeoutFlag, "pg-settings.wal.wal_sender_timeout", 0, "Terminate replication connections that are inactive for longer than this amount of time, in milliseconds.")
	var reqPGSettingsWalWalWriterDelayFlag int
	flagset.IntVar(&reqPGSettingsWalWalWriterDelayFlag, "pg-settings.wal.wal_writer_delay", 0, "WAL flush interval in milliseconds. Note that setting this value to lower than the default 200ms may negatively impact performance")
	var reqPgbouncerSettingsAutodbIdleTimeoutFlag int
	flagset.IntVar(&reqPgbouncerSettingsAutodbIdleTimeoutFlag, "pgbouncer-settings.autodb_idle_timeout", 0, "If the automatically created database pools have been unused this many seconds, they are freed. If 0 then timeout is disabled. [seconds]")
	var reqPgbouncerSettingsAutodbMaxDBConnectionsFlag int
	flagset.IntVar(&reqPgbouncerSettingsAutodbMaxDBConnectionsFlag, "pgbouncer-settings.autodb_max_db_connections", 0, "Do not allow more than this many server connections per database (regardless of user). Setting it to 0 means unlimited.")
	var reqPgbouncerSettingsAutodbPoolModeFlag string
	flagset.StringVar(&reqPgbouncerSettingsAutodbPoolModeFlag, "pgbouncer-settings.autodb_pool_mode", "", "PGBouncer pool mode")
	var reqPgbouncerSettingsAutodbPoolSizeFlag int
	flagset.IntVar(&reqPgbouncerSettingsAutodbPoolSizeFlag, "pgbouncer-settings.autodb_pool_size", 0, "If non-zero then create automatically a pool of that size per user when a pool doesn't exist.")
	var reqPgbouncerSettingsMaxPreparedStatementsFlag int
	flagset.IntVar(&reqPgbouncerSettingsMaxPreparedStatementsFlag, "pgbouncer-settings.max_prepared_statements", 0, "PgBouncer tracks protocol-level named prepared statements related commands sent by the client in transaction and statement pooling modes when max_prepared_statements is set to a non-zero value. Setting it to 0 disables prepared statements. max_prepared_statements defaults to 100, and its maximum is 3000.")
	var reqPgbouncerSettingsMinPoolSizeFlag int
	flagset.IntVar(&reqPgbouncerSettingsMinPoolSizeFlag, "pgbouncer-settings.min_pool_size", 0, "Add more server connections to pool if below this number. Improves behavior when usual load comes suddenly back after period of total inactivity. The value is effectively capped at the pool size.")
	var reqPgbouncerSettingsServerIdleTimeoutFlag int
	flagset.IntVar(&reqPgbouncerSettingsServerIdleTimeoutFlag, "pgbouncer-settings.server_idle_timeout", 0, "If a server connection has been idle more than this many seconds it will be dropped. If 0 then timeout is disabled. [seconds]")
	var reqPgbouncerSettingsServerLifetimeFlag int
	flagset.IntVar(&reqPgbouncerSettingsServerLifetimeFlag, "pgbouncer-settings.server_lifetime", 0, "The pooler will close an unused server connection that has been connected longer than this. [seconds]")
	var reqPgbouncerSettingsServerResetQueryAlwaysFlag bool
	flagset.BoolVar(&reqPgbouncerSettingsServerResetQueryAlwaysFlag, "pgbouncer-settings.server_reset_query_always", false, "Run server_reset_query (DISCARD ALL) in all pooling modes")
	var reqPglookoutSettingsMaxFailoverReplicationTimeLagFlag int
	flagset.IntVar(&reqPglookoutSettingsMaxFailoverReplicationTimeLagFlag, "pglookout-settings.max_failover_replication_time_lag", 0, "Number of seconds of master unavailability before triggering database failover to standby")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqSharedBuffersPercentageFlag int64
	flagset.Int64Var(&reqSharedBuffersPercentageFlag, "shared-buffers-percentage", 0, "Percentage of total RAM that the database server uses for shared memory buffers. Valid range is 20-60 (float), which corresponds to 20% - 60%. This setting adjusts the shared_buffers configuration value.")
	var reqSynchronousReplicationFlag string
	flagset.StringVar(&reqSynchronousReplicationFlag, "synchronous-replication", "", "")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")
	var reqTimescaledbSettingsMaxBackgroundWorkersFlag int
	flagset.IntVar(&reqTimescaledbSettingsMaxBackgroundWorkersFlag, "timescaledb-settings.max_background_workers", 0, "The number of background workers for timescaledb operations. You should configure this setting to the sum of your number of databases and the total number of concurrent background workers you want running at any given point in time.")
	var reqVariantFlag string
	flagset.StringVar(&reqVariantFlag, "variant", "", "")
	var reqVersionFlag string
	flagset.StringVar(&reqVersionFlag, "version", "", "Version")
	var reqWorkMemFlag int64
	flagset.Int64Var(&reqWorkMemFlag, "work-mem", 0, "Sets the maximum amount of memory to be used by a query operation (such as a sort or hash table) before writing to temporary disk files, in MB. Default is 1MB + 0.075% of total RAM (up to 32MB).")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASServicePGRequest
	req.WorkMem = reqWorkMemFlag
	req.Version = reqVersionFlag
	req.Variant = reqVariantFlag
	if reqTimescaledbSettingsMaxBackgroundWorkersFlag != 0 {
		req.TimescaledbSetti = &v3.JSONSchemaTimescaledb{}
		req.TimescaledbSettings.MaxBackgroundWorkers = reqTimescaledbSettingsMaxBackgroundWorkersFlag
	}
	req.TerminationProtection = reqTerminationProtectionFlag
	req.SynchronousReplication = reqSynchronousReplicationFlag
	req.SharedBuffersPercentage = reqSharedBuffersPercentageFlag
	req.Plan = reqPlanFlag
	if reqPglookoutSettingsMaxFailoverReplicationTimeLagFlag != 0 {
		req.PglookoutSettings = &v3.JSONSchemaPglookout{}
		req.PglookoutSettings.MaxFailoverReplicationTimeLag = reqPglookoutSettingsMaxFailoverReplicationTimeLagFlag
	}
	if reqPgbouncerSettingsServerResetQueryAlwaysFlag != false {
		req.PgbouncerSetting = &v3.JSONSchemaPgbouncer{}
		req.PgbouncerSettings.ServerResetQueryAlways = reqPgbouncerSettingsServerResetQueryAlwaysFlag
	}
	if reqPgbouncerSettingsServerLifetimeFlag != 0 {

		req.PgbouncerSettings.ServerLifetime = reqPgbouncerSettingsServerLifetimeFlag
	}
	if reqPgbouncerSettingsServerIdleTimeoutFlag != 0 {

		req.PgbouncerSettings.ServerIdleTimeout = reqPgbouncerSettingsServerIdleTimeoutFlag
	}
	if reqPgbouncerSettingsMinPoolSizeFlag != 0 {

		req.PgbouncerSettings.MinPoolSize = reqPgbouncerSettingsMinPoolSizeFlag
	}
	if reqPgbouncerSettingsMaxPreparedStatementsFlag != 0 {

		req.PgbouncerSettings.MaxPreparedStatements = reqPgbouncerSettingsMaxPreparedStatementsFlag
	}
	if reqPgbouncerSettingsAutodbPoolSizeFlag != 0 {

		req.PgbouncerSettings.AutodbPoolSize = reqPgbouncerSettingsAutodbPoolSizeFlag
	}
	if reqPgbouncerSettingsAutodbPoolModeFlag != "" {

		req.PgbouncerSettings.AutodbPoolMode = reqPgbouncerSettingsAutodbPoolModeFlag
	}
	if reqPgbouncerSettingsAutodbMaxDBConnectionsFlag != 0 {

		req.PgbouncerSettings.AutodbMaxDBConnections = reqPgbouncerSettingsAutodbMaxDBConnectionsFlag
	}
	if reqPgbouncerSettingsAutodbIdleTimeoutFlag != 0 {

		req.PgbouncerSettings.AutodbIdleTimeout = reqPgbouncerSettingsAutodbIdleTimeoutFlag
	}
	if reqPGSettingsWalWalWriterDelayFlag != 0 {
		req.PGSettings = &v3.PGSettingsWal{}
		req.PGSettingsWal.WalWriterDelay = reqPGSettingsWalWalWriterDelayFlag
	}
	if reqPGSettingsWalWalSenderTimeoutFlag != 0 {

		req.PGSettingsWal.WalSenderTimeout = reqPGSettingsWalWalSenderTimeoutFlag
	}
	if reqPGSettingsWalMaxWalSendersFlag != 0 {

		req.PGSettingsWal.MaxWalSenders = reqPGSettingsWalMaxWalSendersFlag
	}
	if reqPGSettingsWalMaxSlotWalKeepSizeFlag != 0 {

		req.PGSettingsWal.MaxSlotWalKeepSize = reqPGSettingsWalMaxSlotWalKeepSizeFlag
	}
	if reqPGSettingsTrackIoTimingFlag != "" {
		req.PGSettings = &v3.JSONSchemaPG{}
		req.PGSettings.TrackIoTiming = reqPGSettingsTrackIoTimingFlag
	}
	if reqPGSettingsTrackFunctionsFlag != "" {

		req.PGSettings.TrackFunctions = reqPGSettingsTrackFunctionsFlag
	}
	if reqPGSettingsTrackCommitTimestampFlag != "" {

		req.PGSettings.TrackCommitTimestamp = reqPGSettingsTrackCommitTimestampFlag
	}
	if reqPGSettingsTrackActivityQuerySizeFlag != 0 {

		req.PGSettings.TrackActivityQuerySize = reqPGSettingsTrackActivityQuerySizeFlag
	}
	if reqPGSettingsTimezoneFlag != "" {

		req.PGSettings.Timezone = reqPGSettingsTimezoneFlag
	}
	if reqPGSettingsTempFileLimitFlag != 0 {

		req.PGSettings.TempFileLimit = reqPGSettingsTempFileLimitFlag
	}
	if reqPGSettingsPGStatStatementsTrackFlag != "" {

		req.PGSettings.PGStatStatementsTrack = reqPGSettingsPGStatStatementsTrackFlag
	}
	if reqPGSettingsPGStatMonitorPgsmMaxBucketsFlag != 0 {

		req.PGSettings.PGStatMonitorPgsmMaxBuckets = reqPGSettingsPGStatMonitorPgsmMaxBucketsFlag
	}
	if reqPGSettingsPGStatMonitorPgsmEnableQueryPlanFlag != false {

		req.PGSettings.PGStatMonitorPgsmEnableQueryPlan = reqPGSettingsPGStatMonitorPgsmEnableQueryPlanFlag
	}
	if reqPGSettingsPGPartmanBgwRoleFlag != "" {

		req.PGSettings.PGPartmanBgwRole = reqPGSettingsPGPartmanBgwRoleFlag
	}
	if reqPGSettingsPGPartmanBgwIntervalFlag != 0 {

		req.PGSettings.PGPartmanBgwInterval = reqPGSettingsPGPartmanBgwIntervalFlag
	}
	if reqPGSettingsMaxWorkerProcessesFlag != 0 {

		req.PGSettings.MaxWorkerProcesses = reqPGSettingsMaxWorkerProcessesFlag
	}
	if reqPGSettingsMaxStandbyStreamingDelayFlag != 0 {

		req.PGSettings.MaxStandbyStreamingDelay = reqPGSettingsMaxStandbyStreamingDelayFlag
	}
	if reqPGSettingsMaxStandbyArchiveDelayFlag != 0 {

		req.PGSettings.MaxStandbyArchiveDelay = reqPGSettingsMaxStandbyArchiveDelayFlag
	}
	if reqPGSettingsMaxStackDepthFlag != 0 {

		req.PGSettings.MaxStackDepth = reqPGSettingsMaxStackDepthFlag
	}
	if reqPGSettingsMaxReplicationSlotsFlag != 0 {

		req.PGSettings.MaxReplicationSlots = reqPGSettingsMaxReplicationSlotsFlag
	}
	if reqPGSettingsMaxPreparedTransactionsFlag != 0 {

		req.PGSettings.MaxPreparedTransactions = reqPGSettingsMaxPreparedTransactionsFlag
	}
	if reqPGSettingsMaxPredLocksPerTransactionFlag != 0 {

		req.PGSettings.MaxPredLocksPerTransaction = reqPGSettingsMaxPredLocksPerTransactionFlag
	}
	if reqPGSettingsMaxParallelWorkersPerGatherFlag != 0 {

		req.PGSettings.MaxParallelWorkersPerGather = reqPGSettingsMaxParallelWorkersPerGatherFlag
	}
	if reqPGSettingsMaxParallelWorkersFlag != 0 {

		req.PGSettings.MaxParallelWorkers = reqPGSettingsMaxParallelWorkersFlag
	}
	if reqPGSettingsMaxLogicalReplicationWorkersFlag != 0 {

		req.PGSettings.MaxLogicalReplicationWorkers = reqPGSettingsMaxLogicalReplicationWorkersFlag
	}
	if reqPGSettingsMaxLocksPerTransactionFlag != 0 {

		req.PGSettings.MaxLocksPerTransaction = reqPGSettingsMaxLocksPerTransactionFlag
	}
	if reqPGSettingsMaxFilesPerProcessFlag != 0 {

		req.PGSettings.MaxFilesPerProcess = reqPGSettingsMaxFilesPerProcessFlag
	}
	if reqPGSettingsLogTempFilesFlag != 0 {

		req.PGSettings.LogTempFiles = reqPGSettingsLogTempFilesFlag
	}
	if reqPGSettingsLogMinDurationStatementFlag != 0 {

		req.PGSettings.LogMinDurationStatement = reqPGSettingsLogMinDurationStatementFlag
	}
	if reqPGSettingsLogLinePrefixFlag != "" {

		req.PGSettings.LogLinePrefix = reqPGSettingsLogLinePrefixFlag
	}
	if reqPGSettingsLogErrorVerbosityFlag != "" {

		req.PGSettings.LogErrorVerbosity = reqPGSettingsLogErrorVerbosityFlag
	}
	if reqPGSettingsJitFlag != false {

		req.PGSettings.Jit = reqPGSettingsJitFlag
	}
	if reqPGSettingsIdleInTransactionSessionTimeoutFlag != 0 {

		req.PGSettings.IdleInTransactionSessionTimeout = reqPGSettingsIdleInTransactionSessionTimeoutFlag
	}
	if reqPGSettingsDefaultToastCompressionFlag != "" {

		req.PGSettings.DefaultToastCompression = reqPGSettingsDefaultToastCompressionFlag
	}
	if reqPGSettingsDeadlockTimeoutFlag != 0 {

		req.PGSettings.DeadlockTimeout = reqPGSettingsDeadlockTimeoutFlag
	}
	if reqPGSettingsBGWriterBgwriterLruMultiplierFlag != 0 {
		req.PGSettingsBGW = &v3.PGSettingsBGWriter{}
		req.PGSettingsBGWriter.BgwriterLruMultiplier = reqPGSettingsBGWriterBgwriterLruMultiplierFlag
	}
	if reqPGSettingsBGWriterBgwriterLruMaxpagesFlag != 0 {

		req.PGSettingsBGWriter.BgwriterLruMaxpages = reqPGSettingsBGWriterBgwriterLruMaxpagesFlag
	}
	if reqPGSettingsBGWriterBgwriterFlushAfterFlag != 0 {

		req.PGSettingsBGWriter.BgwriterFlushAfter = reqPGSettingsBGWriterBgwriterFlushAfterFlag
	}
	if reqPGSettingsBGWriterBgwriterDelayFlag != 0 {

		req.PGSettingsBGWriter.BgwriterDelay = reqPGSettingsBGWriterBgwriterDelayFlag
	}
	if reqPGSettingsAutovacuumLogAutovacuumMinDurationFlag != 0 {
		req.PGSettings = &v3.PGSettingsAutovacuum{}
		req.PGSettingsAutovacuum.LogAutovacuumMinDuration = reqPGSettingsAutovacuumLogAutovacuumMinDurationFlag
	}
	if reqPGSettingsAutovacuumAutovacuumVacuumThresholdFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumVacuumThreshold = reqPGSettingsAutovacuumAutovacuumVacuumThresholdFlag
	}
	if reqPGSettingsAutovacuumAutovacuumVacuumScaleFactorFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumVacuumScaleFactor = reqPGSettingsAutovacuumAutovacuumVacuumScaleFactorFlag
	}
	if reqPGSettingsAutovacuumAutovacuumVacuumCostLimitFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumVacuumCostLimit = reqPGSettingsAutovacuumAutovacuumVacuumCostLimitFlag
	}
	if reqPGSettingsAutovacuumAutovacuumVacuumCostDelayFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumVacuumCostDelay = reqPGSettingsAutovacuumAutovacuumVacuumCostDelayFlag
	}
	if reqPGSettingsAutovacuumAutovacuumNaptimeFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumNaptime = reqPGSettingsAutovacuumAutovacuumNaptimeFlag
	}
	if reqPGSettingsAutovacuumAutovacuumMaxWorkersFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumMaxWorkers = reqPGSettingsAutovacuumAutovacuumMaxWorkersFlag
	}
	if reqPGSettingsAutovacuumAutovacuumFreezeMaxAgeFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumFreezeMaxAge = reqPGSettingsAutovacuumAutovacuumFreezeMaxAgeFlag
	}
	if reqPGSettingsAutovacuumAutovacuumAnalyzeThresholdFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumAnalyzeThreshold = reqPGSettingsAutovacuumAutovacuumAnalyzeThresholdFlag
	}
	if reqPGSettingsAutovacuumAutovacuumAnalyzeScaleFactorFlag != 0 {

		req.PGSettingsAutovacuum.AutovacuumAnalyzeScaleFactor = reqPGSettingsAutovacuumAutovacuumAnalyzeScaleFactorFlag
	}
	if reqMigrationUsernameFlag != "" {
		req.Migratio = &v3.UpdateDBAASServicePGRequestMigration{}
		req.Migration.Username = reqMigrationUsernameFlag
	}
	if reqMigrationSSLFlag != false {

		req.Migration.SSL = reqMigrationSSLFlag
	}
	if reqMigrationPortFlag != 0 {

		req.Migration.Port = reqMigrationPortFlag
	}
	if reqMigrationPasswordFlag != "" {

		req.Migration.Password = reqMigrationPasswordFlag
	}
	if reqMigrationMethodFlag != "" {

		req.Migration.Method = reqMigrationMethodFlag
	}
	if reqMigrationIgnoreDbsFlag != "" {

		req.Migration.IgnoreDbs = reqMigrationIgnoreDbsFlag
	}
	if reqMigrationHostFlag != "" {

		req.Migration.Host = reqMigrationHostFlag
	}
	if reqMigrationDbnameFlag != "" {

		req.Migration.Dbname = reqMigrationDbnameFlag
	}
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.UpdateDBAASServicePGRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}
	if reqBackupScheduleBackupMinuteFlag != 0 {
		req.BackupSchedul = &v3.UpdateDBAASServicePGRequestBackupSchedule{}
		req.BackupSchedule.BackupMinute = reqBackupScheduleBackupMinuteFlag
	}
	if reqBackupScheduleBackupHourFlag != 0 {

		req.BackupSchedule.BackupHour = reqBackupScheduleBackupHourFlag
	}

	resp, err := client.UpdateDBAASServicePG(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StartDBAASPGMaintenanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("start-dbaas-pg-maintenance", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StartDBAASPGMaintenance(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StopDBAASPGMigrationCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("stop-dbaas-pg-migration", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StopDBAASPGMigration(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASPGConnectionPoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-pg-connection-pool", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var reqDatabaseNameFlag string
	flagset.StringVar(&reqDatabaseNameFlag, "database-name", "", "")
	var reqModeFlag string
	flagset.StringVar(&reqModeFlag, "mode", "", "")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "")
	var reqSizeFlag int64
	flagset.Int64Var(&reqSizeFlag, "size", 0, "")
	var reqUsernameFlag string
	flagset.StringVar(&reqUsernameFlag, "username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASPGConnectionPoolRequest
	req.Username = reqUsernameFlag
	req.Size = reqSizeFlag
	req.Name = reqNameFlag
	req.Mode = reqModeFlag
	req.DatabaseName = reqDatabaseNameFlag

	resp, err := client.CreateDBAASPGConnectionPool(context.Background(), serviceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASPGConnectionPoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-pg-connection-pool", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var connectionPoolNameFlag string
	flagset.StringVar(&connectionPoolNameFlag, "ConnectionPoolName", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASPGConnectionPool(context.Background(), serviceNameFlag, connectionPoolNameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASPGConnectionPoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-pg-connection-pool", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var connectionPoolNameFlag string
	flagset.StringVar(&connectionPoolNameFlag, "ConnectionPoolName", "", "")
	var reqDatabaseNameFlag string
	flagset.StringVar(&reqDatabaseNameFlag, "database-name", "", "")
	var reqModeFlag string
	flagset.StringVar(&reqModeFlag, "mode", "", "")
	var reqSizeFlag int64
	flagset.Int64Var(&reqSizeFlag, "size", 0, "")
	var reqUsernameFlag string
	flagset.StringVar(&reqUsernameFlag, "username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASPGConnectionPoolRequest
	req.Username = reqUsernameFlag
	req.Size = reqSizeFlag
	req.Mode = reqModeFlag
	req.DatabaseName = reqDatabaseNameFlag

	resp, err := client.UpdateDBAASPGConnectionPool(context.Background(), serviceNameFlag, connectionPoolNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASPGDatabaseCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-pg-database", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var reqDatabaseNameFlag string
	flagset.StringVar(&reqDatabaseNameFlag, "database-name", "", "")
	var reqLCCollateFlag string
	flagset.StringVar(&reqLCCollateFlag, "lc-collate", "", "Default string sort order (LC_COLLATE) for PostgreSQL database")
	var reqLCCtypeFlag string
	flagset.StringVar(&reqLCCtypeFlag, "lc-ctype", "", "Default character classification (LC_CTYPE) for PostgreSQL database")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASPGDatabaseRequest
	req.LCCtype = reqLCCtypeFlag
	req.LCCollate = reqLCCollateFlag
	req.DatabaseName = reqDatabaseNameFlag

	resp, err := client.CreateDBAASPGDatabase(context.Background(), serviceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASPGDatabaseCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-pg-database", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var databaseNameFlag string
	flagset.StringVar(&databaseNameFlag, "DatabaseName", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASPGDatabase(context.Background(), serviceNameFlag, databaseNameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASPostgresUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-postgres-user", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var reqAllowReplicationFlag bool
	flagset.BoolVar(&reqAllowReplicationFlag, "allow-replication", false, "")
	var reqUsernameFlag string
	flagset.StringVar(&reqUsernameFlag, "username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASPostgresUserRequest
	req.Username = reqUsernameFlag
	req.AllowReplication = reqAllowReplicationFlag

	resp, err := client.CreateDBAASPostgresUser(context.Background(), serviceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASPostgresUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-postgres-user", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASPostgresUser(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASPostgresAllowReplicationCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-postgres-allow-replication", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")
	var reqAllowReplicationFlag bool
	flagset.BoolVar(&reqAllowReplicationFlag, "allow-replication", false, "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASPostgresAllowReplicationRequest
	req.AllowReplication = reqAllowReplicationFlag

	resp, err := client.UpdateDBAASPostgresAllowReplication(context.Background(), serviceNameFlag, usernameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetDBAASPostgresUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-dbaas-postgres-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")
	var reqPasswordFlag string
	flagset.StringVar(&reqPasswordFlag, "password", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ResetDBAASPostgresUserPasswordRequest
	req.Password = reqPasswordFlag

	resp, err := client.ResetDBAASPostgresUserPassword(context.Background(), serviceNameFlag, usernameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RevealDBAASPostgresUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reveal-dbaas-postgres-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RevealDBAASPostgresUserPassword(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASPGUpgradeCheckCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-pg-upgrade-check", flag.ExitOnError)
	var serviceFlag string
	flagset.StringVar(&serviceFlag, "Service", "", "")
	var reqTargetVersionFlag string
	flagset.StringVar(&reqTargetVersionFlag, "target-version", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASPGUpgradeCheckRequest
	req.TargetVersion = reqTargetVersionFlag

	resp, err := client.CreateDBAASPGUpgradeCheck(context.Background(), serviceFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASServiceRedisCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-service-redis", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASServiceRedis(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASServiceRedisCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-service-redis", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASServiceRedis(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASServiceRedisCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-service-redis", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqForkFromServiceFlag string
	flagset.StringVar(&reqForkFromServiceFlag, "fork-from-service", "", "")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqMigrationDbnameFlag string
	flagset.StringVar(&reqMigrationDbnameFlag, "migration.dbname", "", "Database name for bootstrapping the initial connection")
	var reqMigrationHostFlag string
	flagset.StringVar(&reqMigrationHostFlag, "migration.host", "", "Hostname or IP address of the server where to migrate data from")
	var reqMigrationIgnoreDbsFlag string
	flagset.StringVar(&reqMigrationIgnoreDbsFlag, "migration.ignore-dbs", "", "Comma-separated list of databases, which should be ignored during migration (supported by MySQL only at the moment)")
	var reqMigrationMethodFlag string
	flagset.StringVar(&reqMigrationMethodFlag, "migration.method", "", "")
	var reqMigrationPasswordFlag string
	flagset.StringVar(&reqMigrationPasswordFlag, "migration.password", "", "Password for authentication with the server where to migrate data from")
	var reqMigrationPortFlag int64
	flagset.Int64Var(&reqMigrationPortFlag, "migration.port", 0, "Port number of the server where to migrate data from")
	var reqMigrationSSLFlag bool
	flagset.BoolVar(&reqMigrationSSLFlag, "migration.ssl", false, "The server where to migrate data from is secured with SSL")
	var reqMigrationUsernameFlag string
	flagset.StringVar(&reqMigrationUsernameFlag, "migration.username", "", "User name for authentication with the server where to migrate data from")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqRecoveryBackupNameFlag string
	flagset.StringVar(&reqRecoveryBackupNameFlag, "recovery-backup-name", "", "Name of a backup to recover from for services that support backup names")
	var reqRedisSettingsAclChannelsDefaultFlag string
	flagset.StringVar(&reqRedisSettingsAclChannelsDefaultFlag, "redis-settings.acl_channels_default", "", "Determines default pub/sub channels' ACL for new users if ACL is not supplied. When this option is not defined, all_channels is assumed to keep backward compatibility. This option doesn't affect Redis configuration acl-pubsub-default.")
	var reqRedisSettingsIoThreadsFlag int
	flagset.IntVar(&reqRedisSettingsIoThreadsFlag, "redis-settings.io_threads", 0, "Set Redis IO thread count. Changing this will cause a restart of the Redis service.")
	var reqRedisSettingsLfuDecayTimeFlag int
	flagset.IntVar(&reqRedisSettingsLfuDecayTimeFlag, "redis-settings.lfu_decay_time", 0, "LFU maxmemory-policy counter decay time in minutes")
	var reqRedisSettingsLfuLogFactorFlag int
	flagset.IntVar(&reqRedisSettingsLfuLogFactorFlag, "redis-settings.lfu_log_factor", 0, "Counter logarithm factor for volatile-lfu and allkeys-lfu maxmemory-policies")
	var reqRedisSettingsMaxmemoryPolicyFlag string
	flagset.StringVar(&reqRedisSettingsMaxmemoryPolicyFlag, "redis-settings.maxmemory_policy", "", "Redis maxmemory-policy")
	var reqRedisSettingsNotifyKeyspaceEventsFlag string
	flagset.StringVar(&reqRedisSettingsNotifyKeyspaceEventsFlag, "redis-settings.notify_keyspace_events", "", "Set notify-keyspace-events option")
	var reqRedisSettingsNumberOfDatabasesFlag int
	flagset.IntVar(&reqRedisSettingsNumberOfDatabasesFlag, "redis-settings.number_of_databases", 0, "Set number of Redis databases. Changing this will cause a restart of the Redis service.")
	var reqRedisSettingsPersistenceFlag string
	flagset.StringVar(&reqRedisSettingsPersistenceFlag, "redis-settings.persistence", "", "When persistence is 'rdb', Redis does RDB dumps each 10 minutes if any key is changed. Also RDB dumps are done according to backup schedule for backup purposes. When persistence is 'off', no RDB dumps and backups are done, so data can be lost at any moment if service is restarted for any reason, or if service is powered off. Also service can't be forked.")
	var reqRedisSettingsPubsubClientOutputBufferLimitFlag int
	flagset.IntVar(&reqRedisSettingsPubsubClientOutputBufferLimitFlag, "redis-settings.pubsub_client_output_buffer_limit", 0, "Set output buffer limit for pub / sub clients in MB. The value is the hard limit, the soft limit is 1/4 of the hard limit. When setting the limit, be mindful of the available memory in the selected service plan.")
	var reqRedisSettingsSSLFlag bool
	flagset.BoolVar(&reqRedisSettingsSSLFlag, "redis-settings.ssl", false, "Require SSL to access Redis")
	var reqRedisSettingsTimeoutFlag int
	flagset.IntVar(&reqRedisSettingsTimeoutFlag, "redis-settings.timeout", 0, "Redis idle connection timeout in seconds")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASServiceRedisRequest
	req.TerminationProtection = reqTerminationProtectionFlag
	if reqRedisSettingsTimeoutFlag != 0 {
		req.RedisSettings = &v3.JSONSchemaRedis{}
		req.RedisSettings.Timeout = reqRedisSettingsTimeoutFlag
	}
	if reqRedisSettingsSSLFlag != false {

		req.RedisSettings.SSL = reqRedisSettingsSSLFlag
	}
	if reqRedisSettingsPubsubClientOutputBufferLimitFlag != 0 {

		req.RedisSettings.PubsubClientOutputBufferLimit = reqRedisSettingsPubsubClientOutputBufferLimitFlag
	}
	if reqRedisSettingsPersistenceFlag != "" {

		req.RedisSettings.Persistence = reqRedisSettingsPersistenceFlag
	}
	if reqRedisSettingsNumberOfDatabasesFlag != 0 {

		req.RedisSettings.NumberOfDatabases = reqRedisSettingsNumberOfDatabasesFlag
	}
	if reqRedisSettingsNotifyKeyspaceEventsFlag != "" {

		req.RedisSettings.NotifyKeyspaceEvents = reqRedisSettingsNotifyKeyspaceEventsFlag
	}
	if reqRedisSettingsMaxmemoryPolicyFlag != "" {

		req.RedisSettings.MaxmemoryPolicy = reqRedisSettingsMaxmemoryPolicyFlag
	}
	if reqRedisSettingsLfuLogFactorFlag != 0 {

		req.RedisSettings.LfuLogFactor = reqRedisSettingsLfuLogFactorFlag
	}
	if reqRedisSettingsLfuDecayTimeFlag != 0 {

		req.RedisSettings.LfuDecayTime = reqRedisSettingsLfuDecayTimeFlag
	}
	if reqRedisSettingsIoThreadsFlag != 0 {

		req.RedisSettings.IoThreads = reqRedisSettingsIoThreadsFlag
	}
	if reqRedisSettingsAclChannelsDefaultFlag != "" {

		req.RedisSettings.AclChannelsDefault = reqRedisSettingsAclChannelsDefaultFlag
	}
	req.RecoveryBackupName = reqRecoveryBackupNameFlag
	req.Plan = reqPlanFlag
	if reqMigrationUsernameFlag != "" {
		req.Migratio = &v3.CreateDBAASServiceRedisRequestMigration{}
		req.Migration.Username = reqMigrationUsernameFlag
	}
	if reqMigrationSSLFlag != false {

		req.Migration.SSL = reqMigrationSSLFlag
	}
	if reqMigrationPortFlag != 0 {

		req.Migration.Port = reqMigrationPortFlag
	}
	if reqMigrationPasswordFlag != "" {

		req.Migration.Password = reqMigrationPasswordFlag
	}
	if reqMigrationMethodFlag != "" {

		req.Migration.Method = reqMigrationMethodFlag
	}
	if reqMigrationIgnoreDbsFlag != "" {

		req.Migration.IgnoreDbs = reqMigrationIgnoreDbsFlag
	}
	if reqMigrationHostFlag != "" {

		req.Migration.Host = reqMigrationHostFlag
	}
	if reqMigrationDbnameFlag != "" {

		req.Migration.Dbname = reqMigrationDbnameFlag
	}
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.CreateDBAASServiceRedisRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}
	req.ForkFromService = reqForkFromServiceFlag

	resp, err := client.CreateDBAASServiceRedis(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASServiceRedisCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-service-redis", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqMigrationDbnameFlag string
	flagset.StringVar(&reqMigrationDbnameFlag, "migration.dbname", "", "Database name for bootstrapping the initial connection")
	var reqMigrationHostFlag string
	flagset.StringVar(&reqMigrationHostFlag, "migration.host", "", "Hostname or IP address of the server where to migrate data from")
	var reqMigrationIgnoreDbsFlag string
	flagset.StringVar(&reqMigrationIgnoreDbsFlag, "migration.ignore-dbs", "", "Comma-separated list of databases, which should be ignored during migration (supported by MySQL only at the moment)")
	var reqMigrationMethodFlag string
	flagset.StringVar(&reqMigrationMethodFlag, "migration.method", "", "")
	var reqMigrationPasswordFlag string
	flagset.StringVar(&reqMigrationPasswordFlag, "migration.password", "", "Password for authentication with the server where to migrate data from")
	var reqMigrationPortFlag int64
	flagset.Int64Var(&reqMigrationPortFlag, "migration.port", 0, "Port number of the server where to migrate data from")
	var reqMigrationSSLFlag bool
	flagset.BoolVar(&reqMigrationSSLFlag, "migration.ssl", false, "The server where to migrate data from is secured with SSL")
	var reqMigrationUsernameFlag string
	flagset.StringVar(&reqMigrationUsernameFlag, "migration.username", "", "User name for authentication with the server where to migrate data from")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqRedisSettingsAclChannelsDefaultFlag string
	flagset.StringVar(&reqRedisSettingsAclChannelsDefaultFlag, "redis-settings.acl_channels_default", "", "Determines default pub/sub channels' ACL for new users if ACL is not supplied. When this option is not defined, all_channels is assumed to keep backward compatibility. This option doesn't affect Redis configuration acl-pubsub-default.")
	var reqRedisSettingsIoThreadsFlag int
	flagset.IntVar(&reqRedisSettingsIoThreadsFlag, "redis-settings.io_threads", 0, "Set Redis IO thread count. Changing this will cause a restart of the Redis service.")
	var reqRedisSettingsLfuDecayTimeFlag int
	flagset.IntVar(&reqRedisSettingsLfuDecayTimeFlag, "redis-settings.lfu_decay_time", 0, "LFU maxmemory-policy counter decay time in minutes")
	var reqRedisSettingsLfuLogFactorFlag int
	flagset.IntVar(&reqRedisSettingsLfuLogFactorFlag, "redis-settings.lfu_log_factor", 0, "Counter logarithm factor for volatile-lfu and allkeys-lfu maxmemory-policies")
	var reqRedisSettingsMaxmemoryPolicyFlag string
	flagset.StringVar(&reqRedisSettingsMaxmemoryPolicyFlag, "redis-settings.maxmemory_policy", "", "Redis maxmemory-policy")
	var reqRedisSettingsNotifyKeyspaceEventsFlag string
	flagset.StringVar(&reqRedisSettingsNotifyKeyspaceEventsFlag, "redis-settings.notify_keyspace_events", "", "Set notify-keyspace-events option")
	var reqRedisSettingsNumberOfDatabasesFlag int
	flagset.IntVar(&reqRedisSettingsNumberOfDatabasesFlag, "redis-settings.number_of_databases", 0, "Set number of Redis databases. Changing this will cause a restart of the Redis service.")
	var reqRedisSettingsPersistenceFlag string
	flagset.StringVar(&reqRedisSettingsPersistenceFlag, "redis-settings.persistence", "", "When persistence is 'rdb', Redis does RDB dumps each 10 minutes if any key is changed. Also RDB dumps are done according to backup schedule for backup purposes. When persistence is 'off', no RDB dumps and backups are done, so data can be lost at any moment if service is restarted for any reason, or if service is powered off. Also service can't be forked.")
	var reqRedisSettingsPubsubClientOutputBufferLimitFlag int
	flagset.IntVar(&reqRedisSettingsPubsubClientOutputBufferLimitFlag, "redis-settings.pubsub_client_output_buffer_limit", 0, "Set output buffer limit for pub / sub clients in MB. The value is the hard limit, the soft limit is 1/4 of the hard limit. When setting the limit, be mindful of the available memory in the selected service plan.")
	var reqRedisSettingsSSLFlag bool
	flagset.BoolVar(&reqRedisSettingsSSLFlag, "redis-settings.ssl", false, "Require SSL to access Redis")
	var reqRedisSettingsTimeoutFlag int
	flagset.IntVar(&reqRedisSettingsTimeoutFlag, "redis-settings.timeout", 0, "Redis idle connection timeout in seconds")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASServiceRedisRequest
	req.TerminationProtection = reqTerminationProtectionFlag
	if reqRedisSettingsTimeoutFlag != 0 {
		req.RedisSettings = &v3.JSONSchemaRedis{}
		req.RedisSettings.Timeout = reqRedisSettingsTimeoutFlag
	}
	if reqRedisSettingsSSLFlag != false {

		req.RedisSettings.SSL = reqRedisSettingsSSLFlag
	}
	if reqRedisSettingsPubsubClientOutputBufferLimitFlag != 0 {

		req.RedisSettings.PubsubClientOutputBufferLimit = reqRedisSettingsPubsubClientOutputBufferLimitFlag
	}
	if reqRedisSettingsPersistenceFlag != "" {

		req.RedisSettings.Persistence = reqRedisSettingsPersistenceFlag
	}
	if reqRedisSettingsNumberOfDatabasesFlag != 0 {

		req.RedisSettings.NumberOfDatabases = reqRedisSettingsNumberOfDatabasesFlag
	}
	if reqRedisSettingsNotifyKeyspaceEventsFlag != "" {

		req.RedisSettings.NotifyKeyspaceEvents = reqRedisSettingsNotifyKeyspaceEventsFlag
	}
	if reqRedisSettingsMaxmemoryPolicyFlag != "" {

		req.RedisSettings.MaxmemoryPolicy = reqRedisSettingsMaxmemoryPolicyFlag
	}
	if reqRedisSettingsLfuLogFactorFlag != 0 {

		req.RedisSettings.LfuLogFactor = reqRedisSettingsLfuLogFactorFlag
	}
	if reqRedisSettingsLfuDecayTimeFlag != 0 {

		req.RedisSettings.LfuDecayTime = reqRedisSettingsLfuDecayTimeFlag
	}
	if reqRedisSettingsIoThreadsFlag != 0 {

		req.RedisSettings.IoThreads = reqRedisSettingsIoThreadsFlag
	}
	if reqRedisSettingsAclChannelsDefaultFlag != "" {

		req.RedisSettings.AclChannelsDefault = reqRedisSettingsAclChannelsDefaultFlag
	}
	req.Plan = reqPlanFlag
	if reqMigrationUsernameFlag != "" {
		req.Migratio = &v3.UpdateDBAASServiceRedisRequestMigration{}
		req.Migration.Username = reqMigrationUsernameFlag
	}
	if reqMigrationSSLFlag != false {

		req.Migration.SSL = reqMigrationSSLFlag
	}
	if reqMigrationPortFlag != 0 {

		req.Migration.Port = reqMigrationPortFlag
	}
	if reqMigrationPasswordFlag != "" {

		req.Migration.Password = reqMigrationPasswordFlag
	}
	if reqMigrationMethodFlag != "" {

		req.Migration.Method = reqMigrationMethodFlag
	}
	if reqMigrationIgnoreDbsFlag != "" {

		req.Migration.IgnoreDbs = reqMigrationIgnoreDbsFlag
	}
	if reqMigrationHostFlag != "" {

		req.Migration.Host = reqMigrationHostFlag
	}
	if reqMigrationDbnameFlag != "" {

		req.Migration.Dbname = reqMigrationDbnameFlag
	}
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.UpdateDBAASServiceRedisRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}

	resp, err := client.UpdateDBAASServiceRedis(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StartDBAASRedisMaintenanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("start-dbaas-redis-maintenance", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StartDBAASRedisMaintenance(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StopDBAASRedisMigrationCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("stop-dbaas-redis-migration", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StopDBAASRedisMigration(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StartDBAASRedisToValkeyUpgradeCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("start-dbaas-redis-to-valkey-upgrade", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StartDBAASRedisToValkeyUpgrade(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASRedisUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-redis-user", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var reqUsernameFlag string
	flagset.StringVar(&reqUsernameFlag, "username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASRedisUserRequest
	req.Username = reqUsernameFlag

	resp, err := client.CreateDBAASRedisUser(context.Background(), serviceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASRedisUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-redis-user", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASRedisUser(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetDBAASRedisUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-dbaas-redis-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")
	var reqPasswordFlag string
	flagset.StringVar(&reqPasswordFlag, "password", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ResetDBAASRedisUserPasswordRequest
	req.Password = reqPasswordFlag

	resp, err := client.ResetDBAASRedisUserPassword(context.Background(), serviceNameFlag, usernameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RevealDBAASRedisUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reveal-dbaas-redis-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RevealDBAASRedisUserPassword(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListDBAASServicesCmd(client *v3.Client) {
	resp, err := client.ListDBAASServices(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASServiceLogsCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-service-logs", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var reqLimitFlag int64
	flagset.Int64Var(&reqLimitFlag, "limit", 0, "How many log entries to receive at most, up to 500 (default: 100)")
	var reqOffsetFlag string
	flagset.StringVar(&reqOffsetFlag, "offset", "", "Opaque offset identifier")
	var reqSortOrderFlag string
	flagset.StringVar(&reqSortOrderFlag, "sort-order", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.GetDBAASServiceLogsRequest
	req.SortOrder = reqSortOrderFlag
	req.Offset = reqOffsetFlag
	req.Limit = reqLimitFlag

	resp, err := client.GetDBAASServiceLogs(context.Background(), serviceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASServiceMetricsCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-service-metrics", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var reqPeriodFlag string
	flagset.StringVar(&reqPeriodFlag, "period", "", "Metrics time period (default: hour)")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.GetDBAASServiceMetricsRequest
	req.Period = reqPeriodFlag

	resp, err := client.GetDBAASServiceMetrics(context.Background(), serviceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListDBAASServiceTypesCmd(client *v3.Client) {
	resp, err := client.ListDBAASServiceTypes(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASServiceTypeCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-service-type", flag.ExitOnError)
	var serviceTypeNameFlag string
	flagset.StringVar(&serviceTypeNameFlag, "ServiceTypeName", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASServiceType(context.Background(), serviceTypeNameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASServiceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-service", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASService(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASSettingsGrafanaCmd(client *v3.Client) {
	resp, err := client.GetDBAASSettingsGrafana(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASSettingsKafkaCmd(client *v3.Client) {
	resp, err := client.GetDBAASSettingsKafka(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASSettingsMysqlCmd(client *v3.Client) {
	resp, err := client.GetDBAASSettingsMysql(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASSettingsOpensearchCmd(client *v3.Client) {
	resp, err := client.GetDBAASSettingsOpensearch(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASSettingsPGCmd(client *v3.Client) {
	resp, err := client.GetDBAASSettingsPG(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASSettingsRedisCmd(client *v3.Client) {
	resp, err := client.GetDBAASSettingsRedis(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASSettingsValkeyCmd(client *v3.Client) {
	resp, err := client.GetDBAASSettingsValkey(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASTaskMigrationCheckCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-task-migration-check", flag.ExitOnError)
	var serviceFlag string
	flagset.StringVar(&serviceFlag, "Service", "", "")
	var reqIgnoreDbsFlag string
	flagset.StringVar(&reqIgnoreDbsFlag, "ignore-dbs", "", "Comma-separated list of databases, which should be ignored during migration (supported by MySQL only at the moment)")
	var reqMethodFlag string
	flagset.StringVar(&reqMethodFlag, "method", "", "")
	var reqSourceServiceURIFlag string
	flagset.StringVar(&reqSourceServiceURIFlag, "source-service-uri", "", "Service URI of the source MySQL or PostgreSQL database with admin credentials.")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASTaskMigrationCheckRequest
	req.SourceServiceURI = reqSourceServiceURIFlag
	req.Method = reqMethodFlag
	req.IgnoreDbs = reqIgnoreDbsFlag

	resp, err := client.CreateDBAASTaskMigrationCheck(context.Background(), serviceFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASTaskCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-task", flag.ExitOnError)
	var serviceFlag string
	flagset.StringVar(&serviceFlag, "Service", "", "")
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASTask(context.Background(), serviceFlag, v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASServiceValkeyCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-service-valkey", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASServiceValkey(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDBAASServiceValkeyCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dbaas-service-valkey", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDBAASServiceValkey(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASServiceValkeyCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-service-valkey", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqForkFromServiceFlag string
	flagset.StringVar(&reqForkFromServiceFlag, "fork-from-service", "", "")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqMigrationDbnameFlag string
	flagset.StringVar(&reqMigrationDbnameFlag, "migration.dbname", "", "Database name for bootstrapping the initial connection")
	var reqMigrationHostFlag string
	flagset.StringVar(&reqMigrationHostFlag, "migration.host", "", "Hostname or IP address of the server where to migrate data from")
	var reqMigrationIgnoreDbsFlag string
	flagset.StringVar(&reqMigrationIgnoreDbsFlag, "migration.ignore-dbs", "", "Comma-separated list of databases, which should be ignored during migration (supported by MySQL only at the moment)")
	var reqMigrationMethodFlag string
	flagset.StringVar(&reqMigrationMethodFlag, "migration.method", "", "")
	var reqMigrationPasswordFlag string
	flagset.StringVar(&reqMigrationPasswordFlag, "migration.password", "", "Password for authentication with the server where to migrate data from")
	var reqMigrationPortFlag int64
	flagset.Int64Var(&reqMigrationPortFlag, "migration.port", 0, "Port number of the server where to migrate data from")
	var reqMigrationSSLFlag bool
	flagset.BoolVar(&reqMigrationSSLFlag, "migration.ssl", false, "The server where to migrate data from is secured with SSL")
	var reqMigrationUsernameFlag string
	flagset.StringVar(&reqMigrationUsernameFlag, "migration.username", "", "User name for authentication with the server where to migrate data from")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqRecoveryBackupNameFlag string
	flagset.StringVar(&reqRecoveryBackupNameFlag, "recovery-backup-name", "", "Name of a backup to recover from for services that support backup names")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")
	var reqValkeySettingsAclChannelsDefaultFlag string
	flagset.StringVar(&reqValkeySettingsAclChannelsDefaultFlag, "valkey-settings.acl_channels_default", "", "Determines default pub/sub channels' ACL for new users if ACL is not supplied. When this option is not defined, all_channels is assumed to keep backward compatibility. This option doesn't affect Valkey configuration acl-pubsub-default.")
	var reqValkeySettingsIoThreadsFlag int
	flagset.IntVar(&reqValkeySettingsIoThreadsFlag, "valkey-settings.io_threads", 0, "Set Valkey IO thread count. Changing this will cause a restart of the Valkey service.")
	var reqValkeySettingsLfuDecayTimeFlag int
	flagset.IntVar(&reqValkeySettingsLfuDecayTimeFlag, "valkey-settings.lfu_decay_time", 0, "LFU maxmemory-policy counter decay time in minutes")
	var reqValkeySettingsLfuLogFactorFlag int
	flagset.IntVar(&reqValkeySettingsLfuLogFactorFlag, "valkey-settings.lfu_log_factor", 0, "Counter logarithm factor for volatile-lfu and allkeys-lfu maxmemory-policies")
	var reqValkeySettingsMaxmemoryPolicyFlag string
	flagset.StringVar(&reqValkeySettingsMaxmemoryPolicyFlag, "valkey-settings.maxmemory_policy", "", "Valkey maxmemory-policy")
	var reqValkeySettingsNotifyKeyspaceEventsFlag string
	flagset.StringVar(&reqValkeySettingsNotifyKeyspaceEventsFlag, "valkey-settings.notify_keyspace_events", "", "Set notify-keyspace-events option")
	var reqValkeySettingsNumberOfDatabasesFlag int
	flagset.IntVar(&reqValkeySettingsNumberOfDatabasesFlag, "valkey-settings.number_of_databases", 0, "Set number of Valkey databases. Changing this will cause a restart of the Valkey service.")
	var reqValkeySettingsPersistenceFlag string
	flagset.StringVar(&reqValkeySettingsPersistenceFlag, "valkey-settings.persistence", "", "When persistence is 'rdb', Valkey does RDB dumps each 10 minutes if any key is changed. Also RDB dumps are done according to backup schedule for backup purposes. When persistence is 'off', no RDB dumps and backups are done, so data can be lost at any moment if service is restarted for any reason, or if service is powered off. Also service can't be forked.")
	var reqValkeySettingsPubsubClientOutputBufferLimitFlag int
	flagset.IntVar(&reqValkeySettingsPubsubClientOutputBufferLimitFlag, "valkey-settings.pubsub_client_output_buffer_limit", 0, "Set output buffer limit for pub / sub clients in MB. The value is the hard limit, the soft limit is 1/4 of the hard limit. When setting the limit, be mindful of the available memory in the selected service plan.")
	var reqValkeySettingsSSLFlag bool
	flagset.BoolVar(&reqValkeySettingsSSLFlag, "valkey-settings.ssl", false, "Require SSL to access Valkey")
	var reqValkeySettingsTimeoutFlag int
	flagset.IntVar(&reqValkeySettingsTimeoutFlag, "valkey-settings.timeout", 0, "Valkey idle connection timeout in seconds")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASServiceValkeyRequest
	if reqValkeySettingsTimeoutFlag != 0 {
		req.ValkeySettings = &v3.JSONSchemaValkey{}
		req.ValkeySettings.Timeout = reqValkeySettingsTimeoutFlag
	}
	if reqValkeySettingsSSLFlag != false {

		req.ValkeySettings.SSL = reqValkeySettingsSSLFlag
	}
	if reqValkeySettingsPubsubClientOutputBufferLimitFlag != 0 {

		req.ValkeySettings.PubsubClientOutputBufferLimit = reqValkeySettingsPubsubClientOutputBufferLimitFlag
	}
	if reqValkeySettingsPersistenceFlag != "" {

		req.ValkeySettings.Persistence = reqValkeySettingsPersistenceFlag
	}
	if reqValkeySettingsNumberOfDatabasesFlag != 0 {

		req.ValkeySettings.NumberOfDatabases = reqValkeySettingsNumberOfDatabasesFlag
	}
	if reqValkeySettingsNotifyKeyspaceEventsFlag != "" {

		req.ValkeySettings.NotifyKeyspaceEvents = reqValkeySettingsNotifyKeyspaceEventsFlag
	}
	if reqValkeySettingsMaxmemoryPolicyFlag != "" {

		req.ValkeySettings.MaxmemoryPolicy = reqValkeySettingsMaxmemoryPolicyFlag
	}
	if reqValkeySettingsLfuLogFactorFlag != 0 {

		req.ValkeySettings.LfuLogFactor = reqValkeySettingsLfuLogFactorFlag
	}
	if reqValkeySettingsLfuDecayTimeFlag != 0 {

		req.ValkeySettings.LfuDecayTime = reqValkeySettingsLfuDecayTimeFlag
	}
	if reqValkeySettingsIoThreadsFlag != 0 {

		req.ValkeySettings.IoThreads = reqValkeySettingsIoThreadsFlag
	}
	if reqValkeySettingsAclChannelsDefaultFlag != "" {

		req.ValkeySettings.AclChannelsDefault = reqValkeySettingsAclChannelsDefaultFlag
	}
	req.TerminationProtection = reqTerminationProtectionFlag
	req.RecoveryBackupName = reqRecoveryBackupNameFlag
	req.Plan = reqPlanFlag
	if reqMigrationUsernameFlag != "" {
		req.Migratio = &v3.CreateDBAASServiceValkeyRequestMigration{}
		req.Migration.Username = reqMigrationUsernameFlag
	}
	if reqMigrationSSLFlag != false {

		req.Migration.SSL = reqMigrationSSLFlag
	}
	if reqMigrationPortFlag != 0 {

		req.Migration.Port = reqMigrationPortFlag
	}
	if reqMigrationPasswordFlag != "" {

		req.Migration.Password = reqMigrationPasswordFlag
	}
	if reqMigrationMethodFlag != "" {

		req.Migration.Method = reqMigrationMethodFlag
	}
	if reqMigrationIgnoreDbsFlag != "" {

		req.Migration.IgnoreDbs = reqMigrationIgnoreDbsFlag
	}
	if reqMigrationHostFlag != "" {

		req.Migration.Host = reqMigrationHostFlag
	}
	if reqMigrationDbnameFlag != "" {

		req.Migration.Dbname = reqMigrationDbnameFlag
	}
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.CreateDBAASServiceValkeyRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}
	req.ForkFromService = reqForkFromServiceFlag

	resp, err := client.CreateDBAASServiceValkey(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDBAASServiceValkeyCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dbaas-service-valkey", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")
	var reqMaintenanceDowFlag string
	flagset.StringVar(&reqMaintenanceDowFlag, "maintenance.dow", "", "Day of week for installing updates")
	var reqMaintenanceTimeFlag string
	flagset.StringVar(&reqMaintenanceTimeFlag, "maintenance.time", "", "Time for installing updates, UTC")
	var reqMigrationDbnameFlag string
	flagset.StringVar(&reqMigrationDbnameFlag, "migration.dbname", "", "Database name for bootstrapping the initial connection")
	var reqMigrationHostFlag string
	flagset.StringVar(&reqMigrationHostFlag, "migration.host", "", "Hostname or IP address of the server where to migrate data from")
	var reqMigrationIgnoreDbsFlag string
	flagset.StringVar(&reqMigrationIgnoreDbsFlag, "migration.ignore-dbs", "", "Comma-separated list of databases, which should be ignored during migration (supported by MySQL only at the moment)")
	var reqMigrationMethodFlag string
	flagset.StringVar(&reqMigrationMethodFlag, "migration.method", "", "")
	var reqMigrationPasswordFlag string
	flagset.StringVar(&reqMigrationPasswordFlag, "migration.password", "", "Password for authentication with the server where to migrate data from")
	var reqMigrationPortFlag int64
	flagset.Int64Var(&reqMigrationPortFlag, "migration.port", 0, "Port number of the server where to migrate data from")
	var reqMigrationSSLFlag bool
	flagset.BoolVar(&reqMigrationSSLFlag, "migration.ssl", false, "The server where to migrate data from is secured with SSL")
	var reqMigrationUsernameFlag string
	flagset.StringVar(&reqMigrationUsernameFlag, "migration.username", "", "User name for authentication with the server where to migrate data from")
	var reqPlanFlag string
	flagset.StringVar(&reqPlanFlag, "plan", "", "Subscription plan")
	var reqTerminationProtectionFlag bool
	flagset.BoolVar(&reqTerminationProtectionFlag, "termination-protection", false, "Service is protected against termination and powering off")
	var reqValkeySettingsAclChannelsDefaultFlag string
	flagset.StringVar(&reqValkeySettingsAclChannelsDefaultFlag, "valkey-settings.acl_channels_default", "", "Determines default pub/sub channels' ACL for new users if ACL is not supplied. When this option is not defined, all_channels is assumed to keep backward compatibility. This option doesn't affect Valkey configuration acl-pubsub-default.")
	var reqValkeySettingsIoThreadsFlag int
	flagset.IntVar(&reqValkeySettingsIoThreadsFlag, "valkey-settings.io_threads", 0, "Set Valkey IO thread count. Changing this will cause a restart of the Valkey service.")
	var reqValkeySettingsLfuDecayTimeFlag int
	flagset.IntVar(&reqValkeySettingsLfuDecayTimeFlag, "valkey-settings.lfu_decay_time", 0, "LFU maxmemory-policy counter decay time in minutes")
	var reqValkeySettingsLfuLogFactorFlag int
	flagset.IntVar(&reqValkeySettingsLfuLogFactorFlag, "valkey-settings.lfu_log_factor", 0, "Counter logarithm factor for volatile-lfu and allkeys-lfu maxmemory-policies")
	var reqValkeySettingsMaxmemoryPolicyFlag string
	flagset.StringVar(&reqValkeySettingsMaxmemoryPolicyFlag, "valkey-settings.maxmemory_policy", "", "Valkey maxmemory-policy")
	var reqValkeySettingsNotifyKeyspaceEventsFlag string
	flagset.StringVar(&reqValkeySettingsNotifyKeyspaceEventsFlag, "valkey-settings.notify_keyspace_events", "", "Set notify-keyspace-events option")
	var reqValkeySettingsNumberOfDatabasesFlag int
	flagset.IntVar(&reqValkeySettingsNumberOfDatabasesFlag, "valkey-settings.number_of_databases", 0, "Set number of Valkey databases. Changing this will cause a restart of the Valkey service.")
	var reqValkeySettingsPersistenceFlag string
	flagset.StringVar(&reqValkeySettingsPersistenceFlag, "valkey-settings.persistence", "", "When persistence is 'rdb', Valkey does RDB dumps each 10 minutes if any key is changed. Also RDB dumps are done according to backup schedule for backup purposes. When persistence is 'off', no RDB dumps and backups are done, so data can be lost at any moment if service is restarted for any reason, or if service is powered off. Also service can't be forked.")
	var reqValkeySettingsPubsubClientOutputBufferLimitFlag int
	flagset.IntVar(&reqValkeySettingsPubsubClientOutputBufferLimitFlag, "valkey-settings.pubsub_client_output_buffer_limit", 0, "Set output buffer limit for pub / sub clients in MB. The value is the hard limit, the soft limit is 1/4 of the hard limit. When setting the limit, be mindful of the available memory in the selected service plan.")
	var reqValkeySettingsSSLFlag bool
	flagset.BoolVar(&reqValkeySettingsSSLFlag, "valkey-settings.ssl", false, "Require SSL to access Valkey")
	var reqValkeySettingsTimeoutFlag int
	flagset.IntVar(&reqValkeySettingsTimeoutFlag, "valkey-settings.timeout", 0, "Valkey idle connection timeout in seconds")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDBAASServiceValkeyRequest
	if reqValkeySettingsTimeoutFlag != 0 {
		req.ValkeySettings = &v3.JSONSchemaValkey{}
		req.ValkeySettings.Timeout = reqValkeySettingsTimeoutFlag
	}
	if reqValkeySettingsSSLFlag != false {

		req.ValkeySettings.SSL = reqValkeySettingsSSLFlag
	}
	if reqValkeySettingsPubsubClientOutputBufferLimitFlag != 0 {

		req.ValkeySettings.PubsubClientOutputBufferLimit = reqValkeySettingsPubsubClientOutputBufferLimitFlag
	}
	if reqValkeySettingsPersistenceFlag != "" {

		req.ValkeySettings.Persistence = reqValkeySettingsPersistenceFlag
	}
	if reqValkeySettingsNumberOfDatabasesFlag != 0 {

		req.ValkeySettings.NumberOfDatabases = reqValkeySettingsNumberOfDatabasesFlag
	}
	if reqValkeySettingsNotifyKeyspaceEventsFlag != "" {

		req.ValkeySettings.NotifyKeyspaceEvents = reqValkeySettingsNotifyKeyspaceEventsFlag
	}
	if reqValkeySettingsMaxmemoryPolicyFlag != "" {

		req.ValkeySettings.MaxmemoryPolicy = reqValkeySettingsMaxmemoryPolicyFlag
	}
	if reqValkeySettingsLfuLogFactorFlag != 0 {

		req.ValkeySettings.LfuLogFactor = reqValkeySettingsLfuLogFactorFlag
	}
	if reqValkeySettingsLfuDecayTimeFlag != 0 {

		req.ValkeySettings.LfuDecayTime = reqValkeySettingsLfuDecayTimeFlag
	}
	if reqValkeySettingsIoThreadsFlag != 0 {

		req.ValkeySettings.IoThreads = reqValkeySettingsIoThreadsFlag
	}
	if reqValkeySettingsAclChannelsDefaultFlag != "" {

		req.ValkeySettings.AclChannelsDefault = reqValkeySettingsAclChannelsDefaultFlag
	}
	req.TerminationProtection = reqTerminationProtectionFlag
	req.Plan = reqPlanFlag
	if reqMigrationUsernameFlag != "" {
		req.Migratio = &v3.UpdateDBAASServiceValkeyRequestMigration{}
		req.Migration.Username = reqMigrationUsernameFlag
	}
	if reqMigrationSSLFlag != false {

		req.Migration.SSL = reqMigrationSSLFlag
	}
	if reqMigrationPortFlag != 0 {

		req.Migration.Port = reqMigrationPortFlag
	}
	if reqMigrationPasswordFlag != "" {

		req.Migration.Password = reqMigrationPasswordFlag
	}
	if reqMigrationMethodFlag != "" {

		req.Migration.Method = reqMigrationMethodFlag
	}
	if reqMigrationIgnoreDbsFlag != "" {

		req.Migration.IgnoreDbs = reqMigrationIgnoreDbsFlag
	}
	if reqMigrationHostFlag != "" {

		req.Migration.Host = reqMigrationHostFlag
	}
	if reqMigrationDbnameFlag != "" {

		req.Migration.Dbname = reqMigrationDbnameFlag
	}
	if reqMaintenanceTimeFlag != "" {
		req.Maintenanc = &v3.UpdateDBAASServiceValkeyRequestMaintenance{}
		req.Maintenance.Time = reqMaintenanceTimeFlag
	}
	if reqMaintenanceDowFlag != "" {

		req.Maintenance.Dow = reqMaintenanceDowFlag
	}

	resp, err := client.UpdateDBAASServiceValkey(context.Background(), nameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StartDBAASValkeyMaintenanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("start-dbaas-valkey-maintenance", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StartDBAASValkeyMaintenance(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StopDBAASValkeyMigrationCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("stop-dbaas-valkey-migration", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StopDBAASValkeyMigration(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDBAASValkeyUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dbaas-valkey-user", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var reqUsernameFlag string
	flagset.StringVar(&reqUsernameFlag, "username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDBAASValkeyUserRequest
	req.Username = reqUsernameFlag

	resp, err := client.CreateDBAASValkeyUser(context.Background(), serviceNameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDBAASValkeyUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dbaas-valkey-user", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDBAASValkeyUser(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetDBAASValkeyUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-dbaas-valkey-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")
	var reqPasswordFlag string
	flagset.StringVar(&reqPasswordFlag, "password", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ResetDBAASValkeyUserPasswordRequest
	req.Password = reqPasswordFlag

	resp, err := client.ResetDBAASValkeyUserPassword(context.Background(), serviceNameFlag, usernameFlag, req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RevealDBAASValkeyUserPasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reveal-dbaas-valkey-user-password", flag.ExitOnError)
	var serviceNameFlag string
	flagset.StringVar(&serviceNameFlag, "ServiceName", "", "")
	var usernameFlag string
	flagset.StringVar(&usernameFlag, "Username", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RevealDBAASValkeyUserPassword(context.Background(), serviceNameFlag, usernameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListDeployTargetsCmd(client *v3.Client) {
	resp, err := client.ListDeployTargets(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDeployTargetCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-deploy-target", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDeployTarget(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListDNSDomainsCmd(client *v3.Client) {
	resp, err := client.ListDNSDomains(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDNSDomainCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dns-domain", flag.ExitOnError)
	var reqUnicodeNameFlag string
	flagset.StringVar(&reqUnicodeNameFlag, "unicode-name", "", "Domain name")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDNSDomainRequest
	req.UnicodeName = reqUnicodeNameFlag

	resp, err := client.CreateDNSDomain(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListDNSDomainRecordsCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("list-dns-domain-records", flag.ExitOnError)
	var domainIDFlag string
	flagset.StringVar(&domainIDFlag, "DomainID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ListDNSDomainRecords(context.Background(), v3.UUID(domainIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateDNSDomainRecordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-dns-domain-record", flag.ExitOnError)
	var domainIDFlag string
	flagset.StringVar(&domainIDFlag, "DomainID", "", "")
	var reqContentFlag string
	flagset.StringVar(&reqContentFlag, "content", "", "DNS domain record content")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "DNS domain record name")
	var reqPriorityFlag int64
	flagset.Int64Var(&reqPriorityFlag, "priority", 0, "DNS domain record priority")
	var reqTtlFlag int64
	flagset.Int64Var(&reqTtlFlag, "ttl", 0, "DNS domain record TTL")
	var reqTypeFlag string
	flagset.StringVar(&reqTypeFlag, "type", "", "DNS domain record type")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateDNSDomainRecordRequest
	req.Type = reqTypeFlag
	req.Ttl = reqTtlFlag
	req.Priority = reqPriorityFlag
	req.Name = reqNameFlag
	req.Content = reqContentFlag

	resp, err := client.CreateDNSDomainRecord(context.Background(), v3.UUID(domainIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDNSDomainRecordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dns-domain-record", flag.ExitOnError)
	var domainIDFlag string
	flagset.StringVar(&domainIDFlag, "DomainID", "", "")
	var recordIDFlag string
	flagset.StringVar(&recordIDFlag, "RecordID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDNSDomainRecord(context.Background(), v3.UUID(domainIDFlag), v3.UUID(recordIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDNSDomainRecordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dns-domain-record", flag.ExitOnError)
	var domainIDFlag string
	flagset.StringVar(&domainIDFlag, "DomainID", "", "")
	var recordIDFlag string
	flagset.StringVar(&recordIDFlag, "RecordID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDNSDomainRecord(context.Background(), v3.UUID(domainIDFlag), v3.UUID(recordIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateDNSDomainRecordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-dns-domain-record", flag.ExitOnError)
	var domainIDFlag string
	flagset.StringVar(&domainIDFlag, "DomainID", "", "")
	var recordIDFlag string
	flagset.StringVar(&recordIDFlag, "RecordID", "", "")
	var reqContentFlag string
	flagset.StringVar(&reqContentFlag, "content", "", "DNS domain record content")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "DNS domain record name")
	var reqPriorityFlag int64
	flagset.Int64Var(&reqPriorityFlag, "priority", 0, "DNS domain record priority")
	var reqTtlFlag int64
	flagset.Int64Var(&reqTtlFlag, "ttl", 0, "DNS domain record TTL")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateDNSDomainRecordRequest
	req.Ttl = reqTtlFlag
	req.Priority = reqPriorityFlag
	req.Name = reqNameFlag
	req.Content = reqContentFlag

	resp, err := client.UpdateDNSDomainRecord(context.Background(), v3.UUID(domainIDFlag), v3.UUID(recordIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteDNSDomainCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-dns-domain", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteDNSDomain(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDNSDomainCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dns-domain", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDNSDomain(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetDNSDomainZoneFileCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-dns-domain-zone-file", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetDNSDomainZoneFile(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListElasticIPSCmd(client *v3.Client) {
	resp, err := client.ListElasticIPS(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateElasticIPCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-elastic-ip", flag.ExitOnError)
	var reqAddressfamilyFlag string
	flagset.StringVar(&reqAddressfamilyFlag, "addressfamily", "", "Elastic IP address family (default: :inet4)")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Elastic IP description")
	var reqHealthcheckIntervalFlag int64
	flagset.Int64Var(&reqHealthcheckIntervalFlag, "healthcheck.interval", 0, "Interval between the checks in seconds (default: 10)")
	var reqHealthcheckModeFlag string
	flagset.StringVar(&reqHealthcheckModeFlag, "healthcheck.mode", "", "Health check mode")
	var reqHealthcheckPortFlag int64
	flagset.Int64Var(&reqHealthcheckPortFlag, "healthcheck.port", 0, "Health check port")
	var reqHealthcheckStrikesFailFlag int64
	flagset.Int64Var(&reqHealthcheckStrikesFailFlag, "healthcheck.strikes-fail", 0, "Number of attempts before considering the target unhealthy (default: 3)")
	var reqHealthcheckStrikesOkFlag int64
	flagset.Int64Var(&reqHealthcheckStrikesOkFlag, "healthcheck.strikes-ok", 0, "Number of attempts before considering the target healthy (default: 2)")
	var reqHealthcheckTimeoutFlag int64
	flagset.Int64Var(&reqHealthcheckTimeoutFlag, "healthcheck.timeout", 0, "Health check timeout value in seconds (default: 2)")
	var reqHealthcheckTlsSkipVerifyFlag bool
	flagset.BoolVar(&reqHealthcheckTlsSkipVerifyFlag, "healthcheck.tls-skip-verify", false, "Skip TLS verification")
	var reqHealthcheckTlsSNIFlag string
	flagset.StringVar(&reqHealthcheckTlsSNIFlag, "healthcheck.tls-sni", "", "An optional domain or subdomain to check TLS against")
	var reqHealthcheckURIFlag string
	flagset.StringVar(&reqHealthcheckURIFlag, "healthcheck.uri", "", "An endpoint to use for the health check, for example '/status'")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateElasticIPRequest
	if reqHealthcheckURIFlag != "" {
		req.Healthcheck = &v3.ElasticIPHealthcheck{}
		req.Healthcheck.URI = reqHealthcheckURIFlag
	}
	if reqHealthcheckTlsSNIFlag != "" {

		req.Healthcheck.TlsSNI = reqHealthcheckTlsSNIFlag
	}
	if reqHealthcheckTlsSkipVerifyFlag != false {

		req.Healthcheck.TlsSkipVerify = reqHealthcheckTlsSkipVerifyFlag
	}
	if reqHealthcheckTimeoutFlag != 0 {

		req.Healthcheck.Timeout = reqHealthcheckTimeoutFlag
	}
	if reqHealthcheckStrikesOkFlag != 0 {

		req.Healthcheck.StrikesOk = reqHealthcheckStrikesOkFlag
	}
	if reqHealthcheckStrikesFailFlag != 0 {

		req.Healthcheck.StrikesFail = reqHealthcheckStrikesFailFlag
	}
	if reqHealthcheckPortFlag != 0 {

		req.Healthcheck.Port = reqHealthcheckPortFlag
	}
	if reqHealthcheckModeFlag != "" {

		req.Healthcheck.Mode = reqHealthcheckModeFlag
	}
	if reqHealthcheckIntervalFlag != 0 {

		req.Healthcheck.Interval = reqHealthcheckIntervalFlag
	}
	req.Description = reqDescriptionFlag
	req.Addressfamily = reqAddressfamilyFlag

	resp, err := client.CreateElasticIP(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteElasticIPCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-elastic-ip", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteElasticIP(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetElasticIPCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-elastic-ip", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetElasticIP(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateElasticIPCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-elastic-ip", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Elastic IP description")
	var reqHealthcheckIntervalFlag int64
	flagset.Int64Var(&reqHealthcheckIntervalFlag, "healthcheck.interval", 0, "Interval between the checks in seconds (default: 10)")
	var reqHealthcheckModeFlag string
	flagset.StringVar(&reqHealthcheckModeFlag, "healthcheck.mode", "", "Health check mode")
	var reqHealthcheckPortFlag int64
	flagset.Int64Var(&reqHealthcheckPortFlag, "healthcheck.port", 0, "Health check port")
	var reqHealthcheckStrikesFailFlag int64
	flagset.Int64Var(&reqHealthcheckStrikesFailFlag, "healthcheck.strikes-fail", 0, "Number of attempts before considering the target unhealthy (default: 3)")
	var reqHealthcheckStrikesOkFlag int64
	flagset.Int64Var(&reqHealthcheckStrikesOkFlag, "healthcheck.strikes-ok", 0, "Number of attempts before considering the target healthy (default: 2)")
	var reqHealthcheckTimeoutFlag int64
	flagset.Int64Var(&reqHealthcheckTimeoutFlag, "healthcheck.timeout", 0, "Health check timeout value in seconds (default: 2)")
	var reqHealthcheckTlsSkipVerifyFlag bool
	flagset.BoolVar(&reqHealthcheckTlsSkipVerifyFlag, "healthcheck.tls-skip-verify", false, "Skip TLS verification")
	var reqHealthcheckTlsSNIFlag string
	flagset.StringVar(&reqHealthcheckTlsSNIFlag, "healthcheck.tls-sni", "", "An optional domain or subdomain to check TLS against")
	var reqHealthcheckURIFlag string
	flagset.StringVar(&reqHealthcheckURIFlag, "healthcheck.uri", "", "An endpoint to use for the health check, for example '/status'")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateElasticIPRequest
	if reqHealthcheckURIFlag != "" {
		req.Healthcheck = &v3.ElasticIPHealthcheck{}
		req.Healthcheck.URI = reqHealthcheckURIFlag
	}
	if reqHealthcheckTlsSNIFlag != "" {

		req.Healthcheck.TlsSNI = reqHealthcheckTlsSNIFlag
	}
	if reqHealthcheckTlsSkipVerifyFlag != false {

		req.Healthcheck.TlsSkipVerify = reqHealthcheckTlsSkipVerifyFlag
	}
	if reqHealthcheckTimeoutFlag != 0 {

		req.Healthcheck.Timeout = reqHealthcheckTimeoutFlag
	}
	if reqHealthcheckStrikesOkFlag != 0 {

		req.Healthcheck.StrikesOk = reqHealthcheckStrikesOkFlag
	}
	if reqHealthcheckStrikesFailFlag != 0 {

		req.Healthcheck.StrikesFail = reqHealthcheckStrikesFailFlag
	}
	if reqHealthcheckPortFlag != 0 {

		req.Healthcheck.Port = reqHealthcheckPortFlag
	}
	if reqHealthcheckModeFlag != "" {

		req.Healthcheck.Mode = reqHealthcheckModeFlag
	}
	if reqHealthcheckIntervalFlag != 0 {

		req.Healthcheck.Interval = reqHealthcheckIntervalFlag
	}
	req.Description = reqDescriptionFlag

	resp, err := client.UpdateElasticIP(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetElasticIPFieldCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-elastic-ip-field", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var fieldFlag string
	flagset.StringVar(&fieldFlag, "Field", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ResetElasticIPField(context.Background(), v3.UUID(idFlag), fieldFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func AttachInstanceToElasticIPCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("attach-instance-to-elastic-ip", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqInstanceIDFlag string
	flagset.StringVar(&reqInstanceIDFlag, "instance.id", "", "Instance ID")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.AttachInstanceToElasticIPRequest
	if v3.UUID(reqInstanceIDFlag) != "" {
		req.Instance = &v3.InstanceTarget{}
		req.Instance.ID = v3.UUID(reqInstanceIDFlag)
	}

	resp, err := client.AttachInstanceToElasticIP(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DetachInstanceFromElasticIPCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("detach-instance-from-elastic-ip", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqInstanceIDFlag string
	flagset.StringVar(&reqInstanceIDFlag, "instance.id", "", "Instance ID")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.DetachInstanceFromElasticIPRequest
	if v3.UUID(reqInstanceIDFlag) != "" {
		req.Instance = &v3.InstanceTarget{}
		req.Instance.ID = v3.UUID(reqInstanceIDFlag)
	}

	resp, err := client.DetachInstanceFromElasticIP(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetEnvImpactCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-env-impact", flag.ExitOnError)
	var periodFlag string
	flagset.StringVar(&periodFlag, "Period", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetEnvImpact(context.Background(), periodFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListEventsCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("list-events", flag.ExitOnError)
	var fromFlag time.Time
	flagset.TimeTimeVar(&fromFlag, "From", "", "")
	var toFlag time.Time
	flagset.TimeTimeVar(&toFlag, "To", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ListEvents(context.Background(), fromFlag, toFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetIAMOrganizationPolicyCmd(client *v3.Client) {
	resp, err := client.GetIAMOrganizationPolicy(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateIAMOrganizationPolicyCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-iam-organization-policy", flag.ExitOnError)
	var reqDefaultServiceStrategyFlag string
	flagset.StringVar(&reqDefaultServiceStrategyFlag, "default-service-strategy", "", "IAM default service strategy")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateIAMOrganizationPolicyRequest
	req.DefaultServiceStrategy = reqDefaultServiceStrategyFlag

	resp, err := client.UpdateIAMOrganizationPolicy(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetIAMOrganizationPolicyCmd(client *v3.Client) {
	resp, err := client.ResetIAMOrganizationPolicy(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListIAMRolesCmd(client *v3.Client) {
	resp, err := client.ListIAMRoles(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateIAMRoleCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-iam-role", flag.ExitOnError)
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "IAM Role description")
	var reqEditableFlag bool
	flagset.BoolVar(&reqEditableFlag, "editable", false, "Sets if the IAM Role Policy is editable or not (default: true). This setting cannot be changed after creation")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "IAM Role name")
	var reqPolicyDefaultServiceStrategyFlag string
	flagset.StringVar(&reqPolicyDefaultServiceStrategyFlag, "policy.default-service-strategy", "", "IAM default service strategy")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateIAMRoleRequest
	if reqPolicyDefaultServiceStrategyFlag != "" {
		req.Po = &v3.IAMPolicy{}
		req.Policy.DefaultServiceStrategy = reqPolicyDefaultServiceStrategyFlag
	}
	req.Name = reqNameFlag
	req.Editable = reqEditableFlag
	req.Description = reqDescriptionFlag

	resp, err := client.CreateIAMRole(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteIAMRoleCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-iam-role", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteIAMRole(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetIAMRoleCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-iam-role", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetIAMRole(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateIAMRoleCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-iam-role", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "IAM Role description")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateIAMRoleRequest
	req.Description = reqDescriptionFlag

	resp, err := client.UpdateIAMRole(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateIAMRolePolicyCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-iam-role-policy", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDefaultServiceStrategyFlag string
	flagset.StringVar(&reqDefaultServiceStrategyFlag, "default-service-strategy", "", "IAM default service strategy")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateIAMRolePolicyRequest
	req.DefaultServiceStrategy = reqDefaultServiceStrategyFlag

	resp, err := client.UpdateIAMRolePolicy(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListInstancesCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("list-instances", flag.ExitOnError)
	var managerIDFlag string
	flagset.StringVar(&managerIDFlag, "ManagerID", "", "")
	var managerTypeFlag string
	flagset.StringVar(&managerTypeFlag, "ManagerType", "", "")
	var ipAddressFlag string
	flagset.StringVar(&ipAddressFlag, "IPAddress", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ListInstances(context.Background(), v3.UUID(managerIDFlag), managerTypeFlag, ipAddressFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-instance", flag.ExitOnError)
	var reqAutoStartFlag bool
	flagset.BoolVar(&reqAutoStartFlag, "auto-start", false, "Start Instance on creation (default: true)")
	var reqDeployTargetDescriptionFlag string
	flagset.StringVar(&reqDeployTargetDescriptionFlag, "deploy-target.description", "", "Deploy Target description")
	var reqDeployTargetIDFlag string
	flagset.StringVar(&reqDeployTargetIDFlag, "deploy-target.id", "", "Deploy Target ID")
	var reqDeployTargetNameFlag string
	flagset.StringVar(&reqDeployTargetNameFlag, "deploy-target.name", "", "Deploy Target name")
	var reqDeployTargetTypeFlag string
	flagset.StringVar(&reqDeployTargetTypeFlag, "deploy-target.type", "", "Deploy Target type")
	var reqDiskSizeFlag int64
	flagset.Int64Var(&reqDiskSizeFlag, "disk-size", 0, "Instance disk size in GiB")
	var reqInstanceTypeAuthorizedFlag bool
	flagset.BoolVar(&reqInstanceTypeAuthorizedFlag, "instance-type.authorized", false, "Requires authorization or publicly available")
	var reqInstanceTypeCpusFlag int64
	flagset.Int64Var(&reqInstanceTypeCpusFlag, "instance-type.cpus", 0, "CPU count")
	var reqInstanceTypeFamilyFlag string
	flagset.StringVar(&reqInstanceTypeFamilyFlag, "instance-type.family", "", "Instance type family")
	var reqInstanceTypeGpusFlag int64
	flagset.Int64Var(&reqInstanceTypeGpusFlag, "instance-type.gpus", 0, "GPU count")
	var reqInstanceTypeIDFlag string
	flagset.StringVar(&reqInstanceTypeIDFlag, "instance-type.id", "", "Instance type ID")
	var reqInstanceTypeMemoryFlag int64
	flagset.Int64Var(&reqInstanceTypeMemoryFlag, "instance-type.memory", 0, "Available memory")
	var reqInstanceTypeSizeFlag string
	flagset.StringVar(&reqInstanceTypeSizeFlag, "instance-type.size", "", "Instance type size")
	var reqIpv6EnabledFlag bool
	flagset.BoolVar(&reqIpv6EnabledFlag, "ipv6-enabled", false, "Enable IPv6. DEPRECATED: use `public-ip-assignments`.")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Instance name")
	var reqPublicIPAssignmentFlag string
	flagset.StringVar(&reqPublicIPAssignmentFlag, "public-ip-assignment", "", "")
	var reqSecurebootEnabledFlag bool
	flagset.BoolVar(&reqSecurebootEnabledFlag, "secureboot-enabled", false, "[Beta] Enable secure boot")
	var reqSSHKeyFingerprintFlag string
	flagset.StringVar(&reqSSHKeyFingerprintFlag, "ssh-key.fingerprint", "", "SSH key fingerprint")
	var reqSSHKeyNameFlag string
	flagset.StringVar(&reqSSHKeyNameFlag, "ssh-key.name", "", "SSH key name")
	var reqTemplateBootModeFlag string
	flagset.StringVar(&reqTemplateBootModeFlag, "template.boot-mode", "", "Boot mode (default: legacy)")
	var reqTemplateBuildFlag string
	flagset.StringVar(&reqTemplateBuildFlag, "template.build", "", "Template build")
	var reqTemplateChecksumFlag string
	flagset.StringVar(&reqTemplateChecksumFlag, "template.checksum", "", "Template MD5 checksum")
	var reqTemplateCreatedATFlag time.Time
	flagset.Time.TimeVar(&reqTemplateCreatedATFlag, "template.created-at", "", "Template creation date")
	var reqTemplateDefaultUserFlag string
	flagset.StringVar(&reqTemplateDefaultUserFlag, "template.default-user", "", "Template default user")
	var reqTemplateDescriptionFlag string
	flagset.StringVar(&reqTemplateDescriptionFlag, "template.description", "", "Template description")
	var reqTemplateFamilyFlag string
	flagset.StringVar(&reqTemplateFamilyFlag, "template.family", "", "Template family")
	var reqTemplateIDFlag string
	flagset.StringVar(&reqTemplateIDFlag, "template.id", "", "Template ID")
	var reqTemplateMaintainerFlag string
	flagset.StringVar(&reqTemplateMaintainerFlag, "template.maintainer", "", "Template maintainer")
	var reqTemplateNameFlag string
	flagset.StringVar(&reqTemplateNameFlag, "template.name", "", "Template name")
	var reqTemplatePasswordEnabledFlag bool
	flagset.BoolVar(&reqTemplatePasswordEnabledFlag, "template.password-enabled", false, "Enable password-based login")
	var reqTemplateSizeFlag int64
	flagset.Int64Var(&reqTemplateSizeFlag, "template.size", 0, "Template size")
	var reqTemplateSSHKeyEnabledFlag bool
	flagset.BoolVar(&reqTemplateSSHKeyEnabledFlag, "template.ssh-key-enabled", false, "Enable SSH key-based login")
	var reqTemplateURLFlag string
	flagset.StringVar(&reqTemplateURLFlag, "template.url", "", "Template source URL")
	var reqTemplateVersionFlag string
	flagset.StringVar(&reqTemplateVersionFlag, "template.version", "", "Template version")
	var reqTemplateVisibilityFlag string
	flagset.StringVar(&reqTemplateVisibilityFlag, "template.visibility", "", "Template visibility")
	var reqTpmEnabledFlag bool
	flagset.BoolVar(&reqTpmEnabledFlag, "tpm-enabled", false, "[Beta] Enable Trusted Platform Module (TPM)")
	var reqUserDataFlag string
	flagset.StringVar(&reqUserDataFlag, "user-data", "", "Instance Cloud-init user-data (base64 encoded)")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateInstanceRequest
	req.UserData = reqUserDataFlag
	req.TpmEnabled = reqTpmEnabledFlag
	if reqTemplateVisibilityFlag != "" {
		req.Template = &v3.Template{}
		req.Template.Visibility = reqTemplateVisibilityFlag
	}
	if reqTemplateVersionFlag != "" {

		req.Template.Version = reqTemplateVersionFlag
	}
	if reqTemplateURLFlag != "" {

		req.Template.URL = reqTemplateURLFlag
	}
	if reqTemplateSSHKeyEnabledFlag != false {

		req.Template.SSHKeyEnabled = reqTemplateSSHKeyEnabledFlag
	}
	if reqTemplateSizeFlag != 0 {

		req.Template.Size = reqTemplateSizeFlag
	}
	if reqTemplatePasswordEnabledFlag != false {

		req.Template.PasswordEnabled = reqTemplatePasswordEnabledFlag
	}
	if reqTemplateNameFlag != "" {

		req.Template.Name = reqTemplateNameFlag
	}
	if reqTemplateMaintainerFlag != "" {

		req.Template.Maintainer = reqTemplateMaintainerFlag
	}
	if v3.UUID(reqTemplateIDFlag) != "" {

		req.Template.ID = v3.UUID(reqTemplateIDFlag)
	}
	if reqTemplateFamilyFlag != "" {

		req.Template.Family = reqTemplateFamilyFlag
	}
	if reqTemplateDescriptionFlag != "" {

		req.Template.Description = reqTemplateDescriptionFlag
	}
	if reqTemplateDefaultUserFlag != "" {

		req.Template.DefaultUser = reqTemplateDefaultUserFlag
	}
	if reqTemplateCreatedATFlag != "" {

		req.Template.CreatedAT = reqTemplateCreatedATFlag
	}
	if reqTemplateChecksumFlag != "" {

		req.Template.Checksum = reqTemplateChecksumFlag
	}
	if reqTemplateBuildFlag != "" {

		req.Template.Build = reqTemplateBuildFlag
	}
	if reqTemplateBootModeFlag != "" {

		req.Template.BootMode = reqTemplateBootModeFlag
	}
	if reqSSHKeyNameFlag != "" {
		req.SSHKey = &v3.SSHKey{}
		req.SSHKey.Name = reqSSHKeyNameFlag
	}
	if reqSSHKeyFingerprintFlag != "" {

		req.SSHKey.Fingerprint = reqSSHKeyFingerprintFlag
	}
	req.SecurebootEnabled = reqSecurebootEnabledFlag
	req.PublicIPAssignment = reqPublicIPAssignmentFlag
	req.Name = reqNameFlag
	req.Ipv6Enabled = reqIpv6EnabledFlag
	if reqInstanceTypeSizeFlag != "" {
		req.InstanceTyp = &v3.InstanceType{}
		req.InstanceType.Size = reqInstanceTypeSizeFlag
	}
	if reqInstanceTypeMemoryFlag != 0 {

		req.InstanceType.Memory = reqInstanceTypeMemoryFlag
	}
	if v3.UUID(reqInstanceTypeIDFlag) != "" {

		req.InstanceType.ID = v3.UUID(reqInstanceTypeIDFlag)
	}
	if reqInstanceTypeGpusFlag != 0 {

		req.InstanceType.Gpus = reqInstanceTypeGpusFlag
	}
	if reqInstanceTypeFamilyFlag != "" {

		req.InstanceType.Family = reqInstanceTypeFamilyFlag
	}
	if reqInstanceTypeCpusFlag != 0 {

		req.InstanceType.Cpus = reqInstanceTypeCpusFlag
	}
	if reqInstanceTypeAuthorizedFlag != false {

		req.InstanceType.Authorized = reqInstanceTypeAuthorizedFlag
	}
	req.DiskSize = reqDiskSizeFlag
	if reqDeployTargetTypeFlag != "" {
		req.DeployTarget = &v3.DeployTarget{}
		req.DeployTarget.Type = reqDeployTargetTypeFlag
	}
	if reqDeployTargetNameFlag != "" {

		req.DeployTarget.Name = reqDeployTargetNameFlag
	}
	if v3.UUID(reqDeployTargetIDFlag) != "" {

		req.DeployTarget.ID = v3.UUID(reqDeployTargetIDFlag)
	}
	if reqDeployTargetDescriptionFlag != "" {

		req.DeployTarget.Description = reqDeployTargetDescriptionFlag
	}
	req.AutoStart = reqAutoStartFlag

	resp, err := client.CreateInstance(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListInstancePoolsCmd(client *v3.Client) {
	resp, err := client.ListInstancePools(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateInstancePoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-instance-pool", flag.ExitOnError)
	var reqDeployTargetDescriptionFlag string
	flagset.StringVar(&reqDeployTargetDescriptionFlag, "deploy-target.description", "", "Deploy Target description")
	var reqDeployTargetIDFlag string
	flagset.StringVar(&reqDeployTargetIDFlag, "deploy-target.id", "", "Deploy Target ID")
	var reqDeployTargetNameFlag string
	flagset.StringVar(&reqDeployTargetNameFlag, "deploy-target.name", "", "Deploy Target name")
	var reqDeployTargetTypeFlag string
	flagset.StringVar(&reqDeployTargetTypeFlag, "deploy-target.type", "", "Deploy Target type")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Instance Pool description")
	var reqDiskSizeFlag int64
	flagset.Int64Var(&reqDiskSizeFlag, "disk-size", 0, "Instances disk size in GiB")
	var reqInstancePrefixFlag string
	flagset.StringVar(&reqInstancePrefixFlag, "instance-prefix", "", "Prefix to apply to Instances names (default: pool)")
	var reqInstanceTypeAuthorizedFlag bool
	flagset.BoolVar(&reqInstanceTypeAuthorizedFlag, "instance-type.authorized", false, "Requires authorization or publicly available")
	var reqInstanceTypeCpusFlag int64
	flagset.Int64Var(&reqInstanceTypeCpusFlag, "instance-type.cpus", 0, "CPU count")
	var reqInstanceTypeFamilyFlag string
	flagset.StringVar(&reqInstanceTypeFamilyFlag, "instance-type.family", "", "Instance type family")
	var reqInstanceTypeGpusFlag int64
	flagset.Int64Var(&reqInstanceTypeGpusFlag, "instance-type.gpus", 0, "GPU count")
	var reqInstanceTypeIDFlag string
	flagset.StringVar(&reqInstanceTypeIDFlag, "instance-type.id", "", "Instance type ID")
	var reqInstanceTypeMemoryFlag int64
	flagset.Int64Var(&reqInstanceTypeMemoryFlag, "instance-type.memory", 0, "Available memory")
	var reqInstanceTypeSizeFlag string
	flagset.StringVar(&reqInstanceTypeSizeFlag, "instance-type.size", "", "Instance type size")
	var reqIpv6EnabledFlag bool
	flagset.BoolVar(&reqIpv6EnabledFlag, "ipv6-enabled", false, "Enable IPv6. DEPRECATED: use `public-ip-assignments`.")
	var reqMinAvailableFlag int64
	flagset.Int64Var(&reqMinAvailableFlag, "min-available", 0, "Minimum number of running Instances")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Instance Pool name")
	var reqPublicIPAssignmentFlag string
	flagset.StringVar(&reqPublicIPAssignmentFlag, "public-ip-assignment", "", "Determines public IP assignment of the Instances. Type `none` is final and can't be changed later on.")
	var reqSizeFlag int64
	flagset.Int64Var(&reqSizeFlag, "size", 0, "Number of Instances")
	var reqSSHKeyFingerprintFlag string
	flagset.StringVar(&reqSSHKeyFingerprintFlag, "ssh-key.fingerprint", "", "SSH key fingerprint")
	var reqSSHKeyNameFlag string
	flagset.StringVar(&reqSSHKeyNameFlag, "ssh-key.name", "", "SSH key name")
	var reqTemplateBootModeFlag string
	flagset.StringVar(&reqTemplateBootModeFlag, "template.boot-mode", "", "Boot mode (default: legacy)")
	var reqTemplateBuildFlag string
	flagset.StringVar(&reqTemplateBuildFlag, "template.build", "", "Template build")
	var reqTemplateChecksumFlag string
	flagset.StringVar(&reqTemplateChecksumFlag, "template.checksum", "", "Template MD5 checksum")
	var reqTemplateCreatedATFlag time.Time
	flagset.Time.TimeVar(&reqTemplateCreatedATFlag, "template.created-at", "", "Template creation date")
	var reqTemplateDefaultUserFlag string
	flagset.StringVar(&reqTemplateDefaultUserFlag, "template.default-user", "", "Template default user")
	var reqTemplateDescriptionFlag string
	flagset.StringVar(&reqTemplateDescriptionFlag, "template.description", "", "Template description")
	var reqTemplateFamilyFlag string
	flagset.StringVar(&reqTemplateFamilyFlag, "template.family", "", "Template family")
	var reqTemplateIDFlag string
	flagset.StringVar(&reqTemplateIDFlag, "template.id", "", "Template ID")
	var reqTemplateMaintainerFlag string
	flagset.StringVar(&reqTemplateMaintainerFlag, "template.maintainer", "", "Template maintainer")
	var reqTemplateNameFlag string
	flagset.StringVar(&reqTemplateNameFlag, "template.name", "", "Template name")
	var reqTemplatePasswordEnabledFlag bool
	flagset.BoolVar(&reqTemplatePasswordEnabledFlag, "template.password-enabled", false, "Enable password-based login")
	var reqTemplateSizeFlag int64
	flagset.Int64Var(&reqTemplateSizeFlag, "template.size", 0, "Template size")
	var reqTemplateSSHKeyEnabledFlag bool
	flagset.BoolVar(&reqTemplateSSHKeyEnabledFlag, "template.ssh-key-enabled", false, "Enable SSH key-based login")
	var reqTemplateURLFlag string
	flagset.StringVar(&reqTemplateURLFlag, "template.url", "", "Template source URL")
	var reqTemplateVersionFlag string
	flagset.StringVar(&reqTemplateVersionFlag, "template.version", "", "Template version")
	var reqTemplateVisibilityFlag string
	flagset.StringVar(&reqTemplateVisibilityFlag, "template.visibility", "", "Template visibility")
	var reqUserDataFlag string
	flagset.StringVar(&reqUserDataFlag, "user-data", "", "Instances Cloud-init user-data")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateInstancePoolRequest
	req.UserData = reqUserDataFlag
	if reqTemplateVisibilityFlag != "" {
		req.Template = &v3.Template{}
		req.Template.Visibility = reqTemplateVisibilityFlag
	}
	if reqTemplateVersionFlag != "" {

		req.Template.Version = reqTemplateVersionFlag
	}
	if reqTemplateURLFlag != "" {

		req.Template.URL = reqTemplateURLFlag
	}
	if reqTemplateSSHKeyEnabledFlag != false {

		req.Template.SSHKeyEnabled = reqTemplateSSHKeyEnabledFlag
	}
	if reqTemplateSizeFlag != 0 {

		req.Template.Size = reqTemplateSizeFlag
	}
	if reqTemplatePasswordEnabledFlag != false {

		req.Template.PasswordEnabled = reqTemplatePasswordEnabledFlag
	}
	if reqTemplateNameFlag != "" {

		req.Template.Name = reqTemplateNameFlag
	}
	if reqTemplateMaintainerFlag != "" {

		req.Template.Maintainer = reqTemplateMaintainerFlag
	}
	if v3.UUID(reqTemplateIDFlag) != "" {

		req.Template.ID = v3.UUID(reqTemplateIDFlag)
	}
	if reqTemplateFamilyFlag != "" {

		req.Template.Family = reqTemplateFamilyFlag
	}
	if reqTemplateDescriptionFlag != "" {

		req.Template.Description = reqTemplateDescriptionFlag
	}
	if reqTemplateDefaultUserFlag != "" {

		req.Template.DefaultUser = reqTemplateDefaultUserFlag
	}
	if reqTemplateCreatedATFlag != "" {

		req.Template.CreatedAT = reqTemplateCreatedATFlag
	}
	if reqTemplateChecksumFlag != "" {

		req.Template.Checksum = reqTemplateChecksumFlag
	}
	if reqTemplateBuildFlag != "" {

		req.Template.Build = reqTemplateBuildFlag
	}
	if reqTemplateBootModeFlag != "" {

		req.Template.BootMode = reqTemplateBootModeFlag
	}
	if reqSSHKeyNameFlag != "" {
		req.SSHKey = &v3.SSHKey{}
		req.SSHKey.Name = reqSSHKeyNameFlag
	}
	if reqSSHKeyFingerprintFlag != "" {

		req.SSHKey.Fingerprint = reqSSHKeyFingerprintFlag
	}
	req.Size = reqSizeFlag
	req.PublicIPAssignment = reqPublicIPAssignmentFlag
	req.Name = reqNameFlag
	req.MinAvailable = reqMinAvailableFlag
	req.Ipv6Enabled = reqIpv6EnabledFlag
	if reqInstanceTypeSizeFlag != "" {
		req.InstanceTyp = &v3.InstanceType{}
		req.InstanceType.Size = reqInstanceTypeSizeFlag
	}
	if reqInstanceTypeMemoryFlag != 0 {

		req.InstanceType.Memory = reqInstanceTypeMemoryFlag
	}
	if v3.UUID(reqInstanceTypeIDFlag) != "" {

		req.InstanceType.ID = v3.UUID(reqInstanceTypeIDFlag)
	}
	if reqInstanceTypeGpusFlag != 0 {

		req.InstanceType.Gpus = reqInstanceTypeGpusFlag
	}
	if reqInstanceTypeFamilyFlag != "" {

		req.InstanceType.Family = reqInstanceTypeFamilyFlag
	}
	if reqInstanceTypeCpusFlag != 0 {

		req.InstanceType.Cpus = reqInstanceTypeCpusFlag
	}
	if reqInstanceTypeAuthorizedFlag != false {

		req.InstanceType.Authorized = reqInstanceTypeAuthorizedFlag
	}
	req.InstancePrefix = reqInstancePrefixFlag
	req.DiskSize = reqDiskSizeFlag
	req.Description = reqDescriptionFlag
	if reqDeployTargetTypeFlag != "" {
		req.DeployTarget = &v3.DeployTarget{}
		req.DeployTarget.Type = reqDeployTargetTypeFlag
	}
	if reqDeployTargetNameFlag != "" {

		req.DeployTarget.Name = reqDeployTargetNameFlag
	}
	if v3.UUID(reqDeployTargetIDFlag) != "" {

		req.DeployTarget.ID = v3.UUID(reqDeployTargetIDFlag)
	}
	if reqDeployTargetDescriptionFlag != "" {

		req.DeployTarget.Description = reqDeployTargetDescriptionFlag
	}

	resp, err := client.CreateInstancePool(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteInstancePoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-instance-pool", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteInstancePool(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetInstancePoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-instance-pool", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetInstancePool(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateInstancePoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-instance-pool", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDeployTargetDescriptionFlag string
	flagset.StringVar(&reqDeployTargetDescriptionFlag, "deploy-target.description", "", "Deploy Target description")
	var reqDeployTargetIDFlag string
	flagset.StringVar(&reqDeployTargetIDFlag, "deploy-target.id", "", "Deploy Target ID")
	var reqDeployTargetNameFlag string
	flagset.StringVar(&reqDeployTargetNameFlag, "deploy-target.name", "", "Deploy Target name")
	var reqDeployTargetTypeFlag string
	flagset.StringVar(&reqDeployTargetTypeFlag, "deploy-target.type", "", "Deploy Target type")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Instance Pool description")
	var reqDiskSizeFlag int64
	flagset.Int64Var(&reqDiskSizeFlag, "disk-size", 0, "Instances disk size in GiB")
	var reqInstancePrefixFlag string
	flagset.StringVar(&reqInstancePrefixFlag, "instance-prefix", "", "Prefix to apply to Instances names (default: pool)")
	var reqInstanceTypeAuthorizedFlag bool
	flagset.BoolVar(&reqInstanceTypeAuthorizedFlag, "instance-type.authorized", false, "Requires authorization or publicly available")
	var reqInstanceTypeCpusFlag int64
	flagset.Int64Var(&reqInstanceTypeCpusFlag, "instance-type.cpus", 0, "CPU count")
	var reqInstanceTypeFamilyFlag string
	flagset.StringVar(&reqInstanceTypeFamilyFlag, "instance-type.family", "", "Instance type family")
	var reqInstanceTypeGpusFlag int64
	flagset.Int64Var(&reqInstanceTypeGpusFlag, "instance-type.gpus", 0, "GPU count")
	var reqInstanceTypeIDFlag string
	flagset.StringVar(&reqInstanceTypeIDFlag, "instance-type.id", "", "Instance type ID")
	var reqInstanceTypeMemoryFlag int64
	flagset.Int64Var(&reqInstanceTypeMemoryFlag, "instance-type.memory", 0, "Available memory")
	var reqInstanceTypeSizeFlag string
	flagset.StringVar(&reqInstanceTypeSizeFlag, "instance-type.size", "", "Instance type size")
	var reqIpv6EnabledFlag bool
	flagset.BoolVar(&reqIpv6EnabledFlag, "ipv6-enabled", false, "Enable IPv6. DEPRECATED: use `public-ip-assignments`.")
	var reqMinAvailableFlag int64
	flagset.Int64Var(&reqMinAvailableFlag, "min-available", 0, "Minimum number of running Instances")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Instance Pool name")
	var reqPublicIPAssignmentFlag string
	flagset.StringVar(&reqPublicIPAssignmentFlag, "public-ip-assignment", "", "Determines public IP assignment of the Instances.")
	var reqSSHKeyFingerprintFlag string
	flagset.StringVar(&reqSSHKeyFingerprintFlag, "ssh-key.fingerprint", "", "SSH key fingerprint")
	var reqSSHKeyNameFlag string
	flagset.StringVar(&reqSSHKeyNameFlag, "ssh-key.name", "", "SSH key name")
	var reqTemplateBootModeFlag string
	flagset.StringVar(&reqTemplateBootModeFlag, "template.boot-mode", "", "Boot mode (default: legacy)")
	var reqTemplateBuildFlag string
	flagset.StringVar(&reqTemplateBuildFlag, "template.build", "", "Template build")
	var reqTemplateChecksumFlag string
	flagset.StringVar(&reqTemplateChecksumFlag, "template.checksum", "", "Template MD5 checksum")
	var reqTemplateCreatedATFlag time.Time
	flagset.Time.TimeVar(&reqTemplateCreatedATFlag, "template.created-at", "", "Template creation date")
	var reqTemplateDefaultUserFlag string
	flagset.StringVar(&reqTemplateDefaultUserFlag, "template.default-user", "", "Template default user")
	var reqTemplateDescriptionFlag string
	flagset.StringVar(&reqTemplateDescriptionFlag, "template.description", "", "Template description")
	var reqTemplateFamilyFlag string
	flagset.StringVar(&reqTemplateFamilyFlag, "template.family", "", "Template family")
	var reqTemplateIDFlag string
	flagset.StringVar(&reqTemplateIDFlag, "template.id", "", "Template ID")
	var reqTemplateMaintainerFlag string
	flagset.StringVar(&reqTemplateMaintainerFlag, "template.maintainer", "", "Template maintainer")
	var reqTemplateNameFlag string
	flagset.StringVar(&reqTemplateNameFlag, "template.name", "", "Template name")
	var reqTemplatePasswordEnabledFlag bool
	flagset.BoolVar(&reqTemplatePasswordEnabledFlag, "template.password-enabled", false, "Enable password-based login")
	var reqTemplateSizeFlag int64
	flagset.Int64Var(&reqTemplateSizeFlag, "template.size", 0, "Template size")
	var reqTemplateSSHKeyEnabledFlag bool
	flagset.BoolVar(&reqTemplateSSHKeyEnabledFlag, "template.ssh-key-enabled", false, "Enable SSH key-based login")
	var reqTemplateURLFlag string
	flagset.StringVar(&reqTemplateURLFlag, "template.url", "", "Template source URL")
	var reqTemplateVersionFlag string
	flagset.StringVar(&reqTemplateVersionFlag, "template.version", "", "Template version")
	var reqTemplateVisibilityFlag string
	flagset.StringVar(&reqTemplateVisibilityFlag, "template.visibility", "", "Template visibility")
	var reqUserDataFlag string
	flagset.StringVar(&reqUserDataFlag, "user-data", "", "Instances Cloud-init user-data")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateInstancePoolRequest
	req.UserData = reqUserDataFlag
	if reqTemplateVisibilityFlag != "" {
		req.Template = &v3.Template{}
		req.Template.Visibility = reqTemplateVisibilityFlag
	}
	if reqTemplateVersionFlag != "" {

		req.Template.Version = reqTemplateVersionFlag
	}
	if reqTemplateURLFlag != "" {

		req.Template.URL = reqTemplateURLFlag
	}
	if reqTemplateSSHKeyEnabledFlag != false {

		req.Template.SSHKeyEnabled = reqTemplateSSHKeyEnabledFlag
	}
	if reqTemplateSizeFlag != 0 {

		req.Template.Size = reqTemplateSizeFlag
	}
	if reqTemplatePasswordEnabledFlag != false {

		req.Template.PasswordEnabled = reqTemplatePasswordEnabledFlag
	}
	if reqTemplateNameFlag != "" {

		req.Template.Name = reqTemplateNameFlag
	}
	if reqTemplateMaintainerFlag != "" {

		req.Template.Maintainer = reqTemplateMaintainerFlag
	}
	if v3.UUID(reqTemplateIDFlag) != "" {

		req.Template.ID = v3.UUID(reqTemplateIDFlag)
	}
	if reqTemplateFamilyFlag != "" {

		req.Template.Family = reqTemplateFamilyFlag
	}
	if reqTemplateDescriptionFlag != "" {

		req.Template.Description = reqTemplateDescriptionFlag
	}
	if reqTemplateDefaultUserFlag != "" {

		req.Template.DefaultUser = reqTemplateDefaultUserFlag
	}
	if reqTemplateCreatedATFlag != "" {

		req.Template.CreatedAT = reqTemplateCreatedATFlag
	}
	if reqTemplateChecksumFlag != "" {

		req.Template.Checksum = reqTemplateChecksumFlag
	}
	if reqTemplateBuildFlag != "" {

		req.Template.Build = reqTemplateBuildFlag
	}
	if reqTemplateBootModeFlag != "" {

		req.Template.BootMode = reqTemplateBootModeFlag
	}
	if reqSSHKeyNameFlag != "" {
		req.SSHKey = &v3.SSHKey{}
		req.SSHKey.Name = reqSSHKeyNameFlag
	}
	if reqSSHKeyFingerprintFlag != "" {

		req.SSHKey.Fingerprint = reqSSHKeyFingerprintFlag
	}
	req.PublicIPAssignment = reqPublicIPAssignmentFlag
	req.Name = reqNameFlag
	req.MinAvailable = reqMinAvailableFlag
	req.Ipv6Enabled = reqIpv6EnabledFlag
	if reqInstanceTypeSizeFlag != "" {
		req.InstanceTyp = &v3.InstanceType{}
		req.InstanceType.Size = reqInstanceTypeSizeFlag
	}
	if reqInstanceTypeMemoryFlag != 0 {

		req.InstanceType.Memory = reqInstanceTypeMemoryFlag
	}
	if v3.UUID(reqInstanceTypeIDFlag) != "" {

		req.InstanceType.ID = v3.UUID(reqInstanceTypeIDFlag)
	}
	if reqInstanceTypeGpusFlag != 0 {

		req.InstanceType.Gpus = reqInstanceTypeGpusFlag
	}
	if reqInstanceTypeFamilyFlag != "" {

		req.InstanceType.Family = reqInstanceTypeFamilyFlag
	}
	if reqInstanceTypeCpusFlag != 0 {

		req.InstanceType.Cpus = reqInstanceTypeCpusFlag
	}
	if reqInstanceTypeAuthorizedFlag != false {

		req.InstanceType.Authorized = reqInstanceTypeAuthorizedFlag
	}
	req.InstancePrefix = reqInstancePrefixFlag
	req.DiskSize = reqDiskSizeFlag
	req.Description = reqDescriptionFlag
	if reqDeployTargetTypeFlag != "" {
		req.DeployTarget = &v3.DeployTarget{}
		req.DeployTarget.Type = reqDeployTargetTypeFlag
	}
	if reqDeployTargetNameFlag != "" {

		req.DeployTarget.Name = reqDeployTargetNameFlag
	}
	if v3.UUID(reqDeployTargetIDFlag) != "" {

		req.DeployTarget.ID = v3.UUID(reqDeployTargetIDFlag)
	}
	if reqDeployTargetDescriptionFlag != "" {

		req.DeployTarget.Description = reqDeployTargetDescriptionFlag
	}

	resp, err := client.UpdateInstancePool(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetInstancePoolFieldCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-instance-pool-field", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var fieldFlag string
	flagset.StringVar(&fieldFlag, "Field", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ResetInstancePoolField(context.Background(), v3.UUID(idFlag), fieldFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func EvictInstancePoolMembersCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("evict-instance-pool-members", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.EvictInstancePoolMembersRequest

	resp, err := client.EvictInstancePoolMembers(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ScaleInstancePoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("scale-instance-pool", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqSizeFlag int64
	flagset.Int64Var(&reqSizeFlag, "size", 0, "Number of managed Instances")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ScaleInstancePoolRequest
	req.Size = reqSizeFlag

	resp, err := client.ScaleInstancePool(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListInstanceTypesCmd(client *v3.Client) {
	resp, err := client.ListInstanceTypes(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetInstanceTypeCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-instance-type", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetInstanceType(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-instance", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteInstance(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-instance", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetInstance(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-instance", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Instance name")
	var reqPublicIPAssignmentFlag string
	flagset.StringVar(&reqPublicIPAssignmentFlag, "public-ip-assignment", "", "")
	var reqUserDataFlag string
	flagset.StringVar(&reqUserDataFlag, "user-data", "", "Instance Cloud-init user-data (base64 encoded)")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateInstanceRequest
	req.UserData = reqUserDataFlag
	req.PublicIPAssignment = reqPublicIPAssignmentFlag
	req.Name = reqNameFlag

	resp, err := client.UpdateInstance(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetInstanceFieldCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-instance-field", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var fieldFlag string
	flagset.StringVar(&fieldFlag, "Field", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ResetInstanceField(context.Background(), v3.UUID(idFlag), fieldFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func AddInstanceProtectionCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("add-instance-protection", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.AddInstanceProtection(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateSnapshotCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-snapshot", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.CreateSnapshot(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func EnableTpmCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("enable-tpm", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.EnableTpm(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RevealInstancePasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reveal-instance-password", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RevealInstancePassword(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RebootInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reboot-instance", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RebootInstance(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RemoveInstanceProtectionCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("remove-instance-protection", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RemoveInstanceProtection(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-instance", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDiskSizeFlag int64
	flagset.Int64Var(&reqDiskSizeFlag, "disk-size", 0, "Instance disk size in GiB")
	var reqTemplateBootModeFlag string
	flagset.StringVar(&reqTemplateBootModeFlag, "template.boot-mode", "", "Boot mode (default: legacy)")
	var reqTemplateBuildFlag string
	flagset.StringVar(&reqTemplateBuildFlag, "template.build", "", "Template build")
	var reqTemplateChecksumFlag string
	flagset.StringVar(&reqTemplateChecksumFlag, "template.checksum", "", "Template MD5 checksum")
	var reqTemplateCreatedATFlag time.Time
	flagset.Time.TimeVar(&reqTemplateCreatedATFlag, "template.created-at", "", "Template creation date")
	var reqTemplateDefaultUserFlag string
	flagset.StringVar(&reqTemplateDefaultUserFlag, "template.default-user", "", "Template default user")
	var reqTemplateDescriptionFlag string
	flagset.StringVar(&reqTemplateDescriptionFlag, "template.description", "", "Template description")
	var reqTemplateFamilyFlag string
	flagset.StringVar(&reqTemplateFamilyFlag, "template.family", "", "Template family")
	var reqTemplateIDFlag string
	flagset.StringVar(&reqTemplateIDFlag, "template.id", "", "Template ID")
	var reqTemplateMaintainerFlag string
	flagset.StringVar(&reqTemplateMaintainerFlag, "template.maintainer", "", "Template maintainer")
	var reqTemplateNameFlag string
	flagset.StringVar(&reqTemplateNameFlag, "template.name", "", "Template name")
	var reqTemplatePasswordEnabledFlag bool
	flagset.BoolVar(&reqTemplatePasswordEnabledFlag, "template.password-enabled", false, "Enable password-based login")
	var reqTemplateSizeFlag int64
	flagset.Int64Var(&reqTemplateSizeFlag, "template.size", 0, "Template size")
	var reqTemplateSSHKeyEnabledFlag bool
	flagset.BoolVar(&reqTemplateSSHKeyEnabledFlag, "template.ssh-key-enabled", false, "Enable SSH key-based login")
	var reqTemplateURLFlag string
	flagset.StringVar(&reqTemplateURLFlag, "template.url", "", "Template source URL")
	var reqTemplateVersionFlag string
	flagset.StringVar(&reqTemplateVersionFlag, "template.version", "", "Template version")
	var reqTemplateVisibilityFlag string
	flagset.StringVar(&reqTemplateVisibilityFlag, "template.visibility", "", "Template visibility")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ResetInstanceRequest
	if reqTemplateVisibilityFlag != "" {
		req.Template = &v3.Template{}
		req.Template.Visibility = reqTemplateVisibilityFlag
	}
	if reqTemplateVersionFlag != "" {

		req.Template.Version = reqTemplateVersionFlag
	}
	if reqTemplateURLFlag != "" {

		req.Template.URL = reqTemplateURLFlag
	}
	if reqTemplateSSHKeyEnabledFlag != false {

		req.Template.SSHKeyEnabled = reqTemplateSSHKeyEnabledFlag
	}
	if reqTemplateSizeFlag != 0 {

		req.Template.Size = reqTemplateSizeFlag
	}
	if reqTemplatePasswordEnabledFlag != false {

		req.Template.PasswordEnabled = reqTemplatePasswordEnabledFlag
	}
	if reqTemplateNameFlag != "" {

		req.Template.Name = reqTemplateNameFlag
	}
	if reqTemplateMaintainerFlag != "" {

		req.Template.Maintainer = reqTemplateMaintainerFlag
	}
	if v3.UUID(reqTemplateIDFlag) != "" {

		req.Template.ID = v3.UUID(reqTemplateIDFlag)
	}
	if reqTemplateFamilyFlag != "" {

		req.Template.Family = reqTemplateFamilyFlag
	}
	if reqTemplateDescriptionFlag != "" {

		req.Template.Description = reqTemplateDescriptionFlag
	}
	if reqTemplateDefaultUserFlag != "" {

		req.Template.DefaultUser = reqTemplateDefaultUserFlag
	}
	if reqTemplateCreatedATFlag != "" {

		req.Template.CreatedAT = reqTemplateCreatedATFlag
	}
	if reqTemplateChecksumFlag != "" {

		req.Template.Checksum = reqTemplateChecksumFlag
	}
	if reqTemplateBuildFlag != "" {

		req.Template.Build = reqTemplateBuildFlag
	}
	if reqTemplateBootModeFlag != "" {

		req.Template.BootMode = reqTemplateBootModeFlag
	}
	req.DiskSize = reqDiskSizeFlag

	resp, err := client.ResetInstance(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetInstancePasswordCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-instance-password", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ResetInstancePassword(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResizeInstanceDiskCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("resize-instance-disk", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDiskSizeFlag int64
	flagset.Int64Var(&reqDiskSizeFlag, "disk-size", 0, "Instance disk size in GiB")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ResizeInstanceDiskRequest
	req.DiskSize = reqDiskSizeFlag

	resp, err := client.ResizeInstanceDisk(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ScaleInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("scale-instance", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqInstanceTypeAuthorizedFlag bool
	flagset.BoolVar(&reqInstanceTypeAuthorizedFlag, "instance-type.authorized", false, "Requires authorization or publicly available")
	var reqInstanceTypeCpusFlag int64
	flagset.Int64Var(&reqInstanceTypeCpusFlag, "instance-type.cpus", 0, "CPU count")
	var reqInstanceTypeFamilyFlag string
	flagset.StringVar(&reqInstanceTypeFamilyFlag, "instance-type.family", "", "Instance type family")
	var reqInstanceTypeGpusFlag int64
	flagset.Int64Var(&reqInstanceTypeGpusFlag, "instance-type.gpus", 0, "GPU count")
	var reqInstanceTypeIDFlag string
	flagset.StringVar(&reqInstanceTypeIDFlag, "instance-type.id", "", "Instance type ID")
	var reqInstanceTypeMemoryFlag int64
	flagset.Int64Var(&reqInstanceTypeMemoryFlag, "instance-type.memory", 0, "Available memory")
	var reqInstanceTypeSizeFlag string
	flagset.StringVar(&reqInstanceTypeSizeFlag, "instance-type.size", "", "Instance type size")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ScaleInstanceRequest
	if reqInstanceTypeSizeFlag != "" {
		req.InstanceTyp = &v3.InstanceType{}
		req.InstanceType.Size = reqInstanceTypeSizeFlag
	}
	if reqInstanceTypeMemoryFlag != 0 {

		req.InstanceType.Memory = reqInstanceTypeMemoryFlag
	}
	if v3.UUID(reqInstanceTypeIDFlag) != "" {

		req.InstanceType.ID = v3.UUID(reqInstanceTypeIDFlag)
	}
	if reqInstanceTypeGpusFlag != 0 {

		req.InstanceType.Gpus = reqInstanceTypeGpusFlag
	}
	if reqInstanceTypeFamilyFlag != "" {

		req.InstanceType.Family = reqInstanceTypeFamilyFlag
	}
	if reqInstanceTypeCpusFlag != 0 {

		req.InstanceType.Cpus = reqInstanceTypeCpusFlag
	}
	if reqInstanceTypeAuthorizedFlag != false {

		req.InstanceType.Authorized = reqInstanceTypeAuthorizedFlag
	}

	resp, err := client.ScaleInstance(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StartInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("start-instance", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqRescueProfileFlag string
	flagset.StringVar(&reqRescueProfileFlag, "rescue-profile", "", "Boot in Rescue Mode, using named profile (supported: netboot, netboot-efi)")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.StartInstanceRequest
	req.RescueProfile = reqRescueProfileFlag

	resp, err := client.StartInstance(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func StopInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("stop-instance", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.StopInstance(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RevertInstanceToSnapshotCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("revert-instance-to-snapshot", flag.ExitOnError)
	var instanceIDFlag string
	flagset.StringVar(&instanceIDFlag, "InstanceID", "", "")
	var reqIDFlag string
	flagset.StringVar(&reqIDFlag, "id", "", "Snapshot ID")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.RevertInstanceToSnapshotRequest
	req.ID = reqIDFlag

	resp, err := client.RevertInstanceToSnapshot(context.Background(), v3.UUID(instanceIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListLoadBalancersCmd(client *v3.Client) {
	resp, err := client.ListLoadBalancers(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateLoadBalancerCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-load-balancer", flag.ExitOnError)
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Load Balancer description")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Load Balancer name")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateLoadBalancerRequest
	req.Name = reqNameFlag
	req.Description = reqDescriptionFlag

	resp, err := client.CreateLoadBalancer(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteLoadBalancerCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-load-balancer", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteLoadBalancer(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetLoadBalancerCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-load-balancer", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetLoadBalancer(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateLoadBalancerCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-load-balancer", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Load Balancer description")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Load Balancer name")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateLoadBalancerRequest
	req.Name = reqNameFlag
	req.Description = reqDescriptionFlag

	resp, err := client.UpdateLoadBalancer(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func AddServiceToLoadBalancerCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("add-service-to-load-balancer", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Load Balancer Service description")
	var reqHealthcheckIntervalFlag int64
	flagset.Int64Var(&reqHealthcheckIntervalFlag, "healthcheck.interval", 0, "Healthcheck interval (default: 10). Must be greater than or equal to Timeout")
	var reqHealthcheckModeFlag string
	flagset.StringVar(&reqHealthcheckModeFlag, "healthcheck.mode", "", "Healthcheck mode")
	var reqHealthcheckPortFlag int64
	flagset.Int64Var(&reqHealthcheckPortFlag, "healthcheck.port", 0, "Healthcheck port")
	var reqHealthcheckRetriesFlag int64
	flagset.Int64Var(&reqHealthcheckRetriesFlag, "healthcheck.retries", 0, "Number of retries before considering a Service failed")
	var reqHealthcheckTimeoutFlag int64
	flagset.Int64Var(&reqHealthcheckTimeoutFlag, "healthcheck.timeout", 0, "Healthcheck timeout value (default: 2). Must be lower than or equal to Interval")
	var reqHealthcheckTlsSNIFlag string
	flagset.StringVar(&reqHealthcheckTlsSNIFlag, "healthcheck.tls-sni", "", "SNI domain for HTTPS healthchecks")
	var reqHealthcheckURIFlag string
	flagset.StringVar(&reqHealthcheckURIFlag, "healthcheck.uri", "", "An endpoint to use for the HTTP healthcheck, e.g. '/status'")
	var reqInstancePoolDeployTargetDescriptionFlag string
	flagset.StringVar(&reqInstancePoolDeployTargetDescriptionFlag, "instance-pool.deploy-target.description", "", "Deploy Target description")
	var reqInstancePoolDeployTargetIDFlag string
	flagset.StringVar(&reqInstancePoolDeployTargetIDFlag, "instance-pool.deploy-target.id", "", "Deploy Target ID")
	var reqInstancePoolDeployTargetNameFlag string
	flagset.StringVar(&reqInstancePoolDeployTargetNameFlag, "instance-pool.deploy-target.name", "", "Deploy Target name")
	var reqInstancePoolDeployTargetTypeFlag string
	flagset.StringVar(&reqInstancePoolDeployTargetTypeFlag, "instance-pool.deploy-target.type", "", "Deploy Target type")
	var reqInstancePoolDescriptionFlag string
	flagset.StringVar(&reqInstancePoolDescriptionFlag, "instance-pool.description", "", "Instance Pool description")
	var reqInstancePoolDiskSizeFlag int64
	flagset.Int64Var(&reqInstancePoolDiskSizeFlag, "instance-pool.disk-size", 0, "Instances disk size in GiB")
	var reqInstancePoolIDFlag string
	flagset.StringVar(&reqInstancePoolIDFlag, "instance-pool.id", "", "Instance Pool ID")
	var reqInstancePoolInstancePrefixFlag string
	flagset.StringVar(&reqInstancePoolInstancePrefixFlag, "instance-pool.instance-prefix", "", "The instances created by the Instance Pool will be prefixed with this value (default: pool)")
	var reqInstancePoolInstanceTypeAuthorizedFlag bool
	flagset.BoolVar(&reqInstancePoolInstanceTypeAuthorizedFlag, "instance-pool.instance-type.authorized", false, "Requires authorization or publicly available")
	var reqInstancePoolInstanceTypeCpusFlag int64
	flagset.Int64Var(&reqInstancePoolInstanceTypeCpusFlag, "instance-pool.instance-type.cpus", 0, "CPU count")
	var reqInstancePoolInstanceTypeFamilyFlag string
	flagset.StringVar(&reqInstancePoolInstanceTypeFamilyFlag, "instance-pool.instance-type.family", "", "Instance type family")
	var reqInstancePoolInstanceTypeGpusFlag int64
	flagset.Int64Var(&reqInstancePoolInstanceTypeGpusFlag, "instance-pool.instance-type.gpus", 0, "GPU count")
	var reqInstancePoolInstanceTypeIDFlag string
	flagset.StringVar(&reqInstancePoolInstanceTypeIDFlag, "instance-pool.instance-type.id", "", "Instance type ID")
	var reqInstancePoolInstanceTypeMemoryFlag int64
	flagset.Int64Var(&reqInstancePoolInstanceTypeMemoryFlag, "instance-pool.instance-type.memory", 0, "Available memory")
	var reqInstancePoolInstanceTypeSizeFlag string
	flagset.StringVar(&reqInstancePoolInstanceTypeSizeFlag, "instance-pool.instance-type.size", "", "Instance type size")
	var reqInstancePoolIpv6EnabledFlag bool
	flagset.BoolVar(&reqInstancePoolIpv6EnabledFlag, "instance-pool.ipv6-enabled", false, "Enable IPv6 for instances")
	var reqInstancePoolManagerIDFlag string
	flagset.StringVar(&reqInstancePoolManagerIDFlag, "instance-pool.manager.id", "", "Manager ID")
	var reqInstancePoolManagerTypeFlag string
	flagset.StringVar(&reqInstancePoolManagerTypeFlag, "instance-pool.manager.type", "", "Manager type")
	var reqInstancePoolMinAvailableFlag int64
	flagset.Int64Var(&reqInstancePoolMinAvailableFlag, "instance-pool.min-available", 0, "Minimum number of running instances")
	var reqInstancePoolNameFlag string
	flagset.StringVar(&reqInstancePoolNameFlag, "instance-pool.name", "", "Instance Pool name")
	var reqInstancePoolPublicIPAssignmentFlag string
	flagset.StringVar(&reqInstancePoolPublicIPAssignmentFlag, "instance-pool.public-ip-assignment", "", "")
	var reqInstancePoolSizeFlag int64
	flagset.Int64Var(&reqInstancePoolSizeFlag, "instance-pool.size", 0, "Number of instances")
	var reqInstancePoolSSHKeyFingerprintFlag string
	flagset.StringVar(&reqInstancePoolSSHKeyFingerprintFlag, "instance-pool.ssh-key.fingerprint", "", "SSH key fingerprint")
	var reqInstancePoolSSHKeyNameFlag string
	flagset.StringVar(&reqInstancePoolSSHKeyNameFlag, "instance-pool.ssh-key.name", "", "SSH key name")
	var reqInstancePoolStateFlag string
	flagset.StringVar(&reqInstancePoolStateFlag, "instance-pool.state", "", "Instance Pool state")
	var reqInstancePoolTemplateBootModeFlag string
	flagset.StringVar(&reqInstancePoolTemplateBootModeFlag, "instance-pool.template.boot-mode", "", "Boot mode (default: legacy)")
	var reqInstancePoolTemplateBuildFlag string
	flagset.StringVar(&reqInstancePoolTemplateBuildFlag, "instance-pool.template.build", "", "Template build")
	var reqInstancePoolTemplateChecksumFlag string
	flagset.StringVar(&reqInstancePoolTemplateChecksumFlag, "instance-pool.template.checksum", "", "Template MD5 checksum")
	var reqInstancePoolTemplateCreatedATFlag time.Time
	flagset.Time.TimeVar(&reqInstancePoolTemplateCreatedATFlag, "instance-pool.template.created-at", "", "Template creation date")
	var reqInstancePoolTemplateDefaultUserFlag string
	flagset.StringVar(&reqInstancePoolTemplateDefaultUserFlag, "instance-pool.template.default-user", "", "Template default user")
	var reqInstancePoolTemplateDescriptionFlag string
	flagset.StringVar(&reqInstancePoolTemplateDescriptionFlag, "instance-pool.template.description", "", "Template description")
	var reqInstancePoolTemplateFamilyFlag string
	flagset.StringVar(&reqInstancePoolTemplateFamilyFlag, "instance-pool.template.family", "", "Template family")
	var reqInstancePoolTemplateIDFlag string
	flagset.StringVar(&reqInstancePoolTemplateIDFlag, "instance-pool.template.id", "", "Template ID")
	var reqInstancePoolTemplateMaintainerFlag string
	flagset.StringVar(&reqInstancePoolTemplateMaintainerFlag, "instance-pool.template.maintainer", "", "Template maintainer")
	var reqInstancePoolTemplateNameFlag string
	flagset.StringVar(&reqInstancePoolTemplateNameFlag, "instance-pool.template.name", "", "Template name")
	var reqInstancePoolTemplatePasswordEnabledFlag bool
	flagset.BoolVar(&reqInstancePoolTemplatePasswordEnabledFlag, "instance-pool.template.password-enabled", false, "Enable password-based login")
	var reqInstancePoolTemplateSizeFlag int64
	flagset.Int64Var(&reqInstancePoolTemplateSizeFlag, "instance-pool.template.size", 0, "Template size")
	var reqInstancePoolTemplateSSHKeyEnabledFlag bool
	flagset.BoolVar(&reqInstancePoolTemplateSSHKeyEnabledFlag, "instance-pool.template.ssh-key-enabled", false, "Enable SSH key-based login")
	var reqInstancePoolTemplateURLFlag string
	flagset.StringVar(&reqInstancePoolTemplateURLFlag, "instance-pool.template.url", "", "Template source URL")
	var reqInstancePoolTemplateVersionFlag string
	flagset.StringVar(&reqInstancePoolTemplateVersionFlag, "instance-pool.template.version", "", "Template version")
	var reqInstancePoolTemplateVisibilityFlag string
	flagset.StringVar(&reqInstancePoolTemplateVisibilityFlag, "instance-pool.template.visibility", "", "Template visibility")
	var reqInstancePoolUserDataFlag string
	flagset.StringVar(&reqInstancePoolUserDataFlag, "instance-pool.user-data", "", "Instances Cloud-init user-data")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Load Balancer Service name")
	var reqPortFlag int64
	flagset.Int64Var(&reqPortFlag, "port", 0, "Port exposed on the Load Balancer's public IP")
	var reqProtocolFlag string
	flagset.StringVar(&reqProtocolFlag, "protocol", "", "Network traffic protocol")
	var reqStrategyFlag string
	flagset.StringVar(&reqStrategyFlag, "strategy", "", "Load balancing strategy")
	var reqTargetPortFlag int64
	flagset.Int64Var(&reqTargetPortFlag, "target-port", 0, "Port on which the network traffic will be forwarded to on the receiving instance")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.AddServiceToLoadBalancerRequest
	req.TargetPort = reqTargetPortFlag
	req.Strategy = reqStrategyFlag
	req.Protocol = reqProtocolFlag
	req.Port = reqPortFlag
	req.Name = reqNameFlag
	if reqInstancePoolUserDataFlag != "" {
		req.InstancePool = &v3.InstancePool{}
		req.InstancePool.UserData = reqInstancePoolUserDataFlag
	}
	if reqInstancePoolTemplateVisibilityFlag != "" {
		req.InstancePoolTemplate = &v3.Template{}
		req.InstancePoolTemplate.Visibility = reqInstancePoolTemplateVisibilityFlag
	}
	if reqInstancePoolTemplateVersionFlag != "" {

		req.InstancePoolTemplate.Version = reqInstancePoolTemplateVersionFlag
	}
	if reqInstancePoolTemplateURLFlag != "" {

		req.InstancePoolTemplate.URL = reqInstancePoolTemplateURLFlag
	}
	if reqInstancePoolTemplateSSHKeyEnabledFlag != false {

		req.InstancePoolTemplate.SSHKeyEnabled = reqInstancePoolTemplateSSHKeyEnabledFlag
	}
	if reqInstancePoolTemplateSizeFlag != 0 {

		req.InstancePoolTemplate.Size = reqInstancePoolTemplateSizeFlag
	}
	if reqInstancePoolTemplatePasswordEnabledFlag != false {

		req.InstancePoolTemplate.PasswordEnabled = reqInstancePoolTemplatePasswordEnabledFlag
	}
	if reqInstancePoolTemplateNameFlag != "" {

		req.InstancePoolTemplate.Name = reqInstancePoolTemplateNameFlag
	}
	if reqInstancePoolTemplateMaintainerFlag != "" {

		req.InstancePoolTemplate.Maintainer = reqInstancePoolTemplateMaintainerFlag
	}
	if v3.UUID(reqInstancePoolTemplateIDFlag) != "" {

		req.InstancePoolTemplate.ID = v3.UUID(reqInstancePoolTemplateIDFlag)
	}
	if reqInstancePoolTemplateFamilyFlag != "" {

		req.InstancePoolTemplate.Family = reqInstancePoolTemplateFamilyFlag
	}
	if reqInstancePoolTemplateDescriptionFlag != "" {

		req.InstancePoolTemplate.Description = reqInstancePoolTemplateDescriptionFlag
	}
	if reqInstancePoolTemplateDefaultUserFlag != "" {

		req.InstancePoolTemplate.DefaultUser = reqInstancePoolTemplateDefaultUserFlag
	}
	if reqInstancePoolTemplateCreatedATFlag != "" {

		req.InstancePoolTemplate.CreatedAT = reqInstancePoolTemplateCreatedATFlag
	}
	if reqInstancePoolTemplateChecksumFlag != "" {

		req.InstancePoolTemplate.Checksum = reqInstancePoolTemplateChecksumFlag
	}
	if reqInstancePoolTemplateBuildFlag != "" {

		req.InstancePoolTemplate.Build = reqInstancePoolTemplateBuildFlag
	}
	if reqInstancePoolTemplateBootModeFlag != "" {

		req.InstancePoolTemplate.BootMode = reqInstancePoolTemplateBootModeFlag
	}
	if reqInstancePoolStateFlag != "" {
		req.InstancePool = &v3.InstancePool{}
		req.InstancePool.State = reqInstancePoolStateFlag
	}
	if reqInstancePoolSSHKeyNameFlag != "" {
		req.InstancePoolSSHKey = &v3.SSHKey{}
		req.InstancePoolSSHKey.Name = reqInstancePoolSSHKeyNameFlag
	}
	if reqInstancePoolSSHKeyFingerprintFlag != "" {

		req.InstancePoolSSHKey.Fingerprint = reqInstancePoolSSHKeyFingerprintFlag
	}
	if reqInstancePoolSizeFlag != 0 {
		req.InstancePool = &v3.InstancePool{}
		req.InstancePool.Size = reqInstancePoolSizeFlag
	}
	if reqInstancePoolPublicIPAssignmentFlag != "" {

		req.InstancePool.PublicIPAssignment = reqInstancePoolPublicIPAssignmentFlag
	}
	if reqInstancePoolNameFlag != "" {

		req.InstancePool.Name = reqInstancePoolNameFlag
	}
	if reqInstancePoolMinAvailableFlag != 0 {

		req.InstancePool.MinAvailable = reqInstancePoolMinAvailableFlag
	}
	if reqInstancePoolManagerTypeFlag != "" {
		req.InstancePoolManager = &v3.Manager{}
		req.InstancePoolManager.Type = reqInstancePoolManagerTypeFlag
	}
	if v3.UUID(reqInstancePoolManagerIDFlag) != "" {

		req.InstancePoolManager.ID = v3.UUID(reqInstancePoolManagerIDFlag)
	}
	if reqInstancePoolIpv6EnabledFlag != false {
		req.InstancePoo = &v3.InstancePool{}
		req.InstancePool.Ipv6Enabled = reqInstancePoolIpv6EnabledFlag
	}
	if reqInstancePoolInstanceTypeSizeFlag != "" {
		req.InstancePoolInstanceTyp = &v3.InstanceType{}
		req.InstancePoolInstanceType.Size = reqInstancePoolInstanceTypeSizeFlag
	}
	if reqInstancePoolInstanceTypeMemoryFlag != 0 {

		req.InstancePoolInstanceType.Memory = reqInstancePoolInstanceTypeMemoryFlag
	}
	if v3.UUID(reqInstancePoolInstanceTypeIDFlag) != "" {

		req.InstancePoolInstanceType.ID = v3.UUID(reqInstancePoolInstanceTypeIDFlag)
	}
	if reqInstancePoolInstanceTypeGpusFlag != 0 {

		req.InstancePoolInstanceType.Gpus = reqInstancePoolInstanceTypeGpusFlag
	}
	if reqInstancePoolInstanceTypeFamilyFlag != "" {

		req.InstancePoolInstanceType.Family = reqInstancePoolInstanceTypeFamilyFlag
	}
	if reqInstancePoolInstanceTypeCpusFlag != 0 {

		req.InstancePoolInstanceType.Cpus = reqInstancePoolInstanceTypeCpusFlag
	}
	if reqInstancePoolInstanceTypeAuthorizedFlag != false {

		req.InstancePoolInstanceType.Authorized = reqInstancePoolInstanceTypeAuthorizedFlag
	}
	if reqInstancePoolInstancePrefixFlag != "" {
		req.InstancePool = &v3.InstancePool{}
		req.InstancePool.InstancePrefix = reqInstancePoolInstancePrefixFlag
	}
	if v3.UUID(reqInstancePoolIDFlag) != "" {

		req.InstancePool.ID = v3.UUID(reqInstancePoolIDFlag)
	}
	if reqInstancePoolDiskSizeFlag != 0 {

		req.InstancePool.DiskSize = reqInstancePoolDiskSizeFlag
	}
	if reqInstancePoolDescriptionFlag != "" {

		req.InstancePool.Description = reqInstancePoolDescriptionFlag
	}
	if reqInstancePoolDeployTargetTypeFlag != "" {
		req.InstancePoolDeployTarget = &v3.DeployTarget{}
		req.InstancePoolDeployTarget.Type = reqInstancePoolDeployTargetTypeFlag
	}
	if reqInstancePoolDeployTargetNameFlag != "" {

		req.InstancePoolDeployTarget.Name = reqInstancePoolDeployTargetNameFlag
	}
	if v3.UUID(reqInstancePoolDeployTargetIDFlag) != "" {

		req.InstancePoolDeployTarget.ID = v3.UUID(reqInstancePoolDeployTargetIDFlag)
	}
	if reqInstancePoolDeployTargetDescriptionFlag != "" {

		req.InstancePoolDeployTarget.Description = reqInstancePoolDeployTargetDescriptionFlag
	}
	if reqHealthcheckURIFlag != "" {
		req.Healthcheck = &v3.LoadBalancerServiceHealthcheck{}
		req.Healthcheck.URI = reqHealthcheckURIFlag
	}
	if reqHealthcheckTlsSNIFlag != "" {

		req.Healthcheck.TlsSNI = reqHealthcheckTlsSNIFlag
	}
	if reqHealthcheckTimeoutFlag != 0 {

		req.Healthcheck.Timeout = reqHealthcheckTimeoutFlag
	}
	if reqHealthcheckRetriesFlag != 0 {

		req.Healthcheck.Retries = reqHealthcheckRetriesFlag
	}
	if reqHealthcheckPortFlag != 0 {

		req.Healthcheck.Port = reqHealthcheckPortFlag
	}
	if reqHealthcheckModeFlag != "" {

		req.Healthcheck.Mode = reqHealthcheckModeFlag
	}
	if reqHealthcheckIntervalFlag != 0 {

		req.Healthcheck.Interval = reqHealthcheckIntervalFlag
	}
	req.Description = reqDescriptionFlag

	resp, err := client.AddServiceToLoadBalancer(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteLoadBalancerServiceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-load-balancer-service", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var serviceIDFlag string
	flagset.StringVar(&serviceIDFlag, "ServiceID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteLoadBalancerService(context.Background(), v3.UUID(idFlag), v3.UUID(serviceIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetLoadBalancerServiceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-load-balancer-service", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var serviceIDFlag string
	flagset.StringVar(&serviceIDFlag, "ServiceID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetLoadBalancerService(context.Background(), v3.UUID(idFlag), v3.UUID(serviceIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateLoadBalancerServiceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-load-balancer-service", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var serviceIDFlag string
	flagset.StringVar(&serviceIDFlag, "ServiceID", "", "")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Load Balancer Service description")
	var reqHealthcheckIntervalFlag int64
	flagset.Int64Var(&reqHealthcheckIntervalFlag, "healthcheck.interval", 0, "Healthcheck interval (default: 10). Must be greater than or equal to Timeout")
	var reqHealthcheckModeFlag string
	flagset.StringVar(&reqHealthcheckModeFlag, "healthcheck.mode", "", "Healthcheck mode")
	var reqHealthcheckPortFlag int64
	flagset.Int64Var(&reqHealthcheckPortFlag, "healthcheck.port", 0, "Healthcheck port")
	var reqHealthcheckRetriesFlag int64
	flagset.Int64Var(&reqHealthcheckRetriesFlag, "healthcheck.retries", 0, "Number of retries before considering a Service failed")
	var reqHealthcheckTimeoutFlag int64
	flagset.Int64Var(&reqHealthcheckTimeoutFlag, "healthcheck.timeout", 0, "Healthcheck timeout value (default: 2). Must be lower than or equal to Interval")
	var reqHealthcheckTlsSNIFlag string
	flagset.StringVar(&reqHealthcheckTlsSNIFlag, "healthcheck.tls-sni", "", "SNI domain for HTTPS healthchecks")
	var reqHealthcheckURIFlag string
	flagset.StringVar(&reqHealthcheckURIFlag, "healthcheck.uri", "", "An endpoint to use for the HTTP healthcheck, e.g. '/status'")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Load Balancer Service name")
	var reqPortFlag int64
	flagset.Int64Var(&reqPortFlag, "port", 0, "Port exposed on the Load Balancer's public IP")
	var reqProtocolFlag string
	flagset.StringVar(&reqProtocolFlag, "protocol", "", "Network traffic protocol")
	var reqStrategyFlag string
	flagset.StringVar(&reqStrategyFlag, "strategy", "", "Load balancing strategy")
	var reqTargetPortFlag int64
	flagset.Int64Var(&reqTargetPortFlag, "target-port", 0, "Port on which the network traffic will be forwarded to on the receiving instance")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateLoadBalancerServiceRequest
	req.TargetPort = reqTargetPortFlag
	req.Strategy = reqStrategyFlag
	req.Protocol = reqProtocolFlag
	req.Port = reqPortFlag
	req.Name = reqNameFlag
	if reqHealthcheckURIFlag != "" {
		req.Healthcheck = &v3.LoadBalancerServiceHealthcheck{}
		req.Healthcheck.URI = reqHealthcheckURIFlag
	}
	if reqHealthcheckTlsSNIFlag != "" {

		req.Healthcheck.TlsSNI = reqHealthcheckTlsSNIFlag
	}
	if reqHealthcheckTimeoutFlag != 0 {

		req.Healthcheck.Timeout = reqHealthcheckTimeoutFlag
	}
	if reqHealthcheckRetriesFlag != 0 {

		req.Healthcheck.Retries = reqHealthcheckRetriesFlag
	}
	if reqHealthcheckPortFlag != 0 {

		req.Healthcheck.Port = reqHealthcheckPortFlag
	}
	if reqHealthcheckModeFlag != "" {

		req.Healthcheck.Mode = reqHealthcheckModeFlag
	}
	if reqHealthcheckIntervalFlag != 0 {

		req.Healthcheck.Interval = reqHealthcheckIntervalFlag
	}
	req.Description = reqDescriptionFlag

	resp, err := client.UpdateLoadBalancerService(context.Background(), v3.UUID(idFlag), v3.UUID(serviceIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetLoadBalancerServiceFieldCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-load-balancer-service-field", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var serviceIDFlag string
	flagset.StringVar(&serviceIDFlag, "ServiceID", "", "")
	var fieldFlag string
	flagset.StringVar(&fieldFlag, "Field", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ResetLoadBalancerServiceField(context.Background(), v3.UUID(idFlag), v3.UUID(serviceIDFlag), fieldFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetLoadBalancerFieldCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-load-balancer-field", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var fieldFlag string
	flagset.StringVar(&fieldFlag, "Field", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ResetLoadBalancerField(context.Background(), v3.UUID(idFlag), fieldFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetOperationCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-operation", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetOperation(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetOrganizationCmd(client *v3.Client) {
	resp, err := client.GetOrganization(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListPrivateNetworksCmd(client *v3.Client) {
	resp, err := client.ListPrivateNetworks(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreatePrivateNetworkCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-private-network", flag.ExitOnError)
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Private Network description")
	var reqEndIPFlag net.IP
	flagset.Net.IPVar(&reqEndIPFlag, "end-ip", "", "Private Network end IP address")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Private Network name")
	var reqNetmaskFlag net.IP
	flagset.Net.IPVar(&reqNetmaskFlag, "netmask", "", "Private Network netmask")
	var reqStartIPFlag net.IP
	flagset.Net.IPVar(&reqStartIPFlag, "start-ip", "", "Private Network start IP address")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreatePrivateNetworkRequest
	req.StartIP = reqStartIPFlag
	req.Netmask = reqNetmaskFlag
	req.Name = reqNameFlag
	req.EndIP = reqEndIPFlag
	req.Description = reqDescriptionFlag

	resp, err := client.CreatePrivateNetwork(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeletePrivateNetworkCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-private-network", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeletePrivateNetwork(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetPrivateNetworkCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-private-network", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetPrivateNetwork(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdatePrivateNetworkCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-private-network", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Private Network description")
	var reqEndIPFlag net.IP
	flagset.Net.IPVar(&reqEndIPFlag, "end-ip", "", "Private Network end IP address")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Private Network name")
	var reqNetmaskFlag net.IP
	flagset.Net.IPVar(&reqNetmaskFlag, "netmask", "", "Private Network netmask")
	var reqStartIPFlag net.IP
	flagset.Net.IPVar(&reqStartIPFlag, "start-ip", "", "Private Network start IP address")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdatePrivateNetworkRequest
	req.StartIP = reqStartIPFlag
	req.Netmask = reqNetmaskFlag
	req.Name = reqNameFlag
	req.EndIP = reqEndIPFlag
	req.Description = reqDescriptionFlag

	resp, err := client.UpdatePrivateNetwork(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetPrivateNetworkFieldCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-private-network-field", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var fieldFlag string
	flagset.StringVar(&fieldFlag, "Field", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ResetPrivateNetworkField(context.Background(), v3.UUID(idFlag), fieldFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func AttachInstanceToPrivateNetworkCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("attach-instance-to-private-network", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqInstanceIDFlag string
	flagset.StringVar(&reqInstanceIDFlag, "instance.id", "", "Instance ID")
	var reqIPFlag net.IP
	flagset.Net.IPVar(&reqIPFlag, "ip", "", "Static IP address lease for the corresponding network interface")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.AttachInstanceToPrivateNetworkRequest
	req.IP = reqIPFlag
	if v3.UUID(reqInstanceIDFlag) != "" {
		req.Instance = &v3.AttachInstanceToPrivateNetworkRequestInstance{}
		req.Instance.ID = v3.UUID(reqInstanceIDFlag)
	}

	resp, err := client.AttachInstanceToPrivateNetwork(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DetachInstanceFromPrivateNetworkCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("detach-instance-from-private-network", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqInstanceCreatedATFlag time.Time
	flagset.Time.TimeVar(&reqInstanceCreatedATFlag, "instance.created-at", "", "Instance creation date")
	var reqInstanceDeployTargetDescriptionFlag string
	flagset.StringVar(&reqInstanceDeployTargetDescriptionFlag, "instance.deploy-target.description", "", "Deploy Target description")
	var reqInstanceDeployTargetIDFlag string
	flagset.StringVar(&reqInstanceDeployTargetIDFlag, "instance.deploy-target.id", "", "Deploy Target ID")
	var reqInstanceDeployTargetNameFlag string
	flagset.StringVar(&reqInstanceDeployTargetNameFlag, "instance.deploy-target.name", "", "Deploy Target name")
	var reqInstanceDeployTargetTypeFlag string
	flagset.StringVar(&reqInstanceDeployTargetTypeFlag, "instance.deploy-target.type", "", "Deploy Target type")
	var reqInstanceDiskSizeFlag int64
	flagset.Int64Var(&reqInstanceDiskSizeFlag, "instance.disk-size", 0, "Instance disk size in GiB")
	var reqInstanceIDFlag string
	flagset.StringVar(&reqInstanceIDFlag, "instance.id", "", "Instance ID")
	var reqInstanceInstanceTypeAuthorizedFlag bool
	flagset.BoolVar(&reqInstanceInstanceTypeAuthorizedFlag, "instance.instance-type.authorized", false, "Requires authorization or publicly available")
	var reqInstanceInstanceTypeCpusFlag int64
	flagset.Int64Var(&reqInstanceInstanceTypeCpusFlag, "instance.instance-type.cpus", 0, "CPU count")
	var reqInstanceInstanceTypeFamilyFlag string
	flagset.StringVar(&reqInstanceInstanceTypeFamilyFlag, "instance.instance-type.family", "", "Instance type family")
	var reqInstanceInstanceTypeGpusFlag int64
	flagset.Int64Var(&reqInstanceInstanceTypeGpusFlag, "instance.instance-type.gpus", 0, "GPU count")
	var reqInstanceInstanceTypeIDFlag string
	flagset.StringVar(&reqInstanceInstanceTypeIDFlag, "instance.instance-type.id", "", "Instance type ID")
	var reqInstanceInstanceTypeMemoryFlag int64
	flagset.Int64Var(&reqInstanceInstanceTypeMemoryFlag, "instance.instance-type.memory", 0, "Available memory")
	var reqInstanceInstanceTypeSizeFlag string
	flagset.StringVar(&reqInstanceInstanceTypeSizeFlag, "instance.instance-type.size", "", "Instance type size")
	var reqInstanceIpv6AddressFlag string
	flagset.StringVar(&reqInstanceIpv6AddressFlag, "instance.ipv6-address", "", "Instance IPv6 address")
	var reqInstanceMACAddressFlag string
	flagset.StringVar(&reqInstanceMACAddressFlag, "instance.mac-address", "", "Instance MAC address")
	var reqInstanceManagerIDFlag string
	flagset.StringVar(&reqInstanceManagerIDFlag, "instance.manager.id", "", "Manager ID")
	var reqInstanceManagerTypeFlag string
	flagset.StringVar(&reqInstanceManagerTypeFlag, "instance.manager.type", "", "Manager type")
	var reqInstanceNameFlag string
	flagset.StringVar(&reqInstanceNameFlag, "instance.name", "", "Instance name")
	var reqInstancePublicIPFlag net.IP
	flagset.Net.IPVar(&reqInstancePublicIPFlag, "instance.public-ip", "", "Instance public IPv4 address")
	var reqInstancePublicIPAssignmentFlag string
	flagset.StringVar(&reqInstancePublicIPAssignmentFlag, "instance.public-ip-assignment", "", "")
	var reqInstanceSecurebootEnabledFlag bool
	flagset.BoolVar(&reqInstanceSecurebootEnabledFlag, "instance.secureboot-enabled", false, "Indicates if the instance has secure boot enabled")
	var reqInstanceSSHKeyFingerprintFlag string
	flagset.StringVar(&reqInstanceSSHKeyFingerprintFlag, "instance.ssh-key.fingerprint", "", "SSH key fingerprint")
	var reqInstanceSSHKeyNameFlag string
	flagset.StringVar(&reqInstanceSSHKeyNameFlag, "instance.ssh-key.name", "", "SSH key name")
	var reqInstanceStateFlag string
	flagset.StringVar(&reqInstanceStateFlag, "instance.state", "", "")
	var reqInstanceTemplateBootModeFlag string
	flagset.StringVar(&reqInstanceTemplateBootModeFlag, "instance.template.boot-mode", "", "Boot mode (default: legacy)")
	var reqInstanceTemplateBuildFlag string
	flagset.StringVar(&reqInstanceTemplateBuildFlag, "instance.template.build", "", "Template build")
	var reqInstanceTemplateChecksumFlag string
	flagset.StringVar(&reqInstanceTemplateChecksumFlag, "instance.template.checksum", "", "Template MD5 checksum")
	var reqInstanceTemplateCreatedATFlag time.Time
	flagset.Time.TimeVar(&reqInstanceTemplateCreatedATFlag, "instance.template.created-at", "", "Template creation date")
	var reqInstanceTemplateDefaultUserFlag string
	flagset.StringVar(&reqInstanceTemplateDefaultUserFlag, "instance.template.default-user", "", "Template default user")
	var reqInstanceTemplateDescriptionFlag string
	flagset.StringVar(&reqInstanceTemplateDescriptionFlag, "instance.template.description", "", "Template description")
	var reqInstanceTemplateFamilyFlag string
	flagset.StringVar(&reqInstanceTemplateFamilyFlag, "instance.template.family", "", "Template family")
	var reqInstanceTemplateIDFlag string
	flagset.StringVar(&reqInstanceTemplateIDFlag, "instance.template.id", "", "Template ID")
	var reqInstanceTemplateMaintainerFlag string
	flagset.StringVar(&reqInstanceTemplateMaintainerFlag, "instance.template.maintainer", "", "Template maintainer")
	var reqInstanceTemplateNameFlag string
	flagset.StringVar(&reqInstanceTemplateNameFlag, "instance.template.name", "", "Template name")
	var reqInstanceTemplatePasswordEnabledFlag bool
	flagset.BoolVar(&reqInstanceTemplatePasswordEnabledFlag, "instance.template.password-enabled", false, "Enable password-based login")
	var reqInstanceTemplateSizeFlag int64
	flagset.Int64Var(&reqInstanceTemplateSizeFlag, "instance.template.size", 0, "Template size")
	var reqInstanceTemplateSSHKeyEnabledFlag bool
	flagset.BoolVar(&reqInstanceTemplateSSHKeyEnabledFlag, "instance.template.ssh-key-enabled", false, "Enable SSH key-based login")
	var reqInstanceTemplateURLFlag string
	flagset.StringVar(&reqInstanceTemplateURLFlag, "instance.template.url", "", "Template source URL")
	var reqInstanceTemplateVersionFlag string
	flagset.StringVar(&reqInstanceTemplateVersionFlag, "instance.template.version", "", "Template version")
	var reqInstanceTemplateVisibilityFlag string
	flagset.StringVar(&reqInstanceTemplateVisibilityFlag, "instance.template.visibility", "", "Template visibility")
	var reqInstanceTpmEnabledFlag bool
	flagset.BoolVar(&reqInstanceTpmEnabledFlag, "instance.tpm-enabled", false, "Indicates if the instance has tpm enabled")
	var reqInstanceUserDataFlag string
	flagset.StringVar(&reqInstanceUserDataFlag, "instance.user-data", "", "Instance Cloud-init user-data (base64 encoded)")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.DetachInstanceFromPrivateNetworkRequest
	if reqInstanceUserDataFlag != "" {
		req.Instanc = &v3.Instance{}
		req.Instance.UserData = reqInstanceUserDataFlag
	}
	if reqInstanceTpmEnabledFlag != false {

		req.Instance.TpmEnabled = reqInstanceTpmEnabledFlag
	}
	if reqInstanceTemplateVisibilityFlag != "" {
		req.InstanceTemplate = &v3.Template{}
		req.InstanceTemplate.Visibility = reqInstanceTemplateVisibilityFlag
	}
	if reqInstanceTemplateVersionFlag != "" {

		req.InstanceTemplate.Version = reqInstanceTemplateVersionFlag
	}
	if reqInstanceTemplateURLFlag != "" {

		req.InstanceTemplate.URL = reqInstanceTemplateURLFlag
	}
	if reqInstanceTemplateSSHKeyEnabledFlag != false {

		req.InstanceTemplate.SSHKeyEnabled = reqInstanceTemplateSSHKeyEnabledFlag
	}
	if reqInstanceTemplateSizeFlag != 0 {

		req.InstanceTemplate.Size = reqInstanceTemplateSizeFlag
	}
	if reqInstanceTemplatePasswordEnabledFlag != false {

		req.InstanceTemplate.PasswordEnabled = reqInstanceTemplatePasswordEnabledFlag
	}
	if reqInstanceTemplateNameFlag != "" {

		req.InstanceTemplate.Name = reqInstanceTemplateNameFlag
	}
	if reqInstanceTemplateMaintainerFlag != "" {

		req.InstanceTemplate.Maintainer = reqInstanceTemplateMaintainerFlag
	}
	if v3.UUID(reqInstanceTemplateIDFlag) != "" {

		req.InstanceTemplate.ID = v3.UUID(reqInstanceTemplateIDFlag)
	}
	if reqInstanceTemplateFamilyFlag != "" {

		req.InstanceTemplate.Family = reqInstanceTemplateFamilyFlag
	}
	if reqInstanceTemplateDescriptionFlag != "" {

		req.InstanceTemplate.Description = reqInstanceTemplateDescriptionFlag
	}
	if reqInstanceTemplateDefaultUserFlag != "" {

		req.InstanceTemplate.DefaultUser = reqInstanceTemplateDefaultUserFlag
	}
	if reqInstanceTemplateCreatedATFlag != "" {

		req.InstanceTemplate.CreatedAT = reqInstanceTemplateCreatedATFlag
	}
	if reqInstanceTemplateChecksumFlag != "" {

		req.InstanceTemplate.Checksum = reqInstanceTemplateChecksumFlag
	}
	if reqInstanceTemplateBuildFlag != "" {

		req.InstanceTemplate.Build = reqInstanceTemplateBuildFlag
	}
	if reqInstanceTemplateBootModeFlag != "" {

		req.InstanceTemplate.BootMode = reqInstanceTemplateBootModeFlag
	}
	if reqInstanceStateFlag != "" {
		req.Instanc = &v3.Instance{}
		req.Instance.State = reqInstanceStateFlag
	}
	if reqInstanceSSHKeyNameFlag != "" {
		req.InstanceSSHKey = &v3.SSHKey{}
		req.InstanceSSHKey.Name = reqInstanceSSHKeyNameFlag
	}
	if reqInstanceSSHKeyFingerprintFlag != "" {

		req.InstanceSSHKey.Fingerprint = reqInstanceSSHKeyFingerprintFlag
	}
	if reqInstanceSecurebootEnabledFlag != false {
		req.Ins = &v3.Instance{}
		req.Instance.SecurebootEnabled = reqInstanceSecurebootEnabledFlag
	}
	if reqInstancePublicIPAssignmentFlag != "" {

		req.Instance.PublicIPAssignment = reqInstancePublicIPAssignmentFlag
	}
	if reqInstancePublicIPFlag != "" {

		req.Instance.PublicIP = reqInstancePublicIPFlag
	}
	if reqInstanceNameFlag != "" {

		req.Instance.Name = reqInstanceNameFlag
	}
	if reqInstanceManagerTypeFlag != "" {
		req.InstanceManager = &v3.Manager{}
		req.InstanceManager.Type = reqInstanceManagerTypeFlag
	}
	if v3.UUID(reqInstanceManagerIDFlag) != "" {

		req.InstanceManager.ID = v3.UUID(reqInstanceManagerIDFlag)
	}
	if reqInstanceMACAddressFlag != "" {
		req.Instanc = &v3.Instance{}
		req.Instance.MACAddress = reqInstanceMACAddressFlag
	}
	if reqInstanceIpv6AddressFlag != "" {

		req.Instance.Ipv6Address = reqInstanceIpv6AddressFlag
	}
	if reqInstanceInstanceTypeSizeFlag != "" {
		req.InstanceInstanceTyp = &v3.InstanceType{}
		req.InstanceInstanceType.Size = reqInstanceInstanceTypeSizeFlag
	}
	if reqInstanceInstanceTypeMemoryFlag != 0 {

		req.InstanceInstanceType.Memory = reqInstanceInstanceTypeMemoryFlag
	}
	if v3.UUID(reqInstanceInstanceTypeIDFlag) != "" {

		req.InstanceInstanceType.ID = v3.UUID(reqInstanceInstanceTypeIDFlag)
	}
	if reqInstanceInstanceTypeGpusFlag != 0 {

		req.InstanceInstanceType.Gpus = reqInstanceInstanceTypeGpusFlag
	}
	if reqInstanceInstanceTypeFamilyFlag != "" {

		req.InstanceInstanceType.Family = reqInstanceInstanceTypeFamilyFlag
	}
	if reqInstanceInstanceTypeCpusFlag != 0 {

		req.InstanceInstanceType.Cpus = reqInstanceInstanceTypeCpusFlag
	}
	if reqInstanceInstanceTypeAuthorizedFlag != false {

		req.InstanceInstanceType.Authorized = reqInstanceInstanceTypeAuthorizedFlag
	}
	if v3.UUID(reqInstanceIDFlag) != "" {
		req.Instance = &v3.Instance{}
		req.Instance.ID = v3.UUID(reqInstanceIDFlag)
	}
	if reqInstanceDiskSizeFlag != 0 {

		req.Instance.DiskSize = reqInstanceDiskSizeFlag
	}
	if reqInstanceDeployTargetTypeFlag != "" {
		req.InstanceDeployTarget = &v3.DeployTarget{}
		req.InstanceDeployTarget.Type = reqInstanceDeployTargetTypeFlag
	}
	if reqInstanceDeployTargetNameFlag != "" {

		req.InstanceDeployTarget.Name = reqInstanceDeployTargetNameFlag
	}
	if v3.UUID(reqInstanceDeployTargetIDFlag) != "" {

		req.InstanceDeployTarget.ID = v3.UUID(reqInstanceDeployTargetIDFlag)
	}
	if reqInstanceDeployTargetDescriptionFlag != "" {

		req.InstanceDeployTarget.Description = reqInstanceDeployTargetDescriptionFlag
	}
	if reqInstanceCreatedATFlag != "" {
		req.Instanc = &v3.Instance{}
		req.Instance.CreatedAT = reqInstanceCreatedATFlag
	}

	resp, err := client.DetachInstanceFromPrivateNetwork(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdatePrivateNetworkInstanceIPCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-private-network-instance-ip", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqInstanceIDFlag string
	flagset.StringVar(&reqInstanceIDFlag, "instance.id", "", "Instance ID")
	var reqIPFlag net.IP
	flagset.Net.IPVar(&reqIPFlag, "ip", "", "Static IP address lease for the corresponding network interface")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdatePrivateNetworkInstanceIPRequest
	req.IP = reqIPFlag
	if v3.UUID(reqInstanceIDFlag) != "" {
		req.Instance = &v3.UpdatePrivateNetworkInstanceIPRequestInstance{}
		req.Instance.ID = v3.UUID(reqInstanceIDFlag)
	}

	resp, err := client.UpdatePrivateNetworkInstanceIP(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListQuotasCmd(client *v3.Client) {
	resp, err := client.ListQuotas(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetQuotaCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-quota", flag.ExitOnError)
	var entityFlag string
	flagset.StringVar(&entityFlag, "Entity", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetQuota(context.Background(), entityFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteReverseDNSElasticIPCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-reverse-dns-elastic-ip", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteReverseDNSElasticIP(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetReverseDNSElasticIPCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-reverse-dns-elastic-ip", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetReverseDNSElasticIP(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateReverseDNSElasticIPCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-reverse-dns-elastic-ip", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDomainNameFlag string
	flagset.StringVar(&reqDomainNameFlag, "domain-name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateReverseDNSElasticIPRequest
	req.DomainName = reqDomainNameFlag

	resp, err := client.UpdateReverseDNSElasticIP(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteReverseDNSInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-reverse-dns-instance", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteReverseDNSInstance(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetReverseDNSInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-reverse-dns-instance", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetReverseDNSInstance(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateReverseDNSInstanceCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-reverse-dns-instance", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDomainNameFlag string
	flagset.StringVar(&reqDomainNameFlag, "domain-name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateReverseDNSInstanceRequest
	req.DomainName = reqDomainNameFlag

	resp, err := client.UpdateReverseDNSInstance(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListSecurityGroupsCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("list-security-groups", flag.ExitOnError)
	var visibilityFlag string
	flagset.StringVar(&visibilityFlag, "Visibility", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ListSecurityGroups(context.Background(), visibilityFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateSecurityGroupCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-security-group", flag.ExitOnError)
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Security Group description")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Security Group name")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateSecurityGroupRequest
	req.Name = reqNameFlag
	req.Description = reqDescriptionFlag

	resp, err := client.CreateSecurityGroup(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteSecurityGroupCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-security-group", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteSecurityGroup(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetSecurityGroupCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-security-group", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetSecurityGroup(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func AddRuleToSecurityGroupCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("add-rule-to-security-group", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Security Group rule description")
	var reqEndPortFlag int64
	flagset.Int64Var(&reqEndPortFlag, "end-port", 0, "End port of the range")
	var reqFlowDirectionFlag string
	flagset.StringVar(&reqFlowDirectionFlag, "flow-direction", "", "Network flow direction to match")
	var reqICMPCodeFlag int64
	flagset.Int64Var(&reqICMPCodeFlag, "icmp.code", 0, "")
	var reqICMPTypeFlag int64
	flagset.Int64Var(&reqICMPTypeFlag, "icmp.type", 0, "")
	var reqNetworkFlag string
	flagset.StringVar(&reqNetworkFlag, "network", "", "CIDR-formatted network allowed")
	var reqProtocolFlag string
	flagset.StringVar(&reqProtocolFlag, "protocol", "", "Network protocol")
	var reqSecurityGroupIDFlag string
	flagset.StringVar(&reqSecurityGroupIDFlag, "security-group.id", "", "Security Group ID")
	var reqSecurityGroupNameFlag string
	flagset.StringVar(&reqSecurityGroupNameFlag, "security-group.name", "", "Security Group name")
	var reqSecurityGroupVisibilityFlag string
	flagset.StringVar(&reqSecurityGroupVisibilityFlag, "security-group.visibility", "", "Whether this points to a public security group. This is only valid when in the context of                    a rule addition which uses a public security group as a source or destination.")
	var reqStartPortFlag int64
	flagset.Int64Var(&reqStartPortFlag, "start-port", 0, "Start port of the range")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.AddRuleToSecurityGroupRequest
	req.StartPort = reqStartPortFlag
	if reqSecurityGroupVisibilityFlag != "" {
		req.SecurityGroup = &v3.SecurityGroupResource{}
		req.SecurityGroup.Visibility = reqSecurityGroupVisibilityFlag
	}
	if reqSecurityGroupNameFlag != "" {

		req.SecurityGroup.Name = reqSecurityGroupNameFlag
	}
	if v3.UUID(reqSecurityGroupIDFlag) != "" {

		req.SecurityGroup.ID = v3.UUID(reqSecurityGroupIDFlag)
	}
	req.Protocol = reqProtocolFlag
	req.Network = reqNetworkFlag
	if reqICMPTypeFlag != 0 {
		req.ICMP = &v3.AddRuleToSecurityGroupRequestICMP{}
		req.ICMP.Type = reqICMPTypeFlag
	}
	if reqICMPCodeFlag != 0 {

		req.ICMP.Code = reqICMPCodeFlag
	}
	req.FlowDirection = reqFlowDirectionFlag
	req.EndPort = reqEndPortFlag
	req.Description = reqDescriptionFlag

	resp, err := client.AddRuleToSecurityGroup(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteRuleFromSecurityGroupCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-rule-from-security-group", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var ruleIDFlag string
	flagset.StringVar(&ruleIDFlag, "RuleID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteRuleFromSecurityGroup(context.Background(), v3.UUID(idFlag), v3.UUID(ruleIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func AddExternalSourceToSecurityGroupCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("add-external-source-to-security-group", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqCidrFlag string
	flagset.StringVar(&reqCidrFlag, "cidr", "", "CIDR-formatted network to add")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.AddExternalSourceToSecurityGroupRequest
	req.Cidr = reqCidrFlag

	resp, err := client.AddExternalSourceToSecurityGroup(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func AttachInstanceToSecurityGroupCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("attach-instance-to-security-group", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqInstanceCreatedATFlag time.Time
	flagset.Time.TimeVar(&reqInstanceCreatedATFlag, "instance.created-at", "", "Instance creation date")
	var reqInstanceDeployTargetDescriptionFlag string
	flagset.StringVar(&reqInstanceDeployTargetDescriptionFlag, "instance.deploy-target.description", "", "Deploy Target description")
	var reqInstanceDeployTargetIDFlag string
	flagset.StringVar(&reqInstanceDeployTargetIDFlag, "instance.deploy-target.id", "", "Deploy Target ID")
	var reqInstanceDeployTargetNameFlag string
	flagset.StringVar(&reqInstanceDeployTargetNameFlag, "instance.deploy-target.name", "", "Deploy Target name")
	var reqInstanceDeployTargetTypeFlag string
	flagset.StringVar(&reqInstanceDeployTargetTypeFlag, "instance.deploy-target.type", "", "Deploy Target type")
	var reqInstanceDiskSizeFlag int64
	flagset.Int64Var(&reqInstanceDiskSizeFlag, "instance.disk-size", 0, "Instance disk size in GiB")
	var reqInstanceIDFlag string
	flagset.StringVar(&reqInstanceIDFlag, "instance.id", "", "Instance ID")
	var reqInstanceInstanceTypeAuthorizedFlag bool
	flagset.BoolVar(&reqInstanceInstanceTypeAuthorizedFlag, "instance.instance-type.authorized", false, "Requires authorization or publicly available")
	var reqInstanceInstanceTypeCpusFlag int64
	flagset.Int64Var(&reqInstanceInstanceTypeCpusFlag, "instance.instance-type.cpus", 0, "CPU count")
	var reqInstanceInstanceTypeFamilyFlag string
	flagset.StringVar(&reqInstanceInstanceTypeFamilyFlag, "instance.instance-type.family", "", "Instance type family")
	var reqInstanceInstanceTypeGpusFlag int64
	flagset.Int64Var(&reqInstanceInstanceTypeGpusFlag, "instance.instance-type.gpus", 0, "GPU count")
	var reqInstanceInstanceTypeIDFlag string
	flagset.StringVar(&reqInstanceInstanceTypeIDFlag, "instance.instance-type.id", "", "Instance type ID")
	var reqInstanceInstanceTypeMemoryFlag int64
	flagset.Int64Var(&reqInstanceInstanceTypeMemoryFlag, "instance.instance-type.memory", 0, "Available memory")
	var reqInstanceInstanceTypeSizeFlag string
	flagset.StringVar(&reqInstanceInstanceTypeSizeFlag, "instance.instance-type.size", "", "Instance type size")
	var reqInstanceIpv6AddressFlag string
	flagset.StringVar(&reqInstanceIpv6AddressFlag, "instance.ipv6-address", "", "Instance IPv6 address")
	var reqInstanceMACAddressFlag string
	flagset.StringVar(&reqInstanceMACAddressFlag, "instance.mac-address", "", "Instance MAC address")
	var reqInstanceManagerIDFlag string
	flagset.StringVar(&reqInstanceManagerIDFlag, "instance.manager.id", "", "Manager ID")
	var reqInstanceManagerTypeFlag string
	flagset.StringVar(&reqInstanceManagerTypeFlag, "instance.manager.type", "", "Manager type")
	var reqInstanceNameFlag string
	flagset.StringVar(&reqInstanceNameFlag, "instance.name", "", "Instance name")
	var reqInstancePublicIPFlag net.IP
	flagset.Net.IPVar(&reqInstancePublicIPFlag, "instance.public-ip", "", "Instance public IPv4 address")
	var reqInstancePublicIPAssignmentFlag string
	flagset.StringVar(&reqInstancePublicIPAssignmentFlag, "instance.public-ip-assignment", "", "")
	var reqInstanceSecurebootEnabledFlag bool
	flagset.BoolVar(&reqInstanceSecurebootEnabledFlag, "instance.secureboot-enabled", false, "Indicates if the instance has secure boot enabled")
	var reqInstanceSSHKeyFingerprintFlag string
	flagset.StringVar(&reqInstanceSSHKeyFingerprintFlag, "instance.ssh-key.fingerprint", "", "SSH key fingerprint")
	var reqInstanceSSHKeyNameFlag string
	flagset.StringVar(&reqInstanceSSHKeyNameFlag, "instance.ssh-key.name", "", "SSH key name")
	var reqInstanceStateFlag string
	flagset.StringVar(&reqInstanceStateFlag, "instance.state", "", "")
	var reqInstanceTemplateBootModeFlag string
	flagset.StringVar(&reqInstanceTemplateBootModeFlag, "instance.template.boot-mode", "", "Boot mode (default: legacy)")
	var reqInstanceTemplateBuildFlag string
	flagset.StringVar(&reqInstanceTemplateBuildFlag, "instance.template.build", "", "Template build")
	var reqInstanceTemplateChecksumFlag string
	flagset.StringVar(&reqInstanceTemplateChecksumFlag, "instance.template.checksum", "", "Template MD5 checksum")
	var reqInstanceTemplateCreatedATFlag time.Time
	flagset.Time.TimeVar(&reqInstanceTemplateCreatedATFlag, "instance.template.created-at", "", "Template creation date")
	var reqInstanceTemplateDefaultUserFlag string
	flagset.StringVar(&reqInstanceTemplateDefaultUserFlag, "instance.template.default-user", "", "Template default user")
	var reqInstanceTemplateDescriptionFlag string
	flagset.StringVar(&reqInstanceTemplateDescriptionFlag, "instance.template.description", "", "Template description")
	var reqInstanceTemplateFamilyFlag string
	flagset.StringVar(&reqInstanceTemplateFamilyFlag, "instance.template.family", "", "Template family")
	var reqInstanceTemplateIDFlag string
	flagset.StringVar(&reqInstanceTemplateIDFlag, "instance.template.id", "", "Template ID")
	var reqInstanceTemplateMaintainerFlag string
	flagset.StringVar(&reqInstanceTemplateMaintainerFlag, "instance.template.maintainer", "", "Template maintainer")
	var reqInstanceTemplateNameFlag string
	flagset.StringVar(&reqInstanceTemplateNameFlag, "instance.template.name", "", "Template name")
	var reqInstanceTemplatePasswordEnabledFlag bool
	flagset.BoolVar(&reqInstanceTemplatePasswordEnabledFlag, "instance.template.password-enabled", false, "Enable password-based login")
	var reqInstanceTemplateSizeFlag int64
	flagset.Int64Var(&reqInstanceTemplateSizeFlag, "instance.template.size", 0, "Template size")
	var reqInstanceTemplateSSHKeyEnabledFlag bool
	flagset.BoolVar(&reqInstanceTemplateSSHKeyEnabledFlag, "instance.template.ssh-key-enabled", false, "Enable SSH key-based login")
	var reqInstanceTemplateURLFlag string
	flagset.StringVar(&reqInstanceTemplateURLFlag, "instance.template.url", "", "Template source URL")
	var reqInstanceTemplateVersionFlag string
	flagset.StringVar(&reqInstanceTemplateVersionFlag, "instance.template.version", "", "Template version")
	var reqInstanceTemplateVisibilityFlag string
	flagset.StringVar(&reqInstanceTemplateVisibilityFlag, "instance.template.visibility", "", "Template visibility")
	var reqInstanceTpmEnabledFlag bool
	flagset.BoolVar(&reqInstanceTpmEnabledFlag, "instance.tpm-enabled", false, "Indicates if the instance has tpm enabled")
	var reqInstanceUserDataFlag string
	flagset.StringVar(&reqInstanceUserDataFlag, "instance.user-data", "", "Instance Cloud-init user-data (base64 encoded)")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.AttachInstanceToSecurityGroupRequest
	if reqInstanceUserDataFlag != "" {
		req.Instanc = &v3.Instance{}
		req.Instance.UserData = reqInstanceUserDataFlag
	}
	if reqInstanceTpmEnabledFlag != false {

		req.Instance.TpmEnabled = reqInstanceTpmEnabledFlag
	}
	if reqInstanceTemplateVisibilityFlag != "" {
		req.InstanceTemplate = &v3.Template{}
		req.InstanceTemplate.Visibility = reqInstanceTemplateVisibilityFlag
	}
	if reqInstanceTemplateVersionFlag != "" {

		req.InstanceTemplate.Version = reqInstanceTemplateVersionFlag
	}
	if reqInstanceTemplateURLFlag != "" {

		req.InstanceTemplate.URL = reqInstanceTemplateURLFlag
	}
	if reqInstanceTemplateSSHKeyEnabledFlag != false {

		req.InstanceTemplate.SSHKeyEnabled = reqInstanceTemplateSSHKeyEnabledFlag
	}
	if reqInstanceTemplateSizeFlag != 0 {

		req.InstanceTemplate.Size = reqInstanceTemplateSizeFlag
	}
	if reqInstanceTemplatePasswordEnabledFlag != false {

		req.InstanceTemplate.PasswordEnabled = reqInstanceTemplatePasswordEnabledFlag
	}
	if reqInstanceTemplateNameFlag != "" {

		req.InstanceTemplate.Name = reqInstanceTemplateNameFlag
	}
	if reqInstanceTemplateMaintainerFlag != "" {

		req.InstanceTemplate.Maintainer = reqInstanceTemplateMaintainerFlag
	}
	if v3.UUID(reqInstanceTemplateIDFlag) != "" {

		req.InstanceTemplate.ID = v3.UUID(reqInstanceTemplateIDFlag)
	}
	if reqInstanceTemplateFamilyFlag != "" {

		req.InstanceTemplate.Family = reqInstanceTemplateFamilyFlag
	}
	if reqInstanceTemplateDescriptionFlag != "" {

		req.InstanceTemplate.Description = reqInstanceTemplateDescriptionFlag
	}
	if reqInstanceTemplateDefaultUserFlag != "" {

		req.InstanceTemplate.DefaultUser = reqInstanceTemplateDefaultUserFlag
	}
	if reqInstanceTemplateCreatedATFlag != "" {

		req.InstanceTemplate.CreatedAT = reqInstanceTemplateCreatedATFlag
	}
	if reqInstanceTemplateChecksumFlag != "" {

		req.InstanceTemplate.Checksum = reqInstanceTemplateChecksumFlag
	}
	if reqInstanceTemplateBuildFlag != "" {

		req.InstanceTemplate.Build = reqInstanceTemplateBuildFlag
	}
	if reqInstanceTemplateBootModeFlag != "" {

		req.InstanceTemplate.BootMode = reqInstanceTemplateBootModeFlag
	}
	if reqInstanceStateFlag != "" {
		req.Instanc = &v3.Instance{}
		req.Instance.State = reqInstanceStateFlag
	}
	if reqInstanceSSHKeyNameFlag != "" {
		req.InstanceSSHKey = &v3.SSHKey{}
		req.InstanceSSHKey.Name = reqInstanceSSHKeyNameFlag
	}
	if reqInstanceSSHKeyFingerprintFlag != "" {

		req.InstanceSSHKey.Fingerprint = reqInstanceSSHKeyFingerprintFlag
	}
	if reqInstanceSecurebootEnabledFlag != false {
		req.Ins = &v3.Instance{}
		req.Instance.SecurebootEnabled = reqInstanceSecurebootEnabledFlag
	}
	if reqInstancePublicIPAssignmentFlag != "" {

		req.Instance.PublicIPAssignment = reqInstancePublicIPAssignmentFlag
	}
	if reqInstancePublicIPFlag != "" {

		req.Instance.PublicIP = reqInstancePublicIPFlag
	}
	if reqInstanceNameFlag != "" {

		req.Instance.Name = reqInstanceNameFlag
	}
	if reqInstanceManagerTypeFlag != "" {
		req.InstanceManager = &v3.Manager{}
		req.InstanceManager.Type = reqInstanceManagerTypeFlag
	}
	if v3.UUID(reqInstanceManagerIDFlag) != "" {

		req.InstanceManager.ID = v3.UUID(reqInstanceManagerIDFlag)
	}
	if reqInstanceMACAddressFlag != "" {
		req.Instanc = &v3.Instance{}
		req.Instance.MACAddress = reqInstanceMACAddressFlag
	}
	if reqInstanceIpv6AddressFlag != "" {

		req.Instance.Ipv6Address = reqInstanceIpv6AddressFlag
	}
	if reqInstanceInstanceTypeSizeFlag != "" {
		req.InstanceInstanceTyp = &v3.InstanceType{}
		req.InstanceInstanceType.Size = reqInstanceInstanceTypeSizeFlag
	}
	if reqInstanceInstanceTypeMemoryFlag != 0 {

		req.InstanceInstanceType.Memory = reqInstanceInstanceTypeMemoryFlag
	}
	if v3.UUID(reqInstanceInstanceTypeIDFlag) != "" {

		req.InstanceInstanceType.ID = v3.UUID(reqInstanceInstanceTypeIDFlag)
	}
	if reqInstanceInstanceTypeGpusFlag != 0 {

		req.InstanceInstanceType.Gpus = reqInstanceInstanceTypeGpusFlag
	}
	if reqInstanceInstanceTypeFamilyFlag != "" {

		req.InstanceInstanceType.Family = reqInstanceInstanceTypeFamilyFlag
	}
	if reqInstanceInstanceTypeCpusFlag != 0 {

		req.InstanceInstanceType.Cpus = reqInstanceInstanceTypeCpusFlag
	}
	if reqInstanceInstanceTypeAuthorizedFlag != false {

		req.InstanceInstanceType.Authorized = reqInstanceInstanceTypeAuthorizedFlag
	}
	if v3.UUID(reqInstanceIDFlag) != "" {
		req.Instance = &v3.Instance{}
		req.Instance.ID = v3.UUID(reqInstanceIDFlag)
	}
	if reqInstanceDiskSizeFlag != 0 {

		req.Instance.DiskSize = reqInstanceDiskSizeFlag
	}
	if reqInstanceDeployTargetTypeFlag != "" {
		req.InstanceDeployTarget = &v3.DeployTarget{}
		req.InstanceDeployTarget.Type = reqInstanceDeployTargetTypeFlag
	}
	if reqInstanceDeployTargetNameFlag != "" {

		req.InstanceDeployTarget.Name = reqInstanceDeployTargetNameFlag
	}
	if v3.UUID(reqInstanceDeployTargetIDFlag) != "" {

		req.InstanceDeployTarget.ID = v3.UUID(reqInstanceDeployTargetIDFlag)
	}
	if reqInstanceDeployTargetDescriptionFlag != "" {

		req.InstanceDeployTarget.Description = reqInstanceDeployTargetDescriptionFlag
	}
	if reqInstanceCreatedATFlag != "" {
		req.Instanc = &v3.Instance{}
		req.Instance.CreatedAT = reqInstanceCreatedATFlag
	}

	resp, err := client.AttachInstanceToSecurityGroup(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DetachInstanceFromSecurityGroupCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("detach-instance-from-security-group", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqInstanceCreatedATFlag time.Time
	flagset.Time.TimeVar(&reqInstanceCreatedATFlag, "instance.created-at", "", "Instance creation date")
	var reqInstanceDeployTargetDescriptionFlag string
	flagset.StringVar(&reqInstanceDeployTargetDescriptionFlag, "instance.deploy-target.description", "", "Deploy Target description")
	var reqInstanceDeployTargetIDFlag string
	flagset.StringVar(&reqInstanceDeployTargetIDFlag, "instance.deploy-target.id", "", "Deploy Target ID")
	var reqInstanceDeployTargetNameFlag string
	flagset.StringVar(&reqInstanceDeployTargetNameFlag, "instance.deploy-target.name", "", "Deploy Target name")
	var reqInstanceDeployTargetTypeFlag string
	flagset.StringVar(&reqInstanceDeployTargetTypeFlag, "instance.deploy-target.type", "", "Deploy Target type")
	var reqInstanceDiskSizeFlag int64
	flagset.Int64Var(&reqInstanceDiskSizeFlag, "instance.disk-size", 0, "Instance disk size in GiB")
	var reqInstanceIDFlag string
	flagset.StringVar(&reqInstanceIDFlag, "instance.id", "", "Instance ID")
	var reqInstanceInstanceTypeAuthorizedFlag bool
	flagset.BoolVar(&reqInstanceInstanceTypeAuthorizedFlag, "instance.instance-type.authorized", false, "Requires authorization or publicly available")
	var reqInstanceInstanceTypeCpusFlag int64
	flagset.Int64Var(&reqInstanceInstanceTypeCpusFlag, "instance.instance-type.cpus", 0, "CPU count")
	var reqInstanceInstanceTypeFamilyFlag string
	flagset.StringVar(&reqInstanceInstanceTypeFamilyFlag, "instance.instance-type.family", "", "Instance type family")
	var reqInstanceInstanceTypeGpusFlag int64
	flagset.Int64Var(&reqInstanceInstanceTypeGpusFlag, "instance.instance-type.gpus", 0, "GPU count")
	var reqInstanceInstanceTypeIDFlag string
	flagset.StringVar(&reqInstanceInstanceTypeIDFlag, "instance.instance-type.id", "", "Instance type ID")
	var reqInstanceInstanceTypeMemoryFlag int64
	flagset.Int64Var(&reqInstanceInstanceTypeMemoryFlag, "instance.instance-type.memory", 0, "Available memory")
	var reqInstanceInstanceTypeSizeFlag string
	flagset.StringVar(&reqInstanceInstanceTypeSizeFlag, "instance.instance-type.size", "", "Instance type size")
	var reqInstanceIpv6AddressFlag string
	flagset.StringVar(&reqInstanceIpv6AddressFlag, "instance.ipv6-address", "", "Instance IPv6 address")
	var reqInstanceMACAddressFlag string
	flagset.StringVar(&reqInstanceMACAddressFlag, "instance.mac-address", "", "Instance MAC address")
	var reqInstanceManagerIDFlag string
	flagset.StringVar(&reqInstanceManagerIDFlag, "instance.manager.id", "", "Manager ID")
	var reqInstanceManagerTypeFlag string
	flagset.StringVar(&reqInstanceManagerTypeFlag, "instance.manager.type", "", "Manager type")
	var reqInstanceNameFlag string
	flagset.StringVar(&reqInstanceNameFlag, "instance.name", "", "Instance name")
	var reqInstancePublicIPFlag net.IP
	flagset.Net.IPVar(&reqInstancePublicIPFlag, "instance.public-ip", "", "Instance public IPv4 address")
	var reqInstancePublicIPAssignmentFlag string
	flagset.StringVar(&reqInstancePublicIPAssignmentFlag, "instance.public-ip-assignment", "", "")
	var reqInstanceSecurebootEnabledFlag bool
	flagset.BoolVar(&reqInstanceSecurebootEnabledFlag, "instance.secureboot-enabled", false, "Indicates if the instance has secure boot enabled")
	var reqInstanceSSHKeyFingerprintFlag string
	flagset.StringVar(&reqInstanceSSHKeyFingerprintFlag, "instance.ssh-key.fingerprint", "", "SSH key fingerprint")
	var reqInstanceSSHKeyNameFlag string
	flagset.StringVar(&reqInstanceSSHKeyNameFlag, "instance.ssh-key.name", "", "SSH key name")
	var reqInstanceStateFlag string
	flagset.StringVar(&reqInstanceStateFlag, "instance.state", "", "")
	var reqInstanceTemplateBootModeFlag string
	flagset.StringVar(&reqInstanceTemplateBootModeFlag, "instance.template.boot-mode", "", "Boot mode (default: legacy)")
	var reqInstanceTemplateBuildFlag string
	flagset.StringVar(&reqInstanceTemplateBuildFlag, "instance.template.build", "", "Template build")
	var reqInstanceTemplateChecksumFlag string
	flagset.StringVar(&reqInstanceTemplateChecksumFlag, "instance.template.checksum", "", "Template MD5 checksum")
	var reqInstanceTemplateCreatedATFlag time.Time
	flagset.Time.TimeVar(&reqInstanceTemplateCreatedATFlag, "instance.template.created-at", "", "Template creation date")
	var reqInstanceTemplateDefaultUserFlag string
	flagset.StringVar(&reqInstanceTemplateDefaultUserFlag, "instance.template.default-user", "", "Template default user")
	var reqInstanceTemplateDescriptionFlag string
	flagset.StringVar(&reqInstanceTemplateDescriptionFlag, "instance.template.description", "", "Template description")
	var reqInstanceTemplateFamilyFlag string
	flagset.StringVar(&reqInstanceTemplateFamilyFlag, "instance.template.family", "", "Template family")
	var reqInstanceTemplateIDFlag string
	flagset.StringVar(&reqInstanceTemplateIDFlag, "instance.template.id", "", "Template ID")
	var reqInstanceTemplateMaintainerFlag string
	flagset.StringVar(&reqInstanceTemplateMaintainerFlag, "instance.template.maintainer", "", "Template maintainer")
	var reqInstanceTemplateNameFlag string
	flagset.StringVar(&reqInstanceTemplateNameFlag, "instance.template.name", "", "Template name")
	var reqInstanceTemplatePasswordEnabledFlag bool
	flagset.BoolVar(&reqInstanceTemplatePasswordEnabledFlag, "instance.template.password-enabled", false, "Enable password-based login")
	var reqInstanceTemplateSizeFlag int64
	flagset.Int64Var(&reqInstanceTemplateSizeFlag, "instance.template.size", 0, "Template size")
	var reqInstanceTemplateSSHKeyEnabledFlag bool
	flagset.BoolVar(&reqInstanceTemplateSSHKeyEnabledFlag, "instance.template.ssh-key-enabled", false, "Enable SSH key-based login")
	var reqInstanceTemplateURLFlag string
	flagset.StringVar(&reqInstanceTemplateURLFlag, "instance.template.url", "", "Template source URL")
	var reqInstanceTemplateVersionFlag string
	flagset.StringVar(&reqInstanceTemplateVersionFlag, "instance.template.version", "", "Template version")
	var reqInstanceTemplateVisibilityFlag string
	flagset.StringVar(&reqInstanceTemplateVisibilityFlag, "instance.template.visibility", "", "Template visibility")
	var reqInstanceTpmEnabledFlag bool
	flagset.BoolVar(&reqInstanceTpmEnabledFlag, "instance.tpm-enabled", false, "Indicates if the instance has tpm enabled")
	var reqInstanceUserDataFlag string
	flagset.StringVar(&reqInstanceUserDataFlag, "instance.user-data", "", "Instance Cloud-init user-data (base64 encoded)")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.DetachInstanceFromSecurityGroupRequest
	if reqInstanceUserDataFlag != "" {
		req.Instanc = &v3.Instance{}
		req.Instance.UserData = reqInstanceUserDataFlag
	}
	if reqInstanceTpmEnabledFlag != false {

		req.Instance.TpmEnabled = reqInstanceTpmEnabledFlag
	}
	if reqInstanceTemplateVisibilityFlag != "" {
		req.InstanceTemplate = &v3.Template{}
		req.InstanceTemplate.Visibility = reqInstanceTemplateVisibilityFlag
	}
	if reqInstanceTemplateVersionFlag != "" {

		req.InstanceTemplate.Version = reqInstanceTemplateVersionFlag
	}
	if reqInstanceTemplateURLFlag != "" {

		req.InstanceTemplate.URL = reqInstanceTemplateURLFlag
	}
	if reqInstanceTemplateSSHKeyEnabledFlag != false {

		req.InstanceTemplate.SSHKeyEnabled = reqInstanceTemplateSSHKeyEnabledFlag
	}
	if reqInstanceTemplateSizeFlag != 0 {

		req.InstanceTemplate.Size = reqInstanceTemplateSizeFlag
	}
	if reqInstanceTemplatePasswordEnabledFlag != false {

		req.InstanceTemplate.PasswordEnabled = reqInstanceTemplatePasswordEnabledFlag
	}
	if reqInstanceTemplateNameFlag != "" {

		req.InstanceTemplate.Name = reqInstanceTemplateNameFlag
	}
	if reqInstanceTemplateMaintainerFlag != "" {

		req.InstanceTemplate.Maintainer = reqInstanceTemplateMaintainerFlag
	}
	if v3.UUID(reqInstanceTemplateIDFlag) != "" {

		req.InstanceTemplate.ID = v3.UUID(reqInstanceTemplateIDFlag)
	}
	if reqInstanceTemplateFamilyFlag != "" {

		req.InstanceTemplate.Family = reqInstanceTemplateFamilyFlag
	}
	if reqInstanceTemplateDescriptionFlag != "" {

		req.InstanceTemplate.Description = reqInstanceTemplateDescriptionFlag
	}
	if reqInstanceTemplateDefaultUserFlag != "" {

		req.InstanceTemplate.DefaultUser = reqInstanceTemplateDefaultUserFlag
	}
	if reqInstanceTemplateCreatedATFlag != "" {

		req.InstanceTemplate.CreatedAT = reqInstanceTemplateCreatedATFlag
	}
	if reqInstanceTemplateChecksumFlag != "" {

		req.InstanceTemplate.Checksum = reqInstanceTemplateChecksumFlag
	}
	if reqInstanceTemplateBuildFlag != "" {

		req.InstanceTemplate.Build = reqInstanceTemplateBuildFlag
	}
	if reqInstanceTemplateBootModeFlag != "" {

		req.InstanceTemplate.BootMode = reqInstanceTemplateBootModeFlag
	}
	if reqInstanceStateFlag != "" {
		req.Instanc = &v3.Instance{}
		req.Instance.State = reqInstanceStateFlag
	}
	if reqInstanceSSHKeyNameFlag != "" {
		req.InstanceSSHKey = &v3.SSHKey{}
		req.InstanceSSHKey.Name = reqInstanceSSHKeyNameFlag
	}
	if reqInstanceSSHKeyFingerprintFlag != "" {

		req.InstanceSSHKey.Fingerprint = reqInstanceSSHKeyFingerprintFlag
	}
	if reqInstanceSecurebootEnabledFlag != false {
		req.Ins = &v3.Instance{}
		req.Instance.SecurebootEnabled = reqInstanceSecurebootEnabledFlag
	}
	if reqInstancePublicIPAssignmentFlag != "" {

		req.Instance.PublicIPAssignment = reqInstancePublicIPAssignmentFlag
	}
	if reqInstancePublicIPFlag != "" {

		req.Instance.PublicIP = reqInstancePublicIPFlag
	}
	if reqInstanceNameFlag != "" {

		req.Instance.Name = reqInstanceNameFlag
	}
	if reqInstanceManagerTypeFlag != "" {
		req.InstanceManager = &v3.Manager{}
		req.InstanceManager.Type = reqInstanceManagerTypeFlag
	}
	if v3.UUID(reqInstanceManagerIDFlag) != "" {

		req.InstanceManager.ID = v3.UUID(reqInstanceManagerIDFlag)
	}
	if reqInstanceMACAddressFlag != "" {
		req.Instanc = &v3.Instance{}
		req.Instance.MACAddress = reqInstanceMACAddressFlag
	}
	if reqInstanceIpv6AddressFlag != "" {

		req.Instance.Ipv6Address = reqInstanceIpv6AddressFlag
	}
	if reqInstanceInstanceTypeSizeFlag != "" {
		req.InstanceInstanceTyp = &v3.InstanceType{}
		req.InstanceInstanceType.Size = reqInstanceInstanceTypeSizeFlag
	}
	if reqInstanceInstanceTypeMemoryFlag != 0 {

		req.InstanceInstanceType.Memory = reqInstanceInstanceTypeMemoryFlag
	}
	if v3.UUID(reqInstanceInstanceTypeIDFlag) != "" {

		req.InstanceInstanceType.ID = v3.UUID(reqInstanceInstanceTypeIDFlag)
	}
	if reqInstanceInstanceTypeGpusFlag != 0 {

		req.InstanceInstanceType.Gpus = reqInstanceInstanceTypeGpusFlag
	}
	if reqInstanceInstanceTypeFamilyFlag != "" {

		req.InstanceInstanceType.Family = reqInstanceInstanceTypeFamilyFlag
	}
	if reqInstanceInstanceTypeCpusFlag != 0 {

		req.InstanceInstanceType.Cpus = reqInstanceInstanceTypeCpusFlag
	}
	if reqInstanceInstanceTypeAuthorizedFlag != false {

		req.InstanceInstanceType.Authorized = reqInstanceInstanceTypeAuthorizedFlag
	}
	if v3.UUID(reqInstanceIDFlag) != "" {
		req.Instance = &v3.Instance{}
		req.Instance.ID = v3.UUID(reqInstanceIDFlag)
	}
	if reqInstanceDiskSizeFlag != 0 {

		req.Instance.DiskSize = reqInstanceDiskSizeFlag
	}
	if reqInstanceDeployTargetTypeFlag != "" {
		req.InstanceDeployTarget = &v3.DeployTarget{}
		req.InstanceDeployTarget.Type = reqInstanceDeployTargetTypeFlag
	}
	if reqInstanceDeployTargetNameFlag != "" {

		req.InstanceDeployTarget.Name = reqInstanceDeployTargetNameFlag
	}
	if v3.UUID(reqInstanceDeployTargetIDFlag) != "" {

		req.InstanceDeployTarget.ID = v3.UUID(reqInstanceDeployTargetIDFlag)
	}
	if reqInstanceDeployTargetDescriptionFlag != "" {

		req.InstanceDeployTarget.Description = reqInstanceDeployTargetDescriptionFlag
	}
	if reqInstanceCreatedATFlag != "" {
		req.Instanc = &v3.Instance{}
		req.Instance.CreatedAT = reqInstanceCreatedATFlag
	}

	resp, err := client.DetachInstanceFromSecurityGroup(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RemoveExternalSourceFromSecurityGroupCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("remove-external-source-from-security-group", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqCidrFlag string
	flagset.StringVar(&reqCidrFlag, "cidr", "", "CIDR-formatted network to remove")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.RemoveExternalSourceFromSecurityGroupRequest
	req.Cidr = reqCidrFlag

	resp, err := client.RemoveExternalSourceFromSecurityGroup(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListSKSClustersCmd(client *v3.Client) {
	resp, err := client.ListSKSClusters(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateSKSClusterCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-sks-cluster", flag.ExitOnError)
	var reqAutoUpgradeFlag bool
	flagset.BoolVar(&reqAutoUpgradeFlag, "auto-upgrade", false, "Enable auto upgrade of the control plane to the latest patch version available")
	var reqCniFlag string
	flagset.StringVar(&reqCniFlag, "cni", "", "Cluster CNI")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Cluster description")
	var reqEnableKubeProxyFlag bool
	flagset.BoolVar(&reqEnableKubeProxyFlag, "enable-kube-proxy", false, "Indicates whether to deploy the Kubernetes network proxy. When unspecified, defaults to `true` unless Cilium CNI is selected")
	var reqLevelFlag string
	flagset.StringVar(&reqLevelFlag, "level", "", "Cluster service level")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Cluster name")
	var reqOidcClientIDFlag string
	flagset.StringVar(&reqOidcClientIDFlag, "oidc.client-id", "", "OpenID client ID")
	var reqOidcGroupsClaimFlag string
	flagset.StringVar(&reqOidcGroupsClaimFlag, "oidc.groups-claim", "", "JWT claim to use as the user's group")
	var reqOidcGroupsPrefixFlag string
	flagset.StringVar(&reqOidcGroupsPrefixFlag, "oidc.groups-prefix", "", "Prefix prepended to group claims")
	var reqOidcIssuerURLFlag string
	flagset.StringVar(&reqOidcIssuerURLFlag, "oidc.issuer-url", "", "OpenID provider URL")
	var reqOidcUsernameClaimFlag string
	flagset.StringVar(&reqOidcUsernameClaimFlag, "oidc.username-claim", "", "JWT claim to use as the user name")
	var reqOidcUsernamePrefixFlag string
	flagset.StringVar(&reqOidcUsernamePrefixFlag, "oidc.username-prefix", "", "Prefix prepended to username claims")
	var reqVersionFlag string
	flagset.StringVar(&reqVersionFlag, "version", "", "Control plane Kubernetes version")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateSKSClusterRequest
	req.Version = reqVersionFlag
	if reqOidcUsernamePrefixFlag != "" {
		req.Oidc = &v3.SKSOidc{}
		req.Oidc.UsernamePrefix = reqOidcUsernamePrefixFlag
	}
	if reqOidcUsernameClaimFlag != "" {

		req.Oidc.UsernameClaim = reqOidcUsernameClaimFlag
	}
	if reqOidcIssuerURLFlag != "" {

		req.Oidc.IssuerURL = reqOidcIssuerURLFlag
	}
	if reqOidcGroupsPrefixFlag != "" {

		req.Oidc.GroupsPrefix = reqOidcGroupsPrefixFlag
	}
	if reqOidcGroupsClaimFlag != "" {

		req.Oidc.GroupsClaim = reqOidcGroupsClaimFlag
	}
	if reqOidcClientIDFlag != "" {

		req.Oidc.ClientID = reqOidcClientIDFlag
	}
	req.Name = reqNameFlag
	req.Level = reqLevelFlag
	req.EnableKubeProxy = reqEnableKubeProxyFlag
	req.Description = reqDescriptionFlag
	req.Cni = reqCniFlag
	req.AutoUpgrade = reqAutoUpgradeFlag

	resp, err := client.CreateSKSCluster(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListSKSClusterDeprecatedResourcesCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("list-sks-cluster-deprecated-resources", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ListSKSClusterDeprecatedResources(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GenerateSKSClusterKubeconfigCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("generate-sks-cluster-kubeconfig", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqTtlFlag int64
	flagset.Int64Var(&reqTtlFlag, "ttl", 0, "Validity in seconds of the Kubeconfig user certificate (default: 30 days)")
	var reqUserFlag string
	flagset.StringVar(&reqUserFlag, "user", "", "User name in the generated Kubeconfig. The certificate present in the Kubeconfig will also have this name set for the CN field.")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.GenerateSKSClusterKubeconfigRequest
	req.User = reqUserFlag
	req.Ttl = reqTtlFlag

	resp, err := client.GenerateSKSClusterKubeconfig(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListSKSClusterVersionsCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("list-sks-cluster-versions", flag.ExitOnError)
	var includeDeprecatedFlag string
	flagset.StringVar(&includeDeprecatedFlag, "IncludeDeprecated", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ListSKSClusterVersions(context.Background(), includeDeprecatedFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteSKSClusterCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-sks-cluster", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteSKSCluster(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetSKSClusterCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-sks-cluster", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetSKSCluster(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateSKSClusterCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-sks-cluster", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqAutoUpgradeFlag bool
	flagset.BoolVar(&reqAutoUpgradeFlag, "auto-upgrade", false, "Enable auto upgrade of the control plane to the latest patch version available")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Cluster description")
	var reqEnableOperatorsCAFlag bool
	flagset.BoolVar(&reqEnableOperatorsCAFlag, "enable-operators-ca", false, "Add or remove the operators certificate authority (CA) from the list of trusted CAs of the api server. The default value is true")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Cluster name")
	var reqOidcClientIDFlag string
	flagset.StringVar(&reqOidcClientIDFlag, "oidc.client-id", "", "OpenID client ID")
	var reqOidcGroupsClaimFlag string
	flagset.StringVar(&reqOidcGroupsClaimFlag, "oidc.groups-claim", "", "JWT claim to use as the user's group")
	var reqOidcGroupsPrefixFlag string
	flagset.StringVar(&reqOidcGroupsPrefixFlag, "oidc.groups-prefix", "", "Prefix prepended to group claims")
	var reqOidcIssuerURLFlag string
	flagset.StringVar(&reqOidcIssuerURLFlag, "oidc.issuer-url", "", "OpenID provider URL")
	var reqOidcUsernameClaimFlag string
	flagset.StringVar(&reqOidcUsernameClaimFlag, "oidc.username-claim", "", "JWT claim to use as the user name")
	var reqOidcUsernamePrefixFlag string
	flagset.StringVar(&reqOidcUsernamePrefixFlag, "oidc.username-prefix", "", "Prefix prepended to username claims")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateSKSClusterRequest
	if reqOidcUsernamePrefixFlag != "" {
		req.Oidc = &v3.SKSOidc{}
		req.Oidc.UsernamePrefix = reqOidcUsernamePrefixFlag
	}
	if reqOidcUsernameClaimFlag != "" {

		req.Oidc.UsernameClaim = reqOidcUsernameClaimFlag
	}
	if reqOidcIssuerURLFlag != "" {

		req.Oidc.IssuerURL = reqOidcIssuerURLFlag
	}
	if reqOidcGroupsPrefixFlag != "" {

		req.Oidc.GroupsPrefix = reqOidcGroupsPrefixFlag
	}
	if reqOidcGroupsClaimFlag != "" {

		req.Oidc.GroupsClaim = reqOidcGroupsClaimFlag
	}
	if reqOidcClientIDFlag != "" {

		req.Oidc.ClientID = reqOidcClientIDFlag
	}
	req.Name = reqNameFlag
	req.EnableOperatorsCA = reqEnableOperatorsCAFlag
	req.Description = reqDescriptionFlag
	req.AutoUpgrade = reqAutoUpgradeFlag

	resp, err := client.UpdateSKSCluster(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetSKSClusterAuthorityCertCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-sks-cluster-authority-cert", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var authorityFlag string
	flagset.StringVar(&authorityFlag, "Authority", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetSKSClusterAuthorityCert(context.Background(), v3.UUID(idFlag), authorityFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetSKSClusterInspectionCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-sks-cluster-inspection", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetSKSClusterInspection(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateSKSNodepoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-sks-nodepool", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDeployTargetDescriptionFlag string
	flagset.StringVar(&reqDeployTargetDescriptionFlag, "deploy-target.description", "", "Deploy Target description")
	var reqDeployTargetIDFlag string
	flagset.StringVar(&reqDeployTargetIDFlag, "deploy-target.id", "", "Deploy Target ID")
	var reqDeployTargetNameFlag string
	flagset.StringVar(&reqDeployTargetNameFlag, "deploy-target.name", "", "Deploy Target name")
	var reqDeployTargetTypeFlag string
	flagset.StringVar(&reqDeployTargetTypeFlag, "deploy-target.type", "", "Deploy Target type")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Nodepool description")
	var reqDiskSizeFlag int64
	flagset.Int64Var(&reqDiskSizeFlag, "disk-size", 0, "Nodepool instances disk size in GiB")
	var reqInstancePrefixFlag string
	flagset.StringVar(&reqInstancePrefixFlag, "instance-prefix", "", "Prefix to apply to instances names (default: pool), lowercase only")
	var reqInstanceTypeAuthorizedFlag bool
	flagset.BoolVar(&reqInstanceTypeAuthorizedFlag, "instance-type.authorized", false, "Requires authorization or publicly available")
	var reqInstanceTypeCpusFlag int64
	flagset.Int64Var(&reqInstanceTypeCpusFlag, "instance-type.cpus", 0, "CPU count")
	var reqInstanceTypeFamilyFlag string
	flagset.StringVar(&reqInstanceTypeFamilyFlag, "instance-type.family", "", "Instance type family")
	var reqInstanceTypeGpusFlag int64
	flagset.Int64Var(&reqInstanceTypeGpusFlag, "instance-type.gpus", 0, "GPU count")
	var reqInstanceTypeIDFlag string
	flagset.StringVar(&reqInstanceTypeIDFlag, "instance-type.id", "", "Instance type ID")
	var reqInstanceTypeMemoryFlag int64
	flagset.Int64Var(&reqInstanceTypeMemoryFlag, "instance-type.memory", 0, "Available memory")
	var reqInstanceTypeSizeFlag string
	flagset.StringVar(&reqInstanceTypeSizeFlag, "instance-type.size", "", "Instance type size")
	var reqKubeletImageGCHighThresholdFlag int64
	flagset.Int64Var(&reqKubeletImageGCHighThresholdFlag, "kubelet-image-gc.high-threshold", 0, "")
	var reqKubeletImageGCLowThresholdFlag int64
	flagset.Int64Var(&reqKubeletImageGCLowThresholdFlag, "kubelet-image-gc.low-threshold", 0, "")
	var reqKubeletImageGCMinAgeFlag string
	flagset.StringVar(&reqKubeletImageGCMinAgeFlag, "kubelet-image-gc.min-age", "", "")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Nodepool name, lowercase only")
	var reqPublicIPAssignmentFlag string
	flagset.StringVar(&reqPublicIPAssignmentFlag, "public-ip-assignment", "", "Configures public IP assignment of the Instances with:  * IPv4 (`inet4`) addressing only (default); * both IPv4 and IPv6 (`dual`) addressing.")
	var reqSizeFlag int64
	flagset.Int64Var(&reqSizeFlag, "size", 0, "Number of instances")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateSKSNodepoolRequest
	req.Size = reqSizeFlag
	req.PublicIPAssignment = reqPublicIPAssignmentFlag
	req.Name = reqNameFlag
	if reqKubeletImageGCMinAgeFlag != "" {
		req.KubeletImageGC = &v3.KubeletImageGC{}
		req.KubeletImageGC.MinAge = reqKubeletImageGCMinAgeFlag
	}
	if reqKubeletImageGCLowThresholdFlag != 0 {

		req.KubeletImageGC.LowThreshold = reqKubeletImageGCLowThresholdFlag
	}
	if reqKubeletImageGCHighThresholdFlag != 0 {

		req.KubeletImageGC.HighThreshold = reqKubeletImageGCHighThresholdFlag
	}
	if reqInstanceTypeSizeFlag != "" {
		req.InstanceTyp = &v3.InstanceType{}
		req.InstanceType.Size = reqInstanceTypeSizeFlag
	}
	if reqInstanceTypeMemoryFlag != 0 {

		req.InstanceType.Memory = reqInstanceTypeMemoryFlag
	}
	if v3.UUID(reqInstanceTypeIDFlag) != "" {

		req.InstanceType.ID = v3.UUID(reqInstanceTypeIDFlag)
	}
	if reqInstanceTypeGpusFlag != 0 {

		req.InstanceType.Gpus = reqInstanceTypeGpusFlag
	}
	if reqInstanceTypeFamilyFlag != "" {

		req.InstanceType.Family = reqInstanceTypeFamilyFlag
	}
	if reqInstanceTypeCpusFlag != 0 {

		req.InstanceType.Cpus = reqInstanceTypeCpusFlag
	}
	if reqInstanceTypeAuthorizedFlag != false {

		req.InstanceType.Authorized = reqInstanceTypeAuthorizedFlag
	}
	req.InstancePrefix = reqInstancePrefixFlag
	req.DiskSize = reqDiskSizeFlag
	req.Description = reqDescriptionFlag
	if reqDeployTargetTypeFlag != "" {
		req.DeployTarget = &v3.DeployTarget{}
		req.DeployTarget.Type = reqDeployTargetTypeFlag
	}
	if reqDeployTargetNameFlag != "" {

		req.DeployTarget.Name = reqDeployTargetNameFlag
	}
	if v3.UUID(reqDeployTargetIDFlag) != "" {

		req.DeployTarget.ID = v3.UUID(reqDeployTargetIDFlag)
	}
	if reqDeployTargetDescriptionFlag != "" {

		req.DeployTarget.Description = reqDeployTargetDescriptionFlag
	}

	resp, err := client.CreateSKSNodepool(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteSKSNodepoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-sks-nodepool", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var sksNodepoolIDFlag string
	flagset.StringVar(&sksNodepoolIDFlag, "SKSNodepoolID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteSKSNodepool(context.Background(), v3.UUID(idFlag), v3.UUID(sksNodepoolIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetSKSNodepoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-sks-nodepool", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var sksNodepoolIDFlag string
	flagset.StringVar(&sksNodepoolIDFlag, "SKSNodepoolID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetSKSNodepool(context.Background(), v3.UUID(idFlag), v3.UUID(sksNodepoolIDFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateSKSNodepoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-sks-nodepool", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var sksNodepoolIDFlag string
	flagset.StringVar(&sksNodepoolIDFlag, "SKSNodepoolID", "", "")
	var reqDeployTargetDescriptionFlag string
	flagset.StringVar(&reqDeployTargetDescriptionFlag, "deploy-target.description", "", "Deploy Target description")
	var reqDeployTargetIDFlag string
	flagset.StringVar(&reqDeployTargetIDFlag, "deploy-target.id", "", "Deploy Target ID")
	var reqDeployTargetNameFlag string
	flagset.StringVar(&reqDeployTargetNameFlag, "deploy-target.name", "", "Deploy Target name")
	var reqDeployTargetTypeFlag string
	flagset.StringVar(&reqDeployTargetTypeFlag, "deploy-target.type", "", "Deploy Target type")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Nodepool description")
	var reqDiskSizeFlag int64
	flagset.Int64Var(&reqDiskSizeFlag, "disk-size", 0, "Nodepool instances disk size in GiB")
	var reqInstancePrefixFlag string
	flagset.StringVar(&reqInstancePrefixFlag, "instance-prefix", "", "Prefix to apply to managed instances names (default: pool), lowercase only")
	var reqInstanceTypeAuthorizedFlag bool
	flagset.BoolVar(&reqInstanceTypeAuthorizedFlag, "instance-type.authorized", false, "Requires authorization or publicly available")
	var reqInstanceTypeCpusFlag int64
	flagset.Int64Var(&reqInstanceTypeCpusFlag, "instance-type.cpus", 0, "CPU count")
	var reqInstanceTypeFamilyFlag string
	flagset.StringVar(&reqInstanceTypeFamilyFlag, "instance-type.family", "", "Instance type family")
	var reqInstanceTypeGpusFlag int64
	flagset.Int64Var(&reqInstanceTypeGpusFlag, "instance-type.gpus", 0, "GPU count")
	var reqInstanceTypeIDFlag string
	flagset.StringVar(&reqInstanceTypeIDFlag, "instance-type.id", "", "Instance type ID")
	var reqInstanceTypeMemoryFlag int64
	flagset.Int64Var(&reqInstanceTypeMemoryFlag, "instance-type.memory", 0, "Available memory")
	var reqInstanceTypeSizeFlag string
	flagset.StringVar(&reqInstanceTypeSizeFlag, "instance-type.size", "", "Instance type size")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Nodepool name, lowercase only")
	var reqPublicIPAssignmentFlag string
	flagset.StringVar(&reqPublicIPAssignmentFlag, "public-ip-assignment", "", "Configures public IP assignment of the Instances with:  * IPv4 (`inet4`) addressing only; * both IPv4 and IPv6 (`dual`) addressing.")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateSKSNodepoolRequest
	req.PublicIPAssignment = reqPublicIPAssignmentFlag
	req.Name = reqNameFlag
	if reqInstanceTypeSizeFlag != "" {
		req.InstanceTyp = &v3.InstanceType{}
		req.InstanceType.Size = reqInstanceTypeSizeFlag
	}
	if reqInstanceTypeMemoryFlag != 0 {

		req.InstanceType.Memory = reqInstanceTypeMemoryFlag
	}
	if v3.UUID(reqInstanceTypeIDFlag) != "" {

		req.InstanceType.ID = v3.UUID(reqInstanceTypeIDFlag)
	}
	if reqInstanceTypeGpusFlag != 0 {

		req.InstanceType.Gpus = reqInstanceTypeGpusFlag
	}
	if reqInstanceTypeFamilyFlag != "" {

		req.InstanceType.Family = reqInstanceTypeFamilyFlag
	}
	if reqInstanceTypeCpusFlag != 0 {

		req.InstanceType.Cpus = reqInstanceTypeCpusFlag
	}
	if reqInstanceTypeAuthorizedFlag != false {

		req.InstanceType.Authorized = reqInstanceTypeAuthorizedFlag
	}
	req.InstancePrefix = reqInstancePrefixFlag
	req.DiskSize = reqDiskSizeFlag
	req.Description = reqDescriptionFlag
	if reqDeployTargetTypeFlag != "" {
		req.DeployTarget = &v3.DeployTarget{}
		req.DeployTarget.Type = reqDeployTargetTypeFlag
	}
	if reqDeployTargetNameFlag != "" {

		req.DeployTarget.Name = reqDeployTargetNameFlag
	}
	if v3.UUID(reqDeployTargetIDFlag) != "" {

		req.DeployTarget.ID = v3.UUID(reqDeployTargetIDFlag)
	}
	if reqDeployTargetDescriptionFlag != "" {

		req.DeployTarget.Description = reqDeployTargetDescriptionFlag
	}

	resp, err := client.UpdateSKSNodepool(context.Background(), v3.UUID(idFlag), v3.UUID(sksNodepoolIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetSKSNodepoolFieldCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-sks-nodepool-field", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var sksNodepoolIDFlag string
	flagset.StringVar(&sksNodepoolIDFlag, "SKSNodepoolID", "", "")
	var fieldFlag string
	flagset.StringVar(&fieldFlag, "Field", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ResetSKSNodepoolField(context.Background(), v3.UUID(idFlag), v3.UUID(sksNodepoolIDFlag), fieldFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func EvictSKSNodepoolMembersCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("evict-sks-nodepool-members", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var sksNodepoolIDFlag string
	flagset.StringVar(&sksNodepoolIDFlag, "SKSNodepoolID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.EvictSKSNodepoolMembersRequest

	resp, err := client.EvictSKSNodepoolMembers(context.Background(), v3.UUID(idFlag), v3.UUID(sksNodepoolIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ScaleSKSNodepoolCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("scale-sks-nodepool", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var sksNodepoolIDFlag string
	flagset.StringVar(&sksNodepoolIDFlag, "SKSNodepoolID", "", "")
	var reqSizeFlag int64
	flagset.Int64Var(&reqSizeFlag, "size", 0, "Number of instances")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.ScaleSKSNodepoolRequest
	req.Size = reqSizeFlag

	resp, err := client.ScaleSKSNodepool(context.Background(), v3.UUID(idFlag), v3.UUID(sksNodepoolIDFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RotateSKSCcmCredentialsCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("rotate-sks-ccm-credentials", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RotateSKSCcmCredentials(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RotateSKSCsiCredentialsCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("rotate-sks-csi-credentials", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RotateSKSCsiCredentials(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RotateSKSOperatorsCACmd(client *v3.Client) {
	flagset := flag.NewFlagSet("rotate-sks-operators-ca", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.RotateSKSOperatorsCA(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpgradeSKSClusterCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("upgrade-sks-cluster", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqVersionFlag string
	flagset.StringVar(&reqVersionFlag, "version", "", "Control plane Kubernetes version")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpgradeSKSClusterRequest
	req.Version = reqVersionFlag

	resp, err := client.UpgradeSKSCluster(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpgradeSKSClusterServiceLevelCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("upgrade-sks-cluster-service-level", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.UpgradeSKSClusterServiceLevel(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ResetSKSClusterFieldCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("reset-sks-cluster-field", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var fieldFlag string
	flagset.StringVar(&fieldFlag, "Field", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ResetSKSClusterField(context.Background(), v3.UUID(idFlag), fieldFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListSnapshotsCmd(client *v3.Client) {
	resp, err := client.ListSnapshots(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteSnapshotCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-snapshot", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteSnapshot(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetSnapshotCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-snapshot", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetSnapshot(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ExportSnapshotCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("export-snapshot", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ExportSnapshot(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func PromoteSnapshotToTemplateCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("promote-snapshot-to-template", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDefaultUserFlag string
	flagset.StringVar(&reqDefaultUserFlag, "default-user", "", "Template default user")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Template description")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Template name")
	var reqPasswordEnabledFlag bool
	flagset.BoolVar(&reqPasswordEnabledFlag, "password-enabled", false, "Enable password-based login in the template")
	var reqSSHKeyEnabledFlag bool
	flagset.BoolVar(&reqSSHKeyEnabledFlag, "ssh-key-enabled", false, "Enable SSH key-based login in the template")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.PromoteSnapshotToTemplateRequest
	req.SSHKeyEnabled = reqSSHKeyEnabledFlag
	req.PasswordEnabled = reqPasswordEnabledFlag
	req.Name = reqNameFlag
	req.Description = reqDescriptionFlag
	req.DefaultUser = reqDefaultUserFlag

	resp, err := client.PromoteSnapshotToTemplate(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListSOSBucketsUsageCmd(client *v3.Client) {
	resp, err := client.ListSOSBucketsUsage(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetSOSPresignedURLCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-sos-presigned-url", flag.ExitOnError)
	var bucketFlag string
	flagset.StringVar(&bucketFlag, "Bucket", "", "")
	var keyFlag string
	flagset.StringVar(&keyFlag, "Key", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetSOSPresignedURL(context.Background(), bucketFlag, keyFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListSSHKeysCmd(client *v3.Client) {
	resp, err := client.ListSSHKeys(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RegisterSSHKeyCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("register-ssh-key", flag.ExitOnError)
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "SSH key name")
	var reqPublicKeyFlag string
	flagset.StringVar(&reqPublicKeyFlag, "public-key", "", "Public key value")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.RegisterSSHKeyRequest
	req.PublicKey = reqPublicKeyFlag
	req.Name = reqNameFlag

	resp, err := client.RegisterSSHKey(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteSSHKeyCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-ssh-key", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteSSHKey(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetSSHKeyCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-ssh-key", flag.ExitOnError)
	var nameFlag string
	flagset.StringVar(&nameFlag, "Name", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetSSHKey(context.Background(), nameFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListTemplatesCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("list-templates", flag.ExitOnError)
	var visibilityFlag string
	flagset.StringVar(&visibilityFlag, "Visibility", "", "")
	var familyFlag string
	flagset.StringVar(&familyFlag, "Family", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.ListTemplates(context.Background(), visibilityFlag, familyFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func RegisterTemplateCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("register-template", flag.ExitOnError)
	var reqBootModeFlag string
	flagset.StringVar(&reqBootModeFlag, "boot-mode", "", "Boot mode (default: legacy)")
	var reqBuildFlag string
	flagset.StringVar(&reqBuildFlag, "build", "", "Template build")
	var reqChecksumFlag string
	flagset.StringVar(&reqChecksumFlag, "checksum", "", "Template MD5 checksum")
	var reqDefaultUserFlag string
	flagset.StringVar(&reqDefaultUserFlag, "default-user", "", "Template default user")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Template description")
	var reqMaintainerFlag string
	flagset.StringVar(&reqMaintainerFlag, "maintainer", "", "Template maintainer")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Template name")
	var reqPasswordEnabledFlag bool
	flagset.BoolVar(&reqPasswordEnabledFlag, "password-enabled", false, "Enable password-based login")
	var reqSizeFlag int64
	flagset.Int64Var(&reqSizeFlag, "size", 0, "Template size")
	var reqSSHKeyEnabledFlag bool
	flagset.BoolVar(&reqSSHKeyEnabledFlag, "ssh-key-enabled", false, "Enable SSH key-based login")
	var reqURLFlag string
	flagset.StringVar(&reqURLFlag, "url", "", "Template source URL")
	var reqVersionFlag string
	flagset.StringVar(&reqVersionFlag, "version", "", "Template version")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.RegisterTemplateRequest
	req.Version = reqVersionFlag
	req.URL = reqURLFlag
	req.SSHKeyEnabled = reqSSHKeyEnabledFlag
	req.Size = reqSizeFlag
	req.PasswordEnabled = reqPasswordEnabledFlag
	req.Name = reqNameFlag
	req.Maintainer = reqMaintainerFlag
	req.Description = reqDescriptionFlag
	req.DefaultUser = reqDefaultUserFlag
	req.Checksum = reqChecksumFlag
	req.Build = reqBuildFlag
	req.BootMode = reqBootModeFlag

	resp, err := client.RegisterTemplate(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteTemplateCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-template", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteTemplate(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetTemplateCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-template", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetTemplate(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CopyTemplateCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("copy-template", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqTargetZoneAPIEndpointFlag Endpoint
	flagset.EndpointVar(&reqTargetZoneAPIEndpointFlag, "target-zone.api-endpoint", "", "Zone API endpoint")
	var reqTargetZoneNameFlag string
	flagset.StringVar(&reqTargetZoneNameFlag, "target-zone.name", "", "")
	var reqTargetZoneSOSEndpointFlag Endpoint
	flagset.EndpointVar(&reqTargetZoneSOSEndpointFlag, "target-zone.sos-endpoint", "", "Zone SOS endpoint")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CopyTemplateRequest
	if reqTargetZoneSOSEndpointFlag != "" {
		req.TargetZone = &v3.Zone{}
		req.TargetZone.SOSEndpoint = reqTargetZoneSOSEndpointFlag
	}
	if reqTargetZoneNameFlag != "" {

		req.TargetZone.Name = reqTargetZoneNameFlag
	}
	if reqTargetZoneAPIEndpointFlag != "" {

		req.TargetZone.APIEndpoint = reqTargetZoneAPIEndpointFlag
	}

	resp, err := client.CopyTemplate(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateTemplateCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-template", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqDescriptionFlag string
	flagset.StringVar(&reqDescriptionFlag, "description", "", "Template Description")
	var reqNameFlag string
	flagset.StringVar(&reqNameFlag, "name", "", "Template name")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateTemplateRequest
	req.Name = reqNameFlag
	req.Description = reqDescriptionFlag

	resp, err := client.UpdateTemplate(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func GetUsageReportCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("get-usage-report", flag.ExitOnError)
	var periodFlag string
	flagset.StringVar(&periodFlag, "Period", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.GetUsageReport(context.Background(), periodFlag)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListUsersCmd(client *v3.Client) {
	resp, err := client.ListUsers(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func CreateUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("create-user", flag.ExitOnError)
	var reqEmailFlag string
	flagset.StringVar(&reqEmailFlag, "email", "", "User Email")
	var reqRoleDescriptionFlag string
	flagset.StringVar(&reqRoleDescriptionFlag, "role.description", "", "IAM Role description")
	var reqRoleEditableFlag bool
	flagset.BoolVar(&reqRoleEditableFlag, "role.editable", false, "IAM Role mutability")
	var reqRoleIDFlag string
	flagset.StringVar(&reqRoleIDFlag, "role.id", "", "IAM Role ID")
	var reqRoleNameFlag string
	flagset.StringVar(&reqRoleNameFlag, "role.name", "", "IAM Role name")
	var reqRolePolicyDefaultServiceStrategyFlag string
	flagset.StringVar(&reqRolePolicyDefaultServiceStrategyFlag, "role.policy.default-service-strategy", "", "IAM default service strategy")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.CreateUserRequest
	if reqRolePolicyDefaultServiceStrategyFlag != "" {
		req.RolePo = &v3.IAMPolicy{}
		req.RolePolicy.DefaultServiceStrategy = reqRolePolicyDefaultServiceStrategyFlag
	}
	if reqRoleNameFlag != "" {
		req.Rol = &v3.IAMRole{}
		req.Role.Name = reqRoleNameFlag
	}
	if v3.UUID(reqRoleIDFlag) != "" {

		req.Role.ID = v3.UUID(reqRoleIDFlag)
	}
	if reqRoleEditableFlag != false {

		req.Role.Editable = reqRoleEditableFlag
	}
	if reqRoleDescriptionFlag != "" {

		req.Role.Description = reqRoleDescriptionFlag
	}
	req.Email = reqEmailFlag

	resp, err := client.CreateUser(context.Background(), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func DeleteUserCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("delete-user", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	resp, err := client.DeleteUser(context.Background(), v3.UUID(idFlag))
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func UpdateUserRoleCmd(client *v3.Client) {
	flagset := flag.NewFlagSet("update-user-role", flag.ExitOnError)
	var idFlag string
	flagset.StringVar(&idFlag, "ID", "", "")
	var reqRoleDescriptionFlag string
	flagset.StringVar(&reqRoleDescriptionFlag, "role.description", "", "IAM Role description")
	var reqRoleEditableFlag bool
	flagset.BoolVar(&reqRoleEditableFlag, "role.editable", false, "IAM Role mutability")
	var reqRoleIDFlag string
	flagset.StringVar(&reqRoleIDFlag, "role.id", "", "IAM Role ID")
	var reqRoleNameFlag string
	flagset.StringVar(&reqRoleNameFlag, "role.name", "", "IAM Role name")
	var reqRolePolicyDefaultServiceStrategyFlag string
	flagset.StringVar(&reqRolePolicyDefaultServiceStrategyFlag, "role.policy.default-service-strategy", "", "IAM default service strategy")

	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])

	// Build request body struct from flags
	var req v3.UpdateUserRoleRequest
	if reqRolePolicyDefaultServiceStrategyFlag != "" {
		req.RolePo = &v3.IAMPolicy{}
		req.RolePolicy.DefaultServiceStrategy = reqRolePolicyDefaultServiceStrategyFlag
	}
	if reqRoleNameFlag != "" {
		req.Rol = &v3.IAMRole{}
		req.Role.Name = reqRoleNameFlag
	}
	if v3.UUID(reqRoleIDFlag) != "" {

		req.Role.ID = v3.UUID(reqRoleIDFlag)
	}
	if reqRoleEditableFlag != false {

		req.Role.Editable = reqRoleEditableFlag
	}
	if reqRoleDescriptionFlag != "" {

		req.Role.Description = reqRoleDescriptionFlag
	}

	resp, err := client.UpdateUserRole(context.Background(), v3.UUID(idFlag), req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func ListZonesCmd(client *v3.Client) {
	resp, err := client.ListZones(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}
