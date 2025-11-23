package domStatus

import (
	"sync"
	"testing"
)



func TestVCPUAtomic(t *testing.T){
	dls := &DomainListStatus{}

	wg:= sync.WaitGroup{}
	
	wg.Add(1000)
	for i:=0; i<1000; i++{
		go func(){
			dls.AddAllocatedCPU(4)
			dls.AddSleepingCPU(1)
			defer wg.Done()
		}()
	}
	wg.Wait()

	result:= dls.VcpuAllocated
	if result != 4000{
		t.Errorf("Expected 4000 allocated CPUs, got %d", result)
	}
	result1:=dls.VcpuSleeping
	if result1 != 1000{
		t.Errorf("Expected 1000 sleeping CPUs, got %d", result)
	}
	wg.Add(450)

	for i:=0; i<450; i++{
	go func(){
		dls.TakeAllocatedCPU(4)
		dls.TakeSleepingCPU(1)
		defer wg.Done()
		}()
	}
	wg.Wait()
	finalAllocated:= dls.VcpuAllocated
	if finalAllocated != 2200{
		t.Errorf("Expected 2200 allocated CPUs after taking, got %d", finalAllocated)
	}
	finalSleeping:= dls.VcpuSleeping
	if finalSleeping != 550{
		t.Errorf("Expected 550 sleeping CPUs after taking, got %d", finalSleeping)
	}
}