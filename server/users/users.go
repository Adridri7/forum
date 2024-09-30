package users

import (
	"fmt"
	"forum/server"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type User struct {
	UUID              string
	Username          string
	Email             string
	EncryptedPassword string
	CreatedAt         time.Time
	role              string
	profile_picture   string
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

	var newUser User

	for _, row := range rows {
		newUser = NewUser(row["user_uuid"].(string), row["username"].(string), row["email"].(string), "", row["created_at"].(time.Time), row["role"].(string), row["profile_picture"].(string))
	}
	return newUser, nil
}

// Enregistrer un user complet ( Register )

func RegisterUser(params map[string]interface{}) error {
	// Ajouter la logique de regexp pour les injections ( à faire -> Esteban )
	registerUserQuery := `INSERT INTO users (user_uuid, username, email, password, role, created_at, profile_picture )  VALUES (?, ?, ?, ?, ?, ?, ?)`

	// Ajouter la logique de cryptage du mot de passe ( à faire -> Esteban )

	_, err := server.RunQuery(registerUserQuery, params)

	if err != nil {
		return err
	}

	return nil
}

// Mettre à jour les valeurs d'un utilisateur ( Update )

func (u User) UpdateUser(params map[string]interface{}) error {
	// Ajouter la logique de regexp pour les injections ( à faire -> Esteban )
	updateUserQuery := `UPDATE users SET username = ?, email = ?, password = ?, role = ?, profile_picture = ? WHERE user_uuid = ?`
	_, err := server.RunQuery(updateUserQuery, params)

	if err != nil {
		return err
	}

	return nil
}
