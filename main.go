package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	maxClosestRoom = 3
)

type Room struct {
	name     string
	state    string            // строка состояния комнаты, выводится при команде "осмотреться"
	intro    string            // строка
	items    map[string][]Item // словарь со всеми предметами в этой комнате ключ: место, значение - предмет
	showTask bool
	obj      []string // объекты с которыми можно взаимодействовать с помощью команды применить (только дверь)
}

type Player struct {
	bag         bool
	currentRoom *Room
	inventory   []Item
	task        []string // задачи для выполнения
}

type Item struct {
	name   string
	action func(string) string // функция вызываемая при использовании команды "применить"
}

func JoinStrings(slice []string, del string) string {
	var result string
	for index, str := range slice {
		if index > 0 {
			result += del
		}
		result += str
	}
	return result
}

func (p *Player) LookAt() string {
	/*
		реализация команды "осмотреться"
	*/
	var answer []string
	if p.currentRoom.state != "" {
		answer = append(answer, p.currentRoom.state)
	}
	for obj, items := range p.currentRoom.items { // перечисление всех предметов
		var itemsName []string
		for _, item := range items {
			itemsName = append(itemsName, item.name)
		}
		answer = append(answer, fmt.Sprintf("на %s: %s", obj+"е", strings.Join(itemsName, ", ")))
	}
	if p.currentRoom.showTask {
		// answer = append(answer, "надо "+strings.Join(p.task, " и "))
		answer = append(answer, "надо "+JoinStrings(p.task, " и "))
	}
	if len(answer) > 0 {
		// return fmt.Sprintf("%s. %s", strings.Join(answer, ", "), FindPath(p.currentRoom))
		return fmt.Sprintf("%s. %s", JoinStrings(answer, ", "), FindPath(p.currentRoom))
	}
	return FindPath(p.currentRoom)
}

func (p *Player) PutOn(item string) string {
	/*
		реализация команды "надеть"
	*/
	if item != "рюкзак" {
		return "нельзя надеть"
	}
	k, v := SearchItem(p.currentRoom.items, item)
	if k == "" && v == -1 {
		return "нет такого"
	}
	Get(p.currentRoom.items, k, v)
	p.bag = true
	return "вы надели: " + item
}

func (p *Player) Take(item string) string {
	/*
		реализация команды "взять"
	*/
	if !p.bag {
		return "некуда класть"
	}
	k, v := SearchItem(p.currentRoom.items, item)
	if k == "" && v == -1 {
		return "нет такого"
	}
	sItem := Get(p.currentRoom.items, k, v)
	p.inventory = append(p.inventory, sItem)
	if len(p.currentRoom.items) == 0 && p.currentRoom.name == "комната" {
		p.currentRoom.state = "пустая комната"
	}
	return "предмет добавлен в инвентарь: " + item
}

func (p *Player) GoTo(room string) string {
	/*
		реализация команды "идти"
	*/
	for nextRoom, isOpen := range world[p.currentRoom.name] {
		if nextRoom.name == room {
			if !isOpen {
				return "дверь закрыта"
			}
			p.currentRoom = nextRoom
			return p.currentRoom.intro + " " + FindPath(p.currentRoom)
		}
	}
	return "нет пути в " + room
}

func (p *Player) ApplyTo(what string, to string) string {
	/*
		реализация команды "применить"
	*/
	item, inInventory := func(in string, items []Item) (*Item, bool) { // проверка на существование предмета в инвентаре
		for _, item := range items {
			if item.name == in {
				return &item, true
			}
		}
		return nil, false
	}(what, p.inventory)
	if !inInventory {
		return "нет предмета в инвентаре - " + what
	}
	inRoom := func(in string, objects []string) bool { // проверка на существование объекта, к которому можно применть предмет
		for _, o := range objects {
			if o == in {
				return true
			}
		}
		return false
	}(to, p.currentRoom.obj)
	if !inRoom {
		return "не к чему применить"
	}
	return item.action(to) // выполнение действия, для предмета при успехе
}

func SearchItem(mapItems map[string][]Item, request string) (string, int) {
	/*
		поиск предмета в map'e,
		возвращается ключ и индекс искомого элементав слайсе значений при успехе
	*/
	for obj, items := range mapItems {
		for index, item := range items {
			if item.name == request {
				return obj, index
			}
		}
	}
	return "", -1
}

func Get(mItems map[string][]Item, obj string, index int) Item {
	var item Item
	if mItems[obj][index].name == "рюкзак" { // выполнение одной из цели
		hero.TaskManager(mItems[obj][index].name)
	}
	item = mItems[obj][index]
	if len(mItems[obj]) == 1 {
		delete(mItems, obj)
		return item
	}
	mItems[obj] = append(mItems[obj][:index], mItems[obj][index+1:]...)
	return item
}

func FindPath(r *Room) string {
	/*
		перечислить все доступные комнаты
	*/
	var closestRooms []string
	if world[r.name] == nil {
		return ""
	}
	for room, _ := range world[hero.currentRoom.name] {
		if hero.currentRoom.name == "улица" { // если из улицы идти в коридор, то это теперь дом ¯\_(ツ)_/¯ непонятно почему
			room.name = "домой"
		}
		closestRooms = append(closestRooms, room.name)
	}
	// return "можно пройти - " + strings.Join(closestRooms, ", ")
	return "можно пройти - " + JoinStrings(closestRooms, ", ")
}

func (p *Player) TaskManager(completeTask string) {
	/*
		управление целями
	*/
	for index, quest := range p.task {
		if strings.Contains(quest, completeTask) {
			p.task = append(p.task[:index], p.task[index+1:]...)
		}
	}
}

func NewRoom(n, st, mes string, task bool) Room {
	var room Room
	room.name = n
	room.state = st
	room.intro = mes
	room.showTask = task
	return room
}

func NewPlayer(place *Room) Player {
	var player Player
	player.currentRoom = place
	return player
}

/*
	Глобальные переменные для упраления миром и персонажем
*/
var kitchen Room
var hall Room
var room Room
var street Room
var hero Player
var world = map[string]map[*Room]bool{
	"кухня": map[*Room]bool{
		&hall: true},
	"коридор": map[*Room]bool{
		&kitchen: true,
		&room:    true,
		&street:  true},
	"комната": map[*Room]bool{
		&hall: true},
	"улица": map[*Room]bool{
		&hall: true},
}

func main() {
	initGame()
	for {
		in := bufio.NewReader(os.Stdin)
		command, _ := in.ReadString('\n')
		command = command[:len(command)-1]
		if command == "exit" {
			return
		}
		handleCommand(command)
	}
}

func initGame() {
	/*
		Функция для проверки и открытия двери
	*/
	checkForDoor := func(door string, action func() string) func(string) string {
		return func(object string) string {
			if object == door {
				return action()
			}
			return "нельзя применить"
		}
	}
	openDoor := checkForDoor("дверь", func() string {
		for room, _ := range world[hero.currentRoom.name] {
			if room.name == "улица" {
				world[hero.currentRoom.name][room] = true
			}
		}
		return "дверь открыта"
	})

	/*
		Инициализация стартовых предметов с полями:
		name - имя комнаты
		action - функция вызываемая на команду применить
	*/
	tea := Item{"чай", nil}
	notes := Item{"конспекты", nil}
	backpack := Item{"рюкзак", nil}
	keys := Item{"ключи", openDoor}

	/*
		Инициализация игровых комнат
	*/
	kitchen = NewRoom("кухня", "ты находишься на кухне", "кухня, ничего интересного.", true)
	hall = NewRoom("коридор", "", "ничего интересного.", false)
	room = NewRoom("комната", "", "ты в своей комнате.", false)
	street = NewRoom("улица", "", "на улице весна.", false)
	kitchen.items = map[string][]Item{
		"стол": {tea},
	}
	room.items = map[string][]Item{
		"стол": {keys, notes},
		"стул": {backpack},
	}
	hall.obj = append(hall.obj, "дверь") // поставить дверь в коридоре
	world["коридор"][&street] = false    // закрыть дверь между коридором и улицей

	// Инициализация игрового персонажа
	hero = NewPlayer(&kitchen)
	hero.task = []string{"собрать рюкзак", "идти в универ"} // установка текущий цели
}

func handleCommand(command string) string {
	var answer string
	tokens := strings.Split(command, " ")
	action, params := tokens[0], tokens[1:]
	switch action {
	case "взять":
		answer = hero.Take(params[0])
	case "осмотреться":
		answer = hero.LookAt()
	case "идти":
		answer = hero.GoTo(params[0])
	case "надеть":
		answer = hero.PutOn("рюкзак")
	case "применить":
		answer = hero.ApplyTo(params[0], params[1])
	default:
		answer = "неизвестная команда"
	}
	// fmt.Printf("%s", answer+"\n") // для отображения в stdout
	return answer
}
