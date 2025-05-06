package controller

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/modules/member/adapter/dto"
	"module-clean/internal/modules/member/adapter/mapper"
	http2 "module-clean/internal/modules/member/adapter/presenter/http"
	"module-clean/internal/modules/member/usecase"
	ginerror "module-clean/internal/platform/gin/errordefs"
	ginresponse "module-clean/internal/platform/gin/response"
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
		errCode, msg := http2.MapValidationError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errCode, msg)
		return
	}
	entity := mapper.CreateDTOtoEntity(reqDTO)
	if err := c.memberUseCase.RegisterMember(ctx, entity); err != nil {
		errCode, msg := http2.MapMemberUseCaseError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusInternalServerError, errCode, msg)
		return
	}
	respDTO := http2.PresentCreateMemberDTO(entity)
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
		errCode, msg := http2.MapValidationError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errCode, msg)
		return
	}
	entity := mapper.GetMemberByIDDTOToEntity(reqDTO)
	member, err := c.memberUseCase.GetMemberByID(ctx, entity.ID)
	if err != nil {
		errCode, msg := http2.MapMemberUseCaseError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusNotFound, errCode, msg)
		return
	}
	respDTO := http2.PresentGetMemberByIDDTO(member)
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
		errCode, msg := http2.MapValidationError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errCode, msg)
		return
	}
	entity := mapper.GetMemberByEmailDTOToEntity(reqDTO)
	member, err := c.memberUseCase.GetMemberByEmail(ctx, entity.Email)
	if err != nil {
		errCode, msg := http2.MapMemberUseCaseError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusNotFound, errCode, msg)
		return
	}
	respDTO := http2.PresentGetMemberByEmailDTO(member)
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
		errorCode, msg := http2.MapValidationError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errorCode, msg)
		return
	}
	pagination := mapper.ListMemberToPagination(reqDTO)
	members, total, err := c.memberUseCase.ListMembers(ctx, *pagination)
	if err != nil {
		errorCode, msg := http2.MapMemberUseCaseError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusInternalServerError, errorCode, msg)
		return
	}
	respDTO := http2.PresentListMemberDTO(members)
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
		errorCode, msg := http2.MapValidationError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errorCode, msg)
		return
	}
	inputModel := mapper.UpdateDTOToInputModel(reqDTO)
	member, err := c.memberUseCase.UpdateMember(ctx, inputModel)
	if err != nil {
		errorCode, msg := http2.MapMemberUseCaseError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusInternalServerError, errorCode, msg)
		return
	}
	respDTO := http2.PresentUpdateMemberDTO(member)
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
		errorCode, msg := http2.MapValidationError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusBadRequest, errorCode, msg)
		return
	}
	entity := mapper.DeleteDTOToEntity(reqDTO)
	member, err := c.memberUseCase.DeleteMember(ctx, entity.ID)
	if err != nil {
		errorCode, msg := http2.MapMemberUseCaseError(err)
		ginresponse.FailureAPIResponse(ctx, http.StatusInternalServerError, errorCode, msg)
		return
	}
	respDTO := http2.PresentDeleteMemberDTO(member)
	ginresponse.SuccessAPIResponse(ctx, respDTO)
}
