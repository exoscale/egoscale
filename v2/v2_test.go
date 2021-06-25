package v2

import "testing"

func Test_validateOperationParams(t *testing.T) {
	type testStructA struct {
		test string `req-for:"create"`
	}

	type testStructB struct {
		test *string `req-for:"create"`
	}

	type testStructC struct {
		test *string `req-for:"create,update"`
	}

	type testStructD struct {
		test *string `req-for:"*"`
	}

	type args struct {
		res interface{}
		op  string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil pointer value res",
			args: args{
				res: nil,
				op:  "create",
			},
			wantErr: true,
		},
		{
			name: "non-pointer res",
			args: args{
				res: testStructA{},
				op:  "create",
			},
			wantErr: true,
		},
		{
			name: "no op value",
			args: args{
				res: &testStructA{},
				op:  "",
			},
			wantErr: true,
		},
		{
			name: "non-pointer struct field",
			args: args{
				res: &testStructA{test: "test"},
				op:  "create",
			},
			wantErr: true,
		},
		{
			name: "missing required field for operation",
			args: args{
				res: &testStructB{test: nil},
				op:  "create",
			},
			wantErr: true,
		},
		{
			name: "missing required field with wildcard",
			args: args{
				res: &testStructD{test: nil},
				op:  "any",
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				res: &testStructB{test: func() *string { v := "test"; return &v }()},
				op:  "create",
			},
		},
		{
			name: "ok with multiple operations",
			args: args{
				res: &testStructC{test: func() *string { v := "test"; return &v }()},
				op:  "update",
			},
		},
		{
			name: "ok with wildcard",
			args: args{
				res: &testStructD{test: func() *string { v := "test"; return &v }()},
				op:  "any",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			err := validateOperationParams(tt.args.res, tt.args.op)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOperationParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
