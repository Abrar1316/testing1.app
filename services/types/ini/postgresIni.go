package types

type PostgresIni struct {
	host     string
	username string
	password string
	database string
}

func NewPostgresIni(host, username, password, database string) *PostgresIni {
	return &PostgresIni{host: host, username: username, password: password, database: database}
}

func (m *PostgresIni) GetHost() string {
	if m != nil {
		return m.host
	}
	return ""
}

func (m *PostgresIni) GetUsername() string {
	if m != nil {
		return m.username
	}
	return ""
}

func (m *PostgresIni) GetPassword() string {
	if m != nil {
		return m.password
	}
	return ""
}

func (m *PostgresIni) GetDatabase() string {
	if m != nil {
		return m.database
	}
	return ""
}
