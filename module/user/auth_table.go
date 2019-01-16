package user

const (
	authSqltable = `
	CREATE TABLE IF NOT EXISTS account(
		ID          INTEGER         PRIMARY KEY AUTOINCREMENT,
		username    varchar(20) 	NOT NULL,
		role        INTEGER         NOT NULL,
		psw         varchar(255) 	NOT NULL
	);
	`
	QUERY_isusernameexists string = "SELECT COUNT(*) AS nbr FROM account WHERE username = ?"
	QUERY_isadmincreate    string = "SELECT COUNT(*) AS nbr FROM account WHERE username='admin'"
	QUERY_createaccount    string = "INSERT INTO account (username,role,psw) VALUES (?,?,?)"
	QUERY_accountlogin     string = "SELECT ID,psw,role FROM account WHERE username = ?"
	QUERY_validaccount     string = "SELECT ID FROM account WHERE ID = ? AND psw = ?"
	QUERY_listaccout       string = "SELECT ID,username,role FROM account"
	QUERY_update_pw        string = "UPDATE account SET psw = ? WHERE id = ?"
	QUERY_delete_account   string = "DELETE FROM account WHERE id = ?"
)
