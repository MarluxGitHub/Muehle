package services

import (
	"marluxgithub/muehle/pkg/muehle/domain/entities"
	"testing"
)

func TestIsFieldPartOfClosedMill(t *testing.T) {
	bs := NewBoardService()
	b := bs.CreateBoard()
	b.Fields[0].Color = entities.ColorBlack
	b.Fields[1].Color = entities.ColorBlack
	b.Fields[2].Color = entities.ColorBlack
	if !bs.IsFieldPartOfClosedMill(b, 0, entities.ColorBlack) || !bs.IsFieldPartOfClosedMill(b, 1, entities.ColorBlack) {
		t.Fatal("stones on 0,1,2 should be in mill")
	}
	b.Fields[12].Color = entities.ColorBlack
	if bs.IsFieldPartOfClosedMill(b, 12, entities.ColorBlack) {
		t.Fatal("single black on 12 is not a full mill line alone")
	}
	if bs.IsFieldPartOfClosedMill(b, 0, entities.ColorWhite) {
		t.Fatal("wrong color")
	}
}

func TestEnemyHasStoneOutsideMill(t *testing.T) {
	bs := NewBoardService()
	b := bs.CreateBoard()
	for i := range b.Fields {
		b.Fields[i].Color = entities.ColorUnknown
	}
	b.Fields[9].Color = entities.ColorBlack
	b.Fields[10].Color = entities.ColorBlack
	b.Fields[11].Color = entities.ColorBlack
	b.Fields[12].Color = entities.ColorBlack
	if !bs.EnemyHasStoneOutsideMill(b, entities.ColorBlack) {
		t.Fatal("12 is outside mill 9-10-11")
	}
	for i := range b.Fields {
		b.Fields[i].Color = entities.ColorUnknown
	}
	b.Fields[9].Color = entities.ColorBlack
	b.Fields[10].Color = entities.ColorBlack
	b.Fields[11].Color = entities.ColorBlack
	if bs.EnemyHasStoneOutsideMill(b, entities.ColorBlack) {
		t.Fatal("all three only in one mill, no outside stone")
	}
}
