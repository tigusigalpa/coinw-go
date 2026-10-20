package model

import "fmt"

func (v *FuturesDepthEvent) UnmarshalJSON(raw []byte) error {
	type documentedDepth FuturesDepthEvent
	var next documentedDepth
	if err := Decode(raw, &next); err != nil {
		return err
	}
	if next.Ask == nil && next.Asks == nil {
		return fmt.Errorf("depth requires its documented ask/asks alternative")
	}
	*v = FuturesDepthEvent(next)
	return nil
}
