package util

import (
	"fmt"

	"github.com/invincibot/penn-spark-server/api/models"
)

func WrapString(k, v string) string {
	return fmt.Sprintf("\"%v\":\"%v\"", k, v)
}

func WrapUint(k string, v uint) string {
	return fmt.Sprintf("\"%v\":%v", k, v)
}

func WrapBool(k string, v bool) string {
	return fmt.Sprintf("\"%v\":%v", k, v)
}

func WrapUserRoles(userRoles []models.UserRole) string {
	s := ""
	for i, item := range userRoles {
		s += UserRoleToJSON(item)
		if i < len(userRoles)-1 {
			s += ","
		}
	}
	return fmt.Sprintf("\"user_roles\":[%v]", s)
}

func ParamsToJSON(params []string) string {
	jsonString := "{"
	for i, param := range params {
		if i > 0 {
			jsonString += ","
		}
		jsonString += param
	}
	return jsonString + "}"
}
