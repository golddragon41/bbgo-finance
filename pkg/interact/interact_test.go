package interact

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseFuncArgsAndCall_NoErrorFunction(t *testing.T) {
	noErrorFunc := func(a string, b float64, c bool) error {
		assert.Equal(t, "BTCUSDT", a)
		assert.Equal(t, 0.123, b)
		assert.Equal(t, true, c)
		return nil
	}

	err := parseFuncArgsAndCall(noErrorFunc, []string{"BTCUSDT", "0.123", "true"})
	assert.NoError(t, err)
}

func Test_parseFuncArgsAndCall_ErrorFunction(t *testing.T) {
	errorFunc := func(a string, b float64) error {
		return errors.New("error")
	}

	err := parseFuncArgsAndCall(errorFunc, []string{"BTCUSDT", "0.123"})
	assert.Error(t, err)

}

func Test_parseCommand(t *testing.T) {
	args := parseCommand(`closePosition "BTC USDT" 3.1415926 market`)
	t.Logf("args: %+v", args)
	for i, a := range args {
		t.Logf("args(%d): %#v", i, a)
	}

	assert.Equal(t, 4, len(args))
	assert.Equal(t, "closePosition", args[0])
	assert.Equal(t, "BTC USDT", args[1])
	assert.Equal(t, "3.1415926", args[2])
	assert.Equal(t, "market", args[3])
}


type closePositionTask struct {
	symbol string
	percentage float64
	confirmed bool
}

type TestInteraction struct {
	closePositionTask closePositionTask
}

func (m *TestInteraction) Commands(interact *Interact) {
	interact.Command("closePosition", func() error {
		// send symbol options
		return nil
	}).Next(func(symbol string) error {
		// get symbol from user
		m.closePositionTask.symbol = symbol

		// send percentage options
		return nil
	}).Next(func(percentage float64) error {
		// get percentage from user
		m.closePositionTask.percentage = percentage

		// send confirmation
		return nil
	}).Next(func(confirmed bool) error {
		m.closePositionTask.confirmed = confirmed
		// call position close

		// reply result
		return nil
	})
}

func TestCustomInteraction(t *testing.T) {
	globalInteraction := New()
	testInteraction := &TestInteraction{}
	testInteraction.Commands(globalInteraction)

	err := globalInteraction.init()
	assert.NoError(t, err)

	err = globalInteraction.runCommand("closePosition")
	assert.NoError(t, err)

	assert.Equal(t, "closePosition_1", globalInteraction.curState)

	err = globalInteraction.handleResponse("BTCUSDT")
	assert.NoError(t, err)
	assert.Equal(t, "closePosition_2", globalInteraction.curState)

	err = globalInteraction.handleResponse("0.20")
	assert.NoError(t, err)
	assert.Equal(t, "closePosition_3", globalInteraction.curState)

	err = globalInteraction.handleResponse("true")
	assert.NoError(t, err)
	assert.Equal(t, "closePosition_4", globalInteraction.curState)

	assert.Equal(t, closePositionTask{
		symbol:     "BTCUSDT",
		percentage: 0.2,
		confirmed:  true,
	}, testInteraction.closePositionTask)
}
