package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-redis/redis"
)

var (
	Errusername = errors.New("Incorrect username")
	Errpassword = errors.New("Incorrect password")
)

func LoginUser(username, password string) error {
	hash, err := client.Get("user:" + username).Bytes()
	if err == redis.Nil{
		return Errusername
	} else if err != nil{
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil{
		return Errpassword
	}
	return nil
}

func RegisterUser(username, password string) error {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil{
		return err
	}
	return client.Set("user:" + username, hash, 0).Err()
}