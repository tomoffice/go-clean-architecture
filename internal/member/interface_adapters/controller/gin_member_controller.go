package controller

import (
	"github.com/gin-gonic/gin"
	ginresponse "module-clean/internal/framework/gin"
	ginerror "module-clean/internal/framework/gin/errordefs"
	"module-clean/internal/member/interface_adapters/dto"
	"module-clean/internal/member/interface_adapters/mapper"
	presenterhttp "module-clean/internal/member/interface_adapters/presenter/http"
	"module-clean/internal/member/usecase"
	"net/http"
)

type MemberController struct {
	memberUseCase *usecase.MemberUseCase
}

func NewMemberController(memberUseCase *usecase.MemberUseCase) *MemberController {
	return &MemberController{
		memberUseCase: memberUseCase,
	}
}

func (c *MemberController) Register(ctx *gin.Context) {
	var reqDTO dto.CreateMemberRequestDTO
	if err := ctx.ShouldBindJSON(&reqDTO); err != nil {
		errCode, msg := ginerror.MapGinBindingError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errCode, msg)
		return
	}
	if err := reqDTO.Validate(); err != nil {
		errCode, msg := presenterhttp.MapValidationError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errCode, msg)
		return
	}
	entity := mapper.CreateDTOtoEntity(reqDTO)
	if err := c.memberUseCase.RegisterMember(ctx, entity); err != nil {
		errCode, msg := presenterhttp.MapMemberUseCaseError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusInternalServerError, errCode, msg)
		return
	}
	respDTO := presenterhttp.PresentCreateMemberDTO(entity)
	ginresponse.SuccessAPIResponse(ctx, respDTO)
}
func (c *MemberController) GetByID(ctx *gin.Context) {
	var reqDTO dto.GetMemberByIDRequestDTO
	if err := ctx.ShouldBindUri(&reqDTO); err != nil {
		errCode, msg := ginerror.MapGinBindingError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errCode, msg)
		return
	}
	if err := reqDTO.Validate(); err != nil {
		errCode, msg := presenterhttp.MapValidationError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errCode, msg)
		return
	}
	entity := mapper.GetMemberByIDDTOToEntity(reqDTO)
	member, err := c.memberUseCase.GetMemberByID(ctx, entity.ID)
	if err != nil {
		errCode, msg := presenterhttp.MapMemberUseCaseError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusNotFound, errCode, msg)
		return
	}
	respDTO := presenterhttp.PresentGetMemberByIDDTO(member)
	ginresponse.SuccessAPIResponse(ctx, respDTO)
}
func (c *MemberController) GetByEmail(ctx *gin.Context) {
	var reqDTO dto.GetMemberByEmailRequestDTO
	if err := ctx.ShouldBindQuery(&reqDTO); err != nil {
		errCode, msg := ginerror.MapGinBindingError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errCode, msg)
		return
	}
	if err := reqDTO.Validate(); err != nil {
		errCode, msg := presenterhttp.MapValidationError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errCode, msg)
		return
	}
	entity := mapper.GetMemberByEmailDTOToEntity(reqDTO)
	member, err := c.memberUseCase.GetMemberByEmail(ctx, entity.Email)
	if err != nil {
		errCode, msg := presenterhttp.MapMemberUseCaseError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusNotFound, errCode, msg)
		return
	}
	respDTO := presenterhttp.PresentGetMemberByEmailDTO(member)
	ginresponse.SuccessAPIResponse(ctx, respDTO)
}
func (c *MemberController) List(ctx *gin.Context) {
	var reqDTO dto.ListMemberRequestDTO
	if err := ctx.ShouldBindQuery(&reqDTO); err != nil {
		errCode, msg := ginerror.MapGinBindingError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errCode, msg)
		return
	}
	if err := reqDTO.Validate(); err != nil {
		errorCode, msg := presenterhttp.MapValidationError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errorCode, msg)
		return
	}
	pagination := mapper.ListMemberToPagination(reqDTO)
	members, total, err := c.memberUseCase.ListMembers(ctx, *pagination)
	if err != nil {
		errorCode, msg := presenterhttp.MapMemberUseCaseError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusInternalServerError, errorCode, msg)
		return
	}
	respDTO := presenterhttp.PresentListMemberDTO(members)
	ginresponse.SuccessAPIResponseWithMeta(ctx, respDTO, total, reqDTO.Page, pagination.Limit, pagination.Offset)
}
func (c *MemberController) Update(ctx *gin.Context) {
	var reqDTO dto.UpdateMemberRequestDTO
	if err := ctx.ShouldBindJSON(&reqDTO); err != nil {
		errCode, msg := ginerror.MapGinBindingError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errCode, msg)
		return
	}
	if err := reqDTO.Validate(); err != nil {
		errorCode, msg := presenterhttp.MapValidationError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errorCode, msg)
		return
	}
	inputModel := mapper.UpdateDTOToInputModel(reqDTO)
	member, err := c.memberUseCase.UpdateMember(ctx, inputModel)
	if err != nil {
		errorCode, msg := presenterhttp.MapMemberUseCaseError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusInternalServerError, errorCode, msg)
		return
	}
	respDTO := presenterhttp.PresentUpdateMemberDTO(member)
	ginresponse.SuccessAPIResponse(ctx, respDTO)
}
func (c *MemberController) Delete(ctx *gin.Context) {
	var reqDTO dto.DeleteMemberRequestDTO
	if err := ctx.ShouldBindUri(&reqDTO); err != nil {
		errCode, msg := ginerror.MapGinBindingError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errCode, msg)
		return
	}
	if err := reqDTO.Validate(); err != nil {
		errorCode, msg := presenterhttp.MapValidationError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errorCode, msg)
		return
	}
	entity := mapper.DeleteDTOToEntity(reqDTO)
	member, err := c.memberUseCase.DeleteMember(ctx, entity.ID)
	if err != nil {
		errorCode, msg := presenterhttp.MapMemberUseCaseError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusInternalServerError, errorCode, msg)
		return
	}
	respDTO := presenterhttp.PresentDeleteMemberDTO(member)
	ginresponse.SuccessAPIResponse(ctx, respDTO)
}
