package controller

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/modules/member/usecase/input_port"
	"module-clean/internal/modules/member/usecase/output_port"
	"reflect"
	"testing"
)

func TestMemberController_Delete(t *testing.T) {
	type fields struct {
		useCase   input_port.MemberInputPort
		presenter output_port.MemberPresenter
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
				useCase:   tt.fields.useCase,
				presenter: tt.fields.presenter,
			}
			c.Delete(tt.args.ctx)
		})
	}
}

func TestMemberController_GetByEmail(t *testing.T) {
	type fields struct {
		useCase   input_port.MemberInputPort
		presenter output_port.MemberPresenter
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
				useCase:   tt.fields.useCase,
				presenter: tt.fields.presenter,
			}
			c.GetByEmail(tt.args.ctx)
		})
	}
}

func TestMemberController_GetByID(t *testing.T) {
	type fields struct {
		useCase   input_port.MemberInputPort
		presenter output_port.MemberPresenter
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
				useCase:   tt.fields.useCase,
				presenter: tt.fields.presenter,
			}
			c.GetByID(tt.args.ctx)
		})
	}
}

func TestMemberController_List(t *testing.T) {
	type fields struct {
		useCase   input_port.MemberInputPort
		presenter output_port.MemberPresenter
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
				useCase:   tt.fields.useCase,
				presenter: tt.fields.presenter,
			}
			c.List(tt.args.ctx)
		})
	}
}

func TestMemberController_Register(t *testing.T) {
	type fields struct {
		useCase   input_port.MemberInputPort
		presenter output_port.MemberPresenter
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
				useCase:   tt.fields.useCase,
				presenter: tt.fields.presenter,
			}
			c.Register(tt.args.ctx)
		})
	}
}

func TestMemberController_Update(t *testing.T) {
	type fields struct {
		useCase   input_port.MemberInputPort
		presenter output_port.MemberPresenter
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
				useCase:   tt.fields.useCase,
				presenter: tt.fields.presenter,
			}
			c.Update(tt.args.ctx)
		})
	}
}

func TestNewMemberController(t *testing.T) {
	type args struct {
		memberUseCase input_port.MemberInputPort
		presenter     output_port.MemberPresenter
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

func TestMemberController_Delete1(t *testing.T) {
	type fields struct {
		useCase   input_port.MemberInputPort
		presenter output_port.MemberPresenter
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
				useCase:   tt.fields.useCase,
				presenter: tt.fields.presenter,
			}
			c.Delete(tt.args.ctx)
		})
	}
}

func TestMemberController_GetByEmail1(t *testing.T) {
	type fields struct {
		useCase   input_port.MemberInputPort
		presenter output_port.MemberPresenter
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
				useCase:   tt.fields.useCase,
				presenter: tt.fields.presenter,
			}
			c.GetByEmail(tt.args.ctx)
		})
	}
}

func TestMemberController_GetByID1(t *testing.T) {
	type fields struct {
		useCase   input_port.MemberInputPort
		presenter output_port.MemberPresenter
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
				useCase:   tt.fields.useCase,
				presenter: tt.fields.presenter,
			}
			c.GetByID(tt.args.ctx)
		})
	}
}

func TestMemberController_List1(t *testing.T) {
	type fields struct {
		useCase   input_port.MemberInputPort
		presenter output_port.MemberPresenter
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
				useCase:   tt.fields.useCase,
				presenter: tt.fields.presenter,
			}
			c.List(tt.args.ctx)
		})
	}
}

func TestMemberController_Register1(t *testing.T) {
	type fields struct {
		useCase   input_port.MemberInputPort
		presenter output_port.MemberPresenter
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
				useCase:   tt.fields.useCase,
				presenter: tt.fields.presenter,
			}
			c.Register(tt.args.ctx)
		})
	}
}

func TestMemberController_Update1(t *testing.T) {
	type fields struct {
		useCase   input_port.MemberInputPort
		presenter output_port.MemberPresenter
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
				useCase:   tt.fields.useCase,
				presenter: tt.fields.presenter,
			}
			c.Update(tt.args.ctx)
		})
	}
}

func TestNewMemberController1(t *testing.T) {
	type args struct {
		memberUseCase input_port.MemberInputPort
		presenter     output_port.MemberPresenter
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