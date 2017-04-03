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
    resultsAvg: "#results-average",
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
            title: "Frequency"
          },
          zaxis: {
            title: "Set"
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
      case "num-average-mean":
        drawResultsAverage("")
        break
      case "num-average-ranges":
        drawResultsAverage("/ranges")
        break
      case "graph-freqdist-bar":
        drawResultsGraph("freqdist/bar", layouts.freqdist.bar)
        break
      case "graph-freqdist-scatter":
        drawResultsGraph("freqdist/scatter", layouts.freqdist.scatter)
        break
      case "graph-freqdist-ms-bubble":
        drawResultsGraph("freqdist-ms/bubble", layouts.freqdist.ms)
        break
      case "graph-freqdist-ms-scatter3D":
        drawResultsGraph("freqdist-ms/scatter3D", layouts.freqdist.ms)
        break
      case "graph-timeseries-scatter":
        drawResultsGraph("timeseries/scatter", layouts.timeseries.line)
        break
      case "graph-timeseries-line":
        drawResultsGraph("timeseries/line", layouts.timeseries.line)
        break
      case "graph-3d-scatter":
        drawResultsGraph("3d/scatter", layouts.scatter3D.results)
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

  function drawResultsAverage(type) {
    $.getJSON("/api/results/average" + type, params(), function (d) {
      $(sel.results).empty().append('<h1 id="' + sel.resultsAvg.replace('#', '') + '" class="centered"></h1>')
      $.each(d, function (i, n) {
        $(sel.resultsAvg).append("<span class='num'>" + n + "</span>")
      })
    })
  }

  function drawResultsGraph(type, layout) {
    $.getJSON("/api/results/graph/" + type, params(), function (d) {
      $(sel.results).empty().append('<div id="' + sel.resultsGraph.replace("#", "") + '"></div>')
      Plotly.newPlot(sel.resultsGraph.replace("#", ""), d, layout)
    })
    //getMachinesSetsCombos()
  }

  /*
  // currently unused
  function getMachinesSetsCombos() {
    $.getJSON("/api/machines/sets/combos", params(), function (d) {
      console.log(d)
    })
  }
  */

  //******************************************************* WHAT DO */
  // Draw a frequency distribution
  drawResultsGraph("freqdist/bar", layouts.freqdist.bar)
})