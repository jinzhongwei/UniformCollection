package main

import (
	"errors"
	"fmt"
	"math"
)

type UniformCollection struct {
	itemss         [][]item
	value          float64
	changeOnceCost float64
	original       [][]item
	avg            float64
}

type item struct {
	id    int
	value float64
}

func (uc *UniformCollection) GetResult() [][]item {
	uc.avg = uc.compAvg(uc.itemss)
	rel := uc.judge(uc.itemss)
	small := 9999999999.0
	smallIndex := 0
	for index, relItem := range rel {
		currentCost := uc.cost(relItem)
		if currentCost < small {
			smallIndex = index
			small = currentCost
		}
	}
	if len(rel) == 0 {
		fmt.Println("未找到可行解")
		return [][]item{}
	}
	fmt.Println(small)
	return rel[smallIndex]
}
func (uc *UniformCollection) SetChangeOnceCost(cost float64) {
	uc.changeOnceCost = cost
}
func (uc *UniformCollection) SetItem(items map[int]float64) (err error) {
	if len(items) <= 0 {
		return errors.New("items is empty")
	}
	var thisItems []item
	for id, value := range items {
		var i item
		i.id = id
		i.value = value
		thisItems = append(thisItems, i)
	}
	uc.itemss = append(uc.itemss, thisItems)
	uc.original = append(uc.original, thisItems)
	return nil
}

func (uc *UniformCollection) compAvg(itemss [][]item) float64 {
	total := 0.0
	for _, i := range itemss {
		for _, j := range i {
			total = total + j.value
		}
	}
	avg := total / float64(len(itemss))
	return avg
}

func (uc *UniformCollection) cost(this [][]item) float64 {
	rel := 0.0
	for _, i := range this {
		total := 0.0
		for _, j := range i {
			total = total + j.value
		}
		rel = rel + math.Pow(total-uc.avg, 2)
	}
	changeNumber := 0
	for index, _ := range this {
		changeNumber = changeNumber + uc.setDiffLen(this[index], uc.original[index]) + uc.setDiffLen(uc.original[index], this[index])
	}
	return rel + float64(changeNumber)*uc.changeOnceCost
}

func (uc *UniformCollection) setDiffLen(a, b []item) int {
	aMap := make(map[int]bool)
	bMap := make(map[int]bool)
	var diffSet []int
	for _, i := range a {
		aMap[i.id] = true
	}
	for _, i := range b {
		bMap[i.id] = true
	}
	for key, _ := range aMap {
		if _, ok := bMap[key]; !ok {
			diffSet = append(diffSet, key)
		}
	}
	return len(diffSet)
}

/*
	1.求均值
	2.求当前的代价
	3.把大于均值的和小于均值的放在两个数组中
	4.两层遍历大小数组，从大的数组中任意挑一个数字到小的数组中，计算代价，小于当前代价则返回
*/

func (uc *UniformCollection) judge(a [][]item) [][][]item {
	//先求均值
	var rel [][][]item
	currentCost := uc.cost(a)
	var large, small [][]item
	for _, i := range a {
		total := 0.0
		for _, j := range i {
			total = total + j.value
		}
		if total >= uc.avg {
			large = append(large, i)
		} else {
			small = append(small, i)
		}
	}
	for i, _ := range small {
		for j, _ := range large {
			changeRel := uc.changeAndMerge(small, i, large, j)
			thisCost := uc.cost(changeRel)
			if thisCost < currentCost {
				rel = append(rel, uc.judge(changeRel)...)
			}
		}
	}
	if len(rel) == 0 {
		rel = append(rel, a)
	}
	return rel
}
func (uc *UniformCollection) changeAndMerge(small [][]item, i int, large [][]item, j int) [][]item {
	var rel [][]item
	total1 := 0.0
	total2 := 0.0
	for _, v := range small[i] {
		total1 = total1 + v.value
	}
	for _, v := range large[j] {
		total2 = total2 + v.value
	}
	diff1 := uc.avg - total1
	diff2 := total2 - uc.avg
	smallest := diff1 + diff2
	smallIndex := 0
	for index, value := range large[j] {
		if math.Abs(diff1-value.value)+math.Abs(diff2-value.value) < smallest {
			smallIndex = index
			smallest = math.Abs(diff1-value.value) + math.Abs(diff2-value.value)
		}
	}
	for index, value := range large {
		var value1 []item
		value1 = append(value1, value...)
		if index == j {
			value1 = append(value1[:smallIndex], value1[smallIndex+1:]...)
		}
		rel = append(rel, value1)
	}

	for index, value := range small {
		if index == i {
			value = append(value, large[j][smallIndex])
		}
		rel = append(rel, value)
	}
	return rel
}
func main() {
	valueList := []map[int]float64{map[int]float64{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9, 10: 10}, map[int]float64{11: 11, 12: 12, 13: 13, 14: 14, 15: 15, 16: 16}, map[int]float64{17: 17, 18: 18, 19: 19, 20: 20, 21: 21, 22: 22, 23: 23}}
	var judge UniformCollection
	for _, item := range valueList {
		judge.SetItem(item)
	}
	judge.SetChangeOnceCost(10)
	fmt.Println(judge.GetResult())
}
