$(function () {
  //******************************************************* GLOBALS */
  // Define a map for key ids/classes in case I decide to change them later
  var sel = {
    machineSelect: "#filter-machine",
    setSelect: "#filter-set",
    dateClass: ".form-date",
    dateStart: "#date-start",
    dateEnd: "#date-end",
    queryType: "#query-type",
    querySubmit: "#query-submit",
    results: "#results",
    resultsNum: "#results-numbers",
    resultsGraph: "#results-graph"
  }

  var layouts = {
    freqdist: {
      bar: {
        barmode: "stack",
        xaxis: {
          title: "Result"
        },
        yaxis: {
          title: "Frequency"
        }
      },
      scatter: {
        xaxis: {
          title: "Result"
        },
        yaxis: {
          title: "Frequency"
        }
      },
      ms: {
        xaxis: {
          title: "Machine"
        },
        yaxis: {
          title: "Set"
        },
        hovermode: "closest"
      }
    },
    timeseries: {
      line: {
        xaxis: {
          title: "Date:Set:Machine",
        },
        yaxis: {
          title: "Result"
        }
      }
    },
    scatter3D: {
      results: {
        scene: {
          xaxis: {
            title: "Machine"
          },
          yaxis: {
            title: "Set"
          },
          zaxis: {
            title: "Result"
          }
        }
      },
      frequency: {
        scene: {
          xaxis: {
            title: "Machine"
          },
          yaxis: {
            title: "Set"
          },
          zaxis: {
            title: "Frequency"
          }
        }
      }
    }
  }
  //******************************************************* HANDLERS */
  /// Pickadate.js
  $.getJSON("/api/range", function (d) {
    var pickerOpts = {
      "format": "yyyy-mm-dd",
      "selectYears": true,
      "selectMonths": true,
      "min": new Date(d.first * 1000),
      "max": new Date(d.last * 1000)
    }

    $(sel.dateClass).pickadate(pickerOpts)
  })

  /// Redraw options list for Sets depending on Machine selection
  $(sel.machineSelect).change(function (e) {
    redrawSetsList()
  })

  /// Redraw machine and set lists to reflect available data for date range
  $(sel.dateClass).change(function (e) {
    redrawMachinesList()
    redrawSetsList()
  })


  /// Query Exec
  $(sel.querySubmit).click(function (e) {
    switch ($(sel.queryType).val()) {
      case "num-results-mean":
        drawResultsNumbers("average")
        break
      case "num-results-ranges":
        drawResultsNumbers("ranges")
        break
      case "num-results-frequent":
        drawResultsNumbers("frequent", "Where there is a tie a result is picked pseudo-randomly.")
        break
      case "graph-results-freqdist-bar":
        drawResultsGraph("results/freqdist/bar", layouts.freqdist.bar)
        break
      case "graph-results-freqdist-scatter":
        drawResultsGraph("results/freqdist/scatter", layouts.freqdist.scatter)
        break
      case "graph-results-timeseries-scatter":
        drawResultsGraph("results/timeseries/scatter", layouts.timeseries.line)
        break
      case "graph-results-timeseries-line":
        drawResultsGraph("results/timeseries/line", layouts.timeseries.line)
        break
      case "graph-results-raw-scatter3d":
        drawResultsGraph("results/raw/scatter3d", layouts.scatter3D.results)
        break
      case "graph-ms-freqdist-bubble":
        drawResultsGraph("ms/freqdist/bubble", layouts.freqdist.ms)
        break
      case "graph-ms-freqdist-scatter3d":
        drawResultsGraph("ms/freqdist/scatter3d", layouts.scatter3D.frequency)
        break
      default:
        break
    }
  })

  //******************************************************* FUNCTIONS */
  /// Return query params 
  function params() {
    return {
      start: $(sel.dateStart).val(),
      end: $(sel.dateEnd).val(),
      set: $(sel.setSelect).val(),
      machine: $(sel.machineSelect).val()
    }
  }

  function redrawSetsList() {
    $.getJSON("/api/sets", params(), function (d) {
      var el = $(sel.setSelect)
      el.empty().append('<option value="0">All</option>')
      $.each(d, function (i, n) {
        el.append('<option value="' + n + '">' + n + '</option>')
      })
    })
  }

  function redrawMachinesList() {
    $.getJSON("/api/machines", params(), function (d) {
      var el = $(sel.machineSelect)
      el.empty().append('<option value="all">All</option>')
      $.each(d, function (i, m) {
        el.append('<option value="' + m + '">' + m + '</option>')
      })
    })
  }

  function drawResultsNumbers(type, notes) {
    $.getJSON("/api/numbers/" + type, params(), function (d) {
      $(sel.results).empty().append('<h1 id="' + sel.resultsNum.replace('#', '') + '" class="centered"></h1>')
      $.each(d, function (i, n) {
        $(sel.resultsNum).append("<span class='num'>" + n + "</span>")
      })
      if (notes) {
        $(sel.results).append("<p>" + notes + "</p>")
      }
    })
  }

  function drawResultsGraph(type, layout) {
    $.getJSON("/api/graph/" + type, params(), function (d) {
      $(sel.results).empty().append('<div id="' + sel.resultsGraph.replace("#", "") + '"></div>')
      Plotly.newPlot(sel.resultsGraph.replace("#", ""), d, layout)
    })
  }

  //******************************************************* WHAT DO */
  // Draw a frequency distribution
  drawResultsGraph("results/freqdist/bar", layouts.freqdist.bar)
})