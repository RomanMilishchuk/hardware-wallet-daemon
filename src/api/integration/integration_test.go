package integration

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/skycoin/hardware-wallet-daemon/src/api"

	"github.com/andreyvit/diff"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/hardware-wallet-daemon/src/client"
	"github.com/skycoin/hardware-wallet-daemon/src/client/operations"
	"github.com/skycoin/hardware-wallet-daemon/src/models"
)

const (
	testModeEmulator = "EMULATOR"
	testModeUSB      = "USB"

	testFixturesDir = "testdata"
)

type TestData struct {
	actual   interface{}
	expected interface{}
}

var update = flag.Bool("update", false, "update golden files")
var daemonClient = client.Default

func useCSRF(t *testing.T) bool {
	x := os.Getenv("USE_CSRF")

	if x == "" {
		return false
	}

	useCSRF, err := strconv.ParseBool(x)
	require.NoError(t, err)
	return useCSRF
}

// addCSRFHeader is used to add the CSRF Header param
func addCSRFHeader(t *testing.T, c *client.HardwareWalletDaemon) runtime.ClientAuthInfoWriterFunc {
	return func(req runtime.ClientRequest, _ strfmt.Registry) error {
		if useCSRF(t) {
			csrfResp, err := c.Operations.GetCsrf(nil)
			require.NoError(t, err)
			require.NotNil(t, csrfResp.Payload.Data)
			return req.SetHeaderParam(api.CSRFHeaderName, csrfResp.Payload.Data)
		}

		return nil
	}
}

func mode(t *testing.T) string {
	mode := os.Getenv("HW_DAEMON_INTEGRATION_TEST_MODE")
	switch mode {
	case "":
		mode = testModeEmulator
	case testModeUSB, testModeEmulator:
	default:
		t.Fatalf("Invalid test mode %s, must be emulator or wallet", mode)
	}
	return mode
}

func enabled() bool {
	return os.Getenv("HW_DAEMON_INTEGRATION_TESTS") == "1"
}

func doWallet(t *testing.T) bool {
	if enabled() && mode(t) == testModeUSB {
		return true
	}

	t.Skip("wallet tests disabled")
	return false
}

func doEmulator(t *testing.T) bool {
	if enabled() && mode(t) == testModeEmulator {
		return true
	}

	t.Skip("emulator tests disabled")
	return false
}

func doEmulatorOrWallet(t *testing.T) bool {
	if enabled() {
		switch mode(t) {
		case testModeEmulator, testModeUSB:
			return true
		}
	}

	t.Skip("emulator tests disabled")
	return false
}

func TestVersion(t *testing.T) {
	if !enabled() {
		return
	}

	// check that mode is valid
	mode(t)

	resp, err := daemonClient.Operations.GetVersion(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.NotNil(t, resp.Payload.Data)

	buildInfo := api.BuildInfo{
		Version: resp.Payload.Data.Version,
		Commit:  resp.Payload.Data.Commit,
		Branch:  resp.Payload.Data.Branch,
	}

	_, err = buildInfo.Semver()
	require.NoError(t, err)
}

func TestApplySettings(t *testing.T) {
	if *update {
		t.SkipNow()
	}

	if !doEmulatorOrWallet(t) {
		return
	}

	params := operations.NewPostApplySettingsParams()
	params.ApplySettingsRequest = &models.ApplySettingsRequest{
		Label: "skywallet",
	}

	resp, err := daemonClient.Operations.PostApplySettings(params, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, resp.Payload.Data[0], "ButtonRequest")

	buttonResp, err := daemonClient.Operations.PostIntermediateButton(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, "Settings applied", buttonResp.Payload.Data[0])
}

func TestBackup(t *testing.T) {
	if *update {
		t.SkipNow()
	}

	if !doEmulatorOrWallet(t) {
		return
	}

	bootstrap(t)

	resp, err := daemonClient.Operations.PostBackup(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, resp.Payload.Data[0], "ButtonRequest")

	for {
		buttonResp, err := daemonClient.Operations.PostIntermediateButton(nil, addCSRFHeader(t, daemonClient))
		require.NoError(t, err)
		if buttonResp.Payload.Data[0] != "ButtonRequest" {
			require.Equal(t, "Device backed up!", buttonResp.Payload.Data[0])
			break
		}
	}
}

func TestCheckMessageSignature(t *testing.T) {
	if *update {
		t.SkipNow()
	}

	if !doEmulatorOrWallet(t) {
		return
	}

	params := operations.NewPostCheckMessageSignatureParams()
	params.CheckMessageSignatureRequest = &models.CheckMessageSignatureRequest{
		Address:   newStrPtr("2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw"),
		Message:   newStrPtr("Hello World"),
		Signature: newStrPtr("6ebd63dd5e57cad07b6d229e96b5d2ac7d1bec1466d2a95bd200c21be6a0bf194b5ad5123f6e37c6393ee3635b38b938fcd91bbf1327fc957849a9e5736f6e4300"),
	}

	resp, err := daemonClient.Operations.PostCheckMessageSignature(params, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", resp.Payload.Data[0])
}

func TestEmulatorFeatures(t *testing.T) {
	if !doEmulator(t) {
		return
	}

	bootstrap(t)

	resp, err := daemonClient.Operations.GetFeatures(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)

	var expected models.FeaturesResponse

	// set variable parameters to empty string
	resp.Payload.Data.DeviceID = "foo"
	resp.Payload.Data.Label = "foo"
	resp.Payload.Data.BootloaderHash = "foo"

	checkGoldenFile(t, "features-emulator.golden", TestData{*resp.Payload, &expected})
}

func TestGenerateMnemonic(t *testing.T) {
	if *update {
		t.SkipNow()
	}

	if !doEmulatorOrWallet(t) {
		return
	}

	// wipe existing data
	resp, err := daemonClient.Operations.DeleteWipe(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, resp.Payload.Data[0], "ButtonRequest")

	buttonResp, err := daemonClient.Operations.PostIntermediateButton(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, "Device wiped", buttonResp.Payload.Data[0])

	mnemonicParams := operations.NewPostGenerateMnemonicParams()
	mnemonicParams.GenerateMnemonicRequest = &models.GenerateMnemonicRequest{
		WordCount: newInt64Ptr(12),
	}

	mnemonicResp, err := daemonClient.Operations.PostGenerateMnemonic(mnemonicParams, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, "Mnemonic successfully configured", mnemonicResp.Payload.Data[0])
}

func TestRecovery(t *testing.T) {
	if *update {
		t.SkipNow()
	}

	if !doEmulatorOrWallet(t) {
		return
	}

	// wipe existing data
	resp, err := daemonClient.Operations.DeleteWipe(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, resp.Payload.Data[0], "ButtonRequest")

	buttonResp, err := daemonClient.Operations.PostIntermediateButton(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, "Device wiped", buttonResp.Payload.Data[0])

	params := operations.NewPostRecoveryParams()
	params.RecoveryRequest = &models.RecoveryRequest{
		WordCount: newInt64Ptr(12),
	}

	recoveryResp, err := daemonClient.Operations.PostRecovery(params, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, recoveryResp.Payload.Data[0], "ButtonRequest")

	buttonResp, err = daemonClient.Operations.PostIntermediateButton(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, "WordRequest", buttonResp.Payload.Data[0])

	wordParams := operations.NewPostIntermediateWordParams()
	wordParams.WordRequest = &models.WordRequest{
		Word: newStrPtr("foobar"),
	}

	wordParamsResp, err := daemonClient.Operations.PostIntermediateWord(wordParams, addCSRFHeader(t, daemonClient))
	require.Nil(t, wordParamsResp)
	// match that it contains any of the two available error responses.
	require.Subset(t, [2]string{"Wrong word retyped", "Word not found in a wordlist"}, [1]string{err.Error()})
}

func TestSetMnemonic(t *testing.T) {
	if *update {
		t.SkipNow()
	}

	if !doEmulatorOrWallet(t) {
		return
	}

	mnemonic := "cloud flower upset remain green metal below cup stem infant art thank"
	params := operations.NewPostSetMnemonicParams()
	params.SetMnemonicRequest = &models.SetMnemonicRequest{
		Mnemonic: newStrPtr(mnemonic),
	}

	resp, err := daemonClient.Operations.PostSetMnemonic(params, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, resp.Payload.Data[0], "ButtonRequest")

	buttonResp, err := daemonClient.Operations.PostIntermediateButton(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, mnemonic, buttonResp.Payload.Data[0])
}

func TestConfigurePinCode(t *testing.T) {
	if *update {
		t.SkipNow()
	}

	if !doEmulatorOrWallet(t) {
		return
	}

	pinParams := operations.NewPostConfigurePinCodeParams()
	pinParams.ConfigurePinCodeRequest = &models.ConfigurePinCodeRequest{
		RemovePin: newBoolPtr(false),
	}

	resp, err := daemonClient.Operations.PostConfigurePinCode(pinParams, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, resp.Payload.Data[0], "ButtonRequest")

	buttonResp, err := daemonClient.Operations.PostIntermediateButton(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, "PinMatrixRequest", buttonResp.Payload.Data[0])

	params := operations.NewPostIntermediatePinMatrixParams()
	params.PinMatrixRequest = &models.PinMatrixRequest{
		Pin: newStrPtr("123"),
	}

	pinAckResp, err := daemonClient.Operations.PostIntermediatePinMatrix(params, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, "PinMatrixRequest", pinAckResp.Payload.Data[0])

	params = operations.NewPostIntermediatePinMatrixParams()
	params.PinMatrixRequest = &models.PinMatrixRequest{
		Pin: newStrPtr("123"),
	}

	pinAckResp, err = daemonClient.Operations.PostIntermediatePinMatrix(params, addCSRFHeader(t, daemonClient))
	require.Nil(t, pinAckResp)
	require.Subset(t, [2]string{"PIN mismatch", "Pin invalid"}, [1]string{err.Error()})
}

func TestWipe(t *testing.T) {
	if *update {
		t.SkipNow()
	}

	if !doEmulatorOrWallet(t) {
		return
	}

	resp, err := daemonClient.Operations.DeleteWipe(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, resp.Payload.Data[0], "ButtonRequest")

	buttonResp, err := daemonClient.Operations.PostIntermediateButton(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, "Device wiped", buttonResp.Payload.Data[0])
}

func TestGenerateAddresses(t *testing.T) {
	if !doEmulatorOrWallet(t) {
		return
	}

	bootstrap(t)

	params := operations.NewPostGenerateAddressesParams()
	params.GenerateAddressesRequest = &models.GenerateAddressesRequest{
		AddressN:       newInt64Ptr(2),
		ConfirmAddress: false,
		StartIndex:     0,
	}

	resp, err := daemonClient.Operations.PostGenerateAddresses(params, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)

	var expected models.GenerateAddressesResponse
	checkGoldenFile(t, "generate-addresses.golden", TestData{*resp.Payload, &expected})
}

func TestWalletFeatures(t *testing.T) {
	if !doWallet(t) {
		return
	}

	bootstrap(t)

	resp, err := daemonClient.Operations.GetFeatures(nil, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)

	var expected models.FeaturesResponse

	// set variable parameters to empty string
	resp.Payload.Data.DeviceID = "foo"
	resp.Payload.Data.Label = "foo"
	resp.Payload.Data.BootloaderHash = "foo"

	checkGoldenFile(t, "features-wallet.golden", TestData{*resp.Payload, &expected})
}

func TestTransactionSign(t *testing.T) {
	if *update {
		t.SkipNow()
	}

	if !doEmulatorOrWallet(t) {
		return
	}

	bootstrap(t)

	params := operations.NewPostTransactionSignParams()
	params.TransactionSignRequest = &models.TransactionSignRequest{
		TransactionInputs: []*models.TransactionInput{
			{
				Index: newInt64Ptr(0),
				Hash:  newStrPtr("181bd5656115172fe81451fae4fb56498a97744d89702e73da75ba91ed5200f9"),
			},
		},

		TransactionOutputs: []*models.TransactionOutput{
			{
				Address: newStrPtr("K9TzLrgqz7uXn3QJHGxmzdRByAzH33J2ot"),
				Coins:   newStrPtr("0.1"),
				Hours:   newStrPtr("2"),
			},
		},
	}

	resp, err := daemonClient.Operations.PostTransactionSign(params, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, resp.Payload.Data[0], "ButtonRequest")

	var signature string
	for {
		buttonResp, err := daemonClient.Operations.PostIntermediateButton(nil, addCSRFHeader(t, daemonClient))
		require.NoError(t, err)
		if buttonResp.Payload.Data[0] != "ButtonRequest" {
			require.Len(t, buttonResp.Payload.Data, 1)
			signature = buttonResp.Payload.Data[0]
			break
		}
	}

	// verify the message signature
	fmt.Println(signature)
	verifParams := operations.NewPostCheckMessageSignatureParams()
	verifParams.CheckMessageSignatureRequest = &models.CheckMessageSignatureRequest{
		Address:   newStrPtr("2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw"),
		Message:   newStrPtr("d11c62b1e0e9abf629b1f5f4699cef9fbc504b45ceedf0047ead686979498218"),
		Signature: newStrPtr(signature),
	}

	verifResp, err := daemonClient.Operations.PostCheckMessageSignature(verifParams, addCSRFHeader(t, daemonClient))
	require.NoError(t, err)
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", verifResp.Payload.Data[0])
}

func TestWalletConnected(t *testing.T) {
	if *update {
		t.SkipNow()
	}

	if !doWallet(t) {
		return
	}

	resp, err := daemonClient.Operations.GetAvailable(nil, nil)
	require.NoError(t, err)
	require.Equal(t, resp.Payload.Data[0], true)
}

func bootstrap(t *testing.T) {
	if enabled() {
		// wipe existing data
		resp, err := daemonClient.Operations.DeleteWipe(nil, addCSRFHeader(t, daemonClient))
		require.NoError(t, err)
		require.Equal(t, resp.Payload.Data[0], "ButtonRequest")

		buttonResp, err := daemonClient.Operations.PostIntermediateButton(nil, addCSRFHeader(t, daemonClient))
		require.NoError(t, err)
		require.Equal(t, "Device wiped", buttonResp.Payload.Data[0])

		// set mnemonic
		mnemonic := "cloud flower upset remain green metal below cup stem infant art thank"
		mnemonicParams := operations.NewPostSetMnemonicParams()
		mnemonicParams.SetMnemonicRequest = &models.SetMnemonicRequest{
			Mnemonic: newStrPtr(mnemonic),
		}

		mnemonicResp, err := daemonClient.Operations.PostSetMnemonic(mnemonicParams, addCSRFHeader(t, daemonClient))
		require.NoError(t, err)
		require.Equal(t, mnemonicResp.Payload.Data[0], "ButtonRequest")

		buttonResp, err = daemonClient.Operations.PostIntermediateButton(nil, addCSRFHeader(t, daemonClient))
		require.NoError(t, err)
		require.Equal(t, mnemonic, buttonResp.Payload.Data[0])
	}
}

func newStrPtr(s string) *string {
	return &s
}

func newInt64Ptr(n int64) *int64 {
	return &n
}

func newBoolPtr(b bool) *bool {
	return &b
}

func loadGoldenFile(t *testing.T, filename string, testData TestData) {
	require.NotEmpty(t, filename, "loadGoldenFile golden filename missing")

	goldenFile := filepath.Join(testFixturesDir, filename)

	if *update {
		updateGoldenFile(t, goldenFile, testData.actual)
	}

	f, err := os.Open(goldenFile)
	require.NoError(t, err)
	defer f.Close()

	d := json.NewDecoder(f)
	d.DisallowUnknownFields()

	err = d.Decode(testData.expected)
	require.NoError(t, err, filename)
}

func updateGoldenFile(t *testing.T, filename string, content interface{}) {
	contentJSON, err := json.MarshalIndent(content, "", "\t")
	require.NoError(t, err)
	contentJSON = append(contentJSON, '\n')
	err = ioutil.WriteFile(filename, contentJSON, 0644)
	require.NoError(t, err)
}

func checkGoldenFile(t *testing.T, goldenFile string, td TestData) {
	loadGoldenFile(t, goldenFile, td)
	require.Equal(t, reflect.Indirect(reflect.ValueOf(td.expected)).Interface(), td.actual)

	// Serialize expected to JSON and compare to the goldenFile's contents
	// This will detect field changes that could be missed otherwise
	b, err := json.MarshalIndent(td.expected, "", "\t")
	require.NoError(t, err)

	goldenFile = filepath.Join(testFixturesDir, goldenFile)

	f, err := os.Open(goldenFile)
	require.NoError(t, err)
	defer f.Close()

	c, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	sc := string(c)
	sb := string(b) + "\n"

	require.Equal(t, sc, sb, "JSON struct output differs from golden file, was a field added to the struct?\nDiff:\n"+diff.LineDiff(sc, sb))
}
