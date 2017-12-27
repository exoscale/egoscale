package egoscale

import (
	"log"
	"testing"
)

func TestListAffinityGroups(t *testing.T) {
	ts := newServer(200, `
{
	"listaffinitygroupsresponse": {
		"affinitygroup": [
			{
				"account": "yoan.blanc@exoscale.ch",
				"domain": "yoan.blanc@exoscale.ch",
				"domainid": "2da0d0d3-e7b2-42ef-805d-eb2ea90ae7ef",
				"id": "ade69201-767b-44f4-bd84-0f20a807ce78",
				"name": "",
				"type": "host anti-affinity"
			},
			{
				"account": "yoan.blanc@exoscale.ch",
				"description": "hello world!",
				"domain": "yoan.blanc@exoscale.ch",
				"domainid": "2da0d0d3-e7b2-42ef-805d-eb2ea90ae7ef",
				"id": "eda6b916-5eea-4364-b699-ba1078c94f1c",
				"name": "dummy",
				"type": "host affinity"
			}
		],
		"count": 2
	}
}
	`)
	defer ts.Close()

	cs := NewClient(ts.URL, "TOKEN", "SECRET")
	affinityGroups, err := cs.ListAffinityGroups(&ListAffinityGroupsRequest{})
	if err != nil {
		log.Fatal(err)
	}

	if len(affinityGroups) != 2 {
		t.Errorf("Expected two groups")
	}

	if affinityGroups[0].Type != "host anti-affinity" {
		t.Errorf("Expected host affinity")
	}

	if affinityGroups[1].Type != "host affinity" {
		t.Errorf("Expected host anti-affinity")
	}
}

func TestListAffinityGroupTypes(t *testing.T) {
	ts := newServer(200, `
{
	"listaffinitygrouptypesresponse": {
		"affinityGroupType": [
			{
				"type": "host affinity"
			},
			{
				"type": "host anti-affinity"
			}
		],
		"count": 2
	}
}
	`)
	defer ts.Close()

	cs := NewClient(ts.URL, "TOKEN", "SECRET")
	affinityGroupTypes, err := cs.ListAffinityGroupTypes(&ListAffinityGroupTypesRequest{})
	if err != nil {
		log.Fatal(err)
	}

	if len(affinityGroupTypes) != 2 {
		t.Errorf("Expected two types")
	}

	if affinityGroupTypes[0].Type != "host affinity" {
		t.Errorf("Expected host affinity")
	}

	if affinityGroupTypes[1].Type != "host anti-affinity" {
		t.Errorf("Expected host anti-affinity")
	}
}
