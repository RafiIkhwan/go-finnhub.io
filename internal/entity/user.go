package entity

import "time"

type User struct {
    ID        int       `json:"id"`
    Email     string    `json:"email" validate:"required,email"`
    Password  string    `json:"password" validate:"required,min=8"`  // Hashed
    CreatedAt time.Time `json:"created_at"`
}

// func (u *User) Validate() error {
//     return validate.Struct(u) 
// }