package user

import (
	"context"

	"github.com/lib/pq"
	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func (u *userRepository) registerUserSQL(ctx context.Context, tx *sqlx.Tx, user *entity.User) (*sqlx.Tx, *entity.User, error) {
	query, _ := u.queryLoader.Get("RegisterUser")
	row := tx.QueryRowContext(ctx, query, user.AutdID, user.Email, user.FirstName, user.LastName).Scan(&user.ID)
	if err := row; err != nil {
		zerolog.Ctx(ctx).Error().Str("id", user.ID).Msg("register_user_sql")
		return tx, nil, x.Wrap(err, "register_user_sql")
	}

	query, _ = u.queryLoader.Get("RegisterUserProfile")
	result := tx.MustExecContext(ctx, query, user.ID)
	rows, _ := result.RowsAffected()
	if rows == 0 {
		zerolog.Ctx(ctx).Error().Str("id", user.ID).Msg("register_user_profile_sql")
		return tx, nil, x.NewWithCode(x.CodeSQLCreate, "register_user_profile_sql")
	}

	return tx, user, nil
}

func (u *userRepository) getUserByAuthIDSQL(ctx context.Context, userAuthID string) (*entity.User, error) {
	var user entity.User

	query, _ := u.queryLoader.Get("GetUserByAuthID")
	err := u.db0.QueryRowxContext(ctx, query, userAuthID).StructScan(&user)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_by_auth_id_sql")
		return &user, x.WrapWithCode(err, x.CodeSQLRowScan, "get_user_by_auth_id_sql")
	}

	return &user, nil
}

func (u *userRepository) getUserByIDSQL(ctx context.Context, userID string) (*entity.User, error) {
	var user entity.User

	query, _ := u.queryLoader.Get("GetUserByID")
	err := u.db0.QueryRowxContext(ctx, query, userID).StructScan(&user)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_by_id_sql")
		return &user, x.WrapWithCode(err, x.CodeSQLRowScan, "get_user_by_id_sql")
	}

	return &user, nil
}

func (u *userRepository) getUserAddressesSQL(ctx context.Context, userID string) ([]*entity.UserAddress, error) {
	query, _ := u.queryLoader.Get("GetUserAddresses")
	rows, err := u.db0.QueryContext(ctx, query, userID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_addresses_sql")
		return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_user_addresses_sql")
	}
	defer rows.Close()

	addresses := make([]*entity.UserAddress, 0, 10) // adjust capacity as needed
	for rows.Next() {
		var address entity.UserAddress
		if err := rows.Scan(
			&address.ID,
			&address.AddressType,
			&address.StreetAddress,
			&address.City,
			&address.State, // sql.NullString handles NULL
			&address.PostalCode,
			&address.Country,
			&address.IsDefault,
		); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_addresses_sql row scan")
			return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_user_addresses_sql")
		}

		addresses = append(addresses, &address)
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_addresses_sql rows iteration error")
		return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_user_addresses_sql")
	}

	return addresses, nil
}

func (u *userRepository) createUserAddressSQL(ctx context.Context, userID string, req dto.CreateUserAddress) (string, error) {
	var addressID string

	query, _ := u.queryLoader.Get("CreateUserAddress")
	err := u.db0.QueryRowContext(ctx, query, userID, req.AddressType, req.StreetAddress, req.City, req.State, req.PostalCode, req.Country, req.IsDefault).Scan(&addressID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("create_user_address_sql")

		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code {
			case "23503": // foreign key violation
				return "", x.WrapWithCode(err, x.CodeSQLEmptyRow, "user_not_found")
			case "23514": // check constraint violation (address_type)
				return "", x.WrapWithCode(err, x.CodeSQLQueryBuild, "invalid_address_type")
			}
		}

		return "", x.WrapWithCode(err, x.CodeSQLCreate, "create_user_address_sql")
	}

	return addressID, nil
}
