package server

import (
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	mw_logger "github.com/gofiber/fiber/v2/middleware/logger"
	cm "main/pkg/utils/common"
	"os"
	"time"
)

func (s *server) NewHttpServer() (serve *fiber.App, err error) {
	serve = fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: false,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		ServerHeader:  os.Getenv("SERVER_HEADER"),
		AppName:       os.Getenv("APP_TITLE") + " " + os.Getenv("APP_VERSION"),
		Immutable:     true,
	})

	serve.Use(cors.New())
	serve.Use(mw_logger.New())

	serve.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        5,
		Expiration: 10 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendFile("./toofast.html")
		},
	}))

	serve.Get("/doc/*", swagger.HandlerDefault)

	serve.Get("/", func(c *fiber.Ctx) error {
		s.logger.Infof("Health check RequestID: %d", cm.GenNum())
		return c.SendString(s.cfg.Server.PROJECT_NAME)
	})

	return
}

func (s *server) Listen(serve *fiber.App) (err error) {
	port := os.Getenv("PORT")

	if port == "" {
		port = s.cfg.Http.PORT
		s.logger.Warnf("defaulting to ports %s", port)
	}
	URI := fmt.Sprintf("%s:%s", "", port)
	if s.cfg.Server.APP_ENV == "prod" {
		if err = serve.ListenTLS(URI, s.cfg.Http.SSL_CERT_PATH, s.cfg.Http.SSL_CERT_KEY); err != nil {
			s.logger.Fatalf("Error starting Server with SSL : ", err)
		}
	} else {
		if err = serve.Listen(URI); err != nil {
			s.logger.Fatalf("Error starting Server : ", err)
		}
	}

	s.logger.Infof("Server is listening on PORT: %s", s.cfg.Http.PORT)
	return
}
