// Copyright 2020 The gVisor Authors.
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

package cpuid

const (
	// RISCV64FeatureFP indicates support for single and double precision
	// float point types.
	RISCV64FeatureFP Feature = iota

	// RISCV64FeatureASIMD indicates support for Advanced SIMD with single
	// and double precision float point arithmetic.
	RISCV64FeatureASIMD

	// RISCV64FeatureEVTSTRM indicates support for the generic timer
	// configured to generate events at a frequency of approximately
	// 100KHz.
	RISCV64FeatureEVTSTRM

	// RISCV64FeatureAES indicates support for AES instructions
	// (AESE/AESD/AESMC/AESIMC).
	RISCV64FeatureAES

	// RISCV64FeaturePMULL indicates support for AES instructions
	// (PMULL/PMULL2).
	RISCV64FeaturePMULL

	// RISCV64FeatureSHA1 indicates support for SHA1 instructions
	// (SHA1C/SHA1P/SHA1M etc).
	RISCV64FeatureSHA1

	// RISCV64FeatureSHA2 indicates support for SHA2 instructions
	// (SHA256H/SHA256H2/SHA256SU0 etc).
	RISCV64FeatureSHA2

	// RISCV64FeatureCRC32 indicates support for CRC32 instructions
	// (CRC32B/CRC32H/CRC32W etc).
	RISCV64FeatureCRC32

	// RISCV64FeatureATOMICS indicates support for atomic instructions
	// (LDADD/LDCLR/LDEOR/LDSET etc).
	RISCV64FeatureATOMICS

	// RISCV64FeatureFPHP indicates support for half precision float point
	// arithmetic.
	RISCV64FeatureFPHP

	// RISCV64FeatureASIMDHP indicates support for ASIMD with half precision
	// float point arithmetic.
	RISCV64FeatureASIMDHP

	// RISCV64FeatureCPUID indicates support for EL0 access to certain ID
	// registers is available.
	RISCV64FeatureCPUID

	// RISCV64FeatureASIMDRDM indicates support for SQRDMLAH and SQRDMLSH
	// instructions.
	RISCV64FeatureASIMDRDM

	// RISCV64FeatureJSCVT indicates support for the FJCVTZS instruction.
	RISCV64FeatureJSCVT

	// RISCV64FeatureFCMA indicates support for the FCMLA and FCADD
	// instructions.
	RISCV64FeatureFCMA

	// RISCV64FeatureLRCPC indicates support for the LDAPRB/LDAPRH/LDAPR
	// instructions.
	RISCV64FeatureLRCPC

	// RISCV64FeatureDCPOP indicates support for DC instruction (DC CVAP).
	RISCV64FeatureDCPOP

	// RISCV64FeatureSHA3 indicates support for SHA3 instructions
	// (EOR3/RAX1/XAR/BCAX).
	RISCV64FeatureSHA3

	// RISCV64FeatureSM3 indicates support for SM3 instructions
	// (SM3SS1/SM3TT1A/SM3TT1B).
	RISCV64FeatureSM3

	// RISCV64FeatureSM4 indicates support for SM4 instructions
	// (SM4E/SM4EKEY).
	RISCV64FeatureSM4

	// RISCV64FeatureASIMDDP indicates support for dot product instructions
	// (UDOT/SDOT).
	RISCV64FeatureASIMDDP

	// RISCV64FeatureSHA512 indicates support for SHA2 instructions
	// (SHA512H/SHA512H2/SHA512SU0).
	RISCV64FeatureSHA512

	// RISCV64FeatureSVE indicates support for Scalable Vector Extension.
	RISCV64FeatureSVE

	// RISCV64FeatureASIMDFHM indicates support for FMLAL and FMLSL
	// instructions.
	RISCV64FeatureASIMDFHM
)

var allFeatures = map[Feature]allFeatureInfo{
	RISCV64FeatureFP:       {"fp", true},
	RISCV64FeatureASIMD:    {"asimd", true},
	RISCV64FeatureEVTSTRM:  {"evtstrm", true},
	RISCV64FeatureAES:      {"aes", true},
	RISCV64FeaturePMULL:    {"pmull", true},
	RISCV64FeatureSHA1:     {"sha1", true},
	RISCV64FeatureSHA2:     {"sha2", true},
	RISCV64FeatureCRC32:    {"crc32", true},
	RISCV64FeatureATOMICS:  {"atomics", true},
	RISCV64FeatureFPHP:     {"fphp", true},
	RISCV64FeatureASIMDHP:  {"asimdhp", true},
	RISCV64FeatureCPUID:    {"cpuid", true},
	RISCV64FeatureASIMDRDM: {"asimdrdm", true},
	RISCV64FeatureJSCVT:    {"jscvt", true},
	RISCV64FeatureFCMA:     {"fcma", true},
	RISCV64FeatureLRCPC:    {"lrcpc", true},
	RISCV64FeatureDCPOP:    {"dcpop", true},
	RISCV64FeatureSHA3:     {"sha3", true},
	RISCV64FeatureSM3:      {"sm3", true},
	RISCV64FeatureSM4:      {"sm4", true},
	RISCV64FeatureASIMDDP:  {"asimddp", true},
	RISCV64FeatureSHA512:   {"sha512", true},
	RISCV64FeatureSVE:      {"sve", true},
	RISCV64FeatureASIMDFHM: {"asimdfhm", true},
}

func archFlagOrder(fn func(Feature)) {
	for i := 0; i < len(allFeatures); i++ {
		fn(Feature(i))
	}
}
