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

// Hash takes a password, a salt, and params to generate and return a RawHash.
//
// salt can be nil or a predefined salt. If nil is given, then a random salt
// is generated using Argon2Params.SaltLength given in params.
func Hash(password string, salt []byte, params Argon2Params) (*RawHash, error) {
	if salt == nil {
		newSalt, err := getSalt(params.SaltLength)
		if err != nil {
			return nil, err
		}

		salt = newSalt
	}

	hash := argon2.IDKey(
		[]byte(password),
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
// containing the information that resulted in the hashed password
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
	paramStr := phcSplit[2]
	salt := phcSplit[3]
	hash := phcSplit[4]

	memory, time, threads, err := parseArgon2ParamString(paramStr)
	if err != nil {
		return nil, err
	}

	hashRes := &HashResult{
		Salt: salt,
		Hash: hash,
		Params: Argon2Params{
			SaltLength: len(salt),
			Memory:     memory,
			Time:       time,
			Threads:    threads,
			KeyLength:  uint32(len(hash)),
		},
	}

	return hashRes, nil
}

type RawHash struct {
	Salt   []byte
	Hash   []byte
	Params Argon2Params
}

// Encode encodes the RawHash data into a hash string in the PHC format.
func (rh *RawHash) Encode() string {
	salt := base64.RawStdEncoding.EncodeToString(rh.Salt)
	hash := base64.RawStdEncoding.EncodeToString(rh.Hash)

	// $<algorithm>$v=<version>$m=<memory>,t=<time>,p=<parallelism>$<salt>$<hash>
	// parallelism is equal to the Threads option of Argon2Params
	phcString := "$%s$v=%d$m=%d,t=%d,p=%d$%s$%s"

	argonStr := fmt.Sprintf(
		phcString,
		Argon2ID,
		ArgonVersion,
		rh.Params.Memory,
		rh.Params.Time,
		rh.Params.Threads,
		salt,
		hash,
	)

	return argonStr
}

type HashResult struct {
	Salt   string
	Hash   string
	Params Argon2Params
}

// Decode decodes the hash result back into a raw hash.
func (hr *HashResult) Decode() (*RawHash, error) {
	salt, err := base64.RawStdEncoding.DecodeString(hr.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %v", err)
	}
	hash, err := base64.RawStdEncoding.DecodeString(hr.Hash)
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
