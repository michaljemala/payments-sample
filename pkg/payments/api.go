package payments

import (
	"fmt"
	"log"
	"net/http"

	"github.com/manyminds/api2go"

	"github.com/michaljemala/payments-sample/pkg/domain"
	"github.com/michaljemala/payments-sample/pkg/internal/store/sql"
)

type Config struct {
	Prefix string
	Driver string
	DSN    string
	Logger *log.Logger
}

type API struct {
	config  Config
	db      *sql.DB
	handler http.Handler
}

const jsonApiContentType = "application/vnd.api+json"

func NewAPI(c Config) (*API, error) {
	db, err := sql.Connect(sql.Config{
		Driver: c.Driver,
		DSN:    c.DSN,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to store: %v", err)
	}

	txManager := sql.NewTxManager(db)
	paymentStore := newPaymentStore()
	enumStore := newEnumStore()
	service := newPaymentService(txManager, paymentStore, enumStore, c.Logger)

	api := newAPI(c, service)
	api.db = db

	return api, nil
}

func newAPI(c Config, service paymentService) *API {
	api := api2go.NewAPI(c.Prefix)
	api.ContentType = jsonApiContentType
	api.AddResource(&domain.Payment{}, newResource(service))
	return &API{config: c, handler: api.Handler()}
}

func (api *API) Prefix() string {
	return api.config.Prefix
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.handler.ServeHTTP(w, r)
}

func (api *API) Close() error {
	if api.db != nil {
		err := api.db.Close()
		if err != nil {
			return fmt.Errorf("unable to close store: %v", err)
		}
	}
	return nil
}
