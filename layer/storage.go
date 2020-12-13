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
	"errors"
	"os"
	"strings"
)

type StorageBlock struct {
	name string
	path string
	size int
	file os.File
}

type StoredData struct {
	blocks []StorageBlock
}

type Storage struct {
	shared  StoredData
	private StoredData
}

var storage Storage

//
// --- Storage Data
//

func (storageData StoredData) AddBlock(block StorageBlock) {
	storageData.blocks = append(storageData.blocks, block)
}

func (storageData StoredData) Search(name string) (*StorageBlock, error) {
	for _, block := range storageData.blocks {
		if strings.Compare(block.name, name) == 0 {
			return &block, nil
		}
	}
	return nil, errors.New("Block Not Found")
}

