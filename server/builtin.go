package main

import "sort"

func GetMin(dataset []int) []int {
    sort.Ints(dataset)
    return []int{ dataset[ len(dataset) -1] }
}

func GetMax(dataset []int) []int {
    sort.Ints(dataset)
    return []int{ dataset[0] }
}

//TODO:
//Gt: Greater Than; Lt: Less than
//func CountGtThreshold 
//func GetGtThreshold
//func CountLtThreshold
//func GetLtThreshold
//func CountGtAvg
//func GetGtAvgThreshold
//func CountLtAvgThreshold
//func GetLtAvgThreshold
// How to isolate custom scripts ? 
