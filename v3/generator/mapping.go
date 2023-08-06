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
	},
	"global": nil,
}
