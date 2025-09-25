package service

import (
	"encoding/json"
	"os"

	"github.com/gofrs/flock"
	"github.com/piplabs/story/client/dkg/types"
)

// NOTE: This is a workaround for connecting DKG module with DKG service.
// In v2 , DKG service should be integrated into the DKG module directly.

// These files will be created if not exist.
const dealsFilePath = "/tmp/dkg_deals.json"
const responsesFilePath = "/tmp/dkg_responses.json"

// AddDealsFile appends new deals to the deals file in a thread-safe manner using file locking.
func AddDealsFile(newDeals []*types.Deal) error {
	lock := flock.New(dealsFilePath + ".lock")
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	var deals []*types.Deal
	data, err := os.ReadFile(dealsFilePath)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &deals); err != nil {
			return err
		}
	}

	deals = append(deals, newDeals...)

	out, err := json.Marshal(deals)
	if err != nil {
		return err
	}
	return os.WriteFile(dealsFilePath, out, 0644)
}

// PopDealsFile reads and clears the deals file in a thread-safe manner using file locking.
func PopDealsFile() ([]*types.Deal, error) {
	lock := flock.New(dealsFilePath + ".lock")
	if err := lock.Lock(); err != nil {
		return nil, err
	}
	defer lock.Unlock()

	data, err := os.ReadFile(dealsFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*types.Deal{}, nil
		}
		return nil, err
	}

	var deals []*types.Deal
	if len(data) > 0 {
		if err := json.Unmarshal(data, &deals); err != nil {
			return nil, err
		}
	}

	if err := os.WriteFile(dealsFilePath, []byte{}, 0644); err != nil {
		return nil, err
	}

	return deals, nil
}

// AddResponsesFile appends new responses to the responses file in a thread-safe manner using file locking.
func AddResponsesFile(newResponses []*types.Response) error {
	lock := flock.New(responsesFilePath + ".lock")
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	var responses []*types.Response
	data, err := os.ReadFile(responsesFilePath)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &responses); err != nil {
			return err
		}
	}

	responses = append(responses, newResponses...)

	out, err := json.Marshal(responses)
	if err != nil {
		return err
	}
	return os.WriteFile(responsesFilePath, out, 0644)
}

// PopResponsesFile reads and clears the responses file in a thread-safe manner using file locking.
func PopResponsesFile() ([]*types.Response, error) {
	lock := flock.New(responsesFilePath + ".lock")
	if err := lock.Lock(); err != nil {
		return nil, err
	}
	defer lock.Unlock()

	data, err := os.ReadFile(responsesFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*types.Response{}, nil
		}
		return nil, err
	}

	var responses []*types.Response
	if len(data) > 0 {
		if err := json.Unmarshal(data, &responses); err != nil {
			return nil, err
		}
	}

	if err := os.WriteFile(responsesFilePath, []byte{}, 0644); err != nil {
		return nil, err
	}

	return responses, nil
}
