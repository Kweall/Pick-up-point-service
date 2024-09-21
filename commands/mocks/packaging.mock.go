// Code generated by http://github.com/gojuno/minimock (v3.4.0). DO NOT EDIT.

package mocks

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// PackagingMock implements mm_commands.Packaging
type PackagingMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcCheckWeight          func(weight float64) (b1 bool)
	funcCheckWeightOrigin    string
	inspectFuncCheckWeight   func(weight float64)
	afterCheckWeightCounter  uint64
	beforeCheckWeightCounter uint64
	CheckWeightMock          mPackagingMockCheckWeight

	funcGetPrice          func() (i1 int64)
	funcGetPriceOrigin    string
	inspectFuncGetPrice   func()
	afterGetPriceCounter  uint64
	beforeGetPriceCounter uint64
	GetPriceMock          mPackagingMockGetPrice
}

// NewPackagingMock returns a mock for mm_commands.Packaging
func NewPackagingMock(t minimock.Tester) *PackagingMock {
	m := &PackagingMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CheckWeightMock = mPackagingMockCheckWeight{mock: m}
	m.CheckWeightMock.callArgs = []*PackagingMockCheckWeightParams{}

	m.GetPriceMock = mPackagingMockGetPrice{mock: m}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mPackagingMockCheckWeight struct {
	optional           bool
	mock               *PackagingMock
	defaultExpectation *PackagingMockCheckWeightExpectation
	expectations       []*PackagingMockCheckWeightExpectation

	callArgs []*PackagingMockCheckWeightParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// PackagingMockCheckWeightExpectation specifies expectation struct of the Packaging.CheckWeight
type PackagingMockCheckWeightExpectation struct {
	mock               *PackagingMock
	params             *PackagingMockCheckWeightParams
	paramPtrs          *PackagingMockCheckWeightParamPtrs
	expectationOrigins PackagingMockCheckWeightExpectationOrigins
	results            *PackagingMockCheckWeightResults
	returnOrigin       string
	Counter            uint64
}

// PackagingMockCheckWeightParams contains parameters of the Packaging.CheckWeight
type PackagingMockCheckWeightParams struct {
	weight float64
}

// PackagingMockCheckWeightParamPtrs contains pointers to parameters of the Packaging.CheckWeight
type PackagingMockCheckWeightParamPtrs struct {
	weight *float64
}

// PackagingMockCheckWeightResults contains results of the Packaging.CheckWeight
type PackagingMockCheckWeightResults struct {
	b1 bool
}

// PackagingMockCheckWeightOrigins contains origins of expectations of the Packaging.CheckWeight
type PackagingMockCheckWeightExpectationOrigins struct {
	origin       string
	originWeight string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmCheckWeight *mPackagingMockCheckWeight) Optional() *mPackagingMockCheckWeight {
	mmCheckWeight.optional = true
	return mmCheckWeight
}

// Expect sets up expected params for Packaging.CheckWeight
func (mmCheckWeight *mPackagingMockCheckWeight) Expect(weight float64) *mPackagingMockCheckWeight {
	if mmCheckWeight.mock.funcCheckWeight != nil {
		mmCheckWeight.mock.t.Fatalf("PackagingMock.CheckWeight mock is already set by Set")
	}

	if mmCheckWeight.defaultExpectation == nil {
		mmCheckWeight.defaultExpectation = &PackagingMockCheckWeightExpectation{}
	}

	if mmCheckWeight.defaultExpectation.paramPtrs != nil {
		mmCheckWeight.mock.t.Fatalf("PackagingMock.CheckWeight mock is already set by ExpectParams functions")
	}

	mmCheckWeight.defaultExpectation.params = &PackagingMockCheckWeightParams{weight}
	mmCheckWeight.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmCheckWeight.expectations {
		if minimock.Equal(e.params, mmCheckWeight.defaultExpectation.params) {
			mmCheckWeight.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmCheckWeight.defaultExpectation.params)
		}
	}

	return mmCheckWeight
}

// ExpectWeightParam1 sets up expected param weight for Packaging.CheckWeight
func (mmCheckWeight *mPackagingMockCheckWeight) ExpectWeightParam1(weight float64) *mPackagingMockCheckWeight {
	if mmCheckWeight.mock.funcCheckWeight != nil {
		mmCheckWeight.mock.t.Fatalf("PackagingMock.CheckWeight mock is already set by Set")
	}

	if mmCheckWeight.defaultExpectation == nil {
		mmCheckWeight.defaultExpectation = &PackagingMockCheckWeightExpectation{}
	}

	if mmCheckWeight.defaultExpectation.params != nil {
		mmCheckWeight.mock.t.Fatalf("PackagingMock.CheckWeight mock is already set by Expect")
	}

	if mmCheckWeight.defaultExpectation.paramPtrs == nil {
		mmCheckWeight.defaultExpectation.paramPtrs = &PackagingMockCheckWeightParamPtrs{}
	}
	mmCheckWeight.defaultExpectation.paramPtrs.weight = &weight
	mmCheckWeight.defaultExpectation.expectationOrigins.originWeight = minimock.CallerInfo(1)

	return mmCheckWeight
}

// Inspect accepts an inspector function that has same arguments as the Packaging.CheckWeight
func (mmCheckWeight *mPackagingMockCheckWeight) Inspect(f func(weight float64)) *mPackagingMockCheckWeight {
	if mmCheckWeight.mock.inspectFuncCheckWeight != nil {
		mmCheckWeight.mock.t.Fatalf("Inspect function is already set for PackagingMock.CheckWeight")
	}

	mmCheckWeight.mock.inspectFuncCheckWeight = f

	return mmCheckWeight
}

// Return sets up results that will be returned by Packaging.CheckWeight
func (mmCheckWeight *mPackagingMockCheckWeight) Return(b1 bool) *PackagingMock {
	if mmCheckWeight.mock.funcCheckWeight != nil {
		mmCheckWeight.mock.t.Fatalf("PackagingMock.CheckWeight mock is already set by Set")
	}

	if mmCheckWeight.defaultExpectation == nil {
		mmCheckWeight.defaultExpectation = &PackagingMockCheckWeightExpectation{mock: mmCheckWeight.mock}
	}
	mmCheckWeight.defaultExpectation.results = &PackagingMockCheckWeightResults{b1}
	mmCheckWeight.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmCheckWeight.mock
}

// Set uses given function f to mock the Packaging.CheckWeight method
func (mmCheckWeight *mPackagingMockCheckWeight) Set(f func(weight float64) (b1 bool)) *PackagingMock {
	if mmCheckWeight.defaultExpectation != nil {
		mmCheckWeight.mock.t.Fatalf("Default expectation is already set for the Packaging.CheckWeight method")
	}

	if len(mmCheckWeight.expectations) > 0 {
		mmCheckWeight.mock.t.Fatalf("Some expectations are already set for the Packaging.CheckWeight method")
	}

	mmCheckWeight.mock.funcCheckWeight = f
	mmCheckWeight.mock.funcCheckWeightOrigin = minimock.CallerInfo(1)
	return mmCheckWeight.mock
}

// When sets expectation for the Packaging.CheckWeight which will trigger the result defined by the following
// Then helper
func (mmCheckWeight *mPackagingMockCheckWeight) When(weight float64) *PackagingMockCheckWeightExpectation {
	if mmCheckWeight.mock.funcCheckWeight != nil {
		mmCheckWeight.mock.t.Fatalf("PackagingMock.CheckWeight mock is already set by Set")
	}

	expectation := &PackagingMockCheckWeightExpectation{
		mock:               mmCheckWeight.mock,
		params:             &PackagingMockCheckWeightParams{weight},
		expectationOrigins: PackagingMockCheckWeightExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmCheckWeight.expectations = append(mmCheckWeight.expectations, expectation)
	return expectation
}

// Then sets up Packaging.CheckWeight return parameters for the expectation previously defined by the When method
func (e *PackagingMockCheckWeightExpectation) Then(b1 bool) *PackagingMock {
	e.results = &PackagingMockCheckWeightResults{b1}
	return e.mock
}

// Times sets number of times Packaging.CheckWeight should be invoked
func (mmCheckWeight *mPackagingMockCheckWeight) Times(n uint64) *mPackagingMockCheckWeight {
	if n == 0 {
		mmCheckWeight.mock.t.Fatalf("Times of PackagingMock.CheckWeight mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmCheckWeight.expectedInvocations, n)
	mmCheckWeight.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmCheckWeight
}

func (mmCheckWeight *mPackagingMockCheckWeight) invocationsDone() bool {
	if len(mmCheckWeight.expectations) == 0 && mmCheckWeight.defaultExpectation == nil && mmCheckWeight.mock.funcCheckWeight == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmCheckWeight.mock.afterCheckWeightCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmCheckWeight.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// CheckWeight implements mm_commands.Packaging
func (mmCheckWeight *PackagingMock) CheckWeight(weight float64) (b1 bool) {
	mm_atomic.AddUint64(&mmCheckWeight.beforeCheckWeightCounter, 1)
	defer mm_atomic.AddUint64(&mmCheckWeight.afterCheckWeightCounter, 1)

	mmCheckWeight.t.Helper()

	if mmCheckWeight.inspectFuncCheckWeight != nil {
		mmCheckWeight.inspectFuncCheckWeight(weight)
	}

	mm_params := PackagingMockCheckWeightParams{weight}

	// Record call args
	mmCheckWeight.CheckWeightMock.mutex.Lock()
	mmCheckWeight.CheckWeightMock.callArgs = append(mmCheckWeight.CheckWeightMock.callArgs, &mm_params)
	mmCheckWeight.CheckWeightMock.mutex.Unlock()

	for _, e := range mmCheckWeight.CheckWeightMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.b1
		}
	}

	if mmCheckWeight.CheckWeightMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmCheckWeight.CheckWeightMock.defaultExpectation.Counter, 1)
		mm_want := mmCheckWeight.CheckWeightMock.defaultExpectation.params
		mm_want_ptrs := mmCheckWeight.CheckWeightMock.defaultExpectation.paramPtrs

		mm_got := PackagingMockCheckWeightParams{weight}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.weight != nil && !minimock.Equal(*mm_want_ptrs.weight, mm_got.weight) {
				mmCheckWeight.t.Errorf("PackagingMock.CheckWeight got unexpected parameter weight, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmCheckWeight.CheckWeightMock.defaultExpectation.expectationOrigins.originWeight, *mm_want_ptrs.weight, mm_got.weight, minimock.Diff(*mm_want_ptrs.weight, mm_got.weight))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmCheckWeight.t.Errorf("PackagingMock.CheckWeight got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmCheckWeight.CheckWeightMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmCheckWeight.CheckWeightMock.defaultExpectation.results
		if mm_results == nil {
			mmCheckWeight.t.Fatal("No results are set for the PackagingMock.CheckWeight")
		}
		return (*mm_results).b1
	}
	if mmCheckWeight.funcCheckWeight != nil {
		return mmCheckWeight.funcCheckWeight(weight)
	}
	mmCheckWeight.t.Fatalf("Unexpected call to PackagingMock.CheckWeight. %v", weight)
	return
}

// CheckWeightAfterCounter returns a count of finished PackagingMock.CheckWeight invocations
func (mmCheckWeight *PackagingMock) CheckWeightAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCheckWeight.afterCheckWeightCounter)
}

// CheckWeightBeforeCounter returns a count of PackagingMock.CheckWeight invocations
func (mmCheckWeight *PackagingMock) CheckWeightBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCheckWeight.beforeCheckWeightCounter)
}

// Calls returns a list of arguments used in each call to PackagingMock.CheckWeight.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmCheckWeight *mPackagingMockCheckWeight) Calls() []*PackagingMockCheckWeightParams {
	mmCheckWeight.mutex.RLock()

	argCopy := make([]*PackagingMockCheckWeightParams, len(mmCheckWeight.callArgs))
	copy(argCopy, mmCheckWeight.callArgs)

	mmCheckWeight.mutex.RUnlock()

	return argCopy
}

// MinimockCheckWeightDone returns true if the count of the CheckWeight invocations corresponds
// the number of defined expectations
func (m *PackagingMock) MinimockCheckWeightDone() bool {
	if m.CheckWeightMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.CheckWeightMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.CheckWeightMock.invocationsDone()
}

// MinimockCheckWeightInspect logs each unmet expectation
func (m *PackagingMock) MinimockCheckWeightInspect() {
	for _, e := range m.CheckWeightMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to PackagingMock.CheckWeight at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterCheckWeightCounter := mm_atomic.LoadUint64(&m.afterCheckWeightCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.CheckWeightMock.defaultExpectation != nil && afterCheckWeightCounter < 1 {
		if m.CheckWeightMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to PackagingMock.CheckWeight at\n%s", m.CheckWeightMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to PackagingMock.CheckWeight at\n%s with params: %#v", m.CheckWeightMock.defaultExpectation.expectationOrigins.origin, *m.CheckWeightMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCheckWeight != nil && afterCheckWeightCounter < 1 {
		m.t.Errorf("Expected call to PackagingMock.CheckWeight at\n%s", m.funcCheckWeightOrigin)
	}

	if !m.CheckWeightMock.invocationsDone() && afterCheckWeightCounter > 0 {
		m.t.Errorf("Expected %d calls to PackagingMock.CheckWeight at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.CheckWeightMock.expectedInvocations), m.CheckWeightMock.expectedInvocationsOrigin, afterCheckWeightCounter)
	}
}

type mPackagingMockGetPrice struct {
	optional           bool
	mock               *PackagingMock
	defaultExpectation *PackagingMockGetPriceExpectation
	expectations       []*PackagingMockGetPriceExpectation

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// PackagingMockGetPriceExpectation specifies expectation struct of the Packaging.GetPrice
type PackagingMockGetPriceExpectation struct {
	mock *PackagingMock

	results      *PackagingMockGetPriceResults
	returnOrigin string
	Counter      uint64
}

// PackagingMockGetPriceResults contains results of the Packaging.GetPrice
type PackagingMockGetPriceResults struct {
	i1 int64
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmGetPrice *mPackagingMockGetPrice) Optional() *mPackagingMockGetPrice {
	mmGetPrice.optional = true
	return mmGetPrice
}

// Expect sets up expected params for Packaging.GetPrice
func (mmGetPrice *mPackagingMockGetPrice) Expect() *mPackagingMockGetPrice {
	if mmGetPrice.mock.funcGetPrice != nil {
		mmGetPrice.mock.t.Fatalf("PackagingMock.GetPrice mock is already set by Set")
	}

	if mmGetPrice.defaultExpectation == nil {
		mmGetPrice.defaultExpectation = &PackagingMockGetPriceExpectation{}
	}

	return mmGetPrice
}

// Inspect accepts an inspector function that has same arguments as the Packaging.GetPrice
func (mmGetPrice *mPackagingMockGetPrice) Inspect(f func()) *mPackagingMockGetPrice {
	if mmGetPrice.mock.inspectFuncGetPrice != nil {
		mmGetPrice.mock.t.Fatalf("Inspect function is already set for PackagingMock.GetPrice")
	}

	mmGetPrice.mock.inspectFuncGetPrice = f

	return mmGetPrice
}

// Return sets up results that will be returned by Packaging.GetPrice
func (mmGetPrice *mPackagingMockGetPrice) Return(i1 int64) *PackagingMock {
	if mmGetPrice.mock.funcGetPrice != nil {
		mmGetPrice.mock.t.Fatalf("PackagingMock.GetPrice mock is already set by Set")
	}

	if mmGetPrice.defaultExpectation == nil {
		mmGetPrice.defaultExpectation = &PackagingMockGetPriceExpectation{mock: mmGetPrice.mock}
	}
	mmGetPrice.defaultExpectation.results = &PackagingMockGetPriceResults{i1}
	mmGetPrice.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmGetPrice.mock
}

// Set uses given function f to mock the Packaging.GetPrice method
func (mmGetPrice *mPackagingMockGetPrice) Set(f func() (i1 int64)) *PackagingMock {
	if mmGetPrice.defaultExpectation != nil {
		mmGetPrice.mock.t.Fatalf("Default expectation is already set for the Packaging.GetPrice method")
	}

	if len(mmGetPrice.expectations) > 0 {
		mmGetPrice.mock.t.Fatalf("Some expectations are already set for the Packaging.GetPrice method")
	}

	mmGetPrice.mock.funcGetPrice = f
	mmGetPrice.mock.funcGetPriceOrigin = minimock.CallerInfo(1)
	return mmGetPrice.mock
}

// Times sets number of times Packaging.GetPrice should be invoked
func (mmGetPrice *mPackagingMockGetPrice) Times(n uint64) *mPackagingMockGetPrice {
	if n == 0 {
		mmGetPrice.mock.t.Fatalf("Times of PackagingMock.GetPrice mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmGetPrice.expectedInvocations, n)
	mmGetPrice.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmGetPrice
}

func (mmGetPrice *mPackagingMockGetPrice) invocationsDone() bool {
	if len(mmGetPrice.expectations) == 0 && mmGetPrice.defaultExpectation == nil && mmGetPrice.mock.funcGetPrice == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmGetPrice.mock.afterGetPriceCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmGetPrice.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// GetPrice implements mm_commands.Packaging
func (mmGetPrice *PackagingMock) GetPrice() (i1 int64) {
	mm_atomic.AddUint64(&mmGetPrice.beforeGetPriceCounter, 1)
	defer mm_atomic.AddUint64(&mmGetPrice.afterGetPriceCounter, 1)

	mmGetPrice.t.Helper()

	if mmGetPrice.inspectFuncGetPrice != nil {
		mmGetPrice.inspectFuncGetPrice()
	}

	if mmGetPrice.GetPriceMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetPrice.GetPriceMock.defaultExpectation.Counter, 1)

		mm_results := mmGetPrice.GetPriceMock.defaultExpectation.results
		if mm_results == nil {
			mmGetPrice.t.Fatal("No results are set for the PackagingMock.GetPrice")
		}
		return (*mm_results).i1
	}
	if mmGetPrice.funcGetPrice != nil {
		return mmGetPrice.funcGetPrice()
	}
	mmGetPrice.t.Fatalf("Unexpected call to PackagingMock.GetPrice.")
	return
}

// GetPriceAfterCounter returns a count of finished PackagingMock.GetPrice invocations
func (mmGetPrice *PackagingMock) GetPriceAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetPrice.afterGetPriceCounter)
}

// GetPriceBeforeCounter returns a count of PackagingMock.GetPrice invocations
func (mmGetPrice *PackagingMock) GetPriceBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetPrice.beforeGetPriceCounter)
}

// MinimockGetPriceDone returns true if the count of the GetPrice invocations corresponds
// the number of defined expectations
func (m *PackagingMock) MinimockGetPriceDone() bool {
	if m.GetPriceMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.GetPriceMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.GetPriceMock.invocationsDone()
}

// MinimockGetPriceInspect logs each unmet expectation
func (m *PackagingMock) MinimockGetPriceInspect() {
	for _, e := range m.GetPriceMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to PackagingMock.GetPrice")
		}
	}

	afterGetPriceCounter := mm_atomic.LoadUint64(&m.afterGetPriceCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.GetPriceMock.defaultExpectation != nil && afterGetPriceCounter < 1 {
		m.t.Errorf("Expected call to PackagingMock.GetPrice at\n%s", m.GetPriceMock.defaultExpectation.returnOrigin)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetPrice != nil && afterGetPriceCounter < 1 {
		m.t.Errorf("Expected call to PackagingMock.GetPrice at\n%s", m.funcGetPriceOrigin)
	}

	if !m.GetPriceMock.invocationsDone() && afterGetPriceCounter > 0 {
		m.t.Errorf("Expected %d calls to PackagingMock.GetPrice at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.GetPriceMock.expectedInvocations), m.GetPriceMock.expectedInvocationsOrigin, afterGetPriceCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *PackagingMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockCheckWeightInspect()

			m.MinimockGetPriceInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *PackagingMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *PackagingMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockCheckWeightDone() &&
		m.MinimockGetPriceDone()
}
