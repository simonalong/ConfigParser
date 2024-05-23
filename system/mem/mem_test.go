package mem

import (
	"fmt"
	"testing"

	"github.com/simonalong/gole/system/common"
)

func skipIfNotImplementedErr(t *testing.T, err error) {
	if err == common.ErrNotImplementedError {
		t.Skip("not implemented")
	}
}

func TestVirtual_memory(t *testing.T) {
	v, err := VirtualMemory()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("error %v", err)
	}
	t.Log(v)

}

func TestSwap_memory(t *testing.T) {
	v, err := SwapMemory()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("error %v", err)
	}
	empty := &SwapMemoryStat{}
	if v == empty {
		t.Errorf("error %v", v)
	}

	t.Log(v)
}

func TestVirtualMemoryStat_String(t *testing.T) {
	v := VirtualMemoryStat{
		Total:       10,
		Available:   20,
		Used:        30,
		UsedPercent: 30.1,
		Free:        40,
	}
	e := `{"total":10,"available":20,"used":30,"usedPercent":30.1,"free":40,"active":0,"inactive":0,"wired":0,"laundry":0,"buffers":0,"cached":0,"writeback":0,"dirty":0,"writebacktmp":0,"shared":0,"slab":0,"sreclaimable":0,"sunreclaim":0,"pagetables":0,"swapcached":0,"commitlimit":0,"committedas":0,"hightotal":0,"highfree":0,"lowtotal":0,"lowfree":0,"swaptotal":0,"swapfree":0,"mapped":0,"vmalloctotal":0,"vmallocused":0,"vmallocchunk":0,"hugepagestotal":0,"hugepagesfree":0,"hugepagesize":0}`
	if e != fmt.Sprintf("%v", v) {
		t.Errorf("VirtualMemoryStat string is invalid: %v", v)
	}
}

func TestSwapMemoryStat_String(t *testing.T) {
	v := SwapMemoryStat{
		Total:       10,
		Used:        30,
		Free:        40,
		UsedPercent: 30.1,
		Sin:         1,
		Sout:        2,
		PgIn:        3,
		PgOut:       4,
		PgFault:     5,
		PgMajFault:  6,
	}
	e := `{"total":10,"used":30,"free":40,"usedPercent":30.1,"sin":1,"sout":2,"pgin":3,"pgout":4,"pgfault":5,"pgmajfault":6}`
	if e != fmt.Sprintf("%v", v) {
		t.Errorf("SwapMemoryStat string is invalid: %v", v)
	}
}
