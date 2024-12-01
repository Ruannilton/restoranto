package repositories

import (
	"costumers-api/domain/domain_errors"
	"errors"
	"net"

	"github.com/lib/pq"
)

const (
	pg_connection_exception                              = "08000"
	pg_connection_does_not_exist                         = "08003"
	pg_connection_failure                                = "08006"
	pg_sqlclient_unable_to_establish_sqlconnection       = "08001"
	pg_sqlserver_rejected_establishment_of_sqlconnection = "08004"
)

func isConnectionFailure(code pq.ErrorCode) bool {
	return code == pg_connection_exception || code == pg_connection_does_not_exist || code == pg_connection_failure || code == pg_sqlclient_unable_to_establish_sqlconnection || code == pg_sqlserver_rejected_establishment_of_sqlconnection
}

func GetErrorCode(err error) string {
	var pqErr *pq.Error
	var netErr *net.OpError

	if errors.As(err, &pqErr) {
		if isConnectionFailure(pqErr.Code) {
			return domain_errors.CodeConnectionErr
		}
	}

	if errors.As(err, &netErr) {
		return domain_errors.CodeConnectionErr
	}

	return domain_errors.CodeInvalidOperation
}
