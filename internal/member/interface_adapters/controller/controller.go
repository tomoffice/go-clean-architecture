package controller

import (
	"module-clean/internal/member/interface_adapters/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"module-clean/internal/member/usecase"
)

type MemberController struct {
	memberUseCase *member.UseCase
}

func NewMemberController(memberUseCase *member.UseCase) *MemberController {
	return &MemberController{
		memberUseCase: memberUseCase,
	}
}

// Register handles POST /members
func (c *MemberController) Register(ctx *gin.Context) {
	var req dto.CreateMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "detail": err.Error()})
		return
	}
	entity := req.ToEntity()
	err := c.memberUseCase.RegisterMember(ctx, entity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register member", "detail": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "member registered successfully"})
}
