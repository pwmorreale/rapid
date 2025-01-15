// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"io"
	"net/http"
	"sync"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/service"
)

type FakeService struct {
	CheckContainsStub        func(string, *http.Response, *config.Request) error
	checkContainsMutex       sync.RWMutex
	checkContainsArgsForCall []struct {
		arg1 string
		arg2 *http.Response
		arg3 *config.Request
	}
	checkContainsReturns struct {
		result1 error
	}
	checkContainsReturnsOnCall map[int]struct {
		result1 error
	}
	CreateClientStub        func(*config.Request) (*http.Client, error)
	createClientMutex       sync.RWMutex
	createClientArgsForCall []struct {
		arg1 *config.Request
	}
	createClientReturns struct {
		result1 *http.Client
		result2 error
	}
	createClientReturnsOnCall map[int]struct {
		result1 *http.Client
		result2 error
	}
	CreateRequestStub        func(*config.Request) (*http.Request, error)
	createRequestMutex       sync.RWMutex
	createRequestArgsForCall []struct {
		arg1 *config.Request
	}
	createRequestReturns struct {
		result1 *http.Request
		result2 error
	}
	createRequestReturnsOnCall map[int]struct {
		result1 *http.Request
		result2 error
	}
	ExtractStub        func(string, *http.Response, *config.Request) error
	extractMutex       sync.RWMutex
	extractArgsForCall []struct {
		arg1 string
		arg2 *http.Response
		arg3 *config.Request
	}
	extractReturns struct {
		result1 error
	}
	extractReturnsOnCall map[int]struct {
		result1 error
	}
	GetContentReaderStub        func(*config.Request) io.Reader
	getContentReaderMutex       sync.RWMutex
	getContentReaderArgsForCall []struct {
		arg1 *config.Request
	}
	getContentReaderReturns struct {
		result1 io.Reader
	}
	getContentReaderReturnsOnCall map[int]struct {
		result1 io.Reader
	}
	SendStub        func(*http.Client, *http.Request, *config.Request) (*http.Response, error)
	sendMutex       sync.RWMutex
	sendArgsForCall []struct {
		arg1 *http.Client
		arg2 *http.Request
		arg3 *config.Request
	}
	sendReturns struct {
		result1 *http.Response
		result2 error
	}
	sendReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	ValidateResponseStub        func(*http.Client, *http.Response, *config.Request) error
	validateResponseMutex       sync.RWMutex
	validateResponseArgsForCall []struct {
		arg1 *http.Client
		arg2 *http.Response
		arg3 *config.Request
	}
	validateResponseReturns struct {
		result1 error
	}
	validateResponseReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeService) CheckContains(arg1 string, arg2 *http.Response, arg3 *config.Request) error {
	fake.checkContainsMutex.Lock()
	ret, specificReturn := fake.checkContainsReturnsOnCall[len(fake.checkContainsArgsForCall)]
	fake.checkContainsArgsForCall = append(fake.checkContainsArgsForCall, struct {
		arg1 string
		arg2 *http.Response
		arg3 *config.Request
	}{arg1, arg2, arg3})
	stub := fake.CheckContainsStub
	fakeReturns := fake.checkContainsReturns
	fake.recordInvocation("CheckContains", []interface{}{arg1, arg2, arg3})
	fake.checkContainsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeService) CheckContainsCallCount() int {
	fake.checkContainsMutex.RLock()
	defer fake.checkContainsMutex.RUnlock()
	return len(fake.checkContainsArgsForCall)
}

func (fake *FakeService) CheckContainsCalls(stub func(string, *http.Response, *config.Request) error) {
	fake.checkContainsMutex.Lock()
	defer fake.checkContainsMutex.Unlock()
	fake.CheckContainsStub = stub
}

func (fake *FakeService) CheckContainsArgsForCall(i int) (string, *http.Response, *config.Request) {
	fake.checkContainsMutex.RLock()
	defer fake.checkContainsMutex.RUnlock()
	argsForCall := fake.checkContainsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeService) CheckContainsReturns(result1 error) {
	fake.checkContainsMutex.Lock()
	defer fake.checkContainsMutex.Unlock()
	fake.CheckContainsStub = nil
	fake.checkContainsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeService) CheckContainsReturnsOnCall(i int, result1 error) {
	fake.checkContainsMutex.Lock()
	defer fake.checkContainsMutex.Unlock()
	fake.CheckContainsStub = nil
	if fake.checkContainsReturnsOnCall == nil {
		fake.checkContainsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkContainsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeService) CreateClient(arg1 *config.Request) (*http.Client, error) {
	fake.createClientMutex.Lock()
	ret, specificReturn := fake.createClientReturnsOnCall[len(fake.createClientArgsForCall)]
	fake.createClientArgsForCall = append(fake.createClientArgsForCall, struct {
		arg1 *config.Request
	}{arg1})
	stub := fake.CreateClientStub
	fakeReturns := fake.createClientReturns
	fake.recordInvocation("CreateClient", []interface{}{arg1})
	fake.createClientMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeService) CreateClientCallCount() int {
	fake.createClientMutex.RLock()
	defer fake.createClientMutex.RUnlock()
	return len(fake.createClientArgsForCall)
}

func (fake *FakeService) CreateClientCalls(stub func(*config.Request) (*http.Client, error)) {
	fake.createClientMutex.Lock()
	defer fake.createClientMutex.Unlock()
	fake.CreateClientStub = stub
}

func (fake *FakeService) CreateClientArgsForCall(i int) *config.Request {
	fake.createClientMutex.RLock()
	defer fake.createClientMutex.RUnlock()
	argsForCall := fake.createClientArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeService) CreateClientReturns(result1 *http.Client, result2 error) {
	fake.createClientMutex.Lock()
	defer fake.createClientMutex.Unlock()
	fake.CreateClientStub = nil
	fake.createClientReturns = struct {
		result1 *http.Client
		result2 error
	}{result1, result2}
}

func (fake *FakeService) CreateClientReturnsOnCall(i int, result1 *http.Client, result2 error) {
	fake.createClientMutex.Lock()
	defer fake.createClientMutex.Unlock()
	fake.CreateClientStub = nil
	if fake.createClientReturnsOnCall == nil {
		fake.createClientReturnsOnCall = make(map[int]struct {
			result1 *http.Client
			result2 error
		})
	}
	fake.createClientReturnsOnCall[i] = struct {
		result1 *http.Client
		result2 error
	}{result1, result2}
}

func (fake *FakeService) CreateRequest(arg1 *config.Request) (*http.Request, error) {
	fake.createRequestMutex.Lock()
	ret, specificReturn := fake.createRequestReturnsOnCall[len(fake.createRequestArgsForCall)]
	fake.createRequestArgsForCall = append(fake.createRequestArgsForCall, struct {
		arg1 *config.Request
	}{arg1})
	stub := fake.CreateRequestStub
	fakeReturns := fake.createRequestReturns
	fake.recordInvocation("CreateRequest", []interface{}{arg1})
	fake.createRequestMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeService) CreateRequestCallCount() int {
	fake.createRequestMutex.RLock()
	defer fake.createRequestMutex.RUnlock()
	return len(fake.createRequestArgsForCall)
}

func (fake *FakeService) CreateRequestCalls(stub func(*config.Request) (*http.Request, error)) {
	fake.createRequestMutex.Lock()
	defer fake.createRequestMutex.Unlock()
	fake.CreateRequestStub = stub
}

func (fake *FakeService) CreateRequestArgsForCall(i int) *config.Request {
	fake.createRequestMutex.RLock()
	defer fake.createRequestMutex.RUnlock()
	argsForCall := fake.createRequestArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeService) CreateRequestReturns(result1 *http.Request, result2 error) {
	fake.createRequestMutex.Lock()
	defer fake.createRequestMutex.Unlock()
	fake.CreateRequestStub = nil
	fake.createRequestReturns = struct {
		result1 *http.Request
		result2 error
	}{result1, result2}
}

func (fake *FakeService) CreateRequestReturnsOnCall(i int, result1 *http.Request, result2 error) {
	fake.createRequestMutex.Lock()
	defer fake.createRequestMutex.Unlock()
	fake.CreateRequestStub = nil
	if fake.createRequestReturnsOnCall == nil {
		fake.createRequestReturnsOnCall = make(map[int]struct {
			result1 *http.Request
			result2 error
		})
	}
	fake.createRequestReturnsOnCall[i] = struct {
		result1 *http.Request
		result2 error
	}{result1, result2}
}

func (fake *FakeService) Extract(arg1 string, arg2 *http.Response, arg3 *config.Request) error {
	fake.extractMutex.Lock()
	ret, specificReturn := fake.extractReturnsOnCall[len(fake.extractArgsForCall)]
	fake.extractArgsForCall = append(fake.extractArgsForCall, struct {
		arg1 string
		arg2 *http.Response
		arg3 *config.Request
	}{arg1, arg2, arg3})
	stub := fake.ExtractStub
	fakeReturns := fake.extractReturns
	fake.recordInvocation("Extract", []interface{}{arg1, arg2, arg3})
	fake.extractMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeService) ExtractCallCount() int {
	fake.extractMutex.RLock()
	defer fake.extractMutex.RUnlock()
	return len(fake.extractArgsForCall)
}

func (fake *FakeService) ExtractCalls(stub func(string, *http.Response, *config.Request) error) {
	fake.extractMutex.Lock()
	defer fake.extractMutex.Unlock()
	fake.ExtractStub = stub
}

func (fake *FakeService) ExtractArgsForCall(i int) (string, *http.Response, *config.Request) {
	fake.extractMutex.RLock()
	defer fake.extractMutex.RUnlock()
	argsForCall := fake.extractArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeService) ExtractReturns(result1 error) {
	fake.extractMutex.Lock()
	defer fake.extractMutex.Unlock()
	fake.ExtractStub = nil
	fake.extractReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeService) ExtractReturnsOnCall(i int, result1 error) {
	fake.extractMutex.Lock()
	defer fake.extractMutex.Unlock()
	fake.ExtractStub = nil
	if fake.extractReturnsOnCall == nil {
		fake.extractReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.extractReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeService) GetContentReader(arg1 *config.Request) io.Reader {
	fake.getContentReaderMutex.Lock()
	ret, specificReturn := fake.getContentReaderReturnsOnCall[len(fake.getContentReaderArgsForCall)]
	fake.getContentReaderArgsForCall = append(fake.getContentReaderArgsForCall, struct {
		arg1 *config.Request
	}{arg1})
	stub := fake.GetContentReaderStub
	fakeReturns := fake.getContentReaderReturns
	fake.recordInvocation("GetContentReader", []interface{}{arg1})
	fake.getContentReaderMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeService) GetContentReaderCallCount() int {
	fake.getContentReaderMutex.RLock()
	defer fake.getContentReaderMutex.RUnlock()
	return len(fake.getContentReaderArgsForCall)
}

func (fake *FakeService) GetContentReaderCalls(stub func(*config.Request) io.Reader) {
	fake.getContentReaderMutex.Lock()
	defer fake.getContentReaderMutex.Unlock()
	fake.GetContentReaderStub = stub
}

func (fake *FakeService) GetContentReaderArgsForCall(i int) *config.Request {
	fake.getContentReaderMutex.RLock()
	defer fake.getContentReaderMutex.RUnlock()
	argsForCall := fake.getContentReaderArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeService) GetContentReaderReturns(result1 io.Reader) {
	fake.getContentReaderMutex.Lock()
	defer fake.getContentReaderMutex.Unlock()
	fake.GetContentReaderStub = nil
	fake.getContentReaderReturns = struct {
		result1 io.Reader
	}{result1}
}

func (fake *FakeService) GetContentReaderReturnsOnCall(i int, result1 io.Reader) {
	fake.getContentReaderMutex.Lock()
	defer fake.getContentReaderMutex.Unlock()
	fake.GetContentReaderStub = nil
	if fake.getContentReaderReturnsOnCall == nil {
		fake.getContentReaderReturnsOnCall = make(map[int]struct {
			result1 io.Reader
		})
	}
	fake.getContentReaderReturnsOnCall[i] = struct {
		result1 io.Reader
	}{result1}
}

func (fake *FakeService) Send(arg1 *http.Client, arg2 *http.Request, arg3 *config.Request) (*http.Response, error) {
	fake.sendMutex.Lock()
	ret, specificReturn := fake.sendReturnsOnCall[len(fake.sendArgsForCall)]
	fake.sendArgsForCall = append(fake.sendArgsForCall, struct {
		arg1 *http.Client
		arg2 *http.Request
		arg3 *config.Request
	}{arg1, arg2, arg3})
	stub := fake.SendStub
	fakeReturns := fake.sendReturns
	fake.recordInvocation("Send", []interface{}{arg1, arg2, arg3})
	fake.sendMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeService) SendCallCount() int {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return len(fake.sendArgsForCall)
}

func (fake *FakeService) SendCalls(stub func(*http.Client, *http.Request, *config.Request) (*http.Response, error)) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = stub
}

func (fake *FakeService) SendArgsForCall(i int) (*http.Client, *http.Request, *config.Request) {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	argsForCall := fake.sendArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeService) SendReturns(result1 *http.Response, result2 error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = nil
	fake.sendReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeService) SendReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = nil
	if fake.sendReturnsOnCall == nil {
		fake.sendReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.sendReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeService) ValidateResponse(arg1 *http.Client, arg2 *http.Response, arg3 *config.Request) error {
	fake.validateResponseMutex.Lock()
	ret, specificReturn := fake.validateResponseReturnsOnCall[len(fake.validateResponseArgsForCall)]
	fake.validateResponseArgsForCall = append(fake.validateResponseArgsForCall, struct {
		arg1 *http.Client
		arg2 *http.Response
		arg3 *config.Request
	}{arg1, arg2, arg3})
	stub := fake.ValidateResponseStub
	fakeReturns := fake.validateResponseReturns
	fake.recordInvocation("ValidateResponse", []interface{}{arg1, arg2, arg3})
	fake.validateResponseMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeService) ValidateResponseCallCount() int {
	fake.validateResponseMutex.RLock()
	defer fake.validateResponseMutex.RUnlock()
	return len(fake.validateResponseArgsForCall)
}

func (fake *FakeService) ValidateResponseCalls(stub func(*http.Client, *http.Response, *config.Request) error) {
	fake.validateResponseMutex.Lock()
	defer fake.validateResponseMutex.Unlock()
	fake.ValidateResponseStub = stub
}

func (fake *FakeService) ValidateResponseArgsForCall(i int) (*http.Client, *http.Response, *config.Request) {
	fake.validateResponseMutex.RLock()
	defer fake.validateResponseMutex.RUnlock()
	argsForCall := fake.validateResponseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeService) ValidateResponseReturns(result1 error) {
	fake.validateResponseMutex.Lock()
	defer fake.validateResponseMutex.Unlock()
	fake.ValidateResponseStub = nil
	fake.validateResponseReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeService) ValidateResponseReturnsOnCall(i int, result1 error) {
	fake.validateResponseMutex.Lock()
	defer fake.validateResponseMutex.Unlock()
	fake.ValidateResponseStub = nil
	if fake.validateResponseReturnsOnCall == nil {
		fake.validateResponseReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateResponseReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkContainsMutex.RLock()
	defer fake.checkContainsMutex.RUnlock()
	fake.createClientMutex.RLock()
	defer fake.createClientMutex.RUnlock()
	fake.createRequestMutex.RLock()
	defer fake.createRequestMutex.RUnlock()
	fake.extractMutex.RLock()
	defer fake.extractMutex.RUnlock()
	fake.getContentReaderMutex.RLock()
	defer fake.getContentReaderMutex.RUnlock()
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	fake.validateResponseMutex.RLock()
	defer fake.validateResponseMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeService) recordInvocation(key string, args []interface{}) {
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

var _ service.Service = new(FakeService)
