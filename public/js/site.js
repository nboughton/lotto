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
    $.getJSON("/api/sets", params(), function (d) {
      var el = $(sel.setSelect)
      el.empty().append('<option value="0">All</option>')
      $.each(d, function (i, n) {
        el.append('<option value="' + n + '">' + n + '</option>')
      })
    })
  })
/*
  $(sel.dateClass).change(function (e) {
    var that = this
    $.getJSON("/api/machines", params(), function (d) {
      that.empty().append('<option value="all">All</option>')
      $.each(d, function (i, m) {
        that.append('<option value="' + m + '">' + m + '</option>')
      })
    })
  })
  */

  /// Query Exec
  $(sel.querySubmit).click(function (e) {
    switch ($(sel.queryType).val()) {
      case "num-average":
        drawResultsAverage()
        break
      case "graph-line":
        drawResultsGraph("line")
        break
      case "graph-3d-scatter":
        drawResultsGraph("3d/scatter")
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

  function drawResultsGraph(type) {
    $.getJSON("/api/results/graph/" + type, params(), function (d) {
      $(sel.results).empty().append('<div id="' + sel.resultsGraph.replace("#", "") + '"></div>')
      Plotly.newPlot(sel.resultsGraph.replace("#", ""), d)
    })
  }

})