package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Room struct {
	name  string
	state string
	intro string
	items map[string][]string
	task  []string
	path  []*Room
}

type Player struct {
	bag         bool
	currentRoom *Room
	inventory   []string
}

func (p *Player) LookAt() string {
	var answer []string
	if p.currentRoom.state != "" {
		answer = append(answer, p.currentRoom.state)
	}
	for obj, items := range p.currentRoom.items {
		answer = append(answer, fmt.Sprintf("на %s: %s", obj+"е", strings.Join(items, ", ")))
	}
	if p.currentRoom.task != nil {
		answer = append(answer, "надо "+strings.Join(p.currentRoom.task, " и "))
	}
	return fmt.Sprintf("%s. %s", strings.Join(answer, ", "), FindPath(*p.currentRoom))
}

func (p *Player) PutOn(item string) string {
	if item != "рюкзак" {
		return "нельзя надеть"
	}
	k, v := SearchItem(p.currentRoom.items, item)
	if k == "" && v == 0 {
		return "нет такого"
	}
	Get(p.currentRoom.items, k, v)
	p.bag = true
	return "вы надели: " + item
}

func (p *Player) Take(item string) string {
	if !p.bag {
		return "некуда класть"
	}
	k, v := SearchItem(p.currentRoom.items, item)
	if k == "" && v == 0 {
		return "нет такого"
	}
	Get(p.currentRoom.items, k, v)
	p.inventory = append(p.inventory, item)
	return "предмет добавлен в инвентарь"
}

func (p *Player) GoTo(room string) string {
	for _, nextRoom := range p.currentRoom.path {
		if nextRoom.name == room {
			p.currentRoom = nextRoom
			return p.currentRoom.intro + " " + FindPath(*p.currentRoom)
		}
	}
	return "нет пути в" + room
}

func SearchItem(mapItems map[string][]string, request string) (string, int) {
	for obj, items := range mapItems {
		for index, item := range items {
			if item == request {
				return obj, index
			}
		}
	}
	return "", 0
}

func Get(mapItems map[string][]string, obj string, index int) {
	if len(mapItems[obj]) == 1 {
		delete(mapItems, obj)
		return
	}
	mapItems[obj] = append(mapItems[obj][:index], mapItems[obj][index+1:]...)
}

func FindPath(r Room) string {
	var closestRooms []string
	if r.path == nil {
		return ""
	}
	for _, room := range hero.currentRoom.path {
		closestRooms = append(closestRooms, (*room).name)
	}
	return "можно пройти - " + strings.Join(closestRooms, ", ")
}

func NewRoom(n, st, mes string, i map[string][]string) *Room {
	var room Room
	room.name = n
	room.items = i
	room.state = st
	room.intro = mes
	return &room
}

func NewPlayer(place *Room) *Player {
	var player Player
	player.currentRoom = place
	return &player
}

var kitchen = NewRoom("кухня", "ты находишься на кухне", "кухня, ничего интересного.", map[string][]string{"стол": {"чай"}})
var hall = NewRoom("коридор", "", "ничего интересного.", nil)
var room = NewRoom("комната", "", "ты в своей комнате.", map[string][]string{"стол": {"ключи", "конспекты"}, "стул": {"рюкзак"}})
var street = NewRoom("улица", "", "на улице весна.", nil)
var hero = NewPlayer(kitchen)
var world = map[string]*Room{
	"кухня":   kitchen,
	"коридор": hall,
	"комната": room,
	"улица":   street,
}

func main() {
	/*
		в этой функции можно ничего не писать
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/
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
		эта функция инициализирует игровой мир - все команты
		если что-то было - оно корректно перезатирается
	*/
	kitchen.path = []*Room{hall}
	hall.path = []*Room{kitchen, street, room}
	room.path = []*Room{hall}
	street.path = []*Room{hall}
	kitchen.task = []string{"собрать рюкзак", "идти в универ"}
}

func handleCommand(command string) string {
	/*
		данная функция принимает команду от "пользователя"
		и наверняка вызывает какой-то другой метод или функцию у "мира" - списка комнат
	*/
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
	default:
		fmt.Printf("неизвестная команда: %s\n", action)
	}
	fmt.Printf("%s", answer+"\n")
	return "not implemented"
}
