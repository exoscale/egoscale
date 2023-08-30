package main

// APIMap maps oapi functions to Consumer API.
var APIMap = Map{
	"compute": []Entity{
		{
			RootName: "Instance",
			Fns: []Fn{
				{Name: "ResetPassword", OAPIName: "ResetInstancePassword"},
				{Name: "Get", OAPIName: "GetInstance"},
				{
					Name:                   "List",
					OAPIName:               "ListInstances",
					ResDefOverride:         "[]oapi.InstancesListElement",
					ResPassthroughOverride: "oapi.FromListInstancesResponse(resp)",
				},
			},
		},
		{
			RootName: "InstanceTypes",
			Fns: []Fn{
				{Name: "List", OAPIName: "ListInstanceTypes"},
				{Name: "Get", OAPIName: "GetInstanceType"},
			},
		},
		{
			RootName: "PrivateNetworks",
			Fns: []Fn{
				{Name: "List", OAPIName: "ListPrivateNetworks"},
				{Name: "Get", OAPIName: "GetPrivateNetwork"},
				{Name: "Create", OAPIName: "CreatePrivateNetwork"},
				{Name: "Update", OAPIName: "UpdatePrivateNetwork"},
				{Name: "Delete", OAPIName: "DeletePrivateNetwork"},
			},
		},
		{
			RootName: "SSHKeys",
			Fns: []Fn{
				{Name: "List", OAPIName: "ListSshKeys"},
				{Name: "Register", OAPIName: "RegisterSshKey"},
				{Name: "Delete", OAPIName: "DeleteSshKey"},
				{Name: "Get", OAPIName: "GetSshKey"},
			},
		},
	},
	"dbaas": []Entity{
		{
			RootName: "Integrations",
			Fns: []Fn{
				{
					Name:                   "ListSettings",
					OAPIName:               "ListDbaasIntegrationSettings",
					ResDefOverride:         "*oapi.DBaaSIntegrationSettings",
					ResPassthroughOverride: "oapi.FromListDbaasIntegrationSettingsResponse(resp)",
				},
			},
		},
	},
	"dns": nil,
	"iam": []Entity{
		{
			RootName: "Roles",
			Fns: []Fn{
				{Name: "List", OAPIName: "ListIamRoles"},
				{Name: "Get", OAPIName: "GetIamRole"},
				{Name: "Create", OAPIName: "CreateIamRole"},
				{Name: "Delete", OAPIName: "DeleteIamRole"},
				{Name: "Update", OAPIName: "UpdateIamRole"},
				{Name: "UpdatePolicy", OAPIName: "UpdateIamRolePolicy"},
			},
		},
		{
			RootName: "AccessKey",
			Fns: []Fn{
				{Name: "List", OAPIName: "ListAccessKeys"},
				{Name: "ListKnownOperations", OAPIName: "ListAccessKeyKnownOperations"},
				{Name: "ListOperations", OAPIName: "ListAccessKeyOperations"},
				{Name: "Get", OAPIName: "GetAccessKey"},
				{Name: "Create", OAPIName: "CreateAccessKey"},
				{Name: "Revoke", OAPIName: "RevokeAccessKey"},
			},
		},
		{
			RootName: "OrgPolicy",
			Fns: []Fn{
				{Name: "Get", OAPIName: "GetIamOrganizationPolicy"},
				{Name: "Update", OAPIName: "UpdateIamOrganizationPolicy"},
			},
		},
	},
	"global": []Entity{
		{
			RootName: "Operation",
			Fns: []Fn{
				{Name: "Get", OAPIName: "GetOperation"},
			},
		},
		{
			RootName: "OrgQuotas",
			Fns: []Fn{
				{Name: "List", OAPIName: "ListQuotas"},
				{Name: "Get", OAPIName: "GetQuota"},
			},
		},
	},
}
