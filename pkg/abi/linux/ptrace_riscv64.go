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

package linux

// PtraceRegs is the set of CPU registers exposed by ptrace. Source:
// syscall.PtraceRegs.
//
// +marshal
// +stateify savable
type PtraceRegs struct {
	// TODO: figure out why user_regs_struct miss some regs
	Regs   [36]uint64
}

// InstructionPointer returns the address of the next instruction to be
// executed.
func (p *PtraceRegs) InstructionPointer() uint64 {
	return p.Regs[0]
}

// StackPointer returns the address of the Stack pointer.
func (p *PtraceRegs) StackPointer() uint64 {
	return p.Regs[0];
}

// SetStackPointer sets the stack pointer to the specified value.
func (p *PtraceRegs) SetStackPointer(sp uint64) {
	p.Regs[0] = sp
}
