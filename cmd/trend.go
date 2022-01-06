package cmd

import (
	"bytes"
	"time"

	"github.com/lucklove/naglfar/pkg/client"
	"github.com/lucklove/naglfar/pkg/render"
	"github.com/spf13/cobra"
	chart "github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func newTrendCommand() *cobra.Command {
	field := ""

	cmd := &cobra.Command{
		Use:   "trend",
		Short: "naglfar trend <fragment> [events]",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Help()
			}

			c := client.New()
			defer c.Close()

			n := time.Now()
			if field != "" && len(args) == 2 {
				trends, err := c.GetFieldTrend(cmd.Context(), args[0], n.Add(-time.Hour*24*30), n, args[1], field)
				if err != nil {
					return err
				}
				return renderTrends(trends)
			} else {
				trends, err := c.GetTrend(cmd.Context(), args[0], n.Add(-time.Hour*24*30), n, args[1:]...)
				if err != nil {
					return err
				}
				return renderTrends(trends)
			}
		},
	}

	cmd.Flags().StringVarP(&field, "field", "f", "", "Specify the field to group")

	return cmd
}

func renderTrends(trends []client.Trend) error {
	series := []chart.Series{}
	for _, trend := range trends {
		if len(trend.Points) < 2 {
			trend.Points = append(trend.Points, trend.Points[0])
		}
		xv := []time.Time{}
		yv := []float64{}
		for _, p := range trend.Points {
			xv = append(xv, time.Unix(p.Timestamp, 0))
			yv = append(yv, float64(p.Value))
		}
		series = append(series, chart.TimeSeries{
			XValues: xv,
			YValues: yv,
		})
	}
	graph := chart.Chart{
		Series: series,
		Background: chart.Style{
			FillColor: drawing.ColorBlack,
			Padding:   chart.Box{Top: 40},
		},
		Canvas: chart.Style{FillColor: drawing.ColorBlack},
		XAxis:  chart.XAxis{Style: chart.Style{Hidden: true}},
		YAxis:  chart.YAxis{Style: chart.Style{Hidden: true}},
	}
	buffer := bytes.NewBuffer(nil)
	if err := graph.Render(chart.PNG, buffer); err != nil {
		return err
	}
	return render.Render(buffer)
}
