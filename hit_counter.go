package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type handlerHitCounter struct {
	hits    atomic.Int32
	next    http.Handler
	reset   *handlerHitCounterReset
	metrics *handlerHitCounterMetric
}

func (hc *handlerHitCounter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	hc.hits.Add(1)
	hc.next.ServeHTTP(rw, req)
}

//func (hc *handlerHitCounter) middlewareMetricsInc(next http.Handler) http.Handler {}

type handlerHitCounterReset struct {
	hc *handlerHitCounter
}
type handlerHitCounterMetric struct {
	hc *handlerHitCounter
}

func (hcRst *handlerHitCounterReset) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	hcRst.hc.hits.And(0)
	rw.WriteHeader(200)
}

func (hcMtrc *handlerHitCounterMetric) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	hits := hcMtrc.hc.hits.Load()
	outStr := fmt.Sprintf("Hits: %d\n", hits)
	fmt.Printf(" metric %s \n", outStr)
	rw.Write([]byte(outStr))
	rw.WriteHeader(200)

}

func newHandlerHitCounter(next http.Handler) *handlerHitCounter {
	ret := handlerHitCounter{next: next}
	rst := handlerHitCounterReset{hc: &ret}
	mtrc := handlerHitCounterMetric{hc: &ret}
	ret.reset = &rst
	ret.metrics = &mtrc
	return &ret
}
