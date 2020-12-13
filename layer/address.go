/*
 * Copyright 2020 PeerGum
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package layer

import (
	"fmt"
	"math/rand"
)

const (
	AddressSize = 4
	DistMax     = 1<<15 - 1
)

type Address [AddressSize]uint16
type Distance uint16

var (
	AddressNone      = Address{0}
	BroadcastAddress = Address{65535}
)

//
// --- Distances
//

func (address Address) Distance(address2 Address) Distance {
	var d int16 = DistMax
	//var vLow, v2Low PeerAddress
	for i, v := range address {
		v2 := address2[i]
		//for j := range *peerAddress {
		// Abs returns the absolute value of x.
		di := Abs(Abs(int16(v)-32767) - Abs(int16(v2)-32767))
		//if di > DistMax/2 {
		//	di = DistMax - di
		//}
		if di < d {
			d = di
		}
		//}
	}
	if d > 0 {
	WORMHOLE_TEST:
		for i, v := range address {
			for j, v2 := range address2 {
				if v == v2 && i != j {
					d = 0
					break WORMHOLE_TEST
				}
			}
		}
	}
	return Distance(d)
}

func Abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

//
// --- Utils
//

func generateAddress() (address Address) {
	for i := range address {
		for ; address[i] == 0 || address[i] == 65535; {
			address[i] = uint16(rand.Uint32())
		}
	}
	return address
}

//
// --- Representation
//

func (address Address) String() (value string) {
	for i, a := range address {
		if i > 0 {
			value += ":"
		}
		if a!=0 {
			value += fmt.Sprintf("%4x", a)
		}
	}
	return value
}
