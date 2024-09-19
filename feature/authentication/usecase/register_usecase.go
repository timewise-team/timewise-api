package usecase

import (
	"api/dms"
	"encoding/json"
	"errors"
	"fmt"
	user_register_dto "github.com/timewise-team/timewise-models/dtos/core_dtos/user_register_dtos"
	"github.com/timewise-team/timewise-models/models"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func RegisterUser(RegisterRequestDto user_register_dto.RegisterRequestDto) error {
	// Check if passwords match
	if RegisterRequestDto.Password != RegisterRequestDto.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(RegisterRequestDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error hashing password")
	}

	parts := strings.Fields(RegisterRequestDto.FullName)
	lastName := parts[len(parts)-1]
	firstName := strings.Join(parts[:len(parts)-1], " ")

	newUser := models.TwUser{
		Username:     RegisterRequestDto.UserName,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        RegisterRequestDto.Email,
		PasswordHash: string(hashedPassword),
		LastLoginAt:  time.Now(),
	}

	resp, err := dms.CallAPI(
		"POST",
		"/user",
		newUser,
		nil,
		nil,
		120,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("can not read response body")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	err = json.Unmarshal([]byte(body), &newUser)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		return errors.New("error parsing JSON")
	}

	// Create new record in user_email
	createNewUserEmailRequest := models.TwUserEmail{
		UserId: newUser.ID,
		Email:  RegisterRequestDto.Email,
	}
	resp, err = dms.CallAPI(
		"POST",
		"/user_email",
		createNewUserEmailRequest,
		nil,
		nil,
		120,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("can not read response body")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	return nil
}
