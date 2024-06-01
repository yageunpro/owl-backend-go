package user

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const GoogleUserinfoUrl = "https://www.googleapis.com/oauth2/v1/userinfo"

type Info struct {
	OpenId        string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"verified_email"`
}

func GetUserInfo(accessToken string) (*Info, error) {
	res, err := http.Get(GoogleUserinfoUrl + "?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info, status: " + res.Status)
	}

	raw, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	info := new(Info)
	err = json.Unmarshal(raw, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}
