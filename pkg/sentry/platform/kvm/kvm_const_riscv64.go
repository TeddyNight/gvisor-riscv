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

package kvm

// KVM ioctls for Arm64.
const (
	_KVM_GET_ONE_REG = 0x4010aeab
	_KVM_SET_ONE_REG = 0x4010aeac

	_KVM_RISCV_REG_TYPE_SHIFT = 24
	_KVM_RISCV_REGS		= 0x8000000000000000
	_KVM_RISCV_REGS_CORE	= 0x02 << _KVM_RISCV_REG_TYPE_SHIFT
	_KVM_RISCV64_REG_SIZE	= 1 << 6
	_KVM_RISCV64_REGS_PC	= 0x8000000002000040
	_KVM_RISCV64_REGS_SP	= 0x8000000002000042
	_KVM_RISCV64_REGS_TP	= 0x8000000002000044
	_KVM_RISCV64_REGS_SIE	= 0x8000000003000041
	_KVM_RISCV64_REGS_STVEC	= 0x8000000003000042
	_KVM_RISCV64_REGS_SEPC	= 0x8000000003000044
	_KVM_RISCV64_REGS_SATP	= 0x8000000003000048
)

// Riscv64: Supervisor Interrupt Enable register
const (
	_SIE_SSIE = 1 << 1
	_SIE_UTIE = 1 << 4
	//_SIE_SEIE = 1 << 9
	_SIE_UEIE = 1 << 8
	_SIE_DEFAULT = _SIE_SSIE | _SIE_UTIE | _SIE_UEIE
)
