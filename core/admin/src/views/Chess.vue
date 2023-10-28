<!--
  - Copyright 2023 PeerGum
  -
  - Licensed under the Apache License, Version 2.0 (the "License");
  - you may not use file.value except in compliance with the License.
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
import {nextTick, onMounted, onUpdated, ref, watch} from "vue";

const BLANK = 0, BLACK = 1, WHITE = 2;
const EMPTY = 0, PAWN = 1, ROOK = 2, KNIGHT = 3, BISHOP = 4, QUEEN = 5, KING = 6;
const initialPieces = [2, 3, 4, 5, 6, 4, 3, 2];

const SIMPLE_MOVE = 0, JUMP_MOVE = 1;

const boardWidth = ref(8);
const boardHeight = ref(8);
const pawnRows = ref(4);
const board = ref([]);
const moves = ref([]);
const boardUpdated = ref(false);
const boardReady = ref(false);
const myColor = ref(WHITE);
const myTurn = ref(true);
const moving = ref(false);
const currentRow = ref(-1);
const currentCol = ref(-1);
const singlePlayer = ref(true);
const currentAction = ref("Starting...");

watch(boardReady, (value) => {
  if (value) {
    boardUpdated.value = false;
    // refreshBoard();
  }
})

const pieces = [
  'king',
  'queen',
  'knight',
  'rook',
  'bishop',
  'pawn',
]

// const image = ref({})

onMounted(() => {
  // image.value = {};
  // pieces.forEach((v) => {
  //   let black = 'b_' + v
  //   image.value[black] = new Image()
  //   image.value[black].src = new URL('/images/' + white + ".png", import.meta.url).href;
  //   let white = 'w_' + v
  //   image.value[white] = new Image()
  //   image.value[white].src = new URL('/images/' + white + ".png", import.meta.url).href;
  // })
  // console.log(image.value)
  // setTimeout(200,setupBoard())
  setupBoard();
})


onUpdated(() => {
  nextTick(function () {
    if (boardUpdated.value) {
      console.log("refreshed");
      // refreshBoard();
      boardUpdated.value = false;
    }
  });
})

const setupBoard = () => {
  for (let i = 0; i < boardHeight.value; i++) {
    board.value[i] = [];
    moves.value[i] = [];
    for (let j = 0; j < boardWidth.value; j++) {
      board.value[i][j] = {
        color: BLANK,
        piece: EMPTY,
      }
      moves.value[i][j] = false;
      let color = WHITE;
      if (i > boardHeight.value / 2) {
        color = BLACK;
      }
      if (i == 1 || i == boardHeight.value - 2) {
        board.value[i][j] = {
          color: color,
          piece: PAWN,
        };
      }
      if (i == 0 || i == boardHeight.value - 1) {
        board.value[i][j] = {
          color: color,
          piece: initialPieces[j],
        };
      }
    }
  }
  boardReady.value = true;
  nextTurn();
  // console.log(board.value);
  // boardUpdated.value = true;
  // refreshBoard();
}

const nextTurn = () => {
}

const row = (r) => {
  return myColor.value == WHITE ? r : boardHeight.value - r - 1;
}
const col = (c) => {
  return myColor.value == WHITE ? c : boardWidth.value - c - 1;
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
  return {
    "dark": (i + j) % 2 == 0,
    "light": (i + j) % 2 == 1,
    "piece-to-move": i == currentRow.value && j == currentCol.value,
    "targetable": moves.value[i][j] > 0,
  };
}

const piece = (r, c) => {
  let i = row(r);
  let j = col(c);
  let name = "";
  switch (board.value[i][j].color) {
    case BLACK:
      name += "b-";
      break;
    case WHITE:
      name += "w-";
      break;
    default:
      break;
  }
  switch (board.value[i][j].piece) {
    case PAWN:
      name += "pawn";
      break;
    case ROOK:
      name += "rook";
      break;
    case KNIGHT:
      name += "knight";
      break;
    case BISHOP:
      name += "bishop";
      break;
    case QUEEN:
      name += "queen";
      break;
    case KING:
      name += "king";
      break;
    default:
      name += "blank";
      break;
  }
  // name += ".png";
  console.log(name)
  const image = new Image();
  image.src = new URL('/images/' + name + '.png', import.meta.url).href;
  return image.src;//new URL('/images/'+name+'.png', import.meta.url).href; //image.value[name];
}

const clearMoves = () => {
  for (let i = 0; i < boardHeight.value; i++) {
    for (let j = 0; j < boardWidth.value; j++) {
      moves.value[i][j] = 0;
    }
  }
}

//   move: function (row, col) {
//     console.log(row, col);
//     let square = board.value[row][col];
//     if (!myTurn.value
//         || myColor.value != square.color) {
//       return;
//     }
//     moving.value = row == currentRow.value && col == currentCol.value ? false : true;
//     if (moving.value) {
//       currentRow.value = row(row);
//       currentCol.value = col(col);
//       clearMoves();
//       checkMoves(board.value[row][col], row, col, 1);
//       for (let i = 0; i < boardHeight.value; i++) {
//         for (let j = 0; j < boardWidth.value; j++) {
//           if (moves.value[i][j] > 0) {
//             return;
//           }
//         }
//       }
//     } else {
//       currentCol.value = -1;
//       currentRow.value = -1;
//       clearMoves();
//     }
//   },
//   moveTo: function (row, col) {
//     let i = row(row);
//     let j = col(col);
//     if (moves.value[i][j] == 0 || moves.value[i][j] > 2) {
//       return;
//     }
//     board.value[i][j] = board.value[currentRow.value][currentCol.value];
//     board.value[currentRow.value][currentCol.value] = {color: BLANK, piece: EMPTY};
//     if (moves.value[i][j] > 1) {
//       board.value[Math.floor((i+currentRow.value)/2)][Math.floor((j+currentCol.value)/2)] = {
//         color: BLANK,
//         piece: EMPTY,
//       };
//       currentRow.value = i;
//       currentCol.value = j;
//       clearMoves();
//       if (!checkMoves(board.value[i][j], i, j, 2, true)) {
//         moving.value = false;
//         myTurn.value = false;
//         currentRow.value = currentCol.value = -1;
//       }
//     } else {
//       moving.value = false;
//       myTurn.value = false;
//       currentRow.value = currentCol.value = -1;
//     }
//     boardUpdated.value = true;
//     nextTurn();
//   },
//   nextTurn() {
//     if (!myTurn.value && singlePlayer.value) {
//       currentAction.value = "Thinking...";
//       calculateMove();
//     } else if (myTurn.value) {
//       currentAction.value = "Your turn...";
//     }
//   },
//   checkMoves(square, row, col, level, jumpOnly = false) {
//     // console.log(row,col);
//     let hasMoves = false;
//     let direction = square.color == WHITE ? 1 : -1;
//     if (!jumpOnly) {
//       for (let i = 0; i < 4; i++) {
//         let x = (i % 2) == 0 ? col + 1 : col - 1;
//         if (x < 0 || x >= boardWidth.value) {
//           continue;
//         }
//         let y = (i > 1) ? row + 1 : row - 1;
//         if (y < 0 || y >= boardHeight.value) {
//           continue;
//         }
//         if (((y > row && direction == -1) || (y < row && direction == 1)) && square.piece == PAWN) {
//           continue;
//         }
//         // console.log("testing",y,x);
//         if (board.value[y][x].color == BLANK || board.value[y][x].piece == EMPTY) {
//           moves.value[y][x] = level;
//           hasMoves = true;
//         }
//         // console.log(y, x, board.value[y][x], moves.value[y][x]);
//       }
//       level++;
//     }
//     for (let i = 0; i < 4; i++) {
//       let x = (i % 2) == 0 ? col + 1 : col - 1;
//       if (x < 0 || x >= boardWidth.value) {
//         continue;
//       }
//       let y = (i > 1) ? row + 1 : row - 1;
//       if (y < 0 || y >= boardHeight.value) {
//         continue;
//       }
//       if (board.value[y][x].color == square.color || board.value[y][x].piece == EMPTY) {
//         continue;
//       }
//       x = (i % 2) == 0 ? col + 2 : col - 2;
//       if (x < 0 || x >= boardWidth.value) {
//         continue;
//       }
//       y = (i > 1) ? row + 2 : row - 2;
//       if (y < 0 || y >= boardHeight.value) {
//         continue;
//       }
//       if (board.value[y][x].piece == EMPTY && moves.value[y][x] == 0) {
//         moves.value[y][x] = level;
//         hasMoves |= checkMoves(square, y, x, level + 1, true);
//       }
//       // console.log(y, x, board.value[y][x], moves.value[y][x]);
//     }
//     return hasMoves;
//   },
//   calculateMove: function () {
//     console.log('calculating...');
//     moving.value = true;
//     let bestMove = [];
//     let moveTo = [];
//     let moveSize = 0;
//     for (let i = 0; i < boardHeight.value; i++) {
//       for (let j = 0; j < boardHeight.value; j++) {
//         if ((myColor.value == WHITE && board.value[i][j].color == BLACK)
//             || (myColor.value == BLACK && board.value[i][j].color == WHITE)) {
//           console.log("checking ", i, j);
//           currentRow.value = i;
//           currentCol.value = j;
//           clearMoves();
//           checkMoves(board.value[i][j], i, j, 1);
//           let size = 0;
//           let lastMove = [];
//           for (let k = 0; k < boardHeight.value; k++) {
//             for (let l = 0; l < boardHeight.value; l++) {
//               size += moves.value[k][l] ? 1 : 0;
//               if (moves.value[k][l]) {
//                 lastMove.push([k, l]);
//               }
//             }
//           }
//           if (size > moveSize) {
//             moveSize = size;
//             bestMove = [i, j];
//             let rnd = Math.floor(Math.random() * size);
//             moveTo = lastMove[rnd];
//           }
//         }
//       }
//       // let move = math.ceil(math.random() * moveable.length);
//     }
//     console.log("my play:", bestMove[0], bestMove[1], moveTo[0], moveTo[1])
//     board.value[moveTo[0]][moveTo[1]] = board.value[bestMove[0]][bestMove[1]];
//     board.value[bestMove[0]][bestMove[1]] = {color: BLANK, piece: EMPTY};
//     myTurn.value = true;
//     boardUpdated.value = true;
//     moving.value = false;
//     currentRow.value = currentCol.value = -1;
//     boardUpdated.value = true;
//     nextTurn();
//   },
// },
// refreshBoard: function () {
//   var moveIndex = 0;
//   let vm = this;
//   for (let i = 0; i < boardHeight.value; i++) {
//     for (let j = 0; j < boardWidth.value; j++) {
//       $('#' + id(i, j)).each(function () {
//         $(this).removeClass(['dark', 'light']).addClass((i + j) % 2 == 1 ? 'light' : 'dark');
//         let pos = this;
//         if (vm.board[i][j] == vm.EMPTY) {
//           $(this).attr('tag','');
//         } else if (vm.board[i][j] == vm.WHITE_PAWN) {
//           $('#' + vm.id(i, j) + ' img').each(function () {
//             let canBeDragged = false;
//             if (i > 0 && j > 0) {
//               canBeDragged |= vm.checkPosition(i + 1, j - 1, moveIndex, vm.SIMPLE_MOVE);
//             }
//             if (i > 0 && j < vm.boardWidth - 1) {
//               canBeDragged |= vm.checkPosition(i + 1, j + 1, moveIndex, vm.SIMPLE_MOVE);
//             }
//             if (canBeDragged) {
//               console.log(moveIndex);
//               let el = $(this).draggable({
//                 containment: '.checkers',
//                 scope: 'target_' + moveIndex,
//                 snap: true,
//                 snapMode: "inner",
//                 snapTolerance: 10,
//                 revert: 'invalid',
//                 start: function () {
//                   let tag = $(this).attr('tag')
//                   console.log('start ' + tag);
//                   $('.target_' + tag).each(function () {
//                     // console.log(this);
//                     $(this).droppable({
//                       accept: ".target_" + tag,
//                       scope: "target_" + tag,
//                       drop: function(event, ui) {
//                         let coords = $(ui.draggable).attr('id').split(/[rc]/);
//                         let row = parseInt(pos[1]);
//                         let col = parseInt(pos[2]);
//                         vm.board[row][col] = vm.WHITE_PAWN;
//                         // console.log(row,col);
//                         vm.boardUpdated = true;
//                         vm.boardReady = true;
//                       },
//                       classes: {
//                         "ui-droppable-active": "targetable",
//                         "ui-droppable-hover": "overtarget"
//                       },
//                     });
//                   });
//                 },
//                 stop: function () {
//                   console.log('stop');
//                   let tag = $(this).attr('tag')
//                   let pos = $(this).parent().attr('tag');
//                   let coords = pos.split(",");
//                   let row = parseInt(coords[0]);
//                   let col = parseInt(coords[1]);
//                   console.log(vm.board);
//                   vm.board[row][col] = vm.EMPTY;
//                   // vm.boardUpdated = true;
//                   $('.target_' + tag).each(function () {
//                     $(this).droppable("destroy");
//                   });
//                   vm.boardReady = false;
//                 },
//               });
//               el = $(el).addClass('target_' + moveIndex);
//               $(el).attr('tag', moveIndex);
//               moveIndex++;
//             }
//           });
//         }
//       });
//     }
//   }
//   boardReady.value = true;
// },

const checkPosition = (i, j) => {
  if (moveType == SIMPLE_MOVE && board.value[i][j] == EMPTY) {
    $('#' + id(i, j)).each(function () {
      $(this).addClass('target_' + moveIndex);
    });
    return true;
  }
  let myPawn = myColor.value == WHITE.value ? WHITE.value_PAWN : BLACK.value_PAWN;
  let myKing = myColor.value == WHITE.value ? WHITE.value_KING : BLACK.value_KING;
  let otherPawn = myColor.value == BLACK.value ? WHITE.value_PAWN : BLACK.value_PAWN;
  let otherKing = myColor.value == BLACK.value ? WHITE.value_KING : BLACK.value_KING;
  return false;
  // if (moveType == JUMP_MOVE && (board.value[i, j] == myPawn || board.value[i, j] == myKing)) {
  //
  // }
  // },
}

</script>

<template>
  <div class="flex flex-col lg:flex-row align-items-start justify-content-center" >
    <div class="w-3/4 mx-auto lg:w-1/6 form align-items-start  action-box">{{ currentAction }}</div>
    <div class=" w-full lg:w-1/2 form align-items-start action-box">
      <table class="chess max-w-lg mx-auto justify-content-center align-items-center">
        <thead>
        <td class="v-board-side h-board-side"></td>
        <td v-for="col in boardWidth" class="h-board-side upside-down-text">
          {{ myColor == BLACK ? 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.substring(boardWidth - col, boardWidth - col + 1) : "" }}
        </td>
        <td class="v-board-side h-board-side"></td>
        </thead>
        <tbody>
        <tr v-for="row in boardHeight">
          <td class="v-board-side">{{ myColor == WHITE ? boardWidth - row + 1 : "" }}</td>
          <td v-for="col in boardWidth"
              :id="id(boardHeight-row,col-1)" :class="square(  boardHeight-row,col-1)"
              v-on:click="moveTo(boardHeight-row,col-1)">
            <!--            <div v-if="boardReady">{{ board[boardHeight-row][col-1] }}</div>-->
            <img :class="{'upside-down':row<4}" v-if="boardReady && piece(boardHeight-row,col-1) != ''"
                 :src="piece(boardHeight-row,col-1)" v-on:click="move(boardHeight-row,col-1)"/>
            <!--            <img :class="{movable:row==6}" v-if="row>5 && (col+row)%2 == 1" src="images/white-pawn.png"/>-->
            <div v-if="boardReady && moves[boardHeight-row][col-1]>0" class="moves">{{
                moves[boardHeight - row][col - 1]
              }}
            </div>
          </td>
          <td class="v-board-side upside-down-text">{{ myColor == BLACK ? row : "" }}</td>
        </tr>
        </tbody>
        <thead>
        <td class="v-board-side h-board-side"></td>
        <td v-for="col in boardWidth" class="h-board-side">
          {{ myColor == WHITE ? "ABCDEFGHIJKLMNOPQRSTUVWXYZ".substr(col - 1, 1) : "" }}
        </td>
        <td class="v-board-side h-board-side"></td>
        </thead>
      </table>
    </div>
    <div class="max-md:w-1/2 md:w-1/3 mx-auto p-4 lg:ml-4 lg:border-l border-dotted border-black max-lg:text-center">
      <input v-model="myColor" :value="WHITE" type="radio"> White<br/>
      <input v-model="myColor" :value="BLACK" type="radio"> Black
    </div>
  </div>
</template>

<style scoped>


.action-box {
  @apply text-center mx-auto;

  select {
    width: 100%;

    &.games {
      size: 5px;
    }
  }
}

.move-list {
  /* width: 15 %;*/
  margin: 20px;
  padding: 5px;
  border: 1px solid lightgrey;
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
  @apply w-full mt-4 bg-[url('@/assets/images/leather.png')] bg-repeat;

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
    vertical-align: middle;
    text-align: center;

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