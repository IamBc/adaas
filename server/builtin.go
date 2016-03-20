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

