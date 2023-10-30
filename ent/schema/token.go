package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// Token holds the schema definition for the Token entity.
type Token struct {
	ent.Schema
}

// Fields of the Token.
func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("user", uuid.UUID{}),
		field.String("token").NotEmpty(),
		field.String("refresh_token").NotEmpty(),
		field.String("secret").NotEmpty(),
		field.String("device").Optional().Nillable(),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at").Nillable().Default(time.Now()),
		field.String("ip").Optional(),
	}
}

// Edges of the Token.
func (Token) Edges() []ent.Edge {
	return nil
}
