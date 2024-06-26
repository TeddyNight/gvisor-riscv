// Copyright 2023 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build riscv64
// +build riscv64

package linux

// Only 4K page size is supported on riscv64. In this case, TASK_SIZE can
// be one of three values, corresponding to 3-level, 4-level and
// 5-level paging.
//
// The array has to be sorted in decreasing order.
// SV57 supports 5-level paging
var feasibleTaskSizes = []uintptr{0x100000000000000, 0x800000000000, 0x4000000000}
