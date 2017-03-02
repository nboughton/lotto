{{ template "includes/header" }}
  <div id="app-container">
    <div id="graph-query gq">
      <label for="gq-start-date">Start Date: <input type="date" class="gq-date" id="gq-start-date" /></label>
      <label for="gq-end-date">End Date: <input type="date" class="gq-date" id="gq-end-date" /></label>
      <label for="gq-machine-filter">Filter by Machine: 
        <select id="gq-machine-filter">
          <option value="All">All</option>
          {{ range .Machines }}
          <option value="{{ . }}">{{ . }}</option>
          {{ end }}
        </select>
      </label>
      <label for="gq-set-filter">Filter by Set:
        <select id="gq-set-filter">
          <option value="All">All</option>
          {{ range .Sets }}
          <option value="{{ . }}">{{ . }}</option>
          {{ end }}
        </select>
      </label>
      <label for="gq-submit"><input type="button" id="gq-submit" value="Show Data" /></label>
    </div>
  </div>
{{ template "includes/footer" }}

