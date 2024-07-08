package server

import (
	"context"
	"fmt"
	http2 "github.com/simonalong/gole/http"
	"github.com/simonalong/gole/server/rsp"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/simonalong/gole/bean"
	"github.com/simonalong/gole/debug"
	"github.com/simonalong/gole/listener"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/util"

	"github.com/simonalong/gole/logger"

	"github.com/gin-gonic/gin"
)

type HttpMethod int

const (
	HmAll HttpMethod = iota
	HmGet
	HmPost
	HmPut
	HmDelete
	HmOptions
	HmHead
	HmGetPost
	HmNoRoute
)

var GoleVersion = "1.5.1"
var ApiPrefix = "/api"

var engine *gin.Engine = nil
var pprofHave = false

var loadLock sync.Mutex
var serverLoaded = false

//type methodTrees []methodTree

var ginHandlers []gin.HandlerFunc

func init() {
	util.PrintBanner()
	config.LoadConfig()
	printVersionAndProfile()
}

func GetEngine() *gin.Engine {
	return engine
}

// 提供给外部注册使用
func AddGinHandlers(handler gin.HandlerFunc) {
	if nil == ginHandlers {
		var ginHandlersTem []gin.HandlerFunc
		ginHandlers = ginHandlersTem
	}

	ginHandlers = append(ginHandlers, handler)
}

func InitServer() {
	loadLock.Lock()
	defer loadLock.Unlock()
	if serverLoaded {
		return
	}
	if !config.ExistConfigFile() || !config.GetValueBoolDefault("gole.server.enable", false) {
		return
	}

	if !config.ExistConfigFile() {
		logger.Error("没有找到任何配置文件，服务启动失败")
		return
	}
	mode := config.GetValueStringDefault("gole.server.gin.mode", "release")
	if "debug" == mode {
		gin.SetMode(gin.DebugMode)
	} else if "test" == mode {
		gin.SetMode(gin.TestMode)
	} else if "release" == mode {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	engine = gin.New()

	if config.GetValueBoolDefault("gole.debug.enable", true) {
		// 注册pprof
		if config.GetValueBoolDefault("gole.server.gin.pprof.enable", false) {
			pprofHave = true
			pprof.Register(engine)
		}
	}

	if config.GetValueBoolDefault("gole.server.cors.enable", true) {
		engine.Use(Cors())
	}
	engine.Use(gin.Recovery(), ErrHandler())
	engine.Use(RequestSaveHandler())
	engine.Use(rsp.ResponseHandler())
	for _, handler := range ginHandlers {
		engine.Use(handler)
	}

	// 注册 健康检查endpoint
	if config.GetValueBoolDefault("gole.endpoint.health.enable", false) {
		RegisterHealthCheckEndpoint(apiPreAndModule())
	}

	if config.GetValueBoolDefault("gole.debug.enable", true) {
		// 注册 配置查看和变更功能
		if config.GetValueBoolDefault("gole.endpoint.config.enable", false) {
			RegisterConfigWatchEndpoint(apiPreAndModule())
		}

		// 注册 bean管理的功能
		if config.GetValueBoolDefault("gole.endpoint.bean.enable", false) {
			RegisterBeanWatchEndpoint(apiPreAndModule())
		}

		// 注册 debug的帮助命令
		RegisterHelpEndpoint(apiPreAndModule())
	}

	// 注册 swagger的功能
	if config.GetValueBoolDefault("gole.swagger.enable", false) {
		RegisterSwaggerEndpoint()
	}

	// 添加配置变更事件的监听
	listener.AddListener(listener.EventOfConfigChange, ConfigChangeListener)

	logger.InitLog()
	serverLoaded = true
}

func ConfigChangeListener(event listener.GoleEvent) {
	ev := event.(listener.ConfigChangeEvent)
	if ev.Key == "gole.server.gin.pprof.enable" {
		if util.ToBool(ev.Value) && !pprofHave {
			pprofHave = true
			pprof.Register(engine)
		}
	}
}

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				rsp.Fail(c, 500, http2.NetError{ErrMsg: fmt.Sprintf("业务异常：%v", err)})
				return
			}
		}()
		c.Next()
	}
}

func apiPreAndModule() string {
	ap := config.GetValueStringDefault("gole.api.prefix", "")
	if ap != "" {
		ApiPrefix = ap
	}
	return ApiPrefix + "/" + config.ApiModule
}

func printVersionAndProfile() {
	fmt.Printf("----------------------------- gole: %s --------------------------\n", GoleVersion)
	fmt.Printf("profile：%s\n", config.CurrentProfile)
	fmt.Printf("--------------------------------------------------------------------------\n")
}

func Run() {
	listener.AddListener(listener.EventOfServerRunFinish, func(event listener.GoleEvent) {
		logger.Info("服务启动完成")
	})
	StartServer()
}

func StartServer() {
	if !checkEngine() {
		return
	}

	if engine == nil {
		return
	}

	listener.PublishEvent(listener.ServerRunStartEvent{})

	if !config.GetValueBoolDefault("gole.server.enable", true) {
		return
	}

	logger.Info("开始启动服务")
	port := config.GetValueIntDefault("gole.server.port", 8080)
	logger.Info("服务端口号: %d", port)

	graceRun(port)
}

func graceRun(port int) {
	engineServer := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: engine}
	go func() {
		if err := engineServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("启动服务异常 (%v)", err)
		} else {
			// 发送服务关闭事件
			listener.PublishEvent(listener.ServerStopEvent{})
		}
	}()

	// 发送服务启动事件
	listener.PublishEvent(listener.ServerRunFinishEvent{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn("服务端准备关闭...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := engineServer.Shutdown(ctx); err != nil {
		logger.Warn("服务关闭异常: %v", err.Error())
	}
	logger.Warn("服务端退出")
}

func RegisterStatic(relativePath string, rootPath string) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	engine.Static(relativePath, rootPath)
	return engine
}

func RegisterStaticFile(relativePath string, filePath string) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	engine.StaticFile(relativePath, filePath)
	return engine
}

func RegisterPlugin(plugin gin.HandlerFunc) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	engine.Use(plugin)
	return engine
}

func Engine() *gin.Engine {
	if !checkEngine() {
		return nil
	}
	return engine
}

func RegisterHealthCheckEndpoint(apiGole string) gin.IRoutes {
	if "" == apiGole {
		return nil
	}
	RegisterRoute(apiGole+"/system/status", HmAll, healthSystemStatus)
	RegisterRoute(apiGole+"/system/init", HmAll, healthSystemInit)
	RegisterRoute(apiGole+"/system/destroy", HmAll, healthSystemDestroy)
	return engine
}

func RegisterConfigWatchEndpoint(apiGole string) gin.IRoutes {
	if "" == apiGole {
		return nil
	}
	RegisterRoute(apiGole+"/config/values", HmGet, config.GetConfigValues)
	RegisterRoute(apiGole+"/config/values/yaml", HmGet, config.GetConfigDeepValues)
	RegisterRoute(apiGole+"/config/value/:key", HmGet, config.GetConfigValue)
	RegisterRoute(apiGole+"/config/update", HmPut, config.UpdateConfig)
	return engine
}

func RegisterBeanWatchEndpoint(apiGole string) gin.IRoutes {
	if "" == apiGole {
		return nil
	}
	RegisterRoute(apiGole+"/bean/name/all", HmGet, bean.DebugBeanAll)
	RegisterRoute(apiGole+"/bean/name/list/:name", HmGet, bean.DebugBeanList)
	RegisterRoute(apiGole+"/bean/field/get", HmPost, bean.DebugBeanGetField)
	RegisterRoute(apiGole+"/bean/field/set", HmPut, bean.DebugBeanSetField)
	RegisterRoute(apiGole+"/bean/fun/call", HmPost, bean.DebugBeanFunCall)
	return engine
}

func RegisterSwaggerEndpoint() gin.IRoutes {
	RegisterRoute("/swagger/*any", HmGet, ginSwagger.WrapHandler(swaggerFiles.Handler))
	return engine
}

func RegisterHelpEndpoint(apiGole string) gin.IRoutes {
	if "" == apiGole {
		return nil
	}
	RegisterRoute(apiGole+"/debug/help", HmGet, debug.Help)
	return engine
}

func RegisterCustomHealthCheck(apiGole string, status func() string, init func() string, destroy func() string) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	RegisterRoute(apiGole+"/system/status", HmAll, func(c *gin.Context) {
		c.Data(200, "application/json; charset=utf-8", []byte(status()))
	})
	RegisterRoute(apiGole+"/system/init", HmAll, func(c *gin.Context) {
		c.Data(200, "application/json; charset=utf-8", []byte(init()))
	})
	RegisterRoute(apiGole+"/system/destroy", HmAll, func(c *gin.Context) {
		c.Data(200, "application/json; charset=utf-8", []byte(destroy()))
	})
	return engine
}

func checkEngine() bool {
	if engine == nil {
		InitServer()
		return true
	}
	return true
}

func RegisterRoute(path string, method HttpMethod, handler gin.HandlerFunc) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	if engine == nil {
		logger.Warn("server启动失败，请配置 gole.server.enable 或者查看相关日志")
		return nil
	}
	switch method {
	case HmAll:
		engine.GET(path, handler)
		engine.POST(path, handler)
		engine.PUT(path, handler)
		engine.DELETE(path, handler)
		engine.OPTIONS(path, handler)
		engine.HEAD(path, handler)
	case HmGet:
		engine.GET(path, handler)
	case HmPost:
		engine.POST(path, handler)
	case HmPut:
		engine.PUT(path, handler)
	case HmDelete:
		engine.DELETE(path, handler)
	case HmOptions:
		engine.OPTIONS(path, handler)
	case HmHead:
		engine.HEAD(path, handler)
	case HmGetPost:
		engine.GET(path, handler)
		engine.POST(path, handler)
	case HmNoRoute:
		engine.NoRoute(handler)
	}
	return engine
}

func RegisterRouteWithHeaders(path string, method HttpMethod, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	p := GetApiPath(path, method)
	if p == nil {
		p = NewApiPath(path, method)
		switch method {
		case HmAll:
			engine.GET(path, p.Handler)
			engine.POST(path, p.Handler)
			engine.PUT(path, p.Handler)
			engine.DELETE(path, p.Handler)
			engine.OPTIONS(path, p.Handler)
			engine.HEAD(path, p.Handler)
		case HmGet:
			engine.GET(path, p.Handler)
		case HmPost:
			engine.POST(path, p.Handler)
		case HmPut:
			engine.PUT(path, p.Handler)
		case HmDelete:
			engine.DELETE(path, p.Handler)
		case HmOptions:
			engine.OPTIONS(path, p.Handler)
		case HmHead:
			engine.HEAD(path, p.Handler)
		case HmGetPost:
			engine.GET(path, p.Handler)
			engine.POST(path, p.Handler)
		}
	}
	p.AddVersion(header, versionName, handler)
	return engine
}

func NoRoute(handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute("", HmNoRoute, handler)
}

func Post(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmPost, handler)
}

func Delete(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmDelete, handler)
}

func Put(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmPut, handler)
}

func Head(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmHead, handler)
}

func Get(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmGet, handler)
}

func Options(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmOptions, handler)
}

func GetPost(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmGetPost, handler)
}

func All(path string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRoute(getPathAppendApiModel(path), HmAll, handler)
}

func PostWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmPost, header, versionName, handler)
}

func DeleteWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmDelete, header, versionName, handler)
}

func PutWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmPut, header, versionName, handler)
}

func HeadWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmHead, header, versionName, handler)
}

func GetWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmGet, header, versionName, handler)
}

func OptionsWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmOptions, header, versionName, handler)
}

func GetPostWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmGetPost, header, versionName, handler)
}

func AllWith(path string, header []string, versionName []string, handler gin.HandlerFunc) gin.IRoutes {
	return RegisterRouteWithHeaders(getPathAppendApiModel(path), HmAll, header, versionName, handler)
}

func Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	if !checkEngine() {
		return nil
	}
	engine.Use(middleware...)
	return engine
}

func getPathAppendApiModel(path string) string {
	// 获取 api-module
	apiModel := util.ISCString(config.GetValueString("api-module")).Trim("/")
	// 获取api前缀
	ap := util.ISCString(config.GetValueStringDefault("gole.api.prefix", "")).Trim("/")
	if ap != "" {
		ApiPrefix = "/" + string(ap)
	}
	p2 := util.ISCString(path).Trim("/")
	if strings.HasPrefix(string(p2), "api") {
		return fmt.Sprintf("/%s", p2)
	} else {
		return fmt.Sprintf("/%s/%s/%s", ApiPrefix, apiModel, p2)
	}
}

func RequestSaveHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
