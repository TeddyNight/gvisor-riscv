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

#include "funcdata.h"
#include "textflag.h"

// reference: local_flush_tlb_all_asid, mm/tlbflush.c
TEXT ·FlushTlbByASID(SB),NOSPLIT,$0-8
	MOV asid+0(FP), A1
	WORD $0x12b00073 // sfence.vma x0, a1
	RET

TEXT ·LocalFlushTlbByASID(SB),NOSPLIT,$0-8
	MOV asid+0(FP), A1
	WORD $0x12b00073 // sfence.vma x0, a1
	RET

TEXT ·GetFCSR(SB),NOSPLIT,$0-8
	WORD $0x003025f3 // frcsr a1
	MOV A1, value+0(FP)
	RET

TEXT ·SaveFCSR(SB),NOSPLIT,$0-8
	MOV value+0(FP), A1
	WORD $0x00359073 // fscsr a1 
	RET

TEXT ·SaveFpRegs(SB),NOSPLIT,$0-8
	MOV arg+0(FP), A1
	MOVD F0, 0(A1)
	MOVD F1, 8(A1)
	MOVD F2, 16(A1)
	MOVD F3, 24(A1)
	MOVD F4, 32(A1)
	MOVD F5, 40(A1)
	MOVD F6, 48(A1)
	MOVD F7, 56(A1)
	MOVD F8, 64(A1)
	MOVD F9, 72(A1)
	MOVD F10, 80(A1)
	MOVD F11, 88(A1)
	MOVD F12, 96(A1)
	MOVD F13, 104(A1)
	MOVD F14, 112(A1)
	MOVD F15, 120(A1)
	MOVD F16, 128(A1)
	MOVD F17, 136(A1)
	MOVD F18, 144(A1)
	MOVD F19, 152(A1)
	MOVD F20, 160(A1)
	MOVD F21, 168(A1)
	MOVD F22, 176(A1)
	MOVD F23, 184(A1)
	MOVD F24, 192(A1)
	MOVD F25, 200(A1)
	MOVD F26, 208(A1)
	MOVD F27, 216(A1)
	MOVD F28, 224(A1)
	MOVD F29, 232(A1)
	MOVD F30, 240(A1)
	MOVD F31, 248(A1)
	RET

TEXT ·LoadFpRegs(SB),NOSPLIT,$0-8
	MOV arg+0(FP), A1
	MOVD 0(A1), F0
	MOVD 8(A1), F1
	MOVD 16(A1), F2
	MOVD 24(A1), F3
	MOVD 32(A1), F4
	MOVD 40(A1), F5
	MOVD 48(A1), F6
	MOVD 56(A1), F7
	MOVD 64(A1), F8
	MOVD 72(A1), F9
	MOVD 80(A1), F10
	MOVD 88(A1), F11
	MOVD 96(A1), F12
	MOVD 104(A1), F13
	MOVD 112(A1), F14
	MOVD 120(A1), F15
	MOVD 128(A1), F16
	MOVD 136(A1), F17
	MOVD 144(A1), F18
	MOVD 152(A1), F19
	MOVD 160(A1), F20
	MOVD 168(A1), F21
	MOVD 176(A1), F22
	MOVD 184(A1), F23
	MOVD 192(A1), F24
	MOVD 200(A1), F25
	MOVD 208(A1), F26
	MOVD 216(A1), F27
	MOVD 224(A1), F28
	MOVD 232(A1), F29
	MOVD 240(A1), F30
	MOVD 248(A1), F31
	RET
