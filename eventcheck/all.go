package eventcheck

import (
	"github.com/zilionixx/zilion-base/eventcheck/basiccheck"
	"github.com/zilionixx/zilion-base/eventcheck/epochcheck"
	"github.com/zilionixx/zilion-base/eventcheck/parentscheck"
	"github.com/zilionixx/zilion-base/inter/dag"
)

// Checkers is collection of all the checkers
type Checkers struct {
	Basiccheck   *basiccheck.Checker
	Epochcheck   *epochcheck.Checker
	Parentscheck *parentscheck.Checker
}

// Validate runs all the checks except ZilionBFT-related
func (v *Checkers) Validate(e dag.Event, parents dag.Events) error {
	if err := v.Basiccheck.Validate(e); err != nil {
		return err
	}
	if err := v.Epochcheck.Validate(e); err != nil {
		return err
	}
	if err := v.Parentscheck.Validate(e, parents); err != nil {
		return err
	}
	return nil
}
