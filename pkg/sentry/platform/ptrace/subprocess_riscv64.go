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

package ptrace

import (
	"fmt"
	"strings"

	"golang.org/x/sys/unix"
	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/seccomp"
	"gvisor.dev/gvisor/pkg/sentry/arch"
)

const (
	// initRegsRipAdjustment is the size of the ecall instruction.
	initRegsRipAdjustment = 4
)

// resetSysemuRegs sets up emulation registers.
//
// This should be called prior to calling sysemu.
func (t *thread) resetSysemuRegs(regs *arch.Registers) {
}

// createSyscallRegs sets up syscall registers.
//
// This should be called to generate registers for a system call.
func createSyscallRegs(initRegs *arch.Registers, sysno uintptr, args ...arch.SyscallArgument) arch.Registers {
	// Copy initial registers (Pc, Sp, etc.).
	regs := *initRegs

	// Set our syscall number.
	// a7 for the syscall number.
	// a0-a5 is used to store the parameters.
	regs.Regs[17] = uint64(sysno)
	if len(args) >= 1 {
		regs.Regs[10] = args[0].Uint64()
	}
	if len(args) >= 2 {
		regs.Regs[11] = args[1].Uint64()
	}
	if len(args) >= 3 {
		regs.Regs[12] = args[2].Uint64()
	}
	if len(args) >= 4 {
		regs.Regs[13] = args[3].Uint64()
	}
	if len(args) >= 5 {
		regs.Regs[14] = args[4].Uint64()
	}
	if len(args) >= 6 {
		regs.Regs[15] = args[5].Uint64()
	}

	return regs
}

// isSingleStepping determines if the registers indicate single-stepping.
func isSingleStepping(regs *arch.Registers) bool {
	// return (regs.Pstate & arch.ARMTrapFlag) != 0
	// TODO: figure out where to check this
	return false
}

// updateSyscallRegs updates registers after finishing sysemu.
func updateSyscallRegs(regs *arch.Registers) {
	// No special work is necessary.
	return
}

// syscallReturnValue extracts a sensible return from registers.
func syscallReturnValue(regs *arch.Registers) (uintptr, error) {
	rval := int64(regs.Regs[10])
	if rval < 0 {
		return 0, unix.Errno(-rval)
	}
	return uintptr(rval), nil
}

func dumpRegs(regs *arch.Registers) string {
	var m strings.Builder

	fmt.Fprintf(&m, "Registers:\n")

	for i := 0; i < 32; i++ {
		fmt.Fprintf(&m, "\tRegs[%d]\t = %016x\n", i, regs.Regs[i])
	}

	return m.String()
}

// adjustInitregsRip adjust the current register RIP value to
// be just before the system call instruction execution
func (t *thread) adjustInitRegsRip() {
	t.initRegs.Regs[0] -= initRegsRipAdjustment
}

// Pass the expected PPID to the child via S7 when creating stub process
func initChildProcessPPID(initregs *arch.Registers, ppid int32) {
	initregs.Regs[23] = uint64(ppid)
	// S8 has to be set to 1 when creating stub process.
	initregs.Regs[24] = 1
}

// patchSignalInfo patches the signal info to account for hitting the seccomp
// filters from vsyscall emulation, specified below. We allow for SIGSYS as a
// synchronous trap, but patch the structure to appear like a SIGSEGV with the
// Rip as the faulting address.
//
// Note that this should only be called after verifying that the signalInfo has
// been generated by the kernel.
func patchSignalInfo(regs *arch.Registers, signalInfo *linux.SignalInfo) {
	if linux.Signal(signalInfo.Signo) == linux.SIGSYS {
		signalInfo.Signo = int32(linux.SIGSEGV)

		// Unwind the kernel emulation, if any has occurred. A SIGSYS is delivered
		// with the si_call_addr field pointing to the current RIP. This field
		// aligns with the si_addr field for a SIGSEGV, so we don't need to touch
		// anything there. We do need to unwind emulation however, so we set the
		// instruction pointer to the faulting value, and "unpop" the stack.
		regs.Regs[0] = signalInfo.Addr()
		regs.Regs[2] -= 8
	}
}

// Noop on riscv64.
//
//go:nosplit
func enableCpuidFault() {
}

// appendArchSeccompRules append architecture specific seccomp rules when creating BPF program.
// Ref attachedThread() for more detail.
func appendArchSeccompRules(rules []seccomp.RuleSet, defaultAction linux.BPFAction) []seccomp.RuleSet {
	return rules
}

// probeSeccomp returns true if seccomp is run after ptrace notifications,
// which is generally the case for kernel version >= 4.8.
//
// On riscv64, the support of PTRACE_SYSEMU , so
// probeSeccomp can always return true.
func probeSeccomp() bool {
	return true
}

func (s *subprocess) arm64SyscallWorkaround(t *thread, regs *arch.Registers) {
}
