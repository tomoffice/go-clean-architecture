package controller

import (
	"github.com/gin-gonic/gin"
	gindto "module-clean/internal/framework/http/gin/dto"
	"module-clean/internal/framework/http/gin/errordefs"
	ginmapper "module-clean/internal/framework/http/gin/mapper"
	"module-clean/internal/modules/member/interface_adapter/mapper"
	"module-clean/internal/modules/member/usecase/port/input"
	"module-clean/internal/modules/member/usecase/port/output"
	"net/http"
)

type MemberController struct {
	usecase   input.MemberInputPort
	presenter output.MemberPresenter
}

func NewMemberController(memberUseCase input.MemberInputPort, presenter output.MemberPresenter) *MemberController {
	return &MemberController{
		usecase:   memberUseCase,
		presenter: presenter,
	}
}

func (c *MemberController) Register(ctx *gin.Context) {
	var ginReqDTO gindto.GinCreateMemberRequestDTO
	if err := ctx.ShouldBindJSON(&ginReqDTO); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToCreateMemberDTO(ginReqDTO)
	if err := reqDTO.Validate(); err != nil {
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.CreateMemberDTOToEntity(reqDTO)
	member, err := c.usecase.RegisterMember(ctx, entity)
	if err != nil {
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentCreateMember(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) GetByID(ctx *gin.Context) {
	var ginReqDTO gindto.GinGetMemberByIDRequestDTO
	if err := ctx.ShouldBindUri(&ginReqDTO); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToGetMemberByIDDTO(ginReqDTO)
	if err := reqDTO.Validate(); err != nil {
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.GetMemberByIDDTOToEntity(reqDTO)
	member, err := c.usecase.GetMemberByID(ctx, entity.ID)
	if err != nil {
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentGetMemberByID(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) GetByEmail(ctx *gin.Context) {
	var ginReqDTO gindto.GinGetMemberByEmailRequestDTO
	if err := ctx.ShouldBindQuery(&ginReqDTO); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToGetMemberByEmailDTO(ginReqDTO)
	if err := reqDTO.Validate(); err != nil {
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.GetMemberByEmailDTOToEntity(reqDTO)
	member, err := c.usecase.GetMemberByEmail(ctx, entity.Email)
	if err != nil {
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentGetMemberByEmail(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) List(ctx *gin.Context) {
	var ginReqDTO gindto.GinListMemberRequestDTO
	if err := ctx.ShouldBindQuery(&ginReqDTO); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOtoListMemberDTO(ginReqDTO)
	if err := reqDTO.Validate(); err != nil {
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	pagination := mapper.ListMemberDTOToPagination(reqDTO)
	members, total, err := c.usecase.ListMembers(ctx, *pagination)
	if err != nil {
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentListMembers(members, total)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) Update(ctx *gin.Context) {
	var ginReqDTO gindto.GinUpdateMemberRequestDTO
	if err := ctx.ShouldBindJSON(&ginReqDTO); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToUpdateMemberDTO(ginReqDTO)
	if err := reqDTO.Validate(); err != nil {
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	inputModel := mapper.UpdateMemberDTOToInputModel(reqDTO)
	member, err := c.usecase.UpdateMember(ctx, inputModel)
	if err != nil {
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentUpdateMember(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) Delete(ctx *gin.Context) {
	var ginReqDTO gindto.GinDeleteMemberRequestDTO
	if err := ctx.ShouldBindUri(&ginReqDTO); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToDeleteMemberDTO(ginReqDTO)
	if err := reqDTO.Validate(); err != nil {
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.DeleteMemberDTOToEntity(reqDTO)
	member, err := c.usecase.DeleteMember(ctx, entity.ID)
	if err != nil {
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentDeleteMember(member)
	ctx.JSON(http.StatusOK, resp)
}
