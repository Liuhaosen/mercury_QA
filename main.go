package main

import (
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/mercury/controller/account"
	"modtest/gostudy/lesson2/mercury/controller/answer"
	"modtest/gostudy/lesson2/mercury/controller/category"
	"modtest/gostudy/lesson2/mercury/controller/comment"
	"modtest/gostudy/lesson2/mercury/controller/favorite"
	"modtest/gostudy/lesson2/mercury/controller/question"
	"modtest/gostudy/lesson2/mercury/dal/db"
	"modtest/gostudy/lesson2/mercury/filter"
	"modtest/gostudy/lesson2/mercury/id_gen"
	maccount "modtest/gostudy/lesson2/mercury/middleware/account"
	"modtest/gostudy/lesson2/mercury/util"

	"github.com/DeanThompson/ginpprof"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

//初始化静态资源/模板
func initTemplate(router *gin.Engine) {
	// router.Static("/static/", "./static/vue")
	router.StaticFile("/", "./static/vue/index.html")
	router.StaticFile("/favicon.ico", "./static/vue/favicon.ico")
	router.Static("/css/", "./static/vue/css")
	router.Static("/fonts/", "./static/vue/fonts")
	router.Static("/img/", "./static/vue/img")
	router.Static("/js/", "./static/vue/js")
	//使用vue以后就不需要加载views了
	// router.LoadHTMLGlob("views/*")
}

func initDB() (err error) {
	dsn := "root:@tcp(localhost:3306)/mercury?parseTime=true"
	err = db.Init(dsn)
	if err != nil {
		return
	}
	return
}

func initSession() (err error) {
	err = maccount.InitSession("memory", "")
	return
}

func initFilter() (err error) {
	err = filter.Init("./data/filter.dat.txt")
	if err != nil {
		logger.Error("init filter failed, err : %v", err)
		return
	}
	logger.Debug("init filter success")
	return
}

func consume(message *sarama.ConsumerMessage) {
	logger.Debug("receive from kafka, msg:%v", message)
}

func main() {
	router := gin.Default()
	// router.Use(maccount.StatCost())

	//初始化log日志库
	config := make(map[string]string)
	config["log_level"] = "debug"
	logger.InitLogger("console", config)

	//初始化数据库连接
	err := initDB()
	if err != nil {
		panic(err)
	}

	//初始化敏感词库
	err = initFilter()
	if err != nil {
		panic(err)
	}

	//初始化session
	err = initSession()
	if err != nil {
		panic(err)
	}

	//初始化Id生成器
	err = id_gen.Init(1)
	if err != nil {
		panic(err)
	}

	//初始化kafka
	err = util.InitKafka("127.0.0.1:9092")
	if err != nil {
		panic(err)
	}

	//初始化kafka消费者
	err = util.InitKafkaConsumer("127.0.0.1:9092", "mercury_topic", consume)
	if err != nil {
		panic(err)
	}

	//分析性能
	ginpprof.Wrapper(router)

	//初始化静态资源和模板
	initTemplate(router)

	// router.GET("/user/login", account.LoginViewHandle)
	// router.GET("/user/register", account.RegisterHandle)

	//使用vue以后只写api即可

	//用户账号组
	accountGroup := router.Group("/api/user")
	{
		accountGroup.POST("/register", account.RegisterHandler) //注册接口
		accountGroup.POST("/login", account.LoginHandler)       //登录接口
	}

	//问题分类组 (查看分类和问题暂时无需登录)
	categoryGroup := router.Group("/api/category")
	{
		categoryGroup.GET("/list", category.GetCategoryListHandler) //获取所有问题分类接口
	}

	//问题分组 (查看分类和问题暂时无需登录)
	questionGroup := router.Group("/api/question")
	{
		questionGroup.POST("/submit", maccount.AuthMiddleware, question.QuestionSubmitHandler) //提交问题接口. 使用中间件, 传入即可
		questionGroup.GET("/list", question.GetQuestionListByCategoryidHandler)                //根据分类获取所有问题的接口
		questionGroup.GET("/detail", question.QuestionDetailHandler)                           //问题详情接口
	}

	//问题的回答分组 (查看分类和问题暂时无需登录)
	answerGroup := router.Group("api/answer")
	{
		answerGroup.GET("/list", answer.GetAnswerListHandler) //获取回答列表接口
	}

	//评论分组 (评论需要登录)
	commentGroup := router.Group("/api/comment")
	// commentGroup := router.Group("/api/comment", maccount.AuthMiddleware) //评论路由组 + 登录验证中间件
	{
		commentGroup.POST("/post_comment", comment.PostCommentHandler) //发表评论
		commentGroup.POST("/post_reply", comment.PostReplyHandler)     //发表评论
		commentGroup.GET("/list", comment.CommentListHandler)          //获取评论列表接口
		commentGroup.GET("/reply_list", comment.ReplyListHandler)      //获取回复列表接口
		commentGroup.POST("/like", comment.LikeHandler)                //点赞接口
	}

	//收藏组
	favoriteGroup := router.Group("/api/favorite")
	// favoriteGroup := router.Group("/api/favorite", maccount.AuthMiddleware)
	{
		favoriteGroup.POST("/add_dir", favorite.AddFavoriteDirHandler)     //新增收藏夹接口
		favoriteGroup.POST("/add", favorite.AddFavoriteHandler)            //添加收藏接口
		favoriteGroup.GET("/dir_list", favorite.GetFavoriteDirListHandler) //获取收藏夹列表接口
		favoriteGroup.GET("/list", favorite.GetFavoriteListHandler)        //获取收藏列表接口
	}

	//运行
	router.Run(":9090")
}
