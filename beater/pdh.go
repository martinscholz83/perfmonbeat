package beater

/*
#include <windows.h>
#include <stdio.h>
#include <conio.h>
#include <pdh.h>
#include <pdhmsg.h>
#cgo LDFLAGS: -lpdh
*/
import "C"
import (
	"unsafe"

	"errors"

	"strconv"

	"github.com/elastic/beats/libbeat/common"
	"github.com/maddin2016/perfmonbeat/config"
)

type Handle struct {
	status      C.PDH_STATUS
	query       C.HQUERY
	counterType C.DWORD
	counters    []Counter
}

type Counter struct {
	counterName  string
	counter      C.HCOUNTER
	counterPath  *C.CHAR
	displayValue C.struct__PDH_FMT_COUNTERVALUE
}

func GetHandle(config []config.CounterConfig) (handle *Handle, err error) {
	q := &Handle{query: nil}
	q.status = C.PdhOpenQuery(nil, 1, &q.query)
	counters := make([]Counter, len(config))
	q.counters = counters
	for i, v := range config {
		counters[i] = Counter{counterPath: (*C.CHAR)(C.CString(v.Query)), counterName: v.Alias}
		defer C.free(unsafe.Pointer(counters[i].counterPath))
		q.status = C.PdhAddCounter(q.query, (*C.CHAR)(counters[i].counterPath), 0, &counters[i].counter)
		if q.status != C.ERROR_SUCCESS {
			err := errors.New("PdhAddCounter is failed for " + v.Alias)
			return nil, err
		}
	}

	return q, nil
}

func (q *Handle) ReadData() (data []common.MapStr, err error) {
	result := make([]common.MapStr, len(q.counters))
	q.status = C.PdhCollectQueryData(q.query)

	if q.status != C.ERROR_SUCCESS {
		sP := (*int)(unsafe.Pointer(&q.status))
		s := strconv.Itoa(*sP)
		err := errors.New("PdhCollectQueryData failed with status " + s)
		return nil, err
	}

	for i, v := range q.counters {
		q.status = C.PdhGetFormattedCounterValue(v.counter, C.PDH_FMT_DOUBLE, &q.counterType, &v.displayValue)
		if q.status != C.ERROR_SUCCESS {
			sP := (*int)(unsafe.Pointer(&q.status))
			s := strconv.Itoa(*sP)
			err := errors.New("PdhGetFormattedCounterValue failed with status " + s)
			return nil, err
		}

		doubleValue := (*float64)(unsafe.Pointer(&v.displayValue.anon0))

		val := common.MapStr{
			"name":  v.counterName,
			"value": *doubleValue,
		}
		result[i] = val
	}
	return result, nil
}
