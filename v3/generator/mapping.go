package main

// APIMap maps oapi functions to Consumer API.
var APIMap = Map{
	"Compute": []Entity{
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
	"DBaaS": []Entity{
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
	"DNS": nil,
	"IAM": []Entity{
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
	"Global": []Entity{
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
		{
			RootName: "Zones",
			Fns: []Fn{
				{Name: "List", OAPIName: "ListZones"},
			},
		},
	},
}
