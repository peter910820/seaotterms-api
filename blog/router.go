package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	seaottermsdb "seaotterms-db"
)

// 除了身份驗證表的資料庫，其餘資料庫名稱都定義在各站台router包的main.go中
func BlogRouter(apiGroup fiber.Router, dbm *seaottermsdb.DBModel, store *session.Store) {
	blogGroup := apiGroup.Group("/blog")

	blogGroup.Use(GetUserInfo(store)) // global middleware

	// article
	articleRouter(blogGroup, dbm, store)
	tagRouter(blogGroup, dbm, store)

	todoRouter(blogGroup, dbm, store)
	systemTodoRouter(blogGroup, dbm, store)
	userRouter(blogGroup, dbm, store)
	todoTopicRouter(blogGroup, dbm, store)

	// auth
	authRouter(blogGroup, dbm, store)
}

func articleRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, store *session.Store) {
	articleGroup := blogGroup.Group("/articles")

	articleGroup.Get("/", func(c *fiber.Ctx) error {
		return QueryArticle(c, dbm.DB)
	})

	articleGroup.Get("/:id", func(c *fiber.Ctx) error {
		return QueryArticle(c, dbm.DB)
	})

	// No middleware has been implemented yet
	articleGroup.Post("/", CheckManagement(store), func(c *fiber.Ctx) error {
		return CreateArticle(c, dbm.DB)
	})

	// articleGroup.Post("/:id", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
	// 	return ModifyArticle(c, dbm.DB)
	// })

	articleGroup.Delete("/:id", CheckManagement(store), func(c *fiber.Ctx) error {
		return DeleteArticle(c, dbm.DB)
	})
}

func authRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, store *session.Store) {
	authGroup := blogGroup.Group("/auth")

	// get user info
	authGroup.Get("/", GetUserInfo(store), func(c *fiber.Ctx) error {
		return Auth(c, store)
	})

	authGroup.Post("/login", func(c *fiber.Ctx) error {
		return Login(c, store, dbm.DB)
	})
}

func systemTodoRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, store *session.Store) {
	systemTodoGroup := blogGroup.Group("/system-todos")

	systemTodoGroup.Get("/", func(c *fiber.Ctx) error {
		return QuerySystemTodo(c, dbm.DB)
	})

	systemTodoGroup.Post("/", CheckManagement(store), func(c *fiber.Ctx) error {
		return CreateSystemTodo(c, dbm.DB)
	})

	systemTodoGroup.Patch("/:id", CheckManagement(store), func(c *fiber.Ctx) error {
		return UpdateSystemTodo(c, dbm.DB)
	})

	// quick update
	systemTodoGroup.Patch("/quick/:id", CheckManagement(store), func(c *fiber.Ctx) error {
		return QuickUpdateSystemTodo(c, dbm.DB)
	})

	systemTodoGroup.Delete("/:id", CheckManagement(store), func(c *fiber.Ctx) error {
		return DeleteSystemTodo(c, dbm.DB)
	})
}

func tagRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, store *session.Store) {
	tagGroup := blogGroup.Group("/tags")

	tagGroup.Get("/", func(c *fiber.Ctx) error {
		return QueryTag(c, dbm.DB)
	})

	tagGroup.Get("/:name", func(c *fiber.Ctx) error {
		return QueryArticleForTag(c, dbm.DB)
	})

	// create tag
	tagGroup.Post("/", CheckManagement(store), func(c *fiber.Ctx) error {
		return CreateTag(c, dbm.DB)
	})
}

func todoTopicRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, store *session.Store) {
	todoTopicGroup := blogGroup.Group("/todo-topics")

	todoTopicGroup.Get("/:owner", func(c *fiber.Ctx) error {
		return QueryTodoTopic(c, dbm.DB)
	})
	todoTopicGroup.Post("/", CheckLogin(store), func(c *fiber.Ctx) error {
		return CreateTodoTopic(c, dbm.DB)
	})
}

func todoRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, store *session.Store) {
	todoGroup := blogGroup.Group("/todos")

	todoGroup.Get("/:owner", func(c *fiber.Ctx) error {
		return QueryTodoByOwner(c, dbm.DB)
	})
	todoGroup.Post("/", CheckLogin(store), func(c *fiber.Ctx) error {
		return CreateTodo(c, dbm.DB)
	})
	todoGroup.Patch("/:id", CheckLogin(store), func(c *fiber.Ctx) error {
		return UpdateTodoStatus(c, dbm.DB)
	})
	todoGroup.Delete("/:id", CheckLogin(store), func(c *fiber.Ctx) error {
		return DeleteTodo(c, dbm.DB)
	})
}

func userRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, store *session.Store) {
	userGroup := blogGroup.Group("/users")

	userGroup.Get("/", CheckManagement(store), func(c *fiber.Ctx) error {
		return QueryUser(c, dbm.DB)
	})

	userGroup.Post("/", func(c *fiber.Ctx) error {
		return CreateUser(c, dbm.DB)
	})
	userGroup.Patch("/:id", CheckLogin(store), func(c *fiber.Ctx) error {
		return UpdateUser(c, dbm.DB, store)
	})
}
