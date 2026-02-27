package parsor

import (
	"testing"
)

func Test_UUIDValidator(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	invalidUUID := "123e4567-e89b-12d3-a456-42661417400Z"

	isValid := UUIDValidator(validUUID)
	if !isValid {
		t.Errorf("Expected valid UUID to be recognized as valid, got false")
	}

	isInvalid := UUIDValidator(invalidUUID)
	if isInvalid {
		t.Errorf("Expected invalid UUID to be recognized as invalid, got true")
	}

	// Test with path traversal characters
	pathTraversalUUID := "../123e4567-e89b-12d3-a456-426614174000"
	isPathTraversalValid := UUIDValidator(pathTraversalUUID)
	if isPathTraversalValid {
		t.Errorf("Expected UUID with path traversal characters to be recognized as invalid, got true")
	}

}

func Test_GetSafeFilePath(t *testing.T) {
	baseDir := "/safe/directory"
	invalidDir := "../unsafe/directory"
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	invalidUUID := "123e4567-e89b-12d3-a456-42661417400Z"

	safePath, isValid := GetSafeFilePath(baseDir, validUUID)
	if !isValid {
		t.Errorf("Expected valid UUID to be recognized as valid, got false")
	}
	expectedPath := baseDir + "/" + validUUID
	if safePath != expectedPath {
		t.Errorf("Expected safe path to be %s, got %s", expectedPath, safePath)
	}

	_, isInvalid := GetSafeFilePath(baseDir, invalidUUID)
	if isInvalid {
		t.Errorf("Expected invalid UUID to be recognized as invalid, got true")
	}

	_, isInvalidDir := GetSafeFilePath(invalidDir, validUUID)
	if isInvalidDir {
		t.Errorf("Expected path with invalid base directory to be recognized as invalid, got true")
	}
}
