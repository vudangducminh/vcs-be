package posgresql_query

import (
	"log"
	"net/http"
	"sms/object"
	"sms/server/database/postgresql/connector"
)

func GetAccountPasswordByUsername(username string) string {
	var account object.Account
	has, err := connector.Engine.Table("account").Cols("password").Alias("account").
		Where("username = ?", username).
		Get(&account)

	if err != nil {
		log.Println("Error retrieving account password:", err)
		return ""
	}

	if !has || account.Password == "" {
		log.Println("No account found with username:", username)
		return ""
	}
	return account.Password
}

func AddAccountInfo(account object.Account) int {
	has, err := connector.Engine.Table("account").
		Where("username = ?", account.Username).Count(new(object.Account))
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
	if has > 0 {
		log.Println("Account already exists with username:", account.Username)
		return http.StatusConflict
	}
	affected, err := connector.Engine.Insert(account)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
	if affected > 0 {
		log.Println("Account added successfully:", account.Username)
		return http.StatusCreated
	} else {
		log.Println("Failed to add account:", account.Username)
		return http.StatusInternalServerError
	}
}

func CheckServerExists(IPv4 string) bool {
	var server object.Server
	has, err := connector.Engine.Table("server").Where("IPv4 = ?", IPv4).Get(&server)
	if err != nil {
		log.Println("Error checking if server exists:", err)
		return true
	}
	if has {
		log.Println("Server already exists with IPV4:", IPv4)
		return true
	}
	return false
}

func AddServerInfo(server object.Server) int {
	if CheckServerExists(server.IPv4) {
		log.Println("Server already exists with IPV4:", server.IPv4)
		return http.StatusConflict
	}

	affected, err := connector.Engine.Insert(server)
	if affected == 0 {
		log.Println("Failed to add server:", server.IPv4)
		return http.StatusInternalServerError
	}
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
	log.Println("Server added successfully:", server.IPv4)
	return http.StatusCreated
}

func GetServerBySubstr(substr string) ([]object.Server, int) {
	var servers []object.Server
	err := connector.Engine.Table("server").Where("server_name LIKE ?", "%"+substr+"%").Find(&servers)
	if err != nil {
		log.Println("Error retrieving servers:", err)
		return nil, http.StatusInternalServerError
	}
	if len(servers) == 0 {
		log.Println("No servers found with substring:", substr)
		return nil, http.StatusNotFound
	}
	return servers, http.StatusOK
}

func GetServerById(serverId string) (object.Server, bool) {
	var server object.Server
	has, err := connector.Engine.Table("server").Where("server_id = ?", serverId).Get(&server)
	if err != nil {
		log.Println("Error retrieving server by ID:", err)
		return server, false
	}
	if !has {
		log.Println("No server found with ID:", serverId)
		return server, false
	}
	return server, true
}

func UpdateServerInfo(server object.Server) int {
	affected, err := connector.Engine.Table("server").Where("server_id = ?", server.ServerId).Update(server)
	if err != nil {
		log.Println("Error updating server:", err)
		return http.StatusInternalServerError
	}
	if affected == 0 {
		log.Println("No server found with ID:", server.ServerId)
		return http.StatusNotFound
	}
	log.Println("Server updated successfully:", server.ServerId)
	return http.StatusOK
}
