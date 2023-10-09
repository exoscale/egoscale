package replayer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/exoscale/egoscale/v3/testing/recorder"
)

func InitializeReturnType[T any](myFn any) T {
	fn := reflect.ValueOf(myFn)

	if fn.Kind() != reflect.Func {
		panic(fmt.Sprintf("Not a function: %#v", myFn))
	}

	// Get the return type of the first return value of the function
	returnType := fn.Type().Out(0)

	if returnType.Kind() == reflect.Pointer {
		return reflect.New(returnType.Elem()).Interface().(T)
	}

	return reflect.New(returnType).Elem().Interface().(T)
}

func (replayer *Replayer) AssertArgs(expectedArgs recorder.CallParameters, args ...any) {
	actualArgsMap := recorder.ArgsToMap(args...)

	// to remove type information, we marshal and unmarshal the args
	// we are not printing error messages to panic to prevent leakage of secrets
	jsonValue, err := json.Marshal(actualArgsMap)
	if err != nil {
		panic("failed to marshal args")
	}

	var x recorder.CallParameters
	err = json.Unmarshal(jsonValue, &x)
	if err != nil {
		panic("failed to unmarshal args")
	}

	if !reflect.DeepEqual(expectedArgs, x) {
		panic(fmt.Sprintf("unequal args\n\nexpected:\n%#v\n\n\nactual:\n%#v\n", expectedArgs, actualArgsMap))
	}
}

type Replayer struct {
	Filename string
	CallNr   atomic.Int32
}

func (replayer *Replayer) GetTestCall(callResp interface{}, argsMap *recorder.CallParameters, returnErr *error) error {
	testflow, err := recorder.ReadRecording(replayer.Filename)
	if err != nil {
		return fmt.Errorf("error reading test data %w", err)
	}

	currentCallNr := replayer.CallNr.Load()
	replayer.CallNr.Add(1)
	if len(testflow.Calls) <= int(currentCallNr) {
		return fmt.Errorf("no call %d", currentCallNr)
	}

	callIfval := testflow.Calls[currentCallNr]
	*returnErr = callIfval.ReturnedError
	*argsMap = callIfval.Parameters
	callIfvalRespJson, err := json.Marshal(callIfval.ReturnedValue)
	if err != nil {
		return fmt.Errorf("error marshalling call %w", err)
	}

	err = json.Unmarshal(callIfvalRespJson, callResp)
	if err != nil {
		return fmt.Errorf("error unmarshalling call %w", err)
	}

	return nil
}
