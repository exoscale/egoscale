package recorder

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type CallParameters map[int]any

type Call struct {
	FunctionName  string
	Parameters    CallParameters
	ReturnedValue interface{}
	ReturnedError error
}

type Recording struct {
	Calls []Call
}

func ReadRecording(fileName string) (*Recording, error) {
	recording := &Recording{}

	recFile, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return recording, nil
		}
		return nil, fmt.Errorf("error opening file for reading %w", err)
	}
	defer recFile.Close()

	recContent, err := io.ReadAll(recFile)
	if err != nil {
		return nil, fmt.Errorf("error reading all test data %w", err)
	}

	err = json.Unmarshal(recContent, recording)
	if err != nil {
		return &Recording{}, nil
	}

	return recording, nil
}

func marshalIndent(obj any) ([]byte, error) {
	indentedJSON, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling with ident, not logging error to avoid leaking secrets errid:1")
	}

	return indentedJSON, nil
}

type Recorder struct {
	Filename  string
	mu        sync.Mutex
	FirstCall bool
	Garbler   *Garbler
}

func (recorder *Recorder) RecordCall(funcName string, parameters CallParameters, returnedValue interface{}, returnedError error) error {
	recorder.mu.Lock()
	defer recorder.mu.Unlock()

	var tf *Recording

	if recorder.FirstCall {
		tf = &Recording{}
		recorder.FirstCall = false
	} else {
		var err error
		tf, err = ReadRecording(recorder.Filename)
		if err != nil {
			return fmt.Errorf("error reading test data before writing %w", err)
		}
	}

	tf.Calls = append(tf.Calls, Call{
		FunctionName:  funcName,
		Parameters:    parameters,
		ReturnedValue: returnedValue,
		ReturnedError: returnedError,
	})

	indentedJSON, err := marshalIndent(tf)
	if err != nil {
		return err
	}

	indentedJSON, err = recorder.Garbler.garbleSecrets(indentedJSON)
	if err != nil {
		return err
	}

	f, err := os.Create(recorder.Filename)
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

type Garbler struct {
	garbledInts    map[int]int
	garbledStrings map[string]string
}

func NewGarbler() *Garbler {
	return &Garbler{
		garbledInts:    make(map[int]int),
		garbledStrings: make(map[string]string),
	}
}

func (g *Garbler) garbleSecrets(jsonData []byte) ([]byte, error) {
	readAllGoFilesOnce.Do(readAllGoFiles)

	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	garbledData := g.garble(data)
	if garbledData != nil {
		indentedJSON, err := marshalIndent(garbledData)
		if err != nil {
			return nil, err
		}

		return indentedJSON, nil
	}

	return nil, nil
}

var charSetAlphaNum = "abcdefghijklmnopqrstuvwxyz012346789"

func randIntRange(min int, max int) int {
	return rand.Intn(max-min) + min
}

func randStringFromCharSet(strlen int, charSet string) string {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[randIntRange(0, len(charSet))]
	}
	return string(result)
}

func randString(strlen int) string {
	return randStringFromCharSet(strlen, charSetAlphaNum)
}

func (g *Garbler) garbleInt(input int) int {
	if v, ok := g.garbledInts[input]; ok {
		return v
	}

	ret := rand.Int()
	g.garbledInts[input] = ret
	return ret
}

func (g *Garbler) garbleString(input string) string {
	if v, ok := g.garbledStrings[input]; ok {
		return v
	}

	ret := ""
	if _, err := uuid.Parse(input); err == nil {
		ret = uuid.NewString()
	}

	ret = randString(len(input))

	g.garbledStrings[input] = ret

	return ret
}

func (g *Garbler) garble(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	switch typedData := data.(type) {
	case map[string]interface{}:
		for key, val := range typedData {
			if key != "FunctionName" {
				typedData[key] = g.garble(val)
			}
		}

		return typedData
	case []interface{}:
		for ind, val := range typedData {
			typedData[ind] = g.garble(val)
		}

		return typedData
	case string:
		if !strings.Contains(allGoFiles, typedData) {
			typedData = g.garbleString(typedData)
		}

		return typedData
	case int:
		if !strings.Contains(allGoFiles, fmt.Sprint(typedData)) {
			typedData = g.garbleInt(typedData)
		}

		return typedData
	case bool:
		// booleans cannot be considered secrets

		return typedData
	}

	panic("uninspected value for garbling")
}

func ArgsToMap(args ...any) CallParameters {
	ret := make(CallParameters, 0)
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
