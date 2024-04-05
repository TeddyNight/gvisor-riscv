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
	//"fmt"
	"runtime"

	"golang.org/x/sys/unix"
	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/hostarch"
	"gvisor.dev/gvisor/pkg/ring0"
	"gvisor.dev/gvisor/pkg/ring0/pagetables"
	"gvisor.dev/gvisor/pkg/sentry/platform"
)

type vCPUArchState struct {
	// PCIDs is the set of PCIDs for this vCPU.
	//
	// This starts above fixedKernelPCID.
	PCIDs *pagetables.PCIDs
}

const (
	// fixedKernelPCID is a fixed kernel PCID used for the kernel page
	// tables. We must start allocating user PCIDs above this in order to
	// avoid any conflict (see below).
	fixedKernelPCID = 1

	// poolPCIDs is the number of PCIDs to record in the database. As this
	// grows, assignment can take longer, since it is a simple linear scan.
	// Beyond a relatively small number, there are likely few perform
	// benefits, since the TLB has likely long since lost any translations
	// from more than a few PCIDs past.
	poolPCIDs = 128
)

func (m *machine) mapUpperHalf(pageTable *pagetables.PageTables) {
	applyPhysicalRegions(func(pr physicalRegion) bool {
		pageTable.Map(
			hostarch.Addr(ring0.KernelStartAddress|pr.virtual),
			pr.length,
			pagetables.MapOpts{AccessType: hostarch.AnyAccess, Global: true},
			pr.physical)

		return true // Keep iterating.
	})
}

func archPhysicalRegions(physicalRegions []physicalRegion) []physicalRegion {
	return physicalRegions
}

// nonCanonical generates a canonical address return.
//
//go:nosplit
func nonCanonical(addr uint64, signal int32, info *linux.SignalInfo) (hostarch.AccessType, error) {
	*info = linux.SignalInfo{
		Signo: signal,
		Code:  linux.SI_KERNEL,
	}
	info.SetAddr(addr) // Include address.
	return hostarch.NoAccess, platform.ErrContextSignal
}

// isInstructionAbort returns true if it is an instruction abort.
//
//go:nosplit
func isInstructionAbort(code uint64) bool {
	/*
	value := (code & _ESR_ELx_EC_MASK) >> _ESR_ELx_EC_SHIFT
	return value == _ESR_ELx_EC_IABT_LOW
	*/
	return true;
}

// isWriteFault returns whether it is a write fault.
//
//go:nosplit
func isWriteFault(code uint64) bool {
	return true;
	/*
	if isInstructionAbort(code) {
		return false
	}

	return (code & _ESR_ELx_WNR) != 0
	*/
}

// fault generates an appropriate fault return.
//
//go:nosplit
func (c *vCPU) fault(signal int32, info *linux.SignalInfo) (hostarch.AccessType, error) {
	bluepill(c) // Probably no-op, but may not be.
	faultAddr := c.FaultAddr()
	code, user := c.ErrorCode()
	if !user {
		// The last fault serviced by this CPU was not a user
		// fault, so we can't reliably trust the faultAddr or
		// the code provided here. We need to re-execute.
		return hostarch.NoAccess, platform.ErrContextInterrupt
	}

	// Reset the pointed SignalInfo.
	*info = linux.SignalInfo{Signo: signal}
	info.SetAddr(uint64(faultAddr))
	accessType := hostarch.AccessType{}
	if signal == int32(unix.SIGSEGV) {
		accessType = hostarch.AccessType{
			Read:    !isWriteFault(uint64(code)),
			Write:   isWriteFault(uint64(code)),
			Execute: isInstructionAbort(uint64(code)),
		}
	}

	/*
	ret := code & _ESR_ELx_FSC
	switch ret {
	case _ESR_SEGV_MAPERR_L0, _ESR_SEGV_MAPERR_L1, _ESR_SEGV_MAPERR_L2, _ESR_SEGV_MAPERR_L3:
		info.Code = 1 //SEGV_MAPERR
	case _ESR_SEGV_ACCERR_L1, _ESR_SEGV_ACCERR_L2, _ESR_SEGV_ACCERR_L3, _ESR_SEGV_PEMERR_L1, _ESR_SEGV_PEMERR_L2, _ESR_SEGV_PEMERR_L3:
		info.Code = 2 // SEGV_ACCERR.
	default:
		info.Code = 2
	}
	*/

	return accessType, platform.ErrContextSignal
}

// getMaxVCPU get max vCPU number
func (m *machine) getMaxVCPU() {
	rmaxVCPUs := runtime.NumCPU()
	smaxVCPUs, _, errno := unix.RawSyscall(unix.SYS_IOCTL, uintptr(m.fd), KVM_CHECK_EXTENSION, _KVM_CAP_MAX_VCPUS)
	// compare the max vcpu number from runtime and syscall, use smaller one.
	if errno != 0 {
		m.maxVCPUs = rmaxVCPUs
	} else {
		if rmaxVCPUs < int(smaxVCPUs) {
			m.maxVCPUs = rmaxVCPUs
		} else {
			m.maxVCPUs = int(smaxVCPUs)
		}
	}
}
