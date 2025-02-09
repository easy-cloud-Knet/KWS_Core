package conn

import (
	"fmt"
)

// VirError는 VirError 타입의 문자열 기반 에러 정의
type VirError string

const (
	FaildDeEncoding     VirError = "Error Not Found"
	DomainSearchError VirError = "Error Serching Domain"
	NoSuchDomain     VirError = "Domain Not Found"  
	
	DomainGenerationError VirError ="error Generating Domain"
	LackCapacityRAM  VirError = "Not enough RAM"   // control
	LackCapacityCPU  VirError = "Not Enough CPU" // 
	LackCapacityHD   VirError = "Not Enough HardDisk" // 
	
	InvalidUUID      VirError = "Invalid UUID Provided"
	
	InvalidParameter VirError= "Invalid parameter entered"
	WrongParameter VirError= "Not validated parameter In"
	
	DomainStatusError  VirError ="Error Retreving Domain Status"
	HostStatusError  VirError ="Error Retreving Host Status"

	DeletionDomainError VirError= "Error Deleting Domain"
	DomainShutdownError VirError= "failed in Deleting domain"
)

// VirError는 error 인터페이스를 구현
func (ve VirError) Error() string {
	return string(ve)
}
// err 구조체 정의
type ErrorDescriptor struct {
	ErrorType    VirError `json:"error type"`
	Detail error    `json:"detail"`
}

// err 구조체의 Error() 메서드 구현
func (e ErrorDescriptor) Error() string {
	return fmt.Sprintf("(Error Type= '%s',\n Message='%s')",
		e.ErrorType, e.Detail.Error())
}

func (e ErrorDescriptor) Is(target error) bool {
	// target이 VirError 타입인지 확인
	n, ok := target.(VirError)
	if !ok {
		return false
	}
	return e.ErrorType == n
}
func (e ErrorDescriptor) As(target interface{}) bool {
	v, ok := target.(*VirError)
	if !ok {
		return false
	}
	*v = e.ErrorType
	return true
}

func ErrorGen(baseError VirError, detailError error) error{
	return ErrorDescriptor{
		ErrorType:baseError,
		Detail: detailError,
	}
}

func ErrorJoin(baseError error ,appendingError error) error{
	v,ok := baseError.(*ErrorDescriptor)
	if !ok{
		return ErrorGen(VirError(baseError.Error()), appendingError)
	}
	v.Detail=fmt.Errorf("%w %w", v.Detail, appendingError)
	return v
}

