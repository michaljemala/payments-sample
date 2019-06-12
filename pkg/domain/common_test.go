package domain

import (
	"testing"

	"github.com/michaljemala/payments-sample/pkg/internal/errors"

	"github.com/gofrs/uuid"
)

func TestBaseObject_SetID(t *testing.T) {
	testCases := []struct {
		name    string
		in      string
		out     BaseObject
		errFunc func(*testing.T, error)
	}{
		{
			name: "Valid UUIDv4",
			in:   "276c8bbf-79ca-4ac2-b319-0f1c51463540",
			out:  BaseObject{ID: toID(0x27, 0x6c, 0x8b, 0xbf, 0x79, 0xca, 0x4a, 0xc2, 0xb3, 0x19, 0x0f, 0x1c, 0x51, 0x46, 0x35, 0x40)},
		},
		{
			name: "Invalid ID",
			in:   "I_AM_NOT_VALID",
			errFunc: func(t *testing.T, err error) {
				if err == nil {
					t.Fatalf("expected error, have <nil>")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := BaseObject{}
			err := o.SetID(tc.in)
			if err != nil {
				if tc.errFunc == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				tc.errFunc(t, err)
				return
			}
			if want, have := tc.out, o; want != have {
				t.Fatalf("invalid id: want %v, have %v", want, have)
			}
		})
	}
}

func TestIDFrom(t *testing.T) {
	testCases := []struct {
		name    string
		in      string
		out     ID
		errFunc func(*testing.T, error)
	}{
		{
			name: "Valid UUIDv4",
			in:   "276c8bbf-79ca-4ac2-b319-0f1c51463540",
			out:  toID(0x27, 0x6c, 0x8b, 0xbf, 0x79, 0xca, 0x4a, 0xc2, 0xb3, 0x19, 0x0f, 0x1c, 0x51, 0x46, 0x35, 0x40),
		},
		{
			name: "Invalid UUIDv4",
			in:   "I_AM_NOT_VALID",
			errFunc: func(t *testing.T, err error) {
				if err == nil {
					t.Fatalf("expected error, have <nil>")
				}
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeGenericInvalidArgument, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				default:
					t.Fatalf("unexpected error type: want %T, have %T", errors.Error{}, err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := IDFrom(tc.in)
			if err != nil {
				if tc.errFunc == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				tc.errFunc(t, err)
			}
			if want, have := tc.out.String(), id.String(); want != have {
				t.Fatalf("invalid id: want %s, have %s", want, have)
			}
		})
	}
}

func TestMustIDFrom(t *testing.T) {
	testCases := []struct {
		name      string
		in        string
		panicFunc func(*testing.T, interface{})
	}{
		{
			name: "Valid UUIDv4",
			in:   "276c8bbf-79ca-4ac2-b319-0f1c51463540",
		},
		{
			name: "Invalid UUIDv4",
			in:   "I_AM_NOT_VALID",
			panicFunc: func(t *testing.T, v interface{}) {
				if v == nil {
					t.Fatalf("expected panic value, have <nil>")
				}
				switch err := v.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeGenericInvalidArgument, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				default:
					t.Fatalf("unexpected error type: want %T, have %T", errors.Error{}, err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tc.panicFunc == nil {
						t.Fatalf("unexpected panic: %v", r)
					}
					tc.panicFunc(t, r)
				}
			}()
			MustIDFrom(tc.in)
		})
	}
}

func TestID_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name    string
		in      ID
		out     string
		errFunc func(*testing.T, error)
	}{
		{
			name: "Valid ID",
			in:   toID(0x27, 0x6c, 0x8b, 0xbf, 0x79, 0xca, 0x4a, 0xc2, 0xb3, 0x19, 0x0f, 0x1c, 0x51, 0x46, 0x35, 0x40),
			out:  "\"276c8bbf-79ca-4ac2-b319-0f1c51463540\"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.in.MarshalJSON()
			if err != nil {
				if tc.errFunc == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				tc.errFunc(t, err)
			}
			if want, have := tc.out, string(b); want != have {
				t.Fatalf("invalid JSON marshaled ID: want %v, have %v", want, have)
			}
		})
	}
}

func TestID_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name    string
		in      string
		out     ID
		errFunc func(*testing.T, error)
	}{
		{
			name: "Valid ID",
			in:   "\"276c8bbf-79ca-4ac2-b319-0f1c51463540\"",
			out:  toID(0x27, 0x6c, 0x8b, 0xbf, 0x79, 0xca, 0x4a, 0xc2, 0xb3, 0x19, 0x0f, 0x1c, 0x51, 0x46, 0x35, 0x40),
		},
		{
			name:    "Not quoted",
			in:      "276c8bbf-79ca-4ac2-b319-0f1c51463540",
			errFunc: assertInvalidArgumentError,
		},
		{
			name:    "Not a v4",
			in:      "ff26a2e2-8bbf-11e9-bc42-526af7764f64",
			errFunc: assertInvalidArgumentError,
		},
		{
			name:    "Nonsense",
			in:      "lorem impusLorem ipsum dolor sit amet",
			errFunc: assertInvalidArgumentError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var id ID
			err := id.UnmarshalJSON([]byte(tc.in))
			if err != nil {
				if tc.errFunc == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				tc.errFunc(t, err)
			}
			if want, have := tc.out, id; want != have {
				t.Fatalf("invalid JSON unmarshaled ID: want %v, have %v", want, have)
			}
		})
	}
}

func toID(data ...byte) ID {
	if len(data) != uuid.Size {
		panic("invalid uuid length")
	}
	var id ID
	copy(id[:], data)
	return id
}
