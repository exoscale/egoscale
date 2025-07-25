Changelog
=========

Unreleased
----------

- v3: use go-retryabble as default HTTP client #710

3.1.24
------

- Generator: use a pointer before unmarshalling into non-pointer structs #707

3.1.23
------

- v3: regenerate from new API spec #704

3.1.22
------

- v2: manually add zag1 zone (#705)
- v2: bump Go version to 1.24 (#706)

3.1.21
------

- v3: regenerate from new API spec #703 #702

3.1.20
------

- v3 Generator: correctly build parameter keys for queries #701

3.1.19
------

- v3 Generator: Fix false additional properties in object #700

3.1.18
------

- v3: `secureboot-enabled` and `tpm-enabled` fields for instance

3.1.17
------

- v3: allow null in **update-sks-cluster** `feature-gates`

3.1.16
------

- v3: rotate-csi-credentials

3.1.15
------

- v3: usage-report, env-impact endpoints

3.1.14
------

- v3: Prometheus URI for Grafana DBAAS, feature-gates as new readable attribute for SKS clusters and enable secureboot/tpm for instances (beta)

3.1.13
------

- Disallow `null` in **create-sks-cluster** `feature-gate` (#686)

3.1.12
----------
- v3: DBaaS Valkey Settings #684

3.1.11
----------
- v3: DBaaS Valkey #682

3.1.10
----------
- v3: Support for finding InstanceTypes by family and size #681

3.1.9
----------
- v3: Reduce minimum blockstorage size to 1GiB from 10GiB

3.1.8
----------
- v3: Client introduce a wait timeout with different polling strategy #674
- v3 generator: Findable return error on too many found #669

3.1.7
----------
- generate_v3: dbaas endpoints schema changes
- v3: withHttpClient deprecation

3.1.6
----------
- generate_v3: dbaas endpoints schema changes

3.1.5
----------
- v3: user agent as client options #659

3.1.4
----------
- v3: downgrade to go1.22 (dedicated module for the generator) #658

3.1.3
----------
- generator: Add Custom Extension Findable #656

3.1.2
----------
- generate_v3: DBaaS endpoints and integration schemas #651 #653
- generate_v3: replace deprecated ::set-output #647 
- generate_v3: handle multiple separators #650
- go.mk: upgrade to v2.0.3 #652 

3.1.1
----------

- v3: Add more listable on Name or(not and) ID #642
- v3: automate regeneration when spec changes #643 
- v3: github action to build tooling with new commit #645 

3.1.0
----------

- v3: Add more HTTP error types #641

3.0.0
----------

- v3 meta-data: private Instance fetch metadata from CD-ROM #634
- v3: Add a metadata package to interact with #632
- v3: Libopenapi bump / nullable reference fix #631
- v3: Remove omitempty on nullable complex types #628
- v3: Add the Authorization header value in the dump request #624
- v3: Add the User-Agent header value #623
- go.mk: remove submodule and initialize through make #618 
- Makefile: fix targets #619 
- SKS nodepool: support for kubelet image gc #620 
- v3 trace: print operation IDs #622 
- go.mk: lint with staticcheck #633 
- v3: update generated code #637
- v3: separate v3 into its own go module #638 
- deprecate v1 and v2 and remove alpha notice for v3 #640 

0.102.3
-------

- v2: fix issue with iam role labels #614

0.102.2
-------

- v2: Fixed `SetHTTPClient` #605

0.102.1
-------

- v2: IAMv3 Policy Resources field now noop #604

0.102.0
-------

- v2: implement IAMv3

0.101.1
-------

- fix: give correct permissions to release workflow

0.101.0
-------

- v2: regenerate (2023-08-31) #599

0.100.3
-------

- v2: add validation for UUID

0.100.2
-------

- v2: add private network labels

0.100.1
-------

- v2: refresh openapi generated code (patches generated code to fix codegen bug)

0.100.0
-------

- feature: v2: allow supply security group rules referencing public security groups.

0.99.0
------

- feature: v2: add RevealInstancePassword

0.98.1
------

- feature: v2: add FindSecurityGroups with parameter support

0.98.0
------

- v2: refresh openapi generated code
- feature: v2: add an option to the ListInstances call that allows users to filter by IP address

0.97.0
------

- v2: database migration stop methods for databases supporting it: Redis, PostgreSQL, MySQL

0.96.0
------

- v2: refresh openapi generated code
- feature: v2: rename SKS nodepool addon


0.95.0
------

- feature: v2: add `Client.FindTemplate()` method
- fix: v2: error for `Client.GetTemplateByName()`

0.94.0
------

- feature: v2: add `Client.GetTemplateByName()` method
- feature: v2: implement sort.Interface for []*Template by CreatedAt or by Name

0.93.0
------

- feature: v2: implement publicIpAssignment for Instances

0.92.0
------

- feature: v2: implement reverse DNS management

0.91.0
------

- feature: v2: add `Client.FindDatabaseService()` method

0.90.4
------

- fix: v2: seg fault when listing security groups

0.90.3
------

- feature: v2: add labels support for Elastic IPs

0.90.2
------

- fix: v2: silence retryable HTTP client debug logs

0.90.1
------

- feature: v2: retryable HTTP client by default
- change: no HOST header override for IP address endpoints

0.90.0
------

- feature: v2: add support for Elastic IPv6

0.89.0
------

- feature: v2: add support for DNS management

0.88.2
------

- feature: v2: add support for Build, Maintainer and Version attributes to RegisterTemplate

0.88.1
------

- fix: v2: aligns request signing with public API by urlencoding path in the signature

0.88.0
------

- feature: v2: WaitInstancePoolConverged allows to waits until an instance pool's VMs are provisioned

0.87.0
------

- v2: refresh openapi generated code + fix missing type oapi.Reference

0.86.0
------

- feature: v2: DBaaS: implement migration status command

0.85.0
------

- feature: v2: add support for listing of SKS Cluster deprecated resources via `ListSKSClusterDeprecatedResources`

0.84.3
------

- change: v2: refresh code generated from public API spec

0.84.2
------

- change: v2: refresh code generated from public API spec

0.84.1
------

- fix: v2: fix `CreateSecurityGroupRule()` method (#547)

0.84.0
------

- change: v2: `SKSClusterOIDCConfig` struct field `RequiredClaim` now is a `map[string]string` type instead of a string

0.83.2
------

- change: v2: refresh code generated from public API spec

0.83.1
------

- change: v2: refresh code generated from public API spec

0.83.0
------

- feature: v2: add support for IAM access key resources

0.82.1
------

- change: v2: refresh code generated from public API spec

0.82.0
------

- feature: v2: add support for IAM access keys management
- feature: v2: add support for SKS Cluster OIDC configuration via `CreateSKSClusterOpt` options
- feature: v2: add `ListInstancesOpt` options

0.81.0
------

- change: v2: `DatabaseServiceType`'s `LatestVersion` field has been replaced by `AvailableVersions`
- fix: v2:`Update*()` methods no longer send empty strings to the public API, which does not accept those anymore

0.80.2
------

- change: v2: refresh code generated from public API spec

0.80.1
------

- change: v2: refresh code generated from public API spec

0.80.0
------

- feature: v2: add new `CopyTemplate()`/`UpdateTemplate()` methods

0.79.0
------

- feature: v2: `ListSKSClusterVersions()` method now accepts an additional `ListSKSClusterVersionsOpt` variadic argument

0.78.0
------

- change: v2: type-specific `DatabaseService` methods have been removed, to be re-implemented in a future version; use `v2/oapi` methods for type-specific operations in the meantime

0.77.0
------

- change: v2: switch to func-based options passing

0.76.0
------

- feature: v2: add support for SKS Nodepool taints

0.75.0
------

- feature: v2: add `StartInstanceOpt` options

0.74.0
------

- feature: v2: add a new `Zone` struct field to zone-local API resources

0.73.2
------

- fix: v2: add missing operations params validation

0.73.1
------

- fix: v2: `SecurityGroup`: return external sources when present

0.73.0
------

- feature: v2: add `Client.AddExternalSourceToSecurityGroup()`/`Client.RemoveExternalSourceFromSecurityGroup()` methods

0.72.2
------

- fix: v2: update Exoscale API endpoint prefix

0.72.1
------

- fix: v2: fix `Client.UpdateElasticIP()` method

0.72.0
------

- feature: v2: add `AntiAffinityGroup.InstanceIDs` field

0.71.1
------

- fix: v2: fix `Client.CreateSecurityGroupRule()` method

0.71.0
------

- feature: v2: add `Client.UpgradeSKSClusterServiceLevel()` method

0.70.0
------

- feature: v2: add `DatabaseServiceComponent` struct

0.69.0
------

- feature: v2: add method `Client.GetDatabaseCACertificate()`

0.68.1
------

- fix: v2: add missing `Snapshot.Size` field

0.68.0
------

- feature: v2: add support for quotas management
- change: v2: all API resource-based methods have been relocated to the `Client` struct

0.67.0
------

- feature: v2: add support for SKS Nodepool add-ons

0.66.0
------

- feature: v2: add `Instance.Reset()` method
- feature: v2: add `Instance.Scale()` method
- feature: v2: add `Instance.ResizeDisk()` method

0.65.1
------

- fix: v2: fix `RegisterSSHKey()` method

0.65.0
------

- feature: v2: add support for SSH keys management

0.64.1
------

- tests: v2: add resource API mocks

0.64.0
------

- change: v2: replace `InstancePool.ManagerID` of type `string` with `InstancePool.Manager` field of type `*InstancePoolManager`

0.63.0
------

- feature: v2: add support for Private Networks to SKS Nodepools
- feature: v2: add new `Client.RegisterTemplate()` method
- change: v2: change `DatabaseService.UserConfig` type to pointer

0.62.2
------

- v2: fix a crash in `NetworkLoadBalancer.AddService()` method

0.62.1
------

- fix: v2: fix required params validation for NLB services

0.62.0
------

- feature: v2: add support for Private Networks leases

0.61.0
------

- feature: v2: add `PrivateNetwork.UpdateInstanceIPAddress()` method
- feature: v2: add `Instance.Reboot()` method

0.60.1
------

- fix: v2: don't return pointers to empty maps/slices

0.60.0
------

- change: v2: API resource structs fields are now pointers instead of concrete types

0.59.0
------

- change: v2: `Database*` structs fields are now pointers

0.58.0
------

- feature: v2: add labels support for Compute instances

0.57.0
------

- feature: v2: add support for Database Services

0.56.0
------

- change: the `AuthorizeSecurityGroupIngress` struct now uses an `int` type for `Icmp(Code|Type)` fields (#499)

0.55.0
------

- change: the `IngressRule`/`EgressRule` and `v2.SecurityGroupRule` structs now use an `int` type for the ICMP code/type storage (#498)

0.54.0
------

- change: the `IngressRule`/`EgressRule` and `v2.SecurityGroupRule` structs now use an `int8` type for the ICMP code/type storage (#497)

0.53.1
------

- fix: v2: only point to non-zero struct fields for optional API resource properties (#496)

0.53.0
------

- feature: v2: add `Client.FindInstanceType()` method

0.52.0
------

- feature: v2: make API async polling interval customizable
- feature: v2: add `Client.Find*()` methods
- feature: v2: add `Start`/`Stop` methods to `Instance`
- feature: v2: add labels support for Network Load Balancers
- fix: v2: fix Security Group parsing from API

0.51.0
------

- feature: v2: add new `InstanceType` resource

0.50.0
------

- change: v2: the `Instance.ManagerID` field is replaced with `Instance.Manager` of type `*InstanceManager`

0.49.0
------

- deprecatation: top-level `Version` constant is replaced by `version.Version`
- change: v2: new default HTTP client transport setting request `User-Agent` header to `v2.UserAgent`
- feature: v2: add support for Elastic IP/Private Network/Security Group attachment/detachment to Instances
- feature: v2: add support for Deploy Targets to SKS Nodepools
- feature: v2: add support for Instance prefix to SKS Nodepools

0.48.1
------

- fix: v2: add support for `InstancePool.IPv6Enabled` field resetting

0.48.0
------

- feature: v2: add support for Instance prefix to Instance Pools
- feature: v2: add support for Deploy Targets
- feature: v2: add support for Compute instances management
- feature: v2: add support for Snapshots management
- feature: v2: add support for Templates management
- feature: v2: add getter methods on API resources

0.47.0
------

- feature: v2: add client property setters (#485)

0.46.1
------

- fix: v2: make SKSCluster.RotateCCMCredentials() synchronous (#484)

0.46.0
------

- feature: SKS: add SKSCluster.RotateCCMCredentials() method (#481)
- feature: SKS: add SKSCluster.AuthorityCert() method (#480)

0.45.1
------

- Fix typo in version.go

0.45.0
------

- feature: v2: add support for Elastic IP management
- fix: v2: InstancePool.ManagerID resolution (#479)

0.44.0
------

- feature: v2: add request tracing middleware (#474)
- feature: v2: add support for field resetting (#476)
- feature: v2: add support for Instance Pools management (#471)
- feature: v2: add support for Private Networks management (#472)
- feature: v2: add support for Anti-Affinity Groups management (#473)
- feature: v2: add support for Security Groups management (#475)

0.43.1
------

- change: in `NewClient()`, the `v2.Client` embedded in the `Client` struct doesn't inherit the custom `http.Client` set using `WithHTTPClient()`.

0.43.0
------

- change: [Exoscale API V2](https://openapi-v2.exoscale.com/) related code has been relocated under the `github.com/exoscale/egoscale/v2` package.
  Note: `egoscale.Client` embeds a `v2.Client` initialized implicitly as a convenience.

0.42.0
------

- feature: new `SKSNodepool.AntiAffinityGroupIDs` field
- change: `SKSCluster.Level` field renamed as `SKSCluster.ServiceLevel`

0.41.0
------

- feature: new method `ListZones()`

0.40.1
------

- Improve API v2 async job tests and error reporting (#466)

0.40.0
------

- feature: new method `UpgradeSKSCluster()`
- feature: new fields `SKSCluster.Level` and `SKSCluster.CNI`
- change: `SKSCluster.EnableExoscaleCloudController` replaced with `SKSCluster.AddOns`

0.39.1
------

- fix: add missing `UpdateVirtualMachineSecurityGroups` operation metadata

0.39.0
------

- feature: add `UpdateVirtualMachineSecurityGroups` operation (#464)

0.38.0
------

- feature: add `SKSCluster.EvictNodepoolMembers()` and `ListSKSClusterVersions()` methods

0.37.1
------

- fix: `UpdateIPAddress.HealthcheckTLSSkipVerify` field always set to `false` (#462)

0.37.0
------

- feature: `NewClient()` now accepts options (460)
- fix: NLB service healthcheck TLS SNI bug (#461)

0.36.2
------

- fix: `CreateInstancePool.AntiAffinityGroupIDs` field is optional (#459)

0.36.1
------

- feature: add support for Exoscale Cloud Controller in SKS clusters
- fix: add missing tests for SKS Nodepools Security Groups

0.36.0
------

- feature: add support for Anti-Affinity Groups to Instance Pools
- feature: add support for Security Groups to SKS Nodepools

0.35.3
------

- Fix typo in version.go

0.35.2
------

- Improve API v2 errors handling (#455)

0.35.1
------

- fix: various SKS-related bugs (#454)

0.35.0
------

- feature: add support for SKS resources (#453)

0.34.0
------

- change: `BucketUsage.Usage` is now an `int64` (#451)

0.33.2
------

- fix: make `GetWithContext` return more relevant errors (#450)

0.33.1
------

- fix: `UpdateNetworkLoadBalancer` call panicking following a public API change

0.33.0
------

- feature: add support for Network Load Balancer service HTTPS health checking (#449)

0.32.0
------

- feature: add support for Instance Pool root disk size update (#448)

0.31.2
------

- fix: add missing TLS-specific parameters to `AssociateIPAddress`

0.31.1
------

- fix: Instance Pool IPv6 flag handling

0.31.0
------

- feature: add support for IPv6 in Instance Pools (#446)

0.30.0
------

- feature: add new TLS-specific parameters to managed EIP

0.29.0
------

- feature: `ListVirtualMachines` call to allow searching by `ManagerID` (#442)
- fix: remove duplicate `User-Agent` HTTP header in Runstatus calls
- tests: `*NetworkLoadBalancer*` calls are now tested using HTTP mocks
- codegen: `internal/v2` updated

0.28.1
------

- fix: Fix `ListVolumes` call to allow searching by ID (#440)

0.28.0
------

- feature: add `Manager`/`ManagerID` fields to `VirtualMachine` structure (#438)
- fix: HTTP request User Agent header handling (#439)

0.27.0
------

- feature: Add `evictInstancePoolMembers` call to Instance Pool (#437)

0.26.6
------

- change: Add support for Compute instance templates boot mode (#436)

0.26.5
------

- fix: bug in the ListNetworkLoadBalancers call (#435)

0.26.4
------

- Fixing typo in previous release

0.26.3
------

- change: updated API V2 async operation code (#434)

0.26.2
------

- change: updated OpenAPI code-generated API V2 bindings

0.26.1
------

- change: the `DisplayText` property of `RegisterCustomTemplate` is now optional (#433)

0.26.0
------

- feature: Add support for Network Load Balancer resources (#432)

0.25.0
------

- feature: Add support for `listBucketsUsage` (#431)
- change: Switch CI to Github Actions (#430)

0.24.0
------

- feature: Add export snapshot implementation (#427)
- feature: Add support for public API V2 (#425)
- change: Switch module to Go 1.14 (#429)
- change: Travis CI: set minimum Go version to 1.13
- doc: Annotate API doc regarding use of tags (#423)
- tests: fix request client timeout handling (#422)

0.23.0
------

- change: Add `Resources` field to `APIKey` (#420)

0.22.0
------

- change: Remove all references to Network Offerings (#418)

0.21.0
------

- feature: add const `NotFound` 404 on type `ErrorCode` (#417)

0.20.1
------

- fix: update the `ListAPIKeysResponse` field (#415)

0.20.0
------

- feature: Add Instance pool implementation (#410)
- feature: Add IAM implementation (#411)

0.19.0
------

- feature: add field `Description` on type `IPAddress` (#413)
- change: add Json tag `omitempty` on field  `TemplateFilter` in type `ListTemplates` (#412)

0.18.1
------

- change: make the "User-Agent" HTTP request header more informative and exposed

0.18.0
------

- feature: add method `DeepCopy` on type `AsyncJobResult` (#403)

0.17.2
------

- remove: remove the `IsFeatured` parameter from call `RegisterCustomTemplate` (#402)

0.17.1
------

- feature: add parameter `RescueProfile` to call `StartVirtualMachine` (#401)

0.17.0
------

- feature: add new call `RegisterCustomTemplate` (#400)
- feature: add new call `DeleteTemplate` (#399)

0.16.0
------

- feature: Add `Healthcheck*` parameters to call `UpdateIPAddress`
- change: Replace satori/go.uuid by gofrs/uuid

0.15.0
------

- change: prefix the healthcheck-related params with `Healthcheck` on call `AssociateIPAddress`
- EIP: the healthcheck should be a pointer
- ip addresses: Add the Healthcheck parameters
- readme: point to new lego org (#395)
- dns: user_id is not sent back (#394)

0.14.3
------

- fix: `AffinityGroup` lists virtual machines with `UUID` rather than string

0.14.2
------

- fix: `ListVirtualMachines` by `IDs` to accept `UUID` rather than string

0.14.1
------

- fix: `GetRunstatusPage` to always contain the subresources
- fix: `ListRunstatus*` to fetch all the subresources
- feature: `PaginateRunstatus*` used by list

0.14.0
------

- change: all DNS calls require a context
- fix: `CreateAffinityGroup` allows empty `name`

0.13.3
------

- fix: runstatus unmarshalling errors
- feature: `UUID` implements DeepCopy, DeepCopyInto
- change: export `BooleanResponse`

0.13.2
------

- feat: initial Runstatus API support
- feat: `admin` namespace containing `ListVirtualMachines` for admin usage

0.13.1
------

- feat: `Iso` support `ListIsos`, `AttachIso`, and `DetachIso`

0.13.0
------

- change: `Paginate` to accept `Listable`
- change: `ListCommand` is also `Listable`
- change: `client.Get` doesn't modify the given resource, returns a new one
- change: `Command` and `AsyncCommand` are fully public, thus extensible
- remove: `Gettable`

0.12.5
------

- fix: `AuthorizeSecurityGroupEgress` could return `authorizeSecurityGroupIngress` as name

0.12.4
------

- feat: `Snapshot` is `Listable`

0.12.3
------

- change: replace dep by Go modules
- change: remove domainid,domain,regionid,listall,isrecursive,... fields
- remove: `MigrateVirtualMachine`, `CreateUser`, `EnableAccount`, and other admin calls

0.12.2
------

- fix: `ListNics` has no virtualmachineid limitations anymore
- fix: `PCIDevice` ids are not UUIDs

0.12.1
------

- fix: `UpdateVMNicIP` is async

0.12.0
------

- feat: new VM state `Moving`
- feat: `UpdateNetwork` with `startip`, `endip`, `netmask`
- feat: `NetworkOffering` is `Listable`
- feat: when it fails parsing the body, it shows it
- fix: `Snapshot.State` is a string, rather than an scalar
- change: signature are now using the v3 version with expires by default

0.11.6
------

- fix: `Network.ListRequest` accepts a `Name` argument
- change: `SecurityGroup` and the rules aren't `Taggable` anymore

0.11.5
------

- feat: addition of `UpdateVMNicIP`
- fix: `UpdateVMAffinityGroup` expected response

0.11.4
------

*no changes in the core library*

0.11.3
------

*no changes in the core library*

0.11.2
------

- fix: empty list responses

0.11.1
------

- fix: `client.Sign` handles correctly the brackets (kudos to @stffabi)
- change: `client.Payload` returns a `url.Values`

0.11.0
------

- feat: `listOSCategories` and `OSCategory` type
- feat: `listApis` supports recursive response structures
- feat: `GetRecordsWithFilters` to list records with name or record_type filters
- fix: better `DNSErrorResponse`
- fix: `ListResourceLimits` type
- change: use UUID everywhere

0.10.5
------

- feat: `Client.Logger` to plug in any `*log.Logger`
- feat: `Client.TraceOn`/`ClientTraceOff` to toggle the HTTP tracing

0.10.4
------

- feat: `CIDR` to replace string string
- fix: prevent panic on nil

0.10.3
------

- feat: `Account` is Listable
- feat: `MACAddress` to replace string type
- fix: Go 1.7 support

0.10.2
------

- fix: ActivateIP6 response

0.10.1
------

- feat: expose `SyncRequest` and `SyncRequestWithContext`
- feat: addition of reverse DNS calls
- feat: addition of `SecurityGroup.UserSecurityGroup`

0.10.0
------

- global: cloudstack documentation links are moved into cs
- global: removal of all the `...Response` types
- feat: `Network` is `Listable`
- feat: addition of `deleteUser`
- feat: addition of `listHosts`
- feat: addition of `updateHost`
- feat: exo cmd (kudos to @pierre-emmanuelJ)
- change: refactor `Gettable` to use `ListRequest`

0.9.31
------

- fix: `IPAddress`.`ListRequest` with boolean fields
- fix: `Network`.`ListRequest` with boolean fields
- fix: `ServiceOffering`.`ListRequest` with boolean fields

0.9.30
------

- fix: `VirtualMachine` `PCIDevice` representation was incomplete

0.9.29
------

- change: `DNSErrorResponse` is a proper `error`

0.9.28
------

- feat: addition of `GetDomains`
- fix: `UpdateDomain` may contain more empty fields than `CreateDomain`

0.9.27
------

- fix: expects body to be `application/json`

0.9.26
------

- change: async timeout strategy wait two seconds and not fib(n) seconds

0.9.25
------

- fix: `GetVirtualUserData` response with `Decode` method handling base64 and gzip

0.9.24
------

- feat: `Template` is `Gettable`
- feat: `ServiceOffering` is `Gettable`
- feat: addition of `GetAPILimit`
- feat: addition of `CreateTemplate`, `PrepareTemplate`, `CopyTemplate`, `UpdateTemplate`, `RegisterTemplate`
- feat: addition of `MigrateVirtualMachine`
- feat: cmd cli
- change: remove useless fields related to Project and VPC

0.9.23
------

- feat: `booleanResponse` supports true booleans: https://github.com/apache/cloudstack/pull/2428

0.9.22
------

- feat: `ListUsers`, `CreateUser`, `UpdateUser`
- feat: `ListResourceDetails`
- feat: `SecurityGroup` helper `RuleByID`
- feat: `Sign` signs the payload
- feat: `UpdateNetworkOffering`
- feat: `GetVirtualMachineUserData`
- feat: `EnableAccount` and `DisableAccount` (admin stuff)
- feat: `AsyncRequest` and `AsyncRequestWithContext` to examine the polling
- fix: `AuthorizeSecurityGroupIngress` support for ICMPv6
- change: move `APIName()` into the `Client`, nice godoc
- change: `Payload` doesn't sign the request anymore
- change: `Client` exposes more of its underlying data
- change: requests are sent as GET unless it body size is too big

0.9.21
------

- feat: `Network` is `Listable`
- feat: `Zone` is `Gettable`
- feat: `Client.Payload` to help preview the HTTP parameters
- feat: generate command utility
- fix: `CreateSnapshot` was missing the `Name` attribute
- fix: `ListSnapshots` was missing the `IDs` attribute
- fix: `ListZones` was missing the `NetworkType` attribute
- fix: `ListAsyncJobs` was missing the `ListAll` attribute
- change: ICMP Type/Code are uint8 and TCP/UDP port are uint16

0.9.20
------

- feat: `Template` is `Listable`
- feat: `IPAddress` is `Listable`
- change: `List` and `Paginate` return pointers
- fix: `Template` was missing `tags`

0.9.19
------

- feat: `SSHKeyPair` is `Listable`

0.9.18
------

- feat: `VirtualMachine` is `Listable`
- feat: new `Client.Paginate` and `Client.PaginateWithContext`
- change: the inner logic of `Listable`
- remove: not working `Client.AsyncList`

0.9.17
------

- fix: `AuthorizeSecurityGroup(In|E)gress` startport may be zero

0.9.16
------

- feat: new `Listable` interface
- feat: `Nic` is `Listable`
- feat: `Volume` is `Listable`
- feat: `Zone` is `Listable`
- feat: `AffinityGroup` is `Listable`
- remove: deprecated methods `ListNics`, `AddIPToNic`, and `RemoveIPFromNic`
- remove: deprecated method `GetRootVolumeForVirtualMachine`

0.9.15
------

- feat: `IPAddress` is `Gettable` and `Deletable`
- fix: serialization of *bool

0.9.14
------

- fix: `GetVMPassword` response
- remove: deprecated `GetTopology`, `GetImages`, and al

0.9.13
------

- feat: IP4 and IP6 flags to DeployVirtualMachine
- feat: add ActivateIP6
- fix: error message was gobbled on 40x

0.9.12
------

- feat: add `BooleanRequestWithContext`
- feat: add `client.Get`, `client.GetWithContext` to fetch a resource
- feat: add `cleint.Delete`, `client.DeleteWithContext` to delete a resource
- feat: `SSHKeyPair` is `Gettable` and `Deletable`
- feat: `VirtualMachine` is `Gettable` and `Deletable`
- feat: `AffinityGroup` is `Gettable` and `Deletable`
- feat: `SecurityGroup` is `Gettable` and `Deletable`
- remove: deprecated methods `CreateAffinityGroup`, `DeleteAffinityGroup`
- remove: deprecated methods `CreateKeypair`, `DeleteKeypair`, `RegisterKeypair`
- remove: deprecated method `GetSecurityGroupID`

0.9.11
------

- feat: CloudStack API name is now public `APIName()`
- feat: enforce the mutual exclusivity of some fields
- feat: add `context.Context` to `RequestWithContext`
- change: `AsyncRequest` and `BooleanAsyncRequest` are gone, use `Request` and `BooleanRequest` instead.
- change: `AsyncInfo` is no more

0.9.10
------

- fix: typo made ListAll required in ListPublicIPAddresses
- fix: all bool are now *bool, respecting CS default value
- feat: (*VM).DefaultNic() to obtain the main Nic

0.9.9
-----

- fix: affinity groups virtualmachineIds attribute
- fix: uuidList is not a list of strings

0.9.8
-----

- feat: add RootDiskSize to RestoreVirtualMachine
- fix: monotonic polling using Context

0.9.7
-----

- feat: add Taggable interface to expose ResourceType
- feat: add (Create|Update|Delete|List)InstanceGroup(s)
- feat: add RegisterUserKeys
- feat: add ListResourceLimits
- feat: add ListAccounts

0.9.6
-----

- fix: update UpdateVirtualMachine userdata
- fix: Network's name/displaytext might be empty

0.9.5
-----

- fix: serialization of slice

0.9.4
-----

- fix: constants

0.9.3
-----

- change: userdata expects a string
- change: no pointer in sub-struct's

0.9.2
-----

- bug: createNetwork is a sync call
- bug: typo in listVirtualMachines' domainid
- bug: serialization of map[string], e.g. UpdateVirtualMachine
- change: IPAddress's use net.IP type
- feat: helpers VM.NicsByType, VM.NicByNetworkID, VM.NicByID
- feat: addition of CloudStack ApiErrorCode constants

0.9.1
-----

- bug: sync calls returns succes as a string rather than a bool
- change: unexport BooleanResponse types
- feat: original CloudStack error response can be obtained

0.9.0
-----

Big refactoring, addition of the documentation, compliance to golint.

0.1.0
-----

Initial library
