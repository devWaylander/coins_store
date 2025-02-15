package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/devWaylander/coins_store/pkg/models"
)

func Test_middleware_LoginWithPass(t *testing.T) {
	type fields struct {
		repo   Repository
		jwtKey string
	}
	type args struct {
		ctx context.Context
		qp  models.AuthQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.AuthDTO
		wantErr bool
	}{
		{
			name: "Successfully create new user and generate token",
			fields: fields{
				repo: &MockRepository{
					GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
						return &models.User{ID: 0}, nil
					},
					CreateUserTXFunc: func(ctx context.Context, username, passwordHash string) (int64, error) {
						return 1, nil
					},
				},
				jwtKey: "someKey",
			},
			args: args{
				ctx: context.Background(),
				qp:  models.AuthQuery{Username: "testuser", Password: "Test123@"},
			},
			want:    models.AuthDTO{Token: ""}, // Токен должен быть сгенерирован
			wantErr: false,
		},
		// {
		// 	name: "User already exists and password is correct",
		// 	fields: fields{
		// 		repo: &MockRepository{
		// 			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
		// 				return &models.User{ID: 1}, nil
		// 			},
		// 			GetUserPassHashByUsernameFunc: func(ctx context.Context, username string) (string, error) {
		// 				return "hashedpassword", nil
		// 			},
		// 		},
		// 		jwtKey: "someKey",
		// 	},
		// 	args: args{
		// 		ctx: context.Background(),
		// 		qp:  models.AuthQuery{Username: "testuser", Password: "Test123@"},
		// 	},
		// 	want:    models.AuthDTO{Token: ""}, // Токен должен быть сгенерирован
		// 	wantErr: false,
		// },
		{
			name: "Invalid password format for new user",
			fields: fields{
				repo: &MockRepository{
					GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
						return &models.User{ID: 0}, nil
					},
				},
				jwtKey: "someKey",
			},
			args: args{
				ctx: context.Background(),
				qp:  models.AuthQuery{Username: "testuser", Password: "short"}, // Неправильный формат пароля
			},
			want:    models.AuthDTO{}, // Токен НЕ должен быть сгенерирован
			wantErr: true,
		},
		{
			name: "Error on creating user",
			fields: fields{
				repo: &MockRepository{
					GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
						return &models.User{ID: 0}, nil
					},
					CreateUserTXFunc: func(ctx context.Context, username, passwordHash string) (int64, error) {
						return 0, errors.New("database error") // Ошибка при создании пользователя
					},
				},
				jwtKey: "someKey",
			},
			args: args{
				ctx: context.Background(),
				qp:  models.AuthQuery{Username: "testuser", Password: "Test123@"},
			},
			want:    models.AuthDTO{}, // Токен НЕ должен быть сгенерирован
			wantErr: true,
		},
		// {
		// 	name: "Error on generating JWT token",
		// 	fields: fields{
		// 		repo: &MockRepository{
		// 			GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
		// 				return &models.User{ID: 0}, nil
		// 			},
		// 			CreateUserTXFunc: func(ctx context.Context, username, passwordHash string) (int64, error) {
		// 				return 1, nil
		// 			},
		// 		},
		// 		jwtKey: "someKey",
		// 	},
		// 	args: args{
		// 		ctx: context.Background(),
		// 		qp:  models.AuthQuery{Username: "testuser", Password: "Test123@"},
		// 	},
		// 	want:    models.AuthDTO{}, // Токен НЕ должен быть сгенерирован
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &middleware{
				repo:   tt.fields.repo,
				jwtKey: tt.fields.jwtKey,
			}
			got, err := m.LoginWithPass(tt.args.ctx, tt.args.qp)
			if (err != nil) != tt.wantErr {
				t.Errorf("middleware.LoginWithPass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got.Token == "" {
				t.Errorf("middleware.LoginWithPass() = %v, want Token to be non-empty", got)
			}
		})
	}
}
