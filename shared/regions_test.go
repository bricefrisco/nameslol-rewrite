package shared

import "testing"

func TestNewRegions(t *testing.T) {
	r := NewRegions()
	if r == nil {
		t.Error("NewRegions returned nil")
	}
}

func TestRegions_ValidateValid(t *testing.T) {
	r := NewRegions()
	if !r.Validate("NA") {
		t.Error("NA should be valid")
	}
}

func TestRegions_ValidateInvalid(t *testing.T) {
	r := NewRegions()
	if r.Validate("invalid") {
		t.Error("invalid should be invalid")
	}
}

func TestRegions_GetValid(t *testing.T) {
	r := NewRegions()
	region, err := r.Get("NA")
	if err != nil {
		t.Error("NA should be valid")
	}
	if region != "na1" {
		t.Error("NA should be na1")
	}
}

func TestRegions_GetInvalid(t *testing.T) {
	r := NewRegions()
	_, err := r.Get("invalid")
	if err == nil {
		t.Error("invalid should be invalid")
	}
}

func TestRegions_GetAll(t *testing.T) {
	r := NewRegions()
	regions := r.GetAll()
	if regions == nil {
		t.Error("GetAll returned nil")
	}
	if len(regions) != 5 {
		t.Error("GetAll should return 5 regions")
	}
	if regions["NA"] != "na1" {
		t.Error("NA should be na1")
	}
	if regions["EUW"] != "euw1" {
		t.Error("EUW should be euw1")
	}
	if regions["EUNE"] != "eun1" {
		t.Error("EUNE should be eun1")
	}
	if regions["OCE"] != "oc1" {
		t.Error("OCE should be oc1")
	}
	if regions["LAS"] != "la2" {
		t.Error("LAS should be la2")
	}
}
