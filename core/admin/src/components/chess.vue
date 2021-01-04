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
    <div class="row align-items-center justify-content-center">
      <div class="col">{{ currentAction }}</div>
      <table class="chess justify-content-center align-items-center">
        <thead>
        <td class="v-board-side h-board-side"></td>
        <td v-for="col in boardWidth" class="h-board-side upside-down-text">
          {{ myColor == BLACK ? 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.substr(boardWidth - col, 1) : "" }}
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
            <img v-if="boardReady && piece(boardHeight-row,col-1) != ''"
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
      <div class="col">
        <input v-model="myColor" :value="WHITE" type="radio"> White<br/>
        <input v-model="myColor" :value="BLACK" type="radio"> Black
      </div>
    </div>
  </div>
</template>

<script>
const BLANK = 0, BLACK = 1, WHITE = 2;
const EMPTY = 0, PAWN = 1, ROOK = 2, KNIGHT = 3, BISHOP = 4, QUEEN = 5, KING = 6;
const initialPieces = [2, 3, 4, 5, 6, 4, 3, 2];

export default {
  // name: "messaging",
  props: [],
  data() {
    return {
      BLANK: BLANK,
      BLACK: BLACK,
      WHITE: WHITE,
      SIMPLE_MOVE: 0,
      JUMP_MOVE: 1,
      boardWidth: 8,
      boardHeight: 8,
      pawnRows: 4,
      board: [],
      moves: [],
      boardUpdated: false,
      boardReady: false,
      myColor: WHITE,
      myTurn: this.myColor == WHITE,
      moving: false,
      currentRow: -1,
      currentCol: -1,
      singlePlayer: true,
      currentAction: "Starting...",
    }
  },
  watch: {
    boardReady: function (value) {
      if (value) {
        this.boardUpdated = false;
        // this.refreshBoard();
      }
    },
  },
  mounted() {
    this.setupBoard();
  },
  updated() {
    this.$nextTick(function () {
      if (this.boardUpdated) {
        console.log("refreshed");
        // this.refreshBoard();
        this.boardUpdated = false;
      }
    });
  },
  methods: {
    setupBoard: function () {
      for (let i = 0; i < this.boardHeight; i++) {
        this.board[i] = [];
        this.moves[i] = [];
        for (let j = 0; j < this.boardWidth; j++) {
          this.board[i][j] = {
            color: BLANK,
            piece: EMPTY,
          }
          this.moves[i][j] = false;
          let color = WHITE;
          if (i > this.boardHeight / 2) {
            color = BLACK;
          }
          if (i == 1 || i == this.boardHeight - 2) {
            this.board[i][j] = {
              color: color,
              piece: PAWN,
            };
          }
          if (i == 0 || i == this.boardHeight - 1) {
            this.board[i][j] = {
              color: color,
              piece: initialPieces[j],
            };
          }
        }
      }
      this.boardReady = true;
      this.nextTurn();
      // console.log(this.board);
      // this.boardUpdated = true;
      // this.refreshBoard();
    },
    row: function (r) {
      return this.myColor == WHITE ? r : this.boardHeight - r - 1;
    },
    col: function (c) {
      return this.myColor == WHITE ? c : this.boardWidth - c - 1;
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
      return {
        "dark": (i + j) % 2 == 0,
        "light": (i + j) % 2 == 1,
        "piece-to-move": i == this.currentRow && j == this.currentCol,
        "targetable": this.moves[i][j] > 0,
      };
    },
    piece: function (r, c) {
      let i = this.row(r);
      let j = this.col(c);
      let name = "/images/";
      switch (this.board[i][j].color) {
        case BLACK:
          name += "b-";
          break;
        case WHITE:
          name += "w-";
          break;
        default:
          return "";
      }
      switch (this.board[i][j].piece) {
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
          return "";
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
    //   move: function (row, col) {
    //     console.log(row, col);
    //     let square = this.board[row][col];
    //     if (!this.myTurn
    //         || this.myColor != square.color) {
    //       return;
    //     }
    //     this.moving = row == this.currentRow && col == this.currentCol ? false : true;
    //     if (this.moving) {
    //       this.currentRow = this.row(row);
    //       this.currentCol = this.col(col);
    //       this.clearMoves();
    //       this.checkMoves(this.board[row][col], row, col, 1);
    //       for (let i = 0; i < this.boardHeight; i++) {
    //         for (let j = 0; j < this.boardWidth; j++) {
    //           if (this.moves[i][j] > 0) {
    //             return;
    //           }
    //         }
    //       }
    //     } else {
    //       this.currentCol = -1;
    //       this.currentRow = -1;
    //       this.clearMoves();
    //     }
    //   },
    //   moveTo: function (row, col) {
    //     let i = this.row(row);
    //     let j = this.col(col);
    //     if (this.moves[i][j] == 0 || this.moves[i][j] > 2) {
    //       return;
    //     }
    //     this.board[i][j] = this.board[this.currentRow][this.currentCol];
    //     this.board[this.currentRow][this.currentCol] = {color: BLANK, piece: EMPTY};
    //     if (this.moves[i][j] > 1) {
    //       this.board[Math.floor((i+this.currentRow)/2)][Math.floor((j+this.currentCol)/2)] = {
    //         color: BLANK,
    //         piece: EMPTY,
    //       };
    //       this.currentRow = i;
    //       this.currentCol = j;
    //       this.clearMoves();
    //       if (!this.checkMoves(this.board[i][j], i, j, 2, true)) {
    //         this.moving = false;
    //         this.myTurn = false;
    //         this.currentRow = this.currentCol = -1;
    //       }
    //     } else {
    //       this.moving = false;
    //       this.myTurn = false;
    //       this.currentRow = this.currentCol = -1;
    //     }
    //     this.boardUpdated = true;
    //     this.nextTurn();
    //   },
    //   nextTurn() {
    //     if (!this.myTurn && this.singlePlayer) {
    //       this.currentAction = "Thinking...";
    //       this.calculateMove();
    //     } else if (this.myTurn) {
    //       this.currentAction = "Your turn...";
    //     }
    //   },
    //   checkMoves(square, row, col, level, jumpOnly = false) {
    //     // console.log(row,col);
    //     let hasMoves = false;
    //     let direction = square.color == WHITE ? 1 : -1;
    //     if (!jumpOnly) {
    //       for (let i = 0; i < 4; i++) {
    //         let x = (i % 2) == 0 ? col + 1 : col - 1;
    //         if (x < 0 || x >= this.boardWidth) {
    //           continue;
    //         }
    //         let y = (i > 1) ? row + 1 : row - 1;
    //         if (y < 0 || y >= this.boardHeight) {
    //           continue;
    //         }
    //         if (((y > row && direction == -1) || (y < row && direction == 1)) && square.piece == PAWN) {
    //           continue;
    //         }
    //         // console.log("testing",y,x);
    //         if (this.board[y][x].color == BLANK || this.board[y][x].piece == EMPTY) {
    //           this.moves[y][x] = level;
    //           hasMoves = true;
    //         }
    //         // console.log(y, x, this.board[y][x], this.moves[y][x]);
    //       }
    //       level++;
    //     }
    //     for (let i = 0; i < 4; i++) {
    //       let x = (i % 2) == 0 ? col + 1 : col - 1;
    //       if (x < 0 || x >= this.boardWidth) {
    //         continue;
    //       }
    //       let y = (i > 1) ? row + 1 : row - 1;
    //       if (y < 0 || y >= this.boardHeight) {
    //         continue;
    //       }
    //       if (this.board[y][x].color == square.color || this.board[y][x].piece == EMPTY) {
    //         continue;
    //       }
    //       x = (i % 2) == 0 ? col + 2 : col - 2;
    //       if (x < 0 || x >= this.boardWidth) {
    //         continue;
    //       }
    //       y = (i > 1) ? row + 2 : row - 2;
    //       if (y < 0 || y >= this.boardHeight) {
    //         continue;
    //       }
    //       if (this.board[y][x].piece == EMPTY && this.moves[y][x] == 0) {
    //         this.moves[y][x] = level;
    //         hasMoves |= this.checkMoves(square, y, x, level + 1, true);
    //       }
    //       // console.log(y, x, this.board[y][x], this.moves[y][x]);
    //     }
    //     return hasMoves;
    //   },
    //   calculateMove: function () {
    //     console.log('calculating...');
    //     this.moving = true;
    //     let bestMove = [];
    //     let moveTo = [];
    //     let moveSize = 0;
    //     for (let i = 0; i < this.boardHeight; i++) {
    //       for (let j = 0; j < this.boardHeight; j++) {
    //         if ((this.myColor == WHITE && this.board[i][j].color == BLACK)
    //             || (this.myColor == BLACK && this.board[i][j].color == WHITE)) {
    //           console.log("checking ", i, j);
    //           this.currentRow = i;
    //           this.currentCol = j;
    //           this.clearMoves();
    //           this.checkMoves(this.board[i][j], i, j, 1);
    //           let size = 0;
    //           let lastMove = [];
    //           for (let k = 0; k < this.boardHeight; k++) {
    //             for (let l = 0; l < this.boardHeight; l++) {
    //               size += this.moves[k][l] ? 1 : 0;
    //               if (this.moves[k][l]) {
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
    //     this.board[moveTo[0]][moveTo[1]] = this.board[bestMove[0]][bestMove[1]];
    //     this.board[bestMove[0]][bestMove[1]] = {color: BLANK, piece: EMPTY};
    //     this.myTurn = true;
    //     this.boardUpdated = true;
    //     this.moving = false;
    //     this.currentRow = this.currentCol = -1;
    //     this.boardUpdated = true;
    //     this.nextTurn();
    //   },
    // },
    // refreshBoard: function () {
    //   var moveIndex = 0;
    //   let vm = this;
    //   for (let i = 0; i < this.boardHeight; i++) {
    //     for (let j = 0; j < this.boardWidth; j++) {
    //       $('#' + this.id(i, j)).each(function () {
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
    //   this.boardReady = true;
    // },
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
      // if (moveType == JUMP_MOVE && (this.board[i, j] == myPawn || this.board[i, j] == myKing)) {
      //
      // }
      // },
    },
  }
}
</script>
