package internal

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/segmentio/ksuid"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var source = rand.NewSource(time.Now().UnixNano())

type OnetimecodeType string

const (
	ANumberedCode               OnetimecodeType = "ANumberedCode"
	AnAlphaNumericCode                          = "AnAlphaNumericCode"
	AnAlphaNumericUpperCaseCode                 = "AnAlphaNumericUpperCaseCode"
)

type OnetimecodeConfig func(code *Onetimecode)

type Onetimecode struct {
	codeType      OnetimecodeType
	length        int
	min           int
	max           int
	ulmcase       int
	withoutDashes bool
	groupBy       string
	groupEvery    int
	groupingMode  bool
	code          int64
	stringCode    string
}

func WithLength(length int) OnetimecodeConfig {
	return func(code *Onetimecode) {
		if length > -1 {
			if code.codeType == ANumberedCode && length > 19 {
				length = 19
			}
			code.length = length
			code.min = int(math.Pow(10, float64(length-1)))
			code.max = int(math.Pow(10, float64(length))) - 1
		}
	}
}

func WithMin(min int) OnetimecodeConfig {
	return func(code *Onetimecode) {
		if min > -1 {
			code.min = min
		}
	}
}

func WithMax(max int) OnetimecodeConfig {
	return func(code *Onetimecode) {
		if max > -1 {
			code.max = max
		}
	}
}

func WithMinMax(min, max int) OnetimecodeConfig {
	return func(code *Onetimecode) {
		code.min = min
		code.max = max
	}
}

func WithAlphaNumericCode() OnetimecodeConfig {
	return func(code *Onetimecode) {
		code.codeType = AnAlphaNumericCode
	}
}

func WithCase(caseStr string) OnetimecodeConfig {
	switch caseStr {
	case "lower":
		return WithLowerCase()
	case "upper":
		return WithUpperCase()
	}
	return func(code *Onetimecode) {}
}

func WithUpperCase() OnetimecodeConfig {
	return func(code *Onetimecode) {
		code.ulmcase = 1
	}
}

func WithLowerCase() OnetimecodeConfig {
	return func(code *Onetimecode) {
		code.ulmcase = -1
	}
}

func WithoutDashesFromBoolean(withoutDashes bool) OnetimecodeConfig {
	return func(code *Onetimecode) {
		if withoutDashes == true {
			code.withoutDashes = true
		}
	}
}

func WithoutDashes() OnetimecodeConfig {
	return func(code *Onetimecode) {
		code.withoutDashes = true
	}
}

func WithGrouping(groupEvery int, groupBy string) OnetimecodeConfig {
	return func(code *Onetimecode) {
		if groupEvery > -1 || groupBy != "" {
			if groupEvery < code.length {
				if groupEvery > -1 {
					code.groupEvery = groupEvery
				} else {
					code.groupEvery = 4
				}
				if groupBy != "" {
					code.groupBy = groupBy
				} else {
					code.groupBy = "-"
				}
				code.groupingMode = true
			}
		}
	}
}

func NewNumericalCode(opts ...OnetimecodeConfig) *Onetimecode {
	otc := &Onetimecode{
		codeType: ANumberedCode,
		length:   6,
		min:      1,
		max:      999999,
	}
	for _, opt := range opts {
		opt(otc)
	}
	otc.defineValueNumeric()
	if otc.groupingMode {
		otc.applyGrouping()
	}
	return otc
}

func NewAlphanumericalCode(opts ...OnetimecodeConfig) *Onetimecode {
	otc := &Onetimecode{
		codeType: AnAlphaNumericCode,
		length:   6,
		ulmcase:  0,
	}
	for _, opt := range opts {
		opt(otc)
	}
	otc.defineValueAlphanumeric()
	if otc.groupingMode {
		otc.applyGrouping()
	}
	return otc
}

func NewKsuidCode(opts ...OnetimecodeConfig) *Onetimecode {
	otc := &Onetimecode{}
	for _, opt := range opts {
		opt(otc)
	}
	otc.stringCode = Ksuid()
	return otc
}

func NewUuidCode(opts ...OnetimecodeConfig) *Onetimecode {
	otc := &Onetimecode{
		withoutDashes: false,
	}
	for _, opt := range opts {
		opt(otc)
	}
	otc.stringCode = Uuid()
	if otc.withoutDashes == true {
		otc.stringCode = strings.Replace(otc.stringCode, "-", "", -1)
	}
	return otc
}

func NewUlidCode(opts ...OnetimecodeConfig) *Onetimecode {
	otc := &Onetimecode{}
	for _, opt := range opts {
		opt(otc)
	}
	otc.stringCode = Ulid()
	//if otc.withoutDashes == true {
	//	otc.stringCode = strings.Replace(otc.stringCode, "-", "", -1)
	//}
	return otc
}

func (otc *Onetimecode) ResultAsString() string {
	//if otc.stringCode == "" {
	//	otc.stringCode = string(otc.code)
	//}
	return otc.stringCode
}

func (otc *Onetimecode) NumberCode() int64 {
	return otc.code
}

func (otc *Onetimecode) defineValueNumeric() {
	rand.Seed(time.Now().UnixNano())
	rndNr := rand.Intn(otc.max-otc.min+1) + otc.min
	otc.code = int64(rndNr)
	otc.stringCode = strconv.FormatInt(otc.code, 10)
	if len(otc.stringCode) < otc.length {
		otc.stringCode = strings.Repeat("0", (otc.length-len(otc.stringCode))) + otc.stringCode
	}
}

func (otc *Onetimecode) defineValueAlphanumeric() {
	otc.stringCode = alphaNumberCode(otc.length, otc.ulmcase)
}

func (otc *Onetimecode) applyGrouping() {
	var buffer bytes.Buffer
	before := otc.groupEvery - 1
	last := len(otc.stringCode) - 1
	for i, char := range otc.stringCode {
		buffer.WriteRune(char)
		if i%otc.groupEvery == before && i != last {
			buffer.WriteRune([]rune(otc.groupBy)[0])
		}
	}
	otc.stringCode = buffer.String()
}

// AlphaNumberCode returns an alphanumeric randomized
// code of the given length with numbers, uppercase
// and lowercase characters.
func alphaNumberCode(length int, ulmcase int) string {
	const charsetMixed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const charsetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const charsetLower = "abcdefghijklmnopqrstuvwxyz0123456789"
	var charset string
	switch ulmcase {
	case -1:
		charset = charsetLower
	case 1:
		charset = charsetUpper
	default:
		charset = charsetMixed
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[source.Int63()%int64(len(charset))]
	}
	return string(b)
}

// Ksuid returns a KSUID code.
func Ksuid() string {
	return ksuid.New().String()
}

// Uuid returns an UUID code.
func Uuid() string {
	return uuid.New().String()
}

// Ulid returns an ULID code.
func Ulid() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
