package cmd

import (
	"flag"
	"github.com/Bnei-Baruch/wf-upload/api"
	"github.com/Bnei-Baruch/wf-upload/common"
)


func Init() {
	flag.Parse()
	a := api.App{}
	a.InitClient()
	a.Initialize()
	a.Run(":" + common.PORT)
}
