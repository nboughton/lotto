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
    scatter3D: {
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
        drawResultsAverage()
        break
      case "graph-freqdist-bar":
        drawResultsGraph("freqdist/bar", layouts.bar)
        break
      case "graph-freqdist-scatter":
        drawResultsGraph("freqdist/scatter", layouts.scatter)
        break
      case "graph-timeseries-scatter":
        drawResultsGraph("timeseries/scatter")
        break
      case "graph-3d-scatter":
        drawResultsGraph("3d/scatter", layouts.scatter3D)
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

  function drawResultsAverage() {
    $.getJSON("/api/results/average", params(), function (d) {
      $(sel.results).empty().append('<h1 id="' + sel.resultsAvg.replace('#', '') + '" class="centered"></h1>')
      for (i = 0; i < d.length; i++) {
        $(sel.resultsAvg).append("<span class='num'>" + d[i] + "</span>")
      }
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
  drawResultsGraph("freqdist/bar", layouts.bar)
})