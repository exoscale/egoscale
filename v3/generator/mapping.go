package main

// APIMap maps oapi functions to Consumer API.
var APIMap = Map{
	"compute": nil,
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
	"global": nil,
}
