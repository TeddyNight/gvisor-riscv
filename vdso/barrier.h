// Copyright 2018 The gVisor Authors.
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

#ifndef VDSO_BARRIER_H_
#define VDSO_BARRIER_H_

namespace vdso {

// Compiler Optimization barrier.
inline void barrier(void) { __asm__ __volatile__("" ::: "memory"); }

#if __x86_64__

inline void memory_barrier(void) {
  __asm__ __volatile__("mfence" ::: "memory");
}
inline void read_barrier(void) { barrier(); }
inline void write_barrier(void) { barrier(); }

#elif __aarch64__

inline void memory_barrier(void) {
  __asm__ __volatile__("dmb ish" ::: "memory");
}
inline void read_barrier(void) {
  __asm__ __volatile__("dmb ishld" ::: "memory");
}
inline void write_barrier(void) {
  __asm__ __volatile__("dmb ishst" ::: "memory");
}

#elif __riscv__

#define RISCV_FENCE(p, s) \
	__asm__ __volatile__ ("fence " #p "," #s : : : "memory")

inline void memory_barrier(void) {
    RISCV_FENCE(rw, rw);
}

inline void read_barrier(void) {
    RISCV_FENCE(r, r);
}

inline void write_barrier(void) {
    RISCV_FENCE(w, w);
}

#else
#error "unsupported architecture"
#endif

}  // namespace vdso

#endif  // VDSO_BARRIER_H_
