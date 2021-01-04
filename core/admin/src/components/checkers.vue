<!--
  - Copyright 2020 PeerGum
  -
  - Licensed under the Apache License, Version 2.0 (the "License");
  - you may not use this file except in compliance with the License.
  - You may obtain a copy of the License at
  -
  -     http://www.apache.org/licenses/LICENSE-2.0
  -
  - Unless required by applicable law or agreed to in writing, software
  - distributed under the License is distributed on an "AS IS" BASIS,
  - WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  - See the License for the specific language governing permissions and
  - limitations under the License.
  -->

<template>
  <div>
    <div class="row align-items-start justify-content-center">
      <div class="col form align-items-start action-box">
        <!--        <form class="game-options">-->
        <div class="form-group">
          <label for="games">Games</label>
          <select id="games" class="form-select" v-model="gameId">
            <option value="0" :selected="games.length === 0">New</option>
            <option v-for="game in games" :selected="gameId === game.ID" :value="game.ID">{{ game.name }}</option>
          </select>
        </div>
        <div class="form-row">
          <div class="form-group col">
            <label for="rows">Rows</label>
            <select id="rows" v-model="rows" :disabled="gameId > 0">
              <option v-for="numRows in 4" :selected="(6+numRows*2) === rows" :value="6+numRows*2">{{
                  6 + numRows * 2
                }}
              </option>
            </select>
          </div>
          <div class="form-group col">
            <label for="columns">Columns</label>
            <select id="columns" v-model="columns" :disabled="gameId > 0">
              <option v-for="numColumns in 9" :selected="(7+numColumns) === columns" :value="7 + numColumns">{{
                  7 + numColumns
                }}
              </option>
            </select>
          </div>
        </div>
        <div class="form-group">
          <label for="color">Color</label>
          <select id="color" v-model="game.color" :disabled="gameId > 0">
            <option v-for="color in colors" :selected="color.id === game.color" :value="color.id">{{
                color.name
              }}
            </option>
          </select>
        </div>
        <!--        </form>-->
      </div>
      <table class="checkers justify-content-center align-items-center">
        <thead>
        <td class="v-board-side h-board-side"></td>
        <td v-for="col in game.columns" class="h-board-side upside-down-text">
          {{ game.color == BLACK ? 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.substr(game.columns - col, 1) : " " }}
        </td>
        <td class="v-board-side h-board-side"></td>
        </thead>
        <tbody>
        <tr v-for="row in game.rows">
          <td class="v-board-side">{{ game.color == WHITE ? game.rows - row + 1 : " " }}</td>
          <td v-for="col in game.columns"
              :style="{width:(60-2*(game.rows-8))+'px',height:((60-2*(game.rows-8)))+'px'}"
              :id="id(game.rows-row,col-1)" :class="square(  game.rows-row,col-1)"
              v-on:click="moveTo(game.rows-row,col-1)">
            <!--            <div v-if="boardReady">{{ board[boardHeight-row][col-1] }}</div>-->
            <img
                :style="{width:(60-2*(game.rows-8))+'px',height:((60-2*(game.rows-8)))+'px'}"
                v-if="boardReady/* && piece(game.rows-row,col-1) != ''*/"
                :src="piece(game.rows-row,col-1)"/>
            <!--            <img :class="{movable:row==6}" v-if="row>5 && (col+row)%2 == 1" src="images/white-pawn.png"/>-->
            <!--            <div v-if="boardReady && moves[game.rows-row][col-1]>0" class="moves">{{-->
            <!--                moves[game.rows - row][col - 1]-->
            <!--              }}-->
            <!--            </div>-->
          </td>
          <td class="v-board-side upside-down-text">{{ game.color == BLACK ? row : " " }}</td>
        </tr>
        </tbody>
        <thead>
        <td class="v-board-side h-board-side"></td>
        <td v-for="col in game.columns" class="h-board-side">
          {{ game.color == WHITE ? "ABCDEFGHIJKLMNOPQRSTUVWXYZ".substr(col - 1, 1) : " " }}
        </td>
        <td class="v-board-side h-board-side"></td>
        </thead>
      </table>
      <div class="col move-list">
        <h4>Moves</h4>
        <ul v-on:mouseleave="getBoard(gameId)">
          <li v-for="(move,item) in game.Moves" v-on:mouseover="showBoardAt(item)">{{ move.desc }}
            <div class="color">{{ item % 2 == 0 ? "W" : "B" }}</div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
const BLANK = 0, BLACK = 1, WHITE = 2;
const EMPTY = 0, PAWN = 1, KING = 2;

export default {
  // name: "messaging",
  props: [],
  data() {
    return {
      colors: [
        {
          id: BLACK,
          name: "Black",
        },
        {
          id: WHITE,
          name: "White",
        },
      ],
      BLACK: BLACK,
      WHITE: WHITE,
      SIMPLE_MOVE: 0,
      JUMP_MOVE: 1,
      board: [],
      moves: [],
      boardUpdated: false,
      boardReady: false,
      myColor: WHITE,
      myTurn: true,
      moving: false,
      currentRow: -1,
      currentCol: -1,
      singlePlayer: true,
      currentAction: "Starting...",
      gameId: 0,
      games: [],
      game: {},
      defaultGame: {
        color: WHITE,
        columns: 8,
        rows: 8,
        Moves: [],
      },
      movelist: [],
      rows: 8,
      columns: 8,
      shownMove: -1,
      marked: [],
      currentMarked: 0,
      markedTimer: null,
    }
  },
  watch: {
    boardReady: function (value) {
      if (value) {
        this.boardUpdated = true;
        // this.refreshBoard();
      }
    },
    gameId: function (value) {
      this.getMoves(value);
    },
    rows: function (value) {
      this.game.rows = value;
      this.getBoard(this.gameId);
    },
    columns: function (value) {
      this.game.columns = value;
      this.getBoard(this.gameId);
    },
  },
  mounted() {
    this.getGames();
  },
  updated() {
    this.$nextTick(function () {
      if (this.boardUpdated) {
        console.log("refreshed");
        // this.render()
        this.boardUpdated = false;
        // this.refreshBoard();
      }
    });
  },
  methods: {
    showBoardAt: function(numMove) {
      if (this.boardReady) {
        console.log(numMove);
        this.getBoard(this.gameId,numMove);
        this.shownMove = numMove;
        this.marked = this.game.Moves[this.shownMove].move.Positions;
        var vm = this;
        if (this.markedTimer) {
          clearInterval(this.markedTimer);
          this.markedTimer = null;
        }
        this.markedTimer = setInterval(function() {
          vm.currentMarked = (vm.currentMarked + 1) % vm.marked.length;
        },200);
      }
    },
    getGames: function () {
      var vm = this;
      this.games = [];
      this.gameId = 0;
      axios.get('/checkers')
          .then(function (response) {
            vm.games = response.data;
            if (vm.games.length > 0) {
              vm.gameId = vm.games[vm.games.length - 1].ID;
            }
          })
          .catch(function (error) {
            console.log(error);
          })
          .then(function () {
            // always executed
          });
      // this.setupBoard();
    },
    getMoves: function (gameId) {
      if (gameId == 0) {
        this.game = this.defaultGame;
        this.game.Moves = [];
        this.getBoard(gameId);
        return;
      }
      console.log("requesting game " + gameId);
      var vm = this;
      axios.get('/checkers/' + gameId)
          .then(function (response) {
            vm.game = response.data;
            vm.getBoard(gameId)
          })
          .catch(function (error) {
            console.log(error);
          })
          .then(function () {
            // always executed
          });
    },
    getBoard: function (gameId, moveNum) {
      console.log("requesting board for game " + gameId);
      var vm = this;
      let url = '/checkers/' + gameId + '/board';
      let query = '';
      if (gameId == 0) {
        query += '?r=' + this.rows + '&c=' + this.columns;
      } else if (moveNum !== undefined) {
        url += '/' + moveNum
      } else {
        this.shownMove = -1;
        if (this.markedTimer) {
          clearInterval(this.markedTimer);
          this.markedTimer = null;
        }
      }
      axios.get(url + query)
          .then(function (response) {
            vm.board = response.data.board;
            vm.jail = response.data.jail;
            vm.boardReady = true;
          })
          .catch(function (error) {
            console.log(error);
          })
          .then(function () {
            // always executed
          });
    },
    row: function (r) {
      return this.game.color == WHITE ? r : this.game.rows - r - 1;
    },
    col: function (c) {
      return this.game.color == WHITE ? c : this.game.columns - c - 1;
    },
    id: function (i, j) {
      return 'r' + this.row(i) + 'c' + this.col(j);
    },
    square: function (r, c) {
      let i = this.row(r);
      let j = this.col(c);
      if (!this.boardReady) {
        return {
          "dark": (i + j) % 2 == 0,
          "light": (i + j) % 2 == 1,
        }
      }
      let marked = false;
      if (this.shownMove >= 0) {
        let position = this.marked[this.currentMarked];
        if (i == position.row && j == position.column) {
          marked = true;
        }
      }
      return {
        "dark": (i + j) % 2 == 0,
        "light": (i + j) % 2 == 1,
        "piece-moved": marked,
        // "piece-to-move": i == this.currentRow && j == this.currentCol,
        // "targetable": this.moves[i][j] > 0,
      };
    },
    piece: function (r, c) {
      console.log("piece for " + r + "," + c);
      let i = this.row(r);
      let j = this.col(c);
      let name = "/images/";
      switch (this.board[i][j].color) {
        case BLACK:
          name += "black-";
          break;
        case WHITE:
          name += "white-";
          break;
        default:
          name += "blank";
      }
      switch (this.board[i][j].type) {
        case PAWN:
          name += "pawn";
          break;
        case KING:
          name += "king";
          break;
        default:
          break;
      }
      name += ".png";
      return name;
    },
    clearMoves: function () {
      for (let i = 0; i < this.boardHeight; i++) {
        for (let j = 0; j < this.boardWidth; j++) {
          this.moves[i][j] = 0;
        }
      }
    },
    checkPosition: function (i, j) {
      if (moveType == this.SIMPLE_MOVE && this.board[i][j] == EMPTY) {
        $('#' + this.id(i, j)).each(function () {
          $(this).addClass('target_' + moveIndex);
        });
        return true;
      }
      let myPawn = this.myColor == this.WHITE ? this.WHITE_PAWN : this.BLACK_PAWN;
      let myKing = this.myColor == this.WHITE ? this.WHITE_KING : this.BLACK_KING;
      let otherPawn = this.myColor == this.BLACK ? this.WHITE_PAWN : this.BLACK_PAWN;
      let otherKing = this.myColor == this.BLACK ? this.WHITE_KING : this.BLACK_KING;
      return false;
    },
  }
}
</script>
