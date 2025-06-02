package usecase

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/modules/member/interface_adapter/gateway/repository"
	"module-clean/internal/modules/member/interface_adapter/inputmodel"
	"module-clean/internal/modules/member/usecase/mock"
	"module-clean/internal/modules/member/usecase/port/output"
	"module-clean/internal/shared/common/pagination"
	"reflect"
	"testing"
	"time"
)

func TestMemberUseCase_DeleteMember(t *testing.T) {
	type fields struct {
		MemberRepo output.MemberRepository
	}
	type args struct {
		ctx context.Context
		id  int
	}
	ctrl, ctx, testTime := repoHelper(t)
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *entity.Member
		repoSetup func(*mock.MockMemberRepository)
		wantErr   error
	}{
		{
			name: "normal test",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				id:  0,
			},
			want: &entity.Member{
				ID:        0,
				Name:      "gg",
				Email:     "gg@gmail.com",
				Password:  "",
				CreatedAt: testTime,
			},
			repoSetup: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:        0,
						Name:      "gg",
						Email:     "gg@gmail.com",
						Password:  "",
						CreatedAt: testTime,
					}, nil),
					r.EXPECT().Delete(ctx, gomock.Any()).Return(nil),
				)
			},
			wantErr: nil,
		}, {
			name: "no member found",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				id:  0,
			},
			want: nil,
			repoSetup: func(r *mock.MockMemberRepository) {
				r.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, repository.ErrGatewayMemberNotFound)
			},
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberNotFound),
		}, {
			name: "got member but delete error",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				id:  0,
			},
			want: nil,
			repoSetup: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{}, nil),
					r.EXPECT().Delete(ctx, gomock.Any()).Return(repository.ErrGatewayMemberDBError),
				)
			},
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberDBError),
		}, {
			name: "delete no affect",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				id:  0,
			},
			want: nil,
			repoSetup: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{}, nil),
					r.EXPECT().Delete(ctx, gomock.Any()).Return(repository.ErrGatewayMemberDeleteFailed),
				)
			},
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberDeleteFailed),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, ok := tt.fields.MemberRepo.(*mock.MockMemberRepository)
			if !ok {
				t.Fatalf("expected *mock.MockMemberRepository, got %T", tt.fields.MemberRepo)
			}
			m := &MemberUseCase{
				MemberRepo: mockRepo,
			}
			tt.repoSetup(mockRepo)
			got, err := m.DeleteMember(tt.args.ctx, tt.args.id)
			t.Logf("got = %v, want %v", got, tt.want)
			t.Logf("err = %v, wantErr %v", err, tt.wantErr)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteMember() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberUseCase_GetMemberByEmail(t *testing.T) {
	type fields struct {
		MemberRepo output.MemberRepository
	}
	type args struct {
		ctx   context.Context
		email string
	}
	ctrl, ctx, testTime := repoHelper(t)
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *entity.Member
		repoSetup func(*mock.MockMemberRepository)
		wantErr   error
	}{
		{
			name: "normal test",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx:   ctx,
				email: "",
			},
			want: &entity.Member{
				ID:        0,
				Name:      "gg",
				Email:     "gg@gmail.com",
				Password:  "",
				CreatedAt: testTime,
			},
			repoSetup: func(r *mock.MockMemberRepository) {
				r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(&entity.Member{
					ID:        0,
					Name:      "gg",
					Email:     "gg@gmail.com",
					Password:  "",
					CreatedAt: testTime,
				}, nil)
			},
			wantErr: nil,
		}, {
			name: "no member found",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx:   ctx,
				email: "",
			},
			want: nil,
			repoSetup: func(r *mock.MockMemberRepository) {
				r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, repository.ErrGatewayMemberNotFound)
			},
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberNotFound),
		}, {
			name: "db error",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx:   ctx,
				email: "",
			},
			want: nil,
			repoSetup: func(r *mock.MockMemberRepository) {
				r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, repository.ErrGatewayMemberDBError)
			},
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberDBError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, ok := tt.fields.MemberRepo.(*mock.MockMemberRepository)
			if !ok {
				t.Fatalf("expected *mock.MockMemberRepository, got %T", tt.fields.MemberRepo)
			}
			m := &MemberUseCase{
				MemberRepo: mockRepo,
			}
			tt.repoSetup(mockRepo)
			got, err := m.GetMemberByEmail(tt.args.ctx, tt.args.email)
			t.Logf("got = %v, want %v", got, tt.want)
			t.Logf("err = %v, wantErr %v", err, tt.wantErr)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetMemberByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMemberByEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberUseCase_GetMemberByID(t *testing.T) {
	type fields struct {
		MemberRepo output.MemberRepository
	}
	type args struct {
		ctx context.Context
		id  int
	}
	ctrl, ctx, testTime := repoHelper(t)
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *entity.Member
		repoSetup func(*mock.MockMemberRepository)
		wantErr   error
	}{
		{
			name: "normal test",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				id:  1,
			},
			want: &entity.Member{
				ID:        1,
				Name:      "gg",
				Email:     "gg@gmail.com",
				Password:  "",
				CreatedAt: testTime,
			},
			repoSetup: func(r *mock.MockMemberRepository) {
				r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
					ID:        1,
					Name:      "gg",
					Email:     "gg@gmail.com",
					Password:  "",
					CreatedAt: testTime,
				}, nil)
			},
			wantErr: nil,
		}, {
			name: "no member found",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				id:  0,
			},
			want: nil,
			repoSetup: func(r *mock.MockMemberRepository) {
				r.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, repository.ErrGatewayMemberNotFound)
			},
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberNotFound),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, ok := tt.fields.MemberRepo.(*mock.MockMemberRepository)
			if !ok {
				t.Fatalf("expected *mock.MockMemberRepository, got %T", tt.fields.MemberRepo)
			}
			m := &MemberUseCase{
				MemberRepo: mockRepo,
			}
			tt.repoSetup(mockRepo)
			got, err := m.GetMemberByID(tt.args.ctx, tt.args.id)
			t.Logf("got = %v, want %v", got, tt.want)
			t.Logf("err = %v, wantErr %v", err, tt.wantErr)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetMemberByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMemberByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberUseCase_ListMembers(t *testing.T) {
	type fields struct {
		MemberRepo output.MemberRepository
	}
	type args struct {
		ctx        context.Context
		pagination pagination.Pagination
	}
	ctrl, ctx, testTime := repoHelper(t)
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []*entity.Member
		wantTotal int
		setupRepo func(*mock.MockMemberRepository)
		wantErr   error
	}{
		{
			name: "normal test",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				pagination: pagination.Pagination{
					Limit:   2,
					Offset:  0,
					SortBy:  "id",
					OrderBy: "asc",
					Total:   2,
				},
			},
			want: []*entity.Member{
				{
					ID:        1,
					Name:      "gg",
					Email:     "gg@gmail.com",
					Password:  "",
					CreatedAt: testTime,
				},
				{
					ID:        2,
					Name:      "gg1",
					Email:     "gg1@gmail.com",
					Password:  "",
					CreatedAt: testTime,
				},
			},
			wantTotal: 2,
			setupRepo: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					r.EXPECT().GetAll(ctx, gomock.Any()).Return([]*entity.Member{
						{
							ID:        1,
							Name:      "gg",
							Email:     "gg@gmail.com",
							Password:  "",
							CreatedAt: testTime,
						},
						{
							ID:        2,
							Name:      "gg1",
							Email:     "gg1@gmail.com",
							Password:  "",
							CreatedAt: testTime,
						},
					}, nil),
					r.EXPECT().CountAll(ctx).Return(2, nil),
				)
			},
			wantErr: nil,
		},
		{
			name: "no members",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx:        ctx,
				pagination: pagination.Pagination{},
			},
			want:      nil,
			wantTotal: 0,
			setupRepo: func(r *mock.MockMemberRepository) {
				r.EXPECT().GetAll(ctx, gomock.Any()).Return(nil, repository.ErrGatewayMemberNotFound)
			},
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberNotFound),
		},
		{
			name: "got members but count error",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				pagination: pagination.Pagination{
					Limit:   2,
					Offset:  0,
					SortBy:  "id",
					OrderBy: "asc",
					Total:   2,
				},
			},
			want:      nil,
			wantTotal: 0,
			setupRepo: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					r.EXPECT().GetAll(ctx, gomock.Any()).Return([]*entity.Member{
						{
							ID:        1,
							Name:      "gg",
							Email:     "gg@gmail.com",
							Password:  "",
							CreatedAt: testTime,
						},
						{
							ID:        2,
							Name:      "gg1",
							Email:     "gg1@gmail.com",
							Password:  "",
							CreatedAt: testTime,
						},
					}, nil),
					r.EXPECT().CountAll(ctx).Return(0, repository.ErrGatewayMemberDBError),
				)
			},
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberDBError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.fields.MemberRepo.(*mock.MockMemberRepository)
			m := &MemberUseCase{
				MemberRepo: mockRepo,
			}
			tt.setupRepo(mockRepo)
			got, gotTotal, err := m.ListMembers(tt.args.ctx, tt.args.pagination)
			t.Logf("got = %#v, want %#v", got, tt.want)
			t.Logf("gotTotal = %v, wantTotal %v", gotTotal, tt.wantTotal)
			t.Logf("err = %v, wantErr %v", err, tt.wantErr)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ListMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListMembers() got = %v, want %v", got, tt.want)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("ListMembers() got1 = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}

func TestMemberUseCase_RegisterMember(t *testing.T) {
	type fields struct {
		MemberRepo output.MemberRepository
	}
	type args struct {
		ctx    context.Context
		member *entity.Member
	}
	ctrl, ctx, testTime := repoHelper(t)

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *entity.Member
		wantErr   error
		setupRepo func(*mock.MockMemberRepository)
	}{
		{
			name: "normal test",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx:    ctx,
				member: &entity.Member{},
			},
			want: &entity.Member{
				ID:        1,
				Name:      "gg",
				Email:     "gg@gmail.com",
				Password:  "123455",
				CreatedAt: testTime,
			},
			wantErr: nil,
			setupRepo: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					// 第一次註冊不會回應任何資料
					r.EXPECT().Create(ctx, gomock.Any()).Return(nil),
					// 第二次利用email取得資料
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg",
						Email:     "gg@gmail.com",
						Password:  "123455",
						CreatedAt: testTime,
					}, nil),
				)
			},
		},
		{
			name: "first query already exist",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx:    ctx,
				member: nil,
			},
			want:    nil,
			wantErr: ErrUseCaseMemberAlreadyExists,
			setupRepo: func(r *mock.MockMemberRepository) {
				r.EXPECT().Create(ctx, gomock.Any()).Return(repository.ErrGatewayMemberAlreadyExists)
			},
		},
		{
			name: "second query not found",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx:    ctx,
				member: &entity.Member{},
			},
			want:    nil,
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberNotFound),
			setupRepo: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					// 第一次註冊不會回應任何資料
					r.EXPECT().Create(ctx, gomock.Any()).Return(nil),
					// 第二次利用email取得資料
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, repository.ErrGatewayMemberNotFound),
				)
			},
		},
		{
			name: "first query db error",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx:    ctx,
				member: &entity.Member{},
			},
			want:    nil,
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberDBError),
			setupRepo: func(r *mock.MockMemberRepository) {
				r.EXPECT().Create(ctx, gomock.Any()).Return(repository.ErrGatewayMemberDBError)
			},
		},
		{
			name: "second query db error",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx:    ctx,
				member: &entity.Member{},
			},
			want:    nil,
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberDBError),
			setupRepo: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					// 第一次註冊不會回應任何資料
					r.EXPECT().Create(ctx, gomock.Any()).Return(nil),
					// 第二次利用email取得資料
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, repository.ErrGatewayMemberDBError),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.fields.MemberRepo.(*mock.MockMemberRepository)
			m := &MemberUseCase{
				MemberRepo: mockRepo,
			}
			tt.setupRepo(mockRepo)
			got, err := m.RegisterMember(tt.args.ctx, tt.args.member)
			t.Logf("got = %v, want %v", got, tt.want)
			t.Logf("err = %v, wantErr %v", err, tt.wantErr)
			if !errors.Is(err, tt.wantErr) {

				t.Errorf("RegisterMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterMember() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberUseCase_UpdateMember(t *testing.T) {
	type fields struct {
		MemberRepo output.MemberRepository
	}
	type args struct {
		ctx   context.Context
		patch *inputmodel.PatchUpdateMemberInputModel
	}
	ctrl, ctx, testTime := repoHelper(t)
	stringPtr := func(s string) *string {
		return &s
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *entity.Member
		wantErr   error
		setupRepo func(*mock.MockMemberRepository)
	}{
		{
			name: "normal test",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				patch: &inputmodel.PatchUpdateMemberInputModel{
					ID:       1,
					Name:     stringPtr("gg1"),
					Email:    nil,
					Password: nil,
				},
			},
			want: &entity.Member{
				ID:        1,
				Name:      "gg1",
				Email:     "gg@gmail.com",
				Password:  "",
				CreatedAt: testTime,
			},
			wantErr: nil,
			setupRepo: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg",
						Email:     "gg@gmail.com",
						Password:  "",
						CreatedAt: testTime,
					}, nil),
					r.EXPECT().Update(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg1",
						Email:     "gg@gmail.com",
						Password:  "",
						CreatedAt: testTime,
					}, nil),
				)
			},
		},
		{
			name: "first member not found",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				patch: &inputmodel.PatchUpdateMemberInputModel{
					ID:       1,
					Name:     stringPtr("gg"),
					Email:    nil,
					Password: nil,
				},
			},
			want:    nil,
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberNotFound),
			setupRepo: func(r *mock.MockMemberRepository) {
				r.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, repository.ErrGatewayMemberNotFound)
			},
		},
		{
			name: "second update error",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				patch: &inputmodel.PatchUpdateMemberInputModel{
					ID:       1,
					Name:     stringPtr("gg1"),
					Email:    nil,
					Password: nil,
				},
			},
			want:    nil,
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberUpdateFailed),
			setupRepo: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg",
						Email:     "gg@gmail.com",
						Password:  "",
						CreatedAt: testTime,
					}, nil),
					r.EXPECT().Update(ctx, gomock.Any()).Return(nil, repository.ErrGatewayMemberUpdateFailed),
				)
			},
		},
		{
			name: "db connect error",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				patch: &inputmodel.PatchUpdateMemberInputModel{
					ID:       1,
					Name:     stringPtr("gg1"),
					Email:    nil,
					Password: nil,
				},
			},
			want:    nil,
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberDBError),
			setupRepo: func(r *mock.MockMemberRepository) {
				r.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, repository.ErrGatewayMemberDBError)
			},
		},
		{
			name: "gateway mapping error",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				patch: &inputmodel.PatchUpdateMemberInputModel{
					ID:       1,
					Name:     stringPtr("gg1"),
					Email:    nil,
					Password: nil,
				},
			},
			want:    nil,
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberMappingError),
			setupRepo: func(r *mock.MockMemberRepository) {
				r.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, repository.ErrGatewayMemberMappingError)
			},
		},
		{
			name: "no-update test",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx:   ctx,
				patch: &inputmodel.PatchUpdateMemberInputModel{},
			},
			want:    nil,
			wantErr: MapGatewayErrorToUseCaseError(repository.ErrGatewayMemberUpdateFailed),
			setupRepo: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg",
						Email:     "gg@gmail.com",
						Password:  "",
						CreatedAt: testTime,
					}, nil),
					r.EXPECT().Update(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg1",
						Email:     "gg@gmail.com",
						Password:  "",
						CreatedAt: testTime,
					}, repository.ErrGatewayMemberUpdateFailed),
				)
			},
		}, {
			name: "only email update",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				patch: &inputmodel.PatchUpdateMemberInputModel{
					ID:    1,
					Email: stringPtr("gg1@gmail.com"),
				},
			},
			want: &entity.Member{
				ID:        1,
				Name:      "gg",
				Email:     "gg1@gmial.com",
				Password:  "",
				CreatedAt: testTime,
			},
			wantErr: nil,
			setupRepo: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg",
						Email:     "gg@gmail.com",
						Password:  "",
						CreatedAt: testTime,
					}, nil),
					r.EXPECT().Update(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg",
						Email:     "gg1@gmial.com",
						Password:  "",
						CreatedAt: testTime,
					}, nil),
				)
			},
		}, {
			name: "only password update",
			fields: fields{
				MemberRepo: mock.NewMockMemberRepository(ctrl),
			},
			args: args{
				ctx: ctx,
				patch: &inputmodel.PatchUpdateMemberInputModel{
					ID:       1,
					Password: stringPtr("newpassword"),
				},
			},
			want: &entity.Member{
				ID:        1,
				Name:      "gg",
				Email:     "gg@gmail.com",
				Password:  "newpassword",
				CreatedAt: testTime,
			},
			wantErr: nil,
			setupRepo: func(r *mock.MockMemberRepository) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg",
						Email:     "gg@gmail.com",
						Password:  "oldpassword",
						CreatedAt: testTime,
					}, nil),
					r.EXPECT().Update(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg",
						Email:     "gg@gmail.com",
						Password:  "newpassword",
						CreatedAt: testTime,
					}, nil),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.fields.MemberRepo.(*mock.MockMemberRepository)
			m := &MemberUseCase{
				MemberRepo: mockRepo,
			}
			tt.setupRepo(mockRepo)
			got, err := m.UpdateMember(tt.args.ctx, tt.args.patch)
			t.Logf("got = %v, want %v", got, tt.want)
			t.Logf("err = %v, wantErr %v", err, tt.wantErr)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UpdateMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateMember() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMemberUseCase(t *testing.T) {
	repo := mock.NewMockMemberRepository(gomock.NewController(t))
	got := NewMemberUseCase(repo)
	// 確認got不是nil
	if got == nil {
		t.Errorf("NewMemberUseCase() = %v, want %v", got, repo)
	}
	// 使用型別斷言確認got是*MemberUseCase
	usecase, ok := got.(*MemberUseCase)
	if !ok {
		t.Errorf("NewMemberUseCase() = %v, want %v", got, repo)
	}
	// 確認注入
	if usecase.MemberRepo != repo {
		t.Errorf("NewMemberUseCase() = %v, want %v", usecase.MemberRepo, repo)
	}
}

func repoHelper(t *testing.T) (*gomock.Controller, context.Context, time.Time) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })
	ctx := context.Background()
	testTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	return ctrl, ctx, testTime
}
