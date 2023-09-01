package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

type TestCall struct {
	Req        interface{}
	Resp       interface{}
	RespStatus int
}

type TestFlow struct {
	Calls []TestCall
}

var (
	testdataFilename = "testdata.json"
	mu               sync.Mutex
)

func ReadTestdata(fileName string) (*TestFlow, error) {
	tf := &TestFlow{}

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cntnt, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(cntnt, tf)
	if err != nil {
		return &TestFlow{}, nil
	}

	return tf, nil
}

func GetTestCall(callNr int, callResp interface{}) error {
	testflow, err := ReadTestdata(testdataFilename)
	if err != nil {
		return err
	}

	if len(testflow.Calls) <= callNr {
		return fmt.Errorf("no call %d", callNr)
	}

	callIfval := testflow.Calls[callNr]
	callIfvalJson, err := json.Marshal(callIfval)
	if err != nil {
		return err
	}

	err = json.Unmarshal(callIfvalJson, callResp)
	if err != nil {
		return err
	}

	return nil
}

func WriteTestdata(req, resp interface{}, respStatus int) error {
	mu.Lock()
	defer mu.Unlock()

	tf, err := ReadTestdata(testdataFilename)
	if err != nil {
		return err
	}

	tf.Calls = append(tf.Calls, TestCall{
		RespStatus: respStatus,
		Req:        req,
		Resp:       resp,
	})

	indentedJSON, err := json.MarshalIndent(tf, "", "    ")
	if err != nil {
		return err
	}

	f, err := os.OpenFile(testdataFilename, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(indentedJSON)
	return err
}
