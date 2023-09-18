package replayer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/exoscale/egoscale/v3/recorder"
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

var callNr atomic.Int32

func GetTestCall(callResp interface{}, returnErr *error) error {
	testflow, err := recorder.ReadTestdata(recorder.TestdataFilename)
	if err != nil {
		return fmt.Errorf("error reading test data %w", err)
	}

	currentCallNr := callNr.Load()
	callNr.Add(1)
	if len(testflow.Calls) <= int(currentCallNr) {
		return fmt.Errorf("no call %d", currentCallNr)
	}

	callIfval := testflow.Calls[currentCallNr]
	*returnErr = callIfval.ReturnedError
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
