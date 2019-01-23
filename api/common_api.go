/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package api

import (
    "fmt"
    "strings" 
    "time"
    "strconv"
    "encoding/json"

    "go4api/lib/testcase"

    gjson "github.com/tidwall/gjson"
)

func (tcDataStore *TcDataStore) CommandGroup (cmdGroupOrigin []*testcase.CommandDetails) (string, [][]*testcase.TestMessage) {
    finalResults := "Success"
    var cmdsResults []bool
    var finalTestMessages [][]*testcase.TestMessage

    for i := 0; i < tcDataStore.CmdGroupLength; i ++ {
        var sResults []bool
        var sMessages [][]*testcase.TestMessage

        cmdType := cmdGroupOrigin[i].CmdType

        switch strings.ToLower(cmdType) {
            case "sql":
                sResults, sMessages = tcDataStore.HandleSqlCmd(i)

                cmdsResults = append(cmdsResults, sResults[0:]...)
                finalTestMessages = append(finalTestMessages, sMessages[0:]...)
            case "redis":
                sResults, sMessages = tcDataStore.HandleRedisCmd(i)

                cmdsResults = append(cmdsResults, sResults[0:]...)
                finalTestMessages = append(finalTestMessages, sMessages[0:]...)
            case "init":
                sResults, sMessages = tcDataStore.HandleInitCmd(i)

                cmdsResults = append(cmdsResults, sResults[0:]...)
                finalTestMessages = append(finalTestMessages, sMessages[0:]...)
            default:
                fmt.Println("!! warning, command ", cmdType, " can not be recognized.")
        }
    }

    for i, _ := range cmdsResults {
        if cmdsResults[i] == false {
            finalResults = "Fail"
            break
        }
    }

    return finalResults, finalTestMessages
}

func (tcDataStore *TcDataStore) HandleSqlCmd (i int) ([]bool, [][]*testcase.TestMessage) {
    var sResults []bool
    var sMessages [][]*testcase.TestMessage

    cmdStrPath := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmd"
    tcDataStore.PrepVariablesBuiltins(cmdStrPath)

    tcDataJsonB, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonB)

    cmdStr := gjson.Get(tcDataJson, cmdStrPath).String()

    cmdTgtDb := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmdSource"
    tgtDb := gjson.Get(tcDataJson, cmdTgtDb).String()
    // init
    tcDataStore.CmdType = "sql"
    tcDataStore.CmdExecStatus = ""
    tcDataStore.CmdAffectedCount = -1
    tcDataStore.CmdResults = -1

    // call sql
    if len(tgtDb) == 0 {
        fmt.Println("No target db provided, default to master")
        tgtDb = "master"
    }
    cmdAffectedCount, _, cmdResults, cmdExecStatus := RunSql(tgtDb, cmdStr)
    
    tcDataStore.CmdExecStatus = cmdExecStatus
    tcDataStore.CmdAffectedCount = cmdAffectedCount
    tcDataStore.CmdResults = cmdResults

    sResults, sMessages = tcDataStore.HandleSingleCmdResult(i)

    return sResults, sMessages
}

func (tcDataStore *TcDataStore) HandleRedisCmd (i int) ([]bool, [][]*testcase.TestMessage) {
    var sResults []bool
    var sMessages [][]*testcase.TestMessage

    var cmdStr, cmdKey, cmdValue string

    cmdStrPath := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmd"
    tcDataStore.PrepVariablesBuiltins(cmdStrPath)

    tcDataJsonB, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonB)

    cmdMap := gjson.Get(tcDataJson, cmdStrPath).Map()

    // cmdTgtDb := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmdSource"
    // tgtDb := gjson.Get(tcDataJson, cmdTgtDb).String()

    for k, v := range cmdMap {
        cmdStr = k
        if len(v.Array()) == 1 {
            cmdKey = v.Array()[0].String()
            cmdValue = ""
        }
        if len(v.Array()) > 1 {
            cmdKey = v.Array()[0].String()
            cmdValue = v.Array()[1].String()
        }
    }
    // init
    tcDataStore.CmdType = "redis"
    tcDataStore.CmdExecStatus = ""
    tcDataStore.CmdAffectedCount = -1
    tcDataStore.CmdResults = -1

    cmdAffectedCount, cmdResults, cmdExecStatus := RunRedis(cmdStr, cmdKey, cmdValue)
    
    tcDataStore.CmdExecStatus = cmdExecStatus
    tcDataStore.CmdAffectedCount = cmdAffectedCount
    tcDataStore.CmdResults = cmdResults

    sResults, sMessages = tcDataStore.HandleSingleCmdResult(i)

    return sResults, sMessages
}

func (tcDataStore *TcDataStore) HandleInitCmd (i int) ([]bool, [][]*testcase.TestMessage) {
    var sResults []bool
    var sMessages [][]*testcase.TestMessage

    cmdStrPath := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmd"
    tcDataStore.PrepVariablesBuiltins(cmdStrPath)

    tcDataJsonB, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonB)

    cmdStr := gjson.Get(tcDataJson, cmdStrPath).String()

    s := strings.ToLower(cmdStr)

    if len(cmdStr) > 0 && strings.Contains(s, "sleep") {
        // here may has cmd "sleep xx" for debug purpose
        t := strings.Split(s, " ")
        if len(t) == 1 {
            fmt.Println("No sleep duration provided, no sleep")
        } else {
            tm, err := strconv.Atoi(t[1])
            if err != nil {
                fmt.Println("Provided sleep duration is not number, please fix")
            }

            time.Sleep(time.Duration(tm)*time.Second)
        }
    }

    // as maybe no cmd is executed, the CmdExecStatus is always "cmdSuccess"
    // init
    tcDataStore.CmdType = "init"
    tcDataStore.CmdExecStatus = "cmdSuccess"
    tcDataStore.CmdAffectedCount = -1
    tcDataStore.CmdResults = -1

    sResults, sMessages = tcDataStore.HandleSingleCmdResult(i)

    return sResults, sMessages
}


func (tcDataStore *TcDataStore) HandleSingleCmdResult (i int) ([]bool, [][]*testcase.TestMessage) {
    // --------
    var cmdsResults []bool
    var finalTestMessages = [][]*testcase.TestMessage{}
    var cmdGroup []*testcase.CommandDetails

    if tcDataStore.CmdExecStatus == "cmdSuccess" {
        path := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmdResponse"
        tcDataStore.PrepVariablesBuiltins(path)
        //
        switch tcDataStore.CmdSection {
        case "setUp":
            cmdGroup = tcDataStore.TcData.TestCase.SetUp()
        case "tearDown":
            cmdGroup = tcDataStore.TcData.TestCase.TearDown()
        }

        cmdExpResp := cmdGroup[i].CmdResponse

        singleCmdResults, testMessages := tcDataStore.CompareRespGroup(cmdExpResp)

        fmt.Println("singleCmdResults, testMessages: ", singleCmdResults, testMessages)
        // HandleSingleCommandResults for out
        if singleCmdResults == true {
            tcDataStore.HandleCmdResultsForOut(i)
        }

        cmdsResults = append(cmdsResults, singleCmdResults)
        finalTestMessages = append(finalTestMessages, testMessages)
    } else {
        cmdsResults = append(cmdsResults, false)
    }

    return cmdsResults, finalTestMessages
}

func (tcDataStore *TcDataStore) CompareRespGroup (cmdExpResp map[string]interface{}) (bool, []*testcase.TestMessage){
    //-----------
    singleCmdResults := true
    var testResults []bool
    var testMessages []*testcase.TestMessage

    for key, value := range cmdExpResp {
        cmdExpResp_sub := value.(map[string]interface{})
        for assertionKey, expValueOrigin := range cmdExpResp_sub {
            
            actualValue := tcDataStore.GetResponseValue(key)

            var expValue interface{}
            switch expValueOrigin.(type) {
                case float64, int64: 
                    expValue = expValueOrigin
                default:
                    expValue = tcDataStore.GetResponseValue(expValueOrigin.(string))
            }
            
            testRes, msg := compareCommon(tcDataStore.CmdType, key, assertionKey, actualValue, expValue)
            
            testMessages = append(testMessages, msg)
            testResults = append(testResults, testRes)
        }
    }

    for key := range testResults {
        if testResults[key] == false {
            singleCmdResults = false
            break
        }
    }

    return singleCmdResults, testMessages
}

func (tcDataStore *TcDataStore) HandleCmdResultsForOut (i int) {
    var cmdGroup []*testcase.CommandDetails

    // aa, _ := json.Marshal(tcDataStore)
    // fmt.Println("tcDataStore: ", string(aa))

    // write out session if has
    path := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".session"
    tcDataStore.PrepVariablesBuiltins(path)

    switch tcDataStore.CmdSection {
        case "setUp":
            cmdGroup = tcDataStore.TcData.TestCase.SetUp()
        case "tearDown":
            cmdGroup = tcDataStore.TcData.TestCase.TearDown()
    }

    expTcSession := cmdGroup[i].Session
    tcDataStore.WriteSession(expTcSession)

    // write out global variables if has
    path = "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".outGlobalVariables"
    tcDataStore.PrepVariablesBuiltins(path)

    switch tcDataStore.CmdSection {
        case "setUp":
            cmdGroup = tcDataStore.TcData.TestCase.SetUp()
        case "tearDown":
            cmdGroup = tcDataStore.TcData.TestCase.TearDown()
    }
    
    expOutGlobalVariables := cmdGroup[i].OutGlobalVariables
    tcDataStore.WriteOutGlobalVariables(expOutGlobalVariables)

    // write out tc local variables if has
    path = "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".outLocalVariables"
    tcDataStore.PrepVariablesBuiltins(path)

    switch tcDataStore.CmdSection {
        case "setUp":
            cmdGroup = tcDataStore.TcData.TestCase.SetUp()
        case "tearDown":
            cmdGroup = tcDataStore.TcData.TestCase.TearDown()
    }

    expOutLocalVariables := cmdGroup[i].OutLocalVariables
    tcDataStore.WriteOutTcLocalVariables(expOutLocalVariables)

    // write out files if has
    path = "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".outFiles"
    tcDataStore.PrepVariablesBuiltins(path)

    switch tcDataStore.CmdSection {
        case "setUp":
            cmdGroup = tcDataStore.TcData.TestCase.SetUp()
        case "tearDown":
            cmdGroup = tcDataStore.TcData.TestCase.TearDown()
    }

    expOutFiles := cmdGroup[i].OutFiles
    tcDataStore.HandleOutFiles(expOutFiles)
}

