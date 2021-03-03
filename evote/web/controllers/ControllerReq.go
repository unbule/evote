package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kongyixueyuan.com/evote/dao"
	"github.com/kongyixueyuan.com/evote/model"
	"github.com/kongyixueyuan.com/evote/service"
	"github.com/kongyixueyuan.com/evote/utils"
	"github.com/kongyixueyuan.com/evote/web/middleware"
	"net/http"
)

type Controller struct {
	Fabric *service.ServiceSetup
}

var selectdao dao.SelectDao

func (con *Controller) Router(engine *gin.Engine) {
	con.Fabric = utils.GetFab()
	engine.GET("/index", con.index)
	engine.POST("/login", con.login)
	engine.POST("/register", con.register)
	engine.POST("/init",con.init)
	engine.POST("/vote", con.vote)
	engine.POST("/getvotecandidate", con.getVotecandidate)
	engine.POST("/test",con.test)
	engine.POST("/queryTest",con.queryTest)
}

func (con *Controller) index(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", gin.H{})
}

func (con *Controller) login(context *gin.Context) {
	var postuser model.User
	var user model.User
	err := context.Bind(&postuser)
	if err != nil {
		panic(err.Error())
	}
	db := utils.Getdb()
	defer db.Close()
	db.Scopes(selectdao.SelectByUsername(postuser.Username)).Find(&user)
	if user.Username == "" {
		context.JSON(http.StatusOK, gin.H{"msg": "not found username", "code": "0"})
	} else if postuser.Password != user.Password {
		context.JSON(http.StatusOK, gin.H{"msg": "username or password was wrong", "code": "0"})
	} else {
		token := middleware.Settoken(postuser.Username)
		context.JSON(http.StatusOK, gin.H{"msg": "success", "code": "1", "token": token})
	}
}

func (con *Controller) register(context *gin.Context) {
	var postuser model.Register
	var user model.User
	var register model.Register
	err := context.Bind(&postuser)
	if err != nil {
		panic(err.Error())
	}
	db := utils.Getdb()
	defer db.Close()
	db.Scopes(selectdao.SelectByUsername(postuser.Username)).Find(&user)
	if user.Username != "" {
		context.JSON(http.StatusOK, gin.H{"msg": "already exist this username", "code": "0"})
	} else {
		user.Username = postuser.Username
		user.Password = postuser.Password
		db.Create(&user)

		register.Username = postuser.Username
		register.Password = postuser.Password
		register.Email = postuser.Email
		register.Status = postuser.Status
		db.Create(&register)
		context.JSON(http.StatusOK, gin.H{"msg": "register success", "code": "1"})
	}
}

func (con *Controller) vote(context *gin.Context) {
	VID1 := context.PostForm("VID1")
	CID1 := context.PostForm("CID1")
	AID1 := context.PostForm("AID1")
	msg,err:=con.Fabric.Vote(VID1,CID1,AID1)
	if err != nil {
		fmt.Printf("vote failed",err.Error())
	}else {
		fmt.Printf("success",msg)
	}
}

func (con *Controller) test(context *gin.Context) {
	candidate := context.PostForm("candidate")
	activity := context.PostForm("activity")
	transactionID, err := con.Fabric.SetInfo(activity, candidate)
	if err != nil {
		fmt.Printf("test ",err.Error())
	} else {
		fmt.Printf("操作成功，交易ID: ",transactionID)
	}
}

func (con *Controller) queryTest(context *gin.Context) {
	activity := context.PostForm("activity")
	fmt.Println(activity)
	msg ,err := con.Fabric.GetInfo(activity)
	if err != nil {
		fmt.Printf("queryTest ",err.Error())
	}else {
		fmt.Printf("success ",msg)
	}
}

func (con *Controller) init(context *gin.Context) {
	activity := context.PostForm("activity")
	candidate1 := context.PostForm("candidate1")
	candidate2 := context.PostForm("candidate2")
	candidate3 := context.PostForm("candidate3")
	msg,err:=con.Fabric.Init_Activity_Candiate(activity,candidate1,candidate2,candidate3)
	if err != nil {
		fmt.Printf("init failed",err.Error())
	}else {
		fmt.Printf("success",msg)
	}
}

func (con *Controller) getVotecandidate(context *gin.Context) {
	AID := context.PostForm("AID1")
	CID := context.PostForm("CID1")
	msg,err := con.Fabric.GetVotecandidate(CID,AID)
	if err != nil {
		fmt.Printf("getVotecandidate failed",err.Error())
	}else {
		fmt.Printf("success",msg)
	}
}

