package users

import (
	"fmt"
	"forum/server"
	"os"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	SEPARATOR = "|"
)

type User struct {
	UUID              string    `json:"user_uuid"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	EncryptedPassword string    `json:"password"`
	CreatedAt         time.Time `json:"created_at"`
	Role              string    `json:"role"`
	ProfilePicture    string    `json:"profile_picture"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func NewUser(uuid, username, email, encryptedPassword string, createdAt time.Time, role string, profile_picture string) User {
	newUser := User{uuid, username, email, encryptedPassword, createdAt, role, profile_picture}
	return newUser
}

// Trouver un utilisateur par son email et le renvoyer ( pour login )

func FetchUserByEmail(email string) (User, error) {
	re := regexp.MustCompile(`(?i)<[^>]+>|(SELECT|UPDATE|DELETE|INSERT|DROP|FROM|COUNT|AS|WHERE|--)|^\s|^\s*$|<script.*?>.*?</script.*?>`)

	if re.FindAllString(email, -1) != nil {
		return User{}, fmt.Errorf("injection detected")
	}

	fetchUserQuery := `SELECT * FROM users WHERE email= ?`
	params := []interface{}{email}

	rows, err := server.RunQuery(fetchUserQuery, params...)
	if err != nil {
		return User{}, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	if len(rows) > 1 {
		fmt.Fprintln(os.Stderr, "Y'a plus d'un user avec le même email. C'est normal ça ?")
	} else if len(rows) == 0 {
		return User{}, nil
	}

	newUser := User{}
	result := rows[0]

	// Utiliser des assertions de type avec vérification de valeur nulle
	if v, ok := result["user_uuid"]; ok && v != nil {
		newUser.UUID = v.(string)
	}
	if v, ok := result["username"]; ok && v != nil {
		newUser.Username = v.(string)
	}
	if v, ok := result["password"]; ok && v != nil {
		newUser.EncryptedPassword = v.(string)
	}
	if v, ok := result["profile_picture"]; ok && v != nil {
		newUser.ProfilePicture = v.(string)
	}
	if v, ok := result["email"]; ok && v != nil {
		newUser.Email = v.(string)
	}
	if v, ok := result["role"]; ok && v != nil {
		newUser.Role = v.(string)
	}
	if v, ok := result["created_at"]; ok && v != nil {
		parsedTime, err := time.Parse("2006-01-02", result["created_at"].(string))
		if err != nil {
			fmt.Fprintln(os.Stderr, "dommage")
		}
		newUser.CreatedAt = parsedTime
	}

	return newUser, nil
}

func (u *User) ToMap() map[string]interface{} {
	usrMap := make(map[string]interface{}, 0)

	usrMap["user_uuid"] = u.UUID
	usrMap["username"] = u.Username
	usrMap["email"] = u.Email
	usrMap["password"] = u.EncryptedPassword
	usrMap["created_at"] = u.CreatedAt.Format("2006-01-02")
	usrMap["role"] = u.Role
	usrMap["profile_picture"] = u.ProfilePicture

	return usrMap
}

func (u *User) ToCookieValue() string {
	return u.UUID + SEPARATOR +
		u.Username + SEPARATOR +
		u.Email + SEPARATOR +
		u.Role + SEPARATOR +
		u.ProfilePicture
}

// Enregistrer un user complet ( Register )

func RegisterUser(params map[string]interface{}) error {

	profile_picture, _ := params["profile_picture"].(string)

	re := regexp.MustCompile(`(?i)<[^>]+>|(SELECT|UPDATE|DELETE|INSERT|DROP|FROM|COUNT|AS|WHERE|--)|^\s|^\s*$|<script.*?>.*?</script.*?>`)

	for key, value := range params {
		if (key == "username" || key == "email" || key == "password") && re.FindAllString(value.(string), -1) != nil {
			return fmt.Errorf("injection detected")
		}
	}
	if profile_picture == "" {
		profile_picture = RandomProfilPicture()
	}

	registerUserQuery := `INSERT INTO users (user_uuid, username, email, password, role, created_at, profile_picture )  VALUES (?, ?, ?, ?, ?, ?, ?)`
	var err error

	_, err = server.RunQuery(registerUserQuery, params["user_uuid"], params["username"], params["email"], params["password"], params["role"], params["created_at"], profile_picture)

	if err != nil {
		return err
	}

	return nil
}

// Mettre à jour les valeurs d'un utilisateur ( Update )
func (u *User) UpdateUser(params map[string]interface{}) error {
	re := regexp.MustCompile(`(?i)<[^>]+>|(SELECT|UPDATE|DELETE|INSERT|DROP|FROM|COUNT|AS|WHERE|--)|^\s|^\s*$|<script.*?>.*?</script.*?>`)

	for _, value := range params {
		if re.FindAllString(value.(string), -1) != nil {
			return fmt.Errorf("injection detected")
		}
	}

	updateUserQuery := `UPDATE users SET username = ?, email = ?, password = ?, role = ?, profile_picture = ? WHERE user_uuid = ?`
	_, err := server.RunQuery(updateUserQuery, params)

	if err != nil {
		return err
	}

	return nil
}
