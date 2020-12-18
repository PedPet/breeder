package repository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/PedPet/breeder/model"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

const (
	// InsertBreeder is a sql statement to insert a breeder into the breeders table
	InsertBreeder = `
        INSERT INTO breeders(affix, short_affix, website)
        VALUES(?, ?, ?)
    `
	// InsertOwner is a sql statement to inert an owner into the owners table
	InsertOwner = `
        INSERT INTO owners(breeder_id, forename, surname, address, email)
        VALUES(?, ?, ?, ?, ?)
    `

	// AssociateBreederOwner is a sql statement to associate an owner to a breeder
	AssociateBreederOwner = `
        INSERT INTO breeder_owners(breeder_id, owner_id)
        VALUES(?, ?)
    `

	// UpdateBreeder is a sql statement to update a breeder in the breeders table
	UpdateBreeder = `
        UPDATE breeders SET affix = ?, short_affix = ?, website = ?
        WHERE id = ?
    `
	// UpdateOwner is a sql statement to update an owner in the owners table
	UpdateOwner = `
        UPDATE owners SET forename = ?, surname = ?, address = ?, email = ?
        WHERE id = ?
    `
	// SelectBreeder is a sql statement to get a breeder from the breeders table
	SelectBreeder = `
        SELECT * FROM breeders WHERE id = ?
    `

	// SelectOwner is a sql statement to get an owner from the owners table
	SelectOwner = `
        SELECT * FROM owners WHERE id = ?
    `

	// SelectOwnerByBreeder is a sql statement to get an owner that belongs to a specific breeder
	SelectOwnerByBreeder = `
        SELECT * FROM owners WHERE breeder_id = ?
    `

	// DeleteBreeder is a sql statement to delete a breeder from the breeders table
	DeleteBreeder = `
        UPDATE breeders SET active = 0 WHERE id = ?
    `

	// DeleteOwner is a sql statement to delete an owner from the owners table
	DeleteOwner = `
        UPDATE owners SET active = 0 WHERE id = ?
    `

	// UnassociateBreederOwner is a sql statement to unassociate a breededr and an owner
	UnassociateBreederOwner = `
        DELETE FROM breeder_owners WHERE breeder_id = ? AND owner_id = ?
    `
)

var errRepo = errors.New("Unable to handle Repo Request")

// Breeder interface to define bredder repo
type Breeder interface {
	CreateBreeder(ctx context.Context, breeder *model.Breeder) error
	CreateOwner(ctx context.Context, breederID int, owner *model.Owner) error
	AssociateBreederOwner(ctx context.Context, breederID, ownerID int) error
	UpdateBreeder(ctx context.Context, breeder *model.Breeder) error
	UpdateOwner(ctx context.Context, owner *model.Owner) error
	DeleteBreeder(ctx context.Context, id int) error
	DeleteOwner(ctx context.Context, id int) error
	CheckBreederNeedsUpdating(ctx context.Context, breeder model.Breeder) (bool, error)
	CheckOwnerNeedsUpdating(ctx context.Context, owner model.Owner) (bool, error)
	UnassociateBreederOwner(ctx, breederID, ownerID int) error
	GetBreeder(ctx context.Context, id int) (model.Breeder, error)
	GetOwner(ctx context.Context, id int) (model.Owner, error)
	IsBreederOwnerAssociated(ctx context.Context, breederID, ownerID int) (bool, error)
}

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepo(db *sql.DB, logger log.Logger) Breeder {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (r repo) CreateBreeder(ctx context.Context, breeder *model.Breeder) error {
	logger := log.With(r.logger, "method", "CreateBreeder")
	stmt, err := r.db.PrepareContext(ctx, InsertBreeder)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare insert breeder statement")
	}
	defer stmt.Close()

	if breeder.Affix == "" || len(breeder.Owners) == 0 || breeder.ShortAffix == "" || breeder.Website == "" {
		return errRepo
	}

	result, err := stmt.Exec(breeder.Affix, breeder.ShortAffix, breeder.Website, breeder)
	if err != nil {
		return errors.Wrap(err, "Failed to execute prepared insert statement")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "Failed to get last insert id")
	}

	breeder.ID = int(id)

	for idx := range breeder.Owners {
		r.CreateOwner(ctx, breeder.ID, &breeder.Owners[idx])
	}

	logger.Log("Created breeder", breeder.ID)

	return nil
}

func (r repo) CreateOwner(ctx context.Context, breederID int, owner *model.Owner) error {
	logger := log.With(r.logger, "method", "CreateOwner")
	stmt, err := r.db.PrepareContext(ctx, InsertOwner)
	if err != nil {
		return errors.Wrap(err, "Failed to prepate insert owner statement")
	}
	defer stmt.Close()

	if owner.Forename == "" || owner.Surname == "" || owner.Address == "" || owner.Email == "" {
		return errRepo
	}

	result, err := stmt.Exec(
		breederID,
		owner.Forename,
		owner.Surname,
		owner.Address,
		owner.Email,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to execute prepated statement")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "Failed to get last insert id")
	}

	owner.ID = int(id)

	logger.Log("Created owner", owner.ID)
	return nil
}

func (r repo) AssociateBreederOwner(ctx context.Context, breederID, ownerID int) error {
	logger := log.With(r.logger, "method", "AssociateBreederOwner")
	stmt, err := r.db.PrepareContext(ctx, AssociateBreederOwner)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare associate breeder with owner")
	}
	defer stmt.Close()

	_, err = stmt.Exec(breederID, ownerID)
	if err != nil {
		return errors.Wrap(err, "Failed to execute associate breeder with owner statement")
	}

	logger.Log(fmt.Sprintf("Associate breeder: %d with owner: %d", breederID, ownerID))
	return nil
}

func (r repo) UpdateBreeder(ctx context.Context, breeder *model.Breeder) error {
	logger := log.With(r.logger, "method", "UpdateBreeder")
	stmt, err := r.db.PrepareContext(ctx, UpdateBreeder)
	if err != nil {
		return errors.Wrap(err, "Failed to prepate update breeder statement")
	}
	defer stmt.Close()

	_, err = stmt.Exec(breeder.Affix, breeder.ShortAffix, breeder.Website, breeder.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to execute prepared statement")
	}

	for idx := range breeder.Owners {
		r.UpdateOwner(ctx, &breeder.Owners[idx])
	}

	logger.Log("Updated breeder", breeder.ID)
	return nil
}

func (r repo) UpdateOwner(ctx context.Context, owner *model.Owner) error {
	logger := log.With(r.logger, "method", "UpdateOwner")

	// TODO: add a check if owner needs updating method

	stmt, err := r.db.PrepareContext(ctx, UpdateOwner)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare update owner statement")
	}
	defer stmt.Close()

	_, err = stmt.Exec(owner.Forename, owner.Surname, owner.Address, owner.Email, owner.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to execute prepared update statement")
	}

	logger.Log("Updated owner", owner.ID)
	return nil
}

func (r repo) DeleteBreeder(ctx context.Context, id int) error {
	logger := log.With(r.logger, "method", "DeleteBreeder")
	stmt, err := r.db.PrepareContext(ctx, DeleteBreeder)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare delete breeder statement")
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return errors.Wrap(err, "Failed to execute delete breeder statement")
	}

	logger.Log("Deleted breeder", id)
	return nil
}

func (r repo) DeleteOwner(ctx context.Context, id int) error {
	logger := log.With(r.logger, "method", "DeleteOwner")
	stmt, err := r.db.PrepareContext(ctx, DeleteOwner)
	if err != nil {
		return errors.Wrap(err, "Failed to prepate delete owner statement")
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return errors.Wrap(err, "Failed to execute delete owner")
	}

	logger.Log("Deleted owner", id)
	return nil
}

func (r repo) CheckBreederNeedsUpdating(ctx context.Context, breeder model.Breeder) (bool, error) {
	logger := log.With(r.logger, "method", "CheckBreederNeedsUpdating")
	dbBreeder, err := r.GetBreeder(ctx, breeder.ID)
	if err != nil {
		return false, errors.Wrap(err, "Failed to get breeder")
	}

	equal := reflect.DeepEqual(breeder, dbBreeder)

	logger.Log("Needs updating", equal)
	return equal, nil

}

func (r repo) CheckOwnerNeedsUpdating(ctx context.Context, owner model.Owner) (bool, error) {
	logger := log.With(r.logger, "method", "CheckOwnerNeedsUpdating")

}

// TODO: Add UnassociateBreederOwner method here

func (r repo) GetBreeder(ctx context.Context, id int) (model.Breeder, error) {
	logger := log.With(r.logger, "method", "GetBreeder")
	breeder := model.Breeder{}

	stmt, err := r.db.PrepareContext(ctx, SelectBreeder)
	if err != nil {
		return breeder, errors.Wrap(err, "Failed to prepate get breeder statement")
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return breeder, errors.Wrap(err, "Failed to execute get breeder statement")
	}

	rows.Scan(&breeder)

	logger.Log("Got breeder", breeder)
	return breeder, nil
}
