package v2

import "testing"

func Test_resetFieldName(t *testing.T) {
	type args struct {
		res   interface{}
		field interface{}
	}

	testStruct := struct {
		fieldA string `reset:"api-field-name"`
		fieldB int
	}{}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "nil field value",
			args: args{
				res:   &testStruct,
				field: nil,
			},
			wantErr: true,
		},
		{
			name: "non-pointer field value",
			args: args{
				res:   &testStruct,
				field: testStruct.fieldA,
			},
			wantErr: true,
		},
		{
			name: "field non-resettable",
			args: args{
				res:   &testStruct,
				field: &testStruct.fieldB,
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				res:   &testStruct,
				field: &testStruct.fieldA,
			},
			want: "api-field-name",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got, err := resetFieldName(tt.args.res, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("resetFieldName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("resetFieldName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
