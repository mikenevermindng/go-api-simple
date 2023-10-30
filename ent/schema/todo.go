package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("user", uuid.UUID{}),
		field.String("title").NotEmpty(),
		field.String("description").NotEmpty(),
		field.Int("priority").Min(1).Max(3),
	}
}

// Edges of the Todo.
func (Todo) Edges() []ent.Edge {
	return nil
}
