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

package kvm

import (
	"golang.org/x/sys/unix"
	"gvisor.dev/gvisor/pkg/ring0"
	"gvisor.dev/gvisor/pkg/sentry/arch"
)

var (
	// The action for bluepillSignal is changed by sigaction().
	bluepillSignal = unix.SIGILL
)

// bluepillArchEnter is called during bluepillEnter.
//
//go:nosplit
func bluepillArchEnter(context *arch.SignalContext64) (c *vCPU) {
	c = vCPUPtr(uintptr(context.Regs[8]))
	regs := c.CPU.Registers()
	regs.Regs = context.Regs
	/*
	regs.Regs[32] &^= uint64(ring0.PsrFlagsClear)
	regs.Regs[32] |= ring0.KernelFlagSet
	*/

	return
}

// bluepillArchExit is called during bluepillEnter.
//
//go:nosplit
func bluepillArchExit(c *vCPU, context *arch.SignalContext64) {
	regs := c.CPU.Registers()
	context.Regs = regs.Regs
	/*
	regs.Regs[32] &^= uint64(ring0.PsrFlagsClear)
	regs.Regs[32] |= ring0.UserFlagSet
	*/

	/*
	lazyVfp := c.GetLazyVFP()
	if lazyVfp != 0 {
		fpsimd := fpsimdPtr(c.FloatingPointState().BytePointer()) // escapes: no
		context.Fpsimd64.Fpsr = fpsimd.Fpsr
		context.Fpsimd64.Fpcr = fpsimd.Fpcr
		context.Fpsimd64.Vregs = fpsimd.Vregs
	}
	*/
}

// KernelSyscall handles kernel syscalls.
//
// +checkescape:all
//
//go:nosplit
func (c *vCPU) KernelSyscall() {
	regs := c.Registers()
	if regs.Regs[17] != ^uint64(0) {
		regs.Regs[0] -= 4 // Rewind.
	}

	/*
	fpDisableTrap := ring0.CPACREL1()
	if fpDisableTrap != 0 {
		fpsimd := fpsimdPtr(c.FloatingPointState().BytePointer()) // escapes: no
		fpcr := ring0.GetFPCR()
		fpsr := ring0.GetFPSR()
		fpsimd.Fpcr = uint32(fpcr)
		fpsimd.Fpsr = uint32(fpsr)
		ring0.SaveVRegs(c.FloatingPointState().BytePointer()) // escapes: no
	}
	*/

	ring0.Halt()
}

// KernelException handles kernel exceptions.
//
// +checkescape:all
//
//go:nosplit
func (c *vCPU) KernelException(vector ring0.Vector) {
	regs := c.Registers()
	if vector == ring0.Vector(bounce) {
		regs.Regs[0] = 0
	}

	/*
	fpDisableTrap := ring0.CPACREL1()
	if fpDisableTrap != 0 {
		fpsimd := fpsimdPtr(c.FloatingPointState().BytePointer()) // escapes: no
		fpcr := ring0.GetFPCR()
		fpsr := ring0.GetFPSR()
		fpsimd.Fpcr = uint32(fpcr)
		fpsimd.Fpsr = uint32(fpsr)
		ring0.SaveVRegs(c.FloatingPointState().BytePointer()) // escapes: no
	}
	*/

	ring0.Halt()
}

// hltSanityCheck verifies the current state to detect obvious corruption.
//
//go:nosplit
func (c *vCPU) hltSanityCheck() {
}
