/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package reports

import (
 	// "fmt"
 	"encoding/json"

	// "go4api/texttmpl"
	"go4api/lib/testcase"

    . "github.com/ahmetb/go-linq"
)

type ReportsMStats struct {
    ReportKey interface{}
    Count int
}


func GetMStats_1 () []Group {
    type ReportsStuct struct {
        Path string
        Method string
        // MutationArea string
        // MutationCategory string
        // MutationRule string
        ActualStatusCode int
        // TestResult string
    }

    var query []Group

    From(ExecutionResultSlice).GroupByT(
        func(item *testcase.TcReportResults) ReportsStuct { 
            return ReportsStuct{item.Path, item.Method, item.ActualStatusCode}
        },
        func(item *testcase.TcReportResults) int { return 1 },
    ).ToSlice(&query)

    return query
}

func GetMStats_2 () []Group {
    type ReportsStuct struct {
        Path string
        Method string
        MutationArea string
        // MutationCategory string
        // MutationRule string
        ActualStatusCode int
        // TestResult string
    }

    var query []Group

    From(ExecutionResultSlice).GroupByT(
        func(item *testcase.TcReportResults) ReportsStuct { 
            return ReportsStuct{item.Path, item.Method, item.MutationArea, item.ActualStatusCode}
        },
        func(item *testcase.TcReportResults) int { return 1 },
    ).ToSlice(&query)

    return query
}


func GetMStats_3 () []Group {
    type ReportsStuct struct {
        Path string
        Method string
        MutationArea string
        MutationCategory string
        // MutationRule string
        ActualStatusCode int
        // TestResult string
    }

    var query []Group

    From(ExecutionResultSlice).GroupByT(
        func(item *testcase.TcReportResults) ReportsStuct { 
            return ReportsStuct{item.Path, item.Method, item.MutationArea, item.MutationCategory, item.ActualStatusCode}
        },
        func(item *testcase.TcReportResults) int { return 1 },
    ).ToSlice(&query)

    return query
}


func PrintGroup (query []Group) []ReportsMStats {
    var reportsMStatsSlice []ReportsMStats

    for _, q := range query {
        ii := 0
        for range q.Group {
            ii += 1
        }

        reportsMStats := ReportsMStats {
            ReportKey: q.Key,
            Count: ii,
        }
        reportsMStatsSlice = append(reportsMStatsSlice, reportsMStats)
    }
    return reportsMStatsSlice
}

func GetMutationStatsJson() []string {
    var reJsons []string

    query := GetMStats_1()
    reportsMStatsSlice := PrintGroup(query)

    reJson, _ := json.Marshal(reportsMStatsSlice)
    reJsons = append(reJsons, string(reJson))
    // fmt.Println("=====> reportsMStatsSlice: ", string(reJson))

    query = GetMStats_2()
    reportsMStatsSlice = PrintGroup(query)

    reJson, _ = json.Marshal(reportsMStatsSlice)
    reJsons = append(reJsons, string(reJson))
    // fmt.Println("=====> reportsMStatsSlice: ", string(reJson))

    query = GetMStats_3()
    reportsMStatsSlice = PrintGroup(query)

    reJson, _ = json.Marshal(reportsMStatsSlice)
    reJsons = append(reJsons, string(reJson))
    // fmt.Println("=====> reportsMStatsSlice: ", string(reJson))

    return reJsons
}

