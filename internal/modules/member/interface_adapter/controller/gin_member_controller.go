package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	gindto "github.com/tomoffice/go-clean-architecture/internal/framework/http/gin/dto"
	"github.com/tomoffice/go-clean-architecture/internal/framework/http/gin/errordefs"
	ginmapper "github.com/tomoffice/go-clean-architecture/internal/framework/http/gin/mapper"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/mapper"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/usecase/port/input"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/usecase/port/output"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"github.com/tomoffice/go-clean-architecture/pkg/tracer"
	"net/http"
)

type MemberController struct {
	usecase   input.MemberInputPort
	presenter output.MemberPresenter
	logger    logger.Logger
	tracer    tracer.Tracer
}

func NewMemberController(memberUseCase input.MemberInputPort, presenter output.MemberPresenter, log logger.Logger, tracer tracer.Tracer) *MemberController {
	baseLogger := log.With(logger.NewField("layer", "controller"))
	return &MemberController{
		usecase:   memberUseCase,
		presenter: presenter,
		logger:    baseLogger,
		tracer:    tracer,
	}
}

func (c *MemberController) Register(ctx *gin.Context) {
	// 創建帶有 context 的 logger 用於追蹤
	requestCtx, contextLogger, span := createTracedLogger(ctx.Request.Context(), c.tracer, c.logger)
	defer span.End()

	var ginReqDTO gindto.GinBindingRegisterMemberRequestDTO
	if err := ctx.ShouldBindJSON(&ginReqDTO); err != nil {
		contextLogger.Error("會員註冊參數綁定錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("content_type", ctx.GetHeader("Content-Type")),
		)
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToRegisterMemberDTO(ginReqDTO)
	if err := reqDTO.Validate(); err != nil {
		contextLogger.Error("會員註冊參數驗證錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("request_data", ginReqDTO),
		)
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.RegisterMemberDTOToEntity(reqDTO)
	member, err := c.usecase.RegisterMember(requestCtx, entity)
	if err != nil {
		contextLogger.Error("會員註冊 UseCase 執行錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("member_email", entity.Email),
		)
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentRegisterMember(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) GetByID(ctx *gin.Context) {
	// 創建帶有 context 的 logger 用於追蹤
	requestCtx, contextLogger, span := createTracedLogger(ctx.Request.Context(), c.tracer, c.logger)
	defer span.End()

	var ginReqDTO gindto.GinBindingGetMemberByIDURIRequestDTO
	if err := ctx.ShouldBindUri(&ginReqDTO); err != nil {
		contextLogger.Error("會員查詢(ID)參數綁定錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("uri", ctx.Request.RequestURI),
		)
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToGetMemberByIDDTO(ginReqDTO)
	if err := reqDTO.Validate(); err != nil {
		contextLogger.Error("會員查詢(ID)參數驗證錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("member_id", ginReqDTO.ID),
		)
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.GetMemberByIDDTOToEntity(reqDTO)
	member, err := c.usecase.GetMemberByID(requestCtx, entity.ID)
	if err != nil {
		contextLogger.Error("會員查詢(ID) UseCase 執行錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("member_id", entity.ID),
		)
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentGetMemberByID(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) GetByEmail(ctx *gin.Context) {
	// 創建帶有 context 的 logger 用於追蹤
	requestCtx, contextLogger, span := createTracedLogger(ctx.Request.Context(), c.tracer, c.logger)
	defer span.End()

	var ginReqDTO gindto.GinBindingGetMemberByEmailQueryRequestDTO
	if err := ctx.ShouldBindQuery(&ginReqDTO); err != nil {
		contextLogger.Error("會員查詢(Email)參數綁定錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("query", ctx.Request.URL.RawQuery),
		)
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToGetMemberByEmailDTO(ginReqDTO)
	if err := reqDTO.Validate(); err != nil {
		contextLogger.Error("會員查詢(Email)參數驗證錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("email", ginReqDTO.Email),
		)
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.GetMemberByEmailDTOToEntity(reqDTO)
	member, err := c.usecase.GetMemberByEmail(requestCtx, entity.Email)
	if err != nil {
		contextLogger.Error("會員查詢(Email) UseCase 執行錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("member_email", entity.Email),
		)
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentGetMemberByEmail(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) List(ctx *gin.Context) {
	// 創建帶有 context 的 logger 用於追蹤
	requestCtx, contextLogger, span := createTracedLogger(ctx.Request.Context(), c.tracer, c.logger)
	defer span.End()

	var ginReqDTO gindto.GinBindingListMemberQueryRequestDTO
	if err := ctx.ShouldBindQuery(&ginReqDTO); err != nil {
		contextLogger.Error("會員列表查詢參數綁定錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("query", ctx.Request.URL.RawQuery),
		)
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOtoListMemberDTO(ginReqDTO)
	if err := reqDTO.Validate(); err != nil {
		contextLogger.Error("會員列表查詢參數驗證錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("page", ginReqDTO.Page),
			logger.NewField("limit", ginReqDTO.Limit),
		)
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	pagination := mapper.ListMemberDTOToPagination(reqDTO)
	members, total, err := c.usecase.ListMembers(requestCtx, *pagination)
	if err != nil {
		contextLogger.Error("會員列表查詢 UseCase 執行錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("offset", pagination.Offset),
			logger.NewField("limit", pagination.Limit),
		)
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentListMembers(members, total)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) UpdateProfile(ctx *gin.Context) {
	// 創建帶有 context 的 logger 用於追蹤
	requestCtx, contextLogger, span := createTracedLogger(ctx.Request.Context(), c.tracer, c.logger)
	defer span.End()

	var ginURI gindto.GinBindingUpdateMemberURIRequestDTO
	if err := ctx.ShouldBindUri(&ginURI); err != nil {
		contextLogger.Error("會員資料更新 URI 參數綁定錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("uri", ctx.Request.RequestURI),
		)
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	var ginBody gindto.GinBindingUpdateMemberProfileBodyRequestDTO
	if err := ctx.ShouldBindJSON(&ginBody); err != nil {
		contextLogger.Error("會員資料更新 Body 參數綁定錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("content_type", ctx.GetHeader("Content-Type")),
		)
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToUpdateMemberProfileDTO(ginURI, ginBody)
	if err := reqDTO.Validate(); err != nil {
		contextLogger.Error("會員資料更新參數驗證錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("member_id", ginURI.ID),
		)
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	inputModel := mapper.UpdateMemberProfileDTOToInputModel(reqDTO)
	member, err := c.usecase.UpdateMemberProfile(requestCtx, inputModel)
	if err != nil {
		contextLogger.Error("會員資料更新 UseCase 執行錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("member_id", inputModel.ID),
		)
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentUpdateMemberProfile(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) UpdateEmail(ctx *gin.Context) {
	// 創建帶有 context 的 logger 用於追蹤
	requestCtx, contextLogger, span := createTracedLogger(ctx.Request.Context(), c.tracer, c.logger)
	defer span.End()

	var ginURI gindto.GinBindingUpdateMemberURIRequestDTO
	if err := ctx.ShouldBindUri(&ginURI); err != nil {
		contextLogger.Error("會員 Email 更新 URI 參數綁定錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("uri", ctx.Request.RequestURI),
		)
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	var ginBody gindto.GinBindingUpdateMemberEmailBodyRequestDTO
	if err := ctx.ShouldBindJSON(&ginBody); err != nil {
		contextLogger.Error("會員 Email 更新 Body 參數綁定錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("content_type", ctx.GetHeader("Content-Type")),
		)
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToUpdateMemberEmailDTO(ginURI, ginBody)
	if err := reqDTO.Validate(); err != nil {
		contextLogger.Error("會員 Email 更新參數驗證錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("member_id", ginURI.ID),
		)
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	inputModel := mapper.UpdateMemberEmailDTOToEntity(reqDTO)
	if err := c.usecase.UpdateMemberEmail(requestCtx, inputModel.ID, inputModel.Email, inputModel.Password); err != nil {
		contextLogger.Error("會員 Email 更新 UseCase 執行錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("member_id", inputModel.ID),
			logger.NewField("new_email", inputModel.Email),
		)
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentUpdateMemberEmail()
	ctx.JSON(http.StatusOK, resp)

}
func (c *MemberController) UpdatePassword(ctx *gin.Context) {
	// 創建帶有 context 的 logger 用於追蹤
	requestCtx, contextLogger, span := createTracedLogger(ctx.Request.Context(), c.tracer, c.logger)
	defer span.End()

	var ginURI gindto.GinBindingUpdateMemberURIRequestDTO
	if err := ctx.ShouldBindUri(&ginURI); err != nil {
		contextLogger.Error("會員密碼更新 URI 參數綁定錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("uri", ctx.Request.RequestURI),
		)
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	var ginBody gindto.GinBindingUpdateMemberPasswordBodyRequestDTO
	if err := ctx.ShouldBindJSON(&ginBody); err != nil {
		contextLogger.Error("會員密碼更新 Body 參數綁定錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("content_type", ctx.GetHeader("Content-Type")),
		)
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToUpdateMemberPasswordDTO(ginURI, ginBody)
	if err := reqDTO.Validate(); err != nil {
		contextLogger.Error("會員密碼更新參數驗證錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("member_id", ginURI.ID),
		)
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	inputModel := mapper.UpdateMemberPasswordDTOToInputModel(reqDTO)
	if err := c.usecase.UpdateMemberPassword(requestCtx, inputModel.ID, inputModel.OldPassword, inputModel.NewPassword); err != nil {
		contextLogger.Error("會員密碼更新 UseCase 執行錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("member_id", inputModel.ID),
		)
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentUpdateMemberPassword()
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) Delete(ctx *gin.Context) {
	// 創建帶有 context 的 logger 用於追蹤
	requestCtx, contextLogger, span := createTracedLogger(ctx.Request.Context(), c.tracer, c.logger)
	defer span.End()

	var ginReqDTO gindto.GinBindingDeleteMemberURIRequestDTO
	if err := ctx.ShouldBindUri(&ginReqDTO); err != nil {
		contextLogger.Error("會員刪除參數綁定錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("uri", ctx.Request.RequestURI),
		)
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToDeleteMemberDTO(ginReqDTO)
	if err := reqDTO.Validate(); err != nil {
		contextLogger.Error("會員刪除參數驗證錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("member_id", ginReqDTO.ID),
		)
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.DeleteMemberDTOToEntity(reqDTO)
	member, err := c.usecase.DeleteMember(requestCtx, entity.ID)
	if err != nil {
		contextLogger.Error("會員刪除 UseCase 執行錯誤",
			logger.NewField("error", err.Error()),
			logger.NewField("member_id", entity.ID),
		)
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentDeleteMember(member)
	ctx.JSON(http.StatusOK, resp)
}

func createTracedLogger(ctx context.Context, tr tracer.Tracer, log logger.Logger) (context.Context, logger.Logger, tracer.Span) {
	requestCtx, span := tr.Start(ctx, "")
	lg := log.WithContext(requestCtx)
	return requestCtx, lg, span
}
