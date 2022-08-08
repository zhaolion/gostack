package cputil

import (
	"github.com/jinzhu/copier"
)

func ShallowCopy(toValue interface{}, fromValue interface{}) error {
	return copier.CopyWithOption(toValue, fromValue, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    false,
	})
}

func DeepCopy(toValue interface{}, fromValue interface{}) error {
	return copier.CopyWithOption(toValue, fromValue, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})
}
