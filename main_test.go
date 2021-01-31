package main

import (
	"github.com/penguin-statistics/widget-backend/config"
	"github.com/penguin-statistics/widget-backend/controller/matrix"
	"github.com/penguin-statistics/widget-backend/controller/meta"
	"github.com/penguin-statistics/widget-backend/response"
	"testing"
)

var Controllers *meta.Collection
var Render *response.Assembler

func init() {
	Controllers = meta.NewCollection()
	Render = response.New(Controllers, config.C.Static.Widget.Root)
}

func BenchmarkStage(b *testing.B) {
	query := &matrix.Query{Server: "CN", StageID: "main_01-07"}

	for i := 0; i < b.N; i++ {
		records, err := Controllers.Matrix.Query(query)
		if err != nil {
			b.Error(err)
		}

		Render.MarshalMatrix(records, query)
	}
}

func BenchmarkItem(b *testing.B) {
	query := &matrix.Query{Server: "CN", ItemID: "30012"}

	for i := 0; i < b.N; i++ {
		records, err := Controllers.Matrix.Query(query)
		if err != nil {
			b.Error(err)
		}

		Render.MarshalMatrix(records, query)
	}
}

func BenchmarkExact(b *testing.B) {
	query := &matrix.Query{Server: "CN", StageID: "main_01-07", ItemID: "30012"}

	for i := 0; i < b.N; i++ {
		records, err := Controllers.Matrix.Query(query)
		if err != nil {
			b.Error(err)
		}

		Render.MarshalMatrix(records, query)
	}
}
