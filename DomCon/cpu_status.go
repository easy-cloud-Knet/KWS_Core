package domCon

import "runtime"






type statusUpdate interface {
	allocate(number Status)
	deallocate(number Status)
	status() int
}

func (s * Status) deallocate(number Status){
	*s -= number
}

func (s * Status) allocate(number Status){
	*s += number
}
func (s Status) status() int{
	return int(s)
}

func (dls *domainListStatus) GetVCPUStatus() (int, int){
	dls.mutex.RLock()
	defer dls.mutex.RUnlock()
	return dls.vcpuAll.status(), dls.vcpuAllocated.status()
}

func (dls *domainListStatus) InitializeTotalVCPU(){
	dls.mutex.Lock()
	defer dls.mutex.Unlock()
	dls.vcpuAll= Status(runtime.NumCPU())
}


func (dls *domainListStatus) DeallocateVCPU(number Status){
	dls.mutex.Lock()
	defer dls.mutex.Unlock()
	dls.vcpuAllocated.deallocate(number)
}
func (dls *domainListStatus) AllocateVCPU(number Status	){
	dls.mutex.Lock()
	defer dls.mutex.Unlock()

	dls.vcpuAllocated.allocate(number)
}