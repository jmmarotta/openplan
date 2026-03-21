package plan

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

const crockfordAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// Keep suffix validation derived from the same alphabet used for generation so
// the ID grammar cannot drift from the emitted IDs.
var idPattern = regexp.MustCompile(
	fmt.Sprintf(`^([A-Z][A-Z0-9]*)-([1-9][0-9]*)_([%s]{8})$`, regexp.QuoteMeta(crockfordAlphabet)),
)

var errInvalidID = errors.New("invalid plan id")

type ID struct {
	Prefix string
	Number int
	Suffix string
}

func ParseID(raw string) (ID, error) {
	raw = strings.TrimSpace(raw)
	matches := idPattern.FindStringSubmatch(raw)
	if matches == nil {
		return ID{}, fmt.Errorf("%w %q", errInvalidID, raw)
	}

	number, err := strconv.Atoi(matches[2])
	if err != nil {
		return ID{}, fmt.Errorf("%w %q", errInvalidID, raw)
	}

	return ID{
		Prefix: matches[1],
		Number: number,
		Suffix: matches[3],
	}, nil
}

func FormatID(id ID) string {
	return fmt.Sprintf("%s-%d_%s", id.Prefix, id.Number, id.Suffix)
}

func NewID(prefix string, existing []string, r io.Reader) (ID, error) {
	if r == nil {
		r = rand.Reader
	}

	maxNumber := 0
	seen := make(map[string]struct{}, len(existing))
	for _, raw := range existing {
		parsed, err := ParseID(raw)
		if err != nil {
			continue
		}
		if parsed.Prefix != prefix {
			continue
		}
		seen[raw] = struct{}{}
		if parsed.Number > maxNumber {
			maxNumber = parsed.Number
		}
	}

	// The numeric portion stays repo-local and human-friendly; the random suffix
	// carries the merge-friendly uniqueness.
	for range 256 {
		suffix, err := randomSuffix(r)
		if err != nil {
			return ID{}, err
		}
		candidate := ID{Prefix: prefix, Number: maxNumber + 1, Suffix: suffix}
		if _, ok := seen[FormatID(candidate)]; ok {
			continue
		}
		return candidate, nil
	}

	return ID{}, fmt.Errorf("unable to allocate unique plan id")
}

func FilenameForID(id string) string {
	return id + ".md"
}

func randomSuffix(r io.Reader) (string, error) {
	var buf [5]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return "", err
	}

	bits := uint64(buf[0])<<32 | uint64(buf[1])<<24 | uint64(buf[2])<<16 | uint64(buf[3])<<8 | uint64(buf[4])
	encoded := make([]byte, 8)
	for i := range encoded {
		shift := uint(35 - (i * 5))
		encoded[i] = crockfordAlphabet[(bits>>shift)&31]
	}

	return string(encoded), nil
}
