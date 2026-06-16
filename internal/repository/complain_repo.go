package repository

import (
	"absensi-sppg/internal/model"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ComplainRepository interface {
	GetAll(page int, limit int, status string, types string, parameter_flat string, parameter_zero string, ispu string, mutu_air string, validate string, accuweather_values string, per_page int, all_error string) ([]model.ComplainHandling, int, error)
}

type complainRepository struct {
	db *sqlx.DB
}

func NewComplainRepository(db *sqlx.DB) ComplainRepository {
	return &complainRepository{db: db}
}

func (r *complainRepository) GetAll(page int, limit int, status string, types string, parameter_flat string, parameter_zero string, ispu string, mutu_air string, validate string, accuweather_values string, per_page int, all_error string) ([]model.ComplainHandling, int, error) {
	if per_page > 0 {
		limit = per_page
	}
	offset := (page - 1) * limit

	var results []model.ComplainHandling
	if types == "AQMS" {
		types = "aqms"
	} else if types == "ONLIMO" {
		types = "WQMS"
	} else if types == "SPARING" {
		types = "sparing"
	}

	// Base query
	query := `
		select 
			tch.user_id, td.klh_name, td.type, tch.last_status, tch.last_received, 
			tch.parameter_flat, tch.parameter_zero, tch.parameter_anomaly, tch.ispu, 
			tch.status_mutu_air, tch.validate, tch.accuweather_values, tch.status, 
			tch.note, tch.created_at
		from ticket_complain_handling tch
		JOIN ticket_dataportal td ON td.user_id = tch.user_id
		WHERE DATE(tch.created_at) = CURRENT_DATE
	`

	// Dynamic params
	params := []interface{}{}

	// Add status filter if not "all"
	if status != "" && status != "all" {
		query += ` AND tch.last_status = ?`
		params = append(params, status)
	}

	// Add type filter if not "all"
	if types != "" && types != "all" {
		query += ` AND td.type = ?`
		params = append(params, types)
	}

	// Add parameter_flat filter if not "all"
	if all_error == "yes" {
		query += ` AND (`
		query += ` (tch.parameter_flat IS NOT NULL AND tch.parameter_flat <> '')`
		query += ` OR (tch.parameter_zero IS NOT NULL AND tch.parameter_zero <> '')`
		query += ` OR (tch.ispu IS NOT NULL AND tch.ispu <> '')`
		query += ` OR (tch.status_mutu_air IS NOT NULL AND tch.status_mutu_air <> '')`
		query += ` OR (tch.validate IS NOT NULL AND tch.validate <> '')`
		query += ` OR (tch.accuweather_values IS NOT NULL AND tch.accuweather_values <> '')`
		query += ` )`
		query += ` AND tch.last_status = 'online'`
	} else {
		if parameter_flat == "has" {
			query += ` OR tch.parameter_flat IS NOT NULL AND tch.parameter_flat <> ''`
		}

		// Add parameter_zero filter if not "all"
		if parameter_zero == "has" {
			query += ` AND tch.parameter_zero IS NOT NULL AND tch.parameter_zero <> ''`
		}

		// Add ispu filter if not "all"
		if ispu == "has" {
			query += ` AND tch.ispu IS NOT NULL AND tch.ispu <> ''`
		}

		// Add mutu_air filter if not "all"
		if mutu_air == "has" {
			query += ` AND tch.status_mutu_air IS NOT NULL AND tch.status_mutu_air <> ''`
		}

		// Add validate filter if not "all"
		if validate == "has" {
			query += ` AND tch.validate IS NOT NULL AND tch.validate <> ''`
		}

		// Add accuweather_values filter if not "all"
		if accuweather_values == "has" {
			query += ` AND tch.accuweather_values IS NOT NULL AND tch.accuweather_values <> ''`
		}
	}

	// Order + pagination
	query += `
  	  	GROUP BY tch.user_id
		ORDER BY td.type asc, td.klh_name asc
		LIMIT ? OFFSET ?
	`

	fmt.Println("query", query)
	params = append(params, limit, offset)

	fmt.Println("params", params)

	// Execute query
	err := r.db.Select(&results, query, params...)
	if err != nil {
		return nil, 0, err
	}

	// --- Count total ---
	countQuery := `
		SELECT COUNT(DISTINCT tch.user_id)
		FROM ticket_complain_handling tch
		JOIN ticket_dataportal td ON td.user_id = tch.user_id
		WHERE DATE(tch.created_at) = CURRENT_DATE
	`

	countParams := []interface{}{}

	if status != "" && status != "all" {
		countQuery += ` AND tch.last_status = ?`
		countParams = append(countParams, status)
	}

	if types != "" && types != "all" {
		countQuery += ` AND td.type = ?`
		countParams = append(countParams, types)
	}

	if all_error == "yes" {
		countQuery += ` AND (`
		countQuery += ` (tch.parameter_flat IS NOT NULL AND tch.parameter_flat <> '')`
		countQuery += ` OR (tch.parameter_zero IS NOT NULL AND tch.parameter_zero <> '')`
		countQuery += ` OR (tch.ispu IS NOT NULL AND tch.ispu <> '')`
		countQuery += ` OR (tch.status_mutu_air IS NOT NULL AND tch.status_mutu_air <> '')`
		countQuery += ` OR (tch.validate IS NOT NULL AND tch.validate <> '')`
		countQuery += ` OR (tch.accuweather_values IS NOT NULL AND tch.accuweather_values <> '')`
		countQuery += ` )`
		countQuery += ` AND tch.last_status = 'online'`
	} else {
		if parameter_flat == "has" {
			countQuery += ` AND tch.parameter_flat IS NOT NULL AND tch.parameter_flat <> ''`
		}

		if parameter_zero == "has" {
			countQuery += ` AND tch.parameter_zero IS NOT NULL AND tch.parameter_zero <> ''`
		}

		if ispu == "has" {
			countQuery += ` AND tch.ispu IS NOT NULL AND tch.ispu <> ''`
		}

		if mutu_air == "has" {
			countQuery += ` AND tch.status_mutu_air IS NOT NULL AND tch.status_mutu_air <> ''`
		}

		if validate == "has" {
			countQuery += ` AND tch.validate IS NOT NULL AND tch.validate <> ''`
		}

		if accuweather_values == "has" {
			countQuery += ` AND tch.accuweather_values IS NOT NULL AND tch.accuweather_values <> ''`
		}
	}

	var total int
	err = r.db.Get(&total, countQuery, countParams...)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
