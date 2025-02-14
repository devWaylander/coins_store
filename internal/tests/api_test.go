package tests

// import (
// 	"context"
// 	"errors"
// 	"reflect"
// 	"testing"

// 	"github.com/devWaylander/coins_store/internal/service"
// 	"github.com/devWaylander/coins_store/pkg/models"
// )

// func Test_service_GetUserInfo(t *testing.T) {
// 	type fields struct {
// 		repo service.Repository
// 	}
// 	type args struct {
// 		ctx context.Context
// 		qp  models.InfoQuery
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    models.InfoDTO
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful user fetch",
// 			fields: fields{
// 				repo: &MockRepository{
// 					GetUserInfoFunc: func(ctx context.Context, qp models.InfoQuery) (models.InfoDTO, error) {
// 						return models.InfoDTO{ID: 1, Name: "John Doe"}, nil
// 					},
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				qp:  models.InfoQuery{UserID: 1},
// 			},
// 			want:    models.InfoDTO{ID: 1, Name: "John Doe"},
// 			wantErr: false,
// 		},
// 		{
// 			name: "user not found",
// 			fields: fields{
// 				repo: &MockRepository{
// 					GetUserInfoFunc: func(ctx context.Context, qp models.InfoQuery) (models.InfoDTO, error) {
// 						return models.InfoDTO{}, errors.New("user not found")
// 					},
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				qp:  models.InfoQuery{UserID: 999},
// 			},
// 			want:    models.InfoDTO{},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service.New(tt.fields.repo)
// 			got, err := s.GetUserInfo(tt.args.ctx, tt.args.qp)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("service.GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("service.GetUserInfo() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
