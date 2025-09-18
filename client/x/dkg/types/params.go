package types

import (
	"github.com/piplabs/story/lib/errors"
)

const (
	// periods are in seconds.
	MinDkgStagePeriod            uint32 = 1 * 24 * 60 * 60  // 1 day
	DefaultDkgRegistrationPeriod uint32 = 1 * 24 * 60 * 60  // 1 day
	DefaultDkgChallengePeriod    uint32 = 7 * 24 * 60 * 60  // 7 days
	DefaultDkgDealingPeriod      uint32 = 1 * 24 * 60 * 60  // 1 day
	DefaultDkgFinalizationPeriod uint32 = 1 * 24 * 60 * 60  // 1 day
	DefaultDkgActivePeriod       uint32 = 21 * 24 * 60 * 60 // 21 days
	DefaultDkgComplaintPeriod    uint32 = 2 * 60 * 60       // 2 hours

	// other parameters.
	DefaultMinCommitteeSize uint32 = 3
	ExpectedMrenclaveSize   int    = 32 // 256-bit digest (32 bytes)
)

// NewParams creates a new Params instance.
func NewParams(
	registrationPeriod uint32,
	challengePeriod uint32,
	dealingPeriod uint32,
	finalizationPeriod uint32,
	activePeriod uint32,
	complaintPeriod uint32,
	minCommitteeSize uint32,
) Params {
	return Params{
		RegistrationPeriod: registrationPeriod,
		ChallengePeriod:    challengePeriod,
		DealingPeriod:      dealingPeriod,
		FinalizationPeriod: finalizationPeriod,
		ActivePeriod:       activePeriod,
		ComplaintPeriod:    complaintPeriod,
		MinCommitteeSize:   minCommitteeSize,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultDkgRegistrationPeriod,
		DefaultDkgChallengePeriod,
		DefaultDkgDealingPeriod,
		DefaultDkgFinalizationPeriod,
		DefaultDkgActivePeriod,
		DefaultDkgComplaintPeriod,
		DefaultMinCommitteeSize,
		// no default for mrenclave
	)
}

func (p Params) Validate() error {
	if err := ValidateRegistrationPeriod(p.RegistrationPeriod); err != nil {
		return err
	}

	if err := ValidateChallengePeriod(p.ChallengePeriod); err != nil {
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

	if err := ValidateComplaintPeriod(p.ComplaintPeriod); err != nil {
		return err
	}

	if err := ValidateMinCommitteeSize(p.MinCommitteeSize); err != nil {
		return err
	}

	if err := ValidateMrenclave(p.Mrenclave); err != nil {
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

func ValidateChallengePeriod(challengePeriod uint32) error {
	if challengePeriod == 0 {
		return errors.New("invalid dkg challenge period", "period", challengePeriod)
	}

	if challengePeriod < MinDkgStagePeriod {
		return errors.New("minimum dkg challenge period is 1 day", "period", challengePeriod)
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

func ValidateMinCommitteeSize(minCommitteeSize uint32) error {
	if minCommitteeSize == 0 {
		return errors.New("invalid min committee size", "size", minCommitteeSize)
	}

	return nil
}

func ValidateMrenclave(mrenclave []byte) error {
	if len(mrenclave) != ExpectedMrenclaveSize {
		return errors.New("mrenclave must be a 256-bit digest (32 bytes)",
			"expected_size", ExpectedMrenclaveSize,
			"actual_size", len(mrenclave))
	}

	return nil
}
