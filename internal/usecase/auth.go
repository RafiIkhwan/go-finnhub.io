package usecase

import (
    "context"
    "time"

    "github.com/evrone/go-clean-template/internal/entity"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
    Create(ctx context.Context, user entity.User) error
    GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type AuthClaims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}

type Auth struct {
    repo      UserRepo
    jwtSecret string
}

func NewAuth(repo UserRepo, jwtSecret string) *Auth {
    return &Auth{repo: repo, jwtSecret: jwtSecret}
}

func (a *Auth) Register(ctx context.Context, email, password string) error {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user := entity.User{Email: email, Password: string(hashed)}
		// Validasi user
    // if err := user.Validate(); err != nil {
    //     return err
    // }
    return a.repo.Create(ctx, user)
}

func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
    user, err := a.repo.GetByEmail(ctx, email)
    if err != nil {
        return "", err
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", err
    }
    claims := AuthClaims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(a.jwtSecret))
}