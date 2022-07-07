package config

import (
	"reflect"
	"testing"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/stretchr/testify/assert"
)

// type flagsChecker struct {
// 	fvp *FlagValuesProvider
// }

// func NewUsedFlagsChecker(flagsProvider *FlagValuesProvider) (checker UsedFlagsChecker) {
// 	checker = flagsChecker{
// 		fvp: flagsProvider,
// 	}
// 	return checker
// }

// func (fp flagsChecker) ChangedFlag(flagName string) (isFound bool) {
// 	return fp.fvp.ChangedFlag(flagName)
// }

func getStructElements(epv reflect.Value) (elementList map[string]bool) {
	elementList = make(map[string]bool, expectedMaxYAMLParams)
	for i := 0; i < epv.NumField(); i++ {
		paramTag := epv.Type().Field(i).Tag
		paramName := GetParamName(paramTag)
		elementList[paramName] = true
	}
	return elementList
}

func getFlagFields(reqs []ParameterRequirement) (fieldList map[string]bool) {
	reqList := make(map[string]bool, len(reqs))
	for _, req := range reqs {
		reqList[req.Field] = true
	}
	return reqList
}

var flagExceptions = map[string]bool{
	"no-prompt": true,
	"read-yaml": true,
}

func TestStructHasRequirements(t *testing.T) {
	cases := []struct {
		Name string
		Reqs []ParameterRequirement
	}{
		{
			Name: "Command",
			Reqs: GetCommandFlagSpecs(),
		},
		{
			Name: "Script",
			Reqs: GetScriptFlagSpecs(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			executeParams := &YAMLExecuteParams{}
			ep := reflect.ValueOf(executeParams).Elem()

			elements := getStructElements(ep)
			reqList := getFlagFields(tc.Reqs)

			for name := range reqList {
				if flagExceptions[name] {
					continue
				}
				if !elements[name] {
					t.Fatalf("%s not found", name)
				}
			}
		})
	}
}

func TestNoErrorOnGoodYAML(t *testing.T) {
	testFile := "../../../testdata/test1-ok.yaml"

	rawParams, err := ReadYAMLExecuteParams([]string{testFile}, nil)
	assert.NoError(t, err)

	vp := options.NewMapValuesProvider(rawParams)
	params := options.New(vp)

	assert.True(t, params.ReadBool("conc", false))
	assert.Equal(t, params.ReadString("command", ""), "ls")

	cids, found := params.Read("cids", []string{})
	assert.True(t, found)
	assert.Equal(t, cids, "cdeb33642b4b43caa13b73ce0045d388,7ca5718bd76f1bca7a5ee72660d3120c,42560923b8414a519c7a42047f251fb3")
}

func TestErrorOnBadYAML(t *testing.T) {
	testFile := "../../../testdata/test2-bad.yaml"

	_, err := ReadYAMLExecuteParams([]string{testFile}, nil)
	assert.Error(t, err)
}
