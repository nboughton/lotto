<template>
  <div id="app" class="section">
    <div class="heading">
      <p class="title">
        UK Lottery Data
      </p>
    </div>
    <div class="box">
      <div class="heading">
        <p class="subtitle">
          Parameters
        </p>
      </div>
      <b-field>
        <p class="control">
          <button class="button is-static">FROM</button>
        </p>
        <Datepicker v-model="params.start" v-on:input="adjustFieldData" :disabled="flags.dates" input-class="input"></Datepicker>
        <p class="control">
          <button class="button is-static">TO</button>
        </p>
        <Datepicker v-model="params.end" v-on:input="adjustFieldData" :disabled="flags.dates" input-class="input"></Datepicker>
        <b-select v-model="params.set" expanded>
          <option disabled value="">Please select one</option>
          <option value="0">All Sets</option>
          <option v-for="s in sets" :value="s">Set: {{ s }}</option>
        </b-select>
        <b-select v-model="params.machine" expanded>
          <option disabled value="">Please select one</option>
          <option value="all">All Machines</option>
          <option v-for="m in machines" :value="m">Machine: {{ m }}</option>
        </b-select>
        <p class="control">
          <button class="button" @click="runQuery">Submit</button>
        </p>
      </b-field>
    </div>
    <div class="box">
      <div class="heading">
        <p class="subtitle">
          Data
        </p>
      </div>
      <b-table :data="tables.main.data" :default-sort="[ 'num', 'desc' ]" :bordered="tables.main.cfg.isBordered" :striped="tables.main.cfg.isStriped" :narrowed="tables.main.cfg.isNarrowed" :checkable="tables.main.cfg.isCheckable" :paginated="tables.main.cfg.isPaginated">
        <template scope="props">
          <b-table-column field="label" label="" sortable>
            {{ props.row.label}}
          </b-table-column>
  
          <b-table-column field="num[0]" label="Ball 1" sortable numeric>
            {{ props.row.num[0] }}
          </b-table-column>
  
          <b-table-column field="num[1]" label="Ball 2" sortable numeric>
            {{ props.row.num[1] }}
          </b-table-column>
  
          <b-table-column field="num[2]" label="Ball 3" sortable numeric>
            {{ props.row.num[2] }}
          </b-table-column>
  
          <b-table-column field="num[3]" label="Ball 4" sortable numeric>
            {{ props.row.num[3] }}
          </b-table-column>
  
          <b-table-column field="num[4]" label="Ball 5" sortable numeric>
            {{ props.row.num[4] }}
          </b-table-column>
  
          <b-table-column field="num[5]" label="Ball 6" sortable numeric>
            {{ props.row.num[5] }}
          </b-table-column>
  
          <b-table-column field="num[6]" label="Bonus" sortable numeric>
            {{ props.row.num[6] }}
          </b-table-column>
        </template>
      </b-table>
    </div>
    <div class="box">
      <div class="heading">
        <p class="subtitle">
          Frequency Distribution
        </p>
      </div>
      <BarChart :chart-data="charts.freqDist.data" :options="charts.freqDist.options" :height="200"></BarChart>
    </div>
    <div class="box">
      <div class="heading">
        <p class="subtitle">
          Results Over Time
        </p>
      </div>
      <LineChart :chart-data="charts.timeSeries.data" :options="charts.timeSeries.options" :height="200"></LineChart>
    </div>
  </div>
</template>

<script>
import Vue from 'vue'
import Buefy from 'buefy'
import VueResource from 'vue-resource'
import Datepicker from 'vuejs-datepicker'
import LineChart from '@/components/LineChart'
import BarChart from '@/components/BarChart'

Vue.use(Buefy)
Vue.use(VueResource)

var colors = [
  "rgba(31,119,180,1)",
  "rgba(255,127,14,1)",
  "rgba(44,160,44,1)",
  "rgba(214,39,40,1)",
  "rgba(148,103,189,1)",
  "rgba(140,86,75,1)",
  "rgba(227,119,194,1)",
]

export default {
  name: 'app',

  data: function () {
    return {
      params: {
        start: new Date(2015, 9, 10),
        end: new Date(),
        machine: "all",
        set: "0"
      },
      flags: {
        dates: {
          to: new Date(1994, 10, 19),
          from: new Date()
        }
      },
      sets: [],
      machines: [],
      tables: {
        main: {
          data: [],
          cfg: {
            isBordered: false,
            isStriped: true,
            isNarrowed: true,
            isCheckable: false,
            isPaginated: false,
          }
        }
      },
      charts: {
        timeSeries: {
          data: {},
          options: {
            tooltips: {
              mode: "index"
            }
          }
        },
        freqDist: {
          data: {},
          options: {
            tooltips: {
              mode: "index"
            },
            scales: {
              xAxes: [{
                stacked: true,
                barPercentage: 1,
              }],
              yAxes: [{
                stacked: true
              }]
            }
          }
        }
      }
    }
  },

  computed: {
    qParams: function () {
      var start = this.params.start.toJSON()
      var end = this.params.end.toJSON()

      return {
        start: start,
        end: end,
        machine: this.params.machine,
        set: this.params.set
      }
    }
  },

  created: function () {
    this.getSets()
    this.getMachines()
    this.runQuery()
  },

  methods: {
    runQuery: function () {
      this.$http.get("/api/query", { params: this.qParams }).then(resp => { return resp.json() }).then(d => {
        for (var ball = 0; ball < d.data.timeSeries.datasets.length; ball++) {
          // Line chart
          d.data["timeSeries"].datasets[ball].backgroundColor = colors[ball]
          d.data["timeSeries"].datasets[ball].borderColor = colors[ball]
          d.data["timeSeries"].datasets[ball].borderWidth = 0.3
          d.data["timeSeries"].datasets[ball].pointRadius = 3
          d.data["timeSeries"].datasets[ball].fill = false
          d.data["timeSeries"].datasets[ball].showLines = false

          // Bar chart
          d.data["freqDist"].datasets[ball].backgroundColor = colors[ball]
          d.data["freqDist"].datasets[ball].borderColor = colors[ball]
          d.data["freqDist"].datasets[ball].borderWidth = 0.3
        }

        this.tables.main.data = d.data.mainTable
        this.charts.timeSeries.data = d.data.timeSeries
        this.charts.freqDist.data = d.data.freqDist
      })
    },
    getSets: function () {
      this.$http.get("/api/sets", { params: this.qParams }).then(resp => { return resp.json() }).then(sets => {
        this.sets = sets.data
      })
    },
    getMachines: function () {
      this.$http.get("/api/machines", { params: this.qParams }).then(resp => { return resp.json() }).then(machines => {
        this.machines = machines.data
      })
    },
    adjustFieldData: function () {
      this.getSets()
      this.getMachines()
    },
  },

  components: {
    Datepicker,
    LineChart,
    BarChart
  }
}

</script>

<style lang="scss">
@import "~bulmaswatch/yeti/_variables.scss";
@import "~bulma";
@import "~buefy/src/scss/buefy";
@import "~bulmaswatch/yeti/_overrides.scss";
</style>