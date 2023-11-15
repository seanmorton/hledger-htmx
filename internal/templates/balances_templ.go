// Code generated by templ - DO NOT EDIT.

// templ: version: 0.2.465
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/seanmorton/hledger-htmx/internal/hledger"
)

func graph(entries []hledger.BalanceEntry) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_graph_beca`,
		Function: `function __templ_graph_beca(entries){const canvas = document.getElementById("balanceChart");
  const labels = entries.map((e) => e.account);
  const values = entries.map((e) => e.amount);
  const colors = [
    "#b91d47",
    "#00aba9",
    "#2b5797",
    "#e8c3b9",
    "#1e7145"
  ];
  const chart = new Chart(canvas, {
    type: "pie",
    data: {
      labels: labels,
      datasets: [{
        backgroundColor: colors,
        data: values
      }]
    },
    options: {
      title: {
        display: true,
        text: "Balances"
      }
    }
  });
  canvas.onclick = function(evt) {
    var activePoints = chart.getElementsAtEvent(evt);
    if (activePoints[0]) {
      var chartData = activePoints[0]['_chart'].config.data;
      var idx = activePoints[0]['_index'];

      var account = chartData.labels[idx];
      var accountSel = "#" + account.replaceAll(":", "")
      //var value = chartData.datasets[0].data[idx];
      //var url = "/register?account=" + account
      htmx.trigger(accountSel, "click")
    }
  };}`,
		Call:       templ.SafeScript(`__templ_graph_beca`, entries),
		CallInline: templ.SafeScriptInline(`__templ_graph_beca`, entries),
	}
}

func Balances(entries []hledger.BalanceEntry) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<canvas id=\"balanceChart\" style=\"width:100%;max-width:800px\"></canvas>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = graph(entries).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
