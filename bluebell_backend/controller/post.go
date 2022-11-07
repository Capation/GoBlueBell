package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CreatePostHandler 创建帖子的处理函数
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数的校验

	//c.ShouldBindJSON() // validator --> binding tag
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从 c 取到当前发请求的用户的ID值
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeSeverBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情的处理函数
func GetPostDetailHandler(c *gin.Context) {

	// 1. 获取参数（从url中获取帖子的id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64) // 将字符串转换成数字
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 根据id取出帖子数据(查数据库)
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID(pid) failed", zap.Error(err))
		ResponseError(c, CodeSeverBusy)
		return
	}
	// 3. 返回
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {

	// 获取分页参数
	page, size := getPageInfo(c)

	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeSeverBusy)
		return
	}

	ResponseSuccess(c, data)
	// 返回响应
}

// GetPostListHandler2 升级版帖子列表接口
// 根据前端传来的参数动态的去获取帖子列表
// 按创建时间排序 或者 按分数排序
// 1. 获取参数 query string
// 2. 去redis查询id列表
// 3. 根据id去数据库查询帖子详细信息
func GetPostListHandler2(c *gin.Context) {

	// GET请求参数(query string): /api/v1/posts2?page=1&size=10&order=time
	// 初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	// 获取分页参数
	//c.ShouldBind()     根据请求的数据类型选择相应的方法去获取数据
	//c.ShouldBindJSON() 如果请求中携带的是JSON格式的数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		zap.L().Debug("ShouldBindQuery failed")
		ResponseError(c, CodeInvalidParam)
		return
	}

	fmt.Println("获取请求参数成功")

	// 获取数据
	data, err := logic.GetPostListNew(p) // 更新： 合二为一
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeSeverBusy)
		return
	}

	ResponseSuccess(c, data)
	// 返回响应
}

// GetCommunityPostListHandler 根据社区去查询帖子列表
//func GetCommunityPostListHandler(c *gin.Context) {
//
//	// 初始化结构体时指定初始参数
//	p := &models.ParamCommunityPostList{
//		ParamPostList: models.ParamPostList{
//			Page:  1,
//			Size:  10,
//			Order: models.OrderTime,
//		},
//		//CommunityID:   0,
//	}
//	// 获取分页参数
//	//c.ShouldBind()     根据请求的数据类型选择相应的方法去获取数据
//	//c.ShouldBindJSON() 如果请求中携带的是JSON格式的数据，才能用这个方法获取到数据
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
//		zap.L().Debug("ShouldBindQuery failed")
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	fmt.Println("获取请求参数成功")
//
//	// 获取数据
//	data, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
//		ResponseError(c, CodeSeverBusy)
//		return
//	}
//
//	ResponseSuccess(c, data)
//}
