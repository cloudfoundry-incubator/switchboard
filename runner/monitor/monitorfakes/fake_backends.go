// This file was generated by counterfeiter
package monitorfakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/switchboard/domain"
	"github.com/cloudfoundry-incubator/switchboard/runner/monitor"
)

type FakeBackends struct {
	AllStub        func() <-chan domain.Backend
	allMutex       sync.RWMutex
	allArgsForCall []struct{}
	allReturns     struct {
		result1 <-chan domain.Backend
	}
	SetStateStub        func(backend domain.Backend, state bool)
	setStateMutex       sync.RWMutex
	setStateArgsForCall []struct {
		backend domain.Backend
		state   bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBackends) All() <-chan domain.Backend {
	fake.allMutex.Lock()
	fake.allArgsForCall = append(fake.allArgsForCall, struct{}{})
	fake.recordInvocation("All", []interface{}{})
	fake.allMutex.Unlock()
	if fake.AllStub != nil {
		return fake.AllStub()
	} else {
		return fake.allReturns.result1
	}
}

func (fake *FakeBackends) AllCallCount() int {
	fake.allMutex.RLock()
	defer fake.allMutex.RUnlock()
	return len(fake.allArgsForCall)
}

func (fake *FakeBackends) AllReturns(result1 <-chan domain.Backend) {
	fake.AllStub = nil
	fake.allReturns = struct {
		result1 <-chan domain.Backend
	}{result1}
}

func (fake *FakeBackends) SetState(backend domain.Backend, state bool) {
	fake.setStateMutex.Lock()
	fake.setStateArgsForCall = append(fake.setStateArgsForCall, struct {
		backend domain.Backend
		state   bool
	}{backend, state})
	fake.recordInvocation("SetState", []interface{}{backend, state})
	fake.setStateMutex.Unlock()
	if fake.SetStateStub != nil {
		fake.SetStateStub(backend, state)
	}
}

func (fake *FakeBackends) SetStateCallCount() int {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	return len(fake.setStateArgsForCall)
}

func (fake *FakeBackends) SetStateArgsForCall(i int) (domain.Backend, bool) {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	return fake.setStateArgsForCall[i].backend, fake.setStateArgsForCall[i].state
}

func (fake *FakeBackends) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.allMutex.RLock()
	defer fake.allMutex.RUnlock()
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeBackends) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ monitor.Backends = new(FakeBackends)