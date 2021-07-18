package controllers

import (
	"github.com/mohamadsyukur99/fullstack/api/middlewares"
)

func (s *Server) initializeRoutes() {
	v1 := s.Router.Group("api/v1")

	{
		//Login Router
		v1.POST("login", s.Login)

		// Reset Password:
		v1.POST("password/forgot", s.ForgotPassword)
		v1.POST("password/reset", s.ResetPassword)

		// User Routes
		v1.GET("/cek", middlewares.TokenAuthMiddleware(), s.Cek)
		v1.POST("/users", s.CreateUser)
		v1.GET("/users", middlewares.TokenAuthMiddleware(), s.GetUsers)
		v1.GET("/users/:id", middlewares.TokenAuthMiddleware(), s.GetUser)
		v1.GET("/usersnama/:nama", middlewares.TokenAuthMiddleware(), s.GetUserByName)
		v1.GET("/usersname/:nama", middlewares.TokenAuthMiddleware(), s.GetUserUsername)
		v1.PUT("/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateUser)
		v1.DELETE("/users/:id", middlewares.TokenAuthMiddleware(), s.DeleteUser)

		// Siswa Routes
		v1.POST("/siswa", s.CreateSiswa)
		v1.GET("/siswa", middlewares.TokenAuthMiddleware(), s.GetSiswaAll)
		v1.GET("/siswaid/:id", middlewares.TokenAuthMiddleware(), s.GetSiswa)
		v1.GET("/siswafilter/:nama", middlewares.TokenAuthMiddleware(), s.GetSiswaByName)
		v1.GET("/no_induk/:nama", middlewares.TokenAuthMiddleware(), s.GetSiswaNoInduk)
		v1.PUT("/siswa/:id", middlewares.TokenAuthMiddleware(), s.UpdateSiswa)
		v1.DELETE("/siswa/:id", middlewares.TokenAuthMiddleware(), s.DeleteSiswa)

		// Guru Routes
		v1.POST("/guru", s.CreateGuru)
		v1.GET("/guru", middlewares.TokenAuthMiddleware(), s.GetGuruAll)
		v1.GET("/guruid/:id", middlewares.TokenAuthMiddleware(), s.GetGuru)
		v1.GET("/gurufilter/:nama", middlewares.TokenAuthMiddleware(), s.GetGuruByName)
		v1.GET("/nip/:nama", middlewares.TokenAuthMiddleware(), s.GetGuruNip)
		v1.PUT("/guru/:id", middlewares.TokenAuthMiddleware(), s.UpdateGuru)
		v1.DELETE("/guru/:id", middlewares.TokenAuthMiddleware(), s.DeleteGuru)

		// Post Routes
		v1.POST("/posts", s.CreatePost)
		v1.GET("/posts", s.GetPosts)
		v1.GET("/posts/:id", s.GetPost)
		v1.PUT("/posts/:id", middlewares.TokenAuthMiddleware(), s.UpdatePost)
		v1.DELETE("/posts/:id", middlewares.TokenAuthMiddleware(), s.DeletePost)

	}
}
