package uuid

import "github.com/google/uuid"

func ValidateAndReturnUUID(u string) (*uuid.UUID, error) {
	parsed, err := uuid.Parse(u)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
