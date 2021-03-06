package main

import (
	"strconv"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

// " & " goes in front of a variable when you want to get that variable's memory address.
// " * " goes in front of a variable that holds a memory address and resolves it
// 		(it is therefore the counterpart to the & operator).
// 		It goes and gets the thing that the pointer was pointing at,
type Todo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos = []*Todo{
	{Id: 1, Name: "Clean Car", Completed: false},
	{Id: 2, Name: "Clean Room", Completed: false},
}

func main() {
	app := fiber.New()

	app.Use(middleware.Logger())

	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("Hello World...")
	})

	SetupApiV1(app)

	err := app.Listen(3000)
	if err != nil {
		panic(err)
	}

}

func SetupApiV1(app *fiber.App) {
	v1 := app.Group("/v1")

	SetupTodosRoutes(v1)
}

func SetupTodosRoutes(grp fiber.Router) {
	todosRoutes := grp.Group("/todos")

	todosRoutes.Get("/", GetTodos)
	todosRoutes.Post("/", CreateTodo)
	todosRoutes.Get("/:id", GetTodo)
	todosRoutes.Delete("/:id", DeleteTodo)
	todosRoutes.Patch("/:id", UpdateTodo)
}

func GetTodos(ctx *fiber.Ctx) {
	ctx.Status(fiber.StatusOK).JSON(todos)
}

func GetTodo(ctx *fiber.Ctx) {
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	for _, todo := range todos {
		if todo.Id == id {
			ctx.Status(fiber.StatusOK).JSON(todo)
			return
		}
	}
	ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "todo not found",
	})
}

func CreateTodo(ctx *fiber.Ctx) {
	type request struct {
		Name string `json:"name"`
	}

	var body request

	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return
	}

	todo := &Todo{
		Id:        len(todos) + 1,
		Name:      body.Name,
		Completed: false,
	}

	todos = append(todos, todo)

	ctx.Status(fiber.StatusCreated).JSON(todo)

}

func DeleteTodo(ctx *fiber.Ctx) {
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	for i, todo := range todos {
		if todo.Id == id {
			todos = append(todos[0:i], todos[i+1:]...)
			ctx.Status(fiber.StatusNoContent)
			return
		}
	}

	ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "todo not found",
	})
}

func UpdateTodo(ctx *fiber.Ctx) {

	type request struct {
		Name      *string `json:"name"`
		Completed *bool   `json:"completed"`
	}

	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	var body request
	err = ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse body",
		})
		return
	}

	var todo *Todo
	for _, t := range todos {
		if t.Id == id {
			todo = t
			break
		}
	}

	if todo == nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	if body.Name != nil {
		todo.Name = *body.Name
	}

	if body.Completed != nil {
		todo.Completed = *body.Completed
	}

	ctx.Status(fiber.StatusOK).JSON(todo)
}
