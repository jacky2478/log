package mlog

import (
    "fmt"

    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

func Test_GetOutPutInfo(t *testing.T) {
    Convey("获得输出信息", t, func() {
        strLevel := fmt.Sprintf("%v", LOG_INFO)
        strDirFile := "./logOutPut.go"
        strColor := "[green]"

        getResult := `8 ./logOutPut.go[green]`
        So(GetOutPutInfo(strDirFile, strLevel, strColor), ShouldEqual, getResult)
    })
}

func Test_GetDirFileInfo(t *testing.T) {
    Convey("获得文件目录信息", t, func() {
        getResult := `[convey] [context.go:110]`
        So(GetDirFileInfo(4), ShouldEqual, getResult)
    })
}

func Test_GetLvlColorStr(t *testing.T) {
    Convey("获得日志级别和颜色", t, func() {
        strLevel, strColor := GetLvlColorStr(LOG_INFO, "[info]green")

        getLevel := `[info]`
        getColor := `green`
        So(strLevel, ShouldEqual, getLevel)
        So(strColor, ShouldEqual, getColor)
    })
}
