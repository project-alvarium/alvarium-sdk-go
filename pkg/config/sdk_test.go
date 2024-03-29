/*******************************************************************************
 * Copyright 2024 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package config

import (
	"encoding/json"
	"github.com/project-alvarium/alvarium-sdk-go/test"
	"os"
	"testing"
)

func TestSDKInfo_UnmarshalJSON(t *testing.T) {
	b, err := os.ReadFile("../../test/res/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestSDKInfo_AnnotationInvalid(t *testing.T) {
	b, err := os.ReadFile("../../test/res/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	cfg.Annotators = append(cfg.Annotators, "invalid")
	b, _ = json.Marshal(cfg)

	var x SdkInfo
	err = json.Unmarshal(b, &x)
	test.CheckError(err, true, "test sdk invalid annotation", t)
}
