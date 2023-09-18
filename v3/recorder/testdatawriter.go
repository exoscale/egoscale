package recorder

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type TestCall struct {
	FuncName    string
	Req         interface{}
	ReturnValue interface{}
	ReturnError error
}

type TestFlow struct {
	Calls []TestCall
}

var (
	TestdataFilename = "testdata.json"
	mu               sync.Mutex
)

func ReadTestdata(fileName string) (*TestFlow, error) {
	tf := &TestFlow{}

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return tf, nil
		}
		return nil, fmt.Errorf("error opening file for reading %w", err)
	}
	defer f.Close()

	cntnt, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading all test data %w", err)
	}

	err = json.Unmarshal(cntnt, tf)
	if err != nil {
		return &TestFlow{}, nil
	}

	return tf, nil
}

func WriteTestdata(funcName string, req, resp interface{}, respErr error) error {
	mu.Lock()
	defer mu.Unlock()

	tf, err := ReadTestdata(TestdataFilename)
	if err != nil {
		return fmt.Errorf("error reading test data before writing %w", err)
	}

	fmt.Printf("funcName: %v\n", funcName)
	tf.Calls = append(tf.Calls, TestCall{
		FuncName:    funcName,
		Req:         req,
		ReturnValue: resp,
		ReturnError: respErr,
	})

	indentedJSON, err := json.MarshalIndent(tf, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling with ident %w", err)
	}

	startGarble(indentedJSON)

	f, err := os.Create(TestdataFilename)
	if err != nil {
		return fmt.Errorf("error opening file for writing %w", err)
	}
	defer f.Close()

	_, err = f.Write(indentedJSON)
	if err != nil {
		return fmt.Errorf("error writing testdata %w", err)
	}

	return nil
}

func startGarble(jsonData []byte) ([]byte, error) {
	readAllGoFilesOnce.Do(readAllGoFiles)

	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	garble(data, "")

	return nil, nil
}

func garble(data interface{}, key string) {
	switch value := data.(type) {
	case map[string]interface{}:
		for key, val := range value {
			garble(val, key)
		}
	case []interface{}:
		for _, val := range value {
			garble(val, "")
		}
	case string:
		if !strings.Contains(allGoFiles, value) {
			fmt.Printf("garbling %q: %v\n", key, value)
		}
	}
}

func argsToMap(args ...any) map[int]any {
	ret := make(map[int]any, 0)
	for i, v := range args {
		ret[i] = v
	}

	return ret
}

var (
	readAllGoFilesOnce    sync.Once
	allGoFiles            string
	tokensToKeepUngarbled = []string{
		"success",
		"bucket",
		"restricted",
		"v2",
		"sos",
		"get-sos-object",
		"put-sos-object",
		"list-sos-object",
		"delete-sos-object",
	}
)

func readAllGoFiles() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			fmt.Println(file.Name())
			b, err := os.ReadFile(file.Name())
			if err != nil {
				panic(err)
			}

			allGoFiles += string(b)
		}
	}

	// add tokens we don't want to garble
	for _, token := range tokensToKeepUngarbled {
		allGoFiles += token + "\n"
	}
}
