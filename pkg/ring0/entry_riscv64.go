// Copyright 2019 The gVisor Authors.
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

package ring0

func vectors()
func AddrOfVectors() uintptr

// start is the CPU entrypoint.
//
// The CPU state will be set to c.Registers().
func start()
func AddrOfStart() uintptr
func kernelExitToSupervisor()

func kernelExitToUser()

// Shutdown execution
func Shutdown()

func S_software_interrupt()	
