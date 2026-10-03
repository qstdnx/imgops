package main

type Target struct {
	Name string
	Key  rune
	Id   string
}

var availableTargets = []Target{
	{
		Name: "google",
		Key:  'g',
		Id:   "#t85",
	},
	{
		Name: "bing",
		Key:  'b',
		Id:   "#t101",
	},
	{
		Name: "tineye",
		Key:  't',
		Id:   "#t11",
	},
	{
		Name: "reddit",
		Key:  'r',
		Id:   "#t97",
	},
	{
		Name: "yandex",
		Key:  'y',
		Id:   "#t72",
	},
	{
		Name: "baidu",
		Key:  'a',
		Id:   "#t74",
	},
	{
		Name: "so",
		Key:  's',
		Id:   "#t109",
	},
	{
		Name: "sogou",
		Key:  'u',
		Id:   "#t110",
	},
	{
		Name: "iqdb",
		Key:  'q',
		Id:   "#t78",
	},
	{
		Name: "saucenao",
		Key:  'n',
		Id:   "#t82",
	},
	{
		Name: "tracemoe",
		Key:  'm',
		Id:   "#t201",
	},
	{
		Name: "ascii2d",
		Key:  'd',
		Id:   "#t84",
	},
}

func getKeyToNameTargets(targets []Target) map[rune]string {
	m := make(map[rune]string)

	for _, target := range targets {
		m[target.Key] = target.Name
	}

	return m
}

func getNameToIdTargets(targets []Target) map[string]string {
	m := make(map[string]string)

	for _, target := range targets {
		m[target.Name] = target.Id
	}

	return m
}
