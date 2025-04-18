package controller

import (
	"fmt"
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
	memberUseCase usecase.UseCase
}

func NewMemberController(memberUseCase usecase.UseCase) *MemberController {
	return &MemberController{
		memberUseCase: memberUseCase,
	}
}

// Register handles POST /members
func (c *MemberController) Register(ctx *gin.Context) {
	var reqDTO dto.CreateMemberRequest
	if err := ctx.ShouldBindJSON(&reqDTO); err != nil {
		code, msg := presenterhttp.MapErrorToCode(err)
		ctx.JSON(http.StatusBadRequest, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(code), Message: msg},
		})
		return
	}
	if err := reqDTO.Validate(); err != nil {
		code, msg := presenterhttp.MapErrorToCode(err)
		ctx.JSON(http.StatusBadRequest, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(code), Message: msg},
		})
		return
	}

	entity := mapper.CreateDTOtoEntity(reqDTO)
	if err := c.memberUseCase.RegisterMember(ctx, entity); err != nil {
		code, msg := presenterhttp.MapErrorToCode(err)
		ctx.JSON(http.StatusInternalServerError, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(code), Message: msg},
		})
		return
	}

	ctx.JSON(http.StatusCreated, response.APIResponse[any]{
		Data: &map[string]string{"message": "member registered successfully"},
	})
}
func (c *MemberController) GetByID(ctx *gin.Context) {
	var reqDTO dto.GetMemberByIDRequest
	if err := ctx.ShouldBindUri(&reqDTO); err != nil || reqDTO.Validate() != nil {
		code, msg := presenterhttp.MapErrorToCode(err)
		ctx.JSON(http.StatusBadRequest, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(code), Message: msg},
		})
		return
	}

	entity := mapper.GetMemberByIDToEntity(reqDTO)
	member, err := c.memberUseCase.GetMemberByID(ctx, entity.ID)
	if err != nil {
		code, msg := presenterhttp.MapErrorToCode(err)
		ctx.JSON(http.StatusNotFound, response.APIResponse[any]{
			Error: &response.ErrorPayload{Code: strconv.Itoa(code), Message: msg},
		})
		return
	}

	dto := presenterhttp.ToMemberResponse(member)
	ctx.JSON(http.StatusOK, response.APIResponse[dto.MemberResponseDTO]{Data: &dto})
}
func (c *MemberController) GetByEmail(ctx *gin.Context) {
	var reqDTO dto.GetMemberByEmailRequest
	if err := ctx.ShouldBindQuery(&reqDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "detail": err.Error()})
		return
	}
	entity := mapper.GetMemberByEmailToEntity(reqDTO)
	member, err := c.memberUseCase.GetMemberByEmail(ctx, entity.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get member", "detail": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, member)
}
func (c *MemberController) List(ctx *gin.Context) {
	var reqDTO dto.ListMemberRequest
	if err := ctx.ShouldBindQuery(&reqDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "detail": err.Error()})
		return
	}
	pagination := mapper.ListMemberToPagination(reqDTO)
	members, err := c.memberUseCase.ListMembers(ctx, *pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list members", "detail": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, members)
}
func (c *MemberController) Update(ctx *gin.Context) {
	var reqDTO dto.UpdateMemberRequest
	if err := ctx.ShouldBindJSON(&reqDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "detail": err.Error()})
		return
	}
	entity := mapper.UpdateDTOToInputModel(reqDTO)
	err := c.memberUseCase.UpdateMember(ctx, entity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update member", "detail": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "member updated successfully"})
}
