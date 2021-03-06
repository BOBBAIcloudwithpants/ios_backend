package main

import (
	"github.com/bobbaicloudwithpants/ios_backend/controllers"
	"github.com/bobbaicloudwithpants/ios_backend/middlewares"
	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	//"github.com/service-computing-2020/bbs_backend/middlewares"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/bobbaicloudwithpants/ios_backend/docs" // docs is generated by Swag CLI, you have to import it.
)

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	// api文档自动生成
	url := ginSwagger.URL("http://bobby.run:30085/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	router.MaxMultipartMemory = 5 << 20 // 限制文件大小为5MB

	api := router.Group("/api")
	{
		indexRouter := api.Group("/index")
		{
			indexRouter.GET("", controllers.GetIndexImages)
		}
		infoRouter := api.Group("/info")
		{
			infoRouter.POST("", controllers.Test)
		}

		chatRouter := api.Group("/chats")
		{
			chatRouter.GET("", middlewares.VerifyJWT(),controllers.GetAllChats)
		}

		userRouter := api.Group("/users")
		{

			userRouter.POST("", controllers.UserRegister)
			userRouter.PUT("", controllers.UserLogin)

			userRouter.GET("",middlewares.VerifyJWT(), controllers.GetAllUsers)

			// 单个用户路由
			singleUserRouter := userRouter.Group("/:user_id")
			{
				singleUserRouter.GET("/avatar", controllers.GetAvatar)
				singleUserRouter.Use(middlewares.VerifyJWT())
				singleUserRouter.POST("/avatar", controllers.UploadAvatar)
				singleUserRouter.GET("/subscribe", controllers.GetOneUserSubscribe)
				singleUserRouter.GET("/info", controllers.GetOneUserDetailByUserID)
				singleUserRouter.GET("/posts", controllers.GetOneUserPostsByUserID)
				singleUserRouter.GET("/helped", controllers.GetAllHelpedPeople)
				singleUserRouter.GET("/helps", controllers.GetOneUserHelpByUserID)
			}
		}

		// 获取封面
		cover := api.Group("cover")
		{
			cover.GET("", controllers.GetRecomCover)
		}

		// 获取用户推荐
		recommendation := api.Group("recommendation")
		{
			recommendation.Use(middlewares.VerifyJWT())
			recommendation.GET("", controllers.GetRecommendations)
		}

		notificationRouter := api.Group("/notifications")
		{
			notificationRouter.Use(middlewares.VerifyJWT())
			notificationRouter.GET("", controllers.GetAllUnreadNotification)
		}


		fileRouter := api.Group("/files")
		{
			fileRouter.GET("/:filename", controllers.GetOneFile)
		}

		forumRouter := api.Group("/forums")
		{
			forumRouter.GET("", middlewares.VerifyJWT(), controllers.GetAllPublicFroums)
			forumRouter.POST("", middlewares.VerifyJWT(), controllers.CreateForum)
			// 单个论坛路由
			singleForumRouter := forumRouter.Group("/:forum_id")
			{
				singleForumRouter.GET("/cover", controllers.GetCover)
				singleForumRouter.GET("", middlewares.VerifyJWT(), controllers.GetForumByID)
				singleForumRouter.POST("/cover", middlewares.VerifyJWT(), controllers.UploadCover)
				singleForumRouter.GET("/members", middlewares.VerifyJWT(), controllers.GetForumUsersByForumID)

				// post 路由
				postRouter := singleForumRouter.Group("/posts")
				{
					postRouter.POST("", middlewares.VerifyJWT(), middlewares.CanUserWatchTheForum(),controllers.CreatePost)
					postRouter.GET("", middlewares.VerifyJWT(), middlewares.CanUserWatchTheForum(),controllers.GetAllPostsByForumID)

					singlePostRouter := postRouter.Group("/:post_id")
					{
						singlePostRouter.GET("", middlewares.VerifyJWT(), middlewares.CanUserWatchTheForum(),controllers.GetOnePostDetailByPostID)
						singlePostRouter.POST("/likes", middlewares.VerifyJWT(), middlewares.CanUserWatchTheForum(),controllers.LikeOnePostByPostID)
						singlePostRouter.DELETE("/likes", middlewares.VerifyJWT(), middlewares.CanUserWatchTheForum(), controllers.UnlikeOnePostByPostID)
						singlePostRouter.GET("/files", middlewares.VerifyJWT(), middlewares.CanUserWatchTheForum(),controllers.GetFilesByPostID)

						// comment 路由
						commentRouter := singlePostRouter.Group("/comments")
						commentRouter.Use(middlewares.VerifyJWT(), middlewares.CanUserWatchTheForum())
						{
							commentRouter.POST("", controllers.CreateComment)
							commentRouter.GET("", controllers.GetAllCommentsByPostID)

							singleCommentRouter := commentRouter.Group("/:comment_id")
							{
								singleCommentRouter.GET("", controllers.GetOneCommentDetailByCommentID)
							}
						}

					}
				}

				helpRouter := singleForumRouter.Group("/helps")
				helpRouter.Use(middlewares.VerifyJWT(), middlewares.CanUserWatchTheForum())
				{
					helpRouter.POST("", controllers.CreateHelp)
					helpRouter.GET("/unfinished", controllers.GetAllUnfinishedHelpByForumID)
					helpRouter.GET("/pending", controllers.GetAllPendingHelpByForumID)
					helpRouter.GET("/finished", controllers.GetAllFinishedHelpByForumID)

					singleHelpRouter := helpRouter.Group("/:help_id")
					{
						singleHelpRouter.PATCH("",middlewares.VerifyJWT() ,controllers.ModifyStatusOfOneHelp)
					}
				}

				// role 路由
				roleRouter := singleForumRouter.Group("/role")
				{
					roleRouter.POST("", middlewares.VerifyJWT(), controllers.SubscribeForum)
					roleRouter.DELETE("", middlewares.VerifyJWT(), controllers.UnSubscribeForum)
					roleRouter.PUT("", middlewares.VerifyJWT(), controllers.AddUsersToForum)
					roleRouter.GET("/:user_id", controllers.GetRoleInForum)
					roleRouter.PATCH("/:user_id", middlewares.VerifyJWT(), controllers.UpdateRoleInForum)
				}
			}
		}

	}
	router.Run(":5000")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
