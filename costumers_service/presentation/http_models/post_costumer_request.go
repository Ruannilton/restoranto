package httpmodels

import (
	"encoding/json"
	"fmt"
	"time"
)

type PostCostumerRequest struct {
	Name      string    `json:"name"`
	Cpf       string    `json:"cpf"`
	Birthdate time.Time `json:"birthdate"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
}

func (p *PostCostumerRequest) UnmarshalJSON(data []byte) error {
	// Create a separate struct to handle unmarshalling
	type Alias PostCostumerRequest
	aux := &struct {
		Birthdate string `json:"birthdate"` // Parse birthdate as a string initially
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	// Unmarshal into the auxiliary struct
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Parse the birthdate string using the "2006-01-02" layout
	if aux.Birthdate != "" {
		parsedDate, err := time.Parse("2006-01-02", aux.Birthdate)
		if err != nil {
			return fmt.Errorf("invalid birthdate format: %v", err)
		}
		p.Birthdate = parsedDate
	}

	return nil
}
