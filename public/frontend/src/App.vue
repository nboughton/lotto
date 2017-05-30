<template>
  <div id="app" class="section">
    <div class="columns">
      <div class="column is-one-quarter">
        <div class="box">
          <b-field label="Query">
            <b-select v-model="params.query" expanded>
              <option v-for="q in queries" :value="q.id">{{ q.name }}</option>
            </b-select>
          </b-field>
          <b-field label="From">
            <Datepicker v-model="params.start" :disabled="flags.dates" input-class="input"></Datepicker>
          </b-field>
          <b-field label="To">
            <Datepicker v-model="params.end" :disabled="flags.dates" input-class="input"></Datepicker>
          </b-field>
          <b-field label="Set">
            <b-select v-model="params.set" expanded>
              <option value="0" selected>All</option>
              <option v-for="s in sets" :value="s">{{ s }}</option>
            </b-select>
          </b-field>
          <b-field label="Machine">
            <b-select v-model="params.machine" expanded>
              <option value="all" selected>All</option>
              <option v-for="m in machines" :value="m">{{ m }}</option>
            </b-select>
          </b-field>
          <button class="button" @click="runQuery">Submit</button>
        </div>
      </div>
      <div class="column is-three-quarters">
        <div class="box">
          <LineChart :chart-data="lineChart.data" :options="lineChart.options"></LineChart>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Vue from 'vue'
import Buefy from 'buefy'
import 'buefy/lib/buefy.css'
import VueResource from 'vue-resource'
import Datepicker from 'vuejs-datepicker'
import LineChart from '@/components/LineChart'

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
        query: 1,
        start: new Date(2015, 9, 10),
        end: new Date(),
        machine: "all",
        set: 0
      },
      flags: {
        dates: {
          to: new Date(1994, 10, 19),
          from: new Date()
        }
      },
      queries: [
        { id: 1, name: "Most Recent Results" }
      ],
      sets: [],
      machines: [],
      lineChart: {
        data: {},
        options: {}
      }
    }
  },

  computed: {
    qParams: function () {
      var start = this.params.start.toJSON()
      var end = this.params.end.toJSON()

      return {
        query: this.params.query,
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

        this.lineChart.data = d.data
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
    }
  },

  components: {
    Datepicker,
    LineChart
  }
}
</script>

<style>

</style>
