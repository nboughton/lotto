{{ template "includes/header" }}
<div id="container">
  <section id="query">
    <span class="form-element">
        <label for="query-type">Query Type: </label>
        <select id="query-type">
          <option value="graph-results-freqdist-bar">Frequency Dist (results, bar)</option>
          <option value="graph-results-freqdist-scatter">Frequency Dist (results, scatter)</option>
          <option value="graph-results-timeseries-scatter">Time Series (results, scatter/trend)</option>
          <option value="graph-results-timeseries-line">Time Series (results, line)</option>
          <option value="graph-results-raw-scatter3d">3D Scatter Plot (results)</option>
          <option value="graph-ms-freqdist-bubble">Frequency Dist (machine/set, bubble)</option>
          <option value="graph-ms-freqdist-scatter3d">Frequency Dist (machine/set, scatter3d)</option>
          <option value="num-results">Averages, Ranges and Most Frequent</option>
        </select>
      </span>
    <span class="form-element">
        <label for="date-start">Start: </label>
        <input type="date" class="form-date" id="date-start" value="2015-10-10" />
      </span>
    <span class="form-element">
        <label for="date-end">End: </label>
        <input type="date" class="form-date" id="date-end" value="{{ .End }}" />
      </span>
    <span class="form-element">
        <label for="filter-machine">Machine: </label>
        <select id="filter-machine">
          <option value="all">All</option>
          {{ range .Machines }}
          <option value="{{ . }}">{{ . }}</option>
          {{- end }}
        </select>
      </span>
    <span class="form-element">
        <label for="filter-set">Set: </label>
        <select id="filter-set">
          <option value="0">All</option>
          {{ range .Sets }}
          <option value="{{ . }}">{{ . }}</option>
          {{- end }}
        </select>
      </span>
    <span class="form-element">
        <label for="query-submit"></label>
        <input type="button" id="query-submit" value="Show Data" />
      </span>
  </section>

  <section id="results"></section>
</div>
{{ template "includes/footer" }}