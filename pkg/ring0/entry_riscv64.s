#include "funcdata.h"
#include "go_asm.h"
#include "textflag.h"

#define SRET WORD $0x10500073

TEXT ·start(SB),NOSPLIT,$0
	JMP ·kernelExitToSupervisor(SB)

// func AddrOfStart() uintptr
TEXT ·AddrOfStart(SB), $0-8
	MOV	$·start(SB), A0
	MOV	A0, ret+0(FP)
	RET

// storeAppASID writes the application's asid value.
TEXT ·storeAppASID(SB),NOSPLIT,$0-8
	MOV asid+0(FP), A1
	WORD $0x18059073 // csrw satp,a1
	RET

// See kernel.go.
TEXT ·Halt(SB),NOSPLIT,$0
	WORD $0x10500073 //WFI
	RET

TEXT ·kernelExitToSupervisor(SB),NOSPLIT,$0
	SRET

TEXT ·kernelExitToUser(SB),NOSPLIT,$0
	SRET

TEXT ·S_software_interrupt(SB),NOSPLIT,$0
	SRET

// vectors implements exception vector table.
// The start address of exception vector table should be 12-bits aligned.
TEXT ·vectors(SB),NOSPLIT,$0
	PCALIGN $64
	SRET // USER SOFTWARE INTERRUPT
	PCALIGN $4
	JMP ·S_software_interrupt(SB)	
	PCALIGN $4
	 // RESERVED
	PCALIGN $4
	SRET // MACHINE SOFTWARE INTERRUPT
	PCALIGN $4
	SRET // USER TIMER INTERRUPT
	PCALIGN $4
	SRET // SUPERVISOR TIMER INTERRUPT
	PCALIGN $4
	SRET // RESERVED
	PCALIGN $4
	SRET // MACHINE TIMER INTERRUPT
	PCALIGN $4
	SRET // USER EXTERNAL INTERRUPT
	PCALIGN $4
	SRET // SUPERVISOR EXTERNAL INTERRUPT
	PCALIGN $4
	SRET // RESERVED
	PCALIGN $4
	SRET // MACHINE EXTERNAL INTERRUPT


// func AddrOfVectors() uintptr
TEXT ·AddrOfVectors(SB), $0-8
	MOV    $·vectors(SB), A0
	MOV    A0, ret+0(FP)
	RET
