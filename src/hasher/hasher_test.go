package hasher

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bobllor/assert"
)

var baseHashInfo = struct {
	Password string
	Salt     []byte
	PHC      string
}{
	Password: "anothertestpassword",
	Salt:     []byte("A7iRBwsrtjiNOhnWeAGgng"),
	// generated from a random site online (argon2 online) for comparison
	PHC: "$argon2id$v=19$m=65536,t=2,p=4$QTdpUkJ3c3J0amlOT2huV2VBR2duZw$vzICl8p5CVfpGfypDV4yIVULsYatAmir6B8nHWtcPtE",
}

func TestSalt(t *testing.T) {
	saltSize := 32
	salt, err := getSalt(32)
	assert.Nil(t, err)

	assert.Equal(t, len(salt), saltSize)
}

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	_, err := Hash(password, nil, DefaultArgon2Params)
	assert.Nil(t, err)
}

func TestHashToString(t *testing.T) {
	raw, err := Hash(baseHashInfo.Password, baseHashInfo.Salt, DefaultArgon2Params)
	assert.Nil(t, err)

	res := raw.Encode()

	assert.Equal(t, res.PHC, baseHashInfo.PHC)
}

func TestParsePHC(t *testing.T) {
	hashRes, err := ParsePHC(baseHashInfo.PHC)
	assert.Nil(t, err)

	assert.Equal(t, strings.Contains(baseHashInfo.PHC, hashRes.Hash), true)
	assert.Equal(t, strings.Contains(baseHashInfo.PHC, hashRes.Salt), true)
	assert.Equal(t, strings.Contains(baseHashInfo.PHC, fmt.Sprintf("m=%d", hashRes.Params.Memory)), true)
	assert.Equal(t, strings.Contains(baseHashInfo.PHC, fmt.Sprintf("t=%d", hashRes.Params.Time)), true)
	assert.Equal(t, strings.Contains(baseHashInfo.PHC, fmt.Sprintf("p=%d", hashRes.Params.Threads)), true)
}

func TestTrueCompare(t *testing.T) {
	res, err := ParsePHC(baseHashInfo.PHC)
	assert.Nil(t, err)

	status, err := Compare(baseHashInfo.Password, res)
	assert.Nil(t, err)

	assert.Equal(t, status, true)
}

func TestFalseCompare(t *testing.T) {
	password := "fdsafdsafdsa"

	baseRes, err := ParsePHC(baseHashInfo.PHC)
	assert.Nil(t, err)

	status, err := Compare(password, baseRes)
	assert.Nil(t, err)

	assert.Equal(t, status, false)
}

func TestConvertRawToResToRaw(t *testing.T) {
	password := "password"

	raw, err := Hash(password, nil, DefaultArgon2Params)
	assert.Nil(t, err)

	hashRes := raw.Encode()

	parsedRes, err := ParsePHC(hashRes.PHC)
	assert.Nil(t, err)

	convRaw, err := parsedRes.Decode()
	assert.Nil(t, err)

	assert.Equal(t, convRaw.Params.KeyLength, raw.Params.KeyLength)
	assert.Equal(t, convRaw.Params.SaltLength, raw.Params.SaltLength)
}

func TestNormalParseParamStr(t *testing.T) {
	baseMemory := uint32(65565)
	baseTime := uint32(1)
	baseThreads := uint8(4)

	paramStr := fmt.Sprintf("m=%d,t=%d,p=%d", baseMemory, baseTime, baseThreads)

	memory, time, threads, err := parseArgon2ParamString(paramStr)
	assert.Nil(t, err)

	assert.Equal(t, memory, baseMemory)
	assert.Equal(t, time, baseTime)
	assert.Equal(t, threads, baseThreads)
}

func TestFailParseParamStr(t *testing.T) {
	invalidStrings := []string{
		"m=12345,t=1,p=asdf",
		"asdf",
		"m=123,t=1,p='24",
		"m=,t=,p=",
		"123,4,55",
		"m='3,t=1,p=4",
		"M=34,T=44,P=33",
	}

	for _, paramStr := range invalidStrings {
		_, _, _, err := parseArgon2ParamString(paramStr)
		assert.NotNil(t, err)
	}
}
