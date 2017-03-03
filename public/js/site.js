$(function () {
  // Define a map for key ids/classes in case I decide to change them later
  var defs = {
    machineSelect: "#gq-machine-filter",
    setSelect: "#gq-set-filter",
    dateClass: ".gq-date",
    dateStart: "#gq-start-date",
    dateEnd: "#gq-end-date",
    results: "#gq-query-results",
    resultsAvg: "#gq-average-num-res"
  }

  /// Utility Functions
  // Return a unix timestamp for n number of days ago
  function nDaysAgo(n) {
    var d = new Date()
    d.setDate(d.getDate() - n)
    return Math.floor(d.getTime() / 1000)
  }

  /// Pickadate.js
  $.getJSON("/api/range", function (data) {
    var pickerOpts = {
      "format": "yyyy-mm-dd",
      "selectYears": true,
      "selectMonths": true,
      "min": new Date(data.first * 1000),
      "max": new Date(data.last * 1000)
    }

    $(defs.dateClass).pickadate(pickerOpts)
  })

  /// Redraw options list for Sets depending on Machine selection
  $(defs.machineSelect).change(function (e) {
    var params = {
      machine: $(defs.machineSelect).val()
    }
    $.getJSON("/api/sets", params, function (data) {
      var el = $(defs.setSelect)
      el.empty()
      el.append('<option value="0">All</option>')
      $.each(data, function (i, n) {
        el.append('<option value="' + n + '">' + n + '</option>')
      })
    })
  })

  /// Query Exec
  $("#gq-submit").click(function (e) {
    var params = {
      start: $(defs.dateStart).val(),
      end: $(defs.dateEnd).val(),
      set: $(defs.setSelect).val(),
      machine: $(defs.machineSelect).val()
    }
    switch ($("#gq-query-type").val()) {
      case "results/average":
        $.getJSON("/api/results/average", params, function (data) {
          $(defs.results).empty()
          $(defs.results).append('<h1 id="gq-average-num-res" class="centered"></h1>')
          for (i = 0; i < data.length; i++) {
            $(defs.resultsAvg).append("<span class='num'>" + data[i] + "</span>")
          }
        })
    }
  })

  /// Chart.JS
})