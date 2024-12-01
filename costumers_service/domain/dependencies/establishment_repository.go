package dependencies

import "costumers-api/domain/models"

type IEstablishmentRepository interface {
	Select(id int) (*models.Establishment, error)
	Insert(establishment models.Establishment) (int, error)
	Update(establishment models.Establishment) error
}
