$(function () {
  /// Utility Functions
  // Return a unix timestamp for n number of days ago
  function nDaysAgo(n) {
    var d = new Date()
    d.setDate(d.getDate() - n)
    return Math.floor(d.getTime() / 1000)
  }

  /// Pickadate.js
  var pickerOpts = {
    "format": "yyyy-mm-dd",
    "selectYears": true,
    "selectMonths": true
  }

  $(".gq-date").pickadate(pickerOpts)

  /// Chart.JS

})