// Code generated by counterfeiter. DO NOT EDIT.
package sharedfakes

import (
	"sync"

	"code.cloudfoundry.org/cli/actor/v3action"
	"code.cloudfoundry.org/cli/command/v6/shared"
)

type FakeV3AppSummaryActor struct {
	GetApplicationSummaryByNameAndSpaceStub        func(string, string, bool) (v3action.ApplicationSummary, v3action.Warnings, error)
	getApplicationSummaryByNameAndSpaceMutex       sync.RWMutex
	getApplicationSummaryByNameAndSpaceArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 bool
	}
	getApplicationSummaryByNameAndSpaceReturns struct {
		result1 v3action.ApplicationSummary
		result2 v3action.Warnings
		result3 error
	}
	getApplicationSummaryByNameAndSpaceReturnsOnCall map[int]struct {
		result1 v3action.ApplicationSummary
		result2 v3action.Warnings
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeV3AppSummaryActor) GetApplicationSummaryByNameAndSpace(arg1 string, arg2 string, arg3 bool) (v3action.ApplicationSummary, v3action.Warnings, error) {
	fake.getApplicationSummaryByNameAndSpaceMutex.Lock()
	ret, specificReturn := fake.getApplicationSummaryByNameAndSpaceReturnsOnCall[len(fake.getApplicationSummaryByNameAndSpaceArgsForCall)]
	fake.getApplicationSummaryByNameAndSpaceArgsForCall = append(fake.getApplicationSummaryByNameAndSpaceArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 bool
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetApplicationSummaryByNameAndSpace", []interface{}{arg1, arg2, arg3})
	fake.getApplicationSummaryByNameAndSpaceMutex.Unlock()
	if fake.GetApplicationSummaryByNameAndSpaceStub != nil {
		return fake.GetApplicationSummaryByNameAndSpaceStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.getApplicationSummaryByNameAndSpaceReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeV3AppSummaryActor) GetApplicationSummaryByNameAndSpaceCallCount() int {
	fake.getApplicationSummaryByNameAndSpaceMutex.RLock()
	defer fake.getApplicationSummaryByNameAndSpaceMutex.RUnlock()
	return len(fake.getApplicationSummaryByNameAndSpaceArgsForCall)
}

func (fake *FakeV3AppSummaryActor) GetApplicationSummaryByNameAndSpaceCalls(stub func(string, string, bool) (v3action.ApplicationSummary, v3action.Warnings, error)) {
	fake.getApplicationSummaryByNameAndSpaceMutex.Lock()
	defer fake.getApplicationSummaryByNameAndSpaceMutex.Unlock()
	fake.GetApplicationSummaryByNameAndSpaceStub = stub
}

func (fake *FakeV3AppSummaryActor) GetApplicationSummaryByNameAndSpaceArgsForCall(i int) (string, string, bool) {
	fake.getApplicationSummaryByNameAndSpaceMutex.RLock()
	defer fake.getApplicationSummaryByNameAndSpaceMutex.RUnlock()
	argsForCall := fake.getApplicationSummaryByNameAndSpaceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeV3AppSummaryActor) GetApplicationSummaryByNameAndSpaceReturns(result1 v3action.ApplicationSummary, result2 v3action.Warnings, result3 error) {
	fake.getApplicationSummaryByNameAndSpaceMutex.Lock()
	defer fake.getApplicationSummaryByNameAndSpaceMutex.Unlock()
	fake.GetApplicationSummaryByNameAndSpaceStub = nil
	fake.getApplicationSummaryByNameAndSpaceReturns = struct {
		result1 v3action.ApplicationSummary
		result2 v3action.Warnings
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeV3AppSummaryActor) GetApplicationSummaryByNameAndSpaceReturnsOnCall(i int, result1 v3action.ApplicationSummary, result2 v3action.Warnings, result3 error) {
	fake.getApplicationSummaryByNameAndSpaceMutex.Lock()
	defer fake.getApplicationSummaryByNameAndSpaceMutex.Unlock()
	fake.GetApplicationSummaryByNameAndSpaceStub = nil
	if fake.getApplicationSummaryByNameAndSpaceReturnsOnCall == nil {
		fake.getApplicationSummaryByNameAndSpaceReturnsOnCall = make(map[int]struct {
			result1 v3action.ApplicationSummary
			result2 v3action.Warnings
			result3 error
		})
	}
	fake.getApplicationSummaryByNameAndSpaceReturnsOnCall[i] = struct {
		result1 v3action.ApplicationSummary
		result2 v3action.Warnings
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeV3AppSummaryActor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getApplicationSummaryByNameAndSpaceMutex.RLock()
	defer fake.getApplicationSummaryByNameAndSpaceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeV3AppSummaryActor) recordInvocation(key string, args []interface{}) {
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

var _ shared.V3AppSummaryActor = new(FakeV3AppSummaryActor)