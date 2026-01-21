package services

import (
	"marluxgithub/muehle/pkg/muehle/domain/entities"
	"slices"
)

type BoardService struct{}

func (boardService *BoardService) CreateBoard() entities.Board {
	fields := make([]entities.Field, 24)

	for i := range 24 {
		fields[i] = entities.Field{
			Index:     i,
			Color:     entities.ColorUnknown,
			Neighbors: make([]entities.Field, 0),
		}
	}

	board := entities.Board{
		Fields: fields,
	}

	boardService.CreateNeighbors(board)

	return board
}

func (boardService *BoardService) CreateNeighbors(board entities.Board) entities.Board {
	board.Fields[0].Neighbors = []entities.Field{board.Fields[1], board.Fields[9]}
	board.Fields[1].Neighbors = []entities.Field{board.Fields[0], board.Fields[2], board.Fields[4]}
	board.Fields[2].Neighbors = []entities.Field{board.Fields[1], board.Fields[14]}
	board.Fields[3].Neighbors = []entities.Field{board.Fields[4], board.Fields[10]}
	board.Fields[4].Neighbors = []entities.Field{board.Fields[1], board.Fields[3], board.Fields[5], board.Fields[7]}
	board.Fields[5].Neighbors = []entities.Field{board.Fields[4], board.Fields[13]}
	board.Fields[6].Neighbors = []entities.Field{board.Fields[7], board.Fields[11]}
	board.Fields[7].Neighbors = []entities.Field{board.Fields[4], board.Fields[6], board.Fields[8]}
	board.Fields[8].Neighbors = []entities.Field{board.Fields[7], board.Fields[12]}
	board.Fields[9].Neighbors = []entities.Field{board.Fields[0], board.Fields[10], board.Fields[21]}
	board.Fields[10].Neighbors = []entities.Field{board.Fields[3], board.Fields[9], board.Fields[11], board.Fields[18]}
	board.Fields[11].Neighbors = []entities.Field{board.Fields[6], board.Fields[10], board.Fields[15]}
	board.Fields[12].Neighbors = []entities.Field{board.Fields[8], board.Fields[13], board.Fields[17]}
	board.Fields[13].Neighbors = []entities.Field{board.Fields[5], board.Fields[12], board.Fields[14], board.Fields[20]}
	board.Fields[14].Neighbors = []entities.Field{board.Fields[2], board.Fields[13], board.Fields[23]}
	board.Fields[15].Neighbors = []entities.Field{board.Fields[11], board.Fields[16]}
	board.Fields[16].Neighbors = []entities.Field{board.Fields[15], board.Fields[17], board.Fields[19]}
	board.Fields[17].Neighbors = []entities.Field{board.Fields[12], board.Fields[16]}
	board.Fields[18].Neighbors = []entities.Field{board.Fields[10], board.Fields[19]}
	board.Fields[19].Neighbors = []entities.Field{board.Fields[16], board.Fields[18], board.Fields[20], board.Fields[22]}
	board.Fields[20].Neighbors = []entities.Field{board.Fields[13], board.Fields[19]}
	board.Fields[21].Neighbors = []entities.Field{board.Fields[9], board.Fields[22]}
	board.Fields[22].Neighbors = []entities.Field{board.Fields[19], board.Fields[21], board.Fields[23]}
	board.Fields[23].Neighbors = []entities.Field{board.Fields[14], board.Fields[22]}

	return board
}

func (boardService *BoardService) HasPlayerThreeStones(board entities.Board, fieldIndex int, playerColor entities.Color) bool {
	millCombination := [][]int{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
		{9, 10, 11},
		{12, 13, 14},
		{15, 16, 17},
		{18, 19, 20},
		{21, 22, 23},

		{0, 9, 21},
		{3, 10, 18},
		{6, 11, 15},
		{1, 4, 7},
		{16, 19, 22},
		{8, 12, 17},
		{5, 13, 20},
		{2, 14, 23},
	}

	for _, combination := range millCombination {
		// fieldIndex is in combination
		if slices.Contains(combination, fieldIndex) {
			// check if the combination is a mill
			if board.Fields[combination[0]].Color == playerColor && board.Fields[combination[1]].Color == playerColor && board.Fields[combination[2]].Color == playerColor {
				return true
			}
		}
	}

	return false
}

func NewBoardService() *BoardService {
	boardService := &BoardService{}

	return boardService
}
