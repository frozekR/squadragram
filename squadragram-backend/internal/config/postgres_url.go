package config

import (
	"net"
	"net/url"
)

// PostgresURL escapes credentials so characters such as %, #, / and @ in a
// password cannot change the structure of the connection URL.
func PostgresURL(user, password, host, port, database string) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + database,
	}
	query := url.Values{}
	query.Set("sslmode", "disable")
	u.RawQuery = query.Encode()
	return u.String()
}
