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
	var ginReqDTO gindto.GinRegisterMemberRequestDTO
	if err := ctx.ShouldBindJSON(&ginReqDTO); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToRegisterMemberDTO(ginReqDTO)
	if err := reqDTO.Validate(); err != nil {
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.RegisterMemberDTOToEntity(reqDTO)
	member, err := c.usecase.RegisterMember(ctx, entity)
	if err != nil {
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentRegisterMember(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) GetByID(ctx *gin.Context) {
	var ginReqDTO gindto.GinGetMemberByIDURIRequestDTO
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
	var ginReqDTO gindto.GinGetMemberByEmailQueryRequestDTO
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
	var ginReqDTO gindto.GinListMemberQueryRequestDTO
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
func (c *MemberController) UpdateProfile(ctx *gin.Context) {
	var ginURI gindto.GinUpdateMemberURIRequestDTO
	if err := ctx.ShouldBindUri(&ginURI); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	var ginBody gindto.GinUpdateMemberProfileBodyRequestDTO
	if err := ctx.ShouldBindJSON(&ginBody); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToUpdateMemberProfileDTO(ginURI, ginBody)
	if err := reqDTO.Validate(); err != nil {
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	inputModel := mapper.UpdateMemberProfileDTOToInputModel(reqDTO)
	member, err := c.usecase.UpdateMemberProfile(ctx, inputModel)
	if err != nil {
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentUpdateMemberProfile(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) UpdateEmail(ctx *gin.Context) {
	var ginURI gindto.GinUpdateMemberURIRequestDTO
	if err := ctx.ShouldBindUri(&ginURI); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	var ginBody gindto.GinUpdateMemberEmailBodyRequestDTO
	if err := ctx.ShouldBindJSON(&ginBody); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToUpdateMemberEmailDTO(ginURI, ginBody)
	if err := reqDTO.Validate(); err != nil {
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	inputModel := mapper.UpdateMemberEmailDTOToEntity(reqDTO)
	if err := c.usecase.UpdateMemberEmail(ctx, inputModel.ID, inputModel.Email, inputModel.Password); err != nil {
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentUpdateMemberEmail()
	ctx.JSON(http.StatusOK, resp)

}
func (c *MemberController) UpdatePassword(ctx *gin.Context) {
	var ginURI gindto.GinUpdateMemberURIRequestDTO
	if err := ctx.ShouldBindUri(&ginURI); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	var ginBody gindto.GinUpdateMemberPasswordBodyRequestDTO
	if err := ctx.ShouldBindJSON(&ginBody); err != nil {
		errCode, errMsg := errordefs.MapGinBindingError(err)
		resp := c.presenter.PresentBindingError(errCode, errMsg)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	reqDTO := ginmapper.GinDTOToUpdateMemberPasswordDTO(ginURI, ginBody)
	if err := reqDTO.Validate(); err != nil {
		errCode, resp := c.presenter.PresentValidationError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	inputModel := mapper.UpdateMemberPasswordDTOToInputModel(reqDTO)
	if err := c.usecase.UpdateMemberPassword(ctx, inputModel.ID, inputModel.OldPassword, inputModel.NewPassword); err != nil {
		errCode, resp := c.presenter.PresentUseCaseError(err)
		httpStatus := MapErrorCodeToHTTPStatus(errCode)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentUpdateMemberPassword()
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) Delete(ctx *gin.Context) {
	var ginReqDTO gindto.GinDeleteMemberURIRequestDTO
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
