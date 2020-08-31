/*
Copyright 2020 The Netease Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package curve

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	maxVolIDLen = 128
)

/*
ComposeCSIID composes a CSI ID from passed in parameters.
	[length of user=1:4byte] + [-:1byte]
	[user] + [-:1byte]
	[volName]
*/
func composeCSIID(user, volName string) (string, error) {
	buf16 := make([]byte, 2)

	if (4 + 1 + len(user) + 1 + len(volName)) > maxVolIDLen {
		return "", fmt.Errorf("CSI ID encoding length overflow")
	}

	binary.BigEndian.PutUint16(buf16, uint16(len(user)))
	userLength := hex.EncodeToString(buf16)

	return strings.Join([]string{userLength, user, volName}, "-"), nil
}

func decomposeCSIID(composedCSIID string) (user string, volName string, err error) {
	buf16, err := hex.DecodeString(composedCSIID[0:4])
	if err != nil {
		return "", "", err
	}
	userLength := binary.BigEndian.Uint16(buf16)
	user = composedCSIID[5 : 5+userLength]
	volName = composedCSIID[6+userLength:]
	return user, volName, nil
}
