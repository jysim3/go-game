package controllers

// TODO: refreshing does not reveal current card - done
// TODO: option for additional cards
// TODO: current stage when deciding
// TODO: outcome of decide
// TODO: divider not aligned
// TODO: voting for reveal
// TODO: backdoor

import (
	"encoding/json"
	"fmt"
	"log"

	"jysim/game/models"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type PyramidController struct {
	m          *melody.Melody
	playerData map[string]*Player
	gameData   *PyramidGame
	cards      *models.CardsDeck
}

type Player struct {
	Name    string          `json:"name"`
	Cards   []models.Card   `json:"cards"`
	Id      string          `json:"id"`
	session *melody.Session `json:"-"`
}

type PyramidGame struct {
	Cards     []models.Card
	UsedCards []models.Card
	Round     int
	Target    *Target
}

type Target struct {
	From     string       `json:"from"`
	To       string       `json:"to"`
	Card     *models.Card `json:"card"`
	revealed bool
}

func (p Player) updateState() {
	ret := models.Command{
		"playerUpdate",
		p,
	}
	r, _ := json.Marshal(ret)
	p.session.Write(r)
}

func (h PyramidController) reset() {
	h.m.Close()
	h.m = melody.New()
}

func (h PyramidController) handleConnect(s *melody.Session) {
	sessionId, exists := s.Get("sessionId")
	if !exists || sessionId == nil {
		s.CloseWithMsg(melody.FormatCloseMessage( /* CloseUnsupportedData= */ 1003, "Refresh to get an id"))
		return
	}
	if player, exists := h.playerData[sessionId.(string)]; exists {
		player.session = s
		ret := models.Command{
			Code: "setName",
			Data: player.Name,
		}
		r, _ := json.Marshal(ret)
		s.Write(r)
	} else {
		player := Player{session: s, Id: sessionId.(string)}
		if h.gameData.Round > 0 {
			player.Cards = h.cards.DealCards(3)
		}
		h.playerData[sessionId.(string)] = &player
		ret := models.Command{
			Code: "setName",
			Data: player.Name,
		}
		r, _ := json.Marshal(ret)
		s.Write(r)
	}
	h.playerData[sessionId.(string)].updateState()
	h.updateNames()
	h.updateGame(false)
}

func (h PyramidController) handleDisconnect(s *melody.Session) {
	// h.updateNames()
}

func (h PyramidController) updateNames() {
	names := h.GetNames()
	ret := models.Command{
		Code: "names",
		Data: names,
	}
	r, _ := json.Marshal(ret)
	h.m.Broadcast(r)
}

func (h PyramidController) updateGame(forced bool) {
	if !forced && h.gameData.Round <= 0 {
		return
	}
	var target *Target
	if h.gameData.Target != nil {
		new_target := *h.gameData.Target
		if new_target.revealed != true {
			new_target.Card = nil
		}
		target = &new_target
	}

	ret := models.Command{
		Code: "game",
		Data: struct {
			CurrentCard models.Card `json:"currentCard"`
			Round       int         `json:"round"`
			Target      *Target     `json:"target"`
		}{h.gameData.Cards[h.gameData.Round], h.gameData.Round, target},
	}
	r, _ := json.Marshal(ret)
	h.m.Broadcast(r)
}

func (h PyramidController) GetNames() map[string]string {
	names := make(map[string]string)
	for id, player := range h.playerData {
		if !player.session.IsClosed() {
			names[id] = player.Name
		}
	}
	return names
}

func (h *PyramidController) restartGame() {
	// h.gameData = &PyramidGame{}

	h.cards.RefreshDeck( /* useDiscarded= */ false)
	for _, player := range h.playerData {
		player.Cards = h.cards.DealCards(3)
		player.updateState()
	}
	h.gameData.Cards = h.cards.DealCards(10)
	h.gameData.Round = 0
	h.updateGame( /* forced */ true)
}

func (h PyramidController) setName(id string, name string) {
	player := h.playerData[id]
	player.Name = name
	h.updateNames()
}

func (h PyramidController) sendCard(id string, d interface{}) {
	r, _ := json.Marshal(d)
	var payload struct {
		Target string      `json:"target"`
		Card   models.Card `json:"card"`
	}
	if err := json.Unmarshal(r, &payload); err != nil {
		log.Print(err)
		return
	}
	player := h.playerData[id]
	for index, card := range player.Cards {
		if payload.Card == card {
			player.Cards[index] = h.cards.Deal()
			player.updateState()
			if _, ok := h.playerData[payload.Target]; ok {
				new_card := payload.Card
				target := &Target{
					From:     id,
					To:       payload.Target,
					Card:     &new_card,
					revealed: false,
				}
				fmt.Println(target)
				h.gameData.Target = target
				h.updateGame(false)
			}
		}
	}
}

func (h PyramidController) handleMessage(s *melody.Session, msg []byte) {
	// name, _ := s.Get("name")
	var command models.Command
	id := s.MustGet("sessionId").(string)
	if err := json.Unmarshal(msg, &command); err != nil {
		log.Print(err)
		return
	}
	if command.Code == "reset" {
		h.reset()
	} else if command.Code == "start" {
		h.restartGame()
	} else if command.Code == "setName" {
		h.setName(id, command.Data.(string))
	} else if command.Code == "open" {
		h.gameData.Round += 1
		h.updateGame(false)
	} else if command.Code == "send" {
		h.sendCard(id, command.Data)
	} else if command.Code == "accept" {
		fmt.Println(h.gameData.Target)
		h.gameData.Target.revealed = true
		h.updateGame(false)
		h.gameData.Target = nil
	} else if command.Code == "reject" {
		h.gameData.Target = nil
		h.updateGame(false)
	} else if command.Code == "backdoor" {
	}
}

func (h PyramidController) HandleRequest(c *gin.Context) {
	if val, err := c.Cookie("sessionId"); err == nil {
		h.m.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{"sessionId": val})
		return
	}
	h.m.HandleRequest(c.Writer, c.Request)
}

func (h PyramidController) Close() {
	// h.m.Close()
}

func (h PyramidController) GetCount() int {
	return len(h.GetNames())
}

func NewPyramidController() func() RoomControllerInterface {
	return func() RoomControllerInterface {
		h := PyramidController{
			playerData: make(map[string]*Player),
			gameData:   &PyramidGame{},
			cards:      &models.CardsDeck{},
		}

		m := melody.New()
		h.m = m

		m.HandleConnect(h.handleConnect)
		m.HandleDisconnect(h.handleDisconnect)
		m.HandleMessage(h.handleMessage)

		return h
	}
}
