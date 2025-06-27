package usecase

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/modules/member/usecase/inputmodel"
	"module-clean/internal/modules/member/usecase/mock"
	"module-clean/internal/modules/member/usecase/port/output"
	"module-clean/internal/shared/pagination"
	"reflect"
	"testing"
	"time"
)

func TestMemberUseCase_DeleteMember(t *testing.T) {
	type fields struct {
		MemberRepo output.MemberPersistence
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
		repoSetup func(*mock.MockMemberPersistence)
		wantErr   error
	}{
		{
			name: "normal test",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
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
			repoSetup: func(r *mock.MockMemberPersistence) {
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
		},
		{
			name: "no member found",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx: ctx,
				id:  0,
			},
			want: nil,
			repoSetup: func(r *mock.MockMemberPersistence) {
				r.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, ErrMemberNotFound)
			},
			wantErr: ErrMemberNotFound,
		},
		{
			name: "got member but delete error",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx: ctx,
				id:  0,
			},
			want: nil,
			repoSetup: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{}, nil),
					r.EXPECT().Delete(ctx, gomock.Any()).Return(ErrMemberDBError),
				)
			},
			wantErr: ErrMemberDBError,
		},
		{
			name: "delete no affect",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx: ctx,
				id:  0,
			},
			want: nil,
			repoSetup: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{}, nil),
					r.EXPECT().Delete(ctx, gomock.Any()).Return(ErrMemberNoEffect),
				)
			},
			wantErr: ErrMemberNoEffect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, ok := tt.fields.MemberRepo.(*mock.MockMemberPersistence)
			if !ok {
				t.Fatalf("expected *mock.MockMemberPersistence, got %T", tt.fields.MemberRepo)
			}
			m := &MemberUseCase{
				MemberGateway: mockRepo,
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
		MemberRepo output.MemberPersistence
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
		repoSetup func(*mock.MockMemberPersistence)
		wantErr   error
	}{
		{
			name: "normal test",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
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
			repoSetup: func(r *mock.MockMemberPersistence) {
				r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(&entity.Member{
					ID:        0,
					Name:      "gg",
					Email:     "gg@gmail.com",
					Password:  "",
					CreatedAt: testTime,
				}, nil)
			},
			wantErr: nil,
		},
		{
			name: "no member found",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:   ctx,
				email: "",
			},
			want: nil,
			repoSetup: func(r *mock.MockMemberPersistence) {
				r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, ErrMemberNotFound)
			},
			wantErr: ErrMemberNotFound,
		},
		{
			name: "db error",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:   ctx,
				email: "",
			},
			want: nil,
			repoSetup: func(r *mock.MockMemberPersistence) {
				r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, ErrMemberDBError)
			},
			wantErr: ErrMemberDBError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, ok := tt.fields.MemberRepo.(*mock.MockMemberPersistence)
			if !ok {
				t.Fatalf("expected *mock.MockMemberPersistence, got %T", tt.fields.MemberRepo)
			}
			m := &MemberUseCase{
				MemberGateway: mockRepo,
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
		MemberRepo output.MemberPersistence
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
		repoSetup func(*mock.MockMemberPersistence)
		wantErr   error
	}{
		{
			name: "normal test",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
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
			repoSetup: func(r *mock.MockMemberPersistence) {
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
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx: ctx,
				id:  0,
			},
			want: nil,
			repoSetup: func(r *mock.MockMemberPersistence) {
				r.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, ErrMemberNotFound)
			},
			wantErr: ErrMemberNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, ok := tt.fields.MemberRepo.(*mock.MockMemberPersistence)
			if !ok {
				t.Fatalf("expected *mock.MockMemberPersistence, got %T", tt.fields.MemberRepo)
			}
			m := &MemberUseCase{
				MemberGateway: mockRepo,
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
		MemberRepo output.MemberPersistence
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
		setupRepo func(*mock.MockMemberPersistence)
		wantErr   error
	}{
		{
			name: "normal test",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
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
			setupRepo: func(r *mock.MockMemberPersistence) {
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
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:        ctx,
				pagination: pagination.Pagination{},
			},
			want:      nil,
			wantTotal: 0,
			setupRepo: func(r *mock.MockMemberPersistence) {
				r.EXPECT().GetAll(ctx, gomock.Any()).Return(nil, ErrMemberNotFound)
			},
			wantErr: ErrMemberNotFound,
		},
		{
			name: "got members but count error",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
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
			setupRepo: func(r *mock.MockMemberPersistence) {
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
					r.EXPECT().CountAll(ctx).Return(0, ErrMemberDBError),
				)
			},
			wantErr: ErrMemberDBError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.fields.MemberRepo.(*mock.MockMemberPersistence)
			m := &MemberUseCase{
				MemberGateway: mockRepo,
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
		MemberRepo output.MemberPersistence
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
		setupRepo func(*mock.MockMemberPersistence)
	}{
		{
			name: "normal test",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
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
			setupRepo: func(r *mock.MockMemberPersistence) {
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
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:    ctx,
				member: nil,
			},
			want:    nil,
			wantErr: ErrMemberAlreadyExists,
			setupRepo: func(r *mock.MockMemberPersistence) {
				r.EXPECT().Create(ctx, gomock.Any()).Return(ErrMemberAlreadyExists)
			},
		},
		{
			name: "second query not found",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:    ctx,
				member: &entity.Member{},
			},
			want:    nil,
			wantErr: ErrMemberNotFound,
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					// 第一次註冊不會回應任何資料
					r.EXPECT().Create(ctx, gomock.Any()).Return(nil),
					// 第二次利用email取得資料
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, ErrMemberNotFound),
				)
			},
		},
		{
			name: "first query db error",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:    ctx,
				member: &entity.Member{},
			},
			want:    nil,
			wantErr: ErrMemberDBError,
			setupRepo: func(r *mock.MockMemberPersistence) {
				r.EXPECT().Create(ctx, gomock.Any()).Return(ErrMemberDBError)
			},
		},
		{
			name: "second query db error",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:    ctx,
				member: &entity.Member{},
			},
			want:    nil,
			wantErr: ErrMemberDBError,
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					// 第一次註冊不會回應任何資料
					r.EXPECT().Create(ctx, gomock.Any()).Return(nil),
					// 第二次利用email取得資料
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, ErrMemberDBError),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.fields.MemberRepo.(*mock.MockMemberPersistence)
			m := &MemberUseCase{
				MemberGateway: mockRepo,
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

func TestMemberUseCase_UpdateMemberProfile(t *testing.T) {
	type fields struct {
		MemberRepo output.MemberPersistence
	}
	type args struct {
		ctx   context.Context
		patch *inputmodel.PatchUpdateMemberProfileInputModel
	}
	ctrl, ctx, testTime := repoHelper(t)
	stringPtr := func(s string) *string {
		return &s
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setupRepo func(*mock.MockMemberPersistence)
		want      *entity.Member
		wantErr   error
	}{
		{
			name: "normal test",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx: ctx,
				patch: &inputmodel.PatchUpdateMemberProfileInputModel{
					ID:   1,
					Name: stringPtr("gg1"),
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
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg",
						Email:     "gg@gmail.com",
						Password:  "",
						CreatedAt: testTime,
					}, nil),
					r.EXPECT().UpdateProfile(ctx, gomock.Any()).Return(&entity.Member{
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
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx: ctx,
				patch: &inputmodel.PatchUpdateMemberProfileInputModel{
					ID:   1,
					Name: stringPtr("gg"),
				},
			},
			want:    nil,
			wantErr: ErrMemberNotFound,
			setupRepo: func(r *mock.MockMemberPersistence) {
				r.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, ErrMemberNotFound)
			},
		},
		{
			name: "second update error",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx: ctx,
				patch: &inputmodel.PatchUpdateMemberProfileInputModel{
					ID:   1,
					Name: stringPtr("gg1"),
				},
			},
			want:    nil,
			wantErr: ErrMemberNoEffect,
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg",
						Email:     "gg@gmail.com",
						Password:  "",
						CreatedAt: testTime,
					}, nil),
					r.EXPECT().UpdateProfile(ctx, gomock.Any()).Return(nil, ErrMemberNoEffect),
				)
			},
		},
		{
			name: "db connect error",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx: ctx,
				patch: &inputmodel.PatchUpdateMemberProfileInputModel{
					ID:   1,
					Name: stringPtr("gg1"),
				},
			},
			want:    nil,
			wantErr: ErrMemberDBError,
			setupRepo: func(r *mock.MockMemberPersistence) {
				r.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, ErrMemberDBError)
			},
		},
		{
			name: "gateway mapping error",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx: ctx,
				patch: &inputmodel.PatchUpdateMemberProfileInputModel{
					ID:   1,
					Name: stringPtr("gg1"),
				},
			},
			want:    nil,
			wantErr: ErrMemberMappingError,
			setupRepo: func(r *mock.MockMemberPersistence) {
				r.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, ErrMemberMappingError)
			},
		},
		{
			name: "no-update test",
			fields: fields{
				MemberRepo: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:   ctx,
				patch: &inputmodel.PatchUpdateMemberProfileInputModel{},
			},
			want:    nil,
			wantErr: ErrMemberNoEffect,
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg",
						Email:     "gg@gmail.com",
						Password:  "",
						CreatedAt: testTime,
					}, nil),
					r.EXPECT().UpdateProfile(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "gg1",
						Email:     "gg@gmail.com",
						Password:  "",
						CreatedAt: testTime,
					}, ErrMemberNoEffect),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.fields.MemberRepo.(*mock.MockMemberPersistence)
			m := &MemberUseCase{
				MemberGateway: mockRepo,
			}
			tt.setupRepo(mockRepo)
			got, err := m.UpdateMemberProfile(tt.args.ctx, tt.args.patch)
			t.Logf("got = %v, want %v", got, tt.want)
			t.Logf("err = %v, wantErr %v", err, tt.wantErr)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UpdateMemberProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateMemberProfile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberUseCase_UpdateMemberEmail(t *testing.T) {
	type fields struct {
		MemberGateway output.MemberPersistence
	}
	type args struct {
		ctx      context.Context
		id       int
		newEmail string
		password string
	}
	ctrl, ctx, testTime := repoHelper(t)
	tests := []struct {
		name      string
		fields    fields
		args      args
		setupRepo func(*mock.MockMemberPersistence)
		wantErr   error
	}{
		{
			name: "normal test",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:      ctx,
				id:       0,
				newEmail: "oldemail@gmail.com",
				password: "testpassword",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, ErrMemberNotFound),
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:        0,
						Name:      "test",
						Email:     "test@gmail.com",
						Password:  "testpassword",
						CreatedAt: testTime,
					}, nil),
					r.EXPECT().UpdateEmail(ctx, gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			wantErr: nil,
		},
		{
			name: "email already used",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:      ctx,
				id:       1,
				newEmail: "",
				password: "",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(&entity.Member{
					ID:       2,
					Name:     "test",
					Email:    "",
					Password: "",
				}, nil)
			},
			wantErr: ErrMemberEmailAlreadyExists,
		},
		{
			name: "other error",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:      ctx,
				id:       1,
				newEmail: "",
				password: "",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, ErrMemberDBError),
				)
			},
			wantErr: ErrMemberDBError,
		},
		{
			name: "oldEmail is same as newEmail",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:      ctx,
				id:       1,
				newEmail: "test@gmail.com",
				password: "",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(&entity.Member{
						ID:    1,
						Email: "test@gmail.com",
					}, nil),
				)
			},
			wantErr: ErrMemberUpdateSameEmail,
		},
		{
			name: "GetByID error - member not found",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:      ctx,
				id:       1,
				newEmail: "",
				password: "",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, ErrMemberNotFound),
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, ErrMemberNotFound),
				)
			},
			wantErr: ErrMemberNotFound,
		},
		{
			name: "GetByID error - db error",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:      ctx,
				id:       1,
				password: "testpassword",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, ErrMemberNotFound),
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, ErrMemberDBError),
				)
			},
			wantErr: ErrMemberDBError,
		},
		{
			name: "password mismatch",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:      ctx,
				id:       1,
				password: "wrongpassword",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, ErrMemberNotFound),
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:       1,
						Password: "rightpassword",
					}, nil),
				)
			},
			wantErr: ErrMemberPasswordIncorrect,
		},
		{
			name: "UpdateEmail error - no effect",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx: ctx,
				id:  0,
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, ErrMemberNotFound),
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{}, nil),
					r.EXPECT().UpdateEmail(ctx, gomock.Any(), gomock.Any()).Return(ErrMemberNoEffect),
				)
			},

			wantErr: ErrMemberNoEffect,
		},
		{
			name: "UpdateEmail error - db error",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx: ctx,
				id:  0,
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, ErrMemberNotFound),
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{}, nil),
					r.EXPECT().UpdateEmail(ctx, gomock.Any(), gomock.Any()).Return(ErrMemberDBError),
				)
			},
			wantErr: ErrMemberDBError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemberUseCase{
				MemberGateway: tt.fields.MemberGateway,
			}
			tt.setupRepo(tt.fields.MemberGateway.(*mock.MockMemberPersistence))
			got := m.UpdateMemberEmail(tt.args.ctx, tt.args.id, tt.args.newEmail, tt.args.password)
			if got != nil && tt.wantErr == nil {
				t.Fatalf("UpdateMemberEmail() got unexpected error: %v", got)
			}
			t.Logf("\n\tgot = %v\n\twantErr = %v", got, tt.wantErr)
			assert.Equal(t, got, tt.wantErr, "UpdateMemberEmail() got = %v, wantErr %v", got, tt.wantErr)
		})
	}
}

func TestMemberUseCase_UpdateMemberPassword(t *testing.T) {
	type fields struct {
		MemberGateway output.MemberPersistence
	}
	type args struct {
		ctx         context.Context
		id          int
		newPassword string
		oldPassword string
	}
	ctrl, ctx, testTime := repoHelper(t)
	tests := []struct {
		name      string
		fields    fields
		args      args
		setupRepo func(*mock.MockMemberPersistence)
		wantErr   error
	}{
		{
			name: "normal case",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:         ctx,
				id:          1,
				newPassword: "newpassword",
				oldPassword: "oldpassword",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:        1,
						Name:      "test",
						Password:  "oldpassword",
						CreatedAt: testTime,
					}, nil),
					r.EXPECT().UpdatePassword(ctx, gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			wantErr: nil,
		},
		{
			name: "logic error - update same password",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:         ctx,
				id:          0,
				newPassword: "password",
				oldPassword: "password",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {},
			wantErr:   ErrMemberUpdateSamePassword,
		},
		{
			name: "GetByID error - member not found",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:         ctx,
				id:          0,
				newPassword: "newpassword",
				oldPassword: "oldpassword",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, ErrMemberNotFound),
				)
			},
			wantErr: ErrMemberNotFound,
		},
		{
			name: "GetByID error - db error",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:         ctx,
				id:          0,
				newPassword: "newpassword",
				oldPassword: "oldpassword",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(nil, ErrMemberDBError),
				)
			},
			wantErr: ErrMemberDBError,
		},
		{
			name: "logic error - db password not input password",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:         ctx,
				id:          0,
				oldPassword: "wrongpassword",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:       0,
						Password: "oldpassword",
					}, nil),
				)
			},
			wantErr: ErrMemberPasswordIncorrect,
		},
		{
			name: "UpdatePassword error - no effect",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:         ctx,
				id:          0,
				newPassword: "newpassword",
				oldPassword: "oldpassword",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:       0,
						Password: "oldpassword",
					}, nil),
					r.EXPECT().UpdatePassword(ctx, gomock.Any(), gomock.Any()).Return(ErrMemberNoEffect),
				)
			},
			wantErr: ErrMemberNoEffect,
		},
		{
			name: "UpdatePassword error - db error",
			fields: fields{
				MemberGateway: mock.NewMockMemberPersistence(ctrl),
			},
			args: args{
				ctx:         ctx,
				id:          0,
				newPassword: "newpassword",
				oldPassword: "oldpassword",
			},
			setupRepo: func(r *mock.MockMemberPersistence) {
				gomock.InOrder(
					r.EXPECT().GetByID(ctx, gomock.Any()).Return(&entity.Member{
						ID:       0,
						Password: "oldpassword",
					}, nil),
					r.EXPECT().UpdatePassword(ctx, gomock.Any(), gomock.Any()).Return(ErrMemberDBError),
				)
			},
			wantErr: ErrMemberDBError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemberUseCase{
				MemberGateway: tt.fields.MemberGateway,
			}
			tt.setupRepo(tt.fields.MemberGateway.(*mock.MockMemberPersistence))
			err := m.UpdateMemberPassword(tt.args.ctx, tt.args.id, tt.args.newPassword, tt.args.oldPassword)
			if err != nil && tt.wantErr == nil {
				t.Fatalf("UpdateMemberPassword() got unexpected error: %v", err)
			}
			t.Logf("\n\tgot = %v\n\twantErr = %v", err, tt.wantErr)
			assert.Equal(t, err, tt.wantErr, "UpdateMemberPassword() got = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func TestNewMemberUseCase(t *testing.T) {
	repo := mock.NewMockMemberPersistence(gomock.NewController(t))
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
	if usecase.MemberGateway != repo {
		t.Errorf("NewMemberUseCase() = %v, want %v", usecase.MemberGateway, repo)
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
