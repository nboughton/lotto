{{ template "includes/header" }}
  <div id="container">
    <section id="query">
      <span class="form-element">
        <label for="query-type">Query Type: </label>
        <select id="query-type">
          <option value="results/average">Average Results</option>
          <option value="results/graph">Graph Results</option>
          <option value="results/plotly">Plotly Test</option>
        </select>
      </span>
      <span class="form-element">
        <label for="date-start">Start: </label>
        <input type="date" class="form-date" id="date-start" value="{{ .Start }}" />
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

