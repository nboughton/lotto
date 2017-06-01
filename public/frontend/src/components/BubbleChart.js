import { Bubble, mixins } from 'vue-chartjs'

export default Bubble.extend({
  mixins: [mixins.reactiveProp],
  props: ["chartData", "options"],
  mounted() {
    this.renderChart(this.chartData, this.options)
  }
})