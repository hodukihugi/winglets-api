package utils

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/nullbio/null.v4"
	"regexp"
	"strconv"
	"strings"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

// PascalToSnake eg: given input as "UserID", output will be "user_id"
func PascalToSnake(input string) string {
	// Use a regular expression to match capital letters and insert an underscore before them
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(input, "${1}_${2}")

	// Convert the result to lowercase
	snake = strings.ToLower(snake)

	return snake
}

func EmploymentIDToOrgIDUserEmail(empID string) (uint, string, error) {
	arr := strings.Split(empID, "_")
	if len(arr) != 2 {
		return 0, "", errors.New("invalid employment id")
	}
	orgID, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0, "", err
	}
	return uint(orgID), arr[1], nil
}

func NullUintToUint(u null.Uint) uint {
	if u.Valid {
		return u.Uint
	}
	return 0
}

// CoordinatesStringToPairFloat64 convert coordinates (longitude, latitude) string to pair float64.
func CoordinatesStringToPairFloat64(coordinates string) (float64, float64, error) {
	coords := strings.Split(coordinates, ",")

	longitude, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return 0, 0, err
	}

	latitude, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return 0, 0, err
	}

	return longitude, latitude, nil
}
