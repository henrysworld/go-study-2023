package repository

import (
	"context"
	"github.com/henrysworld/go-study-2023/week02/webook/internal/domain"
	"github.com/henrysworld/go-study-2023/week02/webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// SELECT * FROM `users` WHERE `email`=?
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
	uM, err := r.dao.FindById(ctx, u.Id)
	if err != nil {
		return err
	}
	if u.NickName != "" {
		uM.NickName = u.NickName
	}
	if u.Birthday != "" {
		uM.Birthday = u.Birthday
	}
	if u.Bio != "" {
		uM.Bio = u.Bio
	}

	return r.dao.Update(ctx, uM)
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:  u.Id,
		Bio: u.Bio,
	}, nil

}
