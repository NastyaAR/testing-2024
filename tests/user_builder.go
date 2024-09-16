package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/pkg"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IUserBuilder interface {
	SetUid(uuid2 string)
	SetPassword()
	SetMail()
	SetRole()
	GetUser() domain.User
}

type NormalClientUserBuilder struct {
	uuid     uuid.UUID
	password string
	mail     string
	role     string
}

type NormalModeratorUserBuilder struct {
	uuid     uuid.UUID
	password string
	mail     string
	role     string
}

func (n *NormalClientUserBuilder) SetUid(uuid2 string) {
	n.uuid = uuid.MustParse(uuid2)
}

func (n *NormalClientUserBuilder) SetPassword() {
	encr, _ := pkg.EncryptPassword("password", &zap.Logger{})
	n.password = encr
}

func (n *NormalClientUserBuilder) SetMail() {
	n.mail = "test@mail.ru"
}

func (n *NormalClientUserBuilder) SetRole() {
	n.role = domain.Client
}

func (n *NormalClientUserBuilder) GetUser() domain.User {
	return domain.User{
		UserID:   n.uuid,
		Mail:     n.mail,
		Password: n.password,
		Role:     n.role,
	}
}

func (n *NormalModeratorUserBuilder) SetUid(uuid2 string) {
	n.uuid = uuid.MustParse(uuid2)
}

func (n *NormalModeratorUserBuilder) SetPassword() {
	encr, _ := pkg.EncryptPassword("password", &zap.Logger{})
	n.password = encr
}

func (n *NormalModeratorUserBuilder) SetMail() {
	n.mail = "test@mail.ru"
}

func (n *NormalModeratorUserBuilder) SetRole() {
	n.role = domain.Moderator
}

func (n *NormalModeratorUserBuilder) GetUser() domain.User {
	return domain.User{
		UserID:   n.uuid,
		Mail:     n.mail,
		Password: n.password,
		Role:     n.role,
	}
}
