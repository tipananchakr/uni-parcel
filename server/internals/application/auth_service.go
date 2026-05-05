package application

import (
	"context"
	"strings"
	"time"

	"github.com/tipananchakr/uni-parcel/internals/core/domain"
	"github.com/tipananchakr/uni-parcel/internals/core/ports"
)

type AuthResult struct {
	User  domain.User `json:"user"`
	Token string      `json:"token"`
}

type AuthService struct {
	users   ports.UserRepository
	timeout time.Duration
	hasher  ports.PasswordHasher
	token   ports.TokenManager
}

func NewAuthService(users ports.UserRepository, hasher ports.PasswordHasher, token ports.TokenManager) *AuthService {
	return &AuthService{
		users:   users,
		timeout: 5 * time.Second,
		hasher:  hasher,
		token:   token,
	}
}

// ตอนนี้ service สามารถที่จะเชื่อไปที่ repository ได้แล้วโดยจะเรียกใช้ repository โดยมี method receiver
func (a *AuthService) Register(ctx context.Context, user domain.User) (AuthResult, error) {
	// ส่ง context ของ request register เพื่อ watch ถ้าเกิน timeout ให้ cancel ทิ้งไม่ค้างไว้
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	// Normalize email และ Validation password
	email := normalizeEmail(user.Email)
	password := user.Password
	if email == "" || len(password) < 6 {
		return AuthResult{}, domain.ErrInvalidCredentials
	}

	// Hash password
	passwordHask, err := a.hasher.Hash(password)
	if err != nil {
		return AuthResult{}, err
	}

	// สร้าง User ใหม่ใน database ผ่าน repository CreateUser เป็น interface ที่เราได้กำหนดไว้ใน ports และเราจะเรียกใช้ผ่าน service โดยส่งข้อมูล user ที่เราต้องการจะสร้างไปให้ repository ทำการสร้างให้
	newUser, err := a.users.CreateUser(ctx, domain.User{
		Email:    email,
		Password: passwordHask,
		Name:     user.Name,
		Role:     user.Role,
	})

	if err != nil {
		return AuthResult{}, err
	}

	// สร้าง token ใหม่สำหรับ user ที่เพิ่ง register เสร็จ
	newToken, err := a.token.Generate(newUser.ID.Hex())
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: newUser, Token: newToken}, nil
}

func (a *AuthService) Login(ctx context.Context, email string, password string) (AuthResult, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	// Normalize email และ Validation password
	email = normalizeEmail(email)
	if email == "" || len(password) < 6 {
		return AuthResult{}, domain.ErrInvalidCredentials
	}

	// ค้นหา user ใน database ผ่าน repository GetUserByEmail เป็น interface ที่เราได้กำหนดไว้ใน ports และเราจะเรียกใช้ผ่าน service โดยส่ง email ไปให้ repository ทำการค้นหาให้
	user, err := a.users.GetUserByEmail(ctx, email)
	if err != nil {
		return AuthResult{}, domain.ErrInvalidCredentials
	}

	// เปรียบเทียบ password ที่ user กรอกเข้ามากับ password ที่เก็บไว้ใน database ผ่าน hasher Compare เป็น interface ที่เราได้กำหนดไว้ใน ports และเราจะเรียกใช้ผ่าน service โดยส่ง password ที่ user กรอกเข้ามาและ password hash ที่เก็บไว้ใน database ไปให้ hasher ทำการเปรียบเทียบให้
	err = a.hasher.Compare(user.Password, password)
	if err != nil {
		return AuthResult{}, domain.ErrInvalidCredentials
	}

	// สร้าง token ใหม่สำหรับ user ที่เพิ่ง login เสร็จ
	token, err := a.token.Generate(user.ID.Hex())
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user, Token: token}, nil
}

func (a *AuthService) CurrentUser(ctx context.Context, token string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	// Validate token และดึง userID ออกมา
	userID, err := a.token.Validate(token)
	if err != nil {
		return domain.User{}, domain.ErrInvalidToken
	}

	return a.users.GetUserByID(ctx, userID)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
