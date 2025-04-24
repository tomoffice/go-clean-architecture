package controller

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/member/interface_adapters/dto"
	"module-clean/internal/member/interface_adapters/mapper"
	presenterhttp "module-clean/internal/member/interface_adapters/presenter/http"
	"module-clean/internal/member/usecase"
	"module-clean/internal/shared/response"
	"net/http"
	"strconv"
)

type MemberController struct {
	memberUseCase *usecase.MemberUseCase
}

func NewMemberController(memberUseCase *usecase.MemberUseCase) *MemberController {
	return &MemberController{
		memberUseCase: memberUseCase,
	}
}

// Register handles POST /members
func (c *MemberController) Register(ctx *gin.Context) {
	var reqDTO dto.CreateMemberRequestDTO
	if err := ctx.ShouldBindJSON(&reqDTO); err != nil {
		errCode, msg := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusBadRequest, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errCode), Message: msg},
		})
		return
	}
	if err := reqDTO.Validate(); err != nil {
		errCode, msg := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusBadRequest, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errCode), Message: msg},
		})
		return
	}

	entity := mapper.CreateDTOtoEntity(reqDTO)
	if err := c.memberUseCase.RegisterMember(ctx, entity); err != nil {
		errCode, msg := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusInternalServerError, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errCode), Message: msg},
		})
		return
	}
	respDTO := presenterhttp.PresentCreateMemberDTO(entity)
	ctx.JSON(http.StatusCreated, response.APIResponse[any]{
		Data: respDTO,
	})
}
func (c *MemberController) GetByID(ctx *gin.Context) {
	var reqDTO dto.GetMemberByIDRequestDTO
	if err := ctx.ShouldBindUri(&reqDTO); err != nil || reqDTO.Validate() != nil {
		errCode, msg := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusBadRequest, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errCode), Message: msg},
		})
		return
	}
	entity := mapper.GetMemberByIDDTOToEntity(reqDTO)
	member, err := c.memberUseCase.GetMemberByID(ctx, entity.ID)
	if err != nil {
		errCode, msg := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusNotFound, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errCode), Message: msg},
		})
		return
	}
	respDTO := presenterhttp.PresentGetMemberByIDDTO(member)
	ctx.JSON(http.StatusOK, response.APIResponse[any]{
		Data: respDTO,
	})
}
func (c *MemberController) GetByEmail(ctx *gin.Context) {
	var reqDTO dto.GetMemberByEmailRequestDTO
	if err := ctx.ShouldBindQuery(&reqDTO); err != nil {
		errCode, mag := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusBadRequest, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errCode), Message: mag},
		})
		return
	}
	entity := mapper.GetMemberByEmailDTOToEntity(reqDTO)
	member, err := c.memberUseCase.GetMemberByEmail(ctx, entity.Email)
	if err != nil {
		errCode, msg := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusNotFound, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errCode), Message: msg},
		})
		return
	}
	respDTO := presenterhttp.PresentGetMemberByEmailDTO(member)
	ctx.JSON(http.StatusOK, response.APIResponse[any]{
		Data: respDTO,
	})
}
func (c *MemberController) List(ctx *gin.Context) {
	var reqDTO dto.ListMemberRequestDTO
	if err := ctx.ShouldBindQuery(&reqDTO); err != nil {
		errorCode, msg := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusBadRequest, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errorCode), Message: msg},
		})
		return
	}
	pagination := mapper.ListMemberToPagination(reqDTO)
	members, err := c.memberUseCase.ListMembers(ctx, *pagination)
	if err != nil {
		errorCode, msg := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusInternalServerError, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errorCode), Message: msg},
		})
		return
	}
	respDTO := presenterhttp.PresentListMemberDTO(members)
	ctx.JSON(http.StatusOK, response.APIResponse[any]{
		Data: respDTO,
	})
}
func (c *MemberController) Update(ctx *gin.Context) {
	var reqDTO dto.UpdateMemberRequestDTO
	if err := ctx.ShouldBindJSON(&reqDTO); err != nil {
		errorCode, msg := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusBadRequest, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errorCode), Message: msg},
		})
		return
	}
	inputModel := mapper.UpdateDTOToInputModel(reqDTO)
	member, err := c.memberUseCase.UpdateMember(ctx, inputModel)
	if err != nil {
		errorCode, msg := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusInternalServerError, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errorCode), Message: msg},
		})
		return
	}
	respDTO := presenterhttp.PresentUpdateMemberDTO(member)
	ctx.JSON(http.StatusOK, response.APIResponse[any]{
		Data: respDTO,
	})
}
func (c *MemberController) Delete(ctx *gin.Context) {
	var reqDTO dto.DeleteMemberRequestDTO
	if err := ctx.ShouldBindUri(&reqDTO); err != nil {
		errorCode, msg := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusBadRequest, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errorCode), Message: msg},
		})
		return
	}
	entity := mapper.DeleteDTOToEntity(reqDTO)
	member, err := c.memberUseCase.DeleteMember(ctx, entity.ID)
	if err != nil {
		errorCode, msg := presenterhttp.MapErrorToResponse(err)
		ctx.JSON(http.StatusInternalServerError, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(errorCode), Message: msg},
		})
		return
	}
	respDTO := presenterhttp.PresentDeleteMemberDTO(member)
	ctx.JSON(http.StatusOK, response.APIResponse[any]{
		Data: respDTO,
	})
}
