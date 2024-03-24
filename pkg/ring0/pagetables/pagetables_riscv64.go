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

package pagetables
import (
	"sync/atomic"

	"gvisor.dev/gvisor/pkg/hostarch"
)

// Address constraints.
//
// The lowerTop and upperBottom currently apply to four-level pagetables;
// additional refactoring would be necessary to support five-level pagetables.
const (
	lowerTop    = 0x0000ffffffffffff
	upperBottom = 0xffff000000000000
	pteShift    = 12
	pmdShift    = 21
	pudShift    = 30
	pgdShift    = 39

	pteMask = 0x1ff << pteShift
	pmdMask = 0x1ff << pmdShift
	pudMask = 0x1ff << pudShift
	pgdMask = 0x1ff << pgdShift

	pteSize = 1 << pteShift
	pmdSize = 1 << pmdShift
	pudSize = 1 << pudShift
	pgdSize = 1 << pgdShift

	//ARM64 ttbr
	//ttbrASIDOffset = 48
	//ttbrASIDMask   = 0xff

	entriesPerPage = 512
)

// InitArch does some additional initialization related to the architecture.
//
// +checkescape:hard,stack
//
//go:nosplit
func (p *PageTables) InitArch(allocator Allocator) {
	if p.upperSharedPageTables != nil {
		p.cloneUpperShared()
	} else {
		p.archPageTables.root = p.Allocator.NewPTEs()
		p.archPageTables.rootPhysical = p.Allocator.PhysicalFor(p.archPageTables.root)
	}
}

// cloneUpperShared clone the upper from the upper shared page tables.
//
//go:nosplit
func (p *PageTables) cloneUpperShared() {
	if p.upperStart != upperBottom {
		panic("upperStart should be the same as upperBottom")
	}

	p.archPageTables.root = p.upperSharedPageTables.archPageTables.root
	p.archPageTables.rootPhysical = p.upperSharedPageTables.archPageTables.rootPhysical
}

// PTE is a page table entry.
type PTE uintptr
// PTEs is a collection of entries.
type PTEs [entriesPerPage]PTE
// archPageTables is architecture-specific data.
type archPageTables struct {
	// root is the pagetable root for kernel space.
	root *PTEs

	// rootPhysical is the cached physical address of the root.
	//
	// This is saved only to prevent constant translation.
	rootPhysical uintptr

	asid uint16
}

// MapOpts are x86 options.
type MapOpts struct {
	// AccessType defines permissions.
	AccessType hostarch.AccessType

	// Global indicates the page is globally accessible.
	Global bool

	// User indicates the page is a user page.
	User bool
}

// Bits in page table entries.
// Reference: 
// riscv-privileged-v1.10.pdf
// arch/riscv/include/asm/pgtable-bits.h
const (
	// R/W access permission
	//typeTable   = 0x3 << 1
	typeSect      = 0x1 << 0
	//typePage    = 0x3 << 1
	pteValid      = 0x1 << 0
	present       = pteValid
	//pteTableBit = 0x1 << 1
	//pteTypeMask = 0x3 << 0
	//present     = pteValid | pteTableBit
	readable    = 0x1 << 1
	writable    = 0x1 << 2
	executable  = 0x1 << 3
	user        = 0x1 << 4
	global      = 0x1 << 5
	accessed    = 0x1 << 6
	dirty       = 0x1 << 7

	typeTable   = 0x1 << 0
)

const (
	// include RSW
	optionMask = 0x2ff
	protDefault = present | accessed | user
)

// Address extracts the address. This should only be used if Valid returns true.
//
//go:nosplit
func (p *PTE) Address() uintptr {
	return atomic.LoadUintptr((*uintptr)(p)) &^ optionMask
}

// Set sets this PTE value.
//
// This does not change the sect page property.
//
//go:nosplit
func (p *PTE) Set(addr uintptr, opts MapOpts) {
	v := (addr &^ optionMask) | readable | protDefault

	if !opts.AccessType.Any() {
		// Leave as non-valid if no access is available.
		v &^= pteValid
	}

	if opts.Global {
		v |= global
	}

	if opts.AccessType.Execute {
		v |= executable
	} 

	if opts.AccessType.Write {
		v |= writable
	}

	if opts.User {
		v |= user
	} else {
		v = v &^ user
	}

	atomic.StoreUintptr((*uintptr)(p), v)
}

// Clear clears this PTE, including sect page information.
//
//go:nosplit
func (p *PTE) Clear() {
	atomic.StoreUintptr((*uintptr)(p), 0)
}

// Valid returns true iff this entry is valid.
//
//go:nosplit
func (p *PTE) Valid() bool {
	return atomic.LoadUintptr((*uintptr)(p))&present != 0
}

// Opts returns the PTE options.
//
// These are all options except Valid and Sect.
//
//go:nosplit
func (p *PTE) Opts() MapOpts {
	v := atomic.LoadUintptr((*uintptr)(p))

	return MapOpts{
		AccessType: hostarch.AccessType{
			Read:    true,
			Write:   v&writable != 0,
			Execute: v&executable != 0,
		},
		Global: v&global != 0,
		User:   v&user != 0,
	}
}

// setPageTable sets this PTE value and forces the write bit and sect bit to
// be cleared. This is used explicitly for breaking sect pages.
//
//go:nosplit
func (p *PTE) setPageTable(pt *PageTables, ptes *PTEs) {
	addr := pt.Allocator.PhysicalFor(ptes)
	if addr&^optionMask != addr {
		// This should never happen.
		panic("unaligned physical address!")
	}
	v := addr | typeTable | protDefault
	atomic.StoreUintptr((*uintptr)(p), v)
}
