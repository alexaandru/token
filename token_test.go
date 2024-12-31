package token

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	for tl := MinTokenLength - 3; tl < MaxTokenLength+3; tl++ {
		t.Run(fmt.Sprint("len-", tl), func(t *testing.T) {
			t.Parallel()

			al := tl
			if al < MinTokenLength {
				al = MinTokenLength
			}

			if al > MaxTokenLength {
				al = MaxTokenLength
			}

			tk := New(tl)
			te := tk.Encode()

			if act := len(te); act > al {
				t.Fatalf("Expected encoded token length to be <= %d got %d", al, act)
			}

			nt, err := Decode(te)
			if err != nil {
				t.Fatalf("Expected no error when decoding %s (%d) got %v", te, tk, err)
			}

			if nt != tk {
				t.Fatalf("Expected %d got %d", tk, nt)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	t.Skip("Tested implicitly via TestNew")
}

func TestTokenEncode(t *testing.T) {
	t.Skip("Tested implicitly via TestNew")
}

func TestUnmarshalText(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		tok    string
		exp    Token
		expErr error
	}{
		{"", 0, ErrTokenTooSmall},
		{strings.Repeat("a", MaxTokenLength+1), 0, ErrTokenTooBig},
		{"~", 0, ErrInvalidCharacter},
		{"0123456789", 225557475374453, nil},
		{"ABCDEFGHIJ", 137815617453790883, nil},
		{"KLMNOPQRST", 275405677432207313, nil},
		{"UVWXYZabcd", 412995737410623743, nil},
		{"efghijklmn", 550585797389040173, nil},
		{"opqrstuvwx", 688175857367456603, nil},
		{"yz", 3781, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.tok, func(t *testing.T) {
			t.Parallel()

			nt := Token(0)
			if err := nt.UnmarshalText([]byte(tc.tok)); !errors.Is(err, tc.expErr) {
				t.Fatalf("Expected error %v got %v", tc.expErr, err)
			} else if err == nil && nt != tc.exp {
				t.Fatalf("Expected token %d got %d", tc.exp, nt)
			}
		})
	}
}

func TestMarshalText(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		exp    string
		tok    Token
		expErr error
	}{
		{"", 0, nil},
		{"123456789", 225557475374453, nil},
		{"ABCDEFGHIJ", 137815617453790883, nil},
		{"KLMNOPQRST", 275405677432207313, nil},
		{"UVWXYZabcd", 412995737410623743, nil},
		{"efghijklmn", 550585797389040173, nil},
		{"opqrstuvwx", 688175857367456603, nil},
		{"yz", 3781, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.exp, func(t *testing.T) {
			t.Parallel()

			if act, err := tc.tok.MarshalText(); !errors.Is(err, tc.expErr) {
				t.Fatalf("Expected error %v got %v", tc.expErr, err)
			} else if err == nil && string(act) != tc.exp {
				t.Fatalf("Expected token %q got %q", tc.exp, string(act))
			}
		})
	}
}

func TestMaxHashInt(t *testing.T) {
	t.Skip("Tested implicitly via TestNew")
}
