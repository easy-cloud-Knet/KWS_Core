package virerr

import (
	"errors"
	"fmt"
	"testing"
)

func TestVirError_Error(t *testing.T) {
	if NoSuchDomain.Error() != "Domain Not Found" {
		t.Errorf("unexpected: %s", NoSuchDomain.Error())
	}
}

func TestErrorDescriptor_Error(t *testing.T) {
	err := ErrorGen(NoSuchDomain, fmt.Errorf("detail"))
	want := "(Error Type= 'Domain Not Found',\n Message='detail')"
	if err.Error() != want {
		t.Errorf("expected %q, got %q", want, err.Error())
	}
}

func TestErrorDescriptor_Is(t *testing.T) {
	err := ErrorGen(LackCapacityRAM, fmt.Errorf("oom"))

	if !errors.Is(err, LackCapacityRAM) {
		t.Error("expected errors.Is to match LackCapacityRAM")
	}
	if errors.Is(err, LackCapacityCPU) {
		t.Error("expected errors.Is NOT to match LackCapacityCPU")
	}
}

func TestErrorDescriptor_As(t *testing.T) {
	err := ErrorGen(InvalidUUID, fmt.Errorf("bad uuid"))

	var vt VirError
	if !errors.As(err, &vt) {
		t.Fatal("errors.As should succeed")
	}
	if vt != InvalidUUID {
		t.Errorf("expected InvalidUUID, got %s", vt)
	}
}

func TestErrorGen(t *testing.T) {
	detail := fmt.Errorf("inner")
	err := ErrorGen(SnapshotError, detail)

	desc, ok := err.(ErrorDescriptor)
	if !ok {
		t.Fatal("ErrorGen should return ErrorDescriptor")
	}
	if desc.ErrorType != SnapshotError {
		t.Errorf("wrong ErrorType: %s", desc.ErrorType)
	}
	if desc.Detail != detail {
		t.Errorf("wrong Detail: %v", desc.Detail)
	}
}

func TestErrorJoin_withErrorDescriptor(t *testing.T) {
	base := ErrorGen(DomainSearchError, fmt.Errorf("original"))
	joined := ErrorJoin(base, fmt.Errorf("appended"))

	desc, ok := joined.(ErrorDescriptor)
	if !ok {
		t.Fatal("ErrorJoin should return ErrorDescriptor")
	}
	if desc.ErrorType != DomainSearchError {
		t.Errorf("wrong ErrorType: %s", desc.ErrorType)
	}
	if !errors.Is(joined, DomainSearchError) {
		t.Error("errors.Is should still match after join")
	}
}

func TestErrorJoin_withPlainError(t *testing.T) {
	base := fmt.Errorf("plain error")
	joined := ErrorJoin(base, fmt.Errorf("extra"))

	if joined == nil {
		t.Fatal("ErrorJoin should not return nil")
	}
	if !errors.As(joined, new(VirError)) {
		t.Error("ErrorJoin on plain error should wrap into ErrorDescriptor")
	}
}
