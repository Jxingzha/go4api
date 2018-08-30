/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */
 
package fuzz

import (
    // "fmt"
)

///
func GetCombinationValid(fuzzData FuzzData) [][]interface{} {
    var combins [][]interface{}
    for _, validDataMap := range fuzzData.ValidData {
        for _, validList := range validDataMap {
            // fmt.Println("validList: ", key, validList)
            combins = append(combins, validList)
        }
    }

    // Note: hard-code pairWiseLength = 2 first, as currently the pairwise.go works for 2 well, will improve
    pairWiseLength := 2 
    var validTcData [][]interface{}

    // need to consiber the len(combins) = 1 / = 2 / > 2
    if len(combins) >= pairWiseLength {
        c := make(chan []interface{})

        go func(c chan []interface{}) {
            defer close(c)
            GetPairWise(c, combins, 2)
        }(c)

        for tcData := range c {
            validTcData = append(validTcData, tcData)
        }
    } else if len(combins) == 1{
        for _, item := range combins[0] {
            var itemSlice []interface{}
            itemSlice = append(itemSlice, item)
            validTcData = append(validTcData, itemSlice)
        }
            
    }

    return validTcData
}


// -- for the fuzz data
func GetCombinationInvalid(fuzzData FuzzData) [][]interface{} {
    var validCombins [][]interface{}
    for _, validDataMap := range fuzzData.ValidData {
        for _, validList := range validDataMap {
            // fmt.Println("validList: ", key, validList)
            validCombins = append(validCombins, validList)
        }
    }

    var invalidCombins [][]interface{}
    for _, invalidDataMap := range fuzzData.InvalidData {
        for _, invalidList := range invalidDataMap {
            // fmt.Println("invalidList: ", key, invalidList)
            invalidCombins = append(invalidCombins, invalidList)
        }
    }

    // to ensure each negative value will be combined with each positive value(s)
    var invalidTcData [][]interface{}
    for i, _ := range invalidCombins {
        var tcData []interface{}
        if i == 0 {
            tcData = append(tcData, invalidCombins[0][0])
            for j := i + 1; j < len(validCombins); j++ {
                tcData = append(tcData, validCombins[j][0])
            }
        }

        invalidTcData = append(invalidTcData, tcData)
        break
    }
    

    return invalidTcData
}

// Refer to Python:
// product('ABCD', repeat=2)   AA AB AC AD BA BB BC BD CA CB CC CD DA DB DC DD => cartesian product
// permutations('ABCD', 2)   AB AC AD BA BC BD CA CB CD DA DB DC
// combinations('ABCD', 2)   AB AC AD BC BD CD
// combinations_with_replacement('ABCD', 2)   AA AB AC AD BB BC BD CC CD DD


// func GenerateCombinations(alphabet string, length int) <-chan string {
//     c := make(chan string)

//     // Starting a separate goroutine that will create all the combinations,
//     // feeding them to the channel c
//     go func(c chan string) {
//         defer close(c) // Once the iteration function is finished, we close the channel

//         // This is where the iteration will take place
//         // Your teacher's pseudo code uses recursion
//         // which mean you might want to create a separate function
//         // that can call itself.
//     }(c)

//     return c // Return the channel to the calling function
// }


// combinations([]int{1, 2, 3, 4}, 2) =>
// [1 2]
// [1 3]
// [1 4]
// [2 3]
// [2 4]
// [3 4]
func combinationsInt(list []int, length int) (c chan []int) {
    c = make(chan []int)
    go func() {
        defer close(c)
        switch {
            case length == 0:
                c <- []int{}
            case length == len(list):
                c <- list
            case len(list) < length:
                return
            default:
                for i := 0; i < len(list); i++ {
                    for sub_comb := range combinationsInt(list[i + 1:], length - 1) {
                        c <- append([]int{list[i]}, sub_comb...)
                    }
                }
            }
    }()
    return
}

func combinationsInterface(list []interface{}, length int) (c chan []interface{}) {
    c = make(chan []interface{})
    go func() {
        defer close(c)
        switch {
            case length == 0:
                c <- []interface{}{}
            case length == len(list):
                c <- list
            case len(list) < length:
                return
            default:
                for i := 0; i < len(list); i++ {
                    for sub_comb := range combinationsInterface(list[i + 1:], length - 1) {
                        c <- append([]interface{}{list[i]}, sub_comb...)
                    }
                }
            }
    }()
    return
}


func GenerateCombinationsString(data []string, length int) <-chan []string {  
    c := make(chan []string)
    go func(c chan []string) {
        defer close(c)
        combinsString(c, []string{}, data, length)
    }(c)
    return c
}


func combinsString(c chan []string, combin []string, data []string, length int) {  
    // Check if we reached the length limit
    // If so, we just return without adding anything
    if length <= 0 {
        return
    }
    var newCombin []string
    for _, ch := range data {
        newCombin = append(combin, ch)
        // remove this conditional to return all sets of length <= k
        if(length == 1){
            output := make([]string, len(newCombin))
            copy(output, newCombin)
            c <- output
        }
        combinsString(c, newCombin, data, length - 1)
    }
}


// GenerateCombinationsInt([]int{1,2,3,4}, 2) =>
// [1 1][1 2][1 3][1 4][2 1][2 2][2 3][2 4][3 1][3 2][3 3][3 4][4 1][4 2][4 3][4 4]
func GenerateCombinationsInt(data []int, length int) <-chan []int {  
    c := make(chan []int)
    go func(c chan []int) {
        defer close(c)
        combinsInt(c, []int{}, data, length)
    }(c)
    return c
}


func combinsInt(c chan []int, combin []int, data []int, length int) {  
    // Check if we reached the length limit
    // If so, we just return without adding anything
    if length <= 0 {
        return
    }
    var newCombin []int
    for _, ch := range data {
        newCombin = append(combin, ch)
        // remove this conditional to return all sets of length <=k
        if(length == 1){
            output := make([]int, len(newCombin))
            copy(output, newCombin)
            c <- output
        }
        combinsInt(c, newCombin, data, length - 1)
    }
}


func combinsSliceString(c chan []interface{}, combin []interface{}, data [][]interface{}) {  
    if len(data) > 1 {
        var newCombin []interface{}
        for _, i_v := range data[0] {
            newCombin = append(combin, i_v)

            combinsSliceString(c, newCombin, data[1:])
        }

    } else if len(data) == 1 {
        for _, j_v := range data[0] {
            output := make([]interface{}, len(combin))
            copy(output, combin)

            output = append(output, j_v)
            // fmt.Println("output: ", output)
            c <- output
        }
    }
}

