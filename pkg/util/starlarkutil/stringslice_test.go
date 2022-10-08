// Copyright 2022 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package starlarkutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

func TestToStringSliceNil(t *testing.T) {
	nilList, err := ToStringSlice(nil)
	assert.Empty(t, nilList)
	assert.Nil(t, err)
}

func TestToStringSliceEmptyList(t *testing.T) {
	emptySlice := starlark.NewList(make([]starlark.Value, 0))
	resultSlice, err := ToStringSlice(emptySlice)
	assert.Empty(t, resultSlice)
	assert.Nil(t, err)
}

func TestToStringSliceInvalidType(t *testing.T) {
	boolSlice := []starlark.Value{starlark.False, starlark.True}
	emptyList := starlark.NewList(boolSlice)
	resultSlice, err := ToStringSlice(emptyList)
	assert.Empty(t, resultSlice)
	assert.ErrorContains(t, err, "False as bool")
}

func TestToStringSliceInvalidTypeinMiddle(t *testing.T) {
	mixedInvalidSlice := []starlark.Value{starlark.String(""), starlark.MakeUint(3)}
	mixedInvalidList := starlark.NewList(mixedInvalidSlice)
	resultSlice, err := ToStringSlice(mixedInvalidList)
	assert.Empty(t, resultSlice)
	assert.ErrorContains(t, err, "3 as int")
}

func TestToStringSliceEmptyListofString(t *testing.T) {
	emptyStringSlice := []starlark.Value{starlark.String(""), starlark.String("")}
	emptyStringList := starlark.NewList(emptyStringSlice)
	resultSlice, err := ToStringSlice(emptyStringList)
	assert.Contains(t, resultSlice, "")
	assert.NotEmpty(t, resultSlice)
	assert.Nil(t, err)
}

func TestToStringSliceHappy(t *testing.T) {
	wordsSlice := []starlark.Value{starlark.String("你好"), starlark.String("Χαίρετε")}
	wordsList := starlark.NewList(wordsSlice)
	resultSlice, err := ToStringSlice(wordsList)
	assert.Contains(t, resultSlice[1], "Χ")
	assert.Equal(t, len(resultSlice), 2)
	assert.Nil(t, err)
}
