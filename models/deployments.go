package models

import (
	"fmt"
	"time"

	"github.com/asdutoit/go_backend_template/db"
)

type Deployment struct {
	ID           int64
	Created_at   time.Time
	Deployed_at  time.Time
	Organization string
	Product      string
	System_layer string
	Environment  string
	Version      string
	Url          string
	Status       string
	Order        int
}

func GetAllDeployments() ([]Deployment, error) {
	query := `SELECT * FROM releases`

	rows, err := db.DB.Query(query)

	fmt.Println("rows", rows)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deployments []Deployment

	for rows.Next() {
		var deployment Deployment
		err := rows.Scan(&deployment.ID, &deployment.Created_at, &deployment.Deployed_at, &deployment.Organization, &deployment.Product, &deployment.System_layer, &deployment.Environment, &deployment.Version, &deployment.Url, &deployment.Status, &deployment.Order)

		if err != nil {
			return nil, err
		}

		deployments = append(deployments, deployment)
	}

	return deployments, nil
}

func GetDeploymentByQuery(organization string, product string, systemLayer string, environment string) (*Deployment, error) {
	query := `SELECT * FROM releases WHERE organization = $1 AND product = $2 AND system_layer = $3 AND environment = $4`

	row := db.DB.QueryRow(query, organization, product, systemLayer, environment)

	var deployment Deployment
	err := row.Scan(&deployment.ID, &deployment.Created_at, &deployment.Deployed_at, &deployment.Organization, &deployment.Product, &deployment.System_layer, &deployment.Environment, &deployment.Version, &deployment.Url, &deployment.Status, &deployment.Order)

	if err != nil {
		return nil, err
	}

	return &deployment, nil
}
