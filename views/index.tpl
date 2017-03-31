{{ template "includes/header" }}
<div id="container">
  <section id="query">
    <span class="form-element">
        <label for="query-type">Query Type: </label>
        <select id="query-type">
          <option value="graph-freqdist-bar">Frequency Dist (bar)</option>
          <option value="graph-freqdist-scatter">Frequency Dist (scatter)</option>
          <option value="graph-timeseries-scatter">Time Series (scatter/line)</option>
          <option value="graph-3d-scatter">3D Scatter Plot (pretty!)</option>
          <option value="num-average-mean">Mean Averages</option>
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