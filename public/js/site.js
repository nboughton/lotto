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

  // Chart config
  var chartOptions = {
    tooltips: {
      mode: "label"
    }
  }

  var chartColours = {
    "Ball 1": { "fill": "rgba(255,0,132,0.05)", "stroke": "rgba(255,0,132,1)" },
    "Ball 2": { "fill": "rgba(156,0,255,0.05)", "stroke": "rgba(156,0,255,1)" },
    "Ball 3": { "fill": "rgba(54,0,255,0.05)", "stroke": "rgba(54,0,255,1)" },
    "Ball 4": { "fill": "rgba(0,84,255,0.05)", "stroke": "rgba(0,84,255,1)" },
    "Ball 5": { "fill": "rgba(0,186,255,0.05)", "stroke": "rgba(0,186,255,1)" },
    "Ball 6": { "fill": "rgba(0,255,30,0.05)", "stroke": "rgba(0,255,30,1)" },
    "Bonus Ball": { "fill": "rgba(255,0,0,0.05)", "stroke": "rgba(255,0,0,1)" },
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
    var p = {
      machine: $(sel.machineSelect).val()
    }
    $.getJSON("/api/sets", p, function (d) {
      var el = $(sel.setSelect)
      el.empty().append('<option value="0">All</option>')
      $.each(d, function (i, n) {
        el.append('<option value="' + n + '">' + n + '</option>')
      })
    })
  })

  /// Query Exec
  $(sel.querySubmit).click(function (e) {
    switch ($(sel.queryType).val()) {
      case "results/average":
        drawResultsAverage()
        break
      case "results/graph":
        drawResultsGraph()
        break
      case "results/plotly":
        drawResultsGraphPlotly()
        break
      default:
        break
    }
  })

  //******************************************************* FUNCTIONS */
  // Return a unix timestamp for n number of days ago
  function nDaysAgo(n) {
    var d = new Date()
    d.setDate(d.getDate() - n)
    return Math.floor(d.getTime() / 1000)
  }

  function params() {
    return {
      start: $(sel.dateStart).val(),
      end: $(sel.dateEnd).val(),
      set: $(sel.setSelect).val(),
      machine: $(sel.machineSelect).val()
    }
  }

  function drawResultsAverage() {
    $.getJSON("/api/results/average", params(), function (d) {
      // Empty results dive and append container for data
      $(sel.results).empty().append('<h1 id="' + sel.resultsAvg.replace('#', '') + '" class="centered"></h1>')
      for (i = 0; i < d.length; i++) {
        $(sel.resultsAvg).append("<span class='num'>" + d[i] + "</span>")
      }
    })
  }

  function drawResultsGraph() {
    $.getJSON("/api/results/graph", params(), function (d) {
      // Set graph colours 
      $.each(d.datasets, function (i, o) {
        var c = chartColours[o.label]

        o.backgroundColor = c.fill
        o.borderColor = c.stroke
        o.borderWidth = 1.5
        o.pointBorderColor = "rgba(0,0,0,0)"
        //o.lineTension = 0.3
      })

      // Empty results div and append container for data
      $(sel.results).empty().append('<canvas id="' + sel.resultsGraph.replace("#", "") + '"></canvas>')
      var resultsChart = new Chart($(sel.resultsGraph), {
        type: "line",
        data: d,
        options: chartOptions
      })
    })
  }

  function drawResultsGraphPlotly() {
    $.getJSON("/api/results/plotly", params(), function (d) {
      console.log(d)
      $(sel.results).empty().append('<div id="' + sel.resultsGraph.replace("#", "") + '"></div>')
      Plotly.newPlot(sel.resultsGraph.replace("#", ""), d)
    })
  }
})