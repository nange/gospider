package spider

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestTaskRule(t *testing.T) {
	r1, err := GetTaskRule("test_rule_not_exist")
	assert.Nil(t, r1)
	assert.Equal(t, ErrTaskRuleNotExist, errors.Cause(err))

	assert.Panics(t, func() {
		Register(nil)
	})
	assert.Panics(t, func() {
		Register(&TaskRule{})
	})
	assert.Panics(t, func() {
		Register(&TaskRule{Namespace: ""})
	})
	assert.Panics(t, func() {
		Register(&TaskRule{Namespace: ""})
	})
	assert.Panics(t, func() {
		Register(&TaskRule{
			Name: "test",
			Rule: &Rule{},
		})
	})
	assert.Panics(t, func() {
		Register(&TaskRule{
			Name: "test",
			Rule: &Rule{
				Head: func(ctx *Context) error {
					return nil
				},
				Nodes: map[int]*Node{
					0: {},
					2: {},
				},
			},
		})
	})

	Register(&TaskRule{
		Name:         "test_rule_name",
		Namespace:    "test_namespace",
		OutputFields: []string{"test_field"},
		Rule: &Rule{
			Head: func(ctx *Context) error {
				return nil
			},
			Nodes: map[int]*Node{
				0: {},
			},
		},
	})
	r2, err := GetTaskRule("test_rule_name")
	assert.NoError(t, err)
	assert.NotNil(t, r2)

	keys := GetTaskRuleKeys()
	assert.Contains(t, keys, "test_rule_name")
}
