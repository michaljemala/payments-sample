package domain

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"

	"github.com/michaljemala/payments-sample/pkg/internal/errors"
)

type ID uuid.UUID

func MustIDFrom(s string) ID {
	id, err := IDFrom(s)
	if err != nil {
		panic(err)
	}
	return id
}

func IDFrom(s string) (ID, error) {
	id, err := uuid.FromString(s)
	if err != nil {
		return ID{}, errors.Generic(
			errors.ErrCodeGenericInvalidArgument,
			"invalid format of ID",
			err.Error(),
		)
	}
	return ID(id), nil
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

func (id ID) IsNil() bool {
	return uuid.UUID(id) == uuid.Nil
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(id))
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var uuid uuid.UUID
	err := json.Unmarshal(data, &uuid)
	if err != nil {
		return errors.Generic(
			errors.ErrCodeGenericInvalidArgument,
			"invalid format of ID",
			err.Error(),
		)
	}
	*id = ID(uuid)
	return nil
}

func (id ID) Value() (driver.Value, error) {
	return uuid.UUID(id).Value()
}

func (id *ID) Scan(src interface{}) error {
	uid := new(uuid.UUID)
	err := uid.Scan(src)
	if err != nil {
		return err
	}
	*id = ID(*uid)
	return nil
}

type Decimal decimal.Decimal

func MustDecimalFrom(s string) Decimal {
	d, err := DecimalFrom(s)
	if err != nil {
		panic(err)
	}
	return d
}

func DecimalFrom(s string) (Decimal, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Decimal{}, errors.Generic(
			errors.ErrCodeGenericInvalidArgument,
			"invalid format of decimal",
			err.Error(),
		)
	}
	return Decimal(d), nil
}

func (d Decimal) String() string {
	return decimal.Decimal(d).String()
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(decimal.Decimal(d))
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	var dec decimal.Decimal
	err := json.Unmarshal(data, &dec)
	if err != nil {
		return errors.Generic(
			errors.ErrCodeGenericInvalidArgument,
			"invalid format of Decimal",
			err.Error(),
		)
	}
	*d = Decimal(dec)
	return nil
}

func (d Decimal) Value() (driver.Value, error) {
	return decimal.Decimal(d).Value()
}
func (d *Decimal) Scan(src interface{}) error {
	dec := new(decimal.Decimal)
	err := dec.Scan(src)
	if err != nil {
		return err
	}
	*d = Decimal(*dec)
	return nil
}

type Address struct {
	Line1       string  `json:"line1"`
	Line2       *string `json:"line2,omitempty"`
	City        string  `json:"city"`
	Region      *string `json:"region,omitempty"`
	PostalCode  string  `json:"postal_code"`
	CountryCode string  `json:"country_code"`
}

type Object interface {
	GetID() string
	SetID(s string) error
}

type BaseObject struct {
	ID ID `json:"-"`
}

func (o BaseObject) GetID() string {
	return o.ID.String()
}
func (o *BaseObject) SetID(s string) error {
	id, err := IDFrom(s)
	if err != nil {
		return err
	}
	o.ID = id
	return nil
}

type EnumName string
