/*
 * Copyright 2020 PeerGum
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/go-chi/chi"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "net/http"
    "strconv"
)

var checkerService = Service{
    Info: ServiceInfo{
        Name:        "Checkers",
        Id:          checkerServiceID,
        Version:     0x0100, // version 1.00
        address:     checkerServiceAddress,
        Description: "A basic checker game",
        flags:       0,
    },
    init: initCheckers,
    api: ApiDefinition{
        path:   "/checkers",
        router: checkersRouter,
    },
}

const (
    BLANK = iota
    BLACK
    WHITE
)

const (
    EMPTY = iota
    PAWN
    KING
)

type CKGame struct {
    gorm.Model
    Name    string   `json:"name"`
    Rows    int      `json:"rows"`
    Columns int      `json:"columns"`
    Color   int      `json:"color"`
    Moves   []CKMove `gorm:"foreignKey:GameID",json:"moves"`
}

type CKMove struct {
    gorm.Model
    Color     int          `json:"color"`
    Piece     int          `json:"piece"`
    Positions []CKPosition `gorm:"foreignKey:MoveID",json:"positions"`
    GameID    int
    Game      CKGame `json:"game"`
}

type CKPosition struct {
    gorm.Model
    Row    int `json:"row"`
    Column int `json:"column"`
    Index  int `json:"play"`
    MoveID int
    Move   CKMove `json:"move"`
}

var Games []CKGame

func (position CKPosition) String() string {
    col := bytes.NewBufferString("ABCDEFGHIJKLMNOPQRSTUVWXYZ").Bytes()[position.Column]
    return fmt.Sprintf("%c%d", col, position.Row+1)
}

func (move CKMove) String() (output string) {
    if len(move.Positions) == 0 {
        return ""
    }
    position := move.Positions[0]
    output = position.String()
    for _, nextPosition := range move.Positions[1:] {
        var action string
        if position.distance(nextPosition) == 1 {
            action = "-"
        } else if position.distance(nextPosition) == 2 {
            action = "x"
        } else {
            action = ".."
        }
        var to string
        if move.Piece != KING && ((move.Color == BLACK && nextPosition.Row == 0) || (move.Color == WHITE && nextPosition.Row == nextPosition.Move.Game.Rows-1)) {
            to = "[" + nextPosition.String() + "]"
        } else {
            to = nextPosition.String()
        }
        output += fmt.Sprintf("%s%s", action, to)
        position = nextPosition
    }
    return output
}

type Move CKMove
type MoveDesc struct {
    Move Move   `json:"move"`
    Desc string `json:"desc"`
}

func (move CKMove) MarshalJSON() ([]byte, error) {
    return json.Marshal(MoveDesc{
        Move: Move(move),
        Desc: move.String(),
    })
}

func (pos1 CKPosition) distance(pos2 CKPosition) int {
    return abs(pos1.Row - pos2.Row)
}

func abs(x int) int {
    if x > 0 {
        return x
    }
    return -x
}

func initCheckers() {
    Logln("Checkers Started")

    db.AutoMigrate(&CKGame{}, &CKPosition{}, &CKMove{})

    result := db.Find(&Games)
    if result.RowsAffected == 0 {
        //    var games []CKGames
        //    for i := 0; i < 10; i++ {
        //        messages = append(messages, Mail{
        //            Subject: fmt.Sprintf("Welcome to your new mailbox - msg %d", i),
        //            Body:    "Your new mailboxes were created and are ready for you to enjoy.\nWe hope you'll have good fun.",
        //        })
        //    }
        //    mailboxes := []Mailbox{
        //        Mailbox{Name: "Inbox",
        //            Mails: messages,
        //        },
        //        Mailbox{Name: "Outbox"},
        //        Mailbox{Name: "Sent"},
        //        Mailbox{Name: "Trash"},
        //    }
        //    db.Create(mailboxes)
        //    db.Commit()
        //    _ = db.Find(&mailDB)
        //
    } else {
        //    fmt.Println(inbox)
    }
}

func checkersRouter(r chi.Router) {
    r.Get("/", checkerListGames)
    //r.Post("/",checkerNewGame)
    r.Get("/{gameID}", checkerLoadGame)
    r.Get("/{gameID}/board", checkerLoadBoard)
    r.Get("/{gameID}/board/{moveNum}", checkerLoadBoard)
    //r.Post("/{gameID}/move",checkerMove)
}

func checkerListGames(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-type", "text/html")

    Logln("Game list requested")

    var games []CKGame
    db.Find(&games)
    if err := SendJson(w, games); err != nil {
        err := SendPage(w, "errors/500.html")
        if err != nil {
            w.WriteHeader(500)
            w.Write([]byte(err.Error()))
            Logln("Page Error: ", err)
        }
    }
}

func checkerLoadGame(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-type", "text/html")

    if gameID := chi.URLParam(req, "gameID"); gameID != "" {
        Loglnf("Game %d requested", gameID)
        var game CKGame

        db.Preload("Moves.Positions").Preload(clause.Associations).Find(&game, gameID)
        if err := SendJson(w, game); err != nil {
            err := SendPage(w, "errors/500.html")
            if err != nil {
                w.WriteHeader(500)
                w.Write([]byte(err.Error()))
                Logln("Page Error: ", err)
            }
        }
    } else {
        notFoundHandler(w, req)
    }
}

func checkerLoadBoard(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-type", "text/html")

    if gameID, err := strconv.Atoi(chi.URLParam(req, "gameID")); err == nil {
        Loglnf("Game %d, board requested", gameID)
        var game CKGame

        db.Preload("Moves.Positions").Preload(clause.Associations).Find(&game, gameID)
        if gameID == 0 {
            var rows int
            var columns int
            var err error
            if rows, err = strconv.Atoi(req.URL.Query().Get("r")); err != nil {
                rows = 8
            }
            if columns, err = strconv.Atoi(req.URL.Query().Get("c")); err != nil {
                columns = 8
            }
            game = CKGame{
                Color:   WHITE,
                Rows:    rows,
                Columns: columns,
                Moves:   []CKMove{},
            }
        }
        type Piece struct {
            Color int `json:"color"`
            Type  int `json:"type"`
        }
        type Row []Piece

        // Generate base board
        board := make([][]Piece, game.Rows)
        jail := make(map[int]int, 2)
        for rowNum := range board {
            board[rowNum] = make([]Piece, game.Columns)
            for colNum := range board[rowNum] {
                if rowNum < game.Rows/2-1 && (rowNum+colNum)%2 == 0 {
                    board[rowNum][colNum] = Piece{
                        Color: WHITE,
                        Type:  PAWN,
                    }
                } else if rowNum > game.Rows/2 && (rowNum+colNum)%2 == 0 {
                    board[rowNum][colNum] = Piece{
                        Color: BLACK,
                        Type:  PAWN,
                    }
                } else {
                    board[rowNum][colNum] = Piece{
                        Color: BLANK,
                        Type:  EMPTY,
                    }
                }
            }
        }

        if len(game.Moves) > 0 {
            // update moves

            var moveNum int
            var err error
            if moveNum, err = strconv.Atoi(chi.URLParam(req, "moveNum")); err != nil {
                moveNum = len(game.Moves)
            }
            fmt.Print("movenum=", moveNum)
            for _, move := range game.Moves[:moveNum] {
                newPiece := Piece{
                    Color: move.Color,
                    Type:  move.Piece,
                }
                if len(move.Positions) == 0 {
                    continue
                }
                startPosition := move.Positions[0]
                for _, position := range move.Positions[1:] {
                    // pawn become a king
                    if (position.Row == game.Rows-1 && move.Color == WHITE) ||
                        (position.Row == 0 && move.Color == BLACK) {
                        newPiece.Type = KING
                    }
                    rowIncrease := 1
                    columnIncrease := 1
                    if position.Row < startPosition.Row {
                        rowIncrease = -1
                    }
                    if position.Column < startPosition.Column {
                        columnIncrease = -1
                    }
                    i := 0
                    j := 0
                    for i != position.Row-startPosition.Row {
                        if board[startPosition.Row+i][startPosition.Column+j].Color != move.Color {
                            jail[move.Color]++
                        }
                        board[startPosition.Row+i][startPosition.Column+j] = Piece{
                            Color: BLANK,
                            Type:  EMPTY,
                        }
                        i += rowIncrease
                        j += columnIncrease
                    }
                    board[position.Row][position.Column] = newPiece
                    startPosition = position
                }
            }

        }
        type GameStatus struct {
            Board [][]Piece   `json:"board"`
            Jail  map[int]int `json:"jail"`
        }
        status := GameStatus{
            Board: board,
            Jail:  jail,
        }
        // send board
        if err := SendJson(w, status); err != nil {
            err := SendPage(w, "errors/500.html")
            if err != nil {
                w.WriteHeader(500)
                w.Write([]byte(err.Error()))
                Logln("Page Error: ", err)
            }
        }
    } else {
        notFoundHandler(w, req)
    }
}
