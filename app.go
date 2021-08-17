package main
//
//import (
//	"fmt"
//	. "github.com/obgnail/MarkdownResouceCollecter/global"
//	. "github.com/obgnail/MarkdownResouceCollecter/strategy"
//	"github.com/obgnail/MarkdownResouceCollecter/utils"
//	"sync"
//)
//
//var (
//	handlerOnce      sync.Once
//	mapTypeToHandler map[string]Handler
//	baseHandler      *BaseHandler
//)
//
//func initHandlerMap() {
//	handlerOnce.Do(func() {
//		baseHandler = new(BaseHandler)
//		mapTypeToHandler = make(map[string]Handler)
//		mapTypeToHandler[CollectLocalPicture] = &CollectLocalPictureHandler{BaseHandler: baseHandler}
//		mapTypeToHandler[CollectNetWorkPicture] = &CollectNetWorkPictureHandler{BaseHandler: baseHandler}
//		mapTypeToHandler[UploadLocalPictureToNetWork] = &UploadLocalPictureToNetWorkHandler{BaseHandler: baseHandler}
//		mapTypeToHandler[ExportMarkdown] = &ExportMarkdownHandler{BaseHandler: baseHandler}
//	})
//}
//
//func BuildHandler(Typ string) Handler {
//	return mapTypeToHandler[Typ]
//}
//
//func BuildHandlers(handlerName []string) ([]Handler, error) {
//	var handlers []Handler
//	for _, name := range handlerName {
//		handle := BuildHandler(name)
//		if handle == nil {
//			return nil, fmt.Errorf("[ERROR] No Such Handler: %s", name)
//		}
//		handlers = append(handlers, handle)
//	}
//	return handlers, nil
//}
//
//func Run(hs []Handler) error {
//	fmt.Println("---------------- Start ----------------")
//	if err := baseHandler.Collect(); err != nil {
//		return nil
//	}
//	baseHandler.BaseAdjust()
//
//	for _, h := range hs {
//		h.Adjust()
//		if err := h.Extra(); err != nil {
//			return err
//		}
//	}
//
//	if err := baseHandler.Rewrite(); err != nil {
//		return err
//	}
//	fmt.Println("---------------- END ----------------")
//	baseHandler.Report()
//	fmt.Printf("\nPLEASE CHECK DIR: %s", DirPath.RootDirPath)
//	return nil
//}
//
//func Start() {
//	if err := utils.CheckExecuteDirExist(); err != nil {
//		fmt.Println(err)
//		return
//	}
//	initHandlerMap()
//
//	handlers, err := BuildHandlers(Cfg.Handlers)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if err := Run(handlers); err != nil {
//		fmt.Println(err)
//	}
//	//Exit()
//}
