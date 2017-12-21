package cmd

import (
	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"github.com/hani17/chtq/models"
	"github.com/hani17/chtq/pkg/config"
	"github.com/hani17/chtq/pkg/context"
	"github.com/hani17/chtq/pkg/form"
	"github.com/hani17/chtq/pkg/template"
	"github.com/hani17/chtq/routes"
	"github.com/hani17/chtq/routes/user"
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"
	"log"
	"os"
)

var Web = cli.Command{
	Name:   "web",
	Usage:  "Start web server",
	Action: runWeb,
}

func newMacaron() *macaron.Macaron {
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(cache.Cacher())
	m.Use(captcha.Captchaer())

	m.Use(session.Sessioner(session.Options{
		Provider:       config.Provider,
		ProviderConfig: config.ProviderConfig,
		CookieName:     config.CookieName,
		Maxlifetime:    84600,
	}))

	m.Use(macaron.Static("public", macaron.StaticOptions{
		SkipLogging: false,
	}))

	funcmap := template.NewFuncMap()
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory: "templates",
		Funcs:     funcmap,
	}))

	m.Use(context.Contexter())

	return m
}

func initDB() {
	err := models.NewEngine()
	if err != nil {
		log.Fatal(err)
	}
}

func runWeb(ctx *cli.Context) {

	reqSignOut := context.Toggle(&context.ToggleOptions{SignOutRequired: true})

	m := newMacaron()
	initDB()

	m.Get("/", routes.Home)

	m.Group("/user", func() {
		m.Get("/signup", user.Signup)
		m.Post("/signup", binding.Bind(form.Register{}), user.SignUpPost)

		m.Group("/login", func() {
			m.Combo("").Get(user.Login).
				Post(binding.Bind(form.SignIn{}), user.LoginPost)
		})
	}, reqSignOut)

	m.Group("/user", func() {
		m.Get("/logout", user.Signout)
	})

	m.Group("/:username", func() {
		//m.Get("", user.Profile)
	})

	m.Group("/settings", func() {
		m.Get("", user.Settings)
	})

	m.Get("/new", routes.NewGet)

	m.Run(os.Args)
}
