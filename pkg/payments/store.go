package payments

import (
	"fmt"
	"strings"

	"github.com/michaljemala/payments-sample/pkg/internal/errors"

	"github.com/lib/pq"

	"github.com/michaljemala/payments-sample/pkg/domain"
	"github.com/michaljemala/payments-sample/pkg/internal/store"
	"github.com/michaljemala/payments-sample/pkg/internal/store/sql"
)

func newPaymentStore() paymentStore {
	return &defaultPaymentStore{}
}

type defaultPaymentStore struct{}

func (s *defaultPaymentStore) Count(tx store.Tx, req domain.PaymentSearchRequest) (uint, error) {
	sqlTx := tx.(*sql.Tx)

	query := `SELECT count(*) FROM payment`

	conds, args := s.extractWhereClause(req)
	if len(conds) > 0 {
		where := strings.Join(conds, " AND ")
		query = fmt.Sprintf("%s WHERE %s", query, where)
	}

	var count uint
	err := sqlTx.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, sql.WrapSelectError(err, "unable to count payments")
	}

	return count, nil
}

func (s *defaultPaymentStore) Find(tx store.Tx, req domain.PaymentSearchRequest) ([]*domain.Payment, error) {
	sqlTx := tx.(*sql.Tx)

	query := `
	SELECT 
		id,
		amount_value,
		amount_currency,
		scheme_type,
		creditor_name,
		creditor_account_name,
		creditor_account_number,
		creditor_account_provider_code,
		creditor_account_provider_name,
		creditor_address_line1,
		creditor_address_line2,
		creditor_address_city,
		creditor_address_region,
		creditor_address_postal_code,
		creditor_address_country_code,
		debtor_name,
		debtor_account_name,
		debtor_account_number,
		debtor_account_provider_code,
		debtor_account_provider_name,
		debtor_address_line1,
		debtor_address_line2,
		debtor_address_city,
		debtor_address_region,
		debtor_address_postal_code,
		debtor_address_country_code
	FROM
		payment
	`

	conds, args := s.extractWhereClause(req)
	if len(conds) > 0 {
		where := strings.Join(conds, " AND ")
		query = fmt.Sprintf("%s WHERE %s", query, where)
	}

	if pag := req.SearchPagination; pag != nil {
		query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, pag.Limit(), pag.Offset())
	}

	rows, err := sqlTx.Query(query, args...)
	if err != nil {
		return nil, sql.WrapSelectError(err, "unable to select payments")
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		var payment domain.Payment
		err := rows.Scan(
			&payment.ID,
			&payment.Amount.Value,
			&payment.Amount.Currency,
			&payment.Scheme,
			&payment.Creditor.Name,
			&payment.Creditor.AccountName,
			&payment.Creditor.AccountNumber,
			&payment.Creditor.AccountProvider.Code,
			&payment.Creditor.AccountProvider.Name,
			&payment.Creditor.Address.Line1,
			&payment.Creditor.Address.Line2,
			&payment.Creditor.Address.City,
			&payment.Creditor.Address.Region,
			&payment.Creditor.Address.PostalCode,
			&payment.Creditor.Address.CountryCode,
			&payment.Debtor.Name,
			&payment.Debtor.AccountName,
			&payment.Debtor.AccountNumber,
			&payment.Debtor.AccountProvider.Code,
			&payment.Debtor.AccountProvider.Name,
			&payment.Debtor.Address.Line1,
			&payment.Debtor.Address.Line2,
			&payment.Debtor.Address.City,
			&payment.Debtor.Address.Region,
			&payment.Debtor.Address.PostalCode,
			&payment.Debtor.Address.CountryCode,
		)
		if err != nil {
			return nil, sql.WrapSelectError(err, "unable to scan payment")
		}
		payments = append(payments, &payment)
	}
	err = rows.Err()
	if err != nil {
		return nil, sql.WrapSelectError(err, "unable to select payments")
	}

	return payments, nil
}

func (s *defaultPaymentStore) Get(tx store.Tx, id domain.ID) (*domain.Payment, error) {
	sqlTx := tx.(*sql.Tx)

	query := `
	SELECT 
		id,
		amount_value,
		amount_currency,
		scheme_type,
		creditor_name,
		creditor_account_name,
		creditor_account_number,
		creditor_account_provider_code,
		creditor_account_provider_name,
		creditor_address_line1,
		creditor_address_line2,
		creditor_address_city,
		creditor_address_region,
		creditor_address_postal_code,
		creditor_address_country_code,
		debtor_name,
		debtor_account_name,
		debtor_account_number,
		debtor_account_provider_code,
		debtor_account_provider_name,
		debtor_address_line1,
		debtor_address_line2,
		debtor_address_city,
		debtor_address_region,
		debtor_address_postal_code,
		debtor_address_country_code
	FROM
		payment
	WHERE
		id = ?`

	var payment domain.Payment
	err := sqlTx.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.Amount.Value,
		&payment.Amount.Currency,
		&payment.Scheme,
		&payment.Creditor.Name,
		&payment.Creditor.AccountName,
		&payment.Creditor.AccountNumber,
		&payment.Creditor.AccountProvider.Code,
		&payment.Creditor.AccountProvider.Name,
		&payment.Creditor.Address.Line1,
		&payment.Creditor.Address.Line2,
		&payment.Creditor.Address.City,
		&payment.Creditor.Address.Region,
		&payment.Creditor.Address.PostalCode,
		&payment.Creditor.Address.CountryCode,
		&payment.Debtor.Name,
		&payment.Debtor.AccountName,
		&payment.Debtor.AccountNumber,
		&payment.Debtor.AccountProvider.Code,
		&payment.Debtor.AccountProvider.Name,
		&payment.Debtor.Address.Line1,
		&payment.Debtor.Address.Line2,
		&payment.Debtor.Address.City,
		&payment.Debtor.Address.Region,
		&payment.Debtor.Address.PostalCode,
		&payment.Debtor.Address.CountryCode,
	)
	if err != nil {
		return nil, sql.WrapSelectError(err, "unable to get payment")
	}

	return &payment, nil
}

func (s *defaultPaymentStore) Insert(tx store.Tx, payment *domain.Payment) error {
	sqlTx := tx.(*sql.Tx)

	query := `
	INSERT INTO payment ( 
		id,
		amount_value,
		amount_currency,
		scheme_type,
		creditor_name,
		creditor_account_name,
		creditor_account_number,
		creditor_account_provider_code,
		creditor_account_provider_name,
		creditor_address_line1,
		creditor_address_line2,
		creditor_address_city,
		creditor_address_region,
		creditor_address_postal_code,
		creditor_address_country_code,
		debtor_name,
		debtor_account_name,
		debtor_account_number,
		debtor_account_provider_code,
		debtor_account_provider_name,
		debtor_address_line1,
		debtor_address_line2,
		debtor_address_city,
		debtor_address_region,
		debtor_address_postal_code,
		debtor_address_country_code
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	_, err := sqlTx.Exec(query,
		payment.ID,
		payment.Amount.Value,
		payment.Amount.Currency,
		payment.Scheme,
		payment.Creditor.Name,
		payment.Creditor.AccountName,
		payment.Creditor.AccountNumber,
		payment.Creditor.AccountProvider.Code,
		payment.Creditor.AccountProvider.Name,
		payment.Creditor.Address.Line1,
		payment.Creditor.Address.Line2,
		payment.Creditor.Address.City,
		payment.Creditor.Address.Region,
		payment.Creditor.Address.PostalCode,
		payment.Creditor.Address.CountryCode,
		payment.Debtor.Name,
		payment.Debtor.AccountName,
		payment.Debtor.AccountNumber,
		payment.Debtor.AccountProvider.Code,
		payment.Debtor.AccountProvider.Name,
		payment.Debtor.Address.Line1,
		payment.Debtor.Address.Line2,
		payment.Debtor.Address.City,
		payment.Debtor.Address.Region,
		payment.Debtor.Address.PostalCode,
		payment.Debtor.Address.CountryCode,
	)

	return sql.WrapInsertError(err, "unable to insert payment")
}

func (s *defaultPaymentStore) Delete(tx store.Tx, id domain.ID) error {
	sqlTx := tx.(*sql.Tx)

	query := `DELETE FROM payment WHERE id = ?`

	_, err := sqlTx.Exec(query, id)

	return sql.WrapDeleteError(err, "unable to delete payment")
}

func (s *defaultPaymentStore) Update(tx store.Tx, payment *domain.Payment) error {
	sqlTx := tx.(*sql.Tx)

	query := `
	UPDATE payment
	SET
		amount_value = ?,
		amount_currency = ?,
		scheme_type = ?,
		creditor_name = ?,
		creditor_account_name = ?,
		creditor_account_number = ?,
		creditor_account_provider_code = ?,
		creditor_account_provider_name = ?,
		creditor_address_line1 = ?,
		creditor_address_line2 = ?,
		creditor_address_city = ?,
		creditor_address_region = ?,
		creditor_address_postal_code = ?,
		creditor_address_country_code = ?,
		debtor_name = ?,
		debtor_account_name = ?,
		debtor_account_number = ?,
		debtor_account_provider_code = ?,
		debtor_account_provider_name = ?,
		debtor_address_line1 = ?,
		debtor_address_line2 = ?,
		debtor_address_city = ?,
		debtor_address_region = ?,
		debtor_address_postal_code = ?,
		debtor_address_country_code = ?
	WHERE
		id = ?`

	_, err := sqlTx.Exec(query,
		payment.Amount.Value,
		payment.Amount.Currency,
		payment.Scheme,
		payment.Creditor.Name,
		payment.Creditor.AccountName,
		payment.Creditor.AccountNumber,
		payment.Creditor.AccountProvider.Code,
		payment.Creditor.AccountProvider.Name,
		payment.Creditor.Address.Line1,
		payment.Creditor.Address.Line2,
		payment.Creditor.Address.City,
		payment.Creditor.Address.Region,
		payment.Creditor.Address.PostalCode,
		payment.Creditor.Address.CountryCode,
		payment.Debtor.Name,
		payment.Debtor.AccountName,
		payment.Debtor.AccountNumber,
		payment.Debtor.AccountProvider.Code,
		payment.Debtor.AccountProvider.Name,
		payment.Debtor.Address.Line1,
		payment.Debtor.Address.Line2,
		payment.Debtor.Address.City,
		payment.Debtor.Address.Region,
		payment.Debtor.Address.PostalCode,
		payment.Debtor.Address.CountryCode,
		payment.ID,
	)

	return sql.WrapUpdateError(err, "unable to update payment")
}

func (s *defaultPaymentStore) extractWhereClause(req domain.PaymentSearchRequest) (conds []string, args []interface{}) {
	if list := req.IDs(); len(list) > 0 {
		conds = append(conds, "id = ANY (?)")
		args = append(args, pq.Array(list))
	}
	if list := req.CreditorAccountNumbers(); len(list) > 0 {
		conds = append(conds, "creditor_account_number = ANY (?)")
		args = append(args, pq.Array(list))
	}
	if list := req.DebtorAccountNumbers(); len(list) > 0 {
		conds = append(conds, "debtor_account_number = ANY (?)")
		args = append(args, pq.Array(list))
	}
	return conds, args
}

const (
	enumNameScheme   = domain.EnumName("SCHEME")
	enumNameCountry  = domain.EnumName("COUNTRY")
	enumNameCurrency = domain.EnumName("CURRENCY")
)

func newEnumStore() *defaultEnumStore {
	return &defaultEnumStore{
		enumMapping: map[domain.EnumName]string{
			enumNameScheme:   "enum_scheme",
			enumNameCountry:  "enum_country",
			enumNameCurrency: "enum_currency",
		},
	}
}

type defaultEnumStore struct {
	enumMapping map[domain.EnumName]string
}

func (s *defaultEnumStore) Exists(tx store.Tx, name domain.EnumName, code string) (bool, error) {
	sqlTx := tx.(*sql.Tx)

	tableName, ok := s.enumMapping[name]
	if !ok {
		return false, errors.Generic(errors.ErrCodeGenericInternal, "enum not found", "unable to select enumeration")
	}

	query := fmt.Sprintf(`SELECT count(*) FROM %s WHERE code = ?`, tableName)

	var count uint
	err := sqlTx.QueryRow(query, code).Scan(&count)
	if err != nil {
		return false, sql.WrapSelectError(err, "unable to query enum")
	}

	return count == 1, nil
}
