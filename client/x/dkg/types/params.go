package types

import (
	"cosmossdk.io/math"

	"github.com/piplabs/story/lib/errors"
)

const (
	// Periods are in blocks. Conversion assumes 2.5s block time:
	//   1 day  = 86400s / 2.5s = 34560 blocks
	//   4 days = 345600s / 2.5s = 138240 blocks
	MinDkgStagePeriod            uint32 = 1      // 1 block
	DefaultDkgRegistrationPeriod uint32 = 34560  // 1 day
	DefaultDkgDealingPeriod      uint32 = 34560  // 1 day
	DefaultDkgFinalizationPeriod uint32 = 34560  // 1 day
	DefaultDkgActivePeriod       uint32 = 138240 // 4 days

	ExpectedCodeCommitmentSize int = 32 // 256-bit digest (32 bytes)

	// DKG committee size parameters (sourced from DKG.sol contract events).
	DefaultMinReqRegisteredParticipants uint32 = 5
	DefaultMinReqFinalizedParticipants  uint32 = 5
	DefaultOperationalThreshold         uint32 = 500 // 50% in basis points (out of 1000)
	OperationalThresholdBasis           uint32 = 1000

	// Decrypt request timeout in blocks.
	DefaultDecryptTimeout uint64 = 200
)

// DefaultDkgCommitteeRewardPortion is the default portion of UBI rewards
// allocated to the DKG committee (5%).
var DefaultDkgCommitteeRewardPortion = math.LegacyMustNewDecFromStr("0.05")

// NewParams creates a new Params instance.
func NewParams(
	registrationPeriod uint32,
	dealingPeriod uint32,
	finalizationPeriod uint32,
	activePeriod uint32,
	dkgCommitteeRewardPortion math.LegacyDec,
	minReqRegisteredParticipants uint32,
	minReqFinalizedParticipants uint32,
	operationalThreshold uint32,
	decryptTimeout uint64,
) Params {
	return Params{
		RegistrationPeriod:           registrationPeriod,
		DealingPeriod:                dealingPeriod,
		FinalizationPeriod:           finalizationPeriod,
		ActivePeriod:                 activePeriod,
		DkgCommitteeRewardPortion:    dkgCommitteeRewardPortion,
		MinReqRegisteredParticipants: minReqRegisteredParticipants,
		MinReqFinalizedParticipants:  minReqFinalizedParticipants,
		OperationalThreshold:         operationalThreshold,
		DecryptTimeout:               decryptTimeout,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultDkgRegistrationPeriod,
		DefaultDkgDealingPeriod,
		DefaultDkgFinalizationPeriod,
		DefaultDkgActivePeriod,
		DefaultDkgCommitteeRewardPortion,
		DefaultMinReqRegisteredParticipants,
		DefaultMinReqFinalizedParticipants,
		DefaultOperationalThreshold,
		DefaultDecryptTimeout,
	)
}

func (p Params) Validate() error {
	if err := ValidateRegistrationPeriod(p.RegistrationPeriod); err != nil {
		return err
	}

	if err := ValidateDealingPeriod(p.DealingPeriod); err != nil {
		return err
	}

	if err := ValidateFinalizationPeriod(p.FinalizationPeriod); err != nil {
		return err
	}

	if err := ValidateActivePeriod(p.ActivePeriod); err != nil {
		return err
	}

	if err := ValidateDkgCommitteeRewardPortion(p.DkgCommitteeRewardPortion); err != nil {
		return err
	}

	if err := ValidateMinReqRegisteredParticipants(p.MinReqRegisteredParticipants); err != nil {
		return err
	}

	if err := ValidateMinReqFinalizedParticipants(p.MinReqFinalizedParticipants); err != nil {
		return err
	}

	if err := ValidateOperationalThreshold(p.OperationalThreshold); err != nil {
		return err
	}

	if err := ValidateDecryptTimeout(p.DecryptTimeout); err != nil {
		return err
	}

	return nil
}

func ValidateRegistrationPeriod(registrationPeriod uint32) error {
	if registrationPeriod == 0 {
		return errors.New("invalid dkg registration period", "period", registrationPeriod)
	}

	if registrationPeriod < MinDkgStagePeriod {
		return errors.New("minimum dkg registration period is 1 day", "period", registrationPeriod)
	}

	return nil
}

func ValidateDealingPeriod(dealingPeriod uint32) error {
	if dealingPeriod == 0 {
		return errors.New("invalid dkg dealing period", "period", dealingPeriod)
	}

	if dealingPeriod < MinDkgStagePeriod {
		return errors.New("minimum dkg dealing period is 1 day", "period", dealingPeriod)
	}

	return nil
}

func ValidateFinalizationPeriod(finalizationPeriod uint32) error {
	if finalizationPeriod == 0 {
		return errors.New("invalid dkg finalization period", "period", finalizationPeriod)
	}

	if finalizationPeriod < MinDkgStagePeriod {
		return errors.New("minimum dkg finalization period is 1 day", "period", finalizationPeriod)
	}

	return nil
}

func ValidateActivePeriod(activePeriod uint32) error {
	if activePeriod == 0 {
		return errors.New("invalid dkg active period", "period", activePeriod)
	}

	if activePeriod < MinDkgStagePeriod {
		return errors.New("minimum dkg active period is 1 day", "period", activePeriod)
	}

	return nil
}

func ValidateComplaintPeriod(complaintPeriod uint32) error {
	if complaintPeriod == 0 {
		return errors.New("invalid dkg complaint period", "period", complaintPeriod)
	}

	// complaint period is not conditioned to the min stage period

	return nil
}

func ValidateMinReqRegisteredParticipants(v uint32) error {
	if v == 0 {
		return errors.New("min_req_registered_participants must be greater than zero", "value", v)
	}

	return nil
}

func ValidateMinReqFinalizedParticipants(v uint32) error {
	if v == 0 {
		return errors.New("min_req_finalized_participants must be greater than zero", "value", v)
	}

	return nil
}

func ValidateOperationalThreshold(v uint32) error {
	if v == 0 {
		return errors.New("operational_threshold must be greater than zero", "value", v)
	}

	if v > OperationalThresholdBasis {
		return errors.New("operational_threshold must not exceed basis (1000)", "value", v, "basis", OperationalThresholdBasis)
	}

	return nil
}

func ValidateDkgCommitteeRewardPortion(portion math.LegacyDec) error {
	if portion.IsNegative() {
		return errors.New("dkg committee reward portion must not be negative", "portion", portion.String())
	}

	if portion.GT(math.LegacyOneDec()) {
		return errors.New("dkg committee reward portion must not exceed 1.0", "portion", portion.String())
	}

	return nil
}

func ValidateDecryptTimeout(timeout uint64) error {
	if timeout == 0 {
		return errors.New("decrypt_timeout must be greater than zero", "value", timeout)
	}

	return nil
}

// CalculateThreshold computes the operational threshold from total participants
// and a basis-point threshold value (out of OperationalThresholdBasis = 1000).
// Returns ceil(total * operationalThresholdBps / 1000).
func CalculateThreshold(total, operationalThresholdBps uint32) uint32 {
	if total == 0 || operationalThresholdBps == 0 {
		return 0
	}

	// Use uint64 to avoid overflow: total * operationalThresholdBps could exceed uint32 max
	product := uint64(total) * uint64(operationalThresholdBps)
	threshold := uint32(product / uint64(OperationalThresholdBasis))

	// Round up if there is a remainder (ceiling division)
	if product%uint64(OperationalThresholdBasis) != 0 {
		threshold++
	}

	return threshold
}

func ValidateCodeCommitment(codeCommitment []byte) error {
	if len(codeCommitment) != ExpectedCodeCommitmentSize {
		return errors.New("codeCommitment must be a 256-bit digest (32 bytes)",
			"expected_size", ExpectedCodeCommitmentSize,
			"actual_size", len(codeCommitment))
	}

	return nil
}
