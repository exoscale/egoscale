package main

var APIMap = Map{
	"compute": nil,
	"dbaas":   nil,
	"dns":     nil,
	"iam": []Spec{
		{
			RootName: "Roles",
			Fns: []Fn{
				Fn{Name: "List", OAPIName: "ListIamRoles"},
				Fn{Name: "Get", OAPIName: "GetIamRole"},
				Fn{Name: "Create", OAPIName: "CreateIamRole"},
				Fn{Name: "Delete", OAPIName: "DeleteIamRole"},
				Fn{Name: "Update", OAPIName: "UpdateIamRole"},
			},
		},
	},
	"global": nil,
}
