package mock

import (
	"crypto/rand"
	"fmt"
)

// MockDataGenerator generates mock data for DKG testing.
type MockDataGenerator struct{}

// NewMockDataGenerator creates a new mock data generator.
func NewMockDataGenerator() *MockDataGenerator {
	return &MockDataGenerator{}
}

// GenerateMockMrenclave generates a mock mrenclave (mrenclave).
func (*MockDataGenerator) GenerateMockMrenclave() []byte {
	mrenclave := make([]byte, 32)
	_, err := rand.Read(mrenclave)
	if err != nil {
		for i := range mrenclave {
			mrenclave[i] = byte(i)
		}
	}

	return mrenclave
}

// GenerateMockPubKey generates a mock public key.
func (*MockDataGenerator) GenerateMockPubKey() []byte {
	pubKey := make([]byte, 33)
	pubKey[0] = 0x02
	_, err := rand.Read(pubKey[1:])
	if err != nil {
		for i := 1; i < len(pubKey); i++ {
			pubKey[i] = byte(i)
		}
	}

	return pubKey
}

// GenerateMockRemoteReport generates a mock remote report.
func (*MockDataGenerator) GenerateMockRemoteReport(validatorAddr string, round uint32, pubKey []byte) []byte {
	report := make([]byte, 432)
	copy(report[:len(validatorAddr)], validatorAddr)

	roundBytes := []byte(fmt.Sprintf("round_%d", round))
	copy(report[32:32+len(roundBytes)], roundBytes)

	if len(pubKey) > 0 {
		copy(report[64:64+len(pubKey)], pubKey)
	}

	_, err := rand.Read(report[96:])
	if err != nil {
		for i := 96; i < len(report); i++ {
			report[i] = byte((i * 7) % 256)
		}
	}

	return report
}

// GenerateMockCommitments generates mock VSS commitments.
func (*MockDataGenerator) GenerateMockCommitments(threshold uint32) []byte {
	commitmentSize := 33
	commitmentsLen := int(threshold) * commitmentSize
	commitments := make([]byte, commitmentsLen)

	for i := range int(threshold) {
		offset := i * commitmentSize
		commitments[offset] = 0x02
		_, err := rand.Read(commitments[offset+1 : offset+commitmentSize])
		if err != nil {
			for j := 1; j < commitmentSize; j++ {
				commitments[offset+j] = byte((i*j + 42) % 256)
			}
		}
	}

	return commitments
}

// GenerateMockSignature generates a mock ECDSA signature.
func (*MockDataGenerator) GenerateMockSignature() []byte {
	signature := make([]byte, 64)
	_, err := rand.Read(signature)
	if err != nil {
		for i := range signature {
			signature[i] = byte((i*13 + 67) % 256)
		}
	}

	return signature
}

// GenerateMockComplaintIndexes generates mock complaint indexes.
func (*MockDataGenerator) GenerateMockComplaintIndexes(dealerCount uint32, complaintCount int) []uint32 {
	if complaintCount <= 0 || dealerCount == 0 {
		return []uint32{}
	}

	if complaintCount > int(dealerCount) {
		complaintCount = int(dealerCount)
	}

	complaints := make([]uint32, complaintCount)
	for i := range complaintCount {
		complaints[i] = uint32((i * 3) % int(dealerCount))
	}

	return complaints
}
