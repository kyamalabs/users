package db

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
)

var (
	UserProfileAlreadyExistsError = errors.New("user profile already exists")
	GamerTagAlreadyInUseError     = errors.New("gamer tag already in use")
	UserAlreadyReferredError      = errors.New("user already referred")
	ReferrerDoesNotExistError     = errors.New("referrer does not exist")
	SelfReferralError             = errors.New("self referrals not permitted")
)

type CreateProfileTxParams struct {
	CreateProfileParams
	Referrer    string
	AfterCreate func() error
}

type CreateProfileTxResult struct {
	Profile  Profile
	Referral Referral
}

func (store *SQLStore) CreateProfileTx(ctx context.Context, params CreateProfileTxParams) (CreateProfileTxResult, error) {
	var result CreateProfileTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Profile, err = q.CreateProfile(ctx, params.CreateProfileParams)
		if err != nil {
			log.Error().Err(err).Msg("could not create user profile in db")
			dbError := ParseError(err)

			if dbError.Code == UniqueViolationCode {
				switch dbError.ConstraintName {
				case "profiles_pkey":
					return UserProfileAlreadyExistsError
				case "profiles_gamer_tag_key":
					return GamerTagAlreadyInUseError
				}
			}

			return err
		}

		createReferralParams := CreateReferralParams{
			Referrer: params.Referrer,
			Referee:  params.CreateProfileParams.WalletAddress,
		}

		result.Referral, err = createProfileReferral(ctx, q, createReferralParams)
		if err != nil {
			return err
		}

		return params.AfterCreate()
	})

	return result, err
}

func createProfileReferral(ctx context.Context, q *Queries, createReferralParams CreateReferralParams) (Referral, error) {
	if createReferralParams.Referrer == "" {
		return Referral{}, nil
	}

	if createReferralParams.Referrer == createReferralParams.Referee {
		return Referral{}, SelfReferralError
	}

	_, err := q.GetProfile(ctx, createReferralParams.Referrer)
	if err != nil {
		log.Error().Err(err).Msg("could not get referrer profile in db")

		if err == RecordNotFoundError {
			return Referral{}, ReferrerDoesNotExistError
		}

		return Referral{}, err
	}

	referral, err := q.CreateReferral(ctx, createReferralParams)
	if err != nil {
		log.Error().Err(err).Msg("could not create user referral in db")
		dbError := ParseError(err)

		if dbError.Code == UniqueViolationCode {
			if dbError.ConstraintName == "referrals_referee_key" {
				return Referral{}, UserAlreadyReferredError
			}
		}

		return Referral{}, err
	}

	log.Info().Interface("referral", referral).Msg("successfully created user referral")

	return referral, nil
}
