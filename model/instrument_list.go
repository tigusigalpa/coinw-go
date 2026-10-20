package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

// InstrumentList serializes the documented comma-separated, 20-instrument
// request field while retaining each instrument's settlement distinction.
type InstrumentList []Instrument

func (list InstrumentList) MarshalJSON() ([]byte, error) {
	if len(list) < 1 || len(list) > 20 {
		return nil, fmt.Errorf("instrument list requires 1..20 entries")
	}
	names := make([]string, len(list))
	for i, instrument := range list {
		if _, err := instrument.MarshalJSON(); err != nil {
			return nil, err
		}
		names[i] = instrument.String()
	}
	return json.Marshal(strings.Join(names, ","))
}
