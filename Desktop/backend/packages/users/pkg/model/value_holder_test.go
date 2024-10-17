package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAllValueHolders(t *testing.T) {
	var v valueHolder

	str := StringValueHolder{valueHolder: v}

	str.SetInitialValue("hey")
	assert.Equal(t, "hey", str.GetValue())
	str.SetValue("bye")
	assert.Equal(t, "bye", str.GetValue())
	str.setInitialValue(nil)
	assert.Equal(t, "", str.GetValue())

	boool := BoolValueHolder{valueHolder: v}
	boool.SetInitialValue(false)
	assert.False(t, boool.GetValue())
	boool.SetValue(true)
	assert.True(t, boool.GetValue())
	boool.setInitialValue(nil)
	assert.Equal(t, false, boool.GetValue())

	ints := IntValueHolder{valueHolder: v}
	ints.SetInitialValue(0)
	assert.Equal(t, 0, ints.GetValue())
	ints.SetValue(1)
	assert.Equal(t, 1, ints.GetValue())
	ints.setInitialValue(nil)
	assert.Equal(t, 0, ints.GetValue())

	externalUserOrg := ExternalUserOrganizationValueHolder{valueHolder: v}
	externalUserOrg.setInitialValue(nil)
	assert.Nil(t, externalUserOrg.GetValue())
	externalUserOrg.SetInitialValue(&ExternalUserOrganization{ID: 1})
	assert.Equal(t, ExternalUserOrganization{ID: 1}, *externalUserOrg.GetValue())
	externalUserOrg.SetValue(&ExternalUserOrganization{ID: 2})
	assert.Equal(t, ExternalUserOrganization{ID: 2}, *externalUserOrg.GetValue())
}
