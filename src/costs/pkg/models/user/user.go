package user

import (
	"database/sql"
	"gateway/pkg/models/privilege"
	"gateway/pkg/models/tickets"
)

type UserInfo struct {
	Privilege   *privilege.PrivilegeShortInfo `json:"privilege"`
	TicketsInfo *[]tickets.TicketInfo         `json:"tickets"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	password string
	updated  sql.NullString
}
