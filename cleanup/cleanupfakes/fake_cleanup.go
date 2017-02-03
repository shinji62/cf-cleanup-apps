// This file was generated by counterfeiter
package cleanupfakes

import (
	"sync"

	"github.com/shinji62/cf-cleanup-apps/cleanup"
)

type FakeCleanup struct {
	StopAppStub                func()
	stopAppMutex               sync.RWMutex
	stopAppArgsForCall         []struct{}
	DryRunStub                 func()
	dryRunMutex                sync.RWMutex
	dryRunArgsForCall          []struct{}
	ListExpiredAppsStub        func() ([]cleanup.App, error)
	listExpiredAppsMutex       sync.RWMutex
	listExpiredAppsArgsForCall []struct{}
	listExpiredAppsReturns     struct {
		result1 []cleanup.App
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCleanup) StopApp() {
	fake.stopAppMutex.Lock()
	fake.stopAppArgsForCall = append(fake.stopAppArgsForCall, struct{}{})
	fake.recordInvocation("StopApp", []interface{}{})
	fake.stopAppMutex.Unlock()
	if fake.StopAppStub != nil {
		fake.StopAppStub()
	}
}

func (fake *FakeCleanup) StopAppCallCount() int {
	fake.stopAppMutex.RLock()
	defer fake.stopAppMutex.RUnlock()
	return len(fake.stopAppArgsForCall)
}

func (fake *FakeCleanup) DryRun() {
	fake.dryRunMutex.Lock()
	fake.dryRunArgsForCall = append(fake.dryRunArgsForCall, struct{}{})
	fake.recordInvocation("DryRun", []interface{}{})
	fake.dryRunMutex.Unlock()
	if fake.DryRunStub != nil {
		fake.DryRunStub()
	}
}

func (fake *FakeCleanup) DryRunCallCount() int {
	fake.dryRunMutex.RLock()
	defer fake.dryRunMutex.RUnlock()
	return len(fake.dryRunArgsForCall)
}

func (fake *FakeCleanup) ListExpiredApps() ([]cleanup.App, error) {
	fake.listExpiredAppsMutex.Lock()
	fake.listExpiredAppsArgsForCall = append(fake.listExpiredAppsArgsForCall, struct{}{})
	fake.recordInvocation("ListExpiredApps", []interface{}{})
	fake.listExpiredAppsMutex.Unlock()
	if fake.ListExpiredAppsStub != nil {
		return fake.ListExpiredAppsStub()
	}
	return fake.listExpiredAppsReturns.result1, fake.listExpiredAppsReturns.result2
}

func (fake *FakeCleanup) ListExpiredAppsCallCount() int {
	fake.listExpiredAppsMutex.RLock()
	defer fake.listExpiredAppsMutex.RUnlock()
	return len(fake.listExpiredAppsArgsForCall)
}

func (fake *FakeCleanup) ListExpiredAppsReturns(result1 []cleanup.App, result2 error) {
	fake.ListExpiredAppsStub = nil
	fake.listExpiredAppsReturns = struct {
		result1 []cleanup.App
		result2 error
	}{result1, result2}
}

func (fake *FakeCleanup) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.stopAppMutex.RLock()
	defer fake.stopAppMutex.RUnlock()
	fake.dryRunMutex.RLock()
	defer fake.dryRunMutex.RUnlock()
	fake.listExpiredAppsMutex.RLock()
	defer fake.listExpiredAppsMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeCleanup) recordInvocation(key string, args []interface{}) {
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

var _ cleanup.Cleanup = new(FakeCleanup)
