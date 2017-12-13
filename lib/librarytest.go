// +build e2e_test

// This is a file with e2e tests for C bindings written in library.go.
// As a CGO file, it can't have `_test.go` suffix as it's not allowed by Go.
// At the same time, we don't want this file to be included in the binaries.
// This is why `e2e_test` tag was introduced. Without it, this file is excluded
// from the build. Providing this tag will include this file into the build
// and that's what is done while running e2e tests for `lib/` package.

// Additionaly this file should contain test that mock the Status API.
// Existing test in 'utils.go' that test the Status API will be migrated to the
// e2e package and test that test the C Binding will be migrated to this file

package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/golang/mock/gomock"
	"github.com/status-im/status-go/geth/common"
	"github.com/stretchr/testify/assert"
)

func testCreateAccountWithMock(t *testing.T) {
	realStatusAPI := statusAPI
	defer func() { statusAPI = realStatusAPI }()

	// Setup Mock StatusAPI
	ctrl := gomock.NewController(t)
	status := NewMocklibStatusAPI(ctrl)
	statusAPI = status
	accountInfo1 := common.AccountInfo{Address: "add", Mnemonic: "mne", PubKey: "Pub"}
	accountInfo2 := common.AccountInfo{Error: "Error Message"}
	status.EXPECT().CreateAccount("pass1").Return(accountInfo1, nil)
	status.EXPECT().CreateAccount("").Return(accountInfo1, nil)
	status.EXPECT().CreateAccount(C.GoString(nil)).Return(accountInfo1, nil)
	status.EXPECT().CreateAccount("pass2").Return(accountInfo2, fmt.Errorf("Error Message"))

	// C Strings
	pass1 := C.CString("pass1")
	pass2 := C.CString("pass2")
	empty := C.CString("")
	accountInfo1JSON := C.CString(`{"address":"add","pubkey":"Pub","mnemonic":"mne","error":""}`)
	accountInfo2JSON := C.CString(`{"address":"","pubkey":"","mnemonic":"","error":"Error Message"}`)
	defer func() {
		C.free(unsafe.Pointer(pass1))
		C.free(unsafe.Pointer(pass2))
		C.free(unsafe.Pointer(empty))
		C.free(unsafe.Pointer(accountInfo1JSON))
		C.free(unsafe.Pointer(accountInfo2JSON))
	}()

	tests := []struct {
		name     string
		password *C.char
		want     *C.char
	}{
		{"testCreateAccountWithMock/Normal", pass1, accountInfo1JSON},
		{"testCreateAccountWithMock/EmptyParam", empty, accountInfo1JSON},
		{"testCreateAccountWithMock/NilParam", nil, accountInfo1JSON},
		{"testCreateAccountWithMock/ErrorResult", pass2, accountInfo2JSON},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateAccount(tt.password); C.GoString(got) != C.GoString(tt.want) {
				assert.JSONEq(t, C.GoString(tt.want), C.GoString(got))
			}
		})
	}
}

func testCreateChildAccountWithMock(t *testing.T) {
	realStatusAPI := statusAPI
	defer func() { statusAPI = realStatusAPI }()

	// Setup Mock StatusAPI
	ctrl := gomock.NewController(t)
	status := NewMocklibStatusAPI(ctrl)
	statusAPI = status

	accountInfo1 := common.AccountInfo{Address: "add", PubKey: "Pub"}
	accountInfo2 := common.AccountInfo{Error: "Error Message"}
	status.EXPECT().CreateChildAccount("parent1", "pass1").Return(accountInfo1, nil)
	status.EXPECT().CreateChildAccount("", "").Return(accountInfo1, nil).AnyTimes()
	status.EXPECT().CreateChildAccount("parent2", "pass2").Return(accountInfo2, fmt.Errorf("Error Message"))

	// C Strings
	pass1 := C.CString("pass1")
	pass2 := C.CString("pass2")
	parent1 := C.CString("parent1")
	parent2 := C.CString("parent2")
	empty := C.CString("")
	accountInfo1JSON := C.CString(`{"address":"add","pubkey":"Pub","mnemonic":"","error":""}`)
	accountInfo2JSON := C.CString(`{"address":"","pubkey":"","mnemonic":"","error":"Error Message"}`)
	defer func() {
		C.free(unsafe.Pointer(pass1))
		C.free(unsafe.Pointer(pass2))
		C.free(unsafe.Pointer(parent1))
		C.free(unsafe.Pointer(parent2))
		C.free(unsafe.Pointer(empty))
		C.free(unsafe.Pointer(accountInfo1JSON))
		C.free(unsafe.Pointer(accountInfo2JSON))
	}()

	tests := []struct {
		name     string
		parrent  *C.char
		password *C.char
		want     *C.char
	}{
		{"testCreateChildAccountWithMock/Normal", parent1, pass1, accountInfo1JSON},
		{"testCreateChildAccountWithMock/EmptyParam", empty, empty, accountInfo1JSON},
		{"testCreateChildAccountWithMock/NilParam", nil, nil, accountInfo1JSON},
		{"testCreateChildAccountWithMock/ErrorResult", parent2, pass2, accountInfo2JSON},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateChildAccount(tt.parrent, tt.password); C.GoString(got) != C.GoString(tt.want) {
				assert.JSONEq(t, C.GoString(tt.want), C.GoString(got))
			}
		})
	}

}

func testRecoverAccountWithMock(t *testing.T) {
	realStatusAPI := statusAPI
	defer func() { statusAPI = realStatusAPI }()

	// Setup Mock StatusAPI
	ctrl := gomock.NewController(t)
	status := NewMocklibStatusAPI(ctrl)
	statusAPI = status

	accountInfo1 := common.AccountInfo{Address: "add", PubKey: "Pub", Mnemonic: "mnemonic"}
	accountInfo2 := common.AccountInfo{Error: "Error Message"}
	status.EXPECT().RecoverAccount("pass1", "mnemonic1").Return(accountInfo1, nil)
	status.EXPECT().RecoverAccount("", "").Return(accountInfo1, nil).AnyTimes()
	status.EXPECT().RecoverAccount("pass2", "mnemonic2").Return(accountInfo2, fmt.Errorf("Error Message"))

	// C Strings
	pass1 := C.CString("pass1")
	pass2 := C.CString("pass2")
	mnemonic1 := C.CString("mnemonic1")
	mnemonic2 := C.CString("mnemonic2")
	empty := C.CString("")
	accountInfo1JSON := C.CString(`{"address":"add","pubkey":"Pub","mnemonic":"mnemonic","error":""}`)
	accountInfo2JSON := C.CString(`{"address":"","pubkey":"","mnemonic":"","error":"Error Message"}`)
	defer func() {
		C.free(unsafe.Pointer(pass1))
		C.free(unsafe.Pointer(pass2))
		C.free(unsafe.Pointer(mnemonic1))
		C.free(unsafe.Pointer(mnemonic2))
		C.free(unsafe.Pointer(empty))
		C.free(unsafe.Pointer(accountInfo1JSON))
		C.free(unsafe.Pointer(accountInfo2JSON))
	}()

	tests := []struct {
		name     string
		password *C.char
		mnemonic *C.char
		want     *C.char
	}{
		{"testRecoverAccountWithMock/Normal", pass1, mnemonic1, accountInfo1JSON},
		{"testRecoverAccountWithMock/EmptyParam", empty, empty, accountInfo1JSON},
		{"testRecoverAccountWithMock/NilParam", nil, nil, accountInfo1JSON},
		{"testRecoverAccountWithMock/ErrorResult", pass2, mnemonic2, accountInfo2JSON},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RecoverAccount(tt.password, tt.mnemonic); C.GoString(got) != C.GoString(tt.want) {
				assert.JSONEq(t, C.GoString(tt.want), C.GoString(got))
			}
		})
	}

}

func testValidateNodeConfigWithMock(t *testing.T) {
	realStatusAPI := statusAPI
	defer func() { statusAPI = realStatusAPI }()

	// Setup Mock StatusAPI
	ctrl := gomock.NewController(t)
	status := NewMocklibStatusAPI(ctrl)
	statusAPI = status

	apiDetailedResponse1 := common.APIDetailedResponse{Status: true}
	apiDetailedResponse2 := common.APIDetailedResponse{Status: false, FieldErrors: []common.APIFieldError{
		{Parameter: "param1", Errors: []common.APIError{{Message: "perror1"}, {Message: "perror2"}}},
		{Parameter: "param2", Errors: []common.APIError{{Message: "perror1"}}},
	}}
	apiDetailedResponse3 := common.APIDetailedResponse{}

	status.EXPECT().ValidateJSONConfig("{json1}").Return(apiDetailedResponse1)
	status.EXPECT().ValidateJSONConfig("{json2}").Return(apiDetailedResponse2)
	status.EXPECT().ValidateJSONConfig("").Return(apiDetailedResponse3).AnyTimes()

	// C Strings
	config1 := C.CString("{json1}")
	config2 := C.CString("{json2}")
	empty := C.CString("")
	apiDetailedResponse1JSON := C.CString(`{"status":true}`)
	apiDetailedResponse2JSON := C.CString(`{"status":false,"field_errors":[{"parameter":"param1","errors":[{"message":"perror1"},{"message":"perror2"}]},{"parameter":"param2","errors":[{"message":"perror1"}]}]}`)
	apiDetailedResponse3JSON := C.CString(`{"status":false}`)
	defer func() {
		C.free(unsafe.Pointer(config1))
		C.free(unsafe.Pointer(config2))
		C.free(unsafe.Pointer(empty))
		C.free(unsafe.Pointer(apiDetailedResponse1JSON))
		C.free(unsafe.Pointer(apiDetailedResponse2JSON))
		C.free(unsafe.Pointer(apiDetailedResponse3JSON))
	}()

	tests := []struct {
		name       string
		configJSON *C.char
		want       *C.char
	}{
		{"testValidateNodeConfigWithMock/Normal", config1, apiDetailedResponse1JSON},
		{"testValidateNodeConfigWithMock/ValidationErrors", config2, apiDetailedResponse2JSON},
		{"testValidateNodeConfigWithMock/emptyconfig", empty, apiDetailedResponse3JSON},
		{"testValidateNodeConfigWithMock/nilconfig", nil, apiDetailedResponse3JSON},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateNodeConfig(tt.configJSON); C.GoString(got) != C.GoString(tt.want) {
				assert.JSONEq(t, C.GoString(tt.want), C.GoString(got))
			}
		})
	}

}

func testCompleteTransactionWithMock(t *testing.T) {
	realStatusAPI := statusAPI
	defer func() { statusAPI = realStatusAPI }()

	// Setup Mock StatusAPI
	ctrl := gomock.NewController(t)
	status := NewMocklibStatusAPI(ctrl)
	statusAPI = status
	status.EXPECT().CompleteTransaction(common.QueuedTxID("id1"), "pass1").Return(common.CompleteTransactionResult{ID: "id1", Hash: "0x123"}, nil)
	status.EXPECT().CompleteTransaction(common.QueuedTxID(""), "").Return(common.CompleteTransactionResult{ID: "id1", Hash: "0x123"}, nil).AnyTimes()
	status.EXPECT().CompleteTransaction(common.QueuedTxID("id2"), "pass2").Return(common.CompleteTransactionResult{ID: "id2", Error: "Test Error"}, fmt.Errorf("Test Error"))

	// C Strings
	id1 := C.CString("id1")
	pass1 := C.CString("pass1")
	id2 := C.CString("id2")
	pass2 := C.CString("pass2")
	empty := C.CString("")
	completeTransactionResult1JSON := C.CString(`{"id":"id1","0x123":"hash1","error":""}`)
	completeTransactionResult2JSON := C.CString(`{"id":"id2","0x123":"","error":"Test Error"}`)
	defer func() {
		C.free(unsafe.Pointer(id1))
		C.free(unsafe.Pointer(pass1))
		C.free(unsafe.Pointer(id2))
		C.free(unsafe.Pointer(pass2))
		C.free(unsafe.Pointer(empty))
		C.free(unsafe.Pointer(completeTransactionResult1JSON))
		C.free(unsafe.Pointer(completeTransactionResult2JSON))

	}()

	tests := []struct {
		name     string
		id       *C.char
		password *C.char
		want     *C.char
	}{
		{"testCompleteTransactionWithMock/Normal", id1, pass1, completeTransactionResult1JSON},
		{"testCompleteTransactionWithMock/Empty", empty, empty, completeTransactionResult1JSON},
		{"testCompleteTransactionWithMock/Nil", nil, nil, completeTransactionResult1JSON},
		{"testCompleteTransactionWithMock/Error", id2, pass2, completeTransactionResult2JSON},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompleteTransaction(tt.id, tt.password); C.GoString(got) != C.GoString(tt.want) {
				assert.Equal(t, C.GoString(tt.want), C.GoString(got))
			}
		})
	}
}

func testCompleteTransactionsWithMock(t *testing.T) {
	realStatusAPI := statusAPI
	defer func() { statusAPI = realStatusAPI }()

	// Setup Mock StatusAPI
	ctrl := gomock.NewController(t)
	status := NewMocklibStatusAPI(ctrl)
	statusAPI = status
	status.EXPECT().CompleteTransactions([]common.QueuedTxID{"id1"}, "pass1").Return(
		common.CompleteTransactionsResult{
			Results: map[common.QueuedTxID]common.CompleteTransactionResult{"id1": {ID: "id1", Hash: "0x123"}}})

	status.EXPECT().CompleteTransactions([]common.QueuedTxID{"id2"}, "pass2").Return(
		common.CompleteTransactionsResult{
			Results: map[common.QueuedTxID]common.CompleteTransactionResult{"id2": {ID: "id2", Hash: "", Error: "test error"}}})

	status.EXPECT().CompleteTransactions([]common.QueuedTxID{"id3", "id4"}, "pass2").Return(
		common.CompleteTransactionsResult{
			Results: map[common.QueuedTxID]common.CompleteTransactionResult{"id3": {ID: "id3", Hash: "0x456", Error: ""}, "id4": {ID: "id4", Hash: "0x789", Error: ""}}})

	// C Strings
	id1 := C.CString(`["id1"]`)
	pass1 := C.CString("pass1")
	invalidJSON := C.CString(`id`)
	id2 := C.CString(`["id2"]`)
	pass2 := C.CString("pass2")
	id3 := C.CString(`["id3", "id4"]`)
	empty := C.CString("")
	completeTransactionResult1JSON := C.CString(`{"results": {"id1": {"id":"id1","hash":"0x123","error":""}}}`)
	completeTransactionResult2JSON := C.CString(`{"results": {"none": {"id":"","hash":"","error":"invalid character 'i' looking for beginning of value"}}}`)
	completeTransactionResult3JSON := C.CString(`{"results": {"none": {"id":"","hash":"","error":"unexpected end of JSON input"}}}`)
	completeTransactionResult4JSON := C.CString(`{"results": {"id2": {"id":"id2","hash":"","error":"test error"}}}`)
	completeTransactionResult5JSON := C.CString(`{"results": {"id3": {"id":"id3","hash":"0x456","error":""}, "id4": {"id":"id4","hash":"0x789","error":""}}}`)

	defer func() {
		C.free(unsafe.Pointer(id1))
		C.free(unsafe.Pointer(pass1))
		C.free(unsafe.Pointer(invalidJSON))
		C.free(unsafe.Pointer(id2))
		C.free(unsafe.Pointer(pass2))
		C.free(unsafe.Pointer(id3))
		C.free(unsafe.Pointer(empty))
		C.free(unsafe.Pointer(completeTransactionResult1JSON))
		C.free(unsafe.Pointer(completeTransactionResult2JSON))
		C.free(unsafe.Pointer(completeTransactionResult3JSON))
		C.free(unsafe.Pointer(completeTransactionResult4JSON))
		C.free(unsafe.Pointer(completeTransactionResult5JSON))

	}()

	tests := []struct {
		name     string
		ids      *C.char
		password *C.char
		want     *C.char
	}{
		{"testCompleteTransactionsWithMock/Normal", id1, pass1, completeTransactionResult1JSON},
		{"testCompleteTransactionsWithMock/InvalidJSON", invalidJSON, pass1, completeTransactionResult2JSON},
		{"testCompleteTransactionWithMock/Empty", empty, empty, completeTransactionResult3JSON},
		{"testCompleteTransactionWithMock/Nil", nil, nil, completeTransactionResult3JSON},
		{"testCompleteTransactionWithMock/Error", id2, pass2, completeTransactionResult4JSON},
		{"testCompleteTransactionWithMock/2ID", id3, pass2, completeTransactionResult5JSON},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompleteTransactions(tt.ids, tt.password); C.GoString(got) != C.GoString(tt.want) {
				assert.JSONEq(t, C.GoString(tt.want), C.GoString(got))
			}
		})
	}
}

func testDiscardTransactionWithMock(t *testing.T) {
	realStatusAPI := statusAPI
	defer func() { statusAPI = realStatusAPI }()

	// Setup Mock StatusAPI
	ctrl := gomock.NewController(t)
	status := NewMocklibStatusAPI(ctrl)
	statusAPI = status
	status.EXPECT().DiscardTransaction(common.QueuedTxID("id1")).Return(common.DiscardTransactionResult{ID: "id1"}, nil)
	status.EXPECT().DiscardTransaction(common.QueuedTxID("")).Return(common.DiscardTransactionResult{ID: "id1"}, nil).AnyTimes()
	status.EXPECT().DiscardTransaction(common.QueuedTxID("id2")).Return(common.DiscardTransactionResult{ID: "id2", Error: "Test Error"}, fmt.Errorf("Test Error"))

	// C Strings
	id1 := C.CString("id1")
	id2 := C.CString("id2")
	empty := C.CString("")
	discardTransactionResult1JSON := C.CString(`{"id":"id1","error":""}`)
	discardTransactionResult2JSON := C.CString(`{"id":"id2","error":"Test Error"}`)
	defer func() {
		C.free(unsafe.Pointer(id1))
		C.free(unsafe.Pointer(id2))
		C.free(unsafe.Pointer(empty))
		C.free(unsafe.Pointer(discardTransactionResult1JSON))
		C.free(unsafe.Pointer(discardTransactionResult2JSON))

	}()

	tests := []struct {
		name string
		id   *C.char
		want *C.char
	}{
		{"testDiscardTransactionWithMock/Normal", id1, discardTransactionResult1JSON},
		{"testDiscardTransactionWithMock/Empty", empty, discardTransactionResult1JSON},
		{"testDiscardTransactionWithMock/Nil", nil, discardTransactionResult1JSON},
		{"testDiscardTransactionWithMock/Error", id2, discardTransactionResult2JSON},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiscardTransaction(tt.id); C.GoString(got) != C.GoString(tt.want) {
				assert.Equal(t, C.GoString(tt.want), C.GoString(got))
			}
		})
	}
}
