package main

import "github.com/denisbrodbeck/machineid"

func GetHostId() (string, error) {
	id, err := machineid.ID()
	if err != nil {
		return "", err
	}
	h := GetSHA256([]byte(id))
	return h, nil
}
