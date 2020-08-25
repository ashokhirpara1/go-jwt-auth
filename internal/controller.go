package internal

import (
	"go-jwt/database"
	"go-jwt/database/models"
	"go-jwt/logging"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Controller is the big boss
type Controller struct {
	db  database.Storager
	Log logging.Handler
}

// Get - return new Controller
func Get(db database.Storager, log logging.Handler) *Controller {
	return &Controller{
		db:  db,
		Log: log,
	}
}

// GenerateToken -
func (c *Controller) GenerateToken(user models.User) (string, error) {
	// Create a token for the user
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	tokenRef := &models.Token{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenRef)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Register - user registration
func (c *Controller) Register(user models.User) (interface{}, error) {
	defer c.Log.Exit(c.Log.Enter())

	createdUser, err := c.db.Register(user)

	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

// Login - user login
func (c *Controller) Login(email, password string) (interface{}, error) {
	defer c.Log.Exit(c.Log.Enter())

	user, err := c.db.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return nil, errf
	}

	tokenString, err := c.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	var resp = make(map[string]interface{})
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = user

	return resp, nil
}

// GetUsers -
func (c *Controller) GetUsers() (interface{}, error) {

	users, err := c.db.GetUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}
