package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func main() {
	data, err := ioutil.ReadFile("./Recipes.json")
	if err != nil {
		panic(err)
	}
	r := make([]Recipes, 10)
	err = json.Unmarshal(data, &r)
	if err != nil {
		panic(err)
	}
	for _, s := range r {
		println(s.Name, s.Output)
	}
	res := make([]Recipes, 1)
	var in string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	in = scanner.Text()
	if in == "Devmode" {
		DevCheck(r)
		os.Exit(2)
	}
	var demand float64
	_, err = fmt.Scanf("%f", &demand)
	flag := true
	for _, s := range r {
		if s.Name == in {
			res[0] = s
			flag = false
		}
	}
	if flag {
		println("Not Found")
		os.Exit(-1)
	}
	res = res[0].Search(r, demand / res[0].Output)
	output := GetSum(res)
	fmt.Printf("Belows are need to Create %s amount of %f\n-------------------------------\n", in, demand)
	for _, s := range output {
		fmt.Printf("%s, %.2f\n", s.Name, s.Output)
	}
}

type Recipes struct {
	Name   string  `json:"name"`
	Items  []Item  `json:"item"`
	Output float64 `json:"output"`
	AdditionalOutput []Item `json:"additional_output"`
}

type Item struct {
	Name  string  `json:"name"`
	Input float64 `json:"input"`
}

func (r Recipes) Search(m []Recipes, factor float64) []Recipes {
	res := make([]Recipes, 0)
	res = append(res, Recipes{r.Name, nil, r.Output * factor, nil})
	if r.AdditionalOutput != nil {
		for _, s := range r.AdditionalOutput {
			res = append(res, Recipes{s.Name, nil, s.Input * factor, nil})
		}
	}
	for _, s := range r.Items {
		for _, t := range m {
			if s.Name == t.Name {
				println(s.Name)
				res = append(res, t.Search(m, s.Input/t.Output*factor)...)
			}
		}
	}
	return res
}

func GetSum(r []Recipes)[] Recipes{
	tmp := r[1:]
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].Name < tmp[j].Name  //TODO 変換用のテーブル作って見やすくする
	})
	str := ""
	index := -1
	res := make([]Recipes, 0)
	for _, s := range tmp {
		if s.Name != str {
			str = s.Name
			res = append(res, Recipes{s.Name, nil, 0, nil})
			index++
		}
		res[index].Output += s.Output
	}
	return res
}

func DevCheck(r []Recipes){
	strs := make([]DevChecker, 0)
	f, err := os.Open("./items")
	if err != nil {
		log.Fatal(err)
		return
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		strs = append(strs, DevChecker{s.Text(), false})
	}
	if s.Err() != nil {
		// non-EOF error.
		log.Fatal(s.Err())
	}
	for i, s := range strs {
		for _, t := range r {
			if s.Str == t.Name {
				if s.Flag {
					panic("Duplication Error")
				}
				strs[i].SetFlag(true)
			}
		}
	}
	for _, s := range strs{
		if s.Flag {
			println(s.Str + ":〇")
		}else {
			println(s.Str + ":×")
		}
	}
}

type DevChecker struct {
	Str  string
	Flag bool
}

func (d *DevChecker) SetFlag(i bool) {
	d.Flag = i
}