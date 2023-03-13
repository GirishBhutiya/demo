package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/GirishBhutiya/demo/backend-service/data"
)

type CreatUserRequests struct {
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func newUserResponse(user data.User) userResponse {
	return userResponse{
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.UpdatedAt,
		CreatedAt:         user.CreatedAt,
	}
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}
type userResponse struct {
	FullName          string    `json:"full_name" binding:"required"`
	Email             string    `json:"email" binding:"required,email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	//validate the user against the database
	user, err := app.Repo.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := app.Repo.PasswordMatches(requestPayload.Password, *user)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	accessToken, err := app.tokenMaker.CreateToken(
		user.Email,
		time.Minute*1,
	)
	if err != nil {
		app.errorJSON(w, errors.New("internal Server error"), http.StatusInternalServerError)
		return
	}
	usrRsp := &userResponse{
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.UpdatedAt,
		CreatedAt:         user.CreatedAt,
	}
	loginRsp := &loginUserResponse{
		User:        *usrRsp,
		AccessToken: accessToken,
	}
	/* payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	} */

	app.writeJSON(w, http.StatusAccepted, loginRsp)

}

func (app *Config) CreateUser(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	u := data.User{
		Email:     requestPayload.Email,
		FullName:  requestPayload.FullName,
		Password:  requestPayload.Password,
		Active:    0,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	//validate the user against the database
	usr, err := app.Repo.Insert(u)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	res := newUserResponse(usr)

	/* payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	} */

	app.writeJSON(w, http.StatusAccepted, res)

}
func (app *Config) GetAllPost(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		AccessToken string `json:"access_token"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	_, err = app.tokenMaker.VerifyToken(requestPayload.AccessToken)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	posts, err := app.Repo.GetAllPost()
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusAccepted, posts)

}
