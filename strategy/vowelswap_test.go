// Licensed to Typogenerator under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Typogenerator licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package strategy_test

import (
	"testing"

	"Shahriar-Sazid/typogenerator/strategy"
)

func TestVowelSwap(t *testing.T) {
	in := "zenithar"
	out, err := strategy.VowelSwap.Generate(in, "")
	if err != nil {
		t.Fail()
		t.Fatal("Error should not occurs !", err)
	}

	if len(out) != 15 {
		t.Errorf("invalid permutation count, expected %d, got %d", 15, len(out))
		t.FailNow()
	}

	for _, res := range out {
		if res == in {
			t.Fail()
			t.Fatal("Vowel swap should not swap a letter with itself!")
		}
	}
}
