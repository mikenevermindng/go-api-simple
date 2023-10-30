package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"goApiByGin/db"
	"goApiByGin/ent/todo"
	"goApiByGin/helpers"
	"net/http"
)

type NoteItemInput struct {
	Title       string `validate:"required" json:"title"`
	Description string `validate:"required" json:"description"`
	Priority    int8   `validate:"required" json:"priority"`
}

func AddNote(c *gin.Context) {
	idString, ok := c.Get("userId")

	if !ok {
		c.JSON(http.StatusForbidden, helpers.Response{Code: http.StatusForbidden, Message: "User not found", Data: nil})
		return
	}

	userId := uuid.MustParse(idString.(string))

	var addNoteInput NoteItemInput

	if err := c.ShouldBind(&addNoteInput); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: nil})
		return
	}

	if _, addNoteErr := db.
		Client().Todo.
		Create().
		SetTitle(addNoteInput.Title).
		SetDescription(addNoteInput.Description).
		SetPriority(int(addNoteInput.Priority)).
		SetUser(userId).
		Save(context.Background()); addNoteErr != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: addNoteErr.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, helpers.Response{Code: http.StatusOK, Message: "success", Data: nil})
	return
}

type GetNotesRequest struct {
	Page  int `validate:"required" json:"page"`
	Limit int `validate:"required" json:"limit"`
}

func GetNotes(c *gin.Context) {

	var getNotesRequest GetNotesRequest

	if err := c.ShouldBind(&getNotesRequest); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: nil})
		return
	}

	idString, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusForbidden, helpers.Response{Code: http.StatusForbidden, Message: "User not found", Data: nil})
		return
	}
	userId := uuid.MustParse(idString.(string))
	offset := getNotesRequest.Limit * (getNotesRequest.Page - 1)
	limit := getNotesRequest.Limit
	items, getErr := db.Client().Todo.Query().Where(todo.User(userId)).Limit(limit).Offset(offset).Select("title", "description", "priority").All(context.Background())
	if getErr != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: getErr.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, helpers.Response{Code: http.StatusOK, Message: "success", Data: items})
}
