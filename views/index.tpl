{{ template "includes/header" }}
  <div id="app-container">
    <div id="graph-query" class="gq">
      <span class="gq-input"><label for="gq-query-type">Query Type: </label>
        <select id="gq-query-type">
          <option value="average">Average Results</option>
        </select>
      </span>
      <span class="gq-input">
        <label for="gq-start-date">Start Date: </label>
        <input type="date" class="gq-date" id="gq-start-date" value="{{ .Start }}" />
      </span>
      <span class="gq-input">
        <label for="gq-end-date">End Date: </label>
        <input type="date" class="gq-date" id="gq-end-date" value="{{ .End }}" />
      </span>
      <span class="gq-input">
        <label for="gq-machine-filter">Filter by Machine: </label>
        <select id="gq-machine-filter">
          <option value="all">All</option>
          {{ range .Machines }}
          <option value="{{ . }}">{{ . }}</option>
          {{ end }}
        </select>
      </span>
      <span class="gq-input">
        <label for="gq-set-filter">Filter by Set: </label>
        <select id="gq-set-filter">
          <option value="0">All</option>
          {{ range .Sets }}
          <option value="{{ . }}">{{ . }}</option>
          {{ end }}
        </select>
      </span>
      <span class="gq-input"><label for="gq-submit"></label><input type="button" id="gq-submit" value="Show Data" /></span>
    </div>
    <div id="gq-query-results"></div>
  </div>
{{ template "includes/footer" }}

