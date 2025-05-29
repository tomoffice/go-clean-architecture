package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"module-clean/internal/modules/member/entity"
	"module-clean/internal/modules/member/interface_adapter/controller/mock"
	"module-clean/internal/modules/member/interface_adapter/dto"
	"module-clean/internal/modules/member/interface_adapter/outputmodel"
	"module-clean/internal/modules/member/usecase"
	"module-clean/internal/modules/member/usecase/port/input"
	"module-clean/internal/modules/member/usecase/port/output"
	"module-clean/internal/shared/common/enum"
	"module-clean/internal/shared/common/errorcode"
	sharedviewmodel "module-clean/internal/shared/interface_adapter/viewmodel/http"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestMemberController_Delete(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	ctrl, testTime := portHelper(t)
	tests := []struct {
		name        string
		fields      fields
		args        args
		setupPort   func(*mock.MockMemberInputPort, *mock.MockMemberPresenter)
		setupGinCtx func(ginCtx *gin.Context)
		want        *outputmodel.DeleteMemberResponse
		wantErr     *outputmodel.ErrorResponse
		wantStatus  int
	}{
		{
			name: "normal case",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().DeleteMember(gomock.Any(), gomock.Any()).Return(
					&entity.Member{
						ID:        1,
						Name:      "test",
						Email:     "test@gmail.com",
						CreatedAt: testTime,
					}, nil)
				p.EXPECT().PresentDeleteMember(gomock.Any()).Return(outputmodel.DeleteMemberResponse{
					Data: dto.DeleteMemberResponseDTO{
						ID:        1,
						Name:      "test",
						Email:     "test@gmail.com",
						CreatedAt: testTime.Format(time.RFC3339),
					},
					Error: nil,
					Meta:  nil,
					BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
						UnixTimeStamp:    testTime.Unix(),
						RFC3339TimeStamp: testTime.Format(time.RFC3339),
						Status:           enum.APIStatusSuccess,
					},
				})
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: &outputmodel.DeleteMemberResponse{
				Data: dto.DeleteMemberResponseDTO{
					ID:        1,
					Name:      "test",
					Email:     "test@gmail.com",
					CreatedAt: testTime.Format(time.RFC3339),
				},
				Error: nil,
				Meta:  nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusSuccess,
				},
			},
			wantErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name: "binding error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "test binding error",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					})
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "invalid_id"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "test binding error",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "validation error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrInvalidParams, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "test validation error",
						},
						Meta:             nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "-1"}, // 因為validation會檢查ID是否大於0，所以這裡使用-1來觸發驗證錯誤
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "test validation error",
				},
				Meta:             nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{},
			},
			wantStatus: 400,
		},
		{
			name: "usecase error - member not found",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().DeleteMember(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrUseCaseMemberNotFound)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNotFound,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 404,
		},
		{
			name: "usecase error - delete member failed",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().DeleteMember(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrUseCaseMemberDeleteFailed)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberDeleteFailed,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberDeleteFailed),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)

			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberDeleteFailed),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 500,
		},
		{
			name: "usecase error - db error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().DeleteMember(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrUseCaseMemberDBFailure)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberDBFailure,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberDBFailure),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberDBFailure),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 500,
		},
		{
			name: "usecase error - unexpected error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().DeleteMember(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrUseCaseMemberUnexpectedFail)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrUnexpectedMemberUseCaseFail,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrUnexpectedMemberUseCaseFail),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Params = gin.Params{
					gin.Param{Key: "id", Value: "1"},
				}
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrUnexpectedMemberUseCaseFail),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := tt.fields.usecase.(*mock.MockMemberInputPort)
			mockPresenter := tt.fields.presenter.(*mock.MockMemberPresenter)
			c := &MemberController{
				usecase:   mockUseCase,
				presenter: mockPresenter,
			}
			ginCtx, responseWriter := GinCtxHelper(t)
			tt.setupPort(mockUseCase, mockPresenter)
			tt.setupGinCtx(ginCtx)
			c.Delete(ginCtx)
			response := responseWriter.Body.String()
			gotStatus := responseWriter.Code
			t.Logf("\n\tgotStatus:%d,wantStatus:%d", gotStatus, tt.wantStatus)
			assert.Equal(t, tt.wantStatus, gotStatus)
			if tt.want != nil {
				var got outputmodel.DeleteMemberResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twant %+v", got, *tt.want)
				assert.Equal(t, tt.want, &got)
			} else if tt.wantErr != nil {
				var got outputmodel.ErrorResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twantErr %+v", got, *tt.wantErr)
				assert.Equal(t, tt.wantErr, &got)
			}
		})
	}
}

func TestMemberController_GetByEmail(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	ctrl, testTime := portHelper(t)
	tests := []struct {
		name        string
		fields      fields
		args        args
		setupGinCtx func(ginCtx *gin.Context)
		setupPort   func(*mock.MockMemberInputPort, *mock.MockMemberPresenter)
		want        *outputmodel.GetMemberByEmailResponse
		wantErr     *outputmodel.ErrorResponse
		wantStatus  int
	}{
		{
			name: "normal case",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request, _ = http.NewRequest("GET", "/api/member?email=test@gmail.com", nil)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().GetMemberByEmail(gomock.Any(), gomock.Any()).Return(
					&entity.Member{
						ID:        0,
						Name:      "",
						Email:     "",
						CreatedAt: testTime}, nil)
				p.EXPECT().PresentGetMemberByEmail(gomock.Any()).Return(
					outputmodel.GetMemberByEmailResponse{
						Data: dto.GetMemberByEmailResponseDTO{
							ID:        0,
							Name:      "",
							Email:     "",
							CreatedAt: "",
						},
						Error: nil,
						Meta:  nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusSuccess,
						},
					},
				)
			},
			want: &outputmodel.GetMemberByEmailResponse{
				Data: dto.GetMemberByEmailResponseDTO{
					ID:        0,
					Name:      "",
					Email:     "",
					CreatedAt: "",
				},
				Error: nil,
				Meta:  nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusSuccess,
				},
			},
			wantErr:    nil,
			wantStatus: 200,
		},
		{
			name: "binding error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"GET",
					"/api/member?email=",
					nil,
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentBindingError(gomock.Any(), gomock.Any()).Return(
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					})
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 400,
		},
		{
			name: "validation error",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"GET",
					"/api/member?email=invalid-email",
					nil,
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				p.EXPECT().PresentValidationError(gomock.Any()).Return(
					errorcode.ErrInvalidParams, outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrInvalidParams),
							Message: "test validation error",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrInvalidParams),
					Message: "test validation error",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 400,
		},
		{
			name: "usecase error - member not found",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"GET",
					"/api/member?email=test@gmail.com",
					nil,
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().GetMemberByEmail(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrUseCaseMemberNotFound)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberNotFound,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberNotFound),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 404,
		},
		{
			name: "usecase error - db failure",
			fields: fields{
				usecase:   mock.NewMockMemberInputPort(ctrl),
				presenter: mock.NewMockMemberPresenter(ctrl),
			},
			args: args{
				ctx: nil,
			},
			setupGinCtx: func(ginCtx *gin.Context) {
				ginCtx.Request = httptest.NewRequest(
					"GET",
					"/api/member?email=test@gmail.com",
					nil,
				)
			},
			setupPort: func(uc *mock.MockMemberInputPort, p *mock.MockMemberPresenter) {
				uc.EXPECT().GetMemberByEmail(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrUseCaseMemberDBFailure)
				p.EXPECT().PresentUseCaseError(gomock.Any()).Return(
					errorcode.ErrMemberDBFailure,
					outputmodel.ErrorResponse{
						Data: nil,
						Error: &sharedviewmodel.ErrorPayload{
							Code:    strconv.Itoa(errorcode.ErrMemberDBFailure),
							Message: "",
						},
						Meta: nil,
						BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
							UnixTimeStamp:    testTime.Unix(),
							RFC3339TimeStamp: testTime.Format(time.RFC3339),
							Status:           enum.APIStatusFailed,
						},
					},
				)
			},
			want: nil,
			wantErr: &outputmodel.ErrorResponse{
				Data: nil,
				Error: &sharedviewmodel.ErrorPayload{
					Code:    strconv.Itoa(errorcode.ErrMemberDBFailure),
					Message: "",
				},
				Meta: nil,
				BaseHTTPResponse: sharedviewmodel.BaseHTTPResponse{
					UnixTimeStamp:    testTime.Unix(),
					RFC3339TimeStamp: testTime.Format(time.RFC3339),
					Status:           enum.APIStatusFailed,
				},
			},
			wantStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := tt.fields.usecase.(*mock.MockMemberInputPort)
			mockPresenter := tt.fields.presenter.(*mock.MockMemberPresenter)
			c := &MemberController{
				usecase:   mockUseCase,
				presenter: mockPresenter,
			}
			ginCtx, responseWriter := GinCtxHelper(t)
			tt.setupGinCtx(ginCtx)
			tt.setupPort(mockUseCase, mockPresenter)
			c.GetByEmail(ginCtx)
			response := responseWriter.Body.String()
			gotStatus := responseWriter.Code
			t.Logf("\n\tgotStatus:%d,wantStatus:%d", gotStatus, tt.wantStatus)
			assert.Equal(t, tt.wantStatus, gotStatus)
			if tt.want != nil {
				var got outputmodel.GetMemberByEmailResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twant %+v", got, *tt.want)
				assert.Equal(t, tt.want, &got)
			} else if tt.wantErr != nil {
				var got outputmodel.ErrorResponse
				err := json.Unmarshal([]byte(response), &got)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				t.Logf("\n\tgot %+v\n\twantErr %+v", got, *tt.wantErr)
				assert.Equal(t, tt.wantErr, &got)
			}
		})
	}
}

func TestMemberController_GetByID(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemberController{
				usecase:   tt.fields.usecase,
				presenter: tt.fields.presenter,
			}
			c.GetByID(tt.args.ctx)
		})
	}
}

func TestMemberController_List(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemberController{
				usecase:   tt.fields.usecase,
				presenter: tt.fields.presenter,
			}
			c.List(tt.args.ctx)
		})
	}
}

func TestMemberController_Register(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemberController{
				usecase:   tt.fields.usecase,
				presenter: tt.fields.presenter,
			}
			c.Register(tt.args.ctx)
		})
	}
}

func TestMemberController_Update(t *testing.T) {
	type fields struct {
		usecase   input.MemberInputPort
		presenter output.MemberPresenter
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemberController{
				usecase:   tt.fields.usecase,
				presenter: tt.fields.presenter,
			}
			c.Update(tt.args.ctx)
		})
	}
}

func TestNewMemberController(t *testing.T) {
	type args struct {
		memberUseCase input.MemberInputPort
		presenter     output.MemberPresenter
	}
	tests := []struct {
		name string
		args args
		want *MemberController
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemberController(tt.args.memberUseCase, tt.args.presenter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemberController() = %v, want %v", got, tt.want)
			}
		})
	}
}
func portHelper(t *testing.T) (*gomock.Controller, time.Time) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(
		func() {
			ctrl.Finish()
		},
	)
	testTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) // Set a fixed time for tests
	return ctrl, testTime
}
func GinCtxHelper(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	return ginCtx, w
}
