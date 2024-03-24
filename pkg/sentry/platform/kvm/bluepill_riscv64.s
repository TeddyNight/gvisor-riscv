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

#include "textflag.h"

// VCPU_CPU is the location of the CPU in the vCPU struct.
//
// This is guaranteed to be zero.
#define VCPU_CPU 0x0

// CPU_SELF is the self reference in ring0's percpu.
//
// This is guaranteed to be zero.
#define CPU_SELF 0x0

// Context offsets.
//
// Only limited use of the context is done in the assembly stub below, most is
// done in the Go handlers.
#define SIGINFO_SIGNO 0x0
#define SIGINFO_CODE 0x8
#define CONTEXT_PC  0x1B8
#define CONTEXT_R0 0xB8

#define SYS_MMAP 222

// getTLS returns the value of TPIDR_EL0 register.
TEXT ·getTLS(SB),NOSPLIT,$0-8
	RET

// setTLS writes the TPIDR_EL0 value.
TEXT ·setTLS(SB),NOSPLIT,$0-8
	RET

// See bluepill.go.
TEXT ·bluepill(SB),NOSPLIT,$0
	RET

// sighandler: see bluepill.go for documentation.
//
// The arguments are the following:
//
// 	R0 - The signal number.
// 	R1 - Pointer to siginfo_t structure.
// 	R2 - Pointer to ucontext structure.
//
TEXT ·sighandler(SB),NOSPLIT,$0
	RET

// func addrOfSighandler() uintptr
TEXT ·addrOfSighandler(SB), $0-8
	RET

// The arguments are the following:
//
// 	R0 - The signal number.
// 	R1 - Pointer to siginfo_t structure.
// 	R2 - Pointer to ucontext structure.
//
TEXT ·sigsysHandler(SB),NOSPLIT,$0
	RET

// func addrOfSighandler() uintptr
TEXT ·addrOfSigsysHandler(SB), $0-8
	RET

// dieTrampoline: see bluepill.go, bluepill_arm64_unsafe.go for documentation.
TEXT ·dieTrampoline(SB),NOSPLIT,$0
	RET

// func addrOfDieTrampoline() uintptr
TEXT ·addrOfDieTrampoline(SB), $0-8
	RET
