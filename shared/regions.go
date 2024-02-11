package shared

import "fmt"

type Regions struct {
	regions map[string]string
}

func NewRegions() *Regions {
	return &Regions{
		regions: map[string]string{
			"NA":   "na1",
			"EUW":  "euw1",
			"EUNE": "eun1",
			"OCE":  "oc1",
			"LAS":  "la2",
		},
	}
}

func (r *Regions) Validate(region string) bool {
	_, ok := r.regions[region]
	return ok
}

func (r *Regions) Get(region string) (string, error) {
	if r, ok := r.regions[region]; ok {
		return r, nil
	}
	return "", fmt.Errorf("region not found")
}

func (r *Regions) GetAll() map[string]string {
	return r.regions
}
