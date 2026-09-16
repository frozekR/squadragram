package character

import (
	"github.com/google/uuid"
)

func GetCharacterByIDQuery() {

}

func QuerySelectByUUID(u uuid.UUID) (string, []any) {
	return `SELECT uuid, name, description, role
	 FROM characters
	 WHERE uuid = $1
	 `, []any{u}
}
