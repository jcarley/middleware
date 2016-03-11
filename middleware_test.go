package middleware

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Not recommended (https://golang.org/pkg/reflect/#Value.Pointer),
// but the best we can do.
func funcsEqual(f1, f2 interface{}) bool {
	val1 := reflect.ValueOf(f1)
	val2 := reflect.ValueOf(f2)
	return val1.Pointer() == val2.Pointer()
}

var testApp = MiddlewareHandlerFunc(func(env map[string]interface{}) {
	fmt.Println("")
})

// Tests creating a new chain
func TestNew(t *testing.T) {
	c1 := func(h MiddlewareHandler) MiddlewareHandler {
		return nil
	}
	c2 := func(h MiddlewareHandler) MiddlewareHandler {
		return nil
	}

	slice := []Constructor{c1, c2}

	chain := New(slice...)
	assert.True(t, funcsEqual(chain.constructors[0], slice[0]))
	assert.True(t, funcsEqual(chain.constructors[1], slice[1]))
}

func TestThenWorksWithNoMiddleware(t *testing.T) {
	assert.NotPanics(t, func() {
		chain := New()
		final := chain.Then(testApp)

		assert.True(t, funcsEqual(final, testApp))
	})
}

func TestThenTreatsNilAsDefaultMiddlewareHandler(t *testing.T) {
	chained := New().Then(nil)
	assert.Equal(t, chained, DefaultMiddlewareHandler)
}

func TestThenFuncTreatsNilAsDefaultMiddlewareHandler(t *testing.T) {
	chained := New().ThenFunc(nil)
	assert.Equal(t, chained, DefaultMiddlewareHandler)
}

func TestThenFuncConstructsMiddlewareHandlerFunc(t *testing.T) {
	fn := MiddlewareHandlerFunc(func(env map[string]interface{}) {
		env["result"] = "help"
	})
	env := make(map[string]interface{})
	chained := New().ThenFunc(fn)
	chained.Call(env)
	assert.Equal(t, "help", env["result"])
}
