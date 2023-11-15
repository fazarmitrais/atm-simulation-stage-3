package repository

import "github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"

type AccountCsvRepository interface {
	Import(path string) ([]*entity.Account, error)
}
