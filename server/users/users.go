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
	UUID              string
	Username          string
	Email             string
	EncryptedPassword string
	CreatedAt         time.Time
	role              string
	profile_picture   string
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

	rows, err := server.RunQuery(fetchUserQuery, params)
	if err != nil {
		return User{}, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	if len(rows) > 1 {
		fmt.Fprintln(os.Stderr, "Y'a plus d'un user avec le même email. C'est normal ça ?")
	}

	var newUser User

	for _, row := range rows {
		newUser = NewUser(row["user_uuid"].(string), row["username"].(string), row["email"].(string), "", row["created_at"].(time.Time), row["role"].(string), row["profile_picture"].(string))
	}
	return newUser, nil
}

func (u *User) ToMap() map[string]interface{} {
	usrMap := make(map[string]interface{}, 0)

	usrMap["uuid"] = u.UUID
	usrMap["username"] = u.Username
	usrMap["email"] = u.Email
	usrMap["password"] = u.EncryptedPassword
	usrMap["created_at"] = u.CreatedAt
	usrMap["role"] = u.role
	usrMap["profile_picture"] = u.profile_picture

	return usrMap
}

func (u *User) ToCookieValue() string {
	// A faire valider
	return u.Username + SEPARATOR +
		u.Email + SEPARATOR +
		u.role + SEPARATOR +
		u.profile_picture
}

// Enregistrer un user complet ( Register )

func RegisterUser(params map[string]interface{}) error {
	re := regexp.MustCompile(`(?i)<[^>]+>|(SELECT|UPDATE|DELETE|INSERT|DROP|FROM|COUNT|AS|WHERE|--)|^\s|^\s*$|<script.*?>.*?</script.*?>`)

	for _, value := range params {
		if re.FindAllString(value.(string), -1) != nil {
			return fmt.Errorf("injection detected")
		}
	}

	registerUserQuery := `INSERT INTO users (user_uuid, username, email, password, role, created_at, profile_picture )  VALUES (?, ?, ?, ?, ?, ?, ?)`
	var err error

	// Cryptage du mot de passe
	params["password"], err = HashPassword(params["password"].(string))
	if err != nil {
		return fmt.Errorf("password encryption: %s", err.Error())
	}

	_, err = server.RunQuery(registerUserQuery, params)

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
