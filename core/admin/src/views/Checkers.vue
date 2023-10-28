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

<script setup>
import axios from 'axios'
import {computed, nextTick, onMounted, onUpdated, ref, watch} from "vue";

const BLANK = 0, BLACK = 1, WHITE = 2;
const EMPTY = 0, PAWN = 1, KING = 2;
const SIMPLE_MOVE = 0,
    JUMP_MOVE = 1;

const DEFAULT_GAME = {
  color: WHITE,
  columns: 8,
  rows: 8,
  Moves: [],
}

const colors = ref([
  {
    id: BLACK,
    name: "Black",
  },
  {
    id: WHITE,
    name: "White",
  },
])
// const boardWidth = ref(8);
// const boardHeight = ref(8);
const board = ref([]);
const moves = ref([])
const boardUpdated = ref(false)
const boardReady = ref(false)
const myColor = ref(WHITE)
const myTurn = ref(true)
const moving = ref(false)
const currentRow = ref(-1)
const currentCol = ref(-1)
const singlePlayer = ref(true)
const currentAction = ref("Starting...")
const gameId = ref(0)
const games = ref([])
const game = ref(DEFAULT_GAME)
const jail = ref([])
const moveType = ref(SIMPLE_MOVE)
const movelist = ref([])
const rows = ref(8)
const columns = ref(8)
const shownMove = ref(-1)
const marked = ref([])
const currentMarked = ref(0)
const markedTimer = ref(null)

watch(boardReady, (value) => {
  if (value) {
    boardUpdated.value = true;
    // this.refreshBoard();
  }
})

watch(gameId, (value) => {
  getMoves(value);
})

watch(rows, (value) => {
  game.value.rows = value;
  getBoard(gameId.value);
})

watch(columns, (value) => {
  game.value.columns = value;
  getBoard(gameId.value);
})

watch(myColor, (value) => {
  console.log("color swap")
})

onMounted(() => {
  resetBoard()
  getGames();
  getBoard(gameId.value)
  game.value = {
    columns: columns.value,
    rows: rows.value,
  }
  boardUpdated.value = true
  boardReady.value = false;
})

onUpdated(() => {
  nextTick(() => {
    if (boardUpdated.value) {
      console.log("refreshed");
      // render()
      boardUpdated.value = false;
      boardReady.value = true;
      // this.refreshBoard();
    }
  });
})

const resetBoard = () => {
  for (let row = 0; row < rows.value; row++) {
    if (board.value == null) {
      board.value = [];
    }
    board.value[row] = [];
    for (let col = 0; col < columns.value; col++) {
      board.value[row][col] = {
        color: ((row + col) % 2 == 1 || (row > 2 && row < rows.value - 3)) ? BLANK : (row < 3 ? (3 - myColor.value) : myColor.value),
        type: (row < 3 || row > rows.value - 4) ? PAWN : EMPTY,
      };
    }
  }
  console.log("reset")
  // console.log(board.value)
}
const showBoardAt = (numMove) => {
  if (boardReady.value) {
    console.log(numMove);
    getBoard(gameId, numMove);
    shownMove.value = numMove;
    marked.value = game.value.Moves[shownMove.value].move.Positions;
    if (markedTimer.value) {
      clearInterval(markedTimer.value);
      markedTimer.value = null;
    }
    if (marked.value.length > 0) {
      currentMarked.value = 0;
      markedTimer.value = setInterval(() => {
        currentMarked.value = (currentMarked.value + 1) % marked.value.length;
      }, 200);
    }
  }
}

const getGames = () => {
  games.value = [];
  gameId.value = 0;
  axios.get('/checkers/games')
      .then(function (response) {
        games.value = response.data;
        if (games.value.length > 0) {
          gameId.value = games.value[games.value.length - 1].ID;
        }
      })
      .catch(function (error) {
        console.log(error);
      })
      .then(function () {
        // always executed
      });
  // setupBoard();
}

const getMoves = (gId) => {
  if (gId == 0) {
    game.value = DEFAULT_GAME.value;
    game.value.Moves = [];
    getBoard(gId);
    return;
  }
  console.log("requesting game " + gId);

  axios.get('/checkers/' + gId)
      .then(function (response) {
        game.value = response.data;
        getBoard(gameId)
      })
      .catch(function (error) {
        console.log(error);
      })
      .then(function () {
        // always executed
      });
}

const getBoard = (gId, moveNum = 0) => {
  console.log("requesting board for game " + gId);

  let url = '/checkers/games/' + gId + '/board';
  let query = '';
  if (gId == 0) {
    query += '?r=' + rows.value + '&c=' + columns.value;
  } else if (moveNum !== 0) {
    url += '/' + moveNum
  } else {
    shownMove.value = -1;
    if (markedTimer.value) {
      clearInterval(markedTimer.value);
      markedTimer.value = null;
    }
  }
  axios.get(url + query)
      .then(function (response) {
        board.value = response.data.board;
        jail.value = response.data.jail;
        boardReady.value = true;
      })
      .catch(function (error) {
        console.log(error);
        resetBoard()
      })
      .then(function () {
        // always executed
      });
}

const row = (r) => {
  // console.log(game.value)
  return game.value.color == WHITE ? r : game.value.rows - r - 1;
}

const col = (c) => {
  // console.log(game.value)
  return game.value.color == WHITE ? c : game.value.columns - c - 1;
}

const id = (i, j) => {
  return 'r' + row(i) + 'c' + col(j);
}

const square = (r, c) => {
  let i = row(r);
  let j = col(c);
  if (!boardReady.value) {
    return {
      "dark": (i + j) % 2 == 0,
      "light": (i + j) % 2 == 1,
    }
  }
  let marked = false;
  if (shownMove.value >= 0) {
    let position = marked.value[currentMarked.value];
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
}

const piece = (r, c) => {
  // console.log("piece for " + r + "," + c);
  let i = row(myColor.value == WHITE ? r : (rows.value - r - 1));
  let j = col(myColor.value == WHITE ? c : (columns.value - c - 1));
  let name = "/images/";
  switch (board.value[i][j].color) {
    case BLACK:
      name += "black-";
      break;
    case WHITE:
      name += "white-";
      break;
    default:
      name += "blank";
  }
  if (board.value[i][j].color != BLANK) {
    switch (board.value[i][j].type) {
      case PAWN:
        name += "pawn";
        break;
      case KING:
        name += "king";
        break;
      default:
        break;
    }
  }
  name += ".png";
  return name;
}

const clearMoves = () => {
  for (let i = 0; i < rows.value; i++) {
    for (let j = 0; j < columns.value; j++) {
      moves.value[i][j] = 0;
    }
  }
}

const checkPosition = (i, j) => {
  if (moveType.value == SIMPLE_MOVE && board.value[i][j] == EMPTY) {
    // $('#' + this.id(i, j)).each(function () {
    //   $(this).addClass('target_' + moveIndex.value);
    // });
    return true;
  }
  let myPawn = myColor.value == WHITE ? WHITE_PAWN : BLACK_PAWN;
  let myKing = myColor.value == WHITE ? WHITE_KING : BLACK_KING;
  let otherPawn = myColor.value == BLACK ? WHITE_PAWN : BLACK_PAWN;
  let otherKing = myColor.value == BLACK ? WHITE_KING : BLACK_KING;
  return false;
}

const moveTo = (x, y) => {

}

const checkerClass = computed(() => {
  return (columns.value < 11) ?
      "flex flex-col lg:flex-row align-items-start justify-content-center" :
      "flex lg:flex-row flex-wrap align-items-start justify-content-center";
})
const gamesTopClass = computed(() => {
  return (columns.value < 11) ?
      "w-3/4 mx-auto lg:w-1/6 form align-items-start action-box" :
      "hidden w-1/2 mx-auto form align-items-center action-box";
})
const gamesBoardClass = computed(() => {
  return (columns.value < 11) ?
      "hidden w-full lg:w-1/6 form align-items-start action-box" :
      "w-3/4 mx-auto form align-items-center action-box";
})
const boardClass = computed(() => {
  return (columns.value < 11) ?
      "w-full lg:w-1/2 px-4" :
      "w-full mx-auto max-w-3xl lg:w-2/3 px-4";
})
const sideClass = computed(() => {
  return (columns.value < 11) ?
      "w-full lg:w-1/3 move-list" :
      "w-full lg:w-1/4 move-list";
})

</script>

<template>
  <div :class='checkerClass'>
    <div :class="gamesTopClass">
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
        <select id="color" v-model="myColor" :disabled="gameId > 0">
          <option v-for="color in colors" :selected="color.id === myColor" :value="color.id">{{
              color.name
            }}
          </option>
        </select>
      </div>
      <!--        </form>-->
    </div>
    <div :class="boardClass">
      <div :class="gamesBoardClass">
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
          <select id="color" v-model="myColor" :disabled="gameId > 0">
            <option v-for="color in colors" :selected="color.id === myColor" :value="color.id">{{
                color.name
              }}
            </option>
          </select>
        </div>
        <!--        </form>-->
      </div>
      <table class="w-fit mx-auto checkers justify-content-center align-items-center">
        <thead>
        <td class="v-board-side h-board-side"></td>
        <td v-for="col in game.columns" class="h-board-side upside-down-text">
          {{
            myColor == BLACK ? 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.substring(game.columns - col, game.columns - col + 1) : " "
          }}
        </td>
        <td class="v-board-side h-board-side"></td>
        </thead>
        <tbody>
        <tr class="" v-for="row in game.rows">
          <td class="v-board-side">{{ myColor == WHITE ? game.rows - row + 1 : " " }}</td>
          <td class="flex-1" v-for="col in game.columns"
              :id="id(game.rows-row,col-1)" :class="square(  game.rows-row,col-1)"
              v-on:click="moveTo(game.rows-row,col-1)">
            <!--            <div v-if="boardReady">{{ board[boardHeight-row][col-1] }}</div>-->
            <img
                v-if="boardReady/* && piece(game.rows-row,col-1) != ''*/"
                :src="piece(game.rows-row,col-1)"/>
            <!--            <img :class="{movable:row==6}" v-if="row>5 && (col+row)%2 == 1" src="images/white-pawn.png"/>-->
            <!--            <div v-if="boardReady && moves[game.rows-row][col-1]>0" class="moves">{{-->
            <!--                moves[game.rows - row][col - 1]-->
            <!--              }}-->
            <!--            </div>-->
          </td>
          <td class="v-board-side upside-down-text">{{ myColor == BLACK ? row : " " }}</td>
        </tr>
        </tbody>
        <thead>
        <td class="v-board-side h-board-side"></td>
        <td v-for="col in game.columns" class="h-board-side">
          {{ myColor == WHITE ? "ABCDEFGHIJKLMNOPQRSTUVWXYZ".substr(col - 1, 1) : " " }}
        </td>
        <td class="v-board-side h-board-side"></td>
        </thead>
      </table>
    </div>
    <div :class="sideClass">
      <h4>Moves</h4>
      <ul v-on:mouseleave="getBoard(gameId)">
        <li v-for="(move,item) in game.Moves" v-on:mouseover="showBoardAt(item)">{{ move.desc }}
          <div class="color">{{ item % 2 == 0 ? "W" : "B" }}</div>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>


.action-box {
  @apply p-4;
  /* width: 15 %;*/
  /*
  margin: 20px;
  padding: 5px;
  border: 1px solid lightgrey;
  */
  /* min-height: 100 vh;*/

  select {
    width: 100%;

    &.games {
      size: 5px;
    }
  }
}

.move-list {
  @apply mt-4 pl-4 lg:border-l border-black border-dotted;
  /* width: 15 %;*/
  /*
  margin: 20px;
  padding: 5px;
  border: 1px solid lightgrey;
  */
  /* min-height: 100 vh;*/

  ul {
    list-style: decimal;
    font-family: 'Special Elite';
    color: #0040a0;

    li {
      div.color {
        float: right;
      }

      &:hover {
        color: red;
        background-color: orange;
      }
    }
  }
}

.checkers, .chess {
  @apply w-full max-w-xl mt-4 bg-[url('@/assets/images/leather.png')] bg-repeat;

  .h-board-side {
    @apply h-8 font-sans drop-shadow font-bold text-yellow-200;

    /*
    &.upside-down-text {
      transform: rotate(180deg);
    }
    */
  }

  .v-board-side {
    @apply w-8 font-sans drop-shadow font-bold text-yellow-200;

    /*
    &.upside-down-text {
      transform: rotate(180deg);
    }
    */
  }

  td {
    @apply w-fit h-fit p-0 items-center text-center;
    /* width: 60 px;*/
    /* height: 60 px;*/
    /*
    vertical-align: middle;
    text-align: center;
    */

    &.dark {
      @apply bg-cover bg-[url('@/assets/images/dark-wood.png')];
    }

    &.light {
       @apply bg-cover bg-[url('@/assets/images/light-wood.png')];
    }

    &.targetable {
      background: orange;
    }

    &.piece-to-move {
      background: yellow;
    }

    &.piece-moved {
      background: orange;
    }

    img {
      @apply w-fit h-auto mx-auto;
    //width: 60px;
      /* height: 60 px;*/

      /*
      &.upside-down {
        transform: rotate(180deg);
      }
      */
    }

    div.moves {
      font-size: 20px;
      font-weight: bold;
    }
  }
}

</style>
