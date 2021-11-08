package oapi

import (
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite

	client oapiClient
}

func (ts *testSuite) SetupTest() {
	ts.client = new(oapiClientMock)
}

func (ts *testSuite) TearDownTest() {
	ts.client = nil
}

func (ts *testSuite) randomID() string {
	id, err := uuid.NewV4()
	if err != nil {
		ts.T().Fatalf("unable to generate a new UUID: %s", err)
	}
	return id.String()
}

func TestSuiteOAPITestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func TestNilableString(t *testing.T) {
	type args struct {
		v *string
	}

	testString := "test"

	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "nil pointer",
			args: args{v: nil},
			want: nil,
		},
		{
			name: "non-empty string",
			args: args{v: &testString},
			want: &testString,
		},
		{
			name: "empty string",
			args: args{v: func() *string { v := ""; return &v }()},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NilableString(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilableString() = %v, want %v", got, tt.want)
			}
		})
	}
}
