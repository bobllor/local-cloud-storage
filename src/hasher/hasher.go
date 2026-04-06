package hasher

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	ArgonVersion = 19
	Argon2ID     = "argon2id"
)

var DefaultArgon2Params = Argon2Params{
	SaltLength: 16,
	Time:       2,
	Memory:     1024 * 64,
	Threads:    4,
	KeyLength:  32,
}

type Argon2Params struct {
	SaltLength int
	Time       uint32
	Memory     uint32
	Threads    uint8
	KeyLength  uint32
}

// Hash takes a string, a salt, and params to generate and return a RawHash.
//
// salt can be nil or a predefined salt. If nil is given, then a random salt
// is generated using Argon2Params.SaltLength given in params.
func Hash(str string, salt []byte, params Argon2Params) (*RawHash, error) {
	if salt == nil {
		newSalt, err := getSalt(params.SaltLength)
		if err != nil {
			return nil, err
		}

		salt = newSalt
	}

	hash := argon2.IDKey(
		[]byte(str),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLength,
	)

	rh := &RawHash{
		Salt:   salt,
		Hash:   hash,
		Params: params,
	}

	return rh, nil
}

// ParsePHC parses a PHC string and returns a HashResult
// containing the information that results in the hashed string
// with Argon2ID.
func ParsePHC(phc string) (*HashResult, error) {
	// 6 is the valid length
	phcSplit := strings.Split(phc, "$")

	if len(phcSplit) != 6 {
		return nil, fmt.Errorf(
			"invalid length of arguments found while parsing (got length %d instead of 6)",
			len(phcSplit),
		)
	}

	// dropping the first element due to it being an empty string
	phcSplit = phcSplit[1:]

	// only the last 3 elements are relevant for this project
	// argon2i and argon2d are not supported.
	paramStr := phcSplit[2]
	salt := phcSplit[3]
	hash := phcSplit[4]

	memory, time, threads, err := parseArgon2ParamString(paramStr)
	if err != nil {
		return nil, err
	}

	decodedSalt, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt string: %v", err)
	}
	decodedHash, err := base64.RawStdEncoding.DecodeString(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hash string: %v", err)
	}

	hashRes := &HashResult{
		Salt: salt,
		Hash: hash,
		Params: Argon2Params{
			SaltLength: len(decodedSalt),
			Memory:     memory,
			Time:       time,
			Threads:    threads,
			KeyLength:  uint32(len(decodedHash)),
		},
	}

	return hashRes, nil
}

// Compare hashes a given string and compares it to an existing HashResult.
// It will return true or false, or an error if one occurs.
func Compare(str string, hr *HashResult) (bool, error) {
	convRaw, err := hr.Decode()
	if err != nil {
		return false, err
	}

	raw, err := Hash(str, convRaw.Salt, convRaw.Params)
	if err != nil {
		return false, err
	}

	compareRes := raw.Encode()
	storedHash := hr.Hash

	if compareRes.Hash != storedHash {
		return false, nil
	}

	return true, nil
}

type RawHash struct {
	Salt   []byte
	Hash   []byte
	Params Argon2Params
}

// Encode encodes the RawHash data into a new HashResult.
func (rh *RawHash) Encode() *HashResult {
	salt := base64.RawStdEncoding.EncodeToString(rh.Salt)
	hash := base64.RawStdEncoding.EncodeToString(rh.Hash)

	// $<algorithm>$v=<version>$m=<memory>,t=<time>,p=<parallelism>$<salt>$<hash>
	// parallelism is equal to the Threads option of Argon2Params
	phcString := "$%s$v=%d$m=%d,t=%d,p=%d$%s$%s"

	phc := fmt.Sprintf(
		phcString,
		Argon2ID,
		ArgonVersion,
		rh.Params.Memory,
		rh.Params.Time,
		rh.Params.Threads,
		salt,
		hash,
	)

	res := &HashResult{
		Salt:   salt,
		Hash:   hash,
		Params: rh.Params,
		PHC:    phc,
	}

	return res
}

type HashResult struct {
	Salt   string
	Hash   string
	Params Argon2Params
	PHC    string
}

// Decode decodes the hash result back into a RawHash struct.
func (hr *HashResult) Decode() (*RawHash, error) {
	salt, err := hr.DecodeSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %v", err)
	}
	hash, err := hr.DecodeHash()
	if err != nil {
		return nil, fmt.Errorf("failed to decode hash: %v", err)
	}

	raw := &RawHash{
		Salt:   salt,
		Hash:   hash,
		Params: hr.Params,
	}

	return raw, nil
}

// DecodeSalt decodes the base64 salt string back to its raw form.
func (hr *HashResult) DecodeSalt() ([]byte, error) {
	salt, err := base64.RawStdEncoding.DecodeString(hr.Salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

// DecodeHash decodes the base64 hash string back to its raw form.
func (hr *HashResult) DecodeHash() ([]byte, error) {
	hash, err := base64.RawStdEncoding.DecodeString(hr.Hash)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func getSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

// parseArgon2ParamString parses the parameters to retrieve the
// memory, time, and paralleisms of the string.
func parseArgon2ParamString(paramStr string) (memory uint32, time uint32, threads uint8, err error) {
	paramPattern := "^m=([0-9]+),t=([0-9]+),p=([0-9]+)$"
	reg, err := regexp.Compile(paramPattern)
	if err != nil {
		return 0, 0, 0, err
	}

	bParam := []byte(paramStr)
	isMatch := reg.Match(bParam)
	if !isMatch {
		return 0, 0, 0, fmt.Errorf("%s is not a valid param string (format: m=<digit>,t=<digit>,p=<digit>)", paramStr)
	}

	// will return 4 submatches
	matches := reg.FindSubmatch(bParam)
	if len(matches) != 4 {
		return 0, 0, 0, fmt.Errorf("failed to find submatches of the param string")
	}

	memoryStr := string(matches[1])
	timeStr := string(matches[2])
	threadsStr := string(matches[3])

	intMemory, err := strconv.Atoi(memoryStr)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to convert memory string to a number: %v", err)
	}
	intTime, err := strconv.Atoi(timeStr)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to convert time string to a number: %v", err)
	}
	intThreads, err := strconv.Atoi(threadsStr)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to convert threads string to a number: %v", err)
	}

	memory = uint32(intMemory)
	time = uint32(intTime)
	threads = uint8(intThreads)

	return memory, time, threads, nil
}
