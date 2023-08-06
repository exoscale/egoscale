package main

// APIMap maps oapi functions to Consumer API.
var APIMap = Map{
	"compute": nil,
	"dbaas":   nil,
	"dns":     nil,
	"iam": []Entity{
		{
			RootName: "Roles",
			Fns: []Fn{
				Fn{Name: "List", OAPIName: "ListIamRoles"},
				Fn{Name: "Get", OAPIName: "GetIamRole"},
				Fn{Name: "Create", OAPIName: "CreateIamRole"},
				Fn{Name: "Delete", OAPIName: "DeleteIamRole"},
				Fn{Name: "Update", OAPIName: "UpdateIamRole"},
				Fn{Name: "UpdatePolicy", OAPIName: "UpdateIamRolePolicy"},
			},
		},
		{
			RootName: "AccessKey",
			Fns: []Fn{
				Fn{Name: "List", OAPIName: "ListAccessKeys"},
				Fn{Name: "ListKnownOperations", OAPIName: "ListAccessKeyKnownOperations"},
				Fn{Name: "ListOperations", OAPIName: "ListAccessKeyOperations"},
				Fn{Name: "Get", OAPIName: "GetAccessKey"},
				Fn{Name: "Create", OAPIName: "CreateAccessKey"},
				Fn{Name: "Revoke", OAPIName: "RevokeAccessKey"},
			},
		},
		{
			RootName: "OrgPolicy",
			Fns: []Fn{
				Fn{Name: "Get", OAPIName: "GetIamOrganizationPolicy"},
				Fn{Name: "Update", OAPIName: "UpdateIamOrganizationPolicy"},
			},
		},
	},
	"global": nil,
}
