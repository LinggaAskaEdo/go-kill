package user

import (
	"context"
	"database/sql"
	"strconv"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

func (u *userRepository) createUserSQL(ctx context.Context, req *userpb.CreateUserRequest) (string, error) {
	var userID string

	query, _ := u.queryLoader.Get("RegisterUser")
	err := u.db0.QueryRowContext(ctx, query, req.AuthId, req.Email, req.FirstName, req.LastName).Scan(&userID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("create_user_sql")
		return "", x.WrapWithCode(err, x.CodeSQLCreate, "create_user_sql")
	}

	return userID, nil
}

func (u *userRepository) registerUserSQL(ctx context.Context, tx *sqlx.Tx, user *entity.User) (*sqlx.Tx, *entity.User, error) {
	query, ok := u.queryLoader.Get("RegisterUser")
	if !ok {
		err := x.New("query_loader_register_user")
		zerolog.Ctx(ctx).Error().Err(err).Msg("query_loader_register_user")
		return tx, nil, err
	}
	err := tx.QueryRowContext(ctx, query, user.AuthID, user.Email, user.FirstName, user.LastName).Scan(&user.ID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Str("id", user.ID).Msg("register_user_sql")
		return tx, nil, x.Wrap(err, "register_user_sql")
	}

	return tx, user, nil
}

func (u *userRepository) registerUserProfileSQL(ctx context.Context, tx *sqlx.Tx, userID string) (*sqlx.Tx, error) {
	query, ok := u.queryLoader.Get("RegisterUserProfile")
	if !ok {
		err := x.New("query_loader_register_user_profile")
		zerolog.Ctx(ctx).Error().Err(err).Msg("query_loader_register_user_profile")
		return tx, err
	}
	result, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("userID", userID).Msg("register_user_profile_sql")
		return tx, x.WrapWithCode(err, x.CodeSQLCreate, "register_user_profile_sql")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("userID", userID).Msg("register_user_profile_sql")
		return tx, x.WrapWithCode(err, x.CodeSQLCreate, "register_user_profile_sql")
	}

	if rows == 0 {
		zerolog.Ctx(ctx).Error().Str("id", userID).Msg("register_user_profile_sql")
		return tx, x.NewWithCode(x.CodeSQLCannotRetrieveAffectedRows, "register_user_profile_sql")
	}

	return tx, nil
}

func (u *userRepository) getUserByAuthIDSQL(ctx context.Context, userAuthID string) (*entity.User, error) {
	var user entity.User

	query, _ := u.queryLoader.Get("GetUserByAuthID")
	err := u.db0.QueryRowxContext(ctx, query, userAuthID).StructScan(&user)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_by_auth_id_sql")
		if err == sql.ErrNoRows {
			return &user, x.WrapWithCode(err, x.CodeSQLRecordDoesNotExist, "get_user_by_auth_id_sql")
		}
		return &user, x.WrapWithCode(err, x.CodeSQLRowScan, "get_user_by_auth_id_sql")
	}

	return &user, nil
}

func (u *userRepository) getUserByIDSQL(ctx context.Context, userID string) (*entity.User, error) {
	var user entity.User

	query, ok := u.queryLoader.Get("GetUserByID")
	if !ok {
		err := x.New("query_loader_get_user_by_id")
		zerolog.Ctx(ctx).Error().Err(err).Msg("query_loader_get_user_by_id")
		return &user, err
	}
	err := u.db0.QueryRowxContext(ctx, query, userID).StructScan(&user)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_by_id_sql")
		if err == sql.ErrNoRows {
			return &user, x.WrapWithCode(err, x.CodeSQLRecordDoesNotExist, "get_user_by_id_sql")
		}
		return &user, x.WrapWithCode(err, x.CodeSQLRowScan, "get_user_by_id_sql")
	}

	return &user, nil
}

func (u *userRepository) getUserAddressByIDSQL(ctx context.Context, req *userpb.GetAddressRequest) (*entity.UserAddress, error) {
	var userAddress entity.UserAddress

	query, ok := u.queryLoader.Get("GetUserAddressByID")
	if !ok {
		err := x.New("query_loader_get_user_address_by_id")
		zerolog.Ctx(ctx).Error().Err(err).Msg("query_loader_get_user_address_by_id")
		return &userAddress, err
	}
	err := u.db0.QueryRowxContext(ctx, query, req.AddressId, req.UserId).StructScan(&userAddress)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_address_by_id_sql")

		if err == sql.ErrNoRows {
			return &userAddress, x.WrapWithCode(err, x.CodeSQLRecordDoesNotExist, "get_user_address_by_id_sql")
		}

		return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_user_address_by_id_sql")
	}

	return &userAddress, nil
}

func (u *userRepository) getUserAddressesByUserIDSQL(ctx context.Context, userID string, page string, limit string) ([]*entity.UserAddress, int64, error) {
	if userID == "" {
		err := x.New("empty_user_id")
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_addresses_by_user_id_sql")
		return nil, 0, err
	}

	pageInt := 1
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageInt = p
	}

	limitInt := 20
	if l, err := strconv.Atoi(limit); err == nil && l > 0 {
		limitInt = l
	}
	if limitInt > 100 {
		limitInt = 100
	}

	offset := (pageInt - 1) * limitInt

	countQuery, ok := u.queryLoader.Get("CountUserAddresses")
	if !ok {
		err := x.New("query_loader_count_user_addresses")
		zerolog.Ctx(ctx).Error().Err(err).Msg("query_loader_count_user_addresses")
		return nil, 0, err
	}
	var total int64
	if err := u.db0.QueryRowxContext(ctx, countQuery, userID).Scan(&total); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("count_user_addresses_sql")
		return nil, 0, x.WrapWithCode(err, x.CodeSQLRowScan, "count_user_addresses_sql")
	}

	query, ok := u.queryLoader.Get("GetUserAddresses")
	if !ok {
		err := x.New("query_loader_get_user_addresses")
		zerolog.Ctx(ctx).Error().Err(err).Msg("query_loader_get_user_addresses")
		return nil, 0, err
	}
	rows, err := u.db0.QueryContext(ctx, query, userID, limitInt, offset)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_addresses_sql")
		return nil, 0, x.WrapWithCode(err, x.CodeSQLRowScan, "get_user_addresses_sql")
	}
	defer rows.Close()

	addresses := make([]*entity.UserAddress, 0)
	for rows.Next() {
		var address entity.UserAddress
		if err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.AddressType,
			&address.StreetAddress,
			&address.City,
			&address.State,
			&address.PostalCode,
			&address.Country,
			&address.IsDefault,
		); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_addresses_sql_row_scan")
			return nil, 0, x.WrapWithCode(err, x.CodeSQLRowScan, "get_user_addresses_sql")
		}

		addresses = append(addresses, &address)
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_addresses_sql_rows")
		return nil, 0, x.WrapWithCode(err, x.CodeSQLRead, "get_user_addresses_sql")
	}

	return addresses, total, nil
}

func (u *userRepository) createUserAddressSQL(ctx context.Context, userID string, req dto.CreateUserAddress) (string, error) {
	var addressID string

	query, _ := u.queryLoader.Get("CreateUserAddress")
	err := u.db0.QueryRowContext(ctx, query, userID, req.AddressType, req.StreetAddress, req.City, req.State, req.PostalCode, req.Country, req.IsDefault).Scan(&addressID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("create_user_address_sql")

		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code {
			case "23503":
				return "", x.WrapWithCode(err, x.CodeSQLForeignKeyMissing, "user_not_found")
			case "23514":
				return "", x.WrapWithCode(err, x.CodeSQLQueryBuild, "invalid_address_type")
			}
		}

		return "", x.WrapWithCode(err, x.CodeSQLCreate, "create_user_address_sql")
	}

	return addressID, nil
}
