package hasher

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bobllor/assert"
)

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
	password := "anothertestpassword"
	baseSalt := []byte("A7iRBwsrtjiNOhnWeAGgng")

	raw, err := Hash(password, baseSalt, DefaultArgon2Params)
	assert.Nil(t, err)

	// generated from a random site online (argon2 online) for comparison
	baseEncode := "$argon2id$v=19$m=65536,t=2,p=4$QTdpUkJ3c3J0amlOT2huV2VBR2duZw$vzICl8p5CVfpGfypDV4yIVULsYatAmir6B8nHWtcPtE"

	assert.Equal(t, raw.Encode(), baseEncode)
}

func TestParsePHC(t *testing.T) {
	password := "password"

	raw, err := Hash(password, nil, DefaultArgon2Params)
	assert.Nil(t, err)

	phc := raw.Encode()

	hashRes, err := ParsePHC(phc)
	assert.Nil(t, err)

	assert.Equal(t, strings.Contains(phc, hashRes.Hash), true)
	assert.Equal(t, strings.Contains(phc, hashRes.Salt), true)
	assert.Equal(t, strings.Contains(phc, fmt.Sprintf("m=%d", hashRes.Params.Memory)), true)
	assert.Equal(t, strings.Contains(phc, fmt.Sprintf("t=%d", hashRes.Params.Time)), true)
	assert.Equal(t, strings.Contains(phc, fmt.Sprintf("p=%d", hashRes.Params.Threads)), true)
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
	}

	for _, paramStr := range invalidStrings {
		_, _, _, err := parseArgon2ParamString(paramStr)
		assert.NotNil(t, err)
	}
}
