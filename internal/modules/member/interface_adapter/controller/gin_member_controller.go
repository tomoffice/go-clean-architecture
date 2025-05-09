package controller

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/modules/member/interface_adapter/dto"
	"module-clean/internal/modules/member/interface_adapter/mapper"
	"module-clean/internal/modules/member/usecase"
	"net/http"
)

type MemberController struct {
	useCase   usecase.MemberInputPort
	presenter usecase.MemberOutputPort
}

func NewMemberController(memberUseCase *usecase.MemberUseCase, presenter usecase.MemberOutputPort) *MemberController {
	return &MemberController{
		useCase:   memberUseCase,
		presenter: presenter,
	}
}

func (c *MemberController) Register(ctx *gin.Context) {
	var reqDTO dto.CreateMemberRequestDTO
	if err := ctx.ShouldBindJSON(&reqDTO); err != nil {
		httpStatus, resp := c.presenter.PresentBindingError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	if err := reqDTO.Validate(); err != nil {
		httpStatus, resp := c.presenter.PresentValidationError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.CreateMemberDTOToEntity(reqDTO)
	member, err := c.useCase.RegisterMember(ctx, entity)
	if err != nil {
		httpStatus, resp := c.presenter.PresentUseCaseError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentCreateMember(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) GetByID(ctx *gin.Context) {
	var reqDTO dto.GetMemberByIDRequestDTO
	if err := ctx.ShouldBindUri(&reqDTO); err != nil {
		httpStatus, resp := c.presenter.PresentBindingError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	if err := reqDTO.Validate(); err != nil {
		httpStatus, resp := c.presenter.PresentValidationError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.GetMemberByIDDTOToEntity(reqDTO)
	member, err := c.useCase.GetMemberByID(ctx, entity.ID)
	if err != nil {
		httpStatus, resp := c.presenter.PresentUseCaseError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentCreateMember(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) GetByEmail(ctx *gin.Context) {
	var reqDTO dto.GetMemberByEmailRequestDTO
	if err := ctx.ShouldBindQuery(&reqDTO); err != nil {
		httpStatus, resp := c.presenter.PresentBindingError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	if err := reqDTO.Validate(); err != nil {
		httpStatus, resp := c.presenter.PresentValidationError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.GetMemberByEmailDTOToEntity(reqDTO)
	member, err := c.useCase.GetMemberByEmail(ctx, entity.Email)
	if err != nil {
		httpStatus, resp := c.presenter.PresentUseCaseError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentGetMemberByEmail(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) List(ctx *gin.Context) {
	var reqDTO dto.ListMemberRequestDTO
	if err := ctx.ShouldBindQuery(&reqDTO); err != nil {
		httpStatus, resp := c.presenter.PresentBindingError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	if err := reqDTO.Validate(); err != nil {
		httpStatus, resp := c.presenter.PresentValidationError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	pagination := mapper.ListMemberDTOToPagination(reqDTO)
	members, total, err := c.useCase.ListMembers(ctx, *pagination)
	if err != nil {
		httpStatus, resp := c.presenter.PresentUseCaseError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentListMembers(members, total)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) Update(ctx *gin.Context) {
	var reqDTO dto.UpdateMemberRequestDTO
	if err := ctx.ShouldBindJSON(&reqDTO); err != nil {
		httpStatus, resp := c.presenter.PresentBindingError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	if err := reqDTO.Validate(); err != nil {
		httpStatus, resp := c.presenter.PresentValidationError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	inputModel := mapper.UpdateMemberDTOToInputModel(reqDTO)
	member, err := c.useCase.UpdateMember(ctx, inputModel)
	if err != nil {
		httpStatus, resp := c.presenter.PresentUseCaseError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentUpdateMember(member)
	ctx.JSON(http.StatusOK, resp)
}
func (c *MemberController) Delete(ctx *gin.Context) {
	var reqDTO dto.DeleteMemberRequestDTO
	if err := ctx.ShouldBindUri(&reqDTO); err != nil {
		httpStatus, resp := c.presenter.PresentBindingError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	if err := reqDTO.Validate(); err != nil {
		httpStatus, resp := c.presenter.PresentValidationError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	entity := mapper.DeleteMemberDTOToEntity(reqDTO)
	member, err := c.useCase.DeleteMember(ctx, entity.ID)
	if err != nil {
		httpStatus, resp := c.presenter.PresentUseCaseError(err)
		ctx.JSON(httpStatus, resp)
		return
	}
	resp := c.presenter.PresentDeleteMember(member)
	ctx.JSON(http.StatusOK, resp)
}
