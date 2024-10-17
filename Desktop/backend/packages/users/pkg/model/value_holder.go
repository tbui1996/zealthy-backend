package model

// The default value is only valid on a value holder if HasValue is true

type valueHolder struct {
	value    interface{}
	hasValue bool
	changed  bool
}

func (v *valueHolder) setValue(value interface{}) {
	if v.hasValue && v.value == value {
		return
	}

	v.value = value
	v.hasValue = true
	v.changed = true
}

func (v *valueHolder) IsEqual(other *valueHolder) bool {
	return v.value == other.value && v.hasValue == other.hasValue && v.changed == other.changed
}

func (v *valueHolder) setInitialValue(value interface{}) {
	v.value = value
	v.hasValue = true
}

func (v *valueHolder) HasValue() bool {
	return v.hasValue
}

func (v *valueHolder) Changed() bool {
	return v.changed
}

type StringValueHolder struct {
	valueHolder
}

func (v *StringValueHolder) SetValue(value string) {
	v.setValue(value)
}

func (v *StringValueHolder) SetInitialValue(value string) {
	v.setInitialValue(value)
}

func (v *StringValueHolder) GetValue() string {
	if v.value == nil {
		return ""
	}

	return v.value.(string)
}

type IntValueHolder struct {
	valueHolder
}

func (v *IntValueHolder) SetValue(value int) {
	v.setValue(value)
}

func (v *IntValueHolder) SetInitialValue(value int) {
	v.setInitialValue(value)
}

func (v *IntValueHolder) GetValue() int {
	if v.value == nil {
		return 0
	}

	return v.value.(int)
}

type BoolValueHolder struct {
	valueHolder
}

func (v *BoolValueHolder) SetValue(value bool) {
	v.setValue(value)
}

func (v *BoolValueHolder) SetInitialValue(value bool) {
	v.setInitialValue(value)
}

func (v *BoolValueHolder) GetValue() bool {
	if v.value == nil {
		return false
	}

	return v.value.(bool)
}

type ExternalUserOrganizationValueHolder struct {
	valueHolder
}

func (v *ExternalUserOrganizationValueHolder) SetValue(value *ExternalUserOrganization) {
	v.setValue(value)
}

func (v *ExternalUserOrganizationValueHolder) SetInitialValue(value *ExternalUserOrganization) {
	v.setInitialValue(value)
}

func (v *ExternalUserOrganizationValueHolder) GetValue() *ExternalUserOrganization {
	if v.value == nil {
		return nil
	}

	return v.value.(*ExternalUserOrganization)
}
