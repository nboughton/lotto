$(function () {
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
      "min": new Date(data.first*1000),
      "max": new Date(data.last*1000)
    }

    $(".gq-date").pickadate(pickerOpts)
  })


  /// Chart.JS

})